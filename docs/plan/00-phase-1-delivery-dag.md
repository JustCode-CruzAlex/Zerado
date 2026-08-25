---
title: "Zerado — the Phase 1 delivery DAG, four sprints, thirty-eight work units"
discipline: DELIVERY
doc-no: ZRD-DAG-00
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: implementation-plan
ticket: "#7"
---

# The Phase 1 delivery DAG

**Thirty-eight work units, four sprints, 2026-08-25 → 2026-10-19.**

Every unit carries its **own** duration, stated before any date existed. Every date below is
**derived** — start is when that unit's dependencies finish, end is start plus the padded
duration. Nothing here was typed onto a ticket, and the founder was asked exactly one scheduling
question, which he had already answered: **Sprint 0 starts 2026-08-25.**

---

## 0 · What this was built from

| Input | State | What it contributed |
|---|---|---|
| [`../adr/ADR-0001`](../adr/ADR-0001-zerado-foundational-architecture.md) | **RATIFIED 2026-08-25**, with one amendment | The nine decisions that make a seam a seam. D1–D9 are the reason the units are cut where they are |
| [`../blueprint/`](../blueprint/00-index.md) — 17 documents | ratified | The spine: which screens exist, what each owes the player, where the seams are |
| [`../design/`](../design/00-design-brief.md) — 19 documents | ratified | The design system, the theme contract, the spacing canon |
| [`../design/screens/`](../design/screens/) — 12 specs, 11 572 lines | **11 of 12 ratified**; `Z-10` superseded (§9.2) | The buildable units. Each carries a state table and a key map, which is what makes a screen sizable at all |
| Ticket **#5** — ERD and schema | **in flight, parallel** | Node `#5`. Every unit that touches a table sits behind it |
| Ticket **#6** — interface contracts | **in flight, parallel** | Node `#6`. Every unit that touches a signature sits behind it |
| Ticket **#8** — key-map component visual spec | **being created** | Node `#8`. Gate on `Z2-8` (§9.2) |

`#5` and `#6` are running **right now**, in parallel with this plan. They are the first two nodes
of the graph and nothing waited on them: they were sequenced, not awaited.

---

## 1 · The rule this plan is built on

1. **Duration first, dates never.** Each unit states how long *that* piece of work takes, in days
   or weeks. **Hours are refused by name** — converting an hour to a calendar span needs a
   working-day length nobody has agreed, and a bar wrong by a factor of three is worse than a
   refusal.
2. **+20% breathing room per ticket, never pooled.** `3d` becomes a **4-day bar**; padding rounds
   **up**. Slack pooled on a sprint total does not live where the risk is.
3. **Dates derive from the DAG.** Two units with no edge between them **overlap**, and that
   overlap is the entire reason the chart is worth opening.
4. **A batch that has an order carries it as `Depends-on:` trailers.** Ordering that lives only in
   prose is invisible to the board: waves collapse to one, the dependency column renders blank,
   and nothing warns you.

The formula, so every number below is checkable by hand:

```
padded = ceil(stated × 1.2)
start  = 1 + max(end of every dependency)          # or the sprint's begin, whichever is later
end    = start + padded − 1                        # inclusive calendar days
```

---

## 2 · The work units

`M` = model · `E` = effort · `D` = that unit's **own** stated duration · `→` = padded bar.

### Sprint 0 — Scaffolding and the seams · 11 units

The sprint most likely to finish **early**. It carries no screen, no sync and no ambition: the
repo, the pipeline, the persistence floor, the theme system, the harnesses, and the two seams
whose default implementations are *null* by design.

| # | Title | Kind | M | E | D | → | Depends on |
|---|---|---|---|---|---|---|---|
| **Z0-1** | `[INFRA][GO]` Repo skeleton, module layout, Makefile, lint and the cross-compile matrix | chore | sonnet | medium | **2d** | 3d | — |
| **Z0-2** | `[INFRA][CI]` The pipeline — test, lint, build matrix, coverage floor | chore | sonnet | medium | **2d** | 3d | Z0-1 |
| **Z0-3** | `[I18N]` The catalogue, `x/text`, width-aware measurement, and the no-literals lint that fails the build | feature | opus | high | **4d** | 5d | Z0-2 |
| **Z0-4** | `[DATA]` The migration runner — forward-only, idempotent, transactional — and schema v1 | feature | opus | high | **3d** | 4d | Z0-1, **#5** |
| **Z0-5** | `[DATA]` The `Store` seam — pure-Go SQLite, query layer, derived status, tombstones | feature | opus | xhigh | **5d** | 6d | Z0-4, **#6** |
| **Z0-6** | `[THEME]` The theme engine — token contract, both dialects, SUPPLIED/DERIVED/PINNED, the loader | feature | opus | high | **4d** | 5d | Z0-1 |
| **Z0-7** | `[THEME]` The four-state validator — G1–G5, CIEDE2000 + CVD, deterministic selection | feature | opus | xhigh | **3d** | 4d | Z0-6 |
| **Z0-8** | `[THEME]` **The KITT theme** — `kitt.toml`, plus `delorean` and `retro-82`, with attribution | feature | sonnet | high | **2d** | 3d | Z0-7 |
| **Z0-9** | `[QA]` The golden-file TUI harness and the no-network harness | chore | opus | high | **3d** | 4d | Z0-2, Z0-6 |
| **Z0-10** | `[AUDIO]` The audio seam — `NullAudio` as the default build, the `audio` build tag, `ZERADO_NO_AUDIO` | feature | opus | high | **2d** | 3d | Z0-1, **#6** |
| **Z0-11** | `[IMAGES]` The image seam — capability detection failing closed, `NullImages`, Kitty and iTerm2 placement | feature | opus | high | **3d** | 4d | Z0-1, **#6** |

