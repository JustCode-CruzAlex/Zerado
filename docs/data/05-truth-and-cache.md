---
title: Zerado — truth and cache
discipline: DATA
doc-no: ZRD-DATA-05
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: concept-explainer
ticket: "#5"
---

# Truth and cache

> *"User facts — states, notes, ratings, mood overrides — are truth and must never be lost.
> Provider metadata — covers, synopsis, prices, community scores — is cache and must be
> re-fetchable, expirable and safe to discard. **Draw the line in the schema, not in a
> convention.**"*

This document draws it, says where it is a **table boundary** and where it is only a **column
value**, and gives the test that proves it holds.

---

## 1 · Three classes, not two

*Truth versus cache* is the right instinct and it is one class short. A Steam game's title is
neither: losing it costs nothing permanent, but it is not disposable either, because the library
must render it with the network off.

| Class | Who writes it | If it is lost | Where it lives |
|---|---|---|---|
| **TRUTH** | **Only the player** | **Gone forever.** No third party has a copy | `item_user` · `setting` · `item_mood.source='user'` · **every column of a `physical` row** |
| **REPLICA** | Only a sync | Costs one sync | `item` — `title`, `platform`, `playtime_minutes`, `last_played_at`, `owned_since`, `achievements_*` |
| **CACHE** | A metadata or price fetch | Costs bandwidth | `metadata` · `price_quote` · `item_mood.source='inferred'` · cover blobs (outside the file entirely) |

**The middle class is why the naive design fails.** A schema that puts truth and replica in one
table and calls the whole thing "the library" has no way to stop a sync from overwriting a status,
except by asking every future author to remember not to.

---

## 2 · The line, and where it is drawn

| Boundary | Kind of line | Strength |
|---|---|---|
| TRUTH ↔ REPLICA | **table boundary** — `item_user` vs `item` | **Structural.** The sync's statement cannot name a truth column |
| REPLICA ↔ CACHE | **table boundary** — `item` vs `metadata` / `price_quote` | **Structural.** A discard is `DELETE FROM <table>` and cannot reach `item` |
| TRUTH ↔ CACHE inside `item_mood` | **column value** — `source IN ('user','inferred')` | **Weaker, and named as weaker.** §5 |
| A `physical` row's REPLICA columns are really TRUTH | not a line at all — no sync exists | **Structural, by D1.** §3.2 |

**And the line holds in the other direction too**: deleting one `item` with `foreign_keys=ON`
removed its `item_user`, `metadata`, `price_quote` and `item_mood` rows and **left the other eight
items untouched**, with `foreign_key_check` clean. Cascade reaches every cache table and stops at
the row it was asked about.

---

## 3 · The split, and the argument for it

### 3.1 · What it makes impossible

`UpsertBatch` — the only statement a sync ever runs against the library — is:

```sql
INSERT INTO item (item_uid, item_type, provider_id, provider_ref, acquisition,
                  title, sort_title, search_title, platform,
                  playtime_minutes, last_played_at, owned_since, steam_appid,
                  achievements_total, achievements_unlocked,
                  last_synced_at, created_at, updated_at)
VALUES (…)
ON CONFLICT(provider_id, provider_ref) DO UPDATE SET
  title = excluded.title, …, last_synced_at = excluded.last_synced_at,
  absent_since = NULL                       -- it came back: cleared, silently
RETURNING id, (changes() … );
```

**There is no column in that statement that carries user truth, and there cannot be**, because
`status_manual`, `status_changed_at`, `rating` and `notes` are not columns of `item`. A future
author who tries to "fix" a status during a sync gets `no such column`, at compile time of the
statement, on their own machine.

Compare the alternative: everything on one table, and the guarantee is *"do not put those four
columns in the `DO UPDATE SET` list"* — a convention, in a `SET` list that grows every time a
provider reports a new field, reviewed by whoever is looking that week. Deliverable 7 asks for a
line in the schema. This is the schema saying no.

### 3.2 · Physical rows need no second mechanism

A `physical` row's `title`, `platform` and `owned_since` are the player's own typing —
irreplaceable, and they sit on `item` with the replica columns.

They need no guard, because **`physical` implements `Provider` and not `Syncer`** (ADR-0001 D1).
There is no sync for it. The statement that could overwrite those columns does not exist and
cannot be written, because there is no code path that produces `Item` values for a provider with
`Sync: false`.

**That is D1 paying for a storage guarantee**, which is worth noticing: the interface segregation
was argued on the grounds that *"a physical copy isn't a second-class row"*, and it turns out to be
what makes the truth rule enforceable for the one row class that has no other copy.

### 3.3 · What it costs, measured

| | |
|---|---|
| `Q-LIB-PAGE` | **+9.8 µs** |
| `Q-LIB-COUNTS` | +118 µs |
| **A whole cold open of 400 titles** | **+128 µs** |
| Plan shape | `SEARCH u USING INTEGER PRIMARY KEY (rowid=?)` — a rowid seek into a page already in cache |
| Storage | `item_user` is **12 KiB / 3 pages** at 400 rows |

