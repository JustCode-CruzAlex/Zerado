---
title: Zerado — Z-06 Set status
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-06
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-06 · Set status

> The two-keystroke screen the whole product is named after. *"The moment a game becomes
> zerado is the moment the product exists to create."* (`05-state-machine.md` §4.)

**Canon that governs this screen:** `00-design-brief.md` §10 · `01-design-system.md` §3 (chip) ·
§13 (what counts as destructive — **and this is not it**) · `02-colour-budget.md` §10 ·
`03-designer-manual.md` §3 · `02-composition.md` §2.4 (**binding**: fixed `34 × 11`, centred,
dimmed backdrop where available, **a route at Tiny**) · `03-responsive.md` §3 ·
`04-navigation-and-focus.md` §1 (one overlay slot), §5 (`Esc`) ·
**`05-state-machine.md` §5 — the exact five-item composition** · §2.1 · §4 ·
`01-screen-inventory.md` §5 — *"`Z-06` must not make zerado harder to reach than any other
state, and must not set it automatically."*

---

## 1 · Identity

| | |
|---|---|
| **ID** | `Z-06` |
| **Name** | Set status |
| **Phase** | 1 |
| **Kind** | **Overlay** — the single overlay slot, not on the route stack. **A route at Tiny** (`< 40` cols), where `34 × 11` does not fit inside `30 × 21` |
| **Route in** | `s` from `Z-04` (a row focused) or from `Z-05` (either host) |
| **Route out** | `⏎` applies and dismisses · `Esc` cancels and dismisses. **Nothing else** |
| **Shape** | **Overlay** — shape 5 of the five (`02-composition.md` §3) |
| **Offline class** | **WORKS.** *"Writes immediately. This is the product's core action and it has never needed a network"* (`07-offline-contract.md` §2) |
| **Displayed name** | as an overlay: none — the game's title labels it. As a route at Tiny: breadcrumb `Zerado ✦ Set status`, title `SET STATUS` |

**An overlay cannot open another overlay** (`04-navigation-and-focus.md` §1 rule 2). If a flow
ever needs two steps here, it becomes a route.

---

## 2 · Purpose

**One sentence:** change one game's state in two keystrokes, and — when the player has already
overridden it — offer to hand the decision back to the store, naming exactly what the game will
become.

---

## 3 · Mockup — the overlay, fixed 34 × 11

Three variants, one footprint. The box never resizes and never scrolls.

**A · a manual override IS set — five items** (`05-state-machine.md` §5)

```
┌────────────────────────────────┐
│  Return of the Obra Dinn       │
│                                │
│  ▌ ○  NOT STARTED              │
│    ◐  IN PROGRESS              │
│    ◉  ZERADO                   │
│    ⊘  ABANDONED                │
│  ────────────────────────────  │
│    ×  Clear override           │
│       Steam says IN PROGRESS   │
└────────────────────────────────┘
```

**B · no override — the same box, filled with the fact the player is about to override**

```
┌────────────────────────────────┐
│  Return of the Obra Dinn       │
│                                │
│    ○  NOT STARTED              │
│  ▌ ◐  IN PROGRESS              │
│    ◉  ZERADO                   │
│    ⊘  ABANDONED                │
│  ────────────────────────────  │
│    Steam says IN PROGRESS.     │
│    Your choice overrides it.   │
└────────────────────────────────┘
```

**C · a hand-added copy — no store has an opinion**

```
┌────────────────────────────────┐
│  Chrono Trigger                │
│                                │
│  ▌ ○  NOT STARTED              │
│    ◐  IN PROGRESS              │
│    ◉  ZERADO                   │
│    ⊘  ABANDONED                │
│  ────────────────────────────  │
│    Added by hand. No store to  │
│    ask. Every state is yours.  │
└────────────────────────────────┘
```

### 3.1 · The 34 × 11 budget, to the cell

| | Cols | | Rows |
|---|---|---|---|
| border | 1 + 1 | border | 1 + 1 |
| `BorderInsetX` (D-06-1) | 2 + 2 | vertical inset | **0** |
| **content** | **28** | **content** | **9** |
| | **34** | | **11** |

| Content row | Contents | Cells |
|---|---|---|
| 1 | the game's title, truncated to 28 | ≤ 28 |
| 2 | `InterElementGap` | — |
| 3–6 | the four state items | 16 each |
| 7 | a decorative hairline, 28 cells of `─` | 28 |
| 8–9 | **either** the fifth item + its consequence line **or** the provenance note | ≤ 28 |

**Item geometry — the same as the ledger row's, deliberately** (`01-design-system.md` §3.1):

```
▌ ○  NOT STARTED
│ │  │
│ │  └ label field — 11 cols, the ratified label
│ └ glyph field — 2 cols, width-aware padded
└ focus field — 2 cols
```
`2 + 2 + 1 + 11 = 16` cells. The `×` of *Clear override* sits in the **same 2-column glyph
field**, so all five items share one label edge at content column 6. `×` U+00D7 is **Ambiguous**
width (verified, UCD 16.0.0) — which is exactly why the field is padded rather than counted.

> **D-06-1 · A bordered surface is inset 2 columns each side and 0 rows, fixed at every tier.**
> Read from `01-design-system.md` §6.2's own ratified anatomy, where the first content row sits
> directly under the top border and the left inset is two spaces. Fixed rather than tracking
> `InnerPaddingX` because a bordered box that breathes differently at different terminal widths
> looks resized rather than designed. Shared with `Z-05-game-detail.md` §4.
> **Proposed as `space.BorderInsetX = 2` — see §18.**

