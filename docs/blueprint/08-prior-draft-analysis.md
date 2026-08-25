---
title: Zerado — the prior concept draft, analysed
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-08
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: concept-explainer
ticket: "#2"
---

# The prior concept draft — what survives, what does not, and why

The founder's instruction was a **clean slate**: analyse the earlier "CYBER-DECK" concept pass,
take what survives specialist judgement, and discard the rest **with a reason**. Four choices
need an explicit verdict before anything is inherited.

---

## 0 · What was actually available to analyse — stated plainly

**The draft artifact could not be located.** It was searched for by content
(`Neural Mood Scanner`, `CYBER-DECK`, `cyberdeck`) and by filename across the founder's entire
`~/Documents/projects/cruzalex` tree, including the FlowForge repository and every
`forgeplay-output` directory. There is no match. It is not in the Zerado repository, and issues
#1 and #2 carry no comments containing it.

So this analysis proceeds from the ticket's own enumeration of the draft's choices, which is
specific: *"ASCII mockups of a 'CYBER-DECK' dashboard and a 'Neural Mood Scanner', plus a mobile
sketch"*, and the four named choices — *"emoji glyphs in a terminal, frames drawn wider than 80
columns, an embedded synthwave audio streamer, and a full-width single-view dashboard that has to
hold a status summary, a filter bar, a table and a detail pane at once."*

Every verdict below rests on **canon and measurable constraints**, not on prose we could not
read. None of them would change if the draft were found — but that claim deserves testing, so if
the draft surfaces it should be re-checked against this document.

---

## 1 · Emoji glyphs in a terminal — **REJECTED**

### The verdict
Emoji do not appear anywhere in Zerado: not in state, not in labels, not in copy, not in the
footer, not in log output.

### Why — canon first
The brand manual §8 already settles it as a standing rule: **"No emoji in product copy. Glyphs
carry state; emoji carry nothing."** That alone is dispositive.

### Why — the technical case, which is independent and just as strong

1. **They break the monospaced grid the entire layout is built on.** Most emoji are East Asian
   Wide (two cells), many are ambiguous-width, and several — the ZWJ sequences and anything with
   a variation selector — are rendered at widths that differ *between terminals* running the same
   font. Every column budget in [`02-composition.md`](./02-composition.md) is computed in cells.
   A glyph whose cell count is a per-terminal opinion makes those budgets unenforceable.
2. **They are inherently coloured, which defeats the co-render rule's monochrome guarantee.**
   The whole point of the co-render rule is that removing any one channel still leaves the state
   unambiguous. Emoji are colour *and* shape fused into one channel: strip colour, and most
   terminals render a coloured emoji anyway, or a black-and-white outline that no longer matches
   what the legend showed. Under `NO_COLOR` — where the brand promises *zero* SGR sequences — an
   emoji still arrives in colour, from the font. It is the one glyph class `NO_COLOR` cannot turn off.
3. **They degrade to tofu, not to a fallback.** A missing box-drawing character degrades to a
   recognisable near-miss. A missing emoji degrades to `□`, identically for all of them, which
   turns four distinct states into four identical boxes.
4. **They cannot be typed, searched, or read aloud.** Nobody greps a bug report for a
   game-controller pictograph, and a screen reader announces an emoji's CLDR name, which is not
   the state's name. `ZERADO` is a word: readable, greppable, speakable.

### What replaces it — already decided upstream
The geometric progression from the brand manual §4.3, which was chosen as a **sequence**: an
empty ring, a ring half filled, a ring with a solid core, a ring struck through.

```
  ○  NOT STARTED      [ ]
  ◐  IN PROGRESS      [~]
  ◉  ZERADO           [*]
  ⊘  ABANDONED        [x]
```

All four are single-cell, unambiguous-width, present in every common terminal font, meaningful in
one ink, and carry an ASCII column for the case where even they cannot be relied on. The story
reads left to right without any colour at all.

**Nothing from this choice survives.**

---

## 2 · Frames drawn past 80 columns — **REJECTED as drawn**

