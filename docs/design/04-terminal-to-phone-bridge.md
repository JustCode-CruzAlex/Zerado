---
title: Zerado — The Terminal-to-Phone Bridge
discipline: DESIGN SYSTEM
doc-no: ZRD-DESIGN-05
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: concept-explainer
---

# Zerado — The Terminal-to-Phone Bridge

How the identity holds when the same product renders as a Flutter app in **Phase 4**.

Grounded in brand manual **§9 — one identity, three render targets**. This document decides
**what translates**, **what cannot**, and **what the phone is allowed to do differently**, so
the phone apps are not designed from scratch two phases from now — and so nobody "ports the
terminal to a phone", which would produce a worse product and a worse advertisement for the
brand.

**Scope note:** this is the bridge, not the Phase 4 design. Per the ticket, Flutter app design
beyond this document is out of scope.

---

## 1 · The problem, stated

One brand, three surfaces with almost nothing technically in common:

| | **Web** | **Terminal** | **Phone (Phase 4)** |
|---|---|---|---|
| Colour depth | 24-bit | 256 → 16 → none | 24-bit |
| Layout unit | px, fluid | character cells | dp / pt |
| Type control | full | one monospaced grid | full, platform metrics |
| Motion | CSS, 60 fps | redraw budget | platform animation |
| Logo | `logo.svg` | `[0]` | `logo-mark.svg` |

**The terminal is the product's home.** When web and terminal conflict, the terminal wins
(brand §9). The phone is neither — it is a third expression that must be recognisably the same
product without pretending to be a terminal.

**The public promise this document has to keep**, from `content/landing-copy.md` §10:

> *"Native iOS and Android apps are planned for Phase 4 — the same library, the same states, the
> same mood sort, in your pocket instead of your terminal. The old form and the new form, one
> product."*

Read it precisely. The promise is **"the same library, the same states, the same mood sort"** —
it is *not* a promise of the same interface. That distinction is the whole design brief for
Phase 4.

---

## 2 · The invariants — identical on all three surfaces

Non-negotiable. If one of these differs on the phone, the phone is a different product.

### 2.1 · The four states, with colour and glyph and label

**Co-render is platform-independent.** A phone has richer affordances than a terminal, which
makes it *easier* to drop the label, not harder to keep it. The label stays.

| State | Colour (24-bit) | Mark | Label |
|---|---|---|---|
| Not started | `#A5A29B` — the warm grey | empty ring | `NOT STARTED` |
| In progress | `#FFB000` | ring half filled | `IN PROGRESS` |
| Zerado | `#19E0FF` | ring with a solid core | `ZERADO` |
| Abandoned | `#C77DFF` | ring struck through | `ABANDONED` |

**Use the screen hex values** — the phone is 24-bit. Not the ANSI-256 approximations, which are
a terminal artefact. Not the paper expression, which is for ink. This is easy to get wrong and
easy to catch.

### 2.2 · Amber common, cyan earned

Amber is the ambient voice. **Cyan is spent on completion and on the single most important
action on a screen.** The two-class model from `02-colour-budget.md` §2 carries over unchanged:
*state* cyan is data and unbounded; *chrome* cyan is budgeted at one per screen. A phone screen
with a cyan app bar, a cyan FAB and a cyan tab indicator has thrown the signal away exactly as a
terminal screen would.

### 2.3 · The scanner is the only signature motion

One motion belongs to this brand: a single light travelling across a dark track. **2400 ms full
cycle, `cubic-bezier(0.45, 0, 0.55, 1)`, alternate, infinite.** Reduced motion **parks the pip**
at the centre of its track at full opacity — it is deliberately not hidden, because the lit slot
is the identity and the travel is the decoration.

**No second signature motion may be introduced on the phone.** Platform transitions are the
platform's (§4); the scanner is Zerado's, and it stays the only one.

### 2.4 · The voice

