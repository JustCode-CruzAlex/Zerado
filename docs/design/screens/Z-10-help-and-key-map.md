---
title: Zerado — Z-10 Help and key map
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-10
rev: A
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
§1.5 (readout role) · `02-colour-budget.md` §10 · `03-designer-manual.md` §3 ·
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
> 2. A binding may declare a **precondition** — `s` on `Z-04` describes *set this game's status*
>    when a row exists and *connect Steam* when the library is empty. `Z-10` asks the origin, in
>    its current state, and renders the answer.
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

## 4 · Mockup — 80 × 24, pushed from a populated Library

`tier = Wide` · `leftInset = 3` · **body = 74 × 16** · content begins at **column 4**.
15 content rows + 1 position row. **33 lines in total, so it scrolls.**

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
│                                                                                │
│   IN FILTER MODE                                                               │
│                                                                                │
│   esc           leave the editor, keep the filter                              │
│   ROWS  1–15 of 33                                                             │
│   ↑↓ scroll   esc back   q quit                                                │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Character counts:** body **74 × 16** · key field **12** + gutter **2**, so every description
begins at body column **15** · the widest description row is 50 cells
(`?  this key map — you are looking at it`) · position row **18** · footer **29**.

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
ROWS  1–15 of 33
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

Lines 19–33 of 33. The global and reserved blocks.

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
│   ROWS  19–33 of 33                                                            │
│   ↑↓ scroll   esc back   q quit                                                │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

> **`?  this key map — you are looking at it`.** Pressing `?` here unwinds to a screen already
> on the stack and therefore does nothing visible. Listing it as *"this key map"* with no
> qualifier would be the screen lying about itself on the one screen whose entire job is not
> lying. **The clause is generated from the binding's precondition** (D-10-1 consequence 2), the
> same mechanism that makes `s` describe two different things on `Z-04`.

> **`m` is present only when audio has been enabled.** When it has never been enabled the
> binding is not live, so `Z-10` does not list it and the footer does not carry it — because
> there is nothing to mute. This falls out of D-10-1; it is not a special case.

---

## 6 · Mockup — pushed from an EMPTY Library: a different, honest list

**This is the requirement, demonstrated.** The player came from the same screen ID, `Z-04`, in a
different state — and eight of the eleven bindings are not live, so eight of them are not here.

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
│   s             connect Steam                                                  │
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

**Note what is gone:** `↑ ↓`, `g`/`G`, `^d`/`^u`, `⏎`, `/`, `f`, `r` — every one of them
needs a row, and there are none. **Note what changed:** `s` reads *connect Steam*, not *set this
game's status*, because that is what it does here (`01-design-system.md` §10.1). **`m` is gone**
because audio has never been enabled in this render.

**The prose line is not decoration.** A player who opens help on an empty library and sees two
keys should be told *why* it is two, or the screen reads as broken. One sentence, dry, concrete:

```
Nothing is in the library yet, so most keys do nothing here.
```

> **FLAG** · **`s` describing two different actions on two states of one screen is an upstream
> collision**, not a design of this screen — see `Z-04-library.md` §18 item 1. `Z-10` renders it
> correctly either way, because it asks the origin rather than assuming.

---

## 7 · Mockup — ExtraWide, 120 × 40: two key columns

