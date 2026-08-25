---
title: Zerado — the offline contract
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-07
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: flowchart
ticket: "#2"
---

# The offline contract

Precisely what works with the network off, what degrades, and how each degrade is shown
honestly on screen.

The page promises this publicly — *"Runs with the network off. Once your library is synced,
`Zerado` doesn't need the internet to answer what to play tonight."* This document is what makes
that true rather than usually-true.

---

## 1 · The contract in one sentence

> **Every screen reads from the local file. Nothing renders from a network response.**

The seam architecture already forces it ([`06-data-seams.md`](./06-data-seams.md) §1): providers
write to the store, screens read from the store, and the two never meet. So the offline
behaviour of a screen is not a feature that has to be added to it — it is a consequence of where
its data comes from.

That gives a three-way classification, and every feature in the product is in exactly one of
them:

| Class | Meaning | Shown as |
|---|---|---|
| **WORKS** | Reads only local state. Behaves identically | Nothing. No banner, no notice, no difference |
| **DEGRADES** | Has local data, but it is not fresh | The **degrade banner** — what is stale, and how stale |
| **REFUSES** | Needs the network to do anything at all | An honest refusal naming the reason and the retry key |

**A feature is never in two classes**, and the class does not change with the weather. If a
screen is `WORKS`, it works with the machine in a Faraday cage.

---

## 2 · The contract, feature by feature

### Phase 1

| Feature | Class | Behaviour with the network off |
|---|---|---|
| Browse the library (`Z-04`) | **WORKS** | Identical. Every row, every state, every count |
| Filter and search (`Z-07`) | **WORKS** | Identical — it is a `WHERE` clause |
| Sort | **WORKS** | Identical |
| Read a game's detail (`Z-05`) | **WORKS** | Every locally-known field |
| Set a status (`Z-06`) | **WORKS** | Writes immediately. This is the product's core action and it has never needed a network |
| Add a game by hand (`Z-08`) | **WORKS** | The whole point of a physical shelf |
| Settings (`Z-09`) | **WORKS** | Except connecting a store, which is `REFUSES` |
| Help (`Z-10`) | **WORKS** | It is compiled in |
| **Sync (`Z-03`)** | **REFUSES** | Names the reason, shows when the last sync was, offers `r` to retry |
| **Connect a store (`Z-02`)** | **REFUSES** | Cannot validate a key without reaching the provider. Says so on submit, keeps what was typed |
| First run (`Z-01`) | **WORKS** | The "add a game by hand" door still opens. A player with no network can still start using the product today |

**Nine of eleven Phase 1 screens are fully functional offline.** The two that are not are the two
that are definitionally about reaching somewhere else.

### Phase 2 and later

| Feature | Class | Behaviour |
|---|---|---|
| Cover art, *sinopse* | **DEGRADES** | Cached values render normally. Missing ones render the designed *no-metadata* composition — **not** an error |
| Mood tags the player set | **WORKS** | Local |
| Mood tags inferred | **DEGRADES** | Whatever was inferred last, with its age |
| Prices (`Z-16`) | **DEGRADES** | Last known quote, **always with its age**. Never a bare number |
| Enrichment sync (`Z-12`) | **REFUSES** | Same shape as `Z-03` |
| Tonight (`Z-18`) | **WORKS** | Mood and time are local facts. This is the headline feature and it is offline-complete by design |
| Watchlist (`Z-20`) | **DEGRADES** | Verdicts computed from the last known prices, each stamped |
| Community, profiles, comments (Phase 4) | **REFUSES** | They are definitionally other people |

---

## 3 · How a degrade is shown

One component: the **degrade banner**, one row at the top of the body, above everything
including the pinned summary. It is specified in deliverable B; the spine fixes its **contract**:

```
 ⚠  OFFLINE   Showing your library as of 3 days ago.  r to retry
```

It carries exactly four things, always, in this order:

1. **A glyph and a label** — never colour alone (co-render rule). ASCII fallback `[!]`.
2. **What is unavailable** — the subject, not a generic apology.
3. **How stale what you are seeing is** — in units a person uses, not a timestamp.
4. **The key that retries it** — a degrade with no way out is a dead end.

And it obeys four rules:

- **It never hides content it could have shown.** It is one row, at the top, and the body keeps
  the rest.
- **It is never colour-only and never red.** Red is motion and alarm (brand manual §4.1); being
  offline is neither. The banner is chrome and amber.
- **It never appears when nothing is degraded.** A banner that is always there is furniture.
- **It survives every tier.** Per [`03-responsive.md`](./03-responsive.md) §4 it is on the
  never-hide list: if the product is showing stale data, saying so outranks whatever the row was
  going to be used for.

### 3.1 · The refusal copy

A refusal names what happened, why, and the next action — the brand manual's voice example 3 is
the template, and it is already the copy for the most likely case:

