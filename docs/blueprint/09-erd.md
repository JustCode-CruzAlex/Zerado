---
title: Zerado — the data model
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-09
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: erd
ticket: "#2"
---

# The data model

> **This is the CONCEPTUAL model** — entities, what each is for, relationships and cardinality, and
> what is deliberately *not* in the file. **`fft-database` owns the physical model**: exact types,
> constraints, indexes and the migration sequence. Column lists below are shown at the fidelity
> needed to make the architecture argument, and are illustrative of intent rather than a schema to
> run. [`13-handoffs.md`](./13-handoffs.md) states which claims are load-bearing.

> **Zerado is a games model.** The core entity is `item` — not `games`, and carrying an `item_type`
> constrained to `'game'` — which is the whole of the "door stays open" affordance in
> [`11-media-model.md`](./11-media-model.md). There are no typed extensions, no generic progress and
> no second type. Playtime is playtime.

The schema, drawn. Two sheets, each rendered in both the brand-black and the cyanotype theme.

| Sheet | Drawing | Spec | Rendered |
|---|---|---|---|
| **01 of 10** | Phase 1 core — one SQLite file, and the two things deliberately outside it | [`ZRD-ERD-01.chart.toml`](../adr/charts/ZRD-ERD-01.chart.toml) | [`svg/ZRD-ERD-01.svg`](../adr/charts/svg/ZRD-ERD-01.svg) · [cyanotype](../adr/charts/svg/ZRD-ERD-01.cyanotype.svg) |
| **02 of 10** | Phase 2 enrichment and the Phase 4 sync boundary | [`ZRD-ERD-02.chart.toml`](../adr/charts/ZRD-ERD-02.chart.toml) | [`svg/ZRD-ERD-02.svg`](../adr/charts/svg/ZRD-ERD-02.svg) · [cyanotype](../adr/charts/svg/ZRD-ERD-02.cyanotype.svg) |

```bash
flowforge chart render docs/adr/charts/ZRD-ERD-01.chart.toml   # both themes
flowforge chart render docs/adr/charts/ZRD-ERD-02.chart.toml
```

**Colour key** — the chart renderer draws legend swatches but does not draw their captions, so the
key is here instead: **violet** = the `game` row and its sync bookkeeping · **gold** = Phase 2
enrichment, every row stamped with its age · **blue** = Phase 4, server side · **red** = either
deliberately outside the library file, or the one payload that crosses the sync boundary ·
**grey** = local-only, never syncs.

---

## 1 · Phase 1 — the tables

### `item` — the row everything else hangs off

| Column | Type | Notes |
|---|---|---|
| `id` | `INTEGER PK` | local surrogate |
| `item_uid` | `TEXT NOT NULL` | **indexed, not unique.** A Phase 4 merge *hint* — see §4 |
| **`item_type`** | `TEXT NOT NULL` | `'game'`, **constrained to it**. The door-open affordance, and nothing more |
| `provider_id` | `TEXT NOT NULL` | `steam` · `physical` · … |
| `provider_ref` | `TEXT NOT NULL` | the provider's own id. Steam appid; a UUID for physical |
| `acquisition` | `TEXT NOT NULL` | `digital` · `physical` |
| `title` · `sort_title` | `TEXT NOT NULL` | articles stripped, diacritics folded |
| `platform` | `TEXT NOT NULL` | |
| **`playtime_minutes`** | `INTEGER NOT NULL DEFAULT 0` | provider-reported. **Zero is a real value** and is distinct from *this provider does not report playtime* — see `Capabilities.Playtime` |
| `last_played_at` | `TEXT NULL` | **`NULL` = not reported, not never played** |
| `owned_since` | `TEXT NULL` | |
| `steam_appid` | `TEXT NULL` | |
| `achievements_total` · `achievements_unlocked` | `INTEGER NULL` | may *suggest* `ZERADO` in Phase 2; never sets it |
| `status_manual` | `TEXT NULL` | `not_started` · `in_progress` · `zerado` · `abandoned` · `NULL` |
| `status_changed_at` | `TEXT NULL` | decides the Phase 4 conflict |
| `rating` · `notes` | | the player's own |
| `absent_since` | `TEXT NULL` | Set when a **complete** sync stops returning it; cleared when it returns. **Never a reason to delete** — [`06-data-seams.md`](./06-data-seams.md) §2.4 |
| `merged_into` | `INTEGER NULL` | FK → `item.id`, so a Phase 4 merge never rewrites keys |
| `created_at` `updated_at` | `TEXT NOT NULL` | |

