package providertest

import (
	"context"
	"sync"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
)

// Fake is a scriptable Syncer that reaches no network.
//
// It exists to prove the property the ticket makes a requirement: every seam
// is fakeable with no network, and if it cannot be tested offline the shape is
// wrong. Everything Z-03 has to render is scriptable here — a slow first
// round trip, a known or unknown denominator, a stream that breaks halfway, a
// rejected key, and a successful call that returns nothing at all.
//
// It is a test double and not a reference implementation: a real Syncer's
// interesting parts are HTTP, pagination and rate limits, none of which this
// has. What it does share with a real one is the contract — items on a
// channel, cancellation honoured, terminal error after close, progress readable
// concurrently — which is the part a test needs to exercise.
type Fake struct {
	// Ident and Name are what the provider reports about itself.
	Ident provider.ID
	Name  string

	// Caps is what it claims it can do. Defaults to a Steam-shaped set.
	Caps provider.Capabilities

	// Items is the library it will stream.
	Items []provider.Item

	// FailBefore, when set, makes Sync return this error immediately, before
	// any item — the case that becomes Z-03's FAILED.
	FailBefore error

	// FailAfter, when set, ends the stream with this error after
	// FailAfterCount items — the case that becomes PARTIAL, and the reason
	// Stream carries a terminal error at all.
	FailAfter      error
	FailAfterCount int

	// AnnounceTotal makes the stream publish its denominator, which is what
	// lets Z-03 draw a determinate bar instead of an indeterminate scanner.
	// False models a provider that cannot say how many are coming.
	AnnounceTotal bool

	// Delay is slept between items, so a test can exercise cancellation
	// mid-stream. Zero streams as fast as the consumer reads.
	Delay time.Duration
}

// NewFake returns a Fake with a Steam-shaped capability set.
func NewFake(items ...provider.Item) *Fake {
	return &Fake{
		Ident: "fake",
		Name:  "Fake store",
		Caps: provider.Capabilities{
			Sync: true, Playtime: true, LastPlayed: true, OwnedSince: true,
			Credentials: []provider.CredentialField{
				{Key: "api_key", LabelKey: "cred.api_key.label", Secret: true},
				{Key: "account", LabelKey: "cred.account.label"},
			},
		},
		Items:         items,
		AnnounceTotal: true,
	}
}

// ID returns the fake's identity.
func (f *Fake) ID() provider.ID { return f.Ident }

// Display returns the fake's player-facing name.
func (f *Fake) Display() string { return f.Name }

// Capabilities returns the declared capability set.
func (f *Fake) Capabilities() provider.Capabilities { return f.Caps }

// Sync streams the scripted items.
//
// It honours cancellation at every step, including while sleeping, because a
// provider that only checks ctx between items would pass a fast test and hang
// the product on a slow network — and Z-03 lets the player press Esc at
// exactly that moment.
//
// Missing credentials are rejected before anything is streamed, using the
// provider's own declared field set rather than a hard-coded list, which is
// the same check Z-02's empty-submit state performs.
func (f *Fake) Sync(ctx context.Context, c provider.Credentials) (provider.Stream, error) {
	if f.FailBefore != nil {
		return nil, f.FailBefore
	}
	if missing := provider.Missing(f.Caps, c); len(missing) > 0 {
		return nil, fault.New(fault.KindUnauthorized, "fake.Sync",
			fault.WithSubject(f.Name), fault.WithMessage("fault.unauthorized"))
	}

	s := &stream{ch: make(chan provider.Item)}
	if f.AnnounceTotal {
		s.setTotal(len(f.Items))
	}

	go func() {
		defer close(s.ch)
		for i, it := range f.Items {
			if f.FailAfter != nil && i >= f.FailAfterCount {
				s.setErr(f.FailAfter)
				return
			}
			if f.Delay > 0 {
				select {
				case <-ctx.Done():
					s.setErr(cancelled(ctx))
					return
				case <-time.After(f.Delay):
				}
			}
			select {
			case <-ctx.Done():
				s.setErr(cancelled(ctx))
				return
			case s.ch <- it:
				s.advance()
			}
		}
		if f.FailAfter != nil && len(f.Items) >= f.FailAfterCount {
			s.setErr(f.FailAfter)
		}
	}()
	return s, nil
}

// cancelled turns a context's own error into a classified fault, so that a
// caller switching on fault.Kind sees KindCancelled rather than an
// unclassified error that would render as the fatal screen.
func cancelled(ctx context.Context) error {
	return fault.New(fault.KindCancelled, "fake.Sync", fault.WithCause(ctx.Err()))
}

// stream is the Fake's Stream implementation.
//
// Its progress snapshot is mutex-guarded rather than raced, because the
// contract on provider.Stream.Progress says it is safe to call from another
// goroutine at any time — and a fake that did not honour that would let a
// racy real implementation pass every test.
type stream struct {
	ch chan provider.Item

	mu   sync.Mutex
	prog provider.Progress
	err  error
}

func (s *stream) Items() <-chan provider.Item { return s.ch }

func (s *stream) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

func (s *stream) Progress() provider.Progress {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.prog
}

func (s *stream) setErr(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = err
}

func (s *stream) setTotal(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prog.Total, s.prog.TotalKnown = n, true
}

func (s *stream) advance() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prog.Seen++
	s.prog.LastAt = time.Now()
}

var (
	_ provider.Provider = (*Fake)(nil)
	_ provider.Syncer   = (*Fake)(nil)
	_ provider.Stream   = (*stream)(nil)
)
