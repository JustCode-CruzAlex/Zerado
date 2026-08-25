---
title: Zerado — Z-07 Filter and search
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-07
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-07 · Filter and search

> **A mode of `Z-04`, not a region and not a route.** It gets its own spec because it carries its
> own key map, its own empty state and its own copy — *"a builder cannot build it from the
> library's spec"* (`01-screen-inventory.md` §1).

**Canon that governs this screen:** `00-design-brief.md` §10 · §3.1 (SC 2.1.2 and 2.1.4, the
two that bite hardest here) · `01-design-system.md` §7 (filter bar) · §10.2 (empty result) ·
§3 (chip) · `02-colour-budget.md` §10 · §2.3 (the focus-ring exemption) ·
`03-designer-manual.md` §3 · §7 anti-patterns 4 and 5 ·
`02-composition.md` §2 (`R = 2`; body rows 1–2) · `03-responsive.md` §3 ·
**`04-navigation-and-focus.md` §5 — the two-step `Esc`** · `05-state-machine.md` §7 rule 2 ·
`07-offline-contract.md` §2 — *"Filter and search: **WORKS**. Identical — it is a `WHERE`
clause"* · `01-screen-inventory.md` §5 — *"`Z-07` must not silently return zero rows. An empty
result names the filter that emptied it."*

---

## 1 · Identity

| | |
|---|---|
| **ID** | `Z-07` |
| **Name** | Filter and search |
| **Phase** | 1 |
| **Kind** | **Mode** of `Z-04`. Not on the route stack, not in the overlay slot |
| **Route in** | `/` from `Z-04` — the query editor takes focus · `f` from `Z-04` — the state chips take focus |
| **Route out** | `Esc` (twice — §12) · `⏎` moves focus to the list, keeping the filter · `q` / `Ctrl-C` |
| **Composition** | **Takes the top of the pinned block at every tier.** 2 body rows at Standard and above; 3 at Narrow, where the chips wrap; 2 at Tiny, one chip at a time |
| **`R`** | **2** — the filter bar and the list. This is the only Phase 1 surface below 120 columns where `Tab` exists |
| **Offline class** | **WORKS.** A `WHERE` clause needs no network |
| **Displayed name** | **none of its own.** The band stays `Zerado ✦ Library` / `LIBRARY` |

> **Why the header band does not change.** Announcing the mode in the band would make the band
> the only mutable chrome on screen — the exact objection `01-design-system.md` §2.4 records
> against putting live values there. The mode announces itself in the body, where the live count
> already lives, and the band stays the stable thing it is on every other screen.

---

## 2 · Purpose

**One sentence:** narrow 400 games down to the handful that match, live as the player types, and
— when nothing matches — say which part of the filter emptied it.

---

## 3 · Mockup — 80 × 24, the editor focused

`tier = Wide` · `leftInset = 3` · **body = 74 × 16** · content begins at **column 4**.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   / the▎                                                           23 of 247   │
│   [○  NOT STARTED]  [◐  IN PROGRESS]  [◉  ZERADO]  [⊘  ABANDONED]              │
│                                                                                │
│     STATE           TITLE                                        HOURS   SRC   │
│     ◉  ZERADO       Bastion                                         6h   STM   │
│     ○  NOT STARTED  Death's Door                                    0h   STM   │
│     ◐  IN PROGRESS  Dark Souls III                                 63h   STM   │
│     ⊘  ABANDONED    Nier: Automata                                 28h   STM   │
│   ▌ ◉  ZERADO       Return of the Obra Dinn                         9h   STM   │
│     ○  NOT STARTED  The Forgotten City                              0h   STM   │
│     ◐  IN PROGRESS  The Last Campfire                               3h   STM   │
│     ○  NOT STARTED  The Legend of Zelda: Breath of the Wild          —   PHY   │
│     ◉  ZERADO       The Messenger                                  12h   STM   │
│     ⊘  ABANDONED    The Outer Worlds                               14h   STM   │
│     ○  NOT STARTED  The Talos Principle 2                           0h   STM   │
│   ROWS  1–11 of 23                                                             │
│   esc back to the list   ⏎ first match   tab state chips   ^c quit             │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Character counts:** body **74** on every row · the editor row **74** (sigil + query + caret,
left; `23 of 247`, right-aligned to column 74) · the chip row **63** · the game rows **74**,
unchanged from `Z-04` · scroll position **18** · footer **64**.

### 3.1 · Row map — 80 × 24, filter active

| Body row | Content | Note |
|---|---|---|
| **1** | the query editor + the live ratio | **replaces the pinned summary** (`01-design-system.md` §5.4) |
| **2** | the four state chips | |
| 3 | respiro | never spent (`Z-04-library.md` §8.2, D-04-2) |
| 4 | column header | |
| **5–15** | the scroll region — **11 game rows** | one fewer than `Z-04`'s 12 |
| 16 | scroll position | |

`1 + 1 + 1 + 1 + 11 + 1 = 16` ✓. With a degrade banner also active the banner takes row 1 and
the list drops to **10**.

### 3.2 · The editor row — anatomy

```
/ the▎                                                           23 of 247
│ │  │                                                           │
│ │  └ the caret — the terminal's own, one cell                   └ the live ratio
│ └ the query, --z-text
└ the mode sigil, --z-primary amber
```

**The caret is the terminal's own cursor**, painted by `bubbles/textinput` as one cell. It is
exempt from the cyan budget by name (`02-colour-budget.md` §2.3) and it is drawn `▎` in these
mockups purely as notation for *"the cell the caret occupies"*. `▎` U+258E is **Ambiguous**
width (verified, UCD 16.0.0) and is **not rendered as a literal glyph**.

### 3.3 · The right-hand slot — the live ratio, and the word channel

> **D-07-1 · In filter mode the pinned summary is the ratio, not the four state counts.**
> `01-design-system.md` §5.4 already says the status bar is *replaced by the filter bar*, and
> `05-state-machine.md` §7 rule 2 requires the summary to describe the filtered set **and say
> which set it is**. A ratio does both in one line: `23 of 247` says how many matched, out of
> how many, live as the player types. A state breakdown of a filtered set is either unreadable
> (a 12-row set you can just read) or tautological (filter to `ZERADO` and the breakdown is
> `◉ 412`). R-10(c) still holds — the ratio is pinned outside the scroll region.

> **D-07-2 · A selected state facet is named in WORDS in the same slot.** This is the channel
> that survives `NO_COLOR`. `01-design-system.md` §7.3 specifies a selected chip as *"the state
> colour, bold"* and an unselected one as `--z-text-tertiary` — under `NO_COLOR` that leaves
> **bold alone**, one channel, which is a co-render failure on a control. Naming the facet in
> the ratio slot closes the gap using vocabulary the system already has, and introduces **no new
> glyph**.

| Facets selected | Right-hand slot | Cells |
|---|---|---|
| none | `23 of 247` | 9 |
| one | `ZERADO · 1 of 247` | 17 |
| two | `ZERADO+ABANDONED · 4 of 247` | 27 |
| three or four | `3 STATES · 41 of 247` — the count of facets, because four labels do not fit | ≤ 20 |
| nothing matches | `0 of 247` | 8 |
| **the absent facet** (D-07-8) | `ABSENT · 3 of 250` | 17 |
| **a state facet + absent** | `ZERADO+ABSENT · 2 of 250` | 24 |
| **three or four states + absent** | `3 STATES+ABSENT · 41 of 250` | ≤ 27 |

