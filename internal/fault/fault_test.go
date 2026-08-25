package fault_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
)

// TestTaxonomyIsTotal fails the moment somebody adds a Kind without deciding
// what it looks like on screen. Every member needs a name, a treatment and a
// catalogue key, and no two may share a name.
func TestTaxonomyIsTotal(t *testing.T) {
	seen := map[string]fault.Kind{}
	for _, k := range fault.Kinds() {
		if !k.Valid() {
			t.Errorf("%v is in Kinds() but reports itself invalid", k)
		}
		name := k.String()
		if name == "" || name == "unknown" {
			t.Errorf("kind %d has no stable machine name", k)
		}
		if prev, dup := seen[name]; dup {
			t.Errorf("kinds %v and %v share the machine name %q; the CLI cannot distinguish them", prev, k, name)
		}
		seen[name] = k
		if !k.MessageKey().Valid() {
			t.Errorf("%v has no catalogue key; it would render as a blank message", k)
		}
	}
	if len(fault.Kinds()) != len(seen) {
		t.Fatalf("Kinds() has %d entries but %d distinct names", len(fault.Kinds()), len(seen))
	}
}

// TestAPrivateProfileIsNotANetworkError is the case this whole package exists
// for. Four failures that a bare error string would flatten into one must be
// four distinguishable kinds with three different screen treatments.
func TestAPrivateProfileIsNotANetworkError(t *testing.T) {
	private := fault.New(fault.KindEmpty, "steam.Sync", fault.WithSubject("Steam"))
	offline := fault.New(fault.KindOffline, "steam.Sync", fault.WithSubject("Steam"))
	unreachable := fault.New(fault.KindUnreachable, "steam.Sync", fault.WithSubject("Steam"))
	rejected := fault.New(fault.KindUnauthorized, "steam.Sync", fault.WithSubject("Steam"))

	all := []*fault.Fault{private, offline, unreachable, rejected}
	for i, a := range all {
		for j, b := range all {
			if i != j && a.Kind == b.Kind {
				t.Fatalf("two of the four Z-03 failures share a kind: %v", a.Kind)
			}
		}
	}

	if fault.Is(private, fault.KindOffline) || fault.Is(private, fault.KindUnreachable) {
		t.Fatal("a private profile classified as a network failure")
	}
	if private.Kind.Treatment() != fault.TreatmentRefusal {
		t.Fatalf("a private profile renders as %v; it is an action the player took and is owed a refusal", private.Kind.Treatment())
	}
	if offline.Kind.Treatment() != fault.TreatmentBanner {
		t.Fatalf("being offline renders as %v; it is ambient and belongs in the banner", offline.Kind.Treatment())
	}
	if private.Kind.Retryable() {
		t.Fatal("a private profile is retryable; pressing r changes nothing until the player changes a setting elsewhere")
	}
	if !offline.Kind.Retryable() {
		t.Fatal("being offline is not retryable; the banner would offer no way out")
	}
}

// TestNotFoundIsADesignedEmpty guards 06 §3.1: no metadata is a composition,
// not an error banner.
func TestNotFoundIsADesignedEmpty(t *testing.T) {
	if got := fault.KindNotFound.Treatment(); got != fault.TreatmentDesignedEmpty {
		t.Fatalf("KindNotFound renders as %v; a cartridge no source has heard of is not a failure", got)
	}
}

// TestCancelledIsNotAnError guards Z-03's CANCELLED state: no cue, no red, no
// apology.
func TestCancelledIsNotAnError(t *testing.T) {
	if got := fault.KindCancelled.Treatment(); got != fault.TreatmentNone {
		t.Fatalf("KindCancelled renders as %v; the player did it on purpose", got)
	}
}

