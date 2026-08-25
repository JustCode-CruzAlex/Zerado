import { chromium } from 'playwright';
import { AxeBuilder } from '@axe-core/playwright';
import fs from 'node:fs';

const URL = 'http://localhost:4321/';
const OUT = '/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa';
const VIEWPORTS = [
  { name: '375',  width: 375,  height: 812 },
  { name: '768',  width: 768,  height: 1024 },
  { name: '1280', width: 1280, height: 800 },
  { name: '1920', width: 1920, height: 1080 },
];

const results = { url: URL, generatedAt: new Date().toISOString(), viewports: {} };

const browser = await chromium.launch();

for (const vp of VIEWPORTS) {
  const ctx = await browser.newContext({ viewport: { width: vp.width, height: vp.height }, deviceScaleFactor: 1 });
  const page = await ctx.newPage();
  await page.goto(URL, { waitUntil: 'networkidle' });
  await page.waitForTimeout(400);

  const R = {};

  // --- 1. screenshot
  await page.screenshot({ path: `${OUT}/screenshots/${vp.name}-fullpage.png`, fullPage: true });

  // --- 3. horizontal overflow
  R.overflow = await page.evaluate(() => {
    const de = document.documentElement;
    const offenders = [];
    const vw = de.clientWidth;
    document.querySelectorAll('*').forEach(el => {
      const r = el.getBoundingClientRect();
      // aria-hidden decoration (e.g. the rotated PHASE 4 watermark) is clipped by design
      // and is not content loss — exclude it so the gate reports real defects only.
      if (el.closest('[aria-hidden="true"]')) return;
      if (r.width > 0 && (r.right > vw + 1 || r.left < -1)) {
        offenders.push({
          tag: el.tagName.toLowerCase(),
          cls: (el.className && el.className.baseVal !== undefined ? el.className.baseVal : el.className || '').toString().slice(0,90),
          id: el.id || null, left: Math.round(r.left), right: Math.round(r.right)
        });
      }
    });
    // intentionally scrollable containers
    const scrollers = [...document.querySelectorAll('*')].filter(el => el.scrollWidth > el.clientWidth + 1)
      .map(el => ({ tag: el.tagName.toLowerCase(), cls: (el.className||'').toString().slice(0,90), id: el.id||null,
                    scrollWidth: el.scrollWidth, clientWidth: el.clientWidth,
                    tabindex: el.getAttribute('tabindex'), overflowX: getComputedStyle(el).overflowX }));
    return {
      scrollWidth: de.scrollWidth, clientWidth: de.clientWidth,
      bodyScrollWidth: document.body.scrollWidth, bodyClientWidth: document.body.clientWidth,
      // NOT `de.scrollWidth === de.clientWidth`: `overflow-x: hidden` on html/body pins those
      // two values together, making that assertion a tautology that can never fail. It passed at
      // all four viewports while §02 was destroying 58% of its own text (QA BLOCKING-1).
      // Gate on per-element overflow instead — the thing that actually detects content loss.
      pass: offenders.length === 0,
      documentScrollEqualsClient: de.scrollWidth === de.clientWidth,
      offenders: offenders.slice(0, 25), offenderCount: offenders.length,
      scrollers
    };
  });

  // --- 2. axe
  const axe = await new AxeBuilder({ page })
    .withTags(['wcag2a','wcag2aa','wcag21a','wcag21aa','wcag22aa'])
    .analyze();
  R.axe = {
    violations: axe.violations.map(v => ({
      id: v.id, impact: v.impact, help: v.help, helpUrl: v.helpUrl,
      nodes: v.nodes.map(n => ({ target: n.target, html: n.html.slice(0,220), failureSummary: (n.failureSummary||'').slice(0,300) }))
    })),
    violationCount: axe.violations.length,
    passCount: axe.passes.length,
    incompleteCount: axe.incomplete.length,
    incomplete: axe.incomplete.map(v => ({ id: v.id, impact: v.impact, nodes: v.nodes.length, help: v.help })),
    inapplicableCount: axe.inapplicable.length
  };

  results.viewports[vp.name] = R;
  console.log(`[${vp.name}] overflow pass=${R.overflow.pass} scrollW=${R.overflow.scrollWidth} clientW=${R.overflow.clientWidth} | axe violations=${R.axe.violationCount} passes=${R.axe.passCount} incomplete=${R.axe.incompleteCount}`);
  for (const v of R.axe.violations) console.log(`    VIOLATION ${v.id} [${v.impact}] x${v.nodes.length} :: ${v.nodes[0].target}`);
  if (!R.overflow.pass) console.log(`    OFFENDERS(${R.overflow.offenderCount}):`, JSON.stringify(R.overflow.offenders.slice(0,6)));
  if (R.overflow.scrollers.length) console.log(`    SCROLLERS:`, JSON.stringify(R.overflow.scrollers));

  await ctx.close();
}

await browser.close();
fs.writeFileSync(`${OUT}/qa-raw-results.json`, JSON.stringify(results, null, 2));
console.log('WROTE qa-raw-results.json');
