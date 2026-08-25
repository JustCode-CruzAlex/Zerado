---
archetype: adr-detail
adr: ADR-0001
title: "Zerado foundational architecture — provider seam, persistence, navigation model, Phase 4 sync boundary"
status: PROPOSED
date: 2026-08-25
ticket: "#2"
chart_ref: "docs/adr/charts/"
---

# ADR-0001 · Zerado foundational architecture

Four decisions that are expensive to reverse once Phase 1 code exists. They are recorded together
because they constrain each other: the sync boundary decides what persistence must carry, and
persistence decides what the provider seam can promise.

**Drawings for this ADR** — eight sheets, each rendered in both themes, under
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
