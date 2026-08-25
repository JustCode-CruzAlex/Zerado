---
title: Zerado — Brand Manual
discipline: BRAND IDENTITY
dwg-no: ZRD-BRAND-01
rev: A
date: 2026-08-24
status: in-review
subtitle: An object from 1984, built for 2030
rendered-from: forgeplay-output/landing-page/brand/brand-manual.md
---

# Zerado — Brand Manual

**FlowForgeSoft · zerado.app · Version 1.0 · Revision A**

This manual is the governing document for Zerado's visual identity. Every colour value, type
size and motion curve in it exists as a token in `tokens.css` and `tokens.json` in this same
directory. **Those files are the implementation; this document is the reasoning.** Where the
two disagree, the tokens are what ships — and the disagreement is a bug to file, not a
judgement call to make at the keyboard.

Every contrast ratio printed here was **computed**, not estimated. Every terminal colour index
was **derived** by nearest-neighbour search in CIELAB, not eyeballed. Where a number is tight,
it says so.

---

## 1 · The brand idea

Zerado is a terminal-first game-library tracker. It reads everything you own, sorts it by
**mood** rather than genre, watches the price against the all-time low, and answers one
question: *what should I play tonight?* It is the second product from FlowForgeSoft.

The identity has one job: to look like **an expensive object from 1984 that was designed for
2030.**

### The governing principle — the DeLorean and KITT

Two objects define the register, and they make the same argument. The
**DeLorean DMC-12** (1981–83) is brushed stainless, unpainted, gull-winged — in
a decade of paint and plastic it refused to look like its own decade. The
**Knight Industries Two Thousand** (1982) is a black Trans Am whose interior is
the real statement: a wraparound cockpit of amber readouts, a machine that talks back, and a
single red light sweeping across a dark nose.

Neither is nostalgic. Both are **the eighties' own idea of tomorrow** — and that is the
target. Neon, CRT glow, grid horizon, extruded and chromed type, scanner sweep, amber and red
readouts, the precision of machined metal.

> ### The bar, stated so it can be failed
>
> **Retro-future, never retro-nostalgia.** A faded, dusty, VHS-degraded, sepia,
> scanline-scratched, "remember the eighties?" treatment is a **MISS — not a
> variation, not a stylistic alternative, not a debate.** No distressed
> textures, no fake tape damage, no torn-poster edges, no nostalgia-kitsch.
>
> The eighties in this brand are **new**. The object has just been unboxed. It
> is clean, it is lit, and it is pointed forward.

Zerado is a **text program** — the oldest surviving form of software — with
**native phone apps** coming later: the old form and the new form in one
product. The brand is that argument rendered visually. When a decision is genuinely close, the
tiebreak is *which option makes the product's own argument more clearly?*

### The three things the identity must survive

This brand faces a constraint most brands never do, and it shaped the palette from the first
decision rather than being checked at the end: **a terminal emulator**, where colour is an
index in a 256-entry table that may collapse to 16 or be switched off entirely; **WCAG 2.2
AA**, where saturated neon on black is a minefield rather than a free win; and **colour-vision
deficiency**, which affects roughly 1 in 12 men — a meaningful share of an audience of
developers and sysadmins.

Nothing here was designed and then tested. It was designed *against* those three constraints,
and where a beautiful choice failed one, the choice changed. §4.5 records one such change.

## 2 · The name

*Zerado* is Portuguese for a game **beaten — finished, cleared, 100%'d.** The
word comes from arcade cabinets: the score counter running past its last digit and rolling
back to zero. That origin is **literally an eighties artifact**, and it is the source of the
mark.

The word is also a **status inside the product**. That is deliberate: a tool named after the
act it exists to enable explains itself the first time you use it.

| Object | Form | | Object | Form |
|---|---|---|---|---|
| The product, in prose | `Zerado` | | The status, in prose | `zerado` (italic on first use) |
| The command, binary, domain, code | `zerado` | | The status, in the interface | `ZERADO` |
| The wordmark | `ZERADO` | | | |

The full rationale — why it won, the two runners-up, the three accepted risks and the
collision check — is in **`naming.md`**, which is part of this manual and carries the same
authority.

## 3 · The logo

### 3.1 The idea

