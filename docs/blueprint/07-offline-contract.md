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

> **Your library is always yours. Not every feature is a library.**

Founder direction, 2026-08-25, correcting an earlier draft that drew this too strictly:

> *"working offline doesn't mean if I'm online I have some features that will not fly on offline, so
> radio/streaming only works when online; if offline we can't stream, so an obvious degradation."*

The first draft treated *local-first* as *everything works offline*, and that is a stricter promise
than the page makes and a worse product than the page describes. **Local-first means the library —
the collection, the states, the notes, the ratings — is readable and editable with the network off,
always, completely.** It does not mean a radio stream plays without a network. Some features are
online by nature, they always were, and saying so is honesty rather than retreat.

What does not change: **every screen reads from the local file, and nothing renders from a network
response.** That is still the mechanism.

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
| **NEEDS THE NETWORK** | Online by nature. Offline it stops | It states its own condition where it lives |

**A feature is never in two classes**, and the class does not change with the weather. If a
feature is `WORKS`, it works with the machine in a Faraday cage.

**`NEEDS THE NETWORK` is not a failure class.** How it *states* its condition depends on what kind
of thing it is:

| Kind | Offline behaviour |
|---|---|
| An **action** the player just took — `Z-02` submit, `Z-03` sync | A **refusal**, naming what happened, why, and the next action. The player asked; they are owed an answer |
| An **ambient feature** — the radio | It **stops**, and says so where it lives. A stream that ends when the network does is behaving correctly, and should read the way a song ending reads, not the way an error reads |

The bar is that **every degrade is visible and never silent.** It is not that every degrade is loud.

---

## 2 · The contract, feature by feature

### Phase 1

| Feature | Class | Behaviour with the network off |
|---|---|---|
| Browse the library (`Z-04`) | **WORKS** | Identical. Every row, every state, every count |
| Filter and search (`Z-07`) | **WORKS** | Identical — it is a `WHERE` clause |
| Ordering | **WORKS** | Identical — and **fixed**: title A→Z, with no sort control in Phase 1. A list must have an order, and naming it is honest; an on-screen sort indicator would imply a control that does not exist |
| Read a game's detail (`Z-05`) | **WORKS** | Every locally-known field |
| Set a status (`Z-06`) | **WORKS** | Writes immediately. This is the product's core action and it has never needed a network |
| Add a game by hand (`Z-08`) | **WORKS** | The whole point of a physical shelf |
| Settings (`Z-09`) | **WORKS** | Every dial reads and writes locally. *Connecting* a store is `Z-02`, a separate screen with its own class — Settings only routes there |
| Help (`Z-10`) | **WORKS** | It is compiled in |
| **Sync (`Z-03`)** | **NEEDS THE NETWORK** | Names the reason, shows when the last sync was, offers `r` to retry |
| **Connect a store (`Z-02`)** | **NEEDS THE NETWORK** | Cannot validate a key without reaching the provider. Says so on submit, keeps what was typed |
| First run (`Z-01`) | **WORKS** | The "add a game by hand" door still opens. A player with no network can still start using the product today |
| Interface FX (`Z-09 § Audio`) | **WORKS** | Short cues, local, no network. They keep working offline |
| **Radio** (`Z-09 § Audio`) | **NEEDS THE NETWORK** | It streams, so offline it **stops** — and that is fine. No bundled music, so there is nothing to fall back to and nothing to apologise for ([`12-audio.md`](./12-audio.md)) |
| **Cover art** (`Z-15`) | **DEGRADES** | Fetched online, cached, then available offline. An uncached cover shows the designed no-cover tile, never a broken-image box ([`17-images.md`](./17-images.md)) |
| Fatal error (`Z-11`) | **WORKS** | It reports a local failure and depends on nothing. It is the one screen that is *more* reliable offline, because it reaches for nothing at all |

**Nine of the twelve Phase 1 screens are `WORKS`** — `Z-01`, `Z-04`, `Z-05`, `Z-06`, `Z-07`,
`Z-08`, `Z-09`, `Z-10`, `Z-11`. **Two `NEED THE NETWORK`** — `Z-02` and `Z-03`, the two that are
definitionally about reaching somewhere else, and both state it as a refusal because both are
actions the player just took. **One `DEGRADES`** — `Z-15`, the cover deck, which caches what it
fetched and shows a designed tile for what it has not.

Twelve screens, three classes, every screen in exactly one — which is the §1 invariant holding,
stated as a count so it can be checked rather than asserted.

### Phase 2 and later

| Feature | Class | Behaviour |
|---|---|---|
| Cover art, *sinopse* | **DEGRADES** | Cached values render normally. Missing ones render the designed *no-metadata* composition — **not** an error |
| Mood tags the player set | **WORKS** | Local |
| Mood tags inferred | **DEGRADES** | Whatever was inferred last, with its age |
| Prices (`Z-16`) | **DEGRADES** | Last known quote, **always with its age**. Never a bare number |
| Enrichment sync (`Z-12`) | **NEEDS THE NETWORK** | Same shape as `Z-03` |
| Tonight (`Z-18`) | **WORKS** | Mood and time are local facts. This is the headline feature and it is offline-complete by design |
| Watchlist (`Z-20`) | **DEGRADES** | Verdicts computed from the last known prices, each stamped |
| Community, profiles, comments (Phase 4) | **NEEDS THE NETWORK** | They are definitionally other people |

