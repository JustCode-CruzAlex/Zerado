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
| **`fft-database`** | That there is one SQLite file; that the core entity is **`item`** — not `games` — carrying an `item_type` `CHECK`ed to `'game'`, **with no typed extensions and no generic progress**; which entities exist and how they relate; what crosses the Phase 4 boundary; that `effective_status` is derived and not stored; that credentials and cover art live **outside** the file | The **physical** model: exact column types and constraints, the index strategy, the migration sequence and its tests, query shapes, `VACUUM`/WAL policy, and whether a derived-status generated column is worth its cost |
| **`fft-api-designer`** | That the provider / metadata / price / persistence / credential / **audio** seams exist; what each is responsible for; that `Provider` and `Syncer` are **segregated**; that capabilities are **per provider**; that sync **streams**; that a screen never touches a provider | The **exact shape**: final signatures, generic parameters, error contracts and sentinel errors, context and cancellation semantics, versioning and stability guarantees, package layout, and which concrete audio library satisfies §5 of [`12-audio.md`](./12-audio.md) |

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
- that `Quote` carries **no affiliate URL** — a plain shop link — and that every quote carries its
  **own age**, mandatory and rendered;
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

1. ~~**Music licensing**~~ — **CLOSED by removing its cause.** Nothing is bundled; the music is
   streamed radio the player chooses. There is nothing to license, attribute or weigh.
2. ~~**IGDB's commercial terms**~~ — **ANSWERED, with one action outstanding.** Affiliate links are
   dropped, so Zerado is cleanly **non-commercial**: free software, donation-supported, zero
   revenue. IGDB's published test is whether the *project generates revenue*, and it now does not.
   **This is a reading of IGDB's published rationale, not a legal opinion** — the remaining action
   is a direct confirmation from IGDB that a donation-funded open-source project qualifies for the
   free tier. **A founder action, not a resolved fact.**
   *(The price-intelligence feature is unaffected: current price, all-time low and the "wait or
   buy" verdict all remain. Only the commission tag went.)*
3. **The light state set FAILS the four-state gate — measured, not suspected.** This was previously
   carried as *"never CVD-verified"*. It has now been verified, and it does not pass:

   | Pair | Simulation | ΔE2000 | Floor |
   |---|---|---|---|
   | `not started × zerado` | protanopia | **5.41** | 10.0 |
   | `zerado × abandoned` | deuteranopia | **8.91** | 10.0 |

   **The cause is specific and it is upstream.** Brand manual §4.5 claims the paper expression is
   *"the same hues carried to ink weight"*, and four of the five carry within 9.5° of hue. But light
   `not started` is rotated **173°** and lands **1.1°** from `#9FB0C6` — *the blue-cast steel that
   brand §4.4 explicitly rejected on the dark side, for collapsing against the cyan.* **The dark
   set's own correction was never carried to paper.**

   This is `fft-brand-architect`'s to repair through the §10 governance procedure.
   [`../design/05-theme-system.md`](../design/05-theme-system.md) §3.4 carries a feasibility proof
   of a minimum-motion repair reaching a 10.83 floor, explicitly **non-binding** — it demonstrates
   the problem is fixable, it does not choose the fix.

   Until it is repaired, **Zerado ships no light theme**, because D7's gate is a gate rather than a
   preference. That is the mechanism working, and it is why the gate was worth building.
4. ~~**Whether `playtime` is generic progress or a typed game fact**~~ — **WITHDRAWN.** The
   divergence existed to keep one derivation across media types; the media types were pruned, so the
   argument died with them. `playtime_minutes` is a plain column.
5. **A radio stream has to fit a published promise, and the founder should confirm the reading.**
   The page says *"The only network traffic is `Zerado` reaching out to the services you've
   connected — Steam, price data — using your own keys."* A station the player turned on and chose
   is defensibly *a service they connected*, and [`12-audio.md`](./12-audio.md) §3 argues exactly
   that. But the sentence names Steam and price data specifically, it is **published**, and a radio
   stream uses **no key of the player's**. The reading is sound and it is not mine to ratify.
   *(Raised by `fft-tui-designer`.)*
6. **Brand manual §9's invariant list is incomplete by two — and the rate is the finding.**
   §9 names five invariants carrying identity across web, terminal and phone. There are now
   **seven**: **audio** (the sound register, opt-in posture, and never-the-sole-carrier rule) and
   **the caption rule** (a caption under a cover is always the real title, never decorative).

   Both gaps arrived by the **same mechanism** — a product decision landing after rev A froze — and
   both arrived **within one week**. That is worth seeing as a rate rather than as two tickets: §9
   is currently drifting faster than it is being maintained, and the next amendment will do it
   again unless the governance loop is tightened.

   The caption rule is the one most at risk, because **its reason does not travel between
   surfaces.** There is no phone that cannot draw an image, so a designer who only ever sees the
   phone finds a caption under art that always renders, concludes it is redundant, and removes it —
   in good faith. `fft-design-architect` recorded the rule in §2 of the bridge (where invariants
   live, so it carries) and its non-transferring reason in §3 (what cannot translate), deliberately
   in different sections.

   `fft-brand-architect` closes both through the §10 governance procedure. *(Still open.)*
7. **The nine underived ANSI-256 indices.** `fft-brand-architect` derives them by nearest-neighbour
   search in CIELAB; four components ship an interim uncoloured rendering until then.
   *(Still open.)*
8. **The published page currently makes a money claim the product has decided not to honour.**
   `content/landing-copy.md` §16 carries an **affiliate disclosure** — *"`Zerado` earns a commission
   when you buy a game through a price link on this page"* — and §11 and §14 promise *"a premium
   account or a donation"*. Affiliate is dropped and premium is dropped. Both statements are live on
   a public page.

   **This is not a documentation lag.** `fft-design-architect` put it better than my first phrasing
   did, so it is recorded in its words:

   > *a stale money claim on a public page is not a documentation lag, it is the page telling
   > visitors the product earns a commission it has decided not to earn.*

   My earlier version of this item said the page *"contradicts the product"*, which is accurate and
   far too comfortable — it reads as a sync chore that can queue behind other work. It cannot. A
   consistency problem can wait; **a published untrue statement about money cannot**, and it is the
   only item on this list that is currently visible to people outside the project.

   **It cannot be fixed inside this bundle.** That copy is ticket #1's surface. It needs whoever
   owns the page, and it needs them before the affiliate decision is announced anywhere else.

   *(One knock-on, already handled: the page is authority #1 in the design brief's precedence rule,
   so a design obeying the amendment contradicted authority #1 as written. The rule now reads that
   authority #1 is the **ratified promise set**, not the page's current bytes, and a founder
   amendment postdating the page governs.)*

9. ~~**Whether `film` and `series` are two media types**~~ — **moot.** Media types are pruned; the
   observation survives only as a note in [`11-media-model.md`](./11-media-model.md)'s Appendix —
   [`11-media-model.md`](./11-media-model.md)'s Appendix. Costs nothing now, costs a migration
   later.
