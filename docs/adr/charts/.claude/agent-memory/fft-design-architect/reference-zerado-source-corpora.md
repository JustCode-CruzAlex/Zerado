---
name: reference-zerado-source-corpora
description: Absolute paths to the Zerado brand manual/tokens and the FlowForge omarchy theme corpus — both live OUTSIDE the zerado-2 repo.
metadata:
  type: reference
---

Zerado design work depends on two corpora that are **not in the `zerado-2` repo**. Every design
document cites them, so they must be read at source before asserting any figure.

- **Zerado brand manual + tokens** —
  `/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/brand/`
  (`brand-manual.md`, `naming.md`, `tokens.css`, `tokens.json`). §4.2/§4.3 hold the measured dark
  contrast table, §4.4 the CVD method, §4.5 the paper expression, §5.1–§5.4 the terminal
  representations. The **light/paper state set lives in `tokens.css` §10**
  (`[data-z-surface="paper"]`), not in the manual body.
- **FlowForge omarchy theme corpus** — `.../flowforge/FlowForge/v3/uikit/theme/themes/omarchy/`
  (35 `.toml` + `ATTRIBUTION.md`). The theme machinery is the sibling package: `omarchy.go`
  (palette→Theme, `fillFromSemantic`), `cvd.go`, `registry.go` (the #5369 repair-at-registration /
  audit-at-activation gate), `themelib/import.go`.

The `flowforge/` directory holds **dozens of per-ticket worktrees** (`FlowForge-5622`,
`FF-rev5525`, …). The canonical checkout is the plain `FlowForge/` one; there is no `flowforge/v3`.

Zerado **cannot import** any of this Go — the module does not resolve anonymously and much is
`internal/`. It is inherited as *specification* only (`00-design-brief.md` §7).

Related: [[project-zerado-theme-system]], [[feedback-verify-at-source-not-the-brief]]
