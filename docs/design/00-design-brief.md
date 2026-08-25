---
title: Zerado — Design Brief
discipline: DESIGN SYSTEM
doc-no: ZRD-DESIGN-01
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: project-brief
---

# Zerado — Design Brief

**The governing design document.** It names which guideline corpora rule this product, exactly
where each one applies, which one wins when two disagree, and what "done" means for a Zerado
screen. It is consumed twice: by the designer who builds against it, and by the reviewer who
judges the result. One artifact, two readers — so nobody invents a private bar.

This brief holds guidelines as **pointers**. It does not copy a corpus inline, because a copy
is a stale copy. Where a number appears here it was read at its source and the source is named.

**Companion documents:** `01-design-system.md` (the component vocabulary) ·
`02-colour-budget.md` (the colour rule) · `03-designer-manual.md` (how we work) ·
`04-terminal-to-phone-bridge.md` (Phase 4) · `05-theme-system.md` (what a theme must supply, and
what makes one valid).

---

## 1 · The precedence rule — read this before anything else

Corpora disagree. When they do, this order decides, top wins, no discussion at the keyboard:

| # | Authority | Why it sits here |
|---|---|---|
| **1** | **The ratified public promises** — `ratification/decisions.md` + `content/landing-copy.md` | These are *published*. A design that contradicts them does not lose an argument, it makes the product a liar. Local-first, one SQLite file, no Zerado server before Phase 4, the user's own keys, no named community source, the affiliate disclosure. |
| **2** | **WCAG 2.2 Level AA** | The floor, not the ceiling. Non-negotiable on every surface. The brand already subordinates itself here by construction (brand manual §4.2: where a pairing failed, *the colour changed*), so 1 and 2 collide rarely — but if they ever do, AA wins and the fix goes through brand governance (§10), never through an exception. |
| **3** | **`brand/tokens.css` + `brand/tokens.json`** | The brand manual's own preamble: *"Those files are the implementation; this document is the reasoning. Where the two disagree, the tokens are what ships."* |
| **4** | **`brand/brand-manual.md` + `brand/naming.md`** | The identity authority. Wins on colour, glyph, type, voice, motion, casing. Carries §9's own tiebreak: **the terminal is the product's home — when web and terminal conflict, the terminal wins.** |
| **5** | **FlowForge TUI Spacing Canon (#2435)** | Adopted verbatim (§5 below). Spacing is a solved problem; Zerado does not re-solve it. |
| **6** | **FlowForge TUI Design Manual (#2371)** | The terminal craft bar, minus its FlowForge-internal rules (§4 below). |
| **7** | **Apple HIG · Material 3** | Phase 4 platform surfaces only. **Silent on the terminal** — they have no jurisdiction over a character grid, and citing them there would be cargo cult. |

**The tiebreak when the order does not resolve it** is the brand manual's own (§1):
*which option makes the product's own argument more clearly?* The argument is an expensive
object from 1984 designed for 2030 — the DeLorean and KITT. Retro-**future**, never
retro-nostalgia.

---

## 2 · The corpora, and exactly where each applies

| Corpus | Pointer (verify at source) | Governs |
|---|---|---|
| **WCAG 2.2** | `https://www.w3.org/TR/WCAG22/` — W3C Recommendation, **12 December 2024** (read 2026-08-25) | Every surface, always. AA is the floor. §3 says which criteria actually bite in a terminal. |
| **FlowForge TUI Design Manual #2371** | `documentation/2.0/design/TUI-DESIGN-MANUAL.md` | Terminal craft: library roles, R-1…R-11, the Screen-Quality Gate. §4 records what Zerado inherits. |
| **FlowForge TUI Spacing Canon #2435** | `documentation/2.0/design/TUI-SPACING-CANON.md` | Every spacing value on every terminal screen. Adopted **verbatim** — §5. |
| **Zerado brand manual (rev A)** | `forgeplay-output/landing-page/brand/brand-manual.md` | Identity: palette, the four states, ANSI-256, the 16-colour floor, `NO_COLOR`, type, the scanner, voice. |
| **Zerado tokens (200)** | `brand/tokens.css` · `brand/tokens.json` | The implementation of the above. **No surface may hard-code a hex.** |
| **Zerado naming** | `brand/naming.md` | Casing, everywhere, forever. |
| **Apple HIG** | `https://developer.apple.com/design/human-interface-guidelines` | **Phase 4 only** — §6. |
| **Material 3** | `https://m3.material.io` | **Phase 4 only** — §6. |

> **Anti-fabrication standing rule.** Never assert a clause number, a contrast ratio, a **ΔE
> separation**, or a platform minimum from memory. Every such figure printed in Zerado design
> documents must be one read from the brand manual's measured table, and cited as the manual's
> figure rather than as a measurement of our own. If a pair is not in that table, it is
> **unmeasured** — say so and flag it for measurement. Guessing is a defect, not a shortcut.
> §9 lists the values currently unmeasured.
>
> **The ΔE figures carry one caveat worth knowing:** the manual does not pin the CVD model variant,
> the white point, or the gamut-clamping rule, so its digits are **not independently reproducible**
> — two independent implementations of its own named method disagree with it, and with each other,
> in the third significant figure. The ordering claims those numbers support all hold; only the
> digits move. See `01-design-system.md` §3.2.
>
> **The model is now pinned.** `05-theme-system.md` §2.1 fixes the matrices, the D65 white point,
> the linear-RGB clamp and the CIEDE2000 parameters, because a validation gate cannot be built on
> an unpinned method. Under the pinned method the dark default's tightest pair measures **11.81**
> where the manual records 11.9. **Every ΔE printed in Zerado design documents from rev A onward
> is either the manual's, cited as the manual's, or computed under §2.1 and marked as computed.**
> Reconciling the manual's own digits to the pinned method is `fft-brand-architect`'s, upstream.

---

## 3 · WCAG 2.2 in a terminal — what actually bites

A character grid is not a DOM. Pretending otherwise produces a checklist nobody can fail and
nobody learns from. Below, every criterion is placed in one of three honest buckets.
Titles and levels were read from the W3C Recommendation on 2026-08-25.

### 3.1 · Bites directly — these are real obligations

| SC | Title | Level | What it means on a Zerado screen |
|---|---|---|---|
| **1.4.1** | Use of Color | A | **This SC *is* the co-render rule.** Colour is never the only channel. Every state carries colour **and** glyph **and** label. Non-negotiable, already ratified (brand §4.3). |
| **1.4.2** | Audio Control | A | **Live since audio was added to Phase 1.** Verified at source: *if any audio plays automatically for more than 3 seconds, either a mechanism is available to pause or stop it, or a mechanism is available to control its volume independently from system volume.* Satisfied structurally — see `01-design-system.md` §15.6. |
| **1.4.3** | Contrast (Minimum) | AA | 4.5:1 for text. Every state colour clears it on `--z-surface` (7.21 worst — abandoned) and on the five popular themes brand §5.3 measured. **Honest caveat in §3.4.** |
| **1.4.11** | Non-text Contrast | AA | 3:1 for glyphs, focus indicators and control boundaries. `--z-border` (**1.53**) fails and is therefore **decorative only**; controls use `--z-border-strong` (**4.08**). |
| **2.1.1** | Keyboard | A | Trivially met — the keyboard is the only modality. The inverse risk is what matters: no action may require a pointer. |
| **2.1.2** | No Keyboard Trap | A | **The most under-respected SC in TUI work.** A filter mode, a confirmation, a detail pane or a form that has no way out is a literal keyboard trap. **Escape always leaves; there is no screen without an exit.** |
| **2.1.4** | Character Key Shortcuts | A | Single-key shortcuts are Zerado's whole interaction model. The operative rule: **when a text input holds focus, single-key shortcuts do not fire.** Typing `d` in a search box searches for "d"; it does not delete anything. |
| **2.4.3** | Focus Order | A | Pane traversal order is designed, not incidental, and matches visual order. |
| **2.4.6** | Headings and Labels | AA | The header band's title describes the screen; section readouts describe their block. |
| **2.4.7** | Focus Visible | AA | **The focus ring is never removed, in any state, on any element** (brand §4.2). |
| **2.4.11** | Focus Not Obscured (Minimum) | AA | An overlay — confirmation, command palette, filter mode — may never *entirely* hide the focused row. This is what forbids a full-bleed modal over a ledger. |
| **2.2.2** | Pause, Stop, Hide | A | The scanner is auto-starting motion that can run past five seconds beside other content. Answered by design: it appears **only** for genuinely indeterminate waits (never ambient), and reduced-motion **parks the pip** rather than hiding it (brand §7.3). |
| **2.3.1** | Three Flashes or Below Threshold | A | Satisfied by construction: the scanner's full cycle is **2400 ms** — about 0.42 Hz, far below the 3 Hz threshold. Nothing else in Zerado blinks. |
| **3.2.3** | Consistent Navigation | AA | Global keys mean the same thing on every screen. |
| **3.2.4** | Consistent Identification | AA | A `ZERADO` chip looks and reads identically everywhere it appears. |
| **3.2.6** | Consistent Help | A | `?` opens the key map from every screen, always in the same place. |
| **3.3.1** | Error Identification | A | Errors are named in text, not signalled by a colour change. |
| **3.3.2** | Labels or Instructions | A | Every input says what it wants — e.g. where a Steam API key comes from, at the moment it is asked for. |
| **3.3.3** | Error Suggestion | AA | Errors say what to do next. The ratified voice already demands this: *"Steam returned an empty library. Game details are private on your profile… Settings → Privacy."* |
| **3.3.7** | Redundant Entry | A | Never ask for the Steam key or Steam ID twice in one setup flow. |
| **1.3.2** | Meaningful Sequence | A | **Terminal-specific and easy to get wrong.** A screen's reading order for anything consuming the output stream is *byte order*. Two panes side by side interleave. Any content whose meaning depends on being read in order must be in one column, in that order. |
| **1.3.3** | Sensory Characteristics | A | No copy may say "the cyan one", "the row on the right", or "the round icon". Instructions name the label or the key. |

### 3.2 · Applies only analogically — the intent transfers, the mechanism does not

| SC | Title | Level | The honest terminal analogue |
|---|---|---|---|
| **1.1.1** | Non-text Content | A | Glyphs are non-text content, but a TTY has no `alt`. The text alternative is the **label sitting beside the glyph** — the third leg of co-render. This is the mechanism, and it is why the label may not be dropped for density. |
| **1.3.1** | Info and Relationships | A | There is no programmatic structure to expose. Structure must be carried by *visible* means — position, box drawing, spacing, readout labels — because that is all there is. |
| **1.4.4** | Resize Text | AA | The user resizes by changing their terminal font. Zerado must never assume a cell size and must re-render correctly on `SIGWINCH`. |
| **1.4.10** | Reflow | AA | The 320 CSS px rule does not map. The intent — no loss of content or function on reflow — maps exactly onto the five-tier system and the documented refusal floor. |
| **1.4.13** | Content on Hover or Focus | AA | No hover exists. The focus-triggered analogue is the detail pane: it must be dismissible, must not obscure the focused row (2.4.11), and must persist while focus remains. |
| **2.4.2** | Page Titled | A | Two real mechanisms: the header band's title row, and setting the terminal window title (`OSC 2`) to the current screen. |
| **4.1.2** | Name, Role, Value | A | **There is no accessibility API in a TTY.** A screen reader sees the raw character stream — nothing else. The only way state survives is if it is *written as words*. This is the single strongest argument for the co-render rule, and the reason the label is load-bearing rather than decorative. |
| **4.1.3** | Status Messages | AA | No `aria-live`. The analogue: a status change is **printed as text, in a stable and predictable place** — the pinned status row — never signalled only by a colour changing somewhere on screen. |
| **3.1.2** | Language of Parts | AA | *zerado* and *sinopse* are Portuguese inside English copy. A TTY cannot mark language, so assistive tech will mispronounce them. **Unfixable in the terminal; recorded, not hidden.** It is fixable on web and on Phase 4 phones, and must be done there. |
| **3.3.8** | Accessible Authentication (Minimum) | AA | No cognitive test is imposed — but the concrete duty is real: **the Steam API key field must accept a paste.** Blocking paste would manufacture the exact barrier this SC exists to prevent. |

### 3.3 · Genuinely inapplicable — and why

| SC / group | Why it does not apply |
|---|---|
| **1.2.1 – 1.2.9** (time-based media) | No video, and no speech or informational audio. Phase 1 audio is **non-informational by design** — §15.1 of the design system forbids sound from ever being the only carrier of information — so there is no time-based media whose information needs an alternative. **1.4.2 Audio Control is a different matter and does apply: it has moved to §3.1.** |
| **1.3.4** Orientation | A terminal has no orientation lock. |
| **1.3.5** Identify Input Purpose | No autocomplete metadata exists; the fields are a Steam key and a Steam ID, neither of which is a listed input purpose. |
| **1.4.5 / 1.4.9** Images of Text | No images. Cover art arrives in Phase 2 and is art, not text. |
| **1.4.12** Text Spacing | The SC is **scoped in its own text** to *"content implemented using markup languages that support the following text style properties."* A character grid supports none of them. **Inapplicable as written** — not "waived". The residual duty is 1.4.4's: honour whatever font and line height the user has set. |
| **2.5.1 – 2.5.8** (pointer, gestures, dragging, target size) | No pointer. **2.5.8 Target Size (Minimum) — 24 × 24 CSS px, AA — becomes live in Phase 4** and is carried in `04-terminal-to-phone-bridge.md`. |
| **3.1.1** Language of Page | No markup in which to declare it. |
| **2.4.1** Bypass Blocks | No repeated navigation blocks to bypass; the header band is three rows. |
| **4.1.1** Parsing | **Removed from WCAG 2.2** (obsolete). Verified at source. Do not cite it. |

### 3.4 · The one honest limit on our contrast claim

Zerado's measured ratios are computed against **`--z-surface` `#0B0D14`**. In a terminal the
background belongs to the **user**, not to us. Brand §5.3 measured the four state colours
against five popular grounds — pure black, xterm dark grey, Solarized Dark, Gruvbox Dark, One
Dark — and every state clears **4.5:1** on all five, worst case **5.20** on One Dark.

**What we may claim:** AA on Zerado's own ground and on those five measured grounds.
**What we may not claim:** AA on an arbitrary user theme, because a user may set a light or
mid-tone background we never measured. The product's answer is structural, not a promise:
`NO_COLOR` surrenders colour entirely and the glyph and label carry the state unaided. Say it
this way in any conformance statement. Do not round it up.

---

## 4 · FlowForge TUI Design Manual (#2371) — what Zerado inherits

| Rule | Verdict | Note |
|---|---|---|
| **R-1** Cockpit is the GOLD STANDARD | **No Zerado analogue** | The cockpit is a FlowForge screen in a FlowForge product. Zerado has no cockpit and must not invent a reference to a screen it does not own. **Adapted:** the *principle* — an in-house screen that proves the bar is reachable — is kept, and **the library deck is nominated as Zerado's reference screen**, becoming the internal standard the moment it is ruled GOLDEN. Until then Zerado has no in-house gold standard and says so. |
| **R-2** Shared header template + spacing frame | **Adopted, adapted** | The band, the frame and router-level enforcement are inherited whole. **Adaptation:** Zerado's breadcrumb is **two segments** (`Zerado ✦ Screen`) because Zerado has no project tier — and R-2 explicitly forbids faking an empty segment. Separator decided in `01-design-system.md` §2. |
| **R-3** COLOR by default, theme-composed | **Adopted, adapted** | Zerado has exactly **one** palette, not a user-selectable theme set, so "theme-composed" becomes **token-composed**: zero hex in code, every colour resolved from the Zerado token table. The colour on/off path is not a new toggle — it is **`NO_COLOR`**, already ratified (brand §5.4). |
| **R-4** COMMIT to the design system | **Adopted verbatim in force** | No hand-rolled one-off styles, no magic-number spacing. The *artifact* is Zerado's own token table, because FlowForge's is not importable — see §7. |
| **R-5** Every screen must be a hit | **Adopted verbatim** | Including — especially — first run, the failed sync, and the private Steam profile. |
| **R-6** Co-render law | **Adopted verbatim** | Already the brand's own rule (§4.3). The palettes agree: Zerado's states were CVD-verified with Bang Wong and Paul Tol as reference targets (brand §4.4). |
| **R-7** Responsive breakpoints | **Adopted verbatim** | Tiny `<40` · Narrow `40–59` · Standard `60–79` · Wide `80–119` · ExtraWide `120+`. This is the replacement for the rejected over-80-column frames: **80 is the design width, 120 is progressive enhancement, and no screen may *require* more than 80 columns to be correct.** |
| **R-8** charm.land = the quality bar | **Adopted, adapted** | "Calm, composed, breathing" is inherited. The *visual study* is a FlowForge-internal artifact Zerado cannot cite as its own bar; **Zerado's visual bar is the brand manual's DeLorean-and-KITT register** plus R-8's calmness test. |
| **R-9** Cockpit footer shortcuts (`f` logs, `d` docs) | **No Zerado analogue** | These are affordances of a FlowForge screen. Zerado has no log drawer and no docs tree. **Inheriting `f` and `d` would be cargo cult** — Zerado's key map is its own and is designed in the spine. |
| **R-10** Ledger triad (a) identity column (b) cursor-following scroll (c) pinned totals | **Adopted verbatim as a RULE** | The library deck is a ledger and is not GOLDEN until all three hold **at an overflowing row count** — a 400-game library, not a 12-row golden. The `LedgerTable` package itself is **not importable** (§7), so Zerado builds its own primitive that is correct-by-construction on all three. |
| **R-11** Co-render and contrast survive the screenshot | **Adopted verbatim** | The rendered screenshot is the founder-facing artifact. Colour-alone state fails the screenshot exactly as it fails `NO_COLOR`. |

**The Screen-Quality Gate (#2371 §4b) is adopted whole** as Zerado's review process:
headless render (`freeze` / `vhs`) at the target viewports **and at an overflowing row count** →
specialist review of the shot → founder validation **before** merge → `fft-code-reviewer` marks
GOLDEN only *with* the founder-validated screenshot. No screenshot, no GOLDEN, no merge.

**Charm library roles** (#2371 §2) are adopted as the reuse map: `lipgloss` layout and style ·
`bubbles` components · `glamour` markdown · `huh` forms · `harmonica` spring motion ·
`log` structured logging. Per-component reuse verdicts are in `01-design-system.md`.

---

## 5 · FlowForge Spacing Canon (#2435) — adopted verbatim

Zerado does not re-derive spacing. The tokens and the per-tier matrix below are reproduced from
the canon so a Zerado screen spec can cite them without a cross-repository lookup. **If this
table and the canon ever disagree, the canon is right.**

### 5.1 · The named tokens

| Token | What it does |
|---|---|
| **`OuterMarginX`** | Horizontal inset (cols, each side) between the framed surface and the terminal edge. Applied to the **whole frame, header band included**. |
| **`OuterMarginY`** | Vertical inset (rows, top and bottom). Un-glues the breadcrumb from row 0 and the footer from the last row. |
| **`InnerPaddingX`** | Horizontal padding (cols, each side) inside the frame, between frame boundary and content. |
| **`InnerPaddingY`** | Vertical padding (rows, top and bottom) inside the frame. |
| **`InterElementGap`** | The single respiro row between stacked blocks. Used twice: breadcrumb → title, and header band → body. |
| **`HeaderBandHeight`** | Base rows reserved for the header band **without** a subtitle. |
| **`HeaderBand(tier, hasSubtitle)`** | The effective band. Adds one row at Standard+ **only** when the screen supplies a subtitle. |

### 5.2 · The fixed values (the canon — code must match)

| Token \ Tier | Tiny `<40` | Narrow `40–59` | Standard `60–79` | Wide `80–119` | ExtraWide `120+` |
|---|---|---|---|---|---|
| `OuterMarginX` (cols/side) | **0** | 1 | 2 | **2** | 2 |
| `OuterMarginY` (rows T/B) | **0** | 1 | 1 | **1** | 1 |
| `InnerPaddingX` (cols/side) | 1 | 1 | 1 | **1** | 2 |
| `InnerPaddingY` (rows T/B) | **0** | 1 | 1 | **1** | 1 |
| `InterElementGap` (rows) | 1 | 1 | 1 | **1** | 1 |
| `HeaderBandHeight` (base, no subtitle) | **1** | 3 | 3 | **3** | 3 |
| `HeaderBand(tier, true)` (with subtitle) | **1** | 3 | 4 | **4** | 4 |
| **`leftInset` = `OuterMarginX + InnerPaddingX`** | **1** | 2 | 3 | **3** | 4 |

```
leftInset = OuterMarginX + InnerPaddingX          (per side, cols)

BodyRect.w = width  − 2·leftInset
BodyRect.h = height − HeaderBand(tier, hasSubtitle) − 1 (footer)
                    − 2·OuterMarginY − InnerPaddingY − InterElementGap
```

**The single-sizer invariant holds in Zerado too:** the header's `View` and the body's
`BodyRect` must consult the *same* `HeaderBand(tier, hasSubtitle)`, and `hasSubtitle` must match
whether the screen actually supplies a non-blank subtitle. Disagreement between them is the
defect class that clipped a FlowForge screen once already; it is forbidden by construction.

**Zerado declines the subtitle row on every screen** — a composition decision owned by the spine
(`docs/blueprint/02-composition.md` §1.2). `HeaderBandHeight` is therefore always the 3-row base
and `hasSubtitle` is always `false`, which means the single-sizer desync **cannot occur in Zerado
by construction**: the flag has exactly one value.

---

## 6 · Apple HIG and Material 3 — flagged now, live in Phase 4

Named here so the Phase 4 phone apps are **not designed from scratch later**. They govern
`04-terminal-to-phone-bridge.md` and nothing in Phase 1. Chapter names below were read from
`developer.apple.com` and `m3.material.io` on 2026-08-25; **re-verify at Phase 4 kickoff**,
because both corpora move.

| Concern | Apple HIG | Material 3 |
|---|---|---|
| **Navigation patterns** | Components → *Navigation and search*; Patterns → *Modality* | `components/navigation-bar`, `navigation-rail`, `navigation-drawer`, `app-bars`; `foundations/layout/canonical-examples/list-detail` |
| **Type scales** | Foundations → *Typography* | `styles/typography/type-scale-tokens`, `applying-type`, `fonts` |
| **Dynamic type / font scaling** | Foundations → *Typography* (Dynamic Type) | `foundations/usability` (Android font-scale) |
| **Touch targets** | Foundations → *Layout* | `foundations/layout/grids-spacing/density` — **plus WCAG 2.5.8, 24 × 24 CSS px, AA (verified)** |
| **Dark mode** | Foundations → *Dark Mode*, *Color* | `styles/color/choosing-a-scheme`, `styles/color/roles`, `styles/color/system/how-the-system-works` |
| **Adaptive layout** | Foundations → *Layout* | `foundations/layout/breakpoints/{compact,medium,expanded,large-extra-large}`, `foundations/layout/scaffold` |
| **Platform back behaviour** | Components → *Navigation and search*; Patterns → *Modality* | `foundations/layout/scaffold`; Android system back is a **platform** behaviour — consult Android platform documentation, not M3, for predictive back |
| **Accessibility** | Foundations → *Accessibility*, *Inclusion* | `foundations/building-for-all/{user-needs,co-design}` |
| **Motion** | Foundations → *Motion* | `styles/motion/easing-and-duration/tokens-specs`, `styles/motion/transitions` |
| **Loading / feedback** | Patterns → *Loading*, *Feedback* | `components/progress-indicators`, `components/snackbar` |
| **Onboarding / settings / search** | Patterns → *Onboarding*, *Settings*, *Searching* | `components/search`; `foundations/content-design` |
| **Gestures / input** | Inputs → *Gestures*, *Keyboards*, *Focus and selection* | `foundations/interaction/{gestures,inputs,selection,states}` |
| **Brand colour into a platform scheme** | Foundations → *Color*, *Branding* | `styles/color/static/custom-brand`, `styles/color/advanced/define-new-colors` |

> **A minimum I deliberately did not print.** The commonly-quoted Apple 44 × 44 pt tap target
> is **not stated in the HIG Layout chapter's current published text** (checked at source
> 2026-08-25). Rather than repeat a remembered number, this brief cites the HIG chapter as the
> authority and the **verified** WCAG 2.5.8 figure — 24 × 24 CSS px — as the floor. Confirm
> Apple's current figure at source when Phase 4 begins.

---

## 7 · Reuse-first — and the finding that decides it

The ticket asks for a plain answer on reuse, including where FlowForge's own components do not
fit. The answer turned out not to be aesthetic. It is structural, and it was verified:

| Check | Result (2026-08-25) |
|---|---|
| `github.com/JustCode-CruzAlex/Zerado` — anonymous GitHub API | **HTTP 200** — public |
| `github.com/JustCode-CruzAlex/FlowForge` — anonymous GitHub API | **HTTP 404** — not publicly reachable |
| `proxy.golang.org` version list for `github.com/JustCode-CruzAlex/FlowForge/v3` | **HTTP 404** — module not publicly fetchable |
| `v3/internal/tui/ledgertable` | under `internal/` — **Go forbids cross-module import by construction** |
| `v3/uikit/{space,frame,chrome}` | not `internal/`, but inside the same unfetchable module |

**The verdict: Zerado cannot import a single line of FlowForge Go code.** `LedgerTable` is
doubly barred — it is `internal/`, inside a module that does not resolve anonymously. And a
public Zerado repository depending on a private module would **not build for anyone who cloned
it**, which contradicts the ratified promise that the repository is open (`decisions.md` Q2:
*"the 'it's open' claim on the page is true"*).

**Therefore FlowForge canon is inherited as SPECIFICATION, never as a dependency.** This is
precisely why §5 reproduces the spacing matrix verbatim rather than linking it: the *numbers*
are the transferable artifact. Zerado implements `space`, `frame`, `chrome` and its ledger
primitive natively, against the published values, and is bound by the same acceptance rules.

Reuse-first still governs — it just points at a different shelf. **Zerado's reuse target is the
Charm ecosystem itself** (`bubbles`, `lipgloss`, `viewport`, `huh`, `glamour`), which is public
and is the ratified stack. Per-component verdicts are in `01-design-system.md`; the import
namespace (`charm.land/*/v2` vs `github.com/charmbracelet/*`) is a spine decision for
`fft-tui-architect`, not a design decision, and is deliberately left open here.

---

## 8 · Composition direction

Surface-level intent. It is **not** the layout tree — composition, the pane budget, the focus
state machine, `bubbles` selection and the responsive tier map belong to `fft-tui-architect` in
deliverable A. This section constrains that work; it does not perform it.

1. **80 columns is the design width.** Every screen must be *correct* at 80 and *survive* at 40.
   120+ adds panes; it never rescues a broken 80. A screen that only reads well at 120 has failed.
2. **The ledger is the product.** The library deck is the screen players live in. It obeys
   R-10 (a)(b)(c) at 400 rows, and its identity column carries **a human game title** —
   never an index, never a store ID.
3. **Regions are earned, not assumed.** Per the inherited verdict: the status summary is the
   **one pinned summary row** of R-10(c), not a region. The filter is a **mode** of the list
   below 120 columns. **Two regions only when there is room for both to be correct** — which the
   spine resolves at **120 columns**, not 80: an 80-column split leaves the identity column around
   15 characters and fails R-10(a) on its face. Below 120 the detail is a route, not a pane.
4. **Hierarchy comes from the grid.** No glow, no gradients, no fluid type in the terminal
   (brand §9). Rank is carried by case, weight, colour role, box drawing and spacing — in that
   order — and by nothing else.
5. **Honest states are first-class screens.** First run before anything is synced, a sync that
   fails halfway, a private Steam profile, a 400-game library at 40 columns. These are not edge
   cases; they are the screens that decide whether the product feels solid, and they get the
   same craft as the headline screen.
6. **Offline is normal, not an alarm.** The product publicly promises it works with the network
   off. Offline is therefore rendered in **chrome**, never red, never as a failure — and the
   degrade is always shown, never silent.
7. **One motion.** The scanner is the only signature motion, and only for genuinely
   indeterminate waits. Ambient animation is forbidden — it burns a redraw budget for nothing
   and it is exactly the nostalgia-kitsch the brand rules out.
8. **Cyan is spent, not sprinkled.** The rule is enforceable and countable — `02-colour-budget.md`.

---

## 9 · Values that are unmeasured — flag, do not guess

The brand's derived ANSI-256 table (§5.1) and `tokens.css` cover: surface `232` · raised `234` ·
border `236` · border-strong `67` · text `255` · text-secondary `249` · steel `247` ·
amber `214` · cyan `45` · orchid `177` · scanner `9`.

These tokens are needed by the terminal design system and have **no derived ANSI-256 index**.

> **This list names PRIMITIVES, and that is deliberate.** An ANSI-256 index is derived by
> nearest-neighbour search against a **raw hex value**, and under brand §10's three layers the raw
> value lives in the primitive — the semantic is a `var()` reference to it, carrying meaning
> rather than a value. So a derivation list names primitives by construction.
>
> `01-design-system.md` §1.4 names the **semantic** for the same colours, equally correctly: it is
> a component-facing table, and brand §10 states a primitive is *"never referenced by a
> component."* **The two names are two layers of one token, not two competing vocabularies.**
> The Semantic column below makes the mapping explicit so neither list has to be read against the
> other.

| Primitive (carries the hex) | Semantic (what components reference) | Hex | Needed for |
|---|---|---|---|
| `--z-scanner-300` | *none defined* | `#FF6B6B` | error text (contrast **6.99**, AA — measured) |
| `--z-scanner-900` | `--z-scanner-track` | `#5C1414` | the scanner track |
| `--z-amber-900` | `--z-primary-muted` | `#8A5E00` | unlit/inert progress track |
| `--z-chrome-500` | `--z-text-tertiary` | `#8492A8` | tertiary text (contrast **6.15** — measured) |
| `--z-amber-400` | `--z-primary-hover` | `#FFC94D` | highlight |
| `--z-cyan-300` | `--z-accent-hover` · `--z-text-link-hover` | `#8CF0FF` | highlight |
| `--z-black-700` | `--z-surface-overlay` | `#1D2532` | overlay surfaces |
| `--z-cyan-900` | `--z-accent-muted` | `#0B6C7D` | unlit cyan |
| `--z-orchid-900` | *none defined* | `#4A2A63` | unlit orchid |

Verified against `tokens.css` at source: primitives at lines 46–75, semantics at lines 89–116.
Two rows were corrected while adding this column — `#8492A8` and `#1D2532` were previously listed
under their semantic names (`--z-text-tertiary`, `--z-surface-overlay`), which contradicted the
rule stated just above. Where no semantic exists, the primitive is the only name there is.

**Required action, not optional:** derive each by nearest-neighbour search in CIELAB against the
xterm 256 cube, exactly as brand §5.1 was derived, and add them to `tokens.css` **and**
`tokens.json` in one commit per brand §10. **Nobody may pick these by eye at the keyboard.**
Owner: `fft-brand-architect`. Until they exist, any screen needing one is blocked on the
derivation, not free to improvise.

---

## 10 · The acceptance bar — the yardstick for GOLDEN

Numbered so each line can be failed individually. A Zerado screen is GOLDEN only when **all**
applicable lines pass, evidenced by a rendered artifact rather than a claim.

**Structure and spacing**
1. Nothing is flush against any terminal edge at any tier ≥ Narrow; header-left **equals**
   content-left, verified by column number on an ANSI-stripped render.
2. Every spacing value comes from a named token at its tier value (§5.2). Zero magic numbers.
3. The header band is the **3-row base** with `hasSubtitle` false (Zerado carries no subtitles),
   so no screen pays a body row for a blank fourth line.
4. The breadcrumb is `Zerado ✦ <Screen>` — two segments, no faked empty segment; the
   current-screen segment never truncates.

**Responsiveness**
5. Correct at **80 columns** without horizontal scroll, clipping or overlap.
6. Survives **40 columns**: the body is never starved, the title never disappears.
7. Below its documented floor the screen **refuses and says so in words** — it never renders
   a broken frame.
8. Renders correctly at an **overflowing row count** (≥ 400 games), not only at the golden size.

**The ledger triad (R-10), where the screen is tabular**
9. (a) Every row carries a readable human game title.
10. (b) The selection is always visible; cursor and scroll offset survive a row-set rebuild —
    a background refetch never snaps the view to the top.
11. (c) The summary row is pinned outside the scroll region and is on screen at any row count.

**Colour, state and accessibility**
12. Every state renders **colour and glyph and label** — all three, every row, every tier
    (1.4.1 · R-6).
13. `NO_COLOR` emits **zero** SGR sequences and the screen remains fully unambiguous.
14. The screen is correct at the **16-colour floor**, where `--z-surface` and
    `--z-surface-raised` both collapse to `black` — so **no region is separated by fill**.
15. Every text pair on screen clears **4.5:1**, using ratios read from the brand's measured
    table. No unmeasured pair ships.
16. `--z-border` (1.53) appears only as decoration; every control boundary uses
    `--z-border-strong` (4.08).
17. The focus ring is present on the focused element in **every** state and is never removed.
18. No overlay entirely obscures the focused element (2.4.11).
19. **Escape leaves.** No mode, pane, form or confirmation is a keyboard trap (2.1.2).
20. Single-key shortcuts do **not** fire while a text input holds focus (2.1.4).
21. The chrome-cyan count is **exactly ≤ 1** by the counting method in `02-colour-budget.md`.
22. Red appears only on the closed list in `02-colour-budget.md` §5 — nowhere else.

**Copy and identity**
23. Every string is in the ratified voice: dry, concrete, confident. **No exclamation marks, no
    emoji, the user is never called a "gamer", and the number is stated.**
24. Casing is correct throughout: `Zerado` the product · `zerado` the command · *zerado* the
    status in prose · `ZERADO` the chip.
25. Nothing on screen contradicts `content/landing-copy.md` or `ratification/decisions.md`, and
    no unbuilt capability is presented as working.
26. The screen reads **retro-future**, not retro-nostalgia: clean, lit, pointed forward. No
    distressed texture, no fake CRT damage, no scanline kitsch.
27. It reads **calm, composed, breathing** (R-8). "Calm = done."

**Evidence**
28. A headless render exists at the target viewports **and** at an overflowing row count, has
    been reviewed against this brief, and has been **validated by the founder before merge**.
    No founder-validated screenshot → not GOLDEN.

**Theme** *(added rev A — see `05-theme-system.md`)*
29. Line 15's ratios are measured against **the active theme's own ground**, not against
    `--z-surface`, whenever a theme other than the default is active. A theme that cannot clear
    4.5:1 on its own ground never activates (`05-theme-system.md` §2.2 G2).
30. The screen is correct under **at least one theme from each shipped tier** — brand-true,
    player, monochrome — with co-render intact in all three. Colour is the fast path in every
    tier; it is the answer in none.

---

## 11 · Routing — who builds what

| Work | Owner |
|---|---|
| Guideline authority, this brief, the acceptance bar | `fft-design-architect` |
| **The theme contract, the four-state validation gate, the tiers** (`05-theme-system.md`) | `fft-design-architect` |
| Brand manual, tokens, **the ANSI-256 derivations in §9**, **the light-default CVD repair** (`05-theme-system.md` §3.2) | `fft-brand-architect` |
| The theme **picker** — where it lives, navigation, live preview | `fft-tui-architect` |
| Terminal **composition** — layout tree, pane budget, focus model, Charm adoption, tier map | `fft-tui-architect` |
| Per-screen **visual design** — mockup, hierarchy, co-render, spacing tokens, state tables | `fft-tui-designer` |
| `View` / `Update` line code | `fft-tui` |
| GOLDEN verdict, judged against **this brief** | `fft-code-reviewer` |
| Phase 4 phone **design** | brief by `fft-design-architect` → design by `fft-designer` (interim mobile/native leaf) → build by `fft-flutter`. When Phase 4 design work actually starts, a dedicated `fft-flutter-designer` leaf should be minted rather than assumed. |

**Terminal composition is not this brief's to give.** This document decides *which guidelines
govern* and *what GOLDEN means*; `fft-tui-architect` decides *how the screen is composed*. Two
heads, one product.

---

## 12 · Open for the founder

1. **The ANSI-256 derivation gap (§9).** Nine tokens the terminal needs have no derived index.
   They must be derived, not eyeballed. Confirm `fft-brand-architect` owns this and that it
   lands before the first Phase 1 screen is built.
2. **The contrast claim's honest limit (§3.4).** Confirm that "AA on our ground and on five
   measured popular themes, with `NO_COLOR` as the unconditional fallback" is the claim we
   make publicly — rather than an unqualified "WCAG AA", which would not be true on an
   arbitrary user theme.
3. **Zerado has no in-house gold standard yet (§4, R-1).** Proposal: the library deck becomes
   the reference screen the moment it is ruled GOLDEN, and every later screen rises to it.
   Confirm the nomination.
4. **Portuguese words and assistive technology (§3.2, 3.1.2).** *zerado* and *sinopse* will be
   mispronounced by screen readers in the terminal and this cannot be fixed there. Recorded as
   an accepted cost, consistent with the naming risks already accepted in `naming.md`. Confirm
   it stays accepted.
5. **The theme system (`05-theme-system.md`).** It carries its own founder questions — the
   ΔE floor, the light default's measured failure, a wrong ratio at source, the gate/tier split,
   `retro-82`, and the one exception ever made to "Zerado does not paint a background". Those are
   §8 of that document and are not repeated here.
