import { chromium } from 'playwright';
import AxeBuilder from '@axe-core/playwright';
import { createServer } from 'node:http';
import { readFileSync, existsSync, mkdirSync } from 'node:fs';
import { extname, join } from 'node:path';

const DIST = process.argv[2];
const LABEL = process.argv[3];
const OUT = process.argv[4];
mkdirSync(OUT, { recursive: true });

const MIME = { '.html':'text/html', '.css':'text/css', '.js':'text/javascript', '.svg':'image/svg+xml',
  '.png':'image/png', '.jpg':'image/jpeg', '.webp':'image/webp', '.avif':'image/avif', '.woff2':'font/woff2', '.json':'application/json' };

const server = createServer((req, res) => {
  let p = decodeURIComponent(req.url.split('?')[0]);
  if (p === '/') p = '/index.html';
  const f = join(DIST, p);
  if (!existsSync(f)) { res.writeHead(404); return res.end('404'); }
  res.writeHead(200, { 'Content-Type': MIME[extname(f)] ?? 'application/octet-stream' });
  res.end(readFileSync(f));
});
await new Promise(r => server.listen(4321, r));

const browser = await chromium.launch();
const VIEWPORTS = [[375,812],[768,1024],[1280,800],[1920,1080]];
const report = { label: LABEL, viewports: {} };

for (const [w,h] of VIEWPORTS) {
  const ctx = await browser.newContext({ viewport: { width: w, height: h }, deviceScaleFactor: 1 });
  const page = await ctx.newPage();
  const external = [];
  page.on('request', r => { const u = r.url(); if (!u.startsWith('http://localhost:4321') && !u.startsWith('data:')) external.push(u); });
  await page.goto('http://localhost:4321/', { waitUntil: 'networkidle' });
  await page.evaluate(async () => {
    await new Promise(res => { let y=0; const t=setInterval(()=>{ window.scrollBy(0,900); y+=900; if(y>=document.body.scrollHeight){clearInterval(t);res();} },30); });
    window.scrollTo(0,0);
  });
  await page.waitForTimeout(400);
  // The cover grid is section 6 of 16 — a viewport shot of the fold never
  // contains it, which is how a "screenshot at four viewports" can be evidence
  // of nothing. Capture the section itself.
  await page.screenshot({ path: join(OUT, `${LABEL}-${w}-fold.png`), fullPage: false });
  const section = page.locator('#one-collection');
  await section.scrollIntoViewIfNeeded();
  await page.waitForTimeout(500);
  await section.screenshot({ path: join(OUT, `${LABEL}-${w}-grid.png`) });
  const axe = await new AxeBuilder({ page }).withTags(['wcag2a','wcag2aa','wcag21a','wcag21aa']).analyze();
  const grid = await page.evaluate(() => {
    const g = document.querySelector('.z-cover-grid');
    if (!g) return null;
    const kids = [...g.children].map(c => c.getBoundingClientRect());
    const tops = [...new Set(kids.map(r => Math.round(r.top)))];
    const rows = tops.map(t => kids.filter(r => Math.round(r.top) === t).length);
    return { tiles: kids.length, rows, widths: [...new Set(kids.map(r=>Math.round(r.width)))], heights: [...new Set(kids.map(r=>Math.round(r.height)))] };
  });
  report.viewports[w] = {
    axeViolations: axe.violations.map(v => ({ id: v.id, impact: v.impact, nodes: v.nodes.length })),
    externalRequests: [...new Set(external)],
    grid
  };
  await ctx.close();
}

// CLS + LCP over a real scroll, no throttling — a floor measurement, not Lighthouse.
const ctx = await browser.newContext({ viewport: { width: 1280, height: 800 } });
const page = await ctx.newPage();
await page.addInitScript(() => {
  window.__cls = 0;
  new PerformanceObserver(l => { for (const e of l.getEntries()) if (!e.hadRecentInput) window.__cls += e.value; })
    .observe({ type: 'layout-shift', buffered: true });
});
await page.goto('http://localhost:4321/', { waitUntil: 'networkidle' });
await page.evaluate(async () => {
  await new Promise(res => { let y=0; const t=setInterval(()=>{ window.scrollBy(0,600); y+=600; if(y>=document.body.scrollHeight){clearInterval(t);res();} },40); });
});
await page.waitForTimeout(1200);
report.cls = await page.evaluate(() => window.__cls);
report.transferredBytes = await page.evaluate(() =>
  performance.getEntriesByType('resource').reduce((n, r) => n + (r.transferSize || 0), 0));
report.resourceCount = await page.evaluate(() => performance.getEntriesByType('resource').length);
await ctx.close();
await browser.close();
server.close();
console.log(JSON.stringify(report, null, 2));
