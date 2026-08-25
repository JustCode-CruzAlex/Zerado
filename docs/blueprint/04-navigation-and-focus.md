---
title: Zerado — navigation and focus model
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-04
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: implementation-plan
ticket: "#2"
---

# Navigation and focus

How a player moves between screens and between panes, what the global keys are, how focus is
shown, and how `Esc` behaves everywhere.

---

## 1 · The model: a route stack with one overlay slot

```
        ┌──────────────────────────────┐
        │  overlay slot   (0 or 1)     │   ← Z-06, Z-14, Z-17
        ├──────────────────────────────┤
        │  route stack    (1..n)       │   ← top of stack renders
        │    [3] Z-05 Game detail      │
        │    [2] Z-10 Help             │
        │    [1] Z-04 Library  (root)  │
        └──────────────────────────────┘
```

**Not tabs.** Zerado has one home — the library — and everything else is either reached *from*
something or dismissed *back to* it. A tab bar implies a set of peers, which the screen set is
not: `Z-09 Settings` is not a sibling of `Z-04 Library`, it is a place you go and come back from.
A tab bar would also spend a permanent row of a 24-row terminal on navigation the player uses
a few times a session.

Three rules:

1. **`Z-04 Library` is the root and can never be popped.** The stack bottom is fixed.
2. **At most one overlay.** An overlay cannot open another overlay. If a flow needs two steps,
   it is a route, not an overlay.
3. **Routes do not stack indefinitely.** Pushing a route that is already on the stack **unwinds
   to it** rather than duplicating it, so `?` from help does not build a tower of help screens.

### 1.1 · Routes vs overlays — which to use

| Use a **route** when | Use an **overlay** when |
|---|---|
| The surface needs the whole screen | The decision is small and the context matters |
| The player will spend more than a few seconds | It resolves in one or two keystrokes |
| It has its own scroll region | It has no scroll region |
| It can be arrived at from more than one place | It is always about the thing under it |

At Tiny an overlay becomes a route ([`02-composition.md`](./02-composition.md) §2.4). Behaviour is
identical; only the composition changes.

## 2 · The route graph

```
                       ┌─── Z-10 Help ◄─── (?) from anywhere
                       │
  Z-01 First run ──┬── Z-02 Connect a store ──► Z-03 Sync ──┐
        │          │                                        │
        │          └── Z-08 Add a game by hand ─────────────┤
        │                                                   ▼
        └──────────────────────────────────────────►  Z-04 LIBRARY  (root)
                                                        │  ▲   │
                                     (Enter) ───────────┘  │   └──────── (/) Z-07 Filter [mode]
                                        ▼                  │
                                   Z-05 Game detail        │            (s) Z-06 Set status [overlay]
                                        │                  │             from Z-04 or Z-05
                                        └──── (Esc) ───────┘
                                                           │
                                     (,) Z-09 Settings ────┘
                                     (a) Z-08 Add by hand ─┘
                                     (r) Z-03 Sync ────────┘
```

`Z-11 Fatal error` is not on the graph. It **replaces** everything, from anywhere, and there is
no way back into the program from it.

### 2.1 · First run is a route, not a mode

`Z-01` is pushed as the root **instead of** `Z-04` when the library is empty *and* no provider
has ever been connected. It is not a wizard that owns the session — every door from it lands on a
normal screen, and the player can reach the library from it without doing anything.

The condition is checked once at start-up. It is not re-checked: a player who deletes every game
does not get thrown back into first run, because that would be the program telling them they had
undone their own setup.

## 3 · Global keys

These work on every screen unless a **text input has focus**, in which case every printable key
is literal text and only the modified keys survive.

| Key | Does | Notes |
|---|---|---|
| `?` | Help (`Z-10`) | Push; unwinds if already on the stack |
| `q` | Quit | Immediate. Nothing to confirm — see §3.2 |
| `Ctrl-C` | Quit | Always, including inside a text input |
| `Esc` | Back one level | The full table is §5 |
| `Tab` / `Shift-Tab` | Next / previous region | Only where the screen has ≥2 regions |
| `↑` `↓` / `k` `j` | Move within a region | |
| `←` `→` / `h` `l` | Move between columns, or collapse/expand | Screen-dependent |
| `g` / `G` | First / last item | |
| `Ctrl-D` / `Ctrl-U` | Half page down / up | |
| `Enter` | Activate the focused thing | |
| `r` | Re-sync | From `Z-04` and `Z-05` |
| `s` | Set status (`Z-06`) | From `Z-04` and `Z-05` |
| `a` | Add a game by hand (`Z-08`) | From `Z-04` |
| `/` | Filter and search (`Z-07`) | From `Z-04` |
| `,` | Settings (`Z-09`) | From anywhere |

### 3.1 · Reserved and deliberately unbound from Phase 1

`:` and `Ctrl-K` are **reserved for the command palette** (`Z-17`, Phase 2) and are bound to
nothing in Phase 1. Pressing either does nothing and shows no error.

Reserving them now costs nothing and buys the one thing that is expensive later: **not having to
retrain the player.** A key that means nothing in Phase 1 and something in Phase 2 is a feature
arriving; a key that means one thing in Phase 1 and another in Phase 2 is a betrayal of muscle
memory.

Also reserved, unbound: `1`–`9` (quick filters), `n` / `p` (next / previous *game* from within
the detail view — Phase 2, when there is enough in the detail view to make paging worthwhile).

### 3.2 · Why `q` does not confirm

