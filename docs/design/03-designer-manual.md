---
title: Zerado — Designer Manual
discipline: DESIGN SYSTEM
doc-no: ZRD-DESIGN-04
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: product-manual
---

# Zerado — Designer Manual

**How we work here.** The first document anyone designing for Zerado reads. It says who owns
what, what a screen spec must contain to be buildable, what "done" means, how a design gets
reviewed, which decisions are closed, and how to reopen one if you genuinely must.

---

## 1 · Read these, in this order, before you touch a screen

| # | Document | Why |
|---|---|---|
| 1 | `brand/brand-manual.md` (rev A) — **all of it** | The identity. Not a sketch: the palette in hex **and** ANSI-256, the 16-colour floor, `NO_COLOR`, the measured contrast table, the four states, the scanner, the voice. |
| 2 | `brand/naming.md` | The casing convention. Getting it wrong is the fastest way to make the product read as sloppy. |
| 3 | `brand/tokens.css` · `tokens.json` | **The tokens are the source of truth.** The manual is the reasoning. |
| 4 | `ratification/decisions.md` | The ratified public promises. Binding. |
| 5 | `content/landing-copy.md` | Published, public, binding. Your screens may not contradict a line of it. |
| 6 | `00-design-brief.md` | Which guidelines govern, and the acceptance bar. |
| 7 | `01-design-system.md` | The components you assemble from. |
| 8 | `02-colour-budget.md` | The colour rule you will be failed against. |
| 9 | FlowForge TUI Design Manual (#2371) + Spacing Canon (#2435) | The terminal craft bar and every spacing number. |
| 10 | `PLAY-NOTES.md` | The DeLorean-and-KITT argument in the founder's own words. Read it last, and let it decide the close calls. |

**The bar, so it can be failed:** **retro-*future*, never retro-nostalgia.** Neon, CRT glow,
grid horizon, chrome type, scanner sweep, amber-on-black readout — the eighties' own idea of
**tomorrow**. A faded, dusty, VHS-degraded, "remember the eighties?" treatment is a **MISS, not
a variation.** The object has just been unboxed. It is clean, it is lit, it is pointed forward.

---

## 2 · Who owns which surface

| Surface / decision | Owner | Not theirs |
|---|---|---|
| **Guideline authority** — which corpora govern, the design brief, the acceptance bar | `fft-design-architect` | Does not draw screens |
| **Brand identity** — the manual, the 200 tokens, contrast, CVD, **the ANSI-256 derivations** | `fft-brand-architect` | The only role that may change a token |
| **Terminal composition** — layout tree, pane budget, focus state machine, navigation model, Charm adoption, responsive tier map, the visual QA gate | `fft-tui-architect` | The head of the terminal lane |
| **Per-screen visual design** — the mockup, hierarchy, co-render application, spacing tokens, state tables, copy | `fft-tui-designer` | Does not decide composition or guidelines |
| **`View` / `Update` line code** | `fft-tui` | Does not redesign in code |
| **The GOLDEN verdict** | `fft-code-reviewer` | Judges against `00-design-brief.md`, never a private bar |
| **Phase 4 phone design** | brief by `fft-design-architect` → design by `fft-designer` → build by `fft-flutter` | See `04-terminal-to-phone-bridge.md` |
| **Data model, seams, state machine, ERD** | `fft-tui-architect` (the spine) | Not a design deliverable |

**The chain, stated once:**
`fft-design-architect` (guidelines + acceptance bar) → `fft-tui-architect` (composition) →
`fft-tui-designer` (the screen) → `fft-tui` (the code) → `fft-code-reviewer` (GOLDEN).

**Terminal composition is not the design architect's to give.** The design architect decides
*which guidelines govern* and *what GOLDEN means*; the TUI architect decides *how the screen is
composed*. Two heads, one product. Do not collapse them, and do not route around either.

**When Phase 4 design work actually starts,** a dedicated `fft-flutter-designer` leaf should be
minted rather than assumed. Until then `fft-designer` carries the mobile design load.

---

## 3 · What a screen spec must contain to be buildable

This is the contract `fft-tui-designer` fills for **every** Phase 1 screen. A spec missing any
section is not ready for dispatch — it will generate a follow-up question, and a follow-up
question means the screen gets built twice.

### The required sections

1. **Identity** — screen number, name, phase, and its route in the navigation model.
2. **Purpose** — one sentence: what the player gets from this screen. If it takes two sentences,
   the screen is doing two jobs.
3. **Mockup at 80 columns** — always. 80 is the design width; no screen may *require* more.
4. **Mockup at the screen's primary breakpoint**, if that is not 80.
5. **Visual hierarchy** — what the eye reaches first, second, third, and by which channel (case,
   weight, colour role, box drawing, spacing).
6. **Every applied spacing token, by name**, with its value at the rendered tier —
   `OuterMarginX`, `OuterMarginY`, `InnerPaddingX`, `InnerPaddingY`, `InterElementGap`,
   `HeaderBandHeight` / `HeaderBand(tier, hasSubtitle)`, `leftInset`. **No magic numbers.**
7. **Colour, glyph and label for every state shown**, each resolved to token + hex + ANSI-256 +
   16-colour fallback. If a token's ANSI index is underived, say so and give the interim.
8. **The full state table.** Loading · empty · partial · error · offline · **and the one nobody
   writes down: first run, before anything is synced.** A screen without its first-run row is
   incomplete.
9. **The key map** — every key active on this screen and what it does, including `esc` and `?`.
10. **The exact copy**, in the ratified voice, ready to paste. Not a description of the copy.
11. **40-column behaviour**, and the screen's documented refusal floor below which it declines
    to render and says so in words.
12. **`NO_COLOR` rendering** — the screen with zero SGR sequences, shown, not asserted.
13. **The focus model on this screen** — what can hold focus, the traversal order, and **how
    `esc` behaves**. Every screen has a way out.
14. **The colour budget declaration** — where the one chrome cyan is spent (or that it is
    spent nowhere), and which amber allow-list entries are used.
15. **Reuse verdict per element** — which Charm primitive, or "built fresh" with the reason.
16. **Screen-specific acceptance criteria** — anything beyond `00-design-brief.md` §10 that this
    screen must satisfy.

### The Phase 1 screen set

Confirmed or amended by `fft-tui-architect` in the spine; the amendment is part of what gets
ratified.

First run and setup (Steam key and Steam ID, with the public-profile requirement explained
*where the player hits it*) · sync in progress · **the library deck** · game detail · mark or
change status · search and filter · settings · help and key map · error and offline states.

> **The honest screens are not optional.** First run before any data exists. A sync that fails
> halfway. A private Steam profile that returns nothing. Four hundred games on a forty-column
> terminal. These are the screens that decide whether the product feels solid, and they get the
> same craft as the headline screen. R-5: **every screen must be a hit.**

---

## 4 · The acceptance bar and how a design gets reviewed

### 4.1 · The bar

`00-design-brief.md` §10 — twenty-eight numbered, falsifiable lines. Plus
`02-colour-budget.md` §10, nineteen literal pass/fail boxes. **A reviewer never invents a private
bar; those two lists are the yardstick, and they are the same lists the designer built against.**

### 4.2 · The Screen-Quality Gate

Adopted whole from #2371 §4b. Four steps, in order, no skipping:

1. **Render headless** (`freeze` / `vhs`) at the target viewports — and for any tabular screen
   **at an overflowing row count** (≥ 400 games), not only at the frozen golden size.
2. **Specialist review of the rendered artifact** against this manual, the design brief and the
   colour budget: contrast, co-render, spacing, no clipping, and for a ledger the explicit
   R-10 triad — populated title column, selection stays visible on scroll, summary stays on
   screen.
3. **Founder validation of the screenshot, before merge.** The founder validates the rendered
   artifact — not a description of it.
4. **`fft-code-reviewer` marks GOLDEN only with the founder-validated screenshot.**
   No screenshot → not GOLDEN → no merge.

> **Why a text golden is never enough.** A golden file asserts one render at one size. It never
> exercises "scroll to row 380 and check the selection is still visible", or "is the summary
> still on screen at 400 rows", or "does this survive at 16 colours". Those are only provable
> from a live repro at an overflowing row count and from the rendered pixels.

### 4.3 · Every screen carries its canon

The FlowForge TUI Design Manual and the Spacing Canon travel with every design dispatch — the
same rule a brand manual lives by: *every time you use the logo anywhere, the manual travels
with it.* For Zerado, this designer manual and `00-design-brief.md` travel too. **Nobody
re-derives the canon per ticket.**

---

## 5 · Settled decisions — do not reopen

These are closed. Not "closed for now" — closed. Reopening one costs a founder gate and the full
governance procedure in §6. Designing around one quietly is a defect.

### 5.1 · The four game states

The states, their colours, glyphs, ASCII fallbacks and labels are ratified and **CVD-verified**:

| State | Colour | Hex | ANSI-256 | 16-col | Glyph | ASCII | Label |
|---|---|---|---|---|---|---|---|
| Not started | chrome | `#A5A29B` | 247 | `white` | `○` | `[ ]` | `NOT STARTED` |
| In progress | amber | `#FFB000` | 214 | `bright yellow` | `◐` | `[~]` | `IN PROGRESS` |
| Zerado | cyan | `#19E0FF` | 45 | `bright cyan` | `◉` | `[*]` | `ZERADO` |
| Abandoned | orchid | `#C77DFF` | 177 | `bright magenta` | `⊘` | `[x]` | `ABANDONED` |

- **The co-render rule** — colour **and** glyph **and** label, all three, every one.
- **The warm grey `#A5A29B`** is load-bearing engineering, not taste. The blue-cast `#9FB0C6`
  collapsed against the cyan at **ΔE 8.8** under deuteranopia; the warm grey measures **25.8**.
  **Never "correct" it back toward blue.**
- **The floor to protect is ΔE 11.9** — zerado × abandoned under deuteranopia. The one place
  glyph and label genuinely carry load rather than merely reinforce.

### 5.2 · The palette and the contrast table

Every ratio in brand §4.2 was **computed**, not estimated. `--z-border` is **decorative only**
(1.53); controls use `--z-border-strong` (4.08). **The focus ring is never removed.** There is no
failing text pair in the system, by construction.

### 5.3 · The terminal colour representations

The ANSI-256 indices in brand §5.1 were **derived** by nearest-neighbour search in CIELAB.
The 16-colour floor resolves the four states to four distinct slots with no collisions. **What
is lost is stated honestly:** surface and raised-surface both become `black`, so **elevation is
carried by borders and spacing, never by fill.**

### 5.4 · `NO_COLOR`

When set, Zerado emits **no SGR sequences at all** — not a reduced palette, not a "safe" subset.
Zero. A Zerado screen with colour stripped is not a degraded screen; it is a correct one.

### 5.5 · The type system

Three families, three jobs, all SIL OFL 1.1. **Orbitron is display-only** — never below 23 px,
never for prose, never more than about eight words. 17 px is the floor for prose. Body
line-height 1.65. Left-aligned always. In the terminal the *roles* survive and the *families*
cannot — see `01-design-system.md` §1.5.

### 5.6 · The casing convention

`Zerado` the product · `zerado` the command, binary, domain and code · *zerado* the status in
prose · **`ZERADO`** the state chip in the interface. **Never camel-cased. Never `Zerado.app`.**

### 5.7 · The scanner

One motion belongs to this brand. **2400 ms** full cycle, `cubic-bezier(0.45, 0, 0.55, 1)`,
alternate, infinite; in the terminal a `─` track with a **three-cell** `━` pip at **30 fps**.
Reduced motion **parks the pip**, never hides it. Used **only** for genuinely indeterminate
waits.

### 5.8 · The voice

Dry, confident, concrete. **No exclamation marks. No emoji in product copy. Never call the user
a "gamer". Say the number.** Portuguese words stay Portuguese and are explained once, in a
clause, then used normally. **Never claim what isn't built.**

### 5.9 · The ratified public promises

Local-first · one SQLite file the player owns · no account to use it · works with the network
off · the only network traffic is the services the player connected, using their own keys · no
Zerado-run server before Phase 4 · the affiliate disclosure · **no named community source** ·
English page with Portuguese words kept. The four-phase roadmap and the FAQ in
`content/landing-copy.md` are **published and binding** — no screen may contradict them, and
nothing unbuilt may be presented as working.

### 5.10 · The stack

Go · Bubble Tea · the full Charm ecosystem (lipgloss, bubbles, glamour, huh, log, harmonica).
Flutter carries the Phase 4 phone apps. **Not up for re-litigation.**

### 5.11 · The four inherited verdicts

Rendered by the orchestrator before this manual existed. Inherit them; do not re-decide them.

| # | Verdict | Reason, recorded |
|---|---|---|
| 1 | **Emoji glyphs — REJECTED** | Brand §8 already forbids emoji in product copy. They also break the monospaced grid (ambiguous or double advance width), render as tofu on many terminals, and are **inherently coloured**, which defeats the co-render rule's monochrome guarantee under `NO_COLOR`. The sanctioned replacement is `○ ◐ ◉ ⊘` with the ASCII column `[ ] [~] [*] [x]`. |
| 2 | **Frames drawn past 80 columns — REJECTED as drawn** | **No screen may *require* more than 80 columns to be correct.** 80 is the design width; 120 is progressive enhancement. Replaced by the five tiers: Tiny `<40` · Narrow `40–59` · Standard `60–79` · Wide `80–119` · ExtraWide `120+`. |
| 3 | **The embedded synthwave audio streamer — REJECTED OUTRIGHT, permanently closed** | It contradicts the ratified promises (no background telemetry; works with the network off; the only network traffic is services the player connected); adds an audio-device/cgo dependency to a program whose promise is *"it's a text program, it starts instantly"*; burns the redraw budget for zero product value; carries a music-rights surface a game tracker must not acquire; and is precisely the nostalgia-kitsch brand §1 rules out — synthwave-as-soundtrack is "remember the eighties?", not the era's idea of tomorrow. |
| 4 | **A single view carrying status summary + filter bar + table + detail pane at once — REJECTED as a fixed composition** | Replaced by a responsive one. The status summary is **not a region** — it is the one pinned summary row of R-10(c). The filter bar is a **mode** of the list below 120 columns, not a permanent region. The detail pane is not a default region. **The verdict set its floor at 80 columns; the spine resolved it at 120** on the same R-10(a) arithmetic — an 80-column split leaves the identity column around 15 characters. Below 120 the detail is a route, not a pane. |

### 5.12 · The FlowForge canon Zerado adopts

The Spacing Canon values (#2435 §4) are adopted **verbatim**. R-1…R-11 are inherited per
`00-design-brief.md` §4 — with **R-1 and R-9 having no Zerado analogue**, which is recorded
rather than faked.

---

## 6 · How to propose a change to a settled decision

You may propose one. You may not make one at the keyboard. Brand §10 is the procedure and it is
not ceremonial — each step exists because skipping it has broken something.

1. **Write the reason first.** What failed, on which screen, at which width, with the rendered
   artifact attached. "It would look better" is not a reason; "the state column misaligns at 40
   columns under `ambiguous-width=double`, here is the render" is.
2. **Route it to the owner.** Colour, type, glyph, motion, voice → `fft-brand-architect`.
   Composition → `fft-tui-architect`. Guidelines and the acceptance bar →
   `fft-design-architect`.
3. **Change the token in `tokens.css` and `tokens.json` — same commit.**
4. **Update the affected section of the brand manual — same commit.** A manual that disagrees
   with the tokens is worse than no manual.
5. **If it touches a colour used for text, recompute the contrast ratio.** Do not estimate it.
6. **If it touches a state colour, re-run the CVD simulation on all six pairs** with the
   Viénot / Brettel / Mollon dichromat model, measuring separation as CIEDE2000. **The 11.9
   minimum is the floor to protect.**
7. **If it touches a colour used in the terminal, re-derive the ANSI-256 index** by
   nearest-neighbour search in CIELAB. **Never adjust the old index by hand.**
8. **If it touches a state, the state-chip spec reopens too** (`01-design-system.md` §3) — and
   so does every screen spec that renders one.
9. **Founder ratification** for anything in §5.

**The one rule underneath all of it:** no surface may hard-code a hex value, a font stack, a
spacing step or a duration. If a value is needed and no token exists, **add a token — do not
type a hex.**

---

## 7 · How Zerado screens go wrong — the anti-patterns

Specific, not generic. Each has a section that fixes it.

| # | Anti-pattern | Fix |
|---|---|---|
| 1 | Dropping the state label for density, keeping glyph + colour | Co-render is all three, always. A TTY has **no accessibility API** — the label *is* the text alternative. `01-design-system.md` §3.6 |
| 2 | A screen that only reads well at 120 columns | 80 is the design width. Verdict 2 |
| 3 | Separating panes with a background fill | Invisible at 16 colours. `02-colour-budget.md` §7 |
| 4 | A filter mode, dialog or pane with no `esc` | A literal keyboard trap, WCAG 2.1.2 |
| 5 | Single-key shortcuts firing while a text input has focus | WCAG 2.1.4. Typing `d` must type `d`, not delete |
| 6 | Cyan used for emphasis "because it's the brand colour" | Cyan is **earned**. `02-colour-budget.md` §2 |
| 7 | `OFFLINE` rendered red | Offline is a kept promise, not a fault. `01-design-system.md` §12 |
| 8 | An ambient scanner under the header | Indeterminate waits only. Brand §7.1 |
| 9 | "Something went wrong." | Name what happened, why, the next action — **and what happened to the player's data** |
| 10 | A row index or a store ID in the identity column | R-10(a): a human game title, always |
| 11 | A background refetch snapping the ledger to the top | R-10(b): cursor and offset survive a row-set rebuild |
| 12 | The summary pushed off the bottom at 400 rows | R-10(c): pinned outside the scroll region |
| 13 | Multi-row ASCII-art wordmark as a splash | Brand §3.5: breaks on resize, on narrow terminals, and in any log capture |
| 14 | A field shown for a capability that doesn't exist yet | Never claim what isn't built. Omit the block; don't label it empty |
| 15 | Measuring string width with `len()` or rune count | `01-design-system.md` §1.2 — width-aware measurement, everywhere |
| 16 | Scanline textures, faux CRT damage, sepia, "vintage" wash | Retro-**future**. A MISS, not a variation |

---

## 8 · Open for the founder

1. **The 16-section screen-spec contract (§3).** It is deliberately strict — sections 8, 11, 12
   and 13 (first-run state, the refusal floor, the `NO_COLOR` render, and `esc` behaviour) are
   the ones usually skipped and the ones that cost the most later. Confirm the contract before
   `fft-tui-designer` fills it nine times.
2. **`fft-code-reviewer` judges against `00-design-brief.md` §10 and `02-colour-budget.md` §10 —
   and nothing else.** Confirm those two lists are the whole bar, so no reviewer adds a private
   one and no designer is failed by a rule they could not read in advance.
3. **Zerado has no in-house gold-standard screen yet.** Proposal: the library deck becomes it the
   moment it is ruled GOLDEN. Confirm the nomination.