> **D-06-2 · The box is 11 rows in every variant, and variant B fills it rather than shrinking.**
> `02-composition.md` §2.4 binds `34 × 11` as **fixed**, and the five-item composition needs all
> nine content rows. When there is no override, rows 8–9 carry the provenance note instead. That
> is not padding: it names the fact the player is about to override, which is precisely the fact
> *Clear override* will later restore — so the overlay teaches its own model at the one moment
> the player is looking at it. **A box that changes height between two adjacent presses of the
> same key reads as a glitch; one that changes content reads as an instrument.**

---

## 4 · Mockup — composited over `Z-04` at 80 × 24

The overlay is **horizontally centred in the terminal** and **vertically centred inside
`BodyRect`**, not inside the terminal. At 80 × 24 that is **column 24, terminal row 9**.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   247 games · 198 not started · 12 in progress · 6 zerado · 31 abandoned       │
│                                                                                │
│     STATE           TI┌────────────────────────────────┐         HOURS   SRC   │
│     ◐  IN PROGRESS  Ba│  Return of the Obra Dinn       │           87h   STM   │
│     ○  NOT STARTED  Bl│                                │            0h   STM   │
│     ◉  ZERADO       Ce│  ▌ ○  NOT STARTED              │           14h   STM   │
│     ○  NOT STARTED  Ch│    ◐  IN PROGRESS              │             —   PHY   │
│     ◐  IN PROGRESS  Da│    ◉  ZERADO                   │           63h   STM   │
│     ⊘  ABANDONED    Di│    ⊘  ABANDONED                │           22h   STM   │
│     ◉  ZERADO       Ha│  ────────────────────────────  │           58h   STM   │
│     ○  NOT STARTED  Ho│    ×  Clear override           │           41h   STM   │
│     ◐  IN PROGRESS  Ou│       Steam says IN PROGRESS   │           12h   STM   │
│   ▌ ◉  ZERADO       Re└────────────────────────────────┘            9h   STM   │
│     ⊘  ABANDONED    Sekiro: Shadows Die Twice                       3h   STM   │
│     ○  NOT STARTED  Sid Meier's Civilization VI: Gathering St…      0h   STM   │
│   ROWS  4–15 of 247                                                            │
│   ↑↓ move   ⏎ apply   esc cancel                                               │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

> **D-06-3 · Vertically centred in `BodyRect`, not in the terminal.** Centring in the terminal
> puts the top border on the pinned summary row, and R-10(c) requires that row to be on screen
> at any row count. Centring in `BodyRect` clears it in one line of arithmetic —
> `top = BodyRect.y + (BodyRect.h − 11) / 2` — at every tier, with no special case.

**WCAG 2.4.11, Focus Not Obscured, satisfied twice over:**

1. **By position.** The focused row here — `▌ ◉ ZERADO  Return of the Obra Dinn  9h  STM` at
   terminal row 20 — is entirely clear of the overlay.
2. **By width, always.** The overlay is 34 of 80 columns and centred, so even when the focused
   row falls inside rows 9–19, the row's **focus marker (column 4), its full state chip (columns
   6–19) and its playtime and source (columns 61–77) remain visible.** The row is never
   *entirely* obscured at any scroll position, at any tier ≥ Narrow. At Tiny the overlay becomes
   a route, so the question does not arise.

**The dimmed backdrop:** `02-composition.md` §2.4 permits dimming *where dimming is available*.
At the 16-colour floor and under `NO_COLOR` **there is no dim**, so the overlay is separated by
its **border** — `--z-border-strong` `#64748B`, **4.08:1**, which satisfies WCAG 1.4.11 as a
control boundary. `--z-border` at **1.53** may never do this job. **The border is the separator;
the dim is the courtesy.** A composition that depends on the dim is a composition that vanishes
at 16 colours.

---

## 5 · Mockup — a route at Tiny, 32 × 24

Below 40 columns `34 × 11` does not fit inside a `30 × 21` body, so the overlay is **pushed as a
route**. Behaviourally identical; only the composition changes.

```
┌────────────────────────────────┐
│ SET STATUS                     │
│                                │
│ Return of the Obra Dinn        │
│                                │
│ ▌ ○  NOT STARTED               │
│   ◐  IN PROGRESS               │
│   ◉  ZERADO                    │
│   ⊘  ABANDONED                 │
│                                │
│ ────────────────────────────── │
│                                │
│   ×  Clear override            │
│      Steam says IN PROGRESS    │
│                                │
│                                │
│                                │
│                                │
│                                │
│                                │
│                                │
│                                │
│                                │
│                                │
│ ⏎ apply  esc cancel            │
└────────────────────────────────┘
```

Body **30 × 21** · `leftInset` **1** · `OuterMarginX` and `OuterMarginY` both **0** · the band
collapses to the title row alone. There is no box border here — the frame is the boundary — so
the hairline runs the full 30 columns and the content sits at `leftInset` like every other
route. **`Esc` pops the route** rather than dismissing an overlay; from the player's side the
key does the same thing.

**At Narrow (40 × 24) it is still an overlay**: 34 fits inside a 36-column body with the
2-column outer inset the frame already provides. It is centred in `BodyRect` exactly as at Wide.

---

