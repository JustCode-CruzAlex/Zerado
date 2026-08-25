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

// --- the roadmap -----------------------------------------------------------
// Pins the EXACT per-phase state, not "they all say the same thing". Phase 1
// moved to In progress on 2026-08-25; 2-4 stay Planned; none may ever read a
// done-equivalent. Anything else is a regression, in either direction.
const DONE_WORDS = /\b(shipped|done|complete|completed|released|live|available now)\b/i;
const EXPECTED = ['In progress', 'Planned', 'Planned', 'Planned'];

const roadmap = html.match(/<section[^>]*id="roadmap"[\s\S]*?<\/section>/)?.[0];

check('the roadmap section exists', 'the ratified section order', () =>
  [Boolean(roadmap), roadmap ? `${roadmap.length} bytes` : 'NOT FOUND']);

// Each marker's colour class, glyph and label, in document order.
const markers = roadmap
  ? [...roadmap.matchAll(
      /z-status-marker z-status-marker--([a-z-]+)"[^>]*>\s*<span class="z-status-marker__glyph"[^>]*>([^<]*)<\/span>\s*<span class="z-status-marker__label"[^>]*>([^<]*)</g
    )].map(m => ({ variant: m[1], glyph: m[2].trim(), label: m[3].trim() }))
  : [];

check('the roadmap carries exactly four status markers', 'the four ratified phases', () =>
  [markers.length === 4, `found ${markers.length}: ${JSON.stringify(markers.map(m => m.label))}`]);

check('each phase carries its expected status', 'Phase 1 in progress; 2-4 planned', () => {
  const got = markers.map(m => m.label);
  return [JSON.stringify(got) === JSON.stringify(EXPECTED),
    `expected ${JSON.stringify(EXPECTED)}, got ${JSON.stringify(got)}`];
});

check('no phase is marked done', 'the one line that would cost credibility', () => {
  const all = [...html.matchAll(/z-status-marker__label[^>]*>([^<]*)</g)].map(m => m[1].trim());
  const offenders = all.filter(l => DONE_WORDS.test(l));
  return [offenders.length === 0, `${all.length} marker(s) page-wide; offenders: ${JSON.stringify(offenders)}`];
});

// Colour alone must never carry the meaning: brand-manual.md §3.
check('every status marker co-renders colour, glyph and label', 'the co-render rule', () => {
  const bad = markers.filter(m => !m.variant || !m.glyph || !m.label);
  const distinctGlyphs = new Set(markers.map(m => m.glyph)).size;
  const distinctVariants = new Set(markers.map(m => m.variant)).size;
  return [bad.length === 0 && distinctGlyphs === distinctVariants,
    `${markers.length} marker(s); ${distinctVariants} colour variant(s), ${distinctGlyphs} distinct glyph(s) — ` +
    `${JSON.stringify(markers.map(m => `${m.variant}/${m.glyph}/${m.label}`))}`];
});

// The roadmap stays undated until that separate question is answered.
check('the roadmap renders no dates', 'undated, ordered phases', () => {
  const hits = roadmap
    ? [...roadmap.matchAll(/\b(20\d\d|Q[1-4]|January|February|March|April|May|June|July|August|September|October|November|December)\b/g)].map(m => m[0])
    : [];
  return [hits.length === 0, hits.join(', ') || 'none'];
});

// --- capability honesty -----------------------------------------------------
// The gap that let a false present-tense claim survive two review rounds inside
// a green suite: nothing asserted on what the page says *works*. Zerado has no
// Go code, so nothing on this page may claim, in the present tense, that a
// capability exists. Each entry flips only when the thing actually ships.
const PRESENT_TENSE_CLAIMS = [
  /\bis built and works\b/i,
  /\bsyncs today\b/i,
  /\bworks today\b/i,
  /\bavailable (now|today)\b/i,
  /\bdownload (it |the )?(now|today)\b/i,
  /\byou can install\b/i
];

check('no present-tense capability claim', 'nothing dishonest ships', () => {
  const visible = html.replace(/<(script|style)[^>]*>[\s\S]*?<\/\1>/g, '').replace(/<[^>]+>/g, ' ');
  const hits = PRESENT_TENSE_CLAIMS.map(re => visible.match(re)?.[0]).filter(Boolean);
  return [hits.length === 0, hits.join(', ') || 'none'];
});

check('no store row claims to be live', 'no store syncs until Phase 1 ships', () => {
  const live = (html.match(/z-store-row--live/g) ?? []).length;
  return [live === 0, `${live} live row(s)`];
});

check('every store row exposes its status to assistive technology', 'the co-render rule', () => {
  const rows = [...html.matchAll(/<div class="z-store-row z-store-row--([a-z]+)"[^>]*>([\s\S]*?)<\/div>/g)];
  // The glyph is aria-hidden, so colour would otherwise be the only signal.
  const bare = rows.filter(r => !/z-visually-hidden/.test(r[2]));
  return [rows.length > 0 && bare.length === 0,
    `${rows.length} row(s), ${bare.length} with no accessible status word`];
});

