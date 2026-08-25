import { chromium } from 'playwright';
const b = await chromium.launch();
const ctx = await b.newContext({ viewport:{width:1280,height:900} });
const p = await ctx.newPage();
await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});

const pre = await p.evaluate(()=>{
  const a=document.querySelector('a[href="#main"]'); const cs=getComputedStyle(a); const r=a.getBoundingClientRect();
  return {unfocused:{top:cs.top,left:cs.left,position:cs.position,transform:cs.transform,width:cs.width,height:cs.height,
    clip:cs.clip,clipPath:cs.clipPath,overflow:cs.overflow,opacity:cs.opacity,visibility:cs.visibility,zIndex:cs.zIndex,
    rect:{x:Math.round(r.x),y:Math.round(r.y),w:Math.round(r.width),h:Math.round(r.height)}}};
});
console.log('UNFOCUSED:', JSON.stringify(pre.unfocused,null,1));

await p.keyboard.press('Tab');
const post = await p.evaluate(()=>{
  const a=document.activeElement; const cs=getComputedStyle(a); const r=a.getBoundingClientRect();
  return {tag:a.tagName,href:a.getAttribute('href'),top:cs.top,left:cs.left,position:cs.position,transform:cs.transform,
    opacity:cs.opacity,visibility:cs.visibility,zIndex:cs.zIndex,
    rect:{x:Math.round(r.x),y:Math.round(r.y),w:Math.round(r.width),h:Math.round(r.height)},
    scrollY:Math.round(window.scrollY), viewportH: window.innerHeight,
    intersectsViewport: r.bottom>0 && r.top<window.innerHeight && r.right>0 && r.left<window.innerWidth };
});
console.log('FOCUSED  :', JSON.stringify(post,null,1));
await p.screenshot({path:'/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa/screenshots/DEFECT-skiplink-focused-1280.png'});

// where does focus go after activating skip link
await p.keyboard.press('Enter');
await p.waitForTimeout(300);
await p.keyboard.press('Tab');
const next = await p.evaluate(()=>{const e=document.activeElement;const r=e.getBoundingClientRect();
  return {tag:e.tagName.toLowerCase(),href:e.getAttribute('href'),text:(e.textContent||'').trim().slice(0,40),
          y:Math.round(r.y), inMain: !!e.closest('#main'), inNav: !!e.closest('nav'), scrollY:Math.round(window.scrollY)};});
console.log('AFTER skip+Enter+Tab -> focus:', JSON.stringify(next));

// where is #main relative to nav
const struct = await p.evaluate(()=>{
  const m=document.querySelector('#main'); const nav=document.querySelector('nav');
  return { mainTag:m?.tagName, mainId:m?.id, mainTabindex:m?.getAttribute('tabindex'),
           navInsideMain: !!nav?.closest('#main'),
           navBeforeMain: nav ? !!(m.compareDocumentPosition(nav) & Node.DOCUMENT_POSITION_PRECEDING) : null,
           mainRole:m?.getAttribute('role'), mainTagName:m?.tagName };
});
console.log('STRUCT:', JSON.stringify(struct));
await b.close();
