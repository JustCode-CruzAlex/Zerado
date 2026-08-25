package aged_test

import (
	"testing"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/aged"
)

var now = time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)

// TestAnUnstampedValueIsNeverCurrent is the age rule's floor: a value that
// arrived with no timestamp is rendered as unknown, never as fresh. Dropping
// the age always makes a layout tidier, which is why the type refuses to let
// an unknown age look like a recent one.
func TestAnUnstampedValueIsNeverCurrent(t *testing.T) {
	var v aged.Value[int]
	if v.Known() {
		t.Fatal("the zero Value claims to know when it was observed")
	}
	if got := v.Classify(now, time.Hour); got != aged.FreshnessUnknown {
		t.Fatalf("an unstamped value classified as %v", got)
	}
	if got := v.Age(now); got != 0 {
		t.Fatalf("Age of an unstamped value = %v", got)
	}
}

// TestClassifyBands walks the three behaviours the design system has: nothing,
// the banner as chrome, and the banner as something to act on past ninety days
// (07-offline-contract.md §4.1).
func TestClassifyBands(t *testing.T) {
	cases := []struct {
		name   string
		at     time.Time
		window time.Duration
		want   aged.Freshness
	}{
		{"a price fetched a minute ago", now.Add(-time.Minute), time.Hour, aged.FreshnessCurrent},
		{"a price fetched yesterday", now.Add(-24 * time.Hour), time.Hour, aged.FreshnessStale},
		{"a library synced three days ago", now.Add(-72 * time.Hour), 6 * time.Hour, aged.FreshnessStale},
		{"a library last synced in May", now.Add(-100 * 24 * time.Hour), 6 * time.Hour, aged.FreshnessAncient},
		{"exactly ninety days", now.Add(-aged.WarnAfter), time.Hour, aged.FreshnessAncient},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := aged.New(42, c.at)
			if got := v.Classify(now, c.window); got != c.want {
				t.Fatalf("Classify = %v, want %v", got, c.want)
			}
		})
	}
}

// TestAFutureValueIsNotRenderedAsFuture: a source clock ahead of this machine
// is two machines disagreeing about NTP, not information about the data.
// "In 3 hours" next to a price would be the product reporting somebody else's
// problem as a feature.
func TestAFutureValueIsNotRenderedAsFuture(t *testing.T) {
	v := aged.New("R$ 49,99", now.Add(3*time.Hour))
	if got := v.Age(now); got != 0 {
		t.Fatalf("Age = %v, want 0 for a value stamped in the future", got)
	}
	if got := v.Classify(now, time.Minute); got != aged.FreshnessCurrent {
		t.Fatalf("Classify = %v, want FreshnessCurrent", got)
	}
}

// TestMapKeepsTheObservationTime guards the single way this mechanism is
// defeated: transforming a value and re-stamping it as new, in code that looks
// entirely reasonable.
func TestMapKeepsTheObservationTime(t *testing.T) {
	at := now.Add(-30 * 24 * time.Hour)
	in := aged.New(4999, at)
	out := aged.Map(in, func(cents int) string { return "R$ 49,99" })
	if !out.At.Equal(at) {
		t.Fatalf("Map re-stamped the value: %v, want %v", out.At, at)
	}
	if out.V != "R$ 49,99" {
		t.Fatalf("Map lost the payload: %q", out.V)
	}
}
