// Package audio is the sound seam: two channels, off by default, streamed
// rather than bundled, and removable at runtime and at compile time.
//
// ADR-0001 D6 ships audio in Phase 1 as an opt-in subsystem whose music is
// internet radio the player chooses, with local interface cues as the only
// always-available part. Nothing ships in the binary, so there is nothing to
// license, attribute or weigh — the most expensive open question in the audio
// design was dissolved rather than answered.
//
// # The signatures carry the guarantees
//
// [Audio.Cue] has no error return and no context. That is not an oversight and
// it is not minimalism: there is no audio failure a screen could usefully
// handle, and an error return invites a caller to block on it. A cue that
// cannot play is silence, and silence is not a failure a screen must render.
//
// [Audio.State] exists because Z-09 must say honestly *why* audio is
// unavailable — a device that failed to open, an environment variable that
// overrode it, a build without the subsystem in it. "Off" and "overridden" are
// different facts and Settings shows which.
//
// # Audio is never the only carrier of information
//
// The co-render rule extended to a fourth channel. The test is
// ZERADO_NO_AUDIO=1 and lose nothing — the same test NO_COLOR passes.
//
// # The build-tag split, and what this package does not decide
//
// The default build is pure Go and contains [Null], which keeps D2's
// single-binary cross-compile property and means the silent path is exercised
// by every test run rather than being an untested fallback. The real player
// sits behind a build tag.
//
// 13-handoffs.md §4 assigns the choice of *which* audio library to this
// ticket. It is not made here, and the reason is stated rather than skipped:
// the requirement is a real per-platform check of each candidate's cgo
// dependency, and that is a verification against six toolchains rather than a
// signature decision. It changes no signature in this file — which is what the
// build-tag seam was designed to guarantee — so it is recorded as still open
// in docs/api/07-open-questions.md rather than answered from memory.
package audio

// Audio is the whole surface a screen sees.
type Audio interface {
	// Cue plays a short interface sound. It is fire-and-forget and it MUST
	// NOT block.
	//
	// The implementation is a non-blocking send on a buffered channel: if the
	// buffer is full the cue is dropped, silently, because a missed sound is
	// not worth a dropped frame. A cue is always the second signal, never the
	// first — the visible change happens on the frame it happens on, and the
	// sound follows or does not.
	Cue(c Cue)

	// Music turns the radio stream on or off.
	//
	// Off releases the audio device AND closes the stream. It does not hold a
	// gain of zero and it does not keep pulling bytes: a muted Zerado that
	// kept taking a station's stream would be spending the player's bandwidth
	// on silence, and it would sit badly against the published promise that
	// the only network traffic is Zerado reaching out to services the player
	// connected. Reconnecting on unmute costs a moment of buffering, which is
	// the correct trade.
	Music(on bool)

	// SetVolume sets a channel's volume, 0..100. Out-of-range values are
	// clamped rather than rejected: there is no failure here a screen could
	// present.
	SetVolume(ch Channel, v int)

	// Mute mutes or unmutes one channel. The two channels are independent
	// because someone may want keyclicks without the soundtrack, or the
	// reverse, and neither is the odd request.
	Mute(ch Channel, on bool)

	// Muted reports a channel's mute state, for the footer indicator.
	Muted(ch Channel) bool

	// State is what Z-09 renders, including why audio is unavailable.
	State() State

	// Close stops the owned goroutine and releases the device.
	//
	// It has its own timeout: a stuck audio device does not stop q from
	// quitting. The error is for logs; nothing waits on it.
	Close() error
}

// Channel is one of the two independently controlled outputs.
type Channel uint8

const (
	// ChannelMusic is the radio stream. It NEEDS THE NETWORK: offline it
	// stops, and that is an honest degradation of a feature that is online by
	// nature rather than a broken promise.
	ChannelMusic Channel = iota

	// ChannelFX is local interface cues. They WORK offline, because they are
	// short, local, and reach for nothing.
	ChannelFX
)

// String returns the stable machine name, used as a settings key suffix.
func (c Channel) String() string {
	if c == ChannelMusic {
		return "music"
	}
	return "fx"
}

// Cue names an interface sound.
//
// It is an enum rather than a file path so that a screen names an *event* and
// the audio implementation owns which sound an event makes — which is what
// lets a theme change the sound set without a screen changing.
type Cue uint8

const (
	// CueMove is a list cursor moving.
	CueMove Cue = iota

	// CueSelect is a choice being committed.
	CueSelect

	// CueZerado is the one cue that marks an achievement: a game becoming
	// zerado. It is still never the only carrier — the chip changes colour,
	// glyph and label on the same frame.
	CueZerado

	// CueSyncDone is a sync completing successfully. There is deliberately no
	// cue for PARTIAL or CANCELLED: neither is a failure, and a sound would
	// make them feel like one.
	CueSyncDone

	// CueSyncFailed is a sync failing.
	CueSyncFailed

	// CueRefuse is an action the product declined.
	CueRefuse
)

// State is what Settings reports about audio, honestly.
type State struct {
	// Compiled reports whether this build contains the audio subsystem at
	// all. False for the default pure-Go build.
	Compiled bool

	// Enabled is the player's own opt-in.
	Enabled bool

	// Available reports whether a device was successfully opened.
	//
	// Distinct from Enabled: a player can have turned audio on and have no
	// device, which is a real state on a headless machine and one Settings
	// must name rather than showing "On" over silence.
	Available bool

	// EnvOverride reports that ZERADO_NO_AUDIO is set.
	//
	// It is a separate field from Enabled because Z-09 must show "overridden",
	// not "off" — those are different facts, and collapsing them would tell a
	// player their setting had been changed when it had not.
	EnvOverride bool

	// ReasonKey names the catalogue entry explaining why audio is
	// unavailable, when it is. Empty when it is available.
	ReasonKey string

	// DeviceKey names the catalogue entry describing the output — the
	// "CoreAudio, built-in output" line. A key with arguments rather than a
	// sentence, because the device name is a substitution and the phrasing
	// around it is language.
	DeviceKey string

	// Volume and Muted are per channel, indexed by Channel.
	Volume [2]int
	Muted  [2]bool
}

// Null is the audio implementation that makes no sound.
//
// It is the DEFAULT build, not a fallback, which is the whole point of D6's
// build-tag split: the silent path is exercised by every test run and every CI
// job rather than being the untested branch that only fails on somebody's
// laptop.
type Null struct{}

// Cue does nothing, immediately.
func (Null) Cue(Cue) {}

// Music does nothing.
func (Null) Music(bool) {}

// SetVolume does nothing.
func (Null) SetVolume(Channel, int) {}

// Mute does nothing.
func (Null) Mute(Channel, bool) {}

// Muted reports true for both channels: a build that makes no sound is
// honestly muted, and the footer indicator says so rather than implying a
// volume that does nothing.
func (Null) Muted(Channel) bool { return true }

// State reports a build with no audio in it, and names why.
func (Null) State() State {
	return State{Compiled: false, ReasonKey: "settings.audio.reason.not_compiled"}
}

// Close does nothing and cannot fail.
func (Null) Close() error { return nil }

var _ Audio = Null{}
