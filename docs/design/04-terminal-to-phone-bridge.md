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
| Audio | none | opt-in; local FX + streamed stations; `ZERADO_NO_AUDIO` | opt-in, platform audio session |

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

**Dark-native is not dark-only.** Light is a **first-class second expression**, designed in brand
§4.5 and selected by the order in `05-theme-system.md` §3.5 — `NO_COLOR` first, then an explicit
setting, then background detection, then dark. Dark is the *fallback when nothing answers*, and
that is a measured position rather than a preference: brand §5.3 measured the palette against five
popular **dark** grounds and every state cleared 4.5:1 on all five, so guessing light on a failed
probe would guess against the only evidence there is. On a phone, step 3 is the platform's own
light/dark setting — but see §8.1: the light state set does not yet pass its own gate.

### 2.6 · The casing convention

`Zerado` the product · `zerado` the command · *zerado* the status in prose · **`ZERADO`** the
state chip in the interface. A phone has chips; they read `ZERADO`.

### 2.7 · Audio — the register, the posture, and the co-render extension

Audio became part of the product on **2026-08-25**, and the phone is the surface where a *second*
sound language is easiest to acquire by accident — because the platform hands you a whole sound
design for nothing. Three things are invariant.

**2.7.1 · The sound register translates unchanged.** The reference is **a machine acknowledging
an instruction** — a switch, a relay, a confirmation tone from an instrument panel. Short, dry,
mechanical: the audible sibling of *"motion here is mechanical, not playful: no bounce, no
elastic, no spring overshoot"* (brand §7.2), and the DeLorean-and-KITT bar applied to sound
(`01-design-system.md` §15.4).

**A MISS, not a variation, on any surface:** retro-arcade pastiche, chiptune jingles, coin-drop
or level-up stings, and — specifically on a phone — **the platform's own stock notification and
success sounds.** Those are the platform's identity, not Zerado's. Wearing the platform means
adopting its *conventions*, never its *voice*.

**2.7.2 · Opt-in and off by default is a posture, not a terminal implementation detail.** It
carries to the phone whole: nothing plays until the player turns it on, and the first run is
silent. This is the mobile temptation worth naming — a phone app that greets you with a sound on
first launch has broken the posture even though it kept the register.

**2.7.3 · Audio is never the sole carrier of information — non-negotiable everywhere.** The
co-render rule extended to a fourth channel: colour **and** glyph **and** label **and** —
optionally — sound. A cue is **always the second signal**, confirming what the screen already
said. **Removing all audio must remove no information**, exactly as `NO_COLOR` removes none.
On a phone this is load-bearing in a way it is not in a terminal: the device may be muted by a
hardware switch, in a pocket, or driving a screen reader — so the silent path is the *common*
path, not the exceptional one.

> **Governance flag.** Brand manual §9 lists the cross-surface invariants as *the four states,
> amber-common/cyan-earned, the scanner as the only signature motion, the voice, and dark by
> default.* **Audio is not in that list** — brand rev A is dated 2026-08-24 and the audio
> decision is 2026-08-25, so §9 predates it. **The caption rule below (§2.8) is a second such
> gap**, for the same reason: cover art became Phase 1 after rev A was written, so §9 could not
> have listed it either. These two sections are the working invariants until the manual catches
> up; **§9 should gain audio as a sixth and the caption rule as a seventh** via the brand
> governance procedure (owner: `fft-brand-architect`). Recorded rather than assumed, because a
> bridge document quietly extending the brand's own invariant list is exactly the drift §9 exists
> to prevent.

### 2.8 · A caption is always a real title, never decorative

Non-text content carries a text alternative on every surface — WCAG **1.1.1**, live since cover
art became Phase 1. The caption beside a cover tile is **the game's real title**, never a
decorative flourish, never a truncated slug, never absent.

**Why this is an invariant and not a terminal implementation detail.** In the terminal the
caption does two jobs at once: it is the alternative for a player using a screen reader, **and**
on a terminal with no image protocol it is *the entire content of the tile*. That second job is
what makes it obviously load-bearing — and **that job does not exist on a phone**, because there
is no phone that cannot draw an image (§3).

