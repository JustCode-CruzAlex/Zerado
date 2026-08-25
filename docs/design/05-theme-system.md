---
title: Zerado — the theme system
discipline: DESIGN SYSTEM
doc-no: ZRD-DESIGN-06
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: design-system
ticket: "#2"
---

# Zerado — the theme system

**A theme is data. Zerado is the code.** This document decides what a theme file must supply for
every Zerado component to render correctly, what Zerado derives, what no theme may touch, and
what happens when a palette cannot express something the product needs. It then makes the
founder's rule runnable: **a theme that cannot express four distinguishable states is INVALID and
fails validation rather than shipping broken.**

The direction it implements, verbatim:

> *"bring all the themes that we already build on FlowForge, we will follow that, so we have light
> themes and dark themes, what we will do is create a default theme based on the mocks on our
> site. And we can look on omarchy themes and bring others that will make sense on zerado."*

**Companion documents:** `00-design-brief.md` (which guidelines govern) · `01-design-system.md`
(the components) · `02-colour-budget.md` (the colour rule) · `03-designer-manual.md` (how we work)
· `04-terminal-to-phone-bridge.md` (Phase 4).

**Scope.** This is the theme *contract* and its validation. The theme **picker** — where it lives,
how it is navigated, whether it previews live — is terminal **composition** and belongs to
`fft-tui-architect`, not here. This document constrains that work; it does not perform it.

---

## 0 · Attribution of every number in this document

The standing rule from `00-design-brief.md` §2 applies without exception. Two classes of figure
appear below and they are marked differently.

| Class | Source | How it is marked |
|---|---|---|
| **Contrast ratios of brand tokens** | Read from `brand-manual.md` §4.2 / §4.3 / §4.5 and `tokens.css` §10 | Cited as the manual's |
| **Everything else** — the 35-theme sweep, every ΔE in this document, every ratio on a theme's own ground | **Computed here** by the method pinned in §2.1 | Labelled *computed* |

**One verification worth recording before anything else.** Every contrast ratio in the brand
manual's tables was recomputed here from its hex values using the WCAG 2.2 relative-luminance
formula. **Thirty-four of thirty-five reproduce exactly.** The contrast method is therefore
reproducible and the manual's contrast digits can be relied on as printed. The single exception is
recorded in §3.3.

The manual's **ΔE** digits are a different story, and §2.1 is the response to it.

---

## 1 · The token contract — a theme is DATA, never code

### 1.1 · The load-bearing question

A Zerado component references semantic tokens: `--z-primary`, `--z-accent`, `--z-state-*`,
`--z-border-strong`, `--z-scanner`, `--z-text-secondary`, and about a dozen more. A theme file
carries roughly eighteen colours. **The mapping is not one-to-one and pretending it is would be
the whole defect.**

The contract is a three-way split. Every token Zerado renders belongs to exactly one class, and
which class it belongs to determines what a theme may do to it.

| Class | Definition | Who owns the value | May a theme change it? |
|---|---|---|---|
| **SUPPLIED** | Read from the theme file, or selected from it by the §2.4 rule | The theme author | **Yes** — this is what a theme *is* |
| **DERIVED** | Computed from supplied values by a rule pinned in this document | Zerado's loader | **Indirectly** — it moves when the supplied values move |
| **PINNED** | Theme-invariant. Identical in every theme, forever | The brand | **No.** A theme that tries is ignored, not honoured |

Everything in §1.6 — the repair-then-refuse behaviour — follows from this split, because *the
class decides the remedy*.

### 1.2 · What a theme file actually is — and a defect worth inheriting

**Verified at source, 2026-08-25:** FlowForge ships **35** theme files at
`v3/uikit/theme/themes/omarchy/*.toml` plus `ATTRIBUTION.md`. They are flat key/value TOML — no
tables, no nesting.

**They are not, however, one dialect. There are two, and the difference has already caused a
production defect in FlowForge.**

| Dialect | Keys | Count in the corpus |
|---|---|---|
| **ANSI** | `background`, `foreground`, `accent`, `cursor`, `selection_*`, `color0`…`color15` | **32** of 35 |
| **Semantic** | `background`, `foreground`, `accent`, `mode`, `selection`, `muted`, `red`/`green`/`yellow`/`blue`/`magenta`/`cyan` and their `bright_*` siblings, `*_background`, `*_foreground` | **3** of 35 — `last-horizon`, `lupine`, `solitude` |

