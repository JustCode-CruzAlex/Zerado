---
archetype: adr-detail
adr: ADR-0001
title: "Zerado foundational architecture — provider seam, persistence, navigation, Phase 4 sync, the open door, audio, themes, images, i18n"
status: PROPOSED
date: 2026-08-25
ticket: "#2"
---

# ADR-0001 · Zerado foundational architecture

Nine decisions that are expensive to reverse once Phase 1 code exists. They are recorded together
because they constrain each other: the sync boundary decides what persistence must carry, and
persistence decides what the provider seam can promise.

**Drawings for this ADR** — ten sheets, each rendered in both themes, under
[`charts/`](./charts/):

| Sheet | Drawing | Bears on |
|---|---|---|
| 01–02 | The data model | D2, D4 |
| 03 | First run → a library | D3 |
| 04 | The four states | D4 |
| 05 | Sync, with its failure branches | D1 |
| 06 | Navigation | D3 |
| 07 | The provider seam | D1 |
| 08 | The offline contract | D1, D2 |
| 09 | Screen anatomy at the 80 × 24 design floor | D3 |
| 10 | The audio subsystem and its degrade ladder | D6 |

*(The drawings predate D7–D9; themes, images and i18n are specified in prose and carry no sheet. Recorded rather than left for a reader to notice.)*

---

## Context

Zerado is a terminal-first game library, greenfield as of 2026-08-25 — no Go code exists. The
stack is settled: Go, Bubble Tea, and the Charm ecosystem. Four phases are published and binding
(`content/landing-copy.md` §12), as are six ratified promises (`ratification/decisions.md`):
local-first, one SQLite file, no Zerado-run server before Phase 4, the player's own API keys, an
**donation-only funding with no affiliate commission**, and no named community source.

Those promises are not context. They are **architectural constraints**, and each decision below
names the ones it is paying for.

Two facts sharpen the problem:

1. **Physical copies must be first-class from day one.** The published copy says *"A physical copy
   isn't a second-class row in the list."* The obvious implementation — a Steam-shaped row with an
   `is_physical` flag — makes that sentence false in the second week.
2. **IGDB is free for non-commercial use only.** Zerado's funding model was affiliate commission —
   which is commercial — until founder direction on 2026-08-25 **dropped affiliate links entirely**,
   making the product cleanly non-commercial: free software, donation-supported, **zero revenue**.
   This premise is load-bearing: D1's Consequences and
   [`../blueprint/06-data-seams.md`](../blueprint/06-data-seams.md) §3 close the IGDB question
   *because* it is now false. It must not be restated in the old tense. Cover art and *sinopse* are the visual backbone of Phase 2.
   Designing as if IGDB is guaranteed would be designing on a maybe.

---

## D1 · The provider seam — segregated interfaces, capability-driven screens

### Decision

Two interfaces, not one:

```go
type Provider interface {
    ID() ProviderID
    Display() string
    Capabilities() Capabilities
}

type Syncer interface {
    Provider
    Sync(ctx context.Context, c Credentials) (<-chan Item, error)
}

type Capabilities struct {
    Sync, Playtime, LastPlayed, OwnedSince bool
    Credentials []CredentialField
}
```

- `steam` implements **both**. `physical` implements **`Provider` only** — it is not a Syncer with
  a hole in it; its sync is a human typing into `Z-08`.
- **Everything downstream reads `Capabilities`, never `ProviderID`** — the row renderer, the state
  derivation, the filter, the summary counts, the detail view.
- `Z-02 Connect a store` renders its form from `Capabilities().Credentials`. Adding GOG is:
  implement the interfaces, declare the fields, register. **Zero screens, zero routes, zero schema
  changes.**
- `Sync` **streams**. A cancel mid-sync leaves a valid partial library; memory stays bounded on a
  library of any size; `Z-03` can show honest running counts.
- The metadata and price seams follow the same shape, and both carry an **`Attribution()`
  method** — the credit a source requires is a property *of the source*, so swapping the source
  swaps the credit. **`Quote` carries no affiliate URL** — the link is a plain shop link, because
  affiliate commission is dropped and Zerado earns nothing from a purchase. Every quote does carry
  its **`ObservedAt`**, mandatory and rendered, because a price shown without its age is the product
  giving financial advice from memory.

Pays for: *"A physical copy isn't a second-class row"* and the IGDB risk.

### Alternatives

