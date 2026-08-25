import { chromium } from 'playwright';
import { AxeBuilder } from '@axe-core/playwright';
const b=await chromium.launch();
const ctx=await b.newContext({viewport:{width:375,height:812}});
const p=await ctx.newPage();
await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
const axe=await new AxeBuilder({page:p}).withTags(['wcag2aa','wcag21aa']).analyze();
const inc=axe.incomplete.find(i=>i.id==='color-contrast');
console.log('#### axe color-contrast INCOMPLETE nodes:', inc?inc.nodes.length:0);
const reasons={};
(inc?.nodes||[]).forEach(n=>{
  const msg=(n.any?.[0]?.message)||n.failureSummary||'?';
  reasons[msg]=(reasons[msg]||0)+1;
});
Object.entries(reasons).forEach(([k,v])=>console.log(`   x${v}  ${k.slice(0,150)}`));
console.log('\nsample targets:'); (inc?.nodes||[]).slice(0,10).forEach(n=>console.log('   ',n.target, '::', n.html.slice(0,90)));

// Independent contrast computation on real rendered pixels
const res=await p.evaluate(()=>{
  function parse(c){const m=c.match(/rgba?\(([^)]+)\)/);if(!m)return null;const a=m[1].split(',').map(s=>parseFloat(s));return {r:a[0],g:a[1],b:a[2],a:a.length>3?a[3]:1};}
  function lum({r,g,b}){const f=v=>{v/=255;return v<=0.03928?v/12.92:Math.pow((v+0.055)/1.055,2.4);};return 0.2126*f(r)+0.7152*f(g)+0.0722*f(b);}
  function over(fg,bg){const a=fg.a;return {r:fg.r*a+bg.r*(1-a),g:fg.g*a+bg.g*(1-a),b:fg.b*a+bg.b*(1-a),a:1};}
  function effBg(el){
    let cur=el, stack=[];
    while(cur){const cs=getComputedStyle(cur);const bc=parse(cs.backgroundColor);
      const hasImg=cs.backgroundImage&&cs.backgroundImage!=='none';
      if(hasImg) stack.push({img:cs.backgroundImage.slice(0,40),el:cur.tagName+'.'+(cur.className||'').toString().split(' ')[0]});
      if(bc&&bc.a>0){ if(bc.a>=1) return {color:bc,imgs:stack}; stack.push({semi:bc}); }
      cur=cur.parentElement;}
    return {color:{r:255,g:255,b:255,a:1},imgs:stack};
  }
  const out=[];
  const els=[...document.querySelectorAll('p,span,h1,h2,h3,a,li,summary,div,td,strong,em,figcaption')];
  for(const el of els){
    if(!el.textContent.trim())continue;
    // only elements with a direct text node
    if(![...el.childNodes].some(n=>n.nodeType===3&&n.textContent.trim()))continue;
    const r=el.getBoundingClientRect(); if(r.width<1||r.height<1)continue;
    const cs=getComputedStyle(el);
    if(cs.visibility==='hidden'||cs.display==='none'||cs.opacity==='0')continue;
    if(el.closest('[aria-hidden="true"]'))continue;
    const fg=parse(cs.color); if(!fg)continue;
    const bgInfo=effBg(el);
    const bg=bgInfo.color;
    const fgC=fg.a<1?over(fg,bg):fg;
    const L1=lum(fgC),L2=lum(bg);
    const ratio=(Math.max(L1,L2)+0.05)/(Math.min(L1,L2)+0.05);
    const fs=parseFloat(cs.fontSize), fw=parseInt(cs.fontWeight)||400;
    const large = fs>=24 || (fs>=18.66 && fw>=700);
    const need = large?3.0:4.5;
    out.push({ text:el.textContent.trim().replace(/\s+/g,' ').slice(0,42),
      cls:(el.className||'').toString().split(' ').slice(0,2).join(' '),
      color:cs.color, bg:`rgb(${Math.round(bg.r)},${Math.round(bg.g)},${Math.round(bg.b)})`,
      fs:fs.toFixed(1), fw, large, ratio:+ratio.toFixed(2), need, pass:ratio>=need,
      overGradient:bgInfo.imgs.length>0, imgs:bgInfo.imgs.slice(0,1)});
  }
  return out;
});
console.log(`\n#### independent contrast computation: ${res.length} text elements ####`);
const fails=res.filter(r=>!r.pass);
const grad=res.filter(r=>r.overGradient);
console.log(`FAILS (< required): ${fails.length}`);
fails.forEach(f=>console.log(`   ${f.ratio}:1 (need ${f.need}) ${f.color} on ${f.bg} ${f.fs}px/${f.fw} .${f.cls} :: "${f.text}"`));
console.log(`\nover a gradient/background-image (ratio is approximate): ${grad.length}`);
[...new Set(grad.map(g=>g.cls))].slice(0,12).forEach(c=>console.log('   .'+c));
const worst=[...res].sort((a,b)=>a.ratio-b.ratio).slice(0,12);
console.log('\nLOWEST 12 ratios:');
worst.forEach(w=>console.log(`   ${w.ratio}:1 (need ${w.need}) ${w.pass?'PASS':'FAIL'} ${w.color} on ${w.bg} ${w.fs}px/${w.fw} .${w.cls} :: "${w.text}"`));
await b.close();
