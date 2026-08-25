// Package storetest holds an in-memory store, so every screen and every
// command can be tested with no database and no network.
//
// It is a test double, not a reference implementation: it has no SQL, no
// migrations, no WAL and no file. What it does share with a real store is the
// contract, and it enforces the parts of that contract that are rules rather
// than mechanics — the tombstone guard, the sync-never-writes-a-manual-status
// invariant, and the requirement that the summary counts sum.
//
// A fake that were merely permissive would let a real store violate those and
// still pass every test written against the fake, which is how a contract
// quietly becomes a suggestion.
package storetest

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/aged"
	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
	"github.com/JustCode-CruzAlex/Zerado/internal/library"
	"github.com/JustCode-CruzAlex/Zerado/internal/metadata"
	"github.com/JustCode-CruzAlex/Zerado/internal/pricing"
	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
	"github.com/JustCode-CruzAlex/Zerado/internal/status"
	"github.com/JustCode-CruzAlex/Zerado/internal/store"
)

// Fake is an in-memory store.
//
// Playtime capability is supplied per provider through [Fake.Playtime] rather
// than read from a registry, because a store must not depend on the provider
// seam to answer a query — and because a test wants to flip that capability
// without building a provider.
type Fake struct {
	mu sync.Mutex

	games    map[library.GameID]library.Game
	nextID   library.GameID
	runs     map[store.RunID]store.SyncRun
	nextRun  store.RunID
	conns    map[provider.ID]store.Connection
	settings map[string]string
	meta     map[library.GameID]aged.Value[metadata.Metadata]
	quotes   map[quoteKey]aged.Value[pricing.Quote]

	// Playtime reports whether a provider can report playtime, which the
	// effective-status derivation needs. Absent means false, which is the
	// honest default for a source nothing is known about.
	Playtime map[provider.ID]bool

	// Now supplies the clock.
	Now func() time.Time
}

type quoteKey struct {
	id  library.GameID
	cur pricing.Currency
}

// New returns an empty Fake with a fixed clock.
func New() *Fake {
	base := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	return &Fake{
		games:    map[library.GameID]library.Game{},
		runs:     map[store.RunID]store.SyncRun{},
		conns:    map[provider.ID]store.Connection{},
		settings: map[string]string{},
		meta:     map[library.GameID]aged.Value[metadata.Metadata]{},
		quotes:   map[quoteKey]aged.Value[pricing.Quote]{},
		Playtime: map[provider.ID]bool{},
		Now:      func() time.Time { return base },
	}
}

func (f *Fake) matches(g library.Game, q library.Query) bool {
	if g.MergedInto != nil {
		return false
	}
	switch q.Presence {
	case library.PresentOnly:
		if g.Absent() {
			return false
		}
	case library.AbsentOnly:
		if !g.Absent() {
			return false
		}
	}
	if q.Search != "" && !strings.Contains(fold(g.Title), fold(q.Search)) {
		return false
	}
	if len(q.Sources) > 0 && !containsID(q.Sources, g.Provider) {
		return false
	}
	if len(q.States) > 0 {
		s := g.Status(f.Playtime[g.Provider])
		if !containsStatus(q.States, s) {
			return false
		}
	}
	return true
}

// fold is the fake's stand-in for the store's own case folding.
//
// A real store folds case and diacritics in SQL, with collation that handles
// the accented titles the library already contains. This is deliberately not
// that: a fake that reimplemented collation would be a second implementation
// to keep correct, and the property tests care about is that an empty query
// filters nothing.
func fold(s string) string { return strings.ToLower(s) }

func containsID(xs []provider.ID, x provider.ID) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func containsStatus(xs []status.Status, x status.Status) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// Games returns matching games ordered by sort title.
func (f *Fake) Games(_ context.Context, q library.Query) ([]library.Game, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var out []library.Game
	for _, g := range f.games {
		if f.matches(g, q) {
			out = append(out, g)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].SortTitle != out[j].SortTitle {
			return out[i].SortTitle < out[j].SortTitle
		}
		return out[i].ID < out[j].ID
	})
	if q.Offset > 0 {
		if q.Offset >= len(out) {
			return nil, nil
		}
		out = out[q.Offset:]
	}
	if q.Limit > 0 && q.Limit < len(out) {
		out = out[:q.Limit]
	}
	return out, nil
}

// Game returns one game.
func (f *Fake) Game(_ context.Context, id library.GameID) (library.Game, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	g, ok := f.games[id]
	if !ok || g.MergedInto != nil {
		return library.Game{}, fault.New(fault.KindNotFound, "storetest.Game")
	}
	return g, nil
}

// Counts returns the summary for the same set Games would return.
//
// Limit and Offset are ignored, because the summary describes the filtered set
// rather than the visible page — the rule that stops a list view lying about
// what it is showing.
func (f *Fake) Counts(_ context.Context, q library.Query) (status.Counts, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var c status.Counts
	for _, g := range f.games {
		if f.matches(g, q) {
			c.Add(g.Status(f.Playtime[g.Provider]))
		}
	}
	return c, nil
}

