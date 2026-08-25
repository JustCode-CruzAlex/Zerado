// Package pricing is the price seam: what a game costs now, what it has ever
// cost, and therefore whether to wait.
//
// # The disclosure is structural, and there is no affiliate URL
//
// Founder direction, 2026-08-25: affiliate links are dropped so that Zerado is
// cleanly non-commercial — free software, donation-supported, zero revenue.
// [Quote.URL] is a plain shop link and there is no field that could carry a
// commission tag. That is a decision rather than an omission, and it is
// enforced by the absence: a struct with no place to put an affiliate
// parameter cannot grow one by accident.
//
// The price feature survives intact. Current price, the all-time low, when the
// low was, and the wait-or-buy verdict are all exactly as designed. What went
// is the tag on the outbound link and the disclosure obligation that used to
// travel with it.
//
// # Every quote carries its age, and the age is not optional
//
// A [Quote] is always returned inside an [aged.Value]. 07 §4 is blunt about
// why: Zerado's value proposition is telling a player *not* to buy something,
// and a stale price presented as current is not a cosmetic defect but the
// product giving financial advice from memory. There is no code path that can
// render the number without the age, because they arrive in the same value.
//
// # The verdict is answerable without a second call
//
// The ticket requires that "wait or buy?" be answerable from one response, so
// [Quote] carries the current price, the all-time low and when that low was.
// [Quote.Verdict] computes the answer locally, offline, from a cached quote —
// which is what makes the watchlist a DEGRADES feature rather than one that
// stops.
package pricing

import (
	"context"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/aged"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/metadata"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
)

// Currency is an ISO 4217 code — "BRL", "USD", "EUR".
//
// It is a parameter on every lookup rather than a global setting because the
// budget feature is money in the player's own currency, and a source that can
// answer in one currency and not another must be able to say so per request.
type Currency string

// Money is an amount in a currency's minor units.
//
// Integer minor units, never a float: a price is not a measurement and 0.1
// does not exist in binary floating point. The formatting — symbol placement,
// decimal separator, grouping — is x/text's, at render time, from the player's
// locale, which is also why this type has no String method that could tempt a
// screen into formatting money itself.
type Money struct {
	// Amount is in minor units: 4999 is BRL 49.99.
	Amount int64

	// Currency is the code the amount is denominated in. A Money with an
	// empty Currency is not zero money; it is an unset value, and
	// [Money.Known] is how a caller tells the difference from a genuine free
	// game.
	Currency Currency
}

// Known reports whether this Money carries a currency and is therefore a real
// amount.
func (m Money) Known() bool { return m.Currency != "" }

// Provider is a source of prices.
type Provider interface {
	// ID is the source's stable identity.
	ID() provider.ID

	// Quote returns the current price and the all-time low for one game, in
	// one currency, stamped with when it was observed.
	//
	// A game the source does not sell or does not know is fault.KindNotFound,
	// which renders as no price rather than as an error: a cartridge has no
	// shop page and that is not a failure of anything.
	Quote(ctx context.Context, r library.Ref, cur Currency) (aged.Value[Quote], error)

	// Attribution is the credit this source requires, on the same terms as
	// the metadata seam's: a property of the source, so swapping the source
	// swaps the credit.
	//
	// The published copy names IsThereAnyDeal as the source of price data, and
	// it is nevertheless behind this interface for the same reason a metadata
	// source is: the page's claim is about what the data is, not about who is
	// guaranteed to supply it forever.
	Attribution() metadata.Attribution
}

// Quote is one source's answer about one game's price.
type Quote struct {
	// Current is what it costs now.
	Current Money

	// Low is the all-time low, and LowAt is when it was.
	//
	// LowAt is not optional. A low with no date is not information: "it was
	// R$ 19 once" is a different statement depending on whether once was last
	// month or in 2019, and only one of those two should make a player wait.
	Low   Money
	LowAt time.Time

	// Shop names the retailer, for display.
	Shop string

	// URL is a plain shop link. No affiliate tag, ever — there is no field
	// that could carry one and no parameter that may be appended to this.
	URL string
}

// DiscountPercent returns how far ABOVE the all-time low the current price
// sits, as a percentage of the low, rounded toward zero.
//
// Above, not below: it computes (Current - Low) / Low, so a price at the low
// returns 0, one below it returns a negative, and a full-price game returns a
// large positive. That direction is what [Quote.Verdict] reads — pct <= 0 is
// VerdictAtLow — and the name reads as its opposite, which is why the
// direction is spelled out here rather than left to be inferred from the
// arithmetic.
//
// It returns 0 and false when either amount is unknown or the currencies
// differ, rather than computing across currencies — a percentage derived from
// two different denominations is a number that looks right and means nothing.
func (q Quote) DiscountPercent() (int, bool) {
	if !q.Current.Known() || !q.Low.Known() || q.Current.Currency != q.Low.Currency {
		return 0, false
	}
	if q.Low.Amount <= 0 {
		return 0, false
	}
	diff := q.Current.Amount - q.Low.Amount
	return int(diff * 100 / q.Low.Amount), true
}

// Verdict is the wait-or-buy answer.
type Verdict uint8

const (
	// VerdictUnknown means there is not enough to say. It renders as no
	// verdict, never as "wait" — a product that guesses "wait" when it does
	// not know is a product giving advice it has not earned.
	VerdictUnknown Verdict = iota

	// VerdictAtLow means the current price is at or below the all-time low.
	VerdictAtLow

	// VerdictNearLow means the current price is within the tolerance of the
	// all-time low.
	VerdictNearLow

	// VerdictWait means the current price is meaningfully above the all-time
	// low.
	VerdictWait
)

// String returns the stable machine name, used by the CLI's JSON output.
func (v Verdict) String() string {
	switch v {
	case VerdictAtLow:
		return "at_low"
	case VerdictNearLow:
		return "near_low"
	case VerdictWait:
		return "wait"
	default:
		return "unknown"
	}
}

// Verdict answers wait-or-buy from this quote alone.
//
// tolerancePercent is how far above the all-time low still counts as near it;
// the caller supplies it because it is a product judgement about what "close
// enough" means, and Phase 3 will want to tune it without a signature change.
//
// It is a method on the quote rather than a provider call because the whole
// point is that it is answerable offline, from a cached value, with the
// network down. A verdict computed by the source would be a verdict Zerado
// could not produce for the watchlist when it matters.
func (q Quote) Verdict(tolerancePercent int) Verdict {
	pct, ok := q.DiscountPercent()
	if !ok {
		return VerdictUnknown
	}
	switch {
	case pct <= 0:
		return VerdictAtLow
	case pct <= tolerancePercent:
		return VerdictNearLow
	default:
		return VerdictWait
	}
}
