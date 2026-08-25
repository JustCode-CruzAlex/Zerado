---
title: Zerado — Z-08 Add a game by hand
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-08
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-08 · Add a game by hand

> Fills [`../03-designer-manual.md`](../03-designer-manual.md) §3's 16-section contract.
> Composition binding from [`../../blueprint/02-composition.md`](../../blueprint/02-composition.md) §2 —
> single-pane `huh` form, `R = 1`.
> **This is the screen that makes the published promise true:**
> *"A physical copy isn't a second-class row in the list."*
> (`content/landing-copy.md` §06, in the ratified Zerado canon.)

---

## 1 · Identity

| | |
|---|---|
| **Screen** | `Z-08` · Add a game by hand |
| **Phase** | 1 |
| **Kind** | Route |
| **Routes in** | `Z-01` door 2 · `a` from `Z-04` |
| **Routes out** | pop on `Esc` — to `Z-04`, or to `Z-01` when that is what is beneath |
| **Offline class** | **WORKS** ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2) — *"The whole point of a physical shelf."* No banner, no notice, no difference |
| **Writes** | one `game` row with `provider_id = 'physical'` and a fresh UUID `provider_ref` ([`06-data-seams.md`](../../blueprint/06-data-seams.md) §2.2) |

**`physical` is a `Provider` and deliberately not a `Syncer`.** Its "sync" is this screen. It is
not a Steam-shaped provider with a hole where `Sync` should be — it is a different shape of thing,
and everything downstream reads `Capabilities`, never `ProviderID`.

```go
Capabilities{Sync: false, Playtime: false, LastPlayed: false,
             OwnedSince: true, Credentials: nil}
```

---

## 2 · Purpose

**Enter a disc, a cartridge or anything else the stores do not know about, as a first-class row.**

Screen inventory §5: *"`Z-08` must not require a field Steam happens to have. A cartridge has no
app ID."*

---

## 3 · Mockups at 80 columns

Frame row map as [`Z-01-first-run.md`](./Z-01-first-run.md) §3.1 — content at **column 4**,
body `74 × 16`. The field row budget is **byte-identical to
[`Z-02-connect-a-store.md`](./Z-02-connect-a-store.md) §3.1** — one primitive, two screens.

### 3.1 · The field set — four fields, and only two are required

| # | Field | Column | Required | Why it is what it is |
|---|---|---|---|---|
| 1 | **Title** | `title` (+ derived `sort_title`) | **yes** | R-10(a)'s identity column. A row without one cannot exist |
| 2 | **Platform** | `platform` | **yes** | `NOT NULL` in the schema, and it is half of `game_uid` |
| 3 | **State** | `status_manual` | no — defaults to `NOT STARTED` | see §8.2 |
| 4 | **Owned since** | `owned_since` | no | `Capabilities.OwnedSince = true` is the only optional capability `physical` claims |

**What is deliberately absent, and why each absence is a decision:**

| Not asked for | Why |
|---|---|
| **an app ID / store ID** | *"A cartridge has no app ID."* The `provider_ref` is a **UUID Zerado mints**, never something the player has to find |
| **playtime** | `Capabilities.Progress = false`. A cartridge has no telemetry. Asking for a number nobody can know invites a made-up one, and a made-up number would then drive the derivation |
| **last played** | `Capabilities.LastPlayed = false`, same argument |
| **notes** | the column exists, but **no Phase 1 screen renders it back.** A field the player fills and never sees again is a field that teaches them the product forgets things. It earns its place when `Z-05` has somewhere to show it — §17 |
| **cover art / *sinopse*** | Phase 2. Never claim what is not built (anti-pattern 14) |

### 3.2 · RENDER 80×24 — the empty form

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Add a game by hand

   ADD A GAME BY HAND

   Nothing here syncs. A disc or a cartridge is a row like any other.

   TITLE
   ▌  ┃ ▎                                                                   ┃

   PLATFORM
      │                                                                     │
     Whatever you call it. PS2, Switch, Mega Drive.

   STATE
     ▌ ○  NOT STARTED    ◐  IN PROGRESS    ◉  ZERADO    ⊘  ABANDONED

   OWNED SINCE
      │                                                                     │
     Optional. A year is enough.

   tab next   ⏎ add   esc done   ? help   q quit

```

### 3.3 · RENDER 80×24 — filled, `STATE` focused, and the duplicate note

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Add a game by hand

   ADD A GAME BY HAND

   You already have Shadow of the Colossus on PS2. Adding another is fine.

   TITLE
      │ Shadow of the Colossus                                              │

   PLATFORM
      │ PS2                                                                 │
     Already here: PS2, Switch.

   STATE
   ▌   ○  NOT STARTED    ◐  IN PROGRESS  ▌ ◉  ZERADO    ⊘  ABANDONED

   OWNED SINCE
      │ 2004                                                                │
     Optional. A year is enough.

   ←→ state   tab next   ⏎ add   esc done   ? help   q quit

```

