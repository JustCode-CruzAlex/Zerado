// Package pricingtest holds a price provider backed by a map.
package pricingtest

import (
	"context"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/aged"
	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/metadata"
	"github.com/JustCode-CruzAlex/Zerado/internal/pricing"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
)

// Fake answers quotes from a map, keyed by title, platform and currency.
type Fake struct {
	Ident  provider.ID
	Quotes map[string]pricing.Quote
	At     time.Time
	Credit metadata.Attribution
	Fail   error
}

// New returns a Fake with a fixed observation time.
func New() *Fake {
	return &Fake{
		Ident:  "fake-pricing",
		Quotes: map[string]pricing.Quote{},
		At:     time.Date(2026, 8, 24, 8, 0, 0, 0, time.UTC),
		Credit: metadata.Attribution{TextKey: "attribution.fake_prices"},
	}
}

// Key builds the map key for a ref and currency.
func Key(r library.Ref, cur pricing.Currency) string {
	return r.Title + "|" + r.Platform + "|" + string(cur)
}

// ID returns the fake's identity.
func (f *Fake) ID() provider.ID { return f.Ident }

// Quote answers from the map, stamped.
func (f *Fake) Quote(_ context.Context, r library.Ref, cur pricing.Currency) (aged.Value[pricing.Quote], error) {
	if f.Fail != nil {
		return aged.Value[pricing.Quote]{}, f.Fail
	}
	q, ok := f.Quotes[Key(r, cur)]
	if !ok {
		return aged.Value[pricing.Quote]{}, fault.New(fault.KindNotFound, "pricingtest.Quote")
	}
	return aged.New(q, f.At), nil
}

// Attribution returns the fake's credit line.
func (f *Fake) Attribution() metadata.Attribution { return f.Credit }

var _ pricing.Provider = (*Fake)(nil)