```sql
UNIQUE (provider_id, provider_ref)
CHECK  (item_type = 'game')
INDEX  (item_uid)
INDEX  (sort_title)
```

**Game facts live on the item, not in a typed extension.** Revision A split them into a
`media_game` table to serve a polymorphic core; with books pruned there is no second type to
separate them from, and a join that serves one type is machinery without a purpose. When a second
type is genuinely on the table, that is the moment to decide between an extension table and nullable
columns — with a real second type in hand to decide against.

**`playtime_minutes` is a plain column, not generic progress.** Revision A modelled it as
`progress_value` + `_unit` + `_source`, on the argument that it kept the four-state derivation a
single function across media types. **That argument died with the media types**, so the divergence
is withdrawn and playtime is playtime — which is what the founder's direction said in the first
place.

**`effective_status` is not a column.** It is derived on read:
`status_manual ?? derive(playtime_minutes, capabilities)`. A stored derived value has two ways to be
wrong — stale, or written by the wrong path. A derived one has none.

### `provider_connection`

| Column | Notes |
|---|---|
| `provider_id` `PK` | |
| `account_ref` | the Steam ID, the GOG username — an identifier, never a secret |
| `connected_at` | |
| `last_sync_at` · `last_sync_status` | what the degrade banner reads to say "3 days ago" |

**No credentials.** They are in the Vault (§3).

### `sync_run` — the honest history

| Column | Notes |
|---|---|
| `id` `PK` · `provider_id` `FK` | |
| `started_at` · `finished_at` | `finished_at NULL` = it was killed, not cancelled cleanly |
| `status` | `ok` · `partial` · `failed` · `cancelled` |
| `items_seen` · `items_new` · `items_changed` | what `Z-03 DONE` reports |
| `error` | the classified failure, not a stack trace |

`partial` exists because a cancelled sync leaves a **valid** library. Recording it is what lets
`Z-04` say *"synced 3 days ago, partially"* instead of implying it is complete.

### `setting` · `schema_migration`

`setting(key PK, value)` — everything `Z-09` writes.

`schema_migration(version PK, applied_at, **written_by**)` — forward-only, and a database whose max
version is **higher** than the binary knows about is a `Z-11 Fatal error`, not a silent downgrade.

**`written_by` is the reason this is not just `(version, applied_at)`.** `Z-11` promises to name
*both* versions — the one that wrote the file and the one that is running — and a schema number
alone cannot do that. A number tells the player their file is from the future; the binary version
tells them *which release to go and fetch*, which is the only actionable half. Adding the column
costs nothing now and cannot be backfilled later, because the binary that wrote a row is exactly
what is no longer around to ask.

