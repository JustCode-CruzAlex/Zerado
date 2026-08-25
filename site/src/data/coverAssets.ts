/**
 * What the build actually contains, read from disk at build time.
 *
 * Every statement the page makes about its cover art — the twelve `<img>`
 * elements, the caption under the grid, and the IGDB credit in the footer — is
 * derived from this one function. The page therefore cannot credit a source it
 * did not ship from, and cannot disclose "illustrative artwork" while showing
 * real covers. The two states stay consistent because they read the same fact
 * rather than two flags someone has to remember to flip together.
 *
 * Ref: docs/legal/igdb-image-licence-finding.md §6 (attribution) and §7 (the
 * reversal path this makes cheap).
 */
import { existsSync } from 'node:fs';
import { resolve } from 'node:path';
import { COVER_GRID_SEQUENCE, type CoverTileDef } from './coverGrid';

export interface CoverTileState extends CoverTileDef {
  hasImage: boolean;
}

// Resolved against the working directory, not `import.meta.url`: this module is
// bundled before it runs, so `import.meta.url` points at the emitted chunk and
// a relative walk from it lands nowhere near `public/`. The build runs from
// `site/` (`npm run build`, and `working-directory: site` in CI); the second
// candidate covers a run from the repository root.
const COVERS_DIR =
  [resolve(process.cwd(), 'public/covers'), resolve(process.cwd(), 'site/public/covers')].find(
    (d) => existsSync(d)
  ) ?? resolve(process.cwd(), 'public/covers');

export function coverTiles(): CoverTileState[] {
  return COVER_GRID_SEQUENCE.map((tile) => ({
    ...tile,
    hasImage: existsSync(resolve(COVERS_DIR, `${tile.slug}.jpg`))
  }));
}

/** True when at least one real cover ships — the trigger for crediting IGDB. */
export function hasRealCovers(): boolean {
  return coverTiles().some((t) => t.hasImage);
}
