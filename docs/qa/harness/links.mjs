import { chromium } from 'playwright';
const b=await chromium.launch();
const ctx=await b.newContext({viewport:{width:1280,height:900}});
const p=await ctx.newPage();
const failedReqs=[];
p.on('requestfailed',r=>failedReqs.push({url:r.url(),err:r.failure()?.errorText}));
const responses=[];
p.on('response',r=>responses.push({url:r.url(),status:r.status()}));
await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
// scroll to bottom to trigger any lazy loading
await p.evaluate(async()=>{ for(let y=0;y<document.body.scrollHeight;y+=600){window.scrollTo(0,y);await new Promise(r=>setTimeout(r,40));} window.scrollTo(0,document.body.scrollHeight);});
await p.waitForTimeout(1200);
const imgs=await p.evaluate(()=>[...document.querySelectorAll('img')].map(i=>({
  src:i.getAttribute('src'),alt:i.getAttribute('alt'),loading:i.getAttribute('loading'),
  natural:i.naturalWidth+'x'+i.naturalHeight,complete:i.complete,
  displayed:(()=>{const r=i.getBoundingClientRect();return Math.round(r.width)+'x'+Math.round(r.height);})()})));
console.log('#### IMAGES after full scroll ####');
console.log(JSON.stringify(imgs,null,1));
console.log('\n#### network responses (same-origin assets) ####');
[...new Map(responses.map(r=>[r.url,r])).values()].forEach(r=>console.log(`  ${r.status}  ${r.url}`));
console.log('failed requests:',failedReqs.length?JSON.stringify(failedReqs):'NONE');

// in-page anchors
const anchors=await p.evaluate(()=>{
  const ids=new Set([...document.querySelectorAll('[id]')].map(e=>e.id));
  return [...document.querySelectorAll('a[href^="#"]')].map(a=>({href:a.getAttribute('href'),
    target:a.getAttribute('href').slice(1), exists:ids.has(a.getAttribute('href').slice(1))}));
});
console.log('\n#### in-page anchors ####',JSON.stringify(anchors));
await b.close();
