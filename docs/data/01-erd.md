---
title: Zerado — the physical data model
discipline: DATA
doc-no: ZRD-DATA-01
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: erd
ticket: "#5"
---

# The physical data model

> **This is the PHYSICAL model.** [`../blueprint/09-erd.md`](../blueprint/09-erd.md) is the
> conceptual one and remains correct as such; this document makes it runnable. Where the two
> differ, §6 lists every difference with the reason and the citation that authorises it
> ([`../blueprint/13-handoffs.md`](../blueprint/13-handoffs.md) §3: *"`fft-database` owns the
> physical model and should expect to change types, add constraints, split or merge tables, and
> design indexes against real query shapes"*).

> **Zerado is a games model.** The core entity is `item` — not `games` — carrying an `item_type`
> `CHECK`ed to `'game'`. That is the whole of the door-open affordance
> ([`../blueprint/11-media-model.md`](../blueprint/11-media-model.md), ADR-0001 D5). There are no
> typed extensions, no generic progress and no second type. Playtime is playtime.
> [`07-the-door-three-times.md`](./07-the-door-three-times.md) answers the gate's question —
> *does the model hold for books without a rewrite?* — by writing the migration that opens it.

The DDL is [`schema/migrations/0001_phase1_core.sql`](./schema/migrations/0001_phase1_core.sql)
and it is quoted in full, verbatim, in
[`02-physical-schema.md`](./02-physical-schema.md) §9.

---

## 1 · The drawings

Two new sheets, each rendered in both the brand-black and the cyanotype theme. They sit alongside
sheets 01–02 from the spine, which draw the **conceptual** model and are not superseded.

| Sheet | Drawing | Spec | Rendered |
|---|---|---|---|
| **03 of 12** | The Phase 1 physical model — truth, replica and cache, and the line between them | [`ZRD-ERD-03.chart.toml`](../adr/charts/ZRD-ERD-03.chart.toml) | [brand-black](../adr/charts/svg/ZRD-ERD-03.svg) · [cyanotype](../adr/charts/svg/ZRD-ERD-03.cyanotype.svg) |
| **04 of 12** | The migration decision at start-up — five outcomes, every one of them named | [`ZRD-ERD-04.chart.toml`](../adr/charts/ZRD-ERD-04.chart.toml) | [brand-black](../adr/charts/svg/ZRD-ERD-04.svg) · [cyanotype](../adr/charts/svg/ZRD-ERD-04.cyanotype.svg) |

```bash
flowforge chart render docs/adr/charts/ZRD-ERD-03.chart.toml   # both themes
flowforge chart render docs/adr/charts/ZRD-ERD-04.chart.toml
```

**Colour key for sheet 03** — the renderer draws legend swatches but not their captions:
**violet** = identity and provider replica · **green** = **truth**, only the player writes it ·
**gold** = **cache**, safe to discard, every row stamped with its age · **blue** = sync
bookkeeping · **grey** = local-only, never syncs · **red** = deliberately outside the library file.

---

## 2 · Entities, and what each one is for

Eleven objects across three migrations. Phase 1 creates six tables, one view and five indexes;
Phase 2 and Phase 3 add their own when their binaries ship.

| Entity | Phase | Class | What it is for |
|---|---|---|---|
| **`item`** | 1 | identity + **replica** | One owned title. Everything a *provider* says about it |
| **`item_user`** | 1 | **truth** | Everything the *player* said about it. 1 : 1 with `item` |
| **`item_view`** | 1 | derived | The read surface. `effective_status` is computed here and stored nowhere |
| **`provider_connection`** | 1 | bookkeeping | A credential *relationship*, and **no credential** |
| **`sync_run`** | 1 | bookkeeping | The honest history: what a sync saw, and how it ended |
| **`setting`** | 1 | local truth | Everything `Z-09` writes. Never crosses to Phase 4 |
| **`schema_migration`** | 1 | bookkeeping | The ladder, and **who wrote each rung** |
| **`metadata`** | 2 | **cache** | *Sinopse*, cover path, genres — with its age, `NOT NULL` |
| **`mood_tag`** | 2 | vocabulary | `key` the recommender reasons over, `label` the player reads |
| **`item_mood`** | 2 | mixed | `source = 'user'` is truth; `source = 'inferred'` is cache |
| **`price_quote`** | 3 | **cache** | One row per (item, shop, currency). **There is no history table** |

---

## 3 · Relationships and cardinality

```
                       provider_connection
                        provider_id (PK)
                              │ 1
                              │
                              │ N          ON DELETE CASCADE
                          sync_run          (history follows the relationship;
                                             no item is ever touched)

     ┌──────── NO FOREIGN KEY ────────┐    item.provider_id names a provider;
     │      (§3.1 — deliberate)       │    it does not REFERENCE a connection
     ▼                                ▲
   item ──────────── 1 : 1 ──────────── item_user        CASCADE
    │  id (PK)                             item_id (PK, FK)
    │
    │ 1 : 0..1   ──────────────────────►  metadata        CASCADE   (Phase 2)
    │ 1 : N      ──────────────────────►  price_quote     CASCADE   (Phase 3)
    │ 1 : N      ──────────────────────►  item_mood       CASCADE   (Phase 2)
    │                                          │ N
    │                                          │ 1
    │                                      mood_tag                 (Phase 2)
    │
    └─ merged_into ─► item.id     0..1 : 1     ON DELETE SET NULL
       (self-reference: a Phase 4 merge joins two rows without
        rewriting a single primary key)

   setting            standalone, key/value
   schema_migration   standalone, the ladder
```

| Parent | Child | Cardinality | On delete | Why that action |
|---|---|---|---|---|
| `item` | `item_user` | **1 : 1**, materialised by trigger | `CASCADE` | The player's work about a row that no longer exists is not work |
| `item` | `metadata` | 1 : 0..1 | `CASCADE` | Cache about a gone row is garbage |
| `item` | `price_quote` | 1 : N *(shop × currency)* | `CASCADE` | as above |
| `item` | `item_mood` | 1 : N | `CASCADE` | as above |
| `mood_tag` | `item_mood` | 1 : N | `CASCADE` | Retiring a vocabulary word retires its assignments |
| `item` | `item` via `merged_into` | 0..1 : 1 | `SET NULL` | Deleting a merge *target* must not delete the row merged into it. It un-merges |
| `provider_connection` | `sync_run` | 1 : N | `CASCADE` | A reconnect must not inherit the previous relationship's *"last synced"* age |
| `provider_connection` | `item` | **none** | — | §3.1 |

### 3.1 · The foreign key that is deliberately absent

[`09-erd.md`](../blueprint/09-erd.md) sheet 01 draws `item → provider_connection` as
`provider_id  N : 1`. **The association is real; the referential constraint must not be.**

A connection is a relationship the player can **end**. The library outlives it — that is the whole
point of tombstoning rather than deleting ([`06-data-seams.md`](../blueprint/06-data-seams.md)
§2.4). An FK forces a choice between two wrong behaviours:

| If the FK existed | What disconnecting Steam would do |
|---|---|
| `ON DELETE CASCADE` | **Delete the library.** Every status, note and rating with it |
| `ON DELETE SET NULL` | Impossible — `provider_id` is `NOT NULL`, and a row with no source is not a row |
| `ON DELETE RESTRICT` / no action | **Forbid disconnecting** while any synced item exists, i.e. forever |

There is a fourth reason the FK is wrong even ignoring deletes: **`physical` has no connection
row.** It needs no credentials and never syncs, so a `provider_connection` row for it would be a
synthetic record invented to satisfy a constraint — and a synthetic row is a lie the schema tells
to keep itself consistent.

So `item.provider_id` is a **name**, resolved in code against the compiled-in provider registry.
ADR-0001 D1 already makes that registry the authority — *"everything downstream reads
`Capabilities`, never `ProviderID`"* — and a database copy of a compiled-in fact is a second
source of truth that can drift.

---

## 4 · `item` — attributes

Full types, constraints and defaults are in
[`02-physical-schema.md`](./02-physical-schema.md) §2. This is the shape.

| Column | Type | Null | Note |
|---|---|---|---|
| `id` | `INTEGER` | PK | rowid alias. Local, never crosses a device boundary |
| `item_uid` | `TEXT` | NN | uuidv5. **Indexed, not unique.** A merge *hint* — §5 |
| **`item_type`** | `TEXT` | NN | `CHECK (item_type = 'game')`. The door, and nothing more |
| `provider_id` · `provider_ref` | `TEXT` | NN | **`UNIQUE` together** — the sync's conflict target |
| `acquisition` | `TEXT` | NN | `digital` \| `physical` |
| `title` | `TEXT` | NN | Exactly as the source gave it. Never normalised in place |
| `sort_title` | `TEXT` | NN | Derived: folded, **articles stripped**. `ORDER BY` |
| `search_title` | `TEXT` | NN | Derived: folded, **articles kept**. `LIKE` |
| `platform` | `TEXT` | NN | |
| **`playtime_minutes`** | `INTEGER` | **NULL** | **`NULL` = this source cannot know. `0` = it knows, and it is zero** |
| `last_played_at` | `TEXT` | NULL | `NULL` = not reported, **not** never played |
| `owned_since` | `TEXT` | NULL | `YYYY` \| `YYYY-MM` \| `YYYY-MM-DD`. **Never a clock** |
| `steam_appid` | `TEXT` | NULL | |
| `achievements_total` · `_unlocked` | `INTEGER` | NULL | May *suggest* `ZERADO` in Phase 2. Never sets it |
| `absent_since` | `TEXT` | NULL | Tombstone. **Never a reason to delete** |
| `merged_into` | `INTEGER` | NULL | FK → `item(id)`, so a merge never rewrites keys |
| `last_synced_at` | `TEXT` | NULL | **Per row** — §4.2 |
| `created_at` · `updated_at` | `TEXT` | NN | RFC 3339, UTC, 20 characters exactly |

### 4.1 · `playtime_minutes` is nullable, and that is the change worth arguing for

[`09-erd.md`](../blueprint/09-erd.md) shows `INTEGER NOT NULL DEFAULT 0` with the note *"Zero is a
real value and is distinct from **this provider does not report playtime**."*

**The note is right and the type cannot express it.** With `NOT NULL DEFAULT 0`, both facts are
stored as `0`, and telling them apart requires a second lookup into compiled-in `Capabilities`
keyed by `provider_id` — which is exactly the `switch on ProviderID` that D1 exists to abolish.

Three things say the column should be nullable:

1. **The provider seam already carries the distinction.**
   [`06-data-seams.md`](../blueprint/06-data-seams.md) §2.1 declares `Item.Playtime *int` with
   *"nil = not reported"*. A `NOT NULL DEFAULT 0` column **flattens that pointer at the storage
   boundary** — the seam is careful and the schema throws the care away.
2. **A screen renders the difference from the row alone.** `Z-07`'s own mockup shows `0h` for
   Steam rows and `—` for the `PHY` row in the same list; `Z-05` D4 renders `not tracked`. Both
   read one row and must not consult a registry to draw a column.
3. **It makes the state derivation total.** `NULL > 0` is `NULL`, which falls to the `ELSE`
   branch — `not_started`. That is precisely what `derive()` returns when `!c.Playtime`
   ([`05-state-machine.md`](../blueprint/05-state-machine.md) §2), so `effective_status` becomes
   computable **in SQL, from the row, with nothing from code**. §5 of
   [`02-physical-schema.md`](./02-physical-schema.md) shows the view.

The invariant it creates, which a Phase 1 test must assert:
`playtime_minutes IS NULL` ⟺ the writing provider's `Capabilities.Playtime` is `false`.

**Cost:** `SUM(playtime_minutes)` needs `COALESCE`, and a Go scan needs `*int` or a
`sql.NullInt64`. Both are the honest shape of the fact.

### 4.2 · `last_synced_at` is per row, because `PARTIAL` is a real state

`Z-05` block 3 renders `LAST SYNCED  3 hours ago` for **this game**. A provider-level
`provider_connection.last_sync_at` cannot supply it honestly, because a `partial` run
([`06-data-seams.md`](../blueprint/06-data-seams.md) §2.5, `Z-03` `PARTIAL`) leaves some rows
freshly written and others untouched from days ago. One timestamp for both would be **true for the
rows the run reached and a lie for the rest** — and the rest are invisible, so nobody would find
out.

So: `item.last_synced_at` is stamped by the upsert that wrote the row. `Z-04`'s banner
(*"Last synced 3 days ago"*) keeps reading `provider_connection.last_sync_at`, because that
statement is about the **connection**, not about a row.

*(This column is not in the conceptual model. It was found by reading `Z-05` §3.2 against
`Z-03`'s `PARTIAL` state — the same way `written_by` was found by reading `Z-11`.)*

---

## 5 · `item_uid` — the one Phase 1 column that exists for Phase 4

```
item_uid = uuidv5( NAMESPACE_ZERADO,
                   item_type || '|' || uid_norm(title) || '|' || uid_norm(platform) )

NAMESPACE_ZERADO = uuidv5( NAMESPACE_DNS, "zerado.app" )
                 = 56640350-c577-5522-8dd0-30e65323adf8      ← computed, not invented
```

**The namespace is derived, not random**, so anyone can re-derive it from a sentence rather than
trusting a constant nobody can check.

`uid_norm` is specified with worked examples in
[`02-physical-schema.md`](./02-physical-schema.md) §4 — it was handed to this ticket by
[`13-handoffs.md`](../blueprint/13-handoffs.md) §4.

Three rules travel with it, unchanged from ADR-0001 D4:

1. **Indexed, not unique.** `(provider_id, provider_ref)` remains the uniqueness constraint. Two
   editions of one game may collide, and `seed-edge-cases.sql` contains exactly that pair.
2. **A hint, never an authority.** Phase 4 shows ambiguous matches to the player.
3. **`merged_into` exists from Phase 1**, so joining two rows never rewrites a primary key.

**It is type-scoped** — `item_type` is inside the hash — following
[`09-erd.md`](../blueprint/09-erd.md) §4. [`06-data-seams.md`](../blueprint/06-data-seams.md)
§6.2's earlier formula omits it; the difference is recorded in
[`08-findings.md`](./08-findings.md) F-4 rather than silently resolved.

---

## 6 · Every difference from the conceptual model

Nothing here is a decision reopened. Each row is a physical choice
[`13-handoffs.md`](../blueprint/13-handoffs.md) §1 and §3 assign to this ticket, and each names
what would have to change to reverse it.

| # | Conceptual model says | This model does | Why | To revert |
|---|---|---|---|---|
| **Δ1** | user columns (`status_manual`, `status_changed_at`, `rating`, `notes`) sit on `item` | they sit on **`item_user`**, 1 : 1 | Makes *"a metadata refresh can never destroy user data"* a property of the SQL: the sync's `UPSERT` **cannot name** a truth column. Measured cost: **+128 µs on a whole cold open**. [`05-truth-and-cache.md`](./05-truth-and-cache.md) §3 | Fold four columns back; the guarantee reverts to a convention |
| **Δ2** | `playtime_minutes INTEGER NOT NULL DEFAULT 0` | `INTEGER **NULL**` | §4.1. The seam already carries `*int`; the column was flattening it | `NOT NULL DEFAULT 0` + a capability lookup on every render |
| **Δ3** | `sync_run.status ∈ {ok, partial, failed, cancelled}` | adds **`running`** | A run must be recordable *before* it ends for `finished_at IS NULL` to mean *killed*. A `NULL` enum is a fifth state spelled as an absence | Drop the value; make the column nullable |
| **Δ4** | `sync_run.error` | `error_kind` (CHECKed) + `error_detail` | *"the classified failure, not a stack trace"* made enforceable. The classes are [`07-offline-contract.md`](../blueprint/07-offline-contract.md) §5's classifier, one for one | Merge into one free-text column |
| **Δ5** | *(nothing)* | **`item.last_synced_at`** | §4.2. `Z-05` renders a per-row age and `PARTIAL` makes a per-provider one dishonest | Drop it; `Z-05` block 3 loses its value or gains a lie |
| **Δ6** | `sort_title` *(one derived column)* | **`sort_title` + `search_title`** | They are different functions of the title. Sorting strips articles (`The Witness` under **W**); searching must not (typing `the` must find it — `Z-07` D-07-5) | One column; one of the two behaviours breaks |
| **Δ7** | *(nothing)* | **`STRICT`** on every table | SQLite's default type affinity would accept `'0'` in `playtime_minutes` and `'yesterday'` in a timestamp. Requires SQLite ≥ 3.37.0; **verified at source**: `modernc.org/sqlite` v1.57.0 → **3.53.3** | Remove the keyword; the CHECK constraints stay, the type discipline goes |
| **Δ8** | *(nothing)* | RFC 3339 **`GLOB` CHECKs** on every instant | Fixed 20-character UTC means lexical order **is** chronological order — so `ORDER BY started_at` needs no collation and no index expression | Drop the CHECKs; timestamp shape becomes a convention |
| **Δ9** | `owned_since TEXT` | `YYYY` \| `YYYY-MM` \| `YYYY-MM-DD`, **never an instant** | `Z-08`: *"Optional. A year is enough."* Storing `1998-01-01T00:00:00Z` fabricates five fields the player never gave | Accept a full timestamp and invent the precision |
| **Δ10** | `UNIQUE (provider_id, provider_ref)` inline | a **named** `CREATE UNIQUE INDEX` | Identical semantics; the index gets a name in `EXPLAIN QUERY PLAN` and can be dropped and rebuilt by a migration. An inline constraint becomes `sqlite_autoindex_item_1` | Inline it |
| **Δ11** | `game_uid`, `game_mood` | **`item_uid`**, **`item_mood`** | The table is `item`; a `game_uid` column on it reintroduces the exact word D5 spent a decision removing. [`09-erd.md`](../blueprint/09-erd.md) already says `item_uid`/`item_mood`; **four** other documents still say `game_uid`, and `09-erd.md` itself says `item_mood` in §2 and `game_mood` in §5. [`08-findings.md`](./08-findings.md) F-4 | Rename |
| **Δ12** | *(nothing)* | `item.provider_id` has **no FK** | §3.1 | Add it, and pick which wrong delete behaviour to ship |

**Nothing in ADR-0001 is contradicted by any of these.** Where this ticket believes a canon
document is *wrong* rather than merely under-specified — two places — it says so as a finding and
escalates rather than designing around it: [`08-findings.md`](./08-findings.md) F-1 and F-2.

---

## 7 · What is deliberately NOT in the file

Unchanged from [`09-erd.md`](../blueprint/09-erd.md) §3, restated because a physical model that
only lists what it contains teaches the wrong lesson.

| Not in `library.db` | Where | Why |
|---|---|---|
| **Credentials** | The `Vault` — OS keychain, or `credentials.json` mode `0600` | The library file is a thing the player is invited to back up, move and delete. A key inside it is a key in every backup and every support-ticket attachment |
| **Cover-art blobs** | The XDG **cache** directory | Disposable and re-fetchable, so the OS is allowed to delete it — the correct semantics. `metadata.cover_ref` holds a **local path**, and a `CHECK` refuses an `http(s)` URL, because nothing in Zerado renders from the network |
| **Any user identity, telemetry or analytics table** | Nowhere. It does not exist | There is nothing to phone home with, **by construction**. A reviewer can verify this by reading `sqlite_schema`: eleven objects, and none of them is a person |

---

## 8 · The state model, in the schema

[`05-state-machine.md`](../blueprint/05-state-machine.md) is canon and is not restated. What the
schema adds:

```sql
effective_status = COALESCE(
    status_manual,                                          -- the player's opinion
    CASE WHEN playtime_minutes > 0 THEN 'in_progress'       -- the ONE automatic transition
         ELSE 'not_started' END                             -- NULL playtime lands here too
)
```

| The requirement | How the schema carries it |
|---|---|
| Four states, no fifth | `CHECK (status_manual IN ('not_started','in_progress','zerado','abandoned'))` |
| *No opinion* is representable and distinct from every opinion | `status_manual` is **nullable** — which is what makes `Z-06`'s fifth item, *Clear override*, a different action from choosing `NOT STARTED` |
| `effective_status` is derived, never stored | It is a **view column**. There is no second writer, because there is no writer |
| A sync never changes a status the player set | The sync's statement **cannot name** `item_user`. Δ1 |
| `ZERADO` is never automatic | Nothing in the derivation can produce it |
| **A state may be derived for one source and manual for another** | This is the ticket's deliverable 4, and it needs **no media types** — §8.1 |
| `absent` is not a fifth state | It is `absent_since`, an orthogonal column on `item`. An absent row still carries one of the four |

### 8.1 · One state, two sources — already true, and not because of polymorphism

The ticket asks for *"the fact that a state may be **derived** for one type and **manually set**
for another"*, illustrated by a book having no page counter.

**Games already contain both cases**, and the mechanism that handles it is the capability model,
not a type system:

| Source | `Capabilities.Playtime` | `playtime_minutes` | `NOT STARTED → IN PROGRESS` |
|---|---|---|---|
| `steam` | `true` | `0` or a real count | **Derived.** A sync reporting `> 0` moves it |
| `physical` | `false` | **`NULL`** | **Impossible.** A cartridge has no telemetry, so all four states are manual, always |

That is *one state, two sources, per source* — the requirement, satisfied by a schema with one
media type in it. [`05-state-machine.md`](../blueprint/05-state-machine.md) §2.1 makes the same
argument and names the reason: *"a boolean `is_physical` would have forced the derivation to
special-case one value; a capability set makes the same code correct for every provider that will
ever exist."*

Δ2 is what puts that fact **in the file** rather than only in the code.

---

## 9 · The offline contract, at the storage layer

The ticket's deliverable 8: what is readable with the network off, what is stale, and how
staleness is recorded. Fully specified in
[`05-truth-and-cache.md`](./05-truth-and-cache.md) §6; the storage-level summary:

| | |
|---|---|
| **Readable offline** | **Everything in `library.db`.** Structurally, not as a feature: screens read the `Store` and the `Store` reads a file. There is no code path from a screen to a socket ([`06-data-seams.md`](../blueprint/06-data-seams.md) §1) |
| **How staleness is recorded** | Every network-derived value carries its own age in the same row, and every one of those columns is **`NOT NULL`**: `metadata.fetched_at`, `metadata.cover_fetched_at`, `price_quote.observed_at`, plus `item.last_synced_at` and `provider_connection.last_sync_at` |
| **Why `NOT NULL`** | [`07-offline-contract.md`](../blueprint/07-offline-contract.md) §4 forbids rendering a network-derived value without its age. A nullable age column is a code path that renders one without the other; a `NOT NULL` one is not |
| **What the schema refuses** | A cover reference that is a remote URL. `CHECK (cover_ref NOT LIKE 'http://%' AND cover_ref NOT LIKE 'https://%')` — *"nothing renders from the network"*, enforced |
| **What it cannot know** | Whether the machine is online. It never asks. A failure is classified into `sync_run.error_kind`, and the banner is raised by that failure and cleared by the next success ([`07-offline-contract.md`](../blueprint/07-offline-contract.md) §5) |