### Sprint 1 — The shell, the ledger and the subsystems · 12 units

Everything a screen needs before a screen can exist, plus the one screen that proves it.

| # | Title | Kind | M | E | D | → | Depends on |
|---|---|---|---|---|---|---|---|
| **Z1-1** | `[TUI]` The route stack, the single overlay slot, focus, and one key registry | feature | opus | high | **4d** | 5d | Z0-3, Z0-6 |
| **Z1-2** | `[TUI]` The frame — header band, footer, status bar, spacing tokens, five tiers, the refusal floor | feature | opus | high | **5d** | 6d | Z0-6, Z0-3 |
| **Z1-3** | `[TUI]` The component library — state chip, ledger row, degrade banner, empty, error, progress, the scanner | feature | opus | xhigh | **6d** | 8d | Z1-2, Z0-7 |
| **Z1-4** | `[TUI]` `Z-11 Fatal error` — frameless, silent, and dependent on nothing | feature | sonnet | medium | **1.5d** | 2d | Z0-3 |
| **Z1-5** | `[TUI]` `Z-04 Library` — the ledger, the pinned summary, cursor-following scroll | feature | opus | xhigh | **5d** | 6d | Z1-3, Z0-5 |
| **Z1-6** | `[PROVIDER]` The provider seam — `Provider`/`Syncer` segregated, `Capabilities`, the registry, `physical` | feature | opus | high | **4d** | 5d | Z0-5 |
| **Z1-7** | `[SEC]` The `Vault` — OS keychain, `0600` file fallback, and the honest line about which is in use | feature | sonnet | high | **2.5d** | 3d | Z0-5 |
| **Z1-8** | `[PROVIDER]` The Steam syncer — streaming, tombstoning, batch-transactional upsert, `item_uid` | feature | opus | xhigh | **5d** | 6d | Z1-6 |
| **Z1-9** | `[IMAGES]` The cover cache — XDG cache, bounded, evicting, `Warm` off the render path | feature | opus | high | **3d** | 4d | Z0-11, Z0-5 |
| **Z1-10** | `[AUDIO]` The radio player behind `-tags audio` — two channels, `m`, stations as data | feature | opus | xhigh | **5d** | 6d | Z0-10, Z1-1, Z1-2 |
| **Z1-11** | `[AUDIO]` Interface FX and the fourth co-render channel | feature | sonnet | high | **2.5d** | 3d | Z1-10 |
| **Z1-12** | `[QA]` The offline-contract harness — the no-network and stale-data property runs | chore | opus | high | **3d** | 4d | Z0-9, Z1-3 |

### Sprint 2 — The screens · 9 units

Every remaining Phase 1 surface. All of them are routes or modes hanging off `Z-04`, which is why
they all sit behind `Z1-5`: `Z-04` is the root and can never be popped (ADR-0001 D3), so every
route's `Esc` behaviour is only testable once the root is real.

