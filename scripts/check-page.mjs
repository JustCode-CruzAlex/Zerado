#!/usr/bin/env node
/**
 * Page invariants for the Zerado landing page.
 *
 * Every assertion here is a ratified decision or a load-bearing guarantee, so
 * each one names what it is defending. Run it against a built page:
 *
 *   node scripts/check-page.mjs site/dist/index.html
 *
 * CI runs it too (.github/workflows/site.yml). Exits non-zero on the first
 * failing assertion, after printing all of them.
 *
 * No dependencies on purpose — this must run anywhere Node runs.
 */
import { readFileSync, existsSync } from 'node:fs';
import { dirname, join } from 'node:path';

const file = process.argv[2] ?? 'site/dist/index.html';
if (!existsSync(file)) {
  console.error(`FAIL  ${file} does not exist — the build did not run, or ran somewhere else.`);
  process.exit(1);
}
const html = readFileSync(file, 'utf8');
const dist = dirname(file);

const results = [];
const check = (name, defends, fn) => {
  let ok = false, detail = '';
  try { [ok, detail] = fn(); } catch (e) { ok = false; detail = `threw: ${e.message}`; }
  results.push({ name, defends, ok, detail });
};

// The Q4-banned community-source name is assembled at runtime so the literal
// never appears in this repository. See docs/REDACTIONS.md.
const BANNED_SOURCE = new RegExp(['re', 'ddit'].join(''), 'i');

const ALLOWED_HOSTS = [
  'https://zerado.app',
  'https://github.com/JustCode-CruzAlex/Zerado',
  'https://www.flowforgesoft.com',
  'https://schema.org'
];

check('the page renders at all', 'the 16-section build', () => {
  const n = Buffer.byteLength(html);
  return [n >= 20000, `${n} bytes`];
});

check('_headers ships beside the page', 'the cache policy', () =>
  [existsSync(join(dist, '_headers')), join(dist, '_headers')]);

check('zero client-side JavaScript', 'the page guarantee', () => {
  const tags = html.match(/<script[^>]*>/g) ?? [];
  const offenders = tags.filter(t => !t.includes('type="application/ld+json"'));
  return [offenders.length === 0, `${tags.length} script tag(s); offenders: ${offenders.join(', ') || 'none'}`];
});

check('zero unexpected external requests', 'the page guarantee', () => {
  const urls = [...html.matchAll(/(?:src|href)="(https?:\/\/[^"]*)"/g)].map(m => m[1]);
  const offenders = [...new Set(urls.filter(u => !ALLOWED_HOSTS.some(h => u.startsWith(h))))];
  return [offenders.length === 0, `${urls.length} absolute URL(s); offenders: ${offenders.join(', ') || 'none'}`];
});

check('no funding control', 'ratified decision Q6 — disclosure, never an ask', () => {
  // "donation" in prose is disclosure and is allowed; a control is not.
  const controls = [...html.matchAll(/\b(donate|sponsor|patreon|ko-?fi)\b/gi)].map(m => m[0]);
  return [controls.length === 0, controls.join(', ') || 'none'];
});

check('the banned community-source name is absent', 'ratified decision Q4', () =>
  [!BANNED_SOURCE.test(html), BANNED_SOURCE.test(html) ? 'PRESENT' : 'absent']);

check('no placeholder text', 'the review bar', () => {
  const hits = [...html.matchAll(/lorem ipsum|\[TBD\]|\bFIXME\b|\{\{|\$\{/gi)].map(m => m[0]);
  return [hits.length === 0, hits.join(', ') || 'none'];
});

// --- the roadmap ---------------------------------------------------------
// This replaces a counter that matched the uppercase "PLANNED" inside the OG
// image alt text and never saw the roadmap's own mixed-case "Planned" at all.
// It now reads the actual status markers, in the actual roadmap section.
const DONE_WORDS = /\b(shipped|done|complete|completed|released|live|available now)\b/i;

const roadmap = html.match(/<section[^>]*id="roadmap"[\s\S]*?<\/section>/)?.[0];

check('the roadmap section exists', 'the ratified section order', () =>
  [Boolean(roadmap), roadmap ? `${roadmap.length} bytes` : 'NOT FOUND']);

const labels = roadmap
  ? [...roadmap.matchAll(/z-status-marker__label[^>]*>([^<]*)</g)].map(m => m[1].trim())
  : [];

check('the roadmap carries exactly four status markers', 'the four ratified phases', () =>
  [labels.length === 4, `found ${labels.length}: ${JSON.stringify(labels)}`]);

check('every roadmap phase reads Planned', 'no phase is marked done before it is done', () =>
  [labels.length === 4 && labels.every(l => l === 'Planned'), JSON.stringify(labels)]);

check('no status marker anywhere claims done', 'no phase is marked done before it is done', () => {
  const all = [...html.matchAll(/z-status-marker__label[^>]*>([^<]*)</g)].map(m => m[1].trim());
  const offenders = all.filter(l => DONE_WORDS.test(l));
  return [offenders.length === 0, `${all.length} marker(s) page-wide; offenders: ${JSON.stringify(offenders)}`];
});

// --- report --------------------------------------------------------------
let failed = 0;
for (const r of results) {
  if (!r.ok) failed++;
  console.log(`${r.ok ? 'ok  ' : 'FAIL'}  ${r.name}`);
  console.log(`      defends: ${r.defends}`);
  console.log(`      ${r.detail}`);
}
console.log(`\n${results.length - failed}/${results.length} page invariants hold.`);
if (failed) {
  console.log(`::error::${failed} page invariant(s) failed — each one is a ratified decision.`);
  process.exit(1);
}
