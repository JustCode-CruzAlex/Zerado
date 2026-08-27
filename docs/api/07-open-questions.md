---
title: Zerado — what these contracts did not decide
discipline: API
doc-no: ZRD-API-07
rev: A
date: 2026-08-25
status: draft — for review
archetype: project-brief
ticket: "#6"
---

# Open questions

Named so they are not assumed to have been settled.

---

## 1 · Which audio library keeps the **default** build pure Go

**Assigned to this ticket** by [`../blueprint/13-handoffs.md`](../blueprint/13-handoffs.md) §4, and
**not answered here.**

The handoff itself states the requirement: it *"needs a real per-platform check of cgo dependency,
not an assumption."* That is a verification against six build targets, not a signature decision, and
answering it from memory is exactly the failure this bundle's own review register keeps catching.

**It blocks nothing.** ADR-0001 D6 puts the implementation behind a build tag; `audio.Null` is the
**default** build and is what `internal/audio` ships. Whichever library is chosen changes **no
signature in `audio.Audio`** — which is what the build-tag seam was designed to guarantee, and is the
test of whether the seam is right.

**Owner:** this ticket's follow-up, or the Phase 1 audio implementation ticket. **What it needs:** a
per-platform check of each candidate's cgo dependency on darwin/arm64, darwin/amd64, linux/amd64,
linux/arm64, windows/amd64.

## 2 · The `library.UID` normalisation rules

`fft-database`'s, per §4 of the handoff: it needs testing against a real library's titles. The
**policy** is decided and is carried in the Go doc comment — a merge *hint*, never authoritative,
indexed and not unique, with ambiguous matches shown to the player rather than guessed.

## 3 · The column name for the cross-device identity

[`../blueprint/09-erd.md`](../blueprint/09-erd.md) calls it `item_uid`;
[`../blueprint/06-data-seams.md`](../blueprint/06-data-seams.md) §6.2 calls it `game_uid`. The Go
type is `library.UID` and **nothing above the store sees the column name**, so this is
`fft-database`'s to settle. Flagged only because two ratified documents spell it differently and the
next reader should know that is a known inconsistency rather than a clue.

## 4 · `Capabilities.Progress` in **four** screen specs

`Z-04` (line 1027), `Z-05` (399), `Z-06` (366) and `Z-08` (75, 375) name a field that ADR-0001 D5
withdrew. The contract uses `Playtime`. Those specs should be corrected on their next revision; this
ticket does not own them.

Full list and the reason the count changed: [`06-divergences.md`](./06-divergences.md) §7.

## 5 · The D9 lint's exact implementation

[`../blueprint/16-i18n.md`](../blueprint/16-i18n.md) §1.1 assigns the lint's shape here and requires
that it *exist, run in CI, and fail the build*.

**What is decided and shipped:** the type it enforces against. `i18n.Key` exists; `fault.Fault`
carries a key and never a sentence; `metadata.Attribution` carries a key, with `Verbatim` as the one
sanctioned exception for wording a licence requires unaltered; every CLI verb and flag carries a
`SummaryKey`, asserted by `TestEveryVerbHasACatalogueKey`.

**What is not:** the lint binary itself. It is a build-tooling deliverable — a `go/analysis` pass
over render and format calls with an explicit allow-list — and it needs a render path to analyse.
There is not one yet. **Owner:** the Phase 1 TUI ticket that creates the first render path.

## 6 · The `--json` payload shapes per verb

The envelope is fixed and versioned. What `data` contains for `list` and `doctor` is defined when
those verbs are built, and is **additive** under the §5 stability policy in
[`03-cli-surface.md`](./03-cli-surface.md).

## 7 · Two things ADR-0001 names as not decided, restated so they are not assumed

- **Which metadata provider Phase 2 ships with.** The seam is decided; the provider is not.
- **Whether Phase 4 accounts are email-and-password, OAuth or something else.** Nothing in Phase 1
  depends on it.

## 8 · Withdrawn from this ticket's scope

A shared, product-agnostic **filterable table component contract** was added to this ticket by relay
on 2026-08-25 and **withdrawn by a later relay the same day** (*"forget about the help screen, too
much noise… I'll create that on flowforge side"*), with ticket #8 closed. Nothing of it remains in
this branch — recorded here so that a reader who finds the instruction in the session history knows
it was dropped deliberately rather than forgotten.