Dry, confident, concrete. **No exclamation marks. No emoji. Never call the user a "gamer". Say
the number.** *"Zerado. 41 hours. Sixth this year."* reads identically on a phone, and it should
be the **same string**, not a re-written one (§7).

### 2.5 · Dark by default

Zerado is dark-native on every surface.

### 2.6 · The casing convention

`Zerado` the product · `zerado` the command · *zerado* the status in prose · **`ZERADO`** the
state chip in the interface. A phone has chips; they read `ZERADO`.

---

## 3 · What cannot translate — and why

Naming these prevents the two opposite failures: porting the terminal to a phone, and forgetting
the brand entirely.

| Terminal mechanism | Why it cannot cross | The phone's answer |
|---|---|---|
| **The character-cell grid** | There are no cells. Layout is constraint-based and continuous. | Platform layout. The *rhythm* the grid produced — consistent inset, one respiro gap between blocks — carries as a spacing scale; the grid itself does not. |
| **The fixed 14-column state chip** | Column arithmetic has no meaning in dp. | A chip whose **mark occupies a fixed square box** so a list reads as an aligned column (§5.2). |
| **Box drawing** (`┌ │ ─ ━`) | Not typography on a phone; drawing it would look like a novelty terminal theme. | Platform dividers, surfaces and cards. |
| **The 16-colour floor** | A phone has 24-bit colour. There is no floor. | No analogue is needed — but the *discipline* it forced is worth keeping (§4.3). |
| **`NO_COLOR`** | An environment variable with no phone equivalent. | The platform accessibility settings: high-contrast, grayscale, Smart Invert, `forced-colors`. The co-render rule is what makes all of them safe, exactly as it makes `NO_COLOR` safe. |
| **The `ZERADO_ASCII` fallback** | No font-coverage problem exists. | Not carried. The marks are drawn vectors and always render. |
| **The monospace readout grid as *primary* hierarchy** | On a phone, hierarchy comes from the type scale, weight and space. Setting a whole phone screen in monospace reads as a costume. | **JetBrains Mono survives as a *role*, not as the page** — numerals, hours, prices, commands, and the state labels. Everything else is Space Grotesk. |
| **Orbitron as terminal display type** | It never rendered in a terminal anyway. | Orbitron **returns** on the phone, under its hard rule: display only, never below 23 px equivalent, never prose, never more than about eight words, always uppercase, always tracked out. |
| **The terminal mark `[0]`** | Bracket-and-zero is a terminal artefact. | `logo-mark.svg`. Full colour requires a dark ground; on light use `logo-mono-black.svg`. |
| **The five column tiers (40/60/80/120)** | Columns are not a phone dimension. | M3 window size classes / platform size classes (§6). The *principle* survives: **design for the smallest real surface first, and let the largest add panes rather than rescue them.** |
| **The 24-column / 8-row refusal floor** | No phone is too narrow. | Not carried. |

---

## 4 · What the phone is allowed to do differently

> **Wear the platform; keep the identity.** A phone app that fights iOS or Android to look like
> a terminal is a worse product and a worse advertisement for the brand (brand §9).

### 4.1 · Explicitly permitted, and expected

| Area | The phone may |
|---|---|
| **Navigation** | Use the platform's model — tab bar, navigation bar, rail, drawer, list-detail. Zerado's terminal navigation model does **not** transfer. |
| **Gestures** | Swipe, long-press, pull-to-refresh, edge-swipe back. None exist in the terminal. |
| **Back behaviour** | Follow the platform. Android's system back — including predictive back — is a platform behaviour and must be honoured, not reimplemented. |
| **Type metrics** | Platform line heights, optical sizes and letterspacing. **Real tracking is available on a phone**, so the readout role recovers the 0.18em the terminal could not give it (`01-design-system.md` §1.5). |
| **Dynamic Type / font scale** | Fully supported. Layouts must reflow at the largest accessibility sizes — **the state label must never be the thing that truncates.** |
| **Touch targets** | Sized to the platform minimum, and never below **WCAG 2.5.8 Target Size (Minimum), 24 × 24 CSS px, Level AA** — a criterion that is inapplicable in the terminal and becomes live here. Confirm the current Apple and Material figures at source at Phase 4 kickoff. |
| **Elevation, material, shadow** | Permitted — **this is the clearest case of allowed divergence.** The terminal forbids fill-based elevation because of the 16-colour floor; the phone has no such floor and may use surfaces and shadow normally. |
| **Platform loading idioms** | Skeletons, refresh indicators, snackbars, toasts. |
| **Haptics** | Permitted, sparingly, and never as the sole channel for anything. |
| **Cover art** | Phase 2 metadata is a first-class visual on a phone in a way it never is in a terminal. |

