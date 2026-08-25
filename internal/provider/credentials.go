package provider

// CredentialField is one field on Z-02, declared by the provider that needs
// it.
//
// The screen renders from this slice and knows nothing else. Adding GOG is:
// implement the interfaces, declare the fields, register — zero screens, zero
// routes, zero schema changes. That is the test of whether this seam is right,
// and it is why Z-02 is named "connect a store" rather than "connect Steam".
type CredentialField struct {
	// Key names the field to the provider and to the Vault — "api_key",
	// "steam_id". It is never rendered.
	Key string

	// LabelKey, HelpKey name catalogue entries.
	LabelKey string
	HelpKey  string

	// HelpURL is where the player goes to get this value.
	//
	// A URL rather than a key, because it is not language: it is an address,
	// and a translated address would be a broken one.
	HelpURL string

	// Secret decides the field's destination, not merely its masking.
	//
	// True means the value is masked on screen AND written to the Vault,
	// never to library.db. False means it is an identifier — the Steam ID, a
	// GOG username — and it lives in provider_connection.account_ref.
	//
	// One boolean carrying both facts is deliberate: a masked field that was
	// nonetheless stored in the library file would be the exact failure
	// 06 §5.4 is written to prevent, and the only way to have that failure
	// here is to type two booleans and disagree with yourself.
	Secret bool

	// Validate is the provider's own local check, run before the network is
	// touched at all.
	//
	// It returns a catalogue key naming what is wrong, or "" when the value is
	// acceptable. A key rather than an error because the result is rendered
	// under the field, in the player's language.
	//
	// Z-02 §8.2 constrains what it may reject, and the constraint is
	// deliberately severe: Steam's Validate rejects empty and whitespace-only
	// and nothing else — no length, no character class, no format. A format
	// assertion is a claim about somebody else's API, and if it is wrong, or
	// right today and wrong next year, the product refuses a valid credential
	// on its own authority. The provider's verdict is the only one that
	// matters, and Z-02 §8.1 goes and gets it.
	//
	// nil means nothing is checked locally, which is a legitimate choice and
	// not an omission.
	Validate func(string) string
}

// Credentials is a provider's own key set, keyed by [CredentialField.Key].
//
// It never leaves the machine except to the service it belongs to. Three rules
// hold it to that, and each is a property of this type rather than a habit:
//
//   - it is passed only to the provider that declared the fields;
//   - it is never placed in a [fault], a log line or an [Item];
//   - it never crosses the Phase 4 sync boundary, which is why
//     internal/devicesync has no field that could carry one.
//
// A map rather than a struct, for the same reason [Entry] is: the field set is
// the provider's to declare.
type Credentials map[string]string

// Secrets returns the keys of fields the provider declared as secret.
//
// It exists so that a caller assembling Credentials from the Vault, or a log
// sink redacting them, can do so from the provider's own declaration rather
// than from a list of names it maintains separately — the second list being
// the one that goes stale the day a provider adds a field.
func Secrets(c Capabilities) []string {
	var out []string
	for _, f := range c.Credentials {
		if f.Secret {
			out = append(out, f.Key)
		}
	}
	return out
}

// Missing returns the keys of required credential fields that are absent or
// blank in creds.
//
// Every declared field is required; a provider that wants an optional one
// declares it as an entry field on a form instead. Returned in declaration
// order so Z-02 can focus the first offender, which is what its empty-submit
// state does.
func Missing(c Capabilities, creds Credentials) []string {
	var out []string
	for _, f := range c.Credentials {
		if isBlank(creds[f.Key]) {
			out = append(out, f.Key)
		}
	}
	return out
}

// isBlank reports whether s is empty or only whitespace.
//
// Spelled out rather than deferring to strings.TrimSpace so that the rule is
// visible at the one place the contract states it: empty and whitespace-only
// are the two things every provider's Validate may reject.
func isBlank(s string) bool {
	for _, r := range s {
		switch r {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return false
		}
	}
	return true
}
