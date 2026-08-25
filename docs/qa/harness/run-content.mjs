import { chromium } from 'playwright';
const b=await chromium.launch();
const ctx=await b.newContext({viewport:{width:1280,height:900}});
const p=await ctx.newPage();
await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});

const d=await p.evaluate(()=>{
  const text=document.body.innerText;
  const html=document.documentElement.outerHTML;
  const bad=[
    // The Q4-banned community-source name is assembled at runtime so the literal never
    // appears in this repository (docs/REDACTIONS.md). The check is unchanged.
    ['$Q4_BANNED_SOURCE_NAME',new RegExp(['re','ddit'].join(''),'i')],
    ['lorem ipsum',/lorem\s+ipsum/i],['TODO',/\bTODO\b/],
    ['[TBD]',/\[TBD\]/i],['{{',/\{\{/],['${',/\$\{/],['R$ 00',/R\$\s*00/],
    ['placeholder',/\bplaceholder\b/i],['FIXME',/\bFIXME\b/],['XXX',/\bXXX\b/],
    ['Lorem',/\bLorem\b/],['undefined',/\bundefined\b/],['NaN',/\bNaN\b/]
  ];
  const hits=bad.map(([n,re])=>{
    const inText=re.test(text); const inHtml=re.test(html);
    let sample=null;
    if(inHtml){const m=html.match(new RegExp('.{0,70}'+re.source+'.{0,70}', re.flags.replace('g','')));sample=m?m[0]:null;}
    return {token:n,inVisibleText:inText,inHtml:inHtml,sample};
  }).filter(h=>h.inHtml||h.inVisibleText);

  // lone X placeholders in visible text
  const loneX=[...text.matchAll(/(^|[\s(])X{1,3}([\s).,]|$)/g)].map(m=>m[0]).slice(0,10);

  // roadmap phases
  const roadmap=document.querySelector('#roadmap');
  const phaseCards=[...roadmap.querySelectorAll('.z-phase-card, [class*="phase-card"]')];
  const statusMarkers=[...roadmap.querySelectorAll('[class*="status"], .z-status-marker')].map(e=>e.textContent.trim().replace(/\s+/g,' '));
  const roadmapText=roadmap.innerText;
  const plannedCount=(roadmapText.match(/planned/gi)||[]).length;

  // community section interactivity (§11)
  const community=document.querySelector('#community');
  const commInteractive=[...community.querySelectorAll('a[href],button,input,select,textarea,summary,[tabindex]:not([tabindex="-1"]),[role="button"],[onclick]')]
    .map(e=>({tag:e.tagName.toLowerCase(),href:e.getAttribute('href'),text:e.textContent.trim().slice(0,40)}));

  // donate/sponsor anywhere
  const moneyWords=['donate','sponsor','patreon','ko-fi','kofi','buymeacoffee','github sponsors','paypal','pix','contribute financially'];
  const moneyHits=moneyWords.map(w=>({w,inText:new RegExp(w,'i').test(text),inHtml:new RegExp(w,'i').test(html)})).filter(x=>x.inText||x.inHtml);
  const moneyLinks=[...document.querySelectorAll('a[href]')].filter(a=>/donate|sponsor|patreon|ko-?fi|paypal|buymeacoffee/i.test(a.href+a.textContent))
    .map(a=>({href:a.href,text:a.textContent.trim()}));

  // figures + captions
  const figures=[...document.querySelectorAll('figure')].map(f=>{
    const cap=f.querySelector('figcaption');
    return {section:f.closest('section,footer')?.id||'?', role:f.getAttribute('role'),
      ariaLabelledby:f.getAttribute('aria-labelledby'),
      capId:cap?.id||null,
      caption:cap?cap.textContent.trim().replace(/\s+/g,' '):null,
      captionVisible: cap? (()=>{const cs=getComputedStyle(cap);const r=cap.getBoundingClientRect();
        return cs.display!=='none'&&cs.visibility!=='hidden'&&cs.opacity!=='0'&&r.width>0&&r.height>0;})() : false,
      hasTerminalFrame: !!f.querySelector('[class*="frame"],[class*="terminal"]')};
  });

  // cover grid caption
  const coverCaption=[...document.querySelectorAll('figcaption,p,span,div')]
    .map(e=>e.textContent.trim().replace(/\s+/g,' '))
    .filter(t=>/Cover tiles are illustrative artwork/i.test(t))[0]||null;

  // powered by flowforge
  const pbf=[...document.querySelectorAll('a[href]')].filter(a=>/flowforgesoft\.com/i.test(a.href));
  const pbfInfo=pbf.map(a=>({href:a.getAttribute('href'),text:a.textContent.trim().replace(/\s+/g,' '),
    rel:a.getAttribute('rel'),
    imgs:[...a.querySelectorAll('img')].map(i=>({src:i.getAttribute('src'),alt:i.getAttribute('alt'),
      w:i.naturalWidth,h:i.naturalHeight,complete:i.complete}))}));

  // all links
  const links=[...document.querySelectorAll('a[href]')].map(a=>({
    href:a.getAttribute('href'), resolved:a.href, text:a.textContent.trim().replace(/\s+/g,' ').slice(0,40),
    rel:a.getAttribute('rel'), target:a.getAttribute('target'),
    section:a.closest('section,footer,nav,header')?.id||a.closest('nav')?.tagName||'?'}));

  // images
  const imgs=[...document.querySelectorAll('img')].map(i=>({src:i.getAttribute('src'),alt:i.getAttribute('alt'),
    hasAlt:i.hasAttribute('alt'),natural:i.naturalWidth+'x'+i.naturalHeight,complete:i.complete,
    section:i.closest('section,footer,nav,header')?.id||'?'}));

  // headings
  const heads=[...document.querySelectorAll('h1,h2,h3,h4,h5,h6')].map(h=>({
    lvl:+h.tagName[1], text:h.textContent.trim().replace(/\s+/g,' ').slice(0,60),
    section:h.closest('section,footer')?.id||'?'}));

  // landmarks
  const landmarks={
    header:document.querySelectorAll('header, [role="banner"]').length,
    main:document.querySelectorAll('main, [role="main"]').length,
    footer:document.querySelectorAll('footer, [role="contentinfo"]').length,
    nav:document.querySelectorAll('nav, [role="navigation"]').length
  };

  // sections in DOM order
  const sections=[...document.querySelectorAll('main > section, body > footer, main > *[id]')].map(s=>({tag:s.tagName.toLowerCase(),id:s.id}));
  const allSections=[...document.querySelectorAll('section[id], footer[id]')].map(s=>s.id);

  return {hits,loneX,statusMarkers,plannedCount,roadmapText:roadmapText.slice(0,600),
    phaseCardCount:phaseCards.length,commInteractive,moneyHits,moneyLinks,figures,coverCaption,pbfInfo,links,imgs,heads,landmarks,sections,allSections,
    lang:document.documentElement.lang, title:document.title,
    metaDesc:document.querySelector('meta[name="description"]')?.content?.slice(0,120)};
});

console.log('#### 7 CONTENT HONESTY ####');
console.log('forbidden-token hits:', d.hits.length?JSON.stringify(d.hits,null,1):'NONE');
console.log('lone X placeholders in visible text:', d.loneX.length?JSON.stringify(d.loneX):'NONE');
console.log('\nroadmap status markers:', JSON.stringify(d.statusMarkers));
console.log('roadmap "planned" occurrences:', d.plannedCount, '| phase cards:', d.phaseCardCount);
console.log('roadmap text:', JSON.stringify(d.roadmapText.slice(0,400)));
console.log('\n§11 community interactive elements:', d.commInteractive.length, JSON.stringify(d.commInteractive));
console.log('money words found:', JSON.stringify(d.moneyHits));
console.log('donate/sponsor LINKS:', JSON.stringify(d.moneyLinks));
console.log('\nfigures:', d.figures.length);
d.figures.forEach((f,i)=>console.log(`  fig[${i}] §${f.section} role=${f.role} labelledby=${f.ariaLabelledby} capId=${f.capId} visible=${f.captionVisible}\n     caption="${f.caption}"`));
console.log('\ncover-grid caption:', JSON.stringify(d.coverCaption));
console.log('powered-by-flowforge:', JSON.stringify(d.pbfInfo,null,1));

console.log('\n#### 9 SEMANTICS ####');
console.log('lang=',d.lang,'| title=',JSON.stringify(d.title));
console.log('meta description:',JSON.stringify(d.metaDesc));
console.log('landmarks:',JSON.stringify(d.landmarks));
console.log('headings:',d.heads.length);
d.heads.forEach(h=>console.log(`   h${h.lvl} [§${h.section}] ${h.text}`));
let prev=0,skips=[];
d.heads.forEach(h=>{if(prev&&h.lvl>prev+1)skips.push(`h${prev}->h${h.lvl} at "${h.text}"`);prev=h.lvl;});
console.log('h1 count:',d.heads.filter(h=>h.lvl===1).length,'| skipped levels:',skips.length?JSON.stringify(skips):'NONE');
console.log('images:',d.imgs.length,JSON.stringify(d.imgs,null,1));
console.log('\n#### 10 SECTIONS ####');
console.log('section ids in DOM order:',JSON.stringify(d.allSections));
console.log('\n#### 8 LINKS ####');
d.links.forEach((l,i)=>console.log(`  [${i}] §${l.section} rel=${l.rel} target=${l.target} "${l.text}" -> ${l.href}`));
await b.close();