| Considered | Rejected because |
|---|---|
| **One `LibraryProvider` interface; `physical` returns `ErrManual` from `Sync`** | Every caller must then know which providers lie about their own interface. The error is a runtime rediscovery of a fact the type system could have carried |
| **`is_physical bool` on a Steam-shaped row** | Makes the published promise a convention rather than a structure. Every new feature has to remember the special case, and eventually one will not |
| **A plugin system (subprocess or Go plugins)** | Solves a problem nobody has. Four stores are known and they ship in the binary. `plugin` does not cross-compile, and a subprocess protocol is a network boundary in disguise |
| **Returning `[]Item` instead of a channel** | Cannot show honest progress, cannot survive a cancel with a usable result, and holds a whole library in memory for no gain |

### Consequences

**Easier:** a new store is one package. Physical copies are correct by construction. The offline
contract is structural, because there is exactly one place network I/O can originate. Testing a
screen needs no provider at all.

**Harder:** the upsert must be transactional per *batch* rather than per sync, which makes
"partial" a real state the UI has to name — `sync_run.status` and `Z-03 PARTIAL` exist for this.
Interface segregation also means two names to learn instead of one.

**Cost to reverse:** **high.** Every screen, every sync path and the schema key off this shape.

---

## D2 · Persistence — one SQLite file, a pure-Go driver, credentials outside it

### Decision

1. **One SQLite file**, at `$ZERADO_DB` ?? `$XDG_DATA_HOME/zerado/library.db`. All data access is
   behind a single `Store` interface — the only thing in the program that knows a database exists.
2. **`modernc.org/sqlite`** (pure Go), not `mattn/go-sqlite3` (cgo).
3. **Credentials are not in that file.** They go to a `Vault`: the OS keychain where one exists,
   otherwise `credentials.json` mode `0600` beside the library — and `Z-09 Settings` states which
   backing is in use.
4. **Cover-art blobs are not in that file either.** They go to the XDG *cache* directory.
5. `effective_status` is **not stored**; it is derived on read.
6. Migrations are **forward-only**, idempotent, transactional. A database written by a newer binary
   is a fatal error that names both versions — never a silent downgrade.

### Alternatives

| Considered | Rejected because |
|---|---|
| **cgo SQLite (`mattn/go-sqlite3`)** | Faster writes, smaller binary — both irrelevant at a few hundred rows. It costs the single-binary cross-compile matrix, which is the whole difference between "download one file and run it" and "install a build environment." The page's claim is *"It's a text program. It starts instantly"* |
| **Credentials inside `library.db`** | Simpler, and it puts a secret in a file the player is explicitly invited to back up, move and delete. A key in every backup and every support-ticket attachment |
| **Keychain only, no file fallback** | Breaks on headless Linux and in containers, which is a real part of a terminal-first audience |
| **A separate cache database** | Breaks *"one file"* in the only way the player would notice: two files to copy |
| **A stored `effective_status` column** | A derived value with a second writer. Two ways to be wrong instead of none, for a computation over integers already in memory |
| **Plain files / JSON** | Loses transactions, indexes and `GROUP BY` on a 400-row list that must not lie about its own counts |

### Consequences

**Easier:** one downloadable binary per platform, no toolchain. The backed-up file stays small and
carries no secret. `Store` is the seam Phase 4's sync engine attaches to without touching a screen.

**Harder:** *"one file"* and WAL mode must be reconciled deliberately — `-wal` and `-shm`
companions exist while the process runs and must be checkpointed away on clean shutdown. Two
credential backings means two code paths and an honest line in Settings.

**Cost to reverse:** **high for the file layout, low for the driver.** The driver sits behind
`database/sql` and `Store`; swapping it is a day. Moving credentials into or out of the library
file after people have libraries is not.

---

## D3 · Navigation — a route stack with one overlay slot, not tabs

### Decision

- **A route stack.** `Z-04 Library` is the root and can never be popped. Pushing a route already
  on the stack **unwinds to it** rather than duplicating it.
- **At most one overlay**, in its own slot above the stack. An overlay may not open another
  overlay; a two-step flow is a route.
- **At Tiny (<40 columns) an overlay becomes a route** — identical behaviour, different
  composition.
- **`Esc` = back one level, everywhere**, with one deliberate two-step case (filter mode: the
  first `Esc` leaves the editor and keeps the filter, the second clears it — and the footer says
  so at that moment).
- **`q` = quit, from anywhere, immediately.** No confirmation, because every mutation commits when
  it is made.
