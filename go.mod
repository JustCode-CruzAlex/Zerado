module github.com/JustCode-CruzAlex/Zerado

// A full toolchain version (1.X.Y), not a bare language version.
//
// `go 1.24` is a LANGUAGE version. Go cannot resolve a toolchain from it —
// `GOTOOLCHAIN=go1.24 go env GOROOT` fails with "go1.24 is a language version
// but not a toolchain version" — which stopped the review vehicle before it
// reached any code, and which Sprint 0 #10's acceptance criteria already
// required ("go.mod declares the module and a pinned Go toolchain version").
//
// 1.27.0 is the toolchain this project is developed and reviewed on, so it
// resolves locally with no download on every review and every build. A
// contributor on an older Go switches to it automatically — that is what the
// toolchain directive is for since 1.21 — so raising the floor costs one
// transparent download rather than a build failure.
//
// The CI workflow reads this line via go-version-file rather than repeating
// the number, so there is one place to change it and no drift.
go 1.27.0
