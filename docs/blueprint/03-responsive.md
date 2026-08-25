---
title: Zerado — responsive composition
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-03
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: implementation-plan
ticket: "#2"
---

# Responsive composition — 40 / 60 / 80 / 120

What collapses, what hides, and what **refuses to render and says so**.

---

## 1 · The five tiers, and what each one means to Zerado

Breakpoints are the FlowForge canon's, unchanged: **Tiny `<40` · Narrow `40–59` · Standard
`60–79` · Wide `80–119` · ExtraWide `120+`** (TUI Design Manual R-7).

| Tier | Zerado's posture |
|---|---|
| **ExtraWide** 120+ | Two regions. The library shows the detail pane; the list shows more rows |
| **Wide** 80–119 | **The design target.** One region. Everything is correct here, by definition |
| **Standard** 60–79 | One region, secondary columns shed |
| **Narrow** 40–59 | One region, rows become **two lines** so the state label survives |
| **Tiny** <40 | One region, two-line rows, chrome stripped to a 1-row title |
| **below 24 × 8** | **Refuses.** One honest sentence, and it exits |

## 2 · The rule that governs every degrade

> **Shed columns before you shed meaning. Shed meaning before you shed the state.**

The co-render rule (colour **and** glyph **and** label) is what the state system is built on. A
degrade that drops the *label* leaves colour and glyph — and under `NO_COLOR` that is glyph
alone, which is unambiguous but not *readable* by someone who has not learned the legend.

So Zerado does not drop the label. **It changes the row's shape instead.**

### 2.1 · The game row at each tier

**Wide and ExtraWide — one line, full chip:**
```
▌ ◉ ZERADO       Return of the Obra Dinn                    9h   STM
  ◐ IN PROGRESS  Outer Wilds                               12h   STM
  ○ NOT STARTED  Hollow Knight                             41h   PHY
```

**Standard (60–79) — one line, `source` shed to the detail view:**
```
▌ ◉ ZERADO       Return of the Obra Dinn             9h
  ◐ IN PROGRESS  Outer Wilds                        12h
```

**Narrow and Tiny (<60) — two lines, and nothing is lost:**
```
▌ ◉ Return of the Obra Dinn
    ZERADO · 9h · Steam

  ◐ Outer Wilds
    IN PROGRESS · 12h · Steam
```

> **The glyph field is two columns at every tier**, including the two-line rows above — two of
> the four state glyphs are Ambiguous-width and would otherwise shear the column on terminals
> configured for double-width ambiguity. See [`02-composition.md`](./02-composition.md) §2.2.1.

The two-line row costs half the visible rows (36 × 16 gives **8 games** instead of 16) and buys
back the full label, the source, and an untruncated title on most games. That is the right
trade: at 40 columns a player is reading, not scanning.

## 3 · The per-screen degrade table

| ID | ExtraWide 120+ | Wide 80–119 | Standard 60–79 | Narrow 40–59 | Tiny <40 |
|---|---|---|---|---|---|
| **Z-01** First run | Choice list, prose at 68 cols | Same, prose wraps at body width | Same | Prose to 2 lines per door | Door **titles only**, no prose |
| **Z-02** Connect | Field + inline help beside it | Field, help **below** it | Same | Help collapses to `?` — press to expand | Same; one field visible at a time |
| **Z-03** Sync | Scanner + per-provider lines + running counts | Same | Counts only, no per-provider lines | Scanner + one count line | **Count line only** — the scanner is dropped, see §5 |
| **Z-04** Library | **List ∥ detail**, 28 rows | Single-pane, 12 rows, full chip | `source` shed | **Two-line rows**, 8 games | Two-line rows, title truncated to body width |
| **Z-05** Detail | **A pane** in Z-04 | A route, full width | Same | Same, labels above values instead of beside | Same |
| **Z-06** Set status | Overlay 34 × 11 | Same | Same | Same | **Becomes a route** — the overlay does not fit |
| **Z-07** Filter | Filter bar, **dynamic height** — §5b | Same | Same | Same; chips wrap sooner | One filter at a time, cycled with `Tab` |
| **Z-08** Add by hand | Field + help beside | Field, help below | Same | Help collapses | One field at a time |
| **Z-09** Settings | Single-pane grouped list, **wider gutter** | Single-pane grouped list | Same | Values below labels | Values below labels |
| **Z-10** Help | Two key columns | One key column | Same | Same | Same, scrolls more |
| **Z-11** Fatal | Plain text | Plain text | Plain text | Plain text | Plain text |

## 4 · What hides, and what must never hide

**May be hidden at narrower tiers:**
- the `source` column (Standard and below) — it is in the detail view;
- inline help text next to a form field (Narrow and below) — it collapses to an expandable `?`;
- the column-header row (Narrow and below) — a two-line row labels itself;
- the scanner sweep (Tiny only) — see §5;
- the second region (below 120) — the detail pane.

**Must never hide, at any width:**
- **the state** — glyph and label both, on every row, at every tier;
- **the game's title** — the identity column, R-10(a);
- **the pinned summary row** — R-10(c). It is *the* thing that must survive a 400-game library
  on a short terminal;
- **the focus marker** — the player must always be able to see where they are;
- **the degrade banner**, when one is active. If the product is showing stale or partial data,
  saying so outranks anything the space was going to be used for;
- **the footer key hint line** — at Tiny it shortens to the three keys that matter on that
  screen, but it does not vanish.

## 5 · The scanner sweep and the redraw budget

