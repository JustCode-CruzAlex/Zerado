# Zerado landing page — QA report (accessibility & cross-browser lane)

**Verdict: FAIL — 1 blocking finding.**

| | |
|---|---|
| Artifact | `http://localhost:4321/` (Astro preview of `site/dist/index.html`, 52,880 bytes) |
| Lane | accessibility + cross-browser. **Core Web Vitals / bundle size deliberately NOT run** — that is the `performance` lane's scope. |
| Tooling | Playwright 1.57 (`qa/harness/`), Chromium 151.0.7922.34 · Firefox 153.0 · WebKit 26.5; `@axe-core/playwright` + axe-core 4.11 |
| Raw data | `qa/qa-raw-results.json` (axe + overflow, 4 viewports) · `qa/qa-crossbrowser.json` (3 engines × 4 viewports) |
| Screenshots | `qa/screenshots/` — 30 files |

**Findings: 1 BLOCKING · 2 MAJOR · 7 MINOR.**

The page is, with one exception, unusually clean: **zero axe violations at all four viewports**, zero
contrast failures across 206 rendered text elements, a correct skip link, a complete 16-section DOM in
ratified order, a fully honest content sweep, and near-identical rendering across three engines. The single
blocking defect is a flex-sizing bug in §02 that **cuts a ratified sentence in half on mobile** — and it is
invisible to the exact gate assertion the blueprint prescribed, which is why it survived to QA.

---

## Findings

### 🔴 BLOCKING-1 — §02 `maker-line` text is clipped at 375 and 768; the FlowForgeSoft attribution is unreadable on mobile

**What.** The paragraph `.z-maker-line__text` has an intrinsic width of **762px and never shrinks or wraps**.
Its flex row is `flex-wrap: nowrap` and the paragraph is `flex: 0 0 auto`, so at narrow viewports it simply
runs off the side of the page and is clipped by `overflow-x: hidden` on `<html>`/`<body>`. The text is not
scrollable-to; it is **destroyed**.

**Where.** `#maker-line > .z-maker-line__row > p.z-maker-line__text` — at 375px and 768px. Reproduces
**identically in Chromium, Firefox and WebKit**.

**Evidence** (`qa/harness/mk2.mjs`, character-level `Range.getBoundingClientRect()` measurement of what
actually falls inside the viewport):

```
== 320px ==  text right=806  OVERSHOOT past viewport = 486px   clipped characters: 50
  FULL    : "Zerado is the second product from FlowForgeSoft, the company behind FlowForge."
  VISIBLE : "Zerado is the second product"
== 375px ==  text right=806  OVERSHOOT past viewport = 431px   clipped characters: 45
  VISIBLE : "Zerado is the second product from"
== 768px ==  text right=814  OVERSHOOT past viewport = 46px    clipped characters: 5
  VISIBLE : "Zerado is the second product from FlowForgeSoft, the company behind FlowF"
== 1280px == OVERSHOOT = -210px                                clipped characters: 0
```

Cross-engine (`qa/harness/xb-maker.mjs`):

```
chromium  @375: OVERSHOOT=431px   @768: OVERSHOOT=46px
firefox   @375: OVERSHOOT=431px   @768: OVERSHOOT=46px
webkit    @375: OVERSHOOT=431px   @768: OVERSHOOT=46px
```

Screenshots: `screenshots/DEFECT-maker-line-375.png` (reads *"Zerado is the second product from F"*),
`DEFECT-maker-line-768.png` (reads *"…the company behind FlowF"*), plus `-320`, `-1280`, and
`-firefox-375/768`, `-webkit-375/768`.

**Why it matters.**
1. **It deletes ratified content.** `decisions.md` records permission to name FlowForgeSoft as an
   answered-and-closed item, and the blueprint gives §02 an entire section to carry it. At 375 the visitor
   never sees the words "FlowForgeSoft" or "FlowForge" — the section renders a truncated fragment that says
   nothing. The one job of §02 fails on the most common viewport class.
2. **WCAG 2.1 SC 1.4.10 Reflow (AA) failure.** At 320px CSS width, content is lost with no horizontal
   scroll available to recover it. Clipping is content loss, not reflow.
3. It is not a cosmetic near-miss — 45 of 78 characters (58%) are gone at 375px.

**Root cause** is one declaration: the paragraph is `flex: 0 0 auto` inside a `nowrap` row. Changing it to
permit shrinking (`flex: 0 1 auto` / `min-width: 0`) or allowing the row to wrap fixes both this and MAJOR-2.
*(Reported, not fixed — I do not modify `site/`.)*

---

### 🟠 MAJOR-1 — The blueprint's prescribed overflow assertion cannot detect BLOCKING-1

**What.** Blueprint §5.5 specifies the gate assertion:
`document.documentElement.scrollWidth === document.documentElement.clientWidth`.
That assertion **passes at all four viewports** — while content is being clipped.

