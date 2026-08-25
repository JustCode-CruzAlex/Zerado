import { chromium } from 'playwright';
const b=await chromium.launch();
const ctx=await b.newContext({viewport:{width:1280,height:1000}});
const p=await ctx.newPage(); await p.goto('http://localhost:4321/',{waitUntil:'networkidle'});
// blueprint component -> candidate selectors
const map={
  SkipLink:['.z-skip'], 'CTAButton primary':['.z-cta--primary','.z-cta.z-cta--primary'],
  'CTAButton secondary':['.z-cta--secondary'], NavLink:['.z-nav-link'],
  FAQItem:['details.z-faq-item'], FooterLink:['.z-footer-link'],
  'TerminalFrame screen':['.z-frame__screen'],
  Nav:['nav','.z-nav'], Logotype:['.z-logotype','img[src="/logo.svg"]'],
  SectionEyebrow:['.z-eyebrow'], SectionHeading:['.z-section-heading','h2'],
  Hero:['#hero'], GridHorizon:['.z-grid-horizon'], ScannerRule:['.z-scanner-track'],
  TerminalFrame:['.z-frame'], TerminalChrome:['.z-frame__chrome'], TerminalRow:['.z-trow'],
  MockupCaption:['.z-mockup-caption'], MakerLine:['#maker-line','.z-maker-line'],
  ProblemStack:['#the-problem','.z-problem'], MoodCard:['.z-mood-card'],
  StateChip:['.z-chip'], StateLegend:['.z-state-legend','.z-states__legend'],
  StoreRow:['.z-store-row'], CoverTile:['.z-cover-tile'], PlatformTag:['.z-platform-tag'],
  CoverGrid:['.z-cover-grid'], PriceCard:['.z-price-card'], Sparkline:['.z-sparkline'],
  TrustPoint:['.z-trust-point'], PhoneFrame:['.z-phone-frame','.z-phone'],
  PhaseBadge:['.z-phase-badge'], CommunityPanel:['.z-community-panel'],
  RoadmapTrack:['.z-roadmap-track','.z-roadmap'], PhaseCard:['.z-phase','.z-phase-card'],
  StatusMarker:['.z-status-marker'], SpeculationBlock:['.z-speculation','#after-phase-4'],
  ClosingCTA:['#closing-cta'], Footer:['footer'], PoweredByFlowForge:['.site-footer__powered','a[href*="flowforgesoft.com"]']
};
const res=await p.evaluate(m=>{
  const out={};
  for(const [name,sels] of Object.entries(m)){
    let found=0,used=null;
    for(const s of sels){const n=document.querySelectorAll(s).length; if(n>0){found=n;used=s;break;}}
    out[name]={count:found,selector:used};
  }
  return out;
},map);
let missing=[];
console.log('#### BLUEPRINT COMPONENT INVENTORY (37 named) ####');
Object.entries(res).forEach(([k,v])=>{
  const ok=v.count>0; if(!ok)missing.push(k);
  console.log(`  ${ok?'FOUND':'MISS '}  ${String(v.count).padStart(3)}  ${k.padEnd(24)} ${v.selector||''}`);
});
console.log('\nMISSING:',missing.length?missing.join(', '):'NONE');
await b.close();
