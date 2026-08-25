// Package store is the repository seam: the only thing in Zerado that knows a
// database exists.
//
// # The rule the whole architecture rests on
//
// A screen never talks to a provider. A provider never talks to a screen.
// Between them sits this seam: a sync writes here, screens read from here, and
// nothing renders from a network response, ever (06-data-seams.md §1).
//
// Three properties fall out of that, and each is a published promise:
//
//   - the offline contract is structural — a screen that only reads local
//     state works offline because it cannot do anything else;
//   - a failed sync cannot break a screen, only leave it stale, and stale is a
//     state the design system has a banner for;
//   - "no telemetry running in the background" is provable by inspection,
//     because there is exactly one place network I/O can originate and it is
//     not here.
//
// # There is no freshness parameter on any read, and that is deliberate
//
// The obvious way to express staleness in an API is a policy argument —
// CachedOnly, PreferCached, ForceFresh. This interface has none, because a
// read here cannot reach the network under any policy. What a caller gets is
// what is local, stamped with when it was observed ([aged.Value]), and the
// only way to make it fresher is an explicit sync, which is an action with a
// screen and a key behind it.
//
// That absence is the offline contract expressed in the signatures. A
// ForceFresh option would be a hole through which a render path could start a
// network request, and 07 §7.3's grep rule — no net/http outside the provider
// packages — would then be guarding a door with a window next to it.
//
// # SQL never leaves this package
//
// The interface below takes domain types and returns domain types. It exposes
// no query builder, no expression, no transaction handle and no rows cursor. A
// caller cannot construct SQL through it, which is the house rule and also
// what makes the Phase 4 sync engine attachable at this boundary without
// touching a screen.
package store

import (
	"context"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/aged"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/metadata"
	"github.com/JustCode-CruzAlex/Zerado/internal/pricing"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
	"github.com/JustCode-CruzAlex/Zerado/internal/status"
)

