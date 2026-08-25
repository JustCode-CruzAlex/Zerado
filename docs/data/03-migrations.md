---
title: Zerado — the migration strategy
discipline: DATA
doc-no: ZRD-DATA-03
rev: A
date: 2026-08-25
status: draft — part of the Phase 0 blueprint bundle
archetype: implementation-plan
ticket: "#5"
---

# The migration strategy

> **There is no DBA and there is no ops team.** The program migrates its own file on start, or it
> refuses to run and says why. Those are the only two outcomes, and both of them are named.

Drawn as [sheet 04](../adr/charts/svg/ZRD-ERD-04.svg) ·
[cyanotype](../adr/charts/svg/ZRD-ERD-04.cyanotype.svg).

---

## 1 · The five rules

1. **Forward-only.** There are no down migrations. A mistake is corrected by a **new** migration,
   never by reversing one. A down migration is code that has, by definition, never run on the file
   it is asked to repair.
2. **A released migration is immutable.** Editing `0001` after a release means two files claim
   version 1 and disagree about what it did — and neither knows which one it is.
3. **One transaction per migration**, and the `INSERT INTO schema_migration` is **inside it**.
   There is no half-applied version, because the bookkeeping and the DDL commit together or not at
   all.
4. **A database from a newer binary is a fatal error**, naming both versions. Silently proceeding
   against a schema you do not understand is how one file becomes two incompatible files.
5. **A fresh install runs the same ladder an upgrade runs.** There is no separate
   "create everything" script — §8.

---

## 2 · The bookkeeping table, and why `written_by` is in it

```sql
CREATE TABLE schema_migration (
    version    INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL,
    written_by TEXT NOT NULL      -- e.g. "zerado/0.4.1 (darwin/arm64)"
) STRICT;
```

`Z-11 Fatal error` promises to name **both** versions — the one that wrote the file and the one
that is running. A schema number alone cannot do that: it tells the player their file is from the
future, which is the half they cannot act on. The **binary** version tells them which release to
go and fetch.

**It cannot be backfilled**, because the binary that wrote a row is exactly the thing that is no
longer around to ask. *(Found by `fft-tui-designer` while writing `Z-11`, and carried here
unchanged.)*

### 2.1 · `PRAGMA user_version` is deliberately not used

SQLite offers a built-in integer header field, and mirroring `MAX(version)` into it would let a
curious player read the version with `sqlite3 library.db "pragma user_version"`.

**Rejected: it would be a second source of truth for one fact.** They can disagree — a crash
between the table write and the pragma write, a migration that forgets the mirror — and when they
do, nothing decides which is right. The saving is one command's worth of convenience:
`sqlite3 library.db "select max(version) from schema_migration"` is already a one-liner, and it
returns the answer *and* who wrote it.

---

## 3 · The five outcomes at start-up

| | Condition | What happens |
|---|---|---|
| **A** | The file does not exist | Create it; apply `0001 … HEAD` in order, each in its own transaction, each writing its own row |
| **B** | `MAX(version) == HEAD` | Open it. One page read, ~0.1 ms. The common path |
| **C** | `MAX(version) < HEAD` | Back up if a destructive step is ahead (§5), then apply the gap in order |
| **D** | `MAX(version) > HEAD` | **REFUSE.** `Z-11`, naming the file's version, its `written_by`, and the running binary |
| **E** | No `schema_migration` table, but the file is **not empty** | **REFUSE.** Some other program's SQLite file, or a corrupted one. Never migrate a schema you did not write |

**E is the one that is easy to omit.** *"No bookkeeping table"* looks like *"a fresh file"* until
someone points `$ZERADO_DB` at the wrong path — and then Zerado runs `CREATE TABLE item` inside
somebody else's database. The check is one query: does `sqlite_schema` contain anything at all
that is not `sqlite_%`?

**Before any of the five**, the version precondition: `sqlite_version() >= 3.37.0`, or refuse with
a named error. Without it, an old library fails *inside* migration 0001 with
`unrecognized token: "STRICT"`, which is a message about a keyword rather than about the player's
situation.

---

## 4 · What "idempotent" actually means here

Rule 3 makes atomicity the real guarantee: a crash mid-migration rolls back, and the file is
byte-identical to what it was. Re-running then re-attempts a version that was never recorded.

So the belt is `version NOT IN (SELECT version FROM schema_migration)` — checked before each step
— and the braces are the transaction. `IF NOT EXISTS` is used only where it is *semantically*
right (a `CREATE INDEX` that a repair path may re-run), never as a way to make a half-applied
migration look successful.