So the phone inherits **the rule without the reason.** A designer who only ever sees the phone
will find a caption sitting under an image that always renders, conclude it is redundant, and
optimise it away — at which point screen-reader users lose the only text there was. **A rule
whose justification does not travel is exactly the kind that gets deleted by someone acting in
good faith.** It is recorded here, at invariant level, so it survives that reasoning.

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
| **`ZERADO_NO_AUDIO`** | An environment variable with no phone equivalent. | The platform's own controls: the iOS silent switch, system volume, per-app audio settings — plus Zerado's own in-app opt-in and per-channel mutes, which carry unchanged. |
| **"No device, so silence"** | A phone always has an audio device. | The equivalent condition is *muted, in a pocket, or on a call.* The response is identical and is the point of §2.7.3: **silence loses no information.** |
| **The image-capability degrade** — a terminal that cannot draw, where the **caption is the entire content** | **There is no phone that cannot draw an image.** The condition simply does not exist on the platform. | No analogue, and none is needed. **But the rule the condition justifies still applies — see §2.8.** This is the one place in this document where a rule travels and its reason does not, which is precisely why the rule is written as an invariant rather than left to be re-derived from a constraint the phone never meets. |

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
| **Haptics** | Permitted, sparingly, and never as the sole channel for anything — the §2.7.3 rule covers haptics exactly as it covers sound. A phone-only channel with no terminal analogue. |
| **System volume model** | Use it. The phone has an OS volume model and per-app audio settings the terminal does not; Zerado's own per-channel volumes sit *underneath* it, never override it. |
| **Silent switch / Do Not Disturb** | **Must be honoured.** Zerado's audio is never essential, so it uses an ambient-class audio session that the silent switch silences. Never a playback/priority class that overrides it. |
| **Background audio conventions** | Zerado's audio is **foreground-only.** No background audio entitlement, no lock-screen transport controls, no now-playing metadata — Zerado is not a music player and must not present as one. |
| **Cover art** | Cover art is **first-class on both surfaces** from Phase 1 — what differs is **scale and the fallback path**, a difference of degree rather than of kind. A phone has pixels and no protocol question: tiles are drawn, sized to the layout, and always render. A terminal has a **capability axis** — Kitty and iTerm2 draw the image; everything else takes the designed text degrade. So the phone may show larger art, more of it, and at any aspect it likes; what it may **not** do is assume the caption is optional because its own art always renders (§3). |

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

### 4.4 · Audio on a phone mixes — it never interrupts

The rule with the most practical bite, and it has no terminal analogue at all: **a phone user is
very likely already playing something.** Music, a podcast, a call, a screen reader.

**Zerado never takes exclusive audio focus.** It mixes with whatever is playing, or it stays
quiet. Concretely:

| Situation | Behaviour |
|---|---|
| Another app is playing audio | Zerado's **music bed does not start.** FX may mix at low level, or stay silent — never duck or pause the other app |
| A call is active | Silent |
| A screen reader is speaking | Silent, or mixed below speech — **never over it.** §2.7.3 guarantees nothing is lost by staying quiet |
| Silent switch on / DND | Silent |
| Zerado goes to background | Audio stops |

**Why this is identity and not merely politeness.** Interrupting a player's own music to play
Zerado's would make the product the loudest thing in the room — the precise opposite of *an
expensive object that has just been unboxed.* A machine that talks over you is not KITT; it is a
kiosk.

**The one carried-over engineering constraint:** muting the **music** channel **halts decode and
releases the audio session**, rather than holding a gain of zero (`01-design-system.md` §15.5).
That rule matters more on a phone than in a terminal, because a held audio session there costs
battery and can block another app from playing at all.

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
abandoned under deuteranopia**, where the colours sit at **ΔE 11.81** (the pinned method of
`05-theme-system.md` §2.1; the manual records 11.9) and the mark is doing real work.

**On a phone this review widens, because a phone can theme.** The pair to protect is the tightest
pair **of the active theme**, not of the default — and in the DISTINCT band (10 ≤ ΔE < 15) the
glyph is carrying more load than it does in the default. `05-theme-system.md` §2.6 states the
consequence plainly: co-render is what makes a marginal theme *degraded* rather than *dangerous*,
and a redraw that weakens the glyph spends exactly that margin.

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
5. **The metadata seam stays swappable.** IGDB is free for **non-commercial use only**, and since
   founder direction on 2026-08-25 **dropped affiliate links entirely** — making Zerado cleanly
   non-commercial: free software, donation-supported, **zero revenue** — that specific blocker is
   closed (`ADR-0001` D1). **Do not restate the old premise in the present tense**: a document
   that says Zerado's funding is commercial re-opens a question the bundle reports as answered.
   The seam stays swappable anyway, on engineering grounds — cover art and *sinopse* are the
   visual backbone of a phone app far more than of a terminal, so the phone is the surface most
   exposed if IGDB says no for any other reason. Designing as if IGDB is guaranteed would be
   designing on a maybe.
6. **No named community source.** The ratified decision stands on every surface, including a
   phone screen.
