# Contributing to Zerado

Thanks for looking. Zerado is **pre-alpha** — the landing page and the design
system exist, the Go program does not yet. That shapes what is useful to work on,
so this document starts there rather than burying it under process.

Everyone taking part is expected to follow the
[Code of Conduct](CODE_OF_CONDUCT.md).

---

## What is useful right now

| Useful today | Not yet |
|---|---|
| The landing page under `site/` — accessibility, performance, browser bugs | The CLI/TUI: there is no Go code to change |
| Documentation: this file, the README, `docs/**` | Steam integration, price data, mood tagging |
| The design system and brand docs under `docs/brand/`, `docs/design/` | Anything that assumes a running binary |
| Reproducing and reporting a real defect on [zerado.app](https://zerado.app) | Large refactors of code that has not been written |

Two things are **settled and not open for a drive-by pull request**:

- **The landing page's copy, design and section order.** They were ratified and
  reviewed; the copy in `docs/content/landing-copy.md` is normative and the
  layout contract is `docs/design/blueprint.md`. A change there is a
  conversation first — open an issue.
- **The ratified product decisions** — that Zerado earns no revenue and carries
  no donate button, the English-page-with-Portuguese-words rule, and the "no
  phase is marked done until it is done" rule.

If you are not sure whether something is in that category, open a
[question issue](../../issues/new/choose) and ask. That costs a minute and saves
an afternoon.

## Picking something up

1. Find an issue — [`good first issue`](../../issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
   is the shallow end.
2. **Comment on it before you start** so two people do not build the same thing.
3. Open a **draft pull request early**. Early feedback beats a finished branch
   that went the wrong way.

If nothing is labelled and you have an idea, open an issue describing it before
writing code. An unsolicited large pull request is the one most likely to be
turned down, and it is nobody's favourite outcome.

## Setting up

### The site

```bash
cd site
npm ci
npm run build     # static output → site/dist/
npm run preview   # http://localhost:4321/
```

Node.js **22.12.0+** and npm 10+ — Astro 7 refuses to run on anything older,
and `package.json` declares it under `engines`. CI builds on Node 22 and 24.

`npm run check` runs Astro's type check, but it is **not** wired into CI and it
is not installable out of the box: it needs `@astrojs/check` and `typescript`,
which this project deliberately does not depend on. Install them locally if you
want it (`npm i -D @astrojs/check typescript`); do not commit them without
raising it first, since the page ships zero JavaScript and the dependency tree
is kept deliberately small.

> ⚠️ **Never verify a build by opening `site/dist/index.html` over `file://`.**
> Astro emits root-absolute asset paths; the page will render unstyled and
> font-less and look broken when it is fine. Always use `npm run preview`.

### The page invariants

Every ratified decision that can be checked against a built page is checked by
one script, with no dependencies:

```bash
cd site && npm ci && npm run build && cd ..
node scripts/check-page.mjs site/dist/index.html
```

It asserts zero client-side JavaScript, zero unexpected external requests, no
funding control, no placeholder text, and that the roadmap renders exactly four
status markers all reading *Planned*. Each assertion names the decision it
defends, and each one is mutation-tested — a guard that has never failed is not
a guard. CI runs the same script.

### The QA harness

The forensic scripts that produced
[`docs/qa/qa-report.md`](docs/qa/qa-report.md) are captured from the build
rather than maintained as a suite — nothing in CI runs them, and
[`docs/qa/harness/README.md`](docs/qa/harness/README.md) says what each one
answers. They still run:

```bash
cd docs/qa/harness
npm ci
npx playwright install chromium
# with `npm run preview` serving the site in another terminal:
node run-core.mjs
node run-content.mjs
```

### The Go build

There isn't one yet. When Phase 1 lands, `go build ./cmd/zerado` and the
toolchain version go here.

## Branch naming

Never commit to `main`. Branch from an up-to-date `main` and name the branch for
the issue it closes:

```
feature/<issue>-<short-slug>     feature/14-mood-filter
fix/<issue>-<short-slug>         fix/22-nav-overflow-375
docs/<issue>-<short-slug>        docs/31-steam-key-setup
chore/<issue>-<short-slug>       chore/9-bump-astro
```

Lowercase, hyphens, no spaces. If there is no issue yet, make one — the number is
how the work stays findable a year from now.

## Commit messages

```
<type>(<scope>): <imperative summary, ≤72 chars>

<body: what changed and, more usefully, why>

Refs #<issue>
```

`<type>` is one of `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.
`<scope>` is the area touched — `site`, `brand`, `docs`, `ci`, or later a Go
package.

```
fix(site): stop the maker line clipping below 400px

The line was set in a fixed-width container that did not collapse at the
xs breakpoint, cutting 58% of the text. Switched to the fluid container
the rest of the sections use.

Refs #22
```

Rules that are enforced, not suggested:

- **Imperative mood** — "add", not "added" or "adds".
- **No AI attribution.** No "Generated with", no `Co-Authored-By` for a model,
  no tool footers. The history is a professional record.
- **One logical change per commit.** Formatting-only churn goes in its own
  commit so review can skip it.

## Pull requests

Open it as a **draft** while you work, mark it ready when it is. A pull request
should carry:

- **What changed and why** — the why matters more; the diff already shows what.
- **The issue it closes** — `Closes #14`.
- **Evidence for anything visual.** A change to `site/**` needs screenshots at
  **375, 768, 1280 and 1920 px**. "Looks fine on my machine" is not evidence.
- **A note on accessibility** if you touched markup, colour or focus order.

### The review bar

A pull request is merged when it clears all of these:

1. **CI is green.** Two workflows gate a pull request: `site` builds the page on
   Node 22 and 24 and runs `scripts/check-page.mjs` against the result, and
   `guardrails` runs repository-wide checks on **every** pull request, whatever
   it touched. A broken `npm run build` fails the pull request by design.
2. **It does what it says**, and nothing else. Unrelated changes get split out.
3. **It does not regress the page's guarantees**, which are load-bearing and
   measured: zero client-side JavaScript, zero external network requests, zero
   raster images in the render path, Lighthouse ≥ 95 in every category, and
   **zero axe-core violations** at all four breakpoints.
4. **Colour is never the only signal.** Every state is co-rendered as
   colour + glyph + label. This is a hard rule and it applies to the site and to
   the future TUI equally.
5. **Nothing dishonest ships.** No roadmap phase is marked done before it is
   done, no mockup is presented as a screenshot, and no capability is described
   in the present tense before it exists. This is the rule most likely to fail a
   pull request that is otherwise good.
6. **No placeholder text** — no `lorem ipsum`, no `TODO`, no `[TBD]` in shipped
   output.

Once Go code exists, tests come with the change: the tests are written with the
feature, not after it, and coverage does not go down.

## Reporting a bug

Use the [bug report template](../../issues/new/choose). What actually makes a
report actionable: the URL, the browser and version, the viewport width, what you
expected, what happened, and a screenshot. A defect at a specific breakpoint is
almost always reproducible if you say which breakpoint.

**Do not report security vulnerabilities in the issue tracker.**
See [SECURITY.md](SECURITY.md).

## Licence

By contributing, you agree that your contributions are licensed under the
[MIT Licence](LICENSE) that covers this repository.
