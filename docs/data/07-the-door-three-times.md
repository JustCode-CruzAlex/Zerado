---
title: Zerado — the door, opened three times on paper
discipline: DATA
doc-no: ZRD-DATA-07
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: concept-explainer
ticket: "#5"
---

# The door, opened three times on paper

> **Founder, 2026-08-25:** *"At this point don't even think on books and other media types. What I
> would like is let that door open, because at some point books can be a part of Zerado."*

**Nothing in this document is built.** It exists to answer the gate's one question about this
model — *does it hold for books without a rewrite?* — and to answer it the way a claim about a
migration should be answered: **by running the migration.**

Everything here is the **game** path with one value added to one `CHECK`. If any of it required a
new table, a new column on `item`, or a change to a screen, the affordance would have failed and
this document would say so.

---

## 1 · The rule, stated once

> **A fact goes on `item` if it is what a *source* says about a copy the player owns.
> A fact goes on `item_user` if it is what the *player* says about it.
> Nothing else has a place, and there is no third question to ask.**

The type discriminator does not participate in that rule. It is not a switch; it is a label, and
its entire job is to stop `SELECT` statements from mixing two libraries once there are two.

---

## 2 · Applied to a **game** — the built path

| Fact | Where | Column |
|---|---|---|
| *Steam says you own Hollow Knight on PC* | `item` | `provider_id`, `provider_ref`, `title`, `platform`, `acquisition` |
| *Steam says 63 hours* | `item` | `playtime_minutes = 3780` |
| *A cartridge has no telemetry* | `item` | `playtime_minutes = NULL` |
| *You marked it zerado on 14 July* | `item_user` | `status_manual`, `status_changed_at` |
| *IGDB's cover and sinopse* | `metadata` | cache, with `fetched_at` |
| *This is a game* | `item` | `item_type = 'game'` |

State: `zerado` because the player said so; otherwise derived from playtime.

---

## 3 · Applied to a **book** — the same tables, nothing new

| Fact | Where | Column | New? |
|---|---|---|---|
| *You own Dune, paperback* | `item` | `provider_id='physical'`, `title`, `platform='Paperback'`, `acquisition='physical'` | **no** |
| *A paperback has no page counter* | `item` | `playtime_minutes = NULL` | **no** — the cartridge case, exactly |
| *You finished it* | `item_user` | `status_manual='zerado'`, `status_changed_at` | **no** |
| *Its cover and blurb* | `metadata` | `cover_ref`, `sinopse`, `fetched_at` | **no** |
| *This is a book* | `item` | `item_type='book'` | one value in one `CHECK` |

**`platform` is the only column that reads oddly**, and it reads oddly rather than wrongly: for a
book it holds *the edition or format* — `Paperback`, `Kindle`, `Audiobook`. That is exactly what it
holds for a game (`PC`, `SNES`, `Nintendo Switch`): **the form the copy takes**. It is also half of
`item_uid`, and *Dune* the paperback and *Dune* the audiobook **should** be two items, for the same
reason a Steam copy and a cartridge are.

Whether it should then be *called* `platform` is a naming question for the day books are real. A
rename of one column with no foreign key pointing at it is the cheapest migration in this schema —
which is the whole difference between this and renaming the table every foreign key points at.

**The four states carry.** [`11-media-model.md`](../blueprint/11-media-model.md)'s appendix already
says so: *"A book you finished **is** zerado; the word was never about games."* And *in progress*
is manual for a book for the same reason it is manual for a cartridge — no automatic signal — which
is a case the schema **already handles for games it ships with**
([`01-erd.md`](./01-erd.md) §8.1). That is the ticket's deliverable 4, and it needed no media type
to be true.

---

## 4 · Applied to a **film** — and the one that would not fit

A film is a book-shaped case: one item, no automatic progress signal, four states that mean what
they mean. Nothing new.

**A series is not.** A series is a container of episodes, and *in progress* for it is
*"season 2, episode 4"* — a position, not a boolean. This model has no place for it, and
manufacturing one now is exactly the speculative generality D5 pruned.
[`11-media-model.md`](../blueprint/11-media-model.md)'s appendix already flags this: *"Films and
series are **not** the same kind of thing, and treating them as one type would be the first
mistake."*

