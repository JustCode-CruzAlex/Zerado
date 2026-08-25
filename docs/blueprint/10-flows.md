---
title: Zerado — the process flows
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-10
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: flowchart
ticket: "#2"
---

# The process flows

Seven flows, drawn. Each renders in both the brand-black and the cyanotype theme.

| Sheet | Flow | What it answers | Spec | Rendered |
|---|---|---|---|---|
| **03** | First run → a library | How does a player get from an empty machine to a library — including with no network? | [`ZRD-FLOW-01`](../adr/charts/ZRD-FLOW-01.chart.toml) | [svg](../adr/charts/svg/ZRD-FLOW-01.svg) · [cyanotype](../adr/charts/svg/ZRD-FLOW-01.cyanotype.svg) |
| **04** | The four states | Which transitions are the player's, and which is the machine's? | [`ZRD-FLOW-02`](../adr/charts/ZRD-FLOW-02.chart.toml) | [svg](../adr/charts/svg/ZRD-FLOW-02.svg) · [cyanotype](../adr/charts/svg/ZRD-FLOW-02.cyanotype.svg) |
| **05** | Sync, with its failure branches | What happens when it does not work? | [`ZRD-FLOW-03`](../adr/charts/ZRD-FLOW-03.chart.toml) | [svg](../adr/charts/svg/ZRD-FLOW-03.svg) · [cyanotype](../adr/charts/svg/ZRD-FLOW-03.cyanotype.svg) |
| **06** | Navigation | Where can a player go, and how do they get back? | [`ZRD-FLOW-04`](../adr/charts/ZRD-FLOW-04.chart.toml) | [svg](../adr/charts/svg/ZRD-FLOW-04.svg) · [cyanotype](../adr/charts/svg/ZRD-FLOW-04.cyanotype.svg) |
| **07** | The provider seam | Why does adding GOG add no screens? | [`ZRD-FLOW-05`](../adr/charts/ZRD-FLOW-05.chart.toml) | [svg](../adr/charts/svg/ZRD-FLOW-05.svg) · [cyanotype](../adr/charts/svg/ZRD-FLOW-05.cyanotype.svg) |
| **08** | The offline contract | Works, degrades, or refuses — and how is that decided? | [`ZRD-FLOW-06`](../adr/charts/ZRD-FLOW-06.chart.toml) | [svg](../adr/charts/svg/ZRD-FLOW-06.svg) · [cyanotype](../adr/charts/svg/ZRD-FLOW-06.cyanotype.svg) |
| **10** | The audio subsystem | What happens when a screen emits a cue — and when it cannot? | [`ZRD-FLOW-07`](../adr/charts/ZRD-FLOW-07.chart.toml) | [svg](../adr/charts/svg/ZRD-FLOW-07.svg) · [cyanotype](../adr/charts/svg/ZRD-FLOW-07.cyanotype.svg) |

```bash
for f in docs/adr/charts/ZRD-FLOW-*.chart.toml; do flowforge chart render "$f"; done
```

---

## Sheet 03 · First run → a library

The condition is checked **once**, at start-up: *is the library empty **and** has no provider ever
been connected?* If yes, `Z-01` is pushed as the root instead of `Z-04`.

It is not re-checked. A player who deletes every game does not get thrown back into first run —
that would be the program telling them they had undone their own setup.

Three doors, and **two of them work with no network**. The third is shown, disabled, with its
reason. Hiding it would make the player think Zerado cannot connect to stores at all.

## Sheet 04 · The four states

The whole model is one line:

```
effective_status = status_manual ?? derive(playtime_minutes, capabilities)
```

**Exactly one arrow on the sheet is automatic** — `NOT STARTED → IN PROGRESS`, when a sync reports
playtime on a game whose `status_manual` is `NULL`. Every other arrow is the player.

Two invariants the drawing exists to make unmissable:

- **A sync never writes `status_manual`.** Mark a game `ZERADO`, play three more hours: playtime
  updates, the state does not.
- **`ZERADO` is never inferred.** 100% achievements may *suggest* it in Phase 2 — a dismissible
  line the player accepts. A machine cannot hand someone the moment the product exists to create.

Full detail: [`05-state-machine.md`](./05-state-machine.md).

