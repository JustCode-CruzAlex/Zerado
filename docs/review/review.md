# VERDICT: GOLDEN — 0 blocking findings

*Zerado landing page · review phase · 2026-08-24*
*Reviewed against the ratified outline (SHA-verified), the six ratified decisions, the design blueprint, the copy contract and the brand manual. Every gate below was re-run first-hand against the live preview at `http://localhost:4321/`; the prior QA and performance reports were treated as claims to be tested, not as evidence.*

**Findings: 0 BLOCKING · 0 MAJOR · 7 MINOR.**
None of the seven blocks the page. Five are drift in the *blueprint manifest* (the document describing the page) rather than defects in the page itself; two are visual polish. Ship it.

---

## 1. The chain of contracts

| # | Contract | Result | Evidence I gathered myself |
|---|---|---|---|
| 1 | `ratification/mock.outline.md` — 16 sections, ratified order | **PASS** | `shasum -a 256` = `a2bc87e4…48ecb`, exact match to the declared SHA. DOM: 15 `<section>` in `<main>` in exact ratified id-order (`hero, maker-line, the-problem, moods, states, one-collection, price-intelligence, yours-and-offline, terminal-first, and-on-your-phone, community, roadmap, after-phase-4, faq, closing-cta`) + `<footer>` = §16. Nothing added, dropped or reordered. |
| 2 | `ratification/decisions.md` — six ratified answers | **PASS** | Detailed below. |
| 3 | `design/blueprint.md` + `.tokens.json` | **PASS** (5 manifest-drift MINORs) | Section order, component→DOM resolution, computed tokens, breakpoint behaviour all verified live at 375/768/1280/1920. |
| 4 | `content/landing-copy.md` — verbatim | **PASS** | 23 spot-checked strings across all 16 sections, all verbatim. |
| 5 | `brand/brand-manual.md` + `tokens.css` | **PASS** | Co-render rule, Orbitron floor, contrast all measured live. |

### 1.2 — The six ratified decisions, each verified in the built page

| Decision | Verified | How |
|---|---|---|
| Name **Zerado**, domain `zerado.app` | ✅ | Wordmark, `<title>`, `og:url`, JSON-LD all `Zerado` / `https://zerado.app/`. `dig zerado.app` resolves (`13.248.243.5`, `76.223.105.230`). |
| CTA = waitlist primary, GitHub secondary | ✅ | 5 × `mailto:alex@flowforgesoft.com?subject=Zerado-WaitList` (subject exact), 3 × `https://github.com/JustCode-CruzAlex/Zerado`. Waitlist is the filled cyan button, GitHub the outlined one, at every occurrence. GitHub repo returns **HTTP 200** anonymously — re-checked, not assumed. |
| Community = full section, Phase-4 stamped, states premium/donation, **no donate button** | ✅ | §11 exists as a full section. **Three independent Phase-4 stamps** (badge `⊘ PHASE 4 · NOT AVAILABLE YET`, amber heading clause, diagonal watermark). States verbatim: *"this part will need a premium account or a donation once it exists."* **Interactive elements in §11: 0** — queried `a, button, summary, input, [tabindex]` inside `#community`, empty set. Zero `donate`/`sponsor`/`Patreon`/`Ko-fi` anywhere in text or HTML. |
| The **Q4-banned community-source name** appears nowhere | ✅ | Case-sensitive *and* case-insensitive scan of `body.innerText` **and** full page HTML: **0 hits**. |
| English page, `zerado`/`sinopse` explained once in-line | ✅ | `<html lang="en">`. Exactly two `lang="pt-BR"` elements: `<em>Zerado</em>` (§05) and `<span>sinopse</span>` (§06) — each glossed in-line once, never repeated. |
| Affiliate disclosed in footer, **no sponsor ask anywhere** | ✅ | Footer carries the disclosure verbatim. No funding CTA on the page; the only two money statements are the footer disclosure and the §11 cost disclosure, exactly as ratified. |

---

## 2. The hard gates

