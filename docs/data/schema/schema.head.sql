-- ============================================================================
-- Zerado · schema.head.sql — GENERATED. DO NOT EDIT. NEVER EXECUTED.
-- ----------------------------------------------------------------------------
-- The head-of-ladder schema, produced by applying every migration in
-- migrations/ in order and dumping the result:
--
--   for m in migrations/*.sql; do sqlite3 head.db < "$m"; done
--   sqlite3 head.db '.schema --indent'
--
-- It exists for TWO reasons and neither of them is execution:
--
--   1. REVIEW. A reviewer should be able to read the whole schema on one page
--      without replaying a ladder in their head.
--   2. THE FINGERPRINT TEST. A fresh install runs the SAME ladder an upgrade
--      runs — there is no separate "create everything" path, because the two
--      always drift. This file is what the test compares against, and a diff
--      here on a PR that did not intend one is the drift being caught.
--
-- Regenerate whenever a migration is added. See 03-migrations.md §8.
-- ============================================================================

CREATE TABLE schema_migration(
  version INTEGER PRIMARY KEY,
  applied_at TEXT NOT NULL,
  written_by TEXT NOT NULL,
  CHECK(version > 0),
  CHECK(written_by <> ''),
  CHECK(applied_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;
CREATE TABLE item(
  id INTEGER PRIMARY KEY,
  -- Phase 4 merge HINT. Indexed, never unique, never authoritative.
  -- uuidv5(NAMESPACE_ZERADO, item_type||'|'||uid_norm(title)||'|'||uid_norm(platform))
  -- NAMESPACE_ZERADO = uuidv5(NAMESPACE_DNS, 'zerado.app')
  -- = 56640350-c577-5522-8dd0-30e65323adf8
  item_uid TEXT NOT NULL,
  -- The door. One column, one value, CHECKed to it.
  item_type TEXT NOT NULL DEFAULT 'game',
  provider_id TEXT NOT NULL, -- 'steam' | 'physical' | …
  provider_ref TEXT NOT NULL, -- Steam appid, or a UUID for physical
  acquisition TEXT NOT NULL, -- 'digital' | 'physical'
  title TEXT NOT NULL, -- exactly as the source gave it
  sort_title TEXT NOT NULL, -- derived: folded, articles stripped — ORDER BY
  search_title TEXT NOT NULL, -- derived: folded, articles KEPT — LIKE
  platform TEXT NOT NULL,
  -- NULL = this source does not report playtime at all. 0 = it reported zero.
  -- The provider seam already carries this distinction as `Item.Playtime *int`
  --(06-data-seams §2.1); a NOT NULL DEFAULT 0 column would flatten it at the
  -- storage boundary and Z-04 could no longer render `—` for a cartridge and
  -- `0h` for an unplayed Steam title from the row alone.
  playtime_minutes INTEGER NULL,
  last_played_at TEXT NULL, -- NULL = not reported, NOT never played
  owned_since TEXT NULL, -- reduced precision allowed: YYYY | YYYY-MM | YYYY-MM-DD
  steam_appid TEXT NULL,
  achievements_total INTEGER NULL,
  achievements_unlocked INTEGER NULL,
  -- Tombstone. Set on the first COMPLETE, SUCCESSFUL sync that omits the row;
  -- cleared silently when it returns. NEVER a reason to delete(06 §2.4).
  absent_since TEXT NULL,
  -- Phase 4 merge target, so joining two rows never rewrites a primary key.
  merged_into INTEGER NULL REFERENCES item(id) ON DELETE SET NULL,
  -- PER-ROW sync age. Z-05 block 3 renders `LAST SYNCED` for THIS game,
  and a
  -- PARTIAL run leaves some rows fresh and others untouched — so a single
  -- provider-level timestamp would be a lie about every row the run never
  -- reached. Stamped on every upsert that wrote this row.
  last_synced_at TEXT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  CHECK(item_type = 'game'),
  CHECK(acquisition IN('digital','physical')),
  CHECK(title <> '' AND sort_title <> '' AND search_title <> '' AND platform <> ''),
  CHECK(provider_id <> '' AND provider_ref <> ''),
  CHECK(length(item_uid) = 36
AND item_uid GLOB '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*'),
CHECK(playtime_minutes IS NULL OR playtime_minutes >= 0),
CHECK(achievements_total IS NULL OR achievements_total >= 0),
CHECK(achievements_unlocked IS NULL
OR(achievements_total IS NOT NULL
AND achievements_unlocked BETWEEN 0 AND achievements_total)),
CHECK(merged_into IS NULL OR merged_into <> id),
-- Instants are RFC 3339, UTC, seconds precision, EXACTLY 20 characters, so
  -- that lexical order is chronological order and no index needs a collation.
  CHECK(last_played_at IS NULL OR last_played_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(absent_since IS NULL OR absent_since GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(last_synced_at IS NULL OR last_synced_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(created_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(updated_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
-- An acquisition date has no clock. A player who types "1998" is telling the
  -- truth at year precision; storing 1998-01-01T00:00:00Z would fabricate the
  -- other five fields. Z-08: "Optional. A year is enough."
  CHECK(owned_since IS NULL
OR owned_since GLOB '[0-9][0-9][0-9][0-9]'
OR owned_since GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]'
OR owned_since GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]')
) STRICT;
CREATE TABLE item_user(
  item_id INTEGER PRIMARY KEY REFERENCES item(id) ON DELETE CASCADE,
  -- NULL = the player has never expressed an opinion; derive on read.
  status_manual TEXT NULL,
  -- The timestamp that decides a Phase 4 conflict. NULL only while the player
  -- has NEVER expressed an opinion. CLEARING an override is itself a change
  -- and KEEPS its timestamp — without that, a clear can never win a
  -- last-write-wins merge against an older set on another device, and the
  -- override the player deleted would silently come back.
  status_changed_at TEXT NULL,
  -- Present from 0001 because they are user truth and this is the table for
  -- it. NO Phase 1 screen writes either one(Z-08 §3.1 rejects `notes`; no
-- screen renders a rating). `rating` deliberately carries NO range CHECK:
  -- the scale is the Phase 2 rating screen's to choose, and constraining it
-- now would pick a fight this ticket cannot win.
  rating INTEGER NULL,
  notes TEXT NULL,
  CHECK(status_manual IS NULL
OR status_manual IN('not_started','in_progress','zerado','abandoned')),
  CHECK(status_manual IS NULL OR status_changed_at IS NOT NULL),
  CHECK(status_changed_at IS NULL
OR status_changed_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(notes IS NULL OR notes <> '')
) STRICT;
CREATE TRIGGER item_user_exists AFTER INSERT ON item
BEGIN
    INSERT INTO item_user(item_id) VALUES (new.id);
END;
CREATE TABLE provider_connection(
  provider_id TEXT PRIMARY KEY,
  account_ref TEXT NULL,
  connected_at TEXT NOT NULL,
  last_sync_at TEXT NULL,
  last_sync_status TEXT NULL,
  CHECK(provider_id <> ''),
  CHECK(last_sync_status IS NULL
OR last_sync_status IN('ok','partial','failed','cancelled')),
  CHECK(connected_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(last_sync_at IS NULL OR last_sync_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;
CREATE TABLE sync_run(
  id INTEGER PRIMARY KEY,
  provider_id TEXT NOT NULL REFERENCES provider_connection(provider_id) ON DELETE CASCADE,
  started_at TEXT NOT NULL,
  finished_at TEXT NULL,
  status TEXT NOT NULL DEFAULT 'running',
  items_seen INTEGER NOT NULL DEFAULT 0,
  items_new INTEGER NOT NULL DEFAULT 0,
  items_changed INTEGER NOT NULL DEFAULT 0,
  error_kind TEXT NULL,
  error_detail TEXT NULL,
  CHECK(status IN('running','ok','partial','failed','cancelled')),
  CHECK((status = 'running') =(finished_at IS NULL)),
  CHECK(items_seen >= 0 AND items_new >= 0 AND items_changed >= 0),
  -- Z-03 renders "12 new. 4 changed. 131 unchanged." and its acceptance
  -- criterion is that each set of counts sums. Made a schema fact.
  CHECK(items_seen >= items_new + items_changed),
  CHECK(error_kind IS NULL
OR error_kind IN('no_route','timeout','unauthorized','empty','server','other')),
  CHECK(error_kind IS NOT NULL OR error_detail IS NULL),
  CHECK(status IN('ok','running','cancelled') OR error_kind IS NOT NULL),
  CHECK(started_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(finished_at IS NULL OR finished_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(finished_at IS NULL OR finished_at >= started_at)
) STRICT;
CREATE TRIGGER sync_run_retention AFTER INSERT ON sync_run
BEGIN
    DELETE FROM sync_run
     WHERE provider_id = new.provider_id
       AND id <= (SELECT id FROM sync_run
                   WHERE provider_id = new.provider_id
                   ORDER BY id DESC LIMIT 1 OFFSET 50);
END;
CREATE TABLE setting(
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  CHECK(key <> '')
) STRICT;
CREATE VIEW item_view AS
SELECT
    i.id, i.item_uid, i.item_type, i.provider_id, i.provider_ref, i.acquisition,
    i.title, i.sort_title, i.search_title, i.platform,
    i.playtime_minutes, i.last_played_at, i.owned_since,
    i.steam_appid, i.achievements_total, i.achievements_unlocked,
    i.absent_since, i.merged_into, i.last_synced_at, i.created_at, i.updated_at,
    u.status_manual, u.status_changed_at, u.rating, u.notes,
    COALESCE(
        u.status_manual,
        CASE WHEN i.playtime_minutes > 0 THEN 'in_progress' ELSE 'not_started' END
    )                             AS effective_status,
    (u.status_manual IS NOT NULL) AS status_is_manual,
    (i.playtime_minutes IS NULL)  AS playtime_untracked
FROM item i
JOIN item_user u ON u.item_id = i.id
/* item_view(id,item_uid,item_type,provider_id,provider_ref,acquisition,title,sort_title,search_title,platform,playtime_minutes,last_played_at,owned_since,steam_appid,achievements_total,achievements_unlocked,absent_since,merged_into,last_synced_at,created_at,updated_at,status_manual,status_changed_at,rating,notes,effective_status,status_is_manual,playtime_untracked) */;
CREATE UNIQUE INDEX item_provider_identity ON item(provider_id, provider_ref);
CREATE INDEX item_shelf_order ON item(
  sort_title
)
WHERE absent_since IS NULL 
    AND merged_into IS NULL;