// Connections returns every connected provider.
func (f *Fake) Connections(context.Context) ([]store.Connection, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]store.Connection, 0, len(f.conns))
	for _, c := range f.conns {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Provider < out[j].Provider })
	return out, nil
}

// LastRun returns the most recent run for a provider, or nil when there has
// never been one.
func (f *Fake) LastRun(_ context.Context, p provider.ID) (*store.SyncRun, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var best *store.SyncRun
	for _, r := range f.runs {
		if r.Provider != p {
			continue
		}
		if best == nil || r.StartedAt.After(best.StartedAt) {
			cp := r
			best = &cp
		}
	}
	return best, nil
}

// Setting reads one setting.
func (f *Fake) Setting(_ context.Context, key string) (string, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.settings[key]
	return v, ok, nil
}

// Metadata returns cached enrichment, stamped.
func (f *Fake) Metadata(_ context.Context, id library.GameID) (aged.Value[metadata.Metadata], bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.meta[id]
	return v, ok, nil
}

// Quote returns the last known price, stamped.
func (f *Fake) Quote(_ context.Context, id library.GameID, cur pricing.Currency) (aged.Value[pricing.Quote], bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.quotes[quoteKey{id, cur}]
	return v, ok, nil
}

// SetStatus sets or clears the manual status and stamps the change.
//
// A clear is a change: status_changed_at is written for it, because it is
// something the player did and because Phase 4's last-write-wins needs a
// timestamp for it to compare.
func (f *Fake) SetStatus(_ context.Context, id library.GameID, s *status.Status) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	g, ok := f.games[id]
	if !ok {
		return fault.New(fault.KindNotFound, "storetest.SetStatus")
	}
	if s != nil && !s.Valid() {
		return fault.New(fault.KindMalformed, "storetest.SetStatus")
	}
	now := f.Now()
	g.StatusManual = s
	g.StatusChangedAt = &now
	f.games[id] = g
	return nil
}

// UpsertBatch writes items and reports what changed.
//
// It never touches StatusManual or StatusChangedAt. That omission is the
// invariant: a sync updates playtime and last-played and nothing else, so a
// game marked ZERADO stays ZERADO through any number of syncs.
func (f *Fake) UpsertBatch(_ context.Context, p provider.ID, items []provider.Item) (store.BatchResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var res store.BatchResult
	for _, it := range items {
		if existing, id, ok := f.findRef(p, it.ProviderRef); ok {
			updated := existing
			updated.Title = it.Title
			updated.Platform = it.Platform
			updated.SortTitle = sortTitle(it.Title)
			updated.PlaytimeMinutes = it.PlaytimeMinutes
			updated.LastPlayed = it.LastPlayed
			updated.OwnedSince = it.OwnedSince
			// A game that came back has its absence cleared, silently.
			updated.AbsentSince = nil
			if changed(existing, updated) {
				res.Changed++
			} else {
				res.Unchanged++
			}
			f.games[id] = updated
			continue
		}
		f.nextID++
		f.games[f.nextID] = library.Game{
			ID:              f.nextID,
			UID:             library.UID(sortTitle(it.Title) + "|" + strings.ToLower(it.Platform)),
			Provider:        p,
			ProviderRef:     it.ProviderRef,
			Title:           it.Title,
			Platform:        it.Platform,
			SortTitle:       sortTitle(it.Title),
			PlaytimeMinutes: it.PlaytimeMinutes,
			LastPlayed:      it.LastPlayed,
			OwnedSince:      it.OwnedSince,
			AddedAt:         f.Now(),
		}
		res.New++
	}
	return res, nil
}

func changed(a, b library.Game) bool {
	return a.Title != b.Title ||
		a.Platform != b.Platform ||
		!eqIntPtr(a.PlaytimeMinutes, b.PlaytimeMinutes) ||
		!eqTimePtr(a.LastPlayed, b.LastPlayed) ||
		!eqTimePtr(a.OwnedSince, b.OwnedSince) ||
		(a.AbsentSince != nil) != (b.AbsentSince != nil)
}

