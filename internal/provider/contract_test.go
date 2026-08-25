package provider_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider/providertest"
	"github.com/JustCode-CruzAlex/Zerado/internal/status"
)

// TestManualIsNotASyncer is the decision, checked.
//
// The compiler can assert that a type implements an interface; it cannot
// assert that a type deliberately does not. ADR-0001 D1 rejects the
// alternative — physical implementing Sync and returning ErrManual — because
// every caller would then have to know which providers lie about their own
// interface. This test is the only place that rejection is enforceable.
func TestManualIsNotASyncer(t *testing.T) {
	m := providertest.NewManual()

	if _, isSyncer := any(m).(provider.Syncer); isSyncer {
		t.Fatal("the hand-entry provider implements Syncer; it is not a store with a hole where Sync should be")
	}
	if _, isProvider := any(m).(provider.Provider); !isProvider {
		t.Fatal("the hand-entry provider does not implement Provider; a physical copy would be a second-class row")
	}
	if _, isEnterer := any(m).(provider.Enterer); !isEnterer {
		t.Fatal("the hand-entry provider does not implement Enterer")
	}
}

// TestCapabilitiesAgreeWithTheInterfaces: Capabilities duplicates two facts
// the type system already carries, because screens need them before they have
// a reason to type assert. Duplicated facts drift, so this is the assertion
// that they have not.
func TestCapabilitiesAgreeWithTheInterfaces(t *testing.T) {
	for _, p := range []provider.Provider{providertest.NewManual(), providertest.NewFake()} {
		if problems := provider.Check(p); len(problems) > 0 {
			t.Fatalf("%s: %v", p.ID(), problems)
		}
	}
}

// TestACartridgeCannotClaimTelemetry: a hand-entry provider that reported
// playtime would feed status.Derive a typed number as evidence. Check refuses
// the declaration rather than leaving it to be noticed on a screen.
func TestACartridgeCannotClaimTelemetry(t *testing.T) {
	m := providertest.NewManual()
	liar := &lyingManual{Manual: m}
	if problems := provider.Check(liar); len(problems) == 0 {
		t.Fatal("a hand-entry-only provider claiming Playtime passed Check")
	}
}

type lyingManual struct{ *providertest.Manual }

func (l *lyingManual) Capabilities() provider.Capabilities {
	c := l.Manual.Capabilities()
	c.Playtime = true
	return c
}

// TestHandEntryProducesTheSameValueASyncDoes is the seam's real test. If
// entering a cartridge needed a different type from a Steam row, the contract
// would be Steam-shaped and physical would be bolted on.
func TestHandEntryProducesTheSameValueASyncDoes(t *testing.T) {
	m := providertest.NewManual()
	it, err := m.Compose(provider.Entry{
		"title":       "  Shadow of the Colossus  ",
		"platform":    "PS2",
		"owned_since": "2004-11-15",
	})
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
	if it.Title != "Shadow of the Colossus" {
		t.Fatalf("Title = %q; the form's whitespace reached the library", it.Title)
	}
	if it.ProviderRef == "" {
		t.Fatal("Compose did not mint a provider reference; the player would have to supply an app ID a cartridge does not have")
	}
	if it.PlaytimeMinutes != nil {
		t.Fatal("Compose set a playtime; a cartridge has no telemetry and a zero would be treated as evidence")
	}
	if it.OwnedSince == nil {
		t.Fatal("OwnedSince was dropped; it is the one optional capability physical claims")
	}

	// The same Item type a sync produces, and the derivation over it agrees
	// with the provider's capabilities rather than with its identity.
	if got := status.Derive(0, m.Capabilities().Playtime); got != status.NotStarted {
		t.Fatalf("a hand-entered game derived to %v", got)
	}
}

// TestTwoCartridgesNeverCollide: the ref is minted fresh, so a duplicate is a
// note rather than a block (Z-08 §8.3).
func TestTwoCartridgesNeverCollide(t *testing.T) {
	m := providertest.NewManual()
	e := provider.Entry{"title": "Shadow of the Colossus", "platform": "PS2"}
	a, _ := m.Compose(e)
	b, _ := m.Compose(e)
	if a.ProviderRef == b.ProviderRef {
		t.Fatal("two hand entries share a provider reference; the uniqueness constraint would block a legitimate second copy")
	}
}

