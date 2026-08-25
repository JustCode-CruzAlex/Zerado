package provider

import "time"

// Item is a provider's view of one owned title.
//
// It is what crosses the seam, in both directions of origin: a Steam response
// row becomes an Item, and so does a completed Z-08 form. That symmetry is the
// test the ticket asks for — if hand entry needed a different type, the seam
// would be Steam-shaped.
//
// # Every optional field is a pointer, and that is load-bearing
//
// Z-05 distinguishes three facts that a nullable column and a zero value
// cannot tell apart on their own:
//
//	"not tracked"   — the provider cannot ever know (Capabilities.Playtime false)
//	"—"             — not fetched yet (the field is nil)
//	"0h" / "never"  — the provider answered, and the answer was nothing
//
// The capability answers the first. The pointer answers the second and third.
// Collapsing either distinction produces a screen that tells a player their
// cartridge has been played for zero hours, which is a claim nothing could
// have made.
type Item struct {
	// ProviderRef is the provider's own identifier for this title — a Steam
	// appid, or a UUID Zerado minted for a hand-entered cartridge.
	//
	// Together with the provider ID it is the uniqueness constraint, which is
	// why a hand-entered duplicate is a note and never a block: the ref is
	// fresh every time, so two copies of the same cartridge are two rows and
	// nothing collides (Z-08 §8.3).
	ProviderRef string

	// Title is what the player calls it. Required.
	Title string

	// Platform is required, and it is half of the Phase 4 merge identity.
	//
	// Required even for Steam, where it is unglamorous, because a UID derived
	// from title alone would merge a PS2 disc with its PC remaster on two
	// devices that had every reason to keep them apart.
	Platform string

	// PlaytimeMinutes is provider-reported. nil means not reported.
	//
	// Zero is a real value and is distinct from nil. A provider whose
	// Capabilities.Playtime is false must always send nil here: sending a
	// zero would feed status.Derive a number it would treat as evidence.
	PlaytimeMinutes *int

	// LastPlayed is provider-reported. nil means not reported, which is not
	// the same as never played — Z-05 renders those as "—" and "never
	// played" respectively.
	LastPlayed *time.Time

	// OwnedSince is when the player acquired it, where that is knowable.
	OwnedSince *time.Time

	// Extra carries provider-specific identifiers a later phase will want,
	// keyed by a name the provider owns — "steam_appid", "gog_id".
	//
	// It is a deliberately narrow escape hatch and it has one rule: nothing
	// above this seam may read a key by name. The store persists what it
	// recognises and drops what it does not, so a provider can start emitting
	// a new key without a schema change and without a screen learning about
	// it. A screen reading Extra["steam_appid"] has reintroduced the switch on
	// provider identity that this whole seam exists to remove.
	Extra map[string]string
}

// Playtime returns the reported playtime and whether it was reported.
//
// Callers use this rather than dereferencing, because the two-value form makes
// the "not reported" branch impossible to forget at a call site — which is
// exactly the branch Z-05 renders differently.
func (i Item) Playtime() (int, bool) {
	if i.PlaytimeMinutes == nil {
		return 0, false
	}
	return *i.PlaytimeMinutes, true
}

// EntryField is one field on the hand-entry form, declared by the provider
// that will receive it.
//
// It is the Z-08 twin of [CredentialField], and the two are separate types on
// purpose: a credential is a secret with a destination (the Vault) and a
// help URL, and an entry field is a data field with a default and a
// requirement. Merging them would give each half the other's irrelevant
// fields, and would put Secret on a form where it can only ever be false.
type EntryField struct {
	// Key names the field to the provider — "title", "platform", "state",
	// "owned_since". It is never rendered.
	Key string

	// LabelKey and HelpKey name catalogue entries. They are keys rather than
	// strings because these are user-facing text and ADR-0001 D9 admits no
	// literal.
	LabelKey string
	HelpKey  string

	// Kind tells the screen which editor to render.
	Kind FieldKind

	// Required drives the empty-submit refusal. Z-08 asks for four fields and
	// requires two, because a cartridge has a title and a platform and may
	// have nothing else.
	Required bool

	// Options is the choice set for FieldChoice, in display order. Each entry
	// is a stored value plus a catalogue key for its label.
	Options []FieldOption
}

// FieldKind is the editor a form field needs.
//
// The set is closed and small: it is the set Z-02 and Z-08 between them
// render, and a provider that needs a sixth kind is a provider asking for a
// screen change, which is exactly the event this seam is supposed to make
// visible rather than easy.
type FieldKind uint8

const (
	// FieldText is a single-line text editor.
	FieldText FieldKind = iota

	// FieldChoice is a fixed list rendered as chips, one selected.
	FieldChoice

	// FieldDate is a date, entered as text and parsed by the provider.
	FieldDate
)

// FieldOption is one entry in a [FieldChoice] field.
type FieldOption struct {
	// Value is the stored form — "not_started".
	Value string

	// LabelKey names the catalogue entry that renders it.
	LabelKey string
}

// Entry is a completed hand-entry form, keyed by [EntryField.Key].
//
// A map rather than a struct because the field set is the provider's to
// declare: a struct here would mean the screen and the seam agreeing on a
// shape in advance, which is the coupling Z-08 exists without.
type Entry map[string]string
