---
title: Zerado — Landing Page Design Blueprint
discipline: DESIGN
doc-no: ZRD-DESIGN-01
rev: C
date: 2026-08-24
amended: 2026-08-25
status: normative — the build implements this; the quality gate diffs against it
machine-manifest: design/blueprint.tokens.json
---

# Zerado — landing page design blueprint

**This document is normative.** Where it and an improvisation disagree, this document wins.
Where it and `brand/brand-manual.md` disagree, the brand manual wins — this blueprint *applies*
brand tokens, it never invents colour or type. Where it and `ratification/mock.outline.md`
disagree on structure, the ratified outline wins.

The machine-readable half of this contract is `design/blueprint.tokens.json`. Both were written
together; if they ever drift, the JSON is the one the gate reads and this file is the one a human
reads, and the drift is a defect in this deliverable.

---

## 1 · Conformance header

```
ratifies:  ratification/mock.outline.md
sha256:    a2bc87e434f4035a7dd851dacadb8aee0961d66f5e58465b18f9640b20248ecb   (verified 2026-08-24)
applies:   brand/tokens.css · brand/tokens.json · brand/brand-manual.md   (rev A)
wires:     content/landing-copy.md (rev A, "final — implement verbatim")
head data: content/seo.md (title, meta description, canonical, OG, Twitter) — used verbatim
sections:  16 present · 16 ratified · 0 added · 0 dropped · 0 reordered
```

### 1.0 Amendment record

§9 names the things the build may not change without coming back to this
document. One of them changed. This is that return — the procedure, followed,
rather than a silent edit to the clause that forbade it.

#### Amendment 1 — 2026-08-25 · roadmap status · rev A → rev B

**What changed.** §12's Phase 1 status marker moves from `Planned` to
`In progress`. Phases 2, 3 and 4 are unchanged.

**Authority.** Founder instruction, 2026-08-25, verbatim: *"we need to show
Phase 1 'In progress'"*. Phase 1 work started that day, so the statement became
true and the page was permitted to make it.

**Clauses amended.** §2 (§12) — the ASCII spec and the "all four read Planned"
condition. §8 — the co-render inventory. §9 item 5 — the change-control clause
itself.

**What did NOT change, and remains under §9.**

- No phase may render a **done**-equivalent. `StatusMarker`'s type is
  `'planned' | 'in-progress'`; there is no third value, and adding one is a
  deliberate act requiring another amendment.
- The roadmap stays **undated**.
- Every mockup keeps its **visible caption**.
- The status marker still co-renders **colour + glyph + label**, never colour
  alone.

**Enforcement.** `scripts/check-page.mjs` pins the exact sequence
`['In progress', 'Planned', 'Planned', 'Planned']` and fails CI on any drift, in
either direction. This document and that gate now agree; before this amendment
they did not, which is the defect this record exists to close.

#### Amendment 2 — 2026-08-25 · the money position · rev B → rev C

**What changed.** §11 states the Phase 4 community layer is **donation-supported**
rather than needing "a premium account or a donation", and §16's footer band
states that Zerado **earns nothing** rather than disclosing an affiliate
commission.

**Authority.** Two founder instructions, 2026-08-25. On the premium tier:
*"let's not have premium account for now and only donation, so we're not
commercial app."* On the affiliate model, conditionally: *"if is a must to fit
the non-commercial, drop it, and keep only the donation."* The condition is met —
the metadata provider Phase 2 depends on scopes its free tier by whether the
**project generates revenue**, and an affiliate commission is exactly that.

**Clauses amended.** §2 (§11) — the "nearest link is the affiliate disclosure"
note. §2 (§16) — the disclosure band's contents. §8 — the cover-art rationale,
whose first reason cited the page being commercial.

**What did NOT change, and remains under §9.**

- **Zero funding controls.** No donate button, no sponsor button, no funding CTA.
  Ratified decision Q6's "disclosure is not an ask" is untouched and still
  governs; the page *states* the position and *invites* nothing.
- **§07 price intelligence is untouched** — all-time low, current price and the
  wait-or-buy verdict all render exactly as ratified. Only an affiliate tag on an
  outbound link was dropped, and the blueprint never specified that tag.
- The footer's **disclosure band is retained** as a structural element.

**Enforcement.** `scripts/check-page.mjs` asserts the page makes no revenue claim
and that §11 states donation support, and continues to assert zero funding
controls. All three fail CI on regression.

### 1.1 Section inventory — all sixteen, in the ratified order

| # | id | Ratified heading source | Layout family |
|---|---|---|---|
| 01 | `hero` | H1 from copy §01 | Copy column + full-width terminal frame |
| 02 | `maker-line` | *(no heading — one line)* | Single centred rule-flanked line |
| 03 | `the-problem` | "The problem" (copy §03, optional heading — **rendered**) | Three-line stack, narrow measure |
| 04 | `moods` | "Sorted by mood, not genre" | Intro + 4-card grid + terminal frame |
| 05 | `states` | "Where each game stands" | Intro + 4 state chips in a row |
| 06 | `one-collection` | "One collection, everything you own" | Two text halves + 12-tile cover grid |
| 07 | `price-intelligence` | "When to buy, not just what to buy" | Intro + single price card, centred |
| 08 | `yours-and-offline` | "It's yours, and it's offline" | Three trust points |
| 09 | `terminal-first` | "Built for the terminal, on purpose" | Body + game-detail terminal frame |
| 10 | `and-on-your-phone` | "And on your phone, later" | Body + two phone frames |
| 11 | `community` | "The community layer — Phase 4, not built yet" | Phase badge + body + muted example panel |
| 12 | `roadmap` | "Roadmap" | Horizontal 4-phase track |
| 13 | `after-phase-4` | "After Phase 4 — ideas, not promises" | Quiet single block, dashed frame |
| 14 | `faq` | "Questions" | Six `<details>` rows |
| 15 | `closing-cta` | *(no heading — line + button)* | Centred single-action block |
| 16 | `footer` | *(no heading)* | Three-band footer |

**Order is DOM order.** The build emits these sixteen `<section>` elements as immediate children
of `<main>` (with §16 in `<footer>`), in exactly this sequence, each carrying `id` equal to the
ratified id above. That id is the conformance handle the gate reads.

### 1.2 Global constraints this design commits to

| Constraint | Value | Why it is in the design, not just the build |
|---|---|---|
| **JavaScript** | **none** | Everything on this page is achievable without it: two links, a `mailto:`, a native `<details>` accordion, CSS animation. Zero JS makes INP effectively free and removes a whole class of failure. It is a design decision because it forbids any component that would need script. |
| **Raster images** | **zero** | No `.png`/`.jpg`/`.webp` anywhere. Every "image" on this page is composed HTML/CSS or an inline/linked SVG that already exists in `site/public/`. |
| **External requests** | **zero** | Fonts, logos and favicon are same-origin. No CDN, no analytics, no embed, no web font service. |
| **Horizontal page scroll** | **forbidden at every breakpoint** | `<body>` never scrolls sideways. Anything wider than its container scrolls *inside its own box*. See §5.5. |
| **Scroll-triggered animation** | **none** | Reveal-on-scroll costs CLS risk and reads as 2016 marketing. The page is composed to be good when static. |
| **Colour scheme** | dark-native, not conditional | Per `tokens.css`: the dark palette is `:root` and there is no light mode on this page. |

---

## 2 · Section-by-section layout spec

Notation: `≥1280` / `768–1279` / `<768` are the design's three reflow bands; the four *measured*
viewports are 375 / 768 / 1280 / 1920 (§5). Copy references are to `content/landing-copy.md`
by its section number — the build pulls strings from there and never retypes them.

Grid: a 12-column grid inside `--z-container-max` (1180px) with `--z-gutter` (24px). At 1280 a
single column measures 76.33px. `--z-container-narrow` (760px) is the prose measure column.

### §01 `hero`

```
≥1280
┌──────────────────────────────────────────────────────────────────────┐
│  [nav: lockup ······································ GitHub · CTA ]  │
│                                                                      │
│  cols 1–7 ────────────────────────────────┐    cols 8–12: empty      │
│   PICK A MOOD.                            │    (negative space is    │
│   GET A GAME.                             │     the composition —    │
│   PLAY IT TONIGHT.                        │     see note below)      │
│   <subhead, 19px, max 68ch>               │                          │
│   [ Join the waitlist ] [ View on GitHub ]│                          │
│   <microcopy 15px>                        │                          │
│  ───────────────────────────────────────────────── grid horizon ─────│
│  ┌── cols 1–12 · TerminalFrame TF-1 "library" · two-pane ─────────┐   │
│  │ zerado — library                        412 games · 6 zerado   │   │
│  │ ▓▓▓▓ scanner track (parked pip)                                │   │
│  │ ◐ Hollow Knight    STEAM  42h IN PROGRESS │ TONIGHT            │   │
│  │ ◉ Celeste          STEAM  11h ZERADO      │  ○ NOT STARTED 288 │   │
│  │ ○ Disco Elysium    GOG     0h NOT STARTED │  ◐ IN PROGRESS   9 │   │
│  │ ⊘ Dark Souls III   STEAM   6h ABANDONED   │  ◉ ZERADO        6 │   │
│  │ ◐ Hades            STEAM  28h IN PROGRESS │  ⊘ ABANDONED   109 │   │
│  │ ○ Chrono Trigger   PHYSICAL — NOT STARTED │                    │   │
│  │ ↑↓ move · enter open · m mood · / filter · q quit              │   │
│  └────────────────────────────────────────────────────────────────┘   │
│  <figcaption — the mockup disclosure, verbatim from copy §01 alt>     │
│  ═══════════════ ScannerRule (the page's one animation) ══════════════│
└──────────────────────────────────────────────────────────────────────┘
```

- **Copy column is cols 1–7 and stops there.** Columns 8–12 stay empty. This is deliberate: an
  expensive object is mostly surface. Filling that space with a badge cluster or a second
  graphic is what makes a page look like a template.
- **The image slot is full width and below the copy**, exactly as copy §01 specifies
  ("full width, below the fold line"). It is not a side-by-side split — the ratified copy
  already made that call.