[`evidence/measurements.md`](./evidence/measurements.md) §5 has the full table, including the
unsplit comparison schema it was measured against.

### 3.4 · The alternatives, and why they lost

| Considered | Rejected because |
|---|---|
| **All columns on `item`, plus a review convention** | Deliverable 7 forbids exactly this, and a convention is what fails on the busy week |
| **A `BEFORE UPDATE` trigger that refuses when the writer is a sync** | SQLite triggers cannot know the caller. Simulating it needs a connection-scoped `temp` table the trigger reads — which means every connection must create it or the trigger's subquery breaks, and the guarantee now depends on a setup step nobody sees |
| **A restricted `VIEW` with `INSTEAD OF` triggers for the sync to write through** | Genuinely close: the statement still cannot name a truth column, and there is no join on the read path. It loses on `Z-03`'s counts — the upsert needs per-row *new vs changed*, which comes from `RETURNING` on a real table and gets fiddly and lossy through an `INSTEAD OF` trigger. `Z-03` renders `12 new. 4 changed.` and the schema should not make that number hard to obtain |
| **Splitting three ways (truth / replica / identity)** | A third table for no third writer. The two writers are *the player* and *a sync*; that is the partition, and it has two parts |

### 3.5 · Is this the mistake D5 pruned, wearing a different hat?

Fair question, and the answer is no — for a reason that is worth stating rather than asserting.

D5 pruned a `media_game` extension table because *"a join that serves one type is machinery
without a purpose"* — a table justified by a **media type that does not exist**.

`item_user` is justified by a **constraint that exists today**, named in this ticket's own
non-functionals: *"User data is never destroyed by a metadata refresh."* It serves one writer that
exists (the player) against another writer that exists (a sync). If books never arrive, `item_user`
is still doing its job on the day it ships.

The test that distinguishes them: **remove the future, and does the structure still earn its
keep?** `media_game` does not. `item_user` does.

---

## 4 · The discard test

The property that makes *"cache"* mean something:

```sql
-- Everything below is disposable. Run it, and the product still works.
DELETE FROM metadata;
DELETE FROM price_quote;
DELETE FROM item_mood WHERE source = 'inferred';
-- and, outside the file:
rm -rf "$XDG_CACHE_HOME/zerado/covers"
```

**After that, every one of these must still be true:**

| | |
|---|---|
| Every game is still in the library, with its state | `item` and `item_user` were not named |
| Every note, rating and manual status is byte-identical | as above |
| Every hand-added game is intact, including its title and platform | as above |
| `Z-04`, `Z-05`, `Z-06`, `Z-07`, `Z-08`, `Z-09`, `Z-10`, `Z-11` render **identically** | they read no cache table |
| `Z-15` shows the designed **`not fetched`** tile, not a broken image | `Z-15` §3.3 |
| `Z-05` shows its designed **no-metadata** composition, not an error | `06-data-seams.md` §3.1 |
| Nothing is logged, warned about, or reported to the player | discarding cache is not an event |
| One metadata refresh restores all of it | it is re-fetchable by definition |

**This is a Phase 1 integration test, not a paragraph.** It runs the statements against a seeded
database and diffs `item` ⋈ `item_user` before and after: the diff must be empty. It is cheap, it
is decisive, and it is the only way *"safe to discard"* stays true after a year of changes.

**It has already been run**, against the full ladder plus
[`seed-minimal.sql`](./schema/fixtures/seed-minimal.sql) enriched with metadata, quotes and both
kinds of mood tag:

```
before   9 items · 9 metadata · 9 quotes · 13 mood tags (4 user)
   DELETE FROM metadata; DELETE FROM price_quote; DELETE FROM item_mood WHERE source='inferred';
after    9 items · 0 metadata · 0 quotes ·  4 mood tags (4 user)

item ⋈ item_user, field by field:  IDENTICAL
user-assigned mood tags:           preserved
```

**Its mirror, which matters as much:** a metadata refresh that returns garbage, or nothing, must
leave `item_user` byte-identical. The schema already makes that impossible (§3.1); the test proves
the schema is still shaped that way.

---

## 5 · The one place the line is weaker, said plainly

`item_mood` holds **both** classes in one table, split by `source`:

- `source = 'user'` — the player assigned it. **Truth.** It crosses the Phase 4 boundary.
- `source = 'inferred'` — a Phase 2 recommender guessed it. **Cache.** It never crosses.

This follows [`09-erd.md`](../blueprint/09-erd.md) §2, which specifies one table with a `source`
column, and this ticket does not overrule a canon table shape to win a symmetry argument.

**What that costs, stated:** the discard is `DELETE FROM item_mood WHERE source='inferred'`, and a
missing `WHERE` clause deletes the player's own tags. That is a real hazard and it is exactly the
class of thing the split removes elsewhere.

**What holds it up instead** — weaker than a table boundary, and named as such:

