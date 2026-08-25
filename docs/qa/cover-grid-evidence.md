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
node --test scripts/*.test.mjs
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

## 6 · What shipping the real images surfaced

Three findings. Two were claims that were harmless while the grid was twelve CSS
gradients and became load-bearing the moment real third-party artwork shipped —
both fixed here, because this change is what made them matter. The third is not
this ticket's and is recorded rather than fixed.

### 6.1 · The rights-holder withdrawal path was blocked by the build gate — FIXED (at the second attempt)

`docs/legal/igdb-image-licence-finding.md` §7 names, as a reversal path, "a
rights-holder objecting to a specific cover, in which case that row loses its
image and the tile treatment renders in its place — **a data-only change**."

**It was not a data-only change.** Deleting one cover's files and rebuilding
failed `check-page.mjs`: the assertion handled the all-absent case but not the
partial one, which is exactly the rights-holder case. The two directions were
being treated as one, and they are not symmetric — a rendered alt the copy file
never declared is always drift and always fails, while a declared alt the page
did not render is a failure only if that row's cover actually shipped.

**The first fix was incomplete, and the verification is why it got through.**
The same commit that unblocked `check-page.mjs` added
`scripts/fetch-covers.test.mjs` to the *same* `site.yml` job, six lines below,
carrying `assert.equal(shipped.length, 12)` — which skips at 0, passes at 12 and
**fails at 1–11**. 1–11 is the rights-holder case. Both steps run in the same
`build` job, so the job still went red: the block had moved, not lifted. The
four-direction table published with that commit checked the assertion that had
changed and not the CI job that had changed alongside it, and reported the path
as restored when it was not.

The correction: **`check-page.mjs` is the single authority on partiality.** Two
gates adjudicating "how many covers may ship" by different rules is what
produced the defect; the test file now confines itself to per-file integrity and
takes no view on the count.

**Re-verified as the JOB — both gates, the pair CI actually runs — not one:**

| # | State | check-page | `node --test` | Job | Expected |
|---|---|---|---|---|---|
| 1 | All twelve ship | PASS | PASS | 🟢 GREEN | GREEN ✅ |
| 2 | One cover withdrawn | PASS | PASS | 🟢 GREEN | GREEN ✅ |
| 3 | Cover file present, row not rendered | FAIL | FAIL | 🔴 RED | RED ✅ |
| 4 | Page renders an undeclared alt | FAIL | PASS | 🔴 RED | RED ✅ |
| 5 | Shipped cover missing from provenance | PASS | FAIL | 🔴 RED | RED ✅ |
| 6 | Restored | PASS | PASS | 🟢 GREEN | GREEN ✅ |

Row 2 is the remedy; rows 3–5 are the three ways this could go wrong silently,
and each is still caught. The excuse in `check-page.mjs` resolves alt → slug
with a regex assuming `slug:` sits immediately before `alt:`; reorder those two
fields and the map returns empty, nothing is excused, and the path silently
re-blocks — pinned by a test rather than left to be rediscovered by whoever
needs the remedy in a hurry.

### 6.2 · The licence finding pointed at the wrong file for provenance — FIXED

§8, *Provenance of every image on the page*, said the IGDB game slug and cover
image hash were "recorded per row in `site/src/data/coverGrid.ts`". **They are
not, and never were** — that file has no such fields; it carries the row's
*identity* (the name and release year the fetch must match) and its platform tag.

This is the section a rights-holder question gets answered from, so it stopped
being harmless the moment twelve real covers shipped. §8 now points at
`docs/legal/cover-provenance.json` — IGDB game id, slug, name, release year,
cover image id and exact source URL, per file.

**And the record could be truncated, which §8 asserts cannot happen.**
`fetch-covers.mjs` rebuilt `provenance.covers` from scratch every run and pushed
only rows that resolved, while a failed row's previously-downloaded files were
left on disk untouched — the "pin them by hand" case the script itself prints.
A partial re-run therefore dropped rows from the legal record while their covers
stayed shipped, breaking the one invariant §8 states: *every file under
`site/public/covers/` appears there*. True in the committed state, unenforced in
the code.

Fixed at the source — rows that fail a run but whose files are still shipped are
carried forward — and now asserted rather than trusted (`every shipped cover has
a provenance row`, one direction only: a provenance row with **no** file is a
withdrawn cover, the sanctioned §7 state, not a defect).

Verified end to end, not just by its guard: one row was made deliberately
unresolvable and a real fetch run against the live IGDB API was allowed to fail
on it. The run reported `↻ 1 row(s) carried forward from the previous record
(files still shipped)`, exited 1 as it should, and left the record at all twelve
rows with the failed row's entry intact. Row 5 of the table above is the same
invariant checked from the failing side.

### 6.3 · A pre-existing contrast node — recorded, not fixed

**`.z-frame__watermark` fails Lighthouse 13's `color-contrast` audit** — the
`PHASE 4` watermark inside a terminal frame. It is present identically in the
no-covers build, so it did not arrive with this change; it is inside
`aria-hidden="true"`, so axe-core does not count it and the page still measures
zero WCAG violations at four viewports; and accessibility scores 96, above the
≥ 95 bar. Left alone deliberately: the terminal frames are ratified mockup
surfaces this ticket is explicitly not to touch, and raising the watermark's
contrast is a design decision, not a repair.
