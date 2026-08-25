# Zerado — landing page

A fully static, self-contained landing page for **Zerado**, built with
[Astro](https://astro.build). Zero client-side JavaScript, zero raster
images, zero external network requests — everything on the page is
same-origin HTML, CSS and self-hosted fonts/SVG.

## Requirements

- Node.js **22.12.0+** — required by Astro 7 (built and tested on Node v25.4.0)
- npm 10+ (tested on npm 11.7.0)

The requirement is declared in `package.json` under `engines`, and CI builds on
Node 22 and 24.

## Install

```bash
npm install
```

## Build

```bash
npm run build
```

Output goes to `dist/`. The build is fully static — `dist/` is the entire
deliverable and can be uploaded as-is to any static host (S3 + CloudFront,
Netlify, Vercel static, nginx, etc.).

## Preview the production build locally

```bash
npm run preview
```

This starts a local static file server (default `http://localhost:4321`)
serving the actual contents of `dist/`, with correct MIME types and
root-absolute asset resolution — this is the command to use to verify the
build before shipping, not opening the HTML file directly (see warning
below).

## ⚠️ Do not open `dist/index.html` directly in a browser

Astro emits **root-absolute** asset paths (`/logo.svg`, `/fonts.css`,
`/_astro/…`, etc.). Opening `dist/index.html` via `file://` will load with
**no styling and no fonts**, because the browser resolves `/logo.svg`
against the filesystem root, not against `dist/`. This is expected Astro
behaviour, not a bug in the build.

**Always serve the build over HTTP** instead:

```bash
npm run preview
# then open http://localhost:4321/
```

or, for a quick one-off check without Astro's own server:

```bash
npx serve dist
```

## Deploy

`dist/` is a plain static site — deploy it to any static host or object
store with a web server in front of it (CloudFront, nginx, Netlify, Vercel,
GitHub Pages, etc.). No server runtime, database, or environment variables
are required; there is no backend.

If deploying to a path other than the domain root, note that
`astro.config.mjs` sets `site: 'https://zerado.app'` for canonical/OG URLs —
update that value (and the canonical/OG URLs hard-coded in
`src/layouts/Base.astro`) if the production domain changes.

### Cache policy

`dist/` ships a `public/_headers` file (copied verbatim into `dist/_headers`
by the build) with the recommended `Cache-Control` policy — long-lived
`immutable` caching only for the content-hashed `/_astro/*` bundle, a
revalidating `stale-while-revalidate` policy for the **unhashed** font and
SVG files (a regenerated file under the same name must not be served stale
for a year), and `no-cache` on `index.html` so a new deploy is picked up
immediately.

**Only Netlify and Cloudflare Pages read `_headers` natively.** Other hosts
need the equivalent policy configured at the platform layer instead — the
file's rules are the source of truth to port:

- **Vercel** — add a `headers` array to `vercel.json` (Vercel does not read
  a `_headers` file), one entry per path pattern in `public/_headers`.
- **S3 + CloudFront** — S3 does not evaluate path-glob header rules at all;
  either set per-object `Cache-Control` metadata at upload time (e.g. via
  the AWS CLI `--cache-control` flag, one `aws s3 cp` per path group) or
  attach a CloudFront Function / Lambda@Edge that sets the header by path on
  the way out. A CloudFront distribution in front of either approach is
  still required to get brotli and HTTP/2 (see `../performance/performance-report.md`
  F3) — S3 alone serves neither.
- **nginx** — translate each block into `location` + `add_header
  Cache-Control ...;` directives in the server config.

Whichever host is chosen, also enable brotli compression for `text/html`,
`text/css`, `image/svg+xml` and `application/ld+json` — and explicitly
**exclude** `font/woff2`, which is already Brotli-compressed internally.
Nothing in `dist/` sets this; it is a host/CDN-level configuration step.

## Project structure

```
site/
├── astro.config.mjs
├── package.json
├── public/                  # same-origin static assets (fonts, SVG marks)
└── src/
    ├── layouts/
    │   └── Base.astro       # document shell: head/SEO, skip link, nav, main
    ├── components/          # reusable widgets (buttons, cards, chips, …)
    ├── sections/             # one file per ratified page section (§01–§15)
    ├── data/                 # small typed data tables (states, roadmap, …)
    ├── styles/
    │   ├── tokens.css        # brand design tokens, verbatim from brand/tokens.css
    │   ├── reset.css
    │   ├── base.css          # page-level tokens, typography, containers
    │   ├── motion.css        # scanner-sweep keyframes + reduced-motion collapse
    │   └── global.css        # import entry point
    └── pages/
        └── index.astro       # assembles all 16 sections in ratified order
```

## Design source of truth

This build implements `../design/blueprint.md` and
`../design/blueprint.tokens.json` against the ratified copy in
`../content/landing-copy.md` and `../content/seo.md`, using the brand tokens
in `../brand/tokens.css`. See those files for the normative spec this page
was built against.
