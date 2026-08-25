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
			f := fault.Classify(c.out)
			if f == nil {
				t.Fatalf("Classify returned success for %s", c.name)
			}
			if f.Kind != c.want {
				t.Fatalf("Classify(%s) = %v, want %v", c.name, f.Kind, c.want)
			}
		})
	}

	if f := fault.Classify(fault.Outcome{Op: "steam.Sync", Status: 200}); f != nil {
		t.Fatalf("a 200 classified as %v", f.Kind)
	}
}

// TestTransportOutranksStatus: a rejected key is meaningless if the request
// never arrived, so the transport branch is consulted first.
func TestTransportOutranksStatus(t *testing.T) {
	f := fault.Classify(fault.Outcome{Op: "steam.Sync", Transport: fault.TransportNoRoute, Status: 401})
	if f.Kind != fault.KindOffline {
		t.Fatalf("got %v; a 401 recorded alongside a failed transport is not a credential verdict", f.Kind)
	}
}

// TestAnEmptySyncIsARefusal is the guard behind the ratified copy "Your
// library is unchanged — nothing was lost." A 200 carrying zero items is the
// private-profile case, and treating it as an empty result set would delete a
// 247-game library.
func TestAnEmptySyncIsARefusal(t *testing.T) {
	f := fault.ClassifySync(fault.Outcome{Op: "steam.Sync", Subject: "Steam", Status: 200, Items: 0})
	if f == nil {
		t.Fatal("a successful, empty library sync reported success; the upsert would then have emptied the library")
	}
	if f.Kind != fault.KindEmpty {
		t.Fatalf("got %v, want KindEmpty", f.Kind)
	}
	if f.Subject != "Steam" {
		t.Fatalf("the subject was lost; the copy needs it as a substitution")
	}

	if f := fault.ClassifySync(fault.Outcome{Op: "steam.Sync", Status: 200, Items: 1}); f != nil {
		t.Fatalf("a one-item library classified as %v", f.Kind)
	}
}

// TestEmptinessIsOnlyARefusalForASync: a metadata lookup that finds nothing is
// a designed empty, and only a library fetch may conclude that the player must
// go and change a setting somewhere else.
func TestEmptinessIsOnlyARefusalForASync(t *testing.T) {
	if f := fault.Classify(fault.Outcome{Op: "meta.Lookup", Status: 200, Items: 0}); f != nil {
		t.Fatalf("a non-sync operation treated emptiness as a refusal: %v", f.Kind)
	}
}

// TestRateLimitCarriesTheProvidersWait.
func TestRateLimitCarriesTheProvidersWait(t *testing.T) {
	f := fault.Classify(fault.Outcome{Op: "steam.Sync", Status: 429, RetryAfter: 30 * time.Second})
	if f.RetryAfter != 30*time.Second {
		t.Fatalf("RetryAfter = %v, want 30s", f.RetryAfter)
	}
}
