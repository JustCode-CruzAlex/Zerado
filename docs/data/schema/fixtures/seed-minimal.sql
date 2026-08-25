-- ============================================================================
-- Zerado · seed-minimal.sql — the smallest library that exercises every rule
-- ----------------------------------------------------------------------------
-- Nine rows. Deterministic: fixed item_uids, fixed timestamps, no randomblob.
-- Apply on top of migration 0001. A test that needs Phase 2 rows applies 0002
-- and then seed-minimal-phase2.sql-shaped inserts of its own.
--
-- What each row is FOR is stated on the row, because a fixture whose purpose is
-- not written down is a fixture the next person is afraid to change.
-- ============================================================================
PRAGMA foreign_keys = ON;

INSERT INTO provider_connection(provider_id, account_ref, connected_at, last_sync_at, last_sync_status)
VALUES ('steam', '76561198000000000', '2026-06-01T09:00:00Z', '2026-08-22T04:00:00Z', 'ok');

-- `physical` has NO provider_connection row: it needs no credentials and never
-- syncs. That is not an omission — it is why item.provider_id has no FK.

INSERT INTO item(item_uid, item_type, provider_id, provider_ref, acquisition,
                 title, sort_title, search_title, platform,
                 playtime_minutes, last_played_at, owned_since, steam_appid,
                 achievements_total, achievements_unlocked,
                 absent_since, merged_into, last_synced_at, created_at, updated_at) VALUES
-- 1 · derived NOT STARTED: Steam reports playtime, and it is a real zero.
('11111111-1111-5111-8111-111111111111','game','steam','367520','digital',
 'Hollow Knight','hollow knight','hollow knight','PC',
 0, NULL, NULL, '367520', 63, 0, NULL, NULL, '2026-08-22T04:00:00Z','2026-06-01T09:00:00Z','2026-08-22T04:00:00Z'),
-- 2 · derived IN PROGRESS: the ONLY automatic transition in the product.
('22222222-2222-5222-8222-222222222222','game','steam','374320','digital',
 'Dark Souls III','dark souls iii','dark souls iii','PC',
 3780, '2026-08-20T21:14:00Z', NULL, '374320', 43, 21, NULL, NULL, '2026-08-22T04:00:00Z','2026-06-01T09:00:00Z','2026-08-22T04:00:00Z'),
-- 3 · manual ZERADO on a game still being played: the sync must NOT undo it.
('33333333-3333-5333-8333-333333333333','game','steam','1150690','digital',
 'Return of the Obra Dinn','return of the obra dinn','return of the obra dinn','PC',
 540, '2026-08-19T22:03:00Z', NULL, '1150690', NULL, NULL, NULL, NULL, '2026-08-22T04:00:00Z','2026-06-01T09:00:00Z','2026-08-22T04:00:00Z'),
-- 4 · manual ABANDONED with playtime: manual wins, permanently.
('44444444-4444-5444-8444-444444444444','game','steam','524220','digital',
 'Nier: Automata','nier automata','nier: automata','PC',
 1680, '2026-07-02T19:40:00Z', NULL, '524220', 48, 12, NULL, NULL, '2026-08-22T04:00:00Z','2026-06-01T09:00:00Z','2026-08-22T04:00:00Z'),
-- 5 · a cartridge. playtime_minutes IS NULL: this source CANNOT know.
--     Z-04 renders `—` here and `0h` on row 1, from the row alone.
('55555555-5555-5555-8555-555555555555','game','physical','7f3a1c22-0000-4000-8000-000000000001','physical',
 'The Legend of Zelda: Breath of the Wild','legend of zelda breath of the wild','the legend of zelda: breath of the wild','Nintendo Switch',
 NULL, NULL, '2017', NULL, NULL, NULL, NULL, NULL, NULL, '2026-06-02T10:00:00Z','2026-06-02T10:00:00Z'),
-- 6 · a cartridge whose owned_since is a YEAR. Z-08: "A year is enough."
('66666666-6666-5666-8666-666666666666','game','physical','7f3a1c22-0000-4000-8000-000000000002','physical',
 'Chrono Trigger','chrono trigger','chrono trigger','SNES',
 NULL, NULL, '1995', NULL, NULL, NULL, NULL, NULL, NULL, '2026-06-02T10:05:00Z','2026-06-02T10:05:00Z'),