- **`:` and `Ctrl-K` are reserved and unbound** in Phase 1 for the Phase 2 command palette.
- **Focus is carried by three channels** — position (a `▌` marker), weight (bold), colour (amber)
  — and any two suffice, so it survives `NO_COLOR`. **The focus ring is never removed.**
- **Focus is restored on pop**, and preserved across row-set rebuilds by game identity rather than
  index.

### Alternatives

| Considered | Rejected because |
|---|---|
| **A tab bar** | Implies a set of peers. Settings is not a sibling of the library; it is a place you go and come back from. It also spends a permanent row of a 24-row terminal on navigation used a few times a session |
| **A modal command palette as the primary navigation, from Phase 1** | A palette earns its place when the surface is bigger than the key map. With twelve screens and one home, `?` plus a footer hint is better — and reserving the keys costs nothing |
| **`q` = back, `Ctrl-C` = quit** | Two keys that both mean "leave", disagreeing about how far. `Esc` already means back |
| **A quit confirmation** | Nothing to protect, and it trains people to dismiss confirmations — the exact reflex you need them not to have when something destructive appears |

### Consequences

**Easier:** one mental model; `Esc` is answerable in one table; deep composition is impossible by
construction (no overlay-on-overlay).

**Harder:** no direct jumps — you leave a form before going to settings. That is a real cost, and
the Phase 2 palette is the intended fix. Deep-linking (`zerado game 42`) requires the stack to be
constructible from a route descriptor, not only by pushing.

**Cost to reverse:** **medium.** The code cost is contained; retraining players who have built
muscle memory is the real expense, which is why `:` and `Ctrl-K` are reserved now.

---

## D4 · The Phase 4 sync boundary — only what the player typed crosses

### Decision

**Crosses:** `status_manual` + `status_changed_at`; user-assigned mood tags; manually-entered
(`physical`) games; and a stable `game_uid`.

**Never crosses:** the library itself · `playtime_minutes` · `last_played_at` · cover art ·
*sinopse* · prices · **credentials, ever**. Each device connects with its own keys.

The rule: **only what the player typed crosses. Everything a machine can recompute, each device
recomputes.**

Three things this forces into the **Phase 1** schema:

1. **`game_uid`** — `uuidv5(namespace, normalise(title) + "|" + normalise(platform))`, assigned at
   insert, **indexed but not unique**. A merge *hint*, not an authority: ambiguous matches are
   shown to the player, never guessed.
2. **`merged_into`** — so two rows can be joined later without a migration that rewrites primary
   keys.
3. **`status_changed_at`** — the timestamp that resolves a conflict.

**Conflict resolution: last-write-wins on `status_changed_at`, per game.**

### Alternatives

| Considered | Rejected because |
|---|---|
| **Sync the whole library** | Contradicts local-first in spirit, makes the server expensive enough to break the donation-supported statement, and stores data the player already has two other copies of |
| **Sync credentials so a new device is zero-setup** | A ratified promise says the keys are the player's own. Centralising them makes Zerado a credential custodian, which is a different company |
| **Add `game_uid` in Phase 4 instead** | The migration would have to invent stable identities for rows whose titles the player has since edited. One column and one index now; an unwritable migration later |
| **CRDTs / an operation log** | Correct for concurrent multi-writer editing. The conflicting parties here are one person on two devices who agree about what they did. Real complexity for a case that does not exist |
| **Server-authoritative state** | Requires a server before Phase 4 to be useful, which a ratified promise forbids |

### Consequences

**Easier:** the Phase 4 server is small, cheap and holds almost nothing — which is what makes the
funding statement honest. A device rebuild needs only the player's own keys.

**Harder:** two devices that both change a status while offline will lose the earlier change
silently. That is acceptable for this shape of conflict and unacceptable to leave undocumented, so
`Z-22 Devices and sync` shows the last merge and what it resolved. Manually-entered games are the
one row class with no other copy, so they are the one class whose loss is unrecoverable — they
must be in the first sync payload, not a follow-up.

**Cost to reverse:** **highest of the four.** It decides what the schema carries from the first
migration onward.

---

## D5 · The door stays open — one column and a table name

### Decision

Zerado is a **games** product. Two decisions, both reversible, keep books possible without shaping
anything now: the core table is called **`item`**, not `games`; and it carries an **`item_type`**
column, `CHECK`ed to `'game'`.