## 6 · Visual hierarchy

| # | What | Channel |
|---|---|---|
| **1** | **The four state items as a column** | Four rows of identical geometry with a single label edge — the shape says *choose one* before a word is read |
| **2** | **The focused item** | `▌` in the gutter (position) + **bold** (weight) + `--z-primary` amber on the marker (colour). The only bold row in the box |
| **3** | **The game's title** | Content row 1, `--z-text` (16.65 — the brightest text in the box), sentence case against four uppercase labels |
| **4** | **The four state colours** | Each item in its own state colour. `◉ ZERADO` in cyan is the only cyan in the box, and it is **data**, not emphasis |
| **5** | **The hairline** | `--z-border` — a decorative rule between two blocks of content, which is exactly and only what 1.53:1 is permitted to draw |
| **6** | The fifth item / the provenance note | `--z-text-secondary`, below the rule, subordinate by position and by weight |

**The one thing the player should see first is the four-item column; the second is which one
they are on.** The title is third on purpose: they already know which game they pressed `s` on,
and the row behind the overlay is still visible.

> **`ZERADO` is not made prominent, and that is deliberate.** `01-screen-inventory.md` §5:
> *"`Z-06` must not make zerado harder to reach than any other state."* The inverse is equally
> true — making it **easier** would be the product nudging, and `05-state-machine.md` §4 is
> explicit that this transition is the player's alone. Four items, identical geometry, one
> keystroke apart. The cyan on `◉ ZERADO` is the state's own colour, present on every screen
> that renders that state, not an emphasis this screen adds.

---

## 7 · Every applied spacing token, by name

| Token | Tiny (route) | Narrow | Standard | **Wide** | ExtraWide |
|---|---|---|---|---|---|
| `OuterMarginX` | 0 | 1 | 2 | **2** | 2 |
| `OuterMarginY` | 0 | 1 | 1 | **1** | 1 |
| `InnerPaddingX` | 1 | 1 | 1 | **1** | 2 |
| `InnerPaddingY` | 0 | 1 | 1 | **1** | 1 |
| `InterElementGap` | 1 | 1 | 1 | **1** | 1 |
| `HeaderBand(tier, false)` | **1** | 3 | 3 | **3** | 3 |
| `leftInset` | **1** | 2 | 3 | **3** | 4 |
| **Overlay footprint** | *(route — the full body)* | **34 × 11** | **34 × 11** | **34 × 11** | **34 × 11** |
| **`BorderInsetX`** (D-06-1) | *(none — no border)* | **2** | **2** | **2** | **2** |

**Applied:**

| Surface | Token | Value |
|---|---|---|
| title → items | `InterElementGap` | 1 row (content row 2) |
| items → hairline → the last block | the hairline row itself | 1 row, no extra gap either side — the rule *is* the separation |
| Inside the overlay border | `BorderInsetX` | 2 cols each side, 0 rows |
| Overlay horizontal origin | `(width − 34) / 2` | column 24 at 80 |
| Overlay vertical origin | `BodyRect.y + (BodyRect.h − 11) / 2` | terminal row 9 at 80 × 24 |
| Tiny route: content left | `leftInset` | column 2 |
| Footer while open | the frame's reserved footer row | the **overlay's** keys — §11 |

**Zero magic numbers.** `34 × 11` is the spine's binding footprint; `BorderInsetX` is a named
token proposed in §18; the 16-cell item geometry is the ratified chip plus a 2-column focus
field.

---

## 8 · Colour, glyph and label for every state shown

Every ratio read from the brand manual's measured table (§4.2). None estimated.

| Element | Token | Hex | ANSI-256 | 16-colour | Glyph | ASCII | Label | Ratio |
|---|---|---|---|---|---|---|---|---|
| Item — not started | `--z-state-not-started` | `#A5A29B` | **247** | `white` | `○` U+25CB | `[ ]` | `NOT STARTED` | **7.62** AA |
| Item — in progress | `--z-state-in-progress` | `#FFB000` | **214** | `bright yellow` | `◐` U+25D0 | `[~]` | `IN PROGRESS` | **10.59** AAA |
| Item — zerado | `--z-state-zerado` | `#19E0FF` | **45** | `bright cyan` | `◉` U+25C9 | `[*]` | `ZERADO` | **12.15** AAA |
| Item — abandoned | `--z-state-abandoned` | `#C77DFF` | **177** | `bright magenta` | `⊘` U+2298 | `[x]` | `ABANDONED` | **7.21** AA |
| Item — clear override | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | `×` U+00D7 (**A**) | `x` | `Clear override` | **9.36** AAA |
| Focus marker | `--z-primary` | `#FFB000` | **214** | `bright yellow` | `▌` U+258C (**A**) | `>` | — | **10.59** AAA |
| The game's title | `--z-text` | `#E9EEF5` | **255** | `bright white` | — | — | — | **16.65** AAA |
| The consequence line / provenance note | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | — | — | — | **9.36** AAA |
| **Overlay border** | `--z-border-strong` | `#64748B` | **67** | `bright black` | `┌─┐` | `+-+` | — | **4.08** — meets 1.4.11 |
| The hairline | `--z-border` | `#2A3342` | **236** | `black` | `─` | `-` | — | 1.53 — **decoration between blocks, and it bounds nothing** |
| Backdrop dim, where available | `--z-surface-overlay` | `#1D2532` | ***underived*** | *(vanishes)* | — | — | — | — |

