package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider/providertest"
	"github.com/JustCode-CruzAlex/Zerado/internal/status"
	"github.com/JustCode-CruzAlex/Zerado/internal/store"
	"github.com/JustCode-CruzAlex/Zerado/internal/store/storetest"
)

var at = time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)

func seeded(t *testing.T) (*storetest.Fake, context.Context) {
	t.Helper()
	f := storetest.New()
	f.Playtime["steam"] = true
	f.Playtime["physical"] = false
	ctx := context.Background()

	m87, m0 := 87*60, 0
	_, err := f.UpsertBatch(ctx, "steam", []provider.Item{
		{ProviderRef: "1086940", Title: "Baldur's Gate 3", Platform: "PC", PlaytimeMinutes: &m87},
		{ProviderRef: "774361", Title: "Blasphemous", Platform: "PC", PlaytimeMinutes: &m0},
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	return f, ctx
}

// TestASyncNeverChangesAStatusThePlayerSet is the invariant that makes the
// product trustworthy. Mark a game ZERADO, play three more hours, sync: the
// playtime moves and the state does not.
func TestASyncNeverChangesAStatusThePlayerSet(t *testing.T) {
	f, ctx := seeded(t)
	games, _ := f.Games(ctx, library.Query{})
	id := games[0].ID

	z := status.Zerado
	if err := f.SetStatus(ctx, id, &z); err != nil {
		t.Fatalf("SetStatus: %v", err)
	}

	more := 300 * 60
	if _, err := f.UpsertBatch(ctx, "steam", []provider.Item{
		{ProviderRef: games[0].ProviderRef, Title: games[0].Title, Platform: games[0].Platform, PlaytimeMinutes: &more},
	}); err != nil {
		t.Fatalf("UpsertBatch: %v", err)
	}

	g, _ := f.Game(ctx, id)
	if g.Status(true) != status.Zerado {
		t.Fatalf("a sync changed the player's status to %v; the one action the product is named after was undone by a background job", g.Status(true))
	}
	if pt, _ := g.Playtime(); pt != more {
		t.Fatalf("playtime = %d, want %d — the sync must still update provider facts", pt, more)
	}
}

// TestClearingAnOverrideIsADistinctAction and is stamped, because Phase 4's
// last-write-wins needs a timestamp for it to compare.
func TestClearingAnOverrideIsADistinctAction(t *testing.T) {
	f, ctx := seeded(t)
	games, _ := f.Games(ctx, library.Query{})
	id := games[0].ID

	z := status.Zerado
	_ = f.SetStatus(ctx, id, &z)
	if err := f.SetStatus(ctx, id, nil); err != nil {
		t.Fatalf("clear: %v", err)
	}
	g, _ := f.Game(ctx, id)
	if g.Overridden() {
		t.Fatal("the override survived a clear")
	}
	if g.StatusChangedAt == nil {
		t.Fatal("clearing an override left no timestamp; a Phase 4 merge could not resolve it")
	}
	if g.Status(true) != status.InProgress {
		t.Fatalf("after clearing, the derivation gave %v; 05 §5 requires it to re-derive to IN PROGRESS", g.Status(true))
	}
}

// TestOnlyAnOkRunMayTombstone is the guard that protects a player's own work
// from a truncated stream. In a partial sync, "not returned" and "not reached"
// are indistinguishable.
func TestOnlyAnOkRunMayTombstone(t *testing.T) {
	for _, s := range []store.RunStatus{store.RunPartial, store.RunFailed, store.RunCancelled} {
		t.Run(s.String(), func(t *testing.T) {
			f, ctx := seeded(t)
			run, _ := f.StartRun(ctx, "steam", at)
			_ = f.FinishRun(ctx, run, store.RunResult{FinishedAt: at, Status: s})

			_, err := f.ReconcileAbsence(ctx, run, []string{"1086940"})
			if err == nil {
				t.Fatalf("a %s run was allowed to tombstone; a player's finished game would vanish because a stream broke", s)
			}
			if !fault.Is(err, fault.KindPrecondition) {
				t.Fatalf("got %v, want KindPrecondition", fault.KindOf(err))
			}
			games, _ := f.Games(ctx, library.Query{})
			for _, g := range games {
				if g.Absent() {
					t.Fatal("a row was marked absent despite the refusal")
				}
			}
		})
	}
}

// TestAbsenceIsReversibleAndNeverADeletion walks 06 §2.4 end to end: a
// complete sync that stops returning a game marks it absent, the row stays and
// is findable by filter, and the moment it comes back the flag clears
// silently.
func TestAbsenceIsReversibleAndNeverADeletion(t *testing.T) {
	f, ctx := seeded(t)

	run, _ := f.StartRun(ctx, "steam", at)
	_ = f.FinishRun(ctx, run, store.RunResult{FinishedAt: at, Status: store.RunOK})
	a, err := f.ReconcileAbsence(ctx, run, []string{"1086940"})
	if err != nil {
		t.Fatalf("ReconcileAbsence: %v", err)
	}
	if a.MarkedAbsent != 1 {
		t.Fatalf("MarkedAbsent = %d, want 1", a.MarkedAbsent)
	}

	shown, _ := f.Games(ctx, library.Query{})
	if len(shown) != 1 {
		t.Fatalf("the default view shows %d rows, want 1 — absent rows are excluded from it", len(shown))
	}
	absent, _ := f.Games(ctx, library.Query{Presence: library.AbsentOnly})
	if len(absent) != 1 {
		t.Fatalf("the absent facet found %d rows, want 1 — the row must remain findable by filter", len(absent))
	}
	all, _ := f.Games(ctx, library.Query{Presence: library.Either})
	if len(all) != 2 {
		t.Fatalf("nothing was deleted, so Either must return 2, got %d", len(all))
	}

	m := 0
	if _, err := f.UpsertBatch(ctx, "steam", []provider.Item{
		{ProviderRef: "774361", Title: "Blasphemous", Platform: "PC", PlaytimeMinutes: &m},
	}); err != nil {
		t.Fatalf("UpsertBatch: %v", err)
	}
	back, _ := f.Games(ctx, library.Query{})
	if len(back) != 2 {
		t.Fatalf("the returning game was not restored to the default view: %d rows", len(back))
	}
}

// TestCountsDescribeTheFilteredSet is 05 §7 rule 2: showing whole-library
// counts above a filtered list is the most common way a list view lies.
func TestCountsDescribeTheFilteredSet(t *testing.T) {
	f, ctx := seeded(t)

	manual := providertest.NewManual()
	it, _ := manual.Compose(provider.Entry{"title": "Chrono Trigger", "platform": "SNES"})
	if _, err := f.AddManual(ctx, manual.ID(), it, nil); err != nil {
		t.Fatalf("AddManual: %v", err)
	}

	whole, _ := f.Counts(ctx, library.Query{})
	if whole.Total != 3 || !whole.Sums() {
		t.Fatalf("whole-library counts = %+v", whole)
	}

	q := library.Query{Sources: []provider.ID{"physical"}}
	filtered, _ := f.Counts(ctx, q)
	if filtered.Total != 1 {
		t.Fatalf("filtered Total = %d, want 1", filtered.Total)
	}
	if !filtered.Sums() {
		t.Fatalf("filtered counts do not sum: %+v", filtered)
	}

	rows, _ := f.Games(ctx, q)
	if len(rows) != filtered.Total {
		t.Fatalf("the list shows %d rows and the summary claims %d; they were built from the same Query and must agree", len(rows), filtered.Total)
	}
}

// TestPagingDoesNotShrinkTheSummary: the scroll-position row says "4–15 of
// 247", so the summary describes the filtered set and not the visible page.
func TestPagingDoesNotShrinkTheSummary(t *testing.T) {
	f, ctx := seeded(t)
	q := library.Query{Limit: 1}
	rows, _ := f.Games(ctx, q)
	counts, _ := f.Counts(ctx, q)
	if len(rows) != 1 {
		t.Fatalf("Limit was ignored: %d rows", len(rows))
	}
	if counts.Total != 2 {
		t.Fatalf("Counts honoured Limit: Total = %d, want 2", counts.Total)
	}
}

// TestAHandEnteredGameIsNotSecondClass: it is the same type on the same path,
// it derives correctly with no special case, and its states are all the
// player's.
func TestAHandEnteredGameIsNotSecondClass(t *testing.T) {
	f, ctx := seeded(t)
	manual := providertest.NewManual()
	it, _ := manual.Compose(provider.Entry{"title": "Chrono Trigger", "platform": "SNES"})
	id, err := f.AddManual(ctx, manual.ID(), it, nil)
	if err != nil {
		t.Fatalf("AddManual: %v", err)
	}
	g, _ := f.Game(ctx, id)
	if g.Status(f.Playtime[g.Provider]) != status.NotStarted {
		t.Fatalf("a cartridge derived to %v", g.Status(false))
	}
	if _, reported := g.Playtime(); reported {
		t.Fatal("a cartridge reported a playtime")
	}

	z := status.Zerado
	if err := f.SetStatus(ctx, id, &z); err != nil {
		t.Fatalf("a cartridge could not be marked zerado: %v", err)
	}
	rows, _ := f.Games(ctx, library.Query{States: []status.Status{status.Zerado}})
	if len(rows) != 1 || rows[0].ID != id {
		t.Fatalf("the state filter did not find the hand-entered game: %v", rows)
	}
}

// TestUpsertCountsAreHonest, because Z-03's DONE line states all three and
// they must add up to what was seen.
func TestUpsertCountsAreHonest(t *testing.T) {
	f, ctx := seeded(t)
	m := 0
	res, _ := f.UpsertBatch(ctx, "steam", []provider.Item{
		{ProviderRef: "774361", Title: "Blasphemous", Platform: "PC", PlaytimeMinutes: &m},      // unchanged
		{ProviderRef: "1086940", Title: "Baldur's Gate 3", Platform: "PC", PlaytimeMinutes: &m}, // changed
		{ProviderRef: "413150", Title: "Stardew Valley", Platform: "PC", PlaytimeMinutes: &m},   // new
	})
	if res.New != 1 || res.Changed != 1 || res.Unchanged != 1 {
		t.Fatalf("BatchResult = %+v, want one of each", res)
	}
	if res.Total() != 3 {
		t.Fatalf("Total = %d, want 3", res.Total())
	}
}

// TestAnEmptyBatchIsANoOp: the caller that classified a zero-item sync as a
// refusal must never reach here, and the belt-and-braces case is harmless.
func TestAnEmptyBatchIsANoOp(t *testing.T) {
	f, ctx := seeded(t)
	res, err := f.UpsertBatch(ctx, "steam", nil)
	if err != nil {
		t.Fatalf("UpsertBatch(nil) = %v", err)
	}
	if res.Total() != 0 {
		t.Fatalf("an empty batch reported %+v", res)
	}
	rows, _ := f.Games(ctx, library.Query{Presence: library.Either})
	if len(rows) != 2 {
		t.Fatalf("an empty batch changed the library: %d rows", len(rows))
	}
}

// TestNothingHasEverBeenSynced: LastRun returns nil rather than a zero value,
// so Z-03 renders "Nothing has ever been synced." instead of a fabricated age.
func TestNothingHasEverBeenSynced(t *testing.T) {
	f, ctx := seeded(t)
	run, err := f.LastRun(ctx, "steam")
	if err != nil {
		t.Fatalf("LastRun: %v", err)
	}
	if run != nil {
		t.Fatal("LastRun invented a run")
	}
}

// TestStatusForFaultJoinsTheTaxonomyToTheHistory: the same transport failure
// is PARTIAL if items had arrived and FAILED if none had, and Z-03's copy for
// those two is completely different.
func TestStatusForFaultJoinsTheTaxonomyToTheHistory(t *testing.T) {
	unreachable := fault.New(fault.KindUnreachable, "steam.Sync")
	if got := store.StatusForFault(unreachable, true); got != store.RunPartial {
		t.Fatalf("with items delivered: %v, want partial", got)
	}
	if got := store.StatusForFault(unreachable, false); got != store.RunFailed {
		t.Fatalf("with nothing delivered: %v, want failed", got)
	}
	if got := store.StatusForFault(fault.New(fault.KindCancelled, "steam.Sync"), true); got != store.RunCancelled {
		t.Fatalf("a cancel: %v, want cancelled", got)
	}
	if got := store.StatusForFault(nil, true); got != store.RunOK {
		t.Fatalf("success: %v, want ok", got)
	}
	if !store.RunOK.MayTombstone() {
		t.Fatal("an ok run cannot tombstone")
	}
	for _, s := range []store.RunStatus{store.RunRunning, store.RunPartial, store.RunFailed, store.RunCancelled} {
		if s.MayTombstone() {
			t.Fatalf("%v may tombstone", s)
		}
	}
}

// TestDisconnectDoesNotTouchTheLibrary: disconnecting Steam is not a statement
// about what the player owns, and the rows carry the player's own work.
func TestDisconnectDoesNotTouchTheLibrary(t *testing.T) {
	f, ctx := seeded(t)
	_ = f.SaveConnection(ctx, "steam", "76561198012345678", at)
	if err := f.DeleteConnection(ctx, "steam"); err != nil {
		t.Fatalf("DeleteConnection: %v", err)
	}
	rows, _ := f.Games(ctx, library.Query{Presence: library.Either})
	if len(rows) != 2 {
		t.Fatalf("disconnecting removed rows: %d left", len(rows))
	}
}
