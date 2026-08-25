---
title: Zerado — data seams
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-06
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: implementation-plan
ticket: "#2"
---

# Data seams — interfaces, not implementations

Five seams. Each is named, each has the shape of its contract, and each is swappable without
touching a screen.

The Go below is **shape, not code**. It is here because a seam described in prose is a seam
nobody can hold you to, and because the exact signature is the decision — not the package it
eventually lives in. **No implementation is written by this ticket.**

> **Scope of this document — read before the Go.** The spine decides **that** a seam exists, what it
> is responsible for and where its boundaries are. It does **not** decide the final Go signatures —
> that is `fft-api-designer`'s deliverable, and this document is its brief. Every Go block below is
> **shape, not signature**. [`13-handoffs.md`](./13-handoffs.md) states exactly which decisions are
> load-bearing and which spellings are free to change.

| § | Seam | Interface | Why it must be a seam |
|---|---|---|---|
| 2 | **Store provider** | `Provider` / `Syncer` | Steam first, then PlayStation, GOG, EA — and physical, which is not a store at all |
| 3 | **Metadata** | `MetadataProvider` | IGDB's commercial terms are unresolved. The provider must be replaceable, and *having none* must be a designed state |
| 4 | **Price history** | `PriceProvider` | Every quote carries its age, and a price without one is the product giving financial advice from memory |
| 5 | **Persistence** | `Store` | One SQLite file is a public promise. The interface is what keeps it one file |
| 5.4 | **Credentials** | `Vault` | The player's own keys. They must never be inside the file the player is told to back up and share |
| — | **Audio** | `Audio` | Ships in Phase 1, bundled and off by default. Fully removable at runtime *and* at compile time. Specified in [`12-audio.md`](./12-audio.md) |

---

## 1 · The rule all five obey

> **A screen never talks to a provider. A provider never talks to a screen.**

Between them sits the persistence seam. A sync writes to the store; the screens read from the
store. Nothing renders from a network response, ever.

Three things fall out of that, all of them promises the page already made:

- **the offline contract is structural** — a screen that only reads local state works offline
  because it cannot do anything else;
- **a failed sync cannot break a screen**, only leave it stale — and stale is a state the design
  system has a banner for;
- **"no telemetry running in the background"** is provable by inspection: there is exactly one
  place network I/O can originate.

---

## 2 · The store-provider seam

### 2.1 · Interface segregation is the whole design

```go
type ProviderID string   // "steam" · "physical" · "playstation" · "gog" · "ea"

// Capabilities is what the provider can actually do. Every screen and every
// derivation reads this instead of switching on ProviderID.
type Capabilities struct {
    Sync        bool   // can fetch a library over the network
    Playtime    bool   // reports minutes played
    LastPlayed  bool   // reports a last-played timestamp
    OwnedSince  bool   // reports an acquisition date
    Credentials []CredentialField // what Z-02 renders; empty means none needed
}

// Provider is what EVERY source of games implements — including the ones
// that are a human with a keyboard.
type Provider interface {
    ID() ProviderID
    Display() string          // "Steam", "Physical shelf"
    Capabilities() Capabilities
}

// Syncer is implemented only by providers that can fetch. `physical` does
// not implement it, and that is the point: it is not a Steam-shaped provider
// with a hole where Sync should be.
type Syncer interface {
    Provider
    // Sync streams the provider's current view of the library.
    // ctx cancellation MUST abort in-flight I/O — Z-03 lets the player quit
    // mid-sync and the goroutine budget is zero leaks.
    Sync(ctx context.Context, c Credentials) (<-chan Item, error)
}

// Item is the provider's view of one owned title. Every optional field is a
// pointer, so "not reported" and "reported as zero" are different facts.
type Item struct {
    ProviderRef string      // the provider's own id — Steam appid, or a UUID for physical
    Title       string
    Platform    string
    Playtime    *int        // minutes; nil = not reported
    LastPlayed  *time.Time  // nil = not reported
    OwnedSince  *time.Time
}
```

### 2.2 · Physical copies are first-class by construction

`physical` is a `Provider` and **not** a `Syncer`:

```go
Capabilities{Sync: false, Playtime: false, LastPlayed: false,
             OwnedSince: true, Credentials: nil}
```