> ***underived*** — `--z-surface-overlay` `#1D2532` has **no derived ANSI-256 index**.
> **Interim: the backdrop is not dimmed at all.** The overlay is separated by its
> `--z-border-strong` border, which is the mechanism that has to work at 16 colours and under
> `NO_COLOR` anyway. This invents nothing, and it means the interim render and the 16-colour
> render are the same render. Owner: `fft-brand-architect` (`00-design-brief.md` §9).

**Co-render holds on every item: colour AND glyph AND label, all three.** The tightest CVD pair
in the system — **zerado × abandoned at ΔE 11.9 under deuteranopia** — is on screen here,
adjacent, on rows 5 and 6. This is *the* screen where the glyph and the word carry load rather
than merely reinforce, and it is why neither may ever be dropped for density.

### 8.1 · Audio — the fourth channel

The annunciator lives on the **header band's title row** of whatever is beneath this overlay
(D-A1, `Z-04-library.md` §7.1) — `Z-04`'s or `Z-05`'s. **The overlay never draws it**, and
never covers it: the band is above `BodyRect` and the overlay is centred inside `BodyRect`
(D-06-3). At Tiny, where `Z-06` is a route, it renders on `Z-06`'s own title row.

**This screen owns the product's single most likely sound event: a game becoming *zerado*.**

> **The fourth channel is mandatory here.** If a sound is ever attached to marking a game
> *zerado*, its visible carrier is already specified and must ship with it: the **result line**
> of §10.4, printed in the pinned summary row. Audio may accompany that line; it may never
> replace it, and it may never be the only signal that the write happened. `m` mutes the sound
> and changes **nothing** about what is on screen.

---

## 9 · The full state table

| # | State | Trigger | Composition | Copy |
|---|---|---|---|---|
| **S1** | **First run — unreachable by construction** | Library empty | **`Z-06` cannot be opened.** `s` on an empty library means *connect Steam* (`01-design-system.md` §10.1); there is no game whose status could be set | — |
| **S2** | **Override set — five items** | `status_manual IS NOT NULL` | Variant A. Focus starts on the **current effective state** | §10.1 |
| **S3** | **No override — four items** | `status_manual IS NULL` | Variant B. Rows 8–9 carry the provenance note | §10.2 |
| **S4** | **A hand-added copy** | `Capabilities.Progress = false` | Variant C. The provider has no opinion to name | §10.3 |
| **S5** | **A hand-added copy WITH an override** | `physical` + `status_manual IS NOT NULL` | Variant A, and the consequence line reads `back to NOT STARTED` — there is no provider to attribute it to | §10.1 |
| **S6** | **Applying — the write is in flight** | `⏎` pressed | The overlay **dismisses immediately**; the row's chip beneath reads `ZERADO…` until SQLite confirms | §10.5 |
| **S7** | **The write failed** | SQLite error | The overlay is already gone. The failure is `Z-11 Fatal error` — an unwritable library file is not a screen-level degrade, it is the program's promise broken | — |
| **S8** | **A transition into *zerado* landed** | `⏎` on `◉ ZERADO` | The host's pinned summary is replaced by the result line until the next keypress | §10.4 |
| **S9** | **Cancelled** | `Esc` | Dismiss, discarding the uncommitted choice. **Nothing is written.** Focus returns exactly where it was | — |
| **S10** | **Opened from `Z-05`** | `s` on the detail view | Identical box, drawn over `Z-05` in either host. The overlay does not know or care which screen is beneath it | — |
| **S11** | **Tiny** | `< 40` cols | §5 — a route | — |
| **S12** | **The title overflows 28 columns** | Any long title | Truncated with `…`, inside the padded field | §10.6 |
| **S13** | **Below the refusal floor** | `< 24` cols or `< 8` rows | The program refuses — `Z-04-library.md` §11.3 | — |

**There is no confirmation, in any state.** `01-design-system.md` §13.1 lists exactly three
destructive actions and **marking a game `ABANDONED` is not one of them** — it is a reversible
state change. *"Confirming reversible actions trains players to press `y` without reading, which
is how the irreversible one eventually gets confirmed too."* All twelve ordered transitions are
legal (`05-state-machine.md` §3) and the product never refuses one.

**`ZERADO` is never set automatically, in any state, by anything.** Not by 100 % achievements,
not by a known credits-roll playtime, not by a store's own completed flag. Permanently out of
scope.

---

## 10 · The exact copy — ready to paste

### 10.1 · S2 / S5 — the fifth item and its consequence line

```
  ×  Clear override
     Steam says IN PROGRESS
```

**The second line is mandatory** — `05-state-machine.md` §5: *"clear override without that line
is a button whose effect the player cannot predict."* It has two forms:

| Provider | Second line | Cells |
|---|---|---|
| reports playtime (Steam) | `     <Provider> says <STATE>` — e.g. `     Steam says IN PROGRESS` | 27 |
| does not (`physical`) | `     back to <STATE>` — e.g. `     back to NOT STARTED` | 24 |

Both fit the 28-cell content width. When a provider name would push the line past 28, it drops
to the second form; the **state name never truncates**, because it is the whole point of the
line.

**Be honest about the consequence, in the copy and in the model:** clearing an override on a
game with playtime > 0 makes it `IN PROGRESS` **immediately**. Choosing `NOT STARTED`
explicitly is a *different* action — it stores `not_started` as a manual value and it sticks.
Both are useful; they are not the same; and that is exactly why the model needs a nullable field
and a separate clear.