### The verdict
**No screen may require more than 80 columns to be correct.** 80 is the design width and the
guaranteed floor. Widths above it are *progressive enhancement*, not the baseline.

### Why
A frame drawn at, say, 100 columns has one of two behaviours on an 80-column terminal, and both
are failures: it **wraps**, which turns box drawing into confetti and destroys every alignment in
the screen; or it **clips**, which silently hides the right-hand column — and the right-hand
column is where playtime and source live.

80×24 is not a nostalgic number. It is the default size of a fresh terminal on macOS and on most
Linux desktops, the width of a tmux pane in a vertical split on a laptop, and the size of most
CI and SSH sessions. A product that assumes more than 80 columns is a product that is broken the
first time anyone opens it in the default configuration.

### What survives, and how the intent is honoured
The draft's instinct — that a wide terminal should be *used*, not letterboxed — is correct and is
kept. It is expressed through the tier system instead of through a fixed wide frame:

| Tier | What the extra width buys |
|---|---|
| Standard 60–79 | The `source` column returns |
| **Wide 80–119** | Full one-line rows, full state chips. **Everything is correct here** |
| ExtraWide 120+ | A second region — the detail pane — and 28 visible rows instead of 12 |

So a 140-column terminal gets a genuinely better screen than an 80-column one. It just does not
get a *different* screen, and the 80-column one is never the degraded case.

**The intent survives. The mechanism does not.**

---

## 3 · The embedded synthwave audio streamer — ~~REJECTED OUTRIGHT~~ · **SUPERSEDED**

> ### SUPERSEDED — founder direction, 2026-08-25
>
> **Audio ships in Phase 1.** *"The audio is part of the Phase 1."* The verdict below is no longer
> operative and the design is in [`12-audio.md`](./12-audio.md).
>
> **The reversal is not a change of mind about the same object — the object changed, twice.** What
> was rejected was a *network streamer, always on, that the player never asked for*. What ships is an
> **opt-in subsystem, off by default**, whose music is **internet radio the player chooses and can
> stop in one keystroke**, and whose only always-available part is a handful of **local** interface
> cues.
>
> *(The first reversal specified **bundled tracks**; a further founder direction the same day
> replaced bundling with **streamed stations**. Both moves are recorded because the second one
> **dissolved** the licensing objection rather than answering it.)*
>
> It is kept here rather than deleted, because **four of its five reasons survived the reversal as
> engineering requirements** — and the fifth is the one the founder decided differently. A verdict
> that vanishes teaches nothing; one that shows which of its arguments held is a design brief.
>
> | Original objection | After the reversal |
> |---|---|
> | Contradicts three ratified promises | **Partly retired, and one part is genuinely open.** *"No background telemetry"* and *"works with the network off"* no longer bite: the player starts the stream and can stop it, and the offline promise is about the **library**. But *"the only network traffic is Zerado reaching out to the services you've connected"* **does** bite — a station uses **no key of the player's**, and the sentence names Steam and price data specifically. The reading that a chosen station *is* a connected service is sound and **is on the founder's gate list**, not settled here ([`13-handoffs.md`](./13-handoffs.md) §5) |
> | Costs the pure-Go single binary | **Survives as a requirement** — the player sits behind a build tag; the default build stays pure Go |
> | Costs the 60 fps / no-leak budget | **Survives as a requirement** — non-blocking buffered cue, one owned goroutine, cue dropped before a frame is |
> | Acquires a music-rights surface | **DISSOLVED.** Nothing is bundled, so there is nothing to license, attribute or weigh — and the product is no longer commercial either. The objection was removed rather than answered |
> | Nostalgia-kitsch versus retro-future | **The founder's call, and it is made.** The bar the sound must clear is written down instead |
>
> **Provenance:** the direction reached this session as an agent relay, which by its own header
> carries no ratification authority. It was acted on because this is document revision on a draft
> PR — reversible, and confirmed or restored by the founder at the gate.

