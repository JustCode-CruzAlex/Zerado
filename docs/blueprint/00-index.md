---
title: Zerado — the spine
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-00
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: implementation-plan
ticket: "#2"
---

# Zerado — the spine

**The blueprint bundle for ticket #2. Documents only — no Go was written for this ticket.**

Zerado is a terminal-first game library. It reads a player's collection across stores, marks
each title *not started* / *in progress* / *zerado* / *abandoned*, and answers "what do I play
tonight?" This bundle says what the product **is** — which screens exist, how a player moves
between them, what each one owes the player, and where the seams are that later phases hang
off. The brand manual already said what it looks like.

---

## What is in this bundle

| Part | Document | What it settles |
|---|---|---|
| **A · the spine** | [`01-screen-inventory.md`](./01-screen-inventory.md) | Every screen the product will ever have, named, numbered, phased |
| | [`02-composition.md`](./02-composition.md) | Composition and layout budget per screen, in rows and columns |
| | [`03-responsive.md`](./03-responsive.md) | 40 / 60 / 80 / 120 columns — what collapses, hides, refuses |
| | [`04-navigation-and-focus.md`](./04-navigation-and-focus.md) | The route model, global keys, focus, and Escape everywhere |
| | [`05-state-machine.md`](./05-state-machine.md) | The four states, legal transitions, derived vs user action |
| | [`06-data-seams.md`](./06-data-seams.md) | Provider · metadata · price · persistence · credential seams |
| | [`07-offline-contract.md`](./07-offline-contract.md) | What works with the network off, what degrades, how it is shown |
| | [`08-prior-draft-analysis.md`](./08-prior-draft-analysis.md) | The prior concept draft — what survives, what does not, why |
| | [`09-erd.md`](./09-erd.md) · [`10-flows.md`](./10-flows.md) | The data model and the process flows, **drawn** |
| | [`11-media-model.md`](./11-media-model.md) | **The door stays open** — one column and a table name, and nothing else generalised |
| | [`12-audio.md`](./12-audio.md) | **Audio ships in Phase 1** — streamed radio, off by default, fully removable |
| | [`16-i18n.md`](./16-i18n.md) | **i18n from the first line** — no string literal in code, and the second language is one file |
| | [`17-images.md`](./17-images.md) | **Cover art is foundational** — Kitty/Ghostty, with a supported degrade |
| | [`13-handoffs.md`](./13-handoffs.md) | What this spine decides, and what it hands to `fft-database` and `fft-api-designer` |
| | [`14-contradictions-closed.md`](./14-contradictions-closed.md) | The cross-check register — all 29 findings, enumerated |
| **B · design** | [`../design/`](../design/) | The design system, the designer manual, and [the theme system](../design/05-theme-system.md) |
| **C · screens** | [`../design/screens/`](../design/screens/) | **Twelve** implementation-ready specs — one per Phase 1 screen |

**Deliverable C, in full** — each carries all sixteen required sections from
[`../design/03-designer-manual.md`](../design/03-designer-manual.md) §3, and every mockup is
verified to an exact cell count:

| Screen | Spec |
|---|---|
| `Z-01` First run | [`Z-01-first-run.md`](../design/screens/Z-01-first-run.md) |
| `Z-02` Connect a store | [`Z-02-connect-a-store.md`](../design/screens/Z-02-connect-a-store.md) |
| `Z-03` Sync | [`Z-03-sync.md`](../design/screens/Z-03-sync.md) |
| `Z-04` Library | [`Z-04-library.md`](../design/screens/Z-04-library.md) |
| `Z-05` Game detail | [`Z-05-game-detail.md`](../design/screens/Z-05-game-detail.md) |
| `Z-06` Set status | [`Z-06-set-status.md`](../design/screens/Z-06-set-status.md) |
| `Z-07` Filter and search | [`Z-07-filter-and-search.md`](../design/screens/Z-07-filter-and-search.md) |
| `Z-08` Add a game by hand | [`Z-08-add-a-game-by-hand.md`](../design/screens/Z-08-add-a-game-by-hand.md) |
| `Z-09` Settings | [`Z-09-settings.md`](../design/screens/Z-09-settings.md) |
| `Z-10` Help and key map | [`Z-10-help-and-key-map.md`](../design/screens/Z-10-help-and-key-map.md) |
| `Z-11` Fatal error | [`Z-11-fatal-error.md`](../design/screens/Z-11-fatal-error.md) |
| `Z-15` Cover deck | [`Z-15-cover-deck.md`](../design/screens/Z-15-cover-deck.md) |

The drawings live at [`../adr/charts/`](../adr/charts/) as `.chart.toml` specs and render to
`svg/` in both the brand-black and cyanotype themes.