## Sheet 05 · Sync, with its failure branches

The sheet is mostly failure, because that is where the product is decided. Five classified
outcomes, three destinations:

| Outcome | Destination |
|---|---|
| streamed | `Z-03 DONE` — *"247 games. 6 finished. Last played: 3 weeks ago."* |
| `200` + empty · `401`/`403` | a **refusal** that names what, why, and the next action |
| no route / DNS · timeout / `5xx` | the **degrade banner** on `Z-04`, and the library below it still works completely |

The cancel path is on the sheet deliberately: `q` or `Esc` cancels the context, in-flight I/O
aborts, and **what was written stays written**. A half-synced library is a valid library, which is
why `sync_run.status` has a `partial` value.

## Sheet 06 · Navigation

A route stack with one overlay slot. `Z-04 Library` is the root and can never be popped; at most
one overlay exists at a time; pushing a route already on the stack **unwinds to it** rather than
duplicating it.

`Z-05 Game detail` is drawn as one view with **two hosts** — a route below 120 columns, a pane at
120 and above. That is the composition decision that makes the same spec buildable once and
mounted twice.

`Z-11 Fatal error` is deliberately off the graph: it replaces everything, from anywhere, and
there is no way back into the program from it.

Full key map and the exhaustive `Esc` table:
[`04-navigation-and-focus.md`](./04-navigation-and-focus.md).

## Sheet 07 · The provider seam

Interface segregation, drawn. `steam` is a `Provider` **and** a `Syncer`. `physical` is a
`Provider` **only** — not a Syncer with a hole in it, but a different shape of thing, whose sync is
a human typing.

Everything downstream reads `Capabilities`, never `ProviderID`. An `is_physical` boolean would
force every new feature to remember a special case, and eventually one would not.

The sheet's load-bearing claim: **a screen never talks to a provider, a provider never talks to a
screen, and between them sits the `Store`.** That is why the offline contract is structural rather
than a feature someone had to remember to add.

## Sheet 08 · The offline contract

Two questions, three classes, and **a feature is never in two of them**:

```
is it local?  ── yes ──►  WORKS      render it. No banner, no notice, no difference.
      │
      no
      ▼
is it cached? ── yes ──►  DEGRADES   render it WITH ITS AGE, plus the one-row banner.
      │
      no
      ▼
                          REFUSES    name what, why, and the retry key. Keep what was typed.
```

**Nine of the twelve Phase 1 screens are `WORKS`.** Two **`NEED THE NETWORK`** — `Z-02` and `Z-03`,
the two definitionally about reaching somewhere else — and one **`DEGRADES`**: `Z-15`, the cover
deck, which caches what it fetched and shows a designed tile for what it has not.

Zerado does not probe for connectivity — no heartbeat, no reachability check. That would be
exactly the background traffic the page promises does not exist. A request fails, and the failure
is classified.

Full table: [`07-offline-contract.md`](./07-offline-contract.md).

---

## Sheet 10 · The audio subsystem

Audio ships in Phase 1 ([`12-audio.md`](./12-audio.md)), and the sheet is mostly the ladder down to
silence, because that is where the design is.

A screen emits a **semantic cue** — `sync.done`, never a file name. Three gates follow: is audio
enabled (it is **off by default**), is it available, is the channel muted. Every "no" lands on
silence, and **silence is not an error**: the visible signal already happened, the co-render rule
guarantees it, and Settings says why **once** rather than banner-ing every screen.

Two things the drawing exists to make unmissable:

- **The cue is dropped before a frame is.** `Cue()` is a non-blocking send on a buffered channel;
  a full buffer discards the sound. A missed sound is never worth a dropped frame.
- **Bundled, never streamed.** No CDN, no fetch. That is what keeps three ratified promises intact
  and why audio is `WORKS` in the offline contract rather than an exception to it.

---

## A note on the renderer

`flowforge chart` draws `[[legend]]` swatches but does not render their captions — verified by
grepping the emitted SVG for a caption string (zero matches). Every sheet in this bundle therefore
carries its colour key as `[[free]]` annotation lines, which do render. Recorded so the next
person does not spend the same twenty minutes on it, and so it can be filed against the renderer.
