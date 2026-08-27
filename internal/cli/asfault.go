package cli

import (
	"errors"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
)

// asFault is errors.As specialised to *fault.Fault, kept in its own file so
// the envelope reads as a data shape rather than as error plumbing.
func asFault(err error, target **fault.Fault) bool { return errors.As(err, target) }