**Evidence** (`qa/qa-raw-results.json`):

```
== 375 == docEl 375/375  pass=True | body 806/375 | elements past viewport: 3
== 768 == docEl 768/768  pass=True | body 814/768 | elements past viewport: 3
== 1280 == docEl 1280/1280 pass=True | body 1280/1280 | elements past viewport: 1
== 1920 == docEl 1920/1920 pass=True | body 1920/1920 | elements past viewport: 0
```

`html`/`body` carry `overflow-x: hidden` (measured, all three engines). That pins `documentElement.scrollWidth`
to `clientWidth` **by construction**, so the assertion is a tautology and can never fail.

**Why it matters.** §5.5 ranks its mechanisms — reflow first, reformat second, contained scroll as a
"backstop only" with `tabindex="0"` so keyboard users can reach it. A blanket `overflow-x: hidden` on the
root is **none of the three**: it is a fourth mechanism that silently discards content and simultaneously
disables the only automated check the design specified. Any future width regression is equally undetectable.
The honest assertion is per-element (`element.right <= clientWidth` for non-`aria-hidden` content), which is
what `qa/harness/run-core.mjs` runs and which caught this.

The only other element extending past the viewport is `.z-frame__watermark` — verified `aria-hidden="true"`
decorative repeating "PHASE 4" text inside a deliberate `.z-frame__watermark-clip`. **Not a defect.**

---

### 🟠 MAJOR-2 — §02's flanking rules render at 0px width at *every* viewport

**What.** Blueprint §1.1 assigns §02 the layout family *"Single centred rule-flanked line"*, and §6.3's
MakerLine component consumes `--z-border` for those rules. Both `span.z-maker-line__rule` elements
(`flex: 1 1 0%`) are starved to zero width because the paragraph consumes the entire row.

**Evidence** (`qa/harness/xb-maker.mjs`, all engines): `rule widths=[0,0]` at 375 and 768. At 1280/1920
(`qa/harness/mk.mjs`) both rules measure `left=284 right=284 width=0` and `left=1094 right=1094 width=0`.

**Why it matters.** The ratified layout family is never realised anywhere on the page — §02 renders as a bare
centred line, not the specified rule-flanked one. Same root cause as BLOCKING-1; one fix resolves both. Worth
its own line because fixing only the clipping (e.g. by adding `white-space: normal` alone) would leave the
rules still absent.

---

### 🟡 MINOR (7)

| # | Finding | Evidence |
|---|---|---|
| MINOR-1 | **FAQ expand/collapse motion not implemented.** Blueprint §7.6 motions #7/#8 specify 220ms expand / 140ms collapse. Measured `transition-duration: 0s` on both `.z-faq-item__summary` and `.z-faq-item__answer`; `animation-name: none`. The accordion snaps open. Defensible (native `<details>` cannot animate without `::details-content` + `interpolate-size`, which measured `true` in Chromium but **`false` in Firefox and WebKit**) — but §7.6 claims the nine-entry table "is the complete list", and two of the nine do not exist. | `qa/harness/motion2.mjs`, `run-xb.mjs` |
| MINOR-2 | **Parked scanner track count is 5, blueprint says 4.** §7.6 motion #2 enumerates "TF-1/2/3 chrome bars + §15 top edge (4 instances)". Measured 6 tracks total: 1 live (`#hero`) + parked in `#hero` figure, `#moods`, `#terminal-first`, `#community`, `#closing-cta`. The extra is the §11 community panel's chrome bar. Static in both motion modes, so zero motion harm — an inventory-count drift. | `qa/harness/motion2.mjs` |
| MINOR-3 | **Footer contact link is 19px tall.** `a[href="mailto:alex@flowforgesoft.com"]` measures **180×19px** at both 375 and 1280 (`display: inline`, no line-box padding). Blueprint §5.4 commits to "Every interactive element is ≥ 48 × 48px at every breakpoint" and "Footer links: 44px line-box". **This is NOT a WCAG failure** — SC 2.5.8 (AA) is satisfied via the spacing exception (nearest other target centre is 65px away, ≥24px). It is a deviation from the blueprint's own stated floor, and the only target on the page below it; every other target measured ≥44px. | `qa/harness/targets.mjs`, `nav375.mjs` |
| MINOR-4 | **CoverTile art variety is 4, blueprint says 12.** §6.3 specifies "twelve fixed variants (§7.5)". Measured 12 tiles (`.z-tile`) drawing on **4** art variants (`z-tile--A/B/C/D`), repeated. | `dist/index.html` grep: `z-tile count: 12`, `variants: ['A','B','C','D']` |
| MINOR-5 | **`StateLegend` resolves to no DOM node.** Listed in blueprint §6.3's component table, but §05's own layout spec never places it (it specifies "StateChip × 4 in one row", which the build renders correctly with glyph + LABEL + line). Contract self-inconsistency rather than a build gap — flagging so it is closed deliberately. | `qa/harness/components.mjs` |
| MINOR-6 | **FlowForge footer mark renders at 18px, blueprint §7.8 says 28px.** `img[src="/flowforge-logo.svg"]` displayed `18x18` (natural 150×150, loads fine — `loading="lazy"`, HTTP 200). | `qa/harness/links.mjs` |
| MINOR-7 | **Link inventory exceeds blueprint §2.1** — which states "exactly three destinations and nine link instances … fifteen focusable elements". Build has **10 links / 4 destinations / 16 focusables** at ≥768. The extra is the `https://www.flowforgesoft.com` "Powered by FlowForge" link, which §6.3 *does* require and which my brief requires. **The build is correct; §2.1's table is stale.** Separately, that link uses `target="_blank"` with no external-link cue, where §2.1 reasons that GitHub should open in the same tab. | `qa/harness/run-content.mjs` |

