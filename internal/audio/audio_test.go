package audio_test

import (
	"testing"

	"github.com/JustCode-CruzAlex/Zerado/internal/audio"
	"github.com/JustCode-CruzAlex/Zerado/internal/audio/audiotest"
)

// TestNullIsTheDefaultBuildAndIsHonestAboutIt. NullAudio is not a fallback
// that might be untested; it is the default build, exercised by every test run
// and every CI job — and Settings says "not compiled" rather than "off",
// because those are different facts.
func TestNullIsTheDefaultBuildAndIsHonestAboutIt(t *testing.T) {
	var a audio.Audio = audio.Null{}
	s := a.State()
	if s.Compiled {
		t.Fatal("the default build claims the audio subsystem is compiled in")
	}
	if s.ReasonKey == "" {
		t.Fatal("Settings would show silence with no explanation")
	}
	if !a.Muted(audio.ChannelMusic) || !a.Muted(audio.ChannelFX) {
		t.Fatal("a build that makes no sound reports itself unmuted; the footer indicator would imply a volume that does nothing")
	}
	if err := a.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// TestCueCannotFail is the signature carrying the guarantee. There is no audio
// failure a screen could usefully handle, and an error return invites a caller
// to block on it.
func TestCueCannotFail(t *testing.T) {
	var cue func(audio.Cue) = audio.Null{}.Cue
	for _, c := range []audio.Cue{audio.CueMove, audio.CueSelect, audio.CueZerado, audio.CueSyncDone, audio.CueSyncFailed, audio.CueRefuse} {
		cue(c) // no error to check, and nothing to wait for
	}
}

// TestNoCueForCancelledOrPartial: neither is a failure, and a sound would make
// them feel like one. The cue set is the enumeration of that decision.
func TestNoCueForCancelledOrPartial(t *testing.T) {
	r := audiotest.New()
	r.Cue(audio.CueSyncDone)
	r.Cue(audio.CueSyncFailed)
	if got := len(r.Cues()); got != 2 {
		t.Fatalf("recorded %d cues, want 2", got)
	}
	// There is deliberately no CueSyncCancelled or CueSyncPartial to record.
	// If one were added, this test's comment is where the argument against it
	// lives.
}

// TestTheTwoChannelsAreIndependent: someone may want keyclicks without the
// soundtrack, or the reverse, and neither is the odd request.
func TestTheTwoChannelsAreIndependent(t *testing.T) {
	r := audiotest.New()
	r.Mute(audio.ChannelMusic, true)
	if !r.Muted(audio.ChannelMusic) {
		t.Fatal("muting music did not take")
	}
	if r.Muted(audio.ChannelFX) {
		t.Fatal("muting music also muted the interface cues")
	}
}
