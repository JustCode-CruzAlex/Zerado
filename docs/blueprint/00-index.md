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
| | [`11-media-model.md`](./11-media-model.md) | **Games first, not games only** — the core entity is a media item; the four states verified across four types |
| | [`12-audio.md`](./12-audio.md) | **Audio ships in Phase 1** — bundled, off by default, fully removable |
| | [`13-handoffs.md`](./13-handoffs.md) | What this spine decides, and what it hands to `fft-database` and `fft-api-designer` |
| **B · design** | [`../design/`](../design/) | The design system and the designer manual |
| **C · screens** | [`../design/screens/`](../design/screens/) | One implementation-ready spec per Phase 1 screen |
| **decisions** | [`../adr/`](../adr/) | ADR-0001 — the four expensive-to-reverse decisions |

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

[`../adr/ADR-0001-zerado-foundational-architecture-provider-seam-persistence-n.md`](../adr/)
records the four, each with its alternatives and what it costs to reverse after Phase 1 code
exists:

| Decision | Cost to reverse after Phase 1 |
|---|---|
| **D1 · the provider seam shape** | High — every screen, every sync path and the schema key off it |
| **D2 · one SQLite file, pure-Go driver, credentials outside it** | High for the file; **low** for the driver |
| **D3 · a route stack with overlays, not tabs** | Medium — retraining users is the real cost, not the code |
| **D4 · the Phase 4 sync boundary** | Highest — it decides what the schema must carry from Phase 1 |
| **D5 · a media-polymorphic core** | Highest — retrofitting a type dimension rewrites every table and every seam |
| **D6 · audio, bundled and off by default** | Low for the design; the **licensing** question is the expensive part |

---

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
2. **IGDB's commercial terms are unresolved** and Zerado's funding model is commercial. The
   metadata seam is designed so the provider is swappable and so that *having no metadata at
   all* is a first-class designed state rather than an error state — see
   [`06-data-seams.md`](./06-data-seams.md) §3. What Phase 2 looks like if IGDB says no is
   written down there.
3. **Nine ANSI-256 indices are underived** — `--z-scanner-300` (error text), `--z-scanner-900`
   (the scanner track), `--z-amber-900`, `--z-text-tertiary` and five others have no derived
   terminal index. Four components ship an interim **uncoloured** rendering rather than a guessed
   one. Deriving them is `fft-brand-architect`'s job (brand manual §10 requires a
   nearest-neighbour search in CIELAB, never an eyeballed index), and until it is done those
   components are correct but plainer than designed.
4. **The light-mode state colours have never been CVD-verified.** `tokens.css` §10 defines a
   separate four-colour state set for light grounds; their contrast is recorded, their
   colour-vision separation is not. The dark set's own first draft failed at ΔE 8.8 under
   deuteranopia, so this is not a formality — and Phase 4 will meet system light mode.
5. **The contrast claim is precise, and worth stating precisely.** It is *"AA on Zerado's own
   ground and on five measured popular terminal themes, with `NO_COLOR` as the unconditional
   fallback."* An unqualified "WCAG AA" would not be true on an arbitrary user-chosen theme, and
   the product should not claim it.

6. **Terminal inline-image support is not assumed anywhere.** Cover art is a Phase 2 question
   and is inventoried, not specified. The Phase 1 deck is text, by design and not by omission.

---

## Rendering this bundle

Every document in the bundle carries an `archetype:` line and renders through the FlowForge
docs pipeline:

```bash
flowforge docs render docs/blueprint/09-erd.md      # erd archetype
flowforge docs render docs/blueprint/10-flows.md    # flowchart archetype
flowforge chart render docs/adr/charts/<name>.chart.toml   # both themes
```
