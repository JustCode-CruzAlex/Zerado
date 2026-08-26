package arch_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
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

// TestEveryImportIsStandardLibraryOrLocal asserts what this stage of the
// contracts actually promises: the module's dependency set is the Go standard
// library and nothing else.
//
// It replaces TestTheModuleHasNoThirdPartyDependencies, which could not fail.
// That test read go.mod and either ended without asserting anything (no
// require) or called t.Skip (require present) — there was no input for which
// it reported a failure. A test that cannot fail is a green light nobody
// earned, which is the same defect class as the hand-maintained fault.Kinds()
// slice repaired in the previous round, and it was found the same way: by
// somebody asking what would make it go red.
//
// This one has a real subject. It walks every non-test import in the module
// and fails on any path that is neither module-local nor standard library.
//
// It is not a permanent rule. Bubble Tea, x/text and a SQLite driver are all
// decided and all arrive with the implementation; at that point this test is
// updated to an allow-list of the decided dependencies rather than deleted,
// because "which third-party packages may this module import" stays a question
// worth answering out loud.
//
// The stdlib test is the usual one: a standard-library path has no dot in its
// first segment, because a dot there is a domain name. It is applied after the
// module-local check, since the module path itself contains one.
func TestEveryImportIsStandardLibraryOrLocal(t *testing.T) {
	var checked int
	for pkg, imports := range pkgImports(t) {
		for _, imp := range imports {
			checked++
			if strings.HasPrefix(imp, module) {
				continue
			}
			first, _, _ := strings.Cut(imp, "/")
			if strings.Contains(first, ".") {
				t.Errorf("%s imports %q, which is a third-party dependency.\n"+
					"At this stage the contracts are standard-library-only: that is what lets\n"+
					"every seam be exercised offline with no module download, and why this\n"+
					"module has no go.sum. If the dependency is decided, add it to this test's\n"+
					"allow-list so the set stays something somebody chose.", pkg, imp)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no imports were examined; the walk is broken and this assertion would pass vacuously")
	}
}

// TestEveryDocReferenceResolves fails on a doc comment that cites a document
// which does not exist.
//
// This is the third instance of the same defect in two review rounds: two
// files cited a dossier document under a name it carried before it shipped,
// and both were found by a human reading comments — the most expensive way to
// find a broken link, and the one that stops happening the moment nobody has
// time.
//
// It matters more here than in most codebases. This ticket's deliverable is
// explicitly "the reasoning beside each signature", and that reasoning is
// load-bearing precisely because it points at the ratified document a decision
// came from. A citation that does not resolve is a decision whose provenance
// has quietly been lost, in the one place the product keeps its provenance.
//
// Two citation styles are checked, because the comments use both: a full path
// from the repository root, and the bare filename the blueprint and screen
// specs are referred to by. The bare form resolves against a basename index of
// docs/, which is looser than a path check and still catches the failure that
// actually happens — a document renamed or never written.
func TestEveryDocReferenceResolves(t *testing.T) {
	root := moduleRoot(t)
	docs := filepath.Join(root, "docs")

	// A basename index of everything under docs/, for the bare-filename form.
	byName := map[string]bool{}
	if err := filepath.WalkDir(docs, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			byName[d.Name()] = true
		}
		return nil
	}); err != nil {
		t.Fatalf("indexing docs/: %v", err)
	}
	if len(byName) == 0 {
		t.Fatal("docs/ indexed empty; every bare-filename assertion below would fail spuriously")
	}

	fullPath := regexp.MustCompile(`docs/[A-Za-z0-9._/-]+\.(?:md|toml|svg|json|css)`)
	bareName := regexp.MustCompile(`\b(?:Z-\d{2}|\d{2})-[a-z0-9-]+\.md\b`)

	var checked int
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
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		b, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		rel, _ := filepath.Rel(root, path)

		for _, m := range fullPath.FindAllString(string(b), -1) {
			checked++
			if _, serr := os.Stat(filepath.Join(root, m)); serr != nil {
				t.Errorf("%s cites %s, which does not exist.\n"+
					"A citation that does not resolve is a decision whose provenance has been\n"+
					"lost, in a deliverable whose value is that its reasoning is traceable.", rel, m)
			}
		}
		for _, m := range bareName.FindAllString(string(b), -1) {
			checked++
			if !byName[m] {
				t.Errorf("%s cites %s, and no file of that name exists under docs/.", rel, m)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the module: %v", err)
	}
	if checked == 0 {
		t.Fatal("no doc references were examined; the patterns are broken and this would pass vacuously")
	}
	t.Logf("%d doc references checked against %d files under docs/", checked, len(byName))
}

// TestTheGoDirectiveIsAResolvableToolchainVersion fails on a bare language
// version in go.mod.
//
// `go 1.24` is a language version. Go cannot resolve a toolchain from it —
// `GOTOOLCHAIN=go1.24 go env GOROOT` reports "go1.24 is a language version but
// not a toolchain version" — and anything that resolves the pinned toolchain
// before doing its job stops there. The review vehicle does exactly that, so
// the module never reached the review at all: the failure was upstream of
// every test in this package and invisible to all of them.
//
// It is the same shape as the defects this file already guards. A property the
// project depends on (Sprint 0 #10: "go.mod declares the module and a pinned
// Go toolchain version") held only because somebody typed the right thing, and
// held nowhere that could notice when they did not.
func TestTheGoDirectiveIsAResolvableToolchainVersion(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(moduleRoot(t), "go.mod"))
	if err != nil {
		t.Fatal(err)
	}

	directive := regexp.MustCompile(`(?m)^go\s+(\S+)\s*$`)
	m := directive.FindSubmatch(b)
	if m == nil {
		t.Fatal("go.mod has no go directive")
	}
	version := string(m[1])

	// 1.X.Y — three components. Two is a language version and does not
	// resolve; four is not a Go version at all.
	full := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !full.MatchString(version) {
		t.Fatalf("go.mod pins %q, which is a language version rather than a toolchain version.\n"+
			"Go resolves a toolchain only from a full 1.X.Y, so anything that resolves the pin\n"+
			"before running — the review vehicle, a reproducible build — fails before it reads\n"+
			"a line of code. Sprint 0 #10 requires a pinned toolchain version.", version)
	}
}
