---
archetype: adr-detail
adr: ADR-0001
title: "Zerado foundational architecture — provider seam, persistence, navigation model, Phase 4 sync boundary"
status: PROPOSED
date: 2026-08-25
ticket: "#2"
---

# ADR-0001 · Zerado foundational architecture

Six decisions that are expensive to reverse once Phase 1 code exists. They are recorded together
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

---

## Context

Zerado is a terminal-first game library, greenfield as of 2026-08-25 — no Go code exists. The
stack is settled: Go, Bubble Tea, and the Charm ecosystem. Four phases are published and binding
(`content/landing-copy.md` §12), as are six ratified promises (`ratification/decisions.md`):
local-first, one SQLite file, no Zerado-run server before Phase 4, the player's own API keys, an
affiliate disclosure, and no named community source.

Those promises are not context. They are **architectural constraints**, and each decision below
names the ones it is paying for.

Two facts sharpen the problem:

1. **Physical copies must be first-class from day one.** The published copy says *"A physical copy
   isn't a second-class row in the list."* The obvious implementation — a Steam-shaped row with an
   `is_physical` flag — makes that sentence false in the second week.
2. **IGDB is free for non-commercial use only**, and Zerado's funding model is affiliate
   commission, which is commercial. Cover art and *sinopse* are the visual backbone of Phase 2.
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
  swaps the credit. `Quote` carries `AffiliateURL` and the disclosure obligation **in the same
  struct**, so a refactor cannot separate them.

Pays for: *"A physical copy isn't a second-class row"*, the affiliate disclosure, and the IGDB
risk.

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
| **A modal command palette as the primary navigation, from Phase 1** | A palette earns its place when the surface is bigger than the key map. With eleven screens and one home, `?` plus a footer hint is better — and reserving the keys costs nothing |
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
| **Sync the whole library** | Contradicts local-first in spirit, makes the server expensive enough to break the "premium account or a donation" statement, and stores data the player already has two other copies of |
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

## D5 · The core entity is a **media item**, not a game

### Decision

`media_item` is the core entity, carrying a `media_type`. `game` is the first type; `book` is the
second; `film` and `series` are a plausible third and fourth. **Phase 1 ships games only** — the
column is constrained to `'game'`, there is no `--type` flag, and a Phase 1 player cannot tell this
decision exists.

- **Shared facts live on the core:** ownership, acquisition (digital or physical), the four states,
  progress, mood tags, rating, notes, the source provider, price history.
- **Type-specific facts live in a typed extension:** `media_game` (Steam appid, achievements);
  later `media_book`, `media_film`, `media_series`.
- **Progress is generic** — `value` + `unit` + `source` — rather than a typed `playtime`. This
  diverges from the letter of the founder direction and the reason is in
  [`../blueprint/11-media-model.md`](../blueprint/11-media-model.md) §1: it is what keeps the
  four-state derivation **one function** instead of one per type.
- **Derived versus manual is a capability of the `(provider, media type)` pair, not of the type.**
  A Steam game derives; a cartridge does not; a paper book does not; a Kindle book would. Phase 1
  already needs this for physical copies, so the mechanism is not new work.
- **The provider seam gains `MediaTypes()`** and per-type `Capabilities`. Nothing else reshapes.
- **Mood tags carry a type-neutral `key` and a per-type `label`** — one engine, per-type vocabulary.

**The four states were verified type by type, not assumed** — see
[`../blueprint/11-media-model.md`](../blueprint/11-media-model.md) §2. Two findings, in §4:

- **F-1 · a film has no meaningful `IN PROGRESS`**, so **`film` and `series` should be two types**,
  not one "video" type. Costs nothing now; costs a migration later.
- **F-2 · an ongoing series needs *caught up*, and it must not become a fifth state.** The four are
  ratified and CVD-verified at a measured ΔE floor of 11.9 under deuteranopia; a fifth colour would
  have to be re-verified against that, for one type, in a later phase. Model *caught up* as a derived
  presentation of `IN PROGRESS` from typed facts.

### Alternatives

