// Package devicesync sketches the Phase 4 client/server boundary.
//
// # This is a sketch and it builds nothing
//
// No Phase 1 code calls it. It exists because ADR-0001 D4 decides *what
// crosses* the boundary now rather than in Phase 4, and because that decision
// is the most expensive one in the bundle to reverse: it decides what the
// schema carries from the first migration onward. A sketch that compiles is a
// decision that can be checked; a paragraph is a decision that can be
// misremembered.
//
// # The rule, in one line
//
// Only what the player typed crosses. Everything a machine can recompute,
// each device recomputes.
//
//	Crosses:        the manual status and when it changed; hand-entered games;
//	                a stable identity; mood tags the player assigned (Phase 2).
//	Never crosses:  the library itself; playtime; last-played; cover art;
//	                sinopse; prices; credentials, ever.
//
// # The rule is enforced, not merely documented
//
// [Change] and [ManualGame] have no field that could carry a provider fact or
// a credential, and TestPayloadCarriesOnlyWhatThePlayerTyped walks their field
// sets against an allow-list. A future contributor adding PlaytimeMinutes to
// the payload — for entirely reasonable-looking reasons, in a Phase 4 sprint,
// long after everyone who read D4 has moved on — fails a test that names the
// decision and points at the ADR.
//
// That test is the deliverable here. The types are how it gets something to
// check.
package devicesync

import (
	"context"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/status"
)

// DeviceID identifies one of a player's machines.
//
// Minted locally on first run, stored in the settings table, and never derived
// from hardware: a hardware-derived identity is a fingerprint, and a
// local-first product that promises no telemetry should not be able to
// recognise a machine it was never told about.
type DeviceID string

// Change is one status change, which is the only library fact that crosses.
//
// It carries a UID rather than a GameID because a local surrogate key means
// nothing on another device — which is exactly why library.UID exists in
// Phase 1 and why adding it in Phase 4 would have required a migration that
// invents stable identities for rows whose titles the player had since edited.
type Change struct {
	// UID is the cross-device merge hint. It is a hint and not an authority:
	// an ambiguous match is shown to the player, never guessed.
	UID library.UID

	// Status is the player's manual choice, or nil when they cleared the
	// override. The nil is meaningful and must cross: "I have no opinion" is
	// a change the player made.
	Status *status.Status

	// ChangedAt decides the conflict. Last-write-wins per game.
	//
	// The limits are stated rather than hidden: two devices that both change a
	// status while offline lose the earlier change silently. That is adequate
	// because the conflicting parties are one person on two devices who agree
	// about what they did, and it is unacceptable to leave undocumented, so
	// Z-22 shows the last merge and what it resolved.
	ChangedAt time.Time
}

// ManualGame is a hand-entered game, which crosses because it is the one row
// class with no other copy.
//
// Everything else in the library is re-derivable from the player's own
// sources; a cartridge typed into Z-08 exists nowhere else, so losing it is
// unrecoverable. D4 therefore puts these in the *first* sync payload rather
// than a follow-up.
//
// Note what this struct does not have and cannot grow: no playtime, no
// last-played, no cover, no provider ref, no credentials. A hand-entered
// game's provider ref is a UUID this device minted and it means nothing on
// another machine.
type ManualGame struct {
	UID       library.UID
	Title     string
	Platform  string
	CreatedAt time.Time

	// OwnedSince is the one optional field, because it is the one optional
	// thing the player typed.
	OwnedSince *time.Time
}

// Envelope is one direction of one exchange.
type Envelope struct {
	// Device is which machine this came from, so a merge can say what it
	// resolved and against what.
	Device DeviceID

	// Since is the watermark the sender had before this exchange. The server
	// returns everything after it; there is no full-library transfer because
	// there is no library on the server to transfer.
	Since time.Time

	Changes []Change
	Manual  []ManualGame
}

// Client is the Phase 1 sketch of the Phase 4 boundary.
//
// It is two calls because the payload is small and idempotent, and because
// anything richer — a subscription, a stream, an operation log — would be
// designing for a conflict shape that does not exist. D4 rejects CRDTs and
// operation logs for exactly that reason: they are correct for concurrent
// multi-writer editing, and the parties here are one person on two devices.
//
// What Phase 1 must not make impossible, and does not:
//
//   - a stable identity exists on every row from the first migration;
//   - status_changed_at exists and is written on every change, including a
//     clear;
//   - merged_into exists so two rows can be joined without rewriting keys;
//   - the store is the single writer, so a Phase 4 merge attaches at that
//     seam without touching a screen.
//
// What is deliberately not decided here: whether Phase 4 accounts are
// email-and-password, OAuth or something else. Nothing in Phase 1 depends on
// it, and ADR-0001 names it as not decided.
type Client interface {
	// Push sends this device's changes since the last watermark.
	Push(ctx context.Context, e Envelope) (Receipt, error)

	// Pull returns other devices' changes since the watermark.
	Pull(ctx context.Context, device DeviceID, since time.Time) (Envelope, error)
}

// Receipt is what the server acknowledges.
type Receipt struct {
	// Watermark is the new "since" for this device's next exchange.
	Watermark time.Time

	// Accepted and Superseded let Z-22 report what a merge actually did —
	// which is the price of last-write-wins being an acceptable simplicity
	// choice rather than a silent one.
	Accepted   int
	Superseded int
}
