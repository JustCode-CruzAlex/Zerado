-- ============================================================================
-- Zerado · migration 0001 · Phase 1 core
-- ----------------------------------------------------------------------------
-- ticket   #5   ([BLUEPRINT][DATA] Zerado ERD and physical schema)
-- consumes ADR-0001 D2 (persistence) · D4 (Phase 4 sync boundary) · D5 (the door)
--          docs/blueprint/09-erd.md   — the conceptual model this makes physical
--          docs/blueprint/05-state-machine.md · 06-data-seams.md · 07-offline-contract.md
--
-- This file is EXECUTABLE TRUTH. `docs/data/schema/schema.head.sql` is generated
-- from it and is never executed. A fresh install runs this ladder; it never runs
-- a separate "create everything" script, because the two always drift.
--
-- Runs inside ONE transaction, opened by the migrator, together with the
-- INSERT INTO schema_migration that records it. There is no half-applied 0001.
--
-- Requires SQLite >= 3.37.0 (STRICT). Verified at source 2026-08-25:
-- modernc.org/sqlite v1.57.0 bundles SQLite 3.53.3.
-- ============================================================================


-- ----------------------------------------------------------------------------
-- schema_migration — the ladder's own bookkeeping
-- ----------------------------------------------------------------------------
-- `written_by` is why this is not just (version, applied_at). Z-11 promises to
-- name BOTH versions — the one that wrote the file and the one that is running.
-- A schema number alone cannot do that, and the binary that wrote a row is
-- exactly what is no longer around to ask. It cannot be backfilled later.
CREATE TABLE schema_migration (
    version     INTEGER PRIMARY KEY,
    applied_at  TEXT    NOT NULL,
    written_by  TEXT    NOT NULL,

    CHECK (version > 0),
    CHECK (written_by <> ''),
    CHECK (applied_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;


-- ----------------------------------------------------------------------------
-- item — identity, and what a PROVIDER says about it
-- ----------------------------------------------------------------------------
-- The core entity is `item`, not `games`, carrying an `item_type` CHECKed to
-- 'game'. That is the entire door-open affordance (ADR-0001 D5). Nothing else
-- here is generalised: playtime is playtime.
--
-- EVERY COLUMN IN THIS TABLE IS EITHER IDENTITY OR PROVIDER REPLICA.
-- Nothing the player typed lives here — that is `item_user`. The split is the
-- mechanism that makes "a metadata refresh can never destroy user data" a fact
-- about the SQL rather than a promise about the caller: the sync's UPSERT
-- cannot erase `status_manual`, because the statement cannot name it.
--
-- The one exception is a `physical` row, whose title/platform/owned_since ARE
-- the player's own typing. It needs no separate guard: `physical` implements
-- Provider and NOT Syncer (ADR-0001 D1), so no sync statement ever reaches it.
CREATE TABLE item (
    id                    INTEGER PRIMARY KEY,

    -- Phase 4 merge HINT. Indexed, never unique, never authoritative.
    -- uuidv5(NAMESPACE_ZERADO, item_type||'|'||uid_norm(title)||'|'||uid_norm(platform))
    -- NAMESPACE_ZERADO = uuidv5(NAMESPACE_DNS, 'zerado.app')
    --                  = 56640350-c577-5522-8dd0-30e65323adf8
    item_uid              TEXT    NOT NULL,

    -- The door. One column, one value, CHECKed to it.
    item_type             TEXT    NOT NULL DEFAULT 'game',

    provider_id           TEXT    NOT NULL,   -- 'steam' | 'physical' | …
    provider_ref          TEXT    NOT NULL,   -- Steam appid, or a UUID for physical
    acquisition           TEXT    NOT NULL,   -- 'digital' | 'physical'

    title                 TEXT    NOT NULL,   -- exactly as the source gave it
    sort_title            TEXT    NOT NULL,   -- derived: folded, articles stripped — ORDER BY
    search_title          TEXT    NOT NULL,   -- derived: folded, articles KEPT   — LIKE
    platform              TEXT    NOT NULL,

    -- NULL = this source does not report playtime at all. 0 = it reported zero.
    -- The provider seam already carries this distinction as `Item.Playtime *int`
    -- (06-data-seams §2.1); a NOT NULL DEFAULT 0 column would flatten it at the
    -- storage boundary and Z-04 could no longer render `—` for a cartridge and
    -- `0h` for an unplayed Steam title from the row alone.
    playtime_minutes      INTEGER NULL,
    last_played_at        TEXT    NULL,       -- NULL = not reported, NOT never played
    owned_since           TEXT    NULL,       -- reduced precision allowed: YYYY | YYYY-MM | YYYY-MM-DD

    steam_appid           TEXT    NULL,
    achievements_total    INTEGER NULL,
    achievements_unlocked INTEGER NULL,

    -- Tombstone. Set on the first COMPLETE, SUCCESSFUL sync that omits the row;
    -- cleared silently when it returns. NEVER a reason to delete (06 §2.4).
    absent_since          TEXT    NULL,

    -- Phase 4 merge target, so joining two rows never rewrites a primary key.
    merged_into           INTEGER NULL REFERENCES item(id) ON DELETE SET NULL,

    -- PER-ROW sync age. Z-05 block 3 renders `LAST SYNCED` for THIS game, and a
    -- PARTIAL run leaves some rows fresh and others untouched — so a single
    -- provider-level timestamp would be a lie about every row the run never
    -- reached. Stamped on every upsert that wrote this row.
    last_synced_at        TEXT    NULL,

    created_at            TEXT    NOT NULL,
    updated_at            TEXT    NOT NULL,

    CHECK (item_type = 'game'),
    CHECK (acquisition IN ('digital','physical')),
    CHECK (title <> '' AND sort_title <> '' AND search_title <> '' AND platform <> ''),
    CHECK (provider_id <> '' AND provider_ref <> ''),
    CHECK (length(item_uid) = 36
           AND item_uid GLOB '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*'),
    CHECK (playtime_minutes IS NULL OR playtime_minutes >= 0),
    CHECK (achievements_total IS NULL OR achievements_total >= 0),
    CHECK (achievements_unlocked IS NULL
           OR (achievements_total IS NOT NULL
               AND achievements_unlocked BETWEEN 0 AND achievements_total)),
    CHECK (merged_into IS NULL OR merged_into <> id),

    -- Instants are RFC 3339, UTC, seconds precision, EXACTLY 20 characters, so
    -- that lexical order is chronological order and no index needs a collation.
    CHECK (last_played_at IS NULL OR last_played_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (absent_since   IS NULL OR absent_since   GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (last_synced_at IS NULL OR last_synced_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (created_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (updated_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),

    -- An acquisition date has no clock. A player who types "1998" is telling the
    -- truth at year precision; storing 1998-01-01T00:00:00Z would fabricate the
    -- other five fields. Z-08: "Optional. A year is enough."
    CHECK (owned_since IS NULL
           OR owned_since GLOB '[0-9][0-9][0-9][0-9]'
           OR owned_since GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]'
           OR owned_since GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]')
) STRICT;


-- ----------------------------------------------------------------------------
-- item_user — what the PLAYER typed. Truth. Never re-derivable. Never lost.
-- ----------------------------------------------------------------------------
-- 1:1 with item, materialised by trigger so every read joins and no read can
-- lose a row. This table is the Phase 4 sync payload, the "export my work"
-- selection, and the thing a corrupt `item` table does not take with it.
CREATE TABLE item_user (
    item_id           INTEGER PRIMARY KEY REFERENCES item(id) ON DELETE CASCADE,

    -- NULL = the player has never expressed an opinion; derive on read.
    status_manual     TEXT NULL,

    -- The timestamp that decides a Phase 4 conflict. NULL only while the player
    -- has NEVER expressed an opinion. CLEARING an override is itself a change
    -- and KEEPS its timestamp — without that, a clear can never win a
    -- last-write-wins merge against an older set on another device, and the
    -- override the player deleted would silently come back.
    status_changed_at TEXT NULL,

    -- Present from 0001 because they are user truth and this is the table for
    -- it. NO Phase 1 screen writes either one (Z-08 §3.1 rejects `notes`; no
    -- screen renders a rating). `rating` deliberately carries NO range CHECK:
    -- the scale is the Phase 2 rating screen's to choose, and constraining it
    -- now would pick a fight this ticket cannot win.
    rating            INTEGER NULL,
    notes             TEXT    NULL,

    CHECK (status_manual IS NULL
           OR status_manual IN ('not_started','in_progress','zerado','abandoned')),
    CHECK (status_manual IS NULL OR status_changed_at IS NOT NULL),
    CHECK (status_changed_at IS NULL
           OR status_changed_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (notes IS NULL OR notes <> '')
) STRICT;

-- Pure structure, no clock: allowed. (A trigger that stamped a timestamp would
-- read SQLite's own `now` and bypass the injected clock the seams use for
-- tests — 06-data-seams §7. Every timestamp in this schema is written by the
-- Store from that clock, never by a trigger.)
CREATE TRIGGER item_user_exists AFTER INSERT ON item
BEGIN
    INSERT INTO item_user(item_id) VALUES (new.id);
END;


-- ----------------------------------------------------------------------------
-- provider_connection — a credential RELATIONSHIP, and no credential
-- ----------------------------------------------------------------------------
-- The keys are in the Vault (ADR-0001 D2·3). `account_ref` is an identifier —
-- a Steam ID, a GOG username — never a secret.
--
-- NOTE THE ABSENT FOREIGN KEY. `item.provider_id` does NOT reference this
-- table. 09-erd draws item→connection as `N:1` and it is a real association,
-- but it must not be a referential one: a connection is a relationship the
-- player can END, and the library outlives it. An FK would force the choice
-- between deleting a library on disconnect and forbidding disconnect while
-- items exist. Both are wrong, so there is no FK, and the association is
-- resolved in code against the compiled-in provider registry — which D1 makes
-- the authority anyway (`Capabilities`, never the database).
CREATE TABLE provider_connection (
    provider_id      TEXT PRIMARY KEY,
    account_ref      TEXT NULL,
    connected_at     TEXT NOT NULL,
    last_sync_at     TEXT NULL,
    last_sync_status TEXT NULL,

    CHECK (provider_id <> ''),
    CHECK (last_sync_status IS NULL
           OR last_sync_status IN ('ok','partial','failed','cancelled')),
    CHECK (connected_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (last_sync_at IS NULL OR last_sync_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;


-- ----------------------------------------------------------------------------
-- sync_run — the honest history
-- ----------------------------------------------------------------------------
-- `status = 'running'` exists because a run must be recordable BEFORE it
-- terminates for a kill to be observable — which is exactly the fact 09-erd
-- wants `finished_at IS NULL` to carry. Without it the column would have to be
-- NULL while running, and a NULL enum is a fifth state spelled as an absence.
--
-- `error` is split into a CHECKed classification plus free detail, so
-- "the classified failure, not a stack trace" is enforced rather than intended.
-- The classes are 07-offline-contract §5's classifier, one for one.
--
-- Deleting a connection CASCADES its runs: sync history is bookkeeping, and a
-- reconnect must not inherit a previous relationship's "last synced" age. No
-- item is touched — items are the player's library, not the connection's.
CREATE TABLE sync_run (
    id            INTEGER PRIMARY KEY,
    provider_id   TEXT    NOT NULL REFERENCES provider_connection(provider_id) ON DELETE CASCADE,
    started_at    TEXT    NOT NULL,
    finished_at   TEXT    NULL,
    status        TEXT    NOT NULL DEFAULT 'running',
    items_seen    INTEGER NOT NULL DEFAULT 0,
    items_new     INTEGER NOT NULL DEFAULT 0,
    items_changed INTEGER NOT NULL DEFAULT 0,
    error_kind    TEXT    NULL,
    error_detail  TEXT    NULL,

    CHECK (status IN ('running','ok','partial','failed','cancelled')),
    CHECK ((status = 'running') = (finished_at IS NULL)),
    CHECK (items_seen >= 0 AND items_new >= 0 AND items_changed >= 0),
    -- Z-03 renders "12 new. 4 changed. 131 unchanged." and its acceptance
    -- criterion is that each set of counts sums. Made a schema fact.
    CHECK (items_seen >= items_new + items_changed),
    CHECK (error_kind IS NULL
           OR error_kind IN ('no_route','timeout','unauthorized','empty','server','other')),
    CHECK (error_kind IS NOT NULL OR error_detail IS NULL),
    CHECK (status IN ('ok','running','cancelled') OR error_kind IS NOT NULL),
    CHECK (started_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (finished_at IS NULL OR finished_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (finished_at IS NULL OR finished_at >= started_at)
) STRICT;

-- Retention, made structural. Only the most recent run is ever read; an
-- unbounded log inside a file the player is told to back up is a slow leak.
-- 50 runs per provider is months of history and a few kilobytes.
-- Pure structure, no clock: allowed.
CREATE TRIGGER sync_run_retention AFTER INSERT ON sync_run
BEGIN
    DELETE FROM sync_run
     WHERE provider_id = new.provider_id
       AND id <= (SELECT id FROM sync_run
                   WHERE provider_id = new.provider_id
                   ORDER BY id DESC LIMIT 1 OFFSET 50);
END;


-- ----------------------------------------------------------------------------
-- setting — everything Z-09 writes. Local truth; never crosses to Phase 4.
-- ----------------------------------------------------------------------------
-- Device-local by nature (theme, audio, image capability, the normaliser
-- version). "Only what the player typed crosses" is about the LIBRARY; a
-- preference is about THIS machine.
CREATE TABLE setting (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    CHECK (key <> '')
) STRICT;


-- ----------------------------------------------------------------------------
-- item_view — the read surface. `effective_status` is DERIVED, never stored.
-- ----------------------------------------------------------------------------
-- effective_status = status_manual ?? derive(playtime_minutes, capabilities)
--
-- The capability lookup disappears here, and that is the point of making
-- `playtime_minutes` nullable: NULL (this source cannot know) and 0 (it knows,
-- and it is zero) both fall to 'not_started', which is exactly what
-- derive() returns in both cases — so the derivation is TOTAL over the row and
-- needs nothing from compiled-in code.
--
-- ZERADO is never derived. `in_progress` is the only automatic transition in
-- the product (05-state-machine §2, §4).
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
JOIN item_user u ON u.item_id = i.id;


-- ----------------------------------------------------------------------------
-- Indexes. Every one is justified by a NAMED query in
-- docs/data/04-indexing-plan.md. One is not, and it says so there.
-- ----------------------------------------------------------------------------

-- Q-UPSERT   the sync's ON CONFLICT target, and 09-erd's uniqueness constraint.
CREATE UNIQUE INDEX item_provider_identity ON item(provider_id, provider_ref);

-- Q-LIB-PAGE · Q-LIB-COUNTS · Q-DECK-PAGE
-- Partial, over exactly the default row set: shown = not absent, not merged.
CREATE INDEX item_shelf_order ON item(sort_title)
    WHERE absent_since IS NULL AND merged_into IS NULL;

-- Q-ABSENT-EXISTS · Q-ABSENT-COUNT · Q-ABSENT-LIST
-- Usually an EMPTY index. That is the point: the common case (nothing absent)
-- becomes an empty-B-tree probe instead of a full-table scan, on a query
-- Z-07 runs to decide whether the fifth chip renders at all.
CREATE INDEX item_absent ON item(absent_since) WHERE absent_since IS NOT NULL;

-- ADR-0001 D4 mandates this index in Phase 1. It has NO Phase 1 query.
-- Recorded as such in 04-indexing-plan.md §5 rather than silently kept or
-- silently dropped.
CREATE INDEX item_uid_lookup ON item(item_uid);

-- Q-LAST-SYNC   read by Z-04's banner and every Z-03 terminal state.
CREATE INDEX sync_run_recent ON sync_run(provider_id, started_at DESC);
