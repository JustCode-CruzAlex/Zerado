---
title: Zerado — Terminal Design System
discipline: DESIGN SYSTEM
doc-no: ZRD-DESIGN-02
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: design-system
---

# Zerado — Terminal Design System

The component vocabulary for Zerado's terminal. Every component below carries its anatomy,
its states, its spacing tokens **by name**, its colour roles **resolved to token + hex +
ANSI-256 + 16-colour**, its `NO_COLOR` rendering, and its behaviour at 40 columns.

Governed by `00-design-brief.md`. Colour spend is governed by `02-colour-budget.md`.
Composition — the layout tree, the pane budget, the focus state machine — belongs to
`fft-tui-architect` in the spine; this document supplies the parts that composition arranges.

---

## 1 · The laws every component obeys

### 1.1 · The cell grid is the substrate

Zerado has no glow, no gradients, no fluid type and no drop shadows (brand §9). Hierarchy comes
from **case, weight, colour role, box drawing and spacing — in that order.** A component that
needs a sixth channel is a component that is trying to be a web page.

### 1.2 · Glyph width — a verified hazard, and the fix

Checked against **`EastAsianWidth-17.0.0.txt`, dated 2025-07-24** and **`emoji-data.txt` v17.0,
dated 2025-07-25** (unicode.org), on 2026-08-25:

| Glyph | Codepoint | EAW class | Consequence |
|---|---|---|---|
| `○` not started | U+25CB | **A — Ambiguous** | may render **2 cells** wide |
| `◐` in progress | U+25D0 | **A — Ambiguous** | may render **2 cells** wide |
| `◉` zerado | U+25C9 | N — Neutral | always 1 cell |
| `⊘` abandoned | U+2298 | N — Neutral | always 1 cell |
| `✦` breadcrumb separator | U+2726 | N — Neutral | always 1 cell |
| `▸` palette sigil | U+25B8 | N — Neutral | always 1 cell |
| `▮` `▯` audio indicator | U+25AE · U+25AF | N — Neutral (**both**) | always 1 cell; **not** in the emoji set |
| `─` `━` scanner track / pip | U+2500 · U+2501 | A — Ambiguous (**both**) | consistent with each other |
| `│ ┌ ┐ └ ┘ ├ ┤` box drawing | U+2500 block | A — Ambiguous | the whole family is one class |
| `▌` annunciator | U+258C | A — Ambiguous | |
| `█` vs `░` | U+2588 · U+2591 | **A** vs **N** — mixed | **never pair these** |
| `[ ] [~] [*] [x]` ASCII column | ASCII | Na — Narrow | immune by construction |
| `♪` music note — **rejected** | U+266A | **A — Ambiguous** | would shear the status bar |
| `▪` `▫` small squares — **avoid** | U+25AA · U+25AB | N — Neutral, but **listed as `Emoji`** | may take emoji presentation and arrive coloured |

**The finding that matters:** the four ratified state glyphs are **not all the same width
class**. In a terminal configured to render ambiguous characters as double-width — the norm in
some CJK locales, and a setting many terminals expose — `○` and `◐` occupy two cells while
`◉` and `⊘` occupy one. The state column would misalign on the most-used component in the
product.

**This does not reopen the glyphs.** They are ratified and CVD-verified; they stay. The fix is
in rendering, and it is three rules:

1. **The state glyph sits in a fixed 2-column field.** The renderer measures the glyph's actual
   display width at runtime (an East-Asian-Width-aware width function, not `len()` and not
   `utf8.RuneCountInString`) and pads to exactly 2 columns. A double-width glyph fills the
   field and takes zero padding; a single-width glyph takes one space. **Column alignment is
   then invariant under either terminal setting.**
2. **Every rendered line is measured with the same width-aware function.** Never byte length,
   never rune count. This applies to truncation, padding and centring everywhere.
3. **`ZERADO_ASCII=1` forces the ratified ASCII column** `[ ] [~] [*] [x]`, which is entirely
   narrow and immune. It is also the automatic fallback when the terminal does not report
   Unicode capability.

**The declared requirement:** Zerado is designed for a terminal that treats East-Asian-Ambiguous
characters as **single-width** — the default outside CJK locales. Rule 1 keeps the ledger
correct when that assumption is violated; rule 3 is the escape hatch. Box drawing is Ambiguous
across the whole family, so it is internally consistent and carries no mixed-class risk.

> **Why `█` and `░` are forbidden as a pair.** They are the conventional progress-bar
> characters and they are **different width classes** — the bar changes length as it fills.
> Zerado's progress components use `━` and `─` instead: same class, and the same vocabulary as
> the scanner. See §9.

> **Why the audio indicator is `▮`/`▯` and not `♪`.** Two independent tests, and the note
> records what each one actually proved.
>
> **Width — decisive.** `♪` U+266A is **East-Asian Ambiguous**, so it would render two cells wide
> on exactly the terminals rule 1 above exists to protect, shearing the status bar. `▮` U+25AE
> and `▯` U+25AF are both **Neutral** — one cell, always.
>
> **Emoji — checked, and it does *not* hold for `♪`.** U+266A has **no entry in Unicode's
> `emoji-data.txt`** (its neighbours `♠ ♣ ♥ ♦ ♨` do; the music note does not), so it carries no
> emoji presentation and would not arrive pre-coloured under `NO_COLOR`. Recorded because it is
> the kind of plausible claim that ought to be checked rather than repeated.
>
> **The check did catch a real one, elsewhere.** `▪` U+25AA and `▫` U+25AB *are* listed as
> `Emoji` and `Extended_Pictographic`, so a font may render them in colour — which is precisely
> the failure that had emoji glyphs rejected in the first place. They are therefore on the avoid
> list above, and the Phase 3 price marker was moved off `▪` (§4.6).
>
> Beyond width, `▮`/`▯` are the right pair on their own merits: a **filled/hollow** contrast
> reuses the state system's own visual logic (`◉` against `○`), so the indicator reads as part of
> one family rather than as a borrowed icon, and it survives `NO_COLOR` unaided.

### 1.3 · Zerado does not paint a background — *design decision*

`--z-surface`, `--z-surface-raised` and `--z-surface-overlay` are **reference grounds for
measuring contrast**, not paint. Zerado renders on the user's own terminal background.

**Reasons.** (1) Brand §5.3 designed the palette to survive the user's background and measured
it against five popular themes — painting over them discards that work. (2) At the 16-colour
floor `--z-surface` and `--z-surface-raised` **both collapse to `black`**, so fill carries no
information there anyway. (3) A TUI that paints its own full-screen background looks like it is
squatting in someone's terminal, which is the opposite of an expensive object.

