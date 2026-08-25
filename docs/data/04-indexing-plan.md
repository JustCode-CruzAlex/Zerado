---
title: Zerado — the indexing plan
discipline: DATA
doc-no: ZRD-DATA-04
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: implementation-plan
ticket: "#5"
---

# The indexing plan

> **The rule this document is held to:** *"Every index is justified by a named query it serves —
> an index with no query behind it is removed."*
>
> Eight indexes exist. **Seven** are named against a query traced to a screen spec. **One** is
> not, and §5 says so rather than inventing a justification for it. Six candidate indexes were
> considered and removed, and §4 names each with the reason.

Every number here comes from [`evidence/measurements.md`](./evidence/measurements.md), taken on
`modernc.org/sqlite` v1.57.0 (SQLite 3.53.3) — the driver ADR-0001 D2 selects. Nothing is
estimated.

---

## 1 · The bar

> *"A cold open on a 400-title library must be instant. The indexing plan is judged against that."*

**Measured: 842 µs** — a new connection, the migration check, the root screen's page and its
summary counts. Process start dominates it by roughly an order of magnitude.

The **per-keystroke** budget in filter mode, which is the other place a person can feel a
database, is **313 µs at 400 titles** — 1.9 % of one 60 fps frame.

---

## 2 · The named queries

Every query the Phase 1 screens actually issue. The column is the screen spec that asks for it,
because *"design the indexes against what these screens actually ask for"* is the instruction, and
a query nobody named is a query nobody has.

| Query | Asked for by | Shape | Index used | @400 |
|---|---|---|---|---|
| **`Q-LIB-PAGE`** | `Z-04` §3.1 rows 4–15 · `ROWS 4–15 of 247` | `WHERE item_type='game' AND absent_since IS NULL AND merged_into IS NULL ORDER BY sort_title LIMIT 12 OFFSET n` | `item_shelf_order` | **47.3 µs** |
| **`Q-LIB-COUNTS`** | `Z-04` §3 summary · `05-state-machine` §7 | `GROUP BY effective_status` over the same set | `item_shelf_order` | 439 µs |
| **`Q-SEARCH`** | `Z-07` **D-07-5** | `… AND search_title LIKE '%'‖?‖'%' ORDER BY sort_title LIMIT 12` | `item_shelf_order` *(for the order, not the match)* | 112.8 µs |
| **`Q-SEARCH-RATIO`** | `Z-07` §3.3 · `23 of 247` | `COUNT(*)` over the match | `item_shelf_order` | 200.4 µs |
| **`Q-TOTAL`** | `Z-07` **D-07-9** — the denominator once any row is absent | `COUNT(*)` over the shown set | `item_shelf_order` | 149.5 µs |
| **`Q-ABSENT-EXISTS`** | `Z-07` **F15 / D-07-8** — does the fifth chip render at all? | `EXISTS(SELECT 1 FROM item WHERE absent_since IS NOT NULL)` | `item_absent` **(covering)** | **14.1 µs** |
| **`Q-ABSENT-COUNT`** | `Z-04` §10.3 · **D-04-8** — the summary's `<n> absent` | `COUNT(*) WHERE absent_since IS NOT NULL` | `item_absent` (covering) | **12.6 µs** |
| **`Q-ABSENT-LIST`** | `Z-07` **F16** §5.1 — the row set swaps | `… WHERE absent_since IS NOT NULL ORDER BY sort_title` | `item_absent` | 40.0 µs |
| **`Q-DETAIL`** | `Z-05` | `WHERE id = ?` | PK ×2 | 45.5 µs |
| **`Q-DECK-PAGE`** | `Z-15` §3.4 — the grid | `Q-LIB-PAGE` with the grid's page size | `item_shelf_order` | = `Q-LIB-PAGE` |
| **`Q-LAST-SYNC`** | `Z-04` banners B1/B2 · `Z-03` §3.5 | `WHERE provider_id=? ORDER BY started_at DESC LIMIT 1` | `sync_run_recent` | **15.9 µs** |
| **`Q-UPSERT`** | `Z-03` · `Store.UpsertBatch` | `INSERT … ON CONFLICT(provider_id, provider_ref) DO UPDATE … RETURNING` | `item_provider_identity` | — |
| **`Q-SETSTATUS`** | `Z-06` | `UPDATE item_user SET … WHERE item_id = ?` | PK | — |
| **`Q-ADD-MANUAL`** | `Z-08` | `INSERT INTO item …` (the trigger makes the `item_user` row) | — | — |
| **`Q-SETTING`** | `Z-09` | `WHERE key = ?` | PK | — |
| `Q-MOOD-FILTER` | Phase 2 · `Z-18` | `WHERE mood_id = ?` | `item_mood_by_mood` | Phase 2 |
| `Q-TONIGHT` | Phase 2/3 · `Z-18` — *"what should I play tonight"*, under a budget cap | `Q-LIB-PAGE`'s set ⋈ `item_mood` ⋈ `price_quote` | `item_mood_by_mood` + `price_quote` PK | Phase 2/3 |
| `Q-WATCHLIST` | Phase 3 · `Z-20` | every quote for the watched set | `price_quote` PK | Phase 3 |
| `Q-CACHE-DISCARD` | *Settings* · the discard test | `DELETE FROM item_mood WHERE source='inferred'` | `item_mood_inferred` | Phase 2 |

