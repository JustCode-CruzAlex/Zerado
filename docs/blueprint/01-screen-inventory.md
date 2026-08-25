---
title: Zerado — screen inventory, all four phases
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-01
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: implementation-plan
ticket: "#2"
---

# Screen inventory — all four phases

Every screen the product will ever have, named, numbered, and marked with the phase that earns
it. Later-phase screens are present **so the navigation model is designed for the finished
product rather than retrofitted to it** — they are inventoried and composed, not specified.

---

## 1 · How screens are identified

IDs are **flat and stable**: `Z-NN`, allocated once and never reused. The phase is a *column*,
not part of the ID, because screens move between phases and an ID that encodes a phase would
have to be renumbered when they do — which breaks every link that ever pointed at it.

Three kinds of surface appear in this inventory, and the distinction is load-bearing for
navigation ([`04-navigation-and-focus.md`](./04-navigation-and-focus.md)):

| Kind | Meaning | On the route stack? |
|---|---|---|
| **Route** | A full-screen surface pushed onto the route stack | Yes |
| **Overlay** | A modal drawn above the current route; dismissed with `Esc` | No — a separate overlay slot |
| **Mode** | A state of an existing route, not a surface of its own | No |

A **mode** still gets its own screen spec when it carries its own key map, its own empty
state and its own copy. `Z-07 Filter and search` is the clear case: it is a mode of the
library, and it is fully specified, because a builder cannot build it from the library's spec.

---

## 2 · The full inventory

### Phase 1 — CLI/TUI MVP · *your library, your statuses, stored locally*

| ID | Screen | Kind | What it owes the player |
|---|---|---|---|
| **Z-01** | First run | Route | Say what this program is and what happens next, in three sentences, and offer the three doors out: connect a store, add a game by hand, or look around empty |
| **Z-02** | Connect a store | Route | Take the player's own credentials for one provider, explain each field where it is entered, and fail with the actual reason |
| **Z-03** | Sync | Route | Show honest progress on an indeterminate wait, then a readout of what changed — and survive a failure halfway through |
| **Z-04** | Library | Route (**root**) | The whole collection, one row per game, state visible on every row, and never lie about how many there are |
| **Z-05** | Game detail | Route ≤119 cols · Pane ≥120 | Everything known about one game, and honestly mark what is not known yet |
| **Z-06** | Set status | Overlay | Change one game's state in two keystrokes, and offer to clear a manual override |
| **Z-07** | Filter and search | Mode of Z-04 | Narrow 400 games to the handful that match, and say plainly when nothing does |
| **Z-08** | Add a game by hand | Route | Enter a disc, a cartridge or anything else the stores do not know about, as a first-class row |
| **Z-09** | Settings | Route | Every dial the product has, grouped, with the current value visible without opening anything. **Includes the Audio section** — the opt-in, the two channels, two volumes, and an honest line when audio is unavailable ([`12-audio.md`](./12-audio.md)) |
| **Z-10** | Help and key map | Route | Every key that does anything, on the screen it does it on |
| **Z-11** | Fatal error | Route (**frameless**) | When the program cannot continue, say what broke, where the file is, and what to try — in the plainest possible renderer |
| **Z-15** | Cover deck | Mode of Z-04 | Show the collection as covers. **Moved into Phase 1 by founder direction** — *"without image is not an option"*. A terminal with no image support is supported, not refused: it shows the text deck and a quiet, dismissible note recommending Ghostty or Kitty |

**Twelve Phase 1 screens.** The ticket proposed nine; §3 records the amendments.

`Z-15` keeps its number rather than being renumbered `Z-12`: IDs are **flat and stable** (§1), and
renumbering a screen because its phase moved is exactly what that rule exists to prevent.

### Phase 2 — Enrichment · *covers, synopsis, prices, moods*

