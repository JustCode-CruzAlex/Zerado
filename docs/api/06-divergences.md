---
title: Zerado — divergences from the spine's shapes
discipline: API
doc-no: ZRD-API-06
rev: A
date: 2026-08-25
status: draft — for review
archetype: adr-detail
ticket: "#6"
---

# Divergences

[`../blueprint/13-handoffs.md`](../blueprint/13-handoffs.md) §2 says this deliverable *"should feel
free to change every name in this bundle"* and lists seven decisions it must not change **by
accident**. This document is the accident-prevention mechanism: every departure, named, with the
decision it preserves and the reason it departed.

**First, the seven that are unchanged**, since the list exists to be checked rather than trusted:

| Load-bearing decision | Where it lives now |
|---|---|
| `physical` implements `Provider` and **not** `Syncer` | `providertest.Manual`, asserted by `TestManualIsNotASyncer` |
| everything downstream reads **`Capabilities`**, never an ID | `provider.Capabilities`, asserted by `provider.Check` |
| `Sync` **streams**, so a cancel leaves a valid partial library | `provider.Stream`, asserted by `TestSyncStreams` |
| every network-derived value carries **its own age** | `aged.Value[T]` — strengthened, see §2 |
| `Quote` carries **no affiliate URL**, and every quote carries its age | `pricing.Quote` — there is no field that could hold one |
| `Audio.Cue` **cannot fail and cannot block** | `audio.Audio.Cue` — no error, no context |
| `Store` is the only writer; screens only read | `store.Store`, asserted by `internal/arch` |

---

## 1 · Media-type polymorphism: the ticket asks for a seam, the ADR forbids one

**The ticket's item 9** asks for *"the seams that must be generic so books, and possibly films and
TV, plug in without reshaping anything."*

**ADR-0001 D5 pruned that**, by founder direction on 2026-08-25: *"At this point don't even think on
books and other media types. What I would like is let that door open."*
[`../blueprint/06-data-seams.md`](../blueprint/06-data-seams.md) §7 lists a media-type abstraction as
explicitly **not a seam** — *"an interface parameterised on a type that has one value is machinery
without a purpose"* — and [`../blueprint/11-media-model.md`](../blueprint/11-media-model.md) records
that the speculation had begun to shape Phase 1 tables and Phase 1 states, which is the cost that
makes speculative generality expensive rather than merely unused.

**Resolution: the ratified ADR wins, and the disagreement is surfaced rather than silently
resolved.** `library.Game` has no type parameter and no media-type field.

What answers the ticket's *intent* — that a second media type must not require reshaping anything —
is that the affordance lives in the schema and costs nothing above it: the entity is `item` carrying
an `item_type` constrained to `'game'`. Adding a second type is one `CHECK` constraint. **No
interface changes**, because nothing above the store reads a media type today and nothing would need
to.

*Escalated as a finding, per the brief's instruction to escalate rather than reopen.*

## 2 · The age moves from a field into a wrapper

**Spine shape:** `Metadata.FetchedAt`, `Quote.ObservedAt`.
**Here:** `aged.Value[Metadata]`, `aged.Value[Quote]`.

The load-bearing decision — *every network-derived value carries its own age in the same value* — is
**preserved and strengthened**. A field can be ignored by a renderer; a wrapper cannot, because there
is no way to reach the payload without having the age in hand.

The cost is one `.V` at each call site. The moment it pays for itself is the only moment it is ever
tested: a developer at 1am adding a column to a row, for whom dropping the age makes the layout
tidier.

*(This is a generic type, and it is worth saying why it does not contradict §1: the prohibition is on
parameterising a **seam** over a media type that has one value. `aged.Value[T]` is a utility over
any payload, with several instantiations already.)*

## 3 · `Sync` returns a `Stream`, not a bare channel

**Spine shape:** `Sync(ctx, c) (<-chan Item, error)`.

The `error` return carries a failure known *before* anything arrives and has nowhere to put one that
happens *after* — which is `PARTIAL`, one of `Z-03`'s four terminal states, so it cannot be
inexpressible. `Stream` adds `Err()` after close and `Progress()`, the latter being the denominator
`Z-03` §3.1 needs to draw a determinate bar at all. Reasoned in full in
[`04-concurrency-and-offline.md`](./04-concurrency-and-offline.md) §3.1.

