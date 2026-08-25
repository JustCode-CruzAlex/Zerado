/**
 * The twelve art-directed cover tiles — fixed sequence, treatment × hue,
 * never real game cover art. Ref: design/blueprint.md §7.5,
 * blueprint.tokens.json coverGrid.sequence (verbatim).
 */
export type Treatment = 'A' | 'B' | 'C' | 'D';
export type Hue = 'amber' | 'cyan' | 'steel' | 'orchid';
export type PlatformTagKind = 'STEAM' | 'GOG' | 'EA' | 'PS' | 'PHYSICAL';

export interface CoverTileDef {
  n: number;
  treatment: Treatment;
  hue: Hue;
  platformTag: PlatformTagKind;
}

export const COVER_GRID_SEQUENCE: CoverTileDef[] = [
  { n: 1, treatment: 'A', hue: 'amber', platformTag: 'STEAM' },
  { n: 2, treatment: 'B', hue: 'cyan', platformTag: 'STEAM' },
  { n: 3, treatment: 'C', hue: 'steel', platformTag: 'GOG' },
  { n: 4, treatment: 'D', hue: 'orchid', platformTag: 'STEAM' },
  { n: 5, treatment: 'B', hue: 'amber', platformTag: 'PS' },
  { n: 6, treatment: 'A', hue: 'cyan', platformTag: 'STEAM' },
  { n: 7, treatment: 'D', hue: 'steel', platformTag: 'EA' },
  { n: 8, treatment: 'C', hue: 'orchid', platformTag: 'STEAM' },
  { n: 9, treatment: 'B', hue: 'orchid', platformTag: 'PHYSICAL' },
  { n: 10, treatment: 'A', hue: 'steel', platformTag: 'GOG' },
  { n: 11, treatment: 'D', hue: 'amber', platformTag: 'PS' },
  { n: 12, treatment: 'C', hue: 'cyan', platformTag: 'STEAM' }
];
