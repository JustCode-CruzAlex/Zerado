-- ============================================================================
-- Zerado · migration 0002 · Phase 2 enrichment          NOT SHIPPED IN PHASE 1
-- ----------------------------------------------------------------------------
-- Ships with the Phase 2 BINARY, not before. Creating empty tables for features
-- that do not exist is the same speculative generality ADR-0001 D5 pruned; the
-- ladder shows the future, and the binary applies only what its own phase uses.
--
-- EVERYTHING IN THIS FILE IS CACHE. Every table here can be emptied with a
-- DELETE and the product keeps working, minus pictures and prose. That is the
-- discard test in docs/data/05-truth-and-cache.md §4, and it is why these rows
-- live in their own tables rather than as columns on `item`.
--
-- The one exception is `item_mood.source = 'user'`, which is TRUTH sitting in a
-- cache table. It is the one place in this schema where the line is a column
-- value rather than a table boundary — named as the weaker guarantee it is, in
-- 05-truth-and-cache.md §5.
-- ============================================================================


-- ----------------------------------------------------------------------------
-- metadata — one row per item, from whichever provider Phase 2 ships with
-- ----------------------------------------------------------------------------
CREATE TABLE metadata (
    item_id          INTEGER PRIMARY KEY REFERENCES item(id) ON DELETE CASCADE,
    sinopse          TEXT NULL,

    -- A LOCAL CACHE PATH. Never a remote URL — nothing in Zerado renders from
    -- the network (06-data-seams §1), and the CHECK below is that rule made
    -- structural rather than remembered.
    cover_ref        TEXT NULL,
    cover_fetched_at TEXT NULL,

    released_at      TEXT NULL,   -- reduced precision allowed, as owned_since
    genres           TEXT NULL,   -- JSON array

    source_provider  TEXT NOT NULL,
    -- The credit line is a property OF THE SOURCE (D1), so swapping the source
    -- swaps the credit. Stored with the row it credits, never in config.
    attribution      TEXT NOT NULL,

    -- NOT NULL, because 07-offline-contract §4 forbids rendering a
    -- network-derived value without its age, and the schema is a better place
    -- to enforce that than the renderer's memory.
    fetched_at       TEXT NOT NULL,

    CHECK (source_provider <> '' AND attribution <> ''),
    CHECK (sinopse IS NULL OR sinopse <> ''),
    CHECK (cover_ref IS NULL
           OR (cover_ref NOT LIKE 'http://%' AND cover_ref NOT LIKE 'https://%')),
    CHECK ((cover_ref IS NULL) = (cover_fetched_at IS NULL)),
    CHECK (genres IS NULL OR json_valid(genres)),
    CHECK (fetched_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (cover_fetched_at IS NULL OR cover_fetched_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (released_at IS NULL
           OR released_at GLOB '[0-9][0-9][0-9][0-9]'
           OR released_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]'
           OR released_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]')
) STRICT;


-- ----------------------------------------------------------------------------
-- mood_tag — the vocabulary
-- ----------------------------------------------------------------------------
-- NO `applies_to`. It was per-type machinery from the pruned polymorphic model,
-- and a per-type array in a Phase 1 table is exactly the speculative generality
-- the founder cut (09-erd §2). If a second media type ever arrives, scoping is
-- one nullable column added by a migration — see 07-the-door-three-times.md.
--
-- `key` is the stable identifier the recommender reasons over; `label` is what
-- the player sees, so wording can change without migrating data.
CREATE TABLE mood_tag (
    id    INTEGER PRIMARY KEY,
    key   TEXT NOT NULL,
    label TEXT NOT NULL,
    CHECK (key <> '' AND label <> ''),
    CHECK (key = lower(key) AND key NOT GLOB '* *')
) STRICT;

CREATE UNIQUE INDEX mood_tag_key ON mood_tag(key);


-- ----------------------------------------------------------------------------
-- item_mood — the assignment
-- ----------------------------------------------------------------------------
-- `source` decides whether a tag crosses the Phase 4 boundary. ONLY 'user' does.
-- A 'user' row has no confidence: a person is not 0.8 sure they find something
-- cosy.
CREATE TABLE item_mood (
    item_id     INTEGER NOT NULL REFERENCES item(id)     ON DELETE CASCADE,
    mood_id     INTEGER NOT NULL REFERENCES mood_tag(id) ON DELETE CASCADE,
    source      TEXT    NOT NULL,
    confidence  REAL    NULL,
    assigned_at TEXT    NOT NULL,

    PRIMARY KEY (item_id, mood_id),
    CHECK (source IN ('user','inferred')),
    CHECK ((source = 'inferred') = (confidence IS NOT NULL)),
    CHECK (confidence IS NULL OR (confidence > 0.0 AND confidence <= 1.0)),
    CHECK (assigned_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;

-- Q-MOOD-FILTER: "show me everything tagged cosy" reaches item_mood by mood_id.
-- The table's own PK covers the other direction (this game's tags).
CREATE INDEX item_mood_by_mood ON item_mood(mood_id, item_id);

-- Q-CACHE-DISCARD: the inferred rows are the ones a discard removes, and this
-- partial index is usually a fraction of the table.
CREATE INDEX item_mood_inferred ON item_mood(item_id) WHERE source = 'inferred';