The slot is right-aligned to the body's right edge. At Narrow it shortens to `23/247` and the
facet names drop — at that tier the chips are on their own two rows and the selection is legible
without them.

> **D-07-9 · Once any row is absent, the ratio's denominator is the whole library file, not the
> default view.** With no absent rows the two are the same number and nothing changes. With three
> absent, an unfiltered `/` reads **`247 of 250`** — which is the same fact `Z-04`'s summary
> carries as `3 absent`, said in the vocabulary this row already uses. The alternative — keeping
> the denominator at 247 and rendering `3 of 3` when the absent facet is selected — is true but
> uninterpretable, and it would make the two pinned rows disagree about how big the library is.
> The `ABSENT` name is appended after any state names with `+`, and the `N STATES` collapse
> applies to the four state facets only.

### 3.4 · The state chips

```
[○  NOT STARTED]  [◐  IN PROGRESS]  [◉  ZERADO]  [⊘  ABANDONED]
```

| Part | Cells |
|---|---|
| `[` + glyph field 2 + gap 1 + label + `]` | 16 · 16 · 11 · 14 |
| gaps | 2 · 2 · 2 |
| **total** | **63** ≤ 74 |

The glyph sits in the **same 2-column width-aware field** as the ledger's, so the chips align
with each other whatever a terminal decides an Ambiguous glyph is worth. **Co-render holds even
as a control:** the chip keeps colour **and** glyph **and** label. **The brackets are the control
boundary and are `--z-border-strong` (4.08, WCAG 1.4.11) — never `--z-border` (1.53).**

| Chip state | Rendering |
|---|---|
| Unselected | glyph + label in `--z-text-tertiary` (interim: uncoloured); brackets `--z-border-strong` |
| **Selected** | glyph + label in the **state colour**, **bold**, **and named in the ratio slot** (D-07-2) |
| Focused (the chip row has focus) | the focused chip additionally carries the `▌` marker in the 2-cell gutter to its left, shifting the row right by 2 |

#### The fifth chip — `[ABSENT]`, and why it has no glyph

`06-data-seams.md` §2.4 excludes tombstoned rows from `Z-04`'s default view and says they
*"remain findable by filter"*. **With no facet there is no filter, so the seam's promise would
have no mechanism.** This is the mechanism.

```
[ABSENT]
```