---

## 5 · Backups, and why not a file copy

**Before any migration that drops or rewrites a column**, the file is backed up beside itself as
`library.db.pre-<version>` ([`09-erd.md`](../blueprint/09-erd.md) §6), and `Z-11` names that path
if the migration fails.

```sql
VACUUM INTO 'library.db.pre-0004';
```

**Not `cp`, and not `io.Copy`.** A byte copy of a WAL-mode database *while it is open* copies the
main file without the committed pages still living in `-wal`. The result is a backup that looks
fine and fails to open — discovered at the exact moment it is needed. `VACUUM INTO` produces a
consistent, defragmented, single-file snapshot through SQLite's own transaction machinery.

Measured ([`evidence/measurements.md`](./evidence/measurements.md) §7): **under 10 ms** at both
400 and 10 000 rows, and the output is *smaller* than the live file. Fast enough to do it every
time, which is the only backup policy that survives contact with a deadline.

**Before writing, check free space.** If the backup would not fit, refuse **before** touching
anything — a migration that runs out of disk halfway is precisely the state this whole document
exists to make impossible.

Backups are never deleted automatically. They are the player's file, in the player's directory,
and Zerado does not clean up other people's data.

---

## 6 · Rebuilding a table, when it comes to that

SQLite cannot `ALTER TABLE … DROP CONSTRAINT`, cannot change a column type, and cannot widen a
`CHECK`. Any of those is the twelve-step table rebuild, and the ordering has two traps:

```sql
PRAGMA foreign_keys = OFF;          -- MUST be outside the transaction: SQLite
                                    -- silently ignores this pragma inside one
BEGIN IMMEDIATE;
  DROP VIEW    item_view;           -- dependents FIRST — trap 3
  DROP TRIGGER item_user_exists;
  CREATE TABLE item_new ( … );      -- the new shape
  INSERT INTO item_new SELECT … FROM item;
  DROP TABLE item;
  ALTER TABLE item_new RENAME TO item;
  -- recreate EVERY index, trigger and view: DROP TABLE took the indexes with it
  CREATE UNIQUE INDEX item_provider_identity ON item(provider_id, provider_ref);
  CREATE INDEX item_shelf_order ON item(sort_title)
      WHERE absent_since IS NULL AND merged_into IS NULL;
  CREATE INDEX item_absent ON item(absent_since) WHERE absent_since IS NOT NULL;
  CREATE INDEX item_uid_lookup ON item(item_uid);
  CREATE TRIGGER item_user_exists …;
  CREATE VIEW item_view AS …;
  PRAGMA foreign_key_check;         -- MUST be inside, before COMMIT
COMMIT;
PRAGMA foreign_keys = ON;
```

**Measured**, through `modernc.org/sqlite` on a 400-row library: **9.9 ms**, `integrity_check`
clean, every `item_user` row byte-identical.

**Trap 1:** `PRAGMA foreign_keys` cannot be changed inside a transaction — SQLite ignores it
without error, so a rebuild that sets it inside the transaction runs with FKs still on and fails
on the `DROP`.

**Trap 2:** `DROP TABLE` takes the table's indexes with it. A rebuild that forgets to recreate
`item_shelf_order` leaves a schema that is *correct* and a root screen that is **82× slower** at
10 000 rows ([`evidence/measurements.md`](./evidence/measurements.md) §3) — a performance
regression with no error message. §8's fingerprint test is what catches it.

**Trap 3 — found by running the recipe rather than reading it.** `DROP TABLE item` fails through
the driver while `item_view` still references the table: *`error in view item_view: no such table:
main.item`*. Dependent **views and triggers must be dropped first** and recreated last.

That is the trap. **This is how it hid**, and it is the more useful half: a first attempt at this
recipe appeared to succeed under `sqlite3 db < migration.sql`, with a clean `integrity_check` and
the right row count. It had not. **The `sqlite3` CLI continues past errors unless `-bail` is
given** — it had printed `view item_view already exists` and carried on, and the database happened
to end up correct because `ALTER TABLE … RENAME TO` silently repaired the old view's reference.

Two rules follow, and the second one is the one worth keeping:

1. **Drop dependent views and triggers first.**
2. **Never validate a migration with a bare `sqlite3 db < file.sql`.** Use `sqlite3 -bail`, or the
   driver. A validation that cannot fail is not a validation, and this one told a comfortable lie
   in a document whose whole value is that its facts were checked.

