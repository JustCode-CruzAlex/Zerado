-- ============================================================================
-- Zerado · benchmark fixture generator
-- ----------------------------------------------------------------------------
-- Generates a synthetic library of :N titles against migration 0001, with the
-- state distribution of Z-04's own mockup (247 games · 198 not started ·
-- 12 in progress · 6 zerado · 31 abandoned) held roughly constant as N scales,
-- one physical (playtime-untracked) row in twenty, and a handful of absent rows.
--
-- Substitute :N before running, e.g.
--   sed 's/:N/400/' bench-generate.sql | sqlite3 library.db
--
-- This exists to make the indexing plan EVIDENCE rather than assertion. It is
-- not a test fixture — see seed-minimal.sql and seed-edge-cases.sql for those.
-- item_uid uses randomblob here on purpose: a realistic (unclustered) key
-- distribution is what an index measurement needs. The test fixtures use fixed
-- uids instead, because a test needs determinism.
-- ============================================================================
PRAGMA foreign_keys = ON;

INSERT INTO provider_connection(provider_id, account_ref, connected_at, last_sync_at, last_sync_status)
VALUES ('steam','76561198000000000','2026-06-01T09:00:00Z','2026-08-25T04:00:00Z','ok');

WITH RECURSIVE
n(i) AS (SELECT 1 UNION ALL SELECT i+1 FROM n WHERE i < :N),
w(k,word) AS (VALUES
 (0,'Hollow'),(1,'Celeste'),(2,'Bastion'),(3,'Outer'),(4,'Return'),(5,'Talos'),(6,'Messenger'),
 (7,'Forgotten'),(8,'Campfire'),(9,'Vampire'),(10,'Dark'),(11,'Nier'),(12,'Okami'),(13,'Zelda'),
 (14,'Portal'),(15,'Braid'),(16,'Limbo'),(17,'Inside'),(18,'Journey'),(19,'Gris'),(20,'Sable'),
 (21,'Tunic'),(22,'Hades'),(23,'Dredge'),(24,'Cocoon'),(25,'Chants')),
v(k,word) AS (VALUES
 (0,'Knight'),(1,'Souls'),(2,'Worlds'),(3,'Dinn'),(4,'Principle'),(5,'City'),(6,'Survivors'),
 (7,'Automata'),(8,'Breath'),(9,'Wild'),(10,'Legacy'),(11,'Echo'),(12,'Rift'),(13,'Ashes')),
t(ii,title) AS (
  SELECT i,
         (SELECT word FROM w WHERE k = i % 26) || ' ' ||
         (SELECT word FROM v WHERE k = i % 14) || ' ' || i
  FROM n)
INSERT INTO item(
  item_uid, item_type, provider_id, provider_ref, acquisition,
  title, sort_title, search_title, platform,
  playtime_minutes, last_played_at, owned_since, steam_appid,
  achievements_total, achievements_unlocked, absent_since, merged_into,
  last_synced_at, created_at, updated_at)
SELECT
  lower(substr(hex(randomblob(4)),1,8) || '-' || substr(hex(randomblob(2)),1,4) ||
        '-5' || substr(hex(randomblob(2)),1,3) || '-8' || substr(hex(randomblob(2)),1,3) ||
        '-' || substr(hex(randomblob(6)),1,12)),
  'game',
  CASE WHEN i % 20 = 0 THEN 'physical' ELSE 'steam' END,
  CASE WHEN i % 20 = 0 THEN 'phy-' || i ELSE '' || (220000 + i*7) END,
  CASE WHEN i % 20 = 0 THEN 'physical' ELSE 'digital' END,
  t.title, lower(t.title), lower(t.title),
  CASE WHEN i % 20 = 0 THEN 'Nintendo Switch' ELSE 'PC' END,
  -- physical: NULL (this source cannot know). Most of the rest: a real 0 —
  -- a backlog is mostly unplayed, which is the product's whole premise.
  CASE WHEN i % 20 = 0 THEN NULL WHEN i % 5 <> 1 THEN 0 ELSE (i * 37) % 4000 END,
  CASE WHEN i % 20 = 0 OR i % 5 <> 1 THEN NULL
       ELSE '2026-0' || (1 + (i % 8)) || '-1' || (i % 10) || 'T20:' || printf('%02d', i % 60) || ':00Z' END,
  CASE WHEN i % 20 = 0 THEN '' || (1995 + (i % 30)) ELSE NULL END,
  CASE WHEN i % 20 = 0 THEN NULL ELSE '' || (220000 + i*7) END,
  CASE WHEN i % 7 = 0 THEN 30 + (i % 20) ELSE NULL END,
  CASE WHEN i % 7 = 0 THEN (i % 20) ELSE NULL END,
  CASE WHEN i % 137 = 0 THEN '2026-08-14T04:00:00Z' ELSE NULL END,
  NULL,
  '2026-08-25T04:00:00Z', '2026-06-01T09:00:00Z', '2026-08-25T04:00:00Z'
FROM n JOIN t ON t.ii = n.i;

-- Manual overrides. zerado and abandoned are NEVER derived (05-state-machine §4),
-- so every one of them is a row the player touched — and every one of them
-- therefore carries a status_changed_at.
UPDATE item_user SET status_manual = 'zerado',
                     status_changed_at = '2026-07-0' || (1 + (item_id % 9)) || 'T18:00:00Z'
 WHERE item_id % 40 = 3;
UPDATE item_user SET status_manual = 'abandoned',
                     status_changed_at = '2026-07-1' || (item_id % 9) || 'T18:00:00Z'
 WHERE item_id % 8 = 5;
UPDATE item_user SET status_manual = 'not_started',
                     status_changed_at = '2026-07-20T18:00:00Z'
 WHERE item_id % 97 = 11;

WITH RECURSIVE r(i) AS (SELECT 1 UNION ALL SELECT i+1 FROM r WHERE i < 60)
INSERT INTO sync_run(provider_id, started_at, finished_at, status, items_seen, items_new, items_changed)
SELECT 'steam',
       '2026-08-2' || (i % 5) || 'T0' || (i % 9) || ':00:00Z',
       '2026-08-2' || (i % 5) || 'T0' || (i % 9) || ':02:00Z',
       'ok', 247, 12, 4
FROM r;
