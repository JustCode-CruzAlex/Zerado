---
title: Zerado — the audio subsystem
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-12
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: concept-explainer
ticket: "#2"
---

# The audio subsystem

> **Founder direction, 2026-08-25:** *"The audio is part of the Phase 1."*

---

## 0 · This reverses a verdict in this bundle, and the reversal is recorded rather than hidden

[`08-prior-draft-analysis.md`](./08-prior-draft-analysis.md) §3 rejected the prior draft's audio
proposal outright, with five reasons. **That verdict is superseded.** Audio ships in Phase 1 and my
job is to design it, not to decide whether it exists.

**The reversal is not a change of mind about the same object — the object changed.** What was
rejected was a *network streamer, always on*. What ships is a *local, bundled, opt-in subsystem,
off by default, that makes no network requests at all*. Three of the five objections retire because
of that, not because they were overruled.

The five reasons are not deleted. Two survive the reversal as **binding constraints**, and one of
those is still an open founder decision:

| Original objection | Status after the reversal |
|---|---|
| Contradicts *"no telemetry"*, *"works with the network off"*, *"only traffic is services you connected"* | **Retired — the object changed.** What was rejected was a *network streamer, always on*; what ships is a *local, bundled, opt-in* subsystem making no network requests at all. §3 |
| Costs the pure-Go single-binary distribution (cgo) | **Survives as a requirement** — §5 |
| Costs the 60 fps and zero-goroutine-leak budget | **Survives as a requirement** — §6 |
| Acquires a music-rights surface | **Survives as a founder decision, unresolved** — §7 |
| Nostalgia-kitsch versus retro-future | **The founder's call, and it is made.** §8 records the bar the sound has to clear so the reversal does not quietly become the thing the brand manual rules out |

> **Provenance.** The direction reached this session as an agent relay, which by its own header
> carries no ratification authority. The work is done as directed because it is document revision on
> a draft PR — reversible, and the founder confirms or restores it at the gate. It is recorded this
> way so that confirming it is a decision rather than an assumption.

---

## 1 · What it is

Two independent channels:

| Channel | What it is | Default |
|---|---|---|
| **Music** | A background bed. The retro-future soundtrack | **off** |
| **Interface FX** | Short cues: a completed sync, an error, a state change to `ZERADO` | **off** |

They mute independently and have independent volumes. Someone may want the keyclicks without the
soundtrack, or the soundtrack without the keyclicks, and neither is the odd request.

---

## 2 · Off by default, and that is the whole first-run posture

**Audio is off on first run.** Not quiet — **off**, until an explicit opt-in in
`Z-09 Settings § Audio`.

A terminal program that makes noise the first time someone runs it is a program people uninstall,
and they uninstall it before they have seen anything else. The retro-future feeling is a **reward
the player turns on**, not an ambush.

The cost of off-by-default is that nobody finds the feature, so `Z-01 First run` carries **one
quiet line** naming that audio exists and where to enable it. A line, not a prompt — `Z-01` still
collects nothing ([`01-screen-inventory.md`](./01-screen-inventory.md) §5).

**WCAG 1.4.2 Audio Control (Level A)** is satisfied structurally by this posture rather than by a
control bolted on afterwards: nothing plays automatically, so the three-second threshold is never
reached without the player having asked for it; and once enabled, `m` is a global always-reachable
mute with the two channels independently controllable.

---

## 3 · Bundled, never streamed — and this is what keeps the public promises intact

The published page says *"no telemetry running in the background"*, *"Runs with the network off"*,
and *"The only network traffic is `Zerado` reaching out to the services you've connected."*

**Zerado's audio makes no network requests, ever.** Tracks and cues are bundled with the binary or
read from a local path. There is no streamer, no CDN, no fetch, no cache warm.

That is not a limitation accepted reluctantly. It is the design decision that lets audio ship at
all without amending a ratified promise — and it means audio is **`WORKS`** in the offline contract
([`07-offline-contract.md`](./07-offline-contract.md)), like every other local-only feature.

A streamed soundtrack would be a different product decision requiring a different ratification, and
it is out of scope.

---

## 4 · The seam

Audio is a named seam like every other, and the same rule applies: **a screen never talks to the
audio device.** A screen emits a semantic cue; the seam decides whether anything is audible.

