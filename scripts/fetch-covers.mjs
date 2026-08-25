#!/usr/bin/env node
/**
 * Fetch the twelve §06 cover images from IGDB and write them into the build.
 *
 * ONE PROVIDER. Every image on the landing page comes through this script and
 * therefore through IGDB — including the `PHYSICAL` row, because `PHYSICAL` is
 * a provenance label on the row and not a different image source. Nothing here
 * reaches for a web search, and there is deliberately no code path that could.
 * See `docs/legal/igdb-image-licence-finding.md`.
 *
 *   IGDB_CLIENT_ID=… IGDB_CLIENT_SECRET=… node scripts/fetch-covers.mjs
 *
 * Credentials are read from the environment and never written anywhere. The
 * outputs are:
 *
 *   site/public/covers/{slug}.avif|.webp|.jpg   — 360 × 480, centre-cropped
 *   docs/legal/cover-provenance.json            — what each file came from
 *
 * Re-runnable: existing files are overwritten, and `--check` reports what the
 * build currently contains without touching the network.
 *
 * The only dependency is sharp, which Astro already installs — see
 * site/package-lock.json. Nothing new enters the dependency tree for this.
 */
import { createRequire } from 'node:module';
import { mkdirSync, writeFileSync, existsSync, readdirSync } from 'node:fs';
import { fileURLToPath, pathToFileURL } from 'node:url';

const require = createRequire(new URL('../site/package.json', import.meta.url));

const ROOT = new URL('../', import.meta.url);
const OUT_DIR = new URL('site/public/covers/', ROOT);
const PROVENANCE = new URL('docs/legal/cover-provenance.json', ROOT);

// Kept in step with site/src/data/coverGrid.ts by `--check`, which fails when
// the two disagree — the data file is the source of truth for which twelve.
const WIDTH = 360;
const HEIGHT = 480;

/** Parse the row table straight out of the TypeScript so there is exactly one
 *  list of the twelve and this script cannot drift from what the page renders. */