### 10.2 · S3 — no override

```
  Steam says IN PROGRESS.
  Your choice overrides it.
```

25 and 27 cells. Two facts, no adornment, no exclamation mark.

### 10.3 · S4 — a hand-added copy

```
  Added by hand. No store to
  ask. Every state is yours.
```

28 and 28 cells. It explains the absence rather than apologising for it, and it never suggests
an action that does not exist.

### 10.4 · S8 — the moment the product exists for

Brand manual §8, voice example 4, is the copy. It renders **in the host's pinned summary row**,
not inside the overlay — the overlay is already gone by then.

| Situation | Copy, exact |
|---|---|
| Playtime known | `Zerado. 41 hours. Sixth this year.` |
| Playtime known, first of the year | `Zerado. 9 hours. First this year.` |
| Playtime not reported (a hand-added copy) | `Zerado. Sixth this year.` |
| Past tenth | `Zerado. 41 hours. 14th this year.` |

Three facts. No celebration of a database write. `Sixth this year` is the ordinal of *zerado*
games whose `status_changed_at` falls in the current calendar year; ordinals are spelled to
*tenth*, then numeric with a suffix. It **persists until the next keypress** — no timer, no
dwell constant, no WCAG 2.2.1 time limit — and it renders **only** for a transition into
`zerado` (`Z-04-library.md` §10.5, D-04-3).

### 10.5 · S6 — the write in flight

```
◉  ZERADO…
```

The chip's label carries the `…` until the write confirms. **Never an optimistic silent change**
(`01-design-system.md` §3.4).

### 10.6 · S12 — a title too long for 28 columns

```
  Sid Meier's Civilization VI…
```

Truncated with `…` inside the padded field. **Acceptable here and only here:** in the ledger the
title is the identity column and in `Z-05` it wraps, but in this overlay it is *context* — the
row the player pressed `s` on is still visible behind and beside the box (§4).

### 10.7 · Copy notes

- **Casing** — `ZERADO` on the chip and the item, because this is the interface. `Zerado.` in
  the result line is the **product** speaking, which is why it is capitalised there and nowhere
  else on this screen.
- **Say the number** — `41 hours`, `9h`. Never *a lot of hours*.
- **No exclamation marks. No emoji. Never call the user a "gamer".**
- **Nothing here claims an unbuilt capability.** No "suggest marking this zerado", no
  achievement percentage, no mood tag. `05-state-machine.md` §4 permits a *suggestion* in
  **Phase 2**; Phase 1 does not carry one.
- **Type-neutral where equally natural** — `Clear override`, `Every state is yours.`,
  `Added by hand.` Phase 1 says *game* where a game is what the player is looking at.

---

## 11 · The key map

Small, closed, and printed where the player can see it.

| Key | Does | Note |
|---|---|---|
| `↑` `↓` `k` `j` | Move between the items | Wraps at the ends — five items is short enough that wrapping helps rather than disorients |
| `g` / `G` | First / last item | |
| `⏎` | **Apply the focused item and dismiss** | Writes to SQLite before the overlay closes |
| `Esc` | **Cancel and dismiss.** Nothing is written | Focus returns exactly where it was |
| `Ctrl-C` | Quit | Always |
| `?` | **Inert while the overlay is open** | See below |
| `q` `s` `r` `a` `/` `,` `m` `Tab` | **Inert while the overlay is open** | An overlay owns its keys completely |

> **D-06-4 · `?` is inert while `Z-06` is open as an overlay.** WCAG 3.2.6 Consistent Help asks
> that help be reachable *from every page in the same relative order*; a two-keystroke modal is
> not a page, and it **prints its own key hints in the footer**, which is the stronger form of
> help — the answer is already on screen. Pushing a full-screen help route over an unresolved
> decision would discard the player's place for no gain, and `04-navigation-and-focus.md` §1
> rule 2 forbids stacking a second surface on an overlay anyway.
> **At Tiny, where `Z-06` is a route, `?` works normally** and `Esc` from `Z-10` returns here.
> `Z-10-help-and-key-map.md` §8 lists this screen's keys under that condition.

### 11.1 · The footer while the overlay is open

**The frame's reserved footer row shows the OVERLAY's keys, not the route's.** A footer that
lists `/ filter` while a modal owns the keyboard is a footer that lies
(`04-navigation-and-focus.md` §6).

```
↑↓ move   ⏎ apply   esc cancel
```

**30 cells**, 3-space separators, at every tier ≥ Narrow. At Tiny, where `Z-06` is a route, the
frame's own footer carries `⏎ apply  esc cancel` (**19 cells**, 2-space separators per the
ladder in `Z-04-library.md` §9.2).

**This is also why the overlay needs no key-hint row of its own** — which is what makes all nine
content rows available for the five-item composition.

### 11.2 · No direct-selection keys, and why

`1`–`9` are **reserved and unbound** for Phase 2 quick filters
(`04-navigation-and-focus.md` §3.1) and must not be claimed here — *"a key that means one thing
in Phase 1 and another in Phase 2 is a betrayal of muscle memory."* Letter shortcuts (`n`, `i`,
`z`, `a`) were considered and rejected: five items reachable in at most two `↓` presses do not
need a second input scheme, and `n` and `p` are already reserved for `Z-05`.

---

