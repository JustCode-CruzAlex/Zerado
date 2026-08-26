---
title: Zerado — interface contracts
discipline: API
doc-no: ZRD-API-00
rev: A
date: 2026-08-25
status: draft — for review
archetype: implementation-plan
ticket: "#6"
---

# The contracts

Nine seams, their exact Go shape, and the reasoning that chose each one.

The spine ([`../blueprint/06-data-seams.md`](../blueprint/06-data-seams.md)) decided **that** these
seams exist and what each is responsible for. This deliverable decides their **signatures** — the
job [`../blueprint/13-handoffs.md`](../blueprint/13-handoffs.md) §1 assigns to it: final
signatures, error contracts and sentinel errors, context and cancellation semantics, versioning and
stability guarantees, and package layout.

**No implementation is written.** What ships is contracts, doc comments, three `Null`
implementations the ADR itself designates as first-class states, one hand-entry reference provider,
five test fakes, and the tests that make the contracts checkable rather than assertable.

---

## 1 · Where everything is

| Package | Seam | Ticket item |
|---|---|---|
| `internal/provider` | Store provider — `Provider` · `Syncer` · `Enterer` | 1 |
| `internal/metadata` | Enrichment — covers, *sinopse*, genres, release data | 2 |
| `internal/pricing` | Price — current, all-time low, wait-or-buy | 3 |
| `internal/store` | The repository seam over `fft-database`'s schema | 4 |
| `internal/fault` | The error taxonomy | 5 |
| `internal/aged` | The staleness contract | 6 |
| `internal/cli` | The CLI verb surface as a versioned API | 7 |
| `internal/devicesync` | The Phase 4 boundary, sketched | 8 |
| `internal/library` | Domain types — and where item 9 is answered by *not* generalising | 9 |
| `internal/vault` · `internal/images` · `internal/audio` | Credentials · terminal images · sound | — |
| `internal/status` · `internal/i18n` | The four states · the catalogue key | — |
| `internal/arch` | No code. The tests that keep the import graph honest | — |

The reasoning for each signature is in its doc comment, beside the signature, because a decision
recorded somewhere else is a decision the next person changes without reading. This document holds
what does not belong to a single declaration.

---

## 2 · The gate question, answered by demonstration

> **Can a second store, a second metadata source, and a second media type each be added without
> changing a signature above the seam?** *Show it, do not assert it.*

### 2.1 · A second store

`internal/provider` carries **two implementations already**, and they are as different as the
product will ever have to accommodate:

| | `providertest.Fake` | `providertest.Manual` |
|---|---|---|
| What it is | a network store | a person typing |
| Implements | `Provider` + `Syncer` | `Provider` + `Enterer` |
| Credentials | two fields, one secret | **none** |
| Playtime | reported | **unknowable** |
| Item source | a stream | a form |
| Provider reference | the store's own id | **a UUID Zerado mints** |

Both produce the same `provider.Item`, both travel the same `Store.UpsertBatch` path, and the four
states derive correctly for both from `Capabilities` with no branch on identity anywhere.

`TestManualIsNotASyncer` asserts the negative the compiler cannot: hand entry does **not** implement
`Syncer`. `TestCapabilitiesAgreeWithTheInterfaces` asserts the positive.

**Adding GOG** is: implement the interfaces, declare the credential fields, register. Screens change
by zero, because `Z-02` renders `Capabilities().Credentials` and `Z-04` reads `Capabilities`.

### 2.2 · A second metadata source

`internal/metadata` also carries two — `metadata.Null` and `metadatatest.Fake` — and
`TestTheSeamNamesNoProvider` runs the same call against both through the interface.

The stronger evidence is what is **absent**: the word IGDB appears **only in comments** — seven
occurrences across three files, five of them in `metadata`'s package comment explaining the hedge and
two in `fault`, naming it as the example of a source that may never have heard of a cartridge — and
in **no signature, type, field or constant**. *(An earlier revision of this document said "one
package comment". The load-bearing half was right and the tally was not; corrected after the review
at `c4c8d95` counted them.)* The fake is keyed on **title and platform**, not on a store identifier —
which is the shape a source that has never heard of Steam would have, and therefore the shape a
cartridge needs. A fake keyed on an appid would have agreed with a Steam-shaped seam and proved
nothing.

`Attribution()` is a method, so swapping the source swaps the credit and a licence change is a data
change rather than a redesign.

### 2.3 · A second media type

**The ratified answer is that there must not be a generic seam, and the demonstration is that the
cost of adding one later is two columns.**

The ticket asks for "the seams that must be generic so books plug in". ADR-0001 **D5** pruned that,
and [`../blueprint/06-data-seams.md`](../blueprint/06-data-seams.md) §7 lists a media-type
abstraction as explicitly **not a seam**: *"an interface parameterised on a type that has one value
is machinery without a purpose."*