The mark is **an arcade counter rolled back to zero, with a scanner light travelling through its
slot** — one glyph carrying both reference objects and the name at once. The **zero** is the name and
the arcade origin of the name. The **slot with a travelling red light** is KITT's scanner, and it is
the brand's one signature motion (§7.1). The **rounded-rectangle skeleton** is a digital counter's
readout cell, not a circle — the difference between a machine's zero and a typographic one.

### 3.2 Construction

Drawn on a **64 × 64** grid. All values are grid units.

| Element | Geometry | Value |
|---|---|---|
| Counter ring | rounded rect, `x14 y9 w36 h46 rx12`, stroked | stroke **10** |
| Ring outer bounds | | **46 × 56** (aspect 0.82) |
| Counter aperture | the hole | **26 × 36** |
| Scanner track | rounded rect, `x22 y28.5 w20 h7 rx3.5` | height **7** |
| Scanner pip | rounded rect, `x30 y28.5 w9 h7 rx3.5` | width **9** |

Two rules make the lockup a system rather than an assembly:

1. **The ring is taller than it is wide (0.82).** A square ring reads as a
   button or a warning lamp. The vertical proportion is what makes it a *zero*.
2. **The mark is the wordmark's `O`.** The final letter of `ZERADO` is drawn on
   the same rounded-rectangle skeleton at the same proportion (0.79). Mark and
   logotype are one shape at two scales. **This rhyme is the lockup's whole
   logic and must not be broken by redrawing either half independently.**

The wordmark is **constructed geometry, not set type** — stroked paths, uniform stroke 9, butt
caps, round joins, cap height 48 — so the logo carries no font dependency and no licensing
question. It is *related* to Orbitron but is not Orbitron, and must never be re-set in a
typeface. Letterspacing is **optically corrected, not mechanical**: the open-flanked `A`
tightens against its neighbours and the round `O` opens slightly. Even mechanical advances are
visibly wrong.

### 3.3 Clear space

**Clear space equals the cap height of the logotype** — 48 grid units, written
here as **`Z`** — on all four sides.

```
     ←Z→                                 ←Z→
   ┌───────────────────────────────────────┐
 ↑ │                                       │
 Z │     ┌───┐                             │
 ↓ │     │ ▭ │   Z E R A D O               │
   │     └───┘                             │
 ↑ │                                       │
 Z │                                       │
 ↓ └───────────────────────────────────────┘
```

Nothing enters this box — not text, not a rule, not another logo, not the edge of a
photograph, not a container border. In a header it is measured from the artwork, not from the
container padding.

### 3.4 Minimum sizes — tested, not assumed

Every size below was **rasterised in a browser at the stated pixel size and inspected**. These
are measured floors, not conventions copied from another manual.

| Asset | Minimum | What was verified |
|---|---|---|
| Full lockup | **120 px** wide | All six letters separate and legible |
| Mark alone (`logo-mark.svg`) | **18 px** | Ring, aperture and red slot all present |
| Mark, absolute floor | **16 px** | Reads as a zero with a lit slot; pip and track merge into one dash |
| Favicon (`favicon.svg`) | **16 px** | Simplified geometry — see below |
| Print | **10 mm** wide (mark) | |

> **Why `favicon.svg` is drawn differently.** At 16 px the track/pip distinction
> is 1.75 px tall and collapses. The favicon drops the two-part scanner for a
> **single solid red bar**, thickens the ring, and adds a dark plate so it holds
> on light and dark browser chrome alike. Rendered side by side at 16 px the
> simplified favicon is visibly stronger. **This is the only sanctioned redraw
> of the mark.**

### 3.5 The terminal mark

The identity has to survive where there is no vector renderer at all.

| Form | Cells | Use |
|---|---|---|
| `[0]` | 3 × 1 | The mark. The terminal minimum. |
| `[0] ZERADO` | 10 × 1 | The lockup, for splash and `--help` headers. |
| `0` | 1 × 1 | Absolute floor — a status bar with no room. |

Drawn with `--z-amber` on the brackets and glyph. Never as multi-row ASCII art: it breaks on
resize, on narrow terminals, and in any log capture.

### 3.6 Colour variants

| File | Use |
|---|---|
| `logo.svg` | **Primary.** Full colour, dark backgrounds. The default. |
| `logo-mark.svg` | Mark alone, transparent, square-safe. Avatars, app icons. |
| `logo-mono-white.svg` | One colour, white. Dark, photographic or printed grounds. |
| `logo-mono-black.svg` | One colour, black. Light grounds, fax, single-colour print. |
| `favicon.svg` | Browser tab. Plated, simplified, 16 px-proof. |

