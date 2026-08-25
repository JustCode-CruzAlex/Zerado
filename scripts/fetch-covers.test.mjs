/**
 * Tests for the parts of `fetch-covers.mjs` that can put the WRONG cover on the
 * page, and for the shipped files' intrinsic size.
 *
 * The fetch itself is not tested — it needs credentials and a network, and its
 * failure mode is loud (a non-zero exit and a named row). The failure mode that
 * is SILENT is `pick()` choosing a plausible-but-wrong game: a page that shows
 * the 2023 Dead Space remake under a row asserting the 2008 original looks
 * entirely fine and is wrong. So `pick` is exercised against the real shapes
 * IGDB returns, including the two ambiguities this repository actually hit.
 *
 *   node --test scripts/*.test.mjs
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync, existsSync, statSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { pick, rows } from './fetch-covers.mjs';

const ROOT = new URL('../', import.meta.url);
const COVERS = new URL('site/public/covers/', ROOT);

/** IGDB dates are UNIX seconds; the tests read in years. */
const y = (year) => ({ first_release_date: Date.UTC(year, 5, 1) / 1000 });
const game = (id, name, year, extra = {}) => ({
  id,
  name,
  ...y(year),
  cover: { image_id: `co${id}` },
  ...extra
});

// --- pick(): the match rule ---------------------------------------------

test('a typographic apostrophe in the data file matches IGDB’s ASCII one', () => {
  // The real regression. `coverGrid.ts` sets row 10 with U+2019 because that is
  // how the title reads on the page; IGDB stores U+0027. Byte equality found no
  // named match, fell through to the unnamed pool, and drew the base game AND
  // the Collectors' Edition on year — so the script correctly refused and
  // shipped 11 of 12.
  const results = [
    game(71, "Baldur's Gate II: Shadows of Amn", 2000),
    game(72, "Baldur's Gate II: Shadows of Amn - Collectors' Edition", 2000)
  ];
  const chosen = pick(results, {
    searchName: 'Baldur’s Gate II: Shadows of Amn',
    releaseYear: 2000
  });
  assert.equal(chosen?.id, 71, 'must pick the base game, not the Collectors’ Edition');
});

test('the release year pins the edition a row actually means', () => {
  // Same name, two games, twenty-five years apart. Only the year separates the
  // 2008 original from the 2023 remake, and the row asserts 2008.
  const results = [game(1, 'Dead Space', 2008), game(2, 'Dead Space', 2023)];
  assert.equal(pick(results, { searchName: 'Dead Space', releaseYear: 2008 })?.id, 1);
  assert.equal(pick(results, { searchName: 'Dead Space', releaseYear: 2023 })?.id, 2);
});

test('an exact name beats a longer one that merely contains it', () => {
  const results = [game(1, 'God of War', 2018), game(2, 'God of War: Ragnarök', 2018)];
  assert.equal(pick(results, { searchName: 'God of War', releaseYear: 2018 })?.id, 1);
});

test('the normalisation folds quotes and nothing else', () => {
  // The fold must not become a fuzzy match by degrees. A different title that
  // happens to start with the row's title is still a different title, and with
  // no exact name and no year hit the rule refuses.
  const results = [game(1, 'Dead Space Remake', 2023)];
  assert.equal(pick(results, { searchName: 'Dead Space', releaseYear: 2008 }), null);
});

test('two indistinguishable candidates are refused, not guessed between', () => {
  // A wrong cover is worse than none: the row falls back to its art-directed
  // tile and a human pins the id.
  const results = [game(1, 'Hades', 2020), game(2, 'Hades', 2020)];
  assert.equal(pick(results, { searchName: 'Hades', releaseYear: 2020 }), null);
});

test('among same-name same-year entries the parent game wins', () => {
  // A DLC or an edition carries `parent_game`; the main game does not, and the
  // main game is what a player means by the title.
  const results = [
    game(1, 'Hollow Knight', 2017),
    game(2, 'Hollow Knight', 2017, { parent_game: 1 })
  ];
  assert.equal(pick(results, { searchName: 'Hollow Knight', releaseYear: 2017 })?.id, 1);
});