### 3.4 · RENDER 80×24 — after `⏎ add` (the receipt, and the form ready for the next disc)

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Add a game by hand

   ADD A GAME BY HAND

   Shadow of the Colossus is in the library. 248 games.

   TITLE
   ▌  ┃ ▎                                                                   ┃

   PLATFORM
      │ PS2                                                                 │
     Already here: PS2, Switch.

   STATE
     ▌ ○  NOT STARTED    ◐  IN PROGRESS    ◉  ZERADO    ⊘  ABANDONED

   OWNED SINCE
      │                                                                     │
     Optional. A year is enough.

   tab next   ⏎ add   esc done   ? help   q quit

```

**The form stays open, and `Platform` keeps its value** — *design decision, and it is the one that
decides whether this screen feels finished.* Somebody entering a physical shelf enters **twenty
discs in a row**, and nineteen of them are for the same console. Clearing the platform between
each would be the product forgetting a fact it was just told. `Title`, `State` and `Owned since`
reset; focus returns to `Title`; the intro line is replaced by the receipt and stays replaced.

### 3.5 · The `STATE` field — column budget at 74

The chips reflow **never**: each option reserves its own 2-column marker field, so moving the
cursor moves one `▌` and shifts nothing.

| Field | Cols | Range | Note |
|---|---|---|---|
| field cursor | **2** | 4–5 | `▌` when the `STATE` field itself holds focus |
| option 1 marker | **2** | 6–7 | `▌` on the selected option |
| option 1 glyph field | **2** | 8–9 | fixed, width-aware padded — `○` and `◐` are **Ambiguous** |
| gap | 1 | 10 | |
| option 1 label | **11** | 11–21 | `NOT STARTED` — the longest label sets the field |
| gap | 2 | 22–23 | |
| option 2 block | 16 | 24–39 | `◐  IN PROGRESS` |
| gap | 2 | 40–41 | |
| option 3 block | 11 | 42–52 | `◉  ZERADO` |
| gap | 2 | 53–54 | |
| option 4 block | 14 | 55–68 | `⊘  ABANDONED` |
| | **65** of 74 | | 9 columns spare — the row does not stretch |

The chip's own anatomy is unchanged from
[`../01-design-system.md`](../01-design-system.md) §3.1: **glyph field 2 + gap 1 + label**, and
**the label is never dropped.**

---

## 4 · Mockup at `120 × 40`

`leftInset` **4** · body `112 × 32` · content at **column 5** · `1 + 3 + 1 + 32 + 1 + 1 + 1 = 40` ✓
Responsive table: *"Field + help beside."* Field keeps its 74-column budget; help takes a
34-column column across a 4-column gutter — `74 + 4 + 34 = 112` ✓, identical to `Z-02`.

```text
0.........1.........2.........3.........4.........5.........6.........7.........8.........9.........A.........B.........

    Zerado ✦ Add a game by hand

    ADD A GAME BY HAND

    Nothing here syncs. A disc or a cartridge is a row like any other.

    TITLE
    ▌  ┃ ▎                                                                   ┃

    PLATFORM
       │                                                                     │    Whatever you call it.
                                                                                  PS2, Switch, Mega Drive.

    STATE
      ▌ ○  NOT STARTED    ◐  IN PROGRESS    ◉  ZERADO    ⊘  ABANDONED

    OWNED SINCE
       │                                                                     │    Optional. A year is
                                                                                  enough.