// TestKindSurvivesWrapping is the reason callers use fault.Is rather than a
// type assertion: a package adding context must not reclassify a failure.
func TestKindSurvivesWrapping(t *testing.T) {
	base := fault.New(fault.KindUnauthorized, "steam.Sync")
	wrapped := fmt.Errorf("while syncing the library: %w", base)
	if !fault.Is(wrapped, fault.KindUnauthorized) {
		t.Fatal("wrapping lost the kind")
	}
	if got := fault.KindOf(wrapped); got != fault.KindUnauthorized {
		t.Fatalf("KindOf(wrapped) = %v", got)
	}
}

// TestUnclassifiedIsInternal: an error that reached a screen without ever
// being classified routes to the fatal screen rather than being guessed at.
func TestUnclassifiedIsInternal(t *testing.T) {
	if got := fault.KindOf(errors.New("boom")); got != fault.KindInternal {
		t.Fatalf("an unclassified error reported %v; guessing somebody else's failure is worse than saying it is ours", got)
	}
	if got := fault.KindOf(nil); got != fault.KindUnknown {
		t.Fatalf("KindOf(nil) = %v, want KindUnknown", got)
	}
	if !fault.MessageKeyOf(errors.New("boom")).Valid() {
		t.Fatal("an unclassified error has no catalogue key; it would render blank")
	}
}

// TestErrorDoesNotLeakTheCause is a credential-disclosure guard, not a style
// preference. A provider's transport error stringifies a URL, and a Steam URL
// carries the API key in its query string — straight into the one screen a
// frustrated player screenshots.
func TestErrorDoesNotLeakTheCause(t *testing.T) {
	leak := errors.New(`Get "https://api.steampowered.com/x?key=SECRET-KEY-9F2": dial tcp: no route`)
	f := fault.New(fault.KindOffline, "steam.Sync", fault.WithSubject("Steam"), fault.WithCause(leak))

	if got := f.Error(); contains(got, "SECRET-KEY-9F2") {
		t.Fatalf("Fault.Error() leaked the credential: %q", got)
	}
	if !errors.Is(f, leak) {
		t.Fatal("the cause is unreachable by errors.Is; a log sink could not get it")
	}
}

