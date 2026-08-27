// Package fault is Zerado's error taxonomy: one closed set of kinds, shared by
// every seam, each one distinguishable by the caller because each one renders
// differently on screen.
//
// # Why a taxonomy and not just errors
//
// The offline contract (07-offline-contract.md §5) classifies a failure into a
// screen treatment, and the treatments are not interchangeable:
//
//	no route / DNS   → the OFFLINE banner        — nothing is wrong with anything
//	timeout / 5xx    → the UNREACHABLE banner    — their end, not your key
//	401 / 403        → a credential refusal      — your key, and it is fixable
//	200 + zero items → a private-profile refusal — nothing is broken at all
//
// A private Steam profile is the case this package exists for. It is a
// successful HTTP request, returning a well-formed empty body, from a service
// that is up, using a key that works. Every mechanism that would collapse it
// into "network error" — a bare error string, an HTTP status code, a boolean
// ok — produces a screen that tells the player their connection failed when
// their connection is fine, and sends them to debug the wrong thing. Worse,
// 06-data-seams.md §2.4 requires that this exact case never tombstones a row:
// a caller that cannot tell "the provider returned nothing" from "the provider
// could not be reached" cannot honour that rule.
//
// # The two rules that keep it honest
//
//  1. A Fault carries an i18n.Key, never a sentence. ADR-0001 D9 forbids a
//     user-facing string literal in code, and an error message is user-facing
//     text. The Kind decides which key; the catalogue decides which words.
//  2. A Fault's wrapped cause is for logs, never for a screen. It is where a
//     URL with a query string — and therefore an API key — would otherwise
//     reach a person's terminal. [Fault.Error] is a developer-facing string
//     and is documented as never renderable.
package fault

// Kind is the closed set of ways an operation can fail in Zerado.
//
// The set is closed on purpose. A caller switching on Kind is writing the
// screen treatment for every case the product has, and a new Kind is therefore
// a design decision with a screen behind it — not a convenience for a package
// that met an error it had not thought about. Anything genuinely unclassified
// is [KindInternal], which renders as the fatal screen (Z-11) rather than as a
// new species of banner.
type Kind uint8