```

---

## 5 · Visual hierarchy

**The one thing the player must see first: the `TITLE` field, focused and empty.**

| Rank | Element | Channel | Note |
|---|---|---|---|
| 1 | `ADD A GAME BY HAND` | case + weight + `--z-primary` | the H1 |
| 2 | **the focused field** | `▌` gutter (position) + heavy `┃` (weight) + `--z-focus-ring` cyan (colour) | on first paint this is `TITLE`, which is the only thing the screen actually needs |
| 3 | the intro line — or, after the first add, **the receipt** | `--z-text` at the top of the body | the same row, doing two jobs at two moments |
| 4 | the `STATE` chips | **glyph shape first, colour second** — `○ ◐ ◉ ⊘` is a single visual progression | the only place in this cluster where the ratified state system renders in full |
| 5 | field labels | UPPERCASE `--z-text-secondary` | readout role |
| 6 | field help | sentence case `--z-text-secondary` | |
| 7 | breadcrumb, footer | chrome | |

### 5.1 · Two `▌` on one screen — and here there are briefly **three**

`Z-08` is the densest case of the collision [`Z-02-connect-a-store.md`](./Z-02-connect-a-store.md)
§5.1 records: the **field cursor**, the **selected-option marker**, and — if a write fails — the
**error annunciator**. They stay apart by role and by column, and the two channels that matter
survive `NO_COLOR`:

| | field cursor | option marker | error annunciator |
|---|---|---|---|
| column | **4–5**, the body's left edge | **inside the `STATE` row**, at a fixed option offset | 4–5, below a gap |
| colour | `--z-primary` **214** | the option's own **state colour** when selected | `--z-scanner` **9** |
| what follows | a field label's field | a **state glyph** | an **UPPERCASE heading** |

The field cursor and the option marker are both present on the `STATE` row when it has focus
(§3.3) and that is correct: one says *this field*, the other says *this value*. They are two
columns apart and only one of them is ever on a chip.

---

## 6 · Every applied spacing token, by name

| Token | Wide value | Where Z-08 spends it |
|---|---|---|
| `OuterMarginX` | **2** | frame inset |
| `OuterMarginY` | **1** | rows 1 and 24 |
| `InnerPaddingX` | **1** | inside the frame, **and inside each field boundary** |
| `InnerPaddingY` | **1** | row 22 |
| `InterElementGap` | **1** | breadcrumb→title · band→body · intro→field 1 · between each field block |
| `HeaderBandHeight` | **3** | `hasSubtitle = false` |
| `leftInset` | **3** | header-left **==** content-left at column 4 |
| chip glyph field | **2** cols | fixed, width-aware padded ([`02-composition.md`](../../blueprint/02-composition.md) §2.2.1) |
| chip label field | **11** cols | `NOT STARTED` / `IN PROGRESS` are the longest (§3.1) |
| inter-chip gap | **2** cols | the ledger row's own gutter value (§4.3) |
| ExtraWide help gutter | **4** cols | §4 |

**No in-body key-hint line** — the reserved footer row is the hint block, as
[`Z-03-sync.md`](./Z-03-sync.md) §6 records for the whole cluster.

---

## 7 · Colour, glyph and label for every state shown

### 7.1 · The four game states — ratified, CVD-verified, rendered here in full

Read from [`../01-design-system.md`](../01-design-system.md) §3.2, which reads the brand's measured
table. **Colour AND glyph AND label — all three, every chip.**

| State | Token | Hex | ANSI-256 | 16-colour | Glyph | ASCII | Label | Ratio |
|---|---|---|---|---|---|---|---|---|
| Not started | `--z-state-not-started` | `#A5A29B` | **247** | `white` | `○` U+25CB — **Ambiguous** | `[ ]` | `NOT STARTED` | **7.62** AA |
| In progress | `--z-state-in-progress` | `#FFB000` | **214** | `bright yellow` | `◐` U+25D0 — **Ambiguous** | `[~]` | `IN PROGRESS` | **10.59** AAA |
| Zerado | `--z-state-zerado` | `#19E0FF` | **45** | `bright cyan` | `◉` U+25C9 — Neutral | `[*]` | `ZERADO` | **12.15** AAA |
| Abandoned | `--z-state-abandoned` | `#C77DFF` | **177** | `bright magenta` | `⊘` U+2298 — Neutral | `[x]` | `ABANDONED` | **7.21** AA |

- **The 16-colour floor resolves all four to four distinct slots** — no collisions (brand §5.2).
- **The warm grey `#A5A29B` is load-bearing engineering, not taste.** The blue-cast `#9FB0C6`
  collapsed against the cyan at **ΔE 8.8 under deuteranopia**; the warm grey measures **25.8**.
  Never "correct" it back toward blue.
- **The floor to protect is ΔE 11.9** — zerado × abandoned under deuteranopia. On this row those
  two chips sit **two options apart**, which is exactly where glyph and label carry load rather
  than merely reinforce.
- `ZERADO_ASCII=1` swaps the whole column to `[ ] [~] [*] [x]`, which is entirely **Narrow** and
  immune (§1.2 rule 3).

### 7.2 · Selected versus unselected — and why the unselected chip loses its colour

| Chip state | Rendering | Co-render still holds? |
|---|---|---|
| **selected** | state colour, **bold**, `▌` marker in the state colour | **yes** — colour + glyph + label |
| **unselected** | `--z-text-secondary` `#A9B5C7` **249** / `white`, normal weight, blank marker | **yes** — **glyph + label still carry the state in full** — the colour channel is being spent on *selection*, not on state, and the design system sanctions exactly this: §3.4, *"In a filter, unselected — glyph + label in `--z-text-tertiary`"* |