### The original verdict, as recorded on 2026-08-25 — superseded
Not deferred, not scoped to a later phase, not made optional. **Permanently out of scope**, and
recorded here so it is not re-proposed as a Phase 3 delight feature.

### Why — five independent reasons, any one of which is sufficient

1. **It contradicts three ratified public promises at once.** The page says *"no telemetry
   running in the background"*, *"Runs with the network off"*, and *"The only network traffic is
   `Zerado` reaching out to the services you've connected — Steam, price data — using your own
   keys."* An audio streamer is a persistent outbound connection, to a service the player did not
   connect, using nobody's key, running continuously in the background. It is not adjacent to
   those promises; it is the counterexample to all three.
2. **It contradicts the one-sentence product argument.** *"It's a text program. It starts
   instantly, works offline, and your library is one file you own."* Audio playback in Go means an
   audio-device binding — which in practice means cgo, which means the pure-Go single-binary
   distribution decided in [`06-data-seams.md`](./06-data-seams.md) §5.3 is gone, along with the
   cross-compilation matrix that makes releases cheap.
3. **It costs the performance bar for nothing.** The spine's bar is 60 fps and zero goroutine
   leaks. A streamer adds a decode loop, a network reader, a device writer, and their lifecycles —
   the exact goroutine class hardest to shut down cleanly — in a program whose actual job is
   rendering a list of games.
4. **It acquires a rights surface the product must not have.** Streaming music is a licensing
   question. A game-library tracker has no business acquiring one — and since 2026-08-25 Zerado is
   donation-funded and non-commercial, which makes it less equipped to carry one, not more,
   and "we just embed someone else's stream" is the version of that sentence that ends in a
   takedown.
5. **It is the exact nostalgia-kitsch the brand manual §1 rules out.** This is the reason that
   matters most, because it is the one that would still apply if the other four were solved.
   Synthwave-as-soundtrack is *"remember the eighties?"* — it is the decade quoted back at itself.
   The brand's stated register is *the eighties' own idea of tomorrow*: **the DeLorean was not
   playing synthwave. It was a machine that looked like it had arrived from somewhere.** A
   soundtrack tells the player what era to feel. A machine that just works, instantly, in amber on
   black, tells them what era it came *from*. The second is the brand; the first is a costume.

~~**Nothing from this choice survives.**~~ **Superseded — see the box at the head of this section.**

---

## 4 · One view holding a status summary, a filter bar, a table and a detail pane — **REJECTED as a fixed composition; the intent is kept, responsively**

### The verdict
Four simultaneous regions is not a composition, it is four screens in a trench coat. The **intent**
— density, everything at hand, no hunting — is correct and is preserved. The **mechanism** is
replaced.

### Why the fixed version fails, in numbers
At the 80-column design floor the body is **74 columns**
([`02-composition.md`](./02-composition.md) §1.3). Divide that among a table and a detail pane and
the table gets roughly 44 columns. After the state chip (13) and the playtime figure (6) and their
gutters, the **title column** — which TUI Design Manual **R-10(a)** makes mandatory and
human-readable — falls to about fifteen characters:

```
  ◉ ZERADO       Return of the Ob…      9h │ Return of the Obra Di…
  ◐ IN PROGRESS  The Legend of Ze…     12h │ ◐ IN PROGRESS · 12h
  ○ NOT STARTED  Disco Elysium: T…      0h │
                                           │ (no room for a sinopse)
```

Both halves fail. The table has lost the column that identifies the rows, and the detail pane is
too narrow to hold the thing it exists to show. **Four regions in 74 columns is not a dense
screen; it is two broken ones.**

Vertically it is worse: at 24 rows the body is 16 rows, and a summary block plus a filter bar plus
a column header plus a scroll indicator leaves single digits for the list itself.

### What replaces it — each of the four, reassigned

| The draft's region | What it becomes | Where |
|---|---|---|
| **Status summary** | **Not a region — one pinned row.** R-10(c): pinned outside the scroll region so a 400-game library can never push it off the bottom | `Z-04`, body row 1 |
| **Filter bar** | **A mode, at every tier.** It occupies two body rows only while active, and the active filter stays visible in the summary row after the editor closes | `Z-07` |
| **Table** | The screen. Single-pane below 120 columns | `Z-04` |
| **Detail pane** | **A second region only at ExtraWide.** Below 120 it is a route reached with `Enter` — the same view with a different host | `Z-05` |

