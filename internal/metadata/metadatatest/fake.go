// Package metadatatest holds a metadata provider backed by a map, so the
// enrichment path can be exercised with no network and no source.
package metadatatest

import (
	"context"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/aged"
	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/metadata"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
)

// Fake answers from a map keyed by title and platform.
//
// It keys on title and platform rather than on a store identifier, and that is
// the point rather than a convenience: it is the shape a source that has never
// heard of Steam would have, which is the shape the seam has to support for a
// cartridge. A fake keyed on appid would have quietly agreed with a
// Steam-shaped seam and proved nothing.
type Fake struct {
	Ident   provider.ID
	Records map[string]metadata.Metadata
	At      time.Time
	Credit  metadata.Attribution

	// Fail, when set, is returned instead of a lookup — so a test can
	// exercise the difference between a source that is down (a banner) and a
	// source that has never heard of a game (a designed empty).
	Fail error
}

// New returns a Fake with a fixed observation time.
func New() *Fake {
	return &Fake{
		Ident:   "fake-metadata",
		Records: map[string]metadata.Metadata{},
		At:      time.Date(2026, 8, 20, 9, 0, 0, 0, time.UTC),
		Credit:  metadata.Attribution{TextKey: "attribution.fake"},
	}
}

// Key builds the map key for a ref.
func Key(r library.Ref) string { return r.Title + "|" + r.Platform }

// ID returns the fake's identity.
func (f *Fake) ID() provider.ID { return f.Ident }

// Lookup answers from the map.
func (f *Fake) Lookup(_ context.Context, r library.Ref) (aged.Value[metadata.Metadata], error) {
	if f.Fail != nil {
		return aged.Value[metadata.Metadata]{}, f.Fail
	}
	if !r.Identifiable() {
		return aged.Value[metadata.Metadata]{}, fault.New(fault.KindMalformed, "metadatatest.Lookup")
	}
	m, ok := f.Records[Key(r)]
	if !ok {
		return aged.Value[metadata.Metadata]{}, fault.New(fault.KindNotFound, "metadatatest.Lookup")
	}
	return aged.New(m, f.At), nil
}

// Attribution returns the fake's credit line.
func (f *Fake) Attribution() metadata.Attribution { return f.Credit }

var _ metadata.Provider = (*Fake)(nil)