---

## What this bundle inherits and does not touch

The brand manual (`brand/brand-manual.md`, revision A) is **finished and ratified**. This
bundle applies it. It does not restate it and it does not amend it. In particular the
following are settled upstream and are treated here as facts:

- the four game states, their colours, their glyphs, their ASCII fallbacks and their labels;
- the co-render rule — colour **and** glyph **and** label, all three, every state;
- the colour budget — amber ambient, **cyan earned**, red motion-and-alarm only, chrome structural;
- `NO_COLOR` emits zero SGR sequences, and the 16-colour floor holds;
- the type system, the scanner-sweep motion spec, the voice, and the casing convention.

The six ratified public promises in `ratification/decisions.md` and the published copy in
`content/landing-copy.md` are **binding constraints on the architecture**, not context. Where
a design choice in this bundle exists only because a promise required it, the promise is cited
at the point of the choice.

---

## The three questions at the gate

The ticket asks the founder exactly three things. Here is where each is answered.

> ### Two founder directions landed after the first draft, 2026-08-25
>
> Both fold into this bundle rather than a follow-up, because each changes the spine at its core.
>
> **1 · The data model is media-polymorphic.** The core entity is a **media item**, not a game.
> Games are the first type; books are the second; films and series a plausible third. **Phase 1
> ships games only** — no book screens, no book providers, no `--type` flag. What changed is the
> *shape*, because it is the single most expensive thing in this bundle to retrofit.
> [`11-media-model.md`](./11-media-model.md) verifies the four states across all four types and
> names the two findings the check surfaced.
>
> **2 · Audio ships in Phase 1.** This **reverses** the verdict in
> [`08-prior-draft-analysis.md`](./08-prior-draft-analysis.md) §3. The reversal is recorded with its
> provenance rather than silently applied, and the original reasoning is kept visible: four of its
> five objections survived as engineering requirements, and one — music licensing — is an open
> founder decision. [`12-audio.md`](./12-audio.md).
>
> **Provenance, stated plainly:** both directions reached this session as an **agent relay**, which
> by its own header carries no ratification authority. They were acted on because this is document
> revision on a draft PR — reversible, and confirmed or restored by the founder at the gate. They
> are flagged here so that confirming them is a decision rather than an assumption.

### 1 · Is the Phase 1 screen set right?

Answered in [`01-screen-inventory.md`](./01-screen-inventory.md) §3. The architect's amendment
to the ticket's proposed list, in one paragraph:

**Two screens added, one renamed, one removed.**

- **Added — `Z-08 Add a game by hand`.** The ticket requires physical copies to be first-class
  from day one. The proposed screen set had no screen for entering one. A promise with no
  screen behind it is not a promise. This is the single largest amendment and the one that
  most changes Phase 1's size.
- **Added — `Z-11 Fatal error`.** "Error and offline states" as listed is not one screen; it
  is three different things at three different altitudes. They are separated in §3.2.
- **Renamed — `Z-02 Connect a store`**, not "connect Steam". Steam is the only instance in
  Phase 1, but the screen renders from the provider descriptor, so PlayStation, GOG and EA
  add **zero screens** later. This is the navigation model paying for itself.
- **Removed — a splash screen.** The page promises "it starts instantly." A splash screen is
  a promise being broken on the first frame. Recorded as a rejected screen in §4.

### 2 · Does the design system clear the DeLorean-and-KITT bar?

Deliverable B answers this. The spine's contribution is structural: the terminal is where the
brand is *quietest* by the brand manual's own §9 — no glow, no gradients, no fluid type; the
hierarchy comes from the monospaced grid, box drawing, spacing, and the four-state palette.
The retro-future register in a TUI is carried by **the amber-on-black readout, the tracked-out
uppercase annotation, the scanner sweep as the one motion, and the chassis discipline** — not
by decoration. §5 of [`08-prior-draft-analysis.md`](./08-prior-draft-analysis.md) states the
line between the two.

### 3 · Are the settled decisions the right ones to freeze?

[`ADR-0001-zerado-foundational-architecture.md`](../adr/ADR-0001-zerado-foundational-architecture.md)
records **all nine**, each with its alternatives and what it costs to reverse after Phase 1 code
exists:

