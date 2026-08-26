package fault_test

import (
	"testing"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
)

// TestClassifyFollowsTheOfflineContract walks 07-offline-contract.md §5's
// decision tree, every branch, with no network involved — which is the point:
// the classifier is a pure function over a verdict the provider supplies, so
// it does not import net or net/http and every branch is reachable in a test.
func TestClassifyFollowsTheOfflineContract(t *testing.T) {
	cases := []struct {
		name string
		out  fault.Outcome
		want fault.Kind
	}{
		{"no route", fault.Outcome{Transport: fault.TransportNoRoute}, fault.KindOffline},
		{"DNS did not resolve", fault.Outcome{Transport: fault.TransportDNS}, fault.KindOffline},
		{"timeout", fault.Outcome{Transport: fault.TransportTimeout}, fault.KindUnreachable},
		{"connection reset mid-body", fault.Outcome{Transport: fault.TransportReset}, fault.KindUnreachable},
		{"401", fault.Outcome{Status: 401}, fault.KindUnauthorized},
		{"403", fault.Outcome{Status: 403}, fault.KindUnauthorized},
		{"404", fault.Outcome{Status: 404}, fault.KindNotFound},
		{"429", fault.Outcome{Status: 429}, fault.KindRateLimited},
		{"500", fault.Outcome{Status: 500}, fault.KindUnreachable},
		{"503", fault.Outcome{Status: 503}, fault.KindUnreachable},
		{"400 — our request, not their outage", fault.Outcome{Status: 400}, fault.KindMalformed},
		{"no response and no verdict", fault.Outcome{}, fault.KindInternal},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.out.Op = "steam.Sync"
			err := fault.Classify(c.out)
			if err == nil {
				t.Fatalf("Classify returned success for %s", c.name)
			}
			if got := fault.KindOf(err); got != c.want {
				t.Fatalf("Classify(%s) = %v, want %v", c.name, got, c.want)
			}
		})
	}

	if err := fault.Classify(fault.Outcome{Op: "steam.Sync", Status: 200}); err != nil {
		t.Fatalf("a 200 classified as %v", fault.KindOf(err))
	}
}

// TestTransportOutranksStatus: a rejected key is meaningless if the request
// never arrived, so the transport branch is consulted first.
func TestTransportOutranksStatus(t *testing.T) {
	err := fault.Classify(fault.Outcome{Op: "steam.Sync", Transport: fault.TransportNoRoute, Status: 401})
	if got := fault.KindOf(err); got != fault.KindOffline {
		t.Fatalf("got %v; a 401 recorded alongside a failed transport is not a credential verdict", got)
	}
}

// TestAnEmptySyncIsARefusal is the guard behind the ratified copy "Your
// library is unchanged — nothing was lost." A 200 carrying zero items is the
// private-profile case, and treating it as an empty result set would delete a
// 247-game library.
func TestAnEmptySyncIsARefusal(t *testing.T) {
	err := fault.ClassifySync(fault.Outcome{Op: "steam.Sync", Subject: "Steam", Status: 200, Items: 0})
	if err == nil {
		t.Fatal("a successful, empty library sync reported success; the upsert would then have emptied the library")
	}
	if got := fault.KindOf(err); got != fault.KindEmpty {
		t.Fatalf("got %v, want KindEmpty", got)
	}
	if got := fault.SubjectOf(err); got != "Steam" {
		t.Fatalf("the subject was lost (%q); the copy needs it as a substitution", got)
	}

	if err := fault.ClassifySync(fault.Outcome{Op: "steam.Sync", Status: 200, Items: 1}); err != nil {
		t.Fatalf("a one-item library classified as %v", fault.KindOf(err))
	}
}

// TestEmptinessIsOnlyARefusalForASync: a metadata lookup that finds nothing is
// a designed empty, and only a library fetch may conclude that the player must
// go and change a setting somewhere else.
func TestEmptinessIsOnlyARefusalForASync(t *testing.T) {
	if err := fault.Classify(fault.Outcome{Op: "meta.Lookup", Status: 200, Items: 0}); err != nil {
		t.Fatalf("a non-sync operation treated emptiness as a refusal: %v", fault.KindOf(err))
	}
}

// TestRateLimitCarriesTheProvidersWait.
func TestRateLimitCarriesTheProvidersWait(t *testing.T) {
	err := fault.Classify(fault.Outcome{Op: "steam.Sync", Status: 429, RetryAfter: 30 * time.Second})
	if got := fault.RetryAfterOf(err); got != 30*time.Second {
		t.Fatalf("RetryAfter = %v, want 30s", got)
	}
}

// TestClassifySuccessIsAnUntypedNil is the MINOR-3 repair from the review at
// 4484d9a, pinned.
//
// Classify used to return *Fault, which made the natural provider spelling —
// `return fault.Classify(o)` from a function returning error — hand back a
// non-nil error interface wrapping a typed nil on SUCCESS. KindOf then
// reported KindInternal, so a completely successful sync would have rendered
// the fatal screen.
//
// The check is deliberately the one a caller performs, `err != nil`, on a
// value that has passed through an error-returning function — because that is
// the exact conversion that created the trap, and a direct comparison against
// the concrete return type would not have caught it.
func TestClassifySuccessIsAnUntypedNil(t *testing.T) {
	viaFunc := func(o fault.Outcome) error { return fault.Classify(o) }
	viaSync := func(o fault.Outcome) error { return fault.ClassifySync(o) }

	ok := fault.Outcome{Op: "steam.Sync", Status: 200, Items: 12}
	for name, got := range map[string]error{
		"Classify":     viaFunc(ok),
		"ClassifySync": viaSync(ok),
	} {
		if got != nil {
			t.Fatalf("%s: a successful outcome produced a non-nil error (%v, kind %v).\n"+
				"If this is a typed nil, every caller that returns it as an error reports a\n"+
				"successful sync as KindInternal and renders Z-11.", name, got, fault.KindOf(got))
		}
		if k := fault.KindOf(got); k != fault.KindUnknown {
			t.Fatalf("%s: KindOf(success) = %v, want KindUnknown", name, k)
		}
	}

	// And the failing path still classifies through the same conversion.
	if k := fault.KindOf(viaFunc(fault.Outcome{Op: "steam.Sync", Status: 401})); k != fault.KindUnauthorized {
		t.Fatalf("through an error-returning function, a 401 classified as %v", k)
	}
}
