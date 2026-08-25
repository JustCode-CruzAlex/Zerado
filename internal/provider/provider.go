// Package provider is the seam every source of owned games satisfies —
// including the ones that are a person with a keyboard.
//
// # Designed against the cartridge, not against Steam
//
// The likeliest way to get this seam wrong is to design it around Steam,
// because Steam is the only provider Phase 1 builds. So the shape below was
// written against manual physical entry first and checked against Steam
// second, which is the reverse of the tempting order and the reason the
// contract has two capability interfaces rather than one.
//
// A physical cartridge is not a store. It has no API, no credentials, no
// pagination, no rate limit, no playtime and no last-played date. Its "sync"
// is a person typing into Z-08. Every one of those absences is a place where a
// Steam-shaped interface would have forced physical to lie — by implementing
// Sync and returning an error, by declaring credential fields it does not
// want, or by reporting a playtime of zero that the derivation would then
// treat as a fact.
//
// # Interface segregation is the whole design
//
//	Provider   — every source implements it. Identity, display name, capabilities.
//	Syncer     — a source that can fetch a library over a network.
//	Enterer    — a source whose items arrive because a person entered them.
//
// steam implements Provider + Syncer. physical implements Provider + Enterer.
// A future GOG that both syncs and lets a player add a game the API missed
// implements all three, and nothing above it changes.
//
// ADR-0001 D1 rejects the one-interface alternative for a stated reason: if
// physical implemented Syncer and returned ErrManual, every caller would have
// to know which providers lie about their own interface, and the error would
// be a runtime rediscovery of a fact the type system could have carried.
//
// # Everything downstream reads Capabilities, never the ID
//
// The row renderer, the state derivation, the filter, the summary counts and
// the detail view all read [Capabilities]. A switch on [ID] anywhere above this
// package is the defect this seam exists to prevent, and it is the defect that
// makes adding GOG a six-file change instead of a one-package change.
//
// # No provider constructs its own HTTP client
//
// 06-data-seams.md §7 makes this explicit and 07-offline-contract.md §7.3
// makes it a review rule: a shared *http.Client with a timeout is injected.
// That is what keeps "works offline" from quietly stopping being true, and it
// is why this package exposes no client, no transport and no URL.
package provider

import "context"

// ID is a provider's stable machine identity: "steam", "physical",
// "playstation", "gog", "ea".
//
// It is a database value and a CLI argument, so it is API. It is *not* a
// dispatch key: code above this package that switches on an ID has bypassed
// [Capabilities] and will be wrong about the next provider. The one legitimate
// use above the seam is looking a provider up in a [Registry].
type ID string

// Provider is what every source of games implements, including the ones that
// are a human with a keyboard.
//
// It is deliberately tiny. Everything a screen needs to decide what to render
// is here, and nothing here does I/O — which is what lets Z-02 build its form,
// Z-04 render a row and Z-05 decide which fields exist, all with the network
// off and with no provider having been contacted.
type Provider interface {
	// ID returns the stable machine identity.
	ID() ID

	// Display returns the player-facing name — "Steam", "Physical shelf".
	//
	// It is a plain string rather than an i18n key, and that is a considered
	// exception to ADR-0001 D9 rather than a hole in it. A provider's name is
	// a proper noun: "Steam" is "Steam" in every language, and routing it
	// through a catalogue would invite a translator to translate a brand.
	// Where a provider's name is genuinely a phrase rather than a brand —
	// "Physical shelf" — the provider returns a catalogue-rendered string,
	// because it holds a printer and a screen does not.
	Display() string

	// Capabilities reports what this provider can actually do.
	//
	// It must be a pure function of the provider's own nature — never of its
	// current connection state, never of whether the network is up, and never
	// of the player's settings. A capability that changes when a request fails
	// would make the *shape* of a screen flicker with the weather, and 07 §1
	// is explicit that a feature is never in two offline classes.
	Capabilities() Capabilities
}