7. **Audio settings and cue events are data, not code.** The opt-in flag, the two per-channel
   mutes and the two volumes must serialize as **preferences the phone can read**, and each cue
   must be a **semantic event identifier** (`sync_complete`, `error`, `state_zerado`) rather than
   a sound-file name. The same reasoning as item 1: the *event* is shared, the *asset* is each
   surface's own. A phone that had to parse the terminal's audio filenames would be the same
   defect as a phone parsing its glyphs.
8. **Physical copies are first-class from day one.** A hand-added disc is not a second-class row
   in the terminal and must not become one on a phone.

---

## 8 · Two verified gaps Phase 4 will hit

Recorded now, while they are cheap.

### 8.1 · The light-mode state colours **fail** CVD separation — measured

> **Status changed rev A.** This section previously recorded the light state set as
> *CVD-unverified* and flagged it for measurement. **It has now been measured, and it fails.**
> The measurement, the cause and a minimum-motion repair are in `05-theme-system.md` §3.2–§3.4;
> what follows is the summary and what it means for Phase 4.

Brand §4.4 ran the Viénot / Brettel / Mollon simulation on the **dark** state set. The light
expression of brand §4.5 is implemented in `tokens.css` **§10** (the `[data-z-surface="paper"]`
block), and it defines a **different set of four state colours** for light grounds:

| State | Light value | Ratio on `#FFFFFF` |
|---|---|---|
| Not started | `#5E6A7A` | 5.50 *(the manual's)* |
| In progress | `#8A4F00` | 6.56 *(the manual's)* |
| Zerado | `#0A6070` | 7.19 *(the manual's)* |
| Abandoned | `#6D3D93` | **7.67** *(computed — `tokens.css` §10 records 7.30; see below)* |

Contrast is fine. **Separation is not.** Under the pinned method (`05-theme-system.md` §2.1),
**two of the six pairs fail the ΔE ≥ 10 bar**, both computed:

| Pair | Worst model | ΔE | |
|---|---|---|---|
| not-started × zerado | protanopia | **5.41** | ✗ |
| zerado × abandoned | deuteranopia | **8.91** | ✗ |

**The cause is a single un-carried hue, and it is the one that mattered.** Brand §4.5 claims the
paper colours are *"the same hues carried to ink weight"*, and four of five are — abandoned within
0.12°, scanner within 0.41°, cyan within 3.11°, amber within 9.50°. **Not-started is rotated
173.02°**, from the dark set's warm `h 92.97°` to `h 265.99°` — landing **1.1°** from `#9FB0C6`,
the blue-cast steel that brand §4.4 *explicitly rejected* on the dark side for collapsing against
the cyan. The dark set's hard-won correction was not carried to paper, and the 5.41 above is that
omission measured.

**Two consequences, both binding.**

1. **No Zerado light theme ships — on a phone or in the terminal — until this is repaired.**
   `05-theme-system.md` §3.4 demonstrates a minimum-motion repair along axes canon already
   ratified (carry the hue; restore the lightness spread), reaching a floor of **10.83**. That is
   a feasibility proof, **not** a token change. Owner: `fft-brand-architect`, through brand
   governance §10.
2. **`tokens.css` §10's ratio for `--z-state-abandoned` is wrong.** It records `7.30:1`;
   `#6D3D93` on `#FFFFFF` measures **7.67:1** — recomputed with the same WCAG formula that
   reproduces all thirty-four other brand figures exactly. Nothing shipped broken (the true value
   is higher), but the correction belongs in `tokens.css`, `tokens.json` and the manual in one
   commit.

**And the general point that outlives this particular defect:** *"someone should check the light
palette"* was an open note for as long as this document existed. It is now a **gate** —
`05-theme-system.md` §2.2 — which means no theme, light or dark, default or harvested, can reopen
this class of gap silently again. That is the difference between a flag and a mechanism.

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
4. **Brand §9's invariant list is now incomplete by two (§2.7, §2.8).** Rev A lists five
   cross-surface invariants and predates both later decisions: **audio is a sixth** (2026-08-25)
   and **the caption rule a seventh** (cover art moving to Phase 1). Both are recorded here as
   working invariants. Confirm `fft-brand-architect` folds them into §9 through the governance
   procedure, so the manual and this bridge stop diverging — two gaps in one week is a rate worth
   noticing, not just a pair of items to close.
5. **The glyph substitution rule (§5).** It permits a redraw and forbids a resymbolisation.
   Confirm that the four marks must stay one ring family on the phone — no checkmark for
   *zerado*, however conventional that would be on a mobile list.
