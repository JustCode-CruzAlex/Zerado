import { chromium } from 'playwright';
const OUT='/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa/screenshots';
const b = await chromium.launch();
for (const w of [320,375,768,1280]) {
  const ctx = await b.newContext({ viewport:{width:w,height:900} });
  const p = await ctx.newPage();
  await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
  const el = p.locator('#maker-line');
  await el.screenshot({ path:`${OUT}/DEFECT-maker-line-${w}.png` });
  const d = await p.evaluate(() => {
    const t = document.querySelector('.z-maker-line__text');
    const r = t.getBoundingClientRect();
    const vw = document.documentElement.clientWidth;
    const full = t.textContent.trim().replace(/\s+/g,' ');
    // find how much is visible: walk text nodes / use Range per character
    const range = document.createRange();
    const tn = [];
    const walk = document.createTreeWalker(t, NodeFilter.SHOW_TEXT);
    let n; while ((n = walk.nextNode())) tn.push(n);
    let visible = '';
    for (const node of tn) {
      for (let i=0;i<node.length;i++){
        range.setStart(node,i); range.setEnd(node,i+1);
        const rr = range.getBoundingClientRect();
        if (rr.right <= vw + 0.5) visible += node.textContent[i];
      }
    }
    return { viewportWidth: vw, textLeft: Math.round(r.left), textRight: Math.round(r.right),
             textWidth: Math.round(r.width), overshootPx: Math.round(r.right - vw),
             fullText: full, visibleText: visible.replace(/\s+/g,' ').trim(),
             clippedChars: full.length - visible.replace(/\s+/g,' ').trim().length };
  });
  console.log(`\n== ${w}px ==`);
  console.log(`  text rect: left=${d.textLeft} right=${d.textRight} width=${d.textWidth}  OVERSHOOT past viewport = ${d.overshootPx}px`);
  console.log(`  FULL    : "${d.fullText}"`);
  console.log(`  VISIBLE : "${d.visibleText}"`);
  console.log(`  clipped characters: ${d.clippedChars}`);
  await ctx.close();
}
await b.close();
