// Package images is the terminal-image seam. A screen asks for a cover at a
// size; it never learns which protocol answered, or whether one did.
//
// ADR-0001 D8 makes cover art foundational rather than an enhancement — "a
// game library without cover art is a spreadsheet" — and makes a terminal
// without image support a *supported configuration* rather than a warning
// state. Both halves of that are in this interface: [Images.Cover] answers for
// any terminal, and [Capability] is the only place the protocol is named.
//
// # Two properties carry the whole design
//
// Cover never blocks and never fetches. The render path reads what is already
// present; a cover the cache does not hold is simply not shown this frame.
// This is the same rule as audio's Cue and for the same reason: a missing
// cover is never worth a dropped frame.
//
// A screen never learns the protocol. Adding Sixel later, or dropping iTerm2,
// changes one implementation and no screen.
package images

import (
	"context"

	"github.com/JustCode-CruzAlex/Zerado/internal/library"
)

// Images is the whole surface.
type Images interface {
	// Capability reports what this terminal can draw.
	//
	// Resolved once at start-up and cached for the session, never a config
	// flag the player has to find. Detection must fail closed: an ambiguous
	// or timed-out response means no image support, because guessing yes and
	// emitting escape sequences into a terminal that does not understand them
	// is how a library view turns into garbage on somebody's screen.
	Capability() Capability

	// Cover returns a placement for an already-cached image.
	//
	// It never fetches, never blocks and never returns an error. ok=false
	// means "not this frame", which the caller renders as the designed
	// no-cover tile — never a broken-image box, and never a spinner, because
	// there is nothing to wait for on this path.
	//
	// It takes no context, and that is the signature carrying the guarantee:
	// a function that cannot block has nothing to cancel, and a ctx parameter
	// here would be an invitation to make it blocking later.
	Cover(id library.GameID, cells Rect) (Placement, bool)

	// Warm asks the cache to fetch, off the render path, best-effort.
	//
	// It is the only method that touches the network, and it returns
	// immediately. ids is the set worth having soon — the visible tiles plus
	// a little ahead of the viewport, never the whole library, because
	// warming 247 covers to show 12 is how terminal images come to feel slow.
	//
	// Fetching is one owned worker pool bounded by ctx. The seam owns it; no
	// screen ever starts a goroutine for this.
	Warm(ctx context.Context, ids []library.GameID, cells Rect)

	// Forget drops a cached image, for the Settings action that clears the
	// cache. Deleting the whole cache must cost nothing but bandwidth.
	Forget(id library.GameID)

	// Close stops the worker pool and releases the cache handle.
	Close() error
}

// Capability is what the terminal can draw.
type Capability uint8

const (
	// CapabilityNone is a terminal that cannot draw images, or one where
	// detection was ambiguous, or where ZERADO_NO_IMAGES is set.
	//
	// It is a supported configuration. The full text deck renders, and the
	// player is told once, quietly and dismissibly, that Ghostty or Kitty
	// would show covers — never recurring, never blocking, and never phrased
	// as a fault in their setup.
	CapabilityNone Capability = iota

	// CapabilityKitty is the Kitty graphics protocol: Kitty, Ghostty,
	// WezTerm, Konsole.
	CapabilityKitty

	// CapabilityITerm2 is iTerm2's inline-image protocol.
	CapabilityITerm2
)

// String returns the stable machine name. Z-09 renders these through the
// catalogue; the CLI's JSON emits them directly.
func (c Capability) String() string {
	switch c {
	case CapabilityKitty:
		return "kitty"
	case CapabilityITerm2:
		return "iterm2"
	default:
		return "none"
	}
}

// Draws reports whether this capability can place an image at all.
func (c Capability) Draws() bool { return c != CapabilityNone }

// Rect is a size in terminal cells, not pixels.
//
// Cells because every layout decision in this product is made in cells, and a
// pixel size here would force each screen to know its terminal's cell
// geometry — which is the one measurement a terminal program cannot reliably
// obtain.
type Rect struct{ Cols, Rows int }

// Placement is an opaque instruction to draw an image at a position.
//
// It carries the escape sequence and the footprint, and a screen composes it
// into its frame without interpreting it. Opaque on purpose: the day Sixel is
// added, a screen that had been reading these bytes would be a screen that
// needs changing.
type Placement struct {
	// Sequence is the terminal control sequence that draws the image.
	Sequence string

	// Cells is the footprint the sequence will occupy, so the layout can
	// reserve exactly that much and the rest of the row still adds up.
	Cells Rect

	// ID is the transmitted image's identity in the terminal's own image
	// table, where the protocol has one.
	//
	// Present because re-sending image data every frame is the failure mode
	// that makes terminal images feel slow: Kitty placements are addressed by
	// image id, so an image is transmitted once and placed many times.
	ID uint32
}

// Null is the images implementation that draws nothing.
//
// It is what runs when detection fails, in tests, and under ZERADO_NO_IMAGES —
// so the text path is the well-exercised one rather than an untested fallback.
// It lives in this package rather than in a test helper for that reason: it is
// production code that a large share of users will actually run.
type Null struct{}

// Capability reports that nothing can be drawn.
func (Null) Capability() Capability { return CapabilityNone }

// Cover always declines, which renders the designed no-cover tile.
func (Null) Cover(library.GameID, Rect) (Placement, bool) { return Placement{}, false }

// Warm does nothing. It notably does not fetch: a build that cannot draw an
// image has no reason to spend the player's bandwidth on one.
func (Null) Warm(context.Context, []library.GameID, Rect) {}

// Forget does nothing.
func (Null) Forget(library.GameID) {}

// Close does nothing and cannot fail.
func (Null) Close() error { return nil }

var _ Images = Null{}
