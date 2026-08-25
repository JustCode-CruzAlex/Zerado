---
title: Zerado — the state machine
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-05
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: flowchart
ticket: "#2"
---

# The state machine

The four game states, their legal transitions, and which are user actions versus derived facts.

The states themselves — their colours, glyphs, ASCII fallbacks and labels — are settled in the
brand manual §4.3 and are **not** restated as a design question here. This document is about
what makes a game move between them.

| State | Glyph | ASCII | Label | Colour role |
|---|---|---|---|---|
| Not started | `○` | `[ ]` | `NOT STARTED` | chrome `--z-state-not-started` |
| In progress | `◐` | `[~]` | `IN PROGRESS` | amber `--z-primary` |
| Zerado | `◉` | `[*]` | `ZERADO` | **cyan `--z-accent` — the earned colour** |
| Abandoned | `⊘` | `[x]` | `ABANDONED` | orchid `--z-state-abandoned` |

---

## 1 · The model in one line

```
effective_status = status_manual ?? derive(progress, provider_capabilities)
```

> **These four states are not games-only.** They generalise to books, films and series, and that was
> verified type by type rather than assumed — including the two findings the check surfaced (a film
> has no meaningful `IN PROGRESS`; an ongoing series needs *caught up*, and it must **not** become a
> fifth state). See [`11-media-model.md`](./11-media-model.md) §2 and §4. Phase 1 ships games only.

A game has **one nullable manual status** and a **derivation** that runs when there is none.
That is the whole model, and everything else in this document follows from it.

```go
// The four states. The zero value is deliberately invalid so an
// uninitialised status can never be mistaken for "not started".
type Status uint8

const (
    StatusUnknown Status = iota // zero value — never persisted, never rendered
    StatusNotStarted
    StatusInProgress
    StatusZerado
    StatusAbandoned
)

type Game struct {
    // StatusManual is the player's explicit choice, or nil if they have
    // never expressed one. A sync NEVER writes this field.
    StatusManual *Status
    // Progress is provider-reported and GENERIC — value + unit + source.
    // For a game the unit is minutes and the value is the playtime.
    // Zero is a real value and is distinct from "this provider does not
    // report progress at all" — see Capabilities.Progress.
    Progress Progress
}
```

## 2 · Derived versus user action

**Exactly one transition in the product is ever automatic.**

| Transition | Kind | Trigger |
|---|---|---|
| `NOT STARTED → IN PROGRESS` | **Derived** | A sync reports `progress > 0` on an item with `status_manual = NULL` |
| every other transition | **User action** | The player chooses it in `Z-06` |

The derivation:

```go
func derive(p Progress, c Capabilities) Status {
    if !c.Progress {
        // Nothing is reporting. A cartridge has no telemetry, and
        // neither does a paper book — the same rule covers both.
        return StatusNotStarted
    }
    if p.Value > 0 {
        return StatusInProgress
    }
    return StatusNotStarted
}
```

### 2.1 · The derivation is provider-capability-dependent, and this matters

> This is also the mechanism that makes the media generalisation free: derivation is keyed on the
> **(provider, type) capability**, never on the type itself. A Steam game derives; a cartridge does
> not; a paper book does not; a Kindle book would. One function, every type.
> [`11-media-model.md`](./11-media-model.md) §3.

For a provider that reports progress (Steam reports minutes played), `NOT STARTED` and
`IN PROGRESS` are *facts* until
the player overrides them. For a provider that does not (`physical`, and any store whose API
does not expose progress), **all four states are manual, always** — the derivation has no input
and returns `NOT STARTED` as the honest default.

This is why physical copies are modelled as a **provider with capabilities**
([`06-data-seams.md`](./06-data-seams.md) §2) rather than as a flag on a Steam-shaped row. A
boolean `is_physical` would have forced the derivation to special-case one value; a capability
set makes the same code correct for every provider that will ever exist, including the ones
whose API does not exist yet.

### 2.2 · Manual wins, permanently

> **A sync never changes a status the player set.**

Mark a game `ZERADO`, then play it for three more hours: the next sync updates
`playtime_minutes`, `last_played_at`, and nothing else. The game stays `ZERADO`.

This is the invariant that makes the product trustworthy. The alternative — a sync that
"corrects" the player — means the one action the product is named after can be silently undone
by a background job. That is not a bug to fix later; it is the product failing at its purpose.

## 3 · Legal transitions

**All twelve ordered pairs are legal.** There is no state the player cannot reach from any
other, and the product never refuses a status change.

```
                    ┌───────────────────────────────────────┐
                    │                                       │
                    ▼                                       │
            ┌───────────────┐   playtime > 0 (DERIVED)  ┌────┴──────────┐
            │ ○ NOT STARTED │ ────────────────────────► │ ◐ IN PROGRESS │
            └───────┬───────┘                           └───┬───────┬───┘
                    │  ▲                                    │       │
                    │  │  clear override (re-derive)        │       │
                    │  │                                    ▼       │
                    │  │                            ┌────────────┐  │
                    │  └────────────────────────────┤ ◉ ZERADO   │  │
                    │                               └─────┬──────┘  │
                    ▼                                     │         ▼
            ┌───────────────┐                             │  ┌──────────────┐
            │ ⊘ ABANDONED   │ ◄───────────────────────────┴──┤              │
            └───────────────┘                                └──────────────┘

     Every arrow except the DERIVED one is a player action. All 12 pairs are legal.
```