In the mono variants the scanner **track and pip merge into one solid bar** — a barred zero.
That is the correct degrade: readable in a single ink, with no knockouts and no transparency
tricks.

**Full colour requires a dark ground.** Amber-and-red on white is out of
register with the brand and loses the CRT logic entirely; on light backgrounds use
`logo-mono-black.svg`.

### 3.7 Never

Each of these is a real failure mode, not a generic warning.

- **Never re-set the wordmark in a typeface** — including Orbitron. It is drawn
  geometry; setting it in type breaks the `O`-to-mark rhyme (§3.2).
- **Never scale the mark and the wordmark independently.**
- **Never put the full-colour logo on a light or photographic background.**
- **Never add glow, bevel, drop shadow, gradient or 3-D extrusion to the logo.**
  CRT bloom belongs to *interface* elements (§4.7); the logo stays flat. An
  extruded chrome logo is precisely the nostalgia-kitsch §1 rules out.
- **Never recolour the scanner pip.** It is `--z-scanner` and nothing else;
  making it cyan spends the brand's rarest colour (§4.1) on decoration.
- **Never stretch, condense, rotate or skew** any part.
- **Never rebuild the mark as a perfect circle.** It is a counter cell.
- **Never place anything inside the clear space** (§3.3).
- **Never animate the ring.** Only the pip moves, only along its track.

## 4 · The colour system

### 4.1 The idea: one colour is common, one is earned

The palette carries product meaning, not just aesthetics. **Amber is the ambient voice** — the
phosphor readout, the colour the machine speaks in constantly: headings, the mark, anything running.
It is everywhere, and it is meant to be.

**Cyan is earned.** Reserved for **completion** — the `ZERADO` state — and for the single most
important call to action on a screen. It is the colour of the thing the product exists to let you do,
so it is spent rarely and never decoratively. A screen with cyan in five places has thrown away the
only signal the palette has.

**Red is motion and alarm only** — the scanner sweep, destructive confirmations, errors; never body
text, never a decorative accent. **Chrome is structure** — brushed stainless for text and for the
dormant state. Target distribution on any surface: **60 % chassis (the blue-blacks) · 30 % chrome ·
10 % amber**, with cyan measured in single elements, not percentages.

### 4.2 Roles and values

Measured against `--z-surface` (`#0B0D14`) unless noted. `AAA` = 7:1+, `AA` = 4.5:1+ for body
text.

| Role | Token | Hex | Ratio | Grade |
|---|---|---|---|---|
| Page ground | `--z-bg` | `#05060A` | — | — |
| Surface | `--z-surface` | `#0B0D14` | — | — |
| Raised surface | `--z-surface-raised` | `#141A24` | — | — |
| Overlay surface | `--z-surface-overlay` | `#1D2532` | — | — |
| Primary text | `--z-text` | `#E9EEF5` | **16.65** | AAA |
| Secondary text | `--z-text-secondary` | `#A9B5C7` | **9.36** | AAA |
| Tertiary text | `--z-text-tertiary` | `#8492A8` | **6.15** | AA |
| Primary (amber) | `--z-primary` | `#FFB000` | **10.59** | AAA |
| Primary hover | `--z-primary-hover` | `#FFC94D` | **12.68** | AAA |
| On-primary text | `--z-primary-contrast` | `#05060A` | **11.05** on amber | AAA |
| Accent (cyan) | `--z-accent` | `#19E0FF` | **12.15** | AAA |
| On-accent text | `--z-accent-contrast` | `#05060A` | **12.68** on cyan | AAA |
| Link | `--z-text-link` | `#19E0FF` | **12.15** | AAA |
| Hairline | `--z-border` | `#2A3342` | 1.53 | decorative only |
| UI boundary | `--z-border-strong` | `#64748B` | **4.08** | passes 1.4.11 |
| Focus ring | `--z-focus-ring` | `#19E0FF` | **12.15** | AAA |
| Scanner | `--z-scanner` | `#FF2E2E` | 5.25 | motion/alarm, not text |
| Error text | `--z-scanner-300` | `#FF6B6B` | **6.99** | AA |

**The three pairs the brand explicitly guarantees:**