So `library.Game` has no type parameter and no `MediaType` field. The door is held open by the
schema — the entity is `item`, carrying an `item_type` constrained to `'game'` — which is one
`CHECK` constraint away from a second type and requires **no interface change at all**, because
nothing above the store reads a media type today and nothing would need to.

Recorded as a divergence from the ticket's wording, with the ratified authority named, in
[`06-divergences.md`](./06-divergences.md) §1. **Ratified canon wins over ticket prose; the
disagreement is surfaced rather than silently resolved.**

---

## 3 · Every seam is fakeable with no network

The ticket makes this a correctness test of the shape: *if it cannot be tested offline, the shape is
wrong.*

| Seam | Offline double | Reaches |
|---|---|---|
| Store provider | `providertest.Fake` · `providertest.Manual` | nothing |
| Metadata | `metadata.Null` · `metadatatest.Fake` | nothing |
| Price | `pricingtest.Fake` | nothing |
| Repository | `storetest.Fake` (in-memory) | nothing |
| Credentials | `vaulttest.Fake` | nothing |
| Images | `images.Null` | nothing |
| Audio | `audio.Null` · `audiotest.Recorder` | nothing |
| Error classifier | `fault.Classify` is a **pure function** over a `Transport` verdict | nothing |

`go.mod` pins a **full toolchain version** (`1.27.0`), not a bare language version. `go 1.24` is a
language version and Go cannot resolve a toolchain from it, so anything that resolves the pin before
running — the review vehicle, a reproducible build — fails *upstream of every test in this module*.
Sprint 0 #10's acceptance criteria already required a pinned toolchain;
`TestTheGoDirectiveIsAResolvableToolchainVersion` now makes it checkable, and the CI workflow reads
the number from `go.mod` rather than repeating it.

The whole test suite runs with the machine in a Faraday cage and with **no `go.sum` at all** —
`go.mod` declares no `require`, so there is no third-party dependency to download and no lock file to
write. `TestTheModuleHasNoThirdPartyDependencies` records that as a property of this stage rather
than a permanent rule. *(An earlier revision said "an empty `go.sum`"; there is no such file. The
substance was right and the artefact named was not.)*

**`internal/arch` makes it structural rather than habitual.** `TestOnlyProvidersMayReachTheNetwork`
parses every non-test file in the module and fails if any package outside `internal/provider`
imports `net`, `net/http`, `net/url`, `crypto/tls` or `net/rpc`. That is
[`../blueprint/07-offline-contract.md`](../blueprint/07-offline-contract.md) §7.3's proposed
grep-level review rule, executable — because a rule nobody can check is a convention, and
conventions lose to deadlines. The guard was verified by temporarily adding a violating import and
watching it fail, then reverting.

---

## 4 · The documents

| | |
|---|---|
| [`01-seams.md`](./01-seams.md) | The nine seams, what each decided and why |
| [`02-error-taxonomy.md`](./02-error-taxonomy.md) | Thirteen kinds × treatment × copy × exit code |
| [`03-cli-surface.md`](./03-cli-surface.md) | Verbs, flags, exit codes, output, stability policy |
| [`04-concurrency-and-offline.md`](./04-concurrency-and-offline.md) | Context, cancellation, the streaming shape, the age rule |
| [`05-phase4-sketch.md`](./05-phase4-sketch.md) | What crosses, and what Phase 1 must not make impossible |
| [`06-divergences.md`](./06-divergences.md) | Every departure from the spine's spelling, with its reason |
| [`07-open-questions.md`](./07-open-questions.md) | What this ticket did not decide, and who owns it |

---

## 5 · The assumption this deliverable makes about `fft-database`

Ticket #5 designs the physical schema **in parallel with this one**. Rather than block, the
assumption is stated so it can be checked in one place:

> **The repository interface is written against the *conceptual* model in
> [`../blueprint/09-erd.md`](../blueprint/09-erd.md), and assumes only its load-bearing claims** —
> `status_manual` nullable, `last_played_at NULL` meaning *not reported*, `absent_since` nullable,
> the item identity indexed and **not** unique, `(provider_id, provider_ref)` unique,
> `effective_status` derived and never stored, and `sync_run.status` ∈ `ok · partial · failed ·
> cancelled`.

Everything else — column types, index strategy, migration sequence, whether a generated column is
worth its cost — is `fft-database`'s and **changes nothing in these signatures**, which is the point
of the seam. The one place the two must agree on a name is the cross-device identity: the ERD calls
it `item_uid`, [`06-data-seams.md`](../blueprint/06-data-seams.md) §6.2 calls it `game_uid`, and the
Go type is `library.UID`. The **column name is `fft-database`'s to pick**; nothing above the store
sees it.