**`--z-text-secondary`, not `--z-text-tertiary`** — the cluster-wide substitution recorded in
[`Z-01-first-run.md`](./Z-01-first-run.md) §7.4. `--z-text-tertiary` has **no derived ANSI-256
index**, and its documented interim is *uncoloured*, which on this row would make three unselected
chips render in the terminal's default foreground — brighter, on many themes, than the selected
one. `--z-text-secondary` is derived (**249**), measured (**9.36** AAA) and unambiguously dimmer
than `--z-text`.

### 7.3 · The rest of the screen

| Element | Token | Hex | ANSI-256 | 16-colour | Ratio | `NO_COLOR` |
|---|---|---|---|---|---|---|
| screen title | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | **yes** — case + own row |
| intro line / receipt | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** | **yes** |
| field label | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** | **yes** — case |
| field value | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** | **yes** |
| field boundary, unfocused | `--z-border-strong` | `#64748B` | **67** | `bright black` | **4.08** — meets 1.4.11 | **yes** — `│` drawn |
| field boundary, focused | `--z-focus-ring` | `#19E0FF` | **45** | `bright cyan` | **12.15** | **yes** — **heavy `┃`** + `▌` |
| field cursor `▌` | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | **yes** — position |
| write-error annunciator | `--z-scanner` | `#FF2E2E` | **9** | `bright red` | 5.25 — structure, not text | **yes** — `▌` + case |
| write-error text | **uncoloured + bold** — documented interim (§11.2); `--z-scanner-300` `#FF6B6B` (**6.99** AA) has no derived index | | *underived* | `bright red` | **6.99** | **yes** |
| audio indicator | from [`../01-design-system.md`](../01-design-system.md) §5.2 | | **214** / **249** | | | **yes** — `▮`/`▯` |

---

## 8 · The full state table

| # | State | Trigger | Renders | Copy |
|---|---|---|---|---|
| 1 | **First run** — library empty, reached from `Z-01` door 2 | | §3.2, `PLATFORM` help shows the **generic** examples because there are no existing platforms to list | help: `Whatever you call it. PS2, Switch, Mega Drive.` · receipt after the first add: `Shadow of the Colossus is in the library. 1 game.` — **singular** |
| 2 | **Empty form** — library not empty | `a` from `Z-04` | §3.2 with `PLATFORM` help listing what is already there | `Already here: PS2, Switch.` |
| 3 | **Typing** | | focus ring on the active field; **no single-key shortcut fires** (2.1.4) | — |
| 4 | **`TITLE` empty on submit** | `⏎` | focus jumps to `TITLE`; its help row becomes the inline error | `A title is the one thing Zerado can't invent.` |
| 5 | **`PLATFORM` empty on submit** | `⏎` | focus jumps to `PLATFORM`; help row becomes the error | `Which machine? PS2, Switch, whatever you call it.` |
| 6 | **Duplicate title+platform** | live, as `PLATFORM` is completed | a **note**, not an error — the intro row | §10 · `You already have Shadow of the Colossus on PS2. Adding another is fine.` |
| 7 | **Added** | `⏎` with both required fields | §3.4 — receipt, form cleared except `PLATFORM`, focus back on `TITLE` | §10 |
| 8 | **Write failed** | `Store.AddManual` returns an error | message block replaces the help rows; **the form keeps every value** | §10 · **W** |
| 9 | **Offline** | network off | **nothing changes.** `Z-08` is `WORKS` — no banner, no notice, no difference. Rendering an offline banner on the one screen that never needed a network would be a banner that is furniture ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §3) | — |
| 10 | **Loading** | **N/A** — the only read is the distinct-platform list for the help line, which is local and already resolved before first paint. If it were ever slow, the help line renders **without** the list rather than showing a spinner | — |
| 11 | **Partial** | **N/A** — one row, one transaction | — | — |
| 12 | **Audio** | | rows as [`Z-01-first-run.md`](./Z-01-first-run.md) §8 rows 10–12. **No cue fires on an add** — the closed cue list (§15.3) is sync complete, error, and a game becoming `ZERADO`; adding is none of them | — |
| 13 | **`NO_COLOR`** | env | §12 | identical text |
| 14 | **Below `24 × 8`** | | never renders — [`Z-11-fatal-error.md`](./Z-11-fatal-error.md) | the refusal sentence |

### 8.1 · The result must be indistinguishable in rank — the promise, made checkable

The row `Z-08` writes lands in `Z-04` through **the same renderer, in the same sort, with the same
state chip and the same identity column** as anything Steam reported. Drawn against a synced
neighbour, at Wide:

```text
0.........1.........2.........3.........4.........5.........6.........7.........
   ▌ ◉  ZERADO       Return of the Obra Dinn                       9h    STM
     ○  NOT STARTED  Shadow of the Colossus                              PHY
     ◐  IN PROGRESS  Outer Wilds                                  12h    STM
```

