---
title: Zerado — the nine seams
discipline: API
doc-no: ZRD-API-01
rev: A
date: 2026-08-25
status: draft — for review
archetype: implementation-plan
ticket: "#6"
---

# The nine seams

What each one decided, and the one argument that decided it. The full reasoning for every signature
is in its doc comment; this is the map.

---

## 1 · The store provider — designed against the cartridge, checked against Steam

The likeliest way to get this seam wrong is to design it around Steam, because Steam is the only
provider Phase 1 builds. So it was written against **manual physical entry first** and checked
against a store second — the reverse of the tempting order, and the reason it has **three**
interfaces rather than one.

```go
type Provider interface {                    // every source, including a person with a keyboard
    ID() ID
    Display() string
    Capabilities() Capabilities
}

type Syncer interface {                      // a source that can fetch over a network
    Provider
    Sync(ctx context.Context, c Credentials) (Stream, error)
}

type Enterer interface {                     // a source whose items arrive because a person typed them
    Provider
    Form() []EntryField
    Compose(e Entry) (Item, error)
}
```

Read `providertest.Manual` and note what is **absent**: no credentials, no pagination, no rate limit,
no HTTP client, no `Sync`, no playtime. Each absence is a place a one-interface design would have
forced physical to lie — by returning `ErrManual`, by declaring fields it does not want, or by
reporting a zero the derivation would treat as a fact.

**`Enterer` is a capability, not a provider kind.** A future GOG that syncs *and* lets a player add a
game its API missed implements both, and `Z-08` gains a source picker rather than a special case.

**Everything downstream reads `Capabilities`, never `ID`.** `provider.Check` asserts that the
declaration and the implemented interfaces agree — because `Capabilities` deliberately duplicates two
facts the type system already carries (screens need them before they have a reason to type-assert),
and duplicated facts drift.

**A hand-entry provider may not claim `Playtime`.** `Check` refuses the declaration, because a typed
number would otherwise become evidence in `status.Derive`.

## 2 · Metadata — a hedge whose reason changed and whose shape did not

The seam was designed when IGDB's non-commercial terms sat badly against an affiliate-funded
product. Affiliate links are dropped, so the product is cleanly non-commercial and IGDB's published
test is answered — **as a reading of a published rationale, not a guarantee.** A direct confirmation
is a founder action.

So the hedge stays exactly as designed. Removing it because one risk receded would be the wrong
lesson: a source that is named today can change its terms tomorrow, which is true of every source and
is not a fact about IGDB.

**The test of agnosticism:** the word IGDB appears **only in comments** — seven occurrences across
`metadata/metadata.go`, `fault/kind.go` and `fault/treatment.go` — and in **no signature, type, field
or constant**. The fake is keyed on **title and platform**, which is the shape a source that has
never heard of Steam would have — and therefore the shape a cartridge needs.

**`Null` is production code, not a test helper.** Having no metadata is a designed state: the detail
view shows a designed no-metadata composition, never an error banner. That is the difference between
a product that works without a source and one that is broken without it, and it is why
`KindNotFound`'s treatment is `TreatmentDesignedEmpty`.

## 3 · Price — the disclosure is structural

`Quote` carries the current price, the all-time low, **when the low was**, the shop and a plain shop
link. There is **no field an affiliate tag could occupy**, which is how the decision is enforced
rather than remembered.

`LowAt` is mandatory: *"it was R$ 19 once"* means something different depending on whether once was
last month or in 2019, and only one of those should make a player wait.

`Quote.Verdict` answers wait-or-buy **from one response, locally, offline** — which is what makes the
Phase 3 watchlist a `DEGRADES` feature rather than one that stops. A verdict computed by the source
would be a verdict Zerado could not produce when it matters.

`VerdictUnknown` renders as no verdict, never as "wait". A product that guesses when it does not know
is giving advice it has not earned.

## 4 · The repository — the only thing that knows a database exists

Reads take a `Query` and return domain types. Writes take domain types and return what the screen
needs to say what happened. **No query builder, no expression, no transaction handle, no rows
cursor** — a caller cannot construct SQL through it.

`Games` and `Counts` take **the same `Query`**, which makes
[`../blueprint/05-state-machine.md`](../blueprint/05-state-machine.md) §7 rule 2 enforceable: the
summary describes the filtered set, so the list and the numbers above it cannot describe different
sets. `Counts` ignores `Limit`/`Offset`, because the summary describes the set and not the page.

`SetStatus` takes a `*Status` so that clearing an override is expressible at all — and it stamps
`status_changed_at` on a clear, because a clear is a change the player made and Phase 4 needs a
timestamp for it.

`UpsertBatch` **never writes the manual status**. That invariant lives in the store rather than in a
provider because the store is the only writer and therefore the only place it can be guaranteed.

`ReconcileAbsence` takes a `RunID`. See [`06-divergences.md`](./06-divergences.md) §6 — the argument
type **is** the tombstone guard.

**No query shape exists that no screen needs.** The facets are the three `Z-07` renders plus a
presence mode, because a filter language would be an abstraction with one consumer and the consumer
is a bar with three chips.

`Query.Empty` and `Query.Facets` are a **pair**, and the invariant that binds them — a query that
reports itself filtered always names at least one facet — is asserted over the whole cross-product
rather than case by case. It has to be: the contract's own idiom is
`if !q.Empty() { name(q.Facets()[0]) }`, so a combination that reported itself filtered while naming
nothing would index an empty slice. `Presence: Either` did exactly that until the review at
`4484d9a`, and the repair was the invariant rather than the case.

## 5 · The error taxonomy

Thirteen kinds, each with a screen treatment, a catalogue key and an exit code, all asserted total by
tests. Full reasoning in [`02-error-taxonomy.md`](./02-error-taxonomy.md).

## 6 · Offline and staleness

`aged.Value[T]` on every network-derived value, and **no freshness parameter on any read**. Full
reasoning in [`04-concurrency-and-offline.md`](./04-concurrency-and-offline.md) §4.

## 7 · The CLI

Declared as data, golden-tested, with a written stability policy. Full surface in
[`03-cli-surface.md`](./03-cli-surface.md).

## 8 · The Phase 4 boundary

Sketched, and its rule enforced by a reflection test rather than a paragraph.
[`05-phase4-sketch.md`](./05-phase4-sketch.md).

## 9 · Media types

Not a seam, by ratified decision. [`06-divergences.md`](./06-divergences.md) §1.

---

## The three supporting seams

**`Vault`** — keyed by `(provider, key)`, so forgetting one provider cannot reach another's secrets
and `DeleteProvider` is a complete operation rather than a loop the caller has to get right. Nothing
in the library file can hold a credential, because `SaveConnection` takes an `accountRef` and there
is no method that could take a secret.

**`Images`** — `Cover` takes no context and returns no error, because it never fetches and never
blocks; `Warm` is the only method that touches the network and it returns immediately. A screen never
learns the protocol, so adding Sixel later changes one implementation and no screen. `images.Null` is
a **supported configuration**, not a fallback.

**`Audio`** — `Cue` has no error and no context, for the same reason. `audio.Null` is the **default
build**, so the silent path is exercised by every test run rather than being the branch that only
fails on somebody's laptop. `State` distinguishes *not compiled*, *off*, *unavailable* and
*overridden by the environment*, because `Z-09` must say which — those are four different facts and
collapsing them tells a player their setting changed when it did not.
