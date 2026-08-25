package metadata_test

import (
	"context"
	"testing"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/metadata"
	"github.com/JustCode-CruzAlex/Zerado/internal/metadata/metadatatest"
)

// TestHavingNoSourceIsADesignedState, not an error. This is the difference
// between a product that works without a metadata source and a product that is
// broken without one.
func TestHavingNoSourceIsADesignedState(t *testing.T) {
	_, err := metadata.Null{}.Lookup(context.Background(), library.Ref{Title: "Hades", Platform: "PC"})
	if !fault.Is(err, fault.KindNotFound) {
		t.Fatalf("got %v, want KindNotFound", fault.KindOf(err))
	}
	if got := fault.KindOf(err).Treatment(); got != fault.TreatmentDesignedEmpty {
		t.Fatalf("no metadata renders as %v; an error banner would make the product broken without a source", got)
	}
	if a := (metadata.Null{}).Attribution(); a.TextKey != "" || a.Verbatim != "" {
		t.Fatal("the null source demands a credit")
	}
}

// TestTheSeamNamesNoProvider is the test of whether a seam is genuinely
// agnostic: swapping the source must change one implementation and no
// signature above it. A fake keyed on a title and a platform satisfies exactly
// the same interface a store-aware source would.
func TestTheSeamNamesNoProvider(t *testing.T) {
	var sources []metadata.Provider

	null := metadata.Null{}
	fake := metadatatest.New()
	r := library.Ref{Title: "Chrono Trigger", Platform: "SNES"}
	fake.Records[metadatatest.Key(r)] = metadata.Metadata{Sinopse: "A time-travelling RPG."}

	sources = append(sources, null, fake)

	for _, s := range sources {
		v, err := s.Lookup(context.Background(), r)
		switch s.(type) {
		case metadata.Null:
			if err == nil {
				t.Fatal("the null source answered")
			}
		default:
			if err != nil {
				t.Fatalf("%s: %v", s.ID(), err)
			}
			if v.V.Empty() {
				t.Fatalf("%s returned an empty record", s.ID())
			}
			if !v.Known() {
				t.Fatalf("%s returned a record with no observation time", s.ID())
			}
		}
	}
}

// TestABlankRecordIsAlsoTheDesignedEmpty: a source may legitimately return a
// record with every field blank, and a caller that only checked the error
// would render an empty pane with three blank labels.
func TestABlankRecordIsAlsoTheDesignedEmpty(t *testing.T) {
	if !(metadata.Metadata{}).Empty() {
		t.Fatal("a blank record does not report itself empty")
	}
	if (metadata.Metadata{Genres: []string{"RPG"}}).Empty() {
		t.Fatal("a record with a genre reported itself empty")
	}
}

// TestASourceThatIsDownIsADifferentScreen from a source that has never heard
// of a game. That distinction is the entire reason the taxonomy exists.
func TestASourceThatIsDownIsADifferentScreen(t *testing.T) {
	fake := metadatatest.New()
	fake.Fail = fault.New(fault.KindUnreachable, "metadatatest.Lookup")

	_, err := fake.Lookup(context.Background(), library.Ref{Title: "Hades", Platform: "PC"})
	if fault.KindOf(err).Treatment() == fault.TreatmentDesignedEmpty {
		t.Fatal("a source being down rendered as a designed empty; the player would never learn there was anything to retry")
	}
	if !fault.KindOf(err).Retryable() {
		t.Fatal("a source being down is not retryable; the banner would offer no way out")
	}
}