Its "sync" is `Z-08 Add a game by hand` writing an `Item` directly. Everything downstream — the
row renderer, the state derivation, the filter, the summary counts, the detail view — reads
`Capabilities`, not `ProviderID`, so a physical copy is structurally the same kind of thing as a
Steam game.

The published promise is *"A physical copy isn't a second-class row in the list."* This is what
makes it true rather than aspirational. The failure mode it prevents is concrete: an
`is_physical bool` on a Steam-shaped row means every new feature has to remember to special-case
it, and eventually one will not.

### 2.3 · A new store adds zero screens

`Z-02 Connect a store` renders from `Capabilities().Credentials`:

```go
type CredentialField struct {
    Key         string  // "api_key"
    Label       string  // "Steam API key"
    Help        string  // shown under the field
    Secret      bool    // masked, and stored in the Vault
    Validate    func(string) error
    HelpURL     string  // "where do I get this?"
}
```

Adding GOG is: implement `Provider` + `Syncer`, declare the credential fields, register it.
No screen changes, no route changes, no schema changes. **That is the test of whether this seam
is right**, and it is the reason `Z-02` is named "connect a store" and not "connect Steam".

### 2.4 · A game a sync stops returning is **tombstoned, never deleted**

`fft-tui-designer` refused to write `Z-05` copy about this case because the seam never decided it,
which was the right call — it is a decision about a player's own data. Decided here.

**A sync that no longer returns a game does not delete it.** It sets `absent_since`, and the row
stays.

The reason is that *"the provider stopped returning it"* and *"the player no longer owns it"* are
not the same fact, and only the first one is observable. A game can vanish from a Steam response
because of a delisting, a region change, a family-share expiry, an API hiccup, a truncated page or a
rate limit. Meanwhile the row carries **the player's own work** — the status, the notes, the mood
tags, the fact that it is *zerado*. That is the one thing the product promises is theirs. Destroying
it because a third-party API omitted a line is the product acting on someone else's say-so against
its owner.

Deletion is irreversible. Tombstoning is not. When the difference is that stark and the evidence is
that weak, the reversible option wins.

| Rule | |
|---|---|
| `absent_since` is set | on the first **complete, successful** sync that does not return the game |
| It is cleared | the moment the game comes back — silently, with no notice |
| Absent rows | are excluded from the **default** library view, remain findable by filter, and `Z-05` says plainly why |
| Deletion | happens only when the **player** asks. Never as a consequence of a sync |

**The guard that matters most: only a sync whose `status` is `ok` may tombstone anything.** A
`partial`, `failed` or `cancelled` sync must not — in a truncated stream, *not returned* and *not
reached* are indistinguishable. A sync returning **zero** items is the private-profile case, already
a `REFUSES` ([`07-offline-contract.md`](./07-offline-contract.md) §3.1); it must never be read as
"the player's whole library was removed."

**`absent` is not a fifth state.** It is an orthogonal presence flag, for the same reason *caught
up* is not a fifth state ([`11-media-model.md`](./11-media-model.md), Appendix): the four states are
ratified and CVD-verified, and a game that is absent still has a status — usually the most valuable
one, because a game you finished and no longer own is exactly the row you would be angriest to lose.

A row that is absent **and** carries no player data at all — no manual status, no notes, no mood
tags — may be **offered** for bulk removal. Offered, not done.

### 2.5 · Why the sync streams

`Sync` returns a channel rather than a slice for three reasons:

1. `Z-03` can show honest running counts on a 1,000-title library instead of an indeterminate
   wait followed by a number;
2. a cancel mid-sync leaves a **valid partial** library rather than nothing;
3. memory stays bounded on a library of any size.

The cost is that the store's upsert has to be transactional per batch rather than per sync —
recorded here so it is a known trade, not a surprise.

---

## 3 · The metadata seam — designed around a risk that is not resolved

IGDB is free for **non-commercial** use. Zerado's funding model *was* affiliate commission, which
is commercial. Cover art and *sinopse* are the visual backbone of Phase 2. **This is unresolved
and the architecture must not assume it resolves favourably.**