| Decision | Cost to reverse after Phase 1 |
|---|---|
| **D1 · the provider seam shape** | High — every screen, every sync path and the schema key off it |
| **D2 · one SQLite file, pure-Go driver, credentials outside it** | High for the file; **low** for the driver |
| **D3 · a route stack with overlays, not tabs** | Medium — retraining users is the real cost, not the code |
| **D4 · the Phase 4 sync boundary** | Highest — it decides what the schema must carry from Phase 1 |
| **D5 · the door stays open** | **Low** — one column and a table name. Revision A's polymorphic core was pruned by founder direction |
| **D6 · audio, streamed and off by default** | Low — and the licensing question is **closed**, because nothing is bundled |
| **D7 · themes are data, and must pass the four-state contract** | Medium — the token contract is what every component reads |
| **D8 · terminal images are foundational** | Medium-high — Phase 1 screens are designed around covers existing |
| **D9 · i18n from the first line** | **Highest of the three** — retrofitting i18n is the canonical expensive migration |

---

## The bundle checked itself

The three deliverables were not written in isolation and then stapled together. Each specialist read
the others' output and cross-checked it, and **29 contradictions were found and closed before this
bundle was surfaced** — [enumerated one by one](./14-contradictions-closed.md), so the count is
auditable rather than asserted — a stale row budget, a status bar in a forbidden row, refusal copy that named
a key which cannot fire, a warning glyph whose stated reason was factually wrong, an ASCII fallback
that covered one column of a screen made of the same problem, and a decision about a player's own
data that had simply never been made.

Two are worth naming because they show the mechanism working rather than the documents being tidy:

- **`fft-tui-designer` refused to write copy for a case the seam had not decided** — what happens to
  a game a sync stops returning — rather than inventing an answer about someone's file. It is now
  decided ([`06-data-seams.md`](./06-data-seams.md) §2.4): tombstoned, never deleted.
- **Two width assertions of mine did not survive checking**, and both are recorded as corrections
  rather than quietly repaired. In a bundle whose value is that its facts are checked, a wrong
  premise that happens to reach the right answer is still a defect.

## The honest gaps

Stated here rather than buried, because a blueprint that hides its holes is worse than one
that has none.

1. **The prior "CYBER-DECK" concept draft could not be located on disk.** It was searched for
   across the founder's whole projects tree by content (`Neural Mood Scanner`, `CYBER-DECK`)
   and by filename; there is no match. The analysis in
   [`08-prior-draft-analysis.md`](./08-prior-draft-analysis.md) therefore proceeds from the
   ticket's own enumeration of the draft's choices, which is specific enough to rule on each.
   Every verdict names the evidence it rests on. If the draft surfaces, the verdicts should be
   re-checked against it — none of them depend on prose we could not read, but that claim is
   worth testing.
2. **IGDB is answered, with one founder action outstanding.** Affiliate links are dropped, so
   Zerado is cleanly **non-commercial** — free software, donation-supported, zero revenue — and
   IGDB's published test is whether the *project generates revenue*. **This is a reading of their
   published rationale, not a legal opinion:** the outstanding action is a direct confirmation from
   IGDB that a donation-funded open-source project qualifies for the free tier. The metadata seam
   stays provider-agnostic regardless, because that hedge was right for reasons that never depended
   on IGDB's answer.
3. **Two token vocabularies meet here, and both are correct — which nothing said.** The
   underived-index list names `--z-scanner-900` (`#5C1414`) and `--z-amber-900` (`#8A5E00`), while
   the component-facing colour list names `--z-scanner-track` and `--z-primary-muted` for the same
   two colours. A reviewer reasonably read that as one vocabulary needing to win. **Verified at
   source in `tokens.css`, it is not:** `--z-scanner-900` and `--z-amber-900` are declared
   **primitives** carrying the raw hex (lines 61, 71), and `--z-scanner-track` and
   `--z-primary-muted` are **semantics** declared as `var()` of them (lines 102, 116). Brand manual
   §10 fixes the three layers and says a primitive is *"never referenced by a component"* — so the
   component list is right to name the semantic, and the derive-these list is right to name the
   primitive, because an ANSI-256 index is derived from a **raw value**. Neither should change; the
   layer relationship simply needed stating, and now does.
4. **`tokens.css` §10 records `7.30:1` for `--z-state-abandoned`; it measures `7.67:1`.** One
   figure out of 35 brand contrast values that did not reproduce — the other 34 did, exactly. Small,
   and worth fixing precisely because the other 34 are right: a table that is 34/35 correct is
   trusted, which is what makes the one wrong entry dangerous.
5. **Nine ANSI-256 indices are underived** — the primitives `--z-scanner-300` (error text),
   `--z-scanner-900` (the scanner track), `--z-amber-900` (inert track), `--z-chrome-500`
   (tertiary text) and five others have no derived terminal index. The list is
   [`../design/00-design-brief.md`](../design/00-design-brief.md) §9, which now carries an explicit
   **primitive → semantic** column so neither layer has to be read against the other. Four components ship an interim **uncoloured** rendering rather than a guessed
   one. Deriving them is `fft-brand-architect`'s job (brand manual §10 requires a
   nearest-neighbour search in CIELAB, never an eyeballed index), and until it is done those
   components are correct but plainer than designed.
