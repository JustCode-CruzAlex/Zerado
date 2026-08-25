package cli_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/JustCode-CruzAlex/Zerado/internal/cli"
	"github.com/JustCode-CruzAlex/Zerado/internal/fault"
)

// TestExitCodesAreTotal fails the moment a Kind is added without deciding what
// a script should do about it. A taxonomy the CLI cannot express is a taxonomy
// that reaches a shell as 1.
func TestExitCodesAreTotal(t *testing.T) {
	for _, k := range fault.Kinds() {
		err := fault.New(k, "test.Op")
		code := cli.ExitCode(err)
		if code == cli.ExitOK {
			t.Errorf("%v maps to exit 0; a failure reported as success is the worst outcome for a scripted caller", k)
		}
		if code < 0 || code > 255 {
			t.Errorf("%v maps to %d, which is not a valid process exit status", k, code)
		}
	}
	if got := cli.ExitCode(nil); got != cli.ExitOK {
		t.Fatalf("ExitCode(nil) = %d", got)
	}
}

// TestTheFourNetworkFailuresDoNotCollapse: a cron job wants to retry on
// offline and alert on a rejected key, and it cannot do that if both are 1.
func TestTheFourNetworkFailuresDoNotCollapse(t *testing.T) {
	codes := map[int]fault.Kind{}
	for _, k := range []fault.Kind{fault.KindOffline, fault.KindUnreachable, fault.KindUnauthorized, fault.KindEmpty, fault.KindRateLimited} {
		c := cli.ExitCode(fault.New(k, "steam.Sync"))
		if prev, dup := codes[c]; dup {
			t.Fatalf("%v and %v both exit %d; a script cannot tell them apart", prev, k, c)
		}
		codes[c] = k
	}
}

// TestCancelledUsesTheShellConvention: a player who pressed Ctrl-C sees the
// number they would have seen anyway.
func TestCancelledUsesTheShellConvention(t *testing.T) {
	if got := cli.ExitCode(fault.New(fault.KindCancelled, "steam.Sync")); got != 130 {
		t.Fatalf("cancelled exits %d, want 130", got)
	}
}

// TestUnclassifiedIsNotANetworkCode: a scripted caller retrying forever
// because an internal defect was reported as a timeout is worse than a loud 1.
func TestUnclassifiedIsNotANetworkCode(t *testing.T) {
	if got := cli.ExitCode(errString("boom")); got != cli.ExitInternal {
		t.Fatalf("an unclassified error exits %d, want %d", got, cli.ExitInternal)
	}
}

type errString string

func (e errString) Error() string { return string(e) }

// TestVerbSurfaceIsStable is a golden list. A verb name in somebody's shell
// history has no upgrade path, so renaming one has to be a deliberate act that
// updates this test and says why in the diff.
func TestVerbSurfaceIsStable(t *testing.T) {
	want := []string{
		"", "sync", "list", "mark", "add", "game", "doctor", "version", "help",
		"tonight", "price", "watch", "tag", "enrich", "devices", "export", "import",
	}
	got := make([]string, 0, len(want))
	for _, v := range cli.Surface() {
		got = append(got, v.Name)
	}
	if len(got) != len(want) {
		t.Fatalf("the verb surface has %d entries, the golden list has %d:\n got %v\nwant %v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("verb %d is %q, want %q", i, got[i], want[i])
		}
	}
}

// TestNoDuplicateVerbsAndNoDuplicateFlags.
func TestNoDuplicateVerbsAndNoDuplicateFlags(t *testing.T) {
	seen := map[string]bool{}
	for _, v := range cli.Surface() {
		if seen[v.Name] {
			t.Fatalf("duplicate verb %q", v.Name)
		}
		seen[v.Name] = true

		flags := map[string]bool{}
		for _, f := range append(cli.GlobalFlags(), v.Flags...) {
			if flags[f.Name] {
				t.Fatalf("verb %q declares --%s twice, or shadows a global flag", v.Name, f.Name)
			}
			flags[f.Name] = true
		}
	}
}

// TestReservedVerbsPromiseNothing guards anti-pattern 14: a name is claimed so
// a future feature does not have to settle for a worse word, and nothing about
// it is advertised as working.
func TestReservedVerbsPromiseNothing(t *testing.T) {
	for _, v := range cli.Surface() {
		if v.Phase != cli.PhaseLater {
			continue
		}
		if v.Stability != cli.StabilityReserved {
			t.Errorf("%q is reserved but promises stability", v.Name)
		}
		if len(v.Args) > 0 || len(v.Flags) > 0 {
			t.Errorf("%q is reserved but declares arguments; a reserved verb has no shape yet", v.Name)
		}
	}
	if v, ok := cli.Lookup("tonight"); !ok || v.Phase != cli.PhaseLater {
		t.Fatal("tonight must be reserved: it is Phase 2 and its name must not be squatted")
	}
}

