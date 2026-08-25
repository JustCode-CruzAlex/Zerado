---
title: Zerado — the door stays open
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-11
rev: B
date: 2026-08-25
status: draft — for founder ratification
archetype: concept-explainer
ticket: "#2"
---

# The door stays open

> **Founder, 2026-08-25:** *"At this point don't even think on books and other media types. What I
> would like is let that door open, because at some point books can be a part of Zerado."*

**Zerado is a games product.** This document is one page long on purpose.

---

## 1 · What this replaces, and why the correction is recorded

Revision A of this document was a full media-polymorphic model: a `media_item` core, typed
extensions for books, films and series, a four-type state verification, and two findings about how
films and series differ. It was written in good faith against a direction to *"prepare the database"*
for books and possibly films.

**That over-reached, and the founder pruned it.** The speculation is exactly what he is cutting, and
it was beginning to shape Phase 1 tables and Phase 1 states — which is the cost that makes
speculative generality expensive rather than merely unused.

The correction is recorded rather than silently applied, because revision A is in the PR history and
a reader who finds it should know it was superseded deliberately and by whom.

---

## 2 · The one affordance, and it is cheap

**Two decisions, both reversible, both costing almost nothing today:**

1. **A type discriminator on the core entity.** One column, one value (`game`), constrained to it.
   Adding a second value later is a migration that adds a row to a check constraint.
2. **Do not name the table `games`.** Name it `item` — or anything that is not the word `games` —
   so a future type does not require renaming the table every foreign key points at.

That is the entire affordance. It is not a design for books; it is the absence of a decision that
would make books expensive.

**Everything else is a games model.** Playtime is playtime. Achievements are achievements. The four
states mean what they mean for a game. No column is generalised, no interface takes a type
parameter, no screen has a mode it does not need, and **nothing in Phase 1 is shaped by a media type
that does not exist.**

---

## 3 · What the cost would have been, stated so the trade is visible

Had the type dimension been retrofitted instead, the migration would have had to rename the central
table and every foreign key referencing it, and add a discriminator to rows that predate the concept
— on a file that lives on the player's machine and that Zerado promises never to lose.

Two columns' worth of foresight buys that away. Anything beyond two columns is buying insurance for
a fire that may never happen, in a phase whose job is to ship a games tracker.

---

## Appendix · Not now

Recorded only so the thinking is not re-done from scratch if books are ever actually on the table.
**None of this shapes Phase 1, and none of it should be read as a plan.**

- The four states — *not started* · *in progress* · *zerado* · *abandoned* — do generalise to a book.
  A book you finished **is** *zerado*; the word was never about games. That is a reassuring
  observation, not a design.
- Progress has a different unit and a different source per type, and for a paper book there is no
  automatic signal at all. Zerado already handles that case for physical copies, so the mechanism
  would not be new.
- Films and series are **not** the same kind of thing, and treating them as one type would be the
  first mistake. That is a note for a future ticket, not a finding for this one.

**Nothing above is a commitment, a schema, or a phase.**