1. `CHECK (source IN ('user','inferred'))` and
   `CHECK ((source='inferred') = (confidence IS NOT NULL))` — an inferred row **must** carry a
   confidence and a user row **must not**, so the two are never confusable by inspection;
2. the partial index `item_mood_inferred … WHERE source='inferred'`, which makes *which rows are
   cache* a visible object in the schema rather than a fact in someone's head;
3. the discard test above, which covers this table too.

**Recommendation for the Phase 2 ticket**, so it is decided with a screen in hand rather than
speculatively now: if the recommender turns out to write inferred tags at any volume, split
`item_mood` into `item_mood_user` and `item_mood_inferred`. It is an `ALTER`-free migration — two
`CREATE TABLE`s and one `INSERT … SELECT` — and it converts this section's weakest guarantee into
the strongest one. **Not done now**, because Phase 2 has no recommender and this ticket does not
build for a feature that does not exist.

---

## 6 · The offline contract, at the storage layer

Deliverable 8: what is readable with the network off, what is stale, and how staleness is
recorded so the interface can be honest about it.

### 6.1 · What is readable offline: everything in the file, structurally

The classification is [`07-offline-contract.md`](../blueprint/07-offline-contract.md) §2's and is
not restated. What the storage layer contributes is **why it is not a feature that has to be
built**: a screen reads the `Store`; the `Store` reads a file; there is no code path from a screen
to a socket ([`06-data-seams.md`](../blueprint/06-data-seams.md) §1). Nine of the twelve Phase 1
screens are `WORKS` because of where their data comes from, not because anyone added offline
support to them.

### 6.2 · How staleness is recorded

**Every value that came from the network carries its own age, in the same row, in a `NOT NULL`
column.**

| Column | On | Ages |
|---|---|---|
| `item.last_synced_at` | `item` | this row's provider facts. **Per row** — a `PARTIAL` run leaves rows at different ages ([`01-erd.md`](./01-erd.md) §4.2) |
| `provider_connection.last_sync_at` | the connection | *"Last synced 3 days ago"* on `Z-04`'s banner |
| `sync_run.started_at` / `finished_at` | the run | `Z-03`'s history, and what `Z-04` reads when it needs the last **complete** one |
| **`metadata.fetched_at`** | cache | `NOT NULL` |
| **`metadata.cover_fetched_at`** | cache | `NOT NULL` **when a cover exists** — `CHECK ((cover_ref IS NULL) = (cover_fetched_at IS NULL))` |
| **`price_quote.observed_at`** | cache | `NOT NULL` |
| `price_quote.low_at` | cache | `NOT NULL` **when a low exists** — *"a low with no date is not information"* |

**Why `NOT NULL` rather than a rendering rule.** [`07-offline-contract.md`](../blueprint/07-offline-contract.md)
§4 says *"any value that came from the network is rendered with its age. Always"* and adds that it
is *"the rule most likely to be lost during a build, because dropping the age always makes the
layout tidier."*

A nullable age column **is** a code path that renders a value without one — the renderer has to
handle `NULL`, and the tidiest handling is to omit the age. A `NOT NULL` column deletes that
branch. The rule stops depending on the renderer remembering.

`item.last_synced_at` is the one exception and it is nullable **on purpose**: a hand-added row was
never synced, and `Z-05` D-05-8 says block 3 is **absent** on a hand-added copy rather than showing
an empty age. `NULL` there means *there is no such thing as a sync for this row*, which is a
different fact from *stale*.

### 6.3 · What the storage layer refuses

| | |
|---|---|
| A `cover_ref` that is a remote URL | `CHECK (cover_ref NOT LIKE 'http://%' AND cover_ref NOT LIKE 'https://%')`. *"Nothing renders from the network"* stops being a rule someone can forget |
| A metadata row with no age | `fetched_at NOT NULL` |
| A price with no age | `observed_at NOT NULL` |
| A completed sync that does not account for its own counts | `CHECK (items_seen >= items_new + items_changed)` |
| A failed sync with no classification | `CHECK (status IN ('ok','running','cancelled') OR error_kind IS NOT NULL)` |

### 6.4 · What the schema cannot enforce, and what covers it instead

Two rules are about **which** statement runs, and SQL cannot see that. They are tests, and they
are named here so they are not assumed to be constraints:

1. **Only a sync whose `status` is `ok` may set `absent_since`.**
   [`06-data-seams.md`](../blueprint/06-data-seams.md) §2.4 — in a truncated stream, *not
   returned* and *not reached* are indistinguishable. A `partial`, `failed` or `cancelled` run
   must tombstone nothing. **Test:** run a `partial` sync omitting half the library; assert
   `absent_since` is `NULL` on every row.
2. **A sync returning zero items is the private-profile case, never "the library was removed".**
   **Test:** run a sync that yields no items; assert no row is tombstoned and `sync_run.error_kind`
   is `empty`.

Both are in the Phase 1 suite. Neither is a `CHECK`, and pretending otherwise would be worse than
saying so.
