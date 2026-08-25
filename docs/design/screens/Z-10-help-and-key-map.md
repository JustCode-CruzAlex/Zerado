---
title: Zerado — Z-10 Help and key map
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-10
rev: B
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-10 · Help and key map

> *"Every key that does anything, on the screen it does it on."* (`01-screen-inventory.md` §2.)
> And its inverse, which is the whole design of this screen:
> *"`Z-10` must not list keys that do nothing on the screen the player came from."* (§5.)
> **A help screen that lies is worse than no help screen.**

**Canon that governs this screen:** `00-design-brief.md` §10 · §3.1 (SC 3.2.6 Consistent Help,
1.3.2 Meaningful Sequence, 2.4.6 Headings and Labels) · `01-design-system.md` §2 (band) ·
`01-design-system.md` §1.5 (readout role) · `02-colour-budget.md` §10 · `03-designer-manual.md` §3 ·
`02-composition.md` §2 (single-pane table in a viewport) · `03-responsive.md` §3 ·
**`04-navigation-and-focus.md` §1 rule 3 (unwind, never stack) · §3 (the global key table) ·
§3.1 (reserved and unbound) · §5 (`Esc`) · §6 (the footer)** ·
`07-offline-contract.md` §2 — *"Help (`Z-10`): **WORKS**. It is compiled in."*

---

## 1 · Identity

| | |
|---|---|
| **ID** | `Z-10` |
| **Name** | Help and key map |
| **Phase** | 1 |
| **Kind** | **Route** |
| **Route in** | `?` from any screen where `?` is bound — which is every screen except `Z-11 Fatal error`, and except while an overlay is open at Standard and above (`Z-06-set-status.md` §11, D-06-4) |
| **Route out** | `Esc` → back to the origin, focus restored exactly where it was |
| **Composition** | **Single-pane table in a viewport.** One key column at Wide and below; **two** at ExtraWide |
| **`R`** | **1** — the viewport. No `Tab` |
| **Offline class** | **WORKS.** It is compiled into the binary; it reaches for nothing |
| **Displayed name** | breadcrumb `Zerado ✦ Help` · title `HELP` |

**Pushing `?` while `Z-10` is already on the stack unwinds to it rather than duplicating it**
(`04-navigation-and-focus.md` §1 rule 3), so `?` never builds a tower of help screens. On this
screen `?` therefore does nothing visible — and §11.1 makes that honest rather than mysterious.

---

## 2 · Purpose

**One sentence:** the keys that work **on the screen the player just came from**, then the keys
that work everywhere, then the keys that are reserved and do nothing yet — so that pressing one
and getting silence is an answered question rather than a bug.

---

## 3 · The design — the key map is GENERATED, not written

> **D-10-1 · Every block on this screen is generated from the origin screen's own key registry.
> There is no hand-maintained list of keys anywhere in the product.**
>
> This is the only design that can satisfy `01-screen-inventory.md` §5 permanently. A
> hand-written help table is correct on the day it is written and wrong on the day the next key
> is bound — and nothing fails when it drifts, which is why help screens rot. If the same
> registry that **dispatches** a key is the one that **describes** it, then a key that is not
> bound cannot be listed and a key that is bound cannot be omitted. The same registry also
> composes the **footer** of every screen (`04-navigation-and-focus.md` §6), so the footer and
> the help can never disagree either.
>
> **Three consequences a builder must honour:**
> 1. A binding carries its own one-line description, its display form and its **scope**
>    (`screen` · `mode` · `global` · `reserved`).
> 2. A binding may declare an **availability predicate** — `s` on `Z-04` is live only when a row
>    exists, `c` only when the library is empty, `m` only when audio has been enabled, **`v`
>    only when the deck is reachable, and `x` only while the capability note is showing**.
>    `Z-10` asks the origin, in its current state, and lists what is live.
>    `04-navigation-and-focus.md` §5b names this the registry's own mechanism and §5c uses it
>    for `?`. **Rev B added no mechanism — `v` and `x` arrived and the screen already knew what
>    to do with them**, which is the property D-10-1 was bought for.
> 2b. A binding's **description may differ by state**, and the registry expresses that as **two
>    entries with mutually exclusive predicates**, never as one entry with a conditional string
>    (**D-10-8**). `v` is the first: `see the covers` in list mode, `back to the list` in cover
>    mode, exactly matching the two footer hints `v covers` and `v list`
>    (`Z-15-cover-deck.md` §10.2). One entry with a string chosen at render time would put the
>    footer's wording and the help's wording back in two places — the drift D-10-1 abolishes.
>    **Exactly one of the two is ever live, so block 1 gains exactly one row, never two.**
> 3. **Nothing on this screen is a string literal in `Z-10`'s own code** except the four block
>    headings.

### 3.1 · The four blocks, in this order

| # | Heading | Contains | Sorted by |
|---|---|---|---|
| **1** | `ON THIS SCREEN — <ORIGIN>` | Every binding whose scope is `screen`, for the origin **in the state it is in** | Movement, then activation, then verbs — the order a player learns them |
| **2** | `IN <MODE>` | Present only when the origin has an active mode with its own bindings — today, `Z-07` | as declared |
| **3** | `EVERYWHERE` | Every binding whose scope is `global` and is live on the origin | `?` `,` `m` `esc` `q` `^c` |
| **4** | `RESERVED — NOT BOUND YET` | Every `reserved` binding, each stamped with the phase that will claim it | as declared |

**Block 4 is not padding.** `04-navigation-and-focus.md` §3.1 reserves `:` `Ctrl-K` `1`–`9`
`n` `p` and binds them to nothing: *"a key that means nothing in Phase 1 and something in Phase 2
is a feature arriving; a key that means one thing in Phase 1 and another in Phase 2 is a betrayal
of muscle memory."* A player who presses `:` and gets silence has a question. This block is the
answer, and it is the cheapest possible way to keep the promise.

---

## 4 · Mockup — 80 × 24, pushed from a populated Library with a filter applied

`tier = Wide` · `leftInset = 3` · **body = 74 × 16** · content begins at **column 4**.
15 content rows + 1 position row. **34 lines in total, so it scrolls.**