- **Backdrop:** `GridHorizon` (§7.4) occupies the hero's bottom 38%, masked to zero opacity
  above that line. **No text is ever composited over it** — the terminal frame that overlaps it
  is an opaque `--z-bg` surface with its own border.
- **LCP element is the H1 text.** Nothing above it loads an image, because there are no images.

| Band | Behaviour |
|---|---|
| ≥1280 | As drawn. Frame two-pane, 125 mono columns at 15px (80 list / 3 divider / 42 rail). |
| 768–1279 | Copy column spans full container. Frame single-pane, list only (the TONIGHT rail is dropped, not shrunk), 87 columns at 13px. |
| <768 | Everything full width. Frame rows go **two-line** (§7.3) so the state label is never dropped — the co-render rule (colour + glyph + label) survives 375. 39 columns at 13px. |

**Content slots** — H1 copy §01 headline, rendered as three block-level `<span>`s breaking after
each sentence (measured: this is the only wrap that fits 375 — see §4.2); subhead copy §01;
`CTAButton primary` → `mailto:alex@flowforgesoft.com?subject=Zerado-WaitList`; `CTAButton
secondary` → `https://github.com/JustCode-CruzAlex/Zerado`; microcopy copy §01; figure caption =
copy §01 alt text, verbatim and visible.

### §02 `maker-line`

One line, copy §02, at `--z-text-lg` (19px), `--z-text-secondary`, centred in
`--z-container-narrow`, flanked left and right by a 1px `--z-border` hairline that runs to the
container edge. Vertical rhythm `--z-space-7` (48px) top and bottom rather than the full
`--z-section-y` — this is a credential, not a section, and it should feel like a caption between
two big moments. `Zerado` and `FlowForge` render in `--z-text` (16.65:1) so the two product names
lift out of the sentence without needing colour.

All bands identical; only the horizontal padding changes.

### §03 `the-problem`

Heading "The problem" rendered as the section H2 but at the **small** end of the H2 clamp
(§4.3 `.z-h2--quiet`), because the copy marks it "optional, small". Below it, three lines from
copy §03 as a `ProblemStack`: each line its own block at `--z-text-lg` (19px), 1.6 line-height,
separated by `--z-space-3`, in `--z-text`. The three lines are a crescendo, so line 3 is the only
one that wraps at any breakpoint — set the measure to `--z-container-narrow` and let it.

**A single amber vertical rule** (`2px`, `--z-primary`, full height of the stack) sits at the
left edge with `--z-space-5` of space before the text. That is the whole treatment. No imagery,
per the ratified outline.

All bands: stacked, narrow measure. At <768 the amber rule stays (it costs 2px + 24px = 26px of
a 335px budget and it is the section's only ornament).

### §04 `moods`

H2 copy §04, intro line copy §04 at 19px, then:

- **MoodCard × 4** — grid of 4 at ≥1280, **2×2 at 768–1279**, 1 column at <768. Each card:
  `--z-surface-raised` fill, 1px `--z-border-strong` edge, `--z-radius-md` (4px), padding
  `--z-space-5`. Title at `--z-text-h3` (23px, Space Grotesk 700, `--z-text`), body at
  `--z-text-base` (17px, `--z-text-secondary` — 8.41:1 on raised). Titles and bodies verbatim
  from copy §04 cards 1–4, in order.
- A **mono index** (`01`–`04`, `--z-text-readout`, `--z-text-tertiary`) sits top-right of each
  card. Cockpit annotation, not decoration — it is how a terminal labels a choice list.
- **Cards are not interactive** and carry no hover state. See §6.1.
- **TerminalFrame TF-2 "mood picker"**, `still` variant, below the grid, full container width.

### §05 `states`

H2 copy §05, intro line, then **StateChip × 4 in one row** at ≥768 and a vertical list at <768.
Each chip is `glyph + LABEL + line`, drawn straight from the brand's state quad
(`brand/tokens.json` `color.state`):

| # | Colour token | Glyph | Label | Line source |
|---|---|---|---|---|
| 1 | `--z-state-not-started` `#A5A29B` | `○` | NOT STARTED | copy §05 state 1 |
| 2 | `--z-state-in-progress` `#FFB000` | `◐` | IN PROGRESS | copy §05 state 2 |
| 3 | `--z-state-zerado` `#19E0FF` | `◉` | ZERADO | copy §05 state 3 |
| 4 | `--z-state-abandoned` `#C77DFF` | `⊘` | ABANDONED | copy §05 state 4 |

**State 3 is the payoff and is composed differently on purpose.** Its chip spans wider (it holds
a long explanatory line), gets `--z-glow-cyan` as the page's *only* glow, and its glyph renders
one step larger. This is the moment the product name explains itself; the composition should
notice. The italicised first appearance of *zerado* required by `content/landing-copy.md`'s
casing rule lives in this line.

At <768 the four chips stack; chip 3 keeps the glow and stays last-but-one in reading order — the
ratified order is Not Started · In Progress · Zerado · Abandoned and it does not change.

### §06 `one-collection`

H2 copy §06, then two halves, then the tile grid.

- **Halves:** ≥1280 two columns 50/50 (`connected stores` | `physical shelf`), 768–1279 two
  columns 50/50, <768 stacked. Each half gets a mono `readout` overline. The overlines are the
  ratified half-names from the outline (`CONNECTED STORES`, `THE PHYSICAL SHELF`) rendered
  uppercase in `--z-text-readout`.
- **StoreRow list** inside the first half — Steam marked live, PlayStation/GOG/EA marked planned.
  Live = `◉` + `--z-accent`; planned = `○` + `--z-text-tertiary` (5.53:1 on raised).
  `--z-cyan-900` was rejected for the planned rows: 2.87:1, fails as text (§8.3).
- **CoverGrid** — 12 `CoverTile`s. Full spec in §7.5. Columns: **3 at 375, 4 at 768, 6 at
  ≥1280.** 12 divides cleanly into all three (4 rows / 3 rows / 2 rows) so the shelf never ends
  on a ragged row at any measured viewport.
- Figure caption: copy §06 alt text verbatim, **plus** one design-authored line —
  `Cover tiles are illustrative artwork, not real game covers.` (flagged in §7.5; content phase
  may re-word it, the requirement stands).

### §07 `price-intelligence`

H2 copy §07, intro line, then **one PriceCard, centred, max 560px.** This is the section the
outline calls "the one a skeptical visitor screenshots", so it gets the most compositional care
per pixel of any component on the page.

```
┌─ PRICE CHECK ─────────────────────────────┐   ← mono readout overline
│                                           │
│   $15          TODAY                      │   ← 40px JetBrains Mono 700, --z-primary
│   ───────────────────────────────────     │
│   $8           ALL-TIME LOW · 3 MONTHS AGO│   ← 23px mono, --z-steel
│                                           │
│   ▁▂▃▂▁▁▂▅▇▆▄▂▁▁▂▃   sparkline (CSS bars) │
│                                           │
│   "$15, and in budget. It's been $8.      │   ← the verdict, 19px, --z-primary
│    Maybe wait."                           │
└───────────────────────────────────────────┘
 Illustrative numbers. Price data comes from IsThereAnyDeal once Zerado
 is connected — nothing above is a live quote.        ← caption, copy §07 verbatim
```

- Every figure comes from copy §07's worked example — `$15`, `$8`, `three months ago`, and the
  verdict string verbatim. The design introduces no number the copy did not already state.
- **Verdict colour is `--z-primary` amber (9.53:1 on raised), not red.** `--z-scanner-500`
  measures 4.17:1 on `--z-surface-overlay` and fails as text, and the brand reserves it for
  motion and alarm. "Maybe wait" is a readout, not an alarm. Rejection recorded in §8.3.
- **Sparkline** is 16 CSS bars, `--z-amber-900` for the field with the two extremes lit
  (`--z-primary` for today, `--z-steel` for the low). Decorative, `aria-hidden`, carries no
  information not already in text.
- All bands: identical, centred, card max 560px; at 375 the card is full width (335px) and the
  `$15`/`$8` figures step from 40px/23px to 34px/19px mono.

### §08 `yours-and-offline`

H2 copy §08, then **TrustPoint × 3**: 3 columns at ≥1280, 3 columns at 768–1279 (they are short),
1 column at <768. Each point's first sentence is its own line at `--z-text-h3` (23px,
`--z-text`), the remainder at 17px `--z-text-secondary` — copy §08 supplies both halves in one
paragraph each and the split is at the first period.

Each point is prefixed by a 2px × 24px amber tick rule rather than an icon. **No icon set is
introduced anywhere on this page** (§7.7).

### §09 `terminal-first`

H2 copy §09, body copy §09 at 19px in `--z-container-narrow`, then **TerminalFrame TF-3
"game detail"**, `still`, full container width, with its caption.

At ≥1280 the body sits in cols 1–6 and **cols 8–12 stay empty**, same as the hero. Echoing the
frame's keybar into that space was considered and rejected — it would duplicate the frame's own
footer two inches away. Two sections on this page use deliberate empty right-hand space, and they
are the two sections that talk about the terminal; that is a rhyme, not a repetition.

### §10 `and-on-your-phone`

H2 copy §10, body copy §10, then **PhoneFrame × 2**, side by side and slightly overlapped
(the rear frame offset +32px x, +24px y, behind), centred.

- Frames are CSS only: 1px `--z-border-strong` outline, `--z-surface` fill, a 3px `--z-void`
  bezel inset, and a 28px corner radius. **28px is outside the brand's radius scale on purpose
  and it is declared as an exception** (`--z-phone-radius`): the scale is machined UI chrome, and
  this element depicts *hardware*, not chrome. A 4px-radius phone reads as a bug.
- Inside each: a 5-row list echoing the library — glyph + title + state label, at 13px mono. Same
  co-render rule.
- Both frames carry a `PhaseBadge` "PHASE 4" in the corner.
- Caption: copy §10 alt text, verbatim and visible.
- ≥768: side by side, overlapped. <768: **the rear frame is dropped entirely** (not shrunk) and
  one frame renders at 240px wide, centred. Two overlapping frames at 335px is mush.

### §11 `community`