```go
type MetadataProvider interface {
    ID() ProviderID
    Lookup(ctx context.Context, r Ref) (Metadata, error)
    // Attribution is REQUIRED and is not optional to render. A provider
    // declares the credit line and licence terms its data ships under; the
    // detail view renders whatever it returns.
    Attribution() Attribution
}

type Ref struct {  // enough to identify a game without assuming a store
    Title    string
    Platform string
    Provider ProviderID
    ProviderRef string
}

type Metadata struct {
    Sinopse    string
    CoverRef   string     // a local cache path, never a remote URL at render time
    ReleasedAt *time.Time
    Genres     []string
    FetchedAt  time.Time  // every network-derived value carries its age
}
```

### 3.1 · The two decisions that de-risk it

**`Attribution()` is a method, not a config string.** The credit a source requires is a property
*of the source*, so swapping the source swaps the credit automatically. A licence change becomes
a data change instead of a redesign.

**`NullMetadataProvider` is a first-class, designed state — not an error.** When there is no
metadata provider, or when it returns nothing, the detail view shows a designed *no-metadata*
composition, not an error banner. This is the difference between a product that works without
IGDB and a product that is broken without IGDB.

### 3.2 · What Phase 2 looks like if IGDB says no

Written down now, because deciding it under pressure later is how a product ends up shipping
whatever was easiest that week.

| Need | If IGDB is unavailable |
|---|---|
| **Cover art** | SteamGridDB (community art, its own terms) for synced titles; store-supplied header images where the store's terms allow; **nothing** for physical copies, which is honest and already the majority case for a shelf |
| ***Sinopse*** | The store's own product description where terms allow; Wikidata (CC0) for release data and identity; a hand-entered description field, which `Z-08` already needs for physical copies anyway |
| **Identity / matching** | Wikidata + the store's own ids. This is the part that must not depend on IGDB, because it is what the Phase 4 sync boundary needs (§5) |

**The floor:** Phase 2 ships with **text-only enrichment** if art cannot be licensed. The library
is still sorted by mood, still shows state, still answers what to play tonight. Cover art is an
enhancement to a working product, and the architecture must keep it that way.

---

## 4 · The price seam — the disclosure is structural

```go
type PriceProvider interface {
    ID() ProviderID
    Quote(ctx context.Context, r Ref, cur Currency) (Quote, error)
    Attribution() Attribution
}

type Quote struct {
    Current    Money
    Low        Money      // the all-time low
    LowAt      time.Time  // when. A low with no date is not information
    Shop       string
    URL        string      // a plain shop link. No affiliate tag, ever.
    ObservedAt time.Time   // mandatory, and rendered — see 07 §4
}
```

**There is no affiliate URL, and that is a decision, not an omission.** Founder direction,
2026-08-25: affiliate links are dropped so that Zerado is cleanly **non-commercial** — free
software, donation-supported, zero revenue.

**The price feature survives intact.** Current price, the all-time low, when the low was, and the
*"wait or buy"* verdict all remain exactly as designed. What goes is the commission tag on the
outbound link, and with it the disclosure obligation that used to travel in this struct.

Two consequences worth naming:

- **It answers the IGDB question this bundle has carried since revision A.** IGDB's published test
  is whether the *project generates revenue* — not whether it charges users. With no commission,
  Zerado generates none. **Stated as a reading of IGDB's published rationale, not a legal opinion:**
  the founder should confirm directly with IGDB that a donation-funded open-source project qualifies
  for the free tier. That is a founder action, not a resolved fact, and it is on the gate list.
- **The metadata seam stays provider-agnostic anyway.** The hedge was right for reasons that do not
  depend on IGDB's answer, and removing it because one risk receded would be the wrong lesson.

`ObservedAt` is mandatory and is rendered. **Zerado never shows a price without its age** — that
rule is repeated in [`07-offline-contract.md`](./07-offline-contract.md) §4 because it is the
single easiest way for this product to become dishonest by accident.

IsThereAnyDeal is named in the published copy as the source. It is nevertheless behind this
interface, for the same reason IGDB is: a source that is named today can change its terms
tomorrow, and the page's claim is about *what* the data is, not about who is guaranteed to
supply it forever.

---

