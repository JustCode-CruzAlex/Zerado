package library

import "github.com/JustCode-CruzAlex/Zerado/internal/provider"

// Ref is enough to identify a game to an outside source, without assuming
// that source knows anything about stores.
//
// It is what the metadata and price seams take. Both of them are asked about
// games that may have arrived from Steam, from a shelf, or from a store that
// does not exist yet, so the ref carries what every source can use — a title
// and a platform — plus the store identifiers a source may use if it happens
// to recognise them.
//
// The ordering matters. A ref whose primary key were the Steam appid would
// work beautifully and would be unanswerable for a cartridge, which is the
// majority case for a shelf. Title and platform are the fields that are always
// present because [Game] requires both; everything else is a bonus a source is
// free to ignore.
type Ref struct {
	// Title and Platform are always present.
	Title    string
	Platform string

	// Provider and ProviderRef are the store's own identifiers, when the game
	// came from a store.
	//
	// A metadata source that recognises Steam appids gets a far better match
	// from these than from a title. A source that does not, ignores them. What
	// neither may do is require them, which is what makes this seam answerable
	// for a cartridge.
	Provider    provider.ID
	ProviderRef string

	// Extra carries additional identifiers a provider emitted — the same
	// narrow escape hatch as provider.Item.Extra, and under the same rule: a
	// key is read only by a source that recognises it, never by a screen.
	Extra map[string]string
}

// RefOf builds a Ref from a stored game.
func RefOf(g Game) Ref {
	return Ref{
		Title:       g.Title,
		Platform:    g.Platform,
		Provider:    g.Provider,
		ProviderRef: g.ProviderRef,
	}
}

// Identifiable reports whether this ref carries enough to ask anyone about.
//
// Title and platform are the floor. A ref without them is a bug in the caller,
// and a source asked with one should return fault.KindMalformed rather than a
// guess.
func (r Ref) Identifiable() bool { return r.Title != "" && r.Platform != "" }