Why every pair is legal, stated as the cases that justify each of the surprising ones:

| Transition | The case that makes it legal |
|---|---|
| `NOT STARTED → ZERADO` | You beat it on a console years ago; the library only just learned you own it |
| `NOT STARTED → ABANDONED` | "I am never playing this." Owning something is not a commitment |
| `ZERADO → IN PROGRESS` | New Game+, a replay, a second ending |
| `ZERADO → ABANDONED` | You finished it, started again, and stopped |
| `ABANDONED → ZERADO` | You came back and finished it. This should be the easiest transition in the product |
| `anything → NOT STARTED` | A correction. Somebody else played it on your account, or you mis-keyed |

## 4 · `ZERADO` is never automatic — and this is a product decision, not a technical one

Zerado could infer completion: 100% achievements, a known credits-roll playtime, a store's own
"completed" flag. **It will not.**

The moment a game becomes *zerado* is the moment the product exists to create. The brand manual
writes its copy — *"Zerado. 41 hours. Sixth this year."* — and a machine cannot hand that to
someone. A player who is told they finished something has not finished it; they have been
notified.

What Phase 2+ **may** do is *suggest*: a dismissible line in the detail view reading
`100% achievements. Mark this zerado? (s)`. A suggestion the player accepts is still the
player's action. An automatic transition is not, and it is permanently out of scope.

Recorded as a settled decision so it is not reopened as an optimisation.

## 5 · Clearing an override

`status_manual` is nullable, so "no opinion" is representable and distinct from every opinion.
`Z-06 Set status` therefore offers **five** items, not four:

```
   ○  NOT STARTED
   ◐  IN PROGRESS
   ◉  ZERADO
   ⊘  ABANDONED
   ──────────────────────────
   ×  Clear override
      Steam says IN PROGRESS
```

The fifth item appears **only when `status_manual IS NOT NULL`**, and its second line names what
the game will become — **in at most 27 cells**, because `Z-06`'s overlay is a fixed `34 × 11` with a
28-cell content width. The first draft of this line read `Back to what Steam says: IN PROGRESS`,
which is 36 cells and would have been truncated at exactly the point where it stops being
informative. The requirement is the naming, not the wording — because "clear override" without that line is a button whose effect the
player cannot predict.

Note the consequence, and be honest about it in the copy: clearing an override on a game with
playtime > 0 makes it `IN PROGRESS` immediately. Selecting `NOT STARTED` explicitly is a
*different* thing — it stores `not_started` as a manual value and it sticks. Both are useful and
they are not the same action, which is exactly why the model needs a nullable field and a
separate clear.

## 6 · What is persisted

| Column | Type | Meaning |
|---|---|---|
| `status_manual` | `TEXT NULL` | One of `not_started`, `in_progress`, `zerado`, `abandoned`, or `NULL` |
| `status_changed_at` | `TEXT NULL` | RFC 3339. Set when `status_manual` changes; `NULL` when it has never been set |
| `playtime_minutes` | `INTEGER NOT NULL DEFAULT 0` | Provider-reported |
| `last_played_at` | `TEXT NULL` | Provider-reported; `NULL` means *not reported*, not *never played* |

`effective_status` is **not stored.** It is derived on read.

The reason is the one that decides most caching questions: a stored derived value has two ways
to be wrong (stale, or written by the wrong path) and a derived one has none. The library view
computes it for every visible row — at most a few dozen rows per frame, on integers already in
memory. There is no performance argument for storing it, and there is a correctness argument
against.

If sorting or filtering by state ever needs it in SQL, the right answer is a **generated
column** or an index on the expression — still derived, still single-sourced, never a second
writer.

## 7 · The counts on the summary row

The pinned summary — `247 games · ○ 198 · ◐ 12 · ◉ 6 · ⊘ 31` — is a `GROUP BY` over the derived
status across the **currently filtered** set, and it says which set it is describing.

Two rules:

1. **The counts always sum to the number shown.** If the summary says 247 games, the four
   state counts add to 247. A row that cannot be classified does not exist — `StatusUnknown` is
   never persisted and never rendered.
2. **When a filter is active, the summary describes the filtered set and says so.** Showing
   whole-library counts above a filtered list is the most common way a list view lies.

```
 247 games   ○ 198  ◐ 12  ◉ 6  ⊘ 31
  31 games   ○ 24   ◐ 3   ◉ 1  ⊘ 3      filter: source=physical
```

## 8 · Where the states go in Phase 4

The states are the **only** part of the library that crosses the Phase 4 sync boundary
([`06-data-seams.md`](./06-data-seams.md) §6). Not the library, not the credentials, not the
playtime — those are re-derivable from the player's own sources on any device.

Which means the schema has to carry, from **Phase 1**, enough to merge two devices' state
changes without a server having been involved when they were made:

- `status_manual` — the value;
- `status_changed_at` — the timestamp that decides the winner on conflict (last-write-wins is
  adequate for a single human on two devices, and this is a deliberate simplicity choice, not an
  oversight);
- a **stable cross-device game identity** — see [`06-data-seams.md`](./06-data-seams.md) §5.
  This is the expensive one, and it is why the Phase 4 boundary is decided now rather than in
  Phase 4.
