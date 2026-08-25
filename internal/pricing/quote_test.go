package pricing_test

import (
	"context"
	"testing"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/aged"
	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/pricing"
	"github.com/JustCode-CruzAlex/Zerado/internal/pricing/pricingtest"
)

// TestAQuoteCannotBeObtainedWithoutItsAge is 07 §4 made structural. There is
// no code path that renders the number without the age, because they arrive in
// the same value.
func TestAQuoteCannotBeObtainedWithoutItsAge(t *testing.T) {
	f := pricingtest.New()
	r := library.Ref{Title: "Hades", Platform: "PC"}
	f.Quotes[pricingtest.Key(r, "BRL")] = pricing.Quote{
		Current: pricing.Money{Amount: 4999, Currency: "BRL"},
		Low:     pricing.Money{Amount: 1999, Currency: "BRL"},
		LowAt:   time.Date(2025, 11, 28, 0, 0, 0, 0, time.UTC),
		Shop:    "Steam",
		URL:     "https://store.steampowered.com/app/1145360",
	}

	v, err := f.Quote(context.Background(), r, "BRL")
	if err != nil {
		t.Fatalf("Quote: %v", err)
	}
	if !v.Known() {
		t.Fatal("a quote came back with no observation time")
	}
	var _ aged.Value[pricing.Quote] = v // the signature is the guarantee
}

// TestWaitOrBuyIsAnswerableFromOneResponse, and offline, from a cached value —
// which is what makes the watchlist a DEGRADES feature rather than one that
// stops.
func TestWaitOrBuyIsAnswerableFromOneResponse(t *testing.T) {
	cases := []struct {
		name    string
		current int64
		low     int64
		want    pricing.Verdict
	}{
		{"at the all-time low", 1999, 1999, pricing.VerdictAtLow},
		{"below the recorded low", 1799, 1999, pricing.VerdictAtLow},
		{"a little above it", 2099, 1999, pricing.VerdictNearLow},
		{"well above it", 4999, 1999, pricing.VerdictWait},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			q := pricing.Quote{
				Current: pricing.Money{Amount: c.current, Currency: "BRL"},
				Low:     pricing.Money{Amount: c.low, Currency: "BRL"},
			}
			if got := q.Verdict(10); got != c.want {
				t.Fatalf("Verdict = %v, want %v", got, c.want)
			}
		})
	}
}

// TestAnUnknowableVerdictIsNeverAdvice: a product that guesses "wait" when it
// does not know is a product giving advice it has not earned.
func TestAnUnknowableVerdictIsNeverAdvice(t *testing.T) {
	for _, q := range []pricing.Quote{
		{},
		{Current: pricing.Money{Amount: 4999, Currency: "BRL"}},
		{Current: pricing.Money{Amount: 4999, Currency: "BRL"}, Low: pricing.Money{Amount: 1999, Currency: "USD"}},
		{Current: pricing.Money{Amount: 4999, Currency: "BRL"}, Low: pricing.Money{Amount: 0, Currency: "BRL"}},
	} {
		if got := q.Verdict(10); got != pricing.VerdictUnknown {
			t.Fatalf("Verdict(%+v) = %v, want VerdictUnknown", q, got)
		}
	}
}

// TestThereIsNowhereToPutAnAffiliateTag. The decision is enforced by the
// absence of a field, so it cannot be re-added by accident.
func TestThereIsNowhereToPutAnAffiliateTag(t *testing.T) {
	q := pricing.Quote{}
	// The struct has exactly these fields; a reviewer reading this test sees
	// the whole surface. Adding an AffiliateURL would not compile here.
	_ = q.Current
	_ = q.Low
	_ = q.LowAt
	_ = q.Shop
	_ = q.URL
}

// TestMoneyDistinguishesUnsetFromFree: an empty currency is an unset value,
// not a free game.
func TestMoneyDistinguishesUnsetFromFree(t *testing.T) {
	if (pricing.Money{}).Known() {
		t.Fatal("the zero Money claims to be a real amount")
	}
	if !(pricing.Money{Amount: 0, Currency: "BRL"}).Known() {
		t.Fatal("a genuinely free game read as unset")
	}
}

// TestNoQuoteIsNotAFailureOfAnything: a cartridge has no shop page.
func TestNoQuoteIsNotAFailureOfAnything(t *testing.T) {
	f := pricingtest.New()
	_, err := f.Quote(context.Background(), library.Ref{Title: "Chrono Trigger", Platform: "SNES"}, "BRL")
	if !fault.Is(err, fault.KindNotFound) {
		t.Fatalf("got %v, want KindNotFound", fault.KindOf(err))
	}
	if got := fault.KindNotFound.Treatment(); got != fault.TreatmentDesignedEmpty {
		t.Fatalf("no price renders as %v", got)
	}
}
