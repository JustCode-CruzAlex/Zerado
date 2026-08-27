package fault

import "time"

// Transport is what the network layer observed, in the only four flavours the
// offline contract's classifier distinguishes.
//
// It is an enum rather than an error because this package must not import net
// or net/http. 07-offline-contract.md §7.3 makes "no net/http import outside
// the provider packages" a review rule, and a shared classifier that reached
// for *net.DNSError would be the first, most reasonable-looking violation of
// it — after which the rule is a comment.
//
// So the provider, which is already the only package allowed to know about
// HTTP, does the probing and reports a verdict. The classifier stays a pure
// function over that verdict, which is also what makes every branch of it
// reachable in a test with no network at all.
type Transport uint8

const (
	// TransportOK means a response was received. Its status code decides the
	// rest.
	TransportOK Transport = iota

	// TransportNoRoute means the request could not leave the machine: no
	// route to host, network unreachable, interface down.
	TransportNoRoute

	// TransportDNS means the name did not resolve.
	//
	// Folded with TransportNoRoute at the treatment level — both are the
	// OFFLINE banner — but kept distinct here because they are distinct
	// facts, and a log that says which one is a log worth having.
	TransportDNS

	// TransportTimeout means the request was sent and nothing useful came
	// back in time. This includes a body that stalled mid-stream.
	TransportTimeout

	// TransportReset means the connection was established and then broken.
	TransportReset
)

// Outcome is everything the classifier needs to decide what a player is told.
//
// A provider fills it in and calls [Classify]. It is a struct rather than a
// long parameter list because the fields are not interchangeable and three of
// them are zero in the common case.
type Outcome struct {
	// Op is the operation name for logs — "steam.Sync".
	Op string

	// Subject is the player-facing name of the provider — "Steam".
	Subject string

	// Transport is the network-layer verdict.
	Transport Transport

	// Status is the HTTP status code, or 0 when no response arrived.
	//
	// A plain int precisely so this package needs no HTTP import. A non-HTTP
	// provider reports 200 for success and 0 for "not applicable", and the
	// classifier is unchanged.
	Status int

	// Items is how many items the operation yielded before it ended.
	//
	// This is the field that separates a private profile from a working sync
	// of a small library, and it is why the classifier takes an Outcome
	// rather than a response: a 200 that yielded nothing and a 200 that
	// yielded 247 games are the same response and completely different facts.
	Items int

	// RetryAfter is the provider's own wait instruction, when it sent one.
	RetryAfter time.Duration

	// Cause is the underlying error, for logs. It is never rendered.
	Cause error

	// MessageKey optionally overrides the Kind's default catalogue entry, for
	// a provider that has better copy for its own case.
	MessageKey string
}

// Classify turns an Outcome into an error, implementing the decision tree in
// 07-offline-contract.md §5.
//
// A nil return means success. Whether an empty result is itself a refusal is
// the caller's business, and is answered by [ClassifySync] for the one
// operation where it is — rather than by this function guessing.
//
// # Why it returns error and not *Fault
//
// It returned *Fault once, and that made the natural provider spelling a trap:
//
//	func (s *steam) sync(...) error { return fault.Classify(o) }
//
// A successful outcome returns a nil *Fault, which becomes a NON-nil error
// interface holding a typed nil. [KindOf] then finds err != nil, errors.As
// succeeds, the embedded pointer is nil — and reports KindInternal, so a
// completely successful sync renders the fatal screen.
//
// It fails in the safe direction, loudly and in the "this is ours" direction,
// which is the taxonomy working. But safe is not correct, and the repair is
// the same one this package has now made twice elsewhere: make the guarantee a
// property of the signature rather than of how carefully every caller holds
// it. Returning error means the nil is an untyped nil at every call site,
// including the ones nobody has written yet.
//
// Callers that need the classification ask [KindOf], [Is], [SubjectOf] or
// [RetryAfterOf], which is the documented way to read a fault and works
// through wrapping.
//
// The ordering of the branches is part of the contract:
//
//  1. transport failures first — a rejected key is meaningless if the request
//     never arrived;
//  2. then 401/403, because a credential verdict outranks anything derivable
//     from the body;
//  3. then 404 and 429, each of which has its own screen;
//  4. then 5xx, then any other 4xx, which is our bug and not their outage.
func Classify(o Outcome) error {
	opts := []Option{WithSubject(o.Subject), WithCause(o.Cause), messageOverride(o.MessageKey)}

	switch o.Transport {
	case TransportNoRoute, TransportDNS:
		return New(KindOffline, o.Op, opts...)
	case TransportTimeout, TransportReset:
		return New(KindUnreachable, o.Op, opts...)
	}

	switch {
	case o.Status == 401 || o.Status == 403:
		return New(KindUnauthorized, o.Op, opts...)
	case o.Status == 404:
		return New(KindNotFound, o.Op, opts...)
	case o.Status == 429:
		return New(KindRateLimited, o.Op, append(opts, WithRetryAfter(o.RetryAfter))...)
	case o.Status >= 500:
		return New(KindUnreachable, o.Op, opts...)
	case o.Status >= 400:
		// A 4xx that is not one of the three above means Zerado sent
		// something the provider would not accept. That is our bug, not the
		// player's network, and it must not be dressed as one.
		return New(KindMalformed, o.Op, opts...)
	case o.Status >= 200 && o.Status < 300:
		return nil
	case o.Status == 0:
		// No response and no transport verdict is a provider that did not
		// fill in its Outcome. Refusing to guess is the whole point of the
		// taxonomy.
		return New(KindInternal, o.Op, opts...)
	default:
		return New(KindMalformed, o.Op, opts...)
	}
}

// ClassifySync is Classify with the one rule that applies only to a library
// fetch: a successful response carrying zero items is a refusal, not an empty
// result.
//
// This is the private-profile case, and it is a separate function rather than
// a branch inside Classify because emptiness is a refusal for exactly one
// operation. A metadata lookup that finds nothing renders the designed
// no-metadata composition; a price lookup that finds no quote is simply no
// quote. Only a library sync may conclude, from an empty answer, that the
// player has to go and change a setting somewhere else.
//
// Its consequence is load-bearing rather than cosmetic. 06-data-seams.md §2.4
// forbids tombstoning on any run that is not ok, and Z-03 §8.2 spells out why:
// a 247-game library whose owner has just made their profile private would
// otherwise be deleted by a naive "the provider's view is the truth" upsert,
// after which the screen would honestly report a catastrophe it had caused.
// Returning a Fault here is what makes the ratified copy — "Your library is
// unchanged — nothing was lost." — true at the moment it is printed.
func ClassifySync(o Outcome) error {
	if f := Classify(o); f != nil {
		return f
	}
	if o.Items == 0 {
		return New(KindEmpty, o.Op,
			WithSubject(o.Subject),
			WithCause(o.Cause),
			messageOverride(o.MessageKey),
		)
	}
	return nil
}