func eqIntPtr(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func eqTimePtr(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equal(*b)
}

func sortTitle(s string) string {
	return strings.ToLower(strings.TrimPrefix(strings.ToLower(s), "the "))
}

func (f *Fake) findRef(p provider.ID, ref string) (library.Game, library.GameID, bool) {
	for id, g := range f.games {
		if g.Provider == p && g.ProviderRef == ref {
			return g, id, true
		}
	}
	return library.Game{}, 0, false
}

// AddManual writes a hand-entered item.
func (f *Fake) AddManual(_ context.Context, p provider.ID, it provider.Item, s *status.Status) (library.GameID, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, _, exists := f.findRef(p, it.ProviderRef); exists {
		return 0, fault.New(fault.KindConflict, "storetest.AddManual")
	}
	f.nextID++
	now := f.Now()
	g := library.Game{
		ID:          f.nextID,
		UID:         library.UID(sortTitle(it.Title) + "|" + strings.ToLower(it.Platform)),
		Provider:    p,
		ProviderRef: it.ProviderRef,
		Title:       it.Title,
		Platform:    it.Platform,
		SortTitle:   sortTitle(it.Title),
		OwnedSince:  it.OwnedSince,
		AddedAt:     now,
	}
	if s != nil {
		g.StatusManual = s
		g.StatusChangedAt = &now
	}
	f.games[f.nextID] = g
	return f.nextID, nil
}

// StartRun opens a run.
func (f *Fake) StartRun(_ context.Context, p provider.ID, at time.Time) (store.RunID, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextRun++
	f.runs[f.nextRun] = store.SyncRun{ID: f.nextRun, Provider: p, StartedAt: at, Status: store.RunRunning}
	return f.nextRun, nil
}

// FinishRun closes a run.
func (f *Fake) FinishRun(_ context.Context, id store.RunID, r store.RunResult) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	run, ok := f.runs[id]
	if !ok {
		return fault.New(fault.KindNotFound, "storetest.FinishRun")
	}
	fin := r.FinishedAt
	run.FinishedAt = &fin
	run.Status = r.Status
	run.Seen, run.New, run.Changed, run.Unchanged = r.Seen, r.New, r.Changed, r.Unchanged
	run.FaultKind = r.FaultKind
	f.runs[id] = run

	if c, ok := f.conns[run.Provider]; ok {
		c.LastSyncAt = &fin
		c.LastSyncStatus = r.Status
		f.conns[run.Provider] = c
	}
	return nil
}

// ReconcileAbsence tombstones rows a completed run did not return.
//
// The guard is the whole reason this method takes a RunID: only a run whose
// status is ok may tombstone anything, because in a truncated stream "not
// returned" and "not reached" are indistinguishable. A partial, failed or
// cancelled run is refused here rather than merely discouraged in a comment.
func (f *Fake) ReconcileAbsence(_ context.Context, id store.RunID, seen []string) (store.Absence, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	run, ok := f.runs[id]
	if !ok {
		return store.Absence{}, fault.New(fault.KindNotFound, "storetest.ReconcileAbsence")
	}
	if run.FinishedAt == nil || !run.Status.MayTombstone() {
		return store.Absence{}, fault.New(fault.KindPrecondition, "storetest.ReconcileAbsence",
			fault.WithMessage("fault.precondition"))
	}

	present := make(map[string]bool, len(seen))
	for _, r := range seen {
		present[r] = true
	}

	var a store.Absence
	now := f.Now()
	for id, g := range f.games {
		if g.Provider != run.Provider {
			continue
		}
		switch {
		case present[g.ProviderRef] && g.Absent():
			g.AbsentSince = nil
			f.games[id] = g
			a.Returned++
		case !present[g.ProviderRef] && !g.Absent():
			t := now
			g.AbsentSince = &t
			f.games[id] = g
			a.MarkedAbsent++
		}
	}
	return a, nil
}

// Delete removes a game. It happens only because the player asked.
func (f *Fake) Delete(_ context.Context, id library.GameID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.games[id]; !ok {
		return fault.New(fault.KindNotFound, "storetest.Delete")
	}
	delete(f.games, id)
	return nil
}

// SetSetting writes one setting.
func (f *Fake) SetSetting(_ context.Context, key, value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.settings[key] = value
	return nil
}

// SaveConnection records a connection. accountRef is an identifier, never a
// secret: there is no method here that could accept one.
func (f *Fake) SaveConnection(_ context.Context, p provider.ID, accountRef string, at time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	c := f.conns[p]
	c.Provider, c.AccountRef, c.ConnectedAt = p, accountRef, at
	f.conns[p] = c
	return nil
}

// DeleteConnection forgets a connection without touching the games it
// contributed.
func (f *Fake) DeleteConnection(_ context.Context, p provider.ID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.conns, p)
	return nil
}

// Close does nothing; there is no file.
func (f *Fake) Close() error { return nil }

// PutMetadata and PutQuote seed the Phase 2 and Phase 3 caches, so an offline
// test can exercise the stale-value screens without a metadata or price
// provider existing yet.
func (f *Fake) PutMetadata(id library.GameID, v aged.Value[metadata.Metadata]) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.meta[id] = v
}

// PutQuote seeds a cached price.
func (f *Fake) PutQuote(id library.GameID, cur pricing.Currency, v aged.Value[pricing.Quote]) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.quotes[quoteKey{id, cur}] = v
}

var _ store.Store = (*Fake)(nil)