test('a sole exact name with no year match is still accepted', () => {
  // IGDB's `first_release_date` is the earliest release across all platforms
  // and can disagree with the row by a year. One unambiguous name is enough.
  const results = [game(1, 'Celeste', 2018)];
  assert.equal(pick(results, { searchName: 'Celeste', releaseYear: 2019 })?.id, 1);
});

test('no candidates yields no cover', () => {
  assert.equal(pick([], { searchName: 'Half-Life 2', releaseYear: 2004 }), null);
});

// --- rows(): the script and the page read the same twelve ----------------

test('the row table parses exactly the twelve the page renders', async () => {
  // `rows()` scrapes coverGrid.ts with a regex so there is one list, not two.
  // A refactor of the data file that the regex stops matching would silently
  // shrink the fetch; the script exits on !== 12, and this pins it in CI.
  const list = await rows();
  assert.equal(list.length, 12);
  assert.equal(new Set(list.map((r) => r.slug)).size, 12, 'slugs are distinct');
  for (const r of list) {
    assert.match(r.slug, /^[a-z0-9-]+$/);
    assert.ok(r.searchName.length > 0);
    assert.ok(r.releaseYear >= 1970 && r.releaseYear <= 2100);
  }
});

// --- the shipped files ----------------------------------------------------

/** Intrinsic size straight out of the JPEG SOF marker — no dependency, and it
 *  reads the FILE rather than trusting what the fetch script reported. */
function jpegSize(path) {
  const b = readFileSync(path);
  let i = 2; // skip SOI
  while (i < b.length) {
    if (b[i] !== 0xff) return null;
    const marker = b[i + 1];
    const len = b.readUInt16BE(i + 2);
    // SOF0/1/2/3, 5-7, 9-11, 13-15 all carry height then width at +5.
    if (marker >= 0xc0 && marker <= 0xcf && ![0xc4, 0xc8, 0xcc].includes(marker)) {
      return { height: b.readUInt16BE(i + 5), width: b.readUInt16BE(i + 7) };
    }
    i += 2 + len;
  }
  return null;
}

test('every shipped cover is exactly the box the page reserves', async (t) => {
  // The CLS-is-0.000 claim rests on the <img> width/height being the file's
  // TRUE intrinsic size. check-page.mjs asserts the attributes are present;
  // this asserts the pixels behind them. Covers are a valid absent state
  // (licence finding §7), so skip rather than fail when none shipped.
  //
  // This test deliberately does NOT judge HOW MANY covers ship.
  //
  // It used to: `assert.equal(shipped.length, 12)`. 0 skipped, 12 passed, and
  // 1–11 failed — which is exactly the rights-holder withdrawal of §7. The
  // commit that unblocked that path in check-page.mjs added this assertion to
  // the SAME site.yml job six lines below, so the job still went red and the
  // block moved rather than lifted. Two gates adjudicating partiality by
  // different rules is what produced that; `check-page.mjs` is now the single
  // authority on it, and this file sticks to per-file integrity.
  const list = await rows();
  const shipped = list.filter((r) => existsSync(new URL(`${r.slug}.jpg`, COVERS)));
  if (shipped.length === 0) return t.skip('no covers in this build — the fallback state');

  const grid = readFileSync(fileURLToPath(new URL('site/src/data/coverGrid.ts', ROOT)), 'utf8');
  const w = Number(grid.match(/COVER_WIDTH\s*=\s*(\d+)/)[1]);
  const h = Number(grid.match(/COVER_HEIGHT\s*=\s*(\d+)/)[1]);

  for (const r of shipped) {
    const size = jpegSize(fileURLToPath(new URL(`${r.slug}.jpg`, COVERS)));
    assert.deepEqual(size, { width: w, height: h }, `${r.slug}.jpg is not ${w}×${h}`);
    // The <picture> offers all three; a missing modern format silently costs
    // every visitor the JPEG, which is roughly twice the bytes.
    for (const ext of ['avif', 'webp']) {
      const f = fileURLToPath(new URL(`${r.slug}.${ext}`, COVERS));
      assert.ok(existsSync(f), `${r.slug}.${ext} is missing`);
      assert.ok(statSync(f).size > 0, `${r.slug}.${ext} is empty`);
    }
  }
});

