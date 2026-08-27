// Package audiotest records what a screen asked for, so the co-render rule can
// be tested: a cue may accompany a state change but must never be the only
// carrier of it.
package audiotest

import (
	"sync"

	"github.com/JustCode-CruzAlex/Zerado/internal/audio"
)

// Recorder is an Audio that makes no sound and remembers everything.
//
// It never blocks, because the contract says Cue cannot — and a recorder that
// took a lock a caller could contend on would be a test double that hid the
// one property worth testing. The mutex here is held for the length of an
// append and nothing else.
type Recorder struct {
	mu    sync.Mutex
	cues  []audio.Cue
	muted [2]bool
	on    bool
}

// New returns an empty Recorder.
func New() *Recorder { return &Recorder{} }

// Cue records the cue and returns immediately.
func (r *Recorder) Cue(c audio.Cue) {
	r.mu.Lock()
	r.cues = append(r.cues, c)
	r.mu.Unlock()
}

// Music records the requested state.
func (r *Recorder) Music(on bool) {
	r.mu.Lock()
	r.on = on
	r.mu.Unlock()
}

// SetVolume does nothing.
func (r *Recorder) SetVolume(audio.Channel, int) {}

// Mute records the mute state.
func (r *Recorder) Mute(ch audio.Channel, on bool) {
	r.mu.Lock()
	r.muted[ch] = on
	r.mu.Unlock()
}

// Muted reports the recorded mute state.
func (r *Recorder) Muted(ch audio.Channel) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.muted[ch]
}

// State reports a compiled, available, silent subsystem.
func (r *Recorder) State() audio.State {
	r.mu.Lock()
	defer r.mu.Unlock()
	return audio.State{Compiled: true, Enabled: true, Available: true, Muted: r.muted}
}

// Close does nothing.
func (r *Recorder) Close() error { return nil }

// Cues returns everything cued so far.
func (r *Recorder) Cues() []audio.Cue {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]audio.Cue(nil), r.cues...)
}

var _ audio.Audio = (*Recorder)(nil)