| Purpose | Pair | Ratio | Requirement |
|---|---|---|---|
| **Body text** | `#E9EEF5` on `#0B0D14` | **16.65:1** | 4.5:1 → AAA |
| **Large display text** | `#FFB000` on `#05060A` | **11.05:1** | 3:1 → AAA |
| **Interactive / focus** | `#19E0FF` on `#0B0D14` | **12.15:1** | 3:1 → AAA |

**The worst text pair in the entire system** is `--z-text-tertiary` on
`--z-surface-overlay` at **4.89:1** — still AA. There is no failing text pair in the palette,
by construction: where one appeared, the colour was changed rather than the pairing permitted.

**Two rules that are not negotiable.** `--z-border` is **decorative** — at
1.53:1 it may draw a hairline between blocks of content, but never mark the edge of an input,
a button, or any control whose boundary carries meaning; those use `--z-border-strong`, which
measures ≥ 3.24:1 against every surface and satisfies WCAG 1.4.11. And **the focus ring is
never removed** — 2 px with 2 px offset, on every interactive element, in every state. A
keyboard-first audience is the launch audience; removing focus rings for aesthetics would be a
self-inflicted wound on exactly the people the product is for.

### 4.3 The four game states

The **most-used visual system in the product** — every row of every list.

**The co-render rule: colour AND glyph AND label.** Every state carries all
three; remove any one and the state is still unambiguous. This is not a concession to
accessibility, it is what makes the system survive `NO_COLOR`, a monochrome terminal, a
screenshot in a bug report, and colour-blind users, all by one mechanism.

| State | Colour | Hex | Glyph | ASCII | Label | Ratio |
|---|---|---|---|---|---|---|
| Not started | chrome | `#A5A29B` | `○` U+25CB | `[ ]` | `NOT STARTED` | 7.62 |
| In progress | amber | `#FFB000` | `◐` U+25D0 | `[~]` | `IN PROGRESS` | 10.59 |
| **Zerado** | cyan | `#19E0FF` | `◉` U+25C9 | `[*]` | `ZERADO` | 12.15 |
| Abandoned | orchid | `#C77DFF` | `⊘` U+2298 | `[x]` | `ABANDONED` | 7.21 |

The glyphs are a **single visual progression** — an empty ring, a ring half filled, a ring
with a solid core, a ring struck through — so the sequence reads as a story even in one ink.
`Zerado` takes the cyan because cyan is the earned colour and because a filled core is the
strongest glyph in the set: the state the product is named after is the one that looks like an
achievement.

### 4.4 Colour-vision deficiency — how this was verified

**Method.** Each state colour was simulated for protanopia and deuteranopia with
the **Viénot, Brettel & Mollon (1999)** dichromat model in linear RGB, then all six pairs were
measured for separation as **CIEDE2000** in CIELAB. Reference targets are the Bang Wong and
Paul Tol colour-blind-safe palettes.

| Vision | Worst pair | ΔE2000 |
|---|---|---|
| Normal | not-started × zerado | 28.2 |
| Protanopia | not-started × zerado | 23.7 |
| Deuteranopia | **zerado × abandoned** | **11.9** |

ΔE ≥ 10 is distinct; ≥ 15 is comfortable. **Every pair clears 10 in every model.** The
tightest — cyan against orchid under deuteranopia at 11.9 — is the one place where glyph and
label genuinely carry load rather than merely reinforce. That is what they are there for.

> **The change this testing forced.** The first draft used a **blue-cast** steel
> grey (`#9FB0C6`) for *not started* — the obvious choice beside a blue-black
> chassis. Simulated, it collapsed against the cyan at **ΔE 8.8 under
> deuteranopia**: two states that would look nearly identical to roughly 1 in 12
> male users, on the most-used component in the product.
>
> It was replaced with a **warm** brushed grey, `#A5A29B`, measuring **25.8** —
> a 2.9× improvement. The warmth is load-bearing engineering, not taste, and
> must not be "corrected" back toward blue. It is also more accurate: brushed
> stainless *is* warm-neutral. The constraint and the reference agreed.

### 4.5 The paper expression

Zerado is dark-native, but the manual prints, the PDF renders on a white page, and some
contexts are simply light. Rather than let those fall back to a generic palette, the brand has
a **second designed expression**.