## 5 · The persistence seam

### 5.1 · One file, and the interface is what keeps it one

```go
type Store interface {
    // Reads — everything a screen ever calls.
    Games(ctx context.Context, q Query) ([]Game, error)
    Game(ctx context.Context, id GameID) (Game, error)
    Counts(ctx context.Context, q Query) (StatusCounts, error)
    Connections(ctx context.Context) ([]Connection, error)
    LastSync(ctx context.Context, p ProviderID) (*SyncRun, error)
    Setting(ctx context.Context, key string) (string, bool, error)

    // Writes — everything a command ever calls.
    // SetStatus(nil) clears the manual override (05-state-machine §5).
    SetStatus(ctx context.Context, id GameID, s *Status) error
    UpsertBatch(ctx context.Context, p ProviderID, items []Item) (BatchResult, error)
    AddManual(ctx context.Context, g ManualGame) (GameID, error)
    RecordSyncRun(ctx context.Context, r SyncRun) error
    SetSetting(ctx context.Context, key, value string) error
}
```

Every data access is behind this interface — the house rule, and here it also does two specific
jobs: it is the only thing that knows there is a database at all, and it is the boundary the
Phase 4 sync engine will attach to without touching a screen.

### 5.2 · The file, and where it lives

**One SQLite file**, at the XDG data directory, with an env override:

```
$ZERADO_DB  ??  $XDG_DATA_HOME/zerado/library.db  ??  ~/.local/share/zerado/library.db
macOS:      ~/Library/Application Support/zerado/library.db
```

The published promise is: *"Your whole library lives in a single SQLite file you can back up,
move, or delete — nobody else holds a copy."* Two consequences the architecture must honour:

- **it stays one file.** No sidecar cache database, no separate index. SQLite's WAL journal
  produces `-wal` and `-shm` companions while the process runs; these are checkpointed and
  removed on clean shutdown, and the promise is about what the player copies, so a clean
  shutdown must leave exactly one file. Recorded because "one file" and "WAL mode" have to be
  reconciled deliberately rather than discovered in a bug report;
- **cover-art blobs do not go in it.** Phase 2's image cache is disposable, regenerable, and
  potentially large. It lives in the XDG **cache** directory, where the OS is allowed to delete
  it — which is the correct semantics for something re-fetchable, and it keeps the file the
  player backs up small enough to actually back up.

### 5.3 · Pure-Go SQLite, and why the driver choice is an ADR

**`modernc.org/sqlite`** (pure Go), not `mattn/go-sqlite3` (cgo).

| | Pure Go | cgo |
|---|---|---|
| Cross-compile a release for six platforms | one `GOOS`/`GOARCH` matrix, no toolchains | a C toolchain per target |
| Static single binary | yes | fights you |
| Start-up | no dynamic link | fine, but not better |
| Raw write throughput | slower | faster |
| Binary size | larger | smaller |

The workload is a few hundred rows and a few thousand writes per sync. The throughput advantage
is irrelevant at that size; the distribution advantage is the whole difference between "download
one file and run it" and "install a build environment". The page's claim is *"It's a text
program. It starts instantly"* — that claim is about the player's experience of a **single
downloadable binary**, and cgo is the thing that most often costs one.

**This is the cheapest of the four ADR decisions to reverse**: the driver sits behind
`database/sql` and the `Store` interface. Recorded as a decision anyway, because a driver chosen
by accident in week one is a driver nobody revisits in year two.

### 5.4 · Credentials live outside the library file

The Steam API key is the player's own key, and the library file is a thing the player is
explicitly invited to **back up, move, and delete**. A credential inside it would be a
credential in every backup, every copy, and every support-ticket attachment.

```go
type Vault interface {
    Get(ctx context.Context, p ProviderID, key string) (string, bool, error)
    Set(ctx context.Context, p ProviderID, key, value string) error
    Delete(ctx context.Context, p ProviderID, key string) error
    Backing() string  // "keychain" | "file" — Z-09 shows this, honestly
}
```

- **Preferred:** the OS keychain (macOS Keychain, Secret Service, Windows Credential Manager).
- **Fallback:** `credentials.json` beside the library file, mode `0600`, with `Z-09 Settings`
  stating plainly which backing is in use.

