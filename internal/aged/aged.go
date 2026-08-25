// Package aged carries the mechanism behind the product's single most
// fragile promise: that any value which came from the network is rendered
// with its age, always.
//
// 07-offline-contract.md §4 states the rule and names why it is the one most
// likely to be lost during a build — dropping the age always makes the layout
// tidier. It is also the rule the product's credibility rests on, because
// Zerado's value proposition is telling a player *not* to buy something, and a
// stale price presented as current is not a cosmetic defect but the product
// giving financial advice from memory.
//
// # Why a wrapper type rather than a FetchedAt field
//
// The spine's shape put the age inside the payload — Metadata.FetchedAt,
// Quote.ObservedAt. That satisfies "carries its own age in the same value",
// and 13-handoffs.md §2 lists exactly that decision as load-bearing.
//
// [Value] keeps the decision and strengthens its enforcement. A field can be
// ignored by a renderer; a wrapper cannot, because there is no way to reach
// the payload without having the age in hand. The distinction matters at the
// only moment it is ever tested — a developer at 1am adding a column to a
// row — and the cost is one .V.
//
// Recorded as a deliberate departure in docs/api/06-divergences-from-the-spine.md
// rather than applied quietly, because 13-handoffs.md §2's whole point is that
// changing one of those decisions by accident is the failure to avoid.
package aged

import "time"

// Value pairs a payload with the moment it was observed.
//
// Every network-derived value in Zerado is carried in a Value. A function that
// returns a bare Quote or a bare Metadata is a bug in the contract, not a
// convenience.
//
// The zero Value is meaningful and is not a trap: At is the zero time, which
// [Value.Known] reports as unknown, and every helper on this type treats an
// unknown age as "do not claim freshness". A value that arrived with no
// timestamp is never presented as current.
type Value[T any] struct {
	// V is the payload.
	//
	// One letter, against Go's usual preference for descriptive names,
	// because the type is a wrapper that appears at every network-derived
	// read and the alternative — q.Value.Current — reads worse than
	// q.V.Current at every one of those sites. The name that matters here is
	// the type's.
	V T

	// At is when the payload was observed at its source.
	//
	// Observed, not fetched into the cache and not written to the database:
	// the age a player is owed is the age of the *fact*, and a cache that
	// re-stamps a value on copy would make a six-month-old price look like it
	// arrived this morning.
	At time.Time
}

// New stamps a payload with an observation time.
func New[T any](v T, at time.Time) Value[T] { return Value[T]{V: v, At: at} }

// Known reports whether this value carries a usable observation time.
//
// An unknown age is not an error — a value can legitimately predate the rule,
// or arrive from a source that does not stamp its data. It is a rendering
// instruction: 07 §4 permits no bare number, so a value whose age is unknown
// is rendered as unknown rather than as fresh.
func (v Value[T]) Known() bool { return !v.At.IsZero() }

// Age returns how old the value is as of now.
//
// It returns zero for an unknown age. Callers must consult [Value.Known]
// before treating a zero age as "just now"; every helper in this package that
// classifies age does so already.
//
// A negative age — a source clock ahead of this machine's — is clamped to
// zero. A value from the future is a clock-skew fact about two machines, not
// information about the data, and rendering "in 3 hours" next to a price would
// be the product reporting somebody else's NTP problem as a feature.
func (v Value[T]) Age(now time.Time) time.Duration {
	if !v.Known() {
		return 0
	}
	if d := now.Sub(v.At); d > 0 {
		return d
	}
	return 0
}

// Freshness is the coarse classification a screen switches on.
//
// Three bands, because the design system has three behaviours: nothing, the
// banner as chrome, and the banner as something the player should act on
// (07 §4.1 — past 90 days an age stops being reassuring and becomes a
// warning). A finer scale would be a scale no screen could use.
type Freshness uint8

const (
	// FreshnessUnknown means the value carries no observation time.
	FreshnessUnknown Freshness = iota

	// FreshnessCurrent means the value is within the caller's stated window.
	FreshnessCurrent

	// FreshnessStale means the value is older than the caller's window and
	// younger than the warning threshold. It renders with its age, as chrome.
	FreshnessStale

	// FreshnessAncient means the value is past the warning threshold. The
	// banner turns amber and names what the player should do.
	FreshnessAncient
)

// WarnAfter is the age past which a value stops being merely stale and starts
// being a warning: ninety days, from 07-offline-contract.md §4.1.
//
// It is a constant of the contract rather than a setting, because a threshold
// a caller can lower is a threshold that gets lowered to make a banner go
// away.
const WarnAfter = 90 * 24 * time.Hour

// Classify places the value in one of the three bands.
//
// window is the caller's own notion of current — a library sync is current for
// hours, a price for minutes. The warning threshold is not the caller's to
// choose.
func (v Value[T]) Classify(now time.Time, window time.Duration) Freshness {
	if !v.Known() {
		return FreshnessUnknown
	}
	age := v.Age(now)
	switch {
	case age >= WarnAfter:
		return FreshnessAncient
	case age > window:
		return FreshnessStale
	default:
		return FreshnessCurrent
	}
}

// Map applies f to the payload and keeps the observation time.
//
// It exists so that a caller transforming a value cannot accidentally
// re-stamp it as new — which is the single way this whole mechanism is
// defeated, and it is defeated by code that looks entirely reasonable.
func Map[A, B any](v Value[A], f func(A) B) Value[B] {
	return Value[B]{V: f(v.V), At: v.At}
}