---

## Per-checklist results

### 1 · Screenshots at 375 / 768 / 1280 / 1920 — ✅ done
Full-page PNGs at `qa/screenshots/{375,768,1280,1920}-fullpage.png` (chromium), plus `{firefox,webkit}-*`
for all four widths — 12 full-page captures. Document heights: 375→14,926px · 768→11,700px ·
1280/1920→11,442px. Also captured: `skiplink-focused-{375,1280}.png`, `reduced-motion-{reduce,no-preference}-1280.png`,
`scanner-{reduce,no-preference}.png`, and 8 `DEFECT-maker-line-*` crops.

### 2 · axe-core WCAG 2.1/2.2 AA — ✅ **0 violations at all four viewports**
Tags: `wcag2a, wcag2aa, wcag21a, wcag21aa, wcag22aa`.

```
[375]  violations=0  passes=28  incomplete=1
[768]  violations=0  passes=28  incomplete=1
[1280] violations=0  passes=28  incomplete=1
[1920] violations=0  passes=28  incomplete=1
```

The one `incomplete` is `color-contrast` (42 nodes at 1280), meaning **axe could not compute it** — not that
it passed. Breakdown of axe's own reasons:

```
x31  "Element content contains only non-text characters"   (the ○ ◐ ◉ ⊘ ▸ ▣ glyphs, all aria-hidden="true")
x10  "background color could not be determined because it is overlapped by another element"
x1   "background color could not be determined because it partially overlaps other elements"  -> .z-maker-line__text
```

That last one is axe independently noticing BLOCKING-1.

**Because axe punted, I computed contrast myself** (`qa/harness/contrast.mjs` — WCAG relative-luminance
formula, effective background resolved by walking ancestors and compositing alpha, `aria-hidden` subtrees
excluded, large-text threshold applied per element):

```
@1280: 206 text elements — FAILS: 0 — over a gradient: 0 — lowest 4.89:1 (needs 4.5)
@375 : 188 text elements — FAILS: 0 — over a gradient: 0 — lowest 5.53:1 (needs 4.5)
```

Lowest measured pairs: `#8492A8` on `#1D2532` = 4.89:1 (`.z-trow__platform`, 15px/400) and `#8492A8` on
`#141A24` = 5.53:1. All clear AA. Brand manual §4's co-render rule is satisfied — every state carries
**colour + glyph + text label**, with the glyph `aria-hidden` so AT reads the label once:

```
○ NOT STARTED  rgb(165,162,155) = #A5A29B    ◐ IN PROGRESS rgb(255,176,0)  = #FFB000
◉ ZERADO       rgb(25,224,255)  = #19E0FF    ⊘ ABANDONED   rgb(199,125,255)= #C77DFF
```

Ratified order Not Started · In Progress · Zerado · Abandoned is intact, and chip 3 is composed differently
exactly as §05 requires: width 336px vs 249px, glyph 30px vs 24px, and the page's **only** cyan glow.

### 3 · Horizontal overflow — ⚠️ assertion passes, real overflow present
See MAJOR-1 and BLOCKING-1 above. Per-element check found 3 offenders at 375, 3 at 768, 1 at 1280, 0 at 1920;
after excluding the `aria-hidden` watermark, all are the §02 paragraph and its `<strong>`.

**Intentionally-scrollable containers:** none are keyboard-trapped. `TerminalFrame` never needed its
`overflow-x: auto` backstop — at 375 the terminal frames reflow to two-line rows as §5.5 mechanism 2 intends,
and no `.z-frame__screen` reported `scrollWidth > clientWidth`. The only `overflow: hidden` containers are
`.z-frame__watermark-clip` (decorative, `aria-hidden="true"`) and the root. So there is no scrollable region
requiring `tabindex="0"`, and none is missing.