Both dialects define `background`, `foreground` and `accent`, so **a semantic-dialect file parses
with no error and yields sixteen empty ANSI slots.** FlowForge hit exactly this and documented it
(`v3/uikit/theme/omarchy.go`, the `#5368` note: *"a new-schema file parsed with no error and
simply left all 16 ANSI slots empty — so the theme fell back to generic ANSI and rendered washed
out, silently"*).

> **Design decision — Zerado implements both dialects, and a file that yields fewer than four
> usable colours is a load error, not a silent fallback.** The mapping is FlowForge's
> `fillFromSemantic`, adopted as *specification* per `00-design-brief.md` §7 — Zerado cannot
> import the Go, and does not need to:
> `color0←darker_background · color1←red · color2←green · color3←yellow · color4←blue ·
> color5←magenta · color6←cyan · color7←light_foreground · color8←muted · color9←bright_red ·
> color10←bright_green · color11←bright_yellow · color12←bright_blue · color13←bright_magenta ·
> color14←bright_cyan · color15←bright_foreground`, additive only — an explicit `colorN` always
> wins.
> **Reason:** silence is the failure mode that costs the most. FlowForge's version of this bug
> shipped and was invisible until someone measured. A hard load error costs one line of output;
> a silent fallback costs a wrong screen nobody reports.

The `mode` key (`"dark"` / `"light"`) appears **only** in the three semantic-dialect files.
Polarity therefore cannot be read from the file — see §3.5.

### 1.3 · SUPPLIED — what a theme actually gives Zerado

Nine values. Everything else is derived or pinned.

| Semantic token | Supplied from | Notes |
|---|---|---|
| *the measuring ground* | `background` | **Not paint** — see §3.6. It is the ground every ratio in this theme is measured against |
| `--z-text` | `foreground` | The one value a theme author guarantees is legible on their own ground |
| `--z-state-not-started` | §2.4 selection | \| |
| `--z-state-in-progress` | §2.4 selection | \| the four-state set — the only tokens whose values are *chosen* from the palette |
| `--z-state-zerado` | §2.4 selection | \| rather than read from a fixed key |
| `--z-state-abandoned` | §2.4 selection | \| |
| `--z-primary` | ≡ `--z-state-in-progress` | **Bound, not separately supplied** — see below |
| `--z-accent` | ≡ `--z-state-zerado` | **Bound, not separately supplied** — see below |
| `--z-scanner` | §2.4, red band | The alarm colour. Not a state — see §2.2 |

> **Design decision — `--z-primary` and `--z-accent` are bound to two of the four states, not read
> from the theme's `accent` key.**
>
> **Reason.** Brand §4.1 is the product's own argument: *one colour is common, one is earned* —
> amber is the ambient voice, cyan is spent on completion. In the default palette
> `--z-primary` **is** `--z-state-in-progress` and `--z-accent` **is** `--z-state-zerado`; they are
> the same hex. Binding them keeps that argument structurally true in every theme: whatever a
> theme's warm state is becomes its ambient voice, whatever its earned state is becomes its
> accent. Reading the theme's own `accent` key instead would break it immediately — `retro-82`'s
> `accent` is `#FAA968`, an orange, which would make the *earned* colour and the *ambient* colour
> the same hue family and destroy the only signal the palette has.
>
> **What is lost, stated honestly:** a theme author's declared accent is ignored. That is the
> correct trade. The accent belongs to the desktop; the earned colour belongs to Zerado.

**`--z-text-secondary` is bound to `--z-state-not-started`.** In the default they are two
near-neighbours (`#A9B5C7` and `#A5A29B`) that brand §5.1 already collapses to the same slot —
both resolve to `white` at the 16-colour floor. Since the collision is already ratified at the
floor, making it explicit in a theme costs nothing: co-render means a state is never carried by
colour alone, so secondary text sharing the dormant colour loses no information.

### 1.4 · DERIVED — computed, never authored

These have no key in any theme file and never will. Each is produced by a pinned rule from the
supplied values.

| Token | Derivation rule | Floor it must clear |
|---|---|---|
| `--z-surface`, `--z-surface-raised`, `--z-surface-overlay` | the ground, blended toward its own elevation pole in fixed steps | *(measuring grounds only — §3.6)* |
| `--z-text-tertiary` | `--z-text` blended toward the ground until it first clears its floor | **4.5:1** |
| `--z-border` | the ground blended toward `--z-text` at the decorative step | none — decorative by definition |
| `--z-border-strong` | the ground blended toward `--z-text` until it first clears its floor | **3:1** (WCAG 1.4.11) |
| `--z-primary-muted`, `--z-scanner-track`, `--z-accent-muted` | the corresponding supplied colour blended toward the ground | none — inert tracks |
| `--z-scanner-300` (error text) | `--z-scanner` moved along L\* until it first clears its floor | **4.5:1** |
| **the ANSI-256 index of every supplied and derived colour** | nearest-neighbour search against the xterm 256 cube, measured as CIEDE2000 in CIELAB | — |

> **Design decision — a theme's ANSI-256 indices are computed at load; the default's are
> authored.** Brand §5.1 derived and **ratified** eleven indices for the default palette. Those
> are canon and are never recomputed. A theme, by contrast, cannot carry hand-derived indices —
> it is data, and there is no ratification step for a file a contributor drops in. So the default
> reads its indices from `tokens.css`, and every other theme computes them by brand §5.1's own
> method. Nobody picks one at a keyboard either way (brand §10 rule 5).

> **Design decision — at the 16-colour floor, every theme collapses to the default's ratified
> mapping.** Brand §5.2 pins the four states to `white` (7), `bright yellow` (11), `bright cyan`
> (14), `bright magenta` (13) with no collisions. At sixteen colours there is no palette to
> express — only slots, and those slots are already coloured by the user's own terminal theme.
> **At 16 colours the user's terminal theme IS the theme**, and Zerado's theme selection becomes a
> no-op. This is the honest answer and it costs nothing: `02-colour-budget.md` §7 already forbids
> any layout that depends on fill, so nothing breaks when the palette flattens.

### 1.5 · PINNED — theme-invariant, forever

**No theme may touch any of these.** They are the identity, and they are the reason a Zerado
screen in `hackerman` is still recognisably Zerado.

- **The four glyphs and their ASCII fallbacks** — `○ ◐ ◉ ⊘` and `[ ] [~] [*] [x]` (brand §4.3).
  The single filled/struck progression is the state system's spine.
- **The four labels** — `NOT STARTED` `IN PROGRESS` `ZERADO` `ABANDONED`, and the casing
  convention around them (`naming.md`).
- **The co-render rule itself** — colour AND glyph AND label, every state, every row, every tier.
- **The focus channels** — the `▌` gutter marker, bold, and the ring; and the rule that the ring
  is never removed (`02-colour-budget.md` §8.2).
- **Every spacing value** — the Spacing Canon (#2435), adopted verbatim (`00-design-brief.md` §5).
- **The type roles** — uppercase display, sentence-case voice, the readout register, and the
  no-letterspacing decision (`01-design-system.md` §1.5).
- **The scanner's motion** — its timing, its single-pip form, and the closed list of when it may
  appear (`01-design-system.md` §9).
- **The colour budget** — one chrome-cyan per screen, the amber allow-list, the red closed list
  (`02-colour-budget.md`). A theme changes which hex is "the earned colour"; it never changes how
  many times that colour may be spent.

A theme file that contains a key attempting any of the above is **ignored, and the ignore is
reported** by the validator. Not honoured, not an error — reported, so a contributor learns the
boundary rather than guessing at it.

### 1.6 · When a theme cannot express a token — repair, then refuse

The question the founder's rule actually turns on. The answer is **two stages with different
remedies**, and the class from §1.1 decides which applies.

| Stage | When | What happens | Class it applies to |
|---|---|---|---|
| **1 · REPAIR** at load | A **derived** token misses its floor | The token is moved along its pinned axis by the **minimum** distance that clears the floor, and the movement is recorded | DERIVED only |
| **2 · REFUSE** at activation | A **supplied** token, or the four-state set, misses the §2.2 gate | The theme does not activate. The named reason is printed. **The previously active theme is left untouched** | SUPPLIED only |

**Why the asymmetry, and it is the whole design.** Repairing a *derived* token is not a change —
it is finishing a computation Zerado was already performing; the theme author never authored that
value and cannot be surprised by it. Repairing a *supplied* token is different in kind: it
silently rewrites the theme author's identity, and worse, it hides the defect. A theme whose
palette cannot carry the four states is not improved by Zerado inventing a fifth colour and
calling it the author's — it is a theme that does not fit this product, and saying so is the
service.

**A refusal never falls back.** It leaves the user looking at the theme they could read a moment
ago. Swapping to some "safe" default would be a second surprise stacked on the first.

> This two-stage shape is FlowForge's, adopted as specification. FlowForge repairs at
> **registration** and audits at **activation** (`v3/uikit/theme/registry.go`, the `#5369` gate),
> with the same rationale recorded in its own words: *"Repair and audit are deliberately two
> functions rather than one. Repair returns a theme and cannot fail; audit returns a verdict and
> CAN. A repair that falls short is indistinguishable from one that worked unless something later
> checks."* Zerado inherits the architecture, not the code (`00-design-brief.md` §7).

---

## 2 · The four-state contract — INVALID is a verdict, not a warning

> **The founder's rule, made runnable:** a theme that cannot express four *distinguishable* states
> is **INVALID** and **fails validation rather than shipping broken.**

### 2.1 · The pinned method — this is what makes the rule runnable at all

You cannot set a validation floor with a number produced by an unpinned method. The brand manual
names its CVD model — *Viénot, Brettel & Mollon (1999), dichromat simulation in linear RGB,
separation measured as CIEDE2000 in CIELAB* — but it does not pin the **matrix variant**, the
**white point**, or the **gamut-clamping rule**, and those choices move the third significant
figure. `00-design-brief.md` §2 already recorded that an independent implementation disagreed with
the manual's digits. **Two independent implementations of the manual's own named method, run here
on the same hexes, disagreed with each other too.** That is not a criticism of the manual; it is
the reason this section exists.

**The pinned method. Every ΔE in this document, and every ΔE the validator prints, is produced by
exactly this and by nothing else.**

1. **sRGB → linear RGB** by the IEC 61966-2-1 transfer function (`c/12.92` below 0.04045,
   `((c+0.055)/1.055)^2.4` above).
2. **Linear RGB → LMS** by the Hunt–Pointer–Estévez matrix normalised to D65, in the form Viénot
   et al. publish:
   `L = 17.8824·R + 43.5161·G + 4.11935·B` ·
   `M = 3.45565·R + 27.1554·G + 3.86714·B` ·
   `S = 0.0299566·R + 0.184309·G + 1.46709·B`
3. **Dichromat projection.**
   *Protanopia:* `L′ = 2.02344·M − 2.52581·S`, `M′ = M`, `S′ = S`.
   *Deuteranopia:* `L′ = L`, `M′ = 0.494207·L + 1.24827·S`, `S′ = S`.
4. **LMS → linear RGB** by the inverse:
   `R = 0.080944·L − 0.130504·M + 0.116721·S` ·
   `G = −0.0102485·L + 0.0540194·M − 0.113615·S` ·
   `B = −0.000365294·L − 0.00412163·M + 0.693513·S`
5. **Clamp to \[0,1\] in linear RGB**, then encode back to sRGB. *(The clamp is where two
   implementations most often part company. It is pinned here.)*
6. **sRGB → CIELAB** with the **D65** white point `Xn=95.047, Yn=100.000, Zn=108.883`.
7. **CIEDE2000** with `kL = kC = kH = 1`.

**Tritanopia is deliberately out of scope**, consistent with brand §4.4 which simulated protanopia
and deuteranopia only. Recorded as a known limit, not an oversight.

**Contrast** uses the WCAG 2.2 relative-luminance formula unchanged; §0 records that it reproduces
the manual exactly, so no pinning is needed there.

**Under the pinned method, the dark default's own six pairs measure:**

| Pair | Normal | Protanopia | Deuteranopia | **Worst** |
|---|---|---|---|---|
| zerado × abandoned | 43.44 | 22.83 | **11.81** | **11.81** |
| not-started × zerado | 28.23 | **25.21** | 27.77 | 25.21 |
| not-started × in-progress | 28.11 | **26.42** | 27.95 | 26.42 |
| not-started × abandoned | 33.25 | 32.89 | **31.50** | 31.50 |
| in-progress × zerado | **50.32** | 50.74 | 57.79 | 50.32 |
| in-progress × abandoned | **63.04** | 70.49 | 68.26 | 63.04 |

*All computed here.* The manual records **11.9** for the tightest pair; the pinned method returns
**11.81**. Same pair, same ordering, same conclusion — 0.09 apart. Every ordering claim brand §4.4
makes survives the pinning intact; only digits move, which is exactly what §2 of the design brief
predicted.

### 2.2 · The gate — five lines, each independently failable

A theme is **VALID** only when all five pass. Every line is a measurement, not a judgement.

| # | Line | Threshold |
|---|---|---|
| **G1** | The four state colours are **four distinct values** | exact inequality |
| **G2** | Each state colour clears contrast against **the theme's own ground** | **≥ 4.5:1** |
| **G3** | `--z-text` (the theme's `foreground`) clears contrast against its own ground | **≥ 4.5:1** |
| **G4** | **All six** state pairs clear separation under **all three** vision models | **ΔE2000 ≥ 10.0** |
| **G5** | **No state colour falls in the red band** — CIELAB `C* > 18` and `h ∈ [345°, 45°)` | exact exclusion |

**G5 is the one line that is a rule rather than a measurement, and it is not taste.**
`02-colour-budget.md` §5 is a ratified closed list: red is motion and alarm, and *"never: the
`ABANDONED` state (that is orchid)"*. A theme that assigned a state to the alarm colour would make
the alarm mean "abandoned" on every row of the ledger. Reserving the band is how a settled
decision survives a data file. The band is centred on the brand's own scanner hue — `#FF2E2E`
measures **h 34.99°** and its paper sibling `#C61F1F` measures **h 35.40°**, both computed here.

**The scanner is not a state and does not enter G4.** It is selected from the red band by the same
§2.4 rule, must clear **3:1** on the ground (it is motion and structure, not text — brand §4.2),
and must not be identical to any state colour. It is excluded from the six-pair matrix
deliberately: the scanner never occupies the ledger's state column, so it is never adjacent to a
state glyph the way two states are adjacent to each other. **A theme with no red-band colour at
all is still valid** — `hackerman`, `lumon`, `lupine` and `white` have none — and renders the
scanner and the alarm annunciator uncoloured plus bold, the same interim rendering
`01-design-system.md` §11.2 already specifies for error text.

### 2.3 · The floor is **10.0**, not 11.9 — and the reason matters

The founder's question was whether the dark default's floor becomes the bar for every theme. **It
does not, for two reasons, and the second is the stronger one.**

**First: 11.9 is a figure from an unpinned method.** It cannot be a runnable threshold because a
second implementation of the same named method returns 11.81 (§2.1). A gate whose pass/fail flips
on which library you linked is not a gate.

**Second, and decisive: a floor set at the default's own measurement is a tautology, not a floor.**
It would mean the reference theme passes by *definition*, with zero headroom, and fails the moment
anything about it is re-measured or a token is nudged. Worse, the bar would move every time the
brand touched the default palette — the standard would be defined by the thing being measured.

**10.0 is the brand manual's own stated scale**, not a number invented here: *"ΔE ≥ 10 is distinct;
≥ 15 is comfortable."* Adopting it makes the default clear its own gate with real headroom (11.81
against 10.0) instead of defining it, and it keeps the interpretive language the manual already
taught.

**The two bands are the manual's, and they are load-bearing:**

| Band | Range | What it means | Ships? |
|---|---|---|---|
| **COMFORTABLE** | ΔE ≥ 15 | Colour alone separates the pair for everyone | yes |
| **DISTINCT** | 10 ≤ ΔE < 15 | Colour separates the pair; at the tight end glyph and label carry real load | yes |
| **INVALID** | ΔE < 10 | Two states can look alike. Refused | **no** |

**The dark default sits at DISTINCT (11.81), and that is not an embarrassment — it is the manual's
own position**, stated in brand §4.4: the tightest pair *"is the one place where glyph and label
genuinely carry load rather than merely reinforce. That is what they are there for."* A product
whose reference theme sits in the DISTINCT band is a product that has to keep co-render honest,
which is precisely the discipline this design system is built on.

### 2.4 · The selection rule — deterministic, and checkable by hand

A theme supplies a **palette**, not a state set. The four states are *selected from* it. The rule
is fixed so that two implementations, and a reviewer with a calculator, reach the same answer.

**Candidate pool.** The theme's `color0`…`color15`, plus `accent`, plus `foreground` —
deduplicated by value, in that order. A candidate is **eligible** only if it clears **4.5:1** on
the theme's ground and is **not** in the red band (G5).

**Reference set.** Every role has a reference colour: the **dark default's** four states for a dark
theme, the **light default's** four for a light theme (§3). Distance to the reference is CIEDE2000
in CIELAB.

**Role classes.**

| Role | Class | Test |
|---|---|---|
| `not-started` | **least chromatic** | `C*(not-started) ≤ C*` of each of the other three |
| `in-progress` | warm band | `C* > 18` and `h ∈ [45°, 105°)` |
| `zerado` | cool band | `C* > 18` and `h ∈ [175°, 265°)` |
| `abandoned` | violet band | `C* > 18` and `h ∈ [265°, 345°)` |

The three chromatic bands are derived from the brand's own state hues and the brand sits
comfortably inside all of them — computed here: chrome `C* 4.02`, amber `h 77.28°`, cyan
`h 219.49°`, orchid `h 313.81°`; and on paper `C* 10.45`, `h 67.78°`, `h 222.59°`, `h 313.62°`.
The gaps between the bands are not arbitrary either: **0–45° is reserved for the alarm** (G5) and
**105–175° is green, which has no role in Zerado.**

**The `not-started` rule is a class, not a band, and that is deliberate.** *Not started* is the
absence of activity; a state that out-shouts an active one is a lie the product tells on every row
of the ledger. Stating it as a *relative* test — least chromatic of the four — rather than an
absolute chroma ceiling is what lets the rule hold for a monochrome palette too, where a ceiling
would simply refuse everything.

**Choice.** Among all assignments of four distinct eligible candidates that satisfy G1–G5, take
the one that (a) puts the **most** roles inside their band, then (b) minimises **total** distance
to the reference set, then (c) maximises the **worst** pair separation. Ties break by pool order,
which is fixed. If no assignment satisfies G1–G5, the theme is **INVALID** and the validator
reports *which* line failed and the best value it could reach.

### 2.5 · The three tiers — the gate is measurement, the tier is brand

Passing the gate makes a theme **safe**. It does not make it **Zerado**. Those are different
questions and collapsing them would smuggle taste into a measurement.

| Tier | Condition | What it means |
|---|---|---|
| **BRAND-TRUE** | All four roles in band | The product's own colour argument survives intact: something neutral is dormant, something warm is ambient, something cool is earned, something violet is abandoned |
| **PLAYER** | The three active roles are chromatic (`C* > 18`) but one or more is hue-rotated | A real theme with a legible departure. The deviating roles are **named** in the credits |
| **MONOCHROME** | One or more active role separates by **lightness** rather than hue | Valid, and in one respect the *safest* tier — lightness survives every dichromat model. The colour cue is weak, so glyph and label carry more |

**MONOCHROME is a first-class tier, not a consolation.** Brand §5.4 already states the position:
*"A Zerado screen with colour stripped is not a degraded Zerado screen — it is a correct one."* A
theme whose four states are four greys at ΔE ≥ 10 is that argument with the colour still on.

### 2.6 · What co-render buys — why a marginal theme is *degraded*, not *dangerous*

Glyph and label are PINNED (§1.5). They are identical in `delorean`, in `white`, in the default,
and under `NO_COLOR`. **This is the entire reason a theme in the DISTINCT band, or in the
MONOCHROME tier, is shippable at all.**

Trace the worst case. A deuteranopic user, on a DISTINCT-band theme, at the theme's tightest pair.
Colour has degraded to *nearly the same* — ΔE just over 10. What is on the row:

```
  ◉  ZERADO        Return of the Obra Dinn     9h
  ⊘  ABANDONED     Sekiro                      3h
```

A ring with a solid core against a ring struck through, and two words that share no letters.
**The colour was the fast path. It was never the answer.** So a marginal theme costs *speed of
recognition* and costs **nothing** in correctness — the user reads the row slightly slower, and
reads it right. That is a degradation.

Now delete co-render and the same theme becomes dangerous: two rows that a colour-blind user
cannot tell apart at all, on the most-used component in the product, with no second channel. The
distance between "degraded" and "dangerous" is exactly one pinned glyph.

**This is also why the gate can afford to be generous at 10.0 rather than defensive at 15.** A
product without co-render would need a much higher floor and would still be one palette away from
an unreadable screen.

### 2.7 · What the validator emits

Machine-readable, one record per theme, generated — never hand-written:

- the **verdict** (`VALID` / `INVALID`) and, when invalid, the **failing line** (G1–G5) and the
  best value reached;
- the **tier** and the **band**;
- the **selection**: for each role, the slot name, the hex, its contrast on the theme's ground,
  and its `L* C* h`;
- the **six-pair matrix**: normal / protanopia / deuteranopia and the worst of the three;
- the **derived-token repairs** applied at load (§1.6 stage 1), each with its before and after;
- any **PINNED key** the file attempted, ignored and reported (§1.5).

This record is what §5 requires a contributor to ship, and what makes the tier a *checked* claim
rather than an opinion in a table.

---

## 3 · The default theme, and the light default

### 3.1 · The dark default is the brand — it is not a new decision

Zerado's dark default is authored from the site mocks, which is to say it is the brand manual's
palette. **The theme file is its expression, not a re-decision.** Every value below is the
manual's; nothing here is chosen.

| Role | Token | Hex | Ratio on `--z-surface` | Source |
|---|---|---|---|---|
| Ground | `--z-bg` | `#05060A` | — | brand §4.2 |
| Surface | `--z-surface` | `#0B0D14` | — | brand §4.2 |
| Ink | `--z-text` | `#E9EEF5` | **16.65** AAA | brand §4.2 |
| Ambient / in progress | `--z-primary` | `#FFB000` | **10.59** AAA | brand §4.2 |
| Earned / zerado | `--z-accent` | `#19E0FF` | **12.15** AAA | brand §4.2 |
| Dormant / not started | `--z-state-not-started` | `#A5A29B` | **7.62** AA | brand §4.3 |
| Abandoned | `--z-state-abandoned` | `#C77DFF` | **7.21** AA | brand §4.3 |
| Alarm | `--z-scanner` | `#FF2E2E` | 5.25 — motion, not text | brand §4.2 |
| Control boundary | `--z-border-strong` | `#64748B` | **4.08** — meets 1.4.11 | brand §4.2 |

**Gate verdict: VALID · BRAND-TRUE · DISTINCT (ΔE 11.81).** The default is the only theme whose
ANSI-256 indices are read rather than computed (§1.4), and the only one whose values arrive by
ratification rather than by selection.

### 3.2 · The light default is brand §4.5 — and it currently **FAILS** its own gate

**Light is first-class, not a fallback**, and the brand already designed it: §4.5 is a second
authored expression with measured ratios, implemented in `tokens.css` §10 under
`[data-z-surface="paper"]`. Basing the light default on anything else would waste the work that
made §4.5 exist.

The shipped light state set, read at source from `tokens.css` §10:

| Role | Hex | Ratio on `#FFFFFF` | Ratio on `#F3F5F8` |
|---|---|---|---|
| not started | `#5E6A7A` | 5.50 *(manual)* | 5.04 *(manual, §4.5)* |
| in progress | `#8A4F00` | 6.56 *(manual)* | 6.01 *(computed)* |
| zerado | `#0A6070` | 7.19 *(manual)* | 6.58 *(manual, §4.5)* |
| abandoned | `#6D3D93` | **7.67** *(computed — see §3.3)* | 7.02 *(computed)* |

**Contrast is fine. Separation is not.** `04-terminal-to-phone-bridge.md` §8.1 recorded this set as
CVD-**unverified** and flagged it for measurement. It has now been measured, under the §2.1 pinned
method, and **it fails**:

| Pair | Normal | Protanopia | Deuteranopia | **Worst** | |
|---|---|---|---|---|---|
| not-started × zerado | 15.66 | **5.41** | 8.73 | **5.41** | ✗ **FAIL** |
| zerado × abandoned | 26.88 | 15.47 | **8.91** | **8.91** | ✗ **FAIL** |
| not-started × abandoned | 17.27 | 19.20 | **14.97** | 14.97 | pass |
| not-started × in-progress | **32.89** | 34.91 | 36.47 | 32.89 | pass |
| in-progress × zerado | 39.77 | **36.42** | 42.73 | 36.42 | pass |
| in-progress × abandoned | **47.32** | 55.72 | 54.43 | 47.32 | pass |

*All computed here.* **Floor 5.41 against a required 10.0.** Two of six pairs fail, and the worse
one fails by more than half.

> **This is the point of building the gate.** "Someone should check the light palette" has been an
> open note since `04-terminal-to-phone-bridge.md` §8.1 was written. It is now a **blocking
> validation failure with a name, a number and an owner** — and the general rule that produced it
> is *no theme ships without passing this*, which means the same class of gap cannot reopen
> silently in any future theme.

**Consequence, and it is binding: no Zerado light theme ships until this is repaired.** Owner:
`fft-brand-architect`, through brand governance §10 — not through a screen PR, not through this
document.

### 3.3 · Two errors found at source

**(a) One contrast figure is wrong.** `tokens.css` §10 records
`--z-state-abandoned: #6D3D93; /* 7.30:1 */`. Recomputed against `#FFFFFF` by the same WCAG formula
that reproduces all thirty-four other figures exactly, `#6D3D93` measures **7.67:1**.
`04-terminal-to-phone-bridge.md` §8.1 reproduced the 7.30 in good faith and inherits the
correction. Nothing ships broken — the true value is *higher*, so the pairing was never at risk —
but the figure is wrong and is corrected here and in that document.

**(b) One hue was not carried, and it is the one where carrying it was load-bearing.** Brand §4.5
states its own principle plainly: *"These are not different colours; they are **the same hues
carried to ink weight**."* Measured — every hue below computed here:

| Role | Screen `h` | Paper `h` | Δ | Carried? |
|---|---|---|---|---|
| abandoned | 313.81° | 313.69° | **0.12°** | yes |
| scanner | 34.99° | 35.40° | **0.41°** | yes |
| zerado | 219.49° | 222.59° | **3.11°** | yes *(manual says 3.1° — exact)* |
| in progress | 77.28° | 67.78° | **9.50°** | yes *(manual says 9.7°)* |
| **not started** | **92.97°** | **265.99°** | **173.02°** | **no — rotated to the opposite side of the wheel** |

Four of five carry to within 9.5°. The fifth is inverted — and it lands in the **blue-cast steel
family that brand §4.4 explicitly rejected on the dark side**, after simulation collapsed
`#9FB0C6` (`h 264.88°`, computed) against the cyan. The paper `#5E6A7A` measures `h 265.99°`:
**1.1° from the rejected colour.** The dark set's own hard-won correction — *"The warmth is
load-bearing engineering, not taste, and must not be 'corrected' back toward blue"* — was not
carried to paper. The 5.41 failure in §3.2 is that omission, measured.

### 3.4 · The repair — a feasibility proof, explicitly **non-binding**

A gate nobody can pass is a bug in the gate. This proves the light gate is satisfiable **without
inventing anything**: both moves below are rules already ratified in canon, applied to paper.

**Move 1 — carry the warm-neutral hue, as §4.5 says it does.** Hold `L*` and `C*` of `#5E6A7A`
exactly (44.37 / 10.44) and rotate the hue to the dark set's `92.97°`.
**Move 2 — restore the lightness spread between zerado and abandoned.** The dark pair sits 17.2
`L*` apart; the paper pair sits **1.7** apart. Darken `abandoned` by `ΔL* −8`, hue and chroma held.

| Role | Current | Proposed | Ratio on `#FFFFFF` | Change |
|---|---|---|---|---|
| not started | `#5E6A7A` | `#6F6958` | 5.50 → **5.47** | hue 265.99° → 94.12°; `L*`/`C*` held |
| in progress | `#8A4F00` | *unchanged* | 6.56 | — |
| zerado | `#0A6070` | *unchanged* | 7.19 | — |
| abandoned | `#6D3D93` | `#582A7E` | 7.67 → **10.26** | `ΔL* −8` |

Six-pair result, computed here: **floor 10.83** (was 5.41), worst pair zerado × abandoned under
deuteranopia. `not-started × zerado` moves **5.41 → 22.16**, a 4.1× improvement — the same order of
improvement the dark set's own warm-grey correction produced. Three of four contrasts unchanged;
the fourth **rises**.

> **This is not a token change and must not be treated as one.** No surface may adopt these hexes
> from this document. It is a demonstration that a minimum-motion repair exists along axes canon
> already ratified, handed to `fft-brand-architect` as input to brand governance §10 — which
> requires the change to land in `tokens.css` **and** `tokens.json` **and** the manual in one
> commit, with the ANSI-256 indices re-derived (`00-design-brief.md` §9). Any value that ships will
> be that process's output, not this appendix's.

**One measured note to carry into that work.** In the light set, `in-progress` × `scanner` separate
by only **ΔE 5.68** under deuteranopia (`#8A4F00` × `#C61F1F`, computed; the dark pair is 14.81).
The scanner is excluded from G4 for the reasons in §2.2 and this is not a gate failure — but it is
a real, measured weakness of the paper expression and it should be looked at while the set is open.

### 3.5 · How the default is selected

Four steps, in strict order. The first that answers, wins.

| # | Condition | Result | Why |
|---|---|---|---|
| **1** | **`NO_COLOR` is set** | **No theme at all.** Zero SGR sequences | Brand §5.4 is unconditional: *"not a reduced palette, not a 'safe' subset. Zero."* `NO_COLOR` is not a theme choice, it is the absence of the whole colour layer, and it outranks everything |
| **2** | An **explicit setting** — a persisted preference or `--theme <name>` | That theme, subject to the §2.2 gate | The user asked. A refusal here is loud (§1.6) rather than a silent substitution |
| **3** | **Terminal background detection** — OSC 11, then `COLORFGBG` | The **dark default** or the **light default** — never a named theme | Detection answers one question (is this ground light or dark) and is only trusted to decide that one thing |
| **4** | Nothing answered | **The dark default** | Below |

**Defending step 4.** Zerado is dark-native on every surface (`04-terminal-to-phone-bridge.md`
§2.5); brand §5.3 measured the palette against five popular **dark** grounds and every state
cleared 4.5:1 on all five; and the overwhelming majority of terminals ship dark. Guessing light on
a failed probe would guess against the brand *and* against the only measured evidence there is.
`00-design-brief.md` §3.4 already records the honest limit of that position and this changes
nothing about it: we claim AA on our own ground and on five measured grounds, never on an arbitrary
one, with `NO_COLOR` as the unconditional fallback.

> **Design decision — detection may switch polarity, never identity.** A probe result is allowed to
> choose between the dark and light *defaults*. It is never allowed to pick, or silently swap, a
> named theme. A theme is an identity choice and identity is not something a terminal escape
> sequence gets to make on the user's behalf.

### 3.6 · The paint rule — the one narrow exception to §1.3

`01-design-system.md` §1.3 decided that Zerado paints no background: the surfaces are *reference
grounds for measuring contrast*, not paint, because *"a TUI that paints its own full-screen
background looks like it is squatting in someone's terminal."* A theme system pressures that rule —
an ink-only light theme rendered on a dark terminal is unreadable — so the rule needs one exception,
and exactly one.

> **Design decision — Zerado paints the theme's declared `background` if, and only if, the user
> has explicitly selected a named theme (step 2 above). In every other case it paints nothing.**
>
> **Reason.** §1.3's argument is about *consent*, not about pixels. Painting over a terminal nobody
> asked you to paint is squatting; painting the ground of the theme the user just chose by name is
> delivering what they asked for — an ink-only `delorean` on a `gruvbox` terminal is neither
> `delorean` nor Zerado. The defaults (steps 3 and 4) stay ink-only, so the product's normal
> posture is unchanged and brand §5.3's survive-the-user's-background work is never discarded.
>
> **What is unchanged, and it is binding:** *"elevation is carried by borders and spacing, never by
> fill."* Painting a single flat ground is not elevation. No region may be separated from another
> by fill in any theme, at any colour depth — the rule that survives the 16-colour floor
> (`02-colour-budget.md` §7) survives this too.

---

## 4 · Which omarchy themes suit Zerado

### 4.1 · The sweep

All 35 files were parsed (both dialects, §1.2) and run through the §2.2 gate with the §2.4
selection rule and the §2.1 pinned method. **Everything in this section is computed here.**

**Result: 21 VALID · 14 INVALID.** Of the valid: **5 brand-true · 8 player · 8 monochrome.**

**Polarity, measured — and the corpus is not what a name suggests.** Five of the 35 are light
themes: **`catppuccin-latte`, `flexoki-light`, `lupine`, `rose-pine`, `white`.** Two corrections
worth recording, because both are easy to get wrong from the filename:

- **`solitude` and `ethereal` are DARK**, not light. `solitude` declares `mode = "dark"` and its
  ground `#101315` measures relative luminance **0.006**; `ethereal`'s `#060B1E` measures **0.004**.
- **`lupine` and `rose-pine` are LIGHT.** `lupine` declares `mode = "light"` (`#FAFAFA`, luminance
  0.956). The bundled `rose-pine` is the **Dawn** variant — ground `#FAF4ED`, luminance **0.911**,
  ink `#575279` — despite the name reading like the dark original.

### 4.2 · The recommended set

Six, plus the authored light default. A curated set is a set someone chose; shipping 21 because 21
passed would be a list, not a choice.

| # | Theme | Tier | Band | ΔE floor | Why it earns a place |
|---|---|---|---|---|---|
| **1** | **`delorean`** | brand-true | comfortable | **15.08** | The brand's governing principle is literally the DeLorean, and the palette earns it independently: a near-white neutral, a pale-yellow warm, a cyan, an orchid — structurally Zerado's own set. **The strongest non-default theme in the corpus.** |
| **2** | **`forest-green`** | brand-true | comfortable | **16.35** | The highest floor of any brand-true theme. A warm neutral that is genuinely warm, and a violet that is genuinely violet |
| **3** | **`kanagawa`** | brand-true | distinct | 10.23 | Ink-on-warm-paper feel; the calm option. Sits near the floor — `abandoned` `#957FB8` clears contrast at **4.67**, close to the line |
| **4** | **`catppuccin`** (mocha) | brand-true | distinct | 11.02 | The most widely recognised palette in the set, and it passes on its own merits |
| **5** | **`spectra`** | brand-true | distinct | 10.26 | Pure-white neutral, saturated warm, clean violet — reads closest to the default of any theme here |
| **6** | **`retro-82`** | player | comfortable | **17.84** | The second theme the founder named. Ships **labelled** — see §4.3 |
| **7** | **the authored light default** | — | — | — | §3.2/§3.4. **No harvested light theme reaches brand-true or player tier** (§4.4), so Zerado's light offering has to be its own |

**Optional eighth — `white`, as the monochrome light option.** VALID · monochrome · distinct
(**10.43**). Four greys, separated by lightness alone, which survives *every* dichromat model
intact. It is in one strict sense the safest theme in the corpus, and it is honest about what it
is. Recommended only if a second light option is wanted; it carries no red at all, so the scanner
and the alarm annunciator render uncoloured plus bold throughout.

### 4.3 · `delorean` and `retro-82` — the two named themes split

The founder asked whether the two obvious names actually pass rather than merely sounding right.
**They do not both pass the same way, and that is the more useful answer.**

**`delorean` — VALID · BRAND-TRUE · COMFORTABLE (ΔE 15.08).** It passes the strictest reading of
every rule in this document with no deviation anywhere.

| Role | Slot | Hex | Contrast on `#0D0221` | `L*` | `C*` | `h` | In band |
|---|---|---|---|---|---|---|---|
| not started | `color7` | `#F9F9F9` | 19.04 | 97.9 | 0.0 | — | ✓ neutral |
| in progress | `color11` | `#FFF5B8` | 18.14 | 96.0 | 31.2 | 100.6° | ✓ warm |
| zerado | `color12` | `#52E5FF` | 13.37 | 84.5 | 39.1 | 217.9° | ✓ cool |
| abandoned | `color5` | `#B967FF` | 6.15 | 59.2 | 86.8 | 313.3° | ✓ violet |
| scanner | `color1` | `#FF2A6D` | — | | | 12.4° | ✓ red band |

Worst pair `not-started × zerado` at 15.08 — comfortable. The `abandoned` `#B967FF` sits **0.5°**
from the brand's own orchid `#C77DFF`. This theme is not merely compatible with Zerado; it is
independently arriving at the same colour argument.

**`retro-82` — VALID · PLAYER · COMFORTABLE (ΔE 17.84), with a named deviation.**

| Role | Slot | Hex | Contrast on `#05182E` | `h` | In band |
|---|---|---|---|---|---|
| not started | `color7` | `#A7C9C6` | 10.04 | — (`C* 12.1` → neutral) | ✓ |
| in progress | `color15` | `#F6DCAC` | 13.39 | 85.7° | ✓ warm |
| zerado | `color5` | `#3F8F8A` | **4.68** | 190.5° | ✓ cool |
| **abandoned** | `color3` | `#E97B3C` | 6.26 | **53.8°** | **✗ DEVIATION — warm, not violet** |

**Its palette contains no violet at all.** `retro-82` is a duotone — orange and teal — and its ANSI
assignment does not follow ANSI hue semantics (its "blue" slot `color4` is `#FAA968`, an orange;
its "magenta" slot `color5` is `#3F8F8A`, a teal). With no violet candidate, `abandoned` is forced
into the warm family beside `in progress`. The two are comfortably distinguishable — 19.40, well
clear — but they are the **same hue family**, so the reading "warm means active" stops holding.
Its `zerado` also clears contrast at **4.68**, the tightest margin in the recommended set.

**Verdict: ship it, labelled `player · abandoned deviates (warm)`.** It is the second theme the
founder named, it passes the gate comfortably, and the deviation is visible to the user rather than
hidden. What it is **not** is brand-true, and it should not be described as such.

**The narrower reading, for the record.** Under a strict rule that additionally refused every
out-of-band role, exactly **5 of 35** themes would remain — the five brand-true ones — and **no
light theme at all**. That is the alternative in §6, presented rather than smuggled.

### 4.4 · Unsuitable, with reasons

**(a) The whole palette collides with the reserved alarm colour — 1 theme.**

| Theme | Reason |
|---|---|
| **`mars`** | Exactly **one** non-red palette entry clears 4.5:1 on its own ground. The palette is red and orange end to end, so nearly every candidate falls in the band reserved for the scanner (G5). **This is the clearest case of an accent colliding with the state palette:** the theme's identity colour *is* Zerado's alarm colour, and there is nothing left to build four states from |

**(b) Light palettes that cannot carry ink on their own ground — 3 themes.**

| Theme | Entries clearing 4.5:1 (of 18) |
|---|---|
| **`catppuccin-latte`** | **2** |
| **`flexoki-light`** | **2** |
| **`rose-pine`** (Dawn) | **2** |

This is a **systematic** result, not three coincidences. Omarchy light palettes are authored so
that *syntax highlighting* reads on a light ground — a job that tolerates 2–4:1 — not so that *body
text* clears AA. `catppuccin-latte`'s yellow `#DF8E1D` measures **2.31:1** on its own background.
**It is exactly the problem brand §4.5 already solved by authoring a darkened paper amber** —
`#965600` at 5.79 on white — instead of reusing `#FFB000`. A harvested light palette has not done
that work, and Zerado cannot do it on the author's behalf without inventing their identity (§1.6).
**This is why §4.2's light recommendation is Zerado's own authored expression.**

**(c) Two states can look alike — 10 themes.** Best achievable separation, computed, against the
required 10.0:

| Theme | Best ΔE | Theme | Best ΔE |
|---|---|---|---|
| `blackturq` | **2.59** | `cobalt2` | **6.63** |
| `everforest` | **4.85** | `vantablack` | **7.78** |
| `gruvbox` | **4.91** | `blackgold` | **8.34** |
| `miasma` | 3 eligible entries only | | |

**And three near misses, named because the floor is a line and not a mood:**

| Theme | Best ΔE | Short by |
|---|---|---|
| **`arc-blueberry`** | **9.92** | 0.08 |
| **`solitude`** | **9.49** | 0.51 |
| **`tokyo-night`** | **9.75** | 0.25 |

These are the three cases where the honest answer is *nearly*. They are refused — a gate that bends
for 0.08 is not a gate — but they are recorded so that if an upstream palette shifts, the answer is
recomputed rather than re-argued. **`nord` deserves its own line:** it fails at **7.65**, and its
violet `#B48EAD` measures **4.41:1** on its own ground — **0.09 short of AA**, which removes the
one candidate that would have made it work.

**The valid-but-not-recommended remainder** — `ethereal`, `purplewave`, `hackerman`, `osaka-jade`,
`aetheria`, `ristretto`, `moodpeak` (player); `harbordark`, `inkypinky`, `last-horizon`, `lumon`,
`matte-black`, `nord`… — passes the gate and may be adopted later. They are omitted from the
curated set for taste, not for safety, and this document does not pretend otherwise.

---

## 5 · Attribution

Every palette in the corpus is **MIT** — upstream `basecamp/omarchy` for the distribution themes,
individual authors for the rest, verified in FlowForge's `ATTRIBUTION.md` at source. Adapting them
is legal **with attribution**, and MIT's condition is explicit: *"The above copyright notice and
this permission notice shall be included in all copies or substantial portions of the Software."*

### 5.1 · Where it lives

`themes/ATTRIBUTION.md` in the Zerado repository, **embedded in the binary** and printed by
`zerado theme credits`. FlowForge's reasoning is adopted verbatim and it is the load-bearing part:
a notice that ships *beside* a binary does not travel with it, and MIT's condition is about copies
of the software. Embedding is how the condition is actually met for a distributed executable.

**A test fails the build if a `.toml` lands in the themes directory without a row in
`ATTRIBUTION.md`.** This is FlowForge's `TestEveryBundledThemeIsAttributed` pattern, adopted as
specification. It is the mechanism that stops the register going stale behind the shipped set — the
only mechanism that reliably does.

**Zerado's repository is public and its openness is a ratified promise** (`decisions.md` Q2), which
raises rather than lowers the bar: the notice is readable by the authors it credits.

### 5.2 · What a contributor must supply to add a theme

Six fields. A pull request missing any one is incomplete, and the first four are the licence
condition rather than a house style.

1. **Upstream URL** — the repository or distribution the palette comes from.
2. **Licence** — and it must permit redistribution inside a distributed binary. *No licence file at
   all means no grant*, however freely the palette was published; that is copyright's default, not
   a judgement about the author.
3. **Copyright holder** — exactly as the upstream `LICENSE` states it.
4. **Harvest provenance** — the source path and the date it was taken.
5. **Verbatim or derived** — and if derived, **by what rule**. FlowForge's `delorean` header is the
   model: it records that the palette was converted from `alacritty.toml` by a rule-for-rule port of
   omarchy's own generator, pinned by a conformance test, so that *"the palette is therefore the one
   omarchy itself would generate; no value is invented."* A derivation that cannot say that much
   about itself is an invention wearing a harvest's clothes.
6. **The validation record** — the §2.7 output, **generated, not typed**.

> **Design decision — field 6 is Zerado's addition and it is not optional.** FlowForge's register
> records provenance; Zerado's must also record *verdict*. The tier is printed to users in the
> picker and the credits, which makes it a shipped claim about accessibility. A shipped claim that
> was hand-written is a shipped claim nobody checked. Generated, it is re-checkable on every build,
> and it re-checks itself when the gate changes.

**Themes that are present on a machine but not bundled remain fully usable.** Not redistributing a
palette is a decision about *redistribution*, not about support — FlowForge draws exactly this line
and Zerado draws it in the same place.

---

## 6 · The alternative this document rejected

Per Rule #2, the one genuinely consequential fork, stated rather than buried. It is a real fork:
both options are defensible and they produce visibly different products.

| | **A · Hue-anchored, no deviation** | **B · Gate is measurement, tier is brand** *(adopted)* |
|---|---|---|
| Rule | A theme is valid only if all four states sit in their brand hue band | A theme is valid if it passes G1–G5; band conformance sets a **tier**, not a verdict |
| Themes valid | **5 of 35** | **21 of 35** |
| Light themes valid | **0** | 2 (`white`, `lupine`, both monochrome) |
| `retro-82` | **refused** | valid, labelled `player` |
| Strength | The product's colour argument is identical everywhere. Perfect cross-theme muscle memory | Honours the founder's direction to bring the themes across; keeps taste out of a measurement; makes every departure visible instead of invisible |
| Cost | Refuses 30 themes including every light one, and contradicts the direction that prompted this work | Some themes render `in progress` green or `abandoned` orange. Cross-theme muscle memory is weaker outside the brand-true tier |

**Why B.** Option A collapses two different questions — *is this safe to read* and *is this
recognisably Zerado* — into one verdict, and the moment those are collapsed, taste starts failing
themes on accessibility grounds it cannot actually claim. B keeps them apart: the **gate** is
measured and refuses only what a user genuinely cannot read; the **tier** is a brand judgement,
stated openly, printed in the credits, and revisable without touching the safety rule. It also
delivers what was asked for — the themes come across — while making the five that carry Zerado's
own argument identifiable at a glance.

**A is still available at zero cost:** ship only the brand-true tier. That is a curation decision
the founder can make later without changing a line of this contract.

---

## 7 · The acceptance bar for theme work

Numbered so each line fails individually, in the style of `00-design-brief.md` §10. This is what
`fft-code-reviewer` judges theme work against.

**The contract**
1. Every semantic token is classified **SUPPLIED / DERIVED / PINNED** (§1.1) and no token is in two
   classes.
2. **Both dialects parse** (§1.2). A semantic-dialect file yields a full sixteen-slot palette; a
   file yielding fewer than four usable colours is a **load error**, never a silent fallback.
3. No PINNED token is readable from a theme file. An attempt is **ignored and reported**, never
   honoured.
4. `--z-primary` resolves to the same value as `--z-state-in-progress`, and `--z-accent` to the same
   value as `--z-state-zerado`, in **every** theme.

**The gate**
5. G1–G5 (§2.2) are implemented exactly, each independently failable, each reporting its own name.
6. The CVD method matches §2.1 **step for step** — matrices, D65 white point, the linear-RGB clamp,
   CIEDE2000 with `kL=kC=kH=1`. A test asserts the dark default measures **11.81** at its tightest
   pair; if that number moves, the method drifted.
7. Selection follows §2.4 and is **deterministic**: the same file yields the same four states on
   every run and on every machine.
8. A failing theme **does not activate**, prints its named reason, and **leaves the previously
   active theme untouched** (§1.6).
9. Derived tokens are repaired by **minimum motion**; supplied tokens are **never** repaired.

**Identity under any theme**
10. Every state renders **colour and glyph and label** — all three, every row, every tier, every
    theme (`00-design-brief.md` §10 line 12).
11. `NO_COLOR` outranks every theme and emits **zero** SGR sequences (§3.5 step 1).
12. At the **16-colour floor** every theme collapses to the default's ratified mapping (§1.4) and no
    region is separated by fill.
13. **No region is separated from another by background fill in any theme**, including one whose
    ground Zerado painted under §3.6.
14. The chrome-cyan budget, the amber allow-list and the red closed list are enforced against the
    **active theme's** earned/ambient/alarm colours, not against the default's hexes.

**Evidence**
15. Every ratio and every ΔE printed anywhere is either read from the brand's measured table **or**
    computed and **marked as computed, with the method named**. No estimate ships.
16. Every bundled theme has an `ATTRIBUTION.md` row with all six §5.2 fields, the build fails
    without it, and field 6 is **generated**.
17. A rendered artifact exists for the default **and** at least one theme from each tier, reviewed
    against this document and **validated by the founder before merge**.

---

## 8 · Open for the founder

1. **The floor is 10.0, not 11.9 (§2.3).** 11.9 came from an unpinned method and a floor set at the
   default's own measurement is a tautology. 10.0 is the manual's own stated scale and the default
   clears it with headroom. **Confirm the floor.**
2. **The light default fails its own gate (§3.2) at ΔE 5.41 against 10.0.** Two pairs fail. The
   cause is a single un-carried hue (§3.3b) and a minimum-motion repair exists (§3.4). **Confirm
   that no light theme ships until brand governance repairs it**, and that `fft-brand-architect`
   owns the work.
3. **`tokens.css` §10 records 7.30:1 for `--z-state-abandoned`; it measures 7.67:1 (§3.3a).**
   Nothing shipped broken. **Confirm the correction lands through brand governance**, in
   `tokens.css` and `tokens.json` and the manual in one commit.
4. **The gate/tier split (§6).** Valid means *safe to read*; brand-true means *recognisably Zerado*.
   21 themes are valid, 5 are brand-true. **Confirm this is the split you want** — or choose Option
   A and ship the five.
5. **`retro-82` is valid but not brand-true (§4.3):** its palette has no violet, so `ABANDONED`
   lands in the warm family beside `IN PROGRESS`. Recommendation is to ship it labelled.
   **Confirm.** `delorean` needs no decision — it passes everything cleanly.
6. **The paint rule (§3.6)** is the only exception ever made to *"Zerado does not paint a
   background"*: paint the declared ground **only** when the user explicitly names a theme.
   **Confirm the exception, and that it is the only one.**
7. **`white` as the monochrome light option (§4.2).** Four greys separated by lightness alone —
   arguably the most CVD-robust theme in the corpus, and consistent with brand §5.4's own position
   that a colourless Zerado is a correct one. **In or out?**

---

## 9 · Routing

| Work | Owner |
|---|---|
| This contract, the gate, the tiers, the acceptance bar | `fft-design-architect` |
| **The light-default repair (§3.2, §3.4), the 7.30→7.67 correction (§3.3a), any new token, every ANSI-256 derivation** | `fft-brand-architect` — through brand governance §10, never a screen PR |
| The theme **picker** — where it lives, navigation, live preview, focus | `fft-tui-architect` |
| The picker's per-screen visual design, including how a tier and a deviation are shown | `fft-tui-designer` |
| The loader, the validator, `zerado theme credits` — `View`/`Update` and the pure functions behind them | `fft-tui` |
| GOLDEN verdict, judged against §7 | `fft-code-reviewer` |

**Terminal composition is not this document's to give.** This decides *what a theme must supply* and
*what makes one valid*; `fft-tui-architect` decides *how the picker is composed*. Two heads, one
product.
