import { chromium } from 'playwright';
const b=await chromium.launch();
const ctx=await b.newContext({viewport:{width:1280,height:1000}});
const p=await ctx.newPage(); await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
const d=await p.evaluate(()=>{
  const chips=[...document.querySelectorAll('.z-chip')].map((c,i)=>{
    const g=c.querySelector('[class*="glyph"]'); const cs=getComputedStyle(c);
    return {i,glyph:g?.textContent.trim(),glyphSize:g?getComputedStyle(g).fontSize:null,
      boxShadow:cs.boxShadow.slice(0,50),w:Math.round(c.getBoundingClientRect().width),
      label:(c.querySelector('[class*="label"]')?.textContent||'').trim()};});
  // italicised zerado
  const ital=[...document.querySelectorAll('em,i,[style*="italic"]')].map(e=>({tag:e.tagName,text:e.textContent.trim().slice(0,40),
    style:getComputedStyle(e).fontStyle}));
  // glow count page-wide
  const glows=[...document.querySelectorAll('*')].filter(e=>{const s=getComputedStyle(e).boxShadow;
    return s&&s!=='none';}).map(e=>({cls:(e.className||'').toString().split(' ')[0],shadow:getComputedStyle(e).boxShadow.slice(0,45)}));
  // fonts
  const fonts=[...document.fonts].map(f=>({family:f.family,weight:f.weight,status:f.status}));
  // sinopse / zerado explained once
  const txt=document.body.innerText;
  return {chips,ital,glowCount:glows.length,glows:glows.slice(0,12),fonts,
    zeradoMentions:(txt.match(/zerado/gi)||[]).length,
    sinopseMentions:(txt.match(/sinopse/gi)||[]).length,
    hasPortugueseGloss:/Portuguese for/i.test(txt),
    affiliateDisclosure:/commission|affiliate/i.test(txt),
    premiumLine:/premium account or a donation|premium|donation/i.test(txt)};
});
console.log('STATE CHIPS:');d.chips.forEach(c=>console.log(`  [${c.i}] glyph=${c.glyph} glyphSize=${c.glyphSize} width=${c.w} label="${c.label}" shadow=${c.boxShadow}`));
console.log('\nitalic elements:',JSON.stringify(d.ital));
console.log('page-wide box-shadow users:',d.glowCount);
d.glows.forEach(g=>console.log('   .'+g.cls,g.shadow));
console.log('\nfonts loaded:',JSON.stringify(d.fonts));
console.log('zerado mentions:',d.zeradoMentions,'| sinopse:',d.sinopseMentions,'| PT gloss present:',d.hasPortugueseGloss);
console.log('affiliate disclosure:',d.affiliateDisclosure,'| premium/donation line:',d.premiumLine);
await b.close();