## 12 · The focus model, and `Esc`

### 12.1 · Focus is trapped inside the overlay — correctly

`01-design-system.md` §13.3: focus is trapped *within* the dialog while it is open, **which is
correct, because there is a documented way out.** `Esc` always cancels, `⏎` always applies.
That is what separates a modal from a keyboard trap (WCAG 2.1.2).

### 12.2 · Where focus starts

**On the current effective state**, always. So the player sees where the game *is* before they
move — and the number of keystrokes to any state is the same from every starting point, which is
what keeps `ZERADO` no harder and no easier to reach than the other three.

When the fifth item exists, it is **never** the initial focus: it is a correction, not a
default.

### 12.3 · How focus is shown — three channels, any two sufficient

| Channel | Focused | Not focused |
|---|---|---|
| **Position** | `▌` U+258C in the 2-cell gutter field (ASCII `>`) | two spaces |
| **Weight** | **bold** | normal |
| **Colour** | `--z-primary` amber on the marker | none |

**Amber, not cyan** — this overlay spends no cyan on focus, so the only cyan in the box is the
`◉ ZERADO` item's own state colour. Never by background fill: at 16 colours the surface ramp
collapses to `black` and a highlight bar would disappear. **Never removed.**

### 12.4 · `Esc`, exhaustively

| Context | `Esc` does | Then |
|---|---|---|
| The overlay is open | **Dismiss, discarding the uncommitted choice** | Focus returns exactly where it was — the same row of `Z-04`, or `Z-05` |
| `Z-06` as a route at Tiny | **Pop the route** | Same outcome from the player's side |
| While a write is in flight | Not reachable — the overlay dismissed on `⏎` | — |

**The safe branch is always the default.** There is no pre-selected destructive item, because
there is no destructive item.

---

## 13 · 40-column behaviour, and the refusal floor

**At Narrow (40 × 24) `Z-06` is still an overlay.** The 34-column box fits inside the 36-column
body, centred in `BodyRect`. Nothing about the composition changes: same nine content rows, same
16-cell items, same copy.

**At Tiny (< 40) it becomes a route** — §5. `02-composition.md` §2.4: *"At Tiny the overlay does
not fit inside `30 × 21` with its margins, so at Tiny an overlay becomes a route — pushed,
full-screen, popped with `Esc`. Behaviourally identical; only the composition changes."*

**Never sheds, at any width:** the four items · each item's glyph **and** label · the fifth item
when it exists · **its consequence line** · the focus marker · the key hints.

**The refusal floor is the program's:** below **24 columns or 8 rows**, one frameless sentence
and, at start-up, exit `2`. A running session keeps running.

### 13.1 · The one width hazard on this screen, named

The border and the hairline are **box-drawing characters, and the whole family is East-Asian
Ambiguous** (verified: `┌` U+250C, `─` U+2500, `│` U+2502, `┐` U+2510, `└` U+2514, `┘` U+2518 —
all **A**). On a terminal configured `ambiguous-width=double` the box renders **uniformly**
doubled — 68 cells, not sheared — because every character in it is one class. At Wide that still
fits inside a 74-column body; **at Narrow it does not.**

**The escape hatch is `ZERADO_ASCII=1`**, which `01-design-system.md` §1.2 rule 3 already
defines as the forced-narrow path and as the automatic fallback when the terminal does not
report Unicode capability. **This spec extends it to box drawing and the focus marker**, because
the state column is not the only place the hazard lives:

```
+------------------------------+
|  Return of the Obra Dinn     |
|                              |
| >[ ]  NOT STARTED            |
|  [~]  IN PROGRESS            |
|  [*]  ZERADO                 |
|  [x]  ABANDONED              |
|  --------------------------  |
|  x   Clear override          |
|      Steam says IN PROGRESS  |
+------------------------------+
```

> **D-06-5 · Under `ZERADO_ASCII=1`, box drawing falls back to `+ - |`, the focus marker to `>`,
> the clear-override glyph to `x`, and the state column to the ratified `[ ] [~] [*] [x]`.**
> All are ASCII-narrow and immune by construction. `01-design-system.md` §1.7 already names `>`
> as the marker's ASCII fallback and §1.2 already names the state column; this decision closes
> the remaining two. **This is a cross-cutting design-system addition, not a screen decision —
> see §18 item 2.**

---

## 14 · `NO_COLOR` — rendered, not asserted

Zero SGR sequences. Variant A, character for character:

```
┌────────────────────────────────┐
│  Return of the Obra Dinn       │
│                                │
│  ▌ ○  NOT STARTED              │
│    ◐  IN PROGRESS              │
│    ◉  ZERADO                   │
│    ⊘  ABANDONED                │
│  ────────────────────────────  │
│    ×  Clear override           │
│       Steam says IN PROGRESS   │
└────────────────────────────────┘
```

| Information | Carried without colour by |
|---|---|
| The four states | The **glyph** `○ ◐ ◉ ⊘` **and the label word** — the label *is* the text alternative |
| Which item is focused | The `▌` marker (position) + **bold** (weight) |
| Which is `Clear override` | The `×` glyph, the label, **and its position below the rule** |
| The consequence of clearing | The sentence `Steam says IN PROGRESS` — a **word channel**, which is why it can never be reduced to a colour |
| The overlay's boundary | The **border**, not a dim and not a fill. This is the case that proves the border is the separator |
| The overlay's own keys | The footer line, which is text |

