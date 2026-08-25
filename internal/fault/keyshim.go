package fault

import "github.com/JustCode-CruzAlex/Zerado/internal/i18n"

// messageOverride returns an Option that applies s only when it is non-empty,
// so a caller can pass an unconditional override without first checking it.
//
// Outcome.MessageKey is a plain string rather than an i18n.Key so that a
// provider package building an Outcome does not have to import i18n for one
// field. The conversion is here, once.
func messageOverride(s string) Option {
	return func(f *Fault) {
		if s != "" {
			f.MessageKey = i18n.Key(s)
		}
	}
}
