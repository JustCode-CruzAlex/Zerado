---
title: Zerado — schema measurements
discipline: DATA
doc-no: ZRD-DATA-EV
rev: A
date: 2026-08-25
archetype: implementation-plan
ticket: "#5"
---

# Measurements

Every performance number in [`../04-indexing-plan.md`](../04-indexing-plan.md) comes from here.
Nothing in that document is an estimate.

**Why this file exists.** *"A cold open on a 400-title library must be instant"* is the bar the
ticket sets, and an indexing plan defended by intuition is an indexing plan nobody can check. The
alternative to measuring is asserting, and this bundle's whole value is that its facts were
checked.

---

## 1 · The rig, stated so it can be disbelieved

| | |
|---|---|
| **Engine** | `modernc.org/sqlite` **v1.57.0** → SQLite **3.53.3** — the pure-Go driver ADR-0001 D2 selects, not the system CLI. Verified at source: `SELECT sqlite_version()` through the driver |
| **Host** | Apple Silicon, macOS (`darwin/arm64`), Go 1.27.0 |
| **Schema** | [`../schema/migrations/0001_phase1_core.sql`](../schema/migrations/0001_phase1_core.sql), applied unmodified |
| **Data** | [`../schema/fixtures/bench-generate.sql`](../schema/fixtures/bench-generate.sql) at N = 400 and N = 10 000 |
| **Statistic** | **minimum of 7 samples × 200 iterations**, after 50 warm-up calls. A minimum, not a mean: for a microbenchmark the mean measures the scheduler and the minimum measures the query |
| **Method** | Full row consumption (`rows.Next()` to exhaustion), one connection, `foreign_keys=ON`, `busy_timeout=5000` |

**What is NOT measured here:** process start-up, terminal initialisation, and Bubble Tea's first
frame. Those belong to Phase 1 and are named as such below.

**Honest limits.** These are single-host numbers on a warm page cache. A cold page cache on a
spinning disk adds one seek per page; the whole 400-row file is **204 KiB / 50 pages**, so that
ceiling is bounded and small. The numbers are used to compare *shapes* — indexed against
unindexed, split against unsplit — and comparisons are what they support.

---

## 2 · Query timings

`Q-*` names are defined in [`../04-indexing-plan.md`](../04-indexing-plan.md) §2 and each one is
traced to the screen spec that asks for it.

| Query | **400 rows** | 10 000 rows | Grows with library size? |
|---|---|---|---|
| `Q-LIB-PAGE` — the root screen's 12 rows | **47.3 µs** | 46.8 µs | **No.** Index-ordered, `LIMIT`-bounded |
| `Q-LIB-COUNTS` — the pinned summary | **439 µs** | 17.00 ms | Yes, linear. §4 |
| `Q-SEARCH` — 12 matching rows | **112.8 µs** | 72.7 µs | No (`LIMIT`-bounded) |
| `Q-SEARCH-RATIO` — `23 of 247` | **200.4 µs** | 8.42 ms | Yes, linear |
| `Q-ABSENT-EXISTS` — does the 5th chip render? | **14.1 µs** | 12.5 µs | **No.** Covering index |
| `Q-ABSENT-COUNT` — the summary's `<n> absent` | **12.6 µs** | 14.1 µs | **No.** Covering index |
| `Q-ABSENT-LIST` | **40.0 µs** | 113 µs | Bounded by the absent count |
| `Q-DETAIL` — one game | **45.5 µs** | 45.1 µs | **No.** Two rowid seeks |
| `Q-LAST-SYNC` | **15.9 µs** | 15.9 µs | **No** |
| `Q-TOTAL` — the ratio's denominator | **149.5 µs** | 7.26 ms | Yes, linear |
| **`COLD-OPEN`** — new connection + migration check + `Q-LIB-PAGE` + `Q-LIB-COUNTS` | **842 µs** | 15.20 ms | Dominated by `Q-LIB-COUNTS` |