// Store is the data-access interface. Every read a screen performs and every
// write a command performs goes through it.
//
// # Context on everything
//
// Every method takes a context, including the ones that will always be
// microseconds against a local file. Two reasons, and neither is speculative:
// a filter that re-runs on every keystroke must be abandonable when the next
// keystroke arrives, and a Phase 4 implementation of this same interface will
// be doing I/O that genuinely blocks. An interface that has to grow contexts
// later grows them in every caller at once.
//
// # Reads never mutate, writes never render
//
// The split below is not cosmetic. A read may be called from a render path and
// must not write; a write is called from a command and returns what the screen
// needs to say what happened.
type Store interface {
	// ---- Reads: everything a screen ever calls. ----

	// Games returns the games matching q, ordered by sort title, ascending.
	//
	// The order is fixed rather than a parameter because Phase 1 has no sort
	// control (07 §2) and an order argument would be that control's back end.
	//
	// Rows with MergedInto set are never returned, in any query. A merged row
	// is not a row the player has; it is a tombstone that exists so Phase 4
	// never has to rewrite a key.
	Games(ctx context.Context, q library.Query) ([]library.Game, error)

	// Game returns one game by local id.
	//
	// A missing id is fault.KindNotFound, which on this seam is a real defect
	// — a screen asked for an id it was handed — rather than the routine
	// absence it is on the metadata seam.
	Game(ctx context.Context, id library.GameID) (library.Game, error)

	// Counts returns the pinned summary for the same set Games would return.
	//
	// It takes the same Query for a reason that is a rule rather than a
	// convenience: 05-state-machine.md §7 rule 2 forbids showing whole-library
	// counts above a filtered list, which is the most common way a list view
	// lies. Sharing the argument makes the two impossible to disagree.
	//
	// Limit and Offset are ignored here: the summary describes the filtered
	// set, not the visible page.
	//
	// The four state counts must sum to Total — status.Counts.Sums is the
	// assertion, and a store that cannot satisfy it has a bug in its grouping
	// rather than an unusual library.
	Counts(ctx context.Context, q library.Query) (status.Counts, error)

	// Connections returns every provider the player has connected.
	Connections(ctx context.Context) ([]Connection, error)

	// LastRun returns the most recent sync run for a provider, or nil when
	// there has never been one.
	//
	// The nil case is load-bearing rather than an edge: Z-03 §8.1 renders
	// "Nothing has ever been synced." instead of a fabricated age, and a
	// zero-valued SyncRun would make that fabrication the default.
	LastRun(ctx context.Context, p provider.ID) (*SyncRun, error)

	// Setting reads one setting. ok is false when the key has never been set,
	// which is distinct from a key set to the empty string — Z-09 has dials
	// whose "unset" and "off" are different.
	Setting(ctx context.Context, key string) (value string, ok bool, err error)

	// Metadata returns the cached enrichment for a game, stamped with when it
	// was observed at its source.
	//
	// Phase 2 fills this; Phase 1 implementations return ok=false. It is on
	// the Phase 1 interface because the alternative is a Phase 2 interface
	// change in a seam that eleven screens read, for a method whose Phase 1
	// implementation is one line.
	//
	// ok=false is the designed no-metadata state (06 §3.1), not an error. A
	// hand-added cartridge no source has ever heard of is the normal case for
	// a shelf, and it renders a composition rather than a banner.
	Metadata(ctx context.Context, id library.GameID) (aged.Value[metadata.Metadata], bool, error)

	// Quote returns the last known price for a game in a currency, stamped.
	//
	// Phase 3 fills this. The stamp is not optional and is not droppable: 07
	// §4 forbids rendering a price without its age, and returning an
	// aged.Value is how that becomes structurally true rather than a rule
	// somebody has to remember at 1am while tidying a layout.
	Quote(ctx context.Context, id library.GameID, cur pricing.Currency) (aged.Value[pricing.Quote], bool, error)

	// ---- Writes: everything a command ever calls. ----

	// SetStatus sets or clears the player's manual status.
	//
	// A nil s clears the override, which is a different action from setting
	// NotStarted and both are needed: clearing on a game with playtime makes
	// it IN PROGRESS immediately, while choosing NOT STARTED stores a manual
	// value that sticks (05 §5). A single non-nullable parameter cannot
	// express the difference, which is why this takes a pointer.
	//
	// The store stamps status_changed_at on every change, including a clear.
	// It is the timestamp that decides a Phase 4 conflict, and a clear is a
	// change the player made.
	SetStatus(ctx context.Context, id library.GameID, s *status.Status) error

	// UpsertBatch writes one batch of a provider's items and reports what
	// changed.
	//
	// Batched rather than per-sync because Sync streams: the trade is named in
	// 06 §2.5 and it is what makes a cancelled sync leave a valid partial
	// library rather than nothing.
	//
	// It never writes StatusManual or StatusChangedAt. That is the invariant
	// behind "a sync never changes a status the player set", and it belongs
	// here rather than in a provider because the store is the only writer and
	// therefore the only place it can be guaranteed.
	//
	// An empty batch is a no-op returning a zero BatchResult, not an error.
	// The caller that has just classified a zero-item sync as a refusal
	// (fault.KindEmpty) must not reach here at all, and this makes the
	// belt-and-braces case harmless.
	UpsertBatch(ctx context.Context, p provider.ID, items []provider.Item) (BatchResult, error)

	// AddManual writes one hand-entered item and returns its new local id.
	//
	// It takes a provider.Item rather than a bespoke struct because that is
	// the point of the Enterer seam: hand entry produces the same value a sync
	// does and travels the same path. The provider has already minted the ref.
	AddManual(ctx context.Context, p provider.ID, item provider.Item, s *status.Status) (library.GameID, error)

	// StartRun opens a sync run and returns its id.
	//
	// The run is opened before any item is written, so that a process killed
	// mid-sync leaves a run with a null finished_at — which the ERD names as
	// the difference between "killed" and "cancelled cleanly", and which the
	// next start-up can report honestly rather than inventing.
	StartRun(ctx context.Context, p provider.ID, startedAt time.Time) (RunID, error)

	// FinishRun closes a run with its terminal status and counts.
	FinishRun(ctx context.Context, id RunID, r RunResult) error

	// ReconcileAbsence tombstones the rows a completed sync did not return.
	//
	// seen is every ProviderRef the run delivered. Rows for this provider that
	// are not in seen get absent_since set; rows in seen that carried
	// absent_since have it cleared, silently and with no notice (06 §2.4).
	//
	// It takes a RunID rather than a ProviderID, and that is the guard rather
	// than a stylistic choice. Only a sync whose status is ok may tombstone
	// anything: in a truncated stream "not returned" and "not reached" are
	// indistinguishable, so a partial, failed or cancelled run must not be
	// allowed to conclude that a game is gone. An implementation MUST return
	// fault.KindPrecondition when the named run is unfinished or its status is
	// not ok.
	//
	// Passing a provider id would have made the illegal call spellable, and
	// the illegal call deletes the evidence that a player finished a game.
	// Nothing is ever deleted here: absence is reversible and deletion is not,
	// and 06 §2.4 chooses the reversible option precisely because the evidence
	// for the other one is so weak.
	ReconcileAbsence(ctx context.Context, run RunID, seen []string) (Absence, error)

	// Delete removes a game permanently. It happens only when the player asks,
	// never as a consequence of a sync.
	Delete(ctx context.Context, id library.GameID) error

	// SetSetting writes one setting.
	SetSetting(ctx context.Context, key, value string) error

	// SaveConnection records or replaces a provider connection.
	//
	// accountRef is an identifier and never a secret — the Steam ID, a GOG
	// username. Secrets go to the Vault, which is a different interface in a
	// different package for exactly this reason: the two destinations cannot
	// be confused if there is no method here that could take one.
	SaveConnection(ctx context.Context, p provider.ID, accountRef string, at time.Time) error

	// DeleteConnection forgets a provider connection. It does not touch the
	// games that provider contributed: disconnecting Steam is not a statement
	// about what the player owns, and the rows carry the player's own work.
	DeleteConnection(ctx context.Context, p provider.ID) error

	// Close releases the file.
	//
	// It must checkpoint WAL and leave exactly one file behind. "One SQLite
	// file you can back up, move, or delete" is a published promise, and -wal
	// and -shm companions surviving a clean exit would make it two.
	Close() error
}