| Gate | Result | Evidence |
|---|---|---|
| No dead CTA / broken link | **PASS** | 10 `<a href>`, 5 destinations. `github.com/JustCode-CruzAlex/Zerado` → **200**; `www.flowforgesoft.com` → **200**; 3 × `mailto:` (no network target); `#main` → resolves to `<main id="main">`. **Zero `href="#"`**, zero unbuilt routes, zero download links. |
| No rendered placeholder | **PASS** | Scanned rendered text *and* HTML for `[placeholder]`, `[confirmar]`, `R$ 00`, `lorem ipsum`, `TODO`, `{{`, `${`, `undefined`, `NaN` — **0 hits**. |
| "Powered by FlowForge" footer | **PASS** | Present, links `https://www.flowforgesoft.com` (200), renders `/flowforge-logo.svg` (200, `image/svg+xml`). **It is the official mark, not a redraw** — `shasum -a 256` of the served file is byte-identical (`3f0c54d9…9829`) to `v3/playbooks/landing-page/assets/flowforge-logo.svg`. |
| Blueprint conformance | **PASS** | 16/16 sections present and ordered; every declared component resolves to a real DOM node; **all 12 sampled colour tokens match the blueprint hex exactly**; all three font families genuinely loaded (`document.fonts.check` true for Orbitron, Space Grotesk, JetBrains Mono — not merely declared in a stack). Breakpoints: mood grid 1/2/4/4, cover grid 3/4/6/6, roadmap 1/2/4/4, phone frames **1 at 375, 2 at 768+** — every one matching the manifest. |
| Honest content | **PASS** | All **4** roadmap phases read `○ Planned` — zero marked done. **7 mockup captions** render visibly, each stating it is a mockup and not a screenshot. The cover disclosure *"Cover tiles are illustrative artwork, not real game covers."* renders. No dates on any phase. JSON-LD makes **no** availability/price/download claim. |

---

## 3. The blocking defect and its fix — re-verified first-hand

**The fix is real.** At 320 / 375 / 768 / 1280 / 1920 the §02 maker-line renders **78 of 78 characters**, string-identical to `content/landing-copy.md`, with `scrollWidth − clientWidth = 0` (no clipping). At 375 it wraps to 6 lines inside the flanking rules; I confirmed this visually as well as by measurement.

**The new gate genuinely fails.** This was the claim I most wanted to test, because a gate that cannot fail is worse than no gate. I re-injected the defect at runtime (forcing `white-space: nowrap` on the maker-line) and ran the *exact* gate logic from `qa/harness/run-core.mjs`:

```
w=320   CLEAN pass=true  (offenders 0)  ||  MUTATED pass=false (offenders 1)  ← KILLED
w=375   CLEAN pass=true  (offenders 0)  ||  MUTATED pass=false (offenders 1)  ← KILLED
w=768   CLEAN pass=true  (offenders 0)  ||  MUTATED pass=false (offenders 1)  ← KILLED
w=1280  CLEAN pass=true  (offenders 0)  ||  MUTATED pass=true  (fits: right=1020 < 1280)
w=1920  CLEAN pass=true  (offenders 0)  ||  MUTATED pass=true  (fits: right=1340 < 1920)
```

The 1280/1920 passes are correct, not gate failures: with the bug re-injected the line still fits those viewports, so there is no overflow to detect.

**And I confirmed the old gate was genuinely a tautology.** In *every* mutated run above — including the three where text was overflowing the viewport — `documentElement.scrollWidth === clientWidth` remained **`true`**. That is direct proof the original assertion could never have caught this defect, and that replacing it with `offenders.length === 0` is a real improvement rather than a cosmetic one.

---

## 4. Findings

### BLOCKING — none.
### MAJOR — none.

### MINOR

**MINOR-1 · Four mockup-chrome labels ellipsis-truncate at 320px**
`.z-frame__title`, `.z-frame__counter` (×2 frames), `.z-tile--B`
At a 320px viewport: `"zerado — library"` loses 17px (86.4% visible), `"412 games · 6 zerado"` loses 22px (85.9%), `"zerado — community (preview)"` loses 24px (89.0%), `"phase 4"` loses 7px, and the `▣ PHYSICAL` tile tag loses 4px. **All clear completely by 360px and 375px.**
*Why it is MINOR, not a repeat of the blocking defect:* the truncation is `text-overflow: ellipsis` — a visible, honest affordance — not the silent destruction of 45 characters that blocked. It occurs only *below* the blueprint's declared 375px support floor, on decorative terminal-mockup chrome, and every fact involved is carried redundantly in page prose. It matters because WCAG 2.2 SC 1.4.10 (Reflow) is specified at 320px, and because the QA report states it tested the 320px reflow floor but did not surface these. Worth a one-line fix (shrink the chrome type or hide the counter under 360px) before launch, not worth holding the page.

