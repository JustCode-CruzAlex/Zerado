---
title: Zerado — the physical schema
discipline: DATA
doc-no: ZRD-DATA-02
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: implementation-plan
ticket: "#5"
---

# The physical schema

Every table, column, type, constraint, foreign key and default, with the reasoning attached, and
then the DDL verbatim.

**The executable truth is [`schema/migrations/`](./schema/migrations/).**
[`schema/schema.head.sql`](./schema/schema.head.sql) is generated from it and is never run —
see [`03-migrations.md`](./03-migrations.md) §8 for why that distinction is load-bearing.

---

## 1 · The connection contract

Nothing below is optional, and three of the five are silent failures if forgotten.

| Pragma | Value | What goes wrong without it |
|---|---|---|
| **`foreign_keys`** | `ON` | **SQLite defaults it OFF** — verified at source through `modernc.org/sqlite` v1.57.0: `PRAGMA foreign_keys` returns `0` on a fresh connection. Every `CASCADE` in this schema is inert until it is set, and orphan rows insert without complaint. It is **per connection**, not per database |
| **`journal_mode`** | `WAL` | Crash safety and reader/writer separation. Costs the `-wal`/`-shm` companions — §1.1 |
| **`busy_timeout`** | `5000` | A second connection (a test, a `sqlite3` session the player opened) turns into an instant `SQLITE_BUSY` instead of a short wait |
| **`synchronous`** | `NORMAL` | Safe under WAL, and it is the difference between a sync that writes 400 rows in milliseconds and one that fsyncs 400 times |
| **`trusted_schema`** | `OFF` | Defence in depth: refuses to run functions embedded in a schema, on a file the player is invited to copy from anywhere |

**Version precondition.** `STRICT` tables need SQLite ≥ **3.37.0**. The program checks
`sqlite_version()` at open and **refuses with a named error** rather than failing inside migration
0001 with `unrecognized token: "STRICT"`. Verified at source: `modernc.org/sqlite` **v1.57.0**
bundles SQLite **3.53.3**.

### 1.1 · One connection, and the one-file promise

The published promise is *one SQLite file*. WAL produces `-wal` and `-shm` while the process runs.
Measured through the real driver
([`evidence/measurements.md`](./evidence/measurements.md) §8), both companions are removed when
the **last** connection closes. So the discipline is exactly two rules:

1. **The process owns exactly one connection** to the library file.
2. **Clean shutdown runs `PRAGMA wal_checkpoint(TRUNCATE)` and then closes it.**

After a crash the companions survive and the next open recovers from them — WAL working as
designed. The promise is therefore *"a clean shutdown leaves one file"*, which is what
[`06-data-seams.md`](../blueprint/06-data-seams.md) §5.2 already says.

---

## 2 · `item` — the decisions, column by column

