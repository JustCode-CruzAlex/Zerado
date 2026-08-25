import { chromium } from 'playwright';
const b=await chromium.launch();
const ctx=await b.newContext({viewport:{width:1280,height:1000}});
const p=await ctx.newPage(); await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
const d=await p.evaluate(()=>{
  const rep=[];
  document.querySelectorAll('.z-chip').forEach(c=>{
    const g=c.querySelector('[class*="glyph"]'); const cs=getComputedStyle(c);
    rep.push({kind:'StateChip', glyph:g?g.textContent.trim():null, glyphAriaHidden:g?.getAttribute('aria-hidden'),
      label:c.textContent.replace(g?g.textContent:'','').trim().replace(/\s+/g,' ').slice(0,40),
      color:getComputedStyle(g||c).color});
  });
  document.querySelectorAll('.z-status-marker,[class*="status-marker"]').forEach(c=>{
    const g=c.querySelector('[class*="glyph"]');
    rep.push({kind:'StatusMarker', glyph:g?g.textContent.trim():null, glyphAriaHidden:g?.getAttribute('aria-hidden'),
      label:c.textContent.replace(g?g.textContent:'','').trim().replace(/\s+/g,' ').slice(0,30),
      color:getComputedStyle(g||c).color});
  });
  const rows=[...document.querySelectorAll('.z-trow')].slice(0,6).map(r=>{
    const g=r.querySelector('.z-trow__glyph'); const st=r.querySelector('[class*="state"],[class*="status"]');
    return {glyph:g?.textContent.trim(), ariaHidden:g?.getAttribute('aria-hidden'),
      full:r.innerText.replace(/\s+/g,' ').trim().slice(0,80), color:g?getComputedStyle(g).color:null};
  });
  return {rep,rows};
});
console.log('chips/markers:',d.rep.length);
d.rep.forEach(r=>console.log(`  ${r.kind}: glyph="${r.glyph}" ariaHidden=${r.glyphAriaHidden} color=${r.color} label="${r.label}"`));
console.log('terminal rows (first 6):');
d.rows.forEach(r=>console.log(`  glyph="${r.glyph}" ariaHidden=${r.ariaHidden} color=${r.color} :: ${r.full}`));
await b.close();
