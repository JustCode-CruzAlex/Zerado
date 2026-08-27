// Package i18n carries the one type every other seam needs from the
// catalogue: the key.
//
// ADR-0001 D9 is absolute — no user-facing string literal appears in code.
// That rule has a consequence the seams have to obey rather than merely
// respect: an error, a refusal, a banner label and a CLI message are all
// user-facing text, so none of them may be an English sentence held in a
// struct field. What the seams carry is a Key; what turns a Key into words
// is the catalogue, and the catalogue is a rendering concern that lives
// above every seam in this module.
//
// This package therefore deliberately contains no strings, no catalogue and
// no renderer. It contains the type that makes the rule enforceable at the
// boundary, which is all a contracts package can honestly own.
package i18n

// Key names one entry in the message catalogue.
//
// Keys are dot-separated, lowercase, and read subject-first so that the
// catalogue sorts into the shape of the product:
//
//	fault.steam.key_rejected
//	fault.steam.profile_private
//	sync.done.headline
//	settings.vault.backing.keychain
//
// A Key is not text. It is never rendered, compared for display, lowercased,
// or concatenated with another Key. Two catalogues may render the same Key
// into different languages, and that is the entire point.
//
// The empty Key is invalid. A value that reaches a renderer with an empty Key
// is a programming error, not a blank message: see [Key.Valid].
type Key string

// Valid reports whether k could name a catalogue entry.
//
// It is a shape check, not a lookup — this package does not know which keys
// exist. The catalogue's own loader owns "this key is missing", because that
// is a build-time failure under D9's lint rather than a runtime state.
func (k Key) Valid() bool { return k != "" }

// String returns the key itself, for logs and test failures.
//
// It is spelled out rather than left to Go's default formatting so that a Key
// reaching a log looks like a key and not like a message someone forgot to
// translate.
func (k Key) String() string { return string(k) }

// Args carries the substitutions a message needs.
//
// Named rather than positional, because a translator reordering a sentence is
// the normal case and a positional verb makes that reordering a code change.
// Values are formatted by the catalogue with locale-aware number, currency and
// date rules (golang.org/x/text), so callers pass the value — an int, a
// time.Time, a Money — and never a pre-formatted string.
//
// A pre-formatted string in Args is the most common way D9 is violated while
// appearing to be obeyed: the literal moves out of the render call and into
// the argument, and the lint no longer sees it.
type Args map[string]any
