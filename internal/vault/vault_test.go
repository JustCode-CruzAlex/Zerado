package vault_test

import (
	"context"
	"testing"

	"github.com/JustCode-CruzAlex/Zerado/internal/vault"
	"github.com/JustCode-CruzAlex/Zerado/internal/vault/vaulttest"
)

// TestDisconnectLeavesNothingBehind. DeleteProvider exists because a caller
// looping over the field keys it happens to remember would leave behind
// exactly the field somebody added last week.
func TestDisconnectLeavesNothingBehind(t *testing.T) {
	v := vaulttest.New()
	ctx := context.Background()
	_ = v.Set(ctx, "steam", "api_key", "SECRET")
	_ = v.Set(ctx, "steam", "refresh_token", "ALSO-SECRET")
	_ = v.Set(ctx, "gog", "api_key", "OTHER")

	if err := v.DeleteProvider(ctx, "steam"); err != nil {
		t.Fatalf("DeleteProvider: %v", err)
	}
	if _, ok, _ := v.Get(ctx, "steam", "refresh_token"); ok {
		t.Fatal("a secret survived the disconnect")
	}
	if _, ok, _ := v.Get(ctx, "gog", "api_key"); !ok {
		t.Fatal("forgetting one provider reached another's secrets")
	}
	if v.Len() != 1 {
		t.Fatalf("%d secrets remain, want 1", v.Len())
	}
}

// TestBackingIsNamedSpecificallyEnoughForSettings. The spine's shape returned
// "keychain" or "file"; Z-09's copy names the macOS Keychain, the GNOME
// keyring and Windows Credential Manager, which a two-valued string cannot
// distinguish — and which a screen must not derive from runtime.GOOS, because
// that is right by coincidence and wrong inside a container.
func TestBackingIsNamedSpecificallyEnoughForSettings(t *testing.T) {
	v := vaulttest.New()
	b := v.Backing()
	if b.Kind == vault.BackingUnknown {
		t.Fatal("the vault reports an unresolved backing; a security property the player cannot see is one they cannot rely on")
	}
	if b.NameKey == "" {
		t.Fatal("Settings has no catalogue key for the specific backing in use")
	}
	if b.Kind == vault.BackingFile && b.Path == "" {
		t.Fatal("the file backing does not say where the file is; the player cannot find, inspect or delete it")
	}

	v.Back = vault.Backing{Kind: vault.BackingKeychain, NameKey: "settings.vault.keychain.macos"}
	if got := v.Backing().Kind.String(); got != "keychain" {
		t.Fatalf("BackingKind name = %q", got)
	}
}

// TestMissingIsNotAnError: a first run has no keys and that is the normal
// state, not a failure to report.
func TestMissingIsNotAnError(t *testing.T) {
	v := vaulttest.New()
	got, ok, err := v.Get(context.Background(), "steam", "api_key")
	if err != nil {
		t.Fatalf("Get on an empty vault errored: %v", err)
	}
	if ok || got != "" {
		t.Fatalf("Get = %q, %v", got, ok)
	}
}