The filter bar is a mode **at every tier**, including 120+, and that deserves its reason: at
ExtraWide there *is* room for a second region, and it is worth more spent on the detail pane than
on a filter editor the player uses for a few seconds at a time. A permanent filter rail would be a
region that is idle almost always.

### The rule this produces
> **Two regions only when there is room for both to be correct. Never more than two.**

**The intent survives — density, and everything reachable. The four-region composition does not.**

---

## 5 · What "retro-future, not retro-nostalgia" means in a terminal, specifically

The draft's name — *CYBER-DECK* — and its *Neural Mood Scanner* point at the failure mode the
brand manual §1 rules out, so it is worth writing the line down where a builder will hit it.

The brand manual §9 is explicit that the terminal is where the brand is **quietest**: *"Terminal
trades expression for density and honesty: no glow, no gradients, no fluid type; hierarchy comes
from the monospaced grid, box drawing, spacing and the four-state palette."*

So in a TUI the retro-future register is carried by four things and nothing else:

| Carries the register | Does not |
|---|---|
| **Amber on black.** The phosphor readout, as the ambient voice | Neon *gradients* — there are none in a terminal |
| **The tracked-out uppercase annotation** — `[0] ZERADO ▸ LIBRARY`, the cockpit label | Fake scanlines, faux CRT curvature, simulated flicker |
| **The scanner sweep**, once, for a genuinely indeterminate wait | An ambient sweep that runs for decoration |
| **Chassis discipline** — precise alignment, real margins, nothing flush to an edge | ASCII art, box-drawing "logos", multi-row banners |

And the naming rule that follows: **the product does not talk like a prop.** There is no "Neural
Mood Scanner" in Zerado; there is a **mood picker**, and it asks *"What is tonight for?"*. The
brand manual's voice section is unambiguous — dry, concrete, never performing. A machine that
names its own features dramatically is a machine that does not trust the features.

KITT is the reference precisely because KITT's interface was *plain*: amber readouts, labelled
switches, one moving light. The drama was that it worked.

---

## 6 · What the draft got right, and is inherited

Recorded because a clean slate is not a rejection of everything, and the parts that survive
should be attributed rather than silently reinvented:

1. **The library as the home screen.** Not a dashboard above it, not a menu in front of it. The
   list of games is the product, and it is what opens. Kept — `Z-04` is the root and cannot be
   popped.
2. **Density as a goal.** A terminal product's advantage over a web app is that it can show a
   lot without ceremony. Kept, and paid for by tiering rather than by cramming.
3. **A mood-led way in.** Sorting by what tonight is for, rather than by genre, is the product's
   whole argument and it is in the published copy. Kept as `Z-13`, in Phase 2, where the roadmap
   puts moods.
4. **A phone surface exists.** The draft's mobile sketch was pointing at something real: the old
   form and the new form as one product. Kept as Phase 4, with the identity bridge specified in
   deliverable B rather than sketched.

---

## 7 · The four verdicts, in one table

| Choice | Verdict | Survives? |
|---|---|---|
| Emoji glyphs in a terminal | **REJECTED** — brand manual §8, plus width, colour, tofu and searchability | Nothing |
| Frames drawn past 80 columns | **REJECTED as drawn** — 80 is the floor; wider is enhancement | The intent, via the tier system |
| Embedded synthwave audio streamer | ~~REJECTED OUTRIGHT~~ → **SUPERSEDED 2026-08-25.** Audio ships in Phase 1 — **streamed radio, nothing bundled**, off by default | Three engineering objections as requirements; the licensing one dissolved; one promise question still open |
| One view with summary + filter + table + detail | **REJECTED as fixed** — 74 columns cannot hold four correct regions | The intent, via one pinned row + a mode + two tiers |