### 4.2 · One codebase must not mean one look

Flutter compiles a single codebase to both platforms, and the standing temptation is a single
uniform interface. **The public copy promises "native iOS and Android apps."** Keeping that
promise honest means the app **adapts per platform** — Cupertino idioms and navigation on iOS,
Material 3 on Android — while the **identity layer** (§2) stays byte-identical across both.

**The split, stated so it can be reviewed:** *platform owns structure and interaction; Zerado
owns colour, state, mark, motion and voice.*

### 4.3 · The discipline worth keeping even though the constraint is gone

The 16-colour floor forced the terminal to separate regions by **border and spacing rather than
fill**, and the result reads calm. The phone is free to use fill — but a screen that separates
everything with tinted surfaces will read busier than the terminal it is a sibling to. **Reach
for spacing first, fill second.** This is guidance, not a rule.

---

## 5 · The glyph question — the substitution rule

`○ ◐ ◉ ⊘` are terminal glyphs — codepoints chosen partly for terminal font coverage. On a phone
they become **drawn shapes**. Left unruled, that is where the state system quietly stops being
one system.

**The rule, four clauses:**

### 5.1 · The progression is the invariant, not the codepoint

The four marks must read as **one sequence**: an empty ring → a ring half filled → a ring with a
solid core → a ring struck through. Brand §4.3: *"the sequence reads as a story even in one
ink."*

A phone may redraw them as vectors at any size. It may **not** reorder them, may **not**
resymbolise them, and may **not** break the ring family. **No checkmark for zerado. No trash
icon for abandoned. No progress arc for in-progress.** The moment one mark leaves the family,
the set stops being a progression and becomes four unrelated icons.

`◉` is the strongest mark in the set on purpose: *"the state the product is named after is the
one that looks like an achievement."* That relationship survives the redraw or the redraw is
wrong.

### 5.2 · One silhouette, one box, one stroke weight

All four marks occupy an **identical square bounding box** at any given size and share a single
stroke weight, so a list reads as a column of aligned marks. This is the phone's equivalent of
the terminal's fixed 2-column glyph field, and it exists for the same reason: **a state column
that does not align is a state column nobody scans.**

### 5.3 · The label never disappears

Co-render is platform-independent. A phone list row may present the label compactly, but may
**not** drop it and may **not** rely on an icon-only chip. Under high-contrast, grayscale or
Smart Invert, the **mark and the label** carry the state — exactly as they do under `NO_COLOR`.

The reason is the same one that made the label load-bearing in the terminal, restated for a
platform that *does* have an accessibility API: the label must be the accessible name of the
chip, not a decoration beside it.

### 5.4 · The CVD floor is re-verified, not assumed

If a redraw changes a mark's **fill ratio or stroke weight**, it changes how the four read under
dichromacy — the terminal's ΔE floor was measured on colour, but the glyph is what carries the
remainder at the tight pair. **Any redraw is reviewed against the same worst pair: zerado ×
abandoned under deuteranopia**, where the colours sit at **ΔE 11.9** and the mark is doing real
work.

---

## 6 · The governing corpora for Phase 4

Silent on the terminal; live on the phone. Chapter names read at source **2026-08-25** —
**re-verify at Phase 4 kickoff**, because both corpora move.

