package fault

import "github.com/JustCode-CruzAlex/Zerado/internal/i18n"

// Treatment is how a Kind is allowed to appear on screen.
//
// It exists because "each Kind renders differently" is easy to assert and easy
// to lose: a screen that has three failure paths and one error box will
// quietly collapse the taxonomy back into a single treatment within a month.
// Naming the treatment on the Kind makes the collapse visible in review and
// testable in a table.
//
// It decides the *class* of presentation, never the composition. The design
// system owns the banner (01-design-system.md §12) and each screen owns its
// own refusal block; this only says which of them a given failure is entitled
// to.
type Treatment uint8

const (
	// TreatmentNone is for failures that are not shown as failures at all.
	//
	// KindCancelled is the whole membership: the player stopped it, the
	// screen says what was kept, and nothing about that is an error.
	TreatmentNone Treatment = iota

	// TreatmentBanner is the degrade banner — one row at the top of the body,
	// naming what is unavailable and how stale what is shown is.
	//
	// It is for ambient failure: something in the background could not be
	// refreshed and the screen underneath is still completely usable.
	TreatmentBanner

	// TreatmentRefusal is a block on the screen the player is standing on,
	// naming what happened, why, and the next action.
	//
	// It is for a failure of something the player just asked for. 07 §2 draws
	// exactly this line: an action the player took is owed an answer, an
	// ambient feature merely states its condition.
	TreatmentRefusal

	// TreatmentDesignedEmpty is not a failure presentation at all.
	//
	// KindNotFound on the metadata seam renders the designed no-metadata
	// composition (06 §3.1). This is the difference between a product that
	// works without IGDB and a product that is broken without IGDB, and it is
	// recorded as a treatment so that "no metadata" can never be routed to an
	// error banner by a well-meaning caller.
	TreatmentDesignedEmpty

	// TreatmentFatal is Z-11: the program cannot continue.
	TreatmentFatal
)

// Treatment returns how this Kind is permitted to appear.
//
// The mapping is here, once, rather than in each screen, because the property
// worth protecting is that the twelve Kinds do not silently converge on two
// treatments. A screen may choose different words; it may not choose a
// different class.
func (k Kind) Treatment() Treatment {
	switch k {
	case KindCancelled:
		return TreatmentNone
	case KindOffline, KindUnreachable, KindStale, KindRateLimited:
		return TreatmentBanner
	case KindUnauthorized, KindEmpty, KindMalformed, KindPrecondition, KindConflict, KindUnsupported:
		return TreatmentRefusal
	case KindNotFound:
		return TreatmentDesignedEmpty
	default:
		return TreatmentFatal
	}
}

// MessageKey returns the catalogue key a Kind renders through by default.
//
// Providers override it through [WithMessage] when they have better copy for
// their own case — Steam's private-profile refusal names a specific Steam
// privacy setting, which no provider-agnostic KindEmpty entry could. The
// default exists so that a provider which supplies no copy still renders a
// sentence rather than a blank.
//
// The keys are subject-agnostic: the provider's name arrives as the {subject}
// argument, which is what lets one entry serve Steam today and GOG later, in
// every language, with no new key.
func (k Kind) MessageKey() i18n.Key {
	switch k {
	case KindOffline:
		return "fault.offline"
	case KindUnreachable:
		return "fault.unreachable"
	case KindUnauthorized:
		return "fault.unauthorized"
	case KindRateLimited:
		return "fault.rate_limited"
	case KindEmpty:
		return "fault.empty"
	case KindNotFound:
		return "fault.not_found"
	case KindMalformed:
		return "fault.malformed"
	case KindStale:
		return "fault.stale"
	case KindUnsupported:
		return "fault.unsupported"
	case KindPrecondition:
		return "fault.precondition"
	case KindConflict:
		return "fault.conflict"
	case KindCancelled:
		return "fault.cancelled"
	default:
		return "fault.internal"
	}
}

// Kinds returns every member of the taxonomy, in declaration order.
//
// It is DERIVED from the iota range rather than hand-listed, and that is the
// whole point of it.
//
// The first version of this function returned a literal slice. Every totality
// test iterated that slice, and nothing asserted the slice covered the range —
// so a Kind inserted mid-block was valid, constructible, and invisible to
// every test that existed to catch exactly that. It silently claimed the
// machine name "unknown" on a surface the CLI documents as stable API, and
// rendered as TreatmentFatal with generic internal copy. CI stayed green.
//
// Deriving it removes the second source of truth instead of adding a test to
// reconcile the two. [KindInternal] is the range's upper bound here and in
// [Kind.Valid], so the two cannot disagree.
//
// The totality tests then have something real to iterate: TestTaxonomyIsTotal
// asserts every kind has a distinct machine name, a treatment and a catalogue
// key, and TestNoKindFallsThroughToTheDefaults asserts none of them silently
// inherits the switch defaults. A new Kind fails both until it has been
// thought about everywhere, which is the enforcement this taxonomy needs to
// still be a taxonomy in a year.
func Kinds() []Kind {
	out := make([]Kind, 0, int(KindInternal))
	for k := KindOffline; k <= KindInternal; k++ {
		out = append(out, k)
	}
	return out
}