6. **The light state set has now been verified, and it FAILS.** Previously carried as
   *"never verified"*. Measured under the method pinned in
   [`../design/05-theme-system.md`](../design/05-theme-system.md) §2.1: `not started × zerado`
   **5.41** under protanopia and `zerado × abandoned` **8.91** under deuteranopia, against a floor
   of 10.0. The cause is that brand §4.5's light `not started` is rotated 173° in hue and lands
   1.1° from `#9FB0C6` — **the exact blue-cast steel §4.4 rejected on the dark side.** The dark
   set's correction was never carried to paper. `fft-brand-architect` repairs it; until then no
   light theme ships, which is D7's gate doing its job.
7. **The ΔE figures quoted throughout are the brand manual's measurements, not this bundle's.**
   Every `ΔE` in these documents is quoted from `brand-manual.md` §4.4, which states the model
   (Viénot, Brettel & Mollon 1999 + CIEDE2000) but **does not pin the model variant or the white
   point** — so independent implementations do not land on the same digits, and three of them
   did not: the manual's **11.9**, an independent reviewer's **12.07**, and this bundle's now-pinned
   method's **11.81**. [`../design/05-theme-system.md`](../design/05-theme-system.md) §2.1 **pins
   it** — matrices, D65, the linear-RGB clamp, CIEDE2000 with `kL=kC=kH=1` — because a validation
   gate cannot rest on a method two implementations read differently. An independent
   check of the rejected `#9FB0C6` reproduced **11.02** against the manual's **8.8**, and of the
   shipping floor **12.07** against **11.9** (1.4%). The *load-bearing* claims reproduce: the
   tightest pair really is `zerado × abandoned`, and the warm grey really does separate far better
   than the blue-cast grey it replaced. **Pinning the variant and white point upstream is
   `fft-brand-architect`'s job**, and until it is done these numbers should be read as the
   manual's, not as independently verified.
8. **The contrast claim is precise, and worth stating precisely.** It is *"AA on Zerado's own
   ground and on five measured popular terminal themes, with `NO_COLOR` as the unconditional
   fallback."* An unqualified "WCAG AA" would not be true on an arbitrary user-chosen theme, and
   the product should not claim it.

9. **Cover art is foundational, and the degrade is the design.** Reversed from revision A, which
   deferred it. Kitty and iTerm2 are targeted, Sixel is deferred with reasons, and a terminal
   without image support is a **supported configuration** — the full text deck plus a once-only
   dismissible note. See [`17-images.md`](./17-images.md).

---

## The environment, in one place

Seven variables, specified across seven documents. `fft-tui-designer` pointed out that they had
five homes and no index — and a builder needs the whole set in one look, because they interact:
`NO_COLOR` and `ZERADO_ASCII` both change what a cell contains, and both change what a width
measurement returns.

| Variable | Effect | Specified in |
|---|---|---|
| `NO_COLOR` | **Zero SGR sequences.** Not a reduced palette — none. Glyph and label carry every state | brand manual §5.4 · [`03-responsive.md`](./03-responsive.md) |
| `ZERADO_ASCII` | Switches the **whole glyph vocabulary** to ASCII, not just the state column | [`03-responsive.md`](./03-responsive.md) §5c |
| `ZERADO_NO_AUDIO` | No sound at all. Settings shows *overridden*, not *off* | [`12-audio.md`](./12-audio.md) §10 |
| `ZERADO_NO_IMAGES` | Forces `NullImages`; the text deck renders | [`17-images.md`](./17-images.md) §4 |
| `ZERADO_REDUCED_MOTION` | The scanner pip parks at centre, full weight, and does not travel | [`03-responsive.md`](./03-responsive.md) §5 |
| `ZERADO_DB` | Overrides the library file path | [`06-data-seams.md`](./06-data-seams.md) §5.2 |
| `ZERADO_LANG` | Overrides locale; wins over the setting and over `LANG` | [`16-i18n.md`](./16-i18n.md) §5 |

**Two properties hold across all seven:** each is honoured **without a restart** where that is
possible, and each is **visible in `Z-09 Settings`** as an override rather than silently changing a
value the player thinks they chose. A setting the environment has overridden must never render as
though the player set it.

---

## Rendering this bundle

Every document in the bundle carries an `archetype:` line and renders through the FlowForge
docs pipeline:

```bash
flowforge docs render docs/blueprint/09-erd.md      # erd archetype
flowforge docs render docs/blueprint/10-flows.md    # flowchart archetype
flowforge chart render docs/adr/charts/<name>.chart.toml   # both themes
```