// Capabilities is what a provider can do, in the terms screens actually ask
// about.
//
// Every field answers a question some screen has. Nothing here is
// informational: a capability nobody reads is a capability that will drift out
// of truth unnoticed.
type Capabilities struct {
	// Sync reports whether this provider can fetch a library over a network.
	//
	// True exactly when the provider also implements [Syncer]. The redundancy
	// is deliberate: a screen needs to decide whether to offer "sync" before
	// it has a reason to type-assert anything, and 07's offline classification
	// is a property of the provider rather than of a Go type assertion.
	// [Check] asserts the two agree.
	Sync bool

	// Manual reports whether items may be entered by hand for this provider.
	//
	// True exactly when the provider also implements [Enterer]. This is the
	// field Z-01's "Add a game by hand" door and Z-08 read, and it is the
	// field that makes physical a first-class source rather than a flag.
	Manual bool

	// Playtime reports whether this provider can ever know how long a game
	// has been played.
	//
	// It drives status.Derive, and it is the difference between "0h" and "not
	// tracked" on Z-05 — two facts that look identical in a database and are
	// completely different on a screen. A cartridge has no telemetry, so
	// physical reports false and all four of its states are the player's.
	Playtime bool

	// LastPlayed reports whether this provider can supply a last-played
	// timestamp. When false, Z-03's DONE line omits the third fact entirely
	// rather than rendering it empty.
	LastPlayed bool

	// OwnedSince reports whether an acquisition date is available. It is the
	// one optional capability physical claims, because a player entering a
	// cartridge genuinely knows when they got it.
	OwnedSince bool

	// Credentials lists the fields Z-02 renders, in the order it renders
	// them. Empty means the provider needs none — which is the physical case,
	// and the reason Z-02 is reachable only for providers that have some.
	Credentials []CredentialField
}

// Syncer is implemented only by providers that can fetch a library.
//
// physical does not implement it, and that is the point: it is not a
// Steam-shaped provider with a hole where Sync should be.
type Syncer interface {
	Provider

	// Sync streams this provider's current view of the library.
	//
	// It returns as soon as the request is under way — before any item has
	// arrived — so that Z-03 can paint its indeterminate state immediately
	// rather than after a round trip. A failure that is known at that moment
	// (no route, a rejected key) is returned as the error; a failure that
	// happens later, after items have already been delivered, arrives through
	// [Stream.Err] and is the PARTIAL case.
	//
	// Cancelling ctx MUST abort in-flight I/O and close the stream promptly.
	// Z-03 lets the player press Esc mid-sync and the goroutine budget for the
	// whole product is zero leaks, so a provider that ignores cancellation is
	// not merely impolite — it holds the process open at quit.
	//
	// What arrives before a cancel stays: 06 §2.5 requires that a cancel
	// mid-sync leaves a valid partial library, which is why this streams at
	// all rather than returning a slice.
	Sync(ctx context.Context, c Credentials) (Stream, error)
}

// Enterer is implemented only by providers whose items arrive because a person
// entered them.
//
// This is the interface that keeps physical honest. Without it, hand entry
// would be a screen writing directly into the store, and the two rules that
// make a hand-entered row a real row — that its provider_ref is a UUID Zerado
// mints rather than something the player has to find, and that its optional
// fields respect the provider's own capabilities — would live in the screen
// instead of in the provider.
//
// It is a capability, not a provider kind. A future GOG that syncs *and* lets
// a player add a game its API missed implements both this and [Syncer], and
// Z-08 gains a source picker rather than a special case.
type Enterer interface {
	Provider

	// Form returns the fields Z-08 renders, in order.
	//
	// It is the hand-entry twin of Capabilities().Credentials, and it exists
	// for the same reason: the screen renders what the provider declares, so
	// a new hand-enterable source adds zero screens.
	Form() []EntryField

	// Compose turns a completed form into an [Item].
	//
	// The provider mints the ProviderRef — for physical, a fresh UUID, so a
	// duplicate can never collide and 06 §2.2's uniqueness constraint holds
	// without the player being asked for an app ID a cartridge does not have.
	//
	// It validates and it does not write. The store is the only writer, and
	// Compose returning an Item rather than performing an insert is what keeps
	// hand entry on exactly the same path as a sync: both produce Items, and
	// both go through the same upsert.
	//
	// A validation failure is a fault.KindMalformed carrying the offending
	// field's key, so Z-08 can put the message under the field it belongs to
	// rather than in a banner.
	Compose(e Entry) (Item, error)
}