**MINOR-2 · Roadmap card bottoms are ragged**
`#roadmap .z-phase` — measured heights at 1280: **194 / 194 / 222 / 194 px**. The grid cells (`.z-roadmap__item`) are correctly equal at 222px, but the bordered card inside is not stretched, so Phase 3 — whose title *"Recommendations & Budget"* wraps to two lines — hangs **28px lower** than its three neighbours. Visible in a section whose whole argument is "four phases, one glance." Fix is one line: `.z-phase { height: 100% }`. This is the one place the page reads improvised rather than designed.

**MINOR-3 · Blueprint under-counts the page's own mandatory link**
`design/blueprint.tokens.json` declares `focusableElementCount: 15` and `linkPolicy.totalLinkInstances: 9`. Measured: **16 focusables and 10 links** at ≥768px (15 at 375, where the header GitHub link is hidden). The uncounted element is the **"Powered by FlowForge" footer link** — which is itself a hard requirement of this deliverable, and is absent from `ctaTargets` entirely. The page is right; the manifest is wrong. `linkPolicy.totalDestinations: 3` also disagrees with its own `ctaTargets` list of 4.

**MINOR-4 · Blueprint says there is no OG image; the build shipped one**
Manifest: `"ogImage": null`, reason *"no such asset exists"*, and `rasterImageCount: 0`. Shipped: `og-card.png`, a real **1200×630** PNG referenced by `og:image` and `twitter:image`. This is an *improvement* — and an honest one: its `og:image:alt` states the card is *"MOCKUP OF THE PLANNED INTERFACE — it is not a real screenshot."* The rendered page still loads zero raster images (confirmed by Lighthouse 100/100), so no budget was harmed. Only the manifest is stale. Note the card cannot be fetch-verified until `zerado.app` is deployed, since `og:image` is correctly an absolute URL.

**MINOR-5 · Blueprint names a stylesheet that does not exist**
`inputs.fontCss` and `typography.fontFiles.css` both point to `site/public/fonts.css`. That file **does not exist** anywhere in the project, and `/fonts.css` returns **404**. No impact: Astro inlined all 6 `@font-face` rules into `index.html` and all five `.woff2` files load. Stale path only.

**MINOR-6 · Blueprint still prescribes the disproven gate**
`constraints.bodyScrollWidthAssertion` still reads *"document.documentElement.scrollWidth === document.documentElement.clientWidth at every breakpoint"* — the exact tautology that let the blocking defect ship. `qa/harness/run-core.mjs` was fixed; the blueprint it was derived from was not. **This is the finding with the longest tail**: the blueprint is declared the build's source of truth, so a future rebuild or a second implementer reading it would faithfully re-adopt the broken assertion. One-line edit; worth doing now while the reason is fresh.

**MINOR-7 · One control sits under the blueprint's own touch-target floor**
The manifest declares `minTouchTargetPx: 48`. The header **"View on GitHub"** link measures **25px** tall (it is styled as a text link, not a button). This is **not an accessibility failure** — WCAG 2.2 SC 2.5.8 requires 24×24px, and axe reported zero target-size violations — but it does not meet the blueprint's own stated 48px. Every other CTA measures 44–52px.

---

## 5. What I verified first-hand vs. what I accepted

**Verified first-hand (I ran it myself, against the live page):**
outline SHA; section order and ids; all 10 links + external HTTP status; forbidden-string scan; roadmap statuses; §11 interactivity count; all 7 mockup captions; cover-grid disclosure, tile count, treatment/hue/tag sequence; footer mark + **byte-level SHA comparison to the official asset**; 12 colour tokens; real font loading; Orbitron floor across breakpoints; co-render on all four states; breakpoint layout at 375/768/1280/1920; per-element overflow at 320/375/768/1280/1920; **the maker-line fix**; **the gate mutation test, including proof the old assertion stayed `true` under mutation**; an independent **axe-core** run at four viewports (0 violations); **computed contrast on all 284 rendered text nodes** (0 failures); focusable counts, focus-ring visibility, positive-tabindex check, touch-target sizes; skip-link focus behaviour; reduced-motion behaviour; 23 verbatim copy strings; Lighthouse raw JSON for all 6 runs.

**Accepted from prior reports (with corroboration, not on faith):**
- *Lighthouse 100/100.* I did not re-run Lighthouse, but I parsed all six raw JSON reports (`lh-{mobile,desktop}-{1,2,3}`): performance **100** in all six, CLS **0**, TBT **0 ms**, LCP 0.4s desktop / 1.5s mobile, Lighthouse 13.4.1, all against `http://localhost:4321/`. Consistent across three runs per form factor. Claim holds.
- *Cross-browser (Firefox/WebKit).* I drove **Chromium only**. Firefox and WebKit results come from `qa/qa-crossbrowser.json` and the committed screenshots. **This is my one genuine coverage gap** — I did not personally reproduce the maker-line fix in Firefox or WebKit.
- *Root-cause of the defect* (`flex-shrink: 0` + `flex-basis: 0%`). I verified the *outcome* exhaustively; I did not re-derive the CSS mechanism.