| # | Title | Kind | M | E | D | → | Depends on |
|---|---|---|---|---|---|---|---|
| **Z2-1** | `[TUI]` `Z-02 Connect a store` — the form generated from `Capabilities().Credentials` | feature | opus | high | **3d** | 4d | Z1-5, Z1-6, Z1-7 |
| **Z2-2** | `[TUI]` `Z-03 Sync` — honest progress, the scanner, `PARTIAL`, and the readout | feature | opus | high | **4d** | 5d | Z1-5, Z1-8 |
| **Z2-3** | `[TUI]` `Z-08 Add a game by hand` — a physical copy as a first-class row | feature | sonnet | high | **2.5d** | 3d | Z1-5, Z1-6 |
| **Z2-4** | `[TUI]` `Z-05 Game detail` — and the not-fetched vs known-empty distinction | feature | opus | high | **4.5d** | 6d | Z1-5 |
| **Z2-5** | `[TUI]` `Z-06 Set status` — the overlay, the route at Tiny, and the four-state machine | feature | opus | high | **3d** | 4d | Z2-4 |
| **Z2-6** | `[TUI]` `Z-07 Filter and search` — the mode, the dynamic bar, the two-step `Esc` | feature | opus | high | **4d** | 5d | Z1-5 |
| **Z2-7** | `[TUI]` `Z-15 Cover deck` — the deck, and the once-only dismissible note | feature | opus | xhigh | **5d** | 6d | Z1-5, Z1-9 |
| **Z2-8** | `[UIKIT]` The filterable key-map table component — and Zerado's help as its first consumer | feature | opus | xhigh | **7d** | 9d | Z1-5, Z1-1, **#6**, **#8** |
| **Z2-9** | `[TUI]` `Z-01 First run` — the three doors and the one quiet audio line | feature | sonnet | medium | **2d** | 3d | Z2-1, Z2-3 |

### Sprint 3 — Settings, acceptance and the release · 6 units

| # | Title | Kind | M | E | D | → | Depends on |
|---|---|---|---|---|---|---|---|
| **Z3-1** | `[TUI]` `Z-09 Settings` — every dial visible, the Audio section, the honest availability lines | feature | opus | xhigh | **5d** | 6d | Z1-11, Z1-7, Z0-11, Z2-7 |
| **Z3-2** | `[QA]` The colour acceptance — the 17-line theme bar, `NO_COLOR`, the 16-colour floor, `ZERADO_ASCII` | chore | opus | high | **3d** | 4d | Z0-8, Z2-7 |
| **Z3-3** | `[QA]` The offline, responsive and refusal-floor acceptance runs | chore | opus | high | **3d** | 4d | Z1-12, Z2-9 |
| **Z3-4** | `[I18N]` The one-file proof — add `pt-BR` and change no code | chore | haiku | low | **2d** | 3d | Z0-3, Z2-9 |
| **Z3-5** | `[RELEASE]` The cross-compiled matrix, both build tags, checksums and the install path | chore | sonnet | high | **3d** | 4d | Z3-1, Z3-2, Z3-3, Z3-4 |
| **Z3-6** | `[DOC]` The Phase 1 sign-off pack — rendered artifacts for the default and one theme per tier | documentation | opus | high | **2d** | 3d | Z3-5 |

---

## 3 · The DAG, drawn

Nodes are work units. Edges are `Depends-on`. **The critical path is red.**