| Role | Screen | Paper | Ratio on paper |
|---|---|---|---|
| Ground | `#05060A` | `#F3F5F8` | — |
| Ink | `#E9EEF5` | `#0B0D14` | **17.77** AAA |
| Ink soft | `#A9B5C7` | `#39434F` | **9.20** AAA |
| Ink muted | `#8492A8` | `#5E6A7A` | **5.04** AA |
| Amber | `#FFB000` | `#965600` | **5.30** AA |
| Amber strong | — | `#7A4600` | **7.11** AAA |
| Cyan | `#19E0FF` | `#0A6070` | **6.58** AA |
| Scanner | `#FF2E2E` | `#C61F1F` | **5.32** AA |

These are not different colours; they are **the same hues carried to ink weight** — measured,
the paper amber sits **9.7°** from the screen amber in CIELAB hue and the paper cyan **3.1°**
from the screen cyan. It still reads as Zerado because it *is* Zerado.

The paper values are exposed in `tokens.css` §9 under the unprefixed names the PDF renderer
resolves (`--ink`, `--canvas`, `--accent`, `--font-display` and siblings). That mapping is
what makes a Zerado document print **in Zerado's brand**, and it is scoped so it cannot alter
the dark default.

### 4.6 Glow

CRT bloom is available and used **sparingly** — a focused input, a hovered primary action, the
scanner track: `--z-glow-amber` / `--z-glow-cyan` / `--z-glow-red`, a tight 1 px colour ring
plus a soft short-radius bloom. It is not a drop shadow, it is never stacked, it never appears
on more than one element on screen at a time, and per §3.7 it never touches the logo.

## 5 · Colour in the terminal

Zerado is a TUI first. A palette that only exists as hex is only half a palette.

### 5.1 Three representations

Every brand colour resolves three ways. The ANSI-256 indices were **derived** by
nearest-neighbour search against the xterm 256-colour cube, measured as CIEDE2000 in CIELAB —
not chosen by eye.

| Token | Hex | ANSI-256 | ΔE | 16-colour | Survives 16? |
|---|---|---|---|---|---|
| Surface | `#0B0D14` | **232** | 3.8 | `black` | yes |
| Raised surface | `#141A24` | **234** | 6.7 | `black` | merges with surface |
| Text | `#E9EEF5` | **255** | 3.7 | `bright white` | yes |
| Secondary text | `#A9B5C7` | **249** | 8.5 | `white` | yes |
| UI boundary | `#64748B` | **67** | 9.3 | `bright black` | yes |
| **Amber** | `#FFB000` | **214** | **0.4** | `bright yellow` | yes |
| **Cyan** | `#19E0FF` | **45** | 3.4 | `bright cyan` | yes |
| **Chrome** | `#A5A29B` | **247** | 3.9 | `white` | yes |
| **Orchid** | `#C77DFF` | **177** | 3.8 | `bright magenta` | yes |
| Scanner | `#FF2E2E` | **9** | 5.3 | `bright red` | yes |

The amber match is **near-exact at ΔE 0.4** — xterm 214 is `#FFAF00` against the brand's
`#FFB000`. The brand's primary colour is already a terminal colour.

### 5.2 The 16-colour floor

At 16 colours the four states resolve to **four distinct slots**: not started → `white` (7),
in progress → `bright yellow` (11), zerado → `bright cyan` (14), abandoned → `bright magenta`
(13). No collisions — the state system is fully functional at the lowest depth a terminal
offers.

**What is lost**, stated honestly: the surface ramp flattens, since
`--z-surface` and `--z-surface-raised` both become `black`. Elevation must therefore be
carried by borders and spacing, never by fill. Design terminal layouts so they never depend on
surface fill to separate regions.

### 5.3 Backgrounds the palette must hold

Users bring their own theme. Worst case across the four state colours:

| Background | Worst state ratio |
|---|---|
| Pure black `#000000` | 7.81 |
| xterm dark grey `#1C1C1C` | 6.33 |
| Solarized Dark `#002B36` | 5.58 |
| Gruvbox Dark `#282828` | 5.48 |
| One Dark `#282C34` | 5.20 |

**Every state colour clears AA body text (4.5:1) on every one of them.** The
palette does not require the user to adopt the brand's background.

### 5.4 `NO_COLOR`

Zerado honours **`NO_COLOR`** (no-color.org). When set, the product emits **no SGR sequences
at all** — not a reduced palette, not a "safe" subset. Zero.

**What the identity degrades to: everything.** The four states remain fully
distinct through **glyph and label**, which is the entire reason the co-render rule exists:

```
  ○  NOT STARTED   Hollow Knight              41h
  ◐  IN PROGRESS   Outer Wilds                12h
  ◉  ZERADO        Return of the Obra Dinn     9h
  ⊘  ABANDONED     Sekiro                      3h
```

Structure is then carried by box-drawing characters, spacing, and the hierarchy of a
monospaced grid. A Zerado screen with colour stripped is not a degraded Zerado screen — it is
a correct one. The same applies to `forced-colors: active` on the web, by the same mechanism.
Where a Unicode glyph cannot be relied on, the ASCII column of §4.3 (`[ ]` `[~]` `[*]` `[x]`)
is the fallback.

## 6 · Typography

### 6.1 The three families

Three faces, three jobs, and a hard rule about which does what. **All three are SIL Open Font
License 1.1** — free for commercial use, self-hostable, legally shippable. The site self-hosts
all three; there is no external CDN.

| Role | Family | Licence | Source | Axes |
|---|---|---|---|---|
| **Display** | **Orbitron** | SIL OFL 1.1 (RFN "Orbitron"), Matt McInerney | `github.com/googlefonts/orbitron-vf` | `wght` 400–900 |
| **Voice** | **Space Grotesk** | SIL OFL 1.1, Florian Karsten | `github.com/floriankarsten/space-grotesk` | `wght` 300–700 |
| **Native tongue** | **JetBrains Mono** | SIL OFL 1.1, JetBrains | `github.com/JetBrains/JetBrainsMono` | `wght` 100–800 |

**Orbitron carries the era** — a square geometric face drawn explicitly as
display type "for the future", the closest freely-licensable equivalent to the wide,
high-waisted geometric capitals of eighties corporate-future lettering. It is why a headline
reads as 1984 rather than 2016.

> **Orbitron is display-only, and this is a hard rule.** Never below 23 px,
> never for prose, never more than about eight words, always uppercase, always
> tracked out. Orbitron as body copy is unreadable and instantly reads as a
> games-console UI — the wrong eighties.

**Space Grotesk is the voice** — a proportional face derived from Space Mono, so
technical by heritage and able to pair with a monospace without looking accidental. It carries
every sentence in the product and on the site. It was chosen over a fully neutral face (Inter)
because neutrality would have made Orbitron look like a costume; Space Grotesk shares enough
character to make the pairing read as one voice at two volumes.

**JetBrains Mono is the native tongue.** The product is a terminal application,
so monospace is not a code-block afterthought — it is **a first-class text role** and the
sound of the product itself. It sets every command, path, hours-played figure, price and state
label, and above all the **readout style**: the small tracked-out uppercase label that heads
sections and mimics a cockpit annotation. That role is what carries the KITT reference into
typography.

### 6.2 The ramp

Root 16 px. All values are tokens in `tokens.css` §5.

| Step | Size | Line height | Family | Weight | Tracking | Case |
|---|---|---|---|---|---|---|
| Display | 76 px | 1.02 | Orbitron | 700 | 0.04em | UPPER |
| H1 | 52 px | 1.08 | Orbitron | 700 | 0.03em | UPPER |
| H2 | 34 px | 1.22 | Orbitron | 600 | 0.02em | UPPER |
| H3 | 23 px | 1.22 | Space Grotesk | 700 | 0 | Sentence |
| Body large | 19 px | 1.60 | Space Grotesk | 400 | 0 | Sentence |
| **Body** | **17 px** | **1.65** | Space Grotesk | 400 | 0 | Sentence |
| Body small | 15 px | 1.60 | Space Grotesk | 400 | 0 | Sentence |
| Caption | 13 px | 1.50 | Space Grotesk | 400 | 0 | Sentence |
| **Readout** | 13 px | 1.20 | JetBrains Mono | 500 | **0.18em** | UPPER |

Display and H1 are fluid — `clamp()` tokens hold them between 44 px and 76 px so a long
headline never overflows a phone. **The governing relationship:** as size increases,
line-height tightens and tracking opens. Wide tracking on large capitals is the
eighties-future signature; tight tracking on body text would simply hurt.

### 6.3 Rules

- **17 px is the floor for prose** — never smaller, anywhere, for anything a
  visitor is expected to read as a sentence.
- **Body line-height is 1.65**, above the WCAG 1.4.12 minimum of 1.5.
- **Maximum three weights per family** on any one surface.
- **Measure caps at 68 characters** (`--z-measure`).
- **Left-aligned always.** Never justify — it opens rivers and fights the
  monospaced grid the product's own screens are built on.
