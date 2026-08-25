import { chromium } from 'playwright';
const OUT='/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa/screenshots';
const b = await chromium.launch();
for (const mode of ['no-preference','reduce']) {
  const ctx = await b.newContext({viewport:{width:1280,height:900}, reducedMotion: mode});
  const p = await ctx.newPage();
  await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
  await p.waitForTimeout(800);
  const d = await p.evaluate(()=>{
    const pips=[...document.querySelectorAll('.z-scanner-track__pip')];
    const tracks=[...document.querySelectorAll('.z-scanner-track')];
    return {
      pipCount:pips.length, trackCount:tracks.length,
      pips: pips.map(el=>{
        const cs=getComputedStyle(el); const r=el.getBoundingClientRect();
        const tr=el.parentElement.getBoundingClientRect();
        return { anim:cs.animationName, dur:cs.animationDuration, opacity:cs.opacity,
          visibility:cs.visibility, display:cs.display, left:cs.left, width:cs.width,
          background:cs.background.slice(0,60), boxShadow:cs.boxShadow.slice(0,50),
          rectX:Math.round(r.x), rectW:Math.round(r.width), rectH:Math.round(r.height),
          trackX:Math.round(tr.x), trackW:Math.round(tr.width),
          centreOffsetPx: Math.round((r.x + r.width/2) - (tr.x + tr.width/2)),
          pipPctOfTrack: +( (r.width/tr.width)*100 ).toFixed(1),
          parentClass:(el.parentElement.className||'').toString() };
      })
    };
  });
  console.log(`\n===== reducedMotion=${mode} =====`);
  console.log(`pips=${d.pipCount} tracks=${d.trackCount}`);
  d.pips.forEach((x,i)=>console.log(`  pip[${i}] parent="${x.parentClass}"\n     anim=${x.anim} dur=${x.dur} opacity=${x.opacity} vis=${x.visibility} display=${x.display}\n     left=${x.left} width=${x.width} rectX=${x.rectX} rectW=${x.rectW} rectH=${x.rectH}\n     track x=${x.trackX} w=${x.trackW} | offset-from-track-centre=${x.centreOffsetPx}px | pip=${x.pipPctOfTrack}% of track\n     bg=${x.background} shadow=${x.boxShadow}`));
  const sc = await p.locator('.z-scanner-track').first();
  await sc.screenshot({path:`${OUT}/scanner-${mode}.png`});
  await ctx.close();
}
await b.close();