**Note what is not in this table:** no query sorts by last played, filters by platform, or
searches a *sinopse*. §4 explains why, and it is not an oversight.

---

## 3 · The eight indexes

| # | Index | Serves | Evidence |
|---|---|---|---|
| **1** | `item_provider_identity` **UNIQUE** `(provider_id, provider_ref)` | `Q-UPSERT` | It is also `09-erd.md`'s uniqueness constraint, so it would exist even with no query. Named explicitly rather than inline so it appears in `EXPLAIN QUERY PLAN` and can be rebuilt by a migration (Δ10) |
| **2** | `item_shelf_order` `(sort_title)` **`WHERE absent_since IS NULL AND merged_into IS NULL`** | `Q-LIB-PAGE`, `Q-DECK-PAGE`, `Q-LIB-COUNTS`, `Q-SEARCH`, `Q-SEARCH-RATIO`, `Q-TOTAL` | §3.1 — the most valuable object in the file |
| **3** | `item_absent` `(absent_since)` **`WHERE absent_since IS NOT NULL`** | `Q-ABSENT-EXISTS`, `Q-ABSENT-COUNT`, `Q-ABSENT-LIST` | §3.2 |
| **4** | `item_uid_lookup` `(item_uid)` | **nothing in Phase 1** | §5 |
| **5** | `sync_run_recent` `(provider_id, started_at DESC)` | `Q-LAST-SYNC` | `SEARCH sync_run USING INDEX sync_run_recent`, 15.9 µs, flat in history size |
| 6 | `mood_tag_key` **UNIQUE** `(key)` *(Phase 2)* | the vocabulary's uniqueness | A duplicate `key` would silently split a mood in two |
| 7 | `item_mood_by_mood` `(mood_id, item_id)` *(Phase 2)* | `Q-MOOD-FILTER`, `Q-TONIGHT` | The table's own PK covers *this game's tags*; this covers *this tag's games*, which is the direction the recommender reads |
| 8 | `item_mood_inferred` `(item_id)` **`WHERE source='inferred'`** *(Phase 2)* | `Q-CACHE-DISCARD` | Makes the discard a partial-index scan rather than a full-table one, and makes *which rows are cache* visible in the schema |

### 3.1 · Why `item_shelf_order` is partial, and what it buys

The predicate is exactly the default row set: **shown = not absent, not merged**. SQLite matches a
partial index when the query's `WHERE` contains the index's terms, so `Q-LIB-PAGE` uses it
directly, ordered, with no sort step.

`item_type = 'game'` is **not** in the index predicate. Every row satisfies it (the `CHECK`
guarantees that), so it would narrow nothing — but it stays in every *query*, because ADR-0001 D5
accepts *"every query carries an `item_type` predicate it does not yet need"* as the price of the
door.

| | with | without | |
|---|---|---|---|
| `Q-LIB-PAGE` @ 400 | **47.3 µs** | 189 µs | 4.0× |
| `Q-LIB-PAGE` @ 10 000 | **46.8 µs** | 3 829.6 µs | **81.8×** |