**Rev B: block 1 gained `v`.** [`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md)
§5b's registry gained the cover-deck toggle ([`Z-15-cover-deck.md`](./Z-15-cover-deck.md) §10),
so this screen gained a row without a line of `Z-10` changing — which is D-10-1 working. Every
drawn count below moved with it.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Help                                                                │
│                                                                                │
│   HELP                                                                         │
│                                                                                │
│                                                                                │
│   ON THIS SCREEN — LIBRARY                                                     │
│                                                                                │
│   ↑ ↓  k j      move the cursor                                                │
│   g   G         first row · last row                                           │
│   ^d  ^u        half a page down · up                                          │
│   ⏎             open this game                                                 │
│   s             set this game's status                                         │
│   /             filter and search                                              │
│   f             filter by state                                                │
│   a             add a game by hand                                             │
│   r             sync with Steam                                                │
│   v             see the covers                                                 │
│                                                                                │
│   IN FILTER MODE                                                               │
│                                                                                │
│   ROWS  1–15 of 34                                                             │
│   ↑↓ scroll   esc back   q quit                                                │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Character counts:** body **74 × 16** · key field **12** + gutter **2**, so every description
begins at body column **15** · the widest description row is 50 cells
(`?             this key map — you are looking at it`) · position row **16** · footer **29**.
*(Rev A printed the position row as 18; measured width-aware with Ambiguous = 1 it is 16, and
`–` U+2013 is the only Ambiguous character in it.)*

> **The viewport ends on an orphaned heading, and that is correct — D-10-10.** With `v` present,
> `IN FILTER MODE` is document line 14 and its respiro is line 15, so the last two visible rows
> are a heading and a blank. **A scroll region must never re-flow its cut point to avoid this.**
> The moment the viewport pages by anything other than one line per `↓`, the position readout
> stops describing the document and starts describing a pagination the player cannot predict.
> The heading is not orphaned in the document; it is orphaned in a *window*, and the window is
> labelled — `ROWS  1–15 of 34`, with `↑↓ scroll` in the footer.

### 4.0 · Every count on this screen, and what moves it

The document's length is not a property of `Z-10`. It is the number of live bindings on the
origin, which is why every figure here is derived rather than declared.

| Origin, in the state it is in | Blocks | Lines | Drawn at |
|---|---|---|---|
| Populated Library, audio enabled | 1 · 3 · 4 | **27** | §13.1 minus block 2 |
| …with a filter applied, editor blurred | 1 · 2 · 3 · 4 | **34** | **§4** · §13.1 |
| …with the capability note showing | 1 (+ `x`) · 2 · 3 · 4 | **35** | §13.1's addendum |
| …audio never enabled | as above, minus `m` | **one less** | — |
| Empty library | 1 · 3 · 4, `v` not bound | **20** | §6 · §13.2 |
| From `Z-06` at Tiny, short forms | 1 · 3 · 4 | **21** | §9 |

**None of these numbers may be hard-coded.** They are here so a reviewer can check the drawn
goldens, and every one of them is an assertion the registry can be made to prove (§24 criterion 2).

### 4.1 · Row map — 80 × 24

| Terminal row | Body row | Content |
|---|---|---|
| 1 | — | blank · `OuterMarginY` = 1 |
| 2–4 | — | header band: breadcrumb · `InterElementGap` · title |
| 5–6 | — | blank · `InnerPaddingY` = 1, then `InterElementGap` = 1 |
| **7–21** | **1–15** | **the viewport — 15 lines of a 33-line document** |
| 22 | 16 | the position readout, **pinned outside the viewport** |
| 23 | — | footer, the reserved frame row |
| 24 | — | blank · `OuterMarginY` = 1 |

### 4.2 · The key row — anatomy

```
↑ ↓  k j      move the cursor
│             │
│             └ description field — flex, --z-text, sentence case, begins at column 15
└ key field — 12 cols, --z-primary amber, the key exactly as it is pressed
```

12 is the widest key display form in the product (`esc esc`, `1 – 9`, `:   ^k`, `↑ ↓  k j`).
**The field is width-aware padded**, not counted — `↑` U+2191 and `↓` U+2193 are both
**Ambiguous** width (verified, UCD 16.0.0), as are `–` U+2013 and `—` U+2014, while `⏎` U+23CE
is Neutral. Padding to 12 measured cells keeps the description edge straight on either kind of
terminal.

**Keys are written the way they are pressed.** `^d`, not `Ctrl+D`; `esc`, not `Escape`;
`↑ ↓  k j` on one row because they are one action with four spellings. **`esc esc` is written as
two words** because it is genuinely two presses, and it is the only entry in the product like
that.

### 4.3 · The position readout

```
ROWS  1–15 of 34
```

Pinned outside the viewport, exactly as `Z-04`'s is — the same component, the same wording, the
same place. It is what tells a player there is more below.

> **D-10-2 · Reuse the ledger's position readout rather than invent a scroll indicator.**
> The design system has no scroll-indicator component (`01-design-system.md` §15) and this spec
> composes, it does not invent. `ROWS  n–m of N` already exists on `Z-04`, says the number
> (brand §8), survives `NO_COLOR` and costs one row. A graphical bar would be a new component
> and a second way of saying the same thing.

---

## 5 · Mockup — scrolled to the end

Lines 20–34 of 34. The global and reserved blocks. **Every visible row is byte-identical to
rev A's** — `v` landed in block 1, fifteen lines above the fold, so the only thing that moved
here is the readout. That is worth seeing: a registry change does not repaint a screen, it
repaints the part of the document the change is in.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Help                                                                │
│                                                                                │
│   HELP                                                                         │
│                                                                                │
│                                                                                │
│                                                                                │
│   EVERYWHERE                                                                   │
│                                                                                │
│   ?             this key map — you are looking at it                           │
│   ,             settings                                                       │
│   m             mute or unmute the audio                                       │
│   esc           back one level                                                 │
│   q             quit                                                           │
│   ^c            quit, from anywhere, always                                    │
│                                                                                │
│   RESERVED — NOT BOUND YET                                                     │
│                                                                                │
│   :   ^k        the command palette — Phase 2                                  │
│   1 – 9         quick filters — Phase 2                                        │
│   n   p         next · previous game — Phase 2                                 │
│   ROWS  20–34 of 34                                                            │
│   ↑↓ scroll   esc back   q quit                                                │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

> **`?  this key map — you are looking at it`.** Pressing `?` here unwinds to a screen already
> on the stack and therefore does nothing visible. Listing it as *"this key map"* with no
> qualifier would be the screen lying about itself on the one screen whose entire job is not
> lying. **The clause is generated from the binding's precondition** (D-10-1 consequence 2), the
> same availability predicate that keeps `c` out of a populated library's list and `s` out of an
> empty one's.

> **`m` is present only when audio has been enabled.** When it has never been enabled the
> binding is not live, so `Z-10` does not list it and the footer does not carry it — because
> there is nothing to mute. This falls out of D-10-1; it is not a special case.

> **And `x` is the same mechanism, which answers the question rev B was opened to answer —
> D-10-9.** `x dismiss` ([`Z-15-cover-deck.md`](./Z-15-cover-deck.md) §5.4) is bound **only while
> the capability note is showing**. That is an availability predicate in exactly the sense
> D-10-1 consequence 2 means it, so **`Z-10` excludes `x` from the default render precisely the
> way it excludes `m` when audio has never been enabled** — same field on the same registry
> entry, no second mechanism, no special case in this screen's code. **The default document is
> therefore 34 lines, not 35.**
>
> **The state where it does render is real and must be drawn as such.** The note lives on body
> row 1 of `Z-04`; it is not an overlay and it does not take focus
> ([`Z-15-cover-deck.md`](./Z-15-cover-deck.md) §5.4), so it does not make `?` inert —
> [`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §5c names only a
> text input and an open overlay. A player on Terminal.app can press `v`, see the note, and
> press `?`. In that state block 1 carries **both** `v` and `x` and the document is **35** lines
> (§13.1's addendum, state **H13**).
>
> **The alternative — listing `x` always — is the defect this screen exists to prevent.** It
> would put a key in the map that does nothing on nine renders out of ten, and the player who
> pressed it would get the silence `Z-10` was built to abolish. **The alternative in the other
> direction — never listing it — is the same defect wearing the opposite coat**: a bound key
> omitted, which D-10-1 forbids in the same sentence. The predicate is the only answer that is
> right in both states, and it is already in the registry.

---

## 6 · Mockup — pushed from an EMPTY Library: a different, honest list

**This is the requirement, demonstrated.** The player came from the same screen ID, `Z-04`, in a
different state. **Eight of its nine screen bindings are not live, so eight of them are not
here — and one binding appears that the populated screen does not have.** The registry both drops
and adds, which is the whole point of an availability predicate.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Help                                                                │
│                                                                                │
│   HELP                                                                         │
│                                                                                │
│                                                                                │
│   ON THIS SCREEN — LIBRARY                                                     │
│                                                                                │
│   Nothing is in the library yet, so most keys do nothing here.                 │
│                                                                                │
│   c             connect a store                                                │
│   a             add a game by hand                                             │
│                                                                                │
│   EVERYWHERE                                                                   │
│                                                                                │
│   ?             this key map — you are looking at it                           │
│   ,             settings                                                       │
│   esc           back one level                                                 │
│   q             quit                                                           │
│   ^c            quit, from anywhere, always                                    │
│                                                                                │
│   ROWS  1–15 of 20                                                             │
│   ↑↓ scroll   esc back   q quit                                                │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

> **This render did not move in rev B, and that is the demonstration.** The same amendment that
> took the populated document from 33 lines to 34 leaves this one at **20**, because `v` is not
> bound on an empty library ([`Z-15-cover-deck.md`](./Z-15-cover-deck.md) §10.1 — *"library empty
> → no, and not in the footer"*; C12/C13 there). **There is no deck of nothing.** One registry
> change, two origins, two different right answers, and neither of them written by hand.

**Note what is gone:** `↑ ↓`, `g`/`G`, `^d`/`^u`, `⏎`, **`s`**, `/`, `f`, `r`, **`v`** — every
one of them needs a row, and there are none. **`s` is gone with the rest**: it means *set this game's
status* everywhere, always (`04-navigation-and-focus.md` §3.2), and there is no game.
**Note what appeared:** `c  connect a store`, which is **not** in the populated list — it is live
only here. **`m` is gone** because audio has never been enabled in this render.

**The prose line is not decoration.** A player who opens help on an empty library and sees two
keys should be told *why* it is two, or the screen reads as broken. One sentence, dry, concrete:

```
Nothing is in the library yet, so most keys do nothing here.
```

> **The collision this row used to flag is closed.** Rev A recorded `s` meaning *set status* on
> a populated library and *connect Steam* on an empty one, and routed it to the founder.
> `04-navigation-and-focus.md` §3.2 resolved it — **`s` is set-status everywhere, always;
> connect-a-store is `c`** — and cites this screen as part of the argument: the collision *"forces
> `Z-10 Help` to print different text depending on how full the library is."* It no longer does.
> `Z-10` rendered it correctly either way, because it asks the origin rather than assuming; it
> now renders something simpler.

---

## 7 · Mockup — ExtraWide, 120 × 40: two key columns

`leftInset = 4` · **body = 112 × 32** · two columns of **52** with an **8**-column gutter:
`4 + 52 + 8 + 52 + 4 = 120`, the two 4s being `OuterMarginX` 2 + `InnerPaddingX` 2 on each side.
The whole map fits without scrolling — **21 of the body's 32 rows are used.**

**Every description here is §13.1's long form**, because the 38-cell description field holds all
of them (`this key map — you are looking at it` is 36). **Three rev-A errors are corrected in
this render** and each is the kind only a redraw finds: the frame was **2 columns short** (118
inside, for a 120-column terminal); `f  filter by state` was **missing from block 1** and a
second `f` appeared in block 2, which is not where
[`Z-07-filter-and-search.md`](./Z-07-filter-and-search.md) puts it; and the footer carried
`↑↓ scroll` **on the one render that does not scroll** — the *"footer that lies"* this screen
exists to make impossible. The footer is now `esc back   m mute   q quit`, **26 cells** (§11.1).

```
┌────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                                        │
│    Zerado ✦ Help                                                                                                       │
│                                                                                                                        │
│    HELP                                                                                                                │
│                                                                                                                        │
│                                                                                                                        │
│    ON THIS SCREEN — LIBRARY                                    IN FILTER MODE                                          │
│                                                                                                                        │
│    ↑ ↓  k j      move the cursor                               esc           leave the editor, keep the filter         │
│    g   G         first row · last row                          esc esc       clear the filter                          │
│    ^d  ^u        half a page down · up                         tab           move between editor, chips and list       │
│    ⏎             open this game                                space         toggle the chip you are on                │
│    s             set this game's status                                                                                │
│    /             filter and search                                                                                     │
│    f             filter by state                                                                                       │
│    a             add a game by hand                                                                                    │
│    r             sync with Steam                                                                                       │
│    v             see the covers                                                                                        │
│                                                                                                                        │
│    EVERYWHERE                                                  RESERVED — NOT BOUND YET                                │
│                                                                                                                        │
│    ?             this key map — you are looking at it          :   ^k        the command palette — Phase 2             │
│    ,             settings                                      1 – 9         quick filters — Phase 2                   │
│    m             mute or unmute the audio                      n   p         next · previous game — Phase 2            │
│    esc           back one level                                                                                        │
│    q             quit                                                                                                  │
│    ^c            quit, from anywhere, always                                                                           │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│                                                                                                                        │
│    esc back   m mute   q quit                                                                                          │
│                                                                                                                        │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**The position row is absent** because nothing is scrolled — `ROWS 1–21 of 21` would be
furniture. It appears the moment the content exceeds the viewport. **The `↑↓ scroll` hint goes
with it**, which is why this tier's footer is 26 cells and not 38.

> **The count that would appear is 21, not 34.** At ExtraWide the document is not a 34-line
> stream; it is the same four blocks laid into two columns, and the taller column decides the
> height. A position readout counts **rendered rows**, never document lines — which is the same
> rule [`Z-15-cover-deck.md`](./Z-15-cover-deck.md) §11.5 states for `COVERS` over a grid: the
> label word and the unit change with the composition, because a number that describes nothing
> on screen is worse than no number.

> **D-10-3 · Two columns are safe here under WCAG 1.3.2, and here is the argument.**
> *"A screen's reading order for anything consuming the output stream is byte order. Two panes
> side by side interleave."* (`00-design-brief.md` §3.1.) That is a real hazard and it is why
> `Z-05`'s pane content is a single column. It does **not** bite here, for two reasons:
> **(1)** every cell of this table is a self-contained `key + description` pair, so byte order
> yields `key desc key desc` — interleaved, but never a sentence split across a gutter;
> **(2)** the **section headings span the row**, so the grouping — *on this screen* versus
> *everywhere* versus *reserved* — survives the stream unbroken.
> **The rule this fixes for the future:** a second column is permitted only when every cell in
> it is independently meaningful. A prose column beside a table would fail.

---

## 8 · The origin table — what block 1 says, for every screen

`Z-10` is the same screen eleven times over. This table is the **contract**; each origin's key
list is owned by that origin's spec, not by this one.

| Origin | Block-1 heading | Where its bindings are specified |
|---|---|---|
| `Z-01` First run | `ON THIS SCREEN — FIRST RUN` | `Z-01-first-run.md` §"key map" |
| `Z-02` Connect a store | `ON THIS SCREEN — CONNECT A STORE` | `Z-02-…md` §"key map" |
| `Z-03` Sync | `ON THIS SCREEN — SYNC` | `Z-03-…md` §"key map" |
| **`Z-04` Library, populated** | `ON THIS SCREEN — LIBRARY` | `Z-04-library.md` §9 — **§4 above** |
| **`Z-04` Library, empty** | `ON THIS SCREEN — LIBRARY` + the prose line | `Z-04-library.md` §9 — **§6 above** |
| **`Z-05` Game detail** | `ON THIS SCREEN — GAME DETAIL` | `Z-05-game-detail.md` §14 |
| **`Z-06` Set status** | `ON THIS SCREEN — SET STATUS` | `Z-06-set-status.md` §11 — **Tiny only**, see below |
| **`Z-07` Filter mode** | `ON THIS SCREEN — LIBRARY` **plus block 2** `IN FILTER MODE` | `Z-07-filter-and-search.md` §14 |
| **`Z-15` Cover deck** | `ON THIS SCREEN — LIBRARY` — **block 1's contents change; there is no block 2** | `Z-15-cover-deck.md` §10 |
| `Z-08` Add a game by hand | `ON THIS SCREEN — ADD A GAME BY HAND` | `Z-08-…md` §"key map" |
| `Z-09` Settings | `ON THIS SCREEN — SETTINGS` | `Z-09-…md` §"key map" |
| `Z-11` Fatal error | **unreachable** | `?` is not bound there; only `q` and `Ctrl-C` work |

**Two origins need explaining:**

- **`Z-06`** is an overlay at Narrow and above, and `?` is **inert** while an overlay is open
  (`Z-06-set-status.md` §11, D-06-4) — the overlay prints its own three keys in the footer, which
  is help already on screen. **At Tiny `Z-06` is a route, `?` works normally, and `Esc` from
  `Z-10` returns to it.** So `ON THIS SCREEN — SET STATUS` renders at Tiny and nowhere else.
  §9 draws it.
- **`Z-07`** is a *mode* of `Z-04`, not a screen. Block 1 heading stays `LIBRARY` — because that
  is where the player is — and block 2 appears with the mode's own bindings. **If `?` is pressed
  while the query editor holds focus it types nothing and does nothing**, because single-key
  shortcuts do not fire in a text input (WCAG 2.1.4); the player must leave the editor first,
  and `Z-07`'s footer says how.
- **`Z-15`** is also a mode of `Z-04` — and it gets **no block 2**, which is the distinction
  **D-10-11** draws. Cover mode **swaps the body renderer**, so `↑ ↓` no longer means *move the
  cursor*, it means *move a row of covers*, and `← → h l` become live. Those are not extra keys
  on top of the list's; they are the list's keys **meaning something else**. Block 1 therefore
  renders the deck's descriptions in place of the list's, and `v` reads `back to the list`. The
  heading is still `LIBRARY`, because that is where the player is.

> **D-10-11 · A mode earns block 2 only when it is *additive*.** `Z-07` adds an editor and a
> chip row **on top of** the list: every list key still works and still means what it meant, so
> the mode's own keys are genuinely extra and belong in a block of their own. `Z-15` **replaces**
> the list, so listing `↑ ↓  move the cursor` in block 1 and `↑ ↓  move a row of covers` in
> block 2 would hand the player two contradictory descriptions of one key on one screen. **A
> substitutive mode rewrites block 1; an additive mode adds block 2.** This is the rule for
> every mode Phase 2 brings, and it is decided here because `v` is the first key that forced it.

---

## 9 · Mockup — Tiny, 32 × 24, pushed from `Z-06` as a route

The one composition where the origin is `Set status`. Body **30 × 21**, `leftInset` **1**, both
outer margins **0**, band collapsed to the title row. Key field shortens to **8**.

```
┌────────────────────────────────┐
│ HELP                           │
│                                │
│ ON THIS SCREEN —               │
│ SET STATUS                     │
│                                │
│ ↑↓ k j    move                 │
│ g  G      first · last         │
│ ⏎         apply                │
│ esc       cancel               │
│                                │
│ EVERYWHERE                     │
│                                │
│ ?         this key map         │
│ ,         settings             │
│ esc       back one level       │
│ q         quit                 │
│ ^c        quit, always         │
│                                │
│ RESERVED                       │
│                                │
│ :  ^k     palette, P2          │
│ 1 – 9     filters, P2          │
│ ROWS  1–20 of 21               │
│ ↑↓ scroll  esc back            │
└────────────────────────────────┘
```

**What shortens at Tiny, and nothing else:** the block-1 heading wraps to two lines rather than
truncating; the key field goes 12 → **8**; descriptions shorten to their **declared short form**
(`quit, always`, `palette, P2`, `first · last`); the heading `RESERVED — NOT BOUND YET` shortens
to `RESERVED`.

**`v` and `x` are absent from this render and from every Tiny render**, and neither is a
shortening: `v` is **unbound below 40 columns** (`Z-15-cover-deck.md` **D-15-8**), so the note it
raises cannot exist there either, so `x` cannot either. Nothing is printed about any of it —
which is `Z-15` §12.2's silence, honoured here rather than re-argued.

> **Rev B corrected this golden twice, and the two corrections are the same defect.**
> **`g  G  first · last` was missing.** `Z-06` §11 binds `g` / `G` to *first / last item*, and
> **a bound key that is not listed is the exact failure this screen exists to prevent** —
> D-10-1's second half, caught in `Z-10`'s own drawing. It is restored, which makes the document
> **21** lines. **The readout said `of 22`.** It is now `ROWS  1–20 of 21`, and the arithmetic
> closes for the first time: body **30 × 21**, viewport **20**, position row pinned at body row
> 21, document 21 — so it overflows by exactly one line, `n  p  paging, P2` sits below the fold,
> and the readout and the `↑↓ scroll` hint are both **earned**. At 20 lines they were furniture,
> and §12 **H7** would have required both to disappear.

> **D-10-4 · Every binding carries a long description and a short one; the short one is used
> below Standard.** It is declared at the binding, beside the long one, so it can never be a
> mechanical truncation — `quit, from anyw…` is exactly the kind of string that makes a help
> screen useless at the width it matters most. Neither form is ever truncated at render time.

---

## 10 · Visual hierarchy

| # | What | Channel |
|---|---|---|
| **1** | **`ON THIS SCREEN — LIBRARY`** | UPPERCASE readout heading, `--z-primary` amber, first line of the body. It answers *"is this about where I was?"* before a key is read |
| **2** | **The key column** | A straight vertical block of amber at a fixed 12-cell field — the eye scans it like an index, which is what it is |
| **3** | **The descriptions** | `--z-text`, sentence case, on a straight edge at column 15. The only sentence-case text on screen, so it reads as prose against the keys |
| **4** | **The block headings** | UPPERCASE amber, each preceded by `InterElementGap`. Three ranks — *here* · *everywhere* · *not yet* — carried by position and by the words |
| **5** | Position readout and footer | `--z-text-secondary` chrome |

**The one thing the player should see first is the origin's name in the heading.** They pressed
`?` from somewhere and the first question is whether this screen knows where. **The second is
their key**, which is why the key column is the amber one and the descriptions are not.

**Hierarchy is carried by case, weight, colour role and spacing — in that order — and by no
box drawing at all.** There are **no rules, no boxes and no separators** on this screen: three
blocks, one `InterElementGap` between them, and a straight column edge. A key map is a table
that does not need a table.

---

## 11 · The key map — of the help screen itself

Short by design, and every line true.

| Key | Does | Note |
|---|---|---|
| `↑` `↓` `k` `j` | Scroll the viewport | Listed **only when there is overflow** |
| `^d` / `^u` | Half a page down / up | Same |
| `g` / `G` | Top / bottom | Same |
| `Esc` | **Back to the origin** | Focus restored exactly where it was |
| `,` | Settings → `Z-09` | Global |
| `m` | Mute / unmute | Global; **only when audio has been enabled** |
| `q` | Quit | Immediate |
| `Ctrl-C` | Quit | Always |
| `?` | **Nothing visible** — it unwinds to a route already on the stack | **Not in this screen's footer.** In the `EVERYWHERE` block it is described as `this key map — you are looking at it` |

### 11.1 · Footers, exact

| Tier / state | Footer | Cells |
|---|---|---|
| Wide, scrolling | `↑↓ scroll   esc back   q quit` | **29** |
| Wide, audio enabled | `↑↓ scroll   esc back   m mute   q quit` | **38** |
| Tiny | `↑↓ scroll  esc back` | **19** |
| No overflow, audio never enabled | `esc back   q quit` | 17 |
| **No overflow, audio enabled — §7's ExtraWide render** | **`esc back   m mute   q quit`** | **26** ≤ 112 |

**Rev B struck the row that read `ExtraWide → ↑↓ scroll   esc back   m mute   q quit`.** At
ExtraWide the map fits, so nothing scrolls, so `↑↓ scroll` is a key that does nothing there —
and the rule two paragraphs down is that keys that do nothing here are not listed. The tier does
not decide the footer; **overflow does**, and the two are not the same question.

**`? help` is absent from every one of them.** It is the smallest possible demonstration of
`04-navigation-and-focus.md` §6 — *"keys that do nothing here are not listed"* — applied to the
screen that would be most embarrassing to get wrong. Separator 3 spaces, tightening to 2 before
any hint drops.

---

## 12 · The full state table

| # | State | Trigger | Composition | Copy |
|---|---|---|---|---|
| **H1** | **First run** — pushed from an empty `Z-04`, or from `Z-01` | Library empty | §6. Block 1 lists only the bindings that are live, with the prose line explaining why there are two | §13.2 |
| **H2** | **From a populated Library** | The default | Blocks 1, 3, 4. Block 2 is **absent** — it needs an *active* mode (§3.1). **27 lines** | §13.1 minus block 2 |
| **H3** | **From filter mode** | `?` pressed with the filter applied and the **editor blurred** | Block 1 = `LIBRARY`, block 2 = `IN FILTER MODE`. **34 lines** — **§4 draws this one** | §13.1 |
| **H3c** | **From cover mode** | `?` pressed after `v` | Block 1 = `LIBRARY` with the **deck's** descriptions; **no block 2** (**D-10-11**); `v` reads `back to the list` | §8 |
| **H3b** | **`?` pressed while the query editor has focus** | — | **Nothing happens.** `?` is a literal character the editor does not consume as a shortcut (WCAG 2.1.4). `Z-10` is not pushed | — |
| **H4** | **From `Z-05`** | `?` on the detail view, either host | Block 1 = `GAME DETAIL` | — |
| **H5** | **From `Z-06` at Tiny** | `?` on the set-status **route** | §9 | — |
| **H5b** | **From `Z-06` as an overlay** | — | **Unreachable.** `?` is inert while an overlay is open; the overlay's own footer is the help | — |
| **H6** | **Content taller than the viewport** | Wide and below, always | `bubbles/viewport` scrolls; the position row says where you are | §13.3 |
| **H7** | **Content fits** | ExtraWide | The position row is **absent** and the footer drops `↑↓ scroll` | — |
| **H8** | **Audio never enabled** | The default | The `m` row is absent from the `EVERYWHERE` block and from the footer | — |
| **H9** | **Offline** | Network off | **Identical.** *"It is compiled in."* No banner, no notice, no difference — `Z-10` is `WORKS` | — |
| **H10** | **From `Z-11 Fatal error`** | — | **Unreachable.** `?` is not bound there | — |
| **H11** | **Tiny** | `< 40` cols | §9 — short descriptions, 8-cell key field, wrapped heading | §13.4 |
| **H12** | **Below the refusal floor** | `< 24` cols or `< 8` rows | The program refuses — `Z-04-library.md` §11.3 | — |
| **H13** | **The capability note is showing** | `v` pressed on a terminal that draws no images, note not yet dismissed (`Z-15` §5) | Block 1 carries **both `v` and `x`**. **35 lines.** `?` is not inert — the note is not an overlay and takes no focus | §13.1's addendum |
| **H14** | **The note has been dismissed on a terminal that draws no images** | `covers.note_dismissed` | **`v` is unbound, so it is not listed and neither is `x`** — back to **27** (or 34 in filter mode). See §21 finding 3: this is where `Z-15` §5.4 expects a fact that D-10-1 cannot carry | — |
| **H15** | **`v` is bound below 40 columns** | — | **Unreachable.** `v` is unbound at Tiny (**D-15-8**), so neither it nor `x` can appear in any Tiny render | — |

> **H9 is worth stating explicitly** even though it is a non-event. `07-offline-contract.md` §2
> classifies `Z-10` as **WORKS** and the class rule is *"Nothing. No banner, no notice, no
> difference."* A help screen that grew an `OFFLINE` banner would be a banner appearing when
> nothing is degraded — furniture, and a violation of §3's own rule.

---

## 13 · The exact copy — ready to paste

### 13.1 · The `Z-04` key map, long form

```
ON THIS SCREEN — LIBRARY

↑ ↓  k j      move the cursor
g   G         first row · last row
^d  ^u        half a page down · up
⏎             open this game
s             set this game's status
/             filter and search
f             filter by state
a             add a game by hand
r             sync with Steam
v             see the covers

IN FILTER MODE

esc           leave the editor, keep the filter
esc esc       clear the filter
tab           move between editor, chips and list
space         toggle the chip you are on

EVERYWHERE

?             this key map — you are looking at it
,             settings
m             mute or unmute the audio
esc           back one level
q             quit
^c            quit, from anywhere, always

RESERVED — NOT BOUND YET

:   ^k        the command palette — Phase 2
1 – 9         quick filters — Phase 2
n   p         next · previous game — Phase 2
```

**34 lines.** At Wide the viewport shows 15 of them. **Without block 2 — the plain populated
Library, H2 — it is 27**; block 2 costs seven lines (its heading, its respiro, its four keys and
the gap that separates it from block 3).

**When the capability note is showing (H13), block 1 carries one more row and the document is
35 lines.** Only the tail of block 1 changes; nothing else moves:

```
r             sync with Steam
v             see the covers
x             dismiss the note about cover art
```

**`v`'s description does not change on a terminal that cannot draw images**, and that is
deliberate: pressing `v` there is still how a player asks for the covers, and the note is
Zerado's answer — *"the note answers; it never announces"* (`Z-15-cover-deck.md` **D-15-6**). A
description that hedged with *"if your terminal can"* would put a fact about the terminal in the
one column that is supposed to say what a key **does**, and it would say it on every terminal
that does not need it.

### 13.2 · The empty-library variant

```
ON THIS SCREEN — LIBRARY

Nothing is in the library yet, so most keys do nothing here.

c             connect a store
a             add a game by hand
```

### 13.3 · The position readout

```
ROWS  1–15 of 34
ROWS  20–34 of 34
ROWS  1–20 of 21
```

Absent when nothing is scrolled.

### 13.4 · Short descriptions — the Tiny forms, declared not truncated

| Key | Long | **Short** |
|---|---|---|
| `^c` | `quit, from anywhere, always` | `quit, always` |
| `:   ^k` | `the command palette — Phase 2` | `palette, P2` |
| `1 – 9` | `quick filters — Phase 2` | `filters, P2` |
| `n   p` | `next · previous game — Phase 2` | `paging, P2` |
| `?` | `this key map — you are looking at it` | `this key map` |
| `m` | `mute or unmute the audio` | `mute` |
| `esc esc` | `clear the filter` | `clear filter` |
| `space` | `toggle the chip you are on` | `toggle chip` |
| **`g` / `G`** | `first row · last row` (`Z-04`) · `first item · last item` (`Z-06`) | **`first · last`** |
| **`x`** | `dismiss the note about cover art` | **`dismiss the note`** |
| heading | `RESERVED — NOT BOUND YET` | `RESERVED` |

**Two notes on the rev-B rows.** `g` / `G` shows the mechanism plainly: **the long description
belongs to the origin's binding, not to this screen**, so `Z-04` and `Z-06` declare different
ones and share a short form. And `x`'s short form is used at **Narrow**, not at Tiny — its long
form is 32 cells against Narrow's 24-cell description field, and `x` cannot occur at Tiny at all
(§9). **`v` declares no short form and needs none**: `see the covers` is 14 cells and fits every
field at every tier where `v` is bound.

### 13.5 · Copy notes

- **Descriptions are verb-first and lowercase** — `move the cursor`, not `Moves the cursor` or
  `Cursor movement`. They complete the sentence *"press this to…"*, which is what a player is
  actually asking.
- **Say what it does, not what it is.** `filter and search`, not `opens the filter and search
  screen`. The player does not care that it is a route.
- **Phase-stamped, never vague.** `the command palette — Phase 2`, never *"coming soon"*.
  Brand §8: *"Never claim what isn't built. Unshipped phases are marked as phases."*
- **No copy refers to a colour, a shape or a position** (WCAG 1.3.3). Never *"the cyan chip"*,
  never *"the row on the right"*.
- **No exclamation marks. No emoji. Never call the user a "gamer".**
- **Casing** — `Zerado` in the breadcrumb. No state chip appears on this screen, so `ZERADO`
  does not appear at all; the filter-mode row says `filter by state`, not `filter by ZERADO`.
- **Zerado is a games product and the copy says so.** `open this game`,
  `set this game's status`, `add a game by hand`, `sync with Steam`. `move the cursor`,
  `first row · last row` and `ROWS` are neutral because each names a **movement or a unit** and
  would read the same in any list — not because the copy is leaving room for another kind of
  thing (`11-media-model.md` rev B).

---

## 14 · Every applied spacing token, by name

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
| **Key field** | **8** | 10 | 12 | **12** | **12** |
| **Description field** | 22 | 24 | 40 | **60** | **38** per column |
| **Columns** | 1 | 1 | 1 | **1** | **2** |
| Column gutter | — | — | — | — | **8** |

**Applied:**

| Surface | Token | Value at Wide |
|---|---|---|
| Content left edge | `leftInset` | column 4 — the same column as the breadcrumb and the title |
| Between blocks | `InterElementGap` | 1 row, every time |
| heading → first key row | `InterElementGap` | 1 row |
| key field → description | a 2-column gutter | fixed, so the description edge is straight |
| Between the two columns (ExtraWide) | an 8-column gutter | wide enough that a long description in the left column cannot be misread as belonging to the right |
| Position row | pinned outside the viewport | 1 row, and absent when nothing scrolls |
| Footer | the canon's reserved footer row | 1 row |

**Zero magic numbers.** The key field's 12 (8 at Tiny) is the widest display form in the product
and is stated as such; the ExtraWide column split is `(112 − 8) / 2 = 52`.

---

## 15 · Colour, glyph and label

**There is no state chip on this screen**, so the four-state table does not apply here — the
only screen in the bundle where that is true, and it is worth saying so rather than pasting the
table for form's sake.

| Element | Token | Hex | ANSI-256 | 16-colour | Ratio |
|---|---|---|---|---|---|
| Screen title `HELP` | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA |
| Block headings | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA |
| **The key column** | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA |
| Descriptions | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** AAA |
| The empty-library prose line | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** AAA |
| The `— Phase 2` stamp | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA |
| Breadcrumb, position row, footer | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA |
| Reserved-block key column | `--z-text-tertiary` | `#8492A8` | ***underived*** | `white` | **6.15** AA |

> **The reserved keys are tertiary, not amber, and that is the one colour decision on this
> screen.** Amber is *"the ambient voice — the colour the machine speaks in"*. A key that does
> nothing is not the machine speaking; it is the machine listing what it will say later.
> Rendering `:` in the same amber as `⏎` would tell the eye they are equals, which is the one
> thing this block exists to deny. **Interim: tertiary renders uncoloured** — no derived index —
> which still reads as subordinate against amber, so the interim is not a downgrade here.

### 15.1 · Audio — the fourth channel

The annunciator is right-aligned on this screen's own **title row**
(D-A1, `Z-04-library.md` §7.1) — `HELP` at the left, `▮ AUDIO` or `▯ MUTED` at the right, both
**Neutral**-width glyphs verified against UCD 16.0.0.

**`Z-10` has no sound event.** It reaches for nothing, completes nothing and fails at nothing.
**No sound may be attached to opening or closing help.**

---

## 16 · The focus model, and `Esc`

### 16.1 · Focus

**One region: the viewport.** There is no `Tab`, at any tier — including ExtraWide, where two
*columns* are one region, not two panes. Nothing here is activatable, so there is no item cursor
and no `▌` marker: **the viewport itself is the focused thing**, and the position readout is what
tells the player where they are inside it.

> **D-10-5 · No per-row cursor on this screen.** A key map is read, not operated. A cursor would
> imply `⏎` does something, and it does not — the keys are pressed *elsewhere*, which is the
> entire premise of the screen. `04-navigation-and-focus.md` §4.1 rule 2 is satisfied honestly:
> the region has no selectable items, so focus is on the region.

**The focus ring is not removed** — it has nothing to be on. The screen has exactly one region
and it always has focus, which is the degenerate but correct case.

### 16.2 · `Esc`

| Context | `Esc` does | Then |
|---|---|---|
| `Z-10` on top of the stack | **Pop the route** | The origin, with **its focus restored exactly where it was** — the same ledger row, the same detail view, the same filter with the same query |
| `Z-10` reached from `Z-06` at Tiny | Pop | Back to `Z-06`, with the same item focused |

**`Esc` always leaves. There is no keyboard trap** (WCAG 2.1.2). And because `?` unwinds rather
than stacks (`04-navigation-and-focus.md` §1 rule 3), a player can never end up two help screens
deep and have to press `Esc` twice.

### 16.3 · Consistent Help — WCAG 3.2.6, and the two honest exceptions

`?` opens this screen **from every screen, in the same way, showing the same three blocks in the
same order.** Two places it does not, both recorded rather than hidden:

1. **While a text input has focus** (`Z-07`'s query editor, `Z-02` and `Z-08`'s form fields).
   `?` is a literal character there; single-key shortcuts must not fire (WCAG 2.1.4), and 2.1.4
   is the stronger obligation. Each of those screens' footers names the way out.
2. **While an overlay is open** at Narrow and above (`Z-06`). The overlay prints its own three
   keys in its footer, which is help already on screen and one keystroke closer than a route
   would be. See `Z-06-set-status.md` §18 item 1 — it is open for the founder.

---

## 17 · 40-column behaviour, and the refusal floor

**40 × 24, body 36 × 16.** One column, key field **10**, description field **24**, short
descriptions. `03-responsive.md` §3: *"Narrow: same, Tiny: same, scrolls more."* Blocks and
their order are unchanged; only the field widths and the description forms move.

**32 × 24** is §9: key field **8**, headings wrap rather than truncate.

**Never sheds, at any width:** the block headings and their order · the origin's name · every
live binding · the reserved block · the position readout when there is overflow · the footer.

**Nothing on this screen is ever truncated.** A key map that cuts a description in half has
failed at its only job, exactly as an empty state that gets cut off has
(`01-design-system.md` §10.4). Descriptions have declared short forms (D-10-4); headings wrap.

**The refusal floor is the program's:** below **24 columns or 8 rows**, one frameless sentence
and, at start-up, exit `2`. A running session keeps running.

---

## 18 · `NO_COLOR` — rendered, not asserted

Zero SGR sequences. §4, character for character:

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Help                                                                │
│                                                                                │
│   HELP                                                                         │
│                                                                                │
│                                                                                │
│   ON THIS SCREEN — LIBRARY                                                     │
│                                                                                │
│   ↑ ↓  k j      move the cursor                                                │
│   g   G         first row · last row                                           │
│   ^d  ^u        half a page down · up                                          │
│   ⏎             open this game                                                 │
│   s             set this game's status                                         │
│   /             filter and search                                              │
│   f             filter by state                                                │
│   a             add a game by hand                                             │
│   r             sync with Steam                                                │
│   v             see the covers                                                 │
│                                                                                │
│   IN FILTER MODE                                                               │
│                                                                                │
│   ROWS  1–15 of 34                                                             │
│   ↑↓ scroll   esc back   q quit                                                │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

| Information | Carried without colour by |
|---|---|
| Which block is which | The **UPPERCASE headings**, and the `InterElementGap` above each |
| Key versus description | **Case** — `^d  ^u` against `half a page down · up` — plus the fixed column edge at 15 |
| That the reserved keys are different | The heading `RESERVED — NOT BOUND YET`, the `— Phase 2` stamp on every row, and **position** — they are last |
| Where you are in the document | `ROWS  1–15 of 34` |
| That `?` does nothing here | The words `you are looking at it` |
| Audio | `▮ AUDIO` / `▯ MUTED` — glyph **and** word |

**This screen is almost entirely `NO_COLOR`-native by construction**, because its content is
words about keys. The one thing colour was carrying — the reserved block's subordination — is
carried three other ways as well.

**`ZERADO_ASCII=1`** changes the key column's display forms and nothing else:

```
up/dn k j     move the cursor
enter         open this game
:   ^k        the command palette — Phase 2
1 - 9         quick filters — Phase 2
```

`↑ ↓` become `up/dn`, `⏎` becomes `enter`, `–` becomes `-`, `—` becomes `-`, `·` becomes `,`.
All ASCII-narrow and immune. Part of **D-06-5** (`Z-06-set-status.md` §13.1).

---

## 19 · Colour budget declaration

| State | STATE cyan | Focus ring | **CHROME CYAN** | Verdict |
|---|---|---|---|---|
| H2 from a populated Library | **none — there is no state chip on this screen** | none — nothing here is a control | **0** | **PASS** |
| H1 from an empty Library | none | none | **0** | **PASS** |
| H5 from `Z-06` at Tiny | none | none | **0** | **PASS** |
| H7 ExtraWide | none | none | **0** | **PASS** |

**`Z-10` spends ZERO chrome cyan and contains ZERO cyan of any class.** It is the only screen in
the bundle with no cyan on it at all, and that is correct: a key map has nothing complete and
nothing to urge. **A cyan `?` row, a cyan heading or a cyan `esc back` hint would each be an
automatic fail** — cyan is never a heading, never a title, never emphasis in body text.

**Amber allow-list entries used:** 1 the screen title · **2 readout labels and section heads** —
the three block headings · and the key column, which is the *"key hints"* entry (item 5) applied
at scale. **Not used:** 3 the `IN PROGRESS` state — no chip renders here · 4 progress fill ·
6 the terminal mark · 7 the degrade banner — this screen is `WORKS` and never shows one ·
8 the filter sigil.

**Amber ceiling:** at 80 × 24 = 1920 cells the ceiling is **192**. §4's render spends 4 (`HELP`)
+ 24 + 14 (two headings) + **26** (the key column — the keys themselves, never the padding that
straightens the field) = **68**. *Rev A estimated the key column at "~40" and the figure is now
counted rather than guessed; adding `v` moved it by exactly 1 cell.* The key column is the
largest single amber spend on any screen in this bundle and it is still nowhere near the ceiling
— but it is
the screen to re-measure if a block is ever added.

**Red: none.** No scanner, no annunciator, no error text. **`--z-border` is not used at all** —
there are no boxes, no rules and no separators on this screen. **No region is separated by
fill**, and at ExtraWide the two columns are separated by an 8-column gutter, which survives
every colour depth because it is made of nothing.

---

## 20 · Reuse verdict, per element

| Element | Verdict | Why |
|---|---|---|
| Scroll body | **`bubbles/viewport` — direct fit** | 34 lines into 15. This is exactly what the primitive is for |
| Key rows | Build fresh · `lipgloss` + the width-aware pad | Two padded columns. **`bubbles/table` does not fit** — it is the primitive that independently dropped the title, the scroll and the pinned footer on two FlowForge screens, and a key map needs none of its selection or column machinery |
| **The key registry** | **Build fresh, once, and make it the single source for dispatch, the footer and this screen** | D-10-1. This is the reuse that matters on this screen: not a component, a **source of truth** |
| Two-column layout (ExtraWide) | `lipgloss` join | |
| Position readout | **The same component `Z-04` uses** | Same wording, same place, same behaviour (3.2.4 Consistent Identification) |
| Markdown rendering | **Not `glamour`** | The key map is a generated table, not a document. `glamour` carries its own palette and would have to be restyled to render something that was never markdown |
| Header band | The shared `Frame` | Enforced at the router, so this screen cannot render frameless |
| Error rendering | **Not `charmbracelet/log`** | Developer logging role |

---

## 21 · Upstream findings

**Re-checked against head on 2026-08-25.** Four of rev A's six findings are **closed** and are
struck from this table. Finding 1 — nothing specified where a key's description lives — landed as
`04-navigation-and-focus.md` **§5b, *"One key registry — dispatch, footer and help are the same
source"***, which names the entry's key, scope, description, **availability predicate** and
handler, and states that *"`Z-10` is generated from it"* (`14-contradictions-closed.md` #18).
**D-10-1 is now a spine requirement rather than this screen's proposal.** Findings 3 and 4 — how
help describes itself, and the two places `?` is unavailable — are closed by **§5c**, which
tabulates both *"so the conformance statement can be written once and honestly"*. Finding 6's
audio half is closed (`03-designer-manual.md` §5.11, struck through and marked SUPERSEDED, #4)
and its second half too: **`m` is in the global key table**, marked *"only when audio is
enabled."* And the `s` collision is closed by decision — `04-navigation-and-focus.md` §3.2 binds connect-a-store to `c`, so
`Z-10` no longer prints different text depending on how full the library is (§6).

| # | Finding | Where | Owner |
|---|---|---|---|
| 1 | **`04-navigation-and-focus.md` §3 still lists no `f`.** `01-design-system.md` §7.3 binds it to the state filter chips (*"ExtraWide, or via `f`"*) and `Z-07` composes with it; this spec lists it under `Z-04` as `filter by state`. §5b now settles *where* a description lives but not *which bindings exist*, and **a key that is not in the registry cannot be listed here** — which is D-10-1 working as designed, and is exactly why the omission has to be fixed upstream rather than papered over in this screen | `04-navigation-and-focus.md` §3 | `fft-tui-architect` |
| **2** | **`Z-07`'s absent facet is a control with no key**, by design (`Z-07-filter-and-search.md` D-07-8): `[ABSENT]` is toggled with `space` on the chip row, which is already bound. It therefore has **no row in this screen's block 1 or block 2**, and a player who has absent rows will find the facet on `Z-07` and not in Help. That is correct — Help lists *keys*, and this facet is not one — but it is the first Phase 1 affordance that Help cannot describe, and it is worth knowing before someone reads the omission as a bug | `06-data-seams.md` §2.4 · `Z-07` D-07-8 | `fft-tui-architect` |
| **3** | **`Z-15` §5.4 names this screen as the durable home for a fact D-10-1 forbids it to carry.** Its dismissal table reads *"Nothing is lost — the fact stays available in `Z-10 Help`, which describes what `v` does and which terminals draw covers"*, and `Z-15` §5.3 puts *"the protocol name in the two durable homes — `Z-10 Help`'s description of `v`, and `Z-09 Settings § DISPLAY`"*. But `Z-15` §10.1 unbinds `v` once the note is dismissed — and **a key that is not bound cannot be listed** (D-10-1). So at exactly the moment `Z-15` says the fact survives here, this screen has stopped mentioning `v` at all (state **H14**). **The consequence is that `Z-15` §17 finding 6 — the `Images` row of `Z-09 § DISPLAY` — is not optional; it is the only durable home left**, and `Z-15` §5.4's reassurance should point there instead. `Z-10` will not grow a fifth block for retired keys: block 4 is *reserved*, meaning *not bound yet*, and a key the player has explicitly finished with is not that | `Z-15-cover-deck.md` §5.3 · §5.4 · §17 finding 6 · `Z-09-settings.md` §10.4 | `fft-tui-designer` — **`Z-15` next pass**, then founder on the `Images` row |
| **4** | **`Z-15` §10 gives `g` / `G` and `Ctrl-D` / `Ctrl-U` the scope `global`**; this screen's §3.1 puts movement in **block 1** (`screen` scope) and its `EVERYWHERE` block is declared as exactly `?` `,` `m` `esc` `q` `^c`. Both cannot be true of one registry field. **`screen` is the right answer** — *first* and *last* mean a row on `Z-04`, a tile on `Z-15` and an item on `Z-06`, so the description is owned by the origin, which is what `screen` scope means. The renders in both specs already behave this way; only the label in `Z-15`'s table is wrong | `Z-15-cover-deck.md` §10 | `fft-tui-designer` — **`Z-15` next pass** |
| **5** | **`Z-04` §4's ExtraWide frame is 2 columns short.** Its `120 × 40` render draws a 120-cell frame line, which is **118** columns of terminal; `Z-15` §4's equivalent is 122 and correct. This screen's §7 carried the identical defect and rev B fixed it — which is why it is worth reporting rather than assuming it is a house convention. The internal-consistency checker cannot see it, because every line in the block agrees with every other. Not fixed here: `Z-04` is not this pass's to edit | `Z-04-library.md` §4 (the `120 × 40` render, frame line 148) | `fft-tui-designer` — **`Z-04` next pass** |

---

## 22 · Open for the founder

> **Two items closed since rev A, both adopted rather than answered.** The **key registry**
> (item 1) is now `04-navigation-and-focus.md` §5b, a spine requirement carrying the key, its
> scope, its description, its **availability predicate** and its handler, with *"`Z-10` is
> generated from it"* stated outright. **`?`'s two unavailable contexts** (item 2) are tabulated
> in §5c — a text input (WCAG 2.1.4) and an open overlay — *"so the conformance statement can be
> written once and honestly."* §16.3 is that statement and it now has an upstream to cite.

1. **The reserved block.** Listing keys that deliberately do nothing is unusual, and it is here
   because `04-navigation-and-focus.md` §3.1 reserves five of them and a player who presses one
   deserves an answer rather than silence. **Confirm** — the alternative is that `:` and `1`–`9`
   are simply undocumented until Phase 2.

---

## 23 · Design decisions made in this spec

| # | Decision | Reason |
|---|---|---|
| **D-10-1** | The key map is **generated from the origin screen's key registry**, which is also what dispatches keys and composes footers (§3) | The only design in which help cannot drift. A hand-written table is correct the day it is written and silently wrong afterwards — which is why help screens rot |
| **D-10-2** | Reuse `Z-04`'s `ROWS n–m of N` readout rather than invent a scroll indicator (§4.3) | The design system has no such component and this spec composes; the readout already says the number and survives `NO_COLOR` |
| **D-10-3** | Two columns at ExtraWide are safe under WCAG 1.3.2 **because every cell is self-contained and the headings span the row** (§7) | Byte order yields `key desc key desc`, never a sentence split across a gutter. States the rule for future two-column screens |
| **D-10-4** | Every binding declares a **long and a short** description; the short one is used below Standard (§9) | Mechanical truncation produces `quit, from anyw…` at exactly the width where help matters most |
| **D-10-5** | No per-row cursor; the viewport is the focused thing (§16.1) | A cursor would imply `⏎` does something. The keys are pressed elsewhere — that is the premise of the screen |
| **D-10-6** | The reserved block's key column is `--z-text-tertiary`, not amber (§15) | Amber is the machine speaking; a key that does nothing is the machine listing what it will say later. Equal amber would deny the one distinction the block exists to make |
| **D-10-7** | `?` is described as `this key map — you are looking at it` (§5) | It is a no-op here, and the one screen whose job is not lying cannot list itself without a qualifier |
| **D-10-8** | A binding whose **description differs by state** is **two registry entries with mutually exclusive predicates**, never one entry with a conditional string (§3, consequence 2b) | `v` is the first — `see the covers` / `back to the list`, matching the two footer hints exactly. One entry choosing a string at render time puts the footer's wording and the help's wording back in two places, which is the drift D-10-1 abolishes. Exactly one is ever live, so block 1 gains exactly one row |
| **D-10-9** | **`x dismiss` is excluded from the count by its availability predicate, exactly as `m` is** — default 34 lines, 35 while the capability note is showing (§5, state H13) | Listing it always puts a dead key in the map on nine renders out of ten; never listing it omits a bound key. The predicate is the only answer that is right in both states, and the registry already carries the field |
| **D-10-10** | **The viewport never re-flows its cut point** to avoid an orphaned heading at the fold (§4) | The moment paging is anything other than one line per `↓`, the position readout stops describing the document and starts describing a pagination the player cannot predict. The heading is orphaned in a *window*, and the window is labelled |
| **D-10-11** | **A mode earns block 2 only when it is *additive*; a substitutive mode rewrites block 1** (§8) | `Z-07` adds an editor and chips on top of the list, so its keys are genuinely extra. `Z-15` replaces the list, so `↑ ↓` **means something else** — listing both descriptions would hand the player two contradictory accounts of one key on one screen |

---

## 24 · Screen-specific acceptance criteria

Beyond `00-design-brief.md` §10 and `02-colour-budget.md` §10, both of which apply in full.

1. **Every key listed in block 1 actually fires on the origin, in the state the origin was in.**
   Verified by pressing each one after `Esc`.
2. **Every key that fires on the origin is listed.** The inverse test, and the one that catches
   drift. **Automatable against the registry** — a test that asserts *the set of live bindings
   equals the set of rendered rows* is the acceptance criterion this screen is designed to make
   possible. **It subsumes every line count in §4.0**, which is why none of them may be
   hard-coded: assert `len(rendered) == len(live)`, never `== 34`. *Rev B exists because two
   drawn goldens had drifted from the registry — `v` missing from all of them and `g` / `G`
   missing from §9 — and no test could have failed.*
3. **Opened from an empty library, block 1 lists two keys, not eleven**, and carries the prose
   line.
4. **Opened from filter mode, block 2 is present**; opened from anywhere else, it is absent.
5. **`? help` is not in this screen's footer**, at any tier.
6. **`?` on `Z-10` does not push a second help screen**, and one `Esc` returns to the origin.
7. **`Esc` restores the origin's focus exactly** — the same ledger row at row 380 of 412, the
   same filter query, the same detail view.
8. **No description is ever truncated** at any tier; short forms are used below Standard and are
   the declared strings, not clipped long ones.
9. **The position readout is present when the content overflows and absent when it does not.**
10. **The `m` row and the footer's `m mute` are present if and only if audio has been enabled.**
10b. **The `v` row is present if and only if `v` is bound on the origin** — absent on an empty
    library, absent below 40 columns, absent once the capability note has been dismissed on a
    terminal that draws no images (`Z-15` §10.1's five rows, asserted one at a time).
10c. **The `x` row is present if and only if the capability note is showing** — the same
    assertion as 10, against the same registry field. Press `v` on a terminal with no image
    support, then `?`: block 1 ends `r` · `v` · `x` and the document is **35** lines. Press
    `Esc`, `x`, `?`: `v` and `x` are both gone and it is **27**.
10d. **From cover mode, block 1 carries the deck's descriptions and there is no block 2**
    (D-10-11). Assert that no render ever contains both `move the cursor` and a covers-movement
    description for the same key.
11. **There is zero cyan on this screen**, of any class, in any state — by the
    `02-colour-budget.md` §3.1 machine method.
12. **Amber cells are ≤ 10 % of the viewport**, measured — this screen has the largest amber
    spend in the bundle and is the one to check.
13. **`NO_COLOR=1` loses no information**, including which block is which and that the reserved
    keys are not live.
14. **No banner appears on this screen in any state**, including with the network off.
15. **At ExtraWide the two columns read correctly as a character stream** — headings span, cells
    are self-contained.
16. **Founder-validated screenshot before merge**, at the six viewports of
    `03-responsive.md` §7 plus `NO_COLOR=1` and forced-16-colour at 80 × 24 — **and from at
    least three different origins**, including an empty library. No screenshot → not GOLDEN →
    no merge.
