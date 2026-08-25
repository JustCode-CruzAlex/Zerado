package library_test

import (
	"testing"

	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
	"github.com/JustCode-CruzAlex/Zerado/internal/status"
)

// TestAnEmptyQueryFiltersNothing is Z-07 F2: the editor opens showing "247 of
// 247" and the list is unchanged.
func TestAnEmptyQueryFiltersNothing(t *testing.T) {
	if !(library.Query{}).Empty() {
		t.Fatal("the zero Query reports itself as filtering something")
	}
	if !(library.Query{Limit: 12, Offset: 24}).Empty() {
		t.Fatal("paging counted as a filter; the summary would then claim a filter that does not exist")
	}
	for _, q := range []library.Query{
		{Search: "souls"},
		{States: []status.Status{status.Zerado}},
		{Sources: []provider.ID{"physical"}},
		{Presence: library.AbsentOnly},
	} {
		if q.Empty() {
			t.Fatalf("%+v reported itself empty", q)
		}
	}
}

// TestFacetsNameWhatEmptiedTheSet: an empty result names the filter that
// emptied it rather than apologising — "search was fine, the state filter
// emptied it."
func TestFacetsNameWhatEmptiedTheSet(t *testing.T) {
	q := library.Query{
		Search:   "souls",
		States:   []status.Status{status.Zerado},
		Sources:  []provider.ID{"steam"},
		Presence: library.AbsentOnly,
	}
	got := q.Facets()
	want := []string{"search", "state", "source", "absent"}
	if len(got) != len(want) {
		t.Fatalf("Facets = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Facets = %v, want %v (order is the order Z-07 renders them)", got, want)
		}
	}
	if len((library.Query{}).Facets()) != 0 {
		t.Fatal("an unfiltered query named a facet")
	}
}

// TestAbsenceIsNotAState: the four states are ratified and CVD-verified, and
// an absent game still has one — usually the most valuable one.
func TestAbsenceIsNotAState(t *testing.T) {
	for _, s := range status.All() {
		if int(s) > 4 {
			t.Fatalf("a fifth state exists: %v", s)
		}
	}
	// Presence is a separate axis, which is what lets the ABSENT facet swap
	// the row set and then filter by state within it.
	q := library.Query{Presence: library.AbsentOnly, States: []status.Status{status.Zerado}}
	if len(q.Facets()) != 2 {
		t.Fatalf("state and presence did not compose: %v", q.Facets())
	}
}

// TestGameDistinguishesThreeFactsAboutPlaytime is what Z-05 renders as "not
// tracked", "—" and "0h". Collapsing any two produces a screen that tells a
// player their cartridge has been played for zero hours.
func TestGameDistinguishesThreeFactsAboutPlaytime(t *testing.T) {
	zero := 0

	notTracked := library.Game{Provider: "physical"}
	if _, reported := notTracked.Playtime(); reported {
		t.Fatal("a cartridge reported a playtime")
	}
	if got := notTracked.Status(false); got != status.NotStarted {
		t.Fatalf("a cartridge derived to %v", got)
	}

	notFetched := library.Game{Provider: "steam"}
	if _, reported := notFetched.Playtime(); reported {
		t.Fatal("an unfetched row reported a playtime")
	}

	knownEmpty := library.Game{Provider: "steam", PlaytimeMinutes: &zero}
	pt, reported := knownEmpty.Playtime()
	if !reported || pt != 0 {
		t.Fatalf("a known-zero playtime read as %d, %v; zero is a real value", pt, reported)
	}
	if knownEmpty.Status(true) != status.NotStarted {
		t.Fatal("a known-zero playtime did not derive to NOT STARTED")
	}
}

// TestRefIsAnswerableForACartridge is the shape test for the enrichment
// seams: a ref whose primary key were a store identifier would be unanswerable
// for a shelf, which is the majority case for a shelf.
func TestRefIsAnswerableForACartridge(t *testing.T) {
	cartridge := library.Ref{Title: "Chrono Trigger", Platform: "SNES"}
	if !cartridge.Identifiable() {
		t.Fatal("a title and a platform are not enough to ask a source about a game")
	}
	if (library.Ref{Provider: "steam", ProviderRef: "1086940"}).Identifiable() {
		t.Fatal("a ref with only store identifiers passed; a source that does not know Steam could not use it")
	}

	g := library.Game{Title: "Hades", Platform: "PC", Provider: "steam", ProviderRef: "1145360"}
	r := library.RefOf(g)
	if r.Title != g.Title || r.Platform != g.Platform || r.ProviderRef != g.ProviderRef {
		t.Fatalf("RefOf lost a field: %+v", r)
	}
}

// TestOverriddenDrivesTheFifthMenuItem: Z-06 shows "clear override" only when
// there is one to clear.
func TestOverriddenDrivesTheFifthMenuItem(t *testing.T) {
	z := status.Zerado
	if !(library.Game{StatusManual: &z}).Overridden() {
		t.Fatal("a game with a manual status does not report itself overridden")
	}
	if (library.Game{}).Overridden() {
		t.Fatal("a game with no manual status reports itself overridden; Z-06 would offer to clear nothing")
	}
}
