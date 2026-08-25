import { chromium, firefox, webkit } from 'playwright';
const OUT='/Users/cruzalex/Documents/projects/cruzalex/flowforge/FlowForge/forgeplay-output/landing-page/qa/screenshots';
for (const [name,eng] of Object.entries({chromium,firefox,webkit})) {
  const b=await eng.launch();
  for (const w of [375,768]) {
    const ctx=await b.newContext({viewport:{width:w,height:812}});
    const p=await ctx.newPage(); await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
    const d=await p.evaluate(()=>{
      const t=document.querySelector('.z-maker-line__text'); const r=t.getBoundingClientRect();
      const vw=document.documentElement.clientWidth;
      const rules=[...document.querySelectorAll('.z-maker-line__rule')].map(x=>Math.round(x.getBoundingClientRect().width));
      return {right:Math.round(r.right),vw,overshoot:Math.round(r.right-vw),ruleWidths:rules,
              bodyOverflowX:getComputedStyle(document.body).overflowX};
    });
    console.log(`${name.padEnd(9)} @${w}: text right=${d.right} viewport=${d.vw} OVERSHOOT=${d.overshoot}px  rule widths=${JSON.stringify(d.ruleWidths)} body overflow-x=${d.bodyOverflowX}`);
    if(name!=='chromium'||w!==375){} 
    await p.locator('#maker-line').screenshot({path:`${OUT}/DEFECT-maker-line-${name}-${w}.png`});
    await ctx.close();
  }
  await b.close();
}
