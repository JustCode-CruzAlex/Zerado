# QA harness — captured evidence, re-runnable

These are the scripts that produced [`../qa-report.md`](../qa-report.md),
[`../qa-crossbrowser.json`](../qa-crossbrowser.json) and the screenshots. They
are **forensic scripts from the build**, not a maintained test suite: they were
written to answer specific questions while the page was being built, and they
are kept because the evidence they produced is only meaningful if the method is
inspectable.

**Nothing in CI runs them.** The gate that runs on every pull request is
[`../../../scripts/check-page.mjs`](../../../scripts/check-page.mjs) — 18
mutation-tested invariants. Reach for that first; reach for these when you need
to re-derive a measurement in the QA report.

## Running one

Every script hardcodes `http://localhost:4321/`, so serve the built site first:

```bash
cd ../../../site && npm ci && npm run build && npm run preview &
cd -                       # back to docs/qa/harness
npm ci
npx playwright install chromium
node run-core.mjs
```

To run against the live site instead, edit the URL at the top of the script.

## What each one answers

| Script | Question it was written to answer |
|---|---|
| `run-core.mjs` | axe-core violations + per-element overflow at four viewports |
| `run-content.mjs` | Forbidden tokens, roadmap statuses, money controls, link/image inventory |
| `run-kbd.mjs` | Keyboard traversal and focus order |
| `run-rm.mjs` | `prefers-reduced-motion` — does the scanner park rather than hide |
| `run-xb.mjs` · `xb-maker.mjs` | Cross-browser reproduction (Firefox, WebKit) |
| `contrast.mjs` | Measured contrast of every text node. Takes a viewport: `node contrast.mjs 375 812` |
| `corender.mjs` | The co-render rule — colour never the only signal |
| `links.mjs` | Every link resolves; `mailto:` subjects exact |
| `targets.mjs` | Touch-target sizes |
| `nav375.mjs` | Nav layout at the narrowest breakpoint |
| `skip.mjs` · `skip2.mjs` | Skip-link visibility and focus behaviour |
| `motion2.mjs` | Scanner-sweep animation state |
| `mk.mjs` · `mk2.mjs` | Mockup-caption presence on every figure |
| `pip.mjs` | Phone-frame rendering |
| `final.mjs` | The end-of-build sweep |
| `components.mjs` | Component inventory against the blueprint |

## Known limitations

- **Every check here is programmatic.** Zero axe violations is a real result and
  it is not the same claim as "a screen-reader user has a good time". No testing
  with real assistive technology has been done.
- The names are the ones they were written under. They are not descriptive, and
  renaming them now would break the traceability from `qa-report.md`.