That is the entire affordance. No typed extensions, no generic progress, no interface parameterised
on a type with one value.

**This supersedes revision A of this decision**, which specified a full media-polymorphic core with
typed extensions for books, films and series. Founder direction, 2026-08-25: *"At this point don't
even think on books and other media types. What I would like is let that door open."* The
speculation was beginning to shape Phase 1 tables and Phase 1 states, which is the cost that makes
speculative generality expensive rather than merely unused. The reasoning is kept as a short
appendix in [`../blueprint/11-media-model.md`](../blueprint/11-media-model.md), shaping nothing.

**One consequence worth naming: the generic-progress divergence is withdrawn.** Revision A modelled
playtime as `progress_value` + `_unit` + `_source`, arguing it kept the four-state derivation a
single function across media types. That argument died with the media types. `playtime_minutes` is a
plain column — which is what the original direction said, and it closes a founder-gate item by
removing its cause rather than answering it.

### Alternatives

| Considered | Rejected because |
|---|---|
| **Full polymorphic core** (revision A) | Pruned by the founder. It was shaping Phase 1 for a type that does not exist |
| **Nothing at all — call the table `games`** | The retrofit renames the central table and every foreign key referencing it, on a file the product promises never to lose. Two columns of foresight buy that away |
| **A plugin/registry for types** | An abstraction with one implementation, indefinitely |

### Consequences

**Easier:** Phase 1 reads as a games product, because it is one. **Harder:** every query carries an
`item_type` predicate it does not yet need — one cheap predicate, and the price of not doing the
retrofit. **Cost to reverse: low**, which is the point.

## D6 · Audio ships in Phase 1 — **streamed radio, off by default, fully removable**

> **This reverses a verdict recorded in this same bundle** — see
> [`../blueprint/08-prior-draft-analysis.md`](../blueprint/08-prior-draft-analysis.md) §3, kept
> visible as superseded rather than deleted. The direction arrived as an **agent relay carrying no
> ratification authority**; it was acted on because this is reversible document work on a draft PR,
> and it is flagged so the founder confirms it rather than inherits it.

### Decision

**The reversal is not a change of mind about the same object — the object changed, twice.** What was
originally rejected was a *network streamer, always on, bundled into the product's identity*. What
ships is an **opt-in subsystem, off by default**, whose music is **internet radio the player chooses
and can stop in one keystroke**, and whose only always-available part is a handful of **local**
interface cues.

*(The design moved twice and both moves are recorded rather than folded. The first reversal
specified **bundled tracks**; founder direction on 2026-08-25 then removed bundling entirely in
favour of **streamed stations**. The second move matters because it **dissolved** the licensing
question rather than answering it.)*

- **Streamed, never bundled.** Music is internet radio the player chooses — synthwave, 80s — and
  **nothing ships in the binary.** Interface FX are local and always available. Founder direction,
  2026-08-25: *"let's skip the bundle music, if the user is offline no music, that's it."*
  Radio therefore `NEEDS THE NETWORK` and FX are `WORKS`; a stream stopping when the network does is
  an honest degradation of a feature that is online by nature, not a broken promise.
- **Off by default**, with an explicit opt-in in `Z-09 Settings § Audio`. A terminal program that
  makes noise on first run is one people uninstall before they have seen anything else. This also
  satisfies **WCAG 1.4.2 Audio Control** structurally rather than by a bolted-on control.
- **Two independent channels** — music and interface FX — separately mutable, separately volumed.
- **Global `m` mute**, always reachable, in the footer and key map **only when audio is enabled**.
- **`ZERADO_NO_AUDIO`** mirrors the `NO_COLOR` discipline: when set, no sound at all, and Settings
  shows *overridden*, not *off*.
- **Never errors, never blocks, never hangs.** `Cue()` cannot fail and cannot block; a cue is
  dropped before a frame is. One owned goroutine. Muting music **releases the device** rather than
  holding a gain of zero.
- **The implementation sits behind a build tag**, so the **default** build stays pure-Go and
  cross-compiles with no toolchain — preserving D2's single-binary property. `NullAudio` is the
  default build, so it is exercised by every test run rather than being an untested fallback.
- **Audio is never the only carrier of information** — the co-render rule extended to a fourth
  channel. The test is `ZERADO_NO_AUDIO=1` and lose nothing, the same test `NO_COLOR` passes.

### The licensing question is CLOSED, by removing its cause

