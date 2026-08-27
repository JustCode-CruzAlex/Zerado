package status_test

import (
	"testing"

	"github.com/JustCode-CruzAlex/Zerado/internal/status"
)

// TestZeroValueIsNeverAState guards 05-state-machine.md §1: the zero value is
// deliberately invalid so an uninitialised status can never be mistaken for
// "not started".
func TestZeroValueIsNeverAState(t *testing.T) {
	var s status.Status
	if s.Valid() {
		t.Fatalf("the zero Status reports itself valid; an uninitialised value would render as a state")
	}
	if s != status.Unknown {
		t.Fatalf("the zero Status is %v, expected Unknown", s)
	}
	for _, v := range status.All() {
		if v == status.Unknown {
			t.Fatalf("Unknown appears in All(); it is never rendered and never counted")
		}
	}
}

// TestDeriveIsCapabilityDriven is the test that physical copies are
// first-class. A provider that cannot report playtime gets NOT STARTED
// regardless of what number it was handed, so there is no path by which a
// typed or defaulted zero becomes evidence.
func TestDeriveIsCapabilityDriven(t *testing.T) {
	cases := []struct {
		name     string
		minutes  int
		playtime bool
		want     status.Status
	}{
		{"a store that reports playtime, and there is some", 87 * 60, true, status.InProgress},
		{"a store that reports playtime, and there is none", 0, true, status.NotStarted},
		{"a cartridge, with a number that must be ignored", 9999, false, status.NotStarted},
		{"a cartridge, with no number", 0, false, status.NotStarted},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := status.Derive(c.minutes, c.playtime); got != c.want {
				t.Fatalf("Derive(%d, %v) = %v, want %v", c.minutes, c.playtime, got, c.want)
			}
		})
	}
}

// TestNothingDerivesToZerado is 05-state-machine.md §4 as an executable rule:
// ZERADO is never automatic, and no combination of provider-reported inputs
// may produce it.
func TestNothingDerivesToZerado(t *testing.T) {
	for _, playtime := range []bool{true, false} {
		for _, m := range []int{0, 1, 60, 100000} {
			if got := status.Derive(m, playtime); got == status.Zerado || got == status.Abandoned {
				t.Fatalf("Derive(%d, %v) = %v; only the player may set that state", m, playtime, got)
			}
		}
	}
}

// TestManualWinsPermanently is the invariant that makes the product
// trustworthy: a manual value survives any playtime a sync reports.
func TestManualWinsPermanently(t *testing.T) {
	z := status.Zerado
	if got := status.Effective(&z, 300, true); got != status.Zerado {
		t.Fatalf("playtime overrode a manual ZERADO: got %v", got)
	}
	n := status.NotStarted
	if got := status.Effective(&n, 300, true); got != status.NotStarted {
		t.Fatalf("an explicit NOT STARTED was overridden by playtime: got %v", got)
	}
}

// TestClearingAnOverrideRederives is 05 §5's consequence, stated in the spec
// as something the copy must be honest about: clearing on a game with playtime
// makes it IN PROGRESS immediately, and choosing NOT STARTED explicitly is a
// different action.
func TestClearingAnOverrideRederives(t *testing.T) {
	if got := status.Effective(nil, 300, true); got != status.InProgress {
		t.Fatalf("cleared override with playtime = %v, want InProgress", got)
	}
	n := status.NotStarted
	if status.Effective(&n, 300, true) == status.Effective(nil, 300, true) {
		t.Fatal("setting NOT STARTED and clearing the override produced the same result; the model needs them distinct")
	}
}

// TestEffectiveIsNeverUnknown guards §7 rule 1: every row must be countable,
// so even a corrupt manual value falls back to the derivation.
func TestEffectiveIsNeverUnknown(t *testing.T) {
	bogus := status.Status(200)
	if got := status.Effective(&bogus, 0, false); got != status.NotStarted {
		t.Fatalf("a corrupt manual value produced %v; it must fall back to the derivation", got)
	}
	if !status.Effective(&bogus, 0, false).Valid() {
		t.Fatal("Effective returned an invalid status")
	}
}

// TestParseRoundTrips keeps the stored form stable. These strings are written
// into every player's library file.
func TestParseRoundTrips(t *testing.T) {
	for _, s := range status.All() {
		got, ok := status.Parse(s.String())
		if !ok || got != s {
			t.Fatalf("Parse(%q) = %v, %v; want %v, true", s.String(), got, ok, s)
		}
	}
	if _, ok := status.Parse(""); ok {
		t.Fatal("the empty string parsed as a status; NULL must be represented by a nil pointer, not a value")
	}
	if _, ok := status.Parse("finished"); ok {
		t.Fatal("an unknown token parsed as a status")
	}
}

// TestCountsAlwaysSum is 05 §7 rule 1: if the summary says 247 games, the four
// state counts add to 247.
func TestCountsAlwaysSum(t *testing.T) {
	var c status.Counts
	for i := 0; i < 10; i++ {
		c.Add(status.NotStarted)
	}
	for i := 0; i < 3; i++ {
		c.Add(status.Zerado)
	}
	c.Add(status.Abandoned)
	if !c.Sums() {
		t.Fatalf("counts do not sum: %+v", c)
	}
	if c.Total != 14 {
		t.Fatalf("Total = %d, want 14", c.Total)
	}

	// An unclassifiable row must not increment the total either, or the
	// summary would claim a game the four counts cannot account for.
	c.Add(status.Unknown)
	if !c.Sums() || c.Total != 14 {
		t.Fatalf("an Unknown row disturbed the summary: %+v", c)
	}
}