```
with     |--SCAN i USING INDEX item_shelf_order        stops after 15 rows
without  |--SCAN i  … `--USE TEMP B-TREE FOR ORDER BY  reads all, then sorts all
```

**The number that matters is the flatness, not the ratio.** 47.3 µs at 400 and 46.8 µs at 10 000
is the same query doing the same work whatever the library size, because `LIMIT 12` against an
ordered index stops early. Without it, the root screen's cost is proportional to the library —
the shape that turns *instant* into *it used to be instant*, quietly, on the machine of the one
player with a big library.

Storage: **16 KiB, 4 pages** at 400 rows.

### 3.2 · `item_absent` is an index that is usually empty, and that is the point

`Z-07`'s fifth chip renders **only when at least one row is absent** (D-07-8), so
`Q-ABSENT-EXISTS` runs every time the filter bar is composed. For almost every player, almost
always, the answer is **no**.

Without the index that answer costs a **full table scan** — the worst case is the one that happens
most. With it, the plan is `SEARCH item USING COVERING INDEX item_absent (absent_since>?)` against
a B-tree with **zero entries**: **14.1 µs at 400 rows, 12.5 µs at 10 000.**

Storage: **4 KiB, one page**, because a partial index over a condition nothing satisfies stores
nothing.

It also serves `Q-ABSENT-COUNT` (`Z-04`'s `<n> absent` suffix, D-04-8) and `Q-ABSENT-LIST`
(`Z-07` F16). Three named queries, one page.

---

## 4 · Considered and removed

The ticket names six queries. **Three of them are not Phase 1 queries**, and building indexes for
them would be building against the ticket rather than against the screens — which is the
instruction this plan was given.

| Candidate index | Named by | Removed because |
|---|---|---|
| **`item(search_title)`** | *"full-text search over title"* | `Z-07` **D-07-5** fixes Phase 1 matching as a **substring** match. `LIKE '%x%'` has a leading wildcard and **cannot use a B-tree** — the index would be written on every insert and read never. Measured cost of not having it: `Q-SEARCH` **112.8 µs** at 400 rows and **72.7 µs** at 10 000 (the `LIMIT` bounds it), `Q-SEARCH-RATIO` 200 µs / 8.4 ms |
| **An FTS5 table over `title` + `sinopse`** | *"full-text search over title and **sinopse**"* | Two reasons and either is sufficient. **(1)** `sinopse` does not exist in Phase 1 — it arrives in `metadata` with Phase 2. **(2)** `Z-07` D-07-5 chose substring matching **deliberately**, because a fuzzy matcher *"makes the empty-result diagnostic much harder to write honestly"* (§11.5's `search "souls" 3 games` stops being a fact the player can reproduce). FTS5 is a shadow-table cluster and a write-path cost, for a query no screen issues. **The trigger for revisiting:** Phase 2 shipping a search that includes *sinopse*, or a library past ~100 000 titles. Neither is now |
| **`item(last_played_at DESC)`** | *"sort by last played"* | [`07-offline-contract.md`](../blueprint/07-offline-contract.md) §2 fixes Phase 1's ordering: *"**fixed**: title A→Z, with **no sort control** in Phase 1… an on-screen sort indicator would imply a control that does not exist."* There is no query. When Phase 2 adds the control, this index arrives with it |
| **`item(platform)`** | *"filter by platform"* | `Z-07` has no platform facet. The nearest thing — the `source` facet — is `Z-07` **§19 item 2, open for the founder**: *"the diagnostic block and the ratio slot both support the facet; nothing reaches it."* An index for a facet with no control is an index for a query nobody can issue |
| **`item_user(status_manual)`** | *"filter by state"* | The state chips filter on **`effective_status`**, which is derived across the join and is not a column anywhere. And it would be a poor index regardless: four values over 400 rows is ~25 % selectivity, where SQLite correctly prefers a scan. `Q-LIB-COUNTS` at 439 µs is the measured cost of not having one |
| **`item(merged_into)`** | *(FK child column)* | `ON DELETE SET NULL` scans for children on each delete. Deletes are player-initiated, rare, and over 400 rows. An index that is written on every insert to speed up an operation that happens a few times a year is a bad trade |
| **`price_quote(observed_at)`** | *(staleness)* | No screen asks *which quotes are stalest*. The age is **rendered**, never filtered on. If `Z-20` grows a *refresh-the-stalest-first* behaviour, that is the query that earns it |

---

## 5 · The one index with no query: `item_uid_lookup`

**Stated plainly rather than justified away.** ADR-0001 **D4** requires it:

> *"`game_uid` … assigned at insert, **indexed but not unique**." … "Adding this in Phase 1 costs
> one column and one index."*

By this document's own rule it would be removed. It is kept, because **D4 is founder-ratified
canon and reversing it is not this ticket's call** — and because designing around a ratified
decision quietly is worse than carrying a 24 KiB index.

What the reviewer and the founder should have in front of them:

| | |
|---|---|
| **Phase 1 queries served** | **None.** No Phase 1 screen looks a game up by `item_uid` |
| **Phase 4 queries served** | The merge-candidate lookup, which is the whole reason the *column* exists |
| **Cost, measured** | **24 KiB / 6 pages at 400 rows — the largest index in the file**, larger than the one serving the hottest query in the product. Plus one B-tree write per insert |
| **One incidental use, measured** | SQLite elects it as the smallest covering index for a bare `SELECT COUNT(*) FROM item`. That is real but weak: the count the screens actually need is `Q-TOTAL`, which uses `item_shelf_order` |
| **The counter-argument for keeping it** | An index is the *cheapest thing in this schema to add later* — `CREATE INDEX` on 400 rows is instant and needs no rebuild. The **column** is the expensive part, and the column is not in question |

**The recommendation is to keep it**, on canon. **The alternative, if the founder wants the rule
applied without exception:** move `CREATE INDEX item_uid_lookup` into the Phase 4 migration and
keep the column in 0001. That costs nothing and loses nothing, because unlike the column an index
can be built from data that is already there.

*Recorded as [`08-findings.md`](./08-findings.md) **F-3**, for the founder, not decided here.*

---

## 6 · `effective_status`: no generated column, no expression index

[`13-handoffs.md`](../blueprint/13-handoffs.md) §4 assigns this question here — *"whether derived
status warrants a generated column or an expression index… needs real query shapes and a 400+ row
measurement."* Here is the measurement and the answer.

**No.** Three reasons, in order of weight:

1. **There is nothing to index.** The expensive query is `Q-LIB-COUNTS`, a `GROUP BY` over the
   whole shown set. Every row is in the answer, so an index cannot skip anything — it can only
   avoid re-computing a `COALESCE` over two integers already on the page.
2. **It is impossible as stated.** `effective_status` spans `item` **and** `item_user`. SQLite
   generated columns may not reference another table, and an expression index may not either. A
   generated column would require folding the split back (Δ1) and losing the guarantee it buys.
3. **The measurement does not ask for it.** 439 µs at 400 rows, recomputed only when the row set
   changes — not per keystroke, because filter mode replaces the summary with the ratio
   (`Z-07` D-07-1).

[`05-state-machine.md`](../blueprint/05-state-machine.md) §6's conclusion holds unchanged: *"a
stored derived value has two ways to be wrong and a derived one has none."*

**The threshold, so a future ticket does not have to re-derive it.** `Q-LIB-COUNTS` is 17.0 ms at
10 000 rows and grows linearly, so it becomes a visible hitch somewhere past **~8 000 titles**.
The fix at that point is **not** an index: it is maintaining the four counts incrementally in the
view model and invalidating them on write. Recorded, and **not built**, because building it now is
a cache with two ways to be wrong, for a library nobody has.

---

## 7 · The Phase 1 benchmark this plan requires

An indexing plan that is measured once and never again is an indexing plan that decays. Phase 1
ships a benchmark beside the `Store` package, and these are its **regression floors** — a run
slower than the right-hand column fails:

| Query | @400 measured | Fails above |
|---|---|---|
| `COLD-OPEN` | 842 µs | **5 ms** |
| `Q-LIB-PAGE` | 47.3 µs | 500 µs |
| `Q-LIB-COUNTS` | 439 µs | 3 ms |
| `Q-SEARCH` + `Q-SEARCH-RATIO` (per keystroke) | 313 µs | 2 ms |
| `Q-ABSENT-EXISTS` | 14.1 µs | 200 µs |
| `Q-LAST-SYNC` | 15.9 µs | 200 µs |

The floors are set roughly an order of magnitude above the measurement, so the benchmark catches a
**shape** change — a dropped index after a table rebuild (§3.1: 82× at 10 000 rows), a `LIKE`
that stopped being bounded by its `LIMIT`, a view that stopped flattening — and not the noise of
whichever machine CI ran on that morning.

Two assertions belong with it and are not timings:

1. **`EXPLAIN QUERY PLAN` for `Q-LIB-PAGE` contains `USING INDEX item_shelf_order` and does
   *not* contain `USE TEMP B-TREE FOR ORDER BY`.** That is the real regression test; the timing is
   the symptom.
2. **`EXPLAIN QUERY PLAN` for `Q-ABSENT-EXISTS` contains `COVERING INDEX item_absent`.**

Fixture: [`schema/fixtures/bench-generate.sql`](./schema/fixtures/bench-generate.sql) at N = 400,
with N = 10 000 as a scaling check rather than a gate.