`leftInset = 4` · **body = 112 × 32** · two columns of **52** with an **8**-column gutter.
The whole map fits without scrolling.

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                                      │
│    Zerado ✦ Help                                                                                                     │
│                                                                                                                      │
│    HELP                                                                                                              │
│                                                                                                                      │
│                                                                                                                      │
│    ON THIS SCREEN — LIBRARY                                    IN FILTER MODE                                        │
│                                                                                                                      │
│    ↑ ↓  k j      move the cursor                               esc           leave the editor, keep the filter       │
│    g   G         first row · last row                          esc esc       clear the filter                        │
│    ^d  ^u        half a page down · up                         tab           move to the state chips                 │
│    ⏎             open this game                                f             jump to the state chips                 │
│    s             set this game's status                        space         toggle the chip under the cursor        │
│    /             filter and search                                                                                   │
│    a             add a game by hand                                                                                  │
│    r             sync with Steam                                                                                     │
│                                                                                                                      │
│    EVERYWHERE                                                  RESERVED — NOT BOUND YET                              │
│                                                                                                                      │
│    ?             this key map                                  :   ^k        the command palette — Phase 2           │
│    ,             settings                                      1 – 9         quick filters — Phase 2                 │
│    m             mute or unmute the audio                      n   p         next · previous game — Phase 2          │
│    esc           back one level                                                                                      │
│    q             quit                                                                                                │
│    ^c            quit, from anywhere, always                                                                         │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│                                                                                                                      │
│    ↑↓ scroll   esc back   m mute   q quit                                                                            │
│                                                                                                                      │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**The position row is absent** because nothing is scrolled — `ROWS 1–24 of 24` would be
furniture. It appears the moment the content exceeds the viewport.

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
| **`Z-07` Filter mode** | `ON THIS SCREEN — LIBRARY` **plus block 2** `IN FILTER MODE` | `Z-07-filter-and-search.md` §13 |
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
│ n  p      paging, P2           │
│ ROWS  1–20 of 22               │
│ ↑↓ scroll  esc back            │
└────────────────────────────────┘
```

**What shortens at Tiny, and nothing else:** the block-1 heading wraps to two lines rather than
truncating; the key field goes 12 → **8**; descriptions shorten to their **declared short form**
(`quit, always`, `palette, P2`); the heading `RESERVED — NOT BOUND YET` shortens to `RESERVED`.

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
| ExtraWide | `↑↓ scroll   esc back   m mute   q quit` | 38 ≤ 112 |
| Tiny | `↑↓ scroll  esc back` | **19** |
| No overflow (ExtraWide) | `esc back   q quit` | 17 |

**`? help` is absent from every one of them.** It is the smallest possible demonstration of
`04-navigation-and-focus.md` §6 — *"keys that do nothing here are not listed"* — applied to the
screen that would be most embarrassing to get wrong. Separator 3 spaces, tightening to 2 before
any hint drops.

---

## 12 · The full state table

| # | State | Trigger | Composition | Copy |
|---|---|---|---|---|
| **H1** | **First run** — pushed from an empty `Z-04`, or from `Z-01` | Library empty | §6. Block 1 lists only the bindings that are live, with the prose line explaining why there are two | §13.2 |
| **H2** | **From a populated Library** | The default | §4 | §13.1 |
| **H3** | **From filter mode** | `?` pressed with the filter applied and the **editor blurred** | Block 1 = `LIBRARY`, block 2 = `IN FILTER MODE` | §13.1 |
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

**33 lines.** At Wide the viewport shows 15 of them.

### 13.2 · The empty-library variant

```
ON THIS SCREEN — LIBRARY

Nothing is in the library yet, so most keys do nothing here.

s             connect Steam
a             add a game by hand
```

### 13.3 · The position readout

```
ROWS  1–15 of 33
ROWS  19–33 of 33
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
| heading | `RESERVED — NOT BOUND YET` | `RESERVED` |

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
- **Type-neutral where equally natural** — `move the cursor`, `first row · last row`,
  `filter and search`, `ROWS`. `open this game` and `add a game by hand` say *game* because a
  game is what the player is looking at.

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
   would be. See `Z-06-set-status.md` §18 item 3 — it is open for the founder.

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
│                                                                                │
│   IN FILTER MODE                                                               │
│                                                                                │
│   esc           leave the editor, keep the filter                              │
│   ROWS  1–15 of 33                                                             │
│   ↑↓ scroll   esc back   q quit                                                │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

| Information | Carried without colour by |
|---|---|
| Which block is which | The **UPPERCASE headings**, and the `InterElementGap` above each |
| Key versus description | **Case** — `^d  ^u` against `half a page down · up` — plus the fixed column edge at 15 |
| That the reserved keys are different | The heading `RESERVED — NOT BOUND YET`, the `— Phase 2` stamp on every row, and **position** — they are last |
| Where you are in the document | `ROWS  1–15 of 33` |
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
+ 24 + 14 (two headings) + ~40 (the key column) ≈ **82**. The key column is the largest single
amber spend on any screen in this bundle and it is still well under half the ceiling — but it is
the screen to re-measure if a block is ever added.

**Red: none.** No scanner, no annunciator, no error text. **`--z-border` is not used at all** —
there are no boxes, no rules and no separators on this screen. **No region is separated by
fill**, and at ExtraWide the two columns are separated by an 8-column gutter, which survives
every colour depth because it is made of nothing.

---

## 20 · Reuse verdict, per element

