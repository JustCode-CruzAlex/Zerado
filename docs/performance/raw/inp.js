const puppeteer=require('puppeteer-core');
const CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';
(async()=>{
  const b=await puppeteer.launch({executablePath:CHROME,headless:'new',args:['--no-sandbox','--disable-gpu']});
  const p=await b.newPage();
  await p.setViewport({width:412,height:823,deviceScaleFactor:1.75,isMobile:true,hasTouch:true});
  const c=await p.createCDPSession();
  await c.send('Emulation.setCPUThrottlingRate',{rate:4});   // mid-tier mobile CPU
  await p.evaluateOnNewDocument(`
    window.__ev=[];
    new PerformanceObserver(l=>{for(const e of l.getEntries()){
      window.__ev.push({name:e.name,dur:e.duration,delay:+(e.processingStart-e.startTime).toFixed(2),
        proc:+(e.processingEnd-e.processingStart).toFixed(2)});}})
      .observe({type:'event',buffered:true,durationThreshold:0});
    window.__cls2=0;
    new PerformanceObserver(l=>{for(const e of l.getEntries()){if(!e.hadRecentInput)window.__cls2+=e.value;}})
      .observe({type:'layout-shift',buffered:true});
  `);
  await p.goto('http://localhost:4321/',{waitUntil:'networkidle0'});
  await new Promise(r=>setTimeout(r,500));
  const before=await p.evaluate(()=>window.__cls2);
  // interact with all 6 native <details> FAQ items
  const n=await p.evaluate(()=>document.querySelectorAll('details summary').length);
  for(let i=0;i<n;i++){
    const box=await p.evaluate(i=>{const s=document.querySelectorAll('details summary')[i];
      s.scrollIntoView({block:'center'}); const r=s.getBoundingClientRect();
      return {x:r.x+r.width/2,y:r.y+r.height/2};},i);
    await p.mouse.click(box.x,box.y);
    await new Promise(r=>setTimeout(r,250));
  }
  // hover/click a CTA
  const cta=await p.evaluate(()=>{const a=document.querySelector('a[href], button');a.scrollIntoView({block:'center'});const r=a.getBoundingClientRect();return{x:r.x+r.width/2,y:r.y+r.height/2};});
  await p.mouse.move(cta.x,cta.y); await p.mouse.click(cta.x,cta.y);
  await new Promise(r=>setTimeout(r,600));
  const res=await p.evaluate(()=>({ev:window.__ev,cls:window.__cls2}));
  const durs=res.ev.map(e=>e.dur).filter(d=>d>0).sort((a,b)=>b-a);
  console.log('FAQ <details> toggled: '+n+'  (CPU throttled 4x)');
  console.log('interaction events recorded: '+res.ev.length);
  console.log('worst interaction duration (INP proxy, p100): '+(durs[0]??0)+' ms');
  console.log('top 8 durations: '+durs.slice(0,8).join(', '));
  const worst=res.ev.filter(e=>e.dur===durs[0]).slice(0,3);
  console.log('worst entries: '+JSON.stringify(worst));
  console.log('CLS before interaction: '+before.toFixed(5)+'   CLS after all interactions: '+res.cls.toFixed(5));
  console.log('  (post-interaction shifts are excluded from CLS by hadRecentInput, shown for layout-stability insight only)');
  await b.close();
})();
