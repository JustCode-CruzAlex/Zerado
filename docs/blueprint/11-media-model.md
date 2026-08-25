---
title: Zerado — the media model
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-11
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: concept-explainer
ticket: "#2"
---

# The media model — games first, not games only

> **Founder direction, 2026-08-25:** *"later we will add BOOKS to Zerado, which will work the same
> way — books that people own, indication, books to read by mood or state of the reader. Not now,
> but the database can be prepared for when that happens, maybe even MOVIES and TV SHOWS after."*

**Phase 1 ships games and only games.** No book screens, no book providers, no book commands, no
`--type` flag. What changes is the **shape**, so that adding books later is a new type rather than
a rewrite.

This is the most expensive thing in the bundle to retrofit, which is why it is designed now and
built later.

---

## 1 · The core entity is a **media item**, and `game` is its first type

```
media_item                          ← everything genuinely shared lives here
  ├── media_game     (Phase 1)      ← Steam appid, achievements
  ├── media_book     (Phase N)      ← author, ISBN, page count, publisher
  ├── media_film     (Phase N)      ← runtime, director
  └── media_series   (Phase N)      ← §4 explains why this is not the same type as film
```

### What is shared, and therefore lives on the core

Ownership · acquisition (**digital or physical**) · the four states · progress · mood tags · the
player's rating · the player's notes · the source provider · price history · everything the
recommender reads.

### What is type-specific, and therefore lives in a typed extension

| Type | Facts that are genuinely only true of it |
|---|---|
| `game` | Steam appid, achievements unlocked / total |
| `book` | author, ISBN, page count, publisher, format (print · ebook · audio) |
| `film` | runtime, director |
| `series` | seasons, episodes released, whether the series has ended |

### One deliberate refinement of the direction, with the reason

The direction lists **`playtime_hours` as a game fact**. This blueprint models it as **generic
progress on the core** instead:

```
progress_value   INTEGER NULL     -- 41
progress_unit    TEXT             -- minutes | pages | percent | episodes
progress_source  TEXT             -- derived | manual
```

For a game the unit is `minutes` and the value *is* the playtime — nothing is duplicated, and
`media_game` keeps only what is genuinely game-only.

The reason this is worth diverging on: the direction also requires that **`IN PROGRESS` be derived
for games and settable for books.** If progress is a game column, that rule has to be re-implemented
per type, and the four-state derivation stops being one function. With progress on the core, the
derivation is written **once** and reads a capability. §3 shows it.

If the founder prefers playtime as a typed field, it is a small change — but the derivation then
forks per type, and that fork is the thing worth avoiding.

---

## 2 · The four states generalise — verified, type by type

The direction asked for this to be checked rather than assumed. It was checked, and the answer is
**yes for all four states across all four types**, with two findings in §4.

| State | `game` | `book` | `film` | `series` |
|---|---|---|---|---|
| `○` **NOT STARTED** | owned, never launched | owned, never opened | owned, never watched | owned, never started |
| `◐` **IN PROGRESS** | playing it | reading it | *see finding F-1* | watching it — **its most-used state** |
| `◉` **ZERADO** | beaten, cleared | **finished the book** | watched to the end | watched to the end of what exists |
| `⊘` **ABANDONED** | stopped, not going back | did not finish | stopped watching | dropped |

**`ZERADO` is the strongest evidence the vocabulary is right, not the weakest.** The founder's own
test — *"a book you finished IS zerado"* — holds because the word was never about games. It is
Portuguese for *beaten, cleared, closed out*, and it came from an arcade counter rolling over
(`naming.md`). A finished book is closed out. So is a finished series. The name generalises for the
same reason it was chosen.

`ABANDONED` generalises best of all: *did not finish* is a more familiar idea for books and series
than it is for games.

---

## 3 · Derived versus manual is a **capability**, not a type

This is the mechanism that makes the direction's requirement — *derived for games, settable for
books* — cost nothing.

Derivation is **not** keyed on `media_type`. It is keyed on the **(provider, type) pair's declared
capability**, which is where the truth actually lives:

| Item | Provider reports progress? | `IN PROGRESS` is |
|---|---|---|
| A Steam game | yes, minutes played | **derived** |
| A cartridge on a shelf (`physical`) | no — a cartridge has no telemetry | **manual** |
| A paper book | no — there is no page counter | **manual** |
| A Kindle book, *if* a provider is ever written for it | yes, percent read | **derived** |
| A series tracked by a provider | yes, episodes watched | **derived** |
| A series tracked by hand | no | **manual** |

Two things follow, and both matter:

1. **The same type can be derived for one item and manual for another.** A Steam game and a
   cartridge are both `game`; one derives and one does not. **Phase 1 already needs this**, which
   means the mechanism the founder asked for is not new work — it is the mechanism physical copies
   already forced.
2. **Books are not a special case.** A paper book is manual for exactly the reason a cartridge is
   manual: nothing is reporting. The generic rule already covers it.

The derivation, written once for every type that will ever exist:

```go
func derive(p Progress, c Capabilities) Status {
    if !c.Progress {          // nothing is reporting — the player owns the state
        return StatusNotStarted
    }
    if p.Value > 0 {
        return StatusInProgress
    }
    return StatusNotStarted
}
```

