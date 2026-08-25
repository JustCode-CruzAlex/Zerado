package cli

import "github.com/JustCode-CruzAlex/Zerado/internal/fault"

// Exit codes. These are API: a script branches on them, and a script cannot be
// recompiled.
//
// The mapping is one-to-one with the fault taxonomy wherever a distinction is
// actionable from a shell, and it deliberately does not collapse the four
// network-ish failures into one. A cron job that syncs nightly wants to retry
// on ExitOffline and alert on ExitUnauthorized, and it cannot do that if both
// are 1.
//
// The numbers are small and contiguous rather than sysexits.h-style, with one
// borrowed convention: [ExitCancelled] is 130, which is what a shell reports
// for a process killed by SIGINT, so a player who pressed Ctrl-C sees the
// number they would have seen anyway.
const (
	// ExitOK is success.
	ExitOK = 0

	// ExitInternal is a defect in Zerado.
	ExitInternal = 1

	// ExitUsage is a malformed invocation: an unknown verb, a missing
	// required argument, --json on an interactive verb.
	ExitUsage = 2

	// ExitOffline is no route or no DNS.
	ExitOffline = 3

	// ExitUnreachable is a timeout or a 5xx.
	ExitUnreachable = 4

	// ExitUnauthorized is a rejected credential. Retrying will not help; the
	// key has to change.
	ExitUnauthorized = 5

	// ExitEmpty is a successful request that returned nothing — the private
	// profile. Distinct from ExitUnauthorized because the key is fine, and
	// distinct from ExitOK because nothing was synced.
	ExitEmpty = 6

	// ExitRateLimited is the provider asking for a wait. A scripted caller
	// may back off and retry; the wait, when the provider stated one, is in
	// the JSON envelope.
	ExitRateLimited = 7

	// ExitNotFound is a named game, provider or setting that does not exist.
	ExitNotFound = 8

	// ExitMalformed is an answer Zerado could not read, or a value the player
	// supplied that a provider would not accept.
	ExitMalformed = 9

	// ExitState is an operation that is legal but not possible right now: a
	// precondition, a conflict, a capability the provider does not have, or a
	// value too stale to act on. Grouped because a shell cannot usefully do
	// anything different about them.
	ExitState = 10

	// ExitCancelled is the player stopping it. 130 by shell convention.
	ExitCancelled = 130
)

// ExitCode maps an error to its process exit status.
//
// It is total over fault.Kinds(): every member of the taxonomy has a code,
// asserted by a test rather than by inspection, so adding a Kind without
// deciding what a script should do about it fails CI.
//
// A nil error is ExitOK. An unclassified error is ExitInternal, never a
// plausible-looking network code — a scripted caller retrying forever because
// an internal panic was reported as a timeout is a worse outcome than a loud
// 1.
func ExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	switch fault.KindOf(err) {
	case fault.KindOffline:
		return ExitOffline
	case fault.KindUnreachable:
		return ExitUnreachable
	case fault.KindUnauthorized:
		return ExitUnauthorized
	case fault.KindRateLimited:
		return ExitRateLimited
	case fault.KindEmpty:
		return ExitEmpty
	case fault.KindNotFound:
		return ExitNotFound
	case fault.KindMalformed:
		return ExitMalformed
	case fault.KindStale, fault.KindUnsupported, fault.KindPrecondition, fault.KindConflict:
		return ExitState
	case fault.KindCancelled:
		return ExitCancelled
	default:
		return ExitInternal
	}
}
