---
title: Zerado — the error taxonomy
discipline: API
doc-no: ZRD-API-02
rev: A
date: 2026-08-25
status: draft — for review
archetype: concept-explainer
ticket: "#6"
---

# The error taxonomy

One closed set of kinds, shared by every seam, each distinguishable by the caller **because each one
renders differently on screen**.

The requirement that shaped it, from the ticket: *a private Steam profile is NOT a network error and
must not render as one.*

---

## 1 · Why that case is the whole design

A private Steam profile is:

- a **successful** HTTP request,
- returning a **well-formed** body,
- from a service that is **up**,
- using a key that **works**.

Every mechanism that would collapse it into "network error" — a bare error string, a status code, a
boolean `ok` — produces a screen that tells the player their connection failed when their connection
is fine, and sends them to debug the wrong thing.

And it is worse than a bad message. [`../blueprint/06-data-seams.md`](../blueprint/06-data-seams.md)
§2.4 forbids tombstoning on any run that is not `ok`, and `Z-03` §8.2 spells out the consequence: a
247-game library whose owner has just made their profile private would be deleted by a naive "the
provider's view is the truth" upsert, after which the screen would honestly report a catastrophe it
had caused. **A caller that cannot tell "the provider returned nothing" from "the provider could not
be reached" cannot honour that rule.**

`fault.KindEmpty` exists so it can.

---

## 2 · The thirteen kinds

| Kind | Cause | Treatment | Retryable | Exit | The screen |
|---|---|---|---|---|---|
| `offline` | no route · DNS | banner | **yes** | 3 | `OFFLINE` |
| `unreachable` | timeout · 5xx · reset | banner | **yes** | 4 | `STEAM UNREACHABLE` |
| `unauthorized` | 401 · 403 | refusal | no | 5 | `STEAM KEY REJECTED` |
| `rate_limited` | 429, with the provider's wait | banner | **yes** | 7 | banner + the wait |
| `empty` | 200 with zero items | refusal | no | 6 | `STEAM PROFILE PRIVATE` |
| `not_found` | the source has never heard of it | **designed empty** | no | 8 | the no-metadata composition |
| `malformed` | an answer we cannot read · a 4xx we caused | refusal | no | 9 | a refusal naming what we sent |
| `stale` | a cached value too old to act on | banner | no | 10 | the amber banner |
| `unsupported` | a capability the provider does not have | refusal | no | 10 | a refusal |
| `precondition` | legal operation, wrong state | refusal | no | 10 | a refusal |
| `conflict` | a uniqueness rule | refusal | no | 10 | a refusal |
| `cancelled` | the player stopped it | **none** | no | 130 | `CANCELLED` — no cue, no red |
| `internal` | our defect | **fatal** | no | 1 | `Z-11` |

`TestTaxonomyIsTotal` asserts every kind has a distinct machine name, a treatment and a catalogue
key. `TestExitCodesAreTotal` asserts every kind has an exit code that is not `0`. **Adding a kind
without deciding what it looks like on screen, and what a script should do about it, fails CI.**

### 2.1 · That guarantee was false once, and how it was repaired

The review at `c4c8d95` falsified it, empirically, and the falsification is worth keeping because the
failure was in the enforcement layer rather than in the design.

`Kinds()` was a **hand-maintained slice** while `Valid()` was **derived from the iota range**. Both
totality tests iterated the slice, and nothing asserted the slice covered the range. So a kind
inserted *mid-block* — the natural move, since the kinds are grouped semantically — was valid,
constructible, and **invisible to every test that existed to catch exactly that**. Probed, it
resolved to `String() == "unknown"` on a surface the CLI documents as stable API,
`Treatment() == TreatmentFatal`, `MessageKey() == "fault.internal"` and exit `1`. CI stayed green.

**The repair removes the second source of truth rather than adding a test to reconcile the two.**
`Kinds()` is now derived from the same range as `Valid()`, so they cannot disagree. Four assertions
were added on top, because deriving the set is only half of it — a kind that is *in* the set can
still silently inherit a switch default:

| Test | What it catches |
|---|---|
| `TestKindsCoversTheWholeRange` | a future contributor reintroducing a literal slice |
| `TestNoKindFallsThroughToTheDefaults` | a kind with no case in `String`, `Treatment` or `MessageKey` |
| `TestNoKindSilentlyInheritsTheInternalExitCode` | a kind with no case in `ExitCode` |
| `TestTaxonomyIsTotal` | as before, now over a set that is actually complete |

Verified by re-running the reviewer's own probe: inserting `KindThrottled` between `KindConflict` and
`KindCancelled` now fails **four assertions across two packages**, by number, naming each missing
case.

Two doc comments that asserted the enforcement while it did not exist — `treatment.go` and
`exit.go` — were corrected rather than quietly fixed. A comment claiming a guarantee the code does
not provide is worse than no comment, because it is what a reviewer reads instead of the test.

---

## 3 · The four network-ish failures do not collapse

`TestTheFourNetworkFailuresDoNotCollapse` asserts that `offline`, `unreachable`, `unauthorized`,
`empty` and `rate_limited` map to five **distinct** exit codes.

The reason is a caller nobody thinks about at design time: a nightly `cron` line that runs
`zerado sync`. It wants to retry on `offline`, back off on `rate_limited`, and alert a human on
`unauthorized`. It cannot do any of that if all five are `1`.

---

## 4 · Three rules the type enforces

### 4.1 · A fault carries a key, never a sentence

ADR-0001 **D9** forbids user-facing string literals in code, and an error message is user-facing
text. `Fault.MessageKey` is an `i18n.Key`; the `Kind` supplies a default and a provider may override
it when it has genuinely better copy for its own case — Steam's private-profile refusal names a
specific Steam privacy setting, which no provider-agnostic entry could.

The subject travels as a **substitution**, not as part of the sentence. That is what makes one
catalogue entry correct for Steam today and GOG later, in every language, with no new key.

### 4.2 · The wrapped cause is never rendered

`Fault.Error()` deliberately **does not include its cause**.

This is a credential-disclosure guard, not a style preference. An `*url.Error` from `net/http`
stringifies the whole URL, and Steam's URLs carry the player's API key as a query parameter. An
error string that concatenated its cause would be a disclosure path running straight through the one
screen a frustrated player is most likely to screenshot into a bug report.

The cause stays reachable through `errors.Is`/`As` for a log sink that has been told it may see it.
`TestErrorDoesNotLeakTheCause` and `TestTheEnvelopeCarriesNoMessage` assert both halves — the second
one against the JSON envelope, which is the single most likely thing in this product to be piped
into a file and attached to a bug report.

### 4.3 · An unclassified error is `internal`, never a plausible guess

`fault.KindOf` returns `KindInternal` for any non-`Fault` error. It does not guess `unreachable`.

A scripted caller retrying forever because an internal panic was reported as a timeout is a worse
outcome than a loud `1`, and a player told their network failed when a migration is corrupt has been
sent to fix the wrong thing.

---

## 5 · The classifier is a pure function, and that is deliberate

`fault.Classify` implements
[`../blueprint/07-offline-contract.md`](../blueprint/07-offline-contract.md) §5's decision tree. It
takes a `fault.Outcome` — a `Transport` verdict, a status code, an item count — and returns a
`*Fault`.

It does **not** take an `*http.Response` and does not import `net` or `net/http`. Two consequences,
and the second is the one that mattered:

1. The `fault` package stays outside the network rule's blast radius. §7.3's grep rule would
   otherwise have been violated by its own most reasonable-looking helper, after which the rule is a
   comment.
2. **Every branch is reachable in a test with no network.** `TestClassifyFollowsTheOfflineContract`
   walks all twelve.

The branch order is part of the contract: transport first (a rejected key is meaningless if the
request never arrived), then `401`/`403` (a credential verdict outranks anything derivable from the
body), then `404` and `429`, then `5xx`, then any other `4xx` — which is **our** bug and not their
outage, and must not be dressed as one.

`ClassifySync` is separate from `Classify` because **emptiness is a refusal for exactly one
operation**. A metadata lookup that finds nothing renders a designed composition; a price lookup
that finds no quote is simply no quote. Only a library sync may conclude, from an empty answer, that
the player has to go and change a setting somewhere else.
