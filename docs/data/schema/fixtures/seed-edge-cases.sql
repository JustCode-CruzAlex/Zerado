-- ============================================================================
-- Zerado · seed-edge-cases.sql — the rows that break naive code
-- ----------------------------------------------------------------------------
-- Apply on top of migration 0001 (standalone; does NOT require seed-minimal).
-- Each row is here because something plausible gets it wrong. The comment says
-- what.
-- ============================================================================
PRAGMA foreign_keys = ON;

INSERT INTO item(item_uid, item_type, provider_id, provider_ref, acquisition,
                 title, sort_title, search_title, platform,
                 playtime_minutes, last_played_at, owned_since,
                 created_at, updated_at) VALUES

-- Diacritics. A byte-order sort puts `Ōkami` after `Zelda`; a folded sort_title
-- puts it under O, where a person looks for it. (ADR-0001 D9)
('a0000000-0000-5000-8000-000000000001','game','physical','e-1','physical',
 'Ōkami HD','okami hd','okami hd','PC', NULL, NULL, NULL,
 '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- A leading article. sort_title strips it (sorts under W); search_title KEEPS
-- it, so typing `the` finds it. Two columns, two jobs, one normaliser.
('a0000000-0000-5000-8000-000000000002','game','steam','e-2','digital',
 'The Witness','witness','the witness','PC', 0, NULL, NULL,
 '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- Leading punctuation. `[PROTOTYPE]` sorts under P, not before every letter.
('a0000000-0000-5000-8000-000000000003','game','steam','e-3','digital',
 '[PROTOTYPE]','prototype','[prototype]','PC', 0, NULL, NULL,
 '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- A title that is ONLY punctuation: sort_title must not end up empty (the
-- schema refuses an empty one), so the normaliser falls back to the folded
-- title unchanged.
('a0000000-0000-5000-8000-000000000004','game','steam','e-4','digital',
 '!!!','!!!','!!!','PC', 0, NULL, NULL,
 '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- Full-width and combining characters: terminal CELL WIDTH is not string
-- length. `Ｚｅｒａｄｏ` is 6 runes and 12 cells. (ADR-0001 D9)
('a0000000-0000-5000-8000-000000000005','game','physical','e-5','physical',
 'ＺＥＲＡＤＯ','zerado','ｚｅｒａｄｏ','PC', NULL, NULL, '2001',
 '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- A very long title. Z-04 truncates for display; the STORED title never is.
('a0000000-0000-5000-8000-000000000006','game','steam','e-6','digital',
 'The Legend of the Extremely Long Subtitle That Will Never Fit In A Terminal Column Budget Not Even At ExtraWide',
 'legend of the extremely long subtitle that will never fit in a terminal column budget not even at extrawide',
 'the legend of the extremely long subtitle that will never fit in a terminal column budget not even at extrawide',
 'PC', 0, NULL, NULL, '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- playtime 0 with a last_played_at: the player launched it and quit in under a
-- minute. `0h` and a real last-played date are not contradictory.
('a0000000-0000-5000-8000-000000000007','game','steam','e-7','digital',
 'Vampire Survivors','vampire survivors','vampire survivors','PC',
 0, '2026-08-24T23:59:00Z', NULL, '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- playtime NULL AND last_played NULL: the source reports neither. This row is
-- NOT "never played" — it is "nobody can say". Z-05 renders `not tracked`.
('a0000000-0000-5000-8000-000000000008','game','physical','e-8','physical',
 'Chrono Cross','chrono cross','chrono cross','PlayStation',
 NULL, NULL, '1999', '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- Two editions of one game, same item_uid: the merge HINT collides, and it is
-- allowed to. `item_uid` is indexed and NOT unique, and Phase 4 shows this pair
-- to the player rather than guessing (ADR-0001 D4).
('a0000000-0000-5000-8000-000000000009','game','steam','e-9a','digital',
 'Dark Souls','dark souls','dark souls','PC', 600, NULL, NULL,
 '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),
('a0000000-0000-5000-8000-000000000009','game','physical','e-9b','physical',
 'Dark Souls','dark souls','dark souls','PC', NULL, NULL, '2011',
 '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z'),

-- A row that is absent AND carries NOTHING of the player's. This is the only
-- class that may be OFFERED for bulk removal — offered, never done.
('a0000000-0000-5000-8000-00000000000a','game','steam','e-10','digital',
 'Some Delisted Thing','some delisted thing','some delisted thing','PC',
 0, NULL, NULL, '2026-08-25T04:00:00Z','2026-08-25T04:00:00Z');

UPDATE item SET absent_since = '2026-08-14T04:00:00Z' WHERE provider_ref = 'e-10';

-- A note containing a newline, a quote and a semicolon: it round-trips.
UPDATE item_user
   SET notes = 'Finished on the Switch, 2019.' || char(10) || 'Then again on PC; "same game", different feeling.'
 WHERE item_id = (SELECT id FROM item WHERE provider_ref = 'e-2');
