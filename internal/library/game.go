// Package library holds the domain types the screens read: the game as Zerado
// stores it, the query that selects games, and the identity that will let two
// devices merge in Phase 4.
//
// It sits above the provider seam and below the store seam. A provider
// produces provider.Item; the store turns Items into Games; a screen reads
// Games and nothing else. Nothing in this package does I/O, and that is what
// makes every Phase 1 screen testable with no network and no database.
//
// # Zerado is a games product, and this package is a games model
//
// ADR-0001 D5 pruned the media-polymorphic core and it is worth restating what
// that means at the interface layer, because the ticket asks for
// "media-type polymorphism at the interface layer" and the ratified answer is
// that there must not be any. 06-data-seams.md §7 names a media-type
// abstraction as explicitly not a seam: "an interface parameterised on a type
// that has one value is machinery without a purpose."
//
// The door is held open by two things and neither of them is a type
// parameter: the stored entity is called item rather than games, and it
// carries an item_type constrained to 'game'. That is the whole affordance,
// it costs two columns, and it is why the Go type here is [Game] with no
// generic parameter and no MediaType field. See
// docs/api/06-divergences-from-the-spine.md for the reconciliation with the
// ticket's wording.
package library

import (
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
	"github.com/JustCode-CruzAlex/Zerado/internal/status"
)

// GameID is the local surrogate key.
//
// It is local: it identifies a row in this machine's file and it means nothing
// on another device. Phase 4 merges on [UID] instead, which is exactly why the
// two are different types rather than two int64s that would eventually be
// passed to each other's functions.
type GameID int64

// UID is the stable, content-derived cross-device identity.
//
// It is a merge *hint*, never an authority. 06-data-seams.md §6.2 derives it
// as uuidv5 over a normalised title and platform, and is explicit that the
// normalisation is imperfect — two editions may collide, the same game may
// fail to match across platforms — so Phase 4 presents ambiguous matches to
// the player rather than guessing, and (provider_id, provider_ref) remains the
// unique constraint while this is merely indexed.
//
// Assigned at insert and never changed. Adding it in Phase 1 costs one column
// and one index; adding it in Phase 4 would cost a migration that has to
// invent stable identities for rows whose titles the player has since edited.
//
// The exact normalisation rules are fft-database's, against a real library's
// titles (13-handoffs.md §4). This type carries the policy, not the algorithm.
type UID string

// Game is one row of the library, as the screens see it.
//
// Its optional fields are pointers for the reason given on provider.Item:
// "not tracked", "not fetched yet" and "known to be nothing" are three facts,
// and Z-05 renders all three differently.
type Game struct {
	// ID is the local key. UID is the cross-device hint.
	ID  GameID
	UID UID

	// Provider is where this row came from. It is stored, and it is read for
	// exactly two purposes: looking the provider up in a registry, and
	// rendering the SOURCE row. Any third use is the switch on identity that
	// Capabilities exists to replace.
	Provider    provider.ID
	ProviderRef string

	Title    string
	Platform string

	// SortTitle is the collation key — articles stripped, diacritics folded.
	//
	// It is stored rather than computed on read because Phase 1's only order
	// is title A→Z (07 §2) and the alternative is a full-table sort in Go on
	// every keystroke of a filter. fft-database owns its exact derivation.
	SortTitle string

	// PlaytimeMinutes is provider-reported; nil means not reported.
	PlaytimeMinutes *int

	// LastPlayed is provider-reported; nil means not reported, which is not
	// the same as never played.
	LastPlayed *time.Time

	// OwnedSince is when the player acquired it, where knowable.
	OwnedSince *time.Time

	// AddedAt is when Zerado first saw this row. Always known, because Zerado
	// is the one that wrote it.
	AddedAt time.Time

	// StatusManual is the player's explicit choice, or nil when they have
	// never expressed one.
	//
	// A sync never writes this field. That is the invariant that makes the
	// product trustworthy: mark a game ZERADO, play three more hours, and the
	// next sync updates playtime and last-played and nothing else.
	StatusManual *status.Status

	// StatusChangedAt is when StatusManual last changed, and nil when it never
	// has.
	//
	// It is the timestamp that decides a Phase 4 conflict — last-write-wins
	// per game — and it is in Phase 1's schema for that reason alone, which
	// 05-state-machine.md §8 names as the expensive decision made early.
	StatusChangedAt *time.Time

	// AbsentSince is set when a complete, successful sync stopped returning
	// this game, and cleared the moment it comes back.
	//
	// It is not a fifth state. It is an orthogonal presence flag, and the row
	// still has a status — usually the most valuable one, because a game you
	// finished and no longer own is exactly the row you would be angriest to
	// lose. Absent rows are excluded from the default view, remain findable by
	// filter, and are never deleted as a consequence of a sync.
	AbsentSince *time.Time

	// MergedInto points at the row this one was merged into, in Phase 4.
	//
	// Present in the Phase 1 type because the column is present in the Phase 1
	// schema, so that a merge never has to rewrite primary keys. Nothing in
	// Phase 1 sets it and every Phase 1 query excludes rows that have it.
	MergedInto *GameID
}

// Absent reports whether a sync has stopped returning this game.
func (g Game) Absent() bool { return g.AbsentSince != nil }

// Playtime returns the reported playtime and whether it was reported.
func (g Game) Playtime() (int, bool) {
	if g.PlaytimeMinutes == nil {
		return 0, false
	}
	return *g.PlaytimeMinutes, true
}

// Status returns the effective status: the manual value if the player set one,
// otherwise the derivation.
//
// canReportPlaytime is the provider's Capabilities.Playtime. It is a parameter
// rather than a field because a capability is a property of the provider and
// storing a copy of it on every row would be a second, stale writer of a fact
// the registry already holds — the same argument that keeps effective_status
// out of the schema.
func (g Game) Status(canReportPlaytime bool) status.Status {
	pt, _ := g.Playtime()
	return status.Effective(g.StatusManual, pt, canReportPlaytime)
}

// Overridden reports whether the player has set a manual status.
//
// Z-06 shows a fifth item — "clear override" — only when this is true, and
// Z-05 shows the SET BY and <PROVIDER> SAYS rows only when it is true.
func (g Game) Overridden() bool { return g.StatusManual != nil }