| Column | Decision | Why |
|---|---|---|
| `id` | `INTEGER PRIMARY KEY` | The rowid alias. Every child table's join is then a rowid seek, which is why the `item_user` split costs 9.8 µs and not a millisecond |
| `item_uid` | `TEXT NOT NULL`, shape-CHECKed, **not unique** | §4. A merge hint. The CHECK catches a truncated or upper-cased uuid at insert, where it is one row, instead of at merge time, where it is a mystery |
| `item_type` | `TEXT NOT NULL DEFAULT 'game'` + `CHECK (= 'game')` | The whole door. The `DEFAULT` means no INSERT anywhere has to mention it, which is what keeps the affordance from leaking into Phase 1 code |
| `provider_id`, `provider_ref` | `TEXT NOT NULL`, **`UNIQUE` together** | The sync's `ON CONFLICT` target. A `physical` row uses a generated UUID for `provider_ref`, so hand-entered rows share the identity mechanism rather than dodging it |
| `acquisition` | `CHECK IN ('digital','physical')` | A **fact about the copy**, not about the provider. A digital copy bought on a disc-based console is `physical`; a `physical` provider row is always `physical`. They are related and they are not the same column |
| `title` | `TEXT NOT NULL`, never normalised in place | What the source said. Every normalisation is a *derived* column, so a normaliser change is a regeneration and never a data loss |
| `sort_title` | `TEXT NOT NULL`, derived, **indexed** | §4.2. Folded, lower-cased, leading article and leading punctuation stripped |
| `search_title` | `TEXT NOT NULL`, derived, **not indexed** | §4.2 and [`04-indexing-plan.md`](./04-indexing-plan.md) §4 — a leading-wildcard `LIKE` cannot use a B-tree, and at 400 rows it does not need one |
| `platform` | `TEXT NOT NULL` | Half of `item_uid`, so it cannot be optional |
| `playtime_minutes` | `INTEGER **NULL**`, `>= 0` | [`01-erd.md`](./01-erd.md) §4.1. The one type change worth arguing for |
| `last_played_at` | `TEXT NULL`, RFC 3339 | `NULL` = not reported. Canon |
| `owned_since` | `TEXT NULL`, **three precisions, no clock** | `Z-08` says *"A year is enough."* Three `GLOB` alternatives accept `1998`, `1998-03`, `1998-03-14` and refuse `1998-01-01T00:00:00Z` |
| `steam_appid` | `TEXT NULL` | Provider-specific, and the one place this schema admits a provider by name. It is a **rendering** convenience (a store link); nothing branches on it |
| `achievements_total`, `_unlocked` | `INTEGER NULL`, `unlocked BETWEEN 0 AND total` | The cross-column CHECK is the one that matters: 21 of 43 is a fact, 44 of 43 is a bug, and catching it at insert is catching it in the sync that produced it |
| `absent_since` | `TEXT NULL`, partial-indexed | Tombstone. Set only by a **complete, successful** sync ([`06-data-seams.md`](../blueprint/06-data-seams.md) §2.4); the schema cannot enforce *which* sync writes it, so [`05-truth-and-cache.md`](./05-truth-and-cache.md) §7 makes it a test |
| `merged_into` | `INTEGER NULL` FK → `item(id)` `SET NULL`, `CHECK (<> id)` | Self-reference. The self-CHECK stops the one-row cycle; longer cycles are Phase 4's to prevent, and are named in [`06-phase4-seam.md`](./06-phase4-seam.md) §5 |
| `last_synced_at` | `TEXT NULL` | [`01-erd.md`](./01-erd.md) §4.2 |
| `created_at`, `updated_at` | `TEXT NOT NULL`, RFC 3339 | **Written by the Store from the injected clock, never by a trigger** — §3 |

---

## 3 · Triggers: allowed when pure structure, forbidden when they need a clock

Two triggers exist and a third deliberately does not.

| | What it does | Why it is allowed |
|---|---|---|
| `item_user_exists` | `AFTER INSERT ON item` → insert the 1 : 1 `item_user` row | Pure structure. No clock, no policy. It makes the join total, so no read can lose a row |
| `sync_run_retention` | `AFTER INSERT ON sync_run` → delete all but the newest 50 per provider | Pure structure, ordered by `id`. A retention rule in code is a retention rule that gets forgotten; in the schema it cannot be |

