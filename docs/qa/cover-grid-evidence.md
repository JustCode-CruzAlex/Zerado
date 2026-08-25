---
title: Zerado — §06 cover grid, measured evidence
discipline: QA
doc-no: ZRD-QA-02
rev: A
date: 2026-08-25
ticket: "#16"
---

# §06 cover grid — what was measured, and how to re-measure it

Reproduce any number below with the committed harness:

```bash
cd site && npm ci && npm run build
node scripts/check-page.mjs site/dist/index.html
cd docs/qa/harness && npm install && npx playwright install chromium
node covers-evidence.mjs ../../../site/dist withcovers /tmp/shots
```

Everything here was measured against a local static server on 2026-08-25. **The
live-URL re-run is still owed** and is noted in §5 — `zerado.app` serves the
pre-change page until this merges and deploys, so a "live" number today would be
a measurement of the old page.

---

## 1 · The two states, and why both are measured

The grid renders real IGDB covers when the files are in the build and the
ratified art-directed tiles when they are not
(`docs/legal/igdb-image-licence-finding.md` §7). Both are shipping states, so
both are measured and `scripts/check-page.mjs` asserts all 27 invariants in each.

The with-covers figures were produced with **synthetic stand-in images at the
real dimensions and a realistic compression profile** (360 × 480, AVIF q55 ≈
22 KB each), because the IGDB API credentials had not arrived when this was
written. They are indicative of weight and exact for geometry, markup, ARIA and
layout stability. **They are not real cover art and none of them were
committed.** §5 records what is re-run when the real files land.

---

## 2 · Lighthouse

Lighthouse 12, headless Chrome, local static server.

| | Performance | Accessibility | Best practices | SEO | CLS | LCP |
|---|---|---|---|---|---|---|
| **Desktop · with covers** | **100** | **100** | **100** | **100** | **0** | 0.5 s |
| **Mobile · with covers** | **99** | **100** | **100** | **100** | **0** | 2.0 s |
| **Mobile · fallback (no covers)** | 99 | 100 | 100 | 100 | 0 | 2.0 s |

**Every category is ≥ 95 on both form factors, and CLS is 0.**

The mobile 99 is **not caused by the covers**: the fallback build — byte-identical
to what ships today in the grid — scores 99 with the same 2.0 s LCP under the same
run conditions. It is this machine against the published 100, which was measured
on the live host. The covers are performance-neutral here, which is what
lazy-loading twelve below-the-fold images is supposed to buy: Lighthouse's
`offscreen-images` audit passes and total mobile byte weight is 176 KiB, because
the covers are never fetched during the run at all.

---

## 3 · axe-core, four viewports

`@axe-core/playwright`, tags `wcag2a wcag2aa wcag21a wcag21aa`, after a full
scroll so lazy content is realised.

| Viewport | axe violations (with covers) | axe violations (fallback) | External requests |
|---|---|---|---|
| 375 | **0** | **0** | **0** |
| 768 | **0** | **0** | **0** |
| 1280 | **0** | **0** | **0** |
| 1920 | **0** | **0** | **0** |

**Zero external requests at every viewport in both states** — the covers are
served from this origin, so the page guarantee is intact.

---

## 4 · Layout stability and the aspect-ratio policy

Measured tile geometry, `getBoundingClientRect` on all twelve children:

| Viewport | Rows | Tile width | Tile height | Ratio |
|---|---|---|---|---|
| 375 | 3 · 3 · 3 · 3 | 104 px | 138 px | 0.754 |
| 768 | 4 · 4 · 4 | 168 px | 224 px | 0.750 |
| 1280 | 6 · 6 | 175 px | 234 px | 0.748 |
| 1920 | 6 · 6 | 175 px | 234 px | 0.748 |

One width and one height per viewport, and rows that always divide evenly —
**no ragged row at any breakpoint**, which was the named risk. This is the
ratified geometry unchanged; the covers are cropped to it at fetch time rather
than the grid bending to them.

**CLS, observed over a full scripted scroll** (`layout-shift` PerformanceObserver,
excluding shifts with recent input):

- with covers: `6.98e-6`
- fallback: `4.19e-6`

Both are **0.000** to three decimals. The 2.8e-6 delta is far below anything
reportable, and the residual in *both* states is the pre-existing animated
scanner, not the grid.

**Why the covers cannot shift anything:** each tile reserves its box with
`aspect-ratio: 3 / 4` before any image loads, and every `<img>` carries
`width="360" height="480"` — its true intrinsic size, because the file is written
at exactly those dimensions.

---

## 5 · What is still owed

1. **Re-run against `https://zerado.app` after this deploys.** Today the live URL
   serves the pre-change page.
2. **Re-run §2–§4 with the real IGDB covers** once credentials land. Geometry,
   markup, ARIA and CLS are structural and will not move; the only figure that
   can change is byte weight, and only if real cover art compresses worse than
   the ≈22 KB AVIF stand-in — which would be unusual for flat-ish box art.

Both re-runs are one command each, above.
