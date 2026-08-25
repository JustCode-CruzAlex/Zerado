/**
 * The four ratified roadmap phases, verbatim from content/landing-copy.md §12.
 *
 * `status` is the single source of truth for what the roadmap renders. Phase 1
 * moved to 'in-progress' on 2026-08-25, when Phase 1 work started; phases 2-4
 * remain 'planned'. NOTHING is ever marked done here until it is done — that is
 * the one line on this page that would actually cost credibility.
 */
export type PhaseStatus = 'planned' | 'in-progress';

export const PHASES: readonly {
  phase: number;
  name: string;
  line: string;
  status: PhaseStatus;
}[] = [
  { phase: 1, name: 'CLI/TUI MVP', line: 'Your library, your statuses, stored locally.', status: 'in-progress' },
  { phase: 2, name: 'Enrichment', line: 'Covers, synopsis, prices, moods.', status: 'planned' },
  { phase: 3, name: 'Recommendations & Budget', line: 'What to buy, and whether to wait.', status: 'planned' },
  { phase: 4, name: 'Social & Mobile', line: 'Sync, community, the phone apps.', status: 'planned' }
] as const;
