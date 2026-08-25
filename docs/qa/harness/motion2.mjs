import { chromium } from 'playwright';
const b = await chromium.launch();
for (const mode of ['no-preference','reduce']) {
  const ctx=await b.newContext({viewport:{width:1280,height:900},reducedMotion:mode});
  const p=await ctx.newPage(); await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
  await p.waitForTimeout(400);
  const d=await p.evaluate(()=>{
    const tracks=[...document.querySelectorAll('.z-scanner-track')].map(t=>({
      cls:t.className, section:t.closest('section,footer')?.id||t.closest('footer')?.tagName||'?',
      inFigure: !!t.closest('figure') }));
    const faq=[...document.querySelectorAll('.z-faq-item__summary, details summary')].slice(0,2).map(s=>{
      const cs=getComputedStyle(s); return {transProp:cs.transitionProperty,transDur:cs.transitionDuration};});
    const content=[...document.querySelectorAll('details')].slice(0,1).map(dd=>{
      const inner=dd.querySelector(':scope > *:not(summary)');
      const cs=inner?getComputedStyle(inner):null;
      return {innerTag:inner?.tagName, innerCls:(inner?.className||'').toString(),
              transProp:cs?.transitionProperty, transDur:cs?.transitionDuration,
              animName:cs?.animationName, animDur:cs?.animationDuration,
              height:cs?.height, overflow:cs?.overflow};});
    const supportsDetailsContent = CSS.supports('interpolate-size','allow-keywords');
    return {tracks,faq,content,supportsDetailsContent};
  });
  console.log(`\n== ${mode} ==`);
  console.log('scanner tracks by section:'); d.tracks.forEach((t,i)=>console.log(`  [${i}] section=${t.section} inFigure=${t.inFigure} cls=${t.cls}`));
  console.log('faq summary transitions:',JSON.stringify(d.faq));
  console.log('faq content panel:',JSON.stringify(d.content));
  console.log('interpolate-size supported:',d.supportsDetailsContent);
  // focus ring under this mode
  await p.keyboard.press('Tab'); await p.keyboard.press('Tab'); await p.waitForTimeout(250);
  const fr=await p.evaluate(()=>{const e=document.activeElement;const cs=getComputedStyle(e);
    return{tag:e.tagName,text:(e.textContent||'').trim().slice(0,25),outlineWidth:cs.outlineWidth,outlineStyle:cs.outlineStyle,outlineColor:cs.outlineColor,outlineOffset:cs.outlineOffset};});
  console.log('focus ring on 2nd tab stop:',JSON.stringify(fr));
  await ctx.close();
}
await b.close();