> **D-07-8 · The absent facet is a fifth chip in the existing chip row, it carries no glyph, and
> it renders only when at least one row is absent.**
>
> **Why a chip and not a key.** Every other candidate needs a binding that does not exist —
> which is exactly why the `source` facet has been stuck in §17 since rev A. The chip row is
> already reachable with `f`, `Tab`, `←` `→` and `space`, all bound. **A fifth chip needs zero
> new keys**, so the seam's promise is deliverable on the day this screen is built rather than
> after a spine decision. It also rules out a query syntax like `absent:` — D-07-5 fixes Phase 1
> matching as a plain substring match on the title, and a reserved word would be a query language
> arriving through the back door.
>
> **Why no glyph, stated positively.** `absent` is **not a fifth state** (§2.4). A fifth glyph
> beside `○ ◐ ◉ ⊘` would assert that it is, in the one row where the four are a closed
> vocabulary — and it would need a derivation nobody has run (brand §10 rule 5 — nobody picks a
> glyph at the keyboard). **The missing glyph is itself the channel that says "not a state."**
> Co-render still holds on three counts and none of them is colour: the **word** `ABSENT`, the
> **brackets** that mark it a control, and the **ratio slot** naming it when selected (D-07-2).
> Under `NO_COLOR` it is unchanged.
>
> **Why only when `n ≥ 1`.** A chip that always renders and always yields zero is furniture, and
> it advertises a condition the player does not have — the same reasoning that keeps `m` out of
> the footer until audio is enabled (`01-design-system.md` §5.3). It also buys an invariant worth
> having: **selecting `[ABSENT]` alone can never return zero rows**, because the chip does not
> exist when there are none.
>
> **Colour.** Unselected `--z-text-tertiary` (interim: uncoloured), like every other chip.
> Selected **`--z-text` `#E9EEF5` / 255 / `bright white`, bold — 16.65 AAA** (brand §4.2, measured,
> and it is `Z-04` §7's game-title token) — a weight and brightness step,
> **not a hue.** A hue would enrol it in the state palette, which is the thing it is not. Brackets
> `--z-border-strong` (4.08), like every other control boundary.
>
> **It is orthogonal, not exclusive.** `[ABSENT]` combines with any state facet and with a query.
> Selecting it swaps the row set from *the shown rows* to *the absent rows*; the state chips then
> filter within that set exactly as they do within the shown one.

**The geometry, counted — and the one cell that decides the composition:**

| | Cells |
|---|---|
| the four state chips | 63 |
| gap | 2 |
| `[ABSENT]` | **8** |
| **one row, unfocused** | **73** ≤ 74 ✓ |
| **one row, focused** — the chip-state table above shifts the row right by 2 | **75** > 74 ✗ |

**Five chips fit on one row until the row is focused, and then they are one cell short.** That is
not a thing to fudge by tightening a gap or shaving a bracket, so **the fifth chip wraps to its
own row and the bar becomes three rows** at Wide, Standard and ExtraWide. The cost is one game
row, and it is charged by a rule that already exists rather than a new one: `Z-04-library.md`
D-04-2 — *the pinned chrome block grows downward, the respiro is never spent, the scroll region
absorbs.* Composition: `1 + 1 + 1 + 1 + 1 + 10 + 1 = 16` ✓, **10 game rows.**

| Tier | Body | Chip rows when `n ≥ 1` | Bar rows | Games | Cost |
|---|---|---|---|---|---|
| ExtraWide, list pane | 66 | **2** — `63` then `8` | 3 | 27 | 1 row of 32 |
| **Wide** | **74** | **2** — `63` then `8` | **3** | **10** | 1 game row |
| Standard | 54 | **see the gap below** | | | |
| Narrow | 36 | **3** — `34` / `27` / `8` | 4 | **7** | 1 game row — **§17 finding 7** |
| Tiny | 30 | **1** — one chip at a time, `[ABSENT]` joins the `Tab` cycle | 2 | unchanged | **none** |

> **A gap in this spec that counting the fifth chip exposed, recorded rather than patched.**
> **At Standard the four chips already do not fit and rev A never said what happens.** The chip
> row measures **63** cells and Standard's body is **54**; §1 and `03-responsive.md` §3 both say
> the bar is *"2 body rows at Standard and above"*, which is only true at Wide and ExtraWide.
> The arithmetic is wrong at `n = 0`, before any absent row exists, so **it is not the absent
> facet's to fix** — resolving it changes Standard's `n = 0` composition and §3.1's row map with
> it. Recorded here because this is where the counting happens and a builder reaching Standard
> needs to know the answer is missing.
>
> **The shape of the answer, for whoever takes it:** wrapping 2/2 fits (`34` and `27`, both ≤ 54)
> and makes Standard a 3-row bar at `n = 0` and a 4-row bar with the absent chip — the same
> ladder as Narrow, one tier up. **Route to `fft-tui-architect` with §17 finding 7**, which is
> the same defect seen from the spine's side.

---

## 4 · Mockup — the editor blurred, the filter still applied

The first `Esc` has been pressed. **The caret is gone, the filter is not.** The footer changes —
and that is the whole mechanism that makes the second `Esc` discoverable rather than clever.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   / the                                                            23 of 247   │
│   [○  NOT STARTED]  [◐  IN PROGRESS]  [◉  ZERADO]  [⊘  ABANDONED]              │
│                                                                                │
│     STATE           TITLE                                        HOURS   SRC   │
│     ◉  ZERADO       Bastion                                         6h   STM   │
│     ○  NOT STARTED  Death's Door                                    0h   STM   │
│     ◐  IN PROGRESS  Dark Souls III                                 63h   STM   │
│     ⊘  ABANDONED    Nier: Automata                                 28h   STM   │
│   ▌ ◉  ZERADO       Return of the Obra Dinn                         9h   STM   │
│     ○  NOT STARTED  The Forgotten City                              0h   STM   │
│     ◐  IN PROGRESS  The Last Campfire                               3h   STM   │
│     ○  NOT STARTED  The Legend of Zelda: Breath of the Wild          —   PHY   │
│     ◉  ZERADO       The Messenger                                  12h   STM   │
│     ⊘  ABANDONED    The Outer Worlds                               14h   STM   │
│     ○  NOT STARTED  The Talos Principle 2                           0h   STM   │
│   ROWS  1–11 of 23                                                             │
│   esc clear filter  / search  f state  ↑↓ move  ⏎ open  ? help  q quit         │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Two things changed and nothing else:** the caret disappeared from the editor row, and the
footer went from `esc back to the list` (64 cells, 3-space separators) to
`esc clear filter …` (68 cells, **2-space** separators per the ladder). The list, the ratio, the
chips and the cursor are all exactly where they were.

**That stability is the point.** The player who has just typed a search almost never means to
throw it away; the first `Esc` gives them the arrow keys back and keeps their work, and the
footer says — at the moment it matters — what the next `Esc` will do.

---

## 5 · Mockup — nothing matches

**The honest case.** `01-screen-inventory.md` §5 is explicit: an empty result **names the filter
that emptied it**. `01-design-system.md` §10.2's copy does not name it; this composition does,
and it says the number at each step.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   / souls▎                                                          0 of 247   │
│   [○  NOT STARTED]  [◐  IN PROGRESS]  [◉  ZERADO]  [⊘  ABANDONED]              │
│                                                                                │
│   Nothing matches.                                                             │
│                                                                                │
│     search   "souls"         3 games                                           │
│     state    ZERADO          0 of those 3                                      │
│                                                                                │
│   247 games in the library.                                                    │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│   esc back to the list   ⏎ first match   tab state chips   ^c quit             │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**The column header and the scroll-position row are absent** — there is nothing to head and
nothing to position. The filter bar (rows 1–2) and the respiro (row 3) stay, because the player
needs to see and edit the thing that emptied the list.

**The diagnostic block, exact geometry:** 2-space indent · facet name field **9** · value field
**16** · then the count. Facet names are lowercase (`search`, `state`, `source`); facet *values*
use their interface casing (`"souls"` quoted, `ZERADO` uppercase, `physical` lowercase).

**Not an error, not red, not a banner.** An empty filter result is the filter working. It is
`--z-text` and `--z-text-secondary` throughout.

### 5.1 · Mockup — the absent facet selected, 80 × 24

**The only way to see a tombstoned row.** Three rows are absent, so the fifth chip exists; it is
selected and the chip row holds focus. `tier = Wide` · body **74 × 16** · the bar is **3 rows**
and the scroll region is **10**.

┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   /                                                        ABSENT · 3 of 250   │
│   [○  NOT STARTED]  [◐  IN PROGRESS]  [◉  ZERADO]  [⊘  ABANDONED]              │
│   ▌ [ABSENT]                                                                   │
│                                                                                │
│     STATE           TITLE                                        HOURS   SRC   │
│   ▌ ◉  ZERADO       Alan Wake                                      19h   STM   │
│     ⊘  ABANDONED    Deadpool                                        2h   STM   │
│     ○  NOT STARTED  Scott Pilgrim vs. the World                     0h   STM   │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│   space toggle  ←→ chip  esc back to the list  ⏎ open  ? help  q quit          │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘

**Character counts:** editor row **74** (`/` left; `ABSENT · 3 of 250` right-aligned to column
74) · state-chip row **63** · absent-chip row **10** (`▌` marker + 2-cell gutter + the 8-cell
chip) · column header **74** · every game row **74**, unchanged from `Z-04` · footer **67**
(§11.3, the chip row has focus).

**Read what the rows say.** *Alan Wake* is `◉ ZERADO` — finished, and no longer returned. That is
the row §2.4 was written to protect: *"a game you finished and no longer own is exactly the row
you would be angriest to lose."* The chips are unchanged; **`absent` is not a fifth state and the
ledger renders no fifth glyph.**

**The position row does not render**, because three rows in a ten-row region do not scroll — the
same rule `Z-10-help-and-key-map.md` H7 states for its own readout. It returns the moment the set
overflows.

**The ledger cursor is still visible while the chip row has focus** — R-10(b), and §12.2's row
column, which is about *which row is selected* and not about which region holds focus. The
focused **control** is `[ABSENT]`, carrying its `▌` marker and its `--z-focus-ring` bracket pair
(exempt, singular by definition); the marker in the ledger is where the player lands when they
leave the chip row. **Two markers, two meanings, and §12.2 already distinguishes them.**

---

## 6 · Mockup — Narrow, 40 × 24: the bar takes three rows

`03-responsive.md` §3: *"Narrow 40–59: Filter bar takes 3 rows (chips wrap)."*

```
┌────────────────────────────────────────┐
│                                        │
│  Zerado ✦ Library                      │
│                                        │
│  LIBRARY                               │
│                                        │
│                                        │
│  / the▎                        23/247  │
│  [○  NOT STARTED]  [◐  IN PROGRESS]    │
│  [◉  ZERADO]  [⊘  ABANDONED]           │
│    ◉ Bastion                           │
│      ZERADO · 6h · Steam               │
│    ○ Death's Door                      │
│      NOT STARTED · 0h · Steam          │
│    ◐ Dark Souls III                    │
│      IN PROGRESS · 63h · Steam         │
│  ▌ ◉ Return of the Obra Dinn           │
│      ZERADO · 9h · Steam               │
│    ○ The Forgotten City                │
│      NOT STARTED · 0h · Steam          │
│    ◐ The Last Campfire                 │
│      IN PROGRESS · 3h · Steam          │
│  ROWS  1–6 of 23                       │
│  esc clear  / edit  tab chip           │
│                                        │
└────────────────────────────────────────┘
```

Body **36 × 16**. Chip row 1 = 34 cells, row 2 = 27. `16 − 3 (bar) − 1 (scroll) = 12` → **6
two-line games**, remainder 0, no respiro — the same rule as `Z-04` at this tier
(`Z-04-library.md` §11.1, D-04-5). The ratio shortens to `23/247` (6 cells) and the facet names
drop; the chips are two rows tall and carry the selection themselves.

**At Tiny (< 40)** the chips do not fit two to a row. `03-responsive.md` §3: *"One filter at a
time, cycled with `Tab`."* One chip renders on body row 2 — `[◉  ZERADO]`, 11 cells — and `Tab`
cycles which one is shown and focused. The bar stays **2 rows**; the ratio shortens to `23/247`;
the footer carries `tab chip`.

---

## 7 · Visual hierarchy

| # | What | Channel |
|---|---|---|
| **1** | **What the player is typing** | Body row 1, top-left, preceded by the amber `/` sigil — the only amber mark in the body. The caret is the only moving thing on screen |
| **2** | **The live ratio** | Right-aligned on the same row, `--z-text` numerals — the answer to *"is this working?"*, updating on every keystroke |
| **3** | **The four chips** | Row 2, four bracketed controls of identical geometry; selection carried by colour + weight + the word in the ratio slot |
| **4** | **The result list** | Unchanged from `Z-04` — same rows, same columns, same cursor. **This is deliberate:** filtering must not make the ledger look like a different screen |
| **5** | **The footer's `esc` hint** | The one piece of chrome that changes meaning, so it is the one the player is meant to notice when it does |
| **6** | Chrome | breadcrumb, column header, scroll position |

**The one thing the player should see first is their own query; the second is the ratio.**
Everything else on screen is the library, unchanged — which is what makes a filter feel like a
lens rather than a new screen.

---

## 8 · Every applied spacing token, by name

| Token | Tiny | Narrow | Standard | **Wide** | ExtraWide |
|---|---|---|---|---|---|
| `OuterMarginX` | 0 | 1 | 2 | **2** | 2 |
| `OuterMarginY` | 0 | 1 | 1 | **1** | 1 |
| `InnerPaddingX` | 1 | 1 | 1 | **1** | 2 |
| `InnerPaddingY` | 0 | 1 | 1 | **1** | 1 |
| `InterElementGap` | 1 | 1 | 1 | **1** | 1 |
| `HeaderBand(tier, false)` | **1** | 3 | 3 | **3** | 3 |
| `leftInset` | **1** | 2 | 3 | **3** | 4 |
| `BodyRect.w × h` | 30 × 21 | 36 × 16 | 54 × 16 | **74 × 16** | 112 × 32 |
| **Filter-bar rows** | **2** | **3** | **2** | **2** | **2** |
| **Visible games** | 8 | 6 | 11 | **11** | **27** |

**Applied:**

| Surface | Token | Value at Wide |
|---|---|---|
| Editor row and chip row left edge | `leftInset` | column 4 — the same column as the header and the game rows |
| Ratio slot right edge | `BodyRect.w` | column 77 |
| bar → column header | `InterElementGap` | 1 row, body row 3 — **never spent** |
| Between chips | a 2-column gap | fixed, so the row is stable as labels change |
| Footer | the reserved footer row | 1 row |

**Zero magic numbers.** The chip's 16/16/11/14 cells are the ratified label lengths plus a
2-cell glyph field and two brackets; the diagnostic block's 9/16 fields are this composition's
declared geometry (§5).

---

## 9 · Colour, glyph and label for every state shown

| Element | Token | Hex | ANSI-256 | 16-colour | Glyph | ASCII | Label | Ratio |
|---|---|---|---|---|---|---|---|---|
| Mode sigil `/` | `--z-primary` | `#FFB000` | **214** | `bright yellow` | `/` | `/` | — | **10.59** AAA |
| The query text | `--z-text` | `#E9EEF5` | **255** | `bright white` | — | — | — | **16.65** AAA |
| The caret | the terminal's own | — | — | — | — | — | — | exempt (§2.3) |
| The ratio numerals | `--z-text` | `#E9EEF5` | **255** | `bright white` | — | — | — | **16.65** AAA |
| The ratio's words (`of`, facet names) | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | — | — | — | **9.36** AAA |
| Chip — selected, not started | `--z-state-not-started` | `#A5A29B` | **247** | `white` | `○` U+25CB | `[ ]` | `NOT STARTED` | **7.62** AA |
| Chip — selected, in progress | `--z-state-in-progress` | `#FFB000` | **214** | `bright yellow` | `◐` U+25D0 | `[~]` | `IN PROGRESS` | **10.59** AAA |
| Chip — selected, zerado | `--z-state-zerado` | `#19E0FF` | **45** | `bright cyan` | `◉` U+25C9 | `[*]` | `ZERADO` | **12.15** AAA |
| Chip — selected, abandoned | `--z-state-abandoned` | `#C77DFF` | **177** | `bright magenta` | `⊘` U+2298 | `[x]` | `ABANDONED` | **7.21** AA |
| Chip — **unselected** | `--z-text-tertiary` | `#8492A8` | ***underived*** | `white` | the same glyph | the same | the same label | **6.15** AA |
| **Chip brackets** | `--z-border-strong` | `#64748B` | **67** | `bright black` | `[` `]` | `[` `]` | — | **4.08** — meets 1.4.11 |
| Focus marker on a chip / a row | `--z-primary` | `#FFB000` | **214** | `bright yellow` | `▌` U+258C | `>` | — | **10.59** AAA |
| Empty-result heading and body | `--z-text` | `#E9EEF5` | **255** | `bright white` | — | — | — | **16.65** AAA |
| Diagnostic facet names | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | — | — | — | **9.36** AAA |

> ***underived*** — `--z-text-tertiary` `#8492A8` has **no derived ANSI-256 index**, and this
> screen depends on it more than most: it is the *unselected chip*. **Interim: an unselected
> chip renders uncoloured** — no SGR — while a selected one renders in its state colour and bold
> **and is named in the ratio slot**. The distinction therefore survives the interim intact,
> because D-07-2 does not rely on colour at all. Owner: `fft-brand-architect`.

### 9.1 · Audio — the fourth channel

The annunciator is on `Z-04`'s **header band title row** (D-A1, `Z-04-library.md` §7.1). Filter
mode does not move it and does not draw it. **`m` is not listed in the footer while the editor
holds focus** — see §11.1, and that is not an omission, it is 2.1.4.

**`Z-07` has no sound event.** Nothing here completes, fails or arrives; every change is the
player's own keystroke, already visible in the ratio. **No sound may be attached to a keystroke
on this screen** — a filter that clicks is a filter that is louder than the library it is
searching.

---

## 10 · The full state table

| # | State | Trigger | Composition | Copy |
|---|---|---|---|---|
| **F1** | **First run — unreachable by construction** | Library empty | `/` on an empty library does nothing and shows no error. There is nothing to filter, and `Z-04`'s empty state owns the body. **`/` is not in that state's footer** | — |
| **F2** | **Just opened — the editor is empty** | `/` pressed | Editor shows `/ ` and the caret; the ratio reads `247 of 247`; the chips are all unselected; **the list is unchanged**. An empty query filters nothing | §11.4 |
| **F3** | **Typing — the editor has focus** | Any printable key | §3. Live re-filter on every keystroke. **Single-key shortcuts do not fire** | §11.1 |
| **F4** | **Editor blurred, filter applied** | First `Esc`, or `⏎`, or `Tab` | §4. The caret is gone; the filter is not | §11.2 |
| **F5** | **The chip row has focus** | `f`, or `Tab` from the editor | A `▌` marker sits left of the focused chip; `space` toggles it; `←` `→` move between chips | §11.3 |
| **F6** | **Nothing matches** | 0 rows | §5. The diagnostic block names each facet and its count. **The column header and scroll row are absent** | §11.5 |
| **F7** | **A state facet only, no query** | `f` then `space`, editor empty | The editor row shows `/ ` with no caret; the ratio reads `ZERADO · 6 of 247` | — |
| **F8** | **Filter active + a degrade banner** | Any of `Z-04` §8.1 | Banner takes body row 1; the bar moves to rows 2–3; the list drops to **10** | `Z-04` §8.1 |
| **F9** | **A sync changed the rows underneath** | Return from `Z-03` | The filter **re-runs** and the ratio updates. Cursor and offset are preserved by **game identity**; if the focused game no longer matches, focus moves to the nearest surviving match | — |
| **F10** | **A status changed while filtered** | `Z-06` applied, with a state facet selected | The row may **stop matching**. It is removed on the next frame and focus moves to the nearest surviving row. The ratio updates | §11.6 |
| **F11** | **The query matches everything** | e.g. a single space | Ratio reads `247 of 247`; the list is the whole library. Not an error, not a special case |  — |
| **F12** | **Cleared** | Second `Esc` | The mode ends. `Z-04`'s pinned summary returns and the list is whole again | — |
| **F13** | **Tiny** | `< 40` cols | One chip at a time, cycled with `Tab`; the bar stays 2 rows | §6 |
| **F14** | **Below the refusal floor** | `< 24` cols or `< 8` rows | The program refuses — `Z-04-library.md` §11.3 | — |
| **F15** | **The absent facet exists** | `absent_since IS NOT NULL` on ≥ 1 row (`06-data-seams.md` §2.4) | The fifth chip `[ABSENT]` renders and the chip row wraps; the bar goes to **3 rows** and the list to **10** (§3.4). The ratio's denominator becomes the whole file (D-07-9) | §11.5 |
| **F16** | **The absent facet is selected** | `space` on `[ABSENT]` | §5.1. The row set swaps from the shown rows to the absent ones; state chips then filter **within** that set. **It can never be empty on its own** — the chip does not exist when `n = 0` | §5.1 |
| **F17** | **The last absent row came back** | A sync returned it; `absent_since` cleared | The chip **stops rendering**, the bar returns to 2 rows and the list to 11. If the facet was selected, the mode falls back to the unfiltered shown set. **No notice, no banner, no result line** — §2.4 clears silently | — |

> **F10 is the state most likely to be missed and most likely to feel broken.** A player filters
> to `NOT STARTED`, presses `s`, marks a game `IN PROGRESS` — and the row they were looking at
> vanishes. It **should** vanish; the filter is still true. What must not happen is focus going
> nowhere. `04-navigation-and-focus.md` §4.1 rule 3 governs: focus moves to the nearest surviving
> item, below first, above if there is nothing below. And the *zerado* result line
> (`Z-04-library.md` §10.5) still prints in the pinned row — which in this mode is the editor
> row, so the ratio slot yields to it until the next keypress.

---

## 11 · The exact copy — ready to paste

### 11.1 · The footer while the EDITOR has focus

```
esc back to the list   ⏎ first match   tab state chips   ^c quit
```

**64 cells, 3-space separators.**

> **D-07-3 · While a text input holds focus the footer lists ONLY the keys that actually fire.**
> WCAG 2.1.4 and `03-designer-manual.md` §7 anti-pattern 5: single-key shortcuts must not fire
> while an input has focus — typing `q` types `q`, it does not quit. So `? help`, `q quit`,
> `s status`, `/ search`, `a add`, `r sync` and `m mute` are **all absent from this footer**,
> because listing a key that does not fire is exactly the lie `04-navigation-and-focus.md` §6
> forbids. `^c` **is** listed, because `Ctrl-C` genuinely does work inside a text input.
> *This is the single most valuable line in this spec: it is the footer telling the truth about
> the one state in which most TUIs lie.*

### 11.2 · The footer once the editor is BLURRED — the second `Esc` announced

```
esc clear filter  / search  f state  ↑↓ move  ⏎ open  ? help  q quit
```

**68 cells, 2-space separators** (the ladder in `Z-04-library.md` §9.2: 3 → 2 before any hint
drops; `m mute` is the first casualty when audio is enabled).

**The hint is literally true in both states, and it is different in each.** That is the whole
mechanism `04-navigation-and-focus.md` §5 calls *"discoverable rather than clever"*:

| Focus | Hint | What `Esc` actually does |
|---|---|---|
| the editor | `esc back to the list` | leaves the editor, **keeps** the filter |
| the list, filter applied | `esc clear filter` | **clears** the filter |

A single hint reading `esc clear` in both states would be wrong in one of them, and it is the
one the player is in while typing.

### 11.3 · The footer while the CHIP ROW has focus

```
space toggle  ←→ chip  esc back to the list  ⏎ open  ? help  q quit
```

**67 cells.** `space` is safe here because no text input has focus.

### 11.4 · F2 — just opened

The editor renders the sigil and the caret; the ratio reads the whole library:

```
/                                                                  247 of 247
```

**No placeholder text.** A greyed *"type to search…"* would be a fourth thing on a row that
already says what it is with an amber `/` and a live count, and it would have to be suppressed
the instant a key is pressed — which is a flicker for no information.

### 11.5 · F6 — nothing matches, and the filter that emptied it is named

```
Nothing matches.

  search   "souls"         3 games
  state    ZERADO          0 of those 3

247 games in the library.
```

**The rules that make this block honest:**

- **One line per active facet, in the order they were applied.** A facet that is not active does
  not appear.
- **Each line carries the count *at that step*.** `3 games` then `0 of those 3` tells the player
  exactly where the set collapsed — which is the difference between "nothing matches" and "your
  search was fine, the state filter emptied it."
- **The last line states the library total**, so an empty result never reads as data loss.
- Single-facet forms:

```
Nothing matches.

  search   "quake"         0 games

247 games in the library.
```

```
Nothing matches.

  state    ABANDONED       0 games

247 games in the library.
```

```
Nothing matches.

  source   physical        0 games

247 games in the library.
```

**With the absent facet.** A boolean facet renders its value as `yes` — it never renders `no`,
because an unselected facet does not appear at all:

```
Nothing matches.

  search   "souls"         3 games
  absent   yes             0 of those 3

250 games in the library. 3 are absent.
```

**Two rules that only bite once a row is absent:**

- **The last line states the whole file's total, not the default view's.** It already did — with
  no absent rows the two are the same number. Its job is that *"an empty result never reads as
  data loss"*, and the moment three rows are excluded from the default view, `247 games in the
  library.` would be the sentence doing the opposite of its job. The second sentence says the
  number, which is the brand's own instruction.
- **`absent   yes` alone can never produce this block.** The facet only exists when `n ≥ 1`, so
  selecting it and nothing else always returns `n` rows. Every empty result that names `absent`
  names a second facet beside it, and the *"count at that step"* rule shows exactly which one
  emptied the set.

**Not `"Something went wrong."`** and not `"No results."` — *"Name what happened, why, and the
next action."* The next action is in the footer, where actions live, and `esc clear filter` is
one keystroke away.

### 11.6 · F10 — a row stopped matching

No copy. The row leaves, focus moves to the nearest surviving row, the ratio decrements. If the
change was into *zerado*, the result line prints in the editor row's slot until the next
keypress:

```
/ the                                            Zerado. 9 hours. Sixth this year.
```

At 74 cells the query and the result line together may exceed the row; the **result line wins**
and the query is elided from the right with `…`. It is the more important thing for one beat,
and the query is still in the editor, unchanged, the moment the player presses a key.

### 11.7 · Copy notes

- **Casing** — `ZERADO` in a chip, in the ratio slot and in the diagnostic block, because all
  three are the interface. `247 games in the library.` is prose, so it is lowercase.
- **Say the number** — `23 of 247`, `3 games`, `0 of those 3`. Never *a few*.
- **No exclamation marks. No emoji. Never call the user a "gamer".**
- **No copy refers to a colour, a shape or a position** (WCAG 1.3.3). Never *"the cyan chip"*,
  never *"the row on the right"*. The diagnostic block names facets and the footer names keys.
- **Type-neutral where equally natural** — `Nothing matches.`, `247 games in the library.`,
  `search`, `state`, `source`. Phase 1 says *game* where a game is what the player is looking at.

---

## 12 · The focus model, and the two-step `Esc`

### 12.1 · The two regions

| Region | Takes focus by | Contains |
|---|---|---|
| **The filter bar** | `/` (editor) · `f` (chips) · `Tab` from the list | the query editor **and** the chip row — `Tab` moves between them |
| **The list** | `⏎`, `Esc` from the editor, `Tab` from the chips | the game rows |

**`R = 2`, which is why `Tab` exists here and nowhere else below 120 columns**
(`02-composition.md` §2). Traversal order matches visual order, top to bottom
(WCAG 2.4.3): **editor → chips → list → editor**.

### 12.2 · How focus is shown

| Element | Focused | Not focused |
|---|---|---|
| The query editor | the **caret** is present and blinking; the query is `--z-text` | no caret |
| A chip | `▌` marker in a 2-cell gutter to its left + **bold** + `--z-focus-ring` cyan on the bracket pair | no marker |
| A game row | `▌` in the row's gutter field + **bold** + `--z-primary` amber on the marker | two spaces |

**The chip's focus ring is `--z-focus-ring` cyan and is EXEMPT from the budget** — a focused
*control* is where the brand's ring belongs, and it is *"singular by definition — exactly one
element holds focus — so it cannot multiply"* (`02-colour-budget.md` §2.3). **It is never
removed**, in any state, for any reason.

Under `NO_COLOR` the caret, the `▌` marker and the bold weight all survive. **Three regions,
zero colour, still unambiguous.**

### 12.3 · The `Esc` table for this mode — the only two-step `Esc` in Zerado

| Context | `Esc` does | Then |
|---|---|---|
| **Editor focused** | Exit the **editor**; the filter **stays applied**; **what was typed is kept** | Focus returns to the list on the row it was on. The footer changes to `esc clear filter` |
| **Chip row focused** | Same — leave the bar, keep the filter | Footer changes |
| **Editor blurred, filter applied** | **Clear the filter.** The mode ends | The full library returns; `Z-04`'s pinned summary comes back; focus stays on the same **game** if it survives, else the nearest row |
| **An overlay (`Z-06`) is open over the mode** | Dismiss the overlay only | The filter and the editor state are untouched |

> *"The two-step Escape in filter mode is the only place `Esc` does something different on a
> second press, and it is deliberate: exiting a filter editor and clearing a filter are different
> intentions, and a player who has just typed a search almost never means to throw it away."*
> — `04-navigation-and-focus.md` §5.

**There is no keyboard trap** (WCAG 2.1.2): every focus position has an `Esc` that leaves, and
two presses always return the player to an unfiltered library.

### 12.4 · The rule that protects the player's data

**While the query editor holds focus, no single-key shortcut fires** (WCAG 2.1.4). Typing `d`
types `d`. Typing `a` types `a` — it does not open *Add a game by hand*. Typing `q` types `q` —
it does not quit and lose the session. Only `Esc`, `⏎`, `Tab`/`Shift-Tab`, the arrow keys,
`Backspace`, the editing keys and `Ctrl-C` are interpreted. *"This is the single most common way
a TUI destroys a user's data."*

---

## 13 · The key map

| Key | Does | Editor focused | Chips focused | List focused |
|---|---|---|---|---|
| printable keys | **type into the query** | ✓ literal text | — | — |
| `Backspace` / `Ctrl-W` / `Ctrl-U` | delete character / word / line | ✓ | — | — |
| `←` `→` | move the caret · move between chips | ✓ caret | ✓ chip | — |
| `space` | a space character · **toggle the focused chip** | ✓ literal | ✓ toggle — including `[ABSENT]` when it renders (§3.4) | — |
| `Tab` / `Shift-Tab` | next / previous region | ✓ | ✓ | ✓ |
| `⏎` | leave the editor and focus the **first match** | ✓ | ✓ toggle then leave | opens the game (`Z-05`) |
| `Esc` | §12.3 | ✓ | ✓ | ✓ |
| `↑` `↓` `k` `j` | move the row cursor | — | — | ✓ (`k` `j` are literal text in the editor) |
| `/` | re-focus the editor | — | ✓ | ✓ |
| `f` | focus the chip row | — | — | ✓ |
| `s` `a` `r` `,` `m` `?` `q` | the global and screen keys | **do not fire** | ✓ | ✓ |
| `Ctrl-C` | quit | ✓ | ✓ | ✓ |

**Reserved and inert:** `:` `Ctrl-K` (command palette), `1`–`9` (quick filters — which will one
day *be* this screen's shortcuts, and are unbound now so nobody has to be retrained),
`n` / `p`.

---

## 14 · Colour budget declaration

| State | STATE cyan (uncounted) | Focus ring (exempt) | **CHROME CYAN** | Verdict |
|---|---|---|---|---|
| F3 typing | every `◉` / `ZERADO` in the results, plus a selected `[◉  ZERADO]` chip | the **text caret** — exempt by §2.3 | **0** | **PASS** |
| F4 blurred | as F3 | none | **0** | **PASS** |
| F5 chips focused | as F3 | the focused chip's `--z-focus-ring` bracket pair — exempt, singular by definition | **0** | **PASS** |
| F6 nothing matches | none — there are no rows | the caret | **0** | **PASS** |

**`Z-07` spends ZERO chrome cyan in every state.** The mode sigil `/` is **amber** (allow-list
item 8), the ratio is chrome, the diagnostic block is chrome, and the `esc clear filter` hint is
chrome. **A cyan `/` or a cyan ratio would be the classic failure** — cyan used for emphasis
*"because it's the brand colour"* (`03-designer-manual.md` §7 anti-pattern 6).

**The selected `[◉  ZERADO]` chip is STATE cyan and is explicitly on the uncounted list**
(`02-colour-budget.md` §2.1, fourth bullet: *"the `[ ◉ ZERADO ]` filter chip when selected"*).

**Amber allow-list entries used:** 3 the `IN PROGRESS` state (in results and as a chip) ·
**8 the filter mode sigil `/`** · plus the row and chip focus markers. **Not used:** 1 the screen
title is `Z-04`'s · 2 readout labels are the column header, `--z-text-secondary` · 4 progress ·
5 key hints — the footer here is chrome · 6 the terminal mark · 7 the banner, unless F8.

**Amber ceiling:** at 80 × 24 = 1920 cells the ceiling is **192**; §3's render spends 7
(`LIBRARY`) + 1 (`/`) + 16 (the `IN PROGRESS` chip) + 2 × 14 (`IN PROGRESS` result rows) + 1
(cursor) ≈ **53**.

**Red: none, in any state.** **An empty filter result is not an error and is never red** — it is
the filter working. **Every control boundary is `--z-border-strong`**: the chip brackets at
4.08 satisfy WCAG 1.4.11. **A `--z-border` bracket around a chip would be an automatic fail.**
Nothing on this screen is separated by fill.

---

## 15 · `NO_COLOR` — rendered, not asserted

Zero SGR sequences. §3, character for character, with `[◉  ZERADO]` **selected**:

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   / the▎                                                   ZERADO · 3 of 247   │
│   [○  NOT STARTED]  [◐  IN PROGRESS]  [◉  ZERADO]  [⊘  ABANDONED]              │
│                                                                                │
│     STATE           TITLE                                        HOURS   SRC   │
│     ◉  ZERADO       Bastion                                         6h   STM   │
│   ▌ ◉  ZERADO       Return of the Obra Dinn                         9h   STM   │
│     ◉  ZERADO       The Messenger                                  12h   STM   │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│                                                                                │
│   ROWS  1–3 of 3                                                               │
│   esc back to the list   ⏎ first match   tab state chips   ^c quit             │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

| Information | Carried without colour by |
|---|---|
| That a filter mode is active | The `/` sigil, the ratio, and the `esc` hint in the footer — three text channels |
| **Which chips are selected** | **`ZERADO ·` in the ratio slot** — the word channel of D-07-2. Without it, `NO_COLOR` would leave only bold, and bold alone across four adjacent controls is not legible |
| Which region has focus | The caret (editor) · the `▌` marker (chip or row) · bold |
| The states in the results | glyph **and** label, unchanged |
| What `Esc` will do next | The footer's wording, which changes with focus |
| How many matched | The ratio, in words and numerals |

**This screen is the reason D-07-2 exists.** Strip the colour from a chip row and the selection
is the only thing that could be lost — so it is carried in a word, in a fixed place, on the same
row.

**`ZERADO_ASCII=1`:**

```
/ the                                                       ZERADO · 3 of 247
[[ ]  NOT STARTED]  [[~]  IN PROGRESS]  [[*]  ZERADO]  [[x]  ABANDONED]
```

> **FLAG** · The ASCII state column is itself bracketed, so a bracketed chip becomes `[[*]  ZERADO]` —
> readable, but doubled. **Under `ZERADO_ASCII=1` the chip's outer brackets become `< >`:**
> `<[*]  ZERADO>`. Recorded as part of **D-06-5** (`Z-06-set-status.md` §13.1), which extends
> the ASCII escape hatch beyond the state column; **`<` and `>` are ASCII-narrow and immune.**

---

## 16 · Reuse verdict, per element

| Element | Verdict | Why |
|---|---|---|
| **The query editor** | **`bubbles/textinput` — direct fit, reuse it** | `01-design-system.md` §7.7. It brings the caret, the editing keys, `Ctrl-W`/`Ctrl-U`, and — critically — the **focus discipline that makes 2.1.4 correct by construction**. Restyled to Zerado tokens |
| The state chips | Build fresh · `lipgloss` | Four fixed bracketed cells with the ratified chip geometry. `bubbles/list` brings pagination and a status bar a four-item horizontal row has no use for |
| The ratio slot | Build fresh · `lipgloss` right-align + width-aware truncation | |
| **The result list** | **The shared Zerado ledger primitive** — the same one `Z-04` uses | `01-design-system.md` §4.9. Filtering changes the row set, not the renderer. A second table implementation here is exactly the defect `LedgerTable` was built to cure |
| Empty-result block | Build fresh · `lipgloss` text layout | |
| The diagnostic block | Build fresh · two padded columns | |
| Debounce | **None.** Re-filter on every keystroke | It is a `WHERE` clause over a few hundred rows already in memory. A debounce would add latency to the one thing on this screen that must feel instant |
| Fuzzy matching | **Not in Phase 1** | Nothing in the canon specifies a matcher. Phase 1 is a **case-insensitive substring match on the title** — see §19, D-07-5 |

---

## 17 · Upstream findings

**Re-checked against head on 2026-08-25.** Finding 7 of rev A — `03-designer-manual.md` §5.11
verdict 3 reading as a permanent rejection of audio — is **closed**: the row is struck through and
marked SUPERSEDED with a provenance note (`14-contradictions-closed.md` #4). It is struck from
this table. Findings 1–6 were re-read at source and all six still stand.

| # | Finding | Where | Owner |
|---|---|---|---|
| 1 | **`01-design-system.md` §7.1 says the chips render persistently only at ExtraWide and otherwise appear *"via `f`"*, while `03-responsive.md` §3 budgets the bar at 2 rows from Standard upward and 3 at Narrow.** Two rows means editor + chips. This spec reconciles them: **the chip row is always present in the mode; `f` focuses it rather than creating it** | `01-design-system.md` §7.1 vs `03-responsive.md` §3 | `fft-design-architect` / `fft-tui-architect` |
| 2 | **`01-design-system.md` §7.2 draws `esc clear` in the filter bar itself**, which is wrong while the editor has focus — the first `Esc` does **not** clear (`04-navigation-and-focus.md` §5). The hint is stateful and belongs in the footer, where hints live | `01-design-system.md` §7.2 | `fft-design-architect` |
| 3 | **`01-design-system.md` §7.3 makes a selected chip colour + bold only**, which under `NO_COLOR` is bold alone — a one-channel co-render on a control. Closed by **D-07-2** with a word channel and no new glyph | `01-design-system.md` §7.3 | `fft-design-architect` |
| 4 | **`01-design-system.md` §10.2's empty-result copy does not name the filter that emptied it**, which `01-screen-inventory.md` §5 requires of this screen. §11.5 is the composition that satisfies both | `01-design-system.md` §10.2 | `fft-design-architect` |
| 5 | **No document specifies the match semantics.** Substring? Prefix? Fuzzy? Case sensitivity? Diacritics? This spec decides the minimum honest answer (D-07-5) and flags it | — | `fft-tui-architect` |
| 6 | **`05-state-machine.md` §7 still draws `filter: source=physical`, but no key is bound to a *source* facet** in `04-navigation-and-focus.md` §3, and `07-offline-contract.md` §2 no longer names one. This spec renders the facet if the model supplies it and binds no key — see §18 item 2. **It is the exact problem D-07-8 exists to avoid for `absent`**: a facet with no reachable control is a promise with no mechanism | `05-state-machine.md` §7 | `fft-tui-architect` |
| **7** | **`03-responsive.md` §3 budgets this bar at 2 rows from Standard up, and the four-chip row is 63 cells — which does not fit Standard's 54-cell body at `n = 0`, before any absent row exists.** Counting `06-data-seams.md` §2.4's fifth facet is what surfaced it: the absent chip takes the bar to **3 rows at Wide and ExtraWide (10 and 27 rows)** and **4 at Narrow (7 rows)**, and at Standard there is no `n = 0` figure to add a row to. The composition follows `Z-04-library.md` D-04-2 and closes at 16 everywhere it is defined; **Standard is undefined in both documents and this spec says so rather than inventing it** (§3.4). §2.4 postdates §3, but the Standard hole predates both | `03-responsive.md` §3 · this spec §1 | `fft-tui-architect` |

---

## 18 · Open for the founder

1. **The match semantics** (D-07-5). Phase 1 is specced as a **case-insensitive, accent-folded
   substring match on the title**. Accent folding matters more than usual for a
   Brazilian-by-origin product: a player typing `pokemon` should find `Pokémon`. **Confirm**, or
   nominate a fuzzy matcher — but note that a fuzzy matcher makes the empty-result diagnostic
   (§11.5) much harder to write honestly, because *"3 games matched"* stops being a fact the
   player can reproduce.
2. **Is there a `source` filter in Phase 1?** `05-state-machine.md` §7 still draws
   `filter: source=physical`; `07-offline-contract.md` §2 no longer names one, and no key is
   bound. The diagnostic block and the ratio slot both **support** the facet; nothing
   **reaches** it. Either bind a key (route to `fft-tui-architect`) or confirm it is Phase 2 —
   and if it stays, **D-07-8's answer applies to it too**: a sixth chip needs no new key, at the
   cost of the row arithmetic in §3.4.
3. **The result line versus the query, when both want the same row** (§11.6). When a filtered
   row becomes *zerado*, the pinned row must show both what the player typed and what just
   happened. This spec gives the beat to the result line and elides the query. **Confirm** — the
   alternative is to suppress the result line while filtering, which would mean the product's
   payoff moment is quietest exactly when a player is hunting through a backlog.

---

## 19 · Design decisions made in this spec

| # | Decision | Reason |
|---|---|---|
| **D-07-1** | In filter mode the pinned summary is the **ratio** `<matched> of <total>`, not the four state counts (§3.3) | A ratio describes a subset honestly and in one line; a breakdown of a filtered set is either readable-anyway or tautological. R-10(c) still holds |
| **D-07-2** | A selected state facet is **named in words** in the ratio slot (§3.3) | Colour + bold alone fails `NO_COLOR` on a control. The word channel closes it with no new glyph and no new component |
| **D-07-3** | While the editor has focus the footer lists **only** the keys that actually fire, and `^c` is one of them (§11.1) | WCAG 2.1.4 plus `04-navigation-and-focus.md` §6 — a footer that lists a key that does not fire is the lie both rules forbid |
| **D-07-4** | The chip row is **always present** in the mode; `f` focuses it rather than creating it (§3.1) | Reconciles `01-design-system.md` §7.1 with `03-responsive.md` §3's 2-row budget, and means the bar never changes height when a key is pressed |
| **D-07-5** | Phase 1 matching is a **case-insensitive, accent-folded substring match on the title**; no fuzzy matcher, no debounce | The minimum a player can predict and reproduce, which is what makes the empty-result diagnostic honest. Accent folding because the product is Brazilian by origin and `Pokémon` must be findable by typing `pokemon` |
| **D-07-6** | The mode takes rows **1–2** of the pinned block and the list yields, never the respiro (§3.1) | The same rule as `Z-04`'s D-04-2, so the two compositions cannot drift apart |
| **D-07-7** | On zero results the column header and the scroll-position row are **absent**, but the bar and the respiro stay (§5) | There is nothing to head and nothing to position; the player still needs to see and edit the thing that emptied the list |
| **D-07-8** | The absent facet is a **fifth chip** in the existing chip row, carrying **no glyph**, rendering only when `n ≥ 1`, and wrapping to its own row (§3.4, §5.1) | It is the only mechanism that needs **no new key binding**, so `06-data-seams.md` §2.4's *"remain findable by filter"* is deliverable rather than blocked behind a spine decision — the fate of the `source` facet in §17. No glyph because `absent` is not a fifth state and a fifth glyph beside `○ ◐ ◉ ⊘` would assert that it is; the missing glyph is the channel that says so. Five chips plus the focus marker compose to 75 cells in a 74-cell body, so they wrap rather than cram |
| **D-07-9** | Once any row is absent, the ratio's denominator is the **whole library file** (§3.3) | With `n = 0` the two numbers are identical and nothing changes. With `n ≥ 1`, `247 of 250` is the same fact `Z-04`'s summary carries as `3 absent`, said in the vocabulary this row already uses — and it stops the two pinned rows disagreeing about how big the library is |

---

## 20 · Screen-specific acceptance criteria

Beyond `00-design-brief.md` §10 and `02-colour-budget.md` §10, both of which apply in full.

1. **Typing `q` types `q`.** Verified for `q`, `a`, `s`, `r`, `f`, `/`, `,`, `m` and `?` while
   the editor holds focus. None of them fires its shortcut.
2. **The footer lists no key that does not fire in the current focus state**, in all three focus
   states.
3. **The first `Esc` keeps the filter and keeps what was typed; the second clears it.** Verified
   from both the editor and the chip row.
4. **The footer's `esc` hint is literally true in each state** and changes at the moment focus
   changes.
5. **An empty result names every active facet and its count at that step**, and states the
   library total.
6. **The chip row's selection is legible with `NO_COLOR=1`** — verified by stripping SGR and
   reading only the characters.
7. **The list is byte-identical to `Z-04`'s list** for the same rows — same columns, same
   widths, same cursor. Filtering changes the row set, not the renderer.
8. **The visible game count is 11 at 80 × 24** with the bar active, 10 with a banner as well,
   6 at 40 × 24, 8 at 32 × 24.
9. **The ratio updates on every keystroke** with no debounce and no perceptible lag at 412 rows.
10. **A row that stops matching is removed and focus moves to the nearest surviving row** — never
    nowhere, never to the top.
11. **Returning from a sync re-runs the filter** and preserves the cursor by game identity.
12. **Chrome-cyan count is 0** in every state; the only cyan is state cyan and the focus ring.
13. **Every chip bracket is `--z-border-strong`.** A `--z-border` bracket is an automatic fail.
14. **The empty result is not red and raises no banner.**
15. **At Tiny one chip renders and `Tab` cycles it**, and the bar is still 2 rows.
16. **The `[ABSENT]` chip renders if and only if at least one row is absent**, at every tier —
    and when it renders, the bar is 3 rows and the list is 10 at 80 × 24, 2 chip rows and 27 at
    120 × 40, 4 rows and 7 at 40 × 24, and unchanged at 32 × 24. **Standard is excluded from
    this line until §17 finding 7 is answered** — it has no agreed `n = 0` figure to test against.
17. **Selecting `[ABSENT]` alone never returns zero rows.** Verified with `n` = 1, 3 and 40.
18. **The absent facet's selection is legible with `NO_COLOR=1`** — the chip's word, its
    brackets and the ratio slot's `ABSENT` all survive SGR-stripping. **No fifth glyph appears
    in the chip row or in any ledger row**, at any tier, in any state.
19. **An absent row's ledger row is byte-identical to a shown row of the same state** — same
    chip, same columns, same widths. Absence changes the row *set*, not the row.
20. **`Z-04` and `Z-07` agree on the library's size** in the same render: `Z-04`'s summary says
    `3 absent` and `Z-07`'s slot says `of 250`.
21. **Founder-validated screenshot before merge**, at the six viewports of
    `03-responsive.md` §7 plus `NO_COLOR=1` and forced-16-colour at 80 × 24 — **including the
    zero-result state and a 412-row library filtered to a 380-row result scrolled to the end.**
    No screenshot → not GOLDEN → no merge.