*(Rendered by [`../01-design-system.md`](../01-design-system.md) §4, whose owner is `Z-04`'s spec —
reproduced here only as `Z-08`'s acceptance target.)*

**The playtime cell is empty, not `0h`.** `Capabilities.Progress = false` means *not reported*, and
`Item.Playtime` is a `*int` precisely so that **"not reported" and "reported as zero" stay different
facts**. Printing `0h` would assert the player has never played a disc they may have finished
twice.

What is **not** different: the row height, the chip, the title column, the sort position, the
focusability, the ability to be filtered, and the counts it contributes to on the pinned summary.
What **is** different: three characters in the source column, `PHY` against `STM` — a statement of
origin, not of rank.

### 8.2 · `Z-08` always writes `status_manual` — *design decision*

The state machine's model is `effective_status = status_manual ?? derive(playtime, capabilities)`,
and `status_manual = NULL` is a real, representable *"no opinion"*.

`Z-08` nevertheless writes the field's value **always**, never `NULL`. The reason is
[`05-state-machine.md`](../../blueprint/05-state-machine.md) §2.1 in its own words: for a provider
that does not report playtime, *"all four states are manual, always — the derivation has no input
and returns `NOT STARTED` as the honest default."* For a `physical` row, `NULL` and
`not_started` therefore render identically, forever, and there is no future sync that could ever
make them diverge. Given that, storing what the screen visibly showed is the honest write, and it
is also the one that survives the Phase 4 boundary, where **only what the player typed crosses**
([`06-data-seams.md`](../../blueprint/06-data-seams.md) §6.1).

`Z-06 Set status`'s fifth item — *clear override* — is unaffected: it still exists, and on a
physical row it resolves to `NOT STARTED`, which is what the copy under it will say.

### 8.3 · A duplicate is a note, never a block — *design decision*

`(provider_id, provider_ref)` is the uniqueness constraint and `provider_ref` is a **freshly minted
UUID**, so a duplicate cannot collide. `game_uid` is *indexed, not unique*, exactly so two rows may
legitimately be the same title ([`09-erd.md`](../../blueprint/09-erd.md) §4).

Two copies of the same disc is a real thing a shelf contains — a PAL and an NTSC copy, a boxed and
a loose one, one lent out. Blocking it would be the product telling a player their own shelf is
wrong. So `Z-08` **says what it sees and gets out of the way.**

---

## 9 · The key map

A text input holds focus for three of the four fields, so **every printable key is literal text**
(2.1.4). The `STATE` field is the exception and it is why `←→` appears.

| Key | Does | In the footer? |
|---|---|---|
| any printable key | types into the focused field | no |
| `tab` / `shift-tab` · `↓` `↑` | next / previous **field** | **yes** — `tab next` |
| `←` `→` · `h` `l` | move between `STATE` options — **only while `STATE` has focus** | **yes, only then** — `←→ state` |
| `⏎` | **add** — from any field | **yes** — `⏎ add` |
| `Ctrl-U` / `Ctrl-W` | clear field / delete word | no |
| `esc` | **pop the route.** `Z-08` saves on `⏎`, so there is never uncommitted work older than the current row | **yes** — `esc done` |
| `Ctrl-C` | quit | no |
| `q` `?` `,` `m` `s` `a` `/` `r` | **type their character** while a text field has focus. On the `STATE` field they are **unbound and silent** — never a shortcut, so that the same key never does two things on one screen | `? help`, `q quit` and `m` are listed at screen granularity — see [`Z-02-connect-a-store.md`](./Z-02-connect-a-store.md) §9 |
| `Tab` as a region key | **unbound** — `R = 1` | no |

**Footer per context:**

| Context | Footer |
|---|---|
| a text field has focus | `tab next   ⏎ add   esc done   ? help   q quit` |
| `STATE` has focus | `←→ state   tab next   ⏎ add   esc done   ? help   q quit` |
| a write error is showing | `tab next   ⏎ try again   esc done   ? help   q quit` |

> **`←` `→` are East-Asian-Ambiguous** (U+2190 / U+2192, verified with `unicodedata`, Unicode
> 16.0.0, 2026-08-25), like `↑↓`. The footer is a flowing line, so this changes its length and not
> a column — but it **must** be measured with the width-aware function of §1.2. ASCII fallback:
> `left/right state`.

---

## 10 · The exact copy — ready to paste

| Slot | String |
|---|---|
| breadcrumb | `Zerado ✦ Add a game by hand` |
| title | `ADD A GAME BY HAND` |
| intro line | `Nothing here syncs. A disc or a cartridge is a row like any other.` |
| field 1 label | `TITLE` |
| field 1 empty error | `A title is the one thing Zerado can't invent.` |
| field 2 label | `PLATFORM` |
| field 2 help — library empty | `Whatever you call it. PS2, Switch, Mega Drive.` |
| field 2 help — platforms exist | `Already here: PS2, Switch.` |
| field 2 empty error | `Which machine? PS2, Switch, whatever you call it.` |
| field 3 label | `STATE` |
| field 3 options | `NOT STARTED` · `IN PROGRESS` · `ZERADO` · `ABANDONED` — ratified, never re-worded |
| field 4 label | `OWNED SINCE` |
| field 4 help | `Optional. A year is enough.` |
| duplicate note | `You already have Shadow of the Colossus on PS2. Adding another is fine.` |
| receipt — library has rows | `Shadow of the Colossus is in the library. 248 games.` |
| receipt — first ever row | `Shadow of the Colossus is in the library. 1 game.` |
| footer, default | `tab next   ⏎ add   esc done   ? help   q quit` |
| footer, `STATE` focused | `←→ state   tab next   ⏎ add   esc done   ? help   q quit` |

**W · the write failed** — the one error `Z-08` can raise
```text
   ▌ NOT SAVED

     Zerado could not write to the library file.
     ~/.local/share/zerado/library.db: read-only file system

     Nothing you typed has been lost. Fix the file and press ⏎.
```

Names **what happened** (could not write), **why** (the OS's own reason, verbatim, not
paraphrased), **the next action** (fix the file, press `⏎`) and **what happened to the player's
data** ([`../01-design-system.md`](../01-design-system.md) §11.3). It does not say
*"Something went wrong."*

**On the receipt's plural.** `1 game.` and `248 games.` are two strings, not one with a suffix
bolted on. Say the number, and say it in the language people use.

**On truncation.** A title longer than the receipt's room is truncated with the product's `…`
(U+2026), which is **East-Asian-Ambiguous** — so the truncator must include the ellipsis in its
width-aware measurement, not append it after measuring. §1.2 rule 2 already requires this; it is
restated because the ellipsis is the exact character people forget to measure. **Flagged in §17.**

**Voice check.** No exclamation marks · no emoji · never "gamer" · the number is said, with the
right plural · nothing contradicts `landing-copy.md` §06 (*"title, platform, done"* — which is
exactly the required-field set, and *"A physical copy isn't a second-class row in the list"*, which
§8.1 makes checkable).

---

## 11 · 40-column behaviour, and the refusal floor

### 11.1 · RENDER 40×24 — Narrow · `leftInset` 2 · body `36 × 16`

Responsive table: *"Help collapses."* The boundary glyphs give way to an underline, exactly as
[`Z-02-connect-a-store.md`](./Z-02-connect-a-store.md) §11.1. **The `STATE` chips wrap to two
rows** — they are never compressed, and the label is never dropped (§3.6: *"The chip does not
shrink and the label is never dropped"*).

```text
0.........1.........2.........3.........

  Zerado ✦ Add a game by hand

  ADD A GAME BY HAND

  TITLE                            ?
  ▌ Shadow of the Colossus
    ──────────────────────────────────

  PLATFORM                         ?
    PS2
    ──────────────────────────────────

  STATE
    ▌ ○  NOT STARTED
      ◐  IN PROGRESS

  tab next  ⏎ add  esc  ? help

```

The intro line and `OWNED SINCE` scroll below the fold; `tab` reaches them and the body scrolls to
follow — cursor-following scroll, the same discipline R-10(b) puts on the ledger. `←→` still moves
between all four options; the two off-screen rows scroll into view.

### 11.2 · Standard `60 × 24` · body `54 × 16`

Wide's composition at 54 columns. The `STATE` row is **65 columns** at Wide, so at Standard it
wraps to **two rows of two chips** — `65 > 54`, stated as arithmetic rather than as a judgement.

### 11.3 · Tiny `< 40` — `32 × 24` · body `30 × 21`

Responsive table: *"One field at a time."* Band is the title row only.

```text
0.........1.........2.........3

 ADD A GAME BY HAND

 TITLE                     ?
 ▌ Shadow of the Colossus
   ──────────────────────────

 field 1 of 4

 tab next  ⏎ add  q quit
```

The `STATE` field at Tiny is one option per row, four rows, `←→` or `↑↓` moving between them —
**the label survives at every tier, by growing the shape rather than shrinking the chip.**

### 11.4 · The refusal floor — below `24 × 8`

`Z-08` never renders; see [`Z-11-fatal-error.md`](./Z-11-fatal-error.md) §3.4. Anything typed and
not yet added is lost, because it was never saved — which is the same guarantee `⏎` gives and the
reason `⏎` is one keystroke.

---

## 12 · `NO_COLOR` rendering — shown, not asserted

Zero SGR. Bold goes with it. The characters are unchanged:

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Add a game by hand

   ADD A GAME BY HAND

   You already have Shadow of the Colossus on PS2. Adding another is fine.

   TITLE
      │ Shadow of the Colossus                                              │

   PLATFORM
      │ PS2                                                                 │
     Already here: PS2, Switch.

   STATE
   ▌   ○  NOT STARTED    ◐  IN PROGRESS  ▌ ◉  ZERADO    ⊘  ABANDONED

   OWNED SINCE
      │ 2004                                                                │
     Optional. A year is enough.

   ←→ state   tab next   ⏎ add   esc done   ? help   q quit

```

| Information | Channel that survives |
|---|---|
| **which state is selected** | the `▌` marker before `◉  ZERADO` — **position** |
| what each state *is* | the **glyph** `○ ◐ ◉ ⊘`, a single progression: empty ring, half filled, solid core, struck through — **and the label beside it**, which in a TTY *is* the text alternative (4.1.2) |
| which field has focus | the `▌` at column 4 **and** the heavy `┃` against the light `│` |
| that a field is optional | the word `Optional.` |
| that a duplicate exists | the sentence, naming both the title and the platform |

**No information is lost.** With `ZERADO_ASCII=1` as well, the chips become
`[ ] NOT STARTED   [~] IN PROGRESS   [*] ZERADO   [x] ABANDONED` — entirely Narrow, immune to the
ambiguous-width setting, and still colour **and** glyph **and** label.

---

## 13 · Focus model on this screen

| | |
|---|---|
| **Regions** | **1** — the field group |
| **Items** | four fields, in the order drawn |
| **Initial focus** | `TITLE` — on entry, and again after every successful add |
| **Traversal** | `tab` / `shift-tab`, `↓` / `↑`; wraps. Within `STATE`, `←` / `→` |
| **Focus is never nowhere** | the field list is fixed and non-empty |
| **Focus ring** | present on the focused field in **every** state, including while a write is in flight and while an error is showing. Never removed ([`../02-colour-budget.md`](../02-colour-budget.md) §8.2) |
| **`Esc`** | **pop the route.** One press, one meaning, in every state |
| **Restored on pop** | returning to `Z-04` restores the row the player left; the newly added row is present and, if the sort puts it there, is where they will find it |

`Z-08` has **no running-operation state to cancel** — a write is a single local transaction that
completes in microseconds — so `Esc` has only one job here, unlike on `Z-02` (§13.1 there) and
`Z-03` (§13 there). The same departure from
[`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §5's *"blur the input"*
row applies and for the same reason: `R = 1`, and no screen-level single-key shortcut exists to
return to. **Flagged in [`Z-02-connect-a-store.md`](./Z-02-connect-a-store.md) §17.**

**No keyboard trap** (2.1.2).

---

## 14 · Colour budget declaration

| Class | Count | Where |
|---|---|---|
| **CHROME cyan** | **1** | `⏎ add` in the footer — the one key hint the screen most wants pressed. It becomes `⏎ try again` on a write error and is still the same one mark |
| **STATE cyan** | **unbounded, uncounted** | the `◉` glyph and the `ZERADO` label in the state chip qualify explicitly ([`../02-colour-budget.md`](../02-colour-budget.md) §2.1) — and here the chip is cyan **only when selected**, so the usual count is **0** |
| **Focus-ring cyan** | exempt | the focused field's boundary — §2.3 |
| **Text-cursor cyan** | exempt | the terminal's own caret |
| **Amber** | 3 marks, +1 when `IN PROGRESS` is selected, +1 with audio on | title (allow-list **1**) · `▌` field cursor (**5**) · the `IN PROGRESS` chip (**3**) · `▮ AUDIO` (**9**) |
| **Orchid** | 1 when `ABANDONED` is selected | `--z-state-abandoned` is a **state** colour, not on the amber or cyan budgets at all |
| **Red** | **0**, except the write-error annunciator (§5 list **2**, and see [`Z-02-connect-a-store.md`](./Z-02-connect-a-store.md) §14.1) | no scanner on this screen — there is no indeterminate wait |

**Amber ceiling:** title 18 + `▌` 1 + a selected `IN PROGRESS` chip 14 = **33 cells** of
`80 × 24 = 1920` → **1.7 %**.

**The tempting violation, avoided.** The `STATE` row is the one place in this cluster where four
state colours sit side by side. Cyan appears there as **data** — a state cell — which §2.1 exempts
by name. It would be a violation if the row's *label* `STATE`, or the `⏎ add` hint, borrowed cyan
because the cyan chip is nearby.

---

## 15 · Reuse verdict per element

| Element | Verdict | Note |
|---|---|---|
| The form | **`huh` — fits**, restyled to Zerado tokens | inherits deliberate title→field→help breathing room. **`huh`'s default theme must not ship** |
| Text inputs | **`bubbles/textinput`** — direct fit | §7.7 |
| `STATE` picker | **`huh` Select, restyled — but the option renderer is the Zerado state chip** | the chip is **built fresh** (§3.7 — no `bubbles` primitive fits, and none should be forced). `huh`'s own selected/unselected styling is replaced wholesale, because its indicator is not `▌` and its palette is not Zerado's |
| Field boundary | **Build fresh** — `lipgloss`, `--z-border-strong` / `--z-focus-ring` | shared with `Z-02` |
| Duplicate note / receipt | **Build fresh** — one `lipgloss` row that swaps content | it is the same row as the intro line, doing three jobs at three moments |
| Error block | **Build fresh. Not `charmbracelet/log`** | §11.5 |
| Header band, footer, audio indicator | **Build fresh**, shared | §2.8 · §5.2 |
| Audio cue | **none** | the cue list is closed (§15.3) and an add is not on it. A per-action sound here would be the audible form of the ambient-decoration ban |
| `harmonica` | **not used** | no motion on this screen — and no scanner, because nothing here is an indeterminate wait |

---

## 16 · Screen-specific acceptance criteria

1. **No field asks for anything Steam-shaped.** Grep the render for `app id`, `appid`, `store`,
   `url`, `id`. There must be no match. *A cartridge has no app ID.*
2. **Title and platform are the only required fields.** Submit with those two and nothing else; the
   row must be created.
3. **The result is indistinguishable in rank** (§8.1). Add a row, open `Z-04`, and assert against a
   synced neighbour: same row height, same chip geometry, same title column, same sort key, same
   focusability, same contribution to the pinned counts. The **only** permitted difference is the
   source column.
4. **The playtime cell is empty, not `0h`.** Grep the resulting `Z-04` row for `0h`; there must be
   no match.
5. **The form survives a batch.** Add ten rows in a row; `PLATFORM` retains its value every time,
   `TITLE` is empty and focused every time, and the receipt names each new title with the correct
   running count and the correct plural.
6. **A duplicate is not blocked** (§8.3), and the note names both the title and the platform.
7. **A write error loses nothing.** Force `AddManual` to fail; every field still holds what was
   typed, and the message names the OS's own reason verbatim.
8. **The state label is never dropped**, at any tier — including the two-row wrap at Standard and
   the one-per-row form at Tiny.
9. **`ZERADO_ASCII=1` renders `[ ] [~] [*] [x]`** and the row's column arithmetic is unchanged.
10. **Ambiguous-width proof.** Render at `80 × 24` with the terminal set to
    `ambiguous-width=double`. The `STATE` row's four option blocks must start at the **same
    columns** as in the single-width render — that is what the fixed 2-column glyph field is for.
11. **No sound is emitted on an add** (§8 row 12).
12. **Chrome-cyan count is exactly 1** at every colour depth, in every state — including the state
    where `◉ ZERADO` is selected and cyan appears twice on screen, once as data and once as chrome.
13. Eight artifacts per [`03-responsive.md`](../../blueprint/03-responsive.md) §7, **plus** one at
    `80 × 24` with each of the four states selected — the CVD and 16-colour evidence for the chip
    lives here, because this is the Phase 1 screen that renders all four at once.

---

## 17 · Open for the founder

1. **`notes` is not collected in Phase 1 (§3.1).** The column exists; no Phase 1 screen renders it
   back. Collecting a value the product never shows again is worse than not collecting it. Confirm
   — or add the field here **and** a block in `Z-05` in the same breath, which is what would make
   it honest.
2. **`Z-08` always writes `status_manual`, never `NULL` (§8.2).** For a provider with no playtime
   capability the two are indistinguishable forever, so this stores what the screen showed. Confirm
   against `Z-06`'s *clear override* item.
3. **Physical rows may be double-marked.** The row carries `PHY` in the source column at Wide and
   ExtraWide, **and** [`../01-design-system.md`](../01-design-system.md) §4.6 also specifies a
   ` · physical` title suffix. Both at once is redundant and it puts a second thing in the identity
   column. Proposal: **the suffix appears only where the source column has been shed** (Standard
   and below). This is `Z-04`'s spec to settle — routed there, flagged here because `Z-08`'s
   promise is what it affects.
4. **`STATE` or `STATUS`?** The product says *state* where it renders the four values (the state
   chip, the four game states) and *status* where it stores or commands them (`status_manual`,
   `Z-06 Set status`, `s status`). This field is labelled `STATE` because it shows the four values.
   Confirm the split is deliberate, or pick one word for both.
5. **The truncation ellipsis is Ambiguous-width (§10).** `…` U+2026 must be inside the width-aware
   measurement, not appended after it. Worth a line in
   [`../01-design-system.md`](../01-design-system.md) §1.2, because it is the character everyone
   forgets to measure. Route to `fft-tui-architect`.
