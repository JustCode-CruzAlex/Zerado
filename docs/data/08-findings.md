---
title: Zerado — findings from the physical model
discipline: DATA
doc-no: ZRD-DATA-08
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: implementation-plan
ticket: "#5"
---

# Findings

Everything this ticket found while turning the conceptual model into a runnable one. **Every claim
below was verified at source in this branch's tree**, with the file and line quoted, and every one
of them was re-read rather than remembered.

**Nothing here is fixed in another document's file.** Where canon is wrong or drifting, this
ticket says so and routes it. ADR-0001 is founder-ratified at blob `6aa1fe8`; a specialist does not
quietly reword ratified canon.

| # | Severity | Finding | Owner |
|---|---|---|---|
| **F-1** | **minor** | *"A migration that adds a row to a check constraint"* — SQLite has no such operation | ADR-0001 / `11-media-model.md` |
| **F-2** | **minor** | `status_changed_at`'s meaning is under-specified in a way that has a Phase 4 bug in it | `05-state-machine.md` |
| **F-3** | **founder** | An index exists that this ticket's own rule would remove | founder |
| **F-4** | **minor** | `game_uid` / `item_uid` and `game_mood` / `item_mood` drift across five documents | `fft-design-architect` |
| **F-5** | **minor** | `Capabilities.Progress` vs `Capabilities.Playtime` — the screen specs and the seam spec disagree | `fft-api-designer` (#6) |
| **F-6** | **minor** | Two ratified chart sheets render text past their canvas; `[[free]]` does not wrap | `fft-tui-architect` |
| **F-7** | **trivial** | `Z-07` §3's mockup shows four rows that do not match its own query | `fft-tui-designer` |
| **F-8** | **note** | `Z-07` §19's open `source` facet has no schema cost either way | founder / `fft-tui-architect` |

---

## F-1 · Widening the `item_type` `CHECK` is a table rebuild, not a row

**Where.** ADR-0001 D5 Alternatives, and
[`11-media-model.md`](../blueprint/11-media-model.md) §2 item 1: *"Adding a second value later is a
migration that **adds a row to a check constraint**."*

**The problem.** SQLite has no `ALTER TABLE … DROP CONSTRAINT` and no way to widen a `CHECK`.
Opening the door is the **twelve-step table rebuild** — new table, copy, drop, rename, recreate
every index, trigger and view.

**The decision is right and is not reopened.** The rebuild was written and run
([`07-the-door-three-times.md`](./07-the-door-three-times.md) §5): **9.9 ms** on a 400-row library,
`integrity_check` clean, every `item_user` row byte-identical, zero screens changed. It is a
rebuild of **one** table with **no** foreign key pointing at it — which is exactly the expense D5
bought away by not naming the table `games`.

**Only the sentence is wrong**, and it is wrong about the engine the product ships on. Suggested
wording: *"a migration that widens one `CHECK` on one table — a contained rebuild, measured at
under 10 ms on a 400-title library."*

**Not edited here.** ADR-0001 is founder-ratified; `11-media-model.md` is in the same bundle.

---

## F-2 · A cleared override must keep its timestamp, or Phase 4 resurrects it

**Where.** [`05-state-machine.md`](../blueprint/05-state-machine.md) §6:
*"`status_changed_at` … Set when `status_manual` changes; **`NULL` when it has never been set**."*

**The problem.** The sentence is right and is easy to read as a biconditional — *`NULL` whenever
`status_manual` is `NULL`*. That reading has a Phase 4 bug:

> Device A marks a game `ZERADO` on the 1st; it syncs to device B. Device A **clears the override**
> on the 20th. If the clear drops the timestamp, it carries no time, loses every last-write-wins
> comparison, and **device B's stale `ZERADO` comes back** — the product silently undoing the
> player's most recent decision, which is the failure §2.2 calls *"the product failing at its
> purpose."*

**Closed in this ticket, in the schema:**

```sql
CHECK (status_manual IS NULL OR status_changed_at IS NOT NULL)   -- NOT a biconditional
```

`status_changed_at IS NULL` therefore means exactly one thing: *this player has never had an
opinion about this game.* `seed-minimal.sql` row 2 is the cleared case, present so a test asserts
it before Phase 4 exists to need it.

**Suggested amendment** to §6: *"Set when `status_manual` changes, **including when it is
cleared**. `NULL` only while the player has never expressed an opinion."* One clause.

---

## F-3 · `item_uid_lookup` — an index with no Phase 1 query, kept on canon

**For the founder.** This ticket's brief says *"every index justified by a NAMED query from a
screen spec… an index with no query behind it is removed."* One index has none, and ADR-0001 D4
requires it: *"assigned at insert, **indexed but not unique**… one column and one index now."*

**It is kept**, because reversing a ratified decision is not a specialist's call. What the decision
looks like with numbers on it:

| | |
|---|---|
| Phase 1 queries served | **none** |
| Phase 4 queries served | the merge-candidate lookup |
| Cost, measured | **24 KiB / 6 pages at 400 rows — the largest index in the file**, larger than the one serving the hottest query in the product, plus one B-tree write per insert |
| Incidental use, measured | SQLite elects it as the smallest covering index for a bare `SELECT COUNT(*) FROM item`. Real, but weak — the count the screens need is `Q-TOTAL`, which uses `item_shelf_order` |

**The alternative, if the rule should hold without exception:** move `CREATE INDEX
item_uid_lookup` into the Phase 4 migration and keep the **column** in 0001. That costs nothing
and loses nothing — unlike the column, an index can be built later from data that is already there,
instantly, on 400 rows.

**Recommendation: keep it.** 24 KiB is not worth an exception to canon. Recorded so the choice is
visible rather than assumed. Full detail: [`04-indexing-plan.md`](./04-indexing-plan.md) §5.

---

## F-4 · `game_uid` / `item_uid` and `game_mood` / `item_mood`

**Where**, verified at source in this tree:

| Says `game_uid` | |
|---|---|
| `adr/ADR-0001…md` | lines 216, 226, 241 |
| `blueprint/06-data-seams.md` | lines 429, 433, 440, 444 |
| `blueprint/13-handoffs.md` | lines 68, 84 |
| `design/screens/Z-08-add-a-game-by-hand.md` | lines 66, 405 |

| Says `item_uid` | |
|---|---|
| `blueprint/09-erd.md` | §1 table, §4 in full |

And **`09-erd.md` disagrees with itself** on the mood table: §2 defines **`item_mood`**; §5 refers
to **`game_mood.source = 'user'`** (line 195).

**Resolved here as `item_uid` and `item_mood`**, for one reason: the table is `item`, and a
`game_uid` column on a table called `item` reintroduces the exact word D5 spent a decision
removing. [`09-erd.md`](../blueprint/09-erd.md) is the ERD document and the later statement of the
same decision.

**This is a spelling reconciliation, not a decision reopened.** D4's decision — *a stable
content-derived identity, assigned at insert, indexed and not unique, a hint and never an
authority* — is carried unchanged.

**A second, smaller drift inside the same finding:** the two stated formulas differ.
[`09-erd.md`](../blueprint/09-erd.md) §4 includes `item_type` in the hashed key and argues for it
(*"it is type-scoped: the same title as a game and as a film are two items, and should be"*);
[`06-data-seams.md`](../blueprint/06-data-seams.md) §6.2 omits it. **This ticket follows
`09-erd.md`** — the argument is explicit and correct, and the omission is silent.

---

## F-5 · `Capabilities.Progress` or `Capabilities.Playtime`?

**Where**, verified at source:

| Field name | Documents |
|---|---|
| **`Playtime`** | `blueprint/06-data-seams.md:67` (the declaration) · `blueprint/05-state-machine.md:68` · `blueprint/09-erd.md:59` |
| **`Progress`** | `Z-04:1027` · `Z-05:399` · `Z-06:366` · `Z-08:75, 375` |

**Four screen specs name a field the seam does not declare.** `Progress` is almost certainly a
survivor of the pruned generic-progress model (ADR-0001 D5: *"the generic-progress divergence is
**withdrawn**… `playtime_minutes` is a plain column"*), which would make `Playtime` the surviving
name.

**Not this ticket's to decide** — [`13-handoffs.md`](../blueprint/13-handoffs.md) §1 assigns the
final Go signatures to `fft-api-designer` (#6). It matters to the schema only through the
invariant [`01-erd.md`](./01-erd.md) §4.1 states, which uses whichever name #6 lands on:

> `playtime_minutes IS NULL` ⟺ the writing provider's capability for *"reports minutes played"* is
> `false`.

---

## F-6 · Two ratified chart sheets render text past their canvas

**Found while rendering this ticket's own sheets**, by writing a geometry check for them and then
running it against the existing ones.

| Sheet | Defect |
|---|---|
| `ZRD-ERD-01` | The `THE DOOR-OPEN AFFORDANCE IS TWO DECISIONS…` `[[free]]` block is 296 characters at `font-size 13` from `x=110` — it ends at **x ≈ 2043 on an 1820-wide canvas**. The last ~110 characters are off the sheet, in **both** themes |
| `ZRD-ERD-02` | A `mood_tag` subline (*"key   stable identifier the recommender reasons over…"*) renders **~47 px wider than its node box** |

**The cause is worth recording, because it will happen again.** `flowforge chart`'s `[[free]]`
block has a `w` field that **is not used for wrapping** — the renderer emits one `<text>` element
per block and the text simply runs off. This ticket's sheets work around it by splitting long notes
into one `[[free]]` per line, which is why `ZRD-ERD-03` and `ZRD-ERD-04` check clean in both
themes.

**Deliberately not fixed here.** Both are ticket #2's artifacts and both are ratified. The fix is
two minutes — split one `[[free]]` into two, widen one node — and it is that ticket's to make.

**Offered:** the geometry check is fifteen lines of Python (parse the SVG, estimate each `<text>`
run's width from its `font-size` and `text-anchor`, compare against the canvas and its containing
`<rect>`). It caught three real defects across four sheets on its first run. If `fft-tui-architect`
wants it as a `flowforge chart render` post-check, it is easy to lift out of this ticket's working
notes.

---

## F-7 · `Z-07` §3's mockup shows four rows that its own query excludes

The editor reads `/ the`, the ratio reads `23 of 247`, and the scroll row reads `ROWS 1–11 of 23` —
so the eleven rows are presented as the filtered result. Four of them do not contain the substring
`the`:

`Bastion` · `Death's Door` · `Dark Souls III` · `Nier: Automata`

Under D-07-5 — *a case-insensitive, accent-folded **substring** match on the title* — none of those
match, and D-07-5 is the decision that makes the empty-result diagnostic honest.

**Trivial, and worth one line of a future revision** in a bundle whose stated value is that a
mockup which does not add up is the defect it is trying to prevent. It is also the mockup a builder
would read while implementing the match.

---

## F-8 · The open `source` facet has no schema cost either way — a note, not a finding

[`Z-07`](../design/screens/Z-07-filter-and-search.md) §19 item 2 asks the founder whether a
`source` filter exists in Phase 1. Recorded here because the data model would ordinarily be part of
that decision, and here it is not:

**Either answer costs the schema nothing.** `provider_id` is already `NOT NULL` on `item`, and at
400 rows a `WHERE provider_id = ?` scan is within `Q-SEARCH`'s measured 113 µs — no index, no
column, no migration. **The decision is purely about the interface** (D-07-8's answer applies: a
sixth chip needs no new key, at the cost of the row arithmetic in §3.4).

Said plainly so nobody waits on the database for an answer it does not have.

---

## What this ticket did NOT find

Recorded because a findings list with nothing absent from it is a list nobody trusts.

- **No contradiction in ADR-0001's nine decisions.** Every one of them is carried. F-1 is about a
  sentence describing a mechanic, not about the decision it describes.
- **No place where the conceptual model made the physical model impossible.** Twelve deltas
  ([`01-erd.md`](./01-erd.md) §6), all of them the kind
  [`13-handoffs.md`](../blueprint/13-handoffs.md) §3 explicitly hands over — types, constraints,
  indexes, a table split.
- **No case where the four states failed to carry**, including through the door migration.
- **No performance problem at the ticket's bar.** 842 µs cold open on 400 titles, three orders of
  magnitude inside *instant*.
