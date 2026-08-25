---
title: Zerado — §06 cover grid, measured evidence
discipline: QA
doc-no: ZRD-QA-02
rev: B
date: 2026-08-25
ticket: "#16"
---

# §06 cover grid — what was measured, and how to re-measure it

Reproduce any number below with the committed harness:

```bash
IGDB_CLIENT_ID=… IGDB_CLIENT_SECRET=… node scripts/fetch-covers.mjs
cd site && npm ci && npm run build
node scripts/check-page.mjs site/dist/index.html
node --test "scripts/**/*.test.mjs"
cd docs/qa/harness && npm install && npx playwright install chromium
node covers-evidence.mjs ../../../site/dist withcovers /tmp/shots
```

Everything here was measured against a local static server on 2026-08-25. **The
live-URL re-run is still owed** and is noted in §5 — `zerado.app` serves the
pre-change page until this merges and deploys, so a "live" number today would be
a measurement of the old page.

> **Rev B (2026-08-25).** Rev A's with-covers figures were produced with
> synthetic stand-ins, because the IGDB credentials had not arrived. They have
> now arrived, **all twelve real covers are fetched and committed**, and §2–§4
> below are re-measured against them — which was re-run item 2 of rev A's §5.
> Rev A's prediction held: geometry, markup, ARIA and CLS did not move, and real
> box art compresses *better* than the stand-in (19.6 KiB average AVIF against
> the ≈22 KiB assumed). Every number below is now a measurement of the real
> thing. Item 1 — the live-URL re-run — is still owed and still blocked on the
> same thing.

---

## 1 · The two states, and why both are measured

The grid renders real IGDB covers when the files are in the build and the
ratified art-directed tiles when they are not
(`docs/legal/igdb-image-licence-finding.md` §7). Both are shipping states, so
both are measured and `scripts/check-page.mjs` asserts all 27 invariants in each.

**The with-covers figures below are the real IGDB covers**, fetched by
`scripts/fetch-covers.mjs` and committed at `site/public/covers/`. What each
file came from — IGDB game id, name, release year, cover image id and source
URL — is recorded per row in `docs/legal/cover-provenance.json`.

All twelve resolved to an unambiguous match. Two were near-misses that the
match rule resolved rather than guessed at, and both are now regression tests in
`scripts/fetch-covers.test.mjs`:

- **Dead Space** — the row asserts 2008, and IGDB carries the 2008 original and
  the 2023 remake under the same name. The release-year pin picked the original;
  the shipped file is the severed-arm cover, not the remake's.
- **Baldur's Gate II** — `coverGrid.ts` sets the title with a typographic
  apostrophe (U+2019) because that is how it reads on the page; IGDB stores an
  ASCII one. Byte-equal name matching therefore found no exact match, fell
  through to the unnamed pool, drew the base game **and** the Collectors'
  Edition on year, and correctly **refused** — the first run shipped 11 of 12.
  The fix normalises typographic quotes before comparing, which removes that one
  false distinction and nothing else; the year pin is untouched and the rule
  still refuses when two candidates are genuinely indistinguishable.

The other state is unchanged and still measured: the grid renders the ratified
art-directed tiles when the files are not in the build
(`docs/legal/igdb-image-licence-finding.md` §7). Both are shipping states, so
both are measured and `scripts/check-page.mjs` asserts all 27 invariants in each.

---

## 2 · Lighthouse

Lighthouse 13.4.1, headless Chrome, `astro preview` on this machine. Every run
below is the same build, the same server and the same flags; the only variable is
whether `site/public/covers/` is populated.

| | Performance | Accessibility | Best practices | SEO | CLS | LCP | Page weight |
|---|---|---|---|---|---|---|---|
| **Desktop · real covers** | **100** | **96** | **100** | **100** | **0** | 0.4 s | 94 KiB |
| **Mobile · real covers** | **100** | **96** | **100** | **100** | **0** | 1.5 s | 94 KiB |
| Desktop · fallback (no covers) | 100 | 96 | 100 | 100 | 0 | 0.4 s | 94 KiB |
| Mobile · fallback (no covers) | 100 | 96 | 100 | 100 | 0 | 1.5 s | 94 KiB |

**Every category is ≥ 95 on both form factors, and CLS is 0.**

**The covers are performance-neutral — the four rows are identical.** That is
what lazy-loading twelve below-the-fold images is supposed to buy: Lighthouse
never scrolls to §06, so the covers are never fetched during the run and the page
weight is the same 94 KiB with and without them. `unsized-images` passes, which
is the audit behind the CLS-is-0 result.