// TestComposeRefusesAnIncompleteForm, with the offending field named so Z-08
// can put the message under the field it belongs to.
func TestComposeRefusesAnIncompleteForm(t *testing.T) {
	m := providertest.NewManual()
	for _, e := range []provider.Entry{
		{"platform": "PS2"},
		{"title": "Ico"},
		{"title": "   ", "platform": "PS2"},
	} {
		_, err := m.Compose(e)
		if err == nil {
			t.Fatalf("Compose accepted %v", e)
		}
		if !fault.Is(err, fault.KindMalformed) {
			t.Fatalf("Compose(%v) = %v, want KindMalformed", e, fault.KindOf(err))
		}
		if !fault.MessageKeyOf(err).Valid() {
			t.Fatal("the refusal carries no catalogue key; the field would show a blank error")
		}
	}
}

// TestAStoreNeedsNoScreenChanges: Z-02 renders from the declared credential
// fields, so adding a provider is implement, declare, register.
func TestAStoreNeedsNoScreenChanges(t *testing.T) {
	f := providertest.NewFake()
	fields := f.Capabilities().Credentials
	if len(fields) != 2 {
		t.Fatalf("expected the fake to declare two fields, got %d", len(fields))
	}
	if !fields[0].Secret {
		t.Fatal("the key field is not marked secret; it would be written to library.db")
	}
	if fields[1].Secret {
		t.Fatal("the account field is marked secret; an identifier is not a credential")
	}
	if got := provider.Secrets(f.Capabilities()); len(got) != 1 || got[0] != "api_key" {
		t.Fatalf("Secrets = %v", got)
	}
	if got := provider.Missing(f.Capabilities(), provider.Credentials{"api_key": "  "}); len(got) != 2 {
		t.Fatalf("Missing = %v; whitespace-only must count as absent and the second field is untouched", got)
	}
	// A pasted credential is a plausible source of non-ASCII whitespace, in a
	// product that is internationalised from the first line. A key that is
	// only a NO-BREAK SPACE must be caught here, inline under the field, and
	// not by a network round trip that comes back as a rejected key.
	for _, blank := range []string{"\u00a0", "\u3000", "\u2003\u00a0"} {
		if got := provider.Missing(f.Capabilities(), provider.Credentials{"api_key": blank, "account": "a"}); len(got) != 1 {
			t.Fatalf("Missing with %q = %v; Unicode whitespace read as a real value", blank, got)
		}
	}
	// A shelf declares none, so Z-02 is never reachable for it.
	if len(providertest.NewManual().Capabilities().Credentials) != 0 {
		t.Fatal("the hand-entry provider declares credential fields; Z-02 would open on a screen that can do nothing")
	}
}

// TestSyncStreams: a cancel mid-sync leaves a valid partial library, and the
// stream reports why it stopped after it has closed.
func TestSyncStreams(t *testing.T) {
	f := providertest.NewFake(items(5)...)
	f.Delay = 5 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s, err := f.Sync(ctx, provider.Credentials{"api_key": "k", "account": "a"})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	var got int
	for range s.Items() {
		got++
		if got == 2 {
			cancel()
		}
	}
	if got == 0 {
		t.Fatal("a cancelled sync delivered nothing; what arrived before the cancel must be kept")
	}
	if got == 5 {
		t.Fatal("the cancel was ignored; the goroutine budget is zero leaks and Esc must abort in-flight I/O")
	}
	if !fault.Is(s.Err(), fault.KindCancelled) {
		t.Fatalf("Err after a cancel = %v, want KindCancelled — Z-03's CANCELLED state is not an error", fault.KindOf(s.Err()))
	}
}

