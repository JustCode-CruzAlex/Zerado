package provider

import "time"

// Stream is a sync in flight.
//
// # Why a channel and not a slice
//
// 06-data-seams.md §2.5 gives three reasons and they are all about honesty:
// Z-03 can show running counts on a 1,000-title library instead of an
// indeterminate wait followed by a number; a cancel mid-sync leaves a valid
// partial library rather than nothing; and memory stays bounded on a library
// of any size. The cost — the store's upsert becomes transactional per batch
// rather than per sync — is a known trade rather than a surprise.
//
// # Why a Stream and not a bare channel
//
// The spine's shape was Sync(ctx, c) (<-chan Item, error), which carries the
// failure that is known before anything arrives and has nowhere to put the
// failure that happens after. That second failure is the PARTIAL state — items
// arrived, then the connection broke — and it is one of Z-03's four terminal
// states, so it cannot be inexpressible.
//
// A Stream adds exactly two things to the channel: a terminal error readable
// after the channel closes, which is the sql.Rows pattern, and a progress
// snapshot, which is the thing a render loop needs and cannot get from a
// channel at all.
//
// # The progress snapshot is why Z-03 has two components
//
// Z-03 §3.1 picks between a scanner (indeterminate) and a bar (determinate) by
// whether the denominator is known. A channel of items cannot say "247 are
// coming"; [Progress] can, and it is the only way the screen gets to draw a
// bar rather than an apology.
type Stream interface {
	// Items returns the channel the provider writes to.
	//
	// It is closed exactly once, when the sync ends for any reason —
	// completion, failure, or cancellation. A caller ranges over it and then
	// reads [Stream.Err]; there is no other correct order.
	Items() <-chan Item

	// Err returns the terminal fault, or nil if the sync completed.
	//
	// It is valid only after Items() is closed. Calling it earlier is a
	// programming error and an implementation may return anything; callers
	// range first, always.
	//
	// A non-nil Err after items have already been delivered is the PARTIAL
	// case: what arrived is kept, the sync_run records partial, and — this is
	// the rule that matters — nothing may be tombstoned, because in a
	// truncated stream "not returned" and "not reached" are indistinguishable
	// (06 §2.4).
	Err() error

	// Progress returns a snapshot of how far the sync has got.
	//
	// It is safe to call from another goroutine at any time, including while
	// the sync is running and after it has finished. That is a hard
	// requirement rather than a convenience: the caller is a Bubble Tea render
	// loop that will read it once per frame while the provider's goroutine
	// writes it, and an implementation that does not make it race-free is a
	// data race in the product's most-watched screen.
	Progress() Progress
}

// Progress is a race-free snapshot of a sync in flight.
//
// It is a value, not a pointer, so a screen cannot hold a handle that changes
// underneath it between two lines of a render function.
type Progress struct {
	// Seen is how many items have been delivered so far.
	Seen int

	// Total is how many the provider expects to deliver, and TotalKnown says
	// whether that expectation exists yet.
	//
	// Two fields rather than a sentinel zero, because zero is a real total —
	// it is the private-profile case — and a screen that reads 0 as "unknown"
	// would draw an indeterminate scanner forever for a sync that has already
	// truthfully finished.
	Total      int
	TotalKnown bool

	// LastAt is when the most recent item arrived.
	//
	// It is here so the caller can implement Z-03's stalled state — ten
	// seconds with no progress drops back to the scanner — without the
	// provider having to know what "stalled" means. Whether a pause is
	// alarming is a screen's judgement about a player's patience, not a
	// provider's judgement about a socket.
	LastAt time.Time

	// Batches is how many upsert batches the consumer has committed.
	//
	// Written by the consumer rather than the provider, and present because
	// PARTIAL's copy — "The 138 that arrived are in your library" — is a claim
	// about what was *written*, not about what was received. A provider that
	// delivered 138 items into a consumer that had committed 120 of them
	// makes that sentence false by eighteen.
	Batches int
}

// Determinate reports whether a progress bar can honestly be drawn.
//
// Z-03 §3.1's rule, in one place, so that two screens cannot disagree about
// when the denominator counts as known.
func (p Progress) Determinate() bool { return p.TotalKnown && p.Total > 0 }

// Stalled reports whether nothing has arrived for at least d.
//
// A sync that has not yet delivered anything is not stalled — it is waiting,
// which is Z-03's indeterminate state and already has its own component. The
// distinction matters because "stalled" would otherwise be true for the whole
// of a slow first round trip, and the screen would flip components for no
// reason the player could see.
func (p Progress) Stalled(now time.Time, d time.Duration) bool {
	if p.Seen == 0 || p.LastAt.IsZero() {
		return false
	}
	return now.Sub(p.LastAt) >= d
}