**Consequence, and it is binding:** **elevation is carried by borders and spacing, never by
fill.** No region may be separated from another by a background colour. The two sanctioned
exceptions — the confirmation overlay and the command palette — may paint
`--z-surface-overlay` to establish modality, but must **also** carry a
`--z-border-strong` border so they remain legible at 16 colours where the fill vanishes.

### 1.4 · The colour resolution table

The only colours a Zerado terminal component may use. Ratios are measured against
`--z-surface` and were read from the brand manual's table — none is estimated.

| Role | Token | Hex | ANSI-256 | 16-colour | Ratio |
|---|---|---|---|---|---|
| Primary text | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** AAA |
| Secondary text | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA |
| Tertiary text | `--z-text-tertiary` | `#8492A8` | *underived* | `white` | **6.15** AA |
| Amber — the ambient voice | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA |
| Cyan — earned | `--z-accent` | `#19E0FF` | **45** | `bright cyan` | **12.15** AAA |
| Focus ring | `--z-focus-ring` | `#19E0FF` | **45** | `bright cyan` | **12.15** AAA |
| Chrome / not started | `--z-state-not-started` | `#A5A29B` | **247** | `white` | **7.62** AA |
| In progress | `--z-state-in-progress` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA |
| Zerado | `--z-state-zerado` | `#19E0FF` | **45** | `bright cyan` | **12.15** AAA |
| Abandoned | `--z-state-abandoned` | `#C77DFF` | **177** | `bright magenta` | **7.21** AA |
| Decorative hairline | `--z-border` | `#2A3342` | **236** | `black` | 1.53 — **decoration only** |
| Control boundary | `--z-border-strong` | `#64748B` | **67** | `bright black` | **4.08** — meets 1.4.11 |
| Scanner / alarm | `--z-scanner` | `#FF2E2E` | **9** | `bright red` | 5.25 — **not text** |
| Error text | `--z-scanner-300` | `#FF6B6B` | *underived* | `bright red` | **6.99** AA |
| Scanner track | `--z-scanner-track` | `#5C1414` | *underived* | `black` | — |
| Unlit amber (inert track) | `--z-primary-muted` | `#8A5E00` | *underived* | `bright black` | — |

**`*underived*` means exactly that** — no ANSI-256 index has been derived for that token. Nobody
may pick one at the keyboard (brand §10 rule 5). Until the derivations land, each affected
component below specifies an **interim rendering that invents nothing**. The derivation list is
`00-design-brief.md` §9; the owner is `fft-brand-architect`.

### 1.5 · The type system, translated to a terminal

A terminal has one font — the user's. Orbitron, Space Grotesk and JetBrains Mono cannot render.
The three **roles** survive; the three **families** do not.

| Brand role | Family (web) | Terminal expression |
|---|---|---|
| **Display** | Orbitron | **UPPERCASE + bold + `--z-primary`.** Never faked with multi-row ASCII art — brand §3.5 forbids that for the mark, and it breaks on resize and in any log capture. |
| **Voice** | Space Grotesk | The default text: sentence case, `--z-text`. |
| **Native tongue** | JetBrains Mono | The substrate — the terminal already is this. |
| **Readout** (13 px, 0.18em tracking, UPPER) | JetBrains Mono | **UPPERCASE, `--z-text-secondary`** for block labels, `--z-primary` for a screen-level section head. **No letterspacing** — see below. |

> **Design decision: no letterspacing in the terminal.** The brand's readout style tracks out
> to 0.18em. A character grid has no sub-cell tracking, so "tracking" would mean inserting
> literal space characters — which changes the *text content*. Three costs: it doubles the width
> of every label at exactly the tiers where the canon is protecting body width; it fights the
> monospaced grid that brand §9 names as the terminal's primary hierarchy; and because a TTY has
> **no accessibility API**, anything reading the output stream — assistive technology, a pasted
> bug report, `grep` — receives `L I B R A R Y` instead of `LIBRARY`. Uppercase alone carries
> the readout role. *The founder may overrule this for a single hero label; it should not become
> a pattern.*

### 1.6 · The terminal mark

Per brand §3.5: `[0]` is the mark, `[0] ZERADO` the lockup, `0` the absolute floor. Drawn in
`--z-primary` (`214` / `bright yellow`). **No red anywhere in the terminal mark** — the web
mark's red slot does not translate, and §3.7 forbids recolouring the pip.

**Placement decision:** the mark does **not** repeat on every screen. Repeating a logo in a
three-row header band spends the exact rows the spacing canon fights for. The mark appears in
the first-run splash and in `zerado --help` as `[0] ZERADO`; at Tiny the status bar may carry
bare `0`. Every other screen is identified by the breadcrumb's `Zerado` segment, in text.

### 1.7 · Focus

One element holds focus at a time. Focus is carried by **three** channels so it survives both
`NO_COLOR` and the 16-colour floor — **and any two of them are enough**, which is what makes it
robust. The channels are fixed by the spine (`docs/blueprint/04-navigation-and-focus.md` §4.2):

| Channel | Focused | Not focused |
|---|---|---|
| **Position** | a `▌` (U+258C) marker in the left gutter — ASCII fallback `>` | two spaces |
| **Weight** | bold | normal |
| **Colour** | `--z-primary` amber (`214` / `bright yellow`) on the marker | none |

**The row cursor is amber, not cyan** — so the ledger cursor spends no cyan at all, and the
budget's rarest colour stays reserved for completion and the one primary action.

**For a focusable region** (a pane at ExtraWide): the focused pane's border is
`--z-border-strong` and its title amber; the unfocused pane's border is `--z-border` and its
title `--z-text-secondary`. Under `NO_COLOR` the focused pane uses a **heavier box-drawing
weight** (`┏━┓` against `┌─┐`) so the distinction survives with no colour at all.

**The `--z-focus-ring` cyan token still governs focused *controls*** — a text input, a form
field, a button — where the brand's ring is the indicator. It is exempt from the cyan budget
(`02-colour-budget.md` §2.3) because it is singular by definition.

**Never by background fill** (§1.3). **Never removed** in any state, on any element, for any
reason (brand §4.2, WCAG 2.4.7).

### 1.8 · Reuse-first, stated once

Zerado **cannot import FlowForge Go code** — the module does not resolve anonymously and
`LedgerTable` is `internal/` besides (`00-design-brief.md` §7). FlowForge canon is inherited as
*specification*. The live reuse target is the Charm ecosystem. Per-component verdicts are stated
in each section below, and where a shelf primitive genuinely does not fit, it says so.

---

## 2 · The screen header band