- **Never set body text in Orbitron** (§6.1), and **never set all-caps
  paragraphs** — capitals are for display, readout labels, state chips and
  buttons only.

## 7 · Motion

### 7.1 The scanner sweep — the signature

One motion belongs to this brand: **a single light travelling across a dark track**, taken
directly from KITT. It appears on the hero, as a section divider, as an indeterminate progress
indicator, and inside the mark at display size.

| Property | Value | Token |
|---|---|---|
| Track height | 2 px | `--z-scanner-track-h` |
| Track colour | `#5C1414` (`#6E1818` in the mark) | `--z-scanner-track` |
| Pip width | 18 % of track | `--z-scanner-pip-width` |
| Pip colour | `#FF2E2E` | `--z-scanner` |
| Travel | left edge → right edge − pip width | |
| Full cycle | **2400 ms** (1200 ms each direction) | `--z-duration-scanner` |
| Easing | `cubic-bezier(0.45, 0, 0.55, 1)` | `--z-ease-scanner` |
| Direction | `alternate`, `infinite` | |

> **Why that easing.** KITT's scanner was a **physically oscillating mirror**,
> and a real oscillator is slowest at the extremes and fastest through the
> middle — a sine wave, of which `cubic-bezier(0.45, 0, 0.55, 1)` is the closest
> cubic approximation. Linear travel reads as a loading bar; a bounce or
> overshoot reads as a cartoon. The sinusoid is the only one that reads as
> *hardware*.

```css
.z-scanner {
  position: relative; overflow: hidden;
  block-size: var(--z-scanner-track-h);
  background: var(--z-scanner-track);
}
.z-scanner::after {
  content: ""; position: absolute; inset-block: 0; left: 0;
  inline-size: var(--z-scanner-pip-width);
  background: var(--z-scanner);
  box-shadow: var(--z-glow-red);
  animation: z-sweep var(--z-duration-scanner) var(--z-ease-scanner)
             infinite alternate;
}
@keyframes z-sweep {
  from { left: 0; }
  to   { left: calc(100% - var(--z-scanner-pip-width)); }
}
```

**In the terminal.** The same primitive, one row tall: a rule of `─` (U+2500) in
`--z-scanner-track` with **three consecutive cells** of `━` (U+2501) in `--z-scanner` as the
pip, positioned from elapsed time on the same 2400 ms sinusoid and redrawn at 30 fps. Used
only for genuinely indeterminate waits — never as ambient decoration, which would burn a
redraw budget for nothing.

### 7.2 Interface motion

| Purpose | Duration | Easing |
|---|---|---|
| State change (hover, focus) | 90–140 ms | `--z-ease-standard` |
| Entrances, expansions | 220 ms | `--z-ease-out` |
| Exits, dismissals | 140 ms | `--z-ease-in` |
| Larger transitions | 380 ms | `--z-ease-out` |

Motion here is **mechanical, not playful**: no bounce, no elastic, no spring overshoot
anywhere in the interface. This is a machined object.

### 7.3 Reduced motion

Under `prefers-reduced-motion: reduce` all durations collapse to 1 ms and **the scanner pip
parks at the centre of its track at full opacity.** It is deliberately **not hidden** — the
lit slot is an identity element, the travel is the decoration. Users who ask for less motion
get the brand without the movement, not a blank space where the brand was.

## 8 · Voice and tone

Zerado talks to people who live in terminals. The voice is **dry, confident, concrete**: it
states facts, respects the reader's time, and never performs enthusiasm. Never bro-ish, never
breathless.

| Attribute | We are | We are not |
|---|---|---|
| **Dry** | Understated; the fact is the point | Sarcastic, or bubbly |
| **Concrete** | Numbers, names, specifics | Vague, or jargon-padded |
| **Confident** | Direct; no hedging, no apologising | Boastful, or timid |

**1 · After a sync completes**

- ✗ "🎉 Awesome! You crushed it! 247 games synced!!"
- ✓ "247 games. 6 finished. Last played: 3 weeks ago."
- *The second sentence is the product's actual argument; the first buries it under celebration of a database write.*

**2 · A price is in budget but above its historical low**

- ✗ "🔥 This game is a STEAL at $15! Grab it before it's gone!"
- ✓ "$15, and in budget. It's been $8 — last June. Maybe wait."
- *Zerado's value is telling you **not** to buy. Fake urgency destroys the trust the feature is built on.*