```mermaid
flowchart LR
  classDef crit fill:#2A0A12,stroke:#FF2E2E,stroke-width:2px,color:#FFE9EE
  classDef s0 fill:#08131F,stroke:#19E0FF,color:#DFF8FF
  classDef s1 fill:#0B1520,stroke:#FFB000,color:#FFF0D0
  classDef s2 fill:#120B20,stroke:#C77DFF,color:#F1E4FF
  classDef s3 fill:#0A1A10,stroke:#7CE38B,color:#E2FBE7
  classDef up  fill:#1A1A1A,stroke:#A5A29B,color:#EFEDE8

  T5["#5 ERD + schema"]:::up
  T6["#6 interface contracts"]:::up
  T8["#8 key-map visual spec"]:::up

  Z01["Z0-1 repo skeleton<br/>2d → 3d"]:::crit
  Z02["Z0-2 CI pipeline<br/>2d → 3d"]:::crit
  Z03["Z0-3 i18n spine<br/>4d → 5d"]:::crit
  Z04["Z0-4 migration runner<br/>3d → 4d"]:::s0
  Z05["Z0-5 Store seam<br/>5d → 6d"]:::s0
  Z06["Z0-6 theme engine<br/>4d → 5d"]:::s0
  Z07["Z0-7 four-state validator<br/>3d → 4d"]:::s0
  Z08["Z0-8 KITT theme<br/>2d → 3d"]:::s0
  Z09["Z0-9 golden harness<br/>3d → 4d"]:::s0
  Z010["Z0-10 audio seam<br/>2d → 3d"]:::s0
  Z011["Z0-11 image seam<br/>3d → 4d"]:::s0

  Z11["Z1-1 route stack + keys<br/>4d → 5d"]:::s1
  Z12["Z1-2 frame + tiers<br/>5d → 6d"]:::crit
  Z13["Z1-3 component library<br/>6d → 8d"]:::crit
  Z14["Z1-4 Z-11 fatal<br/>1.5d → 2d"]:::s1
  Z15["Z1-5 Z-04 Library<br/>5d → 6d"]:::crit
  Z16["Z1-6 provider seam<br/>4d → 5d"]:::s1
  Z17["Z1-7 Vault<br/>2.5d → 3d"]:::s1
  Z18["Z1-8 Steam syncer<br/>5d → 6d"]:::s1
  Z19["Z1-9 cover cache<br/>3d → 4d"]:::s1
  Z110["Z1-10 radio player<br/>5d → 6d"]:::s1
  Z111["Z1-11 FX + 4th channel<br/>2.5d → 3d"]:::s1
  Z112["Z1-12 offline harness<br/>3d → 4d"]:::s1

  Z21["Z2-1 Z-02 Connect<br/>3d → 4d"]:::s2
  Z22["Z2-2 Z-03 Sync<br/>4d → 5d"]:::s2
  Z23["Z2-3 Z-08 By hand<br/>2.5d → 3d"]:::s2
  Z24["Z2-4 Z-05 Detail<br/>4.5d → 6d"]:::s2
  Z25["Z2-5 Z-06 Set status<br/>3d → 4d"]:::s2
  Z26["Z2-6 Z-07 Filter<br/>4d → 5d"]:::s2
  Z27["Z2-7 Z-15 Cover deck<br/>5d → 6d"]:::s2
  Z28["Z2-8 key-map component<br/>7d → 9d"]:::crit
  Z29["Z2-9 Z-01 First run<br/>2d → 3d"]:::s2

  Z31["Z3-1 Z-09 Settings<br/>5d → 6d"]:::crit
  Z32["Z3-2 colour acceptance<br/>3d → 4d"]:::s3
  Z33["Z3-3 offline acceptance<br/>3d → 4d"]:::s3
  Z34["Z3-4 pt-BR one-file proof<br/>2d → 3d"]:::s3
  Z35["Z3-5 release matrix<br/>3d → 4d"]:::crit
  Z36["Z3-6 sign-off pack<br/>2d → 3d"]:::crit

  Z01 --> Z02 --> Z03
  Z01 --> Z04 & Z06 & Z010 & Z011
  T5 --> Z04
  T6 --> Z05 & Z010 & Z011
  Z04 --> Z05
  Z06 --> Z07 --> Z08
  Z02 --> Z09
  Z06 --> Z09
  Z03 --> Z11 & Z12 & Z14
  Z06 --> Z11 & Z12
  Z12 --> Z13
  Z07 --> Z13
  Z13 --> Z15 & Z112
  Z05 --> Z15 & Z16 & Z17 & Z19
  Z011 --> Z19
  Z16 --> Z18
  Z010 --> Z110
  Z11 --> Z110
  Z12 --> Z110
  Z110 --> Z111
  Z09 --> Z112
  Z15 --> Z21 & Z22 & Z23 & Z24 & Z26 & Z27 & Z28
  Z11 --> Z28
  T6 --> Z28
  T8 --> Z28
  Z16 --> Z21 & Z23
  Z17 --> Z21
  Z18 --> Z22
  Z19 --> Z27
  Z24 --> Z25
  Z21 --> Z29
  Z23 --> Z29
  Z111 --> Z31
  Z17 --> Z31
  Z011 --> Z31
  Z27 --> Z31 & Z32
  Z08 --> Z32
  Z112 --> Z33
  Z29 --> Z33 & Z34
  Z03 --> Z34
  Z31 --> Z35
  Z32 --> Z35
  Z33 --> Z35
  Z34 --> Z35
  Z35 --> Z36
```

### 3.1 · The critical path

**`Z0-1 → Z0-2 → Z0-3 → Z1-2 → Z1-3 → Z1-5 → Z2-8 → Z3-1 → Z3-5 → Z3-6`**

Fifty-three padded days of work in a fifty-six-day plan — so the plan carries three days of
end-to-end float, all of it at the sprint seams. What it says out loud:

- **The i18n spine is on the critical path from day three.** That is the price of "no user-facing
  string literal in code" and it was chosen deliberately (ADR-0001 D9 — *"we will not make the same
  mistake we did on FlowForge"*). Nothing that renders text can start until it lands.
- **The component library is the widest bar in the plan (8 days) and everything visual queues
  behind it.** All nine Sprint 2 units are its descendants, through `Z1-5`.
- **The KITT theme is NOT on the critical path.** `Z0-8` has float: it finishes 2026-09-08 and
  nothing needs it until the colour acceptance sweep on 2026-10-07 — **twenty-eight days of slack.**
  That is the correct shape for it. It is required, it is sized, and it is not allowed to hold up a
  screen.
- **The three subsystem tracks run genuinely parallel** — persistence (`Z0-4→Z0-5→Z1-6→Z1-8`),
  audio (`Z0-10→Z1-10→Z1-11`) and images (`Z0-11→Z1-9`) each have their own float against the UI
  chain, which is exactly why the chart is worth opening.