const (
	// KindUnknown is the zero value and is never valid.
	//
	// It is deliberately not "internal": a Fault built without a Kind is a
	// construction bug, and giving the zero value a plausible meaning would
	// let that bug ship as a screen nobody questions.
	KindUnknown Kind = iota

	// KindOffline means the machine could not begin the request: no route,
	// or DNS did not resolve.
	//
	// This is the only Kind Zerado may report without having reached anyone,
	// and it is still evidence-driven — 07 §5 forbids probing, so a request
	// has to have been attempted and failed. Screen: the OFFLINE banner.
	KindOffline

	// KindUnreachable means the request was made and the far end did not
	// usefully answer: a timeout, a 5xx, a connection reset mid-body.
	//
	// Distinct from KindOffline because the copy is different and the blame is
	// different: "Steam didn't answer. Not your key — their end, or the
	// connection."
	KindUnreachable

	// KindUnauthorized means the credential was rejected: 401, 403, or a
	// provider's own equivalent.
	//
	// Distinct from KindUnreachable because it is the one network failure the
	// player can actually fix, and the screen routes them to where they fix it.
	KindUnauthorized

	// KindRateLimited means the provider asked Zerado to slow down.
	//
	// It carries a RetryAfter whenever the provider supplied one. It is not
	// KindUnreachable: the service is up, the key is good, and waiting works —
	// which makes it the one failure a caller may retry without asking.
	KindRateLimited

	// KindEmpty means the provider answered successfully and returned nothing.
	//
	// For Steam this is the private-profile case (Z-03 state 11), and it is a
	// refusal rather than an empty result set: 06 §2.4 and Z-03 §8.2 both
	// forbid treating it as "the library was removed". No UpsertBatch, no
	// tombstoning, and the copy may then truthfully say "Your library is
	// unchanged — nothing was lost."
	//
	// Named for what was observed rather than for one provider's reason for
	// it, because GOG and PlayStation will have their own reasons for the same
	// observation and this Kind must survive them.
	KindEmpty

	// KindNotFound means the thing asked for does not exist at the source.
	//
	// On the metadata seam this is routine and is not a failure of anything:
	// a hand-added cartridge that IGDB has never heard of produces
	// KindNotFound, and the detail view renders the designed no-metadata
	// composition (06 §3.1). On the store seam it is a real defect, because
	// a screen asked for a GameID it was handed.
	KindNotFound

	// KindMalformed means the answer arrived and Zerado could not read it.
	//
	// Separated from KindUnreachable because the remedies are opposite: an
	// unreachable provider is retried, a malformed answer is a provider whose
	// contract changed and retrying it forever is how a client hammers an API
	// that will never satisfy it.
	KindMalformed

	// KindStale means a cached value exists but is too old for the caller's
	// purpose, and the caller declined to act on it.
	//
	// Zerado's read path never fetches (see the store seam), so staleness is
	// normally shown rather than refused — a stale price renders with its age.
	// This Kind is for the callers that must not act on an old value at all,
	// the Phase 3 watchlist verdict being the case that motivated it: telling
	// somebody to buy at a price from June is worse than telling them nothing.
	KindStale

	// KindUnsupported means the provider was asked to do something its
	// Capabilities say it cannot.
	//
	// It should be nearly unreachable by construction — physical is not a
	// Syncer, so there is no Sync to call — and that is the point: it is the
	// backstop for a caller that switched on ProviderID instead of reading
	// Capabilities, and it names that mistake instead of returning a plausible
	// zero value.
	KindUnsupported

	// KindPrecondition means the operation is legal but the state is not ready
	// for it.
	//
	// The load-bearing case is absence reconciliation: only a sync whose run
	// status is ok may tombstone anything (06 §2.4), so a store asked to
	// reconcile against a partial, failed or cancelled run must refuse rather
	// than comply.
	KindPrecondition

	// KindConflict means the write collided with a uniqueness rule the store
	// enforces — (provider_id, provider_ref) being the one Phase 1 has.
	KindConflict

	// KindCancelled means the player stopped it: Esc during a sync, q during a
	// check, SIGINT at the CLI.
	//
	// It is never rendered as an error. Z-03's CANCELLED state has no cue, no
	// red and no apology, because a cancelled sync did exactly what it was
	// told and what arrived before the cancel is kept.
	KindCancelled

	// KindInternal is a defect in Zerado: a broken invariant, a corrupt
	// database, a migration from a newer binary.
	//
	// It is the only Kind that is allowed to be vague to the player, because
	// it is the only one where the honest answer is "this is ours". It routes
	// to Z-11 Fatal error, which names both versions and stops.
	KindInternal
)

// String returns the Kind's stable machine name.
//
// These names are API. They appear in the CLI's --json envelope and in logs,
// so a script may switch on them; renaming one is a breaking change to the CLI
// surface (see internal/cli). They are not messages and are never rendered to
// a player — that is what MessageKey is for.
func (k Kind) String() string {
	switch k {
	case KindOffline:
		return "offline"
	case KindUnreachable:
		return "unreachable"
	case KindUnauthorized:
		return "unauthorized"
	case KindRateLimited:
		return "rate_limited"
	case KindEmpty:
		return "empty"
	case KindNotFound:
		return "not_found"
	case KindMalformed:
		return "malformed"
	case KindStale:
		return "stale"
	case KindUnsupported:
		return "unsupported"
	case KindPrecondition:
		return "precondition"
	case KindConflict:
		return "conflict"
	case KindCancelled:
		return "cancelled"
	case KindInternal:
		return "internal"
	default:
		return "unknown"
	}
}

// Valid reports whether k is a member of the taxonomy.
func (k Kind) Valid() bool { return k > KindUnknown && k <= KindInternal }

// Retryable reports whether retrying the same operation unchanged could
// plausibly succeed.
//
// It answers a question the CLI and the sync coordinator both have to ask, and
// it is deliberately conservative: a Kind is retryable only when nothing has
// to change first. KindUnauthorized is not retryable because the key is wrong
// until the player changes it; KindEmpty is not, because the profile is
// private until the player changes it; KindMalformed is not, because the
// provider's contract moved and hammering it will not move it back.
//
// Retryable never means "retry automatically". Zerado retries when the player
// presses r; this reports whether offering r is honest.
func (k Kind) Retryable() bool {
	switch k {
	case KindOffline, KindUnreachable, KindRateLimited:
		return true
	default:
		return false
	}
}
