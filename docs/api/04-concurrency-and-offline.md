---
title: Zerado — context, concurrency and the offline contract in the signatures
discipline: API
doc-no: ZRD-API-04
rev: A
date: 2026-08-25
status: draft — for review
archetype: concept-explainer
ticket: "#6"
---

# Context, concurrency, and staleness

Three properties the ticket asks to be visible **in the signatures** rather than in a convention.

---

## 1 · Context on everything that could ever block

Every method that does I/O takes a `context.Context` as its first parameter — including every read
on the `Store`, which in Phase 1 is microseconds against a local file.

Two reasons, and neither is speculative:

- **A filter re-runs on every keystroke.** `Z-07` filters live as the player types, and a query
  whose successor has already been typed must be abandonable.
- **A Phase 4 implementation of the same interface does real I/O.** An interface that has to grow
  contexts later grows them in every caller at once.

**Two methods deliberately take no context, and the absence is the guarantee:**

| Method | Why no context |
|---|---|
| `images.Cover` | It never fetches and never blocks — it reads what the cache already holds. A function that cannot block has nothing to cancel, and a `ctx` parameter would be an invitation to make it blocking later |
| `audio.Cue` | Fire-and-forget, non-blocking, no error. A cue that cannot play is silence |

Both are the same rule from
[`../blueprint/17-images.md`](../blueprint/17-images.md) §6 and
[`../blueprint/12-audio.md`](../blueprint/12-audio.md) §6: **a missing cover, and a missed sound, are
never worth a dropped frame.**

---

## 2 · Cancellation is a hard requirement, not politeness

`Syncer.Sync` documents that cancelling `ctx` **must** abort in-flight I/O and close the stream
promptly. `Z-03` lets the player press `Esc` mid-sync, `q` quits from anywhere immediately, and the
goroutine budget for the whole product is **zero leaks** — so a provider that ignores cancellation
does not merely delay a screen, it holds the process open at quit.

`TestSyncStreams` cancels after the second of five items and asserts three things: some items
arrived, not all of them did, and `Stream.Err()` reports `KindCancelled` rather than an
unclassified error that would render as the fatal screen.

The fake honours `ctx` **while sleeping as well as between items**, deliberately. A provider that
only checked between items would pass a fast test and hang the product on a slow network — which is
exactly the moment the player reaches for `Esc`.

---

## 3 · The concurrency shape, stated per return type

| Returns | Where | Why |
|---|---|---|
| a **`Stream`** (channel + terminal error + progress) | `Syncer.Sync` | honest running counts; a cancel leaves a valid partial library; bounded memory at any library size |
| a **slice** | `Store.Games`, `Store.Connections` | a bounded page of a local read. A channel here would buy nothing and cost every caller a `range` |
| a **single value** | `metadata.Lookup`, `pricing.Quote` | one question about one game. The Phase 2 enrichment pass is the caller that runs many concurrently, and it owns the worker count — not the interface |
| **nothing** | `images.Warm`, `audio.Cue` | fire-and-forget, off the render path, bounded by a pool the seam owns |

**No interface takes a worker count.** Fan-out is the caller's decision because only the caller knows
whether it is enriching twelve visible rows or all 247, and an interface parameter would push a
rate-limiting policy into a signature that a second provider would need differently.

### 3.1 · Why `Sync` returns a `Stream` and not a bare channel

The spine's shape was `Sync(ctx, c) (<-chan Item, error)`. It carries the failure that is known
before anything arrives and has **nowhere to put the failure that happens after** — which is
`PARTIAL`, one of `Z-03`'s four terminal states.

`Stream` adds exactly two things to the channel:

- **`Err()`**, valid after the channel closes. The `sql.Rows` pattern.
- **`Progress()`**, a race-free snapshot.

`Progress` is what `Z-03` §3.1 needs and a channel cannot supply: the **denominator**. A channel of
items cannot say "247 are coming", so without it the screen draws an indeterminate scanner forever
and never earns its determinate bar.

`Progress()` is documented as safe to call from another goroutine at any time. That is a hard
requirement, not a nicety — the caller is a Bubble Tea render loop reading it once per frame while
the provider's goroutine writes. `TestProgressIsReadableWhileTheSyncRuns` does exactly that and the
suite passes under `-race`.

Two details in `Progress` that a screen would otherwise get wrong:

- **`TotalKnown` is a separate field from `Total`,** because zero is a real total — it is the
  private-profile case — and a screen reading `0` as "unknown" would scan forever for a sync that
  has truthfully finished. `TestAZeroTotalIsNotUnknown`.
- **`LastAt` is provided; "stalled" is not.** Whether a ten-second pause is alarming is a screen's
  judgement about a player's patience, not a provider's judgement about a socket.
  `TestStalledIsNotWaiting` pins the distinction that a sync which has delivered nothing yet is
  *waiting*, which already has its own component.

---

## 4 · The offline and staleness contract, in the signatures

### 4.1 · There is no freshness parameter on any read

The obvious way to express staleness in an API is a policy argument — `CachedOnly`, `PreferCached`,
`ForceFresh`. **The `Store` interface has none.**

A read cannot reach the network under any policy. What a caller gets is what is local, stamped with
when it was observed; the only way to make it fresher is an explicit sync, which is an action with a
screen and a key behind it.

That absence **is** the offline contract expressed in the signatures. A `ForceFresh` option would be
a hole through which a render path could start a network request, and §7.3's grep rule would then be
guarding a door with a window next to it.

### 4.2 · `aged.Value[T]` makes the age impossible to drop

> Any value that came from the network is rendered with its age. Always.

[`../blueprint/07-offline-contract.md`](../blueprint/07-offline-contract.md) §4 names this as the
rule most likely to be lost during a build, because dropping the age always makes the layout tidier.
It is also the rule the product's credibility rests on: Zerado's value proposition is telling a
player **not** to buy something, and a stale price presented as current is the product giving
financial advice from memory.

Every network-derived value is returned inside `aged.Value[T]`. **There is no way to reach the
payload without having the age in hand** — a `FetchedAt` field can be ignored by a renderer; a
wrapper cannot.

`aged.Map` exists to close the one loophole: transforming a value and re-stamping it as new, in code
that looks entirely reasonable. `TestMapKeepsTheObservationTime`.

`Classify` has three bands because the design system has three behaviours — nothing, the banner as
chrome, and amber past **ninety days** (§4.1). `aged.WarnAfter` is a constant of the contract rather
than a setting, because a threshold a caller can lower is a threshold that gets lowered to make a
banner go away.

A value stamped in the **future** — a source clock ahead of this machine — clamps to zero. "In 3
hours" next to a price would be the product reporting somebody else's NTP problem as a feature.

### 4.3 · Three facts about playtime, not two

`Z-05` renders three different things where a nullable column and a zero would give two:

| Screen | Fact | Mechanism |
|---|---|---|
| `not tracked` | the provider cannot ever know | `Capabilities.Playtime == false` |
| `—` | not fetched yet | the pointer is `nil` |
| `0h` · `never played` | the provider answered, and the answer was nothing | the pointer is set, to zero |

The capability answers the first; the pointer answers the second and third. Collapsing either
produces a screen that tells a player their cartridge has been played for zero hours — a claim
nothing could have made. `TestGameDistinguishesThreeFactsAboutPlaytime`.
