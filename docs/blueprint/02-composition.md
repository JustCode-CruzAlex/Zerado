---
title: Zerado — composition and layout budget
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-02
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: implementation-plan
ticket: "#2"
---

# Composition and layout budget

What shape each screen is, and exactly how many rows and columns it has to work with.

---

## 1 · The frame is inherited, and the arithmetic is not negotiable

Zerado adopts the **FlowForge TUI Spacing Canon (#2435)** verbatim — the same named tokens and
the same fixed per-tier values. Nothing here is re-derived and no screen may use a magic number.

> **What is reused is the specification, not the code.** FlowForge's `v3/uikit` carries this
> canon as Go, but the FlowForge repository publishes **no `LICENSE` or `COPYING` file**
> (verified at source, 2026-08-25), so it grants no rights to a third party and its module path
> resolves to a repository an outside contributor cannot fetch. Zerado is a **public**
> repository. It therefore re-implements the token table as its own small `internal/uikit/space`
> package and cites this canon as the source. If FlowForge later publishes `uikit` under an
> open licence, this is a cheap decision to revisit — the values would already match.
>
> This is the ticket's *"reuse beats rebuild where it fits, and say so plainly where it does
> not"* answered plainly. It fits as a **spec**; it does not fit as a **dependency**.

### 1.1 · The token table (Spacing Canon §4, reproduced)

| Token \ Tier | Tiny `<40` | Narrow `40–59` | Standard `60–79` | Wide `80–119` | ExtraWide `120+` |
|---|---|---|---|---|---|
| `OuterMarginX` (cols/side) | 0 | 1 | 2 | 2 | 2 |
| `OuterMarginY` (rows T/B) | 0 | 1 | 1 | 1 | 1 |
| `InnerPaddingX` (cols/side) | 1 | 1 | 1 | 1 | 2 |
| `InnerPaddingY` (rows T/B) | 0 | 1 | 1 | 1 | 1 |
| `InterElementGap` (rows) | 1 | 1 | 1 | 1 | 1 |
| `HeaderBandHeight` (base) | 1 | 3 | 3 | 3 | 3 |
| **`leftInset`** = `OuterMarginX + InnerPaddingX` | **1** | **2** | **3** | **3** | **4** |

### 1.2 · Decision — **Zerado screens carry no header subtitle**

The canon allows a fourth header row for a subtitle at Standard and above. **Zerado declines
it, on every screen.** Two reasons:

1. At the 80×24 floor that row is worth more to the body than to the header — it is a whole
   extra game row.
2. Zerado's context is *always a live count* (`247 games · 6 zerado · filter: mood=grind`) and a
   live value belongs in the pinned summary row, where it can change without redrawing the
   header band. Putting it in the header would make the header the only mutable chrome on screen.

> **The band's *content* — the breadcrumb form, the separator glyph and the title row — is
> deliberately not decided here.** It is a design-system question and is settled in
> [`../design/01-design-system.md`](../design/01-design-system.md) §2: the breadcrumb is
> `Zerado ✦ <Screen>` (two segments, because Zerado has no project tier and R-2 forbids faking an
> empty one), and the separator is `✦` rather than `›` because a chevron carries a false
> affordance in a terminal — it reads as a path you could type. The spine owns *how many rows*;
> the design system owns *what is in them*.

Consequence: `HeaderBand` is always the **3-row base**, `hasSubtitle` is always `false`, and
`BodyRect` never desyncs from the header — the failure mode the canon's single-sizer invariant
exists to forbid cannot occur here because the flag has one value.

### 1.3 · The body canvas — the numbers every screen spec starts from

> **Drawn:** sheet 09, [`ZRD-FRAME-01`](../adr/charts/ZRD-FRAME-01.chart.toml) —
> [svg](../adr/charts/svg/ZRD-FRAME-01.svg) · [cyanotype](../adr/charts/svg/ZRD-FRAME-01.cyanotype.svg).
> The 80 × 24 frame row by row, with every token named and the 74-column game row broken out.


```
BodyRect.w = width  − 2·leftInset
BodyRect.h = height − HeaderBandHeight − 1 (footer)
                    − 2·OuterMarginY − InnerPaddingY − InterElementGap
```

| Terminal | Tier | `leftInset` | **Body `w × h`** |
|---|---|---|---|
| 32 × 24 | Tiny | 1 | **30 × 21** |
| 40 × 24 | Narrow | 2 | **36 × 16** |
| 60 × 24 | Standard | 3 | **54 × 16** |
| **80 × 24** | **Wide** | **3** | **74 × 16** ← *the design floor* |
| 80 × 30 | Wide | 3 | 74 × 22 |
| 100 × 30 | Wide | 3 | 94 × 22 |
| 120 × 40 | ExtraWide | 4 | **112 × 32** |

**80 × 24 is the design width and the guaranteed floor.** Every screen is composed at
`74 × 16` first and enhanced upward. A screen that only reads well at 120 columns is not
finished; a screen that requires more than 80 is a defect.

> Tiny keeps more body rows than Narrow (21 vs 16) because at Tiny the canon collapses the
> header band to a 1-row title and drops both outer margins. That is deliberate starvation
> protection, and it means the **narrowest** terminal is not the **shortest** body.

---

## 2 · Composition per screen — Phase 1

**`R`** = a region that can take focus. The count matters: `Tab` only exists where `R ≥ 2`.

| ID | Screen | Composition | `R` | Regions |
|---|---|---|---|---|
| Z-01 | First run | Single-pane, content block **top-aligned**, not vertically centred | 1 | choice list |
| Z-02 | Connect a store | Single-pane form (`huh`) | 1 | field group |
| Z-03 | Sync | Single-pane readout, scanner sweep on one row | 1 | *(none while running)* |
| **Z-04** | **Library** | **≤119: single-pane list · ≥120: list ∥ detail** | 1 → 2 | list · detail pane |
| Z-05 | Game detail | **≤119: route, single-pane scroll · ≥120: the right pane of Z-04** | 1 | viewport |
| Z-06 | Set status | Overlay, centred, `34 × 11` fixed | 1 | choice list |
| Z-07 | Filter and search | Mode of Z-04 — takes body rows 1–2, focus moves into it | 2 | filter bar · list |
| Z-08 | Add a game by hand | Single-pane form (`huh`) | 1 | field group |
| Z-09 | Settings | Single-pane grouped form | 1 | field group |
| Z-10 | Help and key map | Single-pane table in a viewport | 1 | viewport |
| Z-11 | Fatal error | **Frameless.** Plain text, left-aligned, no chrome | 0 | — |

### 2.1 · Why the detail pane starts at 120 columns and not at 80

At Wide (80 cols) the body is 74 columns. Splitting it gives roughly 44 for a list and 28 for a
detail pane. The list's identity column — the game title, which TUI Design Manual **R-10(a)**
requires to be present and human-readable — would fall to about 15 columns after the state chip
and the playtime figure. Fifteen columns is `Return of the Ob…`, which fails R-10(a) on its face,
**and** the 28-column detail pane cannot hold a *sinopse*.

So the rule is: **two regions only when there is room for both to be correct.**

| Tier | Library composition | Where detail lives |
|---|---|---|
| Tiny · Narrow · Standard · **Wide** | Single-pane list | `Z-05`, a **route** reached with `Enter` |
| **ExtraWide (120+)** | List ∥ detail | `Z-05`, a **pane**; `Enter` moves focus into it instead of routing |

`Z-05` is therefore **one view with two hosts**. This is the single most important composition
decision in the bundle, because it is what makes the same detail spec buildable once and mounted
twice — and it is why `Z-05`'s spec must be written host-agnostic.

### 2.2 · Layout budget — `Z-04 Library` at the 80×24 floor

Body `74 × 16`:

```
 row  1   pinned summary            247 games   ○ 198  ◐ 12  ◉ 6  ⊘ 31
 row  2   (respiro)
 row  3   column header             STATE          TITLE                     HRS
 rows 4–15  the scroll region       12 game rows
 row 16   scroll position           ▄ 4–15 of 247
```

**Twelve game rows visible at the design floor.** Rows 1 and 16 are pinned *outside* the scroll
region, which is TUI Design Manual **R-10(c)** — the summary can never be pushed off the bottom
by a 400-game library.

Row column budget, 74 columns:

| Field | Cols | Note |
|---|---|---|
| focus field | 2 | **fixed width, padded at render time** — see §2.2.1 |
| state chip | 14 | glyph field **2** + space + label; `NOT STARTED` is the longest at 11 |
| gutter | 2 | |
| **title** | **42** | the identity column — R-10(a) |
| gutter | 2 | |
| playtime | 6 | right-aligned, `  41h` |
| gutter | 2 | |
| source | 4 | `STM` / `PHY` |
| | **74** | |

Forty-two columns for the title is comfortable: it holds *Return of the Obra Dinn* (24) and
*The Legend of Zelda: Breath of the Wild* (39) without truncation.

#### 2.2.1 · Why the glyph and focus fields are two columns, not one

**The four state glyphs are not one width class.** Verified against Unicode 16.0
`EastAsianWidth`:

| Glyph | Codepoint | Class |
|---|---|---|
| `○` not started | U+25CB | **Ambiguous** |
| `◐` in progress | U+25D0 | **Ambiguous** |
| `◉` zerado | U+25C9 | Neutral |
| `⊘` abandoned | U+2298 | Neutral |
| `▌` focus marker | U+258C | **Ambiguous** |

An Ambiguous-width character is **one cell in most terminals and two** where the user has set
`ambiguous-width=double` — common in CJK-configured terminals and in several popular emulators'
defaults. So on the most-used component in the product, two of the four states would render one
cell wider than the other two, and every column to their right would shear by one.

**This is fixed in rendering, not by changing the glyphs.** The glyphs are ratified and
CVD-verified; they are not the problem. The row reserves a **fixed two-column field** for the
glyph and a **fixed two-column field** for the focus marker, and pads to that width using a
width-aware measurement at render time. Every downstream column then starts at a fixed offset
whatever the terminal decides an Ambiguous glyph is worth.

The rendering rule, the measurement helper and the `ZERADO_ASCII=1` escape hatch belong to the
design system — [`../design/01-design-system.md`](../design/01-design-system.md) §1.2. The spine's
job here is only to say that the budget **reserves two columns**, and why.

> Found by `fft-design-architect` while writing deliverable B, and independently confirmed
> against `unicodedata` before this budget was corrected. The first draft of this document
> budgeted one column and would have sheared the state column on a real class of terminals.

### 2.3 · Layout budget — `Z-04` at ExtraWide 120×40

Body `112 × 32`, split **66 ∥ 2 ∥ 44**:

```
 ┌ list  66 ────────────────────────────┐ ┌ detail  44 ──────────────────┐
 │ row  1  pinned summary               │ │ title                        │
 │ row  2  (respiro)                    │ │ state · hours · source       │
 │ row  3  column header                │ │ (respiro)                    │
 │ rows 4–31  28 game rows              │ │ sinopse / not-fetched notice │
 │ row 32  scroll position              │ │ …                            │
 └──────────────────────────────────────┘ └──────────────────────────────┘
```

Twenty-eight visible rows. The list drops its `source` column here — the detail pane carries it,
and the title takes the space instead.

### 2.4 · The overlay budget — `Z-06 Set status`

Fixed `34 × 11`, centred on the current route, drawn over a **dimmed** backdrop — where dimming
is available. At the 16-colour floor and under `NO_COLOR` there is no dim, so the overlay is
separated by its **border** instead, in `--z-border-strong` (`#64748B`, 4.08:1 — it satisfies
WCAG 1.4.11 as a control boundary; `--z-border` at 1.53:1 may never do this job).

At Tiny the overlay does not fit inside `30 × 21` with its margins, so **at Tiny an overlay
becomes a route** — pushed, full-screen, popped with `Esc`. Behaviourally identical; only the
composition changes.

---

## 3 · Composition for later phases — inventoried, not specified

Enough to prove the navigation model does not need retrofitting.

| ID | Screen | Composition | Fits the model? |
|---|---|---|---|
| Z-12 | Enrichment sync | Single-pane readout — same shape as Z-03 | Yes, a second instance of one shape |
| Z-13 | Mood picker | ≤119 single-pane list → results route · ≥120 moods ∥ results | Yes, the same two-host pattern as Z-04/Z-05 |
| Z-14 | Mood tags | Overlay `34 × 13` | Yes |
| Z-15 | Cover deck | Mode of Z-04, grid instead of list | Yes — a mode swaps the body renderer, not the frame |
| Z-16 | Price history | Mode of Z-05 — a section inside the detail view | Yes, in both hosts |
| Z-17 | Command palette | Overlay, `min(60, w−8) × 12`, top-anchored | Yes; keys already reserved |
| Z-18 | Tonight | Single-pane, one card at a time | Yes |
| Z-19 | Budget | Single-pane form | Yes |
| Z-20 | Watchlist | Ledger — the same triad as Z-04 | Yes, third instance of the ledger shape |
| Z-21 | Account | Single-pane form | Yes |
| Z-22 | Devices and sync | Ledger + pinned summary | Yes |
| Z-23 | Public profile | Single-pane scroll | Yes |
| Z-24 | Comments and reviews | ≤119 single-pane thread · ≥120 list ∥ thread | Yes |
| Z-25 | Community lists | Ledger | Yes |
| Z-26 | Premium | Single-pane prose (`glamour`) | Yes |

**Every screen in the product is one of five shapes**, and that is the point of doing this now:

1. **Ledger** — a scrolling list with a pinned summary (Z-04, Z-20, Z-22, Z-25)
2. **Detail view** — a scrolling read surface with two possible hosts (Z-05, Z-16, Z-23)
3. **Form** — a `huh` field group (Z-02, Z-08, Z-09, Z-19, Z-21)
4. **Readout** — a centred progress-and-result surface (Z-03, Z-12, Z-18)
5. **Overlay** — a small centred modal, or a route at Tiny (Z-06, Z-14, Z-17)

Plus one deliberate outlier, `Z-11 Fatal error`, which is frameless by design.

Five shapes is what makes the design system finite, and it is what lets more than one worker
build Phase 1 in parallel without the screens drifting apart.