CREATE INDEX item_absent ON item(absent_since) WHERE absent_since IS NOT NULL;
CREATE INDEX item_uid_lookup ON item(item_uid);
CREATE INDEX sync_run_recent ON sync_run(provider_id, started_at DESC);
CREATE TABLE metadata(
  item_id INTEGER PRIMARY KEY REFERENCES item(id) ON DELETE CASCADE,
  sinopse TEXT NULL,
  -- A LOCAL CACHE PATH. Never a remote URL — nothing in Zerado renders from
  -- the network(06-data-seams §1), and the CHECK below is that rule made
  -- structural rather than remembered.
  cover_ref TEXT NULL,
  cover_fetched_at TEXT NULL,
  released_at TEXT NULL, -- reduced precision allowed, as owned_since
  genres TEXT NULL, -- JSON array
  source_provider TEXT NOT NULL,
  -- The credit line is a property OF THE SOURCE(D1), so swapping the source
  -- swaps the credit. Stored with the row it credits, never in config.
  attribution TEXT NOT NULL,
  -- NOT NULL, because 07-offline-contract §4 forbids rendering a
  -- network-derived value without its age, and the schema is a better place
  -- to enforce that than the renderer's memory.
fetched_at TEXT NOT NULL,
CHECK(source_provider <> '' AND attribution <> ''),
CHECK(sinopse IS NULL OR sinopse <> ''),
CHECK(cover_ref IS NULL
OR(cover_ref NOT LIKE 'http://%' AND cover_ref NOT LIKE 'https://%')),
CHECK((cover_ref IS NULL) =(cover_fetched_at IS NULL)),
CHECK(genres IS NULL OR json_valid(genres)),
CHECK(fetched_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(cover_fetched_at IS NULL OR cover_fetched_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(released_at IS NULL
OR released_at GLOB '[0-9][0-9][0-9][0-9]'
OR released_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]'
OR released_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]')
) STRICT;
CREATE TABLE mood_tag(
  id INTEGER PRIMARY KEY,
  key TEXT NOT NULL,
  label TEXT NOT NULL,
  CHECK(key <> '' AND label <> ''),
  CHECK(key = lower(key) AND key NOT GLOB '* *')
) STRICT;
CREATE UNIQUE INDEX mood_tag_key ON mood_tag(key);
CREATE TABLE item_mood(
  item_id INTEGER NOT NULL REFERENCES item(id) ON DELETE CASCADE,
  mood_id INTEGER NOT NULL REFERENCES mood_tag(id) ON DELETE CASCADE,
  source TEXT NOT NULL,
  confidence REAL NULL,
  assigned_at TEXT NOT NULL,
  PRIMARY KEY(item_id, mood_id),
  CHECK(source IN('user','inferred')),
  CHECK((source = 'inferred') =(confidence IS NOT NULL)),
  CHECK(confidence IS NULL OR(confidence > 0.0 AND confidence <= 1.0)),
  CHECK(assigned_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;
CREATE INDEX item_mood_by_mood ON item_mood(mood_id, item_id);
CREATE INDEX item_mood_inferred ON item_mood(
  item_id
) WHERE source = 'inferred';
CREATE TABLE price_quote(
  item_id INTEGER NOT NULL REFERENCES item(id) ON DELETE CASCADE,
  shop TEXT NOT NULL,
  currency TEXT NOT NULL, -- ISO 4217, uppercase
  current_cents INTEGER NOT NULL, -- minor units. Never a float: money.
  low_cents INTEGER NULL, -- the provider's all-time low
low_at TEXT NULL, -- "a low with no date is not information"
  local_low_cents INTEGER NULL, -- watermark, only when the provider gives none
  local_low_at TEXT NULL,
  url TEXT NOT NULL, -- a plain shop link
  source_provider TEXT NOT NULL,
  attribution TEXT NOT NULL,
  -- NOT NULL and rendered. "A price without its age is the product giving
-- financial advice from memory."
observed_at TEXT NOT NULL,
PRIMARY KEY(item_id, shop, currency),
CHECK(shop <> '' AND url <> '' AND source_provider <> '' AND attribution <> ''),
  CHECK(currency GLOB '[A-Z][A-Z][A-Z]'),
CHECK(current_cents >= 0),
CHECK(low_cents IS NULL OR low_cents >= 0),
CHECK((low_cents IS NULL) =(low_at IS NULL)),
CHECK((local_low_cents IS NULL) =(local_low_at IS NULL)),
-- A provider low and a local watermark are alternatives, never both.
  CHECK(low_cents IS NULL OR local_low_cents IS NULL),
  CHECK(observed_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(low_at IS NULL OR low_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
CHECK(local_low_at IS NULL OR local_low_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;