**3 · A Steam sync returns nothing**

- ✗ "Oops! Something went wrong 😅 Please try again later."
- ✓ "Steam returned an empty library. Game details are private on your profile — Steam won't share the list until that's public. Settings → Privacy."
- *Name what happened, why, and the next action. "Something went wrong" is the one sentence a terminal user cannot act on.*

**4 · Marking a game zerado**

- ✗ "Congratulations on this incredible achievement! You're a true champion! 🏆"
- ✓ "Zerado. 41 hours. Sixth this year."
- *The moment carries itself. Three facts land harder than an exclamation mark, and the count is the reward.*

**5 · Introducing the community layer**

- ✗ "We're SO excited to announce our amazing new social features — coming soon to revolutionise how you share your gaming journey!"
- ✓ "Comments, reviews and public profiles. These run on servers that cost money, so this part will need a premium account or a donation. Not built yet — Phase 4."
- *State the cost before anyone asks. Up front reads as respect; later reads as a catch.*

**6 · Explaining why it's a terminal app**

- ✗ "Experience the power of a next-generation, blazing-fast CLI-first gaming ecosystem."
- ✓ "It's a text program. It starts instantly, works offline, and your library is one file you own."
- *Three verifiable claims beat six adjectives. This audience checks.*

### Standing rules

- **No exclamation marks.** No situation in this product needs one.
- **No emoji in product copy.** Glyphs carry state (§4.3); emoji carry nothing.
- **Never call the user a "gamer".**
- **Say the number.** "41 hours", not "a lot of hours".
- **Portuguese words stay Portuguese.** *zerado* and *sinopse* are explained
  once, in a clause, on first appearance — then used normally, never re-glossed.
  A glossary would make the page feel like homework.
- **Never claim what isn't built.** Unshipped phases are marked as phases.

## 9 · One identity, three render targets

The real design problem: this brand must be the same brand in three places with almost nothing
technically in common.

| | **Web** | **Terminal** | **Phone (Phase 4)** |
|---|---|---|---|
| Colour depth | 24-bit | 256 → 16 → none | 24-bit |
| Layout unit | px, fluid | character cells | dp / pt |
| Type control | full | one monospaced grid | full, platform metrics |
| Motion | CSS, 60 fps | redraw budget | platform animation |
| Logo | `logo.svg` | `[0]` | `logo-mark.svg` |

**The invariants — identical across all three:** the four states with colour +
glyph + label; **amber common, cyan earned**; the scanner as the only signature motion; the
voice; and dark by default.

**What is allowed to differ.** **Web** is the most expressive surface — grid
horizon, CRT bloom, fluid display type, the full lockup. This is where the brand performs.
**Terminal** trades expression for density and honesty: no glow, no gradients, no fluid type;
hierarchy comes from the monospaced grid, box drawing, spacing and the four-state palette, and
it must stay correct at 16 colours and at zero. **The terminal is the product's home — when
web and terminal conflict, the terminal wins.** **Phone** follows its platform's conventions —
native navigation, gestures and type metrics — while keeping the palette, the state system and
the voice. A phone app that fights iOS or Android to look like a terminal is a worse product
and a worse advertisement for the brand. **Wear the platform; keep the identity.**

## 10 · Governance

**The tokens are the source of truth.** `tokens.css` and `tokens.json` implement
everything in this manual. No surface may hard-code a hex value, a font stack, a spacing step
or a duration. If a value is needed and no token exists, add a token — do not type a hex.

**Naming has three layers, and they must not be collapsed:**

```
primitive   --z-amber-500        raw value; never referenced by a component
semantic    --z-primary          meaning; what components reference
component   --z-btn-bg-primary   local override; defined at the component
```

**To change something in this brand:**

1. Change the token in `tokens.css` **and** `tokens.json` — same commit.
2. Update the affected section of this manual — same commit. A manual that
   disagrees with the tokens is worse than no manual.
3. If it touches a colour used for text, **recompute the contrast ratio** and
   update the table. Do not estimate it.
4. If it touches a **state** colour, **re-run the CVD simulation** for all six
   pairs. The 11.9 minimum in §4.4 is the floor to protect.
5. If it touches a colour used in the terminal, **re-derive** the ANSI-256 index
   rather than adjusting the old one by hand.

**Revision A · 2026-08-24 · FlowForgeSoft**
