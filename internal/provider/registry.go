package provider

import "sort"

// Registry is the set of providers compiled into this binary.
//
// ADR-0001 D1 rejects a plugin system for stated reasons — four stores are
// known, they ship in the binary, Go plugins do not cross-compile, and a
// subprocess protocol is a network boundary in disguise. So the registry is
// a map, and registration happens once at start-up.
//
// It is the one place above this package where looking a provider up by [ID]
// is legitimate: a stored row carries a provider_id and something has to turn
// that back into a Provider. Every *other* use of an ID above the seam is the
// switch this design exists to remove.
type Registry struct {
	byID  map[ID]Provider
	order []ID
}

// NewRegistry builds a registry from providers, in the order given.
//
// The order is display order — Z-02's picker and Z-09's groups both render it
// — so it is the caller's to choose and is preserved rather than sorted.
// A duplicate ID is a programming error and the later registration wins, which
// is reported by [Registry.Duplicates] rather than silently tolerated.
func NewRegistry(ps ...Provider) *Registry {
	r := &Registry{byID: make(map[ID]Provider, len(ps))}
	for _, p := range ps {
		id := p.ID()
		if _, seen := r.byID[id]; !seen {
			r.order = append(r.order, id)
		}
		r.byID[id] = p
	}
	return r
}

// Get returns the provider with this ID.
func (r *Registry) Get(id ID) (Provider, bool) {
	p, ok := r.byID[id]
	return p, ok
}

// All returns every registered provider in registration order.
func (r *Registry) All() []Provider {
	out := make([]Provider, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.byID[id])
	}
	return out
}

// Syncers returns the providers that can fetch a library.
//
// This is what "sync everything" iterates, and it is why no caller needs to
// know that physical exists in order to skip it.
func (r *Registry) Syncers() []Syncer {
	var out []Syncer
	for _, p := range r.All() {
		if s, ok := p.(Syncer); ok {
			out = append(out, s)
		}
	}
	return out
}

// Enterers returns the providers that accept hand-entered items.
//
// Z-01's "Add a game by hand" door is enabled when this is non-empty, and
// Z-08's source is its single element in Phase 1 — a picker the day it has
// two, with no other change.
func (r *Registry) Enterers() []Enterer {
	var out []Enterer
	for _, p := range r.All() {
		if e, ok := p.(Enterer); ok {
			out = append(out, e)
		}
	}
	return out
}

// Duplicates returns IDs that were registered more than once, sorted.
//
// Registration is a start-up path with no screen behind it, so this is a
// health check rather than an error return: a test asserts it is empty, and a
// developer who registers Steam twice finds out in CI rather than by wondering
// why their credential fields changed.
func (r *Registry) Duplicates(ps ...Provider) []ID {
	seen := map[ID]int{}
	for _, p := range ps {
		seen[p.ID()]++
	}
	var out []ID
	for id, n := range seen {
		if n > 1 {
			out = append(out, id)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// Check reports the ways a provider's declaration disagrees with the
// interfaces it implements.
//
// [Capabilities] duplicates two facts the type system already carries — Sync
// and Manual — because screens need them before they have a reason to type
// assert, and because the offline classification is a property of the provider
// rather than of a Go assertion. Duplicated facts drift, so this is the
// assertion that they have not, and every provider is expected to be run
// through it by a test.
//
// The returned strings are developer-facing and are not rendered anywhere.
func Check(p Provider) []string {
	var problems []string
	c := p.Capabilities()

	_, isSyncer := p.(Syncer)
	if c.Sync != isSyncer {
		problems = append(problems, "Capabilities.Sync disagrees with whether the provider implements Syncer")
	}

	_, isEnterer := p.(Enterer)
	if c.Manual != isEnterer {
		problems = append(problems, "Capabilities.Manual disagrees with whether the provider implements Enterer")
	}

	// A provider that cannot fetch has nothing to authenticate with. This is
	// not pedantry: Z-02 exists only for providers with credential fields, so
	// a non-syncing provider that declares one would put a screen in front of
	// a player that could never do anything.
	if !c.Sync && len(c.Credentials) > 0 {
		problems = append(problems, "a provider that cannot Sync declares credential fields")
	}

	// A hand-entry provider that claims it can report playtime is claiming a
	// telemetry source it does not have; the derivation would then treat a
	// typed number as evidence.
	if c.Manual && !c.Sync && c.Playtime {
		problems = append(problems, "a hand-entry-only provider claims Playtime, which nothing could observe")
	}

	if p.ID() == "" {
		problems = append(problems, "provider ID is empty")
	}
	if p.Display() == "" {
		problems = append(problems, "provider Display is empty")
	}
	return problems
}
