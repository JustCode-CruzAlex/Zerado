---
title: Zerado — Z-04 Library
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-04
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-04 · Library

> The root. The home. The screen a player has open when they are not doing anything else, and
> the one nominated to become Zerado's in-house gold standard the moment it is ruled GOLDEN
> (`00-design-brief.md` §4, R-1).

**Canon that governs this screen, as pointers, not copies:**
`00-design-brief.md` (§10 is the bar) · `01-design-system.md` (§2 band · §3 chip · §4 row ·
§5 status bar · §10 empty · §12 banner) · `02-colour-budget.md` (§10 is the second bar) ·
`03-designer-manual.md` (§3 is this document's contract, §5 is closed, §7 is the defect list) ·
`docs/blueprint/02-composition.md` §1.3 · §2.2 · §2.2.1 · §2.3 (**binding**) ·
`03-responsive.md` · `04-navigation-and-focus.md` · `05-state-machine.md` §7 ·
`07-offline-contract.md` §2–§5 · FlowForge TUI Design Manual **R-10** · Spacing Canon #2435.

---

## 1 · Identity

| | |
|---|---|
| **ID** | `Z-04` |
| **Name** | Library |
| **Phase** | 1 |
| **Kind** | **Route — the root.** Stack position `[1]`. Can never be popped (`04-navigation-and-focus.md` §1 rule 1) |
| **Route in** | Start-up, when a library exists **or** a provider has ever been connected. From `Z-01`, `Z-03`, `Z-05` (`Esc`), `Z-08`, `Z-09`, `Z-10` |
| **Route out** | `⏎` → `Z-05` · `s` → `Z-06` · `/` → `Z-07` (a **mode**, not a route) · `a` → `Z-08` · `r` → `Z-03` · `,` → `Z-09` · `?` → `Z-10` |
| **Shape** | **Ledger** — shape 1 of the five (`02-composition.md` §3) |
| **Composition** | ≤ 119 cols: single-pane list · ≥ 120 cols: **list ∥ detail** |
| **Offline class** | **WORKS** (`07-offline-contract.md` §2). Identical with the network off — every row, every state, every count |

---

## 2 · Purpose

**One sentence:** every game the player owns, one per row, with its state readable on every row
at every width, and a count that never lies about how many there are.

---

## 3 · Mockup — 80 × 24, the design floor and the primary breakpoint

`tier = Wide` · `leftInset = 3` · **body = 74 × 16** · every visible row begins at **column 4**.
Header-left equals content-left. Drawn to exact cell count; the outer rule is the terminal edge.

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
│     STATE           TITLE                                        HOURS   SRC   │
│     ◐  IN PROGRESS  Baldur's Gate 3                                87h   STM   │
│     ○  NOT STARTED  Blasphemous                                     0h   STM   │
│     ◉  ZERADO       Celeste                                        14h   STM   │
│     ○  NOT STARTED  Chrono Trigger                                   —   PHY   │
│     ◐  IN PROGRESS  Dark Souls III                                 63h   STM   │
│     ⊘  ABANDONED    Disco Elysium                                  22h   STM   │
│     ◉  ZERADO       Hades                                          58h   STM   │
│     ○  NOT STARTED  Hollow Knight                                  41h   STM   │
│     ◐  IN PROGRESS  Outer Wilds                                    12h   STM   │
│   ▌ ◉  ZERADO       Return of the Obra Dinn                         9h   STM   │
│     ⊘  ABANDONED    Sekiro: Shadows Die Twice                       3h   STM   │
│     ○  NOT STARTED  Sid Meier's Civilization VI: Gathering St…      0h   STM   │
│   ROWS  4–15 of 247                                                            │
│   ↑↓ move   ⏎ open   s status   / filter   a add   r sync   ? help   q quit    │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Character counts, stated because a mockup that does not add up is the defect this spec
exists to prevent:**

| Line | Cells |
|---|---|
| Terminal | **80 × 24** |
| Body content, every body row | **74** |
| Summary row `247 games · … · 31 abandoned` | 70 |
| Column-header row | 74 |
| Every game row | **74**, invariant |
| Scroll-position row `ROWS  4–15 of 247` | 17 |
| Footer key line | 73 |
| Header band | rows 2–4 (breadcrumb · gap · title) |

### 3.1 · Row map — 80 × 24

| Terminal row | Content | Token that puts it there |
|---|---|---|
| 1 | blank | `OuterMarginY` = 1 |
| 2 | breadcrumb `Zerado ✦ Library` | `HeaderBandHeight` row 1 |
| 3 | blank | `InterElementGap` = 1 (inside the band) |
| 4 | title `LIBRARY` | `HeaderBandHeight` row 3 |
| 5 | blank | `InnerPaddingY` = 1 |
| 6 | blank | `InterElementGap` = 1 (band → body) |
| **7** | **body 1** — pinned summary | `BodyRect` |
| 8 | body 2 — respiro | `02-composition.md` §2.2 |
| 9 | body 3 — column header | |
| 10–21 | **body 4–15 — the scroll region, 12 game rows** | |
| 22 | body 16 — scroll position | |
| 23 | footer key line | the canon's **reserved footer row** |
| 24 | blank | `OuterMarginY` = 1 |

### 3.2 · The 74-column game row — the field budget (binding, `02-composition.md` §2.2)

| Field | Cols | Cell range | Note |
|---|---|---|---|
| focus field | **2** | 1–2 | `▌` U+258C is **Ambiguous** width — the field is padded to 2 measured cells (§2.2.1) |
| glyph field | **2** | 3–4 | `○` and `◐` are **Ambiguous**; `◉` and `⊘` are Neutral. Padded to 2 measured cells |
| chip gap | 1 | 5 | |
| state label field | **11** | 6–16 | `NOT STARTED` and `IN PROGRESS` are the longest, at 11 |
| gutter | 2 | 17–18 | |
| **title — the identity column, R-10(a)** | **42** | **19–60** | |
| gutter | 2 | 61–62 | |
| playtime | 6 | 63–68 | right-aligned |
| gutter | 2 | 69–70 | |
| source | 4 | 71–74 | right-aligned — see D-04-6 |
| | **74** | | |

The glyph field + gap + label field is the ratified **14-column state chip**
(`01-design-system.md` §3.1). Verified at source with `unicodedata` (UCD 16.0.0, 2026-08-25):
`○` U+25CB = **A**, `◐` U+25D0 = **A**, `◉` U+25C9 = **N**, `⊘` U+2298 = **N**,
`▌` U+258C = **A**, `✦` U+2726 = **N**, `—` U+2014 = **A**, `·` U+00B7 = **A**,
`…` U+2026 = **A**, `↑`/`↓` U+2191/U+2193 = **A**, `⏎` U+23CE = **N**.
**Every one of those lives inside a padded field or on a line that is measured, never counted.**

---

## 4 · Mockup — 120 × 40, ExtraWide: list ∥ detail

`leftInset = 4` · **body = 112 × 32** · split **66 ∥ 2 ∥ 44** (`02-composition.md` §2.3).
The list sheds its `source` column; the pane carries it. **28 game rows.**

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                                      │
│    Zerado ✦ Library                                                                                                  │
│                                                                                                                      │
│    LIBRARY                                                                                                           │
│                                                                                                                      │
│                                                                                                                      │
│    247 games  198 not started  12 in progress  6 zerado  31 abandoned  ┌ DETAIL ──────────────────────────────────┐  │
│                                                                        │  Return of the Obra Dinn                 │  │
│      STATE           TITLE                                      HOURS  │                                          │  │
│      ◐  IN PROGRESS  Baldur's Gate 3                              87h  │  ◉  ZERADO                               │  │
│      ○  NOT STARTED  Blasphemous                                   0h  │                                          │  │
│      ◉  ZERADO       Celeste                                      14h  │  PLAYTIME     9h                         │  │
│      ○  NOT STARTED  Chrono Trigger                                 —  │  LAST PLAYED  2 Aug 2026                 │  │
│      ◐  IN PROGRESS  Dark Souls III                               63h  │  ADDED        12 Mar 2026                │  │
│      ⊘  ABANDONED    Disco Elysium                                22h  │  SOURCE       Steam                      │  │
│      ◉  ZERADO       Hades                                        58h  │                                          │  │
│      ○  NOT STARTED  Hollow Knight                                41h  │  SET BY       you, 12 Aug 2026           │  │
│      ◉  ZERADO       Inscryption                                  11h  │  STEAM SAYS   IN PROGRESS                │  │
│      ◐  IN PROGRESS  Kentucky Route Zero                           5h  │                                          │  │
│      ○  NOT STARTED  Metroid Dread                                  —  │                                          │  │
│      ⊘  ABANDONED    Nier: Automata                               28h  │                                          │  │
│      ◐  IN PROGRESS  Outer Wilds                                  12h  │                                          │  │
│      ○  NOT STARTED  Pentiment                                     0h  │                                          │  │
│      ◉  ZERADO       Portal 2                                     13h  │                                          │  │
│      ○  NOT STARTED  Prey                                          0h  │                                          │  │
│    ▌ ◉  ZERADO       Return of the Obra Dinn                       9h  │                                          │  │
│      ⊘  ABANDONED    Sekiro: Shadows Die Twice                     3h  │                                          │  │
│      ○  NOT STARTED  Sid Meier's Civilization VI: Gathering …      0h  │                                          │  │
│      ◐  IN PROGRESS  Signalis                                      7h  │                                          │  │
│      ◉  ZERADO       Slay the Spire                               44h  │                                          │  │
│      ○  NOT STARTED  Stardew Valley                                0h  │                                          │  │
│      ⊘  ABANDONED    Subnautica                                   16h  │                                          │  │
│      ◐  IN PROGRESS  Tunic                                         6h  │                                          │  │
│      ○  NOT STARTED  The Legend of Zelda: Breath of the Wild        —  │                                          │  │
│      ◉  ZERADO       Undertale                                     8h  │                                          │  │
│      ○  NOT STARTED  Vampire Survivors                             0h  │  LAST SYNCED  3 hours ago                │  │
│      ⊘  ABANDONED    Wasteland 3                                   9h  │                                          │  │
│    ROWS  4–31 of 247                                                   └──────────────────────────────────────────┘  │
│    ↑↓ move   ⏎ detail   tab pane   s status   / filter   a add   r sync   ? help   q quit                            │
│                                                                                                                      │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Counts:** list column **66** · gutter **2** · detail pane **44** (border 1 + inset 2 + content
**38** + inset 2 + border 1) = **112**. Footer 86 ≤ 112. ExtraWide game row: focus 2 · chip 14 ·
gutter 2 · **title 40** · gutter 2 · playtime 6 = 66.

> **D-04-4 · At ExtraWide the ledger is unbordered; only the detail pane carries a border.**
> Three reasons. (1) A border around the ledger would push its content 3 columns right of the
> title row and break the canon's headline invariant — header-left **equals** content-left
> (Spacing Canon **#2435 §5** — §5.1 names the tokens, §5.2 fixes the values;
> `00-design-brief.md` §10 line 1). (2) It would cost 2 columns of the identity column and
> 2 rows of the visible row count for zero information. (3) `02-colour-budget.md` §7.1 already
> names the sanctioned separator: *a border and a two-column gutter* — the pane's own border and
> the gutter carry it; a second border is redundant. Elevation is still never carried by fill.

> **Consequence for the title column, recorded honestly.** The ExtraWide title field is **40**
> columns — *narrower* than Wide's 42 — because the list pane is narrower than the full Wide
> body even after shedding `source`. R-10(a) still holds at the composition level: the pane
> beside it renders the selected game's title **untruncated and wrapped**.
> `01-design-system.md` §4.4 now carries the same figure — **40** — and names the trap in the
> same words (*"ExtraWide does not widen the title, it narrows it"*). The two documents agree.

---

## 5 · Visual hierarchy — what the eye reaches, in order

| # | What | Channel that carries it (case → weight → colour role → box drawing → spacing) |
|---|---|---|
| **1** | **`LIBRARY`** | UPPERCASE + bold + `--z-primary` amber, alone on its row with a respiro above and below. The only amber block above the body |
| **2** | **The numbers in the pinned summary** | Numerals in `--z-text` (16.65) against `--z-text-secondary` (9.36) words — a weight step inside one line, no colour spend |
| **3** | **The focused row** | The `▌` gutter marker (position) + **bold** (weight) + `--z-primary` on the marker (colour). The only bold row in the scroll region |
| **4** | **The state column** | A fixed vertical stripe at cells 3–16 on every row. The `◉ ZERADO` chips are the only cyan on screen, so they read by rarity, not by loudness |
| **5** | **The title column** | The widest field, `--z-text`, left-aligned at a fixed column — the identity column you actually read |
| **6** | Chrome | Breadcrumb, column header, scroll position, footer — all `--z-text-secondary`, all quiet |

**The one thing a player should see first is the title `LIBRARY`; the one thing they should see
second is the count.** Everything below that is a ruled column block that rewards scanning, not
attention.

**No hairline rule under the column header — considered and rejected.** `--z-border` may draw a
decorative hairline (`02-colour-budget.md` §8.1) and one would be legitimate here. It is declined
because it costs a game row at the 16-row floor and the respiro above the header already
separates the pinned block from the ledger. Spacing before ink (`01-design-system.md` §1.1).

**No graphical scroll indicator — considered and rejected.** The design system has no
scroll-indicator component (`01-design-system.md` §15) and this spec composes, it does not
invent. `ROWS  4–15 of 247` says the number, which is the brand's own instruction (§8).

---

## 6 · Every applied spacing token, by name, with its value at the rendered tier

| Token | Tiny `<40` | Narrow `40–59` | Standard `60–79` | **Wide `80–119`** | ExtraWide `120+` |
|---|---|---|---|---|---|
| `OuterMarginX` | 0 | 1 | 2 | **2** | 2 |
| `OuterMarginY` | 0 | 1 | 1 | **1** | 1 |
| `InnerPaddingX` | 1 | 1 | 1 | **1** | 2 |
| `InnerPaddingY` | 0 | 1 | 1 | **1** | 1 |
| `InterElementGap` | 1 | 1 | 1 | **1** | 1 |
| `HeaderBand(tier, false)` | **1** | 3 | 3 | **3** | 3 |
| **`leftInset` = `OuterMarginX + InnerPaddingX`** | **1** | 2 | 3 | **3** | 4 |
| `BodyRect.w` = `width − 2·leftInset` | 30 | 36 | 54 | **74** | 112 |
| `BodyRect.h` | 21 | 16 | 16 | **16** | 32 |

`hasSubtitle` is **`false`** on this screen and on every Zerado screen
(`02-composition.md` §1.2). The band is always the **3-row base**. The single-sizer desync
cannot occur here because the flag has exactly one value.

**Applied, not merely quoted:**

| Surface | Token | Value at Wide |
|---|---|---|
| Frame inset, all four sides | `OuterMarginX` / `OuterMarginY` | 2 cols / 1 row |
| Inside the frame | `InnerPaddingX` / `InnerPaddingY` | 1 col / 1 row |
| breadcrumb → title | `InterElementGap` | 1 row |
| band → body | `InterElementGap` | 1 row |
| Header-left **and** content-left | `leftInset` | **column 4** |
| Footer | the canon's reserved footer row | 1 row, not stolen from `BodyRect` |
| Detail pane inset (ExtraWide) | D-06-1 bordered-surface inset | 2 cols each side, 0 rows |

**Zero magic numbers.** The only non-token numbers on this screen are the field widths of §3.2,
which are the spine's binding column budget, and the summary/footer degrade thresholds of §10.3
and §9.2, which are declared ladders rather than constants.

---

## 7 · Colour, glyph and label for every state shown

Ratios are read from the brand manual's measured table (§4.2). **None is estimated.**

| Shown as | Token | Hex | ANSI-256 | 16-colour | Glyph | ASCII | Label | Ratio |
|---|---|---|---|---|---|---|---|---|
| Not started | `--z-state-not-started` | `#A5A29B` | **247** | `white` | `○` U+25CB | `[ ]` | `NOT STARTED` | **7.62** AA |
| In progress | `--z-state-in-progress` | `#FFB000` | **214** | `bright yellow` | `◐` U+25D0 | `[~]` | `IN PROGRESS` | **10.59** AAA |
| Zerado | `--z-state-zerado` | `#19E0FF` | **45** | `bright cyan` | `◉` U+25C9 | `[*]` | `ZERADO` | **12.15** AAA |
| Abandoned | `--z-state-abandoned` | `#C77DFF` | **177** | `bright magenta` | `⊘` U+2298 | `[x]` | `ABANDONED` | **7.21** AA |
| Screen title `LIBRARY` | `--z-primary` | `#FFB000` | **214** | `bright yellow` | — | — | the word | **10.59** AAA |
| Focus marker `▌` | `--z-primary` | `#FFB000` | **214** | `bright yellow` | `▌` U+258C | `>` | — | **10.59** AAA |
| Game title | `--z-text` | `#E9EEF5` | **255** | `bright white` | — | — | the title | **16.65** AAA |
| Breadcrumb, column header, summary words, scroll row, footer | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | — | — | — | **9.36** AAA |
| Summary numerals | `--z-text` | `#E9EEF5` | **255** | `bright white` | — | — | — | **16.65** AAA |
| Unreported playtime `—`, truncation `…` | `--z-text-tertiary` | `#8492A8` | ***underived*** | `white` | — | — | — | **6.15** AA |
| Banner `▌`, informational | `--z-border-strong` | `#64748B` | **67** | `bright black` | `▌` | `>` | the label word | **4.08** (1.4.11) |
| Banner `▌`, action required | `--z-primary` | `#FFB000` | **214** | `bright yellow` | `▌` | `>` | the label word | **10.59** AAA |
| Detail-pane border, focused (ExtraWide) | `--z-border-strong` | `#64748B` | **67** | `bright black` | `┏━┓` | `+-+` | — | **4.08** |
| Detail-pane border, unfocused | `--z-border` | `#2A3342` | **236** | `black` | `┌─┐` | `+-+` | — | 1.53 — **decoration only, and it bounds no control here** |
| Audio annunciator, unmuted | `--z-primary` | `#FFB000` | **214** | `bright yellow` | `▮` U+25AE (**N**) | *none — label alone* | `AUDIO` | **10.59** AAA |
| Audio annunciator, muted | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | `▯` U+25AF (**N**) | *none — label alone* | `MUTED` | **9.36** AAA |

> **The annunciator has no ASCII glyph, deliberately.** `[ ]` and `[*]` are the *state column's*
> tokens — *not started* and *zerado* — and on an ASCII-mode library screen they would appear twice
> on one row meaning two different things. Co-render survives, because the labels differ, but
> confusability on the most-used screen is what this bundle rejects everywhere else. The annunciator
> already carries a label word, so the glyph is its third channel and the label is unambiguous
> alone. `03-responsive.md` §5c is the authority.

> ***underived*** means exactly that: **no ANSI-256 index has been derived for
> `--z-text-tertiary` `#8492A8`.** Nobody may pick one at the keyboard (brand §10 rule 5).
> **Interim rendering on this screen: tertiary text renders uncoloured** — no SGR at all — which
> invents nothing and degrades identically to the `NO_COLOR` path. Owner:
> `fft-brand-architect`; tracked in `00-design-brief.md` §9.

**The co-render rule holds on every row at every tier: colour AND glyph AND label, all three.**
The label is never dropped for density; below 60 columns the *row* grows a second line instead
(§11). The tightest CVD pair in the system is **zerado × abandoned at ΔE 11.9 under
deuteranopia** — the floor to protect, and the one place the glyph and the word genuinely carry
load rather than merely reinforce.

### 7.1 · The fourth channel — audio (spine delta, 2026-08-25)

Audio ships in Phase 1. `03-designer-manual.md` §5.11 verdict 3 is **superseded by founder
direction relayed 2026-08-25 through the spine**, and the manual now carries the strike-through
and the provenance note. The two documents agree.

| Audio state | Glyph | Label | Colour | In the footer? |
|---|---|---|---|---|
| Never enabled — **the default** | — | — | **no indicator at all** | **`m` is not listed** |
| Enabled, unmuted | `▮` U+25AE | `AUDIO` | `--z-primary` `#FFB000` / 214 / `bright yellow` | `m mute` |
| Enabled, muted | `▯` U+25AF | `MUTED` | `--z-text-secondary` `#A9B5C7` / 249 / `white` | `m unmute` |

`▮` and `▯` were verified **Neutral** width against UCD 16.0.0 on 2026-08-25 — unlike `♪` U+266A
(**Ambiguous**), they cannot shear a line. They carry the same filled/hollow logic as the state
glyphs, so the vocabulary is already learned.

> **D-A1 · The annunciator is right-aligned on the header band's TITLE row, on every screen,
> at every tier.** *Alternative considered and rejected:* the pinned summary row. It was rejected
> because `Z-05` as a route, `Z-06` as an overlay and `Z-10` have **no** pinned summary, which
> would force a second position and destroy the one property the indicator needs — WCAG 4.1.3's
> *stable and predictable place*. The title row is the one surface every framed screen has at
> every tier (`HeaderBand` ≥ 1 always). It does not repurpose the footer, which
> `04-navigation-and-focus.md` §6 forbids from carrying status; and it does not compete with the
> summary's counts, which R-10(c) forbids truncating.
> `01-design-system.md` §2.4's objection — *"a live value would make the band the only mutable
> chrome"* — is aimed at a self-changing count; this is a two-state annunciator that changes
> only on an explicit `m`.

**Audio is never the only carrier.** Anything that makes a sound is also visible, as a fourth
co-render channel. On `Z-04` there are exactly two candidate sound events and both already have
a visible carrier: a completed sync (the summary's numbers change, and a transition into
*zerado* prints the result line of §10.5) and a failed sync (the degrade banner). **No sound may
be introduced on this screen without a visible carrier already in this spec.**

```
   LIBRARY                                                            ▮ AUDIO
```
```
   ↑↓ move  ⏎ open  s status  / filter  a add  r sync  m mute  ? help  q quit
```
*(74 cells — with `m` present the footer separator tightens from 3 to 2 per §9.2. Default,
audio never enabled: the 73-cell 3-space line of §3.)*

---

## 8 · The full state table

Every row is a real screen, including the one nobody writes down.

| # | State | Trigger | Body composition | Rows lost | Copy |
|---|---|---|---|---|---|
| **L1** | **First run — nothing synced, nothing connected** | Library empty **and** no provider ever connected. `Z-01` is pushed as root instead, but the player reaches `Z-04` from it | Empty state (a), top-aligned in `BodyRect` | — | §10.1 |
| **L2** | **Empty because the provider returned nothing** | A sync succeeded and returned **0** items | Degrade banner + empty state (c) | — | §10.2 |
| **L3** | **Loading** | Reading SQLite on first paint | Summary renders `— games`; the scroll region renders **nothing**, not a spinner. **No scanner** — this is a local read, not an indeterminate wait | — | §10.6 |
| **L4** | **Populated — the default** | ≥ 1 row | §3 | — | §10.3 |
| **L5** | **Populated + a degrade banner** | Any of B1–B7 (§8.1) | Banner takes body row 1; the pinned block grows down | **1 game row** (12 → 11) | §10.4 |
| **L6** | **Filter active** | `/` or `f` | The **mode** — see `Z-07-filter-and-search.md`. The filter bar replaces the pinned summary | 1 | Z-07 |
| **L7** | **Filter active, zero matches** | Filter matches 0 | Z-07's zero-result block; the column header and scroll row are **absent** | — | Z-07 |
| **L8** | **A pending write** | `Z-06` applied, SQLite not yet confirmed | The row's chip label is suffixed `…` — `ZERADO…` — until the write confirms. **Never an optimistic silent change** (`01-design-system.md` §3.4) | — | §10.5 |
| **L9** | **A transition into *zerado* just landed** | `Z-06` wrote `zerado` | The pinned summary is replaced by the result line until the next keypress | — | §10.5 |
| **L10** | **412 rows, scrolled to the end** | Overflow | R-10(a)(b)(c) all hold — §12 | — | — |
| **L11** | **Detail pane focused** (ExtraWide) | `Tab` or `⏎` | Pane border goes heavy; the ledger cursor stays visible in chrome, not amber | — | §13 |
| **L12** | **Below the refusal floor** | `< 24` cols or `< 8` rows | **Frameless.** One sentence, `exit 2` at start-up; a running session keeps running | all | §11.3 |
| **L13** | **A sync stopped returning some rows — they are tombstoned and out of the default view** | `absent_since IS NOT NULL` on ≥ 1 row (`06-data-seams.md` §2.4) | The ledger renders **exactly as L4** — the absent rows are simply not in the default row set. The pinned summary appends `<n> absent` and degrades one step to make room (§10.3) | — | §10.3 |

> **L13 · The default-view exclusion rule, stated because a row that is invisible has to be
> accounted for somewhere.** `06-data-seams.md` §2.4 decides it: a game a sync stops returning is
> **tombstoned, never deleted** — `absent_since` is set, the row stays, and it carries the
> player's own work. Three consequences bind this screen:
>
> 1. **Absent rows are excluded from the default row set.** `Z-04`'s ledger, its scroll
>    arithmetic, its `ROWS n–m of N` readout and its cursor all operate on the shown set. Nothing
>    on this screen renders an absent row.
> 2. **The summary's counts therefore describe the set actually shown.** `05-state-machine.md`
>    §7 rule 1 holds unchanged — the four counts sum to the number shown — because both figures
>    are the shown set.
> 3. **And the summary says so.** See §10.3 and **D-04-8**: the moment the default view stops
>    being the whole file, silence would be the *"list view lies"* failure rule 2 names.
>
> **`absent` is not a fifth state and never renders as one.** It is an orthogonal presence flag
> (§2.4). An absent game still carries one of the four states — usually the most valuable one —
> and if it is ever shown (via `Z-07`'s facet) it renders with its own chip, unchanged. **No
> fifth glyph, no fifth colour, no fifth chip in the ledger.**
>
> **Only a sync whose status is `ok` may tombstone.** A `partial`, `failed` or `cancelled` sync
> must not, because *not returned* and *not reached* are indistinguishable in a truncated stream.
> **This screen never initiates that** — `r` pushes `Z-03` — but it is the reason `Z-04` can
> trust that a drop in the count is a fact rather than an artefact.
>
> **`absent_since` is cleared silently when the game comes back.** No banner, no result line, no
> celebration: the count in the clause simply goes down, and at zero the clause stops rendering.
> The *zerado* result line (§10.5) is the only moment this screen announces, and this is not it.

**`Z-04` is never a sync screen.** `r` pushes `Z-03`. `07-offline-contract.md` §5 is explicit —
*"There is no connectivity check, no background ping… Between the two, nothing is running."*
So **no progress readout and no scanner ever render on this screen**, one at a time or
otherwise. `01-design-system.md` §5.4's *"Syncing → replaced by the progress readout"* applies
to `Z-03`'s status bar, not to this one.

### 8.1 · The banner set — Phase 1, and only one at a time

Raised by a **failure**, cleared by the next success (`07-offline-contract.md` §5). It never
appears when nothing is degraded — a banner that is always there is furniture.

| # | Trigger (the classified failure) | Class | Copy, exact | Cells |
|---|---|---|---|---|
| B1 | no route / DNS | informational · chrome `▌` | `▌ OFFLINE   Last synced 3 hours ago. Everything here still works.` | 65 |
| B2 | timeout / 5xx | informational · chrome `▌` | `▌ UNREACHABLE   Steam didn't answer. Last synced 3 hours ago. Press r to retry.` | 79 → **see note** |
| B3 | 401 / 403 | **action** · amber `▌` | `▌ STEAM REJECTED THAT KEY   Check it hasn't been regenerated.` | 61 |
| B4 | 200 + empty | **action** · amber `▌` | `▌ STEAM PROFILE PRIVATE   Steam won't share the list until it's public.` | 71 |
| B5 | no credential stored | **action** · amber `▌` | `▌ STEAM KEY MISSING   Press s to add it.` | 40 |
| B6 | a sync ended mid-stream | **action** · amber `▌` | `▌ SYNC INCOMPLETE   147 of 247 synced. Nothing lost. Press r to finish.` | 71 |
| B7 | last success ≥ 90 days | **action** · amber `▌` | `▌ LAST SYNCED IN MAY   Playtime and new games are that old. Press r to sync.` | 76 → **see note** |

> **B2 and B7 exceed 74 cells at Wide.** They truncate from the right **at a word boundary**,
> never mid-word (`01-design-system.md` §12.5), and the **label word never truncates**. At Wide
> they render as `…Steam didn't answer. Last synced 3 hours ago. Press r to` and
> `…Playtime and new games are that old. Press r to` respectively — which is why B2's and B7's
> action key is *also* in the footer. The label word plus the footer key survive every width.
> An alternative that fits — dropping the age from B2 — was rejected because
> `07-offline-contract.md` §4 makes the age mandatory on any network-derived value.

**B5 and B7 carry Phase-1-correct copy.** `01-design-system.md` §12.3's `STEAM KEY MISSING`
row is reproduced verbatim. `07-offline-contract.md` §4.1 now carries the Phase 1 wording
itself — *`Last synced in May. Anything you have played since then is missing. r to sync.`* —
and defers the price wording to Phase 3 with its reason. B7 above is that wording.

**Priority when two conditions hold — at most one banner ever renders:**
`B6 > B5 > B3 > B4 > B7 > B2 > B1`. Action-required outranks informational; the most specific
outranks the most general.

### 8.2 · Mockup — L5, a sync that stopped halfway

The list is **real but partial**. 11 game rows, not 12: the banner is on the never-hide list
(`03-responsive.md` §4) and the respiro is not spent (D-04-2).

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   ▌ SYNC INCOMPLETE   147 of 247 synced. Nothing lost. Press r to finish.      │
│   247 games · 198 not started · 12 in progress · 6 zerado · 31 abandoned       │
│                                                                                │
│     STATE           TITLE                                        HOURS   SRC   │
│     ◐  IN PROGRESS  Baldur's Gate 3                                87h   STM   │
│     ○  NOT STARTED  Blasphemous                                     0h   STM   │
│     ◉  ZERADO       Celeste                                        14h   STM   │
│     ○  NOT STARTED  Chrono Trigger                                   —   PHY   │
│     ◐  IN PROGRESS  Dark Souls III                                 63h   STM   │
│     ⊘  ABANDONED    Disco Elysium                                  22h   STM   │
│     ◉  ZERADO       Hades                                          58h   STM   │
│     ○  NOT STARTED  Hollow Knight                                  41h   STM   │
│     ◐  IN PROGRESS  Outer Wilds                                    12h   STM   │
│   ▌ ◉  ZERADO       Return of the Obra Dinn                         9h   STM   │
│     ⊘  ABANDONED    Sekiro: Shadows Die Twice                       3h   STM   │
│   ROWS  4–14 of 247                                                            │
│   ↑↓ move   ⏎ open   s status   / filter   a add   r sync   ? help   q quit    │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

> **D-04-2 · The pinned chrome block grows downward; the respiro row below it is never spent;
> the scroll region absorbs every additional chrome row.** One rule, and the arithmetic closes
> in every case at the 16-row floor:
>
> | Composition | rows | Games |
> |---|---|---|
> | summary · respiro · header · **list** · scroll | 1+1+1+**12**+1 = 16 | 12 |
> | banner · summary · respiro · header · **list** · scroll | 1+1+1+1+**11**+1 = 16 | 11 |
> | filter · summary · respiro · header · **list** · scroll | 1+1+1+1+**11**+1 = 16 | 11 |
> | banner · filter · summary · respiro · header · **list** · scroll | 1+1+1+1+1+**10**+1 = 16 | 10 |
>
> The respiro is the one row that tells the eye the pinned block and the ledger are different
> things. Spending it to save one row of twelve would be a craft failure on **every** render to
> buy 8 % on one. Below Standard the rule changes shape — see §11.1.

### 8.3 · Mockup — L1, first run, before anything is synced

**The row nobody writes down, and the first screen every new player sees.**

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   No library yet.                                                              │
│                                                                                │
│   Zerado reads your Steam library once you add a key.                          │
│   Physical discs and cartridges can be added by hand.                          │
│                                                                                │
│   c  connect a store      a  add a game by hand                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│   c connect a store   a add by hand   , settings   ? help   q quit             │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

Top-aligned in `BodyRect`, left-aligned at `leftInset`, `InterElementGap` between blocks — never
vertically centred, because a centred block moves as soon as a line is added.

**`c  connect a store` is where this screen spends its one chrome cyan** (§14). `a  add a game
by hand` is `--z-primary` amber, per `02-colour-budget.md` §4.1 item 5.

> **The `s` collision is closed.** An earlier revision of this spec flagged `s` meaning *set
> status* on a populated library and *connect Steam* on an empty one, and recommended `c`.
> `04-navigation-and-focus.md` §3.2 adopted it: **`s` means set-status everywhere, always, in
> every state; connect-a-store is `c`.** `01-design-system.md` §10.1 and §10.3 carry the
> corrected copy and it is reproduced here verbatim. **`s` is not live on an empty library at
> all** — there is no game whose status could be set — so it is absent from this footer.

### 8.4 · Mockup — L2, Steam returned an empty library

Ratified copy, `01-design-system.md` §10.3, **verbatim**.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   ▌ STEAM PROFILE PRIVATE   Steam won't share the list until it's public.      │
│                                                                                │
│   Steam returned an empty library.                                             │
│                                                                                │
│   Game details are private on your profile — Steam won't                       │
│   share the list until that's public.                                          │
│   Settings → Privacy.                                                          │
│                                                                                │
│   r  try again        c  Steam settings                                        │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│   r try again   c Steam settings   , settings   ? help   q quit                │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

> **FLAG** · `→` U+2192 is **Ambiguous** width (verified, UCD 16.0.0). It sits on a free-flowing prose
> line that is measured, not counted, so it cannot shear a column. Recorded because it is the
> only Ambiguous glyph in ratified copy.

---

## 9 · The key map

Every key active on this screen, and nothing that is not. **A footer that lies is worse than no
footer** (`04-navigation-and-focus.md` §6).

| Key | Does | Scope | Note |
|---|---|---|---|
| `↑` `↓` `k` `j` | Move the row cursor | screen | The viewport follows — R-10(b) |
| `g` / `G` | First row / last row | global | |
| `Ctrl-D` / `Ctrl-U` | Half a page down / up | global | |
| `⏎` | Open the focused game | screen | ≤ 119: push `Z-05` · ≥ 120: **move focus into the pane** |
| `Tab` / `Shift-Tab` | Next / previous region | global | **ExtraWide only** — `R = 2`. Below 120 there is one region and `Tab` does nothing and is not listed |
| `s` | Set this game's status → `Z-06` | screen | **Set-status everywhere, always** (`04-navigation-and-focus.md` §3.2). Not live on an empty library — there is no row — and therefore not in that footer |
| `c` | Connect a store → `Z-02` | screen | **Only on an empty library** (§8.3). Not live once a row exists |
| `/` | Filter and search → `Z-07` (a mode) | screen | |
| `f` | Jump to the state chips of `Z-07` | screen | Enters filter mode with the chip row focused |
| `a` | Add a game by hand → `Z-08` | screen | |
| `r` | Re-sync → `Z-03` | screen | |
| `,` | Settings → `Z-09` | global | |
| `m` | Mute / unmute the audio | global | **Only bound and only listed when audio has been enabled** (§7.1) |
| `?` | Help → `Z-10` | global | Unwinds if already on the stack |
| `Esc` | **Nothing.** This is the root, at rest | global | The footer hints `q quit · ? help` |
| `q` | Quit, immediately | global | Nothing to confirm — every mutation already committed (`04-navigation-and-focus.md` §3.2) |
| `Ctrl-C` | Quit | global | Always, including inside a text input |

**Reserved and deliberately inert in Phase 1:** `:` and `Ctrl-K` (command palette, `Z-17`),
`1`–`9` (quick filters), `n` / `p` (next / previous game). Pressing one does **nothing** and
shows **no error**.

### 9.1 · Footer strings, exact

| Screen state | Footer | Cells |
|---|---|---|
| Populated, audio never enabled | `↑↓ move   ⏎ open   s status   / filter   a add   r sync   ? help   q quit` | **73** |
| Populated, audio enabled | `↑↓ move  ⏎ open  s status  / filter  a add  r sync  m mute  ? help  q quit` | **74** |
| ExtraWide | `↑↓ move   ⏎ detail   tab pane   s status   / filter   a add   r sync   ? help   q quit` | 86 ≤ 112 |
| First run (L1) | `c connect a store   a add by hand   , settings   ? help   q quit` | **64** |
| Empty from provider (L2) | `r try again   c Steam settings   , settings   ? help   q quit` | 61 |
| Narrow 40 | `↑↓ ⏎ / s   ? help   q quit` | 26 |
| Tiny 32 | `⏎ open  ? help  q quit` | 22 |

### 9.2 · The footer degrade ladder — declared, not improvised

1. Separator is **3 spaces**.
2. If the composed width exceeds `BodyRect.w`, the separator tightens to **2**.
3. If it still exceeds, hints drop in this order:
   `a add` → `r sync` → `m mute` → `s status` → `/ filter` → `⏎ open` → `↑↓ move`.
4. **`? help` and `q quit` are the last two to go** (`04-navigation-and-focus.md` §6).

The line is **measured** with the width-aware function (`01-design-system.md` §1.2 rule 2),
never counted — which is what makes `↑` and `↓` (both **Ambiguous**) safe here: on a terminal
that renders them double-width the composed line is 2 cells longer and the ladder simply drops
one more hint. Under `ZERADO_ASCII=1` the pair renders `up/dn`.

---

## 10 · The exact copy — ready to paste

### 10.1 · First run (L1) — `01-design-system.md` §10.1

```
No library yet.

Zerado reads your Steam library once you add a key.
Physical discs and cartridges can be added by hand.

c  connect a store      a  add a game by hand
```

### 10.2 · Provider returned nothing (L2) — `01-design-system.md` §10.3, **verbatim**

```
Steam returned an empty library.

Game details are private on your profile — Steam won't
share the list until that's public.
Settings → Privacy.

r  try again        c  Steam settings
```

### 10.3 · The pinned summary — the three forms, and the ladder

> **D-04-1 · The summary spells the state words; the glyph key is the degrade, not the default.**
> Three reasons. (1) `02-colour-budget.md` §2.4 **forbids colouring the summary** — *"the summary
> is prose about state, not a state cell"* — so a glyph there would be an uncoloured glyph, the
> weakest co-render the system permits anywhere. (2) Prose casing is already ratified for this
> row (`6 zerado`, not `6 ZERADO` — `01-design-system.md` §3.3), and prose spells words.
> (3) It removes two Ambiguous glyphs from a line that is not column-aligned.
> `05-state-machine.md` §7's rule 1 — **the four counts always sum to the number shown** — holds
> in all three forms; no fact is ever dropped.

| Form | Composed | At 247 games | Used when |
|---|---|---|---|
| **1 · full prose** | `<n> games · <a> not started · <b> in progress · <c> zerado · <d> abandoned` | `247 games · 198 not started · 12 in progress · 6 zerado · 31 abandoned` — **70** | fits `BodyRect.w` (Wide 74 ✓) |
| **2 · tight separators** | as form 1, ` · ` → `  ` | `247 games  198 not started  12 in progress  6 zerado  31 abandoned` — **66** | form 1 does not fit (ExtraWide list pane 66 ✓) |
| **3 · glyph key** | `<n> games  ○ <a>  ◐ <b>  ◉ <c>  ⊘ <d>` | `247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31` — **33** | forms 1 and 2 do not fit (Standard 54, Narrow 36) |
| **3T · Tiny** | `<n>  ○ <a>  ◐ <b>  ◉ <c>  ⊘ <d>` | `247  ○ 198  ◐ 12  ◉ 6  ⊘ 31` — **27** | Tiny 30 |

**The rule that keeps form 3 honest:** it may render only on a tier where the full state
**label** is visible in the list below. That is true at every tier — Standard and above keep the
14-column chip, and Narrow and Tiny put the label on the row's second line — so the glyph key
never becomes the only carrier of a state's name.

**Filter active:** the summary is replaced by `Z-07`'s bar
(`01-design-system.md` §5.4). See `Z-07-filter-and-search.md` §10.

**Absent rows present (L13): the summary appends `<n> absent`, and it says the number.**

> **D-04-8 · When at least one row is absent, the summary carries a trailing `<n> absent`
> clause; the clause is exempt from the degrade ladder and the base form degrades around it.**
>
> **Whether to mention the hidden count at all was the decision, and the answer is yes.**
> `05-state-machine.md` §7's preamble requires the summary to say *"which set it is describing"*,
> and rule 2 gives the filter case. A default view with absent rows is the *second* case, and it
> is rule 2's failure **in mirror image**: rule 2 forbids *"showing whole-library counts above a
> filtered list"*; silence here would show **filtered counts and call them the whole library.**
> Same lie, opposite direction. A player whose family share expires loses forty rows from view,
> and a summary that quietly reads `207 games` where it read `247 games` yesterday has told them
> their library shrank. **Silence was considered and rejected on that sentence.**
>
> *The alternative placements considered and rejected.* **A banner** — it would cost a game row
> on every render, compete for the one-banner-at-a-time slot (§8.1), and claim urgency the seam
> deliberately refuses (`absent_since` is set and cleared silently). **The scroll-position row** —
> it has 57 spare cells at Wide and looks tempting, but it does not render when nothing scrolls,
> and a library small enough not to scroll can still have absent rows. The carrier has to be the
> one row that is always there, which is the summary (R-10(c)).
>
> **The four rules that keep it honest:**
> 1. The clause renders only when `n ≥ 1`. At `n = 0` the summary is byte-identical to forms 1–3T
>    above — the common case pays nothing.
> 2. The total and the four state counts describe the **shown** set, so §7 rule 1 holds unchanged.
>    The file total is `shown + n` and is never printed as a separate figure here.
> 3. **The clause is never dropped.** It is not a rung on the ladder; the ladder runs beneath it.
>    A screen too narrow to carry both the clause and four state buckets does not exist — step 5
>    composes to 30 cells at the 30-column Tiny floor.
> 4. **`absent` is the same word on all four surfaces** — this clause, `Z-07`'s `[ABSENT]` chip,
>    `Z-07`'s `absent   yes` diagnostic line, and `Z-05`'s `SOURCE` value. One word, one meaning,
>    learned once. It is chosen over *missing*, *removed* and *hidden* because each of those
>    asserts something the product cannot know: that something is lost, that someone took it, or
>    that Zerado is keeping it from them.

**The composed forms, at 247 shown and 3 absent.** Counted, as everything on this screen is:

| Form | Composed | Cells | Fits |
|---|---|---|---|
| 1 + clause | `247 games · 198 not started · 12 in progress · 6 zerado · 31 abandoned · 3 absent` | **81** | ✗ at Wide 74 |
| 2 + clause | the same with `  ` separators | **76** | ✗ at Wide 74 |
| **3 + clause** | `247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31  3 absent` | **43** | ✓ ExtraWide pane 66 · Wide 74 · Standard 54 |
| 3T + clause | `247  ○ 198  ◐ 12  ◉ 6  ⊘ 31  3 absent` | **37** | ✗ at Narrow 36 |
| **unspaced + clause** (`01-design-system.md` §5.5's Tiny form) | `247  ○198 ◐12 ◉6 ⊘31  3 absent` | **30** | ✓ Narrow 36 · Tiny 30 |

**So at Wide, a 247-game library with 3 absent rows renders the glyph key rather than the prose.**
That is the ladder working, not a regression: the row changes shape on the rare day it has
something extra to say, and changes back the moment the game returns. A smaller library keeps the
prose — `84 games  61 not started  4 in progress  3 zerado  16 abandoned  3 absent` is **73** and
fits Wide in form 2.

**The ladder, declared rather than improvised** (the same shape as the footer's, §9.2):

1. Compose form 1 **plus the clause**. If it measures ≤ `BodyRect.w`, done.
2. Separators tighten ` · ` → `  `.
3. The state words become the glyph key.
4. The total sheds its noun — `247 games` → `247`.
5. The glyph counts unspace — `○198 ◐12 ◉6 ⊘31`.
6. **No state bucket is ever dropped** (`01-design-system.md` §5.6 — a summary missing a bucket no
   longer sums) **and the clause is never dropped.**

Steps 3–5 are legend-dependent, and D-04-1's guard still binds: the glyph key may render only
where the full state **label** is visible in the list below, which is true at every tier. **The
clause is not legend-dependent** — it is a word.

### 10.4 · Degrade banners

Verbatim strings are in §8.1. Every one carries the four mandatory parts of
`07-offline-contract.md` §3 — an **uppercase label word**, **what is unavailable**, **how stale
it is**, and **the key that retries it**. `▌` is **structure**, not a state channel; there is no
`⚠` anywhere in Zerado.

### 10.5 · The moment the product exists for

The brand writes this copy (§8, voice example 4). It lands here.

| Situation | Copy, exact |
|---|---|
| A write is in flight | the chip reads `ZERADO…` (label + `…`) until SQLite confirms |
| Marked *zerado*, playtime known | `Zerado. 41 hours. Sixth this year.` |
| Marked *zerado*, first this year | `Zerado. 9 hours. First this year.` |
| Marked *zerado*, playtime not reported (a hand-added copy) | `Zerado. Sixth this year.` |
| Marked *zerado*, past tenth | `Zerado. 41 hours. 14th this year.` |

> **D-04-3 · The result line replaces the pinned summary row and persists until the next
> keypress.** No timer, no dwell constant, no token gap — and therefore no WCAG 2.2.1 time
> limit at all. It is printed **as text in a stable and predictable place**, which is 4.1.3's
> terminal analogue exactly (`00-design-brief.md` §3.2).
>
> **It renders for a transition into `zerado` only.** The other three states are visible in the
> row's own chip and in the summary's counts; announcing every edit would make the summary row
> flicker and would spend this moment's weight on routine corrections.
> `05-state-machine.md` §4 is the reason it exists at all: *"a player who is told they finished
> something has not finished it; they have been notified."* The line reports what **they** just
> did.
>
> **Ordinals** are spelled to *tenth*, then numeric with a suffix — `11th`, `14th`. `Sixth this
> year.` is the ratified example and it is a word.

### 10.6 · Loading (L3)

```
— games
```
No spinner, no scanner, no per-row placeholder. Reading a local SQLite file is not an
indeterminate wait; the scanner is reserved for those and never used ambiently
(`01-design-system.md` §9.3).

### 10.7 · Copy notes

- **Casing** — `Zerado` in the breadcrumb (the product); `ZERADO` in the chip, the column and
  the filter; `zerado` in the summary's prose. Never camel-cased, never `Zerado.app`.
- **Say the number** — `247 games`, `147 of 247 synced`, `41 hours`. Never *a lot of*.
- **No exclamation marks. No emoji. The user is never a "gamer".**
- **Type-neutral where a neutral word is equally natural** — `TITLE` not `GAME`,
  `ROWS  4–15 of 247`, `247 games in the library.` This costs nothing today and means the
  media-item generalisation does not force a copy rewrite later. Phase 1 says **game** where a
  game is what the player is looking at, and the copy is not contorted to avoid it.

---

## 11 · 40-column behaviour, and the refusal floor

### 11.1 · Narrow — 40 × 24, body 36 × 16, two-line rows

> **Shed columns before you shed meaning. Shed meaning before you shed the state.**
> (`03-responsive.md` §2.) Zerado does not drop the label. It changes the row's shape.

```
┌────────────────────────────────────────┐
│                                        │
│  Zerado ✦ Library                      │
│                                        │
│  LIBRARY                               │
│                                        │
│                                        │
│  247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31     │
│    ◐ Baldur's Gate 3                   │
│      IN PROGRESS · 87h · Steam         │
│    ○ Blasphemous                       │
│      NOT STARTED · 0h · Steam          │
│    ◉ Celeste                           │
│      ZERADO · 14h · Steam              │
│    ○ Chrono Trigger                    │
│      NOT STARTED · — · Physical        │
│    ◐ Dark Souls III                    │
│      IN PROGRESS · 63h · Steam         │
│    ⊘ Disco Elysium                     │
│      ABANDONED · 22h · Steam           │
│  ▌ ◉ Return of the Obra Dinn           │
│      ZERADO · 9h · Steam               │
│  ROWS  4–10 of 247                     │
│  ↑↓ ⏎ / s   ? help   q quit            │
│                                        │
└────────────────────────────────────────┘
```

**Counts:** body 36. Line 1 = focus field 2 + glyph field 2 + **title 32**. Line 2 = 4-space
indent + `LABEL · Nh · Source` ≤ 32. Summary form 3 = 33 ≤ 36. **7 games.**

**Sheds:** the `source` column (into line 2, spelled out) · the column-header row (a two-line
row labels itself) · the respiro (see below).
**Keeps, because `03-responsive.md` §4 forbids hiding them:** the state — glyph **and** label ·
the title · the pinned summary · the focus marker · the degrade banner · the footer.

> **D-04-5 · Below Standard the scroll region takes `floor(remaining / 2)` games and any single
> leftover row becomes the respiro.** Narrow: `16 − summary(1) − scroll(1) = 14` → **7 games**,
> remainder 0, no respiro. Tiny: `21 − 1 − 1 = 19` → **9 games**, remainder 1 → the respiro
> takes it. This is the one place D-04-2 yields, and it always closes the arithmetic with no
> orphan half-row — which is the failure a fixed respiro would produce here.

### 11.2 · Tiny — 32 × 24, body 30 × 21

`OuterMarginX` and `OuterMarginY` shed to **0** and the band collapses to the **title row only**
— the canon's deliberate starvation guard. The narrowest terminal is not the shortest body.

```
┌────────────────────────────────┐
│ LIBRARY                        │
│                                │
│ 247  ○ 198  ◐ 12  ◉ 6  ⊘ 31    │
│                                │
│   ◐ Baldur's Gate 3            │
│     IN PROGRESS · 87h · Steam  │
│   ○ Blasphemous                │
│     NOT STARTED · 0h · Steam   │
│   ◉ Celeste                    │
│     ZERADO · 14h · Steam       │
│   ○ Chrono Trigger             │
│     NOT STARTED · — · Physical │
│   ◐ Dark Souls III             │
│     IN PROGRESS · 63h · Steam  │
│   ⊘ Disco Elysium              │
│     ABANDONED · 22h · Steam    │
│   ◉ Hades                      │
│     ZERADO · 58h · Steam       │
│   ○ Hollow Knight              │
│     NOT STARTED · 41h · Steam  │
│ ▌ ◐ Outer Wilds                │
│     IN PROGRESS · 12h · Steam  │
│ ROWS  4–12 of 247              │
│ ⏎ open  ? help  q quit         │
└────────────────────────────────┘
```

Title field **26**. Summary form 3T = 27 ≤ 30. **9 games.** Nothing is flush at Narrow and
above; at Tiny the frame is deliberately flush, which acceptance line 1 explicitly exempts
(*"at any tier ≥ Narrow"*).

### 11.3 · The refusal floor — 24 columns or 8 rows

Below **24 columns or 8 rows** `Z-04` does not render a degraded interface. **Frameless**, one
line:

```
Zerado needs at least 24 columns and 8 rows. This terminal is 20 x 6.
```

Exit status **2** at start-up. A **running** session resized below the floor replaces the screen
with the same sentence and **keeps running** — the player is probably dragging a divider and
will drag it back. Checked at start and on every `WindowSizeMsg`.

Three reasons this is a refusal and not a degrade: at 24 columns the title field is about 16
characters, so every row would be an ellipsis and a list of ellipses is a broken view, not a
smaller one; a degrade that cannot show the state has abandoned the co-render rule, which is the
product's accessibility mechanism and not a style; and refusing is honest, actionable and one
line, while rendering garbage is none of the three.

---

## 12 · The ledger triad — R-10, and how each leg is met at 412 rows

**A frozen golden cannot prove any of this.** Each leg below is stated as something a reviewer
can falsify from a live repro at an overflowing row count.

### (a) A populated, human-readable identity column

The title field is **42** cells at Wide, **40** at ExtraWide, **32** at Narrow, **26** at Tiny —
**never** a row index, **never** a store ID (anti-pattern 10). A row with no title cannot exist:
`Z-08` requires one, and a provider row that arrives without a name is **dropped by the sync and
logged**, never rendered with a placeholder identity. Truncation is `…` inside the padded field,
so an Ambiguous ellipsis costs one more character and never shears the column.

### (b) Cursor-following scroll — the selection is ALWAYS visible

The viewport follows the cursor. Cursor and scroll offset are preserved **by game identity, not
by index**, across a row-set rebuild (`04-navigation-and-focus.md` §4.1 rule 5), so returning
from `Z-03` after a sync puts the player on the same game, not at row 1. If the focused game did
not survive the rebuild, focus moves to the nearest surviving row — below first, above if there
is nothing below. **Focus is never nowhere.**

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   412 games · 331 not started · 18 in progress · 9 zerado · 54 abandoned       │
│                                                                                │
│     STATE           TITLE                                        HOURS   SRC   │
│     ◉  ZERADO       Yakuza 0                                       96h   STM   │
│     ○  NOT STARTED  Yoku's Island Express                           0h   STM   │
│     ◐  IN PROGRESS  Yooka-Laylee                                    4h   STM   │
│     ⊘  ABANDONED    Ys VIII: Lacrimosa of DANA                     11h   STM   │
│     ○  NOT STARTED  Zelda II: The Adventure of Link                  —   PHY   │
│     ◉  ZERADO       Zeno Clash                                      6h   STM   │
│     ○  NOT STARTED  Zeno Clash II                                   0h   STM   │
│     ⊘  ABANDONED    Zero Escape: Virtue's Last Reward              19h   STM   │
│     ◉  ZERADO       Zero Escape: Zero Time Dilemma                 23h   STM   │
│     ○  NOT STARTED  Zombie Army 4: Dead War                         0h   STM   │
│     ◐  IN PROGRESS  Zoo Tycoon: Ultimate Animal Collection          2h   STM   │
│   ▌ ◉  ZERADO       Zuma's Revenge                                 31h   STM   │
│   ROWS  401–412 of 412                                                         │
│   ↑↓ move   ⏎ open   s status   / filter   a add   r sync   ? help   q quit    │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

331 + 18 + 9 + 54 = **412**. The counts sum to the number shown, at any row count
(`05-state-machine.md` §7 rule 1).

**And they sum to the number shown when rows are absent, which is the case rule 1 is easiest to
break in.** With 412 shown and 20 tombstoned, the summary reads
`412 games  ○ 331  ◐ 18  ◉ 9  ⊘ 54  20 absent` — **412**, not 432. The four counts are a
`GROUP BY` over the **shown** set and the clause names the rest; a total of 432 above a
412-row ledger would break the leg this section exists to prove. See §10.3, D-04-8.

### (c) The summary is pinned OUTSIDE the scroll region

Body row 1 — or row 2 when a banner is present — never inside the viewport, never scrollable,
on screen at 12 rows and at 412. **`Z-06`'s overlay is vertically centred inside `BodyRect`
rather than inside the terminal precisely so it never covers this row** — see
`Z-06-set-status.md` §4.

### 12.1 · Row order — stated so it is not picked arbitrarily

**Phase 1 renders exactly one order: title, A → Z.** Deterministic, stable across syncs (which
is what makes (b)'s identity-preserving cursor work), and it never lies. `last_played_at` is
nullable and *"`NULL` means not reported, not never played"* (`05-state-machine.md` §6), so a
recency sort would have a large undefined tail. **No sort control is bound in Phase 1**, and no
sort indicator renders — an on-screen indicator would imply a control that does not exist
(anti-pattern 14). `07-offline-contract.md` §2 now lists **`Ordering`** rather than *Sort*, and
describes it in the same terms.

### 12.2 · Reading the not-fetched / known-empty distinction in the ledger

The playtime column carries the same distinction `Z-05` spells out in words:

| Fact | Renders | Meaning |
|---|---|---|
| Playtime reported, zero | `0h` in `--z-text-secondary` | The provider answered: nothing played |
| Playtime reported, non-zero | `41h` in `--z-text-secondary` | |
| **Not reported** — the provider has no such capability, or has not been asked | `—` in `--z-text-tertiary` (interim: uncoloured) | **Nothing to fetch, or nothing fetched yet — and only one of those is actionable** |

A hand-added copy is never shown `0h`. That would be a lie: `physical` is a provider with
`Capabilities.Progress = false` (`05-state-machine.md` §2.1), so there is no zero to report.

---

## 13 · The focus model, and how `Esc` behaves

### 13.1 · Regions

| Tier | `R` | Regions | `Tab` |
|---|---|---|---|
| Tiny · Narrow · Standard · Wide | **1** | the list | **does not exist** and is not listed in the footer |
| ExtraWide `120+` | **2** | list · detail pane | `Tab` / `Shift-Tab` move between them |

### 13.2 · How focus is shown — three channels, any two sufficient

| Channel | Focused row | Unfocused row |
|---|---|---|
| **Position** | `▌` U+258C in the 2-cell gutter field (ASCII `>`) | two spaces |
| **Weight** | **bold** | normal |
| **Colour** | `--z-primary` amber `#FFB000` / 214 / `bright yellow` on the marker | none |

**The row cursor is amber, not cyan** — so the ledger spends **zero** cyan on focus and the
budget's rarest colour stays reserved for completion. Never by background fill
(`01-design-system.md` §1.3): at the 16-colour floor `--z-surface` and `--z-surface-raised` both
collapse to `black`, so a highlight bar would vanish. **Never removed**, in any state, for any
reason.

### 13.3 · Selection vs focus at ExtraWide (D-04-6)

When focus is in the detail pane, the ledger's marked row is a **selection**, not the focus.

| | Focus in the ledger | Focus in the pane |
|---|---|---|
| Ledger gutter | `▌` in `--z-primary` amber | `▌` in `--z-text-secondary` chrome |
| Ledger row weight | **bold** | normal |
| Pane border | `┌─┐` `--z-border` | **`┏━┓`** `--z-border-strong` |
| Pane title | `--z-text-secondary` | `--z-primary` amber |

Under `NO_COLOR` the amber/chrome distinction is gone, and **two channels still carry it**: the
ledger row's weight, and the pane's **box-drawing weight** (`04-navigation-and-focus.md` §4.2).
No new glyph is introduced to solve this — the canon already had the mechanism.

### 13.4 · `Esc`, exhaustively, on this screen

| Context | `Esc` does | Then |
|---|---|---|
| Root, at rest, no filter, no overlay | **Nothing** | The footer hints `? help   q quit` |
| An overlay (`Z-06`) is open | Dismiss it, discarding any uncommitted choice | Focus returns exactly where it was |
| Filter mode, editor focused | Leave the **editor**; the filter **stays applied** | The footer changes to `esc clear filter` |
| Filter mode, editor blurred, filter applied | **Clear** the filter | The full library returns |
| Focus in the detail pane (ExtraWide) | Return focus to the ledger | The pane stays rendered — it is not a route |

This is `04-navigation-and-focus.md` §5, unchanged. **There is no keyboard trap on this screen**
(WCAG 2.1.2) and **single-key shortcuts do not fire while the filter editor holds focus**
(2.1.4) — typing `a` types `a`, it does not open `Z-08`.

---

## 14 · Colour budget declaration

**Counted by `02-colour-budget.md` §3.1, from the ANSI stream, classified by payload.**

| Screen state | STATE cyan (uncounted) | Focus ring (exempt) | **CHROME CYAN — the budget** | Verdict |
|---|---|---|---|---|
| **L4 populated** | every `◉` and `ZERADO` — 3 chips in §3's render, 40 in a library with 40 finished games | none — the row cursor is **amber** | **0** | **PASS.** Zero is a pass; a browsing screen has nothing to urge |
| **L1 first run** | none | none | **1** — `c  connect a store` | **PASS** |
| **L2 empty from provider** | none | none | **1** — `r  try again` | **PASS** |
| **L5 banner** | as L4 | none | **0** — the banner's action key is **amber**, the action-required class (`01-design-system.md` §12.2) | **PASS** |
| **L6 filter mode** | as L4 | the text input's ring, **exempt** by §2.3 | **0** | **PASS** — see `Z-07` §14 |

**Amber allow-list entries used** (`02-colour-budget.md` §4.1): 1 the screen title · 2 readout
labels (`STATE`, `TITLE`, `HOURS`, `SRC`, `ROWS`) · 3 the `IN PROGRESS` state · 5 key hints in
the status bar and the empty state, **except the one that spends the cyan** · 7 the
action-required degrade banner · 8 the filter mode sigil `/`. Plus the row cursor `▌` and the
audio annunciator. **Items 4 and 6 are not used** — no progress bar and no terminal mark render
here.

**Amber ceiling:** at 80 × 24 = 1920 cells the ceiling is **192**. §3's render spends
7 (`LIBRARY`) + 34 (column headers) + 4 (`ROWS`) + 3 × 14 (`IN PROGRESS` chips) + 1 (cursor)
≈ **88** — well under. A screen approaching the ceiling has stopped using amber as a voice.

**Red:** **none, anywhere on this screen, in any state.** No scanner (there is no indeterminate
wait here), no destructive confirmation (`ABANDONED` is reversible and must never raise one —
`01-design-system.md` §13.1), no error text. **`OFFLINE` is chrome, never red** — colouring a
publicly promised behaviour red would call it a fault.

**`--z-border`** appears exactly once, as the unfocused detail pane's border at ExtraWide, where
it bounds **no control**. Every control boundary uses `--z-border-strong` (4.08, WCAG 1.4.11).

**No region is separated by fill.** At 16 colours the ledger and the pane are still two regions:
a border and a two-column gutter, both of which survive.

---

## 15 · `NO_COLOR` — rendered, not asserted

With `NO_COLOR` set, Zerado emits **zero SGR sequences**. Not a reduced palette, not a "safe"
subset. The render below is the §3 screen with colour stripped, character for character.

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
│     STATE           TITLE                                        HOURS   SRC   │
│     ◐  IN PROGRESS  Baldur's Gate 3                                87h   STM   │
│     ○  NOT STARTED  Blasphemous                                     0h   STM   │
│     ◉  ZERADO       Celeste                                        14h   STM   │
│     ○  NOT STARTED  Chrono Trigger                                   —   PHY   │
│     ◐  IN PROGRESS  Dark Souls III                                 63h   STM   │
│     ⊘  ABANDONED    Disco Elysium                                  22h   STM   │
│     ◉  ZERADO       Hades                                          58h   STM   │
│     ○  NOT STARTED  Hollow Knight                                  41h   STM   │
│     ◐  IN PROGRESS  Outer Wilds                                    12h   STM   │
│   ▌ ◉  ZERADO       Return of the Obra Dinn                         9h   STM   │
│     ⊘  ABANDONED    Sekiro: Shadows Die Twice                       3h   STM   │
│     ○  NOT STARTED  Sid Meier's Civilization VI: Gathering St…      0h   STM   │
│   ROWS  4–15 of 247                                                            │
│   ↑↓ move   ⏎ open   s status   / filter   a add   r sync   ? help   q quit    │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**The point is that nothing changed.** Channel by channel:

| Information | Carried without colour by |
|---|---|
| The four states | The **glyph** `○ ◐ ◉ ⊘` **and the label word** — the label *is* the text alternative; a TTY has no accessibility API |
| The focused row | The `▌` gutter marker (position) + **bold** (weight) — two of the three channels, which is the point of having three |
| The screen title | UPPERCASE + bold, alone on its row with respiro above and below |
| Hierarchy | Case, weight, box drawing and spacing — the grid, exactly as `01-design-system.md` §1.1 intends |
| Which region has focus (ExtraWide) | `┏━┓` against `┌─┐` — box-drawing **weight** |
| The banner class | The **label word** — `OFFLINE`, `SYNC INCOMPLETE`. Colour distinguishes only informational from action, and the word already says which |
| Audio | `▮ AUDIO` / `▯ MUTED` — filled vs hollow glyph **and** the word |
| **An absent row** | It is not on screen at all, and the summary's **`3 absent`** clause says how many are not. **A word and a number, with no colour and no glyph in it** — the one channel on this screen that never depended on either (§10.3, D-04-8) |

**Cross-check (`02-colour-budget.md` §3.3):** run the same screen with `NO_COLOR=1`. **If any
information disappears, the screen was encoding meaning in colour.** Nothing above does.

**`ZERADO_ASCII=1`** substitutes the entirely-narrow ASCII column and marker:

```
  [~]  IN PROGRESS  Baldur's Gate 3                                87h   STM
  [ ]  NOT STARTED  Blasphemous                                     0h   STM
  [*]  ZERADO       Celeste                                        14h   STM
> [*]  ZERADO       Return of the Obra Dinn                         9h   STM
```

---

## 16 · Reuse verdict, per element

| Element | Verdict | Why |
|---|---|---|
| Header band | **Build fresh** against #2435 · `lipgloss` | ~40 lines behind a single `Frame` wrapper enforced at the router, so a screen **cannot render frameless by construction** |
| **The ledger** | **Build the Zerado ledger primitive, once** · `bubbles/viewport` + `lipgloss` | `bubbles/table` is the primitive that, hand-rolled per screen inside FlowForge, independently dropped the title, the scroll and the pinned footer on two different screens. `LedgerTable` is **not importable** — `internal/`, inside a module that does not resolve anonymously (`00-design-brief.md` §7). Build one primitive that is correct-by-construction on R-10 (a)(b)(c); **every list screen uses it and no screen hand-rolls a table** |
| State chip | Build fresh · `lipgloss` + the width-aware pad | No `bubbles` primitive fits and none should be forced |
| Status bar / pinned summary | Build fresh · `lipgloss` join + width-aware truncation | |
| Detail pane (ExtraWide) | See `Z-05-game-detail.md` §16 | One view, two hosts — built once, mounted twice |
| Filter bar | **`bubbles/textinput` — direct fit** | See `Z-07-filter-and-search.md` §16 |
| Empty state | Build fresh · `lipgloss` text layout | Trivial, and worth doing well |
| Degrade banner | Build fresh · one `lipgloss` row | |
| Scroll-position row | Build fresh · one `lipgloss` row | |
| Progress bar / scanner | **Neither renders on this screen** | No indeterminate or determinate wait exists here (§8) |
| Error rendering | **Not `charmbracelet/log`** | Its role is structured developer logging; its level colours are its own palette, not Zerado's |
| Width measurement | An East-Asian-Width-aware function, **everywhere** | Never `len()`, never `utf8.RuneCountInString` — truncation, padding and centring alike |

**Zerado cannot import a single line of FlowForge Go code.** FlowForge canon is inherited as
**specification**; the live reuse target is the Charm ecosystem, which is public and is the
ratified stack.

---

## 17 · Upstream findings — contradictions and stale rows

Recorded rather than silently designed around. None is mine to fix; each belongs to the document
named.

**Re-checked against head on 2026-08-25.** Eight findings were recorded in rev A; **seven are
closed** and are struck from this table rather than left to rot — `01-design-system.md` §4.3
(row budget), §4.4 (ExtraWide title **40**), §5.4 (the summary is body row 1, *"it is not the
footer row"*), §5.3 (the three-facts guideline reconciled, *"where they genuinely conflict, the
falsifiable rule wins"*), §5.11 of the manual (audio, struck through and marked SUPERSEDED),
`07-offline-contract.md` §4.1 (Phase 1 copy names the library, the price wording deferred to
Phase 3 with its reason) and §2 (*Sort* → **`Ordering`**, described honestly). They are
enumerated as **#7–#12** and **#4** in `14-contradictions-closed.md`.

| # | Finding | Where | Owner |
|---|---|---|---|
| 1 | **Two documents cite different Unicode versions for the same width table.** `01-design-system.md` §1.2 is checked against `EastAsianWidth-17.0.0.txt`; `02-composition.md` §2.2.1 cites Unicode 16.0. This spec re-verified every glyph it uses against **UCD 16.0.0** on 2026-08-25 and every class matched both documents, so nothing is wrong today — but a version that is stated twice and differently is a version nobody owns | `01-design-system.md` §1.2 vs `02-composition.md` §2.2.1 | `fft-design-architect` |
| 2 | **`01-design-system.md` §6.1 splits ExtraWide as `ledger 64 · gutter 2 · pane 46`; `02-composition.md` §2.3 splits it `66 ∥ 2 ∥ 44`.** Both sum to 112, which is the harder kind of disagreement — a builder following one and a reviewer following the other would each think the screen correct. **This spec follows the spine (66 ∥ 2 ∥ 44), which is binding**, and §4's 40-column ExtraWide title is derived from it. `14-contradictions-closed.md` #15 records this as closed with *"real pane is 44 wide, 38 content"*, so the design system is the document that has not caught up | `01-design-system.md` §6.1 | `fft-design-architect` |

---

## 18 · Open for the founder

> **Closed since rev A.** The `s` collision — *set status* on a populated library, *connect
> Steam* on an empty one — was item 1 here. `04-navigation-and-focus.md` §3.2 adopted the
> recommendation: **`s` means set-status everywhere, always; connect-a-store is `c`.**
> `01-design-system.md` §10.1 and §10.3 carry the corrected copy and this spec reproduces it
> verbatim (§8.3, §8.4, §10.1, §10.2). Nothing is open.

1. **Nine ANSI-256 indices are still underived** (`00-design-brief.md` §9). This screen ships
   `--z-text-tertiary` **uncoloured** as its documented interim — correct and honest, but not
   the designed colour. Confirm the derivation lands before `Z-04` is built.
2. **The library deck as Zerado's in-house gold standard.** `00-design-brief.md` §4 (R-1)
   proposes that this screen becomes the reference every later screen rises to, the moment it is
   ruled GOLDEN. Confirm the nomination; it changes how the next eight screens are reviewed.

---

## 19 · Design decisions made in this spec

Marked, with reasons, because the canon did not settle them.

| # | Decision | Reason |
|---|---|---|
| **D-04-1** | The pinned summary spells the state words; the glyph key is the degrade, not the default (§10.3) | The colour budget forbids colouring the summary, so a glyph there would be a one-channel co-render; prose casing is already ratified for this row |
| **D-04-2** | The pinned chrome block grows downward; the respiro is never spent; the scroll region absorbs (§8.2) | One rule that closes the arithmetic in all four compositions at the 16-row floor, and never sacrifices the row that separates chrome from ledger |
| **D-04-3** | The *zerado* result line replaces the summary until the next keypress, and only for a transition **into** `zerado` (§10.5) | No timer, no dwell constant, no WCAG 2.2.1 time limit; and the moment the product is named after should not share its weight with routine edits |
| **D-04-4** | At ExtraWide the ledger is unbordered; only the detail pane carries a border (§4) | A ledger border would push content 3 columns right of the title and break header-left == content-left, the canon's headline invariant |
| **D-04-5** | Below Standard the scroll region takes `floor(remaining/2)` games; the leftover row becomes the respiro (§11.1) | Always closes the arithmetic with no orphan half-row |
| **D-04-6** | Source is **right-aligned** in its 4-column field; selection and focus are distinguished by the ledger row's weight plus the pane's box weight (§3.2, §13.3) | Every ledger row then ends flush with the body's right edge, so the block reads as ruled rather than ragged; and the focus/selection distinction survives `NO_COLOR` using only vocabulary the canon already has |
| **D-04-7** | Row order is title A→Z, no sort control, no sort indicator (§12.1) | A list must have an order and naming it is honest; an on-screen indicator would imply a control that does not exist |
| **D-04-8** | When at least one row is absent the summary appends `<n> absent`; the clause is exempt from the degrade ladder and the base form degrades around it (§8 L13, §10.3) | `06-data-seams.md` §2.4 excludes absent rows from the default view, so the summary is describing a subset — and `05-state-machine.md` §7 requires it to say which set it describes. A banner would cost a game row and claim an urgency the seam refuses; the scroll row does not render when nothing scrolls. The summary is the only always-present carrier |
| **D-A1** | The audio annunciator is right-aligned on the header band's title row, on every screen (§7.1) | The only surface every framed screen has at every tier; keeps 4.1.3's stable place with no per-screen exception, without repurposing the footer or truncating the summary |

---

## 20 · Screen-specific acceptance criteria

Beyond `00-design-brief.md` §10 (28 lines) and `02-colour-budget.md` §10 (19 boxes), both of
which apply in full. Each line below is falsifiable from a rendered artifact.

1. **Every game row measures exactly `BodyRect.w` cells** at Standard and above, on a terminal
   configured `ambiguous-width=single` **and** on one configured `double`. Verified by measuring
   an ANSI-stripped render, not by counting runes.
2. **Header-left equals content-left**, by column number, at every tier: column 4 at Wide,
   column 2 at Narrow, column 1 at Tiny, column 5 at ExtraWide.
3. **The pinned summary is on screen at 412 rows**, scrolled to the last row, at 80 × 24 and at
   32 × 24.
4. **The selection is visible after scrolling to row 380** and **after a row-set rebuild** — the
   cursor lands on the same *game*, not the same index.
5. **Every row carries a human title.** No index, no store ID, no empty identity cell — at 412
   rows, at every tier.
6. **The four state counts sum to the total shown**, in every summary form and with a filter
   active.
7. **Chrome-cyan count is 0** on the populated screen and **exactly 1** on each empty state, by
   the §3.1 machine method.
8. **`NO_COLOR=1` loses no information** — same characters, same meaning, zero SGR.
9. **At the 16-colour floor the ledger and the detail pane are still two regions.** No fill
   separates anything.
10. **The degrade banner never hides a row's state, its title, the summary, the focus marker or
    the footer** — it costs a game row and nothing else.
11. **The footer lists no key that does nothing here**, and lists no key at all while a text
    input holds focus except `esc`, `⏎`, `tab` and `^c`.
12. **Below 24 × 8 the screen refuses in one sentence** and does not render a frame.
13. **No scanner and no progress bar appear on this screen in any state.**
14. **No absent row appears in the default view**, at 412 rows with 20 tombstoned, at every
    tier — verified by row count and by title against the file.
15. **The summary's four state counts sum to the number it shows**, with absent rows present,
    in every one of the five composed forms of §10.3.
16. **`<n> absent` is on screen whenever `n ≥ 1`** at 80 × 24, 60 × 24, 40 × 24 and 32 × 24,
    and is **byte-absent** when `n = 0`. It is the last clause standing at 32 columns.
17. **`absent_since` clearing produces no banner, no result line and no announcement** — the
    clause's count decrements, and at zero the clause stops rendering.
18. **Audio, when enabled, is announced on the title row and in the footer**; when it has never
    been enabled, neither renders and `m` is unbound.
19. **The screenshot is founder-validated before merge.** Eight artifacts per
    `03-responsive.md` §7: `24×8` · `32×24` · `40×24` · `60×24` · **`80×24`** · `120×40`, plus
    `NO_COLOR=1` at `80×24` and a forced 16-colour depth at `80×24` — **and every tabular one at
    412 rows.** No founder-validated screenshot → not GOLDEN → no merge.