| Concern | Apple HIG | Material 3 |
|---|---|---|
| Navigation | Components → *Navigation and search*; Patterns → *Modality* | `components/navigation-bar` · `navigation-rail` · `navigation-drawer` · `app-bars`; `foundations/layout/canonical-examples/list-detail` |
| Type scale | Foundations → *Typography* | `styles/typography/type-scale-tokens` · `applying-type` · `fonts` |
| Dynamic type / font scaling | Foundations → *Typography* | `foundations/usability` |
| Touch targets | Foundations → *Layout* | `foundations/layout/grids-spacing/density` — plus **WCAG 2.5.8, 24 × 24 CSS px, AA** |
| Dark mode | Foundations → *Dark Mode*, *Color* | `styles/color/choosing-a-scheme` · `styles/color/roles` · `styles/color/system/how-the-system-works` |
| Brand colour into a platform scheme | Foundations → *Color*, *Branding* | `styles/color/static/custom-brand` · `styles/color/advanced/define-new-colors` |
| Adaptive layout | Foundations → *Layout* | `foundations/layout/breakpoints/{compact,medium,expanded,large-extra-large}` · `foundations/layout/scaffold` |
| Motion | Foundations → *Motion* | `styles/motion/easing-and-duration/tokens-specs` · `styles/motion/transitions` |
| Loading & feedback | Patterns → *Loading*, *Feedback* | `components/progress-indicators` · `components/snackbar` |
| Onboarding · settings · search | Patterns → *Onboarding*, *Settings*, *Searching* | `components/search` · `foundations/content-design` |
| Gestures & input | Inputs → *Gestures*, *Keyboards*, *Focus and selection* | `foundations/interaction/{gestures,inputs,selection,states}` |
| Accessibility | Foundations → *Accessibility*, *Inclusion* | `foundations/building-for-all/{user-needs,co-design}` |
| Back behaviour | Components → *Navigation and search*; Patterns → *Modality* | `foundations/layout/scaffold`; **Android system/predictive back is platform documentation, not M3** |

**WCAG 2.2 stays the floor**, and criteria that were inapplicable in the terminal become live:
**2.5.8** Target Size · **2.5.1** Pointer Gestures · **2.5.2** Pointer Cancellation ·
**2.5.4** Motion Actuation · **1.3.4** Orientation · **1.3.5** Identify Input Purpose ·
**4.1.2** Name, Role, Value (a real accessibility API exists — the state chip's accessible name
is the label from §5.3) · **3.1.1/3.1.2** Language (so *zerado* and *sinopse* can finally be
marked as Portuguese, which the terminal could not do).

---

## 7 · What Phase 1 must do now so Phase 4 is possible

Forward-compatibility debt is cheap to avoid now and expensive later. These are constraints on
the **Phase 1 spine**, and belong in `fft-tui-architect`'s deliverable.

1. **States serialize as stable identifiers, never as glyphs or labels.**
   `not_started` · `in_progress` · `zerado` · `abandoned` in the SQLite schema and on every
   seam. The glyph and the label are **presentation**, chosen by each surface. Storing `◉` or
   `ZERADO` would make the phone parse the terminal's rendering.
2. **Copy lives apart from render code.** The ratified voice is an asset both surfaces reuse
   verbatim — *"Zerado. 41 hours. Sixth this year."* should be one string, not two
   independently-drifting ones.
3. **Formatting rules are data, not baked strings.** Hours, prices, dates and counts need one
   definition of how they render, shared across surfaces.
4. **The SQLite file is the Phase 4 sync boundary.** Phase 4 must not require re-modelling
   states or re-deriving history — and the ratified promise stands until then: **no Zerado-run
   server before Phase 4**, and the library remains one file the player owns.
5. **The metadata seam stays swappable.** IGDB is free for **non-commercial use only** and
   Zerado's funding model is affiliate commission, which is commercial. Cover art and *sinopse*
   are the visual backbone of a phone app far more than of a terminal — so the phone is the
   surface most exposed if IGDB says no. Design the seam so the provider can change.