-- 7 · ABSENT: tombstoned, never deleted. Carries the player's ZERADO, which is
--     exactly the row they would be angriest to lose (06-data-seams §2.4).
('77777777-7777-5777-8777-777777777777','game','steam','108710','digital',
 'Alan Wake','alan wake','alan wake','PC',
 1320, '2026-05-11T23:12:00Z', NULL, '108710', NULL, NULL,
 '2026-08-14T04:00:00Z', NULL, '2026-08-14T04:00:00Z','2026-06-01T09:00:00Z','2026-08-14T04:00:00Z'),
-- 8 · a title that is not ASCII. `Pokémon` must be findable by typing `pokemon`.
('88888888-8888-5888-8888-888888888888','game','physical','7f3a1c22-0000-4000-8000-000000000003','physical',
 'Pokémon Red','pokemon red','pokemon red','Game Boy',
 NULL, NULL, '1998', NULL, NULL, NULL, NULL, NULL, NULL, '2026-06-02T10:10:00Z','2026-06-02T10:10:00Z'),
-- 9 · MERGED into row 5. Out of the default view; its keys were never rewritten.
('99999999-9999-5999-8999-999999999999','game','steam','9990001','digital',
 'Zelda: Breath of the Wild','zelda breath of the wild','zelda: breath of the wild','PC',
 60, '2026-04-01T20:00:00Z', NULL, '9990001', NULL, NULL, NULL, NULL, '2026-08-22T04:00:00Z','2026-06-01T09:00:00Z','2026-08-22T04:00:00Z');

UPDATE item SET merged_into = 5 WHERE provider_ref = '9990001';

-- The player's own work. item_user rows already exist (materialised by trigger).
-- Every manual status carries a status_changed_at — the schema refuses otherwise.
UPDATE item_user SET status_manual='zerado',      status_changed_at='2026-07-14T22:10:00Z' WHERE item_id = 3;
UPDATE item_user SET status_manual='abandoned',   status_changed_at='2026-07-02T20:01:00Z' WHERE item_id = 4;
UPDATE item_user SET status_manual='zerado',      status_changed_at='2026-05-12T00:04:00Z' WHERE item_id = 7;
UPDATE item_user SET status_manual='not_started', status_changed_at='2026-08-01T11:00:00Z' WHERE item_id = 6;
-- Row 2: the player CLEARED an override. status_manual is NULL again but the
-- timestamp REMAINS, so a Phase 4 merge can tell "I cleared this on the 21st"
-- from "I never had an opinion".
UPDATE item_user SET status_manual='in_progress', status_changed_at='2026-08-21T18:00:00Z' WHERE item_id = 2;
UPDATE item_user SET status_manual=NULL                                                    WHERE item_id = 2;

INSERT INTO sync_run(provider_id, started_at, finished_at, status, items_seen, items_new, items_changed, error_kind, error_detail) VALUES
('steam','2026-08-22T04:00:00Z','2026-08-22T04:00:41Z','ok',       7, 0, 2, NULL, NULL),
('steam','2026-08-23T09:12:00Z','2026-08-23T09:12:09Z','partial',  3, 0, 0, 'timeout','read tcp: i/o timeout after 3 pages'),
('steam','2026-08-24T08:00:00Z','2026-08-24T08:00:02Z','failed',   0, 0, 0, 'no_route','dial tcp: no route to host');

INSERT INTO setting(key, value) VALUES
('theme','zerado-dark'),
('audio.enabled','false'),
('normaliser_version','1');

-- Expected shape, so a test can assert it without re-deriving it:
--   shown rows (default view)     7   (9 total - 1 absent - 1 merged)
--   NOT STARTED  4   (1 derived-zero, 5 and 8 derived-untracked, 6 manual)
--   IN PROGRESS  1   (2 — derived again after the clear)
--   ZERADO       1   (3;  7 is absent and out of the default view)
--   ABANDONED    1   (4)
--                    4 + 1 + 1 + 1 = 7, and the four counts sum to the number
--                    shown — 05-state-machine §7 rule 1, asserted by a test.
--   absent       1   ·  merged 1  ·  whole file 9
