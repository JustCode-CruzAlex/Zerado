const puppeteer=require('puppeteer-core');
const CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';
const URL='http://localhost:4321/';

const OBSERVER=`
window.__cls=0; window.__shifts=[];
new PerformanceObserver((l)=>{for(const e of l.getEntries()){
  if(e.hadRecentInput) continue;
  window.__cls+=e.value;
  window.__shifts.push({value:e.value,time:e.startTime,
    sources:(e.sources||[]).map(s=>({
      node: s.node? (s.node.nodeName+(s.node.className&&typeof s.node.className==='string'?'.'+s.node.className.trim().split(/\\s+/).join('.'):'')) : 'null',
      prev:s.previousRect, cur:s.currentRect}))});
}}).observe({type:'layout-shift',buffered:true});
window.__lcp=null;
new PerformanceObserver((l)=>{const es=l.getEntries();const e=es[es.length-1];
  window.__lcp={time:e.startTime,size:e.size,el:e.element?(e.element.nodeName+(e.element.className&&typeof e.element.className==='string'?'.'+e.element.className.trim().split(/\\s+/).join('.'):'')):'n/a',
  text:e.element?(e.element.textContent||'').slice(0,70):''};
}).observe({type:'largest-contentful-paint',buffered:true});
`;

async function run(label, opts){
  const browser=await puppeteer.launch({executablePath:CHROME,headless:'new',
    args:['--no-sandbox','--disable-gpu','--disable-background-networking','--disable-extensions','--force-device-scale-factor=1']});
  const page=await browser.newPage();
  await page.setViewport(opts.viewport);
  await page.setCacheEnabled(false);
  const client=await page.createCDPSession();
  await client.send('Network.enable');
  await client.send('Network.setCacheDisabled',{cacheDisabled:true});
  if(opts.throttle) await client.send('Network.emulateNetworkConditions',opts.throttle);
  if(opts.cpu) await client.send('Emulation.setCPUThrottlingRate',{rate:opts.cpu});
  if(opts.blockFonts){
    await page.setRequestInterception(true);
    page.on('request',r=>{ r.url().endsWith('.woff2')? r.abort() : r.continue(); });
  }
  const reqs=[];
  page.on('request',r=>reqs.push(r.url()));
  await page.evaluateOnNewDocument(OBSERVER);
  await page.goto(URL,{waitUntil:'networkidle0',timeout:120000});
  await new Promise(r=>setTimeout(r,opts.settle||3000));
  const out=await page.evaluate(()=>{
    const g=s=>{const e=document.querySelector(s);if(!e)return null;const r=e.getBoundingClientRect();
      return {w:+r.width.toFixed(1),h:+r.height.toFixed(1),top:+r.top.toFixed(1),font:getComputedStyle(e).fontFamily.split(',')[0]};};
    return {cls:window.__cls, shifts:window.__shifts, lcp:window.__lcp,
      fonts:[...document.fonts].map(f=>({fam:f.family,status:f.status,range:(f.unicodeRange||'').slice(0,24)})),
      h1:g('h1'), subhead:g('.z-hero__subhead'),
      heroFig:g('main section figure'),
      docH:document.documentElement.scrollHeight,
      resources:performance.getEntriesByType('resource').map(r=>({n:r.name,size:r.transferSize,dec:r.decodedBodySize,start:+r.startTime.toFixed(0),end:+r.responseEnd.toFixed(0)})),
      nav:(()=>{const n=performance.getEntriesByType('navigation')[0];return n?{ttfb:+n.responseStart.toFixed(1),domcl:+n.domContentLoadedEventEnd.toFixed(1),load:+n.loadEventEnd.toFixed(1)}:null;})(),
      fcp:(performance.getEntriesByName('first-contentful-paint')[0]||{}).startTime
    };
  });
  out.label=label; out.allRequests=reqs;
  await browser.close();
  return out;
}

(async()=>{
  const M={width:412,height:823,deviceScaleFactor:1.75,isMobile:true,hasTouch:true};
  const D={width:1350,height:940,deviceScaleFactor:1};
  const SLOW={offline:false,latency:400,downloadThroughput:400*1024/8,uploadThroughput:400*1024/8};
  const G4={offline:false,latency:150,downloadThroughput:1638.4*1024/8,uploadThroughput:750*1024/8};
  const scenarios=[
    ['mobile-cold-nothrottle',{viewport:M}],
    ['mobile-cold-slow4g',{viewport:M,throttle:G4,cpu:4}],
    ['mobile-cold-400kbps-cpu6',{viewport:M,throttle:SLOW,cpu:6,settle:5000}],
    ['desktop-cold-nothrottle',{viewport:D}],
    ['mobile-fonts-blocked',{viewport:M,blockFonts:true}],
  ];
  const all=[];
  for(const [l,o] of scenarios){ const r=await run(l,o); all.push(r);
    console.log('\n========== '+l);
    console.log(' CLS='+r.cls.toFixed(5)+'  shifts='+r.shifts.length+'  FCP='+(r.fcp?r.fcp.toFixed(0):'n/a')+'ms  TTFB='+(r.nav?r.nav.ttfb:'?')+'ms  load='+(r.nav?r.nav.load:'?')+'ms');
    if(r.lcp) console.log(' LCP@'+r.lcp.time.toFixed(0)+'ms size='+r.lcp.size+' el='+r.lcp.el+' :: '+JSON.stringify(r.lcp.text));
    console.log(' h1='+JSON.stringify(r.h1)+'\n subhead='+JSON.stringify(r.subhead)+'\n heroFig='+JSON.stringify(r.heroFig)+'\n docHeight='+r.docH);
    console.log(' fontRequests='+r.resources.filter(x=>x.n.includes('.woff2')).map(x=>x.n.split('/').pop()+'@'+x.end+'ms').join(', '));
    console.log(' loadedFaces='+r.fonts.filter(f=>f.status==='loaded').map(f=>f.fam).join('|')+'  unloaded='+r.fonts.filter(f=>f.status!=='loaded').map(f=>f.fam+':'+f.status).join('|'));
    const ext=r.allRequests.filter(u=>!u.startsWith('http://localhost:4321')&&!u.startsWith('data:'));
    console.log(' EXTERNAL='+ext.length+(ext.length?' '+JSON.stringify(ext):''));
    if(r.shifts.length) console.log(' SHIFTS='+JSON.stringify(r.shifts,null,1).slice(0,1800));
  }
  require('fs').writeFileSync(process.argv[2],JSON.stringify(all,null,1));
})();
