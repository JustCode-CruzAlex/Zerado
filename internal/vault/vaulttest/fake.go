// Package vaulttest holds an in-memory vault, so a test never touches a real
// keychain and never writes a file.
package vaulttest

import (
	"context"
	"sync"

	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
	"github.com/JustCode-CruzAlex/Zerado/internal/vault"
)

// Fake stores secrets in a map.
//
// Backing is settable so a test can exercise both of Z-09's honest lines —
// "In the OS keychain" and "In credentials.json, mode 0600" — without a
// platform that has one and a platform that does not.
type Fake struct {
	mu   sync.Mutex
	data map[key]string

	// Back is what Backing reports. Defaults to the file backing, which is
	// the pessimistic one: a test that passes against the fallback passes
	// against the keychain.
	Back vault.Backing
}

type key struct {
	p provider.ID
	k string
}

// New returns an empty Fake reporting the file backing.
func New() *Fake {
	return &Fake{
		data: map[key]string{},
		Back: vault.Backing{Kind: vault.BackingFile, NameKey: "settings.vault.file", Path: "/tmp/credentials.json"},
	}
}

// Get returns a stored secret.
func (f *Fake) Get(_ context.Context, p provider.ID, k string) (string, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.data[key{p, k}]
	return v, ok, nil
}

// Set stores a secret.
func (f *Fake) Set(_ context.Context, p provider.ID, k, v string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.data[key{p, k}] = v
	return nil
}

// Delete forgets one secret.
func (f *Fake) Delete(_ context.Context, p provider.ID, k string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.data, key{p, k})
	return nil
}

// DeleteProvider forgets every secret for a provider, which is what a
// disconnect needs and what a caller looping over remembered field names
// would get wrong.
func (f *Fake) DeleteProvider(_ context.Context, p provider.ID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for k := range f.data {
		if k.p == p {
			delete(f.data, k)
		}
	}
	return nil
}

// Backing reports where secrets are kept.
func (f *Fake) Backing() vault.Backing { return f.Back }

// Len reports how many secrets are held, for a test asserting that a
// disconnect left nothing behind.
func (f *Fake) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.data)
}

var _ vault.Vault = (*Fake)(nil)