### 4 · Full keyboard traversal — ✅ pass
Chromium, 1280 and 375.

- **Skip link is first in DOM and first on Tab**, `href="#main"`, target `#main` exists.
- **It becomes visible on focus.** ⚠️ *Correction to my own first measurement:* I initially recorded the
  focused skip link at `top: -900px` (off-screen) and nearly filed it as blocking. That was a measurement
  artifact — I sampled computed style in the same tick as the 140ms `top` transition. Re-measured over time:

  ```
  t=0ms    top=-892.406px  onScreen=false   <- what a naive probe sees
  t=50ms   top=-355.625px  onScreen=false
  t=150ms  top=16px        onScreen=true
  t=400ms  top=16px        onScreen=true
  t=1000ms top=16px        onScreen=true
  @375px   top=16px  rect x=16 y=16 w=176   onScreen=true
  ```
  Visually confirmed in `screenshots/skiplink-focused-375.png` — cyan pill, top-left, focus ring drawn.
- **It works.** Enter → `location.hash=#main`, scrollY 0→97; the next Tab lands on the hero CTA *inside*
  `#main`, correctly skipping the nav's links. (`#main` has no `tabindex="-1"`, so `document.activeElement`
  stays `body` and the browser's sequential-focus starting point does the work — functional in all three
  engines tested, though an explicit `tabindex="-1"` would be more robust.)
- **Focus order matches DOM order exactly** — 16 tab stops at 1280, 15 at 375, compared element-by-element
  against `querySelectorAll` DOM order. No divergence.
- **Every stop has a visible indicator:** all 16 report `:focus-visible` true with
  `outline: 2px solid rgb(25,224,255)`, `outline-offset: 2px` — constant geometry, as §6.2 specifies. Ring is
  present under `prefers-reduced-motion: reduce` too (only its transition collapses).
- **No focus trap** — Tab exits the document after the last stop.
- **No unreachable interactive element, and no zero-size focusable.** At 375 the nav GitHub link measures
  0×0 — verified `display: none`, i.e. correctly removed from the tab order (15 stops, `zero-size focusable: 0`).
  This is the ratified §5.2 no-hamburger behaviour ("the GitHub link drops to the footer"), **not a defect**.
- **`<details>` open and close by keyboard: 6/6 pass.**
  ```
  details[0..5] closed->false  Enter->true  Enter->false   OK
  ```

### 5 · `prefers-reduced-motion: reduce` — ✅ pass, including the pip park
The critical assertion holds. Measured `.z-scanner-track--live .z-scanner-track__pip`:

```
no-preference : anim=z-scanner-sweep dur=2.4s iter=infinite  offset-from-track-centre=72px  opacity=1
reduce        : anim=none            dur=0s                  offset-from-track-centre=0px   opacity=1
                visibility=visible  display=block  pip=18% of track  bg=rgb(255,46,46)  glow intact
```

**The pip parks dead-centre at full opacity and is never hidden** — blueprint §7.6's explicit requirement.
Animated-element count across the page drops from 1 → 0; all transitions collapse to `0.001s`.

| # | Motion | Reduced-motion state observed | Verdict |
|---|---|---|---|
| 1 | ScannerRule sweep (1 live, hero) | `animation: none`, parked at centre offset **0px**, opacity **1**, visible | ✅ per spec |
| 2 | Parked scanner tracks | `animation: none` in **both** modes — static by design | ✅ |
| 3 | CTAButton primary hover (3 found) | `.z-cta` 0.14s → **0.001s** | ✅ |
| 4 | CTAButton secondary hover (1 found) | same `.z-cta` rule, 0.14s → **0.001s** | ✅ |
| 5 | Focus ring appearance | ring **still drawn** (2px solid `rgb(25,224,255)`, offset 2px) in reduce; only transition collapses | ✅ per spec |
| 6 | NavLink / FooterLink underline (4 links) | `.z-nav-link`, `.z-footer-link` ×2, `.site-footer__powered` 0.14s → **0.001s** | ✅ |
| 7 | FAQItem expand (220ms spec) | measured **0s in both modes** — never implemented | ⚠️ MINOR-1 |
| 8 | FAQItem collapse (140ms spec) | measured **0s in both modes** — never implemented | ⚠️ MINOR-1 |
| 9 | GridHorizon | static in both modes | ✅ |

`scroll-behavior: auto` in both modes. No scroll-triggered animation, parallax, marquee, or autoplay found.

### 6 · Cross-browser — ✅ pass, no meaningful divergence
Chromium 151.0.7922.34 · Firefox 153.0 · WebKit 26.5, each at all four viewports (`qa/qa-crossbrowser.json`).