// --- the money position ------------------------------------------------------
// Amended 2026-08-25 (content/landing-copy.md, Amendments 3 and 4): the affiliate
// model is dropped and there is no premium account. The page states that the
// Phase 4 layer will be donation-supported and asks for nothing — ratified
// decision Q6's "disclosure is not an ask" still governs. These pin the positive
// claims specifically, so the honest negations ("there is no commission…") pass.
const REVENUE_CLAIMS = [
  /\baffiliate\b/i,
  /earns? a commission/i,
  /premium account or/i,
  /\bcommission when\b/i
];

check('the page claims no revenue', 'the affiliate model is dropped', () => {
  const visible = html.replace(/<(script|style)[^>]*>[\s\S]*?<\/\1>/g, '').replace(/<[^>]+>/g, ' ');
  const hits = REVENUE_CLAIMS.map(re => visible.match(re)?.[0]).filter(Boolean);
  return [hits.length === 0, hits.join(', ') || 'none'];
});

check('the community layer is stated as donation-supported', 'the amended Q3 position', () => {
  const community = html.match(/<section[^>]*id="community"[\s\S]*?<\/section>/)?.[0] ?? '';
  const text = community.replace(/<[^>]+>/g, ' ');
  return [/donation/i.test(text) && !/premium account or/i.test(text),
    /donation/i.test(text) ? 'states donation-supported' : 'NO donation statement'];
});


