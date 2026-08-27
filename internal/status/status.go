// Package status holds Zerado's four states and the one rule that decides
// which of them a game is in.
//
// # State versus status
//
// 05-state-machine.md opens by drawing this distinction and keeping it:
// **state** is what the player sees — the four values, the chip, the column,
// the filter; **status** is what the machine stores and commands —
// status_manual, SetStatus, Z-06 Set status. This package is named for the
// machine's half because that is what it is: a stored value and a derivation
// over it.
//
// # The four are ratified
//
// The values, their colours, glyphs, ASCII fallbacks and labels are settled in
// the brand manual §4.3 and are not a design question here. What this package
// owns is the model: one nullable manual value, one derivation when there is
// none, and the invariant that a sync never writes the manual value.
package status

// Status is one of Zerado's four game states.
//
// It is a uint8 with an invalid zero value, deliberately. 05-state-machine.md
// §1 requires that an uninitialised status can never be mistaken for "not
// started", and a string type would make the empty string a plausible-looking
// fifth value that renders as a blank chip. [Unknown] renders as nothing
// because it is never persisted and never reaches a screen — the counts on the
// summary row have to sum to the number shown (§7 rule 1), and a row that
// cannot be classified would break that silently.
type Status uint8

const (
	// Unknown is the zero value. It is never persisted and never rendered.
	Unknown Status = iota

	// NotStarted is the honest default: owned, and no evidence of play.
	//
	// It is also what [Derive] returns for a provider that cannot report
	// playtime at all — a cartridge has no telemetry, and "not started" is
	// the truthful thing to say about a game nothing can observe.
	NotStarted

	// InProgress is the one state a machine may ever assign on its own, and
	// only from playtime the provider reported.
	InProgress

	// Zerado is the earned state — the moment the product exists to create.
	//
	// It is never automatic. 05-state-machine.md §4 records that as a product
	// decision rather than a technical one: a player who is *told* they
	// finished something has not finished it, they have been notified. Phase 2
	// may suggest; only the player may set.
	Zerado

	// Abandoned is a choice, never an inference. Nothing about a save file
	// distinguishes "stopped" from "busy this month".
	Abandoned
)

// String returns the stored, machine-readable form.
//
// These are the values written to status_manual and read back from it, and
// they are the values the CLI accepts and emits. They are API: renaming one
// breaks every existing library file and every script. They are not labels —
// the uppercase words a screen shows come from the catalogue, because
// ADR-0001 D9 leaves no user-facing string in code.
func (s Status) String() string {
	switch s {
	case NotStarted:
		return "not_started"
	case InProgress:
		return "in_progress"
	case Zerado:
		return "zerado"
	case Abandoned:
		return "abandoned"
	default:
		return "unknown"
	}
}

// Valid reports whether s is one of the four.
func (s Status) Valid() bool { return s >= NotStarted && s <= Abandoned }

// Parse converts the stored form back to a Status.
//
// ok is false for anything unrecognised, including the empty string. A caller
// reading a NULL status_manual must not route it through Parse — nil and
// "no opinion" are the same fact and are represented by a nil *Status, which
// is the distinction 05 §5 requires so that "clear override" is expressible at
// all.
func Parse(s string) (Status, bool) {
	switch s {
	case "not_started":
		return NotStarted, true
	case "in_progress":
		return InProgress, true
	case "zerado":
		return Zerado, true
	case "abandoned":
		return Abandoned, true
	default:
		return Unknown, false
	}
}

// All returns the four states in ratified display order.
//
// The order is the order Z-06 lists them and the order the summary row counts
// them, so it is part of the contract rather than a convenience. Returned as a
// fresh slice so a caller cannot reorder the product's own vocabulary.
func All() []Status { return []Status{NotStarted, InProgress, Zerado, Abandoned} }

// Derive computes a status from what the provider reported, for a game whose
// player has never expressed an opinion.
//
// It is the whole automatic half of the model, and it has exactly one
// transition in it: NOT STARTED becomes IN PROGRESS when a provider that can
// report playtime reports some.
//
// canReportPlaytime is Capabilities.Playtime, passed as a bool rather than as
// the Capabilities struct so that this package does not import the provider
// seam — the derivation is a rule about a number and a capability, not about a
// provider, and keeping it importless is what lets the store, the CLI and the
// screens all call it.
//
// The capability argument is the load-bearing one. For Steam, NOT STARTED and
// IN PROGRESS are facts until the player overrides them. For physical, and for
// any store whose API does not expose playtime, all four states are manual
// always — the derivation has no input and returns the honest default. That is
// precisely why physical is modelled as a provider with capabilities rather
// than as an is_physical flag: a boolean would have forced this function to
// special-case one value, and a capability makes it correct for every provider
// that will ever exist, including the ones whose API does not exist yet.
func Derive(playtimeMinutes int, canReportPlaytime bool) Status {
	if !canReportPlaytime {
		return NotStarted
	}
	if playtimeMinutes > 0 {
		return InProgress
	}
	return NotStarted
}

// Effective resolves the model in one line:
//
//	effective = manual ?? derive(playtime, capability)
//
// manual is nil when the player has never expressed an opinion. A sync never
// writes it: mark a game ZERADO, play it for three more hours, and the next
// sync updates playtime and last-played and nothing else. That invariant is
// what makes the product trustworthy, because the alternative is a background
// job silently undoing the one action the product is named after.
//
// The result is never [Unknown] — an invalid manual value is ignored rather
// than propagated, because a corrupt row must still render as a countable game
// and §7 rule 1 requires the counts to sum.
func Effective(manual *Status, playtimeMinutes int, canReportPlaytime bool) Status {
	if manual != nil && manual.Valid() {
		return *manual
	}
	return Derive(playtimeMinutes, canReportPlaytime)
}

// Counts is the pinned summary row: a total, and the four states beneath it.
//
// It is a struct rather than a map because 05 §7 rule 1 is that the four
// counts always sum to the total, and a map has no shape that can be asserted.
// [Counts.Sums] is the assertion, and it is cheap enough to run in a test for
// every filter the product offers.
type Counts struct {
	// Total is the number of games in the set being described.
	Total int

	// NotStarted, InProgress, Zerado, Abandoned count by *effective* status,
	// which is why this type lives beside the derivation and not beside the
	// store.
	NotStarted, InProgress, Zerado, Abandoned int
}

// Sums reports whether the four state counts add to Total.
//
// 05 §7 rule 1: "If the summary says 247 games, the four state counts add to
// 247." A row that cannot be classified does not exist. This is the executable
// form of that sentence, and a store implementation that cannot satisfy it has
// a bug in its GROUP BY rather than an unusual library.
func (c Counts) Sums() bool {
	return c.NotStarted+c.InProgress+c.Zerado+c.Abandoned == c.Total
}

// Of returns the count for one status.
func (c Counts) Of(s Status) int {
	switch s {
	case NotStarted:
		return c.NotStarted
	case InProgress:
		return c.InProgress
	case Zerado:
		return c.Zerado
	case Abandoned:
		return c.Abandoned
	default:
		return 0
	}
}

// Add increments the count for s and the total.
//
// Provided so that a store or a fake building Counts row by row cannot get the
// total out of step with the parts — the only way [Counts.Sums] fails in
// practice.
func (c *Counts) Add(s Status) {
	switch s {
	case NotStarted:
		c.NotStarted++
	case InProgress:
		c.InProgress++
	case Zerado:
		c.Zerado++
	case Abandoned:
		c.Abandoned++
	default:
		return
	}
	c.Total++
}