**One correction to my own process, recorded for honesty:** my first skip-link measurement reported `top: -15` (appearing to fail focus-visibility). That was my error — I measured mid-transition. Re-measured with settle time, the skip link is fully in viewport when focused (16,16, 176×52, cyan on dark, z-index 100). **No defect. I am reporting the correction because a review that hides its own false starts is not auditable.**

---

## 6. Design judgment — does it look designed, or improvised?

**It looks designed. This clears the retro-FUTURE bar, and it clears it on the hard axis — the one most attempts miss.**

The brief was specific and difficult: 1980s *retro-future*, the DeLorean and KITT — 80s objects pointed at the future — and explicitly **not** retro-nostalgia. The failure mode is well-known: VHS grain, faded sepia, dusty chrome, a wink at the viewer. **None of that is here.** There is no grain overlay, no degraded-tape treatment, no "remember the 80s?" wink anywhere on the page. The palette is near-black (`#05060A`) with cyan (`#19E0FF`) and amber (`#FFB000`) at full saturation — these are *lit* colours, emitting rather than fading. That is the correct read of the brief, and getting it right was not guaranteed.

**What is genuinely excellent, specifically:**

- **The terminal mockup (§01).** This is the best thing on the page. It is a convincing TUI — state glyphs and colours per row, a selected row with a `▶` marker, a right-hand rail whose counts actually sum correctly (288 + 9 + 6 + 109 = 412, matching the "412 games" in the chrome bar), and a keybar. Someone who lives in a terminal will believe it. That internal numerical consistency is a detail almost everyone skips, and it is what separates a mockup from a prop.
- **The price card (§07).** The outline predicted this is "the section a skeptical visitor screenshots," and the build earned that. Big amber CRT figures, a bar sparkline, the verdict line in the product's own voice, and an honest caption underneath. It is the most product-convincing component on the page.
- **The `zerado` payoff (§05).** The one cyan glow in the entire page budget is spent on the ZERADO card — the moment the product's name explains itself. Spending your single loudest visual effect on your single most important idea is editorial discipline, not decoration.
- **Restraint.** Zero JavaScript, zero raster images on the rendered page, one animated element on the whole page, and 100/100 Lighthouse on both form factors. The page is fast *because* it is disciplined, not despite it.

**Where it is merely competent, said plainly:**

- **The hero has a large empty right half at ≥1280.** The headline occupies columns 1–7 by design and the space to its right stays empty down to the terminal frame. It reads as intentional asymmetry rather than a mistake, but it is the one composition where the desktop layout is carrying less than it could. A grid-horizon or scanner element in that void would earn it.
- **Several cover tiles read as empty.** Of the 12, roughly four (the brushed-panel diagonals and the darker phosphor arcs) are close to flat dark rectangles. A shelf should feel dense and various; a third of it currently feels unfinished rather than art-directed. See the ruling below.
- **MINOR-2's ragged roadmap row** is the single place a careful eye will catch improvisation.

**Bottom line for the customer:** this is a page you can put in front of strangers today. It is not a template with a colour swap — the terminal frame, the price card and the state chips are bespoke components built for this product, and the aesthetic is a correct and disciplined reading of a brief that is easy to get wrong.

---

## 7. Ruling on the two judgment calls

### Call 1 — No real game cover art. **CORRECT, and I would have made the same call. Execution is the weaker half.**

Using real covers would have put third-party copyrighted artwork on a commercial page carrying an affiliate model — while `decisions.md` already flags that IGDB's licence is free for **non-commercial** use only and that commercial use is an unresolved risk you must settle directly with IGDB. Shipping real covers would have manufactured a legal exposure on the very page whose footer discloses commercial intent. It would also have broken the zero-raster-image budget that produces the 100/100 scores. The call protects you, and the visible disclosure keeps it honest.

**But it only half-satisfies your actual request.** You asked for "nice game images," and what shipped is *nice images that are not game images*. The tiles are attractive abstract artwork; they do not read as a game shelf. The gap is not the decision — it is the execution: about four of the twelve tiles are near-empty dark rectangles. **Recommendation:** keep the decision, raise the floor on the four weakest tiles (more treatment density, a suggestion of title-block geometry) so the grid reads as a *collection* rather than a gradient study. If you want true covers later, the honest routes are settling IGDB commercial licensing, or using store-provided assets under each store's affiliate-programme terms.