The scanner sweep is the brand's one signature motion: a `─` (U+2500) track carrying **three
consecutive `━` (U+2501) cells** as the pip, driven from elapsed time on the 2400 ms sinusoid
`cubic-bezier(0.45, 0, 0.55, 1)`, redrawn at 30 fps (brand manual §7.1).

Three rules the spine imposes on it:

1. **Indeterminate waits only.** When Zerado knows the count — "412 of 1,140" — it shows the
   count and a determinate bar. The sweep is for when the answer is honestly *"we don't know yet."*
2. **One at a time.** Never two sweeps on one screen. Never a sweep as ambient decoration.
3. **At Tiny it is dropped**, not shrunk. Below 40 columns a three-cell pip on a ~28-cell track
   is a blinking dash, which reads as a glitch rather than as hardware. The count line replaces
   it, which is more information in less space.

Under `prefers-reduced-motion` — or its terminal equivalent, which Zerado takes to be
`NO_COLOR` **or** an explicit `ZERADO_REDUCED_MOTION` setting — the pip **parks at the centre of
the track at full weight** and does not travel. It is deliberately not hidden: the lit slot is an
identity element, the travel is the decoration (brand manual §7.3).

## 5b · The filter bar's height is dynamic, not a per-tier constant

The first draft fixed the filter bar at **2 body rows at Standard and above**. `fft-tui-designer`
measured it and the constant does not survive: five active chips are **73–75 cells**, and the
**Standard** body is **54**. The bar was specified to a height it cannot render in.

**Rule: the filter bar takes the rows its active facets need, and the list absorbs the difference.**
That is already how it behaved at Narrow, so this makes the existing behaviour the rule instead of
a per-tier constant that happened to hold at two tiers.

| | |
|---|---|
| Minimum | **1** row — the editor, with no chips |
| Growth | chips wrap; each wrapped line is one more row |
| **Cap** | **4** rows. Beyond that the chips collapse to a count — `4 filters · ⏎ to edit` — and the bar returns to 2 |
| Floor it protects | the list keeps **at least 12** rows at the 80 × 24 body |

The cap is what stops a filter bar from eating the thing it filters. A player with enough facets
active to overflow four rows is better served by a count than by a wall of chips, and the editor is
one keystroke away.

*(Found by counting: the defect predates the `absent` facet and was surfaced by adding the fifth
chip. The designer recorded it rather than patching it, because the fix changes a tier's row map —
which is the spine's to change, not a screen's.)*

## 5c · `ZERADO_ASCII` covers every Ambiguous glyph, not just the state column

The state glyphs have an ASCII column (`[ ] [~] [*] [x]`) for terminals where they cannot be relied
on. But the state glyphs are **not the only Ambiguous-width characters on a Zerado screen** — the
box-drawing family, the focus marker `▌`, and the scanner's `─`/`━` are all Ambiguous too, and they
carry the frame itself.

**`ZERADO_ASCII=1` therefore switches the whole glyph vocabulary, not one column:**

| Role | Default | ASCII fallback |
|---|---|---|
| State | `○ ◐ ◉ ⊘` | `[ ] [~] [*] [x]` |
| Focus marker | `▌` | `>` |
| Box drawing | `┌─┐│└┘` | `+-+|+ +` |
| Scanner track / pip | `─` / `━` | `-` / `=` |
| Breadcrumb separator | `✦` | `>` |
| Audio annunciator | `▮` / `▯` | `[*]` / `[ ]` |

A fallback that rescues the state column while leaving the frame to shear would fix the smallest
problem on the screen and leave the largest.

*(Found by `fft-tui-designer`: the original scope covered only the state column.)*

## 6 · The refusal floor

Below **24 columns or 8 rows**, Zerado does not render a degraded interface. It prints one line
and exits with status `2`:

```
Zerado needs at least 24 columns and 8 rows. This terminal is 20 x 6.
```

Three reasons this is a refusal and not a degrade:

1. At 24 columns the body is 22 wide, and after the two-column focus field, the two-column glyph
   field and their space, the *title* — the one column R-10(a) makes mandatory — has **17
   characters**. `The Legend of Zel…`, `Return of the Ob…`, `Disco Elysium: T…`: every row becomes
   an ellipsis. A list of ellipses is not a smaller library view; it is a broken one.
2. A degrade that cannot show the state has abandoned the co-render rule, and the co-render rule
   is the product's accessibility mechanism, not a style.
3. Refusing is honest, actionable and one line. Rendering garbage is none of the three.

The floor is checked once at start and on every `WindowSizeMsg`. If a running terminal is
resized below it, Zerado replaces the screen with the same sentence and **keeps running** — it
does not exit mid-session, because the player is probably dragging a divider and will drag it
back. Exit is only for the start-up case.

## 7 · How this gets verified

Responsiveness is not provable from a single golden file — the whole point of TUI Design Manual
**R-10** and the `flowforge-live-repro-over-freeze` discipline. Every Phase 1 screen must be
rendered headless (`vhs` / `freeze`) at **all six** of these, and at an **overflowing row count**
where the screen is a ledger:

`24×8` (the refusal) · `32×24` (Tiny) · `40×24` (Narrow) · `60×24` (Standard) ·
**`80×24` (the floor)** · `120×40` (ExtraWide)

and again with `NO_COLOR=1` at `80×24`, and again at a forced 16-colour depth at `80×24`.

That is **eight artifacts per screen.** They are the founder-facing evidence in the
Screen-Quality Gate, and a screen without them is not GOLDEN.
