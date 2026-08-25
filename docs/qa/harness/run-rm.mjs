import { chromium } from 'playwright';
const OUT='/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa/screenshots';
const b = await chromium.launch();

for (const mode of ['no-preference','reduce']) {
  const ctx = await b.newContext({ viewport:{width:1280,height:900}, reducedMotion: mode==='reduce'?'reduce':'no-preference' });
  const p = await ctx.newPage();
  await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
  await p.waitForTimeout(600);
  const d = await p.evaluate(()=>{
    const out = {};
    out.mqReduce = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    // every element with an animation
    const anim = [];
    document.querySelectorAll('*').forEach(el=>{
      for (const pseudo of ['', '::before', '::after']) {
        const cs = getComputedStyle(el, pseudo||null);
        if (cs.animationName && cs.animationName !== 'none') {
          anim.push({ sel:(el.tagName.toLowerCase()+'.'+((el.className||'').toString().split(' ')[0]||''))+pseudo,
            name:cs.animationName, dur:cs.animationDuration, iter:cs.animationIterationCount,
            dir:cs.animationDirection, opacity:cs.opacity, left:cs.left, timing:cs.animationTimingFunction });
        }
      }
    });
    out.animations = anim;
    // transitions on interactive elements
    const trans = [];
    document.querySelectorAll('a, summary, .z-btn, .z-cta').forEach(el=>{
      const cs=getComputedStyle(el);
      if (cs.transitionDuration && cs.transitionDuration!=='0s')
        trans.push({ tag:el.tagName.toLowerCase(), cls:(el.className||'').toString().split(' ')[0],
                     prop:cs.transitionProperty, dur:cs.transitionDuration });
    });
    out.transitions = trans;
    // scanner pip specifically
    const sc = document.querySelector('.z-scanner');
    if (sc) {
      const cs = getComputedStyle(sc,'::after');
      const scr = sc.getBoundingClientRect();
      out.scannerPip = { animationName:cs.animationName, animationDuration:cs.animationDuration,
        opacity:cs.opacity, left:cs.left, width:cs.width, trackWidth:Math.round(scr.width),
        iter:cs.animationIterationCount };
      out.scannerCount = document.querySelectorAll('.z-scanner').length;
    }
    out.scrollBehavior = getComputedStyle(document.documentElement).scrollBehavior;
    return out;
  });
  console.log(`\n########## prefers-reduced-motion: ${mode} ##########`);
  console.log('media query matches reduce:', d.mqReduce, '| html scroll-behavior:', d.scrollBehavior);
  console.log('scanner instances (.z-scanner):', d.scannerCount);
  console.log('SCANNER PIP (::after):', JSON.stringify(d.scannerPip,null,1));
  console.log(`animated elements: ${d.animations.length}`);
  d.animations.forEach(a=>console.log(`   ${a.sel} anim=${a.name} dur=${a.dur} iter=${a.iter} dir=${a.dir} left=${a.left} opacity=${a.opacity}`));
  console.log(`transition-bearing interactive elements: ${d.transitions.length}`);
  const uniq={}; d.transitions.forEach(t=>{const k=t.cls+'|'+t.prop+'|'+t.dur; uniq[k]=(uniq[k]||0)+1;});
  Object.entries(uniq).forEach(([k,v])=>console.log(`   x${v}  ${k}`));
  await p.screenshot({path:`${OUT}/reduced-motion-${mode}-1280.png`, fullPage:false});
  await ctx.close();
}
await b.close();