```go
type Cue string   // "sync.done" · "sync.failed" · "status.zerado" · "error"

type Audio interface {
    // Cue is FIRE-AND-FORGET and MUST NOT BLOCK. It never returns an error:
    // a cue that cannot play is silence, not a failure a screen must handle.
    Cue(c Cue)

    Music(on bool)
    SetVolume(ch Channel, v int)   // 0..100
    Mute(ch Channel, on bool)
    Muted(ch Channel) bool

    // State is what Z-09 renders honestly — including WHY audio is unavailable.
    State() AudioState
    Close() error
}

type AudioState struct {
    Enabled   bool
    Available bool
    Reason    Unavailable   // none | no_device | ssh | ci | env_disabled | init_failed
    Backing   string        // the device/driver actually in use, or ""
}
```

`NullAudio` implements the whole interface and does nothing. It is the default, it is what runs
whenever audio is disabled or unavailable, and it is what the tests use. **There is no code path
where a screen has to know whether audio exists.**

---

## 5 · The distribution problem, stated honestly

Audio playback in Go means an audio-device binding, and the common ones use **cgo** on at least
some platforms. [`06-data-seams.md`](./06-data-seams.md) §5.3 chose a **pure-Go** SQLite driver
specifically to keep a single cross-compiled static binary, because the page's claim is *"It's a
text program. It starts instantly."*

**Decision: the audio implementation sits behind a build tag.**

| Build | Contains | Audio |
|---|---|---|
| `zerado` (default, **pure Go**, cross-compiles to every target with no toolchain) | `NullAudio` | silent — and Settings says why |
| `zerado` (release build, `-tags audio`) | the real player | full |

Three things this buys:

1. The **cross-compile matrix stays free**. A platform with no audio toolchain still gets a binary.
2. *"Completely removable at runtime"* becomes **removable at compile time as well** — a strictly
   stronger guarantee.
3. `NullAudio` is not a fallback that might be untested; it is the **default build**, exercised by
   every test run and every CI job.

> **Deliberately not decided here:** *which* audio library. That is `fft-api-designer`'s call
> against the seam above, and it needs a real check of which players are pure-Go on which
> platforms. This blueprint requires only that the choice does not make the **default** build
> impure. See [`13-handoffs.md`](./13-handoffs.md).

---

## 6 · It never blocks, never leaks, never hangs the TUI

The performance bar is 60 fps and zero goroutine leaks, and audio is the classic way to lose both.

| Rule | Mechanism |
|---|---|
| **Never blocks the render loop** | `Cue()` is a non-blocking send on a **buffered** channel. If the buffer is full the cue is **dropped**, silently — a missed sound is not worth a dropped frame |
| **Never leaks** | One owned goroutine, started at init, stopped by `Close()` on a `context`. The seam owns it; no screen ever starts one |
| **Never hangs on device failure** | Device init happens **once, off the render path**, with a timeout. A failure sets `Available=false` and `Reason`, and the program continues silent |
| **Never errors at a screen** | `Cue()` has no error return. There is no failure a screen could usefully handle |
| **Never blocks shutdown** | `Close()` has its own timeout. A stuck audio device does not stop `q` from quitting |
| **Muting music RELEASES the device** | Muting the music channel **halts decode and releases the audio device** — it does not hold a gain of zero. Unmuting restarts it. A muted bed holding a goroutine and a device handle open contradicts both the no-leak bar and *"it starts instantly"*, for nothing anyone can hear. FX need no equivalent: they hold nothing between cues |

**A cue is always the second signal, never the first.** The visible change happens on the frame it
happens on; the sound follows or does not. Nothing waits for audio.

---

## 7 · Licensing — unresolved, and it is a founder decision

**Bundled music must be DRM-free and licensed for commercial redistribution, or it does not ship.**

The constraint is sharper here than for most projects, and both halves bite:

- **The repository is public and open-source.** Bundled tracks are redistributed by every clone,
  every fork and every release artifact.
- **The product is commercial.** The funding model is affiliate commission, so a
  "non-commercial use" licence does not cover it — the same trap already identified for IGDB
  ([`06-data-seams.md`](./06-data-seams.md) §3).

**Do not assume a track is usable.** A licence that permits "free use with attribution" frequently
does not permit redistribution inside a binary, and rarely contemplates commercial use.

| Needs a decision | Note |
|---|---|
| Which tracks, under which licence | Named, with the licence text archived in-repo |
| Whether tracks are **bundled** or **user-supplied** (point Zerado at your own directory) | User-supplied sidesteps licensing entirely and is worth considering as the Phase 1 answer |
| Attribution surface | Where the credit renders — `Z-10 Help` is the natural home |

