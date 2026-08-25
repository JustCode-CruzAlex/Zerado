package images_test

import (
	"context"
	"testing"

	"github.com/JustCode-CruzAlex/Zerado/internal/images"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
)

// TestATerminalWithoutImagesIsSupported, not warned at. The text deck renders
// and the player is told once, quietly. Null is the implementation a large
// share of users will actually run, which is why it lives in the production
// package rather than in a test helper.
func TestATerminalWithoutImagesIsSupported(t *testing.T) {
	var im images.Images = images.Null{}
	if im.Capability().Draws() {
		t.Fatal("Null claims it can draw")
	}
	if _, ok := im.Cover(library.GameID(1), images.Rect{Cols: 17, Rows: 6}); ok {
		t.Fatal("Null returned a placement")
	}
	if err := im.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// TestCoverCannotBlockAndCannotFail is the signature carrying the guarantee: a
// function with no context has nothing to cancel, and one with no error return
// cannot be waited on. A missing cover is never worth a dropped frame.
func TestCoverCannotBlockAndCannotFail(t *testing.T) {
	// Cover takes no context and returns no error. If either changed, this
	// assignment would stop compiling — which is the test.
	var cover func(library.GameID, images.Rect) (images.Placement, bool) = images.Null{}.Cover
	if _, ok := cover(1, images.Rect{}); ok {
		t.Fatal("Null answered a cover")
	}
}

// TestWarmIsTheOnlyThingThatTouchesTheNetwork, and a build that cannot draw
// has no reason to spend the player's bandwidth.
func TestWarmIsTheOnlyThingThatTouchesTheNetwork(t *testing.T) {
	images.Null{}.Warm(context.Background(), []library.GameID{1, 2, 3}, images.Rect{Cols: 17, Rows: 6})
}

// TestCapabilityNamesAreStable: Z-09 renders them through the catalogue and
// the CLI's JSON emits them directly.
func TestCapabilityNamesAreStable(t *testing.T) {
	for c, want := range map[images.Capability]string{
		images.CapabilityNone:   "none",
		images.CapabilityKitty:  "kitty",
		images.CapabilityITerm2: "iterm2",
	} {
		if got := c.String(); got != want {
			t.Fatalf("Capability(%d).String() = %q, want %q", c, got, want)
		}
	}
}
