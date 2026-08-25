package arch_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const module = "github.com/JustCode-CruzAlex/Zerado"

// pkgImports maps each package's directory (relative to the module root) to
// the set of paths its non-test files import.
//
// Test files are excluded deliberately: a test may import a fake, and a fake
// may do things production code must not. What is being asserted here is what
// ships.
func pkgImports(t *testing.T) map[string][]string {
	t.Helper()
	root := moduleRoot(t)
	out := map[string][]string{}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "site" || d.Name() == "docs" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, perr := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if perr != nil {
			return perr
		}
		rel, _ := filepath.Rel(root, filepath.Dir(path))
		for _, imp := range f.Imports {
			p, uerr := strconv.Unquote(imp.Path.Value)
			if uerr != nil {
				return uerr
			}
			out[rel] = append(out[rel], p)
		}
		if _, seen := out[rel]; !seen {
			out[rel] = nil
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the module: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("no packages found; the walk is broken and every assertion below would pass vacuously")
	}
	return out
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found above the working directory")
		}
		dir = parent
	}
}

// TestOnlyProvidersMayReachTheNetwork is 07-offline-contract.md §7.3's grep
// rule, executable.
//
// It is the mechanism behind two published promises at once. A screen that
// only reads local state works offline because it cannot do anything else, and
// "the only network traffic is Zerado reaching out to the services you've
// connected" is checkable rather than asserted — there is exactly one kind of
// package where a packet can originate.
func TestOnlyProvidersMayReachTheNetwork(t *testing.T) {
	networking := []string{"net/http", "net/url", "net", "crypto/tls", "net/rpc"}
	for pkg, imports := range pkgImports(t) {
		if strings.HasPrefix(pkg, filepath.Join("internal", "provider")) {
			continue // the one place a request may originate
		}
		for _, imp := range imports {
			for _, n := range networking {
				if imp == n {
					t.Errorf("%s imports %q.\n"+
						"Only provider packages may reach the network (07-offline-contract.md §7.3).\n"+
						"If this is error classification rather than I/O, hand fault.Classify a\n"+
						"fault.Transport verdict instead — that is exactly why it takes one.", pkg, imp)
				}
			}
		}
	}
}

// TestNoProviderConstructsItsOwnClient. 06-data-seams.md §7 declines to make
// HTTP a seam and requires a shared client with a timeout, injected — because
// a provider that builds its own is how a "works offline" claim quietly stops
// being true, and how a request without a timeout hangs the sync screen
// forever.
//
// The contracts stage of this ticket ships no provider that does I/O, so this
// currently asserts a vacuous truth. It is here so that the first real
// provider inherits the rule rather than discovering it in review.
func TestNoProviderConstructsItsOwnClient(t *testing.T) {
	root := moduleRoot(t)
	dir := filepath.Join(root, "internal", "provider")
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		b, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		for _, bad := range []string{"http.Client{", "http.DefaultClient", "http.Transport{"} {
			if strings.Contains(string(b), bad) {
				rel, _ := filepath.Rel(root, path)
				t.Errorf("%s constructs %s. The client is shared and injected, with a timeout "+
					"(06-data-seams.md §7).", rel, bad)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the provider packages: %v", err)
	}
}

// TestOnlyTheStoreKnowsThereIsADatabase. The persistence seam is the only
// thing in the program that knows a database exists, and that is what makes it
// the boundary the Phase 4 sync engine attaches to without touching a screen.
func TestOnlyTheStoreKnowsThereIsADatabase(t *testing.T) {
	dbish := []string{"database/sql", "modernc.org/sqlite", "github.com/mattn/go-sqlite3"}
	for pkg, imports := range pkgImports(t) {
		if strings.HasPrefix(pkg, filepath.Join("internal", "store")) {
			continue
		}
		for _, imp := range imports {
			for _, d := range dbish {
				if imp == d {
					t.Errorf("%s imports %q; every data access is behind the store interface", pkg, imp)
				}
			}
		}
	}
}

// TestAProviderNeverTalksToAScreenOrTheStore. 06-data-seams.md §1: a sync
// writes to the store and the screens read from the store, and the two never
// meet. If the provider seam could import the store, that rule would be a
// habit rather than a structure.
func TestAProviderNeverTalksToAScreenOrTheStore(t *testing.T) {
	forbidden := map[string][]string{
		filepath.Join("internal", "provider"): {
			module + "/internal/store",
			module + "/internal/library",
		},
		filepath.Join("internal", "library"): {
			module + "/internal/store",
		},
	}
	for pkg, bans := range forbidden {
		for _, imp := range pkgImports(t)[pkg] {
			for _, b := range bans {
				if imp == b {
					t.Errorf("%s imports %s; the seam only points one way", pkg, b)
				}
			}
		}
	}
}

// TestTheLeafPackagesStayLeaves. fault, i18n, status and aged are imported by
// everything, so a dependency added to one of them is a dependency added to
// the whole product — and a cycle in the seam graph is a refactor rather than
// a fix.
func TestTheLeafPackagesStayLeaves(t *testing.T) {
	allowed := map[string]map[string]bool{
		filepath.Join("internal", "i18n"):   {},
		filepath.Join("internal", "status"): {},
		filepath.Join("internal", "aged"):   {},
		filepath.Join("internal", "fault"):  {module + "/internal/i18n": true},
	}
	for pkg, ok := range allowed {
		for _, imp := range pkgImports(t)[pkg] {
			if !strings.HasPrefix(imp, module) {
				continue // the standard library is fine
			}
			if !ok[imp] {
				t.Errorf("%s imports %s; it is a leaf and everything depends on it", pkg, imp)
			}
		}
	}
}

// TestNoShippedCodeImportsAFake. A fake reaching production is how an
// in-memory store ends up behind a real screen, and it is silent when it
// happens.
func TestNoShippedCodeImportsAFake(t *testing.T) {
	for pkg, imports := range pkgImports(t) {
		for _, imp := range imports {
			if !strings.HasPrefix(imp, module) {
				continue
			}
			base := imp[strings.LastIndex(imp, "/")+1:]
			if strings.HasSuffix(base, "test") && base != "test" {
				t.Errorf("%s imports the test double %s from non-test code", pkg, imp)
			}
		}
	}
}

// TestTheModuleHasNoThirdPartyDependencies, for now.
//
// Not a permanent rule — Bubble Tea, x/text and a SQLite driver are all
// decided and all arrive with the implementation. It is an assertion about
// *this* stage: the contracts are pure Go with a standard-library-only
// dependency set, which is what lets every seam be exercised offline with no
// module download and means this module has no go.sum at all — go.mod
// declares no require, so there is nothing to lock.
func TestTheModuleHasNoThirdPartyDependencies(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(moduleRoot(t), "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "require") {
		t.Skip("dependencies have arrived with the implementation; this stage-assertion has served its purpose")
	}
}
