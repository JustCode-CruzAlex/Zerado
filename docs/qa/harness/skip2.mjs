import { chromium } from 'playwright';
const b = await chromium.launch();
const ctx = await b.newContext({ viewport:{width:1280,height:900} });
const p = await ctx.newPage();
await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
await p.keyboard.press('Tab');
for (const ms of [0,50,150,400,1000]) {
  if (ms) await p.waitForTimeout(ms === 50 ? 50 : ms - (ms===150?50:(ms===400?150:400)));
  const s = await p.evaluate(()=>{
    const a=document.activeElement; const cs=getComputedStyle(a); const r=a.getBoundingClientRect();
    return { top:cs.top, y:Math.round(r.y), matchesFV:a.matches(':focus-visible'), matchesF:a.matches(':focus'),
             intersects: r.bottom>0 && r.top<window.innerHeight };
  });
  console.log(`t=${ms}ms -> computed top=${s.top} rectY=${s.y} :focus-visible=${s.matchesFV} :focus=${s.matchesF} onScreen=${s.intersects}`);
}
await p.screenshot({path:'/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa/screenshots/skiplink-focused-1280.png'});
// also at 375
await ctx.close();
const c2 = await b.newContext({viewport:{width:375,height:812}});
const p2 = await c2.newPage();
await p2.goto('http://localhost:4321/',{waitUntil:'networkidle'});
await p2.keyboard.press('Tab'); await p2.waitForTimeout(500);
const s2 = await p2.evaluate(()=>{const a=document.activeElement;const r=a.getBoundingClientRect();
  return {top:getComputedStyle(a).top,y:Math.round(r.y),x:Math.round(r.x),w:Math.round(r.width),
          onScreen: r.bottom>0&&r.top<window.innerHeight&&r.right>0&&r.left<window.innerWidth};});
console.log('375px focused skip link:', JSON.stringify(s2));
await p2.screenshot({path:'/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa/screenshots/skiplink-focused-375.png'});
await b.close();