// TestEveryVerbHasACatalogueKey: help text is user-facing, and D9 admits no
// literal even there.
func TestEveryVerbHasACatalogueKey(t *testing.T) {
	for _, v := range cli.Surface() {
		if v.SummaryKey == "" {
			t.Errorf("verb %q has no summary key", v.Name)
		}
		for _, f := range v.Flags {
			if f.SummaryKey == "" {
				t.Errorf("verb %q flag --%s has no summary key", v.Name, f.Name)
			}
		}
	}
	for _, f := range cli.GlobalFlags() {
		if f.SummaryKey == "" {
			t.Errorf("global flag --%s has no summary key", f.Name)
		}
	}
}

// TestOfflineClassIsDeclared: a scripted caller can tell, before running
// anything, which verbs refuse on a plane. Only the two verbs that are
// definitionally about reaching somewhere else may claim it.
func TestOfflineClassIsDeclared(t *testing.T) {
	for _, v := range cli.Surface() {
		if v.NeedsNetwork && v.Name != "sync" {
			t.Errorf("verb %q claims it needs the network; every Phase 1 verb but sync reads local state", v.Name)
		}
	}
	if v, _ := cli.Lookup("list"); v.NeedsNetwork {
		t.Fatal("list needs the network; it is a WHERE clause")
	}
}

// TestTheEnvelopeCarriesNoMessage. A pre-rendered English sentence in a JSON
// field is the hardest D9 violation to spot, and a JSON envelope is the single
// most likely thing in this product to be piped into a bug report.
func TestTheEnvelopeCarriesNoMessage(t *testing.T) {
	leak := errString(`Get "https://api.steampowered.com/x?key=SECRET-KEY-9F2": no route`)
	f := fault.New(fault.KindOffline, "steam.Sync",
		fault.WithSubject("Steam"), fault.WithCause(leak))

	b, err := json.Marshal(cli.Err(f))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, forbidden := range []string{"SECRET-KEY-9F2", "api.steampowered.com", "no route"} {
		if contains(s, forbidden) {
			t.Fatalf("the JSON envelope leaked %q: %s", forbidden, s)
		}
	}

	var env cli.Envelope
	if err := json.Unmarshal(b, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.API != cli.APIVersion {
		t.Fatalf("API = %d, want %d", env.API, cli.APIVersion)
	}
	if env.OK {
		t.Fatal("a failure envelope reports ok")
	}
	if env.Error == nil || env.Error.Kind != "offline" {
		t.Fatalf("Error = %+v", env.Error)
	}
	if env.Error.MessageKey == "" {
		t.Fatal("no catalogue key; a consumer with the catalogue could not render the same sentence Zerado would")
	}
}

// TestRetryAfterReachesTheEnvelopeAsSeconds: every JSON consumer can compare a
// number and not all of them can parse "1m30s".
func TestRetryAfterReachesTheEnvelopeAsSeconds(t *testing.T) {
	f := fault.New(fault.KindRateLimited, "steam.Sync", fault.WithRetryAfter(90*time.Second))
	env := cli.Err(f)
	if env.Error.RetryAfterSeconds != 90 {
		t.Fatalf("RetryAfterSeconds = %d, want 90", env.Error.RetryAfterSeconds)
	}
}

// TestSuccessAndFailureShareOneShape: a caller that has to parse two different
// top-level structures depending on an exit code will get it wrong on the day
// it matters.
func TestSuccessAndFailureShareOneShape(t *testing.T) {
	ok := cli.Ok(map[string]int{"games": 247})
	if !ok.OK || ok.Error != nil || ok.API != cli.APIVersion {
		t.Fatalf("success envelope = %+v", ok)
	}
	bad := cli.Err(fault.New(fault.KindNotFound, "store.Game"))
	if bad.OK || bad.Data != nil {
		t.Fatalf("failure envelope = %+v", bad)
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

// TestNoKindSilentlyInheritsTheInternalExitCode is the CLI half of the
// taxonomy-totality repair.
//
// ExitCode's default is ExitInternal, so a Kind nobody gave a code to still
// exits 1 — which is a plausible-looking answer and therefore the dangerous
// one. Only KindInternal is entitled to it.
func TestNoKindSilentlyInheritsTheInternalExitCode(t *testing.T) {
	for _, k := range fault.Kinds() {
		if k == fault.KindInternal {
			continue
		}
		if got := cli.ExitCode(fault.New(k, "test.Op")); got == cli.ExitInternal {
			t.Errorf("kind %v exits %d; it has no case in ExitCode and a script would read it as our defect", k, got)
		}
	}
}