**Named as the honest limit, per the gate's instruction to say what would have to change:** a
series would need a child table (`item_progress(item_id, season, episode, watched_at)`) and a
different derivation. It would **add** a table; it would not rewrite `item`. That is a good outcome
for a type nobody has asked for, and it is not a plan.

---

## 5 · The migration, written and run

**Not a description of a migration — the SQL, executed against a seeded database, with the
result.**

```sql
-- 000N_open_the_door.sql   (DOES NOT SHIP. Written to prove the claim.)
PRAGMA foreign_keys = OFF;
BEGIN IMMEDIATE;
  DROP VIEW    item_view;                      -- dependents first (03-migrations §6, trap 3)
  DROP TRIGGER item_user_exists;
  CREATE TABLE item_new ( … CHECK (item_type IN ('game','book')) … );  -- the ONLY change
  INSERT INTO item_new SELECT * FROM item;     -- every column, unchanged
  DROP TABLE item;
  ALTER TABLE item_new RENAME TO item;
  -- recreate the four indexes, the trigger and the view, unchanged
  PRAGMA foreign_key_check;
COMMIT;
PRAGMA foreign_keys = ON;
```

**What was measured**, through `modernc.org/sqlite` v1.57.0:

| | |
|---|---|
| Time, 400-row library | **9.9 ms** |
| `PRAGMA integrity_check` | `ok` |
| `PRAGMA foreign_key_check` | clean |
| Every `item_user` row after | **byte-identical** to before, asserted field by field |
| Schema objects after | **14** — the same 14 |
| Inserting `item_type='book'` before | `CHECK constraint failed: item_type = 'game'` |
| Inserting `item_type='book'` after | accepted |
| A book's state, through the unchanged view | `zerado`, `playtime_untracked = 1` |
| The games library after | unchanged — 4 / 1 / 1 / 1 across the four states |

**Rows changed on the games side: zero. Screens changed: zero. Columns added: zero.**

---

## 6 · The correction the running turned up

[`11-media-model.md`](../blueprint/11-media-model.md) §2 and ADR-0001 D5 both describe the cost as
*"a migration that **adds a row to a check constraint**"*.

**SQLite has no such operation.** There is no `ALTER TABLE … DROP CONSTRAINT` and no way to widen a
`CHECK`. Opening the door is the **twelve-step table rebuild** above.

**The decision is unaffected and is not reopened.** The rebuild is 9.9 ms on a 400-row library, and
it is a rebuild of **one** table with **no** foreign key pointing at it — which is precisely the
expense D5 bought away by not calling the table `games`. Had it been called `games`, the same
migration would have had to rename the table *and* rewrite every reference to it.

What is affected is the **sentence**, and a bundle whose value is that its facts were checked
should not carry a mechanic that is wrong on the engine it ships on. Filed as
[`08-findings.md`](./08-findings.md) **F-1**, for the ADR's owner. **Not edited here**: ADR-0001 is
founder-ratified at blob `6aa1fe8`, and a specialist does not quietly reword ratified canon.

### 6.1 · A lookup table would avoid the rebuild, and is still the wrong trade

`item_type TEXT NOT NULL REFERENCES item_type(code)`, seeded with one row, would make opening the
door a single `INSERT` and no rebuild at all.

**Rejected**, on three counts: it contradicts D5's explicit *"`CHECK`ed to `'game'` — that is the
entire affordance"*; it adds a table and a foreign key **today** to save **9.9 ms** on a migration
that may never run; and it moves the vocabulary out of `sqlite_schema`, where a reader sees it, and
into data, where they have to query for it. Recorded so it is not re-proposed as an improvement.

---

## 7 · What this document is not

It is not a book design, a phase, a commitment, or a plan. It contains no book table, no book
column, no book screen and no book code.

**It is the receipt for one claim:** that the affordance ADR-0001 D5 bought is real, costs one
value in one `CHECK` plus a 9.9 ms rebuild, and takes the player's own data through it untouched.

The gate asked whether the model holds for books without a rewrite. **It does, and the migration
that proves it has been run.**
