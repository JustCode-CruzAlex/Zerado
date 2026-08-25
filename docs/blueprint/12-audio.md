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
rejected was a *network streamer, always on, bundled into the product's identity*. What ships is an
**opt-in subsystem, off by default**, whose music is somebody else's radio station that the player
chooses and can stop in one keystroke, and whose only always-available part is a handful of local
interface cues. Three of the five objections retire because of that, not because they were
overruled.

*(The design moved twice: the first reversal specified bundled tracks; founder direction on
2026-08-25 then removed bundling entirely in favour of streamed stations. Both moves are recorded
rather than folded, because the second one **dissolved** the licensing question rather than
answering it, and that is worth being able to see.)*

The five reasons are not deleted. Two survive the reversal as **binding constraints**, and one of
those is still an open founder decision:

| Original objection | Status after the reversal |
|---|---|
| Contradicts *"no telemetry"*, *"works with the network off"*, *"only traffic is services you connected"* | **Retired — the object changed.** A station the player chose and can stop in one keystroke is a *connected service*, not background telemetry; and *"works with the network off"* is a promise about the **library**, not about every feature. §3 |
| Costs the pure-Go single-binary distribution (cgo) | **Survives as a requirement** — §5 |
| Costs the 60 fps and zero-goroutine-leak budget | **Survives as a requirement** — §6 |
| Acquires a music-rights surface | **DISSOLVED.** Nothing is bundled, so there is nothing to license — §7 |
| Nostalgia-kitsch versus retro-future | **The founder's call, and it is made.** §8 records the bar the sound has to clear so the reversal does not quietly become the thing the brand manual rules out |

> **Provenance.** The direction reached this session as an agent relay, which by its own header
> carries no ratification authority. The work is done as directed because it is document revision on
> a draft PR — reversible, and the founder confirms or restores it at the gate. It is recorded this
> way so that confirming it is a decision rather than an assumption.

---

## 1 · What it is

> **Founder, 2026-08-25:** *"we can have radio stations that play synthwave 24/7, and others that
> play only 80's… To me the audio is done."* And decisively: *"let's skip the bundle music, if the
> user is offline no music, that's it. Easy peasy, and we don't worry about that."*

Two independent channels:

| Channel | What it is | Online? | Default |
|---|---|---|---|
| **Radio** | Streamed stations — synthwave 24/7, 80s | **Yes.** Offline it stops, and that is fine | **off** |
| **Interface FX** | Short local cues: a completed sync, an error, a state change to `ZERADO` | No — local, always available | **off** |

**Nothing is bundled.** That single decision removes an entire class of problem: no licensing
question, no attribution surface, no repo weight, and no argument about which tracks. The music is
somebody else's stream, played when there is a network and silent when there is not.

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

## 3 · Streamed, never bundled — and the promises still hold

The published page says *"no telemetry running in the background"*, *"Runs with the network off"*,
and *"The only network traffic is `Zerado` reaching out to the services you've connected."*

**A radio stream does not violate any of those, and the reason is the word *connected*.** The player
turns audio on and chooses a station; from then on the stream is a service they connected, exactly
like Steam. It is not background telemetry, because the player started it and can stop it in one
keystroke. It is not traffic they did not ask for.

*"Runs with the network off"* is the promise the earlier draft read too strictly. It means **the
library works offline** — see [`07-offline-contract.md`](./07-offline-contract.md) §1. A stream that
stops when the network does is an obvious, honest degradation of a feature that is online by nature,
not a broken promise.

So: **interface FX are `WORKS`** — local, always. **Radio `NEEDS THE NETWORK`** — offline it stops,
says so where it lives, and there is nothing to fall back to because nothing is bundled.

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
| **Muting music RELEASES the device AND the connection** | Muting the music channel **halts decode, releases the audio device, and closes the stream** — it does not hold a gain of zero, and it does not keep pulling bytes. Unmuting reconnects. FX need no equivalent: they hold nothing between cues |

> **The connection half was added when music became a stream, and it matters more than the device
> half.** The original rule was written when music was a local bed, where "release" could only mean
> the audio device. A muted Zerado that keeps pulling bytes from someone else's server is **spending
> the player's bandwidth on silence** — and it sits badly against the published promise that *"the
> only network traffic is `Zerado` reaching out to the services you've connected"*: the player
> connected that station, and has just told Zerado to stop it.
>
> Reconnecting on unmute costs a moment of buffering. That is the correct trade: a stream the player
> has muted is a stream they are not listening to, and the honest thing is to stop taking it.
>
> *(Raised by `fft-design-architect`, which noticed that the spine had adopted its device-half rule
> verbatim and that the rule's scope had moved underneath it.)*

**A cue is always the second signal, never the first.** The visible change happens on the frame it
happens on; the sound follows or does not. Nothing waits for audio.

---

## 7 · Stations are data — and the licensing question is closed

**The music-licensing gate item is closed, by removing its cause.** Nothing ships in the binary, so
there is nothing to license, nothing to attribute and nothing to argue about. This was the most
expensive open question in the audio design and it was dissolved rather than answered.

**The default station list ships as data, not compiled in**, and is user-editable — a plain file the
player can add to, reorder or replace. Two rules:

1. **Every URL in the default list must be verified to resolve before it ships.** A dead station in
   the default list is a broken first impression, and it is the kind of rot that arrives silently
   months later. The list needs a check that can be re-run, not a one-time look.
2. **A station that fails to play is a station, not an error.** It reports that it could not connect
   and the player picks another. It never becomes a modal, and it never stops the library.

Zerado streams; it does not host, cache or redistribute. That keeps it clearly on the side of a
client, which is the same posture any podcast or radio client takes.

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
| **`Z-10 Help`** | `m` in the key map. **No attribution surface** — nothing is bundled, so there is nothing to attribute |
| **`Z-11 Fatal error`** | **Silent, always.** `Z-11` depends on nothing, and that includes audio. A crash screen must not try to play a sound |