---

## 3 · How a degrade is shown

One component: the **degrade banner**, one row at the top of the body, above everything
including the pinned summary. It is specified in deliverable B; the spine fixes its **contract**:

```
 ▌ OFFLINE   Last synced 3 days ago. Everything here still works.
```

It carries exactly four things, always, in this order:

1. **An uppercase label word** — `OFFLINE`, `PRICES OFFLINE`, `STEAM PROFILE PRIVATE`. The word
   carries the state, so it survives with zero colour and needs no warning glyph. The leading
   `▌` is **structural**, not a state channel.
2. **What is unavailable** — the subject, not a generic apology.
3. **How stale what you are seeing is** — in units a person uses, not a timestamp.
4. **The key that retries it**, whenever there is one — a degrade with no way out is a dead end.

> **No `⚠`.** U+26A0 is listed **`Emoji`** in Unicode's `emoji-data.txt`, so it carries an emoji
> presentation: many terminals render it in colour, from the font, which `NO_COLOR` cannot switch
> off — and a coloured emoji presentation is typically double-width, which would break the column
> budget. It is exactly the class of glyph [`08-prior-draft-analysis.md`](./08-prior-draft-analysis.md)
> §1 rejects.
>
> **A correction, recorded rather than quietly fixed:** an earlier draft of this note also called
> U+26A0 *Ambiguous* width. It is **Neutral** (UAX #11, verified against UCD at source). The
> conclusion held, but one of its two premises did not — and in a bundle whose value is that its
> facts are checked, a wrong premise that happens to reach the right answer is still a defect.
> *(Caught by `fft-tui-designer`.)*
>
> The label word does the job better anyway: it is readable, greppable, and speakable by a screen
> reader.

And it obeys four rules:

- **It never hides content it could have shown.** It is one row, at the top, and the body keeps
  the rest.
- **It is never colour-only and never red.** Red is motion and alarm (brand manual §4.1); being
  offline is neither — red would call a promised behaviour a fault. The banner is **chrome when
  it is informational** and **amber only when the player must do something**. That distinction is
  the rule, and the component spec is
  [`../design/01-design-system.md`](../design/01-design-system.md) §12.
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
| No network at all | `No network. Last synced 3 days ago — everything below still works.` |
| The provider is unreachable | `Steam didn't answer. Not your key — their end, or the connection.` |
| The key is rejected | `Steam rejected that key. Check it hasn't been regenerated.` |
| The library came back empty | `Steam returned an empty library. Game details are private on your profile — Steam won't share the list until that's public. Settings → Privacy.` |

**Three rules the copy above obeys, each because the first draft broke it:**

1. **The retry key is not baked into the sentence.** The first draft ended each line with `r to
   retry.`, which is wrong on any screen with a text input focused: WCAG **2.1.4** means `r` types a
   literal `r` there, so the copy would name a key that does nothing. The **screen** supplies the
   affordance — `r to retry` on `Z-04`, `⏎ to try again` on `Z-02` — and the message supplies only
   the fact.
2. **A message never points at the screen the player is standing on.** `Settings → Steam` is useful
   from the library and absurd on `Z-02`, which *is* the Steam settings. The destination is the
   screen's to fill in, not the message's.
3. **A message never describes a composition it cannot see.** The first draft's `— everything below
   still works` names a library beneath the banner. That is true on `Z-04`, false on `Z-02` (nothing
   is synced yet) and shapeless on `Z-03`, where the body *is* the message. The clause is now the
   screen's to add where it is true.

*(All three found by `fft-tui-designer` while writing `Z-02` and `Z-03`.)*

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
chrome to amber. **In Phase 1** it reads `Last synced in May. Anything you have played since then
is missing. r to sync.` — the library is the only thing there is to be stale about.

From Phase 3, when there are prices, the copy names them instead, because a stale price is worse
than a stale playtime count. **It must not name them before then:** copy that mentions a capability
the build does not have presents something unbuilt as working.

---

## 5 · How the product knows it is offline

**It does not guess, and it never probes.**

> **One clarification, because §6 draws a first-run screen with the Connect door already disabled
> before anything has failed.** Reading the operating system's **local routing table** — is there a
> default route? — is *not* probing: it is a syscall against state the kernel already holds, it
> emits **no packet**, it contacts nothing, and it is instant. That is what lets `Z-01` disable a
> door at first paint without contradicting this section.
>
> The distinction that matters is **packet versus no packet**, not *knowing* versus *not knowing*. A
> reachability check sends traffic to somebody; reading your own routing table does not. What
> Zerado still refuses to do is send anything to anyone to find out whether it can.
>
> It is also honest about its own limits: a default route is not the internet. A machine can have
> one and still not reach Steam — which is why the *classification* below is still driven by real
> failures, and the routing-table read only ever downgrades an affordance, never claims success.

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