| Element | Verdict | Why |
|---|---|---|
| Scroll body | **`bubbles/viewport` — direct fit** | 33 lines into 15. This is exactly what the primitive is for |
| Key rows | Build fresh · `lipgloss` + the width-aware pad | Two padded columns. **`bubbles/table` does not fit** — it is the primitive that independently dropped the title, the scroll and the pinned footer on two FlowForge screens, and a key map needs none of its selection or column machinery |
| **The key registry** | **Build fresh, once, and make it the single source for dispatch, the footer and this screen** | D-10-1. This is the reuse that matters on this screen: not a component, a **source of truth** |
| Two-column layout (ExtraWide) | `lipgloss` join | |
| Position readout | **The same component `Z-04` uses** | Same wording, same place, same behaviour (3.2.4 Consistent Identification) |
| Markdown rendering | **Not `glamour`** | The key map is a generated table, not a document. `glamour` carries its own palette and would have to be restyled to render something that was never markdown |
| Header band | The shared `Frame` | Enforced at the router, so this screen cannot render frameless |
| Error rendering | **Not `charmbracelet/log`** | Developer logging role |

---

## 21 · Upstream findings

| # | Finding | Where | Owner |
|---|---|---|---|
| 1 | **No document specifies where a key's description lives.** `04-navigation-and-focus.md` §3 has a global key table and §6 has a footer rule, but nothing binds them together — so today a key could be dispatched, listed in a footer and described in help by three separate strings. **D-10-1 proposes the registry that makes them one** | `04-navigation-and-focus.md` §3, §6 | `fft-tui-architect` |
| 2 | **`04-navigation-and-focus.md` §3 lists `f` under no screen**, but `01-design-system.md` §7.3 binds it to the state filter chips *"(ExtraWide, or via `f`)"*. This spec lists it under `Z-04` as `filter by state`; it must exist in the registry or `Z-10` cannot show it | `04-navigation-and-focus.md` §3 | `fft-tui-architect` |
| 3 | **`04-navigation-and-focus.md` §3 says `?` "unwinds if already on the stack"**, which means `?` on `Z-10` is a no-op. Nothing says how help should describe itself in that case. **§5 resolves it in copy** | `04-navigation-and-focus.md` §1 rule 3 | — |
| 4 | **`?` is unavailable in two documented cases** — inside a text input (2.1.4) and while an overlay is open (D-06-4) — which slightly qualifies WCAG 3.2.6's *"from every screen"*. Both are recorded in §16.3 rather than glossed | `00-design-brief.md` §3.1 | `fft-design-architect` |
| 5 | **`s` describes two different actions on two states of `Z-04`**, which this screen renders correctly by asking the origin. It remains an upstream collision | `01-design-system.md` §10.1 vs `04-navigation-and-focus.md` §3 | founder — `Z-04-library.md` §18 item 1 |
| 6 | `03-designer-manual.md` §5.11 verdict 3 still reads as a permanent rejection of the audio subsystem, superseded by founder direction relayed 2026-08-25. **`m` must be in the global key table before `Z-10` can list it** | `03-designer-manual.md` §5.11 · `04-navigation-and-focus.md` §3 | `fft-design-architect` / `fft-tui-architect` |

---

## 22 · Open for the founder

1. **The key registry (D-10-1) is a spine-level ask, not a screen decision.** It is the only way
   `01-screen-inventory.md` §5 stays true after the tenth key is bound, and it makes the footer
   and the help provably consistent for free. It is small — a binding struct with a key, a long
   description, a short description, a scope and an optional precondition — but it has to be
   decided by `fft-tui-architect` before `Z-10` is built, because otherwise `Z-10` gets a
   hand-written table and this screen's promise is already broken on day one.
2. **`?` is unavailable in two places** (§16.3): inside a text input, and while an overlay is
   open. The first is forced by WCAG 2.1.4 and is not negotiable. The second is a design decision
   (`Z-06-set-status.md` D-06-4) and is open there. **Confirm both, so the conformance statement
   for 3.2.6 is written once and honestly.**
3. **The reserved block.** Listing keys that deliberately do nothing is unusual, and it is here
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

---

## 24 · Screen-specific acceptance criteria

Beyond `00-design-brief.md` §10 and `02-colour-budget.md` §10, both of which apply in full.

1. **Every key listed in block 1 actually fires on the origin, in the state the origin was in.**
   Verified by pressing each one after `Esc`.
2. **Every key that fires on the origin is listed.** The inverse test, and the one that catches
   drift. **Automatable against the registry** — a test that asserts *the set of live bindings
   equals the set of rendered rows* is the acceptance criterion this screen is designed to make
   possible.
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