> *"Steam returned an empty library. Game details are private on your profile — Steam won't
> share the list until that's public. Settings → Privacy."*

The four Phase 1 refusals, in the ratified voice:

| Situation | Copy |
|---|---|
| No network at all | `No network. Last synced 3 days ago — everything below still works. r to retry.` |
| The provider is unreachable | `Steam didn't answer. Not your key — their end, or the connection. r to retry.` |
| The key is rejected | `Steam rejected that key. Check it hasn't been regenerated. Settings → Steam.` |
| The library came back empty | `Steam returned an empty library. Game details are private on your profile — Steam won't share the list until that's public. Settings → Privacy.` |

Note what none of them say: *"Something went wrong."* That is the one sentence a terminal user
cannot act on.

---

## 4 · The age rule

> **Any value that came from the network is rendered with its age. Always.**

```
  $15                      ← forbidden
  $15 · 3 days ago         ← required
  Last synced 3 days ago   ← required on the library, whenever a sync is not current
```

This is the rule most likely to be lost during a build, because dropping the age always makes the
layout tidier. It is also the rule the product's credibility rests on: Zerado's value proposition
is telling a player **not** to buy something. A stale price presented as current is not a cosmetic
defect — it is the product giving financial advice from memory.

Practically, every network-derived struct carries a `FetchedAt` / `ObservedAt` and the renderer
takes both together. There is no code path that can render one without the other, because they
arrive in the same value.

### 4.1 · Ages are said the way people say them

`3 days ago` · `2 hours ago` · `just now` · `last June`. Not `2026-08-22T04:11:09Z`. The exact
timestamp is available in the detail view for anyone who wants it.

Past **90 days**, an age stops being reassuring and becomes a warning; the banner changes from
chrome to amber and reads `Last synced in May. Prices this old are not useful. r to sync.`

---

## 5 · How the product knows it is offline

**It does not guess, and it never probes.**

There is no connectivity check, no background ping, no "am I online" heartbeat. Two reasons: a
reachability probe is exactly the background network traffic the page promises does not exist,
and it is unreliable anyway — a machine can be on a network and still unable to reach Steam.

Instead: **a request fails, and the failure is classified.**

```
  attempt ──► success ──────────────────► fresh data, no banner
      │
      └────► failure ──► classify ──┬──► no route / DNS      → OFFLINE banner
                                    ├──► timeout             → UNREACHABLE banner
                                    ├──► 401 / 403           → CREDENTIAL refusal (Z-02)
                                    ├──► 200 + empty         → PRIVATE PROFILE refusal
                                    └──► 5xx / other         → UNREACHABLE banner
```

The banner is raised by the failure and cleared by the next success. Between the two, nothing is
running.

The consequence is honest and worth stating: **Zerado does not know it is offline until you ask
it to do something that needs the network.** Opening the app on a plane shows no banner, because
nothing has failed — and nothing needs to. The banner appears when you press `r`. That is the
correct behaviour for a local-first product, and it is a better experience than a permanent
"you are offline" strip on a screen that does not need to be online.

---

## 6 · First run with no network

The screen nobody writes down, and one of the two the ticket names as deciding whether the
product feels solid.

A player who opens Zerado for the first time with no connection **still has a working product**:

```
  Zerado ✦ First run

  FIRST RUN

  Zerado keeps your game library in one file on this machine.
  Nothing here needs an account, and nothing is sent anywhere.

    ▌ Add a game by hand              Works now — discs, cartridges, anything
      Connect a store                 Needs a connection. Not available right now
      Look around first               An empty library, and every key that works on it

  ↑↓ move   ⏎ choose   ? help   q quit
```

Three properties that make it right:

- **the unavailable door is shown, disabled, with its reason** — not hidden. Hiding it would
  make the player think Zerado cannot connect to stores at all;
- **the available doors are genuinely useful** — a player can enter their shelf and have a real
  library before they ever connect anything;
- **nothing is presented as broken.** Nothing is broken. One door is closed and it says why.

---

## 7 · How this gets verified

The offline contract is the easiest thing in this bundle to claim and the easiest to get wrong,
so it is tested as a **property**, not as a demo:

1. **A no-network integration run.** The whole Phase 1 suite executes with the HTTP client
   replaced by one that fails every request. Every screen in the `WORKS` class must render
   byte-identically to its online golden, **including the absence of a banner**.
2. **A no-network-with-stale-data run.** Same, with a library that was last synced 100 days ago.
   Every `DEGRADES` screen must show its banner, with the age, and every stale value must carry
   one.
3. **A grep-level rule in review:** no `net/http` import outside the provider packages, and no
   provider constructing its own `*http.Client`. That is what keeps class membership from drifting
   quietly over a year of changes.