func contains(hay, needle string) bool {
	for i := 0; i+len(needle) <= len(hay); i++ {
		if hay[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// TestRetryAfterTravels: a rate limit is a distinct kind precisely because it
// carries a wait, and the wait has to survive to the CLI envelope.
func TestRetryAfterTravels(t *testing.T) {
	f := fault.New(fault.KindRateLimited, "steam.Sync", fault.WithRetryAfter(90*time.Second))
	if got := fault.RetryAfterOf(f); got != 90*time.Second {
		t.Fatalf("RetryAfterOf = %v, want 90s", got)
	}
	if got := fault.RetryAfterOf(fault.New(fault.KindOffline, "x")); got != 0 {
		t.Fatalf("an unstated wait reported %v; zero must mean 'not stated', never 'retry now'", got)
	}
}

// TestNewRefusesAnInvalidKind: an unclassified failure is a hole in the
// taxonomy, and a panic at the construction site is cheaper than a banner
// nobody designed.
func TestNewRefusesAnInvalidKind(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("fault.New accepted KindUnknown")
		}
	}()
	_ = fault.New(fault.KindUnknown, "somewhere")
}

// TestProviderOverridesCopy: Steam's private-profile refusal names a specific
// Steam privacy setting, which no provider-agnostic entry could.
func TestProviderOverridesCopy(t *testing.T) {
	f := fault.New(fault.KindEmpty, "steam.Sync", fault.WithMessage("fault.steam.profile_private"))
	if got := fault.MessageKeyOf(f); got != "fault.steam.profile_private" {
		t.Fatalf("MessageKeyOf = %q; a provider must be able to supply better copy for its own case", got)
	}
	if fault.MessageKeyOf(fault.New(fault.KindEmpty, "gog.Sync")) != fault.KindEmpty.MessageKey() {
		t.Fatal("a provider that supplies no copy must still render a sentence")
	}
}

// TestKindsCoversTheWholeRange is the assertion that MAJOR-1 of the review at
// c4c8d95 falsified.
//
// Kinds() used to be a hand-maintained slice while Valid() was derived from
// the iota range, so the two could disagree — and a Kind inserted mid-block
// was valid, constructible, and absent from every totality test. Kinds() is
// now derived from the same range, and this asserts the two agree so that a
// future contributor who reintroduces a literal slice fails here.
func TestKindsCoversTheWholeRange(t *testing.T) {
	listed := map[fault.Kind]bool{}
	for _, k := range fault.Kinds() {
		listed[k] = true
	}
	for k := fault.KindOffline; k <= fault.KindInternal; k++ {
		if !listed[k] {
			t.Errorf("kind %d is valid and constructible but is not in Kinds(); "+
				"every totality test iterates Kinds() and would not see it", k)
		}
	}
	if len(fault.Kinds()) != int(fault.KindInternal) {
		t.Fatalf("Kinds() has %d entries but the taxonomy runs to %d",
			len(fault.Kinds()), int(fault.KindInternal))
	}
	if listed[fault.KindUnknown] {
		t.Fatal("KindUnknown is in Kinds(); it is never a member of the taxonomy")
	}
}

// TestNoKindFallsThroughToTheDefaults is the other half of the same repair.
//
// String, Treatment and MessageKey are switches with defaults, and a Kind that
// nobody added a case for inherits all three silently: the machine name
// "unknown" on a surface the CLI documents as stable API, the fatal treatment,
// and generic internal copy. Only KindInternal is entitled to those, so
// everything else claiming one is an unfinished addition.
func TestNoKindFallsThroughToTheDefaults(t *testing.T) {
	for _, k := range fault.Kinds() {
		if k == fault.KindInternal {
			continue // the defaults are its real answers
		}
		if got := k.String(); got == "unknown" {
			t.Errorf("kind %d has no case in String(); it would appear in the CLI's JSON as %q", k, got)
		}
		if got := k.MessageKey(); got == fault.KindInternal.MessageKey() {
			t.Errorf("kind %d has no case in MessageKey(); it would render the internal-failure copy", k)
		}
		if got := k.Treatment(); got == fault.TreatmentFatal {
			t.Errorf("kind %d has no case in Treatment(); it would stop the program", k)
		}
	}
}

// TestMarshallingAFaultDropsTheCause makes the redaction a property of the
// type rather than of every call site.
//
// Fault.Error and the CLI envelope both redact and both are tested, but
// neither protects a caller that has not been written yet — a log sink, a
// future debug verb, anything that reaches for json.Marshal. An *url.Error has
// an exported URL field, and a Steam URL carries the player's API key in its
// query string.
func TestMarshallingAFaultDropsTheCause(t *testing.T) {
	leak := &leakyError{URL: "https://api.steampowered.com/x?key=SECRET-KEY-9F2"}
	f := fault.New(fault.KindRateLimited, "steam.Sync",
		fault.WithSubject("Steam"), fault.WithCause(leak), fault.WithRetryAfter(30*time.Second))

	b, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, forbidden := range []string{"SECRET-KEY-9F2", "api.steampowered.com", "Cause", "cause"} {
		if contains(string(b), forbidden) {
			t.Fatalf("marshalling a Fault leaked %q: %s", forbidden, b)
		}
	}

	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["kind"] != "rate_limited" {
		t.Fatalf("kind = %v", got["kind"])
	}
	if got["retry_after_seconds"] != float64(30) {
		t.Fatalf("retry_after_seconds = %v; the one field a scripted caller acts on was dropped", got["retry_after_seconds"])
	}
	if _, err := json.Marshal((*fault.Fault)(nil)); err != nil {
		t.Fatalf("marshalling a nil Fault: %v", err)
	}
}

// leakyError mimics *url.Error: an exported field carrying a credential-bearing
// URL, which is what makes a naive json.Marshal a disclosure path.
type leakyError struct {
	URL string
}

func (e *leakyError) Error() string { return "Get " + e.URL + ": dial tcp: no route" }
