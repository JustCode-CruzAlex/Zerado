package devicesync_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/JustCode-CruzAlex/Zerado/internal/devicesync"
)

// TestPayloadCarriesOnlyWhatThePlayerTyped is the point of this package.
//
// ADR-0001 D4 decides what crosses the Phase 4 boundary, and it is the most
// expensive decision in the bundle to reverse because it decides what the
// schema carries from the first migration onward. A paragraph stating that
// rule is a paragraph somebody will not have read in 2027, when adding
// playtime to the payload will look like an obvious improvement — the server
// already has the rows, and showing hours on another device is a nice feature.
//
// So the rule is a test. A field added to the payload fails here, by name,
// with the decision attached.
func TestPayloadCarriesOnlyWhatThePlayerTyped(t *testing.T) {
	allowed := map[reflect.Type]map[string]bool{
		reflect.TypeOf(devicesync.Change{}): {
			"UID":       true, // the cross-device merge hint
			"Status":    true, // the manual value, or nil for a clear
			"ChangedAt": true, // decides the conflict
		},
		reflect.TypeOf(devicesync.ManualGame{}): {
			"UID":        true,
			"Title":      true, // the player typed it
			"Platform":   true, // the player typed it
			"CreatedAt":  true,
			"OwnedSince": true, // the one optional thing the player typed
		},
		reflect.TypeOf(devicesync.Envelope{}): {
			"Device":  true,
			"Since":   true,
			"Changes": true,
			"Manual":  true,
		},
	}

	for typ, ok := range allowed {
		for i := 0; i < typ.NumField(); i++ {
			name := typ.Field(i).Name
			if !ok[name] {
				t.Errorf("%s.%s crosses the Phase 4 sync boundary and is not on the allow-list.\n"+
					"ADR-0001 D4: only what the player typed crosses. Everything a machine can\n"+
					"recompute, each device recomputes. If this field is genuinely something the\n"+
					"player typed, add it to the allow-list in this test with a reason. If it is a\n"+
					"provider fact — playtime, last-played, a cover, a price — it must not cross.",
					typ.Name(), name)
			}
		}
	}
}

// TestNothingResemblingACredentialCanCross is the same guard aimed at the one
// class of field that would be a breach rather than a design regression. A
// ratified promise says the keys are the player's own; centralising them makes
// Zerado a credential custodian, which is a different company.
func TestNothingResemblingACredentialCanCross(t *testing.T) {
	banned := []string{"key", "token", "secret", "password", "credential", "auth", "cookie", "session"}
	types := []reflect.Type{
		reflect.TypeOf(devicesync.Change{}),
		reflect.TypeOf(devicesync.ManualGame{}),
		reflect.TypeOf(devicesync.Envelope{}),
		reflect.TypeOf(devicesync.Receipt{}),
	}
	for _, typ := range types {
		for i := 0; i < typ.NumField(); i++ {
			name := strings.ToLower(typ.Field(i).Name)
			for _, b := range banned {
				if strings.Contains(name, b) {
					t.Fatalf("%s.%s looks like a credential. Credentials never cross the Phase 4 "+
						"boundary — each device connects with its own keys (ADR-0001 D4).",
						typ.Name(), typ.Field(i).Name)
				}
			}
		}
	}
}

// TestAClearedOverrideCanCross: "I have no opinion" is a change the player
// made, so Status is a pointer and nil is meaningful. A non-nullable field
// would make a clear unsyncable, and the player's device would silently
// disagree with itself forever.
func TestAClearedOverrideCanCross(t *testing.T) {
	f, ok := reflect.TypeOf(devicesync.Change{}).FieldByName("Status")
	if !ok {
		t.Fatal("Change has no Status field")
	}
	if f.Type.Kind() != reflect.Pointer {
		t.Fatalf("Change.Status is %v; it must be nullable so that clearing an override can cross", f.Type.Kind())
	}
}

// TestTheMergeKeyIsTheCrossDeviceIdentity: a local surrogate key means nothing
// on another machine, which is exactly why library.UID exists in Phase 1.
func TestTheMergeKeyIsTheCrossDeviceIdentity(t *testing.T) {
	f, _ := reflect.TypeOf(devicesync.Change{}).FieldByName("UID")
	if got := f.Type.Name(); got != "UID" {
		t.Fatalf("Change.UID is a %s; a merge keyed on a local id would match the wrong rows", got)
	}
}

// TestPayloadTypesStayFlat closes the gap the review at c4c8d95 named: the
// two guards above walk depth 1, because reflect.NumField does not descend.
//
// A future field whose type is a nested struct would pass both — the
// allow-list sees only the outer name and the credential-name ban sees only
// the outer name — while carrying anything at all inside it. That is not a
// hypothetical shape for a sync payload: a `Device DeviceInfo` or a
// `Source ProviderRef` is exactly what somebody adds in good faith.
//
// So rather than note the limitation, this pins the *types* as well as the
// names. Every field of every payload struct must be one of a small allow-list
// of leaf types, which makes nesting impossible without editing this test —
// and editing this test is the moment somebody reads D4.
func TestPayloadTypesStayFlat(t *testing.T) {
	allowed := map[string]bool{
		"string":                  true,
		"devicesync.DeviceID":     true,
		"library.UID":             true,
		"time.Time":               true,
		"*time.Time":              true,
		"*status.Status":          true,
		"[]devicesync.Change":     true,
		"[]devicesync.ManualGame": true,
		"int":                     true,
	}
	for _, typ := range []reflect.Type{
		reflect.TypeOf(devicesync.Change{}),
		reflect.TypeOf(devicesync.ManualGame{}),
		reflect.TypeOf(devicesync.Envelope{}),
		reflect.TypeOf(devicesync.Receipt{}),
	} {
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if !allowed[f.Type.String()] {
				t.Errorf("%s.%s is a %s, which is not an allowed payload type.\n"+
					"The allow-list is leaf types only, so that a nested struct cannot smuggle\n"+
					"fields past the depth-1 checks above. If this type is genuinely something\n"+
					"the player typed, add it here — and read ADR-0001 D4 while you are in this\n"+
					"file, because that is what editing this list is for.",
					typ.Name(), f.Name, f.Type)
			}
		}
	}
}
