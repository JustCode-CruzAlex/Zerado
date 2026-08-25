import { chromium, firefox, webkit } from 'playwright';
import fs from 'node:fs';
const OUT='/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa';
const engines={chromium,firefox,webkit};
const report={};
for (const [name,eng] of Object.entries(engines)) {
  let b;
  try { b = await eng.launch(); } catch(e){ console.log(`${name}: LAUNCH FAILED ${e.message}`); report[name]={launchError:e.message}; continue; }
  const v = b.version();
  report[name]={version:v,viewports:{}};
  console.log(`\n################ ${name} v${v} ################`);
  for (const w of [375,768,1280,1920]) {
    const ctx=await b.newContext({viewport:{width:w,height:w===375?812:1000},deviceScaleFactor:1});
    const p=await ctx.newPage();
    const errs=[]; p.on('pageerror',e=>errs.push(e.message));
    const bad=[]; p.on('response',r=>{if(r.status()>=400)bad.push(r.status()+' '+r.url());});
    await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
    await p.waitForTimeout(500);
    await p.screenshot({path:`${OUT}/screenshots/${name}-${w}-fullpage.png`,fullPage:true});
    const d=await p.evaluate(()=>{
      const de=document.documentElement;
      const supp=(prop,val)=>{try{return CSS.supports(prop,val);}catch(_){return null;}};
      const h1=document.querySelector('h1'); const h1cs=h1?getComputedStyle(h1):null;
      const det=document.querySelector('details'); const sum=det?.querySelector('summary');
      const sumcs=sum?getComputedStyle(sum):null;
      const tile=document.querySelector('[class*="cover-tile"]');
      const tilecs=tile?getComputedStyle(tile):null;
      const bdf=[...document.querySelectorAll('*')].filter(e=>{const c=getComputedStyle(e);
        return (c.backdropFilter&&c.backdropFilter!=='none')||(c.webkitBackdropFilter&&c.webkitBackdropFilter!=='none');}).length;
      const grads=[...document.querySelectorAll('*')].filter(e=>/gradient/.test(getComputedStyle(e).backgroundImage)).length;
      const scripts=document.querySelectorAll('script').length;
      // section geometry fingerprint
      const secs=[...document.querySelectorAll('section[id],footer[id]')].map(s=>{
        const r=s.getBoundingClientRect(); return {id:s.id,h:Math.round(r.height),w:Math.round(r.width)};});
      return {
        scrollWidth:de.scrollWidth, clientWidth:de.clientWidth,
        overflowPass: de.scrollWidth===de.clientWidth,
        docHeight: Math.round(document.body.scrollHeight),
        supports:{clamp:supp('width','clamp(1px,2vw,3px)'),aspectRatio:supp('aspect-ratio','1/1'),
          backdropFilter:supp('backdrop-filter','blur(2px)'),
          interpolateSize:supp('interpolate-size','allow-keywords'),
          detailsContent:supp('selector(::details-content)')||null,
          hasSel:(()=>{try{document.querySelector(':has(a)');return true;}catch(_){return false;}})()},
        h1:{fontFamily:h1cs?.fontFamily?.slice(0,40),fontSize:h1cs?.fontSize,lineHeight:h1cs?.lineHeight,
            fontWeight:h1cs?.fontWeight, w:Math.round(h1.getBoundingClientRect().width), h:Math.round(h1.getBoundingClientRect().height)},
        details:{listStyle:sumcs?.listStyle?.slice(0,30),display:sumcs?.display,
                 markerHidden:(()=>{try{return getComputedStyle(sum,'::-webkit-details-marker').display;}catch(_){return 'n/a';}})(),
                 h:Math.round(sum.getBoundingClientRect().height)},
        coverTile:tilecs?{aspectRatio:tilecs.aspectRatio,w:Math.round(tile.getBoundingClientRect().width),h:Math.round(tile.getBoundingClientRect().height)}:null,
        backdropFilterUsers:bdf, gradientUsers:grads, scriptTags:scripts,
        fontsLoaded:[...document.fonts].length,
        secs};
    });
    report[name].viewports[w]={...d,pageErrors:errs,badResponses:bad};
    console.log(`  [${w}] overflowPass=${d.overflowPass} scrollW=${d.scrollWidth}/${d.clientWidth} docH=${d.docHeight} h1=${d.h1.fontSize}/${d.h1.fontFamily} h1box=${d.h1.w}x${d.h1.h} summaryH=${d.details.h} tile=${d.coverTile?d.coverTile.w+'x'+d.coverTile.h+' ar='+d.coverTile.aspectRatio:'-'} grads=${d.gradientUsers} bdf=${d.backdropFilterUsers} scripts=${d.scriptTags} fonts=${d.fontsLoaded}`);
    if(errs.length)console.log(`     PAGE ERRORS: ${JSON.stringify(errs)}`);
    if(bad.length)console.log(`     BAD RESPONSES: ${JSON.stringify(bad)}`);
    if(w===1280)console.log(`     supports: ${JSON.stringify(d.supports)}`);
    await ctx.close();
  }
  await b.close();
}
fs.writeFileSync(`${OUT}/qa-crossbrowser.json`,JSON.stringify(report,null,2));
console.log('\nWROTE qa-crossbrowser.json');
// section-height divergence
const eng=Object.keys(report).filter(k=>!report[k].launchError);
for (const w of [375,768,1280,1920]) {
  const base=report[eng[0]].viewports[w].secs;
  console.log(`\n--- section height divergence @${w} (baseline ${eng[0]}) ---`);
  base.forEach((s,i)=>{
    const row=eng.map(e=>report[e].viewports[w].secs[i]?.h);
    const max=Math.max(...row),min=Math.min(...row);
    const delta=max-min;
    if(delta>4) console.log(`  ${s.id.padEnd(20)} ${eng.map((e,j)=>e+'='+row[j]).join('  ')}   Δ=${delta}px`);
  });
}
