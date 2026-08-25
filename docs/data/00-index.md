---
title: Zerado — the physical data model, index
discipline: DATA
doc-no: ZRD-DATA-00
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: project-brief
ticket: "#5"
---

# The physical data model

Ticket **#5**. The spine decided *that* a persistence seam exists and what it is responsible for
([`../blueprint/09-erd.md`](../blueprint/09-erd.md), ADR-0001 D2). This is its exact shape: the
ERD drawn, the DDL that runs, the migration strategy, the indexing plan, and the line between what
only the player can supply and what a machine can fetch again.

**Ratified as part of the Phase 0 blueprint bundle. Not a separate founder gate.**

---

## Read in this order

| | Document | For |
|---|---|---|
| **1** | [`01-erd.md`](./01-erd.md) | The ERD — entities, attributes, relationships, cardinality, both drawn sheets, and **every difference from the conceptual model** with its reason |
| **2** | [`02-physical-schema.md`](./02-physical-schema.md) | Every column decision, the connection contract, the normalisers, and **the DDL verbatim** |
| **3** | [`03-migrations.md`](./03-migrations.md) | How the program migrates its own file — or refuses, and says why |
| **4** | [`04-indexing-plan.md`](./04-indexing-plan.md) | Eight indexes, seven named queries, one confession, six removals |
| **5** | [`05-truth-and-cache.md`](./05-truth-and-cache.md) | The line in the schema, the discard test, and the offline contract at the storage layer |
| **6** | [`06-phase4-seam.md`](./06-phase4-seam.md) | What changes when a sync server exists. Designed, not built |
| **7** | [`07-the-door-three-times.md`](./07-the-door-three-times.md) | The gate's question — *does it hold for books?* — answered by **running the migration** |
| **8** | [`08-findings.md`](./08-findings.md) | Eight findings, each verified at source, each routed. Nothing fixed in another ticket's file |
| — | [`evidence/measurements.md`](./evidence/measurements.md) | Every performance number in this bundle, and how it was taken |

## The artifacts

| | |
|---|---|
| **The DDL that runs** | [`schema/migrations/0001_phase1_core.sql`](./schema/migrations/0001_phase1_core.sql) — Phase 1. [`0002`](./schema/migrations/0002_phase2_enrichment.sql) and [`0003`](./schema/migrations/0003_phase3_prices.sql) are written and **do not ship with Phase 1** |
| **The head schema** | [`schema/schema.head.sql`](./schema/schema.head.sql) — **generated, never executed**. It exists to be reviewed and to be diffed |
| **Fixtures** | [`seed-minimal.sql`](./schema/fixtures/seed-minimal.sql) (nine rows, every rule) · [`seed-edge-cases.sql`](./schema/fixtures/seed-edge-cases.sql) (the rows that break naive code) · [`bench-generate.sql`](./schema/fixtures/bench-generate.sql) |
| **The drawings** | Sheet 03 [brand-black](../adr/charts/svg/ZRD-ERD-03.svg) · [cyanotype](../adr/charts/svg/ZRD-ERD-03.cyanotype.svg) — the physical model. Sheet 04 [brand-black](../adr/charts/svg/ZRD-ERD-04.svg) · [cyanotype](../adr/charts/svg/ZRD-ERD-04.cyanotype.svg) — the migration decision |

---

## The seven decisions worth arguing about

Everything else follows from the spine. These are the calls this ticket made, each with the
measurement or the citation that decided it.

1. **`item` and `item_user` are two tables.** *"A metadata refresh can never destroy user data"*
   becomes a property of the SQL — the sync's `UPSERT` **cannot name** a truth column — instead of a
   convention the next author has to remember. **Costs +128 µs on a whole cold open.**
   [`05-truth-and-cache.md`](./05-truth-and-cache.md) §3.
2. **`playtime_minutes` is nullable.** `NULL` = *this source cannot know*; `0` = *it knows, and it
   is zero*. The provider seam already carries that distinction as `*int`; a `NOT NULL DEFAULT 0`
   column was flattening it, and `Z-04` renders `—` and `0h` from the row alone.
   [`01-erd.md`](./01-erd.md) §4.1.
3. **`item.provider_id` has no foreign key.** An FK would force a choice between deleting a library
   when the player disconnects a store and forbidding the disconnect. Both are wrong.
   [`01-erd.md`](./01-erd.md) §3.1.
4. **A cleared override keeps its timestamp.** Without it a clear loses every last-write-wins
   comparison and a stale status from another device silently comes back.
   [`02-physical-schema.md`](./02-physical-schema.md) §6.1.
5. **`last_synced_at` is per row.** A `PARTIAL` sync leaves rows at different ages, and `Z-05`
   renders the age of *this* game. One provider-level timestamp would be true for the rows the run
   reached and a lie for the rest. [`01-erd.md`](./01-erd.md) §4.2.
6. **There is no price-history table**, so there is no retention rule to forget. The price seam
   already supplies the all-time low, so one row per (item, shop, currency), UPSERTed in place,
   answers *"how far above the low is today?"* without re-fetching and without growing.
   [`schema/migrations/0003_phase3_prices.sql`](./schema/migrations/0003_phase3_prices.sql).
7. **Six candidate indexes were removed and one confessed.** Three of the queries the ticket names
   — sort by last played, filter by platform, full-text over *sinopse* — are not Phase 1 queries,
   and the screen specs say so. [`04-indexing-plan.md`](./04-indexing-plan.md) §4, §5.

---

## The bar, and where it landed

> *"A cold open on a 400-title library must be instant. The indexing plan is judged against that."*

| | Measured | |
|---|---|---|
| **Cold open, 400 titles** | **842 µs** | connection + migration check + the root screen's page and counts |
| The root screen's page | 47.3 µs | and **flat** — 46.8 µs at 10 000 titles |
| Per keystroke while filtering | 313 µs | **1.9 %** of one frame at 60 fps |
| Does the fifth chip render? | 14.1 µs | a covering index that is usually empty |
| The whole file, 400 titles | **204 KiB / 51 pages** | it fits in the page cache many times over |
| Opening the door to books | 9.9 ms | measured by running the migration, not by describing it |

Taken on `modernc.org/sqlite` v1.57.0 → SQLite 3.53.3 — the driver ADR-0001 D2 selects, not the
system CLI. [`evidence/measurements.md`](./evidence/measurements.md) states the rig, the method and
the limits.

---

## Scope, honoured

**Books were de-scoped by the founder**, and this model contains no book table, no book column, no
book screen and no book code. The affordance is what ADR-0001 D5 specifies and nothing more: the
table is `item`, and it carries an `item_type` `CHECK`ed to `'game'`.
[`07-the-door-three-times.md`](./07-the-door-three-times.md) is a **receipt**, not a design — it
opens the door once, on a copy, to prove the claim, and reports what it cost.

**Out of scope, and not done:** no Go, no ORM or query-builder choice, no repository interface
(that is `fft-api-designer`, #6), and no Phase 4 sync server.