**The accessibility 96 is pre-existing and is not this grid.** The single failing
audit is `color-contrast`, and its one node is `.z-frame__watermark` — the
`PHASE 4` watermark inside a §14 terminal frame, which ticket #16 does not touch.
The fallback rows above score the identical 96 with the identical node, which is
the proof it did not arrive with the covers. It is also not a real barrier: the
watermark sits inside `aria-hidden="true"`, and axe-core run directly (§3) reports
**zero** violations at all four viewports for exactly that reason. Lighthouse 13
scores the node anyway. Recorded here rather than fixed, because the rev-A run on
Lighthouse 12 reported 100 and a reviewer comparing the two documents deserves to
know which of the two changed — the tool did, not the page.

---

## 3 · axe-core, four viewports

`@axe-core/playwright`, tags `wcag2a wcag2aa wcag21a wcag21aa`, after a full
scroll so lazy content is realised.

| Viewport | axe violations (real covers) | axe violations (fallback) | External requests |
|---|---|---|---|
| 375 | **0** | **0** | **0** |
| 768 | **0** | **0** | **0** |
| 1280 | **0** | **0** | **0** |
| 1920 | **0** | **0** | **0** |

**Zero external requests at every viewport in both states** — every cover is
served from this origin. Nothing hotlinks IGDB's CDN, which is what keeps the
page's zero-external-request guarantee true now that it carries third-party
imagery.

**Shipped byte weight**, three formats per cover, 360 × 480 each:

| Format | 12 covers | Average | Served to |
|---|---|---|---|
| AVIF | 235.7 KiB | 19.6 KiB | every current browser |
| WebP | 365.2 KiB | 30.4 KiB | the fallback before that |
| JPEG | 412.3 KiB | 34.4 KiB | the `<img>` of last resort |

A visitor downloads **one** of the three per tile, and only on scrolling to §06 —
≈236 KiB of AVIF for the whole grid, below the ≈264 KiB rev A projected.

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
**no ragged row at any breakpoint**, which was the named risk. **These numbers
are unchanged from rev A**, measured now against the real covers: the composition
is the ratified one, untouched. The covers are cropped to the grid at fetch time
rather than the grid bending to them, so the twelve source images — which arrive
from IGDB at 0.706 and would otherwise have gone ragged — cannot move it.

**CLS, observed over a full scripted scroll** (`layout-shift` PerformanceObserver,
excluding shifts with recent input):

- real covers: `5.82e-6`
- fallback: `6.61e-6`

Both are **0.000** to three decimals, and Lighthouse reports a flat `0` for CLS
in all four runs of §2. The real covers measure *below* the fallback, which is
the tell that the residual in both states is the pre-existing animated scanner
and not the grid — twelve images cannot make a page more stable than no images,
so the difference is run-to-run noise on a number six orders of magnitude below
the 0.1 "good" threshold.

**What it actually looks like**, captured from the same runs — the section, not
the fold, because a viewport shot of the fold is evidence of nothing here:

- [`screenshots/cover-grid-real-1280.png`](screenshots/cover-grid-real-1280.png) — 6 · 6
- [`screenshots/cover-grid-real-375.png`](screenshots/cover-grid-real-375.png) — 3 · 3 · 3 · 3

Both show twelve real covers, one uniform tile size, the platform plates legible
over the art, and the caption reading *"Cover art from IGDB.com"* — which the
build derived from the files being present, not from a flag.

**Why the covers cannot shift anything:** each tile reserves its box with
`aspect-ratio: 3 / 4` before any image loads, and every `<img>` carries
`width="360" height="480"` — its true intrinsic size, because the file is written
at exactly those dimensions.

---

## 5 · What is still owed

1. **Re-run against `https://zerado.app` after this deploys.** Still owed, and
   still for the same reason. It is now a two-step wait rather than one: the site
   deploys from `main`, and `main` is at the first cut — the covers reach the
   public page only after the release line is cut into `main` and redeployed.
   Until then the live URL serves the pre-change page, and any "live" number
   would be a measurement of that page rather than of this change.

2. ~~**Re-run §2–§4 with the real IGDB covers** once credentials land.~~
   **Done, rev B.** All twelve fetched, committed and re-measured above.

One command, above. Nothing else is outstanding.

---

## 6 · The pre-existing finding this run surfaced

Not this ticket's, and recorded so it is not lost:

**`.z-frame__watermark` fails Lighthouse 13's `color-contrast` audit** — the
`PHASE 4` watermark inside a terminal frame. It is present identically in the
no-covers build, so it did not arrive with this change; it is inside
`aria-hidden="true"`, so axe-core does not count it and the page still measures
zero WCAG violations at four viewports; and accessibility scores 96, above the
≥ 95 bar. Left alone deliberately: the terminal frames are ratified mockup
surfaces this ticket is explicitly not to touch, and raising the watermark's
contrast is a design decision, not a repair.