Revision A carried an open founder decision: bundled tracks would need commercial redistribution
rights from a public repository. **Nothing is bundled, so there is nothing to license, nothing to
attribute and no repo weight.** The most expensive open question in the audio design was dissolved
rather than answered.

**Stations ship as data, user-editable, not compiled in** — with the requirement that every default
URL is verified to resolve before it ships, and re-checkable thereafter. A dead station in the
default list is a broken first impression that arrives silently months later.

### Alternatives

| Considered | Rejected because |
|---|---|
| **Bundled tracks** (this decision's own revision A) | Needs commercial redistribution rights from a **public** repository, which is a rights surface a game tracker should not acquire — and it puts weight in the binary for a feature that is off by default. Superseded 2026-08-25 |
| **A streamer that runs unprompted, always on** (the original draft) | Contradicts three ratified promises. What ships is not this: the player turns it on and chooses the station |
| **On by default** | A terminal program that makes noise unprompted is uninstalled before it is understood |
| **One channel with one mute** | Someone may want keyclicks without the soundtrack, or the reverse. Neither is the odd request |
| **Audio compiled into the default build** | Costs D2's pure-Go cross-compile matrix — the property that makes releases one file and no toolchain |
| **`Cue()` returning an error** | There is no audio failure a *screen* could usefully handle. An error return invites a caller to block on it |
| **A second key to pause music, distinct from mute** | A redundant control on a scarce surface, for no accessibility gain — `m` already satisfies 1.4.2 |

### Consequences

**Easier:** the retro-future feeling becomes something the player turns on. `NullAudio` as the
default build means the silent path is the well-tested one.

**Harder:** a build-tag split means two binaries to produce and a Settings line that must explain
which one is running. And the **station list has to stay alive** — every default URL verified to
resolve before it ships, and re-checkable, because a dead station is a broken first impression that
arrives silently months later.

**Cost to reverse:** **low.** The seam is removable at compile time and the stations are data. There
is **no licensing commitment to undo**, because nothing ships in the binary — that was revision A's
expensive part, and removing bundling removed it.

---

## D7 · Themes are data, and every theme must pass the four-state contract

### Decision

A theme is a **data file**, never code. Zerado adopts FlowForge's theme approach and its omarchy
corpus — **35 files verified at source**, including `delorean` and `retro-82`, under MIT with
attribution.

- **The dark default is the brand palette** from the site mocks; **light is first-class**, based on
  the brand manual's §4.5 paper expression rather than invented.
- **Every theme must satisfy the four-state contract or FAIL VALIDATION** — the four states
  distinguishable by colour with **measured CVD separation**, and every text pair clearing 4.5:1 on
  that theme's own ground. **A theme that cannot express four distinguishable states does not ship
  broken; it does not ship.**
- Glyph and label are theme-invariant, so co-render holds even where colour degrades — which is what
  makes a marginal theme *degraded* rather than *dangerous*.

The token contract and the validation are specified in
[`../design/05-theme-system.md`](../design/05-theme-system.md).

**This closes the light-mode CVD gap properly.** It stops being *"someone should check the light
palette"* and becomes *"no theme ships without passing this"* — a standing gate rather than a
one-time task.

### Alternatives

| Considered | Rejected because |
|---|---|
| **One theme, the brand's** | The founder asked for the FlowForge set. And a terminal product whose users have strong theme opinions is a product that will be re-themed with or without permission |
| **Themes as Go code** | A theme becomes a release. It must be a file a player can drop in |
| **Accept any palette** | The four states are the product's most-used visual system. A theme that collapses two of them silently breaks the one thing co-render exists to protect |

### Consequences

**Easier:** a theme is a contribution, not a change. **Harder:** validation must exist and run, and
some attractive palettes will fail it — which is the decision working, not misfiring.
**Cost to reverse: medium** — the token contract is what every component reads.

---

## D8 · Terminal images are foundational, with an honest degrade

### Decision

Cover art is **Phase 1**, reversing revision A which deferred it and declined to assume inline
images. Founder direction: *"without image is not an option."*

**Kitty graphics protocol** (Kitty, Ghostty) plus **iTerm2**; **Sixel deferred** with its reasons.
**Capability detection at startup, failing closed** — never a config flag the player must find, and
never guessing yes.

**A terminal without image support is a supported configuration**, not a warning state: it renders
the full text deck, plus a **once-only, dismissible** note that Ghostty or Kitty would show covers.
Never recurring, never blocking, never phrased as a fault in the player's setup.

Covers are **cache, never truth** — XDG cache, bounded, evicting; deleting it costs only bandwidth.
`Cover()` **never fetches and never blocks**: a missing cover is never worth a dropped frame.

### Alternatives

| Considered | Rejected because |
|---|---|
| **Defer to Phase 2** (revision A) | Optimised for the terminal instead of the player. A game library without covers is a spreadsheet |
| **Require an image-capable terminal** | Refuses a supported configuration to avoid designing a degrade |
| **Half-block / ASCII art covers** | A picture of a picture, and precisely the retro-kitsch the brand rules out |
| **Sixel in Phase 1** | Poor colour fidelity on photographic content and slow at deck sizes. Revisit on evidence |

### Consequences

**Easier:** the product looks like what it is. **Harder:** a whole capability axis to detect, cache
and bound, and a degrade path that must stay genuinely first-class rather than becoming a
second-class citizen once covers look good. **Cost to reverse: medium-high** — it is a seam, but
Phase 1 screens are designed around covers existing.

---

## D9 · Internationalisation from the first line

### Decision

**No user-facing string literal in code.** Every string comes from a catalogue, by key, enforced by
a **lint that fails the build** — not a convention. English is the only language in v1; adding
`pt-BR` must be **adding one file**, with no code change.

`golang.org/x/text`, already in the dependency set, for plural rules, locale-aware number, currency
and date formatting, collation, and locale negotiation.

**Four things naive i18n misses, all specified**: terminal **cell width** (the TUI-specific one, and
already biting *today* — `Pokémon` and `Ōkami` are not ASCII); **currency**, because the budget
feature is money in the player's own currency; **plurals** via CLDR, exercised by the English
catalogue so the mechanism is live from day one; and **collation**, because byte order already
mis-sorts accented titles in English.

**`zerado` and `sinopse` stay Portuguese** — ratified brand vocabulary carried with a translator
note, so no future translator "fixes" them.

Founder direction, naming it as a scar: *"we will not make the same mistake we did on FlowForge."*

### Alternatives

| Considered | Rejected because |
|---|---|
| **Retrofit when a second language is wanted** | The strings a retrofit misses are the ones nobody looks at — errors, empty states, the fatal screen. The worst screens to meet an unreadable language on |
| **`map[string]string`** | No plurals, no CLDR. It becomes `x/text` eventually, having lost the intervening strings |
| **`go-i18n`** | Capable, but a second dependency for what `x/text` already does |
| **A convention instead of a lint** | A rule nobody can check is a convention, and conventions lose to deadlines |

### Consequences

**Easier:** the second language is a file. Width-correct rendering is a *side effect* of doing this
properly, and it fixes a defect that exists today. **Harder:** every string costs a key, and the lint
will occasionally be wrong and need an allow-list entry with a reason. **Cost to reverse: highest of
the three added here** — retrofitting i18n into a finished product is the canonical expensive
migration, which is exactly why it starts now.

---

## What this ADR does not decide

Named so they are not assumed to have been settled here:

- **Which metadata provider Phase 2 ships with.** The seam is decided; the provider is not, and
  cannot be until IGDB answers. [`../blueprint/06-data-seams.md`](../blueprint/06-data-seams.md) §3.2
  records what Phase 2 looks like if the answer is no.
- **Terminal inline-image support** for cover art. Not assumed anywhere; the text deck is the
  default, and the cover deck is a **Phase 1** mode (D8) that degrades honestly where the terminal
  cannot draw.
- **The community source.** It is not named anywhere in the product while the ratified decision
  stands, and the recommendation engine must not hard-code a source into its shape.
- **Whether Phase 4 accounts are email-and-password, OAuth, or something else.** A Phase 4
  question; nothing in Phase 1 depends on it.

---

## Status

**PROPOSED.** **All nine** — D1 through D9 — are the third of the founder's three ratification
questions. Ratifying the bundle ratifies them; that ratification is itself the authorization to
emit the Phase 1 implementation tickets, with no second approval step.

The clause originally read *"these four"*, written when there were four. **D5** (the
media-polymorphic core) and **D6** (audio) arrived by founder direction on 2026-08-25 and are the
two [](../blueprint/00-index.md) marks **Highest** and
**Low-but-licensing** reversal cost — so on the document's own words, ratifying "them" would not
have reached the two most expensive decisions in it. Corrected.
