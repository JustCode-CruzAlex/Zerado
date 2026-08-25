import { chromium } from 'playwright';
const b=await chromium.launch();
const ctx=await b.newContext({viewport:{width:375,height:812}});
const p=await ctx.newPage(); await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
const navLink=await p.evaluate(()=>{
  const a=[...document.querySelectorAll('a')].find(x=>/github/i.test(x.href)&&x.closest('header,nav'));
  if(!a)return null; const cs=getComputedStyle(a); const r=a.getBoundingClientRect();
  return {display:cs.display,visibility:cs.visibility,opacity:cs.opacity,w:r.width,h:r.height,
    clip:cs.clip,clipPath:cs.clipPath,position:cs.position,parentDisplay:getComputedStyle(a.parentElement).display,
    parentCls:(a.parentElement.className||'').toString()};
});
console.log('nav GitHub link @375:',JSON.stringify(navLink,null,1));
// tab order at 375
await p.evaluate(()=>window.scrollTo(0,0));
const stops=[];
for(let i=0;i<20;i++){
  await p.keyboard.press('Tab'); await p.waitForTimeout(160);
  const s=await p.evaluate(()=>{const e=document.activeElement; if(!e||e===document.body)return{none:true};
    const r=e.getBoundingClientRect(); const cs=getComputedStyle(e);
    return{tag:e.tagName.toLowerCase(),text:(e.textContent||'').trim().replace(/\s+/g,' ').slice(0,30),
      href:e.getAttribute('href'),w:Math.round(r.width),h:Math.round(r.height),
      y:Math.round(r.y),display:cs.display,
      onScreen:r.width>0&&r.height>0&&r.bottom>0&&r.top<window.innerHeight};});
  if(s.none){console.log(`tab[${i}] left document`);break;}
  stops.push(s);
  console.log(`tab[${i}] ${s.w}x${s.h} onScreen=${s.onScreen} display=${s.display} <${s.tag}> "${s.text}"`);
  if(stops.length>17)break;
}
console.log('total stops @375:',stops.length,'| zero-size focusable:',stops.filter(s=>s.w===0||s.h===0).length);

// footer contact spacing (WCAG 2.5.8 spacing exception)
const sp=await p.evaluate(()=>{
  const targets=[...document.querySelectorAll('a[href],summary,button')].map(e=>{const r=e.getBoundingClientRect();
    return {t:(e.textContent||'').trim().slice(0,24),cx:r.x+r.width/2,cy:r.y+window.scrollY+r.height/2,w:r.width,h:r.height};});
  const me=targets.find(t=>/alex@flowforgesoft/.test(t.t));
  const others=targets.filter(t=>t!==me&&t.w>0);
  const dists=others.map(o=>({t:o.t,d:Math.round(Math.hypot(o.cx-me.cx,o.cy-me.cy))})).sort((a,b)=>a.d-b.d).slice(0,3);
  return {me:{w:Math.round(me.w),h:Math.round(me.h)},nearest:dists};
});
console.log('footer contact link:',JSON.stringify(sp));
await b.close();