// --- the withdrawal path stays open --------------------------------------

test('every row’s alt maps back to its slug', async () => {
  // `check-page.mjs` excuses a declared-but-unrendered alt string only when
  // that row's cover file is genuinely absent — the rights-holder withdrawal
  // of `docs/legal/igdb-image-licence-finding.md` §7. It resolves alt -> slug
  // with a regex over coverGrid.ts that assumes `slug:` sits immediately before
  // `alt:`. Reorder those two fields and the map comes back EMPTY, nothing is
  // ever excused, and the withdrawal path silently re-blocks — a build gate
  // standing in front of a legal remedy, failing the way it used to.
  //
  // This pins the shape that assumption rests on. It is cheap here and would
  // cost a second full build to catch anywhere else.
  const grid = readFileSync(fileURLToPath(new URL('site/src/data/coverGrid.ts', ROOT)), 'utf8');
  const pairs = [...grid.matchAll(/slug:\s*'([a-z0-9-]+)',?\s*\n?\s*alt:\s*'((?:[^'\\]|\\.)*)'/g)];
  assert.equal(pairs.length, 12, 'the alt -> slug regex in check-page.mjs no longer matches all twelve');

  const slugs = new Set((await rows()).map((r) => r.slug));
  for (const [, slug, alt] of pairs) {
    assert.ok(slugs.has(slug), `${slug} is not one of the twelve rows`);
    assert.ok(alt.trim().length > 0, `${slug} has an empty alt`);
  }
  assert.equal(new Set(pairs.map((m) => m[2])).size, 12, 'alt strings are distinct');
});

// --- §8: the legal record covers everything that ships --------------------

test('every shipped cover has a provenance row', async () => {
  // `docs/legal/igdb-image-licence-finding.md` §8 states it as an invariant:
  // "Every file under `site/public/covers/` appears there." It is the section a
  // rights-holder question gets answered from, so a shipped image with no
  // recorded origin is the one state that must never reach main.
  //
  // It was reachable. `fetch-covers.mjs` rebuilds `provenance.covers` from
  // scratch each run and pushes only rows that resolve, while a failed row's
  // previously-downloaded files stay on disk untouched — so a partial re-run
  // (the "pin them by hand" case the script itself prints) silently dropped
  // rows from the record while their covers stayed shipped. The script now
  // carries those rows forward; this asserts the result rather than trusting it.
  //
  // One direction only, on purpose: a provenance row with no file is a
  // WITHDRAWN cover, which is the sanctioned §7 state and not a defect.
  const provPath = fileURLToPath(new URL('docs/legal/cover-provenance.json', ROOT));
  const list = await rows();
  const shipped = list.filter((r) => existsSync(new URL(`${r.slug}.jpg`, COVERS)));
  if (shipped.length === 0) return;

  assert.ok(existsSync(provPath), 'covers ship but docs/legal/cover-provenance.json is missing');
  const recorded = new Map(
    JSON.parse(readFileSync(provPath, 'utf8')).covers.map((c) => [c.slug, c])
  );

  const undocumented = shipped.map((r) => r.slug).filter((slug) => !recorded.has(slug));
  assert.deepEqual(undocumented, [], 'shipped cover(s) with no recorded IGDB origin');

  for (const r of shipped) {
    const c = recorded.get(r.slug);
    for (const field of ['igdbGameId', 'igdbSlug', 'igdbName', 'igdbCoverImageId', 'sourceUrl']) {
      assert.ok(c[field] !== undefined && c[field] !== null && c[field] !== '',
        `${r.slug}: provenance is missing ${field}`);
    }
    assert.match(c.sourceUrl, /^https:\/\/images\.igdb\.com\//,
      `${r.slug}: provenance sourceUrl is not an IGDB image URL`);
  }
});
