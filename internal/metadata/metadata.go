// Package metadata is the enrichment seam: cover art, sinopse, genres and
// release data, from a source Zerado deliberately does not name in its own
// types.
//
// # The seam is a hedge, and the hedge outlived its original reason
//
// It was designed when IGDB's free-for-non-commercial terms sat badly against
// an affiliate-funded product. Founder direction on 2026-08-25 dropped
// affiliate links entirely, so Zerado is cleanly non-commercial — free
// software, donation-supported, zero revenue — and IGDB's published test,
// whether the project generates revenue, is now answered.
//
// That reading is a reading of a published rationale and not a guarantee: a
// direct confirmation from IGDB is a founder action, not a resolved fact
// (13-handoffs.md §5.2). The seam therefore stays exactly as provider-agnostic
// as it was designed to be, and removing the hedge because one risk receded
// would be the wrong lesson — a source that is named today can change its
// terms tomorrow, which is true of every source and not a fact about IGDB.
//
// Nothing in this package names a provider. The word IGDB appears in this
// comment and nowhere in the API, which is the test of whether a seam is
// genuinely agnostic: swapping the source must change one implementation and
// no signature above it.
//
// # Having none is a designed state, not an error
//
// [Null] is a first-class implementation. When there is no metadata provider,
// or when the one there is returns nothing, the detail view shows a designed
// no-metadata composition rather than an error banner. That is the difference
// between a product that works without a metadata source and a product that is
// broken without one, and it is why fault.KindNotFound's treatment is
// TreatmentDesignedEmpty rather than a refusal.
package metadata

import (
	"context"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/aged"
	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
)

// Provider is a source of enrichment for a game.
//
// It is a lookup, not a sync: enrichment is per-game and on demand, because
// the alternative is fetching 247 records to render twelve rows. The Phase 2
// enrichment pass that walks the library is a caller of this, and it is the
// caller that owns the worker count and the rate limiting — not this
// interface, which stays a single question about a single game.
type Provider interface {
	// ID is the source's stable identity, stored on the metadata row so a
	// library enriched by one source can be re-enriched by another without
	// guessing which rows came from where.
	ID() provider.ID

	// Lookup returns what this source knows about one game.
	//
	// The result is stamped. There is no way to obtain a Metadata without its
	// observation time, which is 07 §4's age rule made structural rather than
	// remembered.
	//
	// A game the source has never heard of is fault.KindNotFound, and that is
	// routine rather than a failure: a hand-added cartridge is the majority
	// case for a shelf and the designed no-metadata composition is what it
	// renders. A source that is down is KindUnreachable, which is a different
	// screen, which is the entire reason the taxonomy exists.
	//
	// ctx cancellation must abort in-flight I/O. The Phase 2 enrichment pass
	// runs many of these concurrently and abandons them when the player
	// leaves the screen.
	Lookup(ctx context.Context, r library.Ref) (aged.Value[Metadata], error)

	// Attribution is the credit this source requires, and it is a method
	// rather than a configuration string.
	//
	// The credit a source requires is a property of the source, so swapping
	// the source swaps the credit automatically and a licence change becomes a
	// data change instead of a redesign. It is required and it is not optional
	// to render: the detail view renders whatever this returns.
	Attribution() Attribution
}

// Metadata is what a source knows about a game.
//
// It carries no age of its own: it is always handed back inside an
// [aged.Value], which is the mechanism rather than a convention. A FetchedAt
// field here would be a second, ignorable copy of the same fact.
type Metadata struct {
	// Sinopse is the description. The Portuguese word is ratified brand
	// vocabulary and is carried with a translator note so no future
	// translation "fixes" it.
	Sinopse string

	// CoverRef is a local cache path, never a remote URL.
	//
	// Nothing renders from the network, so a URL here would be a value no
	// renderer could legally use — and the one a renderer would eventually
	// use anyway. The cache lives in the XDG cache directory, outside
	// library.db, because covers are disposable, refetchable and potentially
	// large, and the file the player backs up has to stay small enough that
	// they actually back it up.
	CoverRef string

	// ReleasedAt is nil when the source did not say.
	ReleasedAt *time.Time

	// Genres is whatever taxonomy the source uses, unmapped.
	//
	// Deliberately not normalised into a Zerado genre enum. A shared
	// vocabulary across sources is a mapping table that would have to be
	// maintained against somebody else's taxonomy forever, and Phase 2 renders
	// these as text.
	Genres []string
}

// Empty reports whether this metadata carries nothing worth rendering.
//
// A source may legitimately return a record with every field blank. That is
// the designed no-metadata composition's trigger just as much as a KindNotFound
// is, and a caller that only checked the error would render an empty pane with
// three blank labels.
func (m Metadata) Empty() bool {
	return m.Sinopse == "" && m.CoverRef == "" && m.ReleasedAt == nil && len(m.Genres) == 0
}

// Attribution is the credit a source requires.
type Attribution struct {
	// TextKey names the catalogue entry that renders the credit line, with
	// the source's name as an argument.
	//
	// A key rather than a sentence, for the same D9 reason as everything else
	// — and with one deliberate consequence: a source whose licence requires
	// a specific untranslated wording supplies that wording through
	// [Attribution.Verbatim] instead, so the catalogue is never asked to
	// translate a legal obligation.
	TextKey string

	// Verbatim is a credit line that must appear exactly as given, for a
	// licence that requires specific wording.
	//
	// When it is non-empty it wins over TextKey. This is the one sanctioned
	// user-facing string in the seams, and it is sanctioned because it is not
	// language: it is a term of somebody else's licence, and translating it
	// would breach the thing it exists to satisfy.
	Verbatim string

	// URL is the source's own link, where its terms require one.
	URL string
}

// Null is the metadata provider that knows nothing.
//
// It is a designed state and not a fallback, which is why it lives in this
// package rather than in a test helper: it is the implementation Phase 1 runs
// with, and it is what the product uses if a source is ever withdrawn. Because
// it is the default rather than an emergency, the no-metadata composition is
// the well-exercised path instead of an untested one — the same argument that
// makes NullAudio the default build.
type Null struct{}

// ID returns the reserved identity of the no-source source.
func (Null) ID() provider.ID { return "none" }

// Lookup always reports that nothing is known.
//
// It returns KindNotFound rather than an empty success so that a caller
// distinguishing "nothing is known" from "the record is blank" gets the same
// answer from Null as it would from a real source that had never heard of the
// game.
func (Null) Lookup(context.Context, library.Ref) (aged.Value[Metadata], error) {
	return aged.Value[Metadata]{}, fault.New(fault.KindNotFound, "metadata.Null.Lookup")
}

// Attribution returns no credit, because there is no source to credit.
func (Null) Attribution() Attribution { return Attribution{} }

// Null satisfies Provider — asserted here so the claim is checked by the
// compiler rather than by a reader.
var _ Provider = Null{}
