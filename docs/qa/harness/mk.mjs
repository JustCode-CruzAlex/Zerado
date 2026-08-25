import { chromium } from 'playwright';
const b = await chromium.launch();
for (const w of [375,768,1280,1920]) {
  const ctx = await b.newContext({ viewport:{width:w,height:900} });
  const p = await ctx.newPage();
  await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
  const info = await p.evaluate(() => {
    const sec = document.querySelector('#maker-line');
    const row = sec.querySelector('.z-maker-line__row');
    const kids = [...row.children].map(c => {
      const r = c.getBoundingClientRect();
      const cs = getComputedStyle(c);
      return { tag:c.tagName.toLowerCase(), cls:(c.className||'').toString(), text:(c.textContent||'').trim().slice(0,60),
               left:Math.round(r.left), right:Math.round(r.right), width:Math.round(r.width),
               flex:cs.flex, whiteSpace:cs.whiteSpace };
    });
    const rr = row.getBoundingClientRect();
    return { html: sec.outerHTML.slice(0,700), rowLeft:Math.round(rr.left), rowRight:Math.round(rr.right), rowW:Math.round(rr.width),
             rowScrollW: row.scrollWidth, rowClientW: row.clientWidth,
             rowCS: { display:getComputedStyle(row).display, flexWrap:getComputedStyle(row).flexWrap, gap:getComputedStyle(row).gap },
             kids, bodyOverflowX: getComputedStyle(document.body).overflowX,
             htmlOverflowX: getComputedStyle(document.documentElement).overflowX,
             bodyScrollW: document.body.scrollWidth, bodyClientW: document.body.clientWidth };
  });
  console.log(`\n===== ${w}px =====`);
  console.log('body overflowX=', info.bodyOverflowX, ' html overflowX=', info.htmlOverflowX, ' bodyScrollW=', info.bodyScrollW, 'bodyClientW=', info.bodyClientW);
  console.log('row rect left/right/width:', info.rowLeft, info.rowRight, info.rowW, 'scrollW', info.rowScrollW, 'clientW', info.rowClientW, JSON.stringify(info.rowCS));
  console.log('children:', JSON.stringify(info.kids, null, 1));
  if (w===375) console.log('HTML:', info.html);
  await ctx.close();
}
await b.close();