| ID | Screen | Kind | Note |
|---|---|---|---|
| Z-12 | Enrichment sync | Route | Fetching covers and *sinopse*; a second, separate sync from the store sync |
| Z-13 | Mood picker | Route | "What is tonight for?" — the mood, then the games that match |
| Z-14 | Mood tags | Overlay | Assign or correct a game's mood tags by hand |
| ~~Z-15~~ | *(moved to Phase 1 — see above)* | | Cover art was Phase 2 on the assumption that inline images could not be relied on. Founder direction 2026-08-25 reversed it: images are foundational, with a graceful degrade for terminals that lack them |
| Z-16 | Price history | Mode of Z-05 | Current price, all-time low, and when the low was |
| Z-17 | Command palette | Overlay | Earned here, not in Phase 1 — see §4 |

### Phase 3 — Recommendations and budget · *what to buy, and whether to wait*

| ID | Screen | Kind | Note |
|---|---|---|---|
| Z-18 | Tonight | Route | The recommendation: mood plus the time you actually have |
| Z-19 | Budget | Route | What you can spend, and in which currency |
| Z-20 | Watchlist | Route | What is in budget, what is near its all-time low, and the "maybe wait" verdict |

### Phase 4 — Social and mobile · *sync, community, the phone apps*

| ID | Screen | Kind | Note |
|---|---|---|---|
| Z-21 | Account | Route | The **first** time Zerado has an account at all. Before Phase 4 there is none, and that is a ratified promise |
| Z-22 | Devices and sync | Route | What syncs, to where, and what does not — see [`06-data-seams.md`](./06-data-seams.md) §6 |
| Z-23 | Public profile | Route | The profile as other people see it |
| Z-24 | Comments and reviews | Route | Per-game threads |
| Z-25 | Community lists | Route | Top lists from an **unnamed** source. The source is not named anywhere in the product while the ratified decision stands |
| Z-26 | Premium | Route | Disclosure of what the community layer costs to run. Not a purchase surface until there is something to purchase |
| Z-27 | Phone app | *(Flutter)* | Inventoried so navigation is designed for the finished product. Composition is Flutter's, not this document's |

**Twenty-seven screens over four phases. Eleven of them are Phase 1.**

---

## 3 · The amendment to the Phase 1 set, and why

The ticket proposed: *first-run and setup · sync in progress · the library deck · game detail ·
mark or change status · search and filter · settings · help and key map · error and offline
states.* The architect's amendment:

### 3.1 · Added — `Z-08 Add a game by hand`

The ticket says: *"Physical copies must be first-class from day one, not an afterthought bolted
onto a Steam-shaped model."* The published copy says: *"A physical copy isn't a second-class row
in the list."* The proposed screen set contained **no screen through which a physical copy can
enter the library.**

This is not a missing nicety. Without it, the promise is architecture-only — true in the schema
and false in the product. Adding the screen is the smaller half of the fix; the larger half is
in [`06-data-seams.md`](./06-data-seams.md) §2, where `physical` is modelled as a **provider in
its own right** whose acquisition method happens to be a human typing, rather than as a flag on
a Steam-shaped row.

### 3.2 · Added — `Z-11 Fatal error`; "error and offline states" split three ways

"Error and offline states" as written is one line covering three different things that live at
three different altitudes and cannot share a spec:

| Altitude | What it is | Where it is specified |
|---|---|---|
| **In the row** | A field could not be fetched | The component's own missing-value rendering — deliverable B |
| **In the screen** | Something is unavailable but the screen still works — no network, stale data, a partial sync | The **degrade banner**, one row at the top of the body — deliverable B, applied in every screen's state table |
| **Above the screen** | The program cannot continue — the library file is unreadable, the schema is from a newer version | **`Z-11 Fatal error`**, a route |

`Z-11` renders **without the frame**. A fatal error may be a failure of the very subsystem the
frame is built on, so it uses the plainest renderer that can exist: plain text to stdout, no
layout engine, no theme lookup, no colour. A crash screen that itself crashes tells the player
nothing.

