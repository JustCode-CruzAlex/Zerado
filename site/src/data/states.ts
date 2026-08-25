/**
 * The brand's four game-state definitions — colour token, glyph, label.
 * Single source for StateChip, TerminalRow and any other co-render surface.
 * Ref: brand/tokens.css §3, design/blueprint.md §2 (§05), blueprint.tokens.json color.state.
 */
export type StateKey = 'notStarted' | 'inProgress' | 'zerado' | 'abandoned';

export interface StateDef {
  key: StateKey;
  token: string;
  glyph: string;
  label: string;
}

export const STATES: Record<StateKey, StateDef> = {
  notStarted: { key: 'notStarted', token: '--z-state-not-started', glyph: '○', label: 'NOT STARTED' },
  inProgress: { key: 'inProgress', token: '--z-state-in-progress', glyph: '◐', label: 'IN PROGRESS' },
  zerado: { key: 'zerado', token: '--z-state-zerado', glyph: '◉', label: 'ZERADO' },
  abandoned: { key: 'abandoned', token: '--z-state-abandoned', glyph: '⊘', label: 'ABANDONED' }
};

export const STATE_ORDER: StateKey[] = ['notStarted', 'inProgress', 'zerado', 'abandoned'];