*(Found by `fft-tui-designer`: `Z-11`'s copy required a fact the schema could not supply.)*

---

## 2 · Phase 2 — the enrichment tables

| Table | Key columns | The rule it carries |
|---|---|---|
| `metadata` | `item_id PK/FK` · `sinopse` · `cover_ref` · `released_at` · `genres` · `source_provider` · `attribution` · **`fetched_at NOT NULL`** | `cover_ref` is a **local cache path, never a remote URL** — nothing renders from the network. `attribution` comes from the provider, so swapping the source swaps the credit |
| `price_quote` | `item_id FK` · `shop` · `currency` · `current_cents` · `low_cents` · `low_at` · `url` · `affiliate_url` · **`observed_at NOT NULL`** | `affiliate_url` and the disclosure obligation live in the same row, so a refactor cannot separate them |
| `mood_tag` | `id PK` · **`key`** · **`applies_to[]`** · `label` | The engine reasons over the type-neutral `key`; the interface shows the per-type `label`. One engine, per-type vocabulary — [`11-media-model.md`](./11-media-model.md) §6 |
| `item_mood` | `item_id FK` · `mood_id FK` · `source` (`user` \| `inferred`) · `confidence` | `source` is what decides whether a tag crosses the Phase 4 boundary. **Only `user` does** |

**Every enrichment row carries its age, and the column is `NOT NULL`.** That is the schema
enforcing [`07-offline-contract.md`](./07-offline-contract.md) §4 rather than the renderer
remembering to.

---

## 3 · What is deliberately NOT in the file

| Not in `library.db` | Where it lives | Why |
|---|---|---|
| **Credentials** | The Vault — OS keychain, or `credentials.json` mode `0600` | The library file is a thing the player is invited to back up, move and delete. A key inside it would be a key in every backup, every copy and every support-ticket attachment |
| **Cover-art blobs** | The XDG **cache** directory | Disposable and refetchable, so the OS is allowed to delete it — the correct semantics — and it keeps the backed-up file small enough that people actually back it up |

Both are drawn on sheet 01 as boxes **outside** the `library.db` zone, because a diagram that
only shows what is inside teaches the wrong lesson.

---

## 4 · `item_uid` — the one Phase 1 column that exists for Phase 4

```
item_uid = uuidv5( namespace_zerado, item_type + "|" + normalise(title) + "|" + normalise(platform) )
```

`normalise` lowercases, strips punctuation and leading articles, and folds diacritics.

It is **stable, not authoritative.** Two editions of one game may collide; the same game may fail
to match across platforms. So three rules travel with it:

1. `item_uid` is **indexed, not unique**. `(provider_id, provider_ref)` remains the uniqueness
   constraint. It is **type-scoped**: the same title as a game and as a film are two items, and
   should be.
2. Phase 4's merge treats it as a **hint** and shows ambiguous matches to the player rather than
   guessing.
3. `merged_into` exists **from Phase 1**, so two rows can be joined later without a migration that
   rewrites primary keys.

Adding this now costs one column and one index. Adding it in Phase 4 costs a migration that has to
invent stable identities for rows whose titles the player has since edited — which is the migration
nobody wants to write. That is why it is in ADR-0001 rather than in a Phase 4 ticket.

---

## 5 · Phase 4 — the server side, and the boundary

Sheet 02's right-hand zone. The single red box is the point of the drawing:

| Crosses the boundary | Never crosses |
|---|---|
| `state_change` — `item_uid`, `status_manual`, `status_changed_at` | the library itself |
| user-assigned mood tags (`game_mood.source = 'user'`) | `playtime_minutes` · `last_played_at` |
| manually-entered games — **the only rows not re-derivable from a store** | cover art · *sinopse* · prices |
| | **credentials. Ever.** |

The rule underneath: **only what the player typed crosses.** Everything a machine can recompute,
each device recomputes from the player's own sources with the player's own keys.

That is a privacy posture and also an economics one: it is what keeps the Phase 4 server small
enough to be supportable by **donations alone**, which is itself a published
statement.

Conflict resolution is **last-write-wins on `status_changed_at`**, per game. The limit is stated
rather than hidden: two devices that both change a status while offline will lose the earlier
change silently. That is acceptable for one person on two devices and unacceptable to leave
undocumented, so `Z-22 Devices and sync` shows the last merge and what it resolved.

---

## 6 · Migration policy

- **Forward-only.** No down migrations. A mistake is corrected by a new migration, not by
  reversing one.
- **Every migration is idempotent** and runs inside a transaction.
- **A database from a newer binary is a fatal error**, named as such — `Z-11` says which version
  wrote it and which version is running. Silently proceeding against a schema you do not
  understand is how one file becomes two incompatible files.
- **The file is backed up before any migration that drops or rewrites a column**, to
  `library.db.pre-<version>` beside it, and `Z-11` names that path if the migration fails.
