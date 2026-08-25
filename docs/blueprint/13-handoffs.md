---
title: Zerado — what this spine decides, and what it hands on
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-13
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: project-brief
ticket: "#2"
---

# Handoffs — the boundary of this document

Two specialists work downstream of this spine on their own tickets. **This document is their
brief.** It decides *that* a seam exists, what it is responsible for, and where its boundaries
are. It does not decide the final SQL or the final Go signatures.

Stated explicitly, because a blueprint that quietly over-reaches produces two specifications that
disagree — and the one that gets built is whichever the builder read last.

---

## 1 · The split

| | **This spine decides** | **The specialist decides** |
|---|---|---|
| **`fft-database`** | That there is one SQLite file; that the core entity is a **media item** with typed extensions; which entities exist and how they relate; what crosses the Phase 4 boundary; that `effective_status` is derived and not stored; that credentials and cover art live **outside** the file | The **physical** model: exact column types and constraints, the index strategy, the migration sequence and its tests, query shapes, `VACUUM`/WAL policy, and whether a derived-status generated column is worth its cost |
| **`fft-api-designer`** | That the provider / metadata / price / persistence / credential / **audio** seams exist; what each is responsible for; that `Provider` and `Syncer` are **segregated**; that capabilities are per `(provider, media type)`; that sync **streams**; that a screen never touches a provider | The **exact shape**: final signatures, generic parameters, error contracts and sentinel errors, context and cancellation semantics, versioning and stability guarantees, package layout, and which concrete audio library satisfies §5 of [`12-audio.md`](./12-audio.md) |

---

## 2 · What the Go in this bundle is, and is not

Every Go block in [`06-data-seams.md`](./06-data-seams.md),
[`05-state-machine.md`](./05-state-machine.md), [`11-media-model.md`](./11-media-model.md) and
[`12-audio.md`](./12-audio.md) is **shape, not signature.**

It is there because a seam described only in prose is a seam nobody can be held to — the argument
for `Provider` and `Syncer` being two interfaces cannot be made without showing that `physical`
does not implement the second one. What is load-bearing is **the decision the shape encodes**, not
the identifier spelling, the package, or the parameter order.

`fft-api-designer` should feel free to change every name in this bundle. It should not silently
change:

- that `physical` implements `Provider` and **not** `Syncer`;
- that everything downstream reads **`Capabilities`**, never `ProviderID` or `MediaType`;
- that `Sync` **streams**, so a cancel leaves a valid partial library;
- that every network-derived value carries **its own age** in the same value;
- that `Quote` carries the affiliate URL and the disclosure obligation **in one struct**;
- that `Audio.Cue` **cannot fail and cannot block**;
- that `Store` is the only writer, and screens only read.

Each of those is a decision with a reason attached in the document that states it. Changing one is
fine; changing one **by accident** is the failure this list exists to prevent.

---

## 3 · What the ERD in this bundle is, and is not

[`09-erd.md`](./09-erd.md) and its two drawn sheets are the **conceptual** model: entities, what
each is for, relationships and cardinality, and — critically — **what is deliberately not in the
file**.

Column lists are shown at the fidelity needed to make the *architecture* argument (that
`status_manual` is nullable; that `last_played_at NULL` means *not reported*; that `fetched_at` is
`NOT NULL` so an age can never be missing; that `game_uid` is indexed and **not** unique). They are
**illustrative of intent, not a schema to run**.

`fft-database` owns the physical model and should expect to change types, add constraints, split or
merge tables, and design indexes against real query shapes. The conceptual claims above are the
ones to preserve, and each has its reasoning in [`09-erd.md`](./09-erd.md) or
[`11-media-model.md`](./11-media-model.md).

---

## 4 · Open questions handed on, not answered here

| Question | Owner | Why it is not answered here |
|---|---|---|
| Which audio library keeps the **default** build pure-Go | `fft-api-designer` | Needs a real per-platform check of cgo dependency, not an assumption |
| Whether derived status warrants a generated column or an expression index | `fft-database` | Needs real query shapes and a 400+ row measurement |
| The `game_uid` normalisation function's exact rules | `fft-database` | Needs testing against a real library's titles; the *policy* (a hint, never authoritative) is decided here |
| Sentinel errors and their contracts | `fft-api-designer` | Follows from the final signatures |
| Migration sequencing and rollback posture | `fft-database` | The spine decides forward-only; the sequence is physical |

---

## 5 · Open questions that are the **founder's**, not a specialist's

Collected here so the gate has them in one place.

1. **Music licensing** — [`12-audio.md`](./12-audio.md) §7. Bundled tracks must be DRM-free and
   licensed for **commercial redistribution** from a **public** repository. Recommendation on the
   table: ship Phase 1 with interface FX only and a **user-supplied** music directory, which removes
   the blocker from the critical path.
2. **IGDB's commercial terms** — [`06-data-seams.md`](./06-data-seams.md) §3. Unresolved. The
   architecture is built so the answer can be no.
3. **Nine underived ANSI-256 indices** — [`00-index.md`](./00-index.md). Four components ship an
   interim uncoloured rendering rather than a guessed index. `fft-brand-architect` derives them.
4. **The light-mode state colours have never been CVD-verified** — the dark set's first draft failed
   at ΔE 8.8, so this is not a formality, and Phase 4 meets system light mode.
5. **Whether `playtime` is generic progress or a typed game fact** —
   [`11-media-model.md`](./11-media-model.md) §1. This spine models it generically, against the
   letter of the direction, so the four-state derivation stays one function instead of forking per
   type. The reason is written down; the call is the founder's.
6. **Brand manual §9's invariant list is now incomplete.** §9 names five invariants that carry the
   identity across web, terminal and phone — the four states, amber-common/cyan-earned, the scanner
   as the only signature motion, the voice, and dark by default. **Audio is a sixth**, and the brand
   manual (rev A, 2026-08-24) predates the audio direction by one day. `fft-design-architect`
   recorded the working invariant in the terminal-to-phone bridge rather than editing the brand
   manual, which is correct: a downstream document quietly extending the brand's own invariant list
   is exactly the drift §9 exists to prevent. **`fft-brand-architect` should add it through the
   §10 governance procedure** — token and manual in the same commit.
7. **Whether `film` and `series` are two media types** —
   [`11-media-model.md`](./11-media-model.md) §4, finding F-1. Costs nothing now, costs a migration
   later.
