package cli

import (
	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
)

// Envelope is the shape of every --json response.
//
// One shape for success and failure, because a caller that has to parse two
// different top-level structures depending on an exit code is a caller that
// will get it wrong on the day it matters. The discriminator is [Envelope.OK],
// and exactly one of Data and Error is present.
//
// # The stability policy, stated so it can be held to
//
//   - [APIVersion] appears in every envelope and is the contract's major
//     version.
//   - Fields are added, never removed and never repurposed. A consumer that
//     ignores unknown fields keeps working across every minor change.
//   - A field's meaning is fixed. Changing what "unchanged" counts is a
//     breaking change even though the JSON looks identical, and it bumps the
//     version.
//   - Error kinds are the fault taxonomy's stable names. Adding a kind is
//     additive; a consumer that does not recognise one must treat it as a
//     failure rather than as success, which is why [Envelope.OK] exists
//     separately from the kind.
//   - Exit codes and kinds move together and are documented together.
type Envelope struct {
	// API is the contract's major version. Always present.
	API int `json:"api"`

	// OK is the discriminator. A consumer branches on this, not on the
	// presence of a key and not on the exit code, which it may not have.
	OK bool `json:"ok"`

	// Data is the verb's own payload, present only when OK.
	Data any `json:"data,omitempty"`

	// Error describes the failure, present only when not OK.
	Error *ErrorBody `json:"error,omitempty"`
}

// ErrorBody is the machine-readable half of a failure.
//
// It carries no message. That is not an omission: a message is user-facing
// text, it is rendered from the catalogue in the player's language, and a
// pre-rendered English sentence in a JSON field is exactly the D9 violation
// that is hardest to spot. A consumer that wants a sentence renders one from
// the kind; a consumer that wants to branch uses the kind, which is what it
// should have been using anyway.
type ErrorBody struct {
	// Kind is the taxonomy's stable machine name — "offline", "unauthorized".
	Kind string `json:"kind"`

	// Op names the operation, for a bug report.
	Op string `json:"op,omitempty"`

	// Subject is the provider's display name, when the failure was about one.
	Subject string `json:"subject,omitempty"`

	// MessageKey is the catalogue key that would render this, so a consumer
	// with the catalogue can produce the same sentence Zerado would.
	MessageKey string `json:"message_key,omitempty"`

	// RetryAfterSeconds is the provider's own wait instruction, when it gave
	// one. Seconds rather than a duration string, because every JSON consumer
	// can compare a number and not all of them can parse "1m30s".
	RetryAfterSeconds int `json:"retry_after_seconds,omitempty"`
}

// Ok builds a success envelope.
func Ok(data any) Envelope { return Envelope{API: APIVersion, OK: true, Data: data} }

// Err builds a failure envelope from any error.
//
// It never carries the wrapped cause. The cause may contain a URL with the
// player's API key in its query string, and a JSON envelope is the single most
// likely thing in this product to be piped into a file and attached to a bug
// report.
func Err(err error) Envelope {
	k := fault.KindOf(err)
	body := &ErrorBody{
		Kind:       k.String(),
		MessageKey: fault.MessageKeyOf(err).String(),
		Subject:    fault.SubjectOf(err),
	}
	if d := fault.RetryAfterOf(err); d > 0 {
		body.RetryAfterSeconds = int(d.Seconds())
	}
	var f *fault.Fault
	if asFault(err, &f) {
		body.Op = f.Op
	}
	return Envelope{API: APIVersion, OK: false, Error: body}
}