- **Zero page errors, zero 4xx/5xx responses, in every engine at every viewport.**
- Overflow assertion result identical in all three (and BLOCKING-1 reproduces identically in all three).
- Document heights agree to within 0.4%: 1280 → chromium 11,442 / firefox 11,443 / webkit 11,442.
- **`<details>` styling** — consistent: `summary` height 56px in all three; `list-style: none` +
  `::-webkit-details-marker { display: none }` suppresses the native marker everywhere; the custom `+`/`−`
  marker renders in all three.
- **`clamp()`** — supported everywhere; H1 resolves to 30.0007px (Chromium) / 30px (FF) / 30.000675px (WebKit)
  at 375, and 52px in all three at 1280. H1 box identical at 335×97 (375) and 650×168 (1280).
- **`aspect-ratio`** — `CSS.supports` true in all three; cover tiles lay out identically (12 tiles, uniform grid).
- **Gradients** — 12 gradient-bearing elements in all three, no fallback divergence. **No text sits on a
  gradient** (independently verified: `over a gradient: 0`), so no contrast risk from them.
- **`backdrop-filter`** — supported in all three but **used by zero elements**; not a risk on this page.
- **Variable fonts** — all three self-hosted `woff2` files load in every engine (Orbitron 400–900,
  Space Grotesk 300–700, JetBrains Mono 100–800, `status: loaded`), served same-origin, HTTP 200. Fallback
  stacks (`Michroma`, `Arial Black`, `system-ui`) present.

Only divergence found, and it is benign — one extra wrapped text line in Firefox/WebKit at narrow widths:

```
@375  moods           chromium=1655  firefox=1683  webkit=1655   Δ=28px
@375  one-collection  chromium=1732  firefox=1760  webkit=1760   Δ=28px
@768  states          chromium=710   firefox=735   webkit=710    Δ=25px
@1280 and @1920: no section differs by more than 4px.
```

### 7 · Content honesty sweep — ✅ **PASS, no violations**

- **Forbidden tokens: zero occurrences.** Scanned both visible `innerText` and full serialised HTML for
  the Q4-banned community-source name, `lorem ipsum`, `TODO`, `[TBD]`, `{{`, `${`, `R$ 00`, plus `FIXME`, `XXX`, `placeholder`, `Lorem`,
  `undefined`, `NaN`. Result: `forbidden-token hits: NONE`. Lone-`X` placeholder regex over visible text:
  `NONE`. The Q4 constraint (the banned community-source name appears nowhere — copy, icon, screenshot or FAQ) **holds**.
- **All four roadmap phases read "Planned".** Rendered: `PHASE 1 CLI/TUI MVP … ○ PLANNED`, `PHASE 2
  Enrichment … ○ PLANNED`, `PHASE 3 Recommendations & Budget … ○ PLANNED`, `PHASE 4 Social & Mo… ○ PLANNED`.
  4 `.z-status-marker` nodes, each `○` + `Planned`, colour `#A5A29B` — none marked done.
- **§11 community has zero interactive elements.** Queried `a[href], button, input, select, textarea,
  summary, [tabindex]:not([tabindex="-1"]), [role="button"], [onclick]` inside `#community` → **0 results**.
- **No donate/sponsor anywhere.** Searched text and HTML for `donate, sponsor, patreon, ko-fi, kofi,
  buymeacoffee, github sponsors, paypal, pix` → `money words found: []`, `donate/sponsor LINKS: []`.
  The two ratified money *statements* are both present and are disclosure only: the §11 line *"servers cost
  money — so this part will need a premium account or a donation once it exists. Nothing is decided about
  price, tiers, or a date, and nothing on this page is asking…"* and the footer's affiliate-commission
  disclosure. Exactly the Q3/Q6 resolution in `decisions.md`.
- **All 6 mockup figures carry a visible `<figcaption>`**, each `role="img"` + `aria-labelledby` → its own
  caption id, each verbatim from `content/landing-copy.md`:

  | § | caption visible | opening words |
  |---|---|---|
  | hero | ✅ | "An illustration of Zerado's terminal library view. This is a mockup of the planned interface, not a screenshot…" |
  | moods | ✅ | "Mockup of Zerado's mood picker in the terminal…" |
  | one-collection | ✅ | "Mockup of a cover-art grid… **Cover tiles are illustrative artwork, not real game covers.**" |
  | terminal-first | ✅ | "Mockup of a game's detail page in the terminal…" |
  | and-on-your-phone | ✅ | "An illustration of two phone frames… Not built yet — Phase 4." |
  | community | ✅ | "Example mockup of the community layer… watermarked 'Phase 4, not available yet.'" |

- **The cover-grid caption renders**, appended to the §06 figcaption, verbatim:
  *"Cover tiles are illustrative artwork, not real game covers."*
