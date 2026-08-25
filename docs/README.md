# Zerado — documentation

Everything the landing page was built from, and the evidence that it works.
The page is reproducible from this repository alone: build it with the commands
in [`../site/README.md`](../site/README.md).

## Brand — the identity

| | |
|---|---|
| [`brand/brand-manual.md`](brand/brand-manual.md) | The flagship document. Logo construction and misuse, the colour system in hex **and ANSI-256** (the palette has to survive in a terminal), the measured contrast table, the type system, the scanner-sweep motion spec, voice and tone. |
| [`brand/naming.md`](brand/naming.md) | Why the name is Zerado, the two runners-up, the risks accepted — and the **casing convention** (`Zerado` the product · `zerado` the command · *zerado* the status · `ZERADO` the chip). |
| [`brand/tokens.css`](brand/tokens.css) · [`brand/tokens.json`](brand/tokens.json) | 200 design tokens. The source of truth for the site and for the branded PDFs. |
| [`brand/tokens-print.css`](brand/tokens-print.css) | The same tokens plus the three faces embedded as base64, for PDF rendering. Do not use it for the site. |
| `brand/logo*.svg` · `brand/favicon.svg` | The marks. |

## Content — the words

| | |
|---|---|
| [`content/landing-copy.md`](content/landing-copy.md) | Every section's final copy, wired verbatim into the page. **Normative** — the page implements this, not the other way round. |
| [`content/seo.md`](content/seo.md) | Title, meta, OG/Twitter, JSON-LD. No fabricated ratings, prices or version. |
| [`content/voice-and-tone.md`](content/voice-and-tone.md) | How Zerado talks, with do/don't pairs. |

## Design — the contract

| | |
|---|---|
| [`design/blueprint.md`](design/blueprint.md) | Per-section layout at four breakpoints, component inventory, applied tokens, motion inventory. |
| [`design/blueprint.tokens.json`](design/blueprint.tokens.json) | The machine-readable half — the assertions the QA harness checks the built page against. |

## Evidence

| | |
|---|---|
| [`qa/qa-report.md`](qa/qa-report.md) | axe-core **0 violations** at 375/768/1280/1920, 0 contrast failures across 284 text nodes, cross-browser results. |
| [`qa/screenshots/`](qa/screenshots) | 30 captures — four viewports × three engines, plus reduced-motion, focus and defect-history shots. |
| [`qa/harness/`](qa/harness) | The checks, **re-runnable**. See [`../CONTRIBUTING.md`](../CONTRIBUTING.md). |
| [`performance/performance-report.md`](performance/performance-report.md) | Lighthouse **100/100** mobile and desktop · LCP 1504 ms / 366 ms · **CLS 0.000**, over six runs. |
| [`review/review.md`](review/review.md) | The verdict: **GOLDEN**, 0 blocking, 0 major. |

Every accessibility result here is **programmatic**. Zero axe violations is a
real result, and it is not the same claim as "a screen-reader user has a good
time" — real assistive-technology testing by a human has not been done.

## Deploying

| | |
|---|---|
| [`deploy/digitalocean.md`](deploy/digitalocean.md) | The runbook: confirmed price, hosting decision and its trade-offs, DNS, TLS, the verification checklist. |
| [`deploy/nginx/zerado.app.conf`](deploy/nginx/zerado.app.conf) | The fallback path's nginx config — ports `_headers` verbatim, adds HSTS and a CSP. |
| [`../.do/app.yaml`](../.do/app.yaml) | The App Platform spec for the recommended path. |

## PDFs

[`pdf/`](pdf) carries nine of the documents above as branded engineering sheets,
each a framed sheet with its own title block, DWG number, revision and legend.
Regenerate with:

```bash
flowforge docs pdf --out pdf --brand-css brand/tokens-print.css <files…>
```

Note `tokens-print.css`, not `tokens.css`: `tokens.css` deliberately ships no
`@font-face`, which is correct for the site but would silently fall the PDFs
back to a system typeface.

## One more thing

[`REDACTIONS.md`](REDACTIONS.md) records what was changed on the way from the
private deliverable into this public repository, and what was left out. Nothing
about the page itself was changed.
