// Package vault is where the player's own keys live, and it is deliberately
// not the library file.
//
// The Steam API key is the player's own key, and library.db is a file the
// player is explicitly invited to back up, move and delete. A credential
// inside it would be a credential in every backup, every copy, and every
// support-ticket attachment.
//
// This strengthens the one-file promise rather than violating it: the library
// file stays purely a library, so sharing it is safe. Settings shows which
// backing is in use because a security property the player cannot see is a
// security property they cannot rely on.
package vault

import (
	"context"

	"github.com/JustCode-CruzAlex/Zerado/internal/provider"
)

// Vault stores and retrieves per-provider secrets.
//
// It is keyed by (provider, key) rather than by a single string, so that
// forgetting one provider cannot reach another's secrets and so that
// [Vault.DeleteProvider] is a complete operation rather than a loop the caller
// has to get right.
type Vault interface {
	// Get returns a stored secret. ok is false when nothing is stored, which
	// is not an error — a first run has no keys and that is the normal state.
	Get(ctx context.Context, p provider.ID, key string) (value string, ok bool, err error)

	// Set stores a secret, replacing any previous value.
	Set(ctx context.Context, p provider.ID, key, value string) error

	// Delete forgets one secret.
	Delete(ctx context.Context, p provider.ID, key string) error

	// DeleteProvider forgets every secret for a provider.
	//
	// It exists because Z-09's disconnect must leave nothing behind, and a
	// caller looping over the field keys it happens to remember would leave
	// behind exactly the field somebody added last week.
	DeleteProvider(ctx context.Context, p provider.ID) error

	// Backing reports where secrets are actually being kept.
	//
	// Never a guess and never aspirational: it reports what this process
	// resolved at start-up, so a keychain that was unavailable reads as the
	// file backing rather than as the keychain it wanted to be.
	Backing() Backing
}

// Backing is where the vault is storing secrets, in enough detail for Z-09 to
// name it.
//
// The spine's shape returned a string, "keychain" or "file". Z-09 needs more
// than that: its copy is "In the macOS Keychain", "In the GNOME keyring", "In
// Windows Credential Manager" — whichever the vault reports — and a two-valued
// string cannot distinguish those three. A struct with a kind and a platform
// name can, without the screen switching on runtime.GOOS, which would put a
// platform assumption in a renderer.
type Backing struct {
	// Kind is which of the two mechanisms is in use.
	Kind BackingKind

	// NameKey names the catalogue entry that renders the specific backing —
	// the macOS Keychain, the GNOME keyring, Windows Credential Manager, or
	// the credentials file with its mode.
	//
	// The vault supplies the key because the vault is the only thing that
	// knows which service it actually opened. A screen that derived this from
	// the operating system would be right by coincidence and wrong inside a
	// container.
	NameKey string

	// Path is the file's location when Kind is BackingFile, and empty
	// otherwise. Shown in Settings so the player can find, inspect and delete
	// it.
	Path string
}

// BackingKind is the storage mechanism.
type BackingKind uint8

const (
	// BackingUnknown is the zero value, before resolution. A vault that
	// reports it has not started up correctly.
	BackingUnknown BackingKind = iota

	// BackingKeychain is the operating system's own secret store: macOS
	// Keychain, Secret Service, Windows Credential Manager. Preferred.
	BackingKeychain

	// BackingFile is credentials.json beside the library, mode 0600.
	//
	// It is a fallback and not a failure. Headless Linux and containers are a
	// real part of a terminal-first audience, and a keychain-only design would
	// simply not run there — so Settings states the backing plainly rather
	// than the product pretending both are the same.
	BackingFile
)

// String returns the stable machine name, used in the CLI's JSON output.
func (k BackingKind) String() string {
	switch k {
	case BackingKeychain:
		return "keychain"
	case BackingFile:
		return "file"
	default:
		return "unknown"
	}
}