---

## 4 · The derived schedule

### Sprint 0 — Scaffolding and the seams · **2026-08-25 → 2026-09-08**

| Unit | stated | +20% | start | end | wave |
|---|---|---|---|---|---|
| Z0-1 | 2d | **3d** | 2026-08-25 | 2026-08-27 | 0 |
| Z0-2 | 2d | **3d** | 2026-08-28 | 2026-08-30 | 1 |
| Z0-4 | 3d | **4d** | 2026-08-28 | 2026-08-31 | 1 |
| Z0-6 | 4d | **5d** | 2026-08-28 | 2026-09-01 | 1 |
| Z0-10 | 2d | **3d** | 2026-08-28 | 2026-08-30 | 1 |
| Z0-11 | 3d | **4d** | 2026-08-28 | 2026-08-31 | 1 |
| Z0-3 | 4d | **5d** | 2026-08-31 | 2026-09-04 | 2 |
| Z0-5 | 5d | **6d** | 2026-09-01 | 2026-09-06 | 2 |
| Z0-7 | 3d | **4d** | 2026-09-02 | 2026-09-05 | 2 |
| Z0-9 | 3d | **4d** | 2026-09-02 | 2026-09-05 | 2 |
| Z0-8 | 2d | **3d** | 2026-09-06 | 2026-09-08 | 3 |

**Five bars start on 2026-08-28 and none of them ends on the same day.** Eleven units, four waves,
nine distinct spans.

### Sprint 1 — The shell, the ledger and the subsystems · **2026-09-07 → 2026-09-26**

| Unit | stated | +20% | start | end | wave |
|---|---|---|---|---|---|
| Z1-1 | 4d | **5d** | 2026-09-07 | 2026-09-11 | 0 |
| Z1-2 | 5d | **6d** | 2026-09-07 | 2026-09-12 | 0 |
| Z1-4 | 1.5d | **2d** | 2026-09-07 | 2026-09-08 | 0 |
| Z1-6 | 4d | **5d** | 2026-09-07 | 2026-09-11 | 0 |
| Z1-7 | 2.5d | **3d** | 2026-09-07 | 2026-09-09 | 0 |
| Z1-9 | 3d | **4d** | 2026-09-07 | 2026-09-10 | 0 |
| Z1-8 | 5d | **6d** | 2026-09-12 | 2026-09-17 | 1 |
| Z1-3 | 6d | **8d** | 2026-09-13 | 2026-09-20 | 1 |
| Z1-10 | 5d | **6d** | 2026-09-13 | 2026-09-18 | 1 |
| Z1-11 | 2.5d | **3d** | 2026-09-19 | 2026-09-21 | 2 |
| Z1-5 | 5d | **6d** | 2026-09-21 | 2026-09-26 | 2 |
| Z1-12 | 3d | **4d** | 2026-09-21 | 2026-09-24 | 2 |

### Sprint 2 — The screens · **2026-09-27 → 2026-10-06**

| Unit | stated | +20% | start | end | wave |
|---|---|---|---|---|---|
| Z2-3 | 2.5d | **3d** | 2026-09-27 | 2026-09-29 | 0 |
| Z2-1 | 3d | **4d** | 2026-09-27 | 2026-09-30 | 0 |
| Z2-2 | 4d | **5d** | 2026-09-27 | 2026-10-01 | 0 |
| Z2-6 | 4d | **5d** | 2026-09-27 | 2026-10-01 | 0 |
| Z2-4 | 4.5d | **6d** | 2026-09-27 | 2026-10-02 | 0 |
| Z2-7 | 5d | **6d** | 2026-09-27 | 2026-10-02 | 0 |
| Z2-8 | 7d | **9d** | 2026-09-27 | 2026-10-05 | 0 |
| Z2-9 | 2d | **3d** | 2026-10-01 | 2026-10-03 | 1 |
| Z2-5 | 3d | **4d** | 2026-10-03 | 2026-10-06 | 1 |

### Sprint 3 — Settings, acceptance and the release · **2026-10-07 → 2026-10-19**

| Unit | stated | +20% | start | end | wave |
|---|---|---|---|---|---|
| Z3-4 | 2d | **3d** | 2026-10-07 | 2026-10-09 | 0 |
| Z3-2 | 3d | **4d** | 2026-10-07 | 2026-10-10 | 0 |
| Z3-3 | 3d | **4d** | 2026-10-07 | 2026-10-10 | 0 |
| Z3-1 | 5d | **6d** | 2026-10-07 | 2026-10-12 | 0 |
| Z3-5 | 3d | **4d** | 2026-10-13 | 2026-10-16 | 1 |
| Z3-6 | 2d | **3d** | 2026-10-17 | 2026-10-19 | 2 |