// TestPartialIsExpressible is why Stream exists rather than a bare channel:
// items arrived, then the connection broke, and that is one of Z-03's four
// terminal states.
func TestPartialIsExpressible(t *testing.T) {
	f := providertest.NewFake(items(5)...)
	f.FailAfterCount = 3
	f.FailAfter = fault.New(fault.KindUnreachable, "fake.Sync", fault.WithSubject("Fake store"))

	s, err := f.Sync(context.Background(), provider.Credentials{"api_key": "k", "account": "a"})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	var got int
	for range s.Items() {
		got++
	}
	if got != 3 {
		t.Fatalf("delivered %d items before the break, want 3", got)
	}
	if !fault.Is(s.Err(), fault.KindUnreachable) {
		t.Fatalf("Err = %v; a failure after items have arrived has nowhere to go without Stream.Err", fault.KindOf(s.Err()))
	}
}

// TestProgressIsReadableWhileTheSyncRuns is the property Z-03 needs and a bare
// channel cannot supply: the denominator, so a bar can be drawn instead of an
// apology, read from the render goroutine while the provider writes.
//
// Run with -race, this is also the guard on the contract's "safe to call from
// another goroutine at any time".
func TestProgressIsReadableWhileTheSyncRuns(t *testing.T) {
	f := providertest.NewFake(items(20)...)
	f.Delay = time.Millisecond

	s, err := f.Sync(context.Background(), provider.Credentials{"api_key": "k", "account": "a"})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if p := s.Progress(); !p.Determinate() {
		t.Fatal("the denominator is unknown before the first item; Z-03 could not draw a determinate bar")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for range s.Items() {
		}
	}()
	for {
		p := s.Progress() // concurrent with the producer, on purpose
		if p.Seen == 20 {
			break
		}
		select {
		case <-done:
			if s.Progress().Seen != 20 {
				t.Errorf("Seen = %d after the stream closed, want 20", s.Progress().Seen)
			}
			return
		default:
		}
	}
	<-done
}

// TestAnUnknownDenominatorDrawsTheScanner: a provider that cannot say how many
// are coming must not be forced to invent a total.
func TestAnUnknownDenominatorDrawsTheScanner(t *testing.T) {
	f := providertest.NewFake(items(3)...)
	f.AnnounceTotal = false
	s, _ := f.Sync(context.Background(), provider.Credentials{"api_key": "k", "account": "a"})
	if s.Progress().Determinate() {
		t.Fatal("Determinate is true with no announced total")
	}
	for range s.Items() {
	}
}

// TestAZeroTotalIsNotUnknown: zero is a real total — it is the private-profile
// case — and a screen reading it as "unknown" would scan forever for a sync
// that has truthfully finished.
func TestAZeroTotalIsNotUnknown(t *testing.T) {
	f := providertest.NewFake()
	s, _ := f.Sync(context.Background(), provider.Credentials{"api_key": "k", "account": "a"})
	for range s.Items() {
	}
	p := s.Progress()
	if !p.TotalKnown {
		t.Fatal("an announced total of zero reported itself unknown")
	}
	if p.Determinate() {
		t.Fatal("a total of zero drew a determinate bar; there is nothing to fill")
	}
}

// TestStalledIsNotWaiting: a sync that has delivered nothing yet is waiting,
// which already has its own component, and calling it stalled would flip the
// screen for no reason the player could see.
func TestStalledIsNotWaiting(t *testing.T) {
	base := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	waiting := provider.Progress{Seen: 0}
	if waiting.Stalled(base.Add(time.Minute), 10*time.Second) {
		t.Fatal("a sync that has not yet delivered anything reported itself stalled")
	}
	moving := provider.Progress{Seen: 5, LastAt: base}
	if !moving.Stalled(base.Add(11*time.Second), 10*time.Second) {
		t.Fatal("eleven seconds without an item was not reported as stalled")
	}
	if moving.Stalled(base.Add(3*time.Second), 10*time.Second) {
		t.Fatal("three seconds without an item was reported as stalled")
	}
}

// TestARejectedKeyFailsBeforeAnythingStreams.
func TestARejectedKeyFailsBeforeAnythingStreams(t *testing.T) {
	f := providertest.NewFake(items(3)...)
	_, err := f.Sync(context.Background(), provider.Credentials{})
	if err == nil {
		t.Fatal("Sync accepted empty credentials")
	}
	if !fault.Is(err, fault.KindUnauthorized) {
		t.Fatalf("got %v, want KindUnauthorized", fault.KindOf(err))
	}
}

