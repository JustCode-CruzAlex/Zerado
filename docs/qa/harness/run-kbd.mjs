import { chromium } from 'playwright';
const b = await chromium.launch();
const ctx = await b.newContext({ viewport:{width:1280,height:900} });
const p = await ctx.newPage();
await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});

// DOM order of focusables
const domOrder = await p.evaluate(() => {
  const sel = 'a[href], button, input, select, textarea, summary, [tabindex]:not([tabindex="-1"])';
  return [...document.querySelectorAll(sel)].map((e,i)=>({
    i, tag:e.tagName.toLowerCase(), text:(e.textContent||'').trim().replace(/\s+/g,' ').slice(0,45),
    href:e.getAttribute('href')||null, tabindex:e.getAttribute('tabindex')
  }));
});
console.log(`DOM focusable candidates: ${domOrder.length}`);
domOrder.forEach(d=>console.log(`  dom[${d.i}] <${d.tag}> ti=${d.tabindex} href=${d.href} :: "${d.text}"`));

// Tab traversal
await p.evaluate(()=>document.body.focus());
const seen = []; const MAX = 40;
for (let i=0;i<MAX;i++){
  await p.keyboard.press('Tab');
  const info = await p.evaluate(()=>{
    const e = document.activeElement;
    if(!e || e===document.body) return {none:true};
    const cs = getComputedStyle(e);
    const r = e.getBoundingClientRect();
    // focus-visible detection
    let fv=false; try { fv = e.matches(':focus-visible'); } catch(_){}
    return { tag:e.tagName.toLowerCase(), text:(e.textContent||'').trim().replace(/\s+/g,' ').slice(0,45),
      href:e.getAttribute('href')||null, tabindex:e.getAttribute('tabindex'),
      focusVisible: fv,
      outlineStyle: cs.outlineStyle, outlineWidth: cs.outlineWidth, outlineColor: cs.outlineColor,
      outlineOffset: cs.outlineOffset, boxShadow: cs.boxShadow.slice(0,60),
      rect:{x:Math.round(r.x),y:Math.round(r.y),w:Math.round(r.width),h:Math.round(r.height)},
      inViewport: r.width>0 && r.height>0 };
  });
  if (info.none) { console.log(`tab[${i}] -> left document (body)`); break; }
  const key = info.tag+'|'+info.href+'|'+info.text;
  if (seen.length && seen[0].key === key && i>2) { console.log(`tab[${i}] -> WRAPPED to first (no trap)`); break; }
  seen.push({key, info});
  console.log(`tab[${i}] <${info.tag}> ti=${info.tabindex} href=${info.href} fv=${info.focusVisible} outline=${info.outlineWidth} ${info.outlineStyle} ${info.outlineColor} off=${info.outlineOffset} rect=${JSON.stringify(info.rect)} :: "${info.text}"`);
}
console.log(`\nTotal tab stops before wrap: ${seen.length}`);

// skip link behaviour
await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
await p.keyboard.press('Tab');
const skip = await p.evaluate(()=>{
  const e=document.activeElement; const r=e.getBoundingClientRect(); const cs=getComputedStyle(e);
  return { tag:e.tagName.toLowerCase(), href:e.getAttribute('href'), text:e.textContent.trim(),
           rect:{x:Math.round(r.x),y:Math.round(r.y),w:Math.round(r.width),h:Math.round(r.height)},
           visible: r.width>0&&r.height>0&&cs.visibility!=='hidden'&&cs.opacity!=='0',
           clip:cs.clip, clipPath:cs.clipPath, position:cs.position };
});
console.log('\nFIRST TAB STOP (skip link):', JSON.stringify(skip));
await p.keyboard.press('Enter');
await p.waitForTimeout(300);
const afterSkip = await p.evaluate(()=>({ hash:location.hash, scrollY:Math.round(window.scrollY),
  active:document.activeElement.tagName.toLowerCase()+'#'+(document.activeElement.id||''),
  mainExists: !!document.querySelector('#main'), mainTabindex: document.querySelector('#main')?.getAttribute('tabindex') }));
console.log('AFTER Enter on skip link:', JSON.stringify(afterSkip));

// details keyboard
await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
const detailsRes = await p.evaluate(()=>({count:document.querySelectorAll('details').length}));
console.log(`\n<details> count: ${detailsRes.count}`);
const sums = await p.locator('details > summary').all();
let dOK=0, dFail=[];
for (let i=0;i<sums.length;i++){
  await sums[i].focus();
  const before = await p.evaluate(i=>document.querySelectorAll('details')[i].open, i);
  await p.keyboard.press('Enter');
  await p.waitForTimeout(280);
  const afterOpen = await p.evaluate(i=>document.querySelectorAll('details')[i].open, i);
  await p.keyboard.press('Enter');
  await p.waitForTimeout(280);
  const afterClose = await p.evaluate(i=>document.querySelectorAll('details')[i].open, i);
  const ok = before===false && afterOpen===true && afterClose===false;
  if(ok) dOK++; else dFail.push({i,before,afterOpen,afterClose});
  console.log(`  details[${i}] closed->${before} Enter->${afterOpen} Enter->${afterClose}  ${ok?'OK':'FAIL'}`);
}
console.log(`details keyboard: ${dOK}/${sums.length} pass`, dFail.length?JSON.stringify(dFail):'');
await b.close();
