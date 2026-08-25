package fault

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/i18n"
)

// Fault is the one error type every Zerado seam returns.
//
// Callers are not expected to type-assert it. They ask [Is] or [KindOf], both
// of which see through wrapping, so a package may wrap a Fault in context
// without breaking the caller's switch. The struct is exported because a
// provider package has to build one, not because a screen should read one.
//
// # What a Fault may and may not carry
//
// It may carry: a Kind, a machine-readable operation name, the subject the
// message is about, an i18n key, a retry hint, and a wrapped cause.
//
// It may not carry a sentence. ADR-0001 D9 forbids user-facing string
// literals in code, and D9's lint cannot see a literal that has been parked in
// an error struct. The Kind chooses the treatment, the MessageKey chooses the
// words, and the catalogue owns every word.
//
// It may not carry a credential, a URL with a query string, or anything else
// derived from the player's keys. That rule is about Cause: an HTTP client's
// own *url.Error stringifies the full URL, and a Steam URL carries the API key
// in it. See [Fault.Error].
type Fault struct {
	// Kind is the taxonomy member. It is required; a Fault with KindUnknown is
	// a construction bug and [New] refuses to make one.
	Kind Kind

	// Op names the operation that failed, in package.Method form —
	// "steam.Sync", "store.UpsertBatch", "vault.Get".
	//
	// It is for logs, for the CLI's --json envelope and for test failures. It
	// is never rendered to a player: a person who has just been told their
	// Steam profile is private does not need to know which Go method noticed.
	Op string

	// Subject is what the message is about, in the player's terms — the value
	// a Provider's Display() returns, so "Steam" and not "steam".
	//
	// It is a substitution for the catalogue, not a message. Copy such as
	// "Steam didn't answer" is one catalogue entry with a {subject} argument,
	// which is what makes the same entry correct for GOG on the day GOG
	// exists, in every language, with no new key.
	Subject string

	// MessageKey names the catalogue entry that renders this fault.
	//
	// A Kind has a default key (see [Kind.MessageKey]); a provider may
	// override it when it has genuinely better copy for its own case. The
	// override is why this is a field rather than a pure function of Kind:
	// Steam's private-profile refusal names a specific Steam privacy setting,
	// which no generic KindEmpty copy could.
	MessageKey i18n.Key

	// RetryAfter is the provider's own instruction to wait, when it gave one.
	//
	// Zero means "not stated", never "retry immediately". Only meaningful
	// with KindRateLimited, and it is the reason KindRateLimited is a distinct
	// Kind rather than a flavour of KindUnreachable.
	RetryAfter time.Duration

	// Cause is the underlying error, for logs and for errors.Is/As by the
	// package that built the Fault.
	//
	// It is never rendered. See [Fault.Error] for why that is a redaction
	// requirement and not a style preference.
	Cause error
}