## 4 · `Enterer` is a new interface the spine did not name

The spine has `Provider` and `Syncer`, and describes hand entry as *"`Z-08` writing an `Item`
directly"*.

Written that way, two rules that make a hand-entered row a real row would live in a **screen**: that
its `provider_ref` is a UUID Zerado mints rather than something the player has to find, and that its
optional fields respect the provider's own capabilities. `Enterer` puts them in the provider, which
is where they belong and where a second hand-enterable source inherits them.

It is a **capability, not a provider kind**: a future GOG that syncs *and* lets a player add a game
its API missed implements both, and `Z-08` gains a source picker rather than a special case.

## 5 · `Vault.Backing()` returns a struct, not a string

**Spine shape:** `Backing() string` — `"keychain"` or `"file"`.

`Z-09` §10's copy is *"In the macOS Keychain"* · *"In the GNOME keyring"* · *"In Windows Credential
Manager"* — three values a two-valued string cannot distinguish. The alternative is a screen
switching on `runtime.GOOS`, which is right by coincidence and **wrong inside a container**.

`vault.Backing{Kind, NameKey, Path}` lets the vault report what it actually opened. `Path` is there
because Settings shows the file's location so the player can find, inspect and delete it.

## 6 · `Store` gains `StartRun`/`FinishRun`, and `ReconcileAbsence` takes a run

**Spine shape:** `RecordSyncRun(ctx, r SyncRun)`.

A single record-at-the-end call cannot represent a sync that was **killed** — the ERD distinguishes
that from *cancelled* by `finished_at` being null, and `Z-11`-adjacent honesty depends on it. So the
run is opened before the first item is written.

`ReconcileAbsence` takes a `RunID` rather than a `ProviderID`, and that is **the guard rather than a
style choice**. Only a run whose status is `ok` may tombstone anything (§2.4). Passing a provider id
would have made the illegal call *spellable* — and the illegal call deletes the evidence that a
player finished a game. `TestOnlyAnOkRunMayTombstone` runs it for `partial`, `failed` and
`cancelled` and asserts the refusal and that nothing was marked.

## 7 · `Capabilities.Playtime`, not `Capabilities.Progress`

`Capabilities.Progress` is the **withdrawn** generic-progress model from revision A: ADR-0001 D5
states *"the generic-progress divergence is withdrawn"* and the ERD makes `playtime_minutes` a plain
column. The blueprint spells it `Playtime`. **`Playtime` wins**, being the ratified one.

**Four screen specs still carry the withdrawn name, at five sites.** Enumerated rather than
summarised, so the follow-up catches all of them:

| Spec | Line |
|---|---|
| `../design/screens/Z-04-library.md` | 1027 |
| `../design/screens/Z-05-game-detail.md` | 399 |
| `../design/screens/Z-06-set-status.md` | 366 |
| `../design/screens/Z-08-add-a-game-by-hand.md` | 75, 375 |

Recorded here rather than edited from this ticket, which does not own those specs.

*(An earlier revision of this section recorded **two**. The review at `c4c8d95` grepped and found
four, which left two stale specs unrecorded — a completeness defect in the one document whose whole
job is completeness. The list above is the grep output, not a recollection.)*

## 8 · `Item.PlaytimeMinutes` is a pointer

The spine's `Item` already made optional fields pointers, and this keeps that. It is listed because
one consequence is a rule rather than a type: **a provider whose `Capabilities.Playtime` is false
must always send `nil`.** Sending a zero would feed `status.Derive` a number it would treat as
evidence — which is precisely the failure the capability model exists to prevent.

## 9 · `Store` reads have no freshness parameter

Not a departure from a spelling but from an expectation, recorded because a reviewer will look for
one. Reasoned in [`04-concurrency-and-offline.md`](./04-concurrency-and-offline.md) §4.1.

## 10 · Package layout

`internal/`, one package per seam, with `fault`, `i18n`, `status` and `aged` as leaves that import
nothing but the standard library and each other in one direction.

`internal/` because ADR-0001 D1 rejects plugins: providers ship in the binary, so nothing outside
this module ever needs to implement these interfaces. Making them public would be publishing an API
with no consumer and no upgrade path.

`internal/arch` holds no code — only the tests that keep the import graph honest.
