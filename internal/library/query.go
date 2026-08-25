package library

import (
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
	"github.com/JustCode-CruzAlex/Zerado/internal/status"
)

// Query selects a subset of the library.
//
// It is the shape Z-07's filter bar produces and the shape the store consumes,
// and it exists as a value type so that "the summary describes the filtered
// set" (05-state-machine.md §7 rule 2) is enforceable: the same Query goes to
// [store.Store.Games] and to [store.Store.Counts], so the list and the numbers
// above it cannot describe different sets.
//
// # Why this and not a generic CRUD surface
//
// The ticket asks for the query shapes the interface actually needs rather
// than a generic surface, and this is the whole of it. Phase 1 has exactly
// three facets — a search, a set of states, a set of sources — plus a presence
// mode. There is no builder, no expression tree and no operator: a filter
// language would be an abstraction with one consumer, and the consumer is a
// filter bar with three chips.
//
// # There is no sort field, and that is a decision
//
// 07-offline-contract.md §2 fixes Phase 1's order at title A→Z with no sort
// control, and is explicit that an on-screen sort indicator would imply a
// control that does not exist. A Sort field here would be that indicator's
// back end: present, unused, and the obvious thing to wire up. It is left out
// so that adding sorting is a deliberate act with a screen behind it.
type Query struct {
	// Search is the free-text facet. Empty matches everything.
	//
	// Matching is the store's — a case-folding, diacritic-folding substring
	// match over the title, which is what SortTitle exists for. It is not
	// specified further here because it is a SQL detail (fft-database's), and
	// because the one property a screen depends on is that an empty query
	// filters nothing (Z-07 F2).
	Search string

	// States restricts to these effective statuses. Empty means all four.
	//
	// Effective, not manual: a player filtering for IN PROGRESS means the
	// games that are in progress, including the ones that are in progress
	// because Steam said so.
	States []status.Status

	// Sources restricts to these providers. Empty means all.
	Sources []provider.ID

	// Presence decides whether absent rows are included, excluded, or the
	// only thing shown.
	Presence Presence

	// Limit and Offset page the result. Zero Limit means no limit.
	//
	// Present because Z-04 renders twelve rows of a 247-row library at 80×24
	// and reading all 247 to show 12 is a habit that stops being free at
	// 10,000. The screen's scroll position row — "ROWS 4–15 of 247" — needs
	// the total as well, which is why [store.Store.Counts] takes the same
	// Query and is not merely a convenience over Games.
	Limit  int
	Offset int
}

// Presence is the absent-row mode of a query.
//
// Absence is not a fifth state (06-data-seams.md §2.4), so it is not a member
// of States. It is an orthogonal axis, and modelling it as one is what lets
// Z-07's ABSENT chip swap the row set and then filter by state *within* it.
type Presence uint8

const (
	// PresentOnly excludes absent rows. This is the default, and it is the
	// default library view.
	PresentOnly Presence = iota

	// AbsentOnly returns only absent rows — Z-07's ABSENT facet, which is how
	// 06 §2.4's promise that absent rows "remain findable by filter" is kept.
	AbsentOnly

	// Either includes both, for the CLI's --all and for a Phase 4 merge that
	// must see everything.
	Either
)

// Empty reports whether the query filters nothing.
//
// Z-07 F2 requires that an empty editor filters nothing and shows the full
// ratio, and Z-04's unfiltered summary is the same predicate read the other
// way: a summary says "247 games" without naming a filter exactly when this is
// true. Paging is not a filter, so Limit and Offset are not consulted.
func (q Query) Empty() bool {
	return q.Search == "" &&
		len(q.States) == 0 &&
		len(q.Sources) == 0 &&
		q.Presence == PresentOnly
}

// Facets returns the names of the active facets, in the order Z-07 renders
// them.
//
// It exists for the empty-result line, which must name the facet that emptied
// the set rather than apologising — "search was fine, the state filter emptied
// it". Names are stable machine tokens; the screen renders them through the
// catalogue.
func (q Query) Facets() []string {
	var out []string
	if q.Search != "" {
		out = append(out, "search")
	}
	if len(q.States) > 0 {
		out = append(out, "state")
	}
	if len(q.Sources) > 0 {
		out = append(out, "source")
	}
	if q.Presence == AbsentOnly {
		out = append(out, "absent")
	}
	return out
}