// --- the cover grid ----------------------------------------------------------
// Ticket #16 replaced twelve art-directed CSS tiles with real IGDB cover art.
// Both states are legitimate — a cover withdrawn for any reason falls back to
// the tile it replaced (docs/legal/igdb-image-licence-finding.md §7) — so every
// assertion below holds in BOTH, and the pair of them pins the two states to
// each other so the page can never show one and say the other.
const COVER_IMGS = [...html.matchAll(/<img\b[^>]*\bsrc="\/covers\/[^"]+"[^>]*>/g)].map(m => m[0]);
const attr = (tag, name) => tag.match(new RegExp(`\\b${name}="([^"]*)"`))?.[1];

check('every cover declares its intrinsic width and height', 'CLS stays 0.000', () => {
  const bad = COVER_IMGS.filter(t => !/\bwidth="\d+"/.test(t) || !/\bheight="\d+"/.test(t));
  return [bad.length === 0, `${COVER_IMGS.length} cover(s); ${bad.length} missing dimensions`];
});

check('every cover is lazy-loaded', 'the grid is section 6 of 16 — never above the fold', () => {
  const bad = COVER_IMGS.filter(t => attr(t, 'loading') !== 'lazy');
  return [bad.length === 0, `${bad.length} cover(s) not lazy`];
});

check('every cover is served from this origin', 'zero external requests', () => {
  const bad = COVER_IMGS.filter(t => !(attr(t, 'src') ?? '').startsWith('/covers/'));
  const srcsets = [...html.matchAll(/<source\b[^>]*\bsrcset="([^"]*)"/g)].map(m => m[1]);
  const badSet = srcsets.filter(u => !u.startsWith('/covers/'));
  return [bad.length === 0 && badSet.length === 0,
    `${bad.length} off-origin img, ${badSet.length} off-origin source`];
});

check('cover alt text names each game and is distinct', 'twelve rows must not reach a reader as one string', () => {
  const alts = COVER_IMGS.map(t => attr(t, 'alt') ?? '');
  const empty = alts.filter(a => a.trim() === '').length;
  const distinct = new Set(alts).size;
  return [empty === 0 && distinct === alts.length,
    `${alts.length} alt(s), ${distinct} distinct, ${empty} empty`];
});

check('no list element carries role="img"', 'ARIA-in-HTML — the aria-allowed-role fix of REDACTIONS §3.5', () => {
  const bad = [...html.matchAll(/<(ul|ol|li)\b[^>]*\brole="img"/g)].map(m => m[0]);
  return [bad.length === 0, bad.join(', ') || 'none'];
});

check('the cover disclosure matches what actually shipped', 'the page never credits a source it did not use', () => {
  const shipped = COVER_IMGS.length > 0;
  const saysIllustrative = /Cover tiles are illustrative artwork/.test(html);
  const creditsIgdb = /Cover art from IGDB\.com/.test(html);
  const ok = shipped ? (creditsIgdb && !saysIllustrative) : (saysIllustrative && !creditsIgdb);
  return [ok, shipped
    ? `${COVER_IMGS.length} real cover(s): credits IGDB ${creditsIgdb}, still says illustrative ${saysIllustrative}`
    : `no real covers: says illustrative ${saysIllustrative}, credits IGDB ${creditsIgdb}`];
});

check('real cover art carries its attribution in the footer', 'IGDB fair attribution — visible, static location', () => {
  const shipped = COVER_IMGS.length > 0;
  const credited = /cover art on this page is from IGDB\.com/i.test(html);
  return [shipped ? credited : !credited, shipped ? `credited: ${credited}` : 'no covers, no credit — correct'];
});

check('every mockup caption still discloses that it is not real', 'there is no runnable TUI to screenshot yet', () => {
  // Deliberately broader than the word "mockup": §11's phone frames say
  // "illustration … Not built yet", which is the same disclosure in different
  // words and is not this ticket's to touch. What must never happen is a
  // caption losing its disclosure entirely.
  const captions = [...html.matchAll(/class="[^"]*z-mockup-caption[^"]*"[^>]*>([\s\S]*?)<\/figcaption>/g)]
    .map(m => m[1].replace(/<[^>]+>/g, ' '));
  const HONEST = /mockup|illustration|not a screenshot|not built yet|not available yet/i;
  const bare = captions.filter(c => !HONEST.test(c));
  return [captions.length > 0 && bare.length === 0,
    `${captions.length} caption(s), ${bare.length} with no disclosure`];
});

// --- the copy file and the built page agree ---------------------------------
// docs/content/landing-copy.md is the contract; the page is the artefact. A
// reviewer checks them against each other by hand, so check them here instead —
// every §06 cover alt string quoted in the copy file must appear verbatim in the
// rendered page, and vice versa. Skipped (not failed) when the copy file is not
// beside the build, so this stays runnable against a bare dist/.
check('the copy file and the page agree on every cover alt string', 'the alt-text contract', () => {
  const copyPath = 'docs/content/landing-copy.md';
  if (!existsSync(copyPath)) return [true, 'copy file not present — skipped'];
  const copy = readFileSync(copyPath, 'utf8');
  const six = copy.slice(copy.indexOf('## 06 · one-collection'), copy.indexOf('## 07 ·'));
  const declared = [...six.matchAll(/^\d{1,2}\. "([^"]+)"$/gm)].map(m => m[1]);
  const decode = t => t.replace(/&#39;/g, "'").replace(/&amp;/g, '&').replace(/&quot;/g, '"');
  const rendered = COVER_IMGS.map(t => decode(attr(t, 'alt') ?? ''));
  if (rendered.length === 0) return [declared.length === 12, `no covers shipped; ${declared.length} declared`];
  // The two directions are NOT symmetric, and treating them as one is what
  // broke the withdrawal path (see below).
  //
  // `extra` — the page renders an alt string the copy file never declared — is
  // always drift, and is always a failure. The page may not invent alt text.
  //
  // `missing` — the copy file declares a row the page did not render — is a
  // failure ONLY if that row actually shipped an image. A row whose cover is
  // absent from the build renders its art-directed tile and contributes no
  // <img> and no alt, which is the CORRECT behaviour and a documented state:
  // `docs/legal/igdb-image-licence-finding.md` §7 names a rights-holder
  // objecting to one specific cover, "in which case that row loses its image
  // and the tile treatment renders in its place — a data-only change".
  //
  // It was not a data-only change. Deleting one cover's files failed this
  // assertion, so the one remedy the project would need to execute FASTEST,
  // under legal pressure, was blocked by its own build gate. The all-absent
  // case was already handled above; only the partial case — which is exactly
  // the rights-holder case — was not.
  //
  // So a declared-but-unrendered row is excused precisely when its files are
  // genuinely not in the build, and not otherwise. The contract keeps its
  // teeth: the copy file cannot silently lose a row whose cover still ships.
  const extra = rendered.filter(r => !declared.includes(r));
  const missing = declared.filter(d => !rendered.includes(d));

  const gridPath = 'site/src/data/coverGrid.ts';
  const coversDir = ['site/public/covers', join(dist, 'covers')].find(existsSync);
  let unexplained = missing;
  let excused = [];
  if (missing.length && existsSync(gridPath) && coversDir) {
    // alt -> slug, straight out of the row table, so this cannot drift either.
    const grid = readFileSync(gridPath, 'utf8');
    const bySlug = new Map(
      [...grid.matchAll(/slug:\s*'([a-z0-9-]+)',?\s*\n?\s*alt:\s*'((?:[^'\\]|\\.)*)'/g)]
        .map(m => [m[2].replace(/\\'/g, "'"), m[1]])
    );
    excused = missing.filter(d => {
      const slug = bySlug.get(d);
      return slug && !existsSync(join(coversDir, `${slug}.jpg`));
    });
    unexplained = missing.filter(d => !excused.includes(d));
  }

  return [unexplained.length === 0 && extra.length === 0,
    `${declared.length} declared, ${rendered.length} rendered; ` +
    `unexplained: ${unexplained.join(' | ') || 'none'}; ` +
    `withdrawn (no cover file, renders its tile): ${excused.length}; ` +
    `unlisted: ${extra.join(' | ') || 'none'}`];
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