**Nothing is lost.** Run the `02-colour-budget.md` §3.3 cross-check — `NO_COLOR=1` — and the
same nine content rows carry the same nine facts.

**The 16-colour floor is the same picture.** `--z-surface-overlay` collapses to `black`, so the
dim vanishes; the border and the spacing carry the modality, exactly as designed. The four
states resolve to four distinct slots — `white` · `bright yellow` · `bright cyan` ·
`bright magenta` — with **no collisions**.

---

## 15 · Colour budget declaration

| State | STATE cyan (uncounted) | Focus ring (exempt) | **CHROME CYAN** | Verdict |
|---|---|---|---|---|
| S2 five items | **1** — `◉` + `ZERADO` on its own item | none — the item cursor is **amber** | **0** | **PASS** |
| S3 four items | **1** | none | **0** | **PASS** |
| S4 hand-added | **1** | none | **0** | **PASS** |
| Tiny route | **1** | none | **0** | **PASS** |

**`Z-06` spends ZERO chrome cyan, in every variant.** The temptation this budget exists to catch
is exactly here: `◉ ZERADO` is the product's name and its payoff, and it would be very easy to
give it a cyan border, a cyan heading, or a cyan `⏎ apply` hint. **All three would be fails.**
Cyan is never a border, never a heading, never a title. The one cyan in this box is the
`ZERADO` **state cell**, which is data and is not budgeted.

**Amber allow-list entries used:** 3 the `IN PROGRESS` state, plus the focus marker. At Tiny,
also 1 the screen title. **Not used:** 2 readout labels (there are none — the item labels are
state labels, not readouts) · 4 progress fill · 5 key hints — the footer's `⏎ apply` is
`--z-text-secondary` **chrome, not amber, and certainly not cyan** · 6 the terminal mark ·
7 the degrade banner · 8 the filter sigil.

**Amber ceiling:** the box is 34 × 11 = 374 cells; amber spends 16 (the `IN PROGRESS` item) + 1
(the marker) = **17**, about 4.5 %.

**Red: none.** No scanner, no annunciator, no error text. **This screen raises no confirmation,
so it has no destructive annunciator to draw** — which is the point of §13.1's closed list.

**`--z-border` draws the hairline only**, between two blocks of content, bounding nothing. The
**overlay's own border is `--z-border-strong` at 4.08**, because it is the boundary of a control
and WCAG 1.4.11 applies. A `--z-border` box around this overlay would be an automatic fail.

---

## 16 · Reuse verdict, per element

| Element | Verdict | Why |
|---|---|---|
| **The overlay chrome** (border, modality, focus trap, `esc` cancel) | **`huh` — fits**, restyled to Zerado tokens | `01-design-system.md` §13.5. It inherits `huh`'s built-in breathing room. **Its default theme carries its own palette and must not ship** |
| The five-item list | Build fresh · `lipgloss` | It is five fixed rows with the ratified chip geometry. `bubbles/list` brings pagination, filtering and a status bar this box has neither room nor need for |
| State chip | **The shared chip** from `01-design-system.md` §3 — not a second implementation | A `ZERADO` chip must look and read identically everywhere (WCAG 3.2.4) |
| The hairline | `lipgloss` — one repeated rune, width-aware | |
| Backdrop dim | `lipgloss`, **where available**; the composition must not depend on it | It vanishes at 16 colours and under `NO_COLOR` |
| Overlay positioning | `lipgloss` place, against `BodyRect` | D-06-3 |
| Result line | The **host's** status bar — `Z-04` or `Z-05` renders it | The overlay is gone by then |
| Confirmation dialog | **Not built, not needed** | `01-design-system.md` §13.1 — nothing here is destructive |
| Spring animation | **Not `harmonica`, and no animation at all** | There is no motion on this screen. The scanner is the brand's one motion and it is for indeterminate waits |

---

## 17 · Upstream findings

| # | Finding | Where | Owner |
|---|---|---|---|
| 1 | **`05-state-machine.md` §5's second line — `Back to what Steam says: IN PROGRESS` — is 36 characters and does not fit the 28-cell content width** of the binding `34 × 11` overlay. This spec keeps the *requirement* (name what the game will become) and shortens the wording to `Steam says IN PROGRESS` (27) | `05-state-machine.md` §5 | `fft-tui-architect` |
| 2 | **`02-composition.md` §2.4 binds `34 × 11` as fixed but the five-item composition needs all nine content rows**, leaving none for a key-hint line. Resolved by putting the overlay's keys in the frame's reserved footer row (§11.1). Recorded because a builder reading §2.4 alone would look for room that is not there | `02-composition.md` §2.4 | `fft-tui-architect` |
| 3 | **`01-design-system.md` §6.2's bordered-surface inset is drawn (2 cols / 0 rows) but never named**, and it does not equal `InnerPaddingX` at any tier except ExtraWide. Adopted as D-06-1; proposed as a token in §18 | `01-design-system.md` §6.2 | `fft-design-architect` |
| 4 | **`ZERADO_ASCII=1` is defined for the state column only** (`01-design-system.md` §1.2 rule 3), but the box-drawing family is equally Ambiguous and an overlay is the place it bites hardest. Extended by D-06-5; needs adopting into the design system | `01-design-system.md` §1.2 | `fft-design-architect` |
| 5 | **`s` is bound to *set status* globally and to *connect Steam* on an empty library** (`01-design-system.md` §10.1). It is why S1 exists as an *unreachable* row rather than a screen. See `Z-04-library.md` §18 item 1 | `01-design-system.md` §10.1 vs `04-navigation-and-focus.md` §3 | founder |
| 6 | `03-designer-manual.md` §5.11 verdict 3 still reads as a permanent rejection of the audio subsystem, superseded by founder direction relayed 2026-08-25 | `03-designer-manual.md` §5.11 | `fft-design-architect` |