### 3.3 · Renamed — `Z-02 Connect a store`, not "connect Steam"

Steam is the only instance in Phase 1. The screen is nevertheless **generated from the provider
descriptor** — the credential fields, their labels, their help text and their validation all
come from the provider, not from the screen. Adding PlayStation, GOG or EA later adds a
provider and **zero screens**.

The public-profile requirement is explained *where the player hits it* — as the failure copy
on an empty-library result, not as a paragraph up front that nobody reads. The brand manual's
voice example 3 is already the exact copy for that moment.

### 3.4 · Removed — a splash screen

Not proposed by the ticket, but present in the prior concept draft and worth closing
explicitly. The published page says: *"It's a text program. It starts instantly."* A splash
screen is that sentence being broken on the first frame the player ever sees. **Rejected
permanently**, not deferred.

### 3.5 · Folded, not added — the sync result

The ticket lists "sync in progress." The moment *after* a sync is one of the two moments the
product exists for, and the brand manual already writes its copy: *"247 games. 6 finished. Last
played: 3 weeks ago."* That is not a separate screen — it is `Z-03`'s **terminal state**, and it
appears in `Z-03`'s state table as `DONE`. Called out here so it cannot be forgotten by being
nobody's screen.

---

## 4 · Screens deliberately NOT in Phase 1

Recorded so they are not re-proposed, and so the reason survives.

| Not built | Why not, and when |
|---|---|
| **Splash screen** | Contradicts "it starts instantly." Rejected permanently — never |
| **Command palette** (`Z-17`) | A palette earns its place when the surface is bigger than the key map. With eleven screens and one home, `?` and a footer hint are better. **Phase 2** — and `:` and `Ctrl-K` are **reserved and unbound from Phase 1** so it can claim them without retraining anyone |
| ~~Cover deck~~ | **No longer deferred.** Founder direction moved it into Phase 1 — *"Images on terminal is a must… without image is not an option"* |
| **Onboarding tour** | The product is eleven screens. A tour is what you build when the product cannot explain itself |
| **A dashboard / home distinct from the library** | The library **is** the home. A dashboard above it would be a screen whose only content is a summary of the screen below it — and that summary is already one pinned row |
| **Account / login** | There is no account before Phase 4. That is a ratified public promise, not a scheduling choice |
| **Quit confirmation** | Every mutation commits when it is made, so there is never unsaved work to protect. A confirmation with nothing to confirm trains people to dismiss confirmations |

---

## 5 · What each Phase 1 screen must not do

The short list of things that would each be a defect, kept here because they are easier to
prevent in a spec than to find in a review.

- **`Z-01`** must not ask for anything. It offers doors; it does not collect. It carries **one**
  quiet line naming that audio exists and where to turn it on — a statement, not a prompt.
- **`Z-02`** must not tell the player their profile is private *before* it has evidence — it
  is the empty result that proves it, and the copy belongs there.
- **`Z-03`** must not report a count it has not finished counting.
- **`Z-04`** must not hide the state of any row to fit more rows. The state is why the screen exists.
- **`Z-05`** must not render an empty field as if it were an empty value. Not-fetched-yet and
  known-to-be-empty are different facts and the player can act on only one of them.
- **`Z-06`** must not make *zerado* harder to reach than any other state, and must not set it
  automatically. See [`05-state-machine.md`](./05-state-machine.md) §4.
- **`Z-07`** must not silently return zero rows. An empty result names the filter that emptied it.
- **`Z-08`** must not require a field Steam happens to have. A cartridge has no app ID.
- **`Z-09`** must not hide a current value behind a submenu, and must say **why** audio is
  unavailable rather than only that it is off — those are different facts.
- **`Z-11`** must be **silent**. It depends on nothing, and that includes the audio subsystem: a
  crash screen must not try to play a sound.
- **`Z-10`** must not list keys that do nothing on the screen the player came from.
- **`Z-11`** must not depend on anything that could be what broke.
