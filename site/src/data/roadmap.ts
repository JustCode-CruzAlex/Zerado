/**
 * The four ratified roadmap phases, verbatim from content/landing-copy.md §12.
 */
export const PHASES = [
  { phase: 1, name: 'CLI/TUI MVP', line: 'Your library, your statuses, stored locally.' },
  { phase: 2, name: 'Enrichment', line: 'Covers, synopsis, prices, moods.' },
  { phase: 3, name: 'Recommendations & Budget', line: 'What to buy, and whether to wait.' },
  { phase: 4, name: 'Social & Mobile', line: 'Sync, community, the phone apps.' }
] as const;