---

## 18 · Open for the founder

1. **The shortened consequence line.** `05-state-machine.md` §5 writes
   `Back to what Steam says: IN PROGRESS`; it does not fit the binding 28-cell width and this
   spec renders `Steam says IN PROGRESS`. The *requirement* — name what the game will become —
   is fully met, but the exact wording changed, and §5 is a spine document.
   **Confirm the shortened form, or widen the overlay.**
2. **`space.BorderInsetX = 2`, fixed at every tier.** Shared by this overlay and `Z-05`'s detail
   pane. It is currently a number read off `01-design-system.md` §6.2's drawing, which is exactly
   the kind of value the canon says must be a token. Owner if approved: `fft-design-architect`.
3. **`?` is inert while an overlay is open** (D-06-4). It is the one place in the product where a
   globally-available key does nothing, and it is a deliberate reading of WCAG 3.2.6 — the
   overlay prints its own keys, so the help *is* on screen. **Confirm**, because the alternative
   (dismiss the overlay, then push help) silently discards an in-progress decision.

---

## 19 · Design decisions made in this spec

| # | Decision | Reason |
|---|---|---|
| **D-06-1** | A bordered surface is inset 2 cols each side, 0 rows, fixed at every tier (§3.1) | Read from `01-design-system.md` §6.2's ratified anatomy; a box that breathes differently at different widths looks resized, not designed |
| **D-06-2** | The box is 11 rows in all three variants; variant B fills rows 8–9 with the provenance note rather than shrinking (§3.1) | `34 × 11` is bound as fixed; and a box that changes height between two presses of the same key reads as a glitch. The note also teaches the override model at the one moment the player is looking at it |
| **D-06-3** | Vertically centred in `BodyRect`, not in the terminal (§4) | Clears the pinned summary row, which R-10(c) requires on screen at any row count — in one line of arithmetic, at every tier, with no special case |
| **D-06-4** | `?` is inert while the overlay is open; it works normally when `Z-06` is a route at Tiny (§11) | A modal is not a page, it prints its own keys, and `04-navigation-and-focus.md` §1 rule 2 forbids stacking on an overlay |
| **D-06-5** | `ZERADO_ASCII=1` also swaps box drawing to `+ - |`, the marker to `>` and `×` to `x` (§13.1) | The box-drawing family is Ambiguous-width and an overlay is where that bites hardest; the escape hatch must cover it |
| **D-06-6** | Focus starts on the current effective state; the fifth item is never the initial focus (§12.2) | The player sees where the game *is* before moving, and every state stays the same distance away — which is what keeps `ZERADO` no harder and no easier than the rest |
| **D-06-7** | The consequence line has two forms — `<Provider> says <STATE>` and `back to <STATE>` (§10.1) | A provider without a playtime capability has no opinion to attribute; naming one would invent a fact |

---

## 20 · Screen-specific acceptance criteria

Beyond `00-design-brief.md` §10 and `02-colour-budget.md` §10, both of which apply in full.

1. **The box measures exactly 34 × 11 in all three variants**, on an ANSI-stripped render, at
   Narrow, Standard, Wide and ExtraWide — and **does not change height** when the fifth item
   appears or disappears.
2. **The fifth item appears if and only if `status_manual IS NOT NULL`**, and its second line
   names the resulting state.
3. **The focused row of the screen beneath is never entirely obscured** (WCAG 2.4.11) — verified
   with the cursor at the top, middle and bottom of the visible list.
4. **The pinned summary row of `Z-04` is never covered** by the overlay, at any tier.
5. **All four states are the same number of keystrokes from every starting point.** No state has
   a direct key that another lacks.
6. **`Esc` writes nothing.** Verified against the SQLite file, not against the render.
7. **`⏎` writes before the overlay closes**, and the row's chip shows `ZERADO…` until confirmed.
8. **No confirmation dialog appears for any of the four states**, including `ABANDONED`.
9. **Chrome-cyan count is 0** in every variant, by `02-colour-budget.md` §3.1 — the only cyan on
   screen is the `◉ ZERADO` item's own state cell.
10. **`NO_COLOR=1` loses no information**, including which item is focused and what clearing the
    override would do.
11. **At 16 colours the overlay is still visibly separate from the screen beneath**, with the
    dim absent. The border does the work.
12. **Under `ZERADO_ASCII=1` the box renders in pure ASCII** and still measures 34 × 11.
13. **At Tiny the overlay is a route**, and `Esc` returns to the same row of `Z-04`.
14. **The footer lists the overlay's three keys and nothing else** while it is open.
15. **A transition into *zerado* prints the result line** in the host's summary row, and that
    line is present whether or not audio is enabled and whether or not it is muted.
16. **Founder-validated screenshot before merge**, at the six viewports of
    `03-responsive.md` §7 plus `NO_COLOR=1` and forced-16-colour at 80 × 24, **in all three
    variants**. No screenshot → not GOLDEN → no merge.
