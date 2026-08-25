---
title: Zerado — the Phase 4 sync seam, from the storage side
discipline: DATA
doc-no: ZRD-DATA-06
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: concept-explainer
ticket: "#5"
---

# The Phase 4 seam, from the storage side

> **Design the seam; do not build it.** Nothing in this document is in migration 0001 except the
> three columns ADR-0001 D4 already requires. What follows is what those three columns are *for*,
> and what Phase 4 will have to add.

The boundary itself is decided and is not reopened here:
ADR-0001 **D4** · [`06-data-seams.md`](../blueprint/06-data-seams.md) §6 ·
[`09-erd.md`](../blueprint/09-erd.md) §5. **Only what the player typed crosses.**

---

## 1 · The payload is a query, and that is the point of the split

```sql
-- Everything that crosses, in Phase 1 terms. There is no fourth clause.
SELECT i.item_uid, u.status_manual, u.status_changed_at
  FROM item i JOIN item_user u ON u.item_id = i.id
 WHERE u.status_changed_at IS NOT NULL;          -- the player has had an opinion

-- Phase 2 adds exactly one more:
SELECT i.item_uid, t.key, m.assigned_at
  FROM item_mood m JOIN item i ON i.id = m.item_id JOIN mood_tag t ON t.id = m.mood_id
 WHERE m.source = 'user';

-- And the one row class with no other copy anywhere:
SELECT item_uid, title, platform, acquisition, owned_since, created_at
  FROM item WHERE provider_id = 'physical';
```

Three statements, and none of them mentions a credential, a playtime, a cover or a price. That is
[`05-truth-and-cache.md`](./05-truth-and-cache.md)'s table boundary paying a second time: *what
crosses* is not a filter someone has to maintain — it is **a table, a column value, and a
provider name**.

---

## 2 · What crosses, and what never does

| Crosses | Never crosses | Why never |
|---|---|---|
| `item_user.status_manual` + `status_changed_at` | the library itself (`item`) | Re-derivable from the player's own sources, on any device |
| `item_mood` rows where `source='user'` | `item_mood` rows where `source='inferred'` | A machine guessed it; each device's machine can guess again |
| `physical` items in full | `playtime_minutes`, `last_played_at`, `achievements_*` | Provider facts. Each device re-fetches them with its own keys |
| `item_uid`, as the join key | `item.id` | A rowid is local and means nothing on another machine |
| | `metadata`, `price_quote`, cover blobs | Cache, with an age. Re-fetch is cheaper than sync |
| | `setting` | A preference is about **this machine** — its theme, its terminal's image support, its audio. "Only what the player typed" is about the **library** |
| | **credentials. Ever.** | A ratified promise. Each device connects with its own keys |

**`setting` is the one worth arguing about**, because a player does type it. It stays local because
half of what it holds describes the *machine* — whether this terminal draws images, whether audio
is available — and syncing a machine's capabilities to a different machine produces a setting that
is wrong on arrival. If Phase 4 ever wants theme to follow the player, that is a **new,
explicitly-scoped table**, not a decision to sync `setting` wholesale.

---

## 3 · Where a conflict is resolved: on the device

**Last-write-wins on `status_changed_at`, per item** (D4, unchanged). Two consequences for the
storage side:

1. **The server stores state changes; it does not adjudicate them.** Merge happens on the device,
   against the local `item_user` row, which is where the merge can be *shown* to the player —
   `Z-22 Devices and sync` shows the last merge and what it resolved.
2. **Timestamps are 20-character UTC RFC 3339, so lexical order is chronological order.** No
   parsing, no timezone, no collation. `MAX(status_changed_at)` is the winner, in SQL, on either
   side of the wire.

### 3.1 · The tie, named rather than left to chance

Two devices changing the same game **in the same second** is vanishingly rare and is not
impossible, and *"whichever arrives first"* is a coin flip nobody can explain afterwards.

**Rule: on an exact tie, the local device's value wins, and the merge is shown in `Z-22` as a
tie.** It is deterministic on each device (a device never silently loses to a stranger's clock),
it needs no sub-second precision, and it is visible. Sub-second timestamps were the alternative;
they cost a variable-width column, break the fixed-20 lexical-ordering property, and buy an
ordering between two events no human distinguished.

### 3.2 · The clear must be able to win, and Phase 1 is where that is decided

Covered in [`02-physical-schema.md`](./02-physical-schema.md) §6.1 and repeated here because it is
a **Phase 4 correctness property that only Phase 1 can provide**:

> Clearing an override is a change. It keeps its `status_changed_at`.

If a clear dropped the timestamp, it would lose every last-write-wins comparison, and a stale
`ZERADO` from another device would silently reappear — the product undoing the player's most
recent decision. `status_changed_at IS NULL` therefore means exactly one thing: *this player has
never had an opinion about this game.*

The schema enforces the half it can — `CHECK (status_manual IS NULL OR status_changed_at IS NOT
NULL)` — and `seed-minimal.sql` row 2 is the cleared case, present so a test can assert it before
Phase 4 exists to need it.

---

## 4 · What Phase 1 already carries, and what Phase 4 must add

| In migration 0001 | Why it could not wait |
|---|---|
| `item.item_uid` | Adding it in Phase 4 means inventing stable identities for rows **whose titles the player has since edited**. That is the migration nobody can write correctly |
| `item.merged_into` | So a merge joins two rows without rewriting primary keys |
| `item_user.status_changed_at` | The conflict is decided by a timestamp that had to be recorded *when the change happened* |

| Phase 4 adds (migrations 0004+) | Shape | Not built now because |
|---|---|---|
| **Device identity** | `setting` rows, or a one-row `device` table: a local device id and label | Nothing in Phase 1 has two devices |
| **A sync cursor** | `sync_state(peer, last_pulled_at, last_pushed_at)` | There is no peer |
| **A merge log** | `merge_event(item_uid, resolved_at, winner, loser_value, reason)` — append-only, what `Z-22` renders | `Z-22` is a Phase 4 screen |
| **A pending-merge queue** | ambiguous `item_uid` matches awaiting the player's decision | The merge that produces them does not exist |

None of these needs a change to an existing table. **That is the test of whether this seam was
designed correctly**, and it is worth stating as the claim it is: *Phase 4 adds tables; it does not
alter `item` or `item_user`.*

---

## 5 · Two things Phase 4 must not get wrong, from the storage side

**`merged_into` chains.** The schema stops the one-row cycle (`CHECK (merged_into <> id)`) and
cannot stop a longer one — A → B → A is two legal rows. The merge must therefore resolve to a
**root** (follow `merged_into` to a row where it is `NULL`) and refuse to create an edge whose
target already leads back to the source. That is Phase 4 code, and it is named here because the
`CHECK` looks like it covers more than it does.

**Manually-entered games are the one class whose loss is unrecoverable**, so they belong in the
**first** sync payload, not a follow-up (D4). The storage side makes them cheap to find —
`WHERE provider_id = 'physical'` — and they carry no provider facts to strip, because they never
had any.
