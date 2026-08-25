/**
 * The twelve rows of the §06 shelf grid.
 *
 * Every image is fetched from IGDB by `scripts/fetch-covers.mjs` and served
 * from this repository — no hotlink, no web-sourced artwork, one provider and
 * therefore one licence question. See
 * `docs/legal/igdb-image-licence-finding.md`.
 *
 * `platformTag` is a PROVENANCE label on the row, not a second image source: a
 * SNES cartridge has an IGDB entry exactly as a Steam title does, and Zerado
 * will fetch a physical copy's cover the same way it fetches a synced one.
 *
 * The per-position `treatment` × `hue` sequence is UNCHANGED from the ratified
 * design and is still the renderer for any row whose image file is absent — a
 * cover pulled for any reason degrades to the tile it replaced, and nothing
 * else on the page moves. Ref: `docs/design/blueprint.md` §7.5 (as amended
 * 2026-08-25), `blueprint.tokens.json` coverGrid.sequence.
 */
export type Treatment = 'A' | 'B' | 'C' | 'D';
export type Hue = 'amber' | 'cyan' | 'steel' | 'orchid';
export type PlatformTagKind = 'STEAM' | 'GOG' | 'EA' | 'PS' | 'PHYSICAL';

export interface CoverTileDef {
  /** Position in the ratified sequence, 1–12. */
  n: number;
  /** Art-directed fallback treatment, used when this row has no image file. */
  treatment: Treatment;
  /** Art-directed fallback hue key. */
  hue: Hue;
  /** The plate rendered bottom-left. A provenance label, not an image source. */
  platformTag: PlatformTagKind;
  /** The game, as a reader would name it. */
  title: string;
  /** What `searchName` IGDB is asked for. Kept separate from `title` so the
   *  displayed name never has to bend to match a database's spelling. */
  searchName: string;
  /** The year of the release we mean. Pins the match: "Dead Space" alone would
   *  as happily resolve to the 2023 remake, and "God of War" to the 2005 one. */
  releaseYear: number;
  /** Basename of the served files: `/covers/{slug}.avif|.webp|.jpg`. */
  slug: string;
  /** Alt text. Names the game AND its platform; all twelve are distinct.
   *  Mirrored verbatim in `docs/content/landing-copy.md` §06. */
  alt: string;
}

/**
 * The twelve, chosen against three constraints: they carry the platform tags
 * the grid already asserts (6 STEAM · 2 GOG · 1 EA · 2 PS · 1 PHYSICAL, with
 * PHYSICAL at position 9 exactly as ratified); they span 1994–2020 so the
 * retro and physical case reads as deliberate rather than accidental; and each
 * one is a title the audience recognises on sight. The per-position
 * platform-tag mapping is unchanged from the ratified sequence.
 */