- **"Powered by FlowForge" is present and correct**: `href="https://www.flowforgesoft.com"`,
  `rel="noopener noreferrer"`, wrapping `img src="/flowforge-logo.svg" alt="FlowForge"` which **resolves
  200** (`image/svg+xml`, natural 150×150). Not a 404. (See MINOR-6 re: 18px vs 28px render size.)
- Bonus checks against `decisions.md` Q5: *zerado* (25 mentions) and *sinopse* (2) are kept, each glossed
  in-line ("Portuguese for…"), and both first appearances are `<em>`-italicised per the copy casing rule.

### 8 · Link integrity — ✅ **PASS, zero 4xx/5xx**

10 `<a href>` on the page. Full inventory:

```
[0] #main                                                    "Skip to content"
[1] https://github.com/JustCode-CruzAlex/Zerado    rel=noopener noreferrer   nav
[2] mailto:alex@flowforgesoft.com?subject=Zerado-WaitList                    nav
[3] mailto:alex@flowforgesoft.com?subject=Zerado-WaitList                    §hero
[4] https://github.com/JustCode-CruzAlex/Zerado    rel=noopener noreferrer   §hero
[5] mailto:alex@flowforgesoft.com?subject=Zerado-WaitList                    §closing-cta
[6] mailto:alex@flowforgesoft.com?subject=Zerado-WaitList                    §footer
[7] https://github.com/JustCode-CruzAlex/Zerado    rel=noopener noreferrer   §footer
[8] mailto:alex@flowforgesoft.com                                            §footer contact
[9] https://www.flowforgesoft.com                  rel=noopener noreferrer target=_blank  §footer
```

- **In-page anchors:** 1 (`#main`) → target id exists. No `href="#"`, no dangling fragment.
- **Outbound HTTP** (curl, no redirect / following redirects):
  ```
  https://github.com/JustCode-CruzAlex/Zerado  -> no-redirect:200  followed:200
  https://www.flowforgesoft.com                -> no-redirect:200  followed:200
  ```
  **No CTA is dead.** Outbound destination set is exactly the expected `{github, flowforgesoft.com, mailto}`.
- **`mailto:` subject verified:** all 4 waitlist CTAs carry `?subject=Zerado-WaitList` exactly; the footer
  contact link is deliberately bare per copy §16.
- **Same-origin assets, all 200:** `/`, `/fonts.css`, `/logo.svg`, `/flowforge-logo.svg`, `/favicon.svg`,
  `/_astro/index.af0tTXyW.css`, and 3 `woff2` files. **`requestfailed` count: 0.** Zero external requests —
  the §1.2 zero-external-request commitment holds.

### 9 · Semantic structure — ✅ pass

```
lang = "en"
title = "Zerado — terminal game library tracker, sorted by mood"
meta description present
landmarks: banner/header=1  main=1  contentinfo/footer=1  nav=1
headings: 13   h1 count: 1   skipped levels: NONE
```

Heading ladder is h1 → h2 ×12, one per section that the blueprint gives a heading; §02 `maker-line`,
§15 `closing-cta` and §16 `footer` correctly carry none, matching §1.1's *"(no heading)"* entries.

**Images: 2 total** (the page ships zero raster images as §1.2 commits — every "image" is composed HTML/CSS
or SVG). Both have `alt`: `/logo.svg` → `alt="Zerado"`, `/flowforge-logo.svg` → `alt="FlowForge"`. The six
mockup figures use `role="img"` + `aria-labelledby` → visible `<figcaption>`, so their alt text is the
on-screen caption and matches `content/landing-copy.md` verbatim (table in §7 above) — announced once, not twice.

**JavaScript:** one `<script>` tag, `type="application/ld+json"` (schema.org `SoftwareApplication`) — structured
data, not executable script. The §1.2 zero-JavaScript commitment holds. Zero page errors in all three engines.

### 10 · Blueprint conformance — ✅ 16/16 sections, 38/40 components

**All 16 sections present, in the ratified DOM order, with the exact ratified ids:**

```
hero · maker-line · the-problem · moods · states · one-collection · price-intelligence ·
yours-and-offline · terminal-first · and-on-your-phone · community · roadmap ·
after-phase-4 · faq · closing-cta · footer
```

0 added · 0 dropped · 0 reordered (15 `<section>` in `<main>` + `<footer id="footer">`).