**The bar is met with three orders of magnitude to spare.** A 400-title cold open costs the
database **0.84 ms**. Process start dominates it by roughly an order of magnitude and the terminal
by more.

**The per-keystroke budget in filter mode** — `Q-SEARCH` + `Q-SEARCH-RATIO`, the two queries
Z-07 re-runs as the player types — is **313 µs at 400 titles**, against a 16.7 ms frame at 60 fps.
It is **1.9 % of one frame**.

---

## 3 · What `item_shelf_order` buys

The same query with the index dropped, everything else identical:

| | with `item_shelf_order` | without | |
|---|---|---|---|
| `Q-LIB-PAGE` @ 400 | **47.3 µs** | 189 µs | **4.0×** |
| `Q-LIB-PAGE` @ 10 000 | **46.8 µs** | 3 829.6 µs | **81.8×** |
| `COLD-OPEN` @ 10 000 | 15.20 ms | 10.91 ms | *(the unindexed build is faster here only because it has one fewer index to open; the query itself is 82× slower)* |

The plan explains it:

```
with     |--SCAN i USING INDEX item_shelf_order      →  stops after 15 rows
         `--SEARCH u USING INTEGER PRIMARY KEY (rowid=?)

without  |--SCAN i                                   →  reads every row,
         |--SEARCH u USING INTEGER PRIMARY KEY (rowid=?)   then sorts all of them
         `--USE TEMP B-TREE FOR ORDER BY
```

**The number that matters is not 4× — it is the flatness.** 47.3 µs at 400 rows and 46.8 µs at
10 000 rows is the same query doing the same work regardless of library size, because `LIMIT 12`
against an ordered index stops after fifteen rows. Without the index the root screen's cost is
proportional to the library, which is the shape that turns *instant* into *it used to be instant*.

---

## 4 · `Q-LIB-COUNTS` is the most expensive Phase 1 query, and it is fine

439 µs at 400 rows, 17.0 ms at 10 000. It is a `GROUP BY` over the derived status across the whole
shown set, so it is linear by nature — there is nothing to index, because every row is in the
answer.

- **At the ticket's bar (400 titles) it is 439 µs** and recomputed only when the row set changes —
  a sync, a status write, entering or leaving filter mode. Not per keystroke: in filter mode the
  summary is **replaced** by the ratio (`Z-07` D-07-1), so the per-keystroke cost is §2's 313 µs.
- **It becomes a visible hitch somewhere past ~8 000 titles.** Recorded with its fix — maintain the
  four counts incrementally in the view model, invalidating on write — and **not built**, because
  building it now is a cache with two ways to be wrong for a library nobody has.

---

## 5 · What the `item` / `item_user` split costs

The same data, same indexes, same queries, against a variant schema with the user columns folded
back onto `item` ([`../05-truth-and-cache.md`](../05-truth-and-cache.md) §3 explains why they are
not).

| Query | split | unsplit | **cost of the split** |
|---|---|---|---|
| `Q-LIB-PAGE` @ 400 | 47.3 µs | 37.5 µs | **+9.8 µs** |
| `Q-LIB-COUNTS` @ 400 | 439 µs | 321 µs | **+118 µs** |
| `Q-SEARCH` @ 400 | 112.8 µs | 103.6 µs | +9.2 µs |
| `Q-DETAIL` @ 400 | 45.5 µs | 35.4 µs | +10.1 µs |
| `Q-ABSENT-LIST` @ 400 | 40.0 µs | 25.5 µs | +14.5 µs |
| `Q-LIB-PAGE` @ 10 000 | 46.8 µs | 37.2 µs | +9.6 µs |
| `Q-LIB-COUNTS` @ 10 000 | 17.00 ms | 12.34 ms | +4.66 ms |

**The whole cold open costs 128 µs more than it would without the split** — 0.13 ms, on a budget
where 1 000 ms would still read as instant. The plan shows why it is so cheap: the join is
`SEARCH u USING INTEGER PRIMARY KEY (rowid=?)`, a rowid seek into a page already in cache.