// New builds a Fault.
//
// It panics on an invalid Kind. That is deliberate and it is the only panic in
// these contracts: an unclassified failure is a hole in the taxonomy, the hole
// is a screen nobody designed, and a construction-site panic in a test run is
// enormously cheaper than shipping a banner that says nothing.
func New(kind Kind, op string, opts ...Option) *Fault {
	if !kind.Valid() {
		panic(fmt.Sprintf("fault.New: invalid kind %d for op %q", kind, op))
	}
	f := &Fault{Kind: kind, Op: op, MessageKey: kind.MessageKey()}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Option customises a Fault at construction.
//
// Options rather than a wide constructor, because the common call is two
// arguments and the rare call is five, and a five-argument constructor with
// three zero values at every site is how Subject silently stops being filled.
type Option func(*Fault)

// WithSubject sets the player-facing subject — a provider's Display() value.
func WithSubject(s string) Option { return func(f *Fault) { f.Subject = s } }

// WithMessage overrides the Kind's default catalogue key.
func WithMessage(k i18n.Key) Option { return func(f *Fault) { f.MessageKey = k } }

// WithRetryAfter records a provider's own wait instruction.
func WithRetryAfter(d time.Duration) Option { return func(f *Fault) { f.RetryAfter = d } }

// WithCause attaches the underlying error for logs and errors.Is.
func WithCause(err error) Option { return func(f *Fault) { f.Cause = err } }

// Error returns a developer-facing description. It is a log line, not a
// message.
//
// It deliberately does not include Cause's text. An *url.Error from net/http
// stringifies the whole URL, and Steam's URLs carry the player's API key as a
// query parameter — so an error string that concatenates its cause is a
// credential-disclosure path that runs straight through the one screen a
// frustrated player is most likely to screenshot and paste into a bug report.
//
// The cause is still reachable by [errors.Unwrap] and by errors.Is/As, which
// is where a log sink that has been told it may see it can get it.
func (f *Fault) Error() string {
	if f == nil {
		return "<nil fault>"
	}
	if f.Subject != "" {
		return fmt.Sprintf("%s: %s: %s", f.Op, f.Subject, f.Kind)
	}
	return fmt.Sprintf("%s: %s", f.Op, f.Kind)
}

// Unwrap returns the wrapped cause, so errors.Is and errors.As reach it.
func (f *Fault) Unwrap() error {
	if f == nil {
		return nil
	}
	return f.Cause
}

// Is reports whether err is, or wraps, a Fault of the given Kind.
//
// This is the call every caller makes. It is a function rather than a method
// so that no caller has to have already established that it holds a Fault, and
// it walks the wrap chain so that a package adding context to a provider's
// error does not silently reclassify it.
func Is(err error, k Kind) bool { return KindOf(err) == k }

// KindOf returns the Kind of err, walking the wrap chain.
//
// It returns KindInternal for a non-nil error that is not a Fault, because an
// error that reached a screen without ever having been classified is exactly
// the case Z-11 exists for — and reporting it as, say, KindUnreachable would
// be a guess about somebody else's failure.
//
// It returns KindUnknown for a nil error, which is the one place KindUnknown
// is a legitimate value: it means "nothing failed".
func KindOf(err error) Kind {
	if err == nil {
		return KindUnknown
	}
	var f *Fault
	if errors.As(err, &f) && f != nil {
		return f.Kind
	}
	return KindInternal
}

// RetryAfterOf returns the wait a provider asked for, or zero if none was
// stated.
func RetryAfterOf(err error) time.Duration {
	var f *Fault
	if errors.As(err, &f) && f != nil {
		return f.RetryAfter
	}
	return 0
}

// MessageKeyOf returns the catalogue key that renders err.
//
// A non-Fault error yields the internal-failure key rather than an empty one,
// so there is no path by which an unclassified error renders as a blank
// screen.
func MessageKeyOf(err error) i18n.Key {
	var f *Fault
	if errors.As(err, &f) && f != nil && f.MessageKey.Valid() {
		return f.MessageKey
	}
	return KindInternal.MessageKey()
}

// MarshalJSON renders a Fault without its cause.
//
// [Fault.Error] already redacts, and the CLI envelope is built field by field
// from the accessors rather than from this type — both paths are tested. This
// exists because neither of those is a property of the *type*: a log sink, a
// future debug verb, or anything else that reaches for json.Marshal on a Fault
// would reopen the exact path the rest of this package closes. An *url.Error
// has an exported URL field, and a Steam URL carries the player's API key in
// its query string.
//
// So the redaction is moved into the type, where it holds for every caller
// including the ones that have not been written yet. The cause stays reachable
// through errors.Is and errors.As for a sink that has been told it may see it.
//
// The shape matches the CLI's ErrorBody deliberately: one machine-readable
// description of a failure, in one shape, wherever it is serialised.
func (f *Fault) MarshalJSON() ([]byte, error) {
	if f == nil {
		return []byte("null"), nil
	}
	out := struct {
		Kind              string `json:"kind"`
		Op                string `json:"op,omitempty"`
		Subject           string `json:"subject,omitempty"`
		MessageKey        string `json:"message_key,omitempty"`
		RetryAfterSeconds int    `json:"retry_after_seconds,omitempty"`
	}{
		Kind:       f.Kind.String(),
		Op:         f.Op,
		Subject:    f.Subject,
		MessageKey: f.MessageKey.String(),
	}
	if f.RetryAfter > 0 {
		out.RetryAfterSeconds = int(f.RetryAfter.Seconds())
	}
	return json.Marshal(out)
}

// SubjectOf returns the player-facing subject carried by err, if any.
func SubjectOf(err error) string {
	var f *Fault
	if errors.As(err, &f) && f != nil {
		return f.Subject
	}
	return ""
}