Component inventory (`qa/harness/components.mjs`) — 38 of 40 resolve to DOM nodes. Counts match spec where
spec gives one: CTAButton primary ×3, secondary ×1, FAQItem ×6, TerminalRow ×6, MoodCard ×4, StateChip ×4,
StoreRow ×4, PhaseCard ×4, StatusMarker ×4, TrustPoint ×3, PhoneFrame ×2, CoverGrid ×1, PriceCard ×1,
Sparkline ×1, CommunityPanel ×1, GridHorizon ×1, `.z-tile` ×12.
Two did not resolve: **`StateLegend`** (MINOR-5 — absent, and unplaced by the blueprint's own §05 spec) and
**`CoverTile`** under my first guessed selector — it resolves as `.z-tile` (12 nodes), so it is present;
only its variant count deviates (MINOR-4).

---

## What I could NOT test, and why

1. **Real assistive technology.** No VoiceOver / NVDA / JAWS run. Everything in §9 is programmatic (roles,
   names, landmark and heading structure). Announcement *quality* — how the six `role="img"` figures and the
   glyph/label state pairs actually sound — is unverified. This is the largest residual gap.
2. **Real devices and real touch.** All four viewports were emulated in desktop engines at
   `deviceScaleFactor: 1`. No physical iOS/Android device, no real Safari-on-iOS (Playwright WebKit ≠ iOS
   Safari), no actual finger-target testing, no on-device font rendering.
3. **The `mailto:` hand-off itself.** I verified the four CTA `href`s carry `?subject=Zerado-WaitList` and
   that the address is well-formed, but I cannot launch a mail client to confirm a composer opens with the
   subject pre-filled. The `decisions.md` Q2 trade-off (mail-client hand-off converts worse; the address is
   scrapable) is a product decision, already recorded, not re-litigated here.
4. **Core Web Vitals, Lighthouse, bundle size.** Deliberately out of scope — the `performance` lane owns
   these and I was told not to duplicate. I report no LCP/CLS/INP number and none should be inferred from
   this document.
5. **Variable-font rasterisation quality.** I verified all three fonts load with the correct weight ranges in
   all three engines; I did not assess hinting, subpixel rendering, or glyph-level quality, which needs a
   human eye on real hardware.
6. **`forced-colors` / Windows High Contrast Mode, and `prefers-contrast`.** Not emulated. A dark-native page
   with no light mode is exactly the kind of surface where forced-colors can go wrong; worth a follow-up.
7. **400% zoom (WCAG 1.4.10's other half).** I tested the 320px-width reflow floor — where BLOCKING-1 is at
   its worst (50 characters lost) — but did not separately exercise browser zoom at 400%, which can expose
   different failures than a narrow viewport.
8. **axe's own contrast verdict.** axe returned `incomplete`, not `pass`, for 42 nodes. The clean contrast
   result in §2 is **my** implementation of the WCAG relative-luminance formula
   (`qa/harness/contrast.mjs`), not axe's. It is auditable and I have shown its numbers, but it is not an
   independent tool corroborating axe — the 31 glyph-only nodes remain formally unverified by axe, though
   they are all `aria-hidden="true"` decoration whose meaning is carried by an adjacent text label.
9. **One correction worth recording.** My first skip-link probe read `top: -900px` on focus and would have
   been filed as a second BLOCKING finding. It was wrong — the sample landed inside a 140ms transition. The
   re-measured time series is in §4. No claim in this report rests on a single-tick computed-style read.

---

## Bottom line

Ship-blocking work is **one CSS declaration in one component**. `.z-maker-line__text` must be allowed to
shrink and wrap; that single fix clears BLOCKING-1 and MAJOR-2 together. Separately, and independent of the
fix, the `overflow-x: hidden` on the document root should be reconsidered (MAJOR-1) — it is what let a
58%-content-loss defect pass the design's own prescribed gate assertion at every viewport.

Everything else on this page is in good shape: zero axe violations, zero contrast failures, zero dead links,
zero placeholder or dishonest content, a correct and visible skip link, a fully keyboard-operable page with a
consistent focus indicator on all 16 stops, an exactly-conforming reduced-motion implementation including the
scanner pip park, and three browser engines that agree to within 28px on an 11,442px page.

*Re-run everything with: `cd qa/harness && node run-core.mjs && node run-kbd.mjs && node run-rm.mjs &&
node run-content.mjs && node run-xb.mjs && node contrast.mjs`*

---

# ADDENDUM — fix round, verified by the orchestrator

*Appended after the fix round. The verdict above (FAIL, 1 blocking) was correct when written; this
records what changed and how it was proven. The original findings are left intact — a report that
edits away its own findings can't be audited.*

## BLOCKING-1 — §02 `maker-line` text clipping · **RESOLVED**

**Root cause, confirmed:** `.z-maker-line__text` carried `flex-shrink: 0`, pinning it at its 762px
intrinsic width forever. The flanking `.z-maker-line__rule` spans had `flex-basis: 0%`, so they could
never receive any share of flex negative-space redistribution — which is why they rendered at 0px at
*every* viewport. Two mechanisms, one trigger.

**Fix:** text → `flex: 0 1 auto; min-width: 0`; rules → `flex: 1 1 0%` with `min-width: 24px`, forcing
the flex resolver to take the remainder from the text.

**Independently re-measured by the orchestrator** (Playwright, live preview, `#maker-line .z-maker-line__text`):

| Width | Chars rendered | Right edge / viewport | Flanking rules |
|---|---|---|---|
| 320 | **78 / 78** | 252 / 320 ✅ | 24px · 24px |
| 375 | **78 / 78** | 307 / 375 ✅ | 24px · 24px |
| 768 | **78 / 78** | 692 / 768 ✅ | 24px · 24px |
| 1280 | **78 / 78** | 948 / 1280 ✅ | 24px · 24px |
| 1920 | **78 / 78** | 1268 / 1920 ✅ | 24px · 24px |

Full string at 375: *"Zerado is the second product from FlowForgeSoft, the company behind FlowForge."*
Character count matches `content/landing-copy.md` exactly at all five widths. The ratified
FlowForgeSoft attribution now survives on a phone, and the rule-flanked layout renders as designed.

## The finding behind the finding — **the tautological gate is CLOSED**

This report's most valuable observation was that its own prescribed gate could never fail:
`documentElement.scrollWidth === clientWidth` is pinned to `true` by `overflow-x: hidden` on the root,
so it reported PASS at all four viewports while 58% of a section's text was being destroyed.

`qa/harness/run-core.mjs` has been changed by the orchestrator (backup at `run-core.mjs.bak`):

- `pass` is now `offenders.length === 0` — the per-element check the harness already computed every
  run but never gated on. The old document-level comparison is retained as
  `documentScrollEqualsClient` for reference, demoted from verdict to datum.
- `aria-hidden` subtrees are excluded from `offenders`, so the by-design clipped `PHASE 4` watermark
  no longer masquerades as content loss.

**Mutation-tested, because a gate nobody has watched fail is just a different tautology:**

```
375px  HEAD (fixed):                 offenders=0   PASS
375px  MUTANT (bug re-injected):     offenders=2   KILLED -> p.z-maker-line__text, strong
768px  HEAD (fixed):                 offenders=0   PASS
768px  MUTANT (bug re-injected):     offenders=2   KILLED -> p.z-maker-line__text, strong
```

The pair is what matters: **green at HEAD and red at the mutant.** The gate now detects the exact
defect it previously slept through.

## Re-run against the fixed build — `run-core.mjs`, live

```
[375]  overflow pass=true | axe violations=0  passes=28  incomplete=1
[768]  overflow pass=true | axe violations=0  passes=28  incomplete=1
[1280] overflow pass=true | axe violations=0  passes=28  incomplete=1
[1920] overflow pass=true | axe violations=0  passes=28  incomplete=1
```

The single reported scroller at every width is `.z-frame__watermark-clip` — `overflow-x: hidden`, not
`auto`, inside an `aria-hidden` subtree. It is a clip, not a scrollable region, so it correctly needs
no `tabindex`. Not a defect.

## No-regression sweep (orchestrator, against the rebuilt `dist/`)

| Check | Result |
|---|---|
| Forbidden tokens (the Q4-banned source name, `lorem`, `TODO`, `placeholder`, `{{`) | **0** |
| Roadmap phases reading "Planned" | **4** |
| `Powered by FlowForge` footer block | **1**, correct href |
| Dead anchors (`href="#"`) | **0** |
| `mailto:` carrying `subject=Zerado-WaitList` | **4** |
| Ratified sections present, in order | **16 / 16** |
| `<h1>` count | **1** |
| JS files in `dist/` | **0** |
| External hosts in built HTML | github.com · flowforgesoft.com · zerado.app · schema.org — no CDN, no analytics |

## Also fixed this round (from the performance lane, and one gap the orchestrator found)

- **Conditional mono preload** — `JetBrainsMono-latin.woff2` (42.7% of page weight) no longer races the
  LCP-critical fonts on mobile, where it has zero above-fold consumers. Verified by CDP capture:
  `isLinkPreload=false` at 375, `true` at 1280.
- **`fonts.css` inlined** — one render-blocking round trip removed; requests 8 → 7; desktop
  render-blocking savings now below Lighthouse's reporting threshold.
- **`public/_headers` added** — `immutable` only for content-hashed `/_astro/*`, revalidating policy for
  the *unhashed* fonts (avoiding a stale-font trap), `no-cache` for `index.html`.
- **OG image, found missing by the orchestrator** — the page declared `twitter:card=summary_large_image`
  with **no image**, so every social share would have rendered blank. A 1200×630 branded card is now
  built, wired via absolute `https://zerado.app/og-card.png`, with honest alt text: the card carries a
  visible *"MOCKUP OF THE PLANNED INTERFACE"* label and the alt text does not imply a real screenshot.

## Verdict after fix round: **PASS**

Zero blocking findings. The largest remaining gap is unchanged and worth restating plainly: **every
accessibility check here is programmatic.** Real screen-reader announcement quality has not been
verified by a human with a screen reader, and axe returning 0 violations is not the same claim.