`Capabilities.Playtime` from [`06-data-seams.md`](./06-data-seams.md) is renamed
**`Capabilities.Progress`**, and gains a `ProgressUnit`. That is the whole change.

---

## 4 · Two findings the check surfaced

The direction asked: *"If a type needs a state the others do not, that is a finding worth naming."*
Two, and neither requires a fifth state.

### F-1 · A film has no meaningful `IN PROGRESS` — which means **film and series are different types**

A 110-minute film is watched or it is not. "In progress" for a film means *"I paused it"*, which is
a **session**, not a state. For a **series**, `IN PROGRESS` is the most-used state in the type.

If film and series were one type, one of them would carry a state it cannot use and the other would
carry a progress unit it does not have (`minutes` versus `episodes`).

> **Recommendation: model `film` and `series` as two types when they arrive, not one "video" type.**
> Costs nothing now. Costs a migration later.

### F-2 · An ongoing series needs *"caught up"*, and it must **not** become a fifth state

You have watched every episode released, but the series has not ended. That is not `ZERADO` — it is
not finished. It is technically `IN PROGRESS`, but the player is not mid-anything; they are waiting.

**Do not add a fifth state.** The four are ratified, CVD-verified at a measured ΔE floor of 11.9
under deuteranopia, and designed around a co-render rule that a fifth colour would have to be
re-verified against — for one type, in a later phase.

> **Recommendation: model *caught up* as a derived presentation of `IN PROGRESS`**, computed from
> typed facts in `media_series` (`episodes_watched == episodes_released && !series_ended`). The
> state stays `◐ IN PROGRESS`; the detail view says *"Caught up. Next episode not out yet."*

Naming this now is the point. It is exactly the pressure that would otherwise arrive in Phase N as
*"can we just add one more state?"* — and the answer needs to be already written down, with the CVD
cost attached.

---

## 5 · The provider seam generalises without reshaping

A provider declares the media types it serves:

```go
type Provider interface {
    ID() ProviderID
    Display() string
    MediaTypes() []MediaType     // steam → [game] · openlibrary → [book] · tmdb → [film, series]
    Capabilities(MediaType) Capabilities
}
```

`Capabilities` is per **(provider, type)** because one provider can serve two types with different
abilities — TMDB reports episode counts for a series and a runtime for a film.

The streamed `Item` carries its type and its typed extension:

```go
type Item struct {
    Type        MediaType
    ProviderRef string
    Title       string
    Progress    *Progress        // nil = not reported
    Acquisition Acquisition      // digital | physical
    Extra       TypeExtension    // the typed payload; opaque to the core
}
```

**Nothing in the Phase 1 seam is reshaped by this** — `MediaTypes()` is added, `Playtime` is renamed
to `Progress`, and `Item` gains two fields. Steam declares `[game]`; `physical` declares `[game]`
today and will declare `[game, book, film]` the day physical books matter, **without a new
provider**, because a shelf is a shelf.

That last point is the cleanest evidence the seam is right: the provider that most obviously
generalises is the one this ticket added.

---

## 6 · Mood tagging is type-aware without duplicating the engine

*"Brain-dead after work"* means a different thing for a game, a book and a film. But it is **the same
axis**: how much attention this asks for.

So mood tags carry a **type-neutral key** and a **per-type label**:

| `key` (the engine reasons over this) | Label for `game` | Label for `book` | Label for `film` |
|---|---|---|---|
| `low_attention` | Mindless grind | Light, short chapters | Something you can half-watch |
| `story_heavy` | Story rich, kind of sad | A novel to sit inside | A film to feel |
| `short_session` | Quick fifteen minutes | A chapter before bed | A short |
| `full_focus` | Tactical, full focus | Dense, worth the effort | Demands the room dark |

```
mood_tag ( id, key, applies_to[], label )
media_mood ( media_item_id, mood_id, source: user|inferred, confidence )
```

**One engine, one axis vocabulary, per-type labels.** The recommender never learns what a book is —
it matches on `key`. The interface never shows a game a book's words.

The Phase 1 labels are already published in `content/landing-copy.md` §04 and are binding; the
`key` column is what lets them stay exactly as written while a second type gets its own.

---

## 7 · What Phase 1 actually builds

To be unambiguous, because "prepare the database" is the sentence most likely to turn into scope:

| Built in Phase 1 | Not built in Phase 1 |
|---|---|
| `media_item` with `media_type`, constrained to `'game'` | Any other type value |
| `media_game` | `media_book` · `media_film` · `media_series` |
| Generic `progress` (value · unit · source) | Any non-`minutes` unit |
| `MediaTypes()` on the provider seam, returning `[game]` | A books or film provider |
| `mood_tag.key` + `applies_to`, all rows `['game']` | Per-type label sets |
| A `media_uid` that is type-scoped | — |

**No screen changes. No new commands. No `--type` flag.** A Phase 1 player cannot tell this document
exists, and that is the correct outcome.

The check that it worked is a thought experiment worth writing down: *adding books should be one
migration, one provider, one typed extension, and one label set — and zero changes to the state
machine, the recommender, the offline contract, or any existing screen.* If a future book ticket has
to touch those, this design failed and the failure is measurable.