**The trigger that is not there: `updated_at`.** The obvious
`AFTER UPDATE ON item SET updated_at = strftime(...,'now')` would read **SQLite's own clock** and
bypass the injected `func() time.Time` the seams use
([`06-data-seams.md`](../blueprint/06-data-seams.md) §7 — *"the clock is not a seam; a
`func() time.Time` field on the structs that need it is enough for tests"*). A golden-file test
whose timestamps come from the wall clock is a test that fails on the machine that runs it
tomorrow. **Every timestamp in this schema is written by the `Store`, from the injected clock.**

That is the rule, stated so it survives: *a trigger may enforce shape; it may never read the
clock, and it may never decide policy.*

---

## 4 · The normalisers, specified

Handed to this ticket by [`13-handoffs.md`](../blueprint/13-handoffs.md) §4 — *"the `game_uid`
normalisation function's exact rules… needs testing against a real library's titles"*.

### 4.1 · `fold(s)` — the shared base

1. Unicode **NFD** decompose; drop every combining mark (category `Mn`); **NFC** recompose.
2. **Unicode case-fold** (`x/text/cases.Fold`, not `strings.ToLower` — Turkish dotted-I and
   German sharp-s both get this wrong under ASCII lowering).
3. Collapse runs of whitespace to one space; trim.

`Pokémon` → `pokemon` · `Ōkami HD` → `okami hd` · `ＺＥＲＡＤＯ` → `zerado`

### 4.2 · Two derived columns, because they are two different functions

| | `sort_title` | `search_title` |
|---|---|---|
| Base | `fold(title)` | `fold(title)` |
| Leading article | **stripped** (`the `, `a `, `an ` in v1) | **kept** |
| Leading punctuation | **stripped** | kept |
| Punctuation inside | kept | kept |
| Empty result | falls back to `fold(title)` — the schema refuses an empty one | — |
| Serves | `ORDER BY` (`Q-LIB-PAGE`) | `LIKE` (`Q-SEARCH`) |
| Indexed | **yes**, partially | **no** — [`04-indexing-plan.md`](./04-indexing-plan.md) §4 |

`The Witness` → sort `witness`, search `the witness`.
`[PROTOTYPE]` → sort `prototype`, search `[prototype]`.
`!!!` → sort `!!!` (the fallback), search `!!!`.

**Why two columns and not one.** `Z-07` D-07-5 fixes Phase 1 matching as a *case-insensitive,
accent-folded substring match on the title*, and `Z-04`'s order is *title A→Z*. If one column
served both, either `The Witness` sorts under **T** (wrong — canon says articles are stripped) or
typing `the` fails to find it (wrong — it is a substring of the title the player is looking at).
Two columns, forty bytes a row, both regenerable.

### 4.3 · Collation is done at **write** time, on purpose

SQLite has three built-in collations — `BINARY`, `NOCASE` (ASCII only), `RTRIM` — and **no
locale-aware collation without ICU**, which the pure-Go driver does not carry (ADR-0001 D2). So
locale-correct ordering cannot live in the `ORDER BY`.

It lives in `sort_title` instead: fold in Go with `x/text`, store the result, sort it `BINARY`.
The index is then a plain B-tree with no collation attached, which is why `Q-LIB-PAGE` costs the
same at 400 rows and at 10 000 ([`evidence/measurements.md`](./evidence/measurements.md) §3).

**The alternative, considered and rejected:** store an `x/text/collate` **sort key** as a `BLOB`.
It is strictly more correct — full CLDR tailoring. It is also opaque (nobody can read a row and
see why it sorted there), pinned to one locale, and must be regenerated when `x/text` changes its
table version. For English v1, with the accented titles that motivate the rule
(`Pokémon`, `Ōkami`) resolved by folding alone, that is a bad trade.

**When a second locale ships**, the mechanism is already built: bump `setting.normaliser_version`
and run the regeneration migration in [`03-migrations.md`](./03-migrations.md) §6. This is
recorded now so pt-BR does not arrive as a surprise — ADR-0001 D9's whole point.

### 4.4 · `uid_norm(s)` — the identity normaliser

`fold(s)`, then strip **all** punctuation to spaces, collapse, trim, then strip a leading article.
More aggressive than `sort_title` because two spellings of one title must collide *on purpose*.

```
NAMESPACE_ZERADO = uuidv5(NAMESPACE_DNS, "zerado.app") = 56640350-c577-5522-8dd0-30e65323adf8

"The Legend of Zelda: Breath of the Wild" / "Nintendo Switch"
   key      = "game|legend of zelda breath of the wild|nintendo switch"
   item_uid = 953fd060-48f0-5d38-8dc0-4ed1c2b83346

"Pokémon Red" / "Game Boy"
   key      = "game|pokemon red|game boy"
   item_uid = 168252b8-65e0-5a2d-afad-3d6e049ea90b

"Ōkami HD" / "PC"
   key      = "game|okami hd|pc"
   item_uid = e7364d01-aa87-5998-9d28-5a1aabb43149
```

These three are computed, not illustrative, and belong in the Phase 1 test corpus as fixed
expectations — because the value of `item_uid` is entirely that it is **the same on two
machines**, and the only way that stays true is a test that fails when the normaliser drifts.

---

## 5 · `item_view` — the read surface

```sql
COALESCE(u.status_manual,
         CASE WHEN i.playtime_minutes > 0 THEN 'in_progress' ELSE 'not_started' END)
    AS effective_status
```

**It flattens.** Verified with `EXPLAIN QUERY PLAN` on the real driver: every read through the
view resolves to `SCAN i USING INDEX item_shelf_order` + `SEARCH u USING INTEGER PRIMARY KEY
(rowid=?)`. The view costs nothing at plan time and the partial index is used *through* it.

Two derived companions ride along because a screen needs them and neither is stored:

| Column | For |
|---|---|
| `status_is_manual` | `Z-06`'s fifth item — *Clear override* — renders **only when `status_manual IS NOT NULL`** |
| `playtime_untracked` | `Z-04` renders `—` and `Z-05` renders `not tracked` when the source cannot know, versus `0h` when it reported zero |

**`effective_status` is not a generated column and not an indexed expression.** Both were
considered; [`04-indexing-plan.md`](./04-indexing-plan.md) §6 has the measurement that decides it
— a `GROUP BY` over 400 rows costs 439 µs and there is nothing to index, because every row is in
the answer. [`13-handoffs.md`](../blueprint/13-handoffs.md) §4 assigns that question here and
this is the answer: **no**, with the number.

---

## 6 · `item_user` — the truth table

Four columns, and two of them are unwritten in Phase 1.

| Column | Phase 1 | Note |
|---|---|---|
| `status_manual` | **written** | `Z-06` |
| `status_changed_at` | **written** | §6.1 |
| `rating` | present, unwritten | No Phase 1 screen renders a rating. **No range CHECK**: the scale is the Phase 2 rating screen's to choose, and pinning it now picks a fight this ticket cannot win |
| `notes` | present, unwritten | `Z-08` §3.1 rejects it as a field *"the column exists, but no Phase 1 screen renders it back"*. `CHECK (notes <> '')` — an empty string is not a note |

**Why keep two unwritten columns when this ticket removes an index that has no query?** Because
the costs are not the same, and the asymmetry is worth naming:

> **A nullable column is free — no write cost, no read cost, and `ALTER TABLE ADD COLUMN` later is
> also free. An index is not: every insert and every update pays it, forever.**

Both are named in [`09-erd.md`](../blueprint/09-erd.md), and `item_user` is where user truth
belongs. Leaving the columns out would risk a Phase 2 screen putting user data somewhere else.

### 6.1 · `status_changed_at` survives a clear, and Phase 4 depends on it

```sql
CHECK (status_manual IS NULL OR status_changed_at IS NOT NULL)
```

Note what this **does not** say: it is not a biconditional. `status_manual IS NULL` with a
**non-null** timestamp is legal and meaningful — *the player expressed an opinion, and then
cleared it.*

[`05-state-machine.md`](../blueprint/05-state-machine.md) §6 reads *"`NULL` when it has never been
set"*, which is right, and a biconditional would over-read it into *"NULL whenever
`status_manual` is NULL"*. That version has a Phase 4 bug:

> Device A: mark `ZERADO` on the 1st. Device B: mark `ZERADO` on the 1st (synced). Device A on the
> 20th: **clear the override.** If clearing drops the timestamp, the clear carries no time, loses
> every last-write-wins comparison, and **device B's stale `ZERADO` comes back** — the product
> silently undoing the player's most recent decision, which is the exact failure §2.2 calls *"the
> product failing at its purpose."*

So: **clearing an override is a state change and carries its own timestamp.** `NULL` on
`status_changed_at` means one thing only — *this player has never had an opinion about this game.*
`seed-minimal.sql` row 2 is that case, and it is there to be asserted.

---

## 7 · `sync_run` — the counts are a schema fact

```sql
CHECK (items_seen >= items_new + items_changed)
CHECK ((status = 'running') = (finished_at IS NULL))
CHECK (status IN ('ok','running','cancelled') OR error_kind IS NOT NULL)
```

`Z-03` renders `12 new. 4 changed. 131 unchanged.` and its own acceptance criterion is that
*"each set of counts sums"*. `unchanged` is `items_seen − items_new − items_changed`, so the
subtraction can go negative and print a nonsense line — unless the schema refuses it. It does.

`error_kind` is CHECKed to the six outcomes of
[`07-offline-contract.md`](../blueprint/07-offline-contract.md) §5's classifier
(`no_route` · `timeout` · `unauthorized` · `empty` · `server` · `other`), so
*"the classified failure, not a stack trace"* is enforced rather than intended. `error_detail`
holds the free text, and the third CHECK makes a `failed` or `partial` run without a
classification impossible.

**Retention: 50 runs per provider**, by trigger. Only the most recent is ever read
(`Q-LAST-SYNC`); the rest exist so a flaky provider is diagnosable. Fifty is months of history
and roughly four kilobytes. An unbounded log inside a file the player is told to back up is a slow
leak that nobody notices until it is large.

---

## 8 · Constraint verification

Every constraint below was executed against the real DDL. **`REFUSED`** means the schema rejected
it — the correct outcome.

| Attempted | Result |
|---|---|
| `item_type = 'book'` | REFUSED · `CHECK constraint failed: item_type = 'game'` |
| `acquisition = 'rental'` | REFUSED |
| `playtime_minutes = -1` | REFUSED |
| 11 achievements unlocked of 10 | REFUSED |
| `created_at = '25/08/2026'` | REFUSED |
| `owned_since = '1998-01-01T00:00:00Z'` (a fabricated clock) | REFUSED |
| `owned_since = '1998'` | **accepted** — correct |
| `playtime_minutes = NULL` on a `physical` row | **accepted** — correct |
| empty `title` | REFUSED |
| `item_uid = 'not-a-uuid'` | REFUSED |
| duplicate `(provider_id, provider_ref)` | REFUSED · `UNIQUE constraint failed` |
| `merged_into = id` | REFUSED |
| `status_manual` set with `status_changed_at` NULL | REFUSED |
| **clearing** `status_manual`, keeping the timestamp | **accepted** — and the timestamp remained. §6.1 |
| `cover_ref = 'https://cdn.example/x.png'` | REFUSED · the *"nothing renders from the network"* rule |
| `cover_ref = 'covers/ab/abcd1234.png'` | **accepted** — correct |
| `items_seen = 10` with `new = 8, changed = 5` | REFUSED |
| `status = 'ok'` with `finished_at` NULL | REFUSED |
| `status = 'failed'` with no `error_kind` | REFUSED |
| a `user` mood tag carrying a `confidence` | REFUSED |
| a price low with no `low_at` | REFUSED |
| a provider low **and** a local watermark on one row | REFUSED |
| orphan `item_mood` **without** `foreign_keys=ON` | **accepted** — which is exactly why §1 makes the pragma mandatory |
| orphan `item_mood` **with** `foreign_keys=ON` | REFUSED |
| `DELETE FROM item WHERE id=1` | cascaded `item_user`, `metadata`, `price_quote` **and** `item_mood`; left the other eight items untouched; `foreign_key_check` clean |

The exact statements are reproducible from [`schema/fixtures/`](./schema/fixtures/) and the
recipe in [`evidence/measurements.md`](./evidence/measurements.md) §9.

---

## 9 · The DDL, verbatim

Migration 0001 in full. Phase 2 and Phase 3 are
[`0002_phase2_enrichment.sql`](./schema/migrations/0002_phase2_enrichment.sql) and
[`0003_phase3_prices.sql`](./schema/migrations/0003_phase3_prices.sql); they are **not** applied by
a Phase 1 binary.

```sql
-- ============================================================================
-- Zerado · migration 0001 · Phase 1 core
-- ----------------------------------------------------------------------------
-- ticket   #5   ([BLUEPRINT][DATA] Zerado ERD and physical schema)
-- consumes ADR-0001 D2 (persistence) · D4 (Phase 4 sync boundary) · D5 (the door)
--          docs/blueprint/09-erd.md   — the conceptual model this makes physical
--          docs/blueprint/05-state-machine.md · 06-data-seams.md · 07-offline-contract.md
--
-- This file is EXECUTABLE TRUTH. `docs/data/schema/schema.head.sql` is generated
-- from it and is never executed. A fresh install runs this ladder; it never runs
-- a separate "create everything" script, because the two always drift.
--
-- Runs inside ONE transaction, opened by the migrator, together with the
-- INSERT INTO schema_migration that records it. There is no half-applied 0001.
--
-- Requires SQLite >= 3.37.0 (STRICT). Verified at source 2026-08-25:
-- modernc.org/sqlite v1.57.0 bundles SQLite 3.53.3.
-- ============================================================================


-- ----------------------------------------------------------------------------
-- schema_migration — the ladder's own bookkeeping
-- ----------------------------------------------------------------------------
-- `written_by` is why this is not just (version, applied_at). Z-11 promises to
-- name BOTH versions — the one that wrote the file and the one that is running.
-- A schema number alone cannot do that, and the binary that wrote a row is
-- exactly what is no longer around to ask. It cannot be backfilled later.
CREATE TABLE schema_migration (
    version     INTEGER PRIMARY KEY,
    applied_at  TEXT    NOT NULL,
    written_by  TEXT    NOT NULL,

    CHECK (version > 0),
    CHECK (written_by <> ''),
    CHECK (applied_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;


-- ----------------------------------------------------------------------------
-- item — identity, and what a PROVIDER says about it
-- ----------------------------------------------------------------------------
-- The core entity is `item`, not `games`, carrying an `item_type` CHECKed to
-- 'game'. That is the entire door-open affordance (ADR-0001 D5). Nothing else
-- here is generalised: playtime is playtime.
--
-- EVERY COLUMN IN THIS TABLE IS EITHER IDENTITY OR PROVIDER REPLICA.
-- Nothing the player typed lives here — that is `item_user`. The split is the
-- mechanism that makes "a metadata refresh can never destroy user data" a fact
-- about the SQL rather than a promise about the caller: the sync's UPSERT
-- cannot erase `status_manual`, because the statement cannot name it.
--
-- The one exception is a `physical` row, whose title/platform/owned_since ARE
-- the player's own typing. It needs no separate guard: `physical` implements
-- Provider and NOT Syncer (ADR-0001 D1), so no sync statement ever reaches it.
CREATE TABLE item (
    id                    INTEGER PRIMARY KEY,

    -- Phase 4 merge HINT. Indexed, never unique, never authoritative.
    -- uuidv5(NAMESPACE_ZERADO, item_type||'|'||uid_norm(title)||'|'||uid_norm(platform))
    -- NAMESPACE_ZERADO = uuidv5(NAMESPACE_DNS, 'zerado.app')
    --                  = 56640350-c577-5522-8dd0-30e65323adf8
    item_uid              TEXT    NOT NULL,

    -- The door. One column, one value, CHECKed to it.
    item_type             TEXT    NOT NULL DEFAULT 'game',

    provider_id           TEXT    NOT NULL,   -- 'steam' | 'physical' | …
    provider_ref          TEXT    NOT NULL,   -- Steam appid, or a UUID for physical
    acquisition           TEXT    NOT NULL,   -- 'digital' | 'physical'

    title                 TEXT    NOT NULL,   -- exactly as the source gave it
    sort_title            TEXT    NOT NULL,   -- derived: folded, articles stripped — ORDER BY
    search_title          TEXT    NOT NULL,   -- derived: folded, articles KEPT   — LIKE
    platform              TEXT    NOT NULL,

    -- NULL = this source does not report playtime at all. 0 = it reported zero.
    -- The provider seam already carries this distinction as `Item.Playtime *int`
    -- (06-data-seams §2.1); a NOT NULL DEFAULT 0 column would flatten it at the
    -- storage boundary and Z-04 could no longer render `—` for a cartridge and
    -- `0h` for an unplayed Steam title from the row alone.
    playtime_minutes      INTEGER NULL,
    last_played_at        TEXT    NULL,       -- NULL = not reported, NOT never played
    owned_since           TEXT    NULL,       -- reduced precision allowed: YYYY | YYYY-MM | YYYY-MM-DD

    steam_appid           TEXT    NULL,
    achievements_total    INTEGER NULL,
    achievements_unlocked INTEGER NULL,

    -- Tombstone. Set on the first COMPLETE, SUCCESSFUL sync that omits the row;
    -- cleared silently when it returns. NEVER a reason to delete (06 §2.4).
    absent_since          TEXT    NULL,

    -- Phase 4 merge target, so joining two rows never rewrites a primary key.
    merged_into           INTEGER NULL REFERENCES item(id) ON DELETE SET NULL,

    -- PER-ROW sync age. Z-05 block 3 renders `LAST SYNCED` for THIS game, and a
    -- PARTIAL run leaves some rows fresh and others untouched — so a single
    -- provider-level timestamp would be a lie about every row the run never
    -- reached. Stamped on every upsert that wrote this row.
    last_synced_at        TEXT    NULL,

    created_at            TEXT    NOT NULL,
    updated_at            TEXT    NOT NULL,

    CHECK (item_type = 'game'),
    CHECK (acquisition IN ('digital','physical')),
    CHECK (title <> '' AND sort_title <> '' AND search_title <> '' AND platform <> ''),
    CHECK (provider_id <> '' AND provider_ref <> ''),
    CHECK (length(item_uid) = 36
           AND item_uid GLOB '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*'),
    CHECK (playtime_minutes IS NULL OR playtime_minutes >= 0),
    CHECK (achievements_total IS NULL OR achievements_total >= 0),
    CHECK (achievements_unlocked IS NULL
           OR (achievements_total IS NOT NULL
               AND achievements_unlocked BETWEEN 0 AND achievements_total)),
    CHECK (merged_into IS NULL OR merged_into <> id),

    -- Instants are RFC 3339, UTC, seconds precision, EXACTLY 20 characters, so
    -- that lexical order is chronological order and no index needs a collation.
    CHECK (last_played_at IS NULL OR last_played_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (absent_since   IS NULL OR absent_since   GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (last_synced_at IS NULL OR last_synced_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (created_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (updated_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),

    -- An acquisition date has no clock. A player who types "1998" is telling the
    -- truth at year precision; storing 1998-01-01T00:00:00Z would fabricate the
    -- other five fields. Z-08: "Optional. A year is enough."
    CHECK (owned_since IS NULL
           OR owned_since GLOB '[0-9][0-9][0-9][0-9]'
           OR owned_since GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]'
           OR owned_since GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]')
) STRICT;


-- ----------------------------------------------------------------------------
-- item_user — what the PLAYER typed. Truth. Never re-derivable. Never lost.
-- ----------------------------------------------------------------------------
-- 1:1 with item, materialised by trigger so every read joins and no read can
-- lose a row. This table is the Phase 4 sync payload, the "export my work"
-- selection, and the thing a corrupt `item` table does not take with it.
CREATE TABLE item_user (
    item_id           INTEGER PRIMARY KEY REFERENCES item(id) ON DELETE CASCADE,

    -- NULL = the player has never expressed an opinion; derive on read.
    status_manual     TEXT NULL,

    -- The timestamp that decides a Phase 4 conflict. NULL only while the player
    -- has NEVER expressed an opinion. CLEARING an override is itself a change
    -- and KEEPS its timestamp — without that, a clear can never win a
    -- last-write-wins merge against an older set on another device, and the
    -- override the player deleted would silently come back.
    status_changed_at TEXT NULL,

    -- Present from 0001 because they are user truth and this is the table for
    -- it. NO Phase 1 screen writes either one (Z-08 §3.1 rejects `notes`; no
    -- screen renders a rating). `rating` deliberately carries NO range CHECK:
    -- the scale is the Phase 2 rating screen's to choose, and constraining it
    -- now would pick a fight this ticket cannot win.
    rating            INTEGER NULL,
    notes             TEXT    NULL,

    CHECK (status_manual IS NULL
           OR status_manual IN ('not_started','in_progress','zerado','abandoned')),
    CHECK (status_manual IS NULL OR status_changed_at IS NOT NULL),
    CHECK (status_changed_at IS NULL
           OR status_changed_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (notes IS NULL OR notes <> '')
) STRICT;

-- Pure structure, no clock: allowed. (A trigger that stamped a timestamp would
-- read SQLite's own `now` and bypass the injected clock the seams use for
-- tests — 06-data-seams §7. Every timestamp in this schema is written by the
-- Store from that clock, never by a trigger.)
CREATE TRIGGER item_user_exists AFTER INSERT ON item
BEGIN
    INSERT INTO item_user(item_id) VALUES (new.id);
END;


-- ----------------------------------------------------------------------------
-- provider_connection — a credential RELATIONSHIP, and no credential
-- ----------------------------------------------------------------------------
-- The keys are in the Vault (ADR-0001 D2·3). `account_ref` is an identifier —
-- a Steam ID, a GOG username — never a secret.
--
-- NOTE THE ABSENT FOREIGN KEY. `item.provider_id` does NOT reference this
-- table. 09-erd draws item→connection as `N:1` and it is a real association,
-- but it must not be a referential one: a connection is a relationship the
-- player can END, and the library outlives it. An FK would force the choice
-- between deleting a library on disconnect and forbidding disconnect while
-- items exist. Both are wrong, so there is no FK, and the association is
-- resolved in code against the compiled-in provider registry — which D1 makes
-- the authority anyway (`Capabilities`, never the database).
CREATE TABLE provider_connection (
    provider_id      TEXT PRIMARY KEY,
    account_ref      TEXT NULL,
    connected_at     TEXT NOT NULL,
    last_sync_at     TEXT NULL,
    last_sync_status TEXT NULL,

    CHECK (provider_id <> ''),
    CHECK (last_sync_status IS NULL
           OR last_sync_status IN ('ok','partial','failed','cancelled')),
    CHECK (connected_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (last_sync_at IS NULL OR last_sync_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;


-- ----------------------------------------------------------------------------
-- sync_run — the honest history
-- ----------------------------------------------------------------------------
-- `status = 'running'` exists because a run must be recordable BEFORE it
-- terminates for a kill to be observable — which is exactly the fact 09-erd
-- wants `finished_at IS NULL` to carry. Without it the column would have to be
-- NULL while running, and a NULL enum is a fifth state spelled as an absence.
--
-- `error` is split into a CHECKed classification plus free detail, so
-- "the classified failure, not a stack trace" is enforced rather than intended.
-- The classes are 07-offline-contract §5's classifier, one for one.
--
-- Deleting a connection CASCADES its runs: sync history is bookkeeping, and a
-- reconnect must not inherit a previous relationship's "last synced" age. No
-- item is touched — items are the player's library, not the connection's.
CREATE TABLE sync_run (
    id            INTEGER PRIMARY KEY,
    provider_id   TEXT    NOT NULL REFERENCES provider_connection(provider_id) ON DELETE CASCADE,
    started_at    TEXT    NOT NULL,
    finished_at   TEXT    NULL,
    status        TEXT    NOT NULL DEFAULT 'running',
    items_seen    INTEGER NOT NULL DEFAULT 0,
    items_new     INTEGER NOT NULL DEFAULT 0,
    items_changed INTEGER NOT NULL DEFAULT 0,
    error_kind    TEXT    NULL,
    error_detail  TEXT    NULL,

    CHECK (status IN ('running','ok','partial','failed','cancelled')),
    CHECK ((status = 'running') = (finished_at IS NULL)),
    CHECK (items_seen >= 0 AND items_new >= 0 AND items_changed >= 0),
    -- Z-03 renders "12 new. 4 changed. 131 unchanged." and its acceptance
    -- criterion is that each set of counts sums. Made a schema fact.
    CHECK (items_seen >= items_new + items_changed),
    CHECK (error_kind IS NULL
           OR error_kind IN ('no_route','timeout','unauthorized','empty','server','other')),
    CHECK (error_kind IS NOT NULL OR error_detail IS NULL),
    CHECK (status IN ('ok','running','cancelled') OR error_kind IS NOT NULL),
    CHECK (started_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (finished_at IS NULL OR finished_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (finished_at IS NULL OR finished_at >= started_at)
) STRICT;

-- Retention, made structural. Only the most recent run is ever read; an
-- unbounded log inside a file the player is told to back up is a slow leak.
-- 50 runs per provider is months of history and a few kilobytes.
-- Pure structure, no clock: allowed.
CREATE TRIGGER sync_run_retention AFTER INSERT ON sync_run
BEGIN
    DELETE FROM sync_run
     WHERE provider_id = new.provider_id
       AND id <= (SELECT id FROM sync_run
                   WHERE provider_id = new.provider_id
                   ORDER BY id DESC LIMIT 1 OFFSET 50);
END;


-- ----------------------------------------------------------------------------
-- setting — everything Z-09 writes. Local truth; never crosses to Phase 4.
-- ----------------------------------------------------------------------------
-- Device-local by nature (theme, audio, image capability, the normaliser
-- version). "Only what the player typed crosses" is about the LIBRARY; a
-- preference is about THIS machine.
CREATE TABLE setting (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    CHECK (key <> '')
) STRICT;


-- ----------------------------------------------------------------------------
-- item_view — the read surface. `effective_status` is DERIVED, never stored.
-- ----------------------------------------------------------------------------
-- effective_status = status_manual ?? derive(playtime_minutes, capabilities)
--
-- The capability lookup disappears here, and that is the point of making
-- `playtime_minutes` nullable: NULL (this source cannot know) and 0 (it knows,
-- and it is zero) both fall to 'not_started', which is exactly what
-- derive() returns in both cases — so the derivation is TOTAL over the row and
-- needs nothing from compiled-in code.
--
-- ZERADO is never derived. `in_progress` is the only automatic transition in
-- the product (05-state-machine §2, §4).
CREATE VIEW item_view AS
SELECT
    i.id, i.item_uid, i.item_type, i.provider_id, i.provider_ref, i.acquisition,
    i.title, i.sort_title, i.search_title, i.platform,
    i.playtime_minutes, i.last_played_at, i.owned_since,
    i.steam_appid, i.achievements_total, i.achievements_unlocked,
    i.absent_since, i.merged_into, i.last_synced_at, i.created_at, i.updated_at,
    u.status_manual, u.status_changed_at, u.rating, u.notes,
    COALESCE(
        u.status_manual,
        CASE WHEN i.playtime_minutes > 0 THEN 'in_progress' ELSE 'not_started' END
    )                             AS effective_status,
    (u.status_manual IS NOT NULL) AS status_is_manual,
    (i.playtime_minutes IS NULL)  AS playtime_untracked
FROM item i
JOIN item_user u ON u.item_id = i.id;


-- ----------------------------------------------------------------------------
-- Indexes. Every one is justified by a NAMED query in
-- docs/data/04-indexing-plan.md. One is not, and it says so there.
-- ----------------------------------------------------------------------------

-- Q-UPSERT   the sync's ON CONFLICT target, and 09-erd's uniqueness constraint.
CREATE UNIQUE INDEX item_provider_identity ON item(provider_id, provider_ref);

-- Q-LIB-PAGE · Q-LIB-COUNTS · Q-DECK-PAGE
-- Partial, over exactly the default row set: shown = not absent, not merged.
CREATE INDEX item_shelf_order ON item(sort_title)
    WHERE absent_since IS NULL AND merged_into IS NULL;

-- Q-ABSENT-EXISTS · Q-ABSENT-COUNT · Q-ABSENT-LIST
-- Usually an EMPTY index. That is the point: the common case (nothing absent)
-- becomes an empty-B-tree probe instead of a full-table scan, on a query
-- Z-07 runs to decide whether the fifth chip renders at all.
CREATE INDEX item_absent ON item(absent_since) WHERE absent_since IS NOT NULL;

-- ADR-0001 D4 mandates this index in Phase 1. It has NO Phase 1 query.
-- Recorded as such in 04-indexing-plan.md §5 rather than silently kept or
-- silently dropped.
CREATE INDEX item_uid_lookup ON item(item_uid);

-- Q-LAST-SYNC   read by Z-04's banner and every Z-03 terminal state.
CREATE INDEX sync_run_recent ON sync_run(provider_id, started_at DESC);
```
