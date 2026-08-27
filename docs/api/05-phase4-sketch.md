---
title: Zerado — the Phase 4 sync boundary, sketched
discipline: API
doc-no: ZRD-API-05
rev: A
date: 2026-08-25
status: draft — for review
archetype: concept-explainer
ticket: "#6"
---

# The Phase 4 boundary

**Sketch only. Nothing is built and no Phase 1 code calls it.**

It exists because ADR-0001 **D4** decides *what crosses* now rather than in Phase 4, and because that
is the most expensive decision in the bundle to reverse: it decides what the schema carries from the
first migration onward.

---

## 1 · The rule

> **Only what the player typed crosses. Everything a machine can recompute, each device
> recomputes.**

| Crosses | Never crosses |
|---|---|
| the manual status + when it changed | the library itself |
| hand-entered games — the one row class with no other copy | `playtime_minutes` · `last_played_at` |
| a stable cross-device identity | cover art · *sinopse* |
| mood tags the player assigned (Phase 2) | prices |
| | **credentials. Ever.** |

---

## 2 · The rule is a test, not a paragraph

`TestPayloadCarriesOnlyWhatThePlayerTyped` walks the payload structs by reflection against an
allow-list. `TestNothingResemblingACredentialCanCross` fails on any field whose name contains `key`,
`token`, `secret`, `password`, `credential`, `auth`, `cookie` or `session`.

The reason is a specific, likely future. In 2027, in a Phase 4 sprint, long after everyone who read
D4 has moved on, adding playtime to the payload will look like an obvious improvement — the server
already has the rows and showing hours on another device is a nice feature. A paragraph will not
stop it. A test that fails by name, with the decision quoted in the failure message, will.

---

## 3 · What Phase 1 must not make impossible — and does not

| Needed in Phase 4 | Present in Phase 1 | Cost of adding it later |
|---|---|---|
| a stable cross-device identity (`library.UID`) | assigned at insert, indexed, **not unique** | a migration that has to invent identities for rows whose titles the player has since edited |
| `status_changed_at`, written on **every** change including a clear | `Store.SetStatus` stamps it, always | conflicts unresolvable for every pre-existing row |
| `merged_into` | on `library.Game`, excluded from every query | a migration that rewrites primary keys |
| a single writer to attach to | `store.Store` is the only writer | a merge engine that has to find every write path |

`TestAClearedOverrideCanCross` asserts `Change.Status` is a pointer: *"I have no opinion"* is a change
the player made, and a non-nullable field would make a clear unsyncable — after which the player's
two devices disagree with themselves forever.

---

## 4 · Two calls, and why not more

`Push` and `Pull`. The payload is small and idempotent, and anything richer — a subscription, a
stream, an operation log, CRDTs — would be designing for a conflict shape that does not exist. D4
rejects those for the stated reason: they are correct for concurrent multi-writer editing, and the
conflicting parties here are **one person on two devices who agree about what they did.**

`Receipt` reports `Accepted` and `Superseded` because last-write-wins is an acceptable simplicity
choice only if it is not a silent one — `Z-22` shows the last merge and what it resolved.

**Not decided here, and named so it is not assumed to have been:** whether Phase 4 accounts are
email-and-password, OAuth or something else. Nothing in Phase 1 depends on it.