**Phase 1 lands 2026-10-19.** Fifty-six calendar days from today, with the +20% already inside
every one of them.

---

## 5 · The Gantt

```mermaid
gantt
    title Zerado Phase 1 — derived from the DAG, +20% per ticket
    dateFormat YYYY-MM-DD
    axisFormat %d %b

    section Sprint 0
    Z0-1 repo skeleton          :crit, z01, 2026-08-25, 3d
    Z0-2 CI pipeline            :crit, z02, 2026-08-28, 3d
    Z0-4 migration runner       :z04, 2026-08-28, 4d
    Z0-6 theme engine           :z06, 2026-08-28, 5d
    Z0-10 audio seam            :z010, 2026-08-28, 3d
    Z0-11 image seam            :z011, 2026-08-28, 4d
    Z0-3 i18n spine             :crit, z03, 2026-08-31, 5d
    Z0-5 Store seam             :z05, 2026-09-01, 6d
    Z0-7 four-state validator   :z07, 2026-09-02, 4d
    Z0-9 golden harness         :z09, 2026-09-02, 4d
    Z0-8 KITT theme             :milestone, z08, 2026-09-06, 3d

    section Sprint 1
    Z1-1 route stack + keys     :z11, 2026-09-07, 5d
    Z1-2 frame + tiers          :crit, z12, 2026-09-07, 6d
    Z1-4 Z-11 fatal error       :z14, 2026-09-07, 2d
    Z1-6 provider seam          :z16, 2026-09-07, 5d
    Z1-7 Vault                  :z17, 2026-09-07, 3d
    Z1-9 cover cache            :z19, 2026-09-07, 4d
    Z1-8 Steam syncer           :z18, 2026-09-12, 6d
    Z1-3 component library      :crit, z13, 2026-09-13, 8d
    Z1-10 radio player          :z110, 2026-09-13, 6d
    Z1-11 FX + 4th channel      :z111, 2026-09-19, 3d
    Z1-5 Z-04 Library           :crit, z15, 2026-09-21, 6d
    Z1-12 offline harness       :z112, 2026-09-21, 4d

    section Sprint 2
    Z2-3 Z-08 add by hand       :z23, 2026-09-27, 3d
    Z2-1 Z-02 connect a store   :z21, 2026-09-27, 4d
    Z2-2 Z-03 sync              :z22, 2026-09-27, 5d
    Z2-6 Z-07 filter and search :z26, 2026-09-27, 5d
    Z2-4 Z-05 game detail       :z24, 2026-09-27, 6d
    Z2-7 Z-15 cover deck        :z27, 2026-09-27, 6d
    Z2-8 key-map component      :crit, z28, 2026-09-27, 9d
    Z2-9 Z-01 first run         :z29, 2026-10-01, 3d
    Z2-5 Z-06 set status        :z25, 2026-10-03, 4d

    section Sprint 3
    Z3-4 pt-BR one-file proof   :z34, 2026-10-07, 3d
    Z3-2 colour acceptance      :z32, 2026-10-07, 4d
    Z3-3 offline acceptance     :z33, 2026-10-07, 4d
    Z3-1 Z-09 settings          :crit, z31, 2026-10-07, 6d
    Z3-5 release matrix         :crit, z35, 2026-10-13, 4d
    Z3-6 sign-off pack          :crit, z36, 2026-10-17, 3d
```

---

## 6 · Wave layering

Waves are computed from the `Depends-on:` edges, not declared. A wave is what can run
**simultaneously**, so the wave count is the sprint's real depth.

| Sprint | Units | Waves | Widest wave | What the width means |
|---|---|---|---|---|
| **0** | 11 | 4 | **5** | Five workers can run on 2026-08-28 |
| **1** | 12 | 3 | **6** | Six workers on 2026-09-07 — the subsystem tracks fan out here |
| **2** | 9 | 2 | **7** | Seven screens in parallel; this is the sprint the fleet exists for |
| **3** | 6 | 3 | **4** | The acceptance sweeps run together, then converge on the release |

---

## 7 · The model and effort budget

Per sprint, so the cost shape is visible before it runs.

| Sprint | units | opus | sonnet | haiku | max | xhigh | high | medium | low |
|---|---|---|---|---|---|---|---|---|---|
| **0** | 11 | 8 | 3 | 0 | 0 | 2 | 7 | 2 | 0 |
| **1** | 12 | 9 | 3 | 0 | 0 | 4 | 7 | 1 | 0 |
| **2** | 9 | 7 | 2 | 0 | 0 | 2 | 6 | 1 | 0 |
| **3** | 6 | 4 | 1 | 1 | 0 | 1 | 4 | 0 | 1 |
| **total** | **38** | **28** | **9** | **1** | **0** | **9** | **24** | **4** | **1** |