This *strengthens* the one-file promise rather than violating it: the library file stays purely
a library, so sharing it is safe. Settings shows the backing because a security property the
player cannot see is a security property they cannot rely on.

---

## 6 · The Phase 4 sync boundary — decided now because it constrains Phase 1

The ratified promise is: **no Zerado-run server before Phase 4.** Phase 4 introduces one. The
expensive question is *what crosses it*, and the answer has to be known in Phase 1 because it
decides what the schema must carry.

### 6.1 · What crosses, and what never does

| Crosses | Never crosses |
|---|---|
| `status_manual` + `status_changed_at` | The library itself — it is re-derivable from the player's own sources |
| Mood tags the player assigned by hand (Phase 2) | `playtime_minutes`, `last_played_at` — provider facts, re-fetchable |
| Manually-entered games (`physical`) — **the only rows not re-derivable from a store** | **Credentials. Ever.** Each device connects with its own keys |
| A stable game identity (§6.2) | Cover art and *sinopse* — cache, refetchable |
| | Prices — always fetched fresh with their age |

The rule underneath: **only what the player typed crosses.** Everything a machine can recompute,
each device recomputes.

That is not only a privacy posture — it is what keeps the Phase 4 server small enough to be
supportable by **donations alone**, which is itself a public statement. A
server that stored everyone's whole library would cost what the page says it will not.

### 6.2 · The one thing Phase 1 must get right: game identity

Merging two devices needs a key that is the same on both. `rowid` is not; `(provider, provider_ref)`
is not, because the same game can arrive from Steam on one machine and as a cartridge on another.

**Decision:** every game carries a **`game_uid`** — a stable, content-derived identity, assigned
at insert and never changed:

```
game_uid = uuidv5(namespace_zerado, normalise(title) + "|" + normalise(platform))
```

`normalise` lowercases, strips punctuation and articles, and folds diacritics. It is imperfect —
two editions of the same game may collide, and the same game may fail to match across platforms —
so:

- `game_uid` is **stable, not authoritative**: it is a merge *hint*, and Phase 4's merge presents
  ambiguous matches to the player rather than guessing;
- a `merged_into` column exists from Phase 1 so two rows can be joined later without a migration
  that rewrites keys;
- `(provider, provider_ref)` remains the **unique** constraint. `game_uid` is indexed, not unique.

Adding this in Phase 1 costs one column and one index. Adding it in Phase 4 costs a migration
that has to invent stable identities for rows whose titles the player has since edited — which
is the migration nobody wants to write, and the reason this decision is in ADR-0001.

### 6.3 · Conflict resolution

**Last-write-wins on `status_changed_at`**, per game.

This is a deliberate simplicity choice and its limits are stated rather than hidden: it is
adequate because the conflicting parties are *one person on two devices*, both of whom agree
about what they did; it would be inadequate for multiple users, which Phase 4's community layer
does not require (comments and reviews are append-only and are a different table entirely).

If two devices set different statuses while both offline, the later timestamp wins and the
earlier change is lost silently. That is acceptable for this shape of conflict and unacceptable
to leave undocumented — so `Z-22 Devices and sync` shows the last merge and what it resolved.

---

## 7 · What is deliberately NOT a seam

Naming these prevents the most expensive kind of over-engineering: an abstraction with one
implementation, forever.

| Not a seam | Why |
|---|---|
| **The renderer** | Bubble Tea is the decision, not an option. An abstraction over the TUI framework would buy nothing and cost every screen |
| **The theme** | The brand manual is the theme. There is no second brand |
| **The clock** | A `func() time.Time` field on the structs that need it is enough for tests. An interface for this is ceremony |
| **The filesystem** | `os` plus a configurable root. `io/fs` where reading is all that is needed |
| **HTTP** | A shared `*http.Client` with a timeout, injected. Providers do not construct their own — that is how a "works offline" claim quietly stops being true |
| **A media-type abstraction** | Zerado is a games product. `item_type` exists and is CHECKed to `'game'`; that is the whole affordance ([`11-media-model.md`](./11-media-model.md)). An interface parameterised on a type that has one value is machinery without a purpose |
