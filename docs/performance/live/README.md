# Lighthouse — the live site

Run against **`https://zerado.app/`** after deploy, not against a local build.
The pre-deploy baseline is [`../performance-report.md`](../performance-report.md);
this is the hosting layer's own result.

| Preset | Performance | Accessibility | Best practices | SEO |
|---|---|---|---|---|
| **desktop** | 100 | 100 | 100 | 100 |
| **mobile** | 100 | 100 | 100 | 100 |

## Core Web Vitals

| Preset | LCP | CLS | TBT | FCP | Speed Index |
|---|---|---|---|---|---|
| **desktop** | 0.5 s | 0 | 0 ms | 0.4 s | 0.9 s |
| **mobile** | 1.1 s | 0 | 0 ms | 0.9 s | 2.1 s |

Captured 2026-08-25. Raw reports: `lh-desktop-live.report.{json,html}`,
`lh-mobile-live.report.{json,html}`.

## What moved between the first live run and this one

The first post-deploy run scored **99 performance / 99 accessibility** on desktop.
Both gaps were real and both were fixed rather than accepted:

- **Accessibility 99** — `aria-allowed-role` failed on six `<figure role="img">`
  elements. ARIA-in-HTML does not permit that role on `<figure>`. The role moved
  to a valid wrapper `<div>`; violations went 6 → 0.
- **Performance 99** — attributable to the run's own variance plus 2.6 KB of HTML
  comments that were being emitted into the shipped page. The comments are now
  Astro comments, which compile away.

## The one audit that is deliberately not perfect

`uses-long-cache-ttl` flags `/logo.svg` at a 7-day cache rather than a year. That
is the ratified `_headers` policy, not an oversight: only **content-hashed** paths
(`/_astro/*`) may take a year-long `immutable` cache, because a file regenerated
under an unchanged name would otherwise be served stale for as long as the cache
lives. See `site/public/_headers` and `../performance-report.md` F5.

## Honest limits

Every result here is programmatic, from one machine, on one network. Zero axe
violations is a real result and is not the same claim as "a screen-reader user
has a good time" — no testing with real assistive technology has been done.
