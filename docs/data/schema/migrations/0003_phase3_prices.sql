-- ============================================================================
-- Zerado · migration 0003 · Phase 3 prices              NOT SHIPPED IN PHASE 1
-- ----------------------------------------------------------------------------
-- THE RETENTION RULE IS THE TABLE'S SHAPE: there is no price HISTORY table, so
-- there is nothing to prune.
--
-- The risk the ticket names — "price history grows without bound if every poll
-- writes a row" — is answered by not writing a row per poll. The price seam's
-- `Quote` already carries `Low` and `LowAt` (06-data-seams §4), so the all-time
-- low arrives from the provider and never has to be reconstructed from Zerado's
-- own observations. One row per (item, shop, currency), UPSERTed in place.
--
-- Bound: items x shops x currencies. At 400 titles, three shops and one
-- currency that is 1 200 rows, forever, no matter how often prices are polled.
--
-- `local_low_*` is the fallback for a provider that does not supply a low. It
-- is a watermark updated IN PLACE, not a log — still one row.
--
-- NO affiliate URL, and its absence is the structure. Affiliate links are
-- dropped (ADR-0001 D1); `url` is a plain shop link and there is no column a
-- commission tag could live in.
-- ============================================================================

CREATE TABLE price_quote (
    item_id         INTEGER NOT NULL REFERENCES item(id) ON DELETE CASCADE,
    shop            TEXT    NOT NULL,
    currency        TEXT    NOT NULL,      -- ISO 4217, uppercase

    current_cents   INTEGER NOT NULL,      -- minor units. Never a float: money.
    low_cents       INTEGER NULL,          -- the provider's all-time low
    low_at          TEXT    NULL,          -- "a low with no date is not information"

    local_low_cents INTEGER NULL,          -- watermark, only when the provider gives none
    local_low_at    TEXT    NULL,

    url             TEXT    NOT NULL,      -- a plain shop link
    source_provider TEXT    NOT NULL,
    attribution     TEXT    NOT NULL,

    -- NOT NULL and rendered. "A price without its age is the product giving
    -- financial advice from memory."
    observed_at     TEXT    NOT NULL,

    PRIMARY KEY (item_id, shop, currency),

    CHECK (shop <> '' AND url <> '' AND source_provider <> '' AND attribution <> ''),
    CHECK (currency GLOB '[A-Z][A-Z][A-Z]'),
    CHECK (current_cents >= 0),
    CHECK (low_cents IS NULL OR low_cents >= 0),
    CHECK ((low_cents IS NULL) = (low_at IS NULL)),
    CHECK ((local_low_cents IS NULL) = (local_low_at IS NULL)),
    -- A provider low and a local watermark are alternatives, never both.
    CHECK (low_cents IS NULL OR local_low_cents IS NULL),
    CHECK (observed_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (low_at       IS NULL OR low_at       GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z'),
    CHECK (local_low_at IS NULL OR local_low_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z')
) STRICT;

-- Q-WATCHLIST (Z-20): "which of my watched titles are near their low?" reaches
-- every quote in one pass and needs no index the PK does not already give.
-- Deliberately no index on observed_at: no screen asks "which quotes are
-- stalest" — the age is RENDERED, not FILTERED ON. If Z-20 ever grows a
-- "refresh the stalest first" behaviour, that is the query that earns it.