### Call 2 — No synthwave audio. **CORRECT, and not a close call.**

All three stated reasons hold independently: browsers block autoplay with sound (so it would frequently not play at all), sound-on-load is a genuine accessibility failure, and an audio asset would dominate the byte budget of a page that currently ships zero raster images and zero JavaScript. Any one of these justifies the decision; together they make it obvious.

The part I want to credit specifically is that the feeling was **carried elsewhere rather than dropped**. The scanner sweep, the grid horizon, the neon-on-black palette and the amber CRT readouts do the work the audio would have done. That is the difference between declining a request and *translating* it — and the register you described is genuinely present on the page without a single byte of audio. If you want the sound later, the correct pattern is a user-initiated, clearly-labelled, off-by-default toggle — never autoplay.

---

## 8. Still required from you before this goes live

> **Addendum, 2026-08-25 — items 1 and 2 have moved since this review was written.**
> This section is left exactly as it was recorded at rev A; a review verdict is
> evidence and is not edited after the fact. What follows is what happened next.
>
> - **Item 1 is resolved, in the opposite direction to the one this review
>   assumed.** It read the tension as an *undersell* and expected it to close by
>   marking Phase 1 complete. The founder confirmed on 2026-08-25 that nothing is
>   built — no TUI, no Go code — so it closed the other way: Phase 1 now reads
>   *In progress* (not done), and §06/§14's *"Steam is built and works" /
>   "Steam syncs today"* moved to the future tense, because with no Go code they
>   were false rather than merely in tension. See `docs/REDACTIONS.md` §3 and
>   `design/blueprint.md` §1.0, Amendment 1.
> - **Item 2 is unchanged:** the roadmap is still undated, and whether it carries
>   dates is still unanswered.
> - **Item 3 is unchanged:** every mockup still carries its visible caption.
>
> The counts and the verdict above are untouched.

These are the same three open values recorded in `decisions.md`. The page currently renders the honest fallback for each, so it is **publishable as-is** — none of these is a launch blocker, but each one currently costs you a little credibility.

1. **Which phase is honestly done today.** All four phases read *"Planned"* — including Phase 1, even though §14 states *"Steam is built and works."* That is an internal tension a careful reader will notice: the FAQ says something works, the roadmap says nothing is started. The fallback is the *honest* choice given no answer, but it undersells you. `blueprint.tokens.json` records `statusChangeIsTokenOnly: true`, so marking a phase complete is a token swap, not a rebuild.
2. **Whether the roadmap carries dates.** Currently undated, ordered phases with a status marker only. Undated is the safer default pre-launch; confirm you want it.
3. **A real screenshot of the TUI running.** Every mockup is currently, correctly, captioned as a mockup — seven captions, all rendering. Replacing §01's frame with a real screenshot is the single highest-value upgrade available to this page, because §01 is doing the persuading.

**Two further items worth your attention, from the contracts rather than the page:**

4. **IGDB commercial licensing** (`decisions.md` risk 1) — unresolved, and it gates cover art and *sinopse* in the product itself. Does not block this page.
5. **`alex@flowforgesoft.com` is a real inbox on a public page** and will be scraped. That is the accepted cost of a backend-free waitlist; a form service later removes it.

---

## VERDICT: **GOLDEN**

**0 blocking · 0 major · 7 minor.**

Every ratified contract holds: the outline SHA matches, all sixteen sections are present and in order, all six ratified decisions are reflected in the built page, the copy is wired verbatim, the tokens match, and the brand rules — co-render, display-type floor, contrast — are met by measurement, not by assertion. Zero axe violations and zero contrast failures across 284 text nodes at four viewports. The blocking defect is genuinely fixed, and — the check that mattered most — **the gate that missed it can now be watched to fail**.

The seven minor findings are worth an hour: two are visual polish (the ragged roadmap row is the one a customer will see), and five are drift in the blueprint manifest rather than defects in the page. **MINOR-6 is the one I would fix today**, because a design contract that still prescribes a gate proven incapable of failing will quietly re-introduce this exact class of defect on the next build.

*Reviewed by fft-code-reviewer. Sole GOLDEN arbiter for this run. Verdict based on first-hand verification of the live page; the one coverage gap — Firefox and WebKit reproduction — is stated explicitly in §5.*
