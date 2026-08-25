const puppeteer=require('puppeteer-core');
const CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';
(async()=>{
for(const [label,vp] of [['MOBILE 412x823',{width:412,height:823,deviceScaleFactor:1.75,isMobile:true}],
                          ['DESKTOP 1350x940',{width:1350,height:940,deviceScaleFactor:1}]]){
  const b=await puppeteer.launch({executablePath:CHROME,headless:'new',args:['--no-sandbox','--disable-gpu']});
  const p=await b.newPage(); await p.setViewport(vp);
  await p.coverage.startCSSCoverage();
  await p.goto('http://localhost:4321/',{waitUntil:'networkidle0'});
  const cov=await p.coverage.stopCSSCoverage();
  const r=await p.evaluate((vh)=>{
    const fams={};
    const hits=[];
    document.querySelectorAll('*').forEach(e=>{
      const b=e.getBoundingClientRect();
      if(b.top < vh && b.bottom > 0 && b.width>0 && b.height>0){
        const f=getComputedStyle(e).fontFamily.split(',')[0].replace(/"/g,'');
        const txt=(e.textContent||'').trim();
        const own=[...e.childNodes].some(n=>n.nodeType===3&&n.textContent.trim());
        if(own){ fams[f]=(fams[f]||0)+1; hits.push({f,tag:e.tagName,cls:(e.className||'').toString().slice(0,40),top:Math.round(b.top),t:txt.slice(0,40)}); }
      }
    });
    const h1=document.querySelector('h1'); const cs=h1?getComputedStyle(h1):null;
    return {fams,hits:hits.slice(0,40),
      h1style: cs? {fam:cs.fontFamily.split(',')[0], fill:cs.webkitTextFillColor, clip:cs.webkitBackgroundClip||cs.backgroundClip, bgImg:cs.backgroundImage.slice(0,60), color:cs.color, opacity:cs.opacity, size:cs.fontSize, lh:cs.lineHeight}:null};
  },vp.height);
  console.log('\n===== '+label);
  console.log(' above-fold text elements by font: '+JSON.stringify(r.fams));
  console.log(' h1 computed: '+JSON.stringify(r.h1style));
  const mono=r.hits.filter(h=>h.f.includes('JetBrains'));
  console.log(' above-fold MONO users: '+(mono.length?JSON.stringify(mono,null,1).slice(0,800):'NONE'));
  let tot=0,used=0;
  for(const c of cov){tot+=c.text.length; used+=c.ranges.reduce((a,x)=>a+(x.end-x.start),0);
    console.log(' CSS '+c.url.split('/').pop()+': total='+c.text.length+' usedAtLoad='+c.ranges.reduce((a,x)=>a+(x.end-x.start),0));}
  console.log(' CSS coverage overall: used='+used+'/'+tot+' = '+(100*used/tot).toFixed(1)+'% ('+(tot-used)+' B unused at initial load)');
  await b.close();
}})();