**The honest reading: this is an opus-heavy plan, and it is not padding.** Twenty-eight of
thirty-eight units are opus because this is greenfield Go against 11 572 lines of binding screen
spec, a colour-science gate, a streaming provider seam and an i18n discipline that fails the build.
The units that are genuinely mechanical got a cheaper model and say so:

- **`Z3-4`** is the only **haiku** unit — adding a `pt-BR` catalogue file and asserting no code
  changed is file-shaped work with a scripted verdict.
- **`Z0-1`, `Z0-2`, `Z0-8`, `Z1-4`, `Z1-7`, `Z1-11`, `Z2-3`, `Z2-9`, `Z3-5`** are the nine **sonnet** units — bounded work
  against a spec that already answers the hard questions.
- **No unit is `max`.** `max` is for a problem whose shape is unknown; every unit here has a
  ratified document telling it what to build.

---

## 8 · The KITT theme — the founder's amendment, sized

**`Z0-8` exists because of a real gap, not a preference.**
[`../design/05-theme-system.md`](../design/05-theme-system.md) names `delorean` **eight times** and
**KITT / Knight Rider zero times** — yet the scanner sweep the design system calls the brand's one
signature motion *is* KITT's oscillating mirror
([`../design/01-design-system.md`](../design/01-design-system.md) §9), the landing page sells that
register, and `../design/00-design-brief.md` §51 states the bar as *"the DeLorean and KITT."*
FlowForge already ships `delorean.toml` among the 35 omarchy themes; **`kitt.toml` is its sibling
and does not exist yet** — verified at source, 2026-08-25.

It is **not** folded into "the theme system". It is its own unit, with its own duration, its own
acceptance and its own row on the chart:

- `kitt.toml` authored against the token contract, in the semantic dialect.
- It must **pass the four-state gate** (`Z0-7`) on its own measurements — G1–G5, CIEDE2000 with
  measured CVD separation, every text pair clearing 4.5:1 on its own ground. A KITT theme that
  cannot express four distinguishable states **does not ship broken; it does not ship.**
- Its `scanner` role comes from the red band, which is the one colour KITT unambiguously owns.
- `delorean` and `retro-82` are adopted in the same unit, `retro-82` labelled
  `player · abandoned deviates (warm)` per `05-theme-system.md` §4.3.
- An `ATTRIBUTION.md` row with all six §5.2 fields for each, field 6 generated.

**It has twenty-eight days of float** (§3.1). That is deliberate: it is required, and it is not
permitted to block a screen.

---

## 9 · Scope changes folded into this plan

### 9.1 · Audio and images are Phase 1

Both were once later phases. Both are ADR-0001 decisions now (D6, D8) and both are sized here:
`Z0-10`, `Z1-10`, `Z1-11` for audio; `Z0-11`, `Z1-9`, `Z2-7` for images. **No music is bundled** —
radio is streamed, off by default, and `NullAudio` is the *default build*.

### 9.2 · `Z-10 Help and key map` is superseded by a reusable component

Founder direction, 2026-08-25, relayed:

> *"instead of create the Help screen of Zerado, I want a component on the design system, that has
> apis to deal with it, and should come with a text filter on it. Then Zerado will be the first
> software to use that screen… Remember the screen should work on any theme so it will work
> perfectly on Zerado and on FlowForge."*

**The old `Z-10` screen unit is withdrawn.** In its place, one chain of three:

| | Owner | Where |
|---|---|---|
| **1 · the contract** — the Go interface, the data/column/filter contracts, the theme-token contract | `fft-api-designer` | added to ticket **#6** |
| **2 · the visual + interaction spec** — mockup, states, spacing, the filter affordance, empty and no-match, key map | `fft-tui-designer` | ticket **#8** |
| **3 · the BUILD, plus Zerado's help/key-map as its first consumer** | Phase 1 | **`Z2-8`**, below |

**It was re-sized, and here is why.** The withdrawn screen was **2d**. `Z2-8` is **7d** — a 3.5×
increase, and the increase is the point rather than an inflation:

| What changed | Why it costs |
|---|---|
| product-agnostic, not a Zerado screen | no Zerado imports, no product vocabulary in a type or field name; a second product must adopt it **without a fork** |
| theme-agnostic across **35 FlowForge themes** plus Zerado's default | no hardcoded hex, no Zerado-specific token, and a render sweep that proves it |
| an **intrinsic** text filter, not one added later | the founder's reasoning is the estimate's justification: *a key map keeps growing, so it is born needing a filter*. Filtering that is retrofitted into a table changes its column model, its focus model and its empty state — all three |
| a generic column model with a width policy | column arithmetic **must not assume one cell per character** — the same width discipline `Z0-3` establishes, applied to a component someone else will feed |
| focus and key pass-through as a contract | the host owns the key map; the component must not swallow a key the host needed |
| co-render survives `NO_COLOR` and 16 colours | colour **and** glyph **and** label, three channels, any two sufficient |
| a second consumer's acceptance bar | reviewed as a library, not as a screen |