Adopted from the Spacing Canon (#2435 §5.1) with two Zerado adaptations.

### 2.1 · Anatomy — Wide tier (80–119 cols), no subtitle

```
row 1   (blank)                                    ← OuterMarginY = 1
row 2      Zerado ✦ Library                        ← breadcrumb, begins at col 4
row 3   (blank)                                    ← InterElementGap = 1
row 4      LIBRARY                                 ← title, begins at col 4
row 5   (blank)                                    ← InnerPaddingY + InterElementGap
row 6      ▌ ◉ ZERADO       Return of the Obra D…  ← body, begins at col 4
```

`leftInset = OuterMarginX (2) + InnerPaddingX (1) = 3`, so every visible row begins at
**column 4**. Header-left equals content-left. This is the whole point of #2435.

### 2.2 · The breadcrumb — two segments — *design decision*

**`Zerado ✦ <Screen>`.** Zerado has no project tier, and R-2 forbids faking an empty segment,
so the FlowForge three-segment form collapses to two. The `Zerado` segment uses the product
casing from `naming.md`; the screen segment is **sacred and never truncates** (#2435 §6). When
the strip overflows, the leading segment drops first.

**The separator is `✦` (U+2726 BLACK FOUR POINTED STAR).** Three reasons, in order:

1. **It is Neutral width (verified, §1.2)** — one cell in every terminal. The alternative `›`
   (U+203A) is also Neutral, so width does not decide it.
2. **`›` carries a false affordance in a terminal.** A chevron reads as a path or a shell
   prompt — something you could type or `cd` into. `Zerado › Library` invites the reading that
   `Library` is a subcommand. A four-pointed star reads as what it is: a dashboard annunciator
   between two labels.
3. **It is on-register and already proven.** A spark on a dark strip is the KITT cockpit idiom;
   it is also the separator FlowForge's own gold-standard screen uses, so Zerado inherits a
   glyph that has already survived real screens rather than reopening an unresolved question.

### 2.3 · The title row

The screen name, **UPPERCASE, bold, `--z-primary`** (`214` / `bright yellow`) — the display role
per §1.5. At Tiny the band collapses to this row alone.

### 2.4 · Subtitles — declined on every screen

**Composition decision, owned by the spine** (`docs/blueprint/02-composition.md` §1.2):
**Zerado screens carry no header subtitle.** `HeaderBandHeight` is always the **3-row base**,
and `hasSubtitle` is always `false`.

The two reasons, recorded so nobody re-proposes it per screen:

1. At the 80 × 24 floor the fourth row is worth more to the body than to the header — it is a
   whole extra game row.
2. Zerado's context is **always a live count** (`247 games · 6 zerado · filter: mood=grind`), and
   a live value belongs in the pinned summary row where it can change without redrawing the
   header band. In the header it would make the band the only mutable chrome on screen.

**A useful consequence:** because the flag has exactly one value, the canon's single-sizer
desync — the `hasSubtitle` disagreement that once clipped a FlowForge screen's body — **cannot
occur in Zerado by construction.**

Screens that need to explain themselves do it in the **body**, as a first block, not in the band.

### 2.5 · Spacing, per tier

| Token | Tiny | Narrow | Standard | Wide | ExtraWide |
|---|---|---|---|---|---|
| `OuterMarginX` | 0 | 1 | 2 | 2 | 2 |
| `OuterMarginY` | 0 | 1 | 1 | 1 | 1 |
| `InnerPaddingX` | 1 | 1 | 1 | 1 | 2 |
| `InnerPaddingY` | 0 | 1 | 1 | 1 | 1 |
| `InterElementGap` | 1 | 1 | 1 | 1 | 1 |
| `HeaderBandHeight` | **1** | 3 | 3 | 3 | 3 |
| `leftInset` | **1** | 2 | 3 | 3 | 4 |

### 2.6 · `NO_COLOR`

Breadcrumb and title render with **zero SGR sequences**. Bold survives as the non-colour
emphasis channel; the title stays UPPERCASE, so the hierarchy holds unaided.

### 2.7 · At 40 columns

Narrow tier: band stays 3 rows, `leftInset` 2, the subtitle **yields** even on screens that
carry one. Below 40 (Tiny): band collapses to the **title row only**, `OuterMarginX` drops to 0
so the body keeps every usable column.

### 2.8 · Reuse verdict

Build fresh against the #2435 numbers. `lipgloss` supplies the styling and alignment maths; the
band itself is ~40 lines and must live behind a single `Frame` wrapper enforced at the router,
so **a screen cannot render frameless by construction**.

---

## 3 · State chip

The most-used component in the product. Every row of every list.

### 3.1 · Anatomy — full chip, 14 columns

```
◉  ZERADO
│  │
│  └ label field — 11 cols, left-aligned, space-padded
└ glyph field — fixed 2 cols (glyph + runtime pad, §1.2)
```

`2 (glyph field) + 1 (gap) + 11 (label field) = 14 columns`, fixed. The label field is 11
because `NOT STARTED` and `IN PROGRESS` are the longest at 11 characters; the field never
resizes, so state columns align across every row and every screen.

### 3.2 · The four states — ratified, not open

| State | Token | Hex | ANSI-256 | 16-col | Glyph | ASCII | Label | Ratio |
|---|---|---|---|---|---|---|---|---|
| Not started | `--z-state-not-started` | `#A5A29B` | **247** | `white` | `○` U+25CB | `[ ]` | `NOT STARTED` | **7.62** |
| In progress | `--z-state-in-progress` | `#FFB000` | **214** | `bright yellow` | `◐` U+25D0 | `[~]` | `IN PROGRESS` | **10.59** |
| Zerado | `--z-state-zerado` | `#19E0FF` | **45** | `bright cyan` | `◉` U+25C9 | `[*]` | `ZERADO` | **12.15** |
| Abandoned | `--z-state-abandoned` | `#C77DFF` | **177** | `bright magenta` | `⊘` U+2298 | `[x]` | `ABANDONED` | **7.21** |

**The co-render rule: colour AND glyph AND label — all three, every chip, every tier.** Remove
any one and the state is still unambiguous. This is what makes the system survive `NO_COLOR`, a
monochrome terminal, a screenshot in a bug report, and colour-vision deficiency by one mechanism.

**The warm grey `#A5A29B` is load-bearing engineering.** The first draft's blue-cast
`#9FB0C6` collapsed against the cyan at **ΔE 8.8 under deuteranopia**; the warm grey measures
**25.8**. It must never be "corrected" back toward blue. The tightest surviving pair is
**zerado × abandoned at ΔE 11.9 under deuteranopia** — the floor to protect, and the one place
where glyph and label genuinely carry load rather than merely reinforce.

### 3.3 · Casing — where `ZERADO` and `zerado` each belong

Per `naming.md`, and this trips people up:

| Context | Form | Example |
|---|---|---|
| The chip, the filter, a column header | **`ZERADO`** | `◉  ZERADO` |
| A summary sentence of counts | **`zerado`** | `247 games · 6 zerado · 12 in progress` |
| The command | `zerado` | `zerado sync` |
| The product in prose | `Zerado` | "Zerado reads your Steam library." |

Italic-on-first-use is a prose convention and does **not** apply in the terminal — italics are
not reliably supported and the interface is not prose.

### 3.4 · States of the chip itself

| Chip state | Rendering |
|---|---|
| Default | glyph + label in the state colour |
| On the focused row | inherits the row's bold; the `▌` gutter marker is `--z-primary` amber (§1.7) |
| In a filter, unselected | glyph + label in `--z-text-tertiary` |
| In a filter, selected | glyph + label in the state colour, bold |
| Pending write | the label is suffixed `…` until the write confirms — never an optimistic silent change |

### 3.5 · `NO_COLOR`

Zero SGR. The chip is glyph + label, which is already unambiguous:

```
○  NOT STARTED
◐  IN PROGRESS
◉  ZERADO
⊘  ABANDONED
```

### 3.6 · At 40 columns and below

The chip does **not** shrink and the label is **never** dropped. Below 60 columns the *row*
changes shape instead — it becomes two lines, and the chip moves to the second line where it has
room (§4.3). That is the spine's answer, and it is a better one than compressing the chip: at 40
columns a player is reading, not scanning.

> **The alternative I considered and rejected.** A compact chip — glyph plus colour only, with
> the labels carried once by a legend in the status bar — buys back 12 columns per row. It was
> rejected because the ticket states the co-render rule as *"colour AND glyph AND label, all
> three, every one"*, and because a TTY has no accessibility API: the label beside the glyph
> **is** the text alternative (`00-design-brief.md` §3.2, SC 1.1.1 / 4.1.2). Dropping it for
> density would trade the one channel that survives every degrade. Recorded here so it is not
> re-proposed as an optimisation.

### 3.7 · Reuse verdict

Build fresh — it is a `lipgloss` style plus the width-aware pad from §1.2. No `bubbles`
primitive fits, and none should be forced.

---

## 4 · Game row — the ledger row

The library deck's row. Obeys R-10(a): it carries a **human game title**, never an index.

The row's shape per tier is a **composition decision owned by the spine**
(`docs/blueprint/03-responsive.md` §2.1). This section specifies the parts; the spine decides
the arrangement.

### 4.1 · Anatomy — one line, Wide and ExtraWide

```
▌ ◉ ZERADO       Return of the Obra Dinn                    9h   STM
│ │ │            │                                          │    │
│ │ │            │                                          │    └ source — 3 cols
│ │ │            │                                          └ hours — right-aligned
│ │ │            └ title — flex, the identity column (R-10a)
│ │ └ label field — 11 cols
│ └ glyph field — 2 cols, runtime-padded (§1.2)
└ focus gutter — ▌ when focused, blank otherwise (§1.7)
```

At **Standard (60–79)** the `source` column is shed to the detail view; everything else holds.

### 4.2 · Anatomy — two lines, Narrow and Tiny (< 60)

```
▌ ◉ Return of the Obra Dinn
    ZERADO · 9h · Steam

  ◐ Outer Wilds
    IN PROGRESS · 12h · Steam
```

**Nothing is lost.** Line one carries the focus gutter, the state glyph and an untruncated title;
line two carries the full state label, the hours and the source. The cost is half the visible
rows — `36 × 16` shows **8 games** instead of 16 — and it is the right trade, because at 40
columns a player is reading rather than scanning.

**This is what keeps co-render intact at every tier.** The label is never dropped for width; the
row grows a line instead.

### 4.3 · The field budget — one-line form

| Field | Cols | Note |
|---|---|---|
| focus gutter | 1 | `▌` U+258C — Ambiguous width, ASCII fallback `>` |
| gap | 1 | |
| state glyph field | 2 | fixed, runtime-padded |
| gap | 1 | |
| state label field | 11 | fixed — `NOT STARTED` and `IN PROGRESS` are the longest |
| gap | 2 | |
| **title** | **flex** | **the identity column — R-10(a)** |
| gap | 2 | |
| hours | 5 | right-aligned (`   9h` … `1204h`) |
| gap | 3 | Wide and ExtraWide only |
| source | 3 | Wide and ExtraWide only — `STM` · `PHY` |

Fixed left = **18 cols**. Right block (Standard+) = **7 cols**.

### 4.4 · Title width, computed per tier

| Tier | Width | `leftInset` | Body | Hours? | **Title** |
|---|---|---|---|---|---|
| ExtraWide | 120 | 4 | 112 | one line | **81** |
| Wide | 80 | 3 | 74 | one line | **43** |
| Standard | 60 | 3 | 54 | one line, no source | **29** |
| Narrow | 40 | 2 | 36 | **two lines** | **32** on line 1 |
| Tiny | 32 | 1 | 30 | **two lines** | **26** on line 1 |

At Wide and ExtraWide the source column and its gap cost 6 columns; at Standard they are shed.
Below 60 the two-line form spends only 4 columns before the title — gutter, gap and the glyph
field — so the narrow tiers get a **longer** usable title than Standard does. That is the
two-line row paying for itself.

### 4.5 · The refusal floor — 24 columns or 8 rows

**Owned by the spine** (`docs/blueprint/03-responsive.md` §6). Below **24 columns or 8 rows**
Zerado does not render a degraded interface. It prints one line and exits with status `2`:

```
Zerado needs at least 24 columns and 8 rows. This terminal is 20 x 6.
```

Three reasons it is a refusal and not a degrade: at 24 columns the title is down to about 16
characters, so every row would be an ellipsis, and a list of ellipses is not a smaller library
view but a broken one; a degrade that cannot show the state has abandoned the co-render rule,
which is the product's accessibility mechanism and not a style; and refusing is honest,
actionable and one line, while rendering garbage is none of the three.

Checked at start and on every `WindowSizeMsg`. **A running session resized below the floor
replaces the screen with the same sentence and keeps running** — it does not exit mid-session,
because the player is probably dragging a divider and will drag it back. Exit is for start-up
only.

### 4.6 · Row states

| State | Rendering |
|---|---|
| Default | gutter blank; chip in state colour; title `--z-text`; hours `--z-text-secondary` |
| Focused | gutter `▌` in `--z-primary` + bold row (§1.7) |
| Selected (multi-select, Phase 2) | gutter `▌`; title bold — **never a background fill** (§1.3) |
| Physical copy | title suffixed ` · physical` in `--z-text-tertiary`. A hand-added disc is **not** a second-class row (public copy §06). |
| Price-flagged (Phase 3) | a `▬` (U+25AC, Neutral, non-emoji) marker in `--z-primary` after the hours; never red — a good price is not an alarm. **Not `▪`** — see §1.2 |
| Loading placeholder | title renders `—` in `--z-text-tertiary`; the chip renders its real state. **Never a spinner per row.** |

### 4.7 · `NO_COLOR`

Zero SGR. This is the brand manual's own ratified rendering (§5.4):

```
  ○  NOT STARTED   Hollow Knight              41h
  ◐  IN PROGRESS   Outer Wilds                12h
  ◉  ZERADO        Return of the Obra Dinn     9h
  ⊘  ABANDONED     Sekiro                      3h
```

Focus is carried by the `▌` gutter and bold.

### 4.8 · The ledger triad — non-negotiable (R-10)

- **(a)** The title column is populated on every row with a human title.
- **(b)** The viewport follows the cursor; the selection is **always** visible, and cursor plus
  scroll offset **survive a row-set rebuild** so a background re-sync cannot snap the view to
  the top.
- **(c)** The summary row is **pinned outside the scroll region** (§5) and is on screen at any
  row count.

**Proven at 400 rows, not at a 12-row golden.** A frozen golden asserts one render at one size;
it never exercises "scroll to row 380 and check the selection is visible."

### 4.9 · Reuse verdict

**`bubbles/table` does not fit as-is** — it is the primitive that, hand-rolled per screen inside
FlowForge, independently dropped the title, the scroll and the pinned footer on two different
screens. FlowForge's answer, `LedgerTable`, is **not importable** (§1.8).

**Build a Zerado ledger primitive, once, correct-by-construction on (a)(b)(c)**, composed from
`bubbles/viewport` for the scroll region plus `lipgloss` for row styling. Every list screen uses
it; no screen hand-rolls a table. This is reuse-first at the only level available: the pattern
is inherited even though the package cannot be.

---

## 5 · Status bar — the pinned summary

R-10(c). One row, pinned outside the scroll region, always on screen.

### 5.1 · Anatomy — Wide tier, body width 74

```
247 games · 6 zerado · 12 in progress                  ▮ AUDIO   ? help
│                                                      │          │
└ the summary — prose casing (§3.3)                     │          └ help hint, right
                                    audio indicator ────┘  present only when audio is enabled
```

### 5.2 · The audio indicator

Audio is **off until the player opts in** (§15), so the default status bar carries **no audio
indicator at all** and **`m` is absent from the footer** — there is nothing to mute. The
indicator appears only once audio has been enabled.

Co-rendered on three channels, exactly like every other state in the product:

| Audio state | Glyph | Label | Colour role | ANSI-256 / 16-col |
|---|---|---|---|---|
| **Never enabled — the default** | — | — | **no indicator; `m` absent from the footer** | — |
| Enabled, unmuted | `▮` U+25AE | `AUDIO` | `--z-primary` | **214** / `bright yellow` |
| Enabled, muted | `▯` U+25AF | `MUTED` | `--z-text-secondary` | **249** / `white` |

The filled/hollow pair carries the state without colour, so the indicator survives `NO_COLOR`
and the 16-colour floor unaided — the same mechanism as `◉` against `○`. Glyph choice and the
rejected `♪` are recorded in §1.2.

The amber here is an **ambient-voice** use, not an action: it is the machine saying it is on, not
asking to be pressed. It therefore spends no cyan and does not touch the chrome-cyan budget
(`02-colour-budget.md` §4.1).

### 5.3 · Content rules

- **Say the number.** `247 games`, never `a lot of games` (brand §8).
- Prose casing for counts: `6 zerado`, not `6 ZERADO`.
- Three facts maximum. The ratified voice example is the shape: *"247 games. 6 finished. Last
  played: 3 weeks ago."* — three facts, no adornment, no exclamation mark.
- The right-hand hint is **always** `? help` (WCAG 3.2.6 Consistent Help).
- `m` appears in the footer **only when audio is enabled.** A key hint for a subsystem the player
  never turned on is noise, and it advertises a feature as if it were already running.

### 5.4 · Spacing

Occupies the **reserved footer row** of the Spacing Canon frame (#2435 §5.2) — it is not an
extra row and does not steal from `BodyRect`.

### 5.5 · States

| State | Rendering |
|---|---|
| Default | counts in `--z-text-secondary` (`249` / `white`); numerals in `--z-text` |
| Focused row has hours, tier < Standard | replaced by `Hollow Knight · 41h` — where the hours went (§4.4) |
| Syncing | replaced by the progress readout (§9) |
| Offline | prefixed by the degrade banner (§12) |
| Filter active | replaced by the filter bar (§7) |
| Audio enabled | the indicator sits between the summary and the help hint; `m` joins the footer keys |
| Audio muted | indicator becomes `▯ MUTED` in `--z-text-secondary`; nothing else on the row changes |
| Tiny | counts only: `247 · 6 · 12`, and `?` alone as the hint. The audio indicator degrades to the bare glyph `▮` / `▯` — **the glyph is the last thing dropped, never the first** |

### 5.6 · `NO_COLOR` / 40 columns

Zero SGR; the text carries itself. At 40 columns the summary truncates from the right, dropping
whole facts rather than mid-word: `247 games · 6 zerado`. The `? help` hint is the last thing to
go, and at Tiny becomes `?`.

### 5.7 · Reuse verdict

Build fresh — a `lipgloss` join with the width-aware truncation from §1.2.

---

## 6 · Detail pane

One game, expanded. Verdict 4 set the floor for a detail *pane* at 80 columns; the spine
resolved it at **120**, on the same R-10(a) arithmetic. Below 120 the detail is a **route**.

### 6.1 · Presentation, by tier

**Owned by the spine** (`docs/blueprint/02-composition.md` §2.1). The detail is **one view with
two hosts**:

| Tier | Library composition | Where the detail lives |
|---|---|---|
| **ExtraWide `120+`** | list ∥ detail — body 112 → ledger 64 · gutter 2 · pane 46 | a **pane**; `Enter` moves focus into it |
| Tiny · Narrow · Standard · **Wide** | single-pane list | a **route**, reached with `Enter`, left with `Esc` |

**Why not at 80.** Splitting a 74-column body gives roughly 44 for the list and 28 for the pane,
which drops the identity column to about 15 columns after the state chip and the playtime figure
— `Return of the Ob…`, a plain R-10(a) failure — and a 28-column pane cannot hold a *sinopse*
either. **Two regions only when there is room for both to be correct.**

**Consequence for this spec:** the detail component is written **host-agnostic**, built once and
mounted twice. Nothing in it may assume a border, a pane width, or a surrounding route.

### 6.2 · Anatomy — Wide tier, pane 28 cols

```
┌──────────────────────────┐   ← --z-border-strong (67 / bright black)
│  Return of the Obra Dinn │   ← title, --z-text, wrapped, ≤ 24 cols
│                          │
│  ◉  ZERADO               │   ← the full 14-col chip (§3)
│                          │
│  PLAYTIME      9h        │   ← readout label + value
│  ADDED         Mar 2026  │
│  SOURCE        Steam     │
│                          │
│  MOOD                    │   ← Phase 2
│  Story rich, kind of sad │
└──────────────────────────┘
```

Readout labels: UPPERCASE, `--z-text-secondary`, no letterspacing (§1.5). Values: `--z-text`.

### 6.3 · Spacing

`InnerPaddingX` inside the border on each side; `InterElementGap` (1 row) between blocks; a
2-column gutter between ledger and pane. Content never touches the border; the border never
touches the screen edge (#2435 §2, Benchmark A).

### 6.4 · States

| State | Rendering |
|---|---|
| Loading metadata | values render `—`; the state chip is already known and shows immediately |
| No metadata yet (pre-Phase 2) | `MOOD` and `SINOPSE` blocks are **absent**, not empty-labelled. Never show a field for a capability that does not exist. |
| Metadata unavailable | `Not available offline.` in `--z-text-tertiary` |
| Physical copy | `SOURCE   Added by hand` |

### 6.5 · Accessibility

The pane is focus-triggered content: it is **dismissible** with `esc`, it **persists** while
focus remains, and it **may never entirely obscure the focused row** (WCAG 2.4.11) — which is
why it is a side pane and not a centred modal. Because reading order in a terminal is byte order
(SC 1.3.2), the pane's content is a **single column**, top to bottom.

### 6.6 · `NO_COLOR` / 40 columns

Zero SGR; box drawing and spacing carry structure. At 40 columns the pane does not exist — the
detail is a separate screen with the same content in one column.

### 6.7 · Reuse verdict

`bubbles/viewport` for the scroll body when content overflows; `lipgloss` for the border and
layout. In Phase 2, `glamour` renders the *sinopse* — but with a Zerado style, not glamour's
default theme, which carries its own palette and would break the colour budget.

---

## 7 · Filter bar

A **mode of the list**, not a permanent region (inherited verdict 4).

### 7.1 · Presentation, by tier

| Tier | Behaviour |
|---|---|
| ExtraWide `120+` | May render persistently as one line below the header band |
| `< 120` | A **mode**: `/` activates it and it **replaces the status bar row**; `esc` clears and returns |

### 7.2 · Anatomy — mode, body width 74

```
/ obra▎                                        12 of 247      esc clear
│ │   │                                        │              │
│ │   └ cursor                                 │              └ exit hint — always present
│ └ the query, --z-text                        └ live result count
└ mode sigil, --z-primary
```

### 7.3 · State filter chips (ExtraWide, or via `f`)

```
[ ○ NOT STARTED ]  [ ◐ IN PROGRESS ]  [ ◉ ZERADO ]  [ ⊘ ABANDONED ]
```

Unselected: `--z-text-tertiary`. Selected: the state colour, bold. Co-render holds — the chip
keeps colour **and** glyph **and** label even as a filter control. Chip boundaries use
`--z-border-strong` (4.08, meets 1.4.11), never `--z-border` (1.53).

### 7.4 · Accessibility — two rules that are not optional

1. **`esc` always exits the mode.** A filter with no exit is a literal keyboard trap
   (WCAG 2.1.2).
2. **While the query field holds focus, single-key shortcuts do not fire** (WCAG 2.1.4).
   Typing `d` types `d`. This is the single most common way a TUI destroys a user's data.

### 7.5 · Empty result

Not an error. See §10(b).

### 7.6 · `NO_COLOR` / 40 columns

Zero SGR; the `/` sigil and the `esc clear` hint carry the mode in words. At 40 columns the
count shortens to `12/247` and the hint to `esc`; the query field keeps the rest.

### 7.7 · Reuse verdict

**`bubbles/textinput`** — a direct fit, reuse it. The chips are `lipgloss`.

---

## 8 · Progress readout — determinate

For work with a known denominator: syncing 147 of 247 games.

### 8.1 · Anatomy

```
SYNCING STEAM                                          147 / 247
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━────────────────────────────────────
│                             │
└ filled — --z-primary        └ unfilled — --z-primary-muted
```

- **Filled:** `━` U+2501 in `--z-primary` (`214` / `bright yellow`).
- **Unfilled:** `─` U+2500 in `--z-primary-muted` `#8A5E00` — the token is literally documented
  as *"unlit amber — inert tracks"*.
- Label UPPERCASE `--z-text-secondary`; the count `--z-text`. **The number is always shown** —
  the bar is the ornament, the count is the information (brand §8: *say the number*).

**Why `━`/`─` and not `█`/`░`:** verified in §1.2, `█` is Ambiguous and `░` is Neutral — a
mixed-width pair whose rendered length changes as it fills. `━` and `─` are the same class, and
they are already the scanner's vocabulary.

### 8.2 · Interim rendering — the ANSI index is not derived

`--z-primary-muted` `#8A5E00` has **no derived ANSI-256 index**. Until it is derived, the
unfilled track renders **uncoloured** — no SGR at all. That invents nothing, and it degrades
identically to the `NO_COLOR` path. **Do not pick an index by eye.**

### 8.3 · States

| State | Rendering |
|---|---|
| Running | bar advances; count updates |
| Stalled > 10 s with no progress | switches to the scanner (§9) and the label becomes `WAITING ON STEAM` — honest about not knowing |
| Complete | bar full for one beat, then the readout is replaced by the result line: `247 games. 6 zerado.` |
| Partial failure | bar stops at the reached value; an error line is added below (§11). **The bar never jumps to 100 % on failure.** |

### 8.4 · `NO_COLOR` / 40 columns

Zero SGR — the heavier stroke `━` against the lighter `─` carries fill in one ink. At 40 columns
the label and count stack above a full-width bar:

```
SYNCING STEAM      147/247
━━━━━━━━━━━━━━────────────
```

### 8.5 · Reuse verdict

**`bubbles/progress` does not fit.** Its default ramp is a truecolor gradient, which spends
colour Zerado's budget does not have and cannot resolve at the 16-colour floor. Build fresh — a
progress bar is a division and a repeat, and the width-aware fill from §1.2.

---

## 9 · The scanner — indeterminate progress

The brand's **one** signature motion (§7.1). KITT's oscillating mirror, one row tall.

### 9.1 · Anatomy

```
WAITING ON STEAM
─────────━━━──────────────────────────────
         │
         └ pip — exactly 3 cells, --z-scanner
```

- **Track:** `─` U+2500 in `--z-scanner-track` `#5C1414`.
- **Pip:** exactly **three consecutive** `━` U+2501 in `--z-scanner` `#FF2E2E` → ANSI-256 **9**,
  16-colour `bright red`.

### 9.2 · The motion — ratified, reproduce exactly

| Property | Value |
|---|---|
| Full cycle | **2400 ms** (1200 ms each direction), `alternate`, `infinite` |
| Easing | `cubic-bezier(0.45, 0, 0.55, 1)` — the sinusoid |
| Redraw | **30 fps** |
| Travel | pip left position from `0` to `W − 3`, where `W` is track width in cells |

```
p        = (elapsed_ms mod 2400) / 2400          ∈ [0, 1)
u        = 2p            if p < 0.5              ← the "alternate" fold
           2 − 2p        otherwise
e        = bezier_y(0.45, 0, 0.55, 1)(u)
pipLeft  = round(e × (W − 3))
```

**Why that easing:** KITT's scanner was a physically oscillating mirror, and a real oscillator
is slowest at the extremes and fastest through the middle. Linear travel reads as a loading bar;
a bounce reads as a cartoon. The sinusoid is the only one that reads as *hardware*.

### 9.3 · When it may appear — a closed rule

**Only for genuinely indeterminate waits.** Never ambient, never decorative, never as a section
divider in the terminal. Ambient animation burns a redraw budget for nothing and is exactly the
nostalgia-kitsch brand §1 rules out. If the denominator is known, use §8 instead.

Two further rules from the spine: **one at a time** — never two sweeps on one screen; and **at
Tiny the sweep is dropped, not shrunk.** Below 40 columns a three-cell pip on a ~28-cell track is
a blinking dash that reads as a glitch rather than as hardware, so the count line replaces it —
more information in less space.

### 9.4 · Interim rendering

`--z-scanner-track` `#5C1414` has **no derived ANSI-256 index**. Until it does, the track
renders **uncoloured**; the pip renders at index `9`. The pip still reads, and nothing is
invented.

### 9.5 · Reduced motion

The terminal's reduced-motion signal is **`NO_COLOR` or an explicit `ZERADO_REDUCED_MOTION`**
(spine, `docs/blueprint/03-responsive.md` §5). Under either, the pip **parks at the centre of its
track at full weight** and the ticker stops — `pipLeft = round((W − 3) / 2)`. It is deliberately
**not hidden**: the lit slot is an identity element, the travel is the decoration (brand §7.3).

### 9.6 · Performance

The 30 fps ticker runs **only while an indeterminate wait is on screen** and is torn down with
it — no leaked goroutine, no timer surviving the view. A tick must never block input, and each
redraw must complete inside the 60 fps app budget.

### 9.7 · `NO_COLOR` / 40 columns

Zero SGR — the **heavier stroke is the pip**, so the motion survives in one ink. This is why the
primitive uses stroke weight rather than colour as its primary channel. At 40 columns the track
spans the body width; the pip stays exactly 3 cells.

### 9.8 · Reuse verdict — *do not use `harmonica`*

The TUI manual assigns `harmonica` the spring-animation role, and it is the wrong tool here:
harmonica models a **damped spring converging on a target**, while the scanner is an
**undamped, infinite, alternating sinusoid**. Evaluate the bezier directly. Stated plainly
because reaching for the shelf primitive would be a defensible-looking mistake.

---

## 10 · Empty state

**Three different empties.** Collapsing them into one message is how a product tells a player
"no games found" when the real answer is "your Steam profile is private."

### 10.1 · (a) First run — nothing connected yet

```
      No library yet.

      Zerado reads your Steam library once you add a key.
      Physical discs and cartridges can be added by hand.

      s  connect Steam        a  add a game by hand
```

Body `--z-text`; the key hints `--z-primary`. **This is the screen nobody writes down and the
one every new player sees first.**

### 10.2 · (b) Filter matched nothing

```
      No games match this filter.
      247 in the library.

      esc  clear the filter
```

Not an error, not red. Says the number.

### 10.3 · (c) The provider returned nothing — ratified copy, use verbatim

```
      Steam returned an empty library.

      Game details are private on your profile — Steam won't
      share the list until that's public.
      Settings → Privacy.

      r  try again        s  Steam settings
```

Straight from brand §8 and consistent with the public FAQ. **Name what happened, why, and the
next action.** "Something went wrong" is the one sentence a terminal user cannot act on.

### 10.4 · Spacing / `NO_COLOR` / 40 columns

Top-aligned in `BodyRect`, left-aligned text starting at `leftInset`, `InterElementGap`
between blocks. `NO_COLOR`: zero SGR — the copy carries itself entirely. At 40 columns the
copy re-wraps and the key hints stack one per line. **The copy is never truncated** — an empty
state that gets cut off has failed at its only job.

### 10.5 · Reuse verdict

Build fresh — `lipgloss` text layout. Trivial, and worth doing well.

---

## 11 · Error state

### 11.1 · Anatomy

```
▌ SYNC FAILED
  Steam did not answer in 30 seconds.
  Your library is unchanged — nothing was lost.

  r  try again        esc  back
```

- Annunciator `▌` U+258C in `--z-scanner` (`9` / `bright red`) — **structure, not text**, so the
  5.25 ratio is used where the brand permits it.
- Heading `SYNC FAILED`: UPPERCASE, **bold**.
- Body `--z-text`.

### 11.2 · Interim rendering — error text stays uncoloured

`--z-scanner-300` `#FF6B6B` (**6.99**, AA — the readable red) has **no derived ANSI-256 index**.
Until it does, error **text** renders uncoloured and bold, and the red lives only in the `▌`
annunciator at index `9`. This keeps the brand's own rule intact — `--z-scanner` is
motion and alarm, **not text** — and invents nothing.

### 11.3 · Content rules

Every error names three things: **what happened**, **why**, **what to do next** (WCAG 3.3.1,
3.3.3). And a fourth Zerado adds because the product's promise is a local file: **what happened
to the player's data.** `Your library is unchanged — nothing was lost.` is not reassurance, it
is information.

### 11.4 · `NO_COLOR` / 40 columns

Zero SGR — `▌` plus the uppercase heading carry the alarm in one ink. Co-render holds because
the **word** `FAILED` is present; colour is only the third channel. At 40 columns the copy
re-wraps and hints stack; the heading never truncates.

### 11.5 · Reuse verdict

Build fresh. **Do not use `charmbracelet/log` for user-facing errors** — its role is structured
developer logging (#2371 §2), and its level colours are its own palette, not Zerado's.

---

## 12 · Degrade banner — the honest-offline strip

The product publicly promises it *"runs with the network off"*. Offline is therefore a **normal
state**, not a failure — and it must be **shown, never silent**.

### 12.1 · Anatomy — informational

```
▌ OFFLINE   Last synced 3 hours ago. Everything here still works.
```

`▌` in `--z-border-strong` (`67` / `bright black`); text `--z-text-secondary` (`249` / `white`).
**Chrome — not red, not amber.** Red would call a promised behaviour a fault.

### 12.2 · Anatomy — action required

```
▌ STEAM KEY MISSING   Press s to add it.
```

`▌` in `--z-primary` (`214` / `bright yellow`); text `--z-text` (`255`). **Amber appears only
when the player must do something.** That is the whole distinction, and it is the rule.

### 12.3 · The degrade matrix

| Condition | Class | Copy |
|---|---|---|
| Network off, library synced | informational | `OFFLINE   Last synced 3 hours ago. Everything here still works.` |
| Network off, never synced | action | `OFFLINE   Nothing synced yet. Connect when you're back online.` |
| Price data unreachable (Phase 3) | informational | `PRICES OFFLINE   Showing the last prices Zerado saw.` |
| Metadata unreachable (Phase 2) | informational | `COVERS OFFLINE   Titles and states are all local.` |
| Steam key missing | action | `STEAM KEY MISSING   Press s to add it.` |
| Steam profile private | action | `STEAM PROFILE PRIVATE   Steam won't share the list until it's public.` |

**Never invent a degrade that hides a promise.** If something the landing page promises is not
working, the banner says so in words.

### 12.4 · Co-render

The **label word is always present** — `OFFLINE`, `PRICES OFFLINE` — so the state survives with
zero colour and needs no glyph beyond the structural `▌`. Colour is the third channel only, and
it distinguishes only *informational* from *action*.

### 12.5 · Placement, `NO_COLOR`, 40 columns

Pinned directly above the status bar, inside `BodyRect`, one row, never scrolling away. Zero SGR
under `NO_COLOR`. At 40 columns the label survives and the sentence truncates from the right at
a word boundary; at Tiny only the label renders — `▌ OFFLINE`.

### 12.6 · Reuse verdict

Build fresh — one `lipgloss` row.

---

## 13 · Confirmation — destructive

### 13.1 · What counts as destructive — a closed list

**Irreversible loss of the player's data, and nothing else:**

1. Deleting a hand-added physical game (there is no store to sync it back from).
2. Disconnecting a store (drops its synced rows).
3. Resetting the library / deleting the SQLite file.

**Marking a game `ABANDONED` is not destructive** — it is a reversible state change and must
never raise a confirmation. Confirming reversible actions trains players to press `y` without
reading, which is how the irreversible one eventually gets confirmed too.

### 13.2 · Anatomy

```
┌────────────────────────────────────────────┐
│ ▌ DELETE "Hollow Knight"                   │
│                                            │
│   Added by hand. Removing it is permanent  │
│   — there is no store to sync it back      │
│   from.                                    │
│                                            │
│   y  delete            esc  keep           │
└────────────────────────────────────────────┘
```

- Border `--z-border-strong` (`67`) — a control boundary, so 1.4.11 applies and 4.08 satisfies it.
- `▌` annunciator in `--z-scanner` (`9`) — alarm as structure, not as text.
- Heading UPPERCASE bold; body `--z-text`.
- **The default is the safe action.** `esc` keeps. There is no pre-selected `y`.
- May paint `--z-surface-overlay` for modality, but must not depend on it (§1.3).

### 13.3 · Accessibility

`esc` always cancels (2.1.2). The overlay is sized and positioned so it **does not entirely
obscure the focused row** behind it (2.4.11). Focus is trapped *within* the dialog while open —
which is correct, because there is a documented way out.

### 13.4 · `NO_COLOR` / 40 columns

Zero SGR; the box, the uppercase heading and the word `permanent` carry the weight. At 40
columns the box spans `BodyRect` and the key hints stack one per line. **The consequence
sentence is never truncated.**

### 13.5 · Reuse verdict

**`huh` fits** for confirm dialogs and inherits its built-in breathing room (title → field →
help, each with a deliberate gap). Adopt it, restyled to Zerado tokens — `huh`'s default theme
carries its own palette and must not ship.

---

## 14 · Command palette — inventoried, earned in Phase 2

**Not built in Phase 1.** Recorded here so the vocabulary exists and nobody invents a different
one later.

**Why it waits.** A palette earns its place when the verb count exceeds what a key map can hold.
Phase 1 has a small, learnable verb set, and it ships `?` — a key map available from every
screen (WCAG 3.2.6) — which is the honest answer at this size. Shipping a palette over nine
commands is decoration.

**The trigger:** more than roughly a dozen global verbs, or more than roughly eight screens.
Phase 2's enrichment work is the expected crossing point.

**Reserved anatomy** (do not implement yet):

```
┌────────────────────────────────────────────┐
│ ▸ sync▎                                    │
├────────────────────────────────────────────┤
│   sync steam            s                  │
│   sync all              S                  │
│   settings                                 │
└────────────────────────────────────────────┘
```

Border `--z-border-strong`; the sigil `▸` and matched substrings in `--z-primary`; the key
column `--z-text-tertiary`. `esc` closes (2.1.2); it must not entirely obscure the focused row
(2.4.11). Reuse: `bubbles/textinput` + `bubbles/list`.

---

## 15 · Component index

| Component | § | Reuse verdict |
|---|---|---|
| Screen header band | 2 | Build fresh against #2435; `lipgloss` |
| State chip | 3 | Build fresh; `lipgloss` + width-aware pad |
| Game row / ledger | 4 | **Build a Zerado ledger primitive**; `bubbles/viewport` + `lipgloss`. `bubbles/table` does not fit; `LedgerTable` is not importable |
| Status bar | 5 | Build fresh; `lipgloss` |
| Detail pane | 6 | `bubbles/viewport`; `glamour` in Phase 2, restyled |
| Filter bar | 7 | **`bubbles/textinput` — direct fit** |
| Progress readout | 8 | Build fresh; **`bubbles/progress` does not fit** (truecolor gradient) |
| Scanner | 9 | Build fresh; **not `harmonica`** (spring ≠ sinusoid) |
| Empty state | 10 | Build fresh; `lipgloss` |
| Error state | 11 | Build fresh; **not `charmbracelet/log`** (developer logging role) |
| Degrade banner | 12 | Build fresh; `lipgloss` |
| Confirmation | 13 | **`huh` — fits**, restyled to Zerado tokens |
| Command palette | 14 | Phase 2; `bubbles/textinput` + `bubbles/list` |

---

## 16 · Open for the founder

1. **Nine ANSI-256 indices are underived** (`00-design-brief.md` §9). Four components here ship
   an interim uncoloured rendering because of it: the determinate track (§8.2), the scanner
   track (§9.4), error text (§11.2) and tertiary text. They are correct and honest, but they are
   not the designed colour. Confirm the derivation is scheduled before Phase 1 screens are built.
2. **The two-line row below 60 columns** (§4.2), which is the spine's answer to keeping the state
   label at every tier. It halves visible rows at Narrow and Tiny — 8 games instead of 16 on a
   40 × 24 terminal. Confirm that reading a full row beats scanning a truncated one at that width.
3. **No letterspacing in the terminal** (§1.5). The readout role loses its 0.18em tracking.
   Confirm — or nominate the single hero label that may carry it.
4. **The detail pane starts at 120 columns** (§6.1), per the spine — at 80 the split would leave
   the identity column around 15 characters and fail R-10(a). Below 120 the detail is a route,
   not a pane. Confirm the single detail spec mounted two ways.