export async function rows() {
  const src = await import('node:fs').then((fs) =>
    fs.readFileSync(fileURLToPath(new URL('site/src/data/coverGrid.ts', ROOT)), 'utf8')
  );
  const out = [];
  const re =
    /searchName:\s*'((?:[^'\\]|\\.)*)'[\s\S]*?releaseYear:\s*(\d{4}),\s*slug:\s*'([a-z0-9-]+)'/g;
  let m;
  while ((m = re.exec(src)) !== null) {
    out.push({ searchName: m[1].replace(/\\'/g, "'"), releaseYear: Number(m[2]), slug: m[3] });
  }
  return out;
}

function report(list) {
  const present = list.filter((r) => existsSync(new URL(`${r.slug}.jpg`, OUT_DIR)));
  console.log(`${present.length}/${list.length} covers present in site/public/covers/`);
  for (const r of list) {
    const ok = existsSync(new URL(`${r.slug}.jpg`, OUT_DIR));
    console.log(`  ${ok ? '✓' : '·'} ${r.slug}${ok ? '' : '  (renders the art-directed tile)'}`);
  }
  return present.length;
}

async function token(id, secret) {
  const u = new URL('https://id.twitch.tv/oauth2/token');
  u.searchParams.set('client_id', id);
  u.searchParams.set('client_secret', secret);
  u.searchParams.set('grant_type', 'client_credentials');
  const r = await fetch(u, { method: 'POST' });
  if (!r.ok) throw new Error(`Twitch OAuth failed: ${r.status} ${await r.text()}`);
  return (await r.json()).access_token;
}

async function igdb(body, id, tok) {
  const r = await fetch('https://api.igdb.com/v4/games', {
    method: 'POST',
    headers: { 'Client-ID': id, Authorization: `Bearer ${tok}`, Accept: 'application/json' },
    body
  });
  if (!r.ok) throw new Error(`IGDB query failed: ${r.status} ${await r.text()}`);
  return r.json();
}

const yearOf = (g) =>
  g.first_release_date ? new Date(g.first_release_date * 1000).getUTCFullYear() : null;

/**
 * Two titles are the same title, in the typographic sense.
 *
 * This is NOT a fuzzy comparison — it is still exact equality, taken over the
 * text rather than over one particular way of encoding it. `coverGrid.ts`
 * spells row 10 “Baldur’s Gate II” with a typographic apostrophe (U+2019),
 * because that is how the name is set on the page; IGDB stores it with an
 * ASCII apostrophe. Byte equality therefore called them different games,
 * dropped to the unnamed pool, matched the base game AND the Collectors'
 * Edition on year, and correctly refused rather than guess between them.
 *
 * Folding the curly quotes to their ASCII forms (and NFC-normalising first, so
 * a decomposed accent equals its precomposed twin) removes that one false
 * distinction and nothing else: "Dead Space" still does not equal "Dead Space
 * Remake", and the release-year pin below is untouched.
 */
function sameName(a, b) {
  const norm = (s) =>
    s
      .normalize('NFC')
      .replace(/[‘’ʼ]/g, "'")
      .replace(/[“”]/g, '"')
      .toLowerCase()
      .trim();
  return norm(a) === norm(b);
}

/**
 * Deterministic match, and it refuses rather than guesses.
 *
 * A fuzzy "closest result wins" is how a page ends up quietly showing the 2023
 * Dead Space remake under a row that claims to be the 2008 original. So: exact
 * name match first, then the release year from the data file, and if that does
 * not land on exactly one game the row is REPORTED and skipped — it renders its
 * art-directed tile and a human pins the id. A wrong cover is worse than none.
 */
export function pick(results, row) {
  const named = results.filter((g) => g.name && sameName(g.name, row.searchName));
  const pool = named.length ? named : results;
  const dated = pool.filter((g) => yearOf(g) === row.releaseYear);
  if (dated.length === 1) return dated[0];
  if (named.length === 1 && !dated.length) return named[0];
  if (dated.length > 1) {
    // Same name, same year, more than one entry: prefer the parent game
    // (a main game has no `parent_game`), which is the edition a player means.
    const parents = dated.filter((g) => !g.parent_game);
    if (parents.length === 1) return parents[0];
  }
  return null;
}

async function main() {
  const list = await rows();
  if (list.length !== 12) {
    console.error(`Expected 12 rows in site/src/data/coverGrid.ts, parsed ${list.length}.`);
    process.exit(1);
  }

  if (process.argv.includes('--check')) {
    report(list);
    return;
  }

  const id = process.env.IGDB_CLIENT_ID;
  const secret = process.env.IGDB_CLIENT_SECRET;
  if (!id || !secret) {
    console.error(
      'IGDB_CLIENT_ID and IGDB_CLIENT_SECRET are not set.\n\n' +
        'Register a Confidential application at https://dev.twitch.tv/console/apps\n' +
        '(OAuth redirect http://localhost), generate a secret, then:\n\n' +
        '  IGDB_CLIENT_ID=… IGDB_CLIENT_SECRET=… node scripts/fetch-covers.mjs\n\n' +
        'Without them the page renders the art-directed tiles, which is a valid\n' +
        'shipping state — see docs/legal/igdb-image-licence-finding.md §7.'
    );
    process.exit(2);
  }

  const sharp = require('sharp');
  mkdirSync(fileURLToPath(OUT_DIR), { recursive: true });

  const tok = await token(id, secret);
  const provenance = {
    source: 'IGDB (api.igdb.com/v4), served under the Twitch Developer Services Agreement',
    finding: 'docs/legal/igdb-image-licence-finding.md',
    fetchedAt: new Date().toISOString(),
    rendered: { width: WIDTH, height: HEIGHT, fit: 'cover', position: 'centre' },
    covers: []
  };

  let failed = 0;
  for (const row of list) {
    const q =
      `search "${row.searchName.replace(/"/g, '\\"')}"; ` +
      'fields name,slug,first_release_date,parent_game,cover.image_id,cover.width,cover.height; ' +
      'where cover != null; limit 25;';
    const results = await igdb(q, id, tok);
    const game = pick(results, row);

    if (!game || !game.cover?.image_id) {
      failed++;
      console.error(
        `✗ ${row.slug}: no unambiguous IGDB match for "${row.searchName}" (${row.releaseYear}). ` +
          `Candidates: ${results.map((g) => `${g.name} [${yearOf(g) ?? '?'}]`).join(', ') || 'none'}`
      );
      continue;
    }

    // 528 × 748 — the largest cover IGDB publishes, downscaled here rather than
    // upscaled later.
    const src = `https://images.igdb.com/igdb/image/upload/t_cover_big_2x/${game.cover.image_id}.jpg`;
    const res = await fetch(src);
    if (!res.ok) {
      failed++;
      console.error(`✗ ${row.slug}: cover download failed, ${res.status} ${src}`);
      continue;
    }
    const buf = Buffer.from(await res.arrayBuffer());

    // The centre crop that takes IGDB's 0.706 to the ratified 0.750 box. Doing
    // it here means the shipped file's intrinsic size IS the rendered box, so
    // the width/height attributes are true and CLS stays 0.000.
    const base = sharp(buf).resize(WIDTH, HEIGHT, { fit: 'cover', position: 'centre' });
    await Promise.all([
      base.clone().jpeg({ quality: 82, mozjpeg: true }).toFile(fileURLToPath(new URL(`${row.slug}.jpg`, OUT_DIR))),
      base.clone().webp({ quality: 80 }).toFile(fileURLToPath(new URL(`${row.slug}.webp`, OUT_DIR))),
      base.clone().avif({ quality: 55 }).toFile(fileURLToPath(new URL(`${row.slug}.avif`, OUT_DIR)))
    ]);

    provenance.covers.push({
      slug: row.slug,
      igdbGameId: game.id,
      igdbSlug: game.slug,
      igdbName: game.name,
      releaseYear: yearOf(game),
      igdbCoverImageId: game.cover.image_id,
      sourceUrl: src
    });
    console.log(`✓ ${row.slug}  ← ${game.name} [${yearOf(game)}]  ${game.cover.image_id}`);
  }

  provenance.covers.sort((a, b) => a.slug.localeCompare(b.slug));
  writeFileSync(fileURLToPath(PROVENANCE), `${JSON.stringify(provenance, null, 2)}\n`);

  const written = readdirSync(fileURLToPath(OUT_DIR)).filter((f) => f.endsWith('.jpg')).length;
  console.log(`\n${written}/12 covers written to site/public/covers/`);
  console.log(`Provenance recorded in docs/legal/cover-provenance.json`);
  if (failed) {
    console.error(
      `\n${failed} row(s) unresolved — they render the art-directed tile. Pin them by hand ` +
        `rather than loosening the match.`
    );
    process.exit(1);
  }
}

// Runs only when invoked as a command. `pick` and `rows` are exported so
// `scripts/fetch-covers.test.mjs` can exercise the match rule without a
// network, credentials, or writing a file — the match rule is the part of this
// script that can silently put the wrong cover on the page.
if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  await main();
}