6. **No named community source.** The ratified decision stands on every surface, including a
   phone screen.
7. **Physical copies are first-class from day one.** A hand-added disc is not a second-class row
   in the terminal and must not become one on a phone.

---

## 8 · Two verified gaps Phase 4 will hit

Recorded now, while they are cheap.

### 8.1 · The light-mode state colours are **not** CVD-verified

Brand §4.4 ran the Viénot / Brettel / Mollon simulation on the **dark** state set and recorded
the **ΔE 11.9** floor. The light expression of brand §4.5 is implemented in `tokens.css` **§10**
(the `[data-z-surface="paper"]` block), and it defines a **different set of four state colours**
for light grounds:

| State | Light value | Ratio on `#FFFFFF` |
|---|---|---|
| Not started | `#5E6A7A` | 5.50 |
| In progress | `#8A4F00` | 6.56 |
| Zerado | `#0A6070` | 7.19 |
| Abandoned | `#6D3D93` | 7.30 |

Their **contrast** ratios are recorded. Their **CVD separation is not** — no simulation is
recorded for this set anywhere in the manual or the tokens, and brand §4.5 contains no CVD
paragraph at all. Contrast and colour-vision separation are different measurements; clearing one
says nothing about the other.

A phone will meet the system light mode. **Before any Zerado light theme ships on a phone, the
paper/light state set must be simulated for protanopia and deuteranopia across all six pairs, to
the same ΔE ≥ 10 bar.** Owner: `fft-brand-architect`. Do not assume it passes because the dark
set did — the dark set's own first draft failed at ΔE 8.8, which is exactly why this check
exists.

### 8.2 · The dark surface the ratios were measured against

Every ratio in brand §4.2 is measured against **`--z-surface` `#0B0D14`**. If the phone adopts a
platform-default dark surface instead, **the published ratios no longer describe what ships.**

**The rule:** the phone's dark scheme adopts `--z-surface` `#0B0D14` as its surface — inheriting
the measured table intact — **or** every pair is re-measured against whatever surface it does
adopt. There is no third option, and "it looks fine" is not one.

---

## 9 · Phase 4 design routing

| Step | Owner |
|---|---|
| The Phase 4 design brief — HIG/M3 chapters, the acceptance bar, this bridge as the identity contract | `fft-design-architect` |
| The phone design — screens, platform-idiom composition, adaptive iOS/Android structure | `fft-designer` (interim mobile/native design leaf) |
| Implementation | `fft-flutter` |
| Brand-side work — the light-mode CVD run (§8.1), the phone colour scheme | `fft-brand-architect` |
| GOLDEN verdict | `fft-code-reviewer`, against the Phase 4 brief |

**When Phase 4 design work actually starts, mint a dedicated `fft-flutter-designer` leaf** rather
than carrying the interim assignment indefinitely. The trigger is real work, not anticipation.

---

## 10 · Open for the founder

1. **The light-mode CVD gap (§8.1).** The paper/light state colours have never been simulated
   for colour-vision deficiency. This is a real gap, not a formality — the dark set's first draft
   failed. Confirm the simulation is scheduled before any phone light theme.
2. **The dark surface question (§8.2).** Confirm the phone adopts `--z-surface` `#0B0D14` so the
   measured contrast table carries over intact, rather than inheriting a platform default and
   invalidating every published ratio.
3. **"Native iOS and Android apps" (§4.2).** The public copy says native; Flutter is one
   codebase. Confirm the reading that the app **adapts per platform** — Cupertino on iOS,
   Material 3 on Android — with the identity layer identical across both. The alternative, one
   uniform look on both platforms, would be cheaper and would strain the published promise.
4. **The glyph substitution rule (§5).** It permits a redraw and forbids a resymbolisation.
   Confirm that the four marks must stay one ring family on the phone — no checkmark for
   *zerado*, however conventional that would be on a mobile list.