Every mutation Zerado makes is committed when it is made — a status change writes to SQLite
before the overlay closes. There is no unsaved state, so there is nothing for a confirmation to
protect. A confirmation dialog with nothing behind it does active harm: it teaches the player to
dismiss confirmations without reading them, which is exactly the reflex you need them *not* to
have when `Z-06` eventually asks about something destructive.

The one exception: `q` pressed **while a sync is running** cancels the sync and *then* quits,
and says so in one line on the way out. The sync's `context` is cancelled, in-flight I/O is
aborted, and whatever was already written stays written — a half-synced library is a valid
library, and `Z-03`'s `PARTIAL` state exists to say so.

## 4 · Focus

### 4.1 · The invariants

1. **Exactly one region has focus** on any screen with at least one focusable region.
2. **Exactly one item within that region has focus**, unless the region is empty — in which case
   focus is on the region's empty-state action, and if there is none, the region is skipped by
   `Tab` entirely.
3. **Focus is never nowhere.** If the focused item is removed (a filter hides it, a sync drops
   it), focus moves to the nearest surviving item — the one below, or the one above if there is
   nothing below.
4. **Focus is restored on pop.** The route stack stores each frame's focus state, so returning
   from `Z-05` puts the player back on the row they left, **not** at the top of the list. This is
   the single most-felt navigation detail in a 400-row library.
5. **A rebuild never moves focus.** When a sync replaces the row set underneath the player, the
   cursor and the scroll offset are preserved by *game identity*, not by index — TUI Design
   Manual **R-10(b)**. If the focused game survives the rebuild, focus stays on it.

### 4.2 · How focus is shown — never by colour alone

TUI Design Manual **R-6** and the brand's co-render rule apply to focus exactly as they apply to
state. Focus is carried by **three** channels, and any two of them are enough:

| Channel | Focused | Not focused |
|---|---|---|
| **Position** | A `▌` (U+258C) marker in the left gutter | Two spaces |
| **Weight** | Bold | Normal |
| **Colour** | `--z-primary` amber on the marker | none |

Under `NO_COLOR` the marker and the bold weight both survive, which is why they are the primary
two channels and colour is the third. ASCII fallback for the marker is `>`.

For a **region** (a pane, at ExtraWide): the focused pane's border is `--z-border-strong` and its
title is amber; the unfocused pane's border is `--z-border` and its title is
`--z-text-secondary`. Under `NO_COLOR` the focused pane's border uses a **heavier box-drawing
weight** (`┏━┓` vs `┌─┐`) so the distinction survives with no colour at all.

### 4.3 · The focus ring is never removed

Brand manual §4.2, restated because it is the one accessibility rule most often lost in a
redesign: **the focus indicator is present on every interactive element, in every state, always.**
There is no "cleaner look" that justifies removing it. The launch audience is keyboard-first;
removing focus indication for aesthetics would be a self-inflicted wound on exactly the people
the product is for.

## 5 · `Esc`, everywhere

One table, exhaustive, because "what does Escape do here" is the question that makes an interface
feel either solid or improvised.

| Context | `Esc` does | Then |
|---|---|---|
| An **overlay** is open | Dismiss it, discarding any uncommitted choice | Focus returns exactly where it was |
| A **text input** has focus | Blur the input, **keeping** what was typed | Focus returns to the region |
| **`Z-07` filter mode**, editor focused | Exit the filter *editor*; **the filter stays applied** | The summary row shows the active filter; the footer hints `Esc again to clear` |
| **`Z-07` filter mode**, editor blurred, filter applied | **Clear** the filter | The full library returns |
| A **confirmation** is open | **Cancel** — always the safe branch, never the destructive one | |
| A **non-root route** is on top | Pop it | Focus restored from the frame below |
| **`Z-04`**, no filter, no overlay (the root, at rest) | Nothing | The footer hints `q quit · ? help` |
| **`Z-03`** while a sync is running | Cancel the sync | Lands on `Z-03`'s `CANCELLED` state, then `Esc` again to leave |
| **`Z-11` Fatal error** | Nothing | There is nothing to go back to; only `q` and `Ctrl-C` work |

The two-step Escape in filter mode is the only place `Esc` does something different on a second
press, and it is deliberate: exiting a filter *editor* and clearing a filter are different
intentions, and a player who has just typed a search almost never means to throw it away. The
footer says so at the moment it matters, so it is discoverable rather than clever.

## 6 · The footer

One reserved row, on every framed screen (the Spacing Canon reserves it). It carries the keys
that work **on this screen right now**, in a fixed order:

```
 ↑↓ move   ⏎ open   s status   / filter   a add   r sync   ? help   q quit
```

Rules:
- Keys that do nothing here are **not listed**. A footer that lies is worse than no footer.
- It shortens by dropping from the right at narrower tiers; `? help` and `q quit` are the last
  two to go.
- It is not decoration and is never repurposed for status. Status has its own row.

## 7 · What this model costs, honestly

- **No direct jumps.** There is no "go to settings from the middle of a form" — you leave the
  form first. That is a real cost and the command palette (`Z-17`, Phase 2) is the intended fix.
- **Deep-linking is not free.** `zerado game 42` opening straight to a detail view means
  synthesising a stack beneath it. The stack must therefore be constructible from a route
  descriptor rather than only by pushing — a small constraint on the implementation, recorded
  here so it is not discovered late.
- **Two hosts for one view.** `Z-05` being both a route and a pane is more work than picking one.
  It buys the correct composition at both 80 and 120 columns, which is worth it — but it means
  the detail view must never assume its own width or its own key ownership.