export const COVER_GRID_SEQUENCE: CoverTileDef[] = [
  {
    n: 1, treatment: 'A', hue: 'amber', platformTag: 'STEAM',
    title: 'Half-Life 2', searchName: 'Half-Life 2', releaseYear: 2004, slug: 'half-life-2',
    alt: 'Half-Life 2 — cover art, in the Steam part of the library'
  },
  {
    n: 2, treatment: 'B', hue: 'cyan', platformTag: 'STEAM',
    title: 'Hades', searchName: 'Hades', releaseYear: 2020, slug: 'hades',
    alt: 'Hades — cover art, in the Steam part of the library'
  },
  {
    n: 3, treatment: 'C', hue: 'steel', platformTag: 'GOG',
    title: 'The Witcher 3: Wild Hunt', searchName: 'The Witcher 3: Wild Hunt',
    releaseYear: 2015, slug: 'the-witcher-3-wild-hunt',
    alt: 'The Witcher 3: Wild Hunt — cover art, in the GOG part of the library'
  },
  {
    n: 4, treatment: 'D', hue: 'orchid', platformTag: 'STEAM',
    title: 'Portal 2', searchName: 'Portal 2', releaseYear: 2011, slug: 'portal-2',
    alt: 'Portal 2 — cover art, in the Steam part of the library'
  },
  {
    n: 5, treatment: 'B', hue: 'amber', platformTag: 'PS',
    title: 'Bloodborne', searchName: 'Bloodborne', releaseYear: 2015, slug: 'bloodborne',
    alt: 'Bloodborne — cover art, in the PlayStation part of the library'
  },
  {
    n: 6, treatment: 'A', hue: 'cyan', platformTag: 'STEAM',
    title: 'Stardew Valley', searchName: 'Stardew Valley', releaseYear: 2016, slug: 'stardew-valley',
    alt: 'Stardew Valley — cover art, in the Steam part of the library'
  },
  {
    n: 7, treatment: 'D', hue: 'steel', platformTag: 'EA',
    title: 'Dead Space', searchName: 'Dead Space', releaseYear: 2008, slug: 'dead-space',
    alt: 'Dead Space — cover art, in the EA part of the library'
  },
  {
    n: 8, treatment: 'C', hue: 'orchid', platformTag: 'STEAM',
    title: 'Celeste', searchName: 'Celeste', releaseYear: 2018, slug: 'celeste',
    alt: 'Celeste — cover art, in the Steam part of the library'
  },
  {
    // The one hand-added row. Its cover comes from IGDB exactly like the other
    // eleven — PHYSICAL describes where the copy lives, not where the art came
    // from. This is the tile the whole section argues about.
    n: 9, treatment: 'B', hue: 'orchid', platformTag: 'PHYSICAL',
    title: 'Super Metroid', searchName: 'Super Metroid', releaseYear: 1994, slug: 'super-metroid',
    alt: 'Super Metroid — cover art, a SNES cartridge added to the library by hand'
  },
  {
    n: 10, treatment: 'A', hue: 'steel', platformTag: 'GOG',
    title: 'Baldur’s Gate II: Shadows of Amn',
    searchName: 'Baldur’s Gate II: Shadows of Amn',
    releaseYear: 2000, slug: 'baldurs-gate-ii-shadows-of-amn',
    alt: 'Baldur’s Gate II: Shadows of Amn — cover art, in the GOG part of the library'
  },
  {
    n: 11, treatment: 'D', hue: 'amber', platformTag: 'PS',
    title: 'God of War', searchName: 'God of War', releaseYear: 2018, slug: 'god-of-war',
    alt: 'God of War — cover art, in the PlayStation part of the library'
  },
  {
    n: 12, treatment: 'C', hue: 'cyan', platformTag: 'STEAM',
    title: 'Hollow Knight', searchName: 'Hollow Knight', releaseYear: 2017, slug: 'hollow-knight',
    alt: 'Hollow Knight — cover art, in the Steam part of the library'
  }
];

/**
 * Rendered at 3 / 4 — the ratified tile geometry, unchanged.
 *
 * IGDB serves covers at 264 × 374 (`t_cover_big`, or 528 × 748 at `_2x`), which
 * is 0.706 rather than 0.750. **Aspect-ratio policy: one uniform 3 / 4 box, a
 * centred crop taken at fetch time, never a letterbox and never a per-tile
 * ratio.** The crop removes about 3% from the top and 3% from the bottom of the
 * source, which box art tolerates; a letterbox would have put twelve different
 * pillar colours into a six-colour palette, and per-tile ratios would have made
 * the rows ragged — the exact failure the ticket named. The files are written
 * at these dimensions so the `width`/`height` attributes on every `<img>` are
 * the file's true intrinsic size and CLS stays 0.000.
 *
 * **Sized to the RENDERED box, not to the source.** Measured with Playwright,
 * a tile is 104 CSS px wide at 375, 168 at 768 and 175 at 1280 and above — so
 * 350 px covers the widest tile on a DPR 2 display, and 360 × 480 is the
 * nearest exact 3 / 4 above it. Shipping IGDB's full 528 × 748 would have been
 * roughly 2.2× the pixels for something no display can resolve.
 */
export const COVER_WIDTH = 360;
export const COVER_HEIGHT = 480;
