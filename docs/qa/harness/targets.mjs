import { chromium } from 'playwright';
const b=await chromium.launch();
for (const w of [375,1280]) {
  const ctx=await b.newContext({viewport:{width:w,height:812}});
  const p=await ctx.newPage(); await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
  const d=await p.evaluate(()=>[...document.querySelectorAll('a[href],summary,button')].map(e=>{
    const r=e.getBoundingClientRect();
    return {tag:e.tagName.toLowerCase(),text:(e.textContent||'').trim().replace(/\s+/g,' ').slice(0,28),
      w:Math.round(r.width),h:Math.round(r.height)};}));
  console.log(`\n== touch targets @${w} (WCAG 2.5.8 AA min 24x24) ==`);
  d.forEach(t=>{const ok=t.w>=24&&t.h>=24;console.log(`  ${ok?'OK  ':'FAIL'} ${String(t.w).padStart(4)}x${String(t.h).padStart(3)}  <${t.tag}> "${t.text}"`);});
  console.log(`  under 24px: ${d.filter(t=>t.w<24||t.h<24).length} | under 44px: ${d.filter(t=>t.w<44||t.h<44).length}`);
  await ctx.close();
}
await b.close();