| Considered | Rejected because |
|---|---|
| **Keep it games-only; add types when books arrive** | The retrofit rewrites every table, every seam and the whole state machine. This is the single most expensive thing in the bundle to defer, which is why it is the one thing designed ahead of need |
| **One table with nullable type-specific columns** | A `game` row carrying null `isbn` and `page_count` — the table grows a column per type forever and every query learns which nulls are meaningful |
| **JSON blob for type-specific facts** | Unqueryable and unconstrained. The typed facts are exactly what a type-specific screen filters on |
| **A plugin/registry for media types** | An abstraction with one implementation for the whole life of Phase 1. Three types are foreseen and they ship in the binary |
| **Playtime as a typed game fact (the direction's letter)** | Forks the four-state derivation per type. Recorded as an open founder question in [`../blueprint/13-handoffs.md`](../blueprint/13-handoffs.md) §5 |

### Consequences

**Easier:** adding books should be one migration, one provider, one typed extension and one label
set — with **zero** changes to the state machine, the recommender, the offline contract or any
existing screen. That is a measurable test, and it is the one to hold this decision to.

**Harder:** every Phase 1 query carries a `media_type` predicate it does not yet need, and one join
to `media_game` for facts that could have been columns. Both are cheap at this size, and both are
the price of not doing the retrofit.

**Cost to reverse:** **highest, jointly with D4.** It is the shape of every table and every seam.

---

## D6 · Audio ships in Phase 1 — **bundled, off by default, fully removable**

> **This reverses a verdict recorded in this same bundle** — see
> [`../blueprint/08-prior-draft-analysis.md`](../blueprint/08-prior-draft-analysis.md) §3, kept
> visible as superseded rather than deleted. The direction arrived as an **agent relay carrying no
> ratification authority**; it was acted on because this is reversible document work on a draft PR,
> and it is flagged so the founder confirms it rather than inherits it.

### Decision

**The reversal is not a change of mind about the same object — the object changed.** What was
rejected was a *network streamer, always on*. What ships is a *local, bundled, opt-in subsystem, off
by default, that makes no network requests at all*.

- **Bundled, never streamed.** No CDN, no fetch, no cache warm. This is what keeps three ratified
  promises intact — *no background telemetry*, *works with the network off*, *the only traffic is
  services you connected* — and it makes audio `WORKS` in the offline contract.
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

### The open founder decision — music licensing

**Bundled tracks must be DRM-free and licensed for commercial redistribution, or they do not
ship.** Both halves bite: the repository is **public**, so tracks are redistributed by every clone
and every release artifact; and the funding model is **affiliate commission**, so a
"non-commercial use" licence does not cover it — the same trap already identified for IGDB.

> **Recommendation, for the founder to accept or refuse:** ship Phase 1 with **interface FX only**,
> and make the **music bed user-supplied** (point Zerado at a local directory). That delivers the
> feeling, removes the licensing blocker from the critical path entirely, and leaves a bundled
> soundtrack as a later addition once rights are actually cleared. Audio ships either way.

### Alternatives

| Considered | Rejected because |
|---|---|
| **Streamed music** (the prior draft) | Contradicts three ratified public promises, and would need its own ratification |
| **On by default** | A terminal program that makes noise unprompted is uninstalled before it is understood |
| **One channel with one mute** | Someone may want keyclicks without the soundtrack, or the reverse. Neither is the odd request |
| **Audio compiled into the default build** | Costs D2's pure-Go cross-compile matrix — the property that makes releases one file and no toolchain |
| **`Cue()` returning an error** | There is no audio failure a *screen* could usefully handle. An error return invites a caller to block on it |
| **A second key to pause music, distinct from mute** | A redundant control on a scarce surface, for no accessibility gain — `m` already satisfies 1.4.2 |

### Consequences

**Easier:** the retro-future feeling becomes something the player turns on. `NullAudio` as the
default build means the silent path is the well-tested one.

**Harder:** a build-tag split means two binaries to produce and a Settings line that must explain
which one is running. And the licensing question is genuinely unresolved — it is on the founder's
list, not a specialist's.

**Cost to reverse:** **low for the design** — the seam is removable at compile time. The
**licensing** commitment is the part that is expensive to undo once tracks ship in a public repo.

---

## What this ADR does not decide

Named so they are not assumed to have been settled here:

- **Which metadata provider Phase 2 ships with.** The seam is decided; the provider is not, and
  cannot be until IGDB answers. [`../blueprint/06-data-seams.md`](../blueprint/06-data-seams.md) §3.2
  records what Phase 2 looks like if the answer is no.
- **Terminal inline-image support** for cover art. Not assumed anywhere; the text deck is the
  default and the cover deck is a Phase 2 mode that may never be universally available.
- **The community source.** It is not named anywhere in the product while the ratified decision
  stands, and the recommendation engine must not hard-code a source into its shape.
- **Whether Phase 4 accounts are email-and-password, OAuth, or something else.** A Phase 4
  question; nothing in Phase 1 depends on it.

---

## Status

**PROPOSED.** These four are the third of the founder's three ratification questions. Ratifying
the bundle ratifies them; that ratification is itself the authorization to emit the Phase 1
implementation tickets, with no second approval step.
