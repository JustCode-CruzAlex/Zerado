package store

import (
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
)

// RunID identifies one sync run.
type RunID int64

// RunStatus is how a sync ended.
//
// Four values, and they are the four Z-03 renders as terminal states. The set
// is closed because each one has different consequences for what may happen
// next: only ok may tombstone, ok and partial both leave a valid library, and
// failed leaves it untouched.
type RunStatus uint8

const (
	// RunRunning is a run that has been opened and not closed.
	//
	// It is not one of Z-03's terminal states; it is what a row looks like
	// between StartRun and FinishRun, and what it looks like forever if the
	// process was killed. A run found in this state at start-up was killed,
	// which the ERD distinguishes from cancelled by finished_at being null.
	RunRunning RunStatus = iota

	// RunOK is a complete, successful sync. It is the only status that
	// permits absence reconciliation.
	RunOK

	// RunPartial is items arrived, then the stream errored. The library is
	// valid and incomplete.
	RunPartial

	// RunFailed is nothing arrived. The library is untouched, which is what
	// lets the copy say "Your library is unchanged — nothing was lost."
	RunFailed

	// RunCancelled is the player stopped it. What arrived is kept, and none
	// of it is an error.
	RunCancelled
)

// String returns the stored, machine-readable form. These values are written
// to the database and emitted by the CLI, so they are API.
func (s RunStatus) String() string {
	switch s {
	case RunOK:
		return "ok"
	case RunPartial:
		return "partial"
	case RunFailed:
		return "failed"
	case RunCancelled:
		return "cancelled"
	default:
		return "running"
	}
}

// MayTombstone reports whether a run in this status is allowed to mark rows
// absent.
//
// Exactly one status may. This is the executable form of 06-data-seams.md
// §2.4's guard, and it is a method here rather than a condition inside a store
// implementation so that a fake and a real store cannot disagree about it.
func (s RunStatus) MayTombstone() bool { return s == RunOK }

// StatusForFault maps a terminal fault onto the run status it produces.
//
// It is the join between the error taxonomy and the sync history, and it lives
// here so that the two cannot drift. sawItems distinguishes the two cases a
// fault alone cannot: the same transport failure is PARTIAL if items had
// already arrived and FAILED if none had, and Z-03's copy for those two is
// completely different — one says what was saved, the other says the library
// is unchanged.
func StatusForFault(err error, sawItems bool) RunStatus {
	switch {
	case err == nil:
		return RunOK
	case fault.Is(err, fault.KindCancelled):
		return RunCancelled
	case sawItems:
		return RunPartial
	default:
		return RunFailed
	}
}

// SyncRun is one recorded sync.
//
// partial exists as a status because a cancelled or truncated sync leaves a
// valid library, and recording it is what lets Z-04 say "synced 3 days ago,
// partially" instead of implying the library is complete.
type SyncRun struct {
	ID       RunID
	Provider provider.ID

	StartedAt time.Time

	// FinishedAt is nil for a run that was never closed — the process was
	// killed. That is a different fact from cancelled, and the ERD keeps them
	// apart on purpose.
	FinishedAt *time.Time

	Status RunStatus

	// Seen, New, Changed, Unchanged are what Z-03's DONE line reports:
	// "12 new. 4 changed. 231 unchanged."
	//
	// Unchanged is stored rather than derived because it is stated on screen
	// and because Seen minus New minus Changed is only equal to it when
	// nothing was skipped, which is not true of a partial run.
	Seen, New, Changed, Unchanged int

	// FaultKind is the classified failure, or KindUnknown for a run that did
	// not fail.
	//
	// The classification is stored, not a message and not a stack trace: the
	// message is a catalogue key rendered at read time, so a library synced in
	// English and read in Portuguese reports its history in Portuguese.
	FaultKind fault.Kind
}

// Age returns how long ago the run finished, for the degrade banner's
// "Last synced 3 days ago".
//
// An unfinished run has no age, which the caller must render as unknown rather
// than as zero — Z-03 §8.1's rule that a missing previous run reads "Nothing
// has ever been synced." rather than a fabricated duration.
func (r SyncRun) Age(now time.Time) (time.Duration, bool) {
	if r.FinishedAt == nil {
		return 0, false
	}
	if d := now.Sub(*r.FinishedAt); d > 0 {
		return d, true
	}
	return 0, true
}

// RunResult is what FinishRun records.
type RunResult struct {
	FinishedAt time.Time
	Status     RunStatus
	Seen       int
	New        int
	Changed    int
	Unchanged  int
	FaultKind  fault.Kind
}

// BatchResult is what one UpsertBatch changed.
//
// The three counts are summed across batches to produce Z-03's DONE line, and
// they are returned per batch rather than accumulated by the store because the
// store does not know when a sync has ended — only the caller that owns the
// stream does.
type BatchResult struct {
	New       int
	Changed   int
	Unchanged int
}

// Add accumulates another batch's result.
func (b *BatchResult) Add(o BatchResult) {
	b.New += o.New
	b.Changed += o.Changed
	b.Unchanged += o.Unchanged
}

// Total returns how many rows the batch touched.
func (b BatchResult) Total() int { return b.New + b.Changed + b.Unchanged }

// Absence is what one reconciliation changed.
//
// Returned rather than silent because the two numbers have opposite screen
// treatments: rows that came back are cleared silently with no notice (06 §2.4
// is explicit about that), while rows that went absent are the reason Z-05
// grows an "absent since" line. A caller that cannot tell them apart cannot
// honour either rule.
type Absence struct {
	// MarkedAbsent is how many rows the run stopped returning.
	MarkedAbsent int

	// Returned is how many previously-absent rows came back.
	Returned int
}

// Connection is a provider the player has connected.
type Connection struct {
	Provider provider.ID

	// AccountRef is an identifier, never a secret — the Steam ID that Z-02
	// shows back as "Connected as 76561198012345678".
	AccountRef string

	ConnectedAt time.Time

	// LastSyncAt and LastSyncStatus are what the degrade banner reads to say
	// "3 days ago", and nil when nothing has ever been synced.
	LastSyncAt     *time.Time
	LastSyncStatus RunStatus
}