**The regeneration migration**, for when `fold()` changes (a new locale's article list, an
`x/text` upgrade), is a rebuild's cheaper cousin — an `UPDATE` over derived columns, no schema
change:

```sql
-- 000N_regenerate_normalised_titles.sql
UPDATE item SET sort_title = :folded_sort, search_title = :folded_search WHERE id = :id;
INSERT INTO setting(key,value) VALUES('normaliser_version','2')
    ON CONFLICT(key) DO UPDATE SET value = excluded.value;
```

It is driven from Go because `fold()` is Go. **`item_uid` is deliberately *not* regenerated** —
regenerating identities would break the merge hints on every other device, which is the one thing
`item_uid` exists to preserve.

---

## 7 · Failure, and what the player sees

| Failure | What the file looks like afterwards | What the player is told |
|---|---|---|
| A migration statement errors | **Byte-identical to before.** The transaction rolled back | `Z-11`: which migration, the error, and the backup path if one was taken |
| The process is killed mid-migration | Byte-identical. SQLite rolls back on next open | Nothing — the next start retries |
| The disk fills during the backup | Untouched: the check happens **first** | A refusal naming the space needed |
| The disk fills during the migration | Rolled back | `Z-11` |
| A newer binary wrote the file | Untouched | `Z-11`, both versions named (§2) |
| The file is another program's | Untouched | `Z-11`, outcome E |
| `-wal`/`-shm` survive a crash | Recovered on open, then checkpointed | Nothing. This is WAL working |

**Never a partial success, and never a silent one.** Every branch above ends either with an open
library or with a screen that names three things: what happened, why, and the next action — which
is the brand manual's refusal template, and it is why `Z-11` needs `written_by`.

---

## 8 · The fingerprint test, and the drift it catches

> **A fresh install runs the same ladder an upgrade runs.**

The tempting alternative is a `schema.sql` for fresh installs and migrations for upgrades. They
drift — always, and quietly, because nobody runs both paths on the same day. A new column added to
`schema.sql` and not to a migration means every upgraded file lacks it, and the crash arrives
months later on somebody else's machine.

So [`schema/schema.head.sql`](./schema/schema.head.sql) is **generated and never executed**:

```bash
for m in docs/data/schema/migrations/*.sql; do sqlite3 head.db < "$m"; done
sqlite3 head.db ".schema --indent"     # → schema/schema.head.sql
```

The Phase 1 test asserts three things:

1. **The fingerprint.** `SELECT type||' '||name FROM sqlite_schema WHERE name NOT LIKE 'sqlite_%'
   ORDER BY type, name` over a freshly-migrated database equals the recorded list. A **Phase 1**
   binary's HEAD is 0001 and its fingerprint is **14 objects** — 6 tables, 5 indexes, 2 triggers,
   1 view. The full ladder through 0003 is **21** — 10 tables, 8 indexes, 2 triggers, 1 view, which
   is what [`schema/schema.head.sql`](./schema/schema.head.sql) records.
2. **Ladder equality.** A database built by `0001 → 0002 → 0003` and one built by applying the
   same files to an empty file produce **identical** `sqlite_schema` `sql` text. This is what
   catches a forgotten `CREATE INDEX` after a §6 rebuild.
3. **Corpus migration.** For each released version *N*, a checked-in fixture database at version
   *N* migrates to HEAD, and every row of `item_user` is still there with the same values. Truth
   is never lost by a migration, and *"never"* is a test rather than an intention.

A diff in `schema.head.sql` on a pull request that did not intend one **is the drift being
caught**, which is the whole point of committing a generated file.

---

## 9 · The ladder as it stands

| Version | File | Ships with | Destructive? | Backup taken? |
|---|---|---|---|---|
| **0001** | [`0001_phase1_core.sql`](./schema/migrations/0001_phase1_core.sql) | **Phase 1** | no — creates | no |
| 0002 | [`0002_phase2_enrichment.sql`](./schema/migrations/0002_phase2_enrichment.sql) | Phase 2 | no — creates | no |
| 0003 | [`0003_phase3_prices.sql`](./schema/migrations/0003_phase3_prices.sql) | Phase 3 | no — creates | no |

**0002 and 0003 are written but not shipped.** A Phase 1 binary's HEAD is **0001**, and it applies
only that. Creating empty tables for features that do not exist is the same speculative generality
ADR-0001 D5 pruned; the ladder is allowed to show the future, and the binary is not allowed to
build it.

**Numbering:** four digits, zero-padded, monotonic, never reused, never edited after release.
Gaps are allowed (a migration abandoned before release); reuse is not.