The section is stamped three separate ways, because the ratified decision makes this
non-negotiable and one stamp is one thing to miss:

1. **The H2 says it** — copy §11's heading is used verbatim, with the clause "— Phase 4, not
   built yet" rendered in `--z-primary` amber inside the same heading element. Same text,
   emphasised.
2. **A `PhaseBadge` chip** sits above the H2: `⊘ PHASE 4 · NOT AVAILABLE YET`, mono readout,
   amber on `--z-surface-raised` (9.53:1), 1px amber border.
3. **The example panel is the `muted` TerminalFrame variant** — contents at 55% opacity over
   `--z-surface`, with a repeating diagonal watermark reading `PHASE 4` in `--z-chrome-500` at
   0.14 alpha, `aria-hidden`, decorative only. Because the watermark is low-contrast decoration,
   *it is never the only carrier* — stamps 1 and 2 both pass AA on their own.

Body copy §11 verbatim, including the premium-or-donation sentence. **There is no button, no
link, no price and no tier name in this section** — `blueprint.tokens.json` records
`"interactiveElements": 0` for §11 so the gate can assert it. The nearest link is in the footer
and it is the footer's money-disclosure text, not an ask. (Amended 2026-08-25 — §1.0,
Amendment 2: §11 states donation support, not "premium account or donation".)

Panel contents (three mini-panels: a comment thread, a review, a public profile), 2–3 lines each,
generic handles (`@rafa`, `@kim`, `@sol`). ≥768 three across; <768 one column, all three kept.

### §12 `roadmap` — the customer's required section

H2 "Roadmap", intro line copy §12, then the **RoadmapTrack**.

```
≥1280 — one row, four across, left to right
  ●───────────●───────────●───────────●        ← 2px --z-border-strong rail,
  PHASE 1     PHASE 2     PHASE 3     PHASE 4     nodes 10px --z-steel rings
  CLI/TUI MVP Enrichment  Recommend…  Social…
  Your libr…  Covers, s…  What to b…  Sync, c…
  ◐ IN PROG…  ○ PLANNED   ○ PLANNED   ○ PLANNED
```

> Phase 1 reads `◐ IN PROGRESS` in `--z-state-in-progress` amber as of
> **Amendment 1** (§1.0). It read `○ PLANNED` at rev A.

| Band | Layout | Why |
|---|---|---|
| ≥1280 | 4 across, one row, horizontal rail | The ratified shape. 277px per card. |
| 768–1279 | **2 × 2 grid, reading order preserved left-to-right then down** | Measured: 4 across at 768 gives 168px per card, and "Recommendations" alone is 132px at 15px — the cards cramp. 2×2 keeps left-to-right time order (the whole reason the outline chose horizontal) without cramping. |
| <768 | Vertical list, rail rotated to a left-hand vertical spine | The ratified fallback, verbatim: "stacking to a vertical list on narrow screens". |

- **Phase 1 reads `In progress`; phases 2–4 read `Planned`. None is marked done** — and no
  done-equivalent status exists to mark one with. Copy §12 says so and the outline makes it the
  honesty condition of the whole page. Amended 2026-08-25 (§1.0, Amendment 1); at rev A all four
  read `Planned`.
- **No dates render** — decisions.md leaves dates open and the fallback is undated.
- **The status marker is built so that changing one phase is a token change, not a redesign.**
  `StatusMarker` consumes exactly four component tokens:

  ```css
  .z-phase { --z-status-color: var(--z-steel);      /* #A5A29B — 6.85:1 on raised */
             --z-status-track: var(--z-amber-900);  /* unlit; non-text only      */
             --z-status-glyph: "○";
             --z-status-label: "PLANNED"; }
  ```

  Marking Phase 1 shipped later means setting that one block to the brand's `zerado` state quad
  (`--z-cyan-500` / `◉` / `SHIPPED`) on `.z-phase[data-phase="1"]`. Colour, glyph and label all
  move together, so the co-render rule holds by construction and no layout changes.