That is the price of making *"a metadata refresh can never destroy user data"* a property of the
SQL rather than a property of the caller's discipline.

---

## 6 · Storage, per object, at 400 rows

`SELECT name, SUM(pgsize), COUNT(*) FROM dbstat GROUP BY name`:

| Object | Bytes | Pages | |
|---|---|---|---|
| `item` | 94 208 | 23 | |
| **`item_uid_lookup`** | **24 576** | **6** | the largest index in the file, and it has no Phase 1 query — §04-indexing-plan §5 |
| `item_shelf_order` | 16 384 | 4 | serves the hottest query in the product |
| `item_user` | 12 288 | 3 | the entire truth table |
| `item_provider_identity` | 12 288 | 3 | |
| `sync_run_recent` · `sync_run` · `setting` · `provider_connection` · `schema_migration` | 4 096 each | 1 each | |
| **`item_absent`** | **4 096** | **1** | one page, usually empty — which is exactly what makes `Q-ABSENT-EXISTS` 14 µs |
| **Whole file** | **208 896** | **51** | 204 KiB. The entire library fits in the OS page cache many times over |

At 10 000 rows the file is **3.3 MB**. Both are small enough that the player actually backs them
up, which was the reason cover blobs were kept out (ADR-0001 D2·4).

---

## 7 · The pre-migration backup

`VACUUM INTO` — not a file copy. A file copy of an open WAL-mode database copies a torn page set
and produces a backup that fails to open exactly when it is needed.

| | Time | Output |
|---|---|---|
| 400 rows | **< 10 ms** | 208 896 bytes |
| 10 000 rows | **~10 ms** | 3 309 568 bytes — smaller than the live file, because `VACUUM INTO` defragments |

Fast enough that [`../03-migrations.md`](../03-migrations.md) can take a backup before **every**
destructive migration without the player noticing.

---

## 8 · The one-file promise, verified through the driver

The published promise is one SQLite file. WAL mode produces `-wal` and `-shm` companions. Measured
through `modernc.org/sqlite`, one owned connection:

```
  journal_mode = wal
  after write (connection OPEN)
      library.db      208896 bytes
      library.db-wal    8272 bytes
      library.db-shm   32768 bytes
  wal_checkpoint(TRUNCATE) -> busy=0 log=0 checkpointed=0
  after checkpoint (connection still OPEN)
      library.db      208896 bytes
      library.db-wal       0 bytes
      library.db-shm   32768 bytes
  after db.Close()
      library.db      208896 bytes          ← exactly one file
```

**The promise holds, and it holds because of a specific discipline, not by luck.** The companions
are removed by the *last connection closing*, so:

1. the process owns **exactly one** connection to the library file;
2. clean shutdown runs `PRAGMA wal_checkpoint(TRUNCATE)` and then closes it.

After a crash the companions survive and the next open recovers from them — which is WAL working
as designed, and is why the promise is *"a clean shutdown leaves one file"* rather than *"there is
never more than one file"*. [`../03-migrations.md`](../03-migrations.md) §7 makes both a startup
check and a test.

*(For contrast: the `sqlite3` CLI left `-shm` behind across invocations in the same run. That is a
CLI artifact and not the product's path — which is precisely why this was measured through the
driver the product actually ships.)*

---

## 9 · Reproducing this

```bash
sqlite3 library-400.db   < docs/data/schema/migrations/0001_phase1_core.sql
sed 's/:N/400/'   docs/data/schema/fixtures/bench-generate.sql | sqlite3 library-400.db
sqlite3 library-400.db "EXPLAIN QUERY PLAN <the query from 04-indexing-plan.md §2>"
```

Timings need the driver rather than the CLI, so the harness that produced §2 belongs with the
Phase 1 store package: it is specified as a **required Phase 1 benchmark** in
[`../04-indexing-plan.md`](../04-indexing-plan.md) §7, with the numbers above as its regression
floor.