// TestRegistryRoutesWithoutSwitchingOnIdentity: nothing has to know that
// physical exists in order to skip it.
func TestRegistryRoutesWithoutSwitchingOnIdentity(t *testing.T) {
	r := provider.NewRegistry(providertest.NewFake(), providertest.NewManual())
	if got := len(r.Syncers()); got != 1 {
		t.Fatalf("Syncers = %d, want 1", got)
	}
	if got := len(r.Enterers()); got != 1 {
		t.Fatalf("Enterers = %d, want 1", got)
	}
	if _, ok := r.Get("physical"); !ok {
		t.Fatal("the registry cannot resolve a stored provider id back to a provider")
	}
}

// TestDuplicateRegistrationIsVisible is the MINOR-2 repair from the review at
// c4c8d95.
//
// Duplicates used to be a method that ignored its receiver, so the natural
// call returned nil whatever the registry held, and the only call that did
// anything re-passed the slice it had already handed NewRegistry. Both halves
// are now asserted: the package-level check before construction, and what a
// built registry actually swallowed.
func TestDuplicateRegistrationIsVisible(t *testing.T) {
	a, b := providertest.NewFake(), providertest.NewFake()

	if dup := provider.Duplicates(a, b); len(dup) != 1 || dup[0] != a.ID() {
		t.Fatalf("Duplicates = %v, want [%s]", dup, a.ID())
	}
	if dup := provider.Duplicates(a, providertest.NewManual()); len(dup) != 0 {
		t.Fatalf("Duplicates reported %v for two distinct providers", dup)
	}

	// The question a method on Registry has to be able to answer, and could
	// not: what did THIS registry swallow?
	r := provider.NewRegistry(a, b, providertest.NewManual())
	got := r.Collisions()
	if len(got) != 1 || got[0] != a.ID() {
		t.Fatalf("Collisions = %v, want [%s]", got, a.ID())
	}
	if len(r.All()) != 2 {
		t.Fatalf("the registry holds %d providers, want 2 — the duplicate replaced rather than added", len(r.All()))
	}
	if len(provider.NewRegistry(a, providertest.NewManual()).Collisions()) != 0 {
		t.Fatal("a clean registry reported a collision")
	}
}

func items(n int) []provider.Item {
	out := make([]provider.Item, n)
	for i := range out {
		m := i * 10
		out[i] = provider.Item{
			ProviderRef:     string(rune('a' + i)),
			Title:           "Game " + string(rune('A'+i)),
			Platform:        "PC",
			PlaytimeMinutes: &m,
		}
	}
	return out
}

// TestEveryProgressFieldIsActuallyWritten is the MINOR-1 repair from the
// review at c4c8d95, generalised.
//
// Progress carried a Batches field that no interface method allowed anyone to
// write, so it could only ever read zero — and the screen it existed for would
// have rendered "The 0 that arrived are in your library", which is the exact
// sentence PARTIAL exists to make true.
//
// A field on a snapshot that nothing can write is not a contract, it is a
// promise the type cannot keep. This walks Progress by reflection after a real
// stream has run and fails on any field still at its zero value, so the next
// unwritable field is caught at the moment it is added rather than by the
// screen that trusted it.
func TestEveryProgressFieldIsActuallyWritten(t *testing.T) {
	f := providertest.NewFake(items(3)...)
	s, err := f.Sync(context.Background(), provider.Credentials{"api_key": "k", "account": "a"})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	for range s.Items() {
	}

	v := reflect.ValueOf(s.Progress())
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		if v.Field(i).IsZero() {
			t.Errorf("Progress.%s is still its zero value after a complete sync.\n"+
				"Either a conforming implementation cannot write it — in which case it is a\n"+
				"promise the type cannot keep and belongs somewhere a writer exists — or the\n"+
				"fake is not exercising it, which makes every screen that reads it untested.", name)
		}
	}
}
