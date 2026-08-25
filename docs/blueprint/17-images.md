---
title: Zerado — images in the terminal
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-17
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: concept-explainer
ticket: "#2"
---

# Images in the terminal

> **Founder, 2026-08-25:** *"Images on terminal is a must, I know we will rely on kitty/ghostty but
> we need to show images, without image is not an option. Starting on phase 1."*
>
> And on the degrade, asked directly: *"yes it should work on other terminal emulators that will not
> have images, so we show what we can, and a text saying to use ghostty or kitty to better
> experience."*

**Cover art is foundational, not an enhancement.** This reverses revision A, which put the cover
deck in Phase 2 and explicitly declined to assume inline-image support.

---

## 1 · What changed, and the honest reason it was wrong before

Revision A's reasoning was: terminal image support is not universal, so do not depend on it. That is
sound risk-thinking and it produced the wrong answer, because it optimised for the terminal instead
of for the player. A game library without cover art is a spreadsheet — and the product's own pitch
is *"what do I play tonight?"*, a question people answer with their eyes.

The correct shape was available the whole time and revision A did not reach for it: **support the
capability where it exists, degrade honestly where it does not, and let the player know a better
experience is one terminal away** — without ever making that a condition of use.

---

## 2 · The protocol decision

**Target the Kitty graphics protocol.** Supported by **Kitty** and **Ghostty**, the two the founder
named, and by WezTerm and Konsole.

| Protocol | Decision | Why |
|---|---|---|
| **Kitty graphics** | **Adopt** | The most capable and the most widely implemented of the modern protocols; covers both named terminals |
| **iTerm2 inline images** | **Adopt, second** | It is the macOS default terminal for a large part of this audience and the protocol is small. The cost is one encoder |
| **Sixel** | **Do not adopt in Phase 1** | Broad but old, poor colour fidelity on photographic content, and slow enough at deck sizes to threaten the redraw budget. Revisit if a real user is stranded on a Sixel-only terminal — that is a data question, not a design one |
| **Half-block / ASCII art** | **Reject** | It is not the image; it is a picture of a picture, and it fails the brand's own "retro-future, not retro-kitsch" bar |

**Capability detection at startup, never a config flag the player must find.** Query the terminal
and cache the answer for the session. A player should never have to learn the word "protocol" to see
a cover.

Detection must **fail closed**: an ambiguous or timed-out response means *no image support*, and the
text deck renders. Guessing yes and emitting escape sequences into a terminal that does not
understand them is how a library view turns into garbage on somebody's screen.

---

## 3 · The degrade — supported, never refused

**A terminal without image support is a supported configuration.** Not a warning state, not a
degraded tier, not something to apologise for.

| | |
|---|---|
| What renders | The **text deck** — the full library, every state, every column. *"We show what we can"* |
| What the player is told | **Once**, quietly: a dismissible note that Ghostty or Kitty would show cover art |
| What it never is | Recurring, blocking, modal, or phrased as an error |

**The tone is the deliverable here.** *"Your terminal does not support images"* reads as a fault in
the player's setup. *"Cover art needs Ghostty or Kitty"* reads as information. The screen spec
([`../design/screens/Z-15-cover-deck.md`](../design/screens/Z-15-cover-deck.md)) writes the exact
copy; the spine's requirement is that it is **dismissible and never returns** once dismissed, which
means the dismissal is persisted like any other setting.

---

## 4 · The seam

```go
// Images is the whole surface. A screen asks for a cover at a size; it never
// learns which protocol answered, or whether one did.
type Images interface {
    // Capability is resolved once at startup and cached for the session.
    Capability() ImageCapability   // none | kitty | iterm2

    // Cover returns a placement for an already-cached image, or ok=false.
    // It NEVER fetches and NEVER blocks: a cover the cache does not hold is
    // simply not shown this frame.
    Cover(id ItemID, cells Rect) (Placement, bool)

    // Warm asks the cache to fetch, off the render path, best-effort.
    Warm(ctx context.Context, ids []ItemID, cells Rect)
}
```

Two properties carry the whole design:

- **`Cover` never blocks and never fetches.** The render path only ever reads what is already
  present. This is the same rule as `Audio.Cue` and for the same reason: **a missing cover is never
  worth a dropped frame.**
- **A screen never learns the protocol.** Adding Sixel later, or dropping iTerm2, changes one
  implementation and no screen.

`NullImages` reports `none` and answers `false` to everything. It is the implementation used when
detection fails, in tests, and under `ZERADO_NO_IMAGES` — so the text path is the well-exercised one
rather than an untested fallback, exactly as `NullAudio` is.

---

## 5 · Cache, and where it sits in the truth split

Cover art is **cache, never truth** — it belongs in the XDG **cache** directory alongside the
existing image cache decision ([`09-erd.md`](./09-erd.md) §3), not in `library.db`.

That placement is already argued and does not change: covers are disposable, refetchable and
potentially large; the OS is allowed to delete them; and the file the player backs up stays small
enough that they actually back it up. **Deleting the whole cache must cost nothing but bandwidth.**

| Rule | |
|---|---|
| Fetched | online, off the render path, best-effort |
| Stored | XDG cache, keyed by item, at the sizes actually rendered |
| Offline | cached covers render normally; uncached ones show the **designed no-cover tile**, never a broken-image box |
| Never | in `library.db`, and never blocking a frame |

In the offline contract this makes cover art **`DEGRADES`**: it has local data, that data may be
incomplete, and the interface says so by showing a designed tile rather than a gap.

---

## 6 · The performance bar

The same bar as everything else — **60 fps and no goroutine leaks** — and images are the classic way
to lose both.

| Rule | Mechanism |
|---|---|
| Never blocks the render loop | `Cover` reads the cache; `Warm` is the only thing that touches the network, and it is off the render path |
| Never leaks | Fetching is one owned worker pool bounded by a `context`; the seam owns it, no screen starts one |
| Never re-transmits | Kitty placements are addressed by image id — transmit once, place many times. Re-sending image data every frame is the failure mode that makes terminal images feel slow |
| Never stutters on scroll | Only visible tiles are placed; `Warm` runs ahead of the viewport, not across the whole library |
| Bounded | The cache has a size ceiling and evicts least-recently-shown |

---

## 7 · What Phase 1 builds

| Built | Not built |
|---|---|
| Capability detection, failing closed | Sixel |
| Kitty + iTerm2 placement | Image editing, cropping or effects |
| The cover cache, bounded and evicting | Cover art for hand-added physical copies where no source exists — that is a designed empty, not a gap |
| `Z-15` cover deck as a mode of `Z-04` | A separate cover-art screen |
| The dismissible, once-only recommendation note | Any blocking or recurring prompt |

**Where the covers come from is the metadata seam's question, not this one**
([`06-data-seams.md`](./06-data-seams.md) §3) — and that seam is deliberately provider-agnostic,
which matters more now that images are foundational: an image pipeline with one hard-coded source is
a pipeline with a single point of failure the product cannot survive.