**Explicitly deferred, with no date and no work unit:** where the component eventually lives —
moved into the FlowForge design system, or into a shared private component repo. The founder's word
is *"for later"*. **No Phase 1 unit depends on that decision**, and none may be made to.

### 9.3 · Books are out

One cheap affordance and nothing more: the core table is `item`, not `games`, with an `item_type`
column `CHECK`ed to `'game'` (ADR-0001 D5). That is entirely inside `Z0-4`. **There is no
media-polymorphism work unit**, because there is no media polymorphism.

---

## 10 · What is NOT in Phase 1

Recorded so it is not re-proposed, and so the reason survives.

| Not built | Where it lands |
|---|---|
| Cover art *fetching* — the metadata provider | **Phase 2.** The seam is Phase 1 (`Z0-11`); the source is not decided and cannot be until IGDB answers |
| *Sinopse*, mood tags, mood picker, price history | **Phase 2** |
| The command palette (`Z-17`) | **Phase 2.** `:` and `Ctrl-K` are **reserved and unbound** from Phase 1 — that reservation is inside `Z1-1` |
| Recommendations, budget, watchlist | **Phase 3** |
| Accounts, device sync, community, the phone app | **Phase 4.** The one Phase 1 obligation is `item_uid`, `merged_into` and `status_changed_at` in the schema — inside `Z0-4` |
| **A light theme** | **Blocked, not scheduled.** The light state set FAILS its own four-state gate — measured, `not started × zerado` at ΔE 5.41 against a 10.0 floor. Repairing it is `fft-brand-architect`'s, through the theme-system §10 governance procedure. Until it is repaired Zerado ships **no light theme**, and that is the gate working |
| Sixel | Deferred on evidence — poor colour fidelity on photographic content at deck sizes |
| A splash screen | **Never.** It contradicts *"it starts instantly"* on the first frame the player ever sees |
| Moving the key-map component into FlowForge | **Later, no date** (§9.2) |

---

## 11 · Findings — read these before ratifying

1. **Phase 1 lands 2026-10-19, not sooner.** Fifty-six days. The padding is already inside every
   bar and **it will not be shrunk to fit a date.** If that finish is wrong, the answer is to cut
   scope, and §10 is where the cut list starts.
2. **`Z1-3` (the component library) is the single widest bar and the plan's chokepoint.** All nine Sprint 2 units are its
   descendants. If one unit deserves a second pair of eyes at design
   time, it is this one — a defect there costs eight days plus every screen behind it.
3. **`Z2-8` is a 9-day bar inside a 10-day sprint, and it is the newest unit in the plan.** It is
   also the only unit whose upstream (`#6`'s contract amendment and `#8`'s visual spec) does not
   exist yet. It is the plan's largest single estimation risk and it is flagged rather than
   smoothed. If `#8` slips, `Z2-8` slips, and Sprint 3 slips with it.
4. **Two units depend on tickets that are still in flight.** `Z0-4` needs `#5`, and `Z0-5`,
   `Z0-10`, `Z0-11` need `#6`. Both were sequenced rather than awaited, and neither is on the
   critical path — but a `#5` or `#6` that lands late pushes the persistence and subsystem tracks,
   not the whole plan.
5. **Sprint 0 finishes on 2026-09-08 and is deliberately the lightest sprint in the plan.** No
   screen, no sync, no ambition. If any sprint finishes early, this is the one that should.
6. **Nothing here was estimated in hours.** Where the blueprint did not support a duration, it is
   said rather than invented — and it did support all thirty-eight, because twelve of them have a
   1 000-line spec each.

---

## 12 · How the dates got onto the board

Not by typing. The seams, in order:

```
flowforge ticket new … --depends-on <n> --live          # the DAG, written as it is emitted
flowforge backfill-dates --milestone "<sprint>" --durations-only --duration <n>=3d --apply
flowforge backfill-dates --milestone "<sprint>" --start <derived>            # dry run
flowforge backfill-dates --milestone "<sprint>" --start <derived> --apply
flowforge depends-audit  --milestone "<sprint>"
flowforge sprint-audit   --milestone "<sprint>"
```

Sprint 0's start is the founder's one answer: **2026-08-25**. Sprints 1, 2 and 3 were **not** asked
about — each one's begin is derived from the finish of the last unit its own wave-0 members depend
on, so no unit in any sprint is scheduled before its dependencies are done. That is the second
scheduling question **not** being asked.