> **Interim recommendation, for the founder to accept or refuse:** ship Phase 1 with **interface FX
> only** — short cues that are cheap to originate or to license cleanly — and make the **music bed
> user-supplied**, pointed at a local directory. That delivers the feeling, removes the licensing
> blocker from the critical path entirely, and leaves a bundled soundtrack as a later addition once
> the rights are actually cleared. This is a recommendation, not a decision; the founder's
> instruction was that audio ships, and it ships either way.

---

## 8 · What Zerado sounds like — the bar, so it can be failed

The aesthetic objection in the original verdict was the one real judgement call, and the founder has
made it. But the brand manual's bar still governs *what the sound is*, and it needs stating or the
reversal quietly becomes the thing §1 of the brand manual rules out.

**Retro-*future*, never retro-nostalgia — in sound as in everything else.**

| On register | A MISS |
|---|---|
| A machine acknowledging an instruction | A game console jingle |
| Mechanical, precise, short (brand manual §7.2: *"no bounce, no elastic, no spring overshoot"* — the audible sibling) | Chiptune pastiche, coin-drop, fanfare |
| The KITT cockpit: a tone that means *received* | A soundtrack that performs enthusiasm |

**A `ZERADO` cue that celebrates is the audible equivalent of an exclamation mark** — and the voice
section bans those for a reason. The brand manual already writes that moment: *"Zerado. 41 hours.
Sixth this year."* Three facts, landing. The sound should land the same way.

---

## 9 · Accessibility — the co-render rule extended to a fourth channel

> **Audio is never the only carrier of information.**

Every event that makes a sound is **already** visible without it. The sound is redundant by design,
exactly as colour is redundant to glyph and label.

| Event | Visible carrier (always) | Cue (only if enabled) |
|---|---|---|
| Sync complete | `Z-03 DONE` readout with the counts | `sync.done` |
| Sync failed | The refusal copy, naming what and why | `sync.failed` |
| Marked `ZERADO` | The row's state chip changes — colour **and** glyph **and** label | `status.zerado` |
| Error | The error state, in words | `error` |

The test: **run the whole product with `ZERADO_NO_AUDIO=1` and lose no information.** That is the
same test `NO_COLOR` has to pass, and it is the same reason.

---

## 10 · Environment and degrade

`ZERADO_NO_AUDIO` mirrors the `NO_COLOR` discipline already in the brand manual §5.4: **when set,
Zerado emits no sound at all** — not a reduced set, none. Settings shows the setting as
**overridden**, not as off, because those are different facts.

| Condition | Detection | Behaviour |
|---|---|---|
| `ZERADO_NO_AUDIO` set | env, any non-empty value | Silent. Settings: *overridden by the environment* |
| No audio device | device init fails | Silent. Settings names it |
| Over SSH | `$SSH_TTY` / `$SSH_CONNECTION` present | Silent **by default** — the sound would play on the wrong machine. Overridable, because a forwarded session is a real case |
| In CI | `$CI` present | Silent |
| Init failed for any other reason | timeout or error | Silent. Settings names it |

**It never errors, never blocks, never hangs, and never nags.** It goes silent and says so **once**
— a single honest line in `Z-09 Settings`, not a banner on every screen and not a repeated notice.
Audio being unavailable is not a degrade of the product; the product does not depend on it.

> Note that audio has **no degrade banner**, deliberately. The banner exists for *"what you are
> looking at is stale or incomplete"*. Silence is neither.

---

## 11 · Surfaces this adds

| Surface | Change |
|---|---|
| **Global key `m`** | Toggle mute. Added to the global key map ([`04-navigation-and-focus.md`](./04-navigation-and-focus.md) §3) and the footer — **but only when audio is enabled**; there is nothing to mute otherwise |
| **Status bar indicator** | `▮ AUDIO` (amber) / `▯ MUTED` (chrome) / nothing when never enabled. Both glyphs are **Neutral** width, verified — unlike `♪` U+266A, which is Ambiguous and would shear the bar |
| **`Z-09 Settings § Audio`** | The opt-in, the two channels, two volumes, and the honest availability line |
| **`Z-01 First run`** | One quiet line naming that audio exists |
| **`Z-10 Help`** | `m` in the key map; the attribution surface if tracks are bundled |
| **`Z-11 Fatal error`** | **Silent, always.** `Z-11` depends on nothing, and that includes audio. A crash screen must not try to play a sound |