- **`--z-amber-900` (#8A5E00) was rejected as the PLANNED label colour** — 3.06:1 on raised,
  fails body text. It survives as the *track* fill, which carries no text. Recorded in §8.3.

### §13 `after-phase-4`

H2 copy §13, with "— ideas, not promises" in `--z-text-secondary` inside the heading (the same
emphasis mechanic as §11, tuned the other way: §11 lifts its qualifier, §13 lowers it).

Body copy §13 in `--z-container-narrow`, inside a block with a **1px dashed `--z-border-strong`
frame and no fill** — visually the quietest block on the page, and dashed because dashed reads as
provisional in every technical drawing convention there is. No glow, no amber, no card
elevation. It must not be mistakable for the roadmap directly above it.

All bands identical; narrow measure.

### §14 `faq`

H2 "Questions", then **six `FAQItem`s** using native `<details>`/`<summary>` — no JavaScript.

- Question at `--z-text-lg` (19px, `--z-text`, Space Grotesk 600); answer at 17px
  `--z-text-secondary` (9.36:1), measure capped at `--z-measure` (68ch).
- Row separated by a 1px `--z-border` hairline (decorative use — between content blocks — which
  is the one thing `--z-border` is permitted to do).
- Marker is a mono `+` / `−` in `--z-primary`, right-aligned, rotating **not at all** (no
  rotation; the glyph swaps). Mechanical, per brand §7.2.
- **First item closed.** All six closed by default: an open first item implies the others are
  less important, and the questions are of equal weight.
- Questions and answers verbatim from copy §14, in the copy's order.
- All bands identical. At <768 the summary wraps to two lines and the marker stays on line 1.

### §15 `closing-cta`

Line above the button (copy §15) at `--z-text-lg`, `--z-text`, centred, max 560px. One
`CTAButton primary` → the same `mailto:` as §01, same label. Microcopy (copy §15) beneath at
15px `--z-text-tertiary`.

**No secondary CTA here.** The hero offers both; the close offers one. Splitting attention at the
last ask is the standard way to lose it, and the ratified outline says "the same button,
restated. No new promise."

Background: `--z-surface`, with a **static** scanner track (parked pip, no animation) as the
section's top edge — the visual echo of the hero's live rule, at rest. One animation on the page
means every other scanner is furniture.

### §16 `footer`

Three bands inside `<footer role="contentinfo">`, separated by `--z-border` hairlines:

1. **Links band** — `Join the waitlist` · `View on GitHub` (copy §16's repeated links, no new
   strings), plus the contact line (copy §16) with `alex@flowforgesoft.com` as a `mailto:` link
   with no pre-filled subject (it is a contact address, not the waitlist).
2. **Disclosure band** — the money disclosure, copy §16 verbatim, at 15px `--z-text-secondary`.
   This is the page's only funding statement outside §11, and neither is a button.
   (Amended 2026-08-25 — §1.0, Amendment 2: this band carried an affiliate disclosure
   at rev A. The affiliate model is dropped; the band now states that Zerado earns
   nothing. The band itself is retained — its removal would leave a structural gap.)
3. **Company band** — "A FlowForgeSoft product." plus `PoweredByFlowForge`:
   `site/public/flowforge-logo.svg` at 28px (its viewBox is `0 0 36 36`) with the caption
   "Powered by FlowForge." at 13px `--z-text-tertiary`.

The Zerado lockup does **not** repeat in the footer — it is already in the nav, and repeating it
next to the FlowForge mark creates a lockup-adjacency problem the brand's clear-space rule (§3.3)
would then govern. One mark per page, plus the maker's mark.

≥768: three bands as horizontal rows with links left / company right. <768: everything stacks,
links first.

### 2.1 Every link on the page

There are exactly **three** destinations and **nine** link instances. Together with the six FAQ
disclosures that is **fifteen focusable elements on the whole page**. Nothing else is clickable.

| Instance | Destination | Verified |
|---|---|---|
| Nav CTA | `mailto:alex@flowforgesoft.com?subject=Zerado-WaitList` | decisions.md Q2 |
| Nav secondary | `https://github.com/JustCode-CruzAlex/Zerado` | decisions.md: public, HTTP 200 anonymously |
| Hero CTA primary | `mailto:…?subject=Zerado-WaitList` | " |
| Hero CTA secondary | `https://github.com/JustCode-CruzAlex/Zerado` | " |
| §15 CTA | `mailto:…?subject=Zerado-WaitList` | " |
| Footer waitlist | `mailto:…?subject=Zerado-WaitList` | " |
| Footer GitHub | `https://github.com/JustCode-CruzAlex/Zerado` | " |
| Footer contact | `mailto:alex@flowforgesoft.com` (no subject) | copy §16 |
| Skip link | `#main` | element exists in this blueprint |

**No `href="#"`, no unbuilt routes, no anchor to a section that does not exist, no download link,
no social icons, no donate or sponsor link anywhere.** External links carry
`rel="noopener noreferrer"`; the GitHub link opens in the same tab (a new tab is a decision made
for the user, and this audience does not need it made for them).

---

## 3 · Applied brand tokens — colour role → surface mapping

Values are resolved from `brand/tokens.css`. **This blueprint applies; brand owns.** Ratios are
measured with the WCAG 2.2 relative-luminance formula; the method reproduces all sixteen figures
published in `brand-manual.md` §4.2 exactly (verified 2026-08-24), so the extended pairs below
are on the same footing as the manual's own.

### 3.1 Surfaces

| Surface | Token | Value | Used by |
|---|---|---|---|
| Page ground | `--z-bg` | `#05060A` | `<body>`, hero, §09, §13, terminal-frame screens |
| Section surface | `--z-surface` | `#0B0D14` | §02 §04 §06 §08 §12 §14 §15, footer |
| Raised | `--z-surface-raised` | `#141A24` | MoodCard, StateChip, PriceCard, PhaseCard, PhaseBadge, CoverTile plate ground |
| Overlay | `--z-surface-overlay` | `#1D2532` | hovered/focused card fill, terminal selected row |

Sections alternate `--z-bg` / `--z-surface` rather than using dividers everywhere; the difference
is 1.5% luminance, which is exactly the "one machined object, milled in steps" read and nowhere
near a stripe.

### 3.2 Text on surface — the full applied matrix

| Foreground | `#05060A` | `#0B0D14` | `#141A24` | `#1D2532` |
|---|---|---|---|---|
| `--z-text` `#E9EEF5` | 17.37 | **16.65** | 14.97 | 13.22 |
| `--z-text-secondary` `#A9B5C7` | 9.76 | **9.36** | 8.41 | 7.43 |
| `--z-text-tertiary` `#8492A8` | 6.42 | **6.15** | 5.53 | **4.89** ← worst in system |
| `--z-primary` `#FFB000` | **11.05** | **10.59** | 9.53 | 8.41 |
| `--z-primary-hover` `#FFC94D` | 13.23 | **12.68** | 11.40 | 10.06 |
| `--z-accent` / link / focus `#19E0FF` | **12.68** | **12.15** | 10.93 | 9.65 |
| `--z-steel` `#A5A29B` | 7.95 | **7.62** | 6.85 | 6.05 |
| `--z-orchid-500` `#C77DFF` | 7.53 | **7.21** | 6.49 | 5.73 |
| `--z-scanner-300` `#FF6B6B` | 7.30 | **6.99** | 6.29 | 5.55 |

Bold values are the figures published in `brand-manual.md` §4.2 and are cited, not recomputed.
**Every text pair used anywhere on this page is ≥ 4.89:1. There is no AA-marginal text on this
page and no AAA-only-if-large text.**

### 3.3 Role → surface, resolved

| Role on the page | Token | Value |
|---|---|---|
| Body prose | `--z-text` | `#E9EEF5` |
| Supporting prose, FAQ answers, subhead | `--z-text-secondary` | `#A9B5C7` |
| Captions, microcopy, mono meta, planned-store rows | `--z-text-tertiary` | `#8492A8` |
| Headings H1/H2 | `--z-text` on dark; the emphasised clause in `--z-primary` | `#E9EEF5` / `#FFB000` |
| Readout overlines, section eyebrows | `--z-text-tertiary` | `#8492A8` |
| Primary action fill | `--z-accent` | `#19E0FF` |
| Primary action label | `--z-accent-contrast` | `#05060A` (12.68:1) |
| Secondary action label | `--z-text` | `#E9EEF5` |
| Secondary action edge | `--z-border-strong` | `#64748B` (4.08:1, passes 1.4.11) |
| Focus ring | `--z-focus-ring` | `#19E0FF`, 2px, 2px offset — **never removed** |
| Card/control edge | `--z-border-strong` | `#64748B` |
| Content hairline | `--z-border` | `#2A3342` — decorative, **between content blocks only** |
| Terminal amber readout | `--z-primary` | `#FFB000` |
| Scanner track / pip | `--z-scanner-track` / `--z-scanner` | `#5C1414` / `#FF2E2E` |
| Grid horizon lines | `--z-grid-line` | `#1A2740` — decorative, never behind text |

**Why the primary CTA is cyan and the headings are amber.** `tokens.css` names cyan "ACCENT —
completion, the primary CTA" and amber "PRIMARY BRAND COLOUR", the ambient voice. Applied: amber
is everywhere the machine is *talking* (headings, readouts, price figures, in-progress state);
cyan appears exactly where something is *earned or completed* — the ZERADO state, the links, the
focus ring, and the one action the page wants. On a page of amber readouts the single cyan button
is the brightest, rarest thing on screen. That is the KITT cockpit argument, and it is why cyan
does not get spent on decoration anywhere else.

### 3.4 Colour budget — the discipline that keeps it from becoming a light show

| Colour | Permitted instances on the page |
|---|---|
| `--z-accent` cyan as a **fill** | 4 (three CTA buttons + the §05 `ZERADO` chip glyph) |
| `--z-glow-cyan` | **1** (the §05 ZERADO chip) |
| `--z-glow-amber` | **1 at a time** (hover/focus on one control) — per brand §4.6, never stacked |
| `--z-scanner` animated | **1** (the hero's bottom ScannerRule) |
| `--z-scanner` static/parked | 4 (three terminal frame chrome bars + §15's top edge) |
| `--z-orchid-500` | 2 (the `ABANDONED` chip, and abandoned rows inside frames) |

---

## 4 · Type scale

Families are `--z-font-display` (Orbitron), `--z-font-body` (Space Grotesk), `--z-font-mono`
(JetBrains Mono), loaded by the already-staged `site/public/fonts.css` — 5 woff2 files,
same-origin, correct `unicode-range` per subset, `font-display: swap`.

### 4.1 The ramp as applied

| Step | Size | Line-height | Family | Weight | Tracking | Case | Where |
|---|---|---|---|---|---|---|---|
| H1 (hero) | **clamp 30 → 52px** | 1.08 | Orbitron | 700 | 0.03em | UPPER | §01 only |
| H2 | **clamp 24 → 34px** | 1.22 | Orbitron | 600 | 0.02em | UPPER | §03–§14 headings |
| H3 | 23px | 1.22 | Space Grotesk | 700 | 0 | Sentence | MoodCard titles, TrustPoint leads, price figures' labels |
| Body large | 19px | 1.60 | Space Grotesk | 400 | 0 | Sentence | subhead, intros, §03 lines, FAQ questions (600) |
| **Body** | **17px** | **1.65** | Space Grotesk | 400 | 0 | Sentence | all prose — the floor |
| Body small | 15px | 1.60 | Space Grotesk | 400 | 0 | Sentence | microcopy, footer disclosure |
| Caption | 13px | 1.50 | Space Grotesk | 400 | 0 | Sentence | figure captions |
| Readout | 13px | 1.20 | JetBrains Mono | 500 | 0.18em | UPPER | eyebrows, overlines, tags, status labels |
| Terminal | 13 / 15px | 1.45 | JetBrains Mono | 400–500 | 0 | as-is | inside TerminalFrame (§7.3) |
| Price figure | 34 / 40px | 1.0 | JetBrains Mono | 700 | 0 | — | §07 only |

`--z-text-display` (76px) is **not used on this page.** Nothing here earns it, and reaching for
the biggest step because it exists is how a page starts shouting.

### 4.2 The Orbitron rules, and the one measurement that governs the hero

The brand's hard rule (§6.1): **Orbitron is display-only — never below 23px, never for prose,
never more than about eight words, always uppercase, always tracked out.** Applied here:

- Orbitron appears in exactly two places: the H1 and the twelve H2s. Nothing else, ever.
- Smallest Orbitron on the page is **24px** (the H2 clamp floor) — above the 23px floor.
- H1 is 8 words; the longest H2 is 8 words. Both at the stated limit, neither over it.
- No Orbitron paragraph, no Orbitron button, no Orbitron in a terminal frame.

**The hero H1 clamp is a measured value, not a taste value.** Measured from the actual staged
`site/public/fonts/Orbitron-latin.woff2` at `wght 700` with 0.03em tracking:

| Line | @34px | @30px | 375 budget |
|---|---|---|---|
| `PICK A MOOD.` | 274.6px | 242.3px | 335px |
| `GET A GAME.` | 262.5px | 231.6px | 335px |
| `PLAY IT TONIGHT.` | **359.6px — overflows** | **317.3px — fits** | 335px |

The brand's generic `--z-text-h1-fluid` bottoms out at 34px, and at 34px *this specific headline*
blows through a 375px viewport by 24.6px. So the hero gets its own clamp, still far above the
23px Orbitron floor:

```css
--z-h1-hero: clamp(1.875rem, 1.3053rem + 2.4309vw, 3.25rem);
/* resolves: 375 → 30.0px · 768 → 39.6px · 1280 → 52.0px · 1920 → 52.0px (capped) */
```

**And the H1 renders as three block lines**, breaking after each sentence:

```html
<h1><span>Pick a mood.</span><span>Get a game.</span><span>Play it tonight.</span></h1>
```

Three reasons, in order of weight: (1) it is the only wrap that fits every measured viewport —
natural wrapping at 1280 produces the orphan `TONIGHT.` on line 3; (2) the copy is three beats and
should be typeset as three beats; (3) it is deterministic, so the gate can diff it. Screen readers
read the heading as one string; the spans are `display:block` and add no semantics.

`--z-text-h2` clamp `clamp(1.5rem, 1.1rem + 1.7vw, 2.125rem)` resolves 24 / 30.7 / 34 / 34px.
Measured worst case at 24px: the longest single word on the page (`COLLECTION,`) is 188.8px
against a 335px budget. No H2 word overflows at any breakpoint.

### 4.3 Weight usage

Maximum three weights per family per surface, per brand §6.3:

- **Orbitron:** 700 (H1), 600 (H2). Two weights. 900 is never used.
- **Space Grotesk:** 400 (prose), 600 (FAQ questions, CTA labels), 700 (H3). Three weights.
- **JetBrains Mono:** 400 (terminal body), 500 (readouts, labels), 700 (price figures). Three.

Measured: JetBrains Mono's advance is exactly **0.600em** — 7.8px per cell at 13px, 9.0px at 15px.
Every terminal column budget in §7.3 derives from that number.

---

## 5 · Responsive breakpoints — 375 / 768 / 1280 / 1920

These four are the design's breakpoints **and** the quality gate's Playwright viewports. They are
identical by requirement and must not diverge.

```css
/* mobile-first; three media queries, no more */
@media (min-width: 768px)  { … }
@media (min-width: 1280px) { … }
@media (min-width: 1920px) { … }   /* container is already capped; this band only widens padding */
```

### 5.1 Per-breakpoint frame

| | **375** | **768** | **1280** | **1920** |
|---|---|---|---|---|
| Page padding `--z-pad-x` | 20px | 24px | 24px | 24px |
| Content width | 335px | 720px | 1180px (capped) | 1180px (capped) |
| Grid columns | 1 | 2 | 12 | 12 |
| Section rhythm `--z-section-y` | 64px | 80px | 96px | 96px |
| Nav | lockup 120px + primary CTA; **no hamburger** | lockup 140px + both links | same | same |
| Terminal mono | 13px / 39 cols | 13px / 87 cols | 15px / 125 cols | 15px / 125 cols |
| Cover grid | 3 cols | 4 cols | 6 cols | 6 cols |
| Roadmap | vertical list | 2 × 2 | 4 across | 4 across |
| Mood cards | 1 col | 2 × 2 | 4 across | 4 across |
| State chips | stacked | 4 across | 4 across | 4 across |
| Phone frames | 1 frame, 240px | 2, overlapped | 2, overlapped | 2, overlapped |

**1920 is not a new layout.** The container caps at 1180px and centres; the extra 740px is page
ground. A design that keeps growing to fill a 1920 monitor stops being a designed object and
becomes a stretched one. The only 1920-specific behaviour is that the hero's GridHorizon runs
full-bleed to the viewport edge (it is a backdrop, not content), which is the one place the extra
width buys something.

### 5.2 Nav behaviour — the no-hamburger decision

Two links do not earn a menu. At 375 the nav carries the lockup at its **measured 120px minimum**
(`brand-manual.md` §3.4) plus the primary CTA, and the GitHub link drops to the footer, which
already repeats both links per copy §16. Measured fit:

```
lockup 120 + clear-space gap 16 + CTA 171  =  307px   ≤  335px budget      ✓
```

The 16px gap is not arbitrary: the lockup's clear space equals its cap height, which at 120px wide
(viewBox `0 0 364 64`, cap height 48 units) is **15.8px**. 16px satisfies it with 0.2px to spare;
nav vertical padding is 20px, which satisfies it above and below.

A hamburger would also cost the page its zero-JavaScript property for two links. Not worth it.

### 5.3 375 — the hardest viewport, checked element by element

| Element | Risk | Resolution | Measured |
|---|---|---|---|
| H1 | Orbitron overflow | own clamp @30px + 3 explicit lines | 317.3px ≤ 335 ✓ |
| H2 (longest word `COLLECTION,`) | overflow | H2 clamp floor 24px | 188.8px ≤ 335 ✓ |
| Nav | crowding | lockup 120 + CTA, GitHub → footer | 307px ≤ 335 ✓ |
| Terminal frames | horizontal overflow | 39-column budget + two-line rows | see §7.3 ✓ |
| Roadmap track | 4-across cramping | vertical list with left spine | ✓ |
| Cover grid | tile too small for its tag | 3 cols → 106px tile; `PHYSICAL` tag = 62.4px + 12px plate padding | 74px ≤ 106 ✓ |
| Phone frames | two overlapped frames | rear frame dropped, one at 240px | ✓ |
| Price card | 40px mono figures | step to 34px mono | ✓ |
| CTA pair | side-by-side buttons | **stack full-width**, 48px min height each | ✓ |
| FAQ summary | marker collides with wrapped text | marker in its own 32px grid column | ✓ |

### 5.4 Touch targets

Every interactive element is **≥ 48 × 48px** at every breakpoint (WCAG 2.5.8 asks 24×24; 48 is the
comfortable floor, and with only fifteen focusable elements on the page there is no reason to be
tight). CTA buttons: 48px tall at ≥768, 52px at 375. FAQ `<summary>`: 56px minimum row. Footer
links: 44px line-box with 16px spacing between adjacent targets.

### 5.5 The horizontal-overflow contract

**`<body>` never scrolls sideways at 375, 768, 1280 or 1920.** Three mechanisms, in order of
preference:

1. **Reflow** — cards, chips, roadmap, cover grid, phone frames all change column count. First
   choice everywhere it is possible.
2. **Reformat** — terminal-frame rows become two-line at <768 rather than shrinking or dropping
   the state label. The co-render rule is preserved; the layout changes shape instead.
3. **Contained scroll** — as a *backstop only*, `TerminalFrame` sets `overflow-x: auto` on its
   screen element with `overscroll-behavior-x: contain`. If a future frame is authored wider than
   its budget, it scrolls inside its own box and the page still does not. The scroll container is
   keyboard-focusable (`tabindex="0"`) with a visible focus ring, so a keyboard user can reach it —
   a scrollable region that traps its content from keyboard users is a WCAG 2.1.1 failure.

Build assertion the gate can run: at each of the four viewports,
per-element overflow: no element's bounding rect may exceed the viewport (`right <= clientWidth + 1`, `left >= -1`), excluding `aria-hidden` subtrees.

> **Do NOT use `document.documentElement.scrollWidth === document.documentElement.clientWidth`.** `overflow-x: hidden` on the root pins those two values together, so that comparison is a tautology that cannot fail. It reported PASS at all four breakpoints while §02 was destroying 58% of its own text (QA BLOCKING-1). The per-element check above was mutation-tested: green at HEAD, red with the defect re-injected.

---

## 6 · Component inventory

Thirty components. Each row lists its states and the tokens it consumes. **`focus-visible` is
load-bearing** — the brand states the focus ring is never removed, and a keyboard-first audience
is this product's launch audience.

### 6.1 The interaction rule that governs the whole table

**Only interactive elements get hover states.** MoodCards, StateChips, CoverTiles, PhaseCards and
TrustPoints do not lift, glow or change on hover, because they are not clickable and a hover
affordance on a non-clickable element is a small lie that costs a click of trust. This page has
exactly **seven interactive component types** — enumerated in §6.2 — across **fifteen focusable
elements**, and they are the only things on the page that respond to a pointer.

### 6.2 Interactive components

| Component | default | hover | focus-visible | active | Tokens |
|---|---|---|---|---|---|
| **SkipLink** | visually hidden, first in DOM | — | slides to top-left, `--z-accent` fill, `--z-accent-contrast` label, 2px ring | — | `--z-accent`, `--z-accent-contrast`, `--z-focus-ring` |
| **CTAButton `primary`** | `--z-accent` fill, `--z-accent-contrast` label (12.68:1), `--z-radius-md`, 48px tall, Space Grotesk 600 16px, no shadow | fill → `--z-accent-hover` `#8CF0FF` (label 15.45:1), `--z-glow-cyan`, 140ms `--z-ease-standard` | 2px `--z-focus-ring` at 2px offset over page ground (ring vs both neighbours = 12.68:1) — ring is **additive to**, never a replacement for, the hover treatment | translateY(0), fill → `--z-accent`, 90ms | `--z-accent`, `--z-accent-hover`, `--z-accent-contrast`, `--z-glow-cyan`, `--z-radius-md`, `--z-duration-fast` |
| **CTAButton `secondary`** | transparent fill, 1px `--z-border-strong` edge, `--z-text` label (16.65:1) | edge → `--z-primary`, label → `--z-primary` (10.59:1), `--z-glow-amber` | same 2px cyan ring, 2px offset | edge → `--z-primary-hover` | `--z-border-strong`, `--z-primary`, `--z-glow-amber` |
| **NavLink** (GitHub) | `--z-text-secondary` 15px | → `--z-text`, 1px `--z-primary` underline offset 4px | 2px cyan ring | — | `--z-text-secondary`, `--z-primary` |
| **FAQItem** (`details`/`summary`) | closed, question `--z-text` 19px/600, `+` marker `--z-primary`, 1px `--z-border` hairline below | question → `--z-primary`, marker → `--z-primary-hover` | 2px cyan ring around the whole summary row | — | `--z-border`, `--z-primary`, `--z-duration-base` open / `--z-duration-fast` close |
| **FooterLink** | `--z-text-link` `#19E0FF` (12.15:1), no underline | underline 1px, colour → `--z-text-link-hover` `#8CF0FF` (14.81:1) | 2px cyan ring, 2px offset | — | `--z-text-link`, `--z-text-link-hover` |
| **TerminalFrame screen** (scroll backstop only) | `tabindex="0"` only when its content overflows | — | 2px cyan ring inset 2px | — | `--z-focus-ring` |

Focus-ring geometry is constant everywhere: `outline: 2px solid var(--z-focus-ring);
outline-offset: 2px`. On the cyan-filled primary button the 2px offset gap exposes the page ground
between button and ring, so the indicator contrasts 12.68:1 against **both** its neighbours — which
is what WCAG 2.4.13 actually asks. There is no case on this page where the ring sits directly
against a colour it cannot be told apart from.

### 6.3 Non-interactive components

| Component | States | Tokens |
|---|---|---|
| **Nav** | static (not sticky — a fixed bar steals 72px of a 667px phone screen for two links) | `--z-surface`, `--z-border` |
| **Logotype** | one — `site/public/logo.svg`, 140px ≥768 / 120px at 375, clear space = cap height (15.8px @120, 18.5px @140). Never re-set in type, never glowed, never scaled part-wise (brand §3.7) | — |
| **SectionEyebrow** | one — mono readout, zero-padded ratified index + ratified id, e.g. `04 · MOODS` | `--z-text-readout`, `--z-text-tertiary` |
| **SectionHeading** | one, with an optional emphasised trailing clause (`--z-primary` in §11, `--z-text-secondary` in §13) | `--z-font-display`, `--z-text`, `--z-primary` |
| **Hero** | one | `--z-bg`, `--z-h1-hero` |
| **GridHorizon** | static (decorative, `aria-hidden`) | `--z-grid-line`, `--z-primary` @6% |
| **ScannerRule** | `live` (hero only) · `parked` (reduced-motion, and all other instances) | `--z-scanner`, `--z-scanner-track`, `--z-glow-red`, `--z-duration-scanner`, `--z-ease-scanner` |
| **TerminalFrame** | `live` · `still` · `muted` | `--z-bg`, `--z-border-strong`, `--z-font-mono`, state quad |
| **TerminalChrome** | one (title + right-aligned counter + parked scanner track) | `--z-text-tertiary`, `--z-scanner-track` |
| **TerminalRow** | `default` · `selected` (`--z-surface-overlay` fill + `▸` caret) — **presentational only, not focusable** | state quad, `--z-surface-overlay` |
| **MockupCaption** | one | `--z-text-tertiary`, 13px |
| **MakerLine** | one | `--z-text-secondary`, `--z-border` |
| **ProblemStack** | one | `--z-text`, `--z-primary` (rule) |
| **MoodCard** | one (non-interactive) | `--z-surface-raised`, `--z-border-strong`, `--z-text`, `--z-text-secondary`, `--z-text-readout` |
| **StateChip** | four variants, one per state; variant 3 additionally `emphasised` | full state quad + `--z-glow-cyan` (variant 3 only) |
| **StateLegend** | one | — |
| **StoreRow** | `live` · `planned` | `--z-accent`, `--z-text-tertiary` |
| **CoverTile** | twelve fixed variants (§7.5); `physical` is a modifier of one of them | `--z-grid-line`, `--z-primary`, `--z-accent`, `--z-orchid-500`, `--z-steel`, `--z-surface` |
| **PlatformTag** | `store` (steel) · `physical` (amber + `▣`) | `--z-steel`, `--z-primary`, `--z-surface` |
| **CoverGrid** | one (`role="img"`, labelled by its caption) | — |
| **PriceCard** | one | `--z-surface-raised`, `--z-primary`, `--z-steel`, `--z-amber-900` |
| **Sparkline** | static, `aria-hidden` | `--z-amber-900`, `--z-primary`, `--z-steel` |
| **TrustPoint** | one | `--z-text`, `--z-text-secondary`, `--z-primary` |
| **PhoneFrame** | one; `--z-phone-radius: 28px` declared exception | `--z-surface`, `--z-void`, `--z-border-strong` |
| **PhaseBadge** | one | `--z-primary`, `--z-surface-raised` |
| **CommunityPanel** | one (`muted`) | `--z-chrome-500` @0.14 watermark, `--z-surface` |
| **RoadmapTrack** | `4up` · `2x2` · `vertical` | `--z-border-strong` (rail), `--z-steel` (nodes) |
| **PhaseCard** | one | `--z-surface-raised`, `--z-border-strong` |
| **StatusMarker** | driven entirely by `--z-status-color` / `-track` / `-glyph` / `-label`; ships as `PLANNED` on all four | `--z-steel`, `--z-amber-900` |
| **SpeculationBlock** | one | dashed `--z-border-strong`, `--z-text-secondary` |
| **ClosingCTA** | one | `--z-surface`, `--z-scanner-track` (parked) |
| **Footer** | one | `--z-surface`, `--z-border` |
| **PoweredByFlowForge** | one | `site/public/flowforge-logo.svg` @28px, `--z-text-tertiary` |

---

## 7 · Imagery & iconography direction

The page ships **zero raster images**. Everything below is HTML, CSS, or an SVG already staged in
`site/public/`. That is not a compromise forced by an empty asset folder — it is the
right answer for this product, and each of the three decisions below explains why.

### 7.1 Decision one — the terminal views are rendered in HTML/CSS, and every one carries a visible mockup caption

**The problem.** Zerado is not runnable by a visitor yet. There is no screenshot, and
`decisions.md` records that faking one is out of bounds. But a grey placeholder box in the hero
would waste the single most persuasive slot on the page, for an audience that judges a
terminal product by whether its terminal looks right.

**The decision.** All three terminal views are **real HTML/CSS**, amber-on-black, JetBrains Mono,
using the brand's own state colours and glyphs — and each is wrapped in a `<figure>` whose
**visible `<figcaption>` states it is a mockup and not a screenshot.**

**Why this is better than an image even if we had one:**

| | Rendered HTML/CSS | A screenshot PNG |
|---|---|---|
| Bytes | 0 | 80–300 KB, the LCP element |
| Crisp at 1920 and at 375 | yes, it is text | no, it is pixels |
| Selectable / searchable / translatable | yes | no |
| Screen-reader accessible | yes, real text | only via `alt` |
| Honest today | yes, and labelled | no such asset exists |
| Cost to update when the product ships | change a row | re-shoot and re-optimise |

**The caption is load-bearing and non-negotiable.** Without it this is a fake screenshot, which is
precisely what the customer ruled out. And the caption needs no new copy, because
`content/landing-copy.md` already wrote it — each image slot's alt text is exactly the honest
sentence required. So:

> **Rule: every mockup figure's visible `<figcaption>` is that figure's `alt` text from
> `content/landing-copy.md`, verbatim.** The `<figure>` takes `role="img"` and
> `aria-labelledby` pointing at the `<figcaption>`, so a screen reader announces it once rather
> than twice.

For §01 that reads, on screen, in `--z-text-tertiary` at 13px directly under the frame:

> *An illustration of Zerado's terminal library view. This is a mockup of the planned interface,
> not a screenshot — the program isn't runnable by visitors yet.*

Zero authored copy. The disclosure is visible on every mockup on the page, in the content team's
own words.

### 7.2 The line this page draws — names are facts, artwork is property

The terminal frames name real games (`Hollow Knight`, `Celeste`, `Disco Elysium`,
`Dark Souls III`, `Hades`, `Chrono Trigger`) as plain text in a list. The cover grid shows **no
real cover art at all.** That is not an inconsistency; it is the distinction that matters:

- A **title is a fact** about a library, used nominatively, as every storefront, review site and
  backlog tracker does. Naming a game in a list of games is what the product literally is.
- **Cover artwork is a copyrighted work.** Reproducing a dozen of them on a commercial page is a
  different act entirely.

Using real titles also makes the frames credible in a way invented titles never could — and
inventing plausible game titles carries its own collision risk. The six titles chosen span four
states and both acquisition paths (five synced, one `PHYSICAL`), so the mockup proves the
product's claims rather than decorating them.

### 7.3 The three terminal frames, specified

Shared chassis for all three:

```
┌ chrome bar: title (mono 13px --z-text-tertiary) ······· counter (right) ┐
├ parked scanner track: 2px, --z-scanner-track, pip at rest (see §7.6)    ┤
│ screen: --z-bg #05060A, padding 12px @<768 / 16px @768 / 24px @≥1280    │
│ rows …                                                                  │
├ keybar: mono 13px --z-text-tertiary                                     ┤
└ 1px --z-border-strong, --z-radius-md, --z-elevation-2 ──────────────────┘
```

**Column budgets — derived from JetBrains Mono's measured 0.600em advance, not estimated:**

| Viewport | Content | Frame chrome | Interior | Mono | **Columns** |
|---|---|---|---|---|---|
| 375 | 335px | 2px border + 2×12 pad | 309px | 13px (7.80px/cell) | **39** |
| 768 | 720px | 2px + 2×16 | 686px | 13px | **87** |
| 1280 | 1180px | 2px + 2×24 | 1130px | 15px (9.00px/cell) | **125** |
| 1920 | 1180px (capped) | 2px + 2×24 | 1130px | 15px | **125** |

Mono size **steps** (13px below 1280, 15px at and above) rather than fluidly clamping, so the
column budget is an exact integer at every measured viewport and the gate can check it.

**TF-1 — `library`, §01 hero, `live` chassis.**

Row grid, one line at ≥768:
`glyph 2ch │ gap 2 │ title 1fr (min 18ch) │ platform 10ch │ hours 6ch │ state 12ch`

| Glyph | Title | Platform | Hours | State |
|---|---|---|---|---|
| `◐` | Hollow Knight | STEAM | 42h | IN PROGRESS |
| `◉` | Celeste | STEAM | 11h | ZERADO |
| `○` | Disco Elysium | GOG | 0h | NOT STARTED |
| `⊘` | Dark Souls III | STEAM | 6h | ABANDONED |
| `◐` | Hades | STEAM | 28h | IN PROGRESS |
| `○` | Chrono Trigger | PHYSICAL | — | NOT STARTED |

All four states appear, which pre-teaches §05; one row is `PHYSICAL`, which pre-teaches §06. Row 2
(`Celeste`, ZERADO) is the `selected` row — `--z-surface-overlay` fill and a `▸` caret — so the
earned cyan sits at the frame's optical centre.

Chrome: `zerado — library` left, `412 games · 6 zerado` right. Keybar:
`↑↓ move · enter open · m mood · / filter · q quit`.

At ≥1280 the frame is **two-pane**: list 80 cols │ `│` divider 3 cols │ `TONIGHT` rail 42 cols
holding the four state counts (`○ NOT STARTED 288 · ◐ IN PROGRESS 9 · ◉ ZERADO 6 · ⊘ ABANDONED
109` — summing to 412, matching the chrome counter and copy §03's "You own 400 games"
order of magnitude). At 768–1279 the rail is **dropped whole**, not compressed.

At <768 each row becomes **two lines**:

```
◐ Hollow Knight
  IN PROGRESS · STEAM · 42h
```

Measured longest two-line row: `  NOT STARTED · PHYSICAL · —` = 28 cells ≤ 39 ✓. Titles truncate
with `…` as a backstop but none of the six needs it. **The state label never disappears** —
dropping it at narrow widths would break the brand's co-render rule (colour AND glyph AND label)
on the product's most-used component.

**TF-2 — `mood picker`, §04, `still`.** Two columns: mood list │ matches.

```
┌ zerado — pick a mood ────────────────────────── 4 moods ┐
│  ▸ Mindless grind              ◐ Hades                  │
│    Story rich, kind of sad     ◉ Celeste                │
│    Quick fifteen minutes       ○ Chrono Trigger         │
│    Tactical, full focus        ◐ Hollow Knight          │
└ ↑↓ move · enter pick · esc back ────────────────────────┘
```

Mood names verbatim from copy §04's four card titles, in the copy's order. Longest is
`Story rich, kind of sad` = 23 cells; at 39 cols the two columns do not fit side by side, so at
<768 the frame **stacks**: mood list, a `───` rule, then matches. Both parts survive.

**TF-3 — `game detail`, §09, `still`.** One game, one screen:

```
┌ zerado — Hollow Knight ─────────────────────── steam ┐
│  ◐ IN PROGRESS          42h played                   │
│                                                      │
│  MOODS   tactical, full focus · story rich           │
│  PRICE   $15 today   ·   low $8 (3 months ago)       │
│          maybe wait                                  │
│                                                      │
│  SINOPSE  A short plot summary lands here once        │
│           enrichment is wired in.                     │
└ e edit · m mood · p price history · q back ──────────┘
```

Price figures are copy §07's worked example exactly — the page never shows two different example
prices. Mood tags are copy §04's titles. The `SINOPSE` row does double duty: it demonstrates the
detail view and it shows the Portuguese word in situ, which copy §06 explains once. Longest line
is 38 cells — fits 39 at 375 ✓ with no reflow needed.

**Caption on all three:** the figure's own alt text from copy, verbatim (§7.1).

### 7.4 The grid horizon (hero backdrop)

Two stacked CSS layers on the hero's `::before`, both `aria-hidden`, occupying the hero's bottom
38% and masked to zero opacity above that:

1. A perspective plane of `--z-grid-line` `#1A2740` lines — two `repeating-linear-gradient`s,
   `perspective: 320px` + `rotateX(76deg)`, converging to a horizon at the mask boundary.
2. A `radial-gradient` bloom of `--z-primary` at **6% alpha** centred on the horizon, capped at
   10% — the amber phosphor haze.

Verified: even at the 10% cap, the composited ground is `#1E1709` and `--z-text` still measures
**15.24:1** over it, `--z-text-secondary` **8.56:1**. No text is placed there anyway, but the
budget is checked so a later change cannot quietly break it.

**Static. It does not scroll, parallax, or animate.** An animated synthwave grid is the single
fastest way to turn retro-future into retro-kitsch, and it would burn paint budget for it.

### 7.5 Decision two — the cover grid uses art-directed tiles, not real game covers

**The customer asked for "nice game images to illustrate the page." He is not getting a grid of
real game covers, and here is exactly why — and what he is getting instead, which is better.**

**Why not real covers.** Three independent reasons, any one of which would be enough:

1. **Cover art is copyrighted.** Embedding a dozen publishers' artwork is a legal exposure
   the customer would inherit from us, and we do not hand a customer a liability inside a
   deliverable. *(At rev A this reason also cited the page being commercial under an
   affiliate model. That model was dropped on 2026-08-25 — §1.0, Amendment 2 — which
   weakens this particular argument but does not remove it: the artwork is copyrighted
   either way, and reasons 2 and 3 below are untouched.)*
2. **The page is required to be fully self-contained.** No external image requests. Real covers
   would mean either hotlinking a storefront CDN (an external request, and a dependency that
   404s the day someone re-arts a game) or redistributing the files ourselves, which is the same
   copyright problem with extra steps.
3. **`decisions.md` flags IGDB as an open commercial-licensing risk** — cover art and *sinopse*
   depend on a partnership that does not exist yet. Shipping a marketing page built on the exact
   asset class whose licence is unresolved would be building a second problem on top of the first.

**What it is instead, and why it is stronger.** Twelve **art-directed cover tiles**, composed in
CSS in the Zerado palette. A grid of real 2020s marketing covers on an 80s-retro-future page
would fight the page — twelve different studios' art directions, twelve different colour
temperatures, in a palette engineered to be exactly six colours. The tiles carry the *argument*
of the section (shelf rhythm, density, one hand-added item sitting as an equal) without importing
twelve conflicting aesthetics.

**Specification:**

- **Aspect ratio `3 / 4`** — the box-art proportion; it is what makes the grid read as covers
  before a single word is processed.
- **Grid:** 3 columns at 375, 4 at 768, 6 at ≥1280. Gap `--z-space-3` (12px) at 375,
  `--z-space-4` (16px) above. 12 tiles divides evenly into 4 / 3 / 2 rows — the shelf never ends
  ragged.
- **Four treatments × three hue keys = twelve unique tiles.** No two adjacent tiles (in any of
  the three column counts) share both treatment and hue.

| Treatment | Composition | Reads as |
|---|---|---|
| **A · grid horizon** | perspective grid converging to a lit horizon line; hue-keyed glow above it | the synthwave beat, in miniature |
| **B · scanline field** | `repeating-linear-gradient` scanlines over a hue wash, one bright pip off-centre | a CRT holding a signal |
| **C · brushed panel** | a diagonal `--z-steel` band with a machined slot, on `--z-void` | the DeLorean panel |
| **D · phosphor arc** | an off-centre radial disc, hue-keyed, fading into the void | a readout at rest |

Hue keys are `--z-primary` amber, `--z-accent` cyan, `--z-orchid-500` orchid. Fixed sequence,
left to right, top to bottom (also recorded in `blueprint.tokens.json` so the gate can diff it):

```
 1 A-amber   2 B-cyan     3 C-steel    4 D-orchid   5 B-amber    6 A-cyan
 7 D-steel   8 C-orchid   9 B-orchid  10 A-steel   11 D-amber   12 C-cyan
```

- **Every tile carries a `PlatformTag`** at its bottom-left: `STEAM`, `GOG`, `EA`, `PS`, or
  `PHYSICAL`. Distribution: 6 STEAM, 2 GOG, 1 EA, 2 PS, **1 PHYSICAL** (tile 9). The tag is
  **plain text on a solid `--z-surface` plate**, never floating on the artwork — so its contrast
  is fixed at 7.62:1 (steel) regardless of what the tile behind it is doing. That is the whole
  reason for the plate.
- **The `PHYSICAL` tag co-renders** — amber `#FFB000` (10.59:1) + the `▣` glyph + the word, so it
  is distinguishable without colour, exactly like a product state.
- **Tiles carry no titles.** No invented game names, no real ones. The section's argument is
  carried by copy, by the tag distribution, and by density.
- **The grid is `role="img"`** with `aria-labelledby` on its caption, children `aria-hidden`. A
  screen reader gets copy §06's alt text once, not twelve decorative divs.
- Non-interactive: no hover, no link, no lightbox.
- Caption: copy §06 alt text verbatim + one design-authored line —
  `Cover tiles are illustrative artwork, not real game covers.` **Flagged:** this single sentence
  is the one string on the page not drawn from `content/landing-copy.md`. It exists because this
  decision requires disclosure. Content may re-word it; the requirement to disclose stays.

### 7.6 Decision three — motion, and the synthwave question

**The founder wants the feel of 80s synthwave. The page carries that in the visual and motion
register, not in audio.** No audio file, no autoplay, no play button: browsers block autoplay,
sound-on-load is an accessibility failure and a trust failure, and an audio asset would dominate
the byte budget of a page that currently ships zero images. What actually transmits the feeling on
a landing page is the **scanner sweep, the grid horizon, the neon glow and the amber CRT
readout** — and this page has all four.

**Total motion inventory — nine entries, and that is the complete list.**

| # | Motion | Where | Duration / easing | `prefers-reduced-motion: reduce` |
|---|---|---|---|---|
| 1 | **ScannerRule sweep** — the signature | **§01 bottom edge only**, page-wide count: **1** | 2400ms, `--z-ease-scanner`, alternate, infinite, pip 18% of track | **Pip parks at track centre, full opacity, travel stops.** Never hidden — brand §7.3: the lit slot is identity, the travel is decoration. |
| 2 | Parked scanner tracks | TF-1/2/3 chrome bars + §15 top edge (4 instances) | **none — static by design** | unchanged |
| 3 | CTAButton primary hover | 3 buttons | 140ms `--z-ease-standard` (fill, glow) | 1ms |
| 4 | CTAButton secondary hover | 1 button | 140ms `--z-ease-standard` (edge, label, glow) | 1ms |
| 5 | Focus ring appearance | all 15 focusable elements | 90ms `--z-duration-instant` | 1ms — **the ring itself always appears; only its transition collapses** |
| 6 | NavLink / FooterLink underline | 4 links | 140ms `--z-ease-standard` | 1ms |
| 7 | FAQItem expand | 6 items | 220ms `--z-ease-out` | 1ms |
| 8 | FAQItem collapse | 6 items | 140ms `--z-ease-in` | 1ms |
| 9 | GridHorizon | §01 | **static** | unchanged |

**Not present, deliberately:** scroll-triggered reveals, parallax, marquees, typewriter effects,
counting numbers, hover lift/translate on anything, carousels, autoplaying anything, spring or
bounce easing of any kind. Brand §7.2: *mechanical, not playful. This is a machined object.*

The reduced-motion implementation is a global collapse plus one explicit re-park, because a blanket
`animation: none` would leave the pip wherever the cascade dropped it:

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 1ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 1ms !important;
    scroll-behavior: auto !important;
  }
  .z-scanner::after {                 /* re-park, do not hide */
    animation: none !important;
    left: calc(50% - (var(--z-scanner-pip-width) / 2));
    opacity: 1;
  }
}
```

**A scanner sweep that cannot be turned off is an accessibility defect** — and one that is
*hidden* under reduced motion is a brand defect. This does neither.

### 7.7 Iconography — there isn't any

**No icon set is introduced on this page.** Not Feather, not Lucide, not Heroicons, not an emoji.
Everything an icon would have done is done by:

- **the brand's four state glyphs** (`○ ◐ ◉ ⊘`) — already designed, already CVD-verified at
  ΔE ≥ 11.9, already the product's own vocabulary;
- **`▣`** for the `PHYSICAL` tag, drawn from the same geometric family;
- **mono readout labels** — a terminal labels things with words, which is exactly the register;
- **2px amber tick rules** in §08 where a list would otherwise reach for a checkmark.

Adding a general icon set would import a second visual language with a second set of proportions
and a second licence, to solve problems this page does not have. The glyphs come from the font
stack already loaded, so they cost nothing and can never fail to load independently of the text
they annotate.

### 7.8 SVG assets used

| Asset | Where | Size |
|---|---|---|
| `site/public/logo.svg` (viewBox `0 0 364 64`) | Nav | 140px ≥768 · 120px @375 |
| `site/public/favicon.svg` | `<link rel="icon">` | 16–32px |
| `site/public/flowforge-logo.svg` (viewBox `0 0 36 36`) | Footer `PoweredByFlowForge` | 28px |

`logo-mark.svg`, `logo-mono-white.svg` and `logo-mono-black.svg` are **not used on this page** —
the page is dark-native and full-colour on dark is the primary variant (brand §3.6). No OG image
file exists (`content/seo.md` records this); the build emits no `og:image` rather than pointing at
a file that 404s.

---

## 8 · Accessibility & motion notes

Target: **WCAG 2.2 Level AA**, which the brand sets as its standard. Most of this page lands AAA
on contrast because the palette was built that way.

### 8.1 Contrast pairings as token pairs

Cited values (**bold**) are published in `brand-manual.md` §4.2 and are carried, not re-estimated.
Unbolded values are pairs the manual does not publish because they involve a non-default surface;
they are computed with the same method, which was verified to reproduce all sixteen of the
manual's own figures exactly.

| Pairing | Tokens | Ratio | Req. | Verdict |
|---|---|---|---|---|
| Body prose | `--z-text` on `--z-surface` | **16.65** | 4.5 | AAA |
| Body prose on ground | `--z-text` on `--z-bg` | 17.37 | 4.5 | AAA |
| Subhead / FAQ answers | `--z-text-secondary` on `--z-surface` | **9.36** | 4.5 | AAA |
| Card body | `--z-text-secondary` on `--z-surface-raised` | 8.41 | 4.5 | AAA |
| Captions, microcopy | `--z-text-tertiary` on `--z-surface` | **6.15** | 4.5 | AA |
| Caption on raised | `--z-text-tertiary` on `--z-surface-raised` | 5.53 | 4.5 | AA |
| Worst text pair on the page | `--z-text-tertiary` on `--z-surface-overlay` | **4.89** | 4.5 | AA |
| H1 / H2 display | `--z-text` on `--z-bg` | 17.37 | 3.0 (large) | AAA |
| Emphasised heading clause | `--z-primary` on `--z-surface` | **10.59** | 3.0 | AAA |
| Primary CTA label | `--z-accent-contrast` on `--z-accent` | **12.68** | 4.5 | AAA |
| Primary CTA hover label | `--z-accent-contrast` on `#8CF0FF` | 15.45 | 4.5 | AAA |
| Secondary CTA label | `--z-text` on `--z-bg` | 17.37 | 4.5 | AAA |
| Secondary CTA hover label | `--z-primary` on `--z-bg` | **11.05** | 4.5 | AAA |
| Links | `--z-text-link` on `--z-surface` | **12.15** | 4.5 | AAA |
| Link hover | `--z-text-link-hover` on `--z-surface` | **14.81** | 4.5 | AAA |
| Terminal amber readout | `--z-primary` on `--z-bg` | **11.05** | 4.5 | AAA |
| State: NOT STARTED | `--z-steel` on `--z-bg` | 7.95 | 4.5 | AAA |
| State: IN PROGRESS | `--z-primary` on `--z-bg` | **11.05** | 4.5 | AAA |
| State: ZERADO | `--z-accent` on `--z-bg` | **12.68** | 4.5 | AAA |
| State: ABANDONED | `--z-orchid-500` on `--z-bg` | 7.53 | 4.5 | AAA |
| Roadmap `PLANNED` label | `--z-steel` on `--z-surface-raised` | 6.85 | 4.5 | AA |
| Price verdict | `--z-primary` on `--z-surface-raised` | 9.53 | 4.5 | AAA |
| Cover-tile platform tag | `--z-steel` on `--z-surface` plate | **7.62** | 4.5 | AAA |
| Cover-tile PHYSICAL tag | `--z-primary` on `--z-surface` plate | **10.59** | 4.5 | AAA |
| Hero text over amber bloom (worst case, 10%) | `--z-text` on `#1E1709` | 15.24 | 4.5 | AAA |
| **Focus ring** | `--z-focus-ring` on `--z-bg` | **12.68** | 3.0 | AAA |
| Focus ring on raised | `--z-focus-ring` on `--z-surface-raised` | 10.93 | 3.0 | AAA |
| Control edge | `--z-border-strong` on `--z-surface` | **4.08** | 3.0 (1.4.11) | pass |
| Control edge, worst surface | `--z-border-strong` on `--z-surface-overlay` | 3.24 | 3.0 | pass |
| Scanner pip on its track | `--z-scanner` on `--z-scanner-track` | 3.63 | n/a — decorative | pass |

### 8.2 CVD and colour-independence

The four state colours were verified by the brand (Viénot/Brettel/Mollon dichromat simulation,
CIEDE2000): worst pair `zerado × abandoned` under deuteranopia at **ΔE 11.9**, every pair over 10
in every model. This page adds no new meaning-bearing colour, so that verification carries.

**No information on this page is conveyed by colour alone.** Every state carries colour + glyph +
label; the `PHYSICAL` tag carries colour + `▣` + word; the roadmap status carries colour + `○` +
`PLANNED` or colour + `◐` + `IN PROGRESS` (§1.0, Amendment 1); the store rows carry colour +
`◉`/`○` + a per-row status word exposed to assistive technology — position under a labelled
heading is **not** sufficient, because the glyph is `aria-hidden` and a screen reader would
otherwise hear four identical rows. Under `forced-colors: active` the page surrenders colour to the OS and stays fully
readable: glyph and label carry everything, and terminal frames keep a 1px border so they do not
dissolve into the background.

### 8.3 Contrast pairings I rejected before they entered this blueprint

Five pairings were designed, measured, failed, and replaced. They are recorded so nobody
re-derives them:

| Rejected pairing | Ratio | Intended use | Replaced with |
|---|---|---|---|
| `--z-amber-900` `#8A5E00` on `--z-surface-raised` | **3.06** ✗ | The roadmap `PLANNED` label — "unlit amber" is the obvious semantic for *not yet*. | `--z-steel` `#A5A29B` at **6.85**. `--z-amber-900` survives as the status **track** fill, which holds no text. |
| `--z-cyan-900` `#0B6C7D` on `--z-surface-raised` | **2.87** ✗ | The §06 planned-store rows (PlayStation / GOG / EA) — "unlit cyan" for not-connected-yet. | `--z-text-tertiary` `#8492A8` at **5.53**, with the `○` glyph carrying *planned*. |
| `--z-scanner-500` `#FF2E2E` on `--z-surface-overlay` | **4.17** ✗ | The §07 price verdict — red for "don't buy". | `--z-primary` `#FFB000` at **9.53**. It is also semantically right: "maybe wait" is a readout, not an alarm, and the brand reserves scanner red for motion and alarm. |
| `--z-slate-500` `#64748B` on `--z-surface-overlay` | **3.24** ✗ | Caption text inside overlay-filled cards. | `--z-text-secondary` `#A9B5C7` at **7.43**. `--z-slate-500` stays what the brand made it: a UI boundary, where 3.24 satisfies 1.4.11. |
| `--z-border` `#2A3342` as a control edge | **1.53** ✗ | MoodCard / PriceCard / PhaseCard edges — it is the prettier hairline. | `--z-border-strong` `#64748B` (≥3.24 on every surface). `--z-border` is restricted to hairlines *between content blocks* (§02, §14, §16), which is the one use the brand permits. |

### 8.4 Structure, landmarks and focus order

```html
<a class="z-skip" href="#main">Skip to content</a>
<header role="banner">   <nav> lockup · GitHub · Join the waitlist </nav> </header>
<main id="main">         sections 01 … 15 </main>
<footer role="contentinfo"> section 16 </footer>
```

- **One `<h1>`**, in §01. Every other section heading is `<h2>`. No level is skipped; there are no
  `<h3>`s used structurally (MoodCard titles are `<p>` with H3 *styling*, because they are labels
  inside a card, not document structure).
- Sections **with** a visible heading get `<section aria-labelledby="…">` and become named
  regions. Sections **without** one (§02, §15) are plain `<section id="…">` with **no invented
  `aria-label`** — a landmark named with a string that appears nowhere on the page is a
  fabrication, and an unnamed non-region is the correct alternative.
- **Focus order is DOM order**, and DOM order is visual order at every breakpoint. Nothing on this
  page uses `order`, `row-reverse`, or a positive `tabindex`. The full tab sequence is:
  `skip → nav GitHub → nav CTA → hero CTA primary → hero CTA secondary → [any overflowing terminal
  screen] → FAQ ×6 → closing CTA → footer waitlist → footer GitHub → footer contact`.
- **Nothing on this page traps focus** — there is no modal, no dialog, no overlay, no
  auto-focusing element.
- `<html lang="en">`. The two Portuguese words (`zerado`, `sinopse`) are wrapped
  `<span lang="pt-BR">` on the occurrences copy marks as the explained first appearance, so a
  screen reader pronounces them in Portuguese rather than mangling them in English — which is the
  entire point of keeping them (decisions.md Q5-b).
- `<title>`, meta description, canonical, OG and Twitter tags come from `content/seo.md`
  verbatim. No `og:image` is emitted, because no such asset exists.
- Terminal-frame decorative characters (box-drawing runs, the parked scanner track, the sparkline)
  are `aria-hidden`; the rows' text is real text and is read normally.

### 8.5 Performance posture (design-side)

- **LCP candidate is the hero `<h1>` text.** No image can outrank it, because the page has none.
- **Preload three font files** — `Orbitron-latin.woff2` (the LCP element's face),
  `SpaceGrotesk-latin.woff2`, `JetBrainsMono-latin.woff2`. The two `latin-ext` subsets are left to
  lazy `unicode-range` resolution; the page's copy is Latin-1 apart from the box-drawing and state
  glyphs, which resolve from the mono `latin` subset and the system fallback.
- **CLS defence:** terminal frames set an explicit `min-block-size` computed from row count ×
  line-height, cover tiles set `aspect-ratio: 3/4`, and phone frames set fixed dimensions — so a
  font swap or a slow paint cannot shift the page. Nothing on the page is injected after paint,
  because there is no script to inject it.
- **INP** is structurally near-zero: the only interactions are native link activation and a native
  `<details>` toggle.

---

## 9 · What the build may not change without coming back here

1. The sixteen section ids, and their order.
2. Any string from `content/landing-copy.md` — wire by reference, never retype, never re-word.
3. The four breakpoints — 375 / 768 / 1280 / 1920.
4. The rule that every mockup carries a **visible** caption naming it a mockup.
5. The absence of real game cover art, of a donate/sponsor control, of the community-source name banned by ratified decision Q4 (that name is deliberately not written anywhere in this repository — see `docs/REDACTIONS.md`), and of
   any roadmap status other than `Planned` or `In progress` — and **never** a done-equivalent
   (`Done`, `Complete`, `Shipped`, `Released`, `Live`) on any phase, whatever its state.
   `In progress` was added by **Amendment 1** (§1.0); at rev A only `Planned` was permitted.
   Adding a further status means another amendment, not an edit to this line.
6. The focus ring, on any element, in any state.
7. The reduced-motion behaviour of the scanner: **parks, never hides.**
