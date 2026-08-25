// Package providertest holds provider implementations that reach no network,
// so that every screen and every contract test can run with the machine in a
// Faraday cage.
//
// It contains two things, and the distinction between them is the point:
//
//	Manual  — the shape a real hand-entry provider has. It is here rather than
//	          in a provider package because Phase 1 builds no providers; it is
//	          the reference implementation the physical provider will be, and
//	          it is what proves the seam is not Steam-shaped.
//	Fake    — a scriptable Syncer. It is a test double and nothing else.
package providertest

import (
	"fmt"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
)

// Manual is a provider whose items arrive because a person typed them.
//
// It is the case the store-provider seam was designed against first, and it is
// the best available test of whether that seam is secretly Steam-shaped. Read
// it and note what is absent: no credentials, no pagination, no rate limit, no
// HTTP client, no Sync method, and no playtime. Every one of those absences is
// a place a one-interface design would have forced this provider to lie.
//
// It implements [provider.Provider] and [provider.Enterer] and NOT
// [provider.Syncer]. That is asserted at the bottom of this file, and a
// contract test asserts the negative at run time — because the negative is the
// decision, and a compiler cannot check the absence of an implementation.
type Manual struct {
	// NewRef mints a provider reference for a new entry.
	//
	// It is injected so a test gets deterministic refs. A real physical
	// provider mints a UUID, which is why Z-08 never asks the player for an
	// app ID: a cartridge has none, and requiring one would have been the
	// first Steam-shaped assumption to reach a screen.
	NewRef func() string

	// Now supplies the clock. 06-data-seams.md §7 declines to make the clock
	// a seam — a func field on the structs that need it is enough, and an
	// interface for this is ceremony.
	Now func() time.Time
}

// NewManual returns a Manual with deterministic test defaults.
func NewManual() *Manual {
	var n int
	base := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	return &Manual{
		NewRef: func() string { n++; return fmt.Sprintf("manual-%04d", n) },
		Now:    func() time.Time { return base },
	}
}

// ID returns the physical provider's identity.
func (m *Manual) ID() provider.ID { return "physical" }

// Display returns the player-facing name.
//
// A real implementation renders this through the catalogue, because "Physical
// shelf" is a phrase rather than a brand. The literal here is a test double's,
// and it is the one place in this repository a user-facing English string is
// acceptable — nothing ships it.
func (m *Manual) Display() string { return "Physical shelf" }

// Capabilities declares what a shelf can and cannot know.
//
// Sync is false because there is nothing to fetch. Playtime and LastPlayed are
// false because a cartridge has no telemetry — and that single fact is what
// makes all four states manual for this provider, without status.Derive
// containing a special case for it. OwnedSince is true because a player
// entering a cartridge genuinely knows when they got it.
func (m *Manual) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Sync:        false,
		Manual:      true,
		Playtime:    false,
		LastPlayed:  false,
		OwnedSince:  true,
		Credentials: nil, // a shelf has no key
	}
}

// Form returns the four fields Z-08 renders, of which two are required.
func (m *Manual) Form() []provider.EntryField {
	return []provider.EntryField{
		{Key: "title", LabelKey: "entry.title.label", HelpKey: "entry.title.help", Kind: provider.FieldText, Required: true},
		{Key: "platform", LabelKey: "entry.platform.label", HelpKey: "entry.platform.help", Kind: provider.FieldText, Required: true},
		{
			Key: "state", LabelKey: "entry.state.label", Kind: provider.FieldChoice,
			Options: []provider.FieldOption{
				{Value: "not_started", LabelKey: "state.not_started"},
				{Value: "in_progress", LabelKey: "state.in_progress"},
				{Value: "zerado", LabelKey: "state.zerado"},
				{Value: "abandoned", LabelKey: "state.abandoned"},
			},
		},
		{Key: "owned_since", LabelKey: "entry.owned_since.label", HelpKey: "entry.owned_since.help", Kind: provider.FieldDate},
	}
}

// Compose turns a completed form into an item, minting the reference itself.
//
// It validates and does not write: the store is the only writer, and returning
// an Item means hand entry travels the same upsert path a sync does. That
// shared path is what makes a physical copy structurally the same kind of
// thing as a Steam game rather than a flag on one.
//
// Note that it never sets PlaytimeMinutes, not even to zero. Capabilities says
// this provider cannot report playtime, and a zero would be a number
// status.Derive would treat as evidence.
func (m *Manual) Compose(e provider.Entry) (provider.Item, error) {
	title := trim(e["title"])
	platform := trim(e["platform"])
	if title == "" {
		return provider.Item{}, fault.New(fault.KindMalformed, "physical.Compose",
			fault.WithMessage("entry.title.required"))
	}
	if platform == "" {
		return provider.Item{}, fault.New(fault.KindMalformed, "physical.Compose",
			fault.WithMessage("entry.platform.required"))
	}

	it := provider.Item{
		ProviderRef: m.NewRef(),
		Title:       title,
		Platform:    platform,
	}
	if s := trim(e["owned_since"]); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return provider.Item{}, fault.New(fault.KindMalformed, "physical.Compose",
				fault.WithMessage("entry.owned_since.unparsed"), fault.WithCause(err))
		}
		it.OwnedSince = &t
	}
	return it, nil
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && isSpace(s[start]) {
		start++
	}
	for end > start && isSpace(s[end-1]) {
		end--
	}
	return s[start:end]
}

func isSpace(b byte) bool { return b == ' ' || b == '\t' || b == '\n' || b == '\r' }

// Manual implements Provider and Enterer. It does NOT implement Syncer, and
// there is deliberately no assertion here that it does — the absence is the
// design, and TestManualIsNotASyncer checks it at run time.
var (
	_ provider.Provider = (*Manual)(nil)
	_ provider.Enterer  = (*Manual)(nil)
)
