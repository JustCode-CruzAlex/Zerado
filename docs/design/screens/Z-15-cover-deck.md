---
title: Zerado — Z-15 Cover deck
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-15
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-15 · Cover deck

> The visual mode of the library. **Moved into Phase 1 by founder direction, 2026-08-25** —
> *"Images on terminal is a must, I know we will rely on kitty/ghostty but we need to show
> images, without image is not an option. Starting on phase 1."*
>
> And, on the terminals that cannot: *"yes it should work on other terminal emulators that will
> not have images, so we show what we can, and a text saying to use ghostty or kitty to better
> experience."*
>
> **Both halves are the screen.** A terminal that draws pictures gets pictures. A terminal that
> does not gets the whole product, correct, plus one quiet line — said once, dismissible, never
> again.

**Canon that governs this screen, as pointers, not copies:**
[`../00-design-brief.md`](../00-design-brief.md) (§10 is the bar) ·
[`../01-design-system.md`](../01-design-system.md) (§1.2 width · §2 band · §3 chip · §5 status
bar · §12 banner) · [`../02-colour-budget.md`](../02-colour-budget.md) (§10 is the second bar) ·
[`../03-designer-manual.md`](../03-designer-manual.md) §3 (this document's contract) ·
[`./Z-04-library.md`](./Z-04-library.md) (**the host — this is a mode of it**) ·
[`./Z-07-filter-and-search.md`](./Z-07-filter-and-search.md) (the other mode of the same host) ·
[`../../blueprint/02-composition.md`](../../blueprint/02-composition.md) §1.3 · §2.2 · §2.3 ·
§3 (**binding**) · [`../../blueprint/03-responsive.md`](../../blueprint/03-responsive.md) ·
[`../../blueprint/07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2 · §3 · §4 ·
[`../../blueprint/17-images.md`](../../blueprint/17-images.md) (**binding — the spine's image
contract. It owns the protocol, detection, the `Images` seam and the cache; this spec owns only
what the player sees**) · [`../../blueprint/16-i18n.md`](../../blueprint/16-i18n.md) (§11's
strings are catalogue entries, not literals) ·
FlowForge TUI Design Manual **R-10** · Spacing Canon #2435.

---

## 1 · Identity

| | |
|---|---|
| **ID** | `Z-15` |
| **Name** | Cover deck |
| **Phase** | **1** — moved from Phase 2 by founder direction ([`01-screen-inventory.md`](../../blueprint/01-screen-inventory.md) §2) |
| **Kind** | **Mode of `Z-04`**, not a route. Nothing is pushed; the stack does not change; `Esc` still does what `Esc` does on the root |
| **Route in** | `v` on `Z-04` |
| **Route out** | `v` again → the list. `⏎` → `Z-05` (≤ 119) or focus into the pane (≥ 120) · `s` → `Z-06` · `/` → `Z-07` (the two modes compose, §9) · `,` → `Z-09` · `?` → `Z-10` |
| **Shape** | **Ledger** — shape 1 of the five ([`02-composition.md`](../../blueprint/02-composition.md) §3), rendered as a grid instead of a list |
| **Composition** | ≤ 119 cols: single-pane grid · ≥ 120 cols: **grid ∥ detail**, the host's own composition, unchanged |
| **Offline class** | **DEGRADES** ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2). Cached covers render; uncached ones say so. The library itself is `WORKS` and does not change |

### 1.1 · It is a mode of `Z-04` — confirmed against the composition doc, and I agree

[`02-composition.md`](../../blueprint/02-composition.md) §3 already says it, and it is the right
call: *"`Z-15` Cover deck — Mode of Z-04, grid instead of list — Yes, a mode swaps the body
renderer, not the frame."*

**Confirmed. The frame, the header band, the breadcrumb, the title, the pinned row and the
reserved footer row are `Z-04`'s and are not touched.** A mode does not rename the screen: the
breadcrumb stays `Zerado ✦ Library` and the title stays `LIBRARY`, exactly as `Z-07` does. Only
the body renderer swaps.

**One qualification, and it is not a disagreement — it is a completion.** The brief describes the
pinned summary as *unchanged*. It is unchanged in **position, in rank and in what it must say**;
its **form** moves one rung down the ladder `Z-04` §10.3 already declares, because the deck needs
the freed columns for the identity row (**D-15-4**). That is the ladder working as specified, not
a new rule — and `Z-07`, the other mode of this host, already replaces the summary row outright.
A mode that could not touch the pinned row would be a stricter rule than the one `Z-07` ships
under.

**Two things the composition doc has not caught up on, routed rather than fixed here** (§17):
`02-composition.md` §2's Phase 1 table has no `Z-15` row and §3 still lists it under later
phases; and [`03-responsive.md`](../../blueprint/03-responsive.md) §3's per-screen degrade table
has no `Z-15` line. Both are `fft-tui-architect`'s.

---

## 2 · Purpose

**One sentence:** the same library, shown as the art it came with — and on a terminal that cannot
draw pictures, the same library, shown whole, with one line saying where pictures would work.

---

## 3 · Mockup — 80 × 24, the design floor and the primary breakpoint

`tier = Wide` · `leftInset = 3` · **body = 74 × 16** · every visible row begins at **column 4**.
Header-left equals content-left. Drawn to exact cell count; the outer rule is the terminal edge.

> **`░` is a mockup device and never renders.** It marks the cells the Kitty graphics protocol
> places an image over, so the geometry can be counted on paper. The real render is the picture,
> or the type tile of §3.3. Nothing in Zerado ever draws `░`.
>
> `░` U+2591 LIGHT SHADE is **Neutral** width (verified, UCD 17.0) and is used alone. `█` U+2588
> is **Ambiguous** and the two are never paired — [`../01-design-system.md`](../01-design-system.md)
> §1.2 forbids it.

### 3.1 · The tile — and the arithmetic that sizes it

**The co-render rule sizes the tile. The artwork does not.** This is the whole design problem of
the screen and it resolves in one line of arithmetic:

The ratified state chip is **14 columns** — glyph field 2 + gap 1 + label field 11, where
`NOT STARTED` and `IN PROGRESS` are the longest at 11
([`../01-design-system.md`](../01-design-system.md) §3.1). The label is **never dropped for
density** (§3.6, and anti-pattern 1). A tile therefore cannot be narrower than **16** columns —
the chip plus the 2-column focus field — or the state stops co-rendering.

```
 tile — 17 × 6, identical at every tier
 ┌─ 17 ─────────────┐
 │░░░░░░░░░░░░░░░░░ │  ┐
 │░░░░░░░░░░░░░░░░░ │  │  the image box — 17 × 4 cells
 │░░░░░░░░░░░░░░░░░ │  │  17 cells wide × 4 cells tall ≈ 2.125 : 1
 │░░░░░░░░░░░░░░░░░ │  ┘
 │▌ Return of the…  │     caption A — focus field 2 + title 15
 │  ◉  ZERADO       │     caption B — 2 + the ratified 14-column chip
 └──────────────────┘
```

| Field | Cols | Cell range | Note |
|---|---|---|---|
| **image box** | **17 × 4 rows** | rows 1–4 | see the aspect note below |
| caption A · focus field | **2** | 1–2 | `▌` U+258C is **Ambiguous** — padded to 2 measured cells ([`../01-design-system.md`](../01-design-system.md) §1.2 rule 1) |
| caption A · title | **15** | 3–17 | truncated at a measured budget, never a counted one — §3.6 |
| caption B · dead field | 2 | 1–2 | keeps the chip's glyph under the title's first character |
| caption B · **state chip** | **14** | 3–16 | glyph field 2 + gap 1 + label field 11. **Never compressed, never abbreviated** |
| caption B · pad | 1 | 17 | |
| | **17 × 6** | | |

> **D-15-1 · The tile is 17 × 6 at every tier, and only the grid's column and row counts change.**
> A box that breathes differently at different widths looks resized, not designed — the same
> reason `Z-06`'s overlay is bound at `34 × 11` and `Z-05`'s `BorderInsetX` is fixed at 2
> ([`Z-06-set-status.md`](./Z-06-set-status.md) D-06-1). One geometry also means one golden
> file per tier instead of one per tier per tile shape.

> **D-15-2 · The image box is landscape, and the co-render rule is what makes it landscape.**
> A terminal cell is roughly twice as tall as it is wide, so an image rendered `C` columns wide
> occupies about `C × (H/W) × 0.5` rows. Run that against the two candidate artworks:
>
> | Artwork | Aspect `W:H` | Rows at 17 columns | Tile height | Tile rows at the 80 × 24 floor |
> |---|---|---|---|---|
> | **Store header capsule** (460 × 215) | 2.140 : 1 | **4.0** | **6** | **2** — 8 covers |
> | Portrait box art (600 × 900) | 0.667 : 1 | **12.8** | 15 | **0** — it does not fit at all |
>
> The grid region at the floor is **13 rows** (§3.5). Portrait box art at a tile width the chip
> can live in needs **15** rows for a single row of tiles, so a portrait deck **cannot render at
> 80 columns at all** — which inherited verdict 2 forbids outright
> ([`../03-designer-manual.md`](../03-designer-manual.md) §5.11: *"No screen may require more
> than 80 columns to be correct"*).
>
> **A 17 × 4 cell box is 2.125 : 1 — within 0.7 % of the store header capsule.** So the artwork
> Zerado asks for is the artwork that fits, exactly, with no letterbox and no crop.
>
> *The alternative considered and rejected:* portrait art with the caption **beside** it rather
> than below — image 8 × 6 (a true 2 : 3), caption 14 columns, tile 24 × 6. It composes to
> **2 across × 2 down = 4 covers** at the floor against this design's **8**, and it spends 8
> columns of every tile on a gutter. Recorded so it is not re-proposed as "but covers are
> portrait."

### 3.2 · RENDER 80×24 — the deck, every cover fetched

`4 across × 2 down = 8 covers.` Focus is on tile 1.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31  Return of the Obra Dinn                   │
│                                                                                │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ▌ Return of the…     Sekiro: Shadow…    Sid Meier's Ci…    Signalis          │
│     ◉  ZERADO          ⊘  ABANDONED       ○  NOT STARTED     ◐  IN PROGRESS    │
│                                                                                │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│     Slay the Spire     Stardew Valley     Subnautica         Tunic             │
│     ◉  ZERADO          ○  NOT STARTED     ⊘  ABANDONED       ◐  IN PROGRESS    │
│   COVERS  17–24 of 247                                                         │
│   ↑↓←→ move  ⏎ open  s status  / filter  r sync  v list  ? help  q quit        │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Character counts, stated because a mockup that does not add up is the defect this spec exists
to prevent:**

| Line | Cells |
|---|---|
| Terminal | **80 × 24** |
| Body content, every body row | **74** |
| Identity row `247 games  ○ 198 … Return of the Obra Dinn` | 58 |
| — of which the summary (form 3) | **33** |
| — gutter | 2 |
| — the focused game's title, field | **39** |
| Every tile band row | **74** — `4 × 17 + 3 × 2 = 74`, exact |
| Position row `COVERS  17–24 of 247` | 20 |
| Footer key line | **69** |
| Header band | rows 2–4 (breadcrumb · gap · title), `HeaderBand(Wide, false)` = 3 |

### 3.3 · RENDER 80×24 — a mixed deck, and the three type tiles

Three tiles have no picture, for three different reasons, and each says which in a word. This is
the ordinary state of a real library, not an error state.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31  Return of the Obra Dinn                   │
│                                                                                │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░    not fetched      ░░░░░░░░░░░░░░░░░   │
│   ▌ Return of the…     Sekiro: Shadow…    Sid Meier's Ci…    Signalis          │
│     ◉  ZERADO          ⊘  ABANDONED       ○  NOT STARTED     ◐  IN PROGRESS    │
│                                                                                │
│   ░░░░░░░░░░░░░░░░░                                        ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░                                        ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░                                        ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░    nothing found      added by hand    ░░░░░░░░░░░░░░░░░   │
│     Slay the Spire     Stardew Valley     Subnautica         Tunic             │
│     ◉  ZERADO          ○  NOT STARTED     ⊘  ABANDONED       ◐  IN PROGRESS    │
│   COVERS  17–24 of 247                                                         │
│   ↑↓←→ move  ⏎ open  s status  / filter  r sync  v list  ? help  q quit        │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**The word sits on the image box's last row**, at tile column 3, so it clusters with the title
and the chip below it and leaves the air where the picture would go. That air is information: it
is the shape of the missing thing, and it is why a type tile never reads as a broken-image box.

> **D-15-10 · Three words, and `fetch failed` is deliberately not one of them.**
>
> | The fact | The word | Cells | Can `r` change it? |
> |---|---|---|---|
> | Zerado has not asked for this cover yet | `not fetched` | 11 | **yes** |
> | Zerado asked; the provider has no art for this game | `nothing found` | 13 | no |
> | This game came from the player's own shelf | `added by hand` | 13 | no — and nothing is wrong |
>
> A **failed** fetch is not a fourth word. A failure is a property of the network or the
> provider, not of the game, and it is a **screen-wide** condition — so it is carried by the
> banner set `Z-04` §8.1 already owns, and the tile keeps saying the only thing it knows, which
> is `not fetched`. Inventing a per-tile `fetch failed` would print a screen-wide fact 247 times
> and would tell the player the game is broken when the connection is.
>
> This follows [`Z-05-game-detail.md`](./Z-05-game-detail.md) D-05-5 exactly: the distinctions
> that matter are carried **in words**, and only the ones the player can act on differently get
> different words.

### 3.4 · The grid formula — no magic numbers

```
tileW = 17   tileH = 6   gutterX = 2

cols   = max(1, floor((BodyRect.w + gutterX) / (tileW + gutterX)))
gridH  = BodyRect.h − (pinned rows above) − 1 (position row)

rowsWithGap = floor((gridH + 1) / (tileH + 1))
rowsNoGap   = floor( gridH      /  tileH     )
rows        = max(rowsWithGap, rowsNoGap)
gutterY     = 1 if rowsWithGap ≥ rowsNoGap else 0
```

The grid is **left-aligned at `leftInset`** and any remaining columns are left unused, so
header-left still equals content-left — the canon's headline invariant (#2435 §5).

| Tier | `BodyRect` | `cols` | cols used | spare | `gridH` | `rows` | `gutterY` | **covers** |
|---|---|---|---|---|---|---|---|---|
| Narrow `40` | 36 × 16 | 2 | **36** | 0 | 13 | 2 | 1 | **4** |
| Standard `60` | 54 × 16 | 2 | 36 | 18 | 13 | 2 | 1 | **4** |
| Standard `79` | 73 × 16 | 3 | 55 | 18 | 13 | 2 | 1 | **6** |
| **Wide `80`** | **74 × 16** | **4** | **74** | **0** | **13** | **2** | **1** | **8** |
| Wide `80`, banner up | 74 × 16 | 4 | 74 | 0 | **12** | 2 | **0** | **8** |
| ExtraWide `120` | 112 × 32 | 3 *(the deck's 66)* | 55 | 11 | 29 | 4 | 1 | **12** |

> **D-15-5 · The grid's row gutter is the flexible row, and it is the only flexible thing on the
> screen.** `gutterY` is **1** whenever the arithmetic allows it and **0** when spending it buys
> a whole extra row of tiles. It never buys a fraction, and the tile never changes.
>
> Without this rule a degrade banner would cost **half the deck**: 13 grid rows carry two tile
> rows, 12 carry one, and the six rows freed would sit blank at the bottom of the screen. With
> it, **a banner costs zero covers.** The shape of the rule is inherited from
> [`Z-04-library.md`](./Z-04-library.md) **D-04-5**, where below Standard the leftover row
> becomes the respiro rather than an orphan half-row; this is the same trade in the other axis.
>
> At `gutterY = 0` a tile row's chip line sits directly above the next row's image box. A word in
> the state colour against a picture edge is not a collision — the two are different kinds of
> mark and the reading order is unambiguous. It is the tighter of the two renders, and it is only
> ever reached when something more important than air is on screen.

> **The Standard tier's 18 spare columns, named rather than hidden.** At exactly 60 columns the
> body is 54 and the grid uses 36 of it. *The alternative considered and rejected:* stretch the
> tile to `(54 − 2)/2 = 26` so the grid fills the width. It was rejected because it changes the
> tile's shape between tiers, which **D-15-1** forbids for the reason D-06-1 gives, and because
> the stretch would either letterbox the artwork or crop it. The spare columns are honest; a
> resized picture is not. At 79 columns — the same tier — the grid takes 3 across and the spare
> is the same 18, which is the formula behaving consistently rather than a cliff.

### 3.5 · Row map — 80 × 24

| Terminal row | Content | Token that puts it there |
|---|---|---|
| 1 | blank | `OuterMarginY` = 1 |
| 2 | breadcrumb `Zerado ✦ Library` | `HeaderBand` row 1 |
| 3 | blank | `InterElementGap` = 1 (inside the band) |
| 4 | title `LIBRARY` | `HeaderBand` row 3 |
| 5 | blank | `InnerPaddingY` = 1 |
| 6 | blank | `InterElementGap` = 1 (band → body) |
| **7** | **body 1** — pinned summary **+ the identity row** | `BodyRect` · **D-15-4** |
| 8 | body 2 — respiro | `Z-04` D-04-2 |
| 9–14 | body 3–8 — **tile row 1** | 4 image rows + caption A + caption B |
| 15 | body 9 — `gutterY` | |
| 16–21 | body 10–15 — **tile row 2** | |
| 22 | body 16 — position row | pinned **outside** the grid — R-10(c) |
| 23 | footer key line | the canon's **reserved footer row** |
| 24 | blank | `OuterMarginY` = 1 |

> **D-15-4 · The caption is the tile's label; the pinned row is the identity column.**
>
> The caption's title field is **15** columns. [`02-composition.md`](../../blueprint/02-composition.md)
> §2.1 already ruled on that number: *"Fifteen columns is `Return of the Ob…`, which fails
> R-10(a) on its face."* A cover deck whose only title is 15 columns wide is therefore a deck
> that fails the ledger triad — so the untruncated title has to live somewhere, and it has to be
> somewhere that is **always** on screen.
>
> **It shares the pinned row.** In cover mode the summary renders **form 3, the glyph key** —
> `247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31`, **33** cells — and the freed columns carry the focused
> game's full title, left-aligned at body column 36, in a **39**-column field at Wide and **77**
> at ExtraWide.
>
> **Form 3 is licensed here, not assumed.** `Z-04` **D-04-1**'s guard is that the glyph key may
> render *"only on a tier where the full state label is visible in the list below."* Every tile
> carries the ratified 14-column chip with its label spelled out, at every tier, so the guard
> holds by construction — and `05-state-machine.md` §7 rule 1 holds unchanged, because the four
> counts still describe the shown set.
>
> **Where the summary row cannot hold it, the identity row takes the respiro row instead** —
> below `33 + 2 + 26 = 61` columns of body, which is Narrow and Standard-at-60. **26** is not a
> new number: it is `Z-04`'s Tiny title field (§11.2), the narrowest identity column this bundle
> has already ratified. Spending the respiro is the one place `Z-04` **D-04-2** yields, and the
> reason is a rule that outranks it: [`03-responsive.md`](../../blueprint/03-responsive.md) §4
> puts *the game's title* on the **never-hide** list, and the respiro is not on any list.

### 3.6 · Truncation is measured, never counted

The caption's title is cut to `15 − measured_width("…")`, evaluated **at render time** with the
East-Asian-Width-aware function, then a trailing space is stripped and `…` appended.
`…` U+2026 is **Ambiguous**, so on an `ambiguous-width=double` terminal it is **two** cells and
the budget is 13, not 14. Baking the marker in as one cell is the exact failure
[`../01-design-system.md`](../01-design-system.md) §1.2 **rule 3** exists to prevent, and on a
grid it would shear every tile to the right of the truncated one.

Under `ZERADO_ASCII=1` the marker is `...` — three narrow cells, deterministic — and the focus
marker `▌` becomes `>`, per §1.2 rule 4.

---

## 4 · Mockup — 120 × 40, ExtraWide: deck ∥ detail

`leftInset = 4` · **body = 112 × 32** · split **66 ∥ 2 ∥ 44**
([`02-composition.md`](../../blueprint/02-composition.md) §2.3, **unchanged**).
The deck takes **3 across × 4 down = 12 covers**; the pane is `Z-05` in its second host.

**There is no identity row here** — the pane already renders the focused game's title
untruncated and wrapped, which is the property **D-15-4** exists to guarantee. The pinned row is
the summary alone.

```
┌────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                                        │
│    Zerado ✦ Library                                                                                                    │
│                                                                                                                        │
│    LIBRARY                                                                                                             │
│                                                                                                                        │
│                                                                                                                        │
│    247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31                                   ┌ DETAIL ──────────────────────────────────┐    │
│                                                                        │  Return of the Obra Dinn                 │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │  ◉  ZERADO                               │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │  PLAYTIME     9h                         │    │
│    ▌ Return of the…     Sekiro: Shadow…    Sid Meier's Ci…             │  LAST PLAYED  2 Aug 2026                 │    │
│      ◉  ZERADO          ⊘  ABANDONED       ○  NOT STARTED              │  ADDED        12 Mar 2026                │    │
│                                                                        │  SOURCE       Steam                      │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │  SET BY       you, 12 Aug 2026           │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │  STEAM SAYS   IN PROGRESS                │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│      Signalis           Slay the Spire     Stardew Valley              │                                          │    │
│      ◐  IN PROGRESS     ◉  ZERADO          ○  NOT STARTED              │                                          │    │
│                                                                        │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│      Subnautica         Tunic              Undertale                   │                                          │    │
│      ⊘  ABANDONED       ◐  IN PROGRESS     ◉  ZERADO                   │                                          │    │
│                                                                        │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│    ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░             │                                          │    │
│      Vampire Surviv…    Wasteland 3        The Witness                 │                                          │    │
│      ○  NOT STARTED     ⊘  ABANDONED       ○  NOT STARTED              │                                          │    │
│                                                                        │                                          │    │
│                                                                        │  LAST SYNCED  3 hours ago                │    │
│    COVERS  17–28 of 247                                                └──────────────────────────────────────────┘    │
│    ↑↓←→ move   ⏎ detail   tab pane   s status   / filter   a add   r sync   v list   ? help   q quit                   │
│                                                                                                                        │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Counts:** deck column **55** (`3 × 17 + 2 × 2`) · gap **13** · detail pane **44**
(border 1 + `BorderInsetX` 2 + content 38 + `BorderInsetX` 2 + border 1) = **112**.
Grid rows `4 × 6 + 3 × 1 = 27` of `gridH` 29 — **2 rows spare, and they sit below the grid**,
above the pinned position row. Footer **97** ≤ 112.

**The 11 columns the grid does not use land in the gap**, which is why the separation between the
deck and the pane reads as 13 columns rather than the list mode's 2. That is the one rule of
§3.4 — the grid is left-aligned and the pane is right-aligned — producing a wider respiro rather
than a ragged edge. Elevation is still carried by a border and spacing, never by fill
([`../02-colour-budget.md`](../02-colour-budget.md) §7.1).

> **D-15-7 · At ExtraWide the detail pane stays, and the deck takes three columns.**
>
> *The alternative, with its number:* give the deck the whole 112-column body and it composes
> **6 across × 4 down = 24 covers** with 1 spare column, carrying the identity row on the pinned
> row as at Wide. **Twice the covers** — and it is a real option worth the founder's eye.
>
> **It is not this spec's to take.** Dropping the pane takes `Z-04` from `R = 2` to `R = 1` at
> ExtraWide, unbinds `Tab`, and changes what `⏎` means (focus into the pane, versus pushing
> `Z-05` as a route) — three composition decisions, and composition belongs to
> `fft-tui-architect` ([`../03-designer-manual.md`](../03-designer-manual.md) §2). A mode may
> swap the body **renderer**; changing the region count is not a renderer change.
>
> Keeping the pane also pays for itself twice: it carries the untruncated title, so the deck
> needs no identity row at this tier, and it is the one place the cover's own provenance can be
> stated with its age. **Routed to `fft-tui-architect` in §17 and put to the founder in §18.**

---

## 5 · The terminal that cannot draw images — supported, not refused

### 5.1 · Capability detection — at startup, once, and never a flag the player has to find

| Rule | What it means here |
|---|---|
| **Detected, not configured** | Resolved **once at startup**, before first paint, off the render path, with a timeout. There is no setting the player must discover to turn covers on |
| **Target** | **Kitty graphics** (Kitty, Ghostty — the two the founder named — plus WezTerm and Konsole), and **iTerm2 inline images second** ([`../../blueprint/17-images.md`](../../blueprint/17-images.md) §2). Sixel is **not** adopted in Phase 1 and half-block art is **rejected outright** — *"it is not the image; it is a picture of a picture."* The protocol, the query and the cache are the spine's and `fft-api-designer`'s; **this spec does not invent them** and must not be read as specifying them |
| **Detection fails closed** | Silence, ambiguity or a timeout means *no image support*. Zerado never waits on it and never retries it. §2 of the spine is blunt about why: *"Guessing yes and emitting escape sequences into a terminal that does not understand them is how a library view turns into garbage on somebody's screen"* |
| **A missing cover is never worth a dropped frame** | The render path only ever reads what the cache already holds — `Cover` never fetches and never blocks ([`../../blueprint/17-images.md`](../../blueprint/17-images.md) §4, the same rule as `Audio.Cue`). **A cover that has not arrived is simply a type tile this frame**, and that is why §3.3's type tile is the ordinary render rather than an error path |
| **Never re-probed** | The result is a fact about this process. A player who switches terminals starts a new process |
| **`ZERADO_NO_IMAGES`** | The environment override, in the discipline `NO_COLOR`, `ZERADO_ASCII`, `ZERADO_NO_AUDIO` and `ZERADO_REDUCED_MOTION` already established: **when set, Zerado draws no image at all** — not fewer, none. It is the escape hatch for a terminal that claims the protocol and renders garbage, and it is **not** the opt-in; there is no opt-in, because detection is the opt-in. **Adopted by [`17-images.md`](../../blueprint/17-images.md) §4**, where `NullImages` is what it selects — so it is used here as settled, not proposed |

**What the player sees when there is no image support: the text deck, whole, and one line.**
The deck is not degraded, because the deck is not attempted. `v` renders `Z-04`'s ledger — every
row, every state, every count, the ratified render — and adds the note.

### 5.2 · RENDER 80×24 — the text deck and the note

The note occupies body row 1, where the degrade banner goes, and the pinned block grows downward
exactly as `Z-04` **D-04-2** specifies: **11 game rows, not 12**, and the respiro is not spent.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   ▌ COVERS   Cover art needs Ghostty or Kitty. Everything else here works.     │
│   247 games · 198 not started · 12 in progress · 6 zerado · 31 abandoned       │
│                                                                                │
│     STATE           TITLE                                        HOURS   SRC   │
│     ◐  IN PROGRESS  Baldur's Gate 3                                87h   STM   │
│     ○  NOT STARTED  Blasphemous                                     0h   STM   │
│     ◉  ZERADO       Celeste                                        14h   STM   │
│     ○  NOT STARTED  Chrono Trigger                                   —   PHY   │
│     ◐  IN PROGRESS  Dark Souls III                                 63h   STM   │
│     ⊘  ABANDONED    Disco Elysium                                  22h   STM   │
│     ◉  ZERADO       Hades                                          58h   STM   │
│     ○  NOT STARTED  Hollow Knight                                  41h   STM   │
│     ◐  IN PROGRESS  Outer Wilds                                    12h   STM   │
│   ▌ ◉  ZERADO       Return of the Obra Dinn                         9h   STM   │
│     ⊘  ABANDONED    Sekiro: Shadows Die Twice                       3h   STM   │
│   ROWS  4–14 of 247                                                            │
│   ↑↓ move  ⏎ open  s status  / filter  r sync  x dismiss  ? help  q quit       │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Counts:** note **72** · summary form 1 **70** · every game row **74** · position row **17** ·
footer **70**. The summary keeps its **full prose form** here, because this is the list and there
is no identity row to make room for — `Z-04`'s render, unaltered.

### 5.3 · The note — the exact copy, and the tone it has to hold

```
▌ COVERS   Cover art needs Ghostty or Kitty. Everything else here works.
```

**72 cells.** `▌` U+258C is **structure**, not a state channel — the same structural bar the
degrade banner uses and the same reason there is no `⚠` anywhere in Zerado
([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §3).

**Four rules the copy obeys, and each is the reason a draft was thrown away:**

1. **It answers; it never announces.** The note exists only because the player pressed `v`. It
   is not raised at startup, it is not raised by the library, and it is never raised by anything
   the player did not ask for. **This is the whole of "never a recurring scold" (D-15-6).**
2. **The subject is the feature, never the player's terminal.** *Cover art needs* — the sentence
   is about what a picture requires. Nothing in it is a fault, a lack, or a thing the player
   should have done differently, and the words *your*, *support*, *unsupported* and *install*
   never appear.
3. **It ends on the reassurance, and it is the ratified one.** `Everything else here works.` is
   the same clause the `OFFLINE` banner ends on — *"Everything here still works"*
   ([`Z-04-library.md`](./Z-04-library.md) B1). One product, one voice: when Zerado tells you a
   thing is absent, the next thing it tells you is that the rest is not.
4. **It is chrome, never amber and never red.** Amber is reserved for *action required*
   ([`../01-design-system.md`](../01-design-system.md) §12) and this requires no action; red is
   motion and alarm, and a terminal without a graphics protocol is neither
   ([`../02-colour-budget.md`](../02-colour-budget.md) §5.2 — the same reasoning that forbids a
   red `OFFLINE`).

**The dismiss key is not in the sentence.** `x dismiss` is the **screen's** affordance and lives
in the footer, per [`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §3.1 rule 1
— a message that names a key is a message that lies the moment a text input has focus.

**Drafts rejected, so the tone is falsifiable rather than a matter of ear:**

| Draft | Why it was thrown away |
|---|---|
| `Your terminal doesn't support images. Try Ghostty or Kitty.` | *Your* makes it the player's possession and, by a short step, their fault. *Try* is an instruction. *Support* is the vocabulary of a compatibility error, which this is not. [`17-images.md`](../../blueprint/17-images.md) §3 rejects this exact sentence |
| `NO IMAGES   Install Ghostty or Kitty to see cover art.` | *Install* asks for work in a note the player did not open. A label word beginning `NO` reads an absence as a failure |
| `Cover art needs a terminal that draws images, like Ghostty or Kitty.` | **79 cells** — over the 74-column body at Wide, before the label word |
| `This terminal can't draw images. Ghostty and Kitty can.` | **The draft this spec shipped first, and it lost to the spine's.** 66 cells, and its subject is still *this terminal* — a negative statement about the player's setup, however neutrally phrased. See the note below |
| Anything with *unfortunately*, *sorry*, *unsupported*, or an exclamation mark | [`../03-designer-manual.md`](../03-designer-manual.md) §5.8. The voice does not apologise and does not exclaim |

> **A reconciliation, recorded rather than quietly resolved.**
> [`17-images.md`](../../blueprint/17-images.md) §3 says *"The tone is the deliverable here"* and
> gives its own worked contrast: *"Your terminal does not support images"* reads as a fault;
> ***"Cover art needs Ghostty or Kitty"*** reads as information. An earlier draft of this section
> had **rejected** a longer relative of that sentence on two grounds — that it measured 79 cells,
> and that *needs* is a demand.
>
> **The first objection was about the wrong sentence** — the spine's short form is 33 cells and
> fits with room to spare; the 79 was my own padding. **The second did not survive the
> reassurance clause:** followed by `Everything else here works.`, *needs* reads as a requirement
> of the **picture**, not a demand on the **player**. The spine's phrasing is adopted verbatim
> and this spec supplies the second sentence, which is exactly the division
> [`17-images.md`](../../blueprint/17-images.md) §3 asks for — *the screen spec writes the exact
> copy; the spine's requirement is that it is dismissible and never returns.*
>
> Recorded because a spec that silently overrules a spine document it links to is worse than one
> that argues with it in the open.

> **The one genuinely consequential fork, put to the founder (§18).** The note names two
> **programs**. It does not name the **Kitty graphics protocol**, which is the term a player
> would have to search for to understand why one terminal draws covers and another does not.
>
> **This spec chooses program names in the note and the protocol name in the two durable
> homes** — `Z-10 Help`'s description of `v`, and `Z-09 Settings § DISPLAY` if the founder takes
> the `Images` row of §17. The reason is the bundle's own discipline: the note is the
> **interruption**, so it carries the fact a player needs in the second they read it; the places
> a player goes **looking** carry the jargon. Say it once, in the place someone would look
> ([`Z-09-settings.md`](./Z-09-settings.md) §8.2).
>
> **Naming two programs and not four is deliberate, and it is the founder's own list.** WezTerm
> and Konsole also implement Kitty graphics, and iTerm2 is adopted second
> ([`17-images.md`](../../blueprint/17-images.md) §2) — so a player on any of those four never
> sees this note at all. Listing every supported terminal would make the sentence longer and the
> reader's decision harder, and the founder named the two he named.

### 5.4 · Dismissal — once, and it stays dismissed

| | |
|---|---|
| **Key** | `x`, and only while the note is showing. It is in the footer for exactly that long |
| **`Esc`?** | **No.** `Esc` on the root means *nothing* and [`Z-04-library.md`](./Z-04-library.md) §13.4 documents that exhaustively. A note is not a mode and does not take focus, so binding `Esc` to it would make the root's `Esc` conditional — the one thing that table exists to prevent |
| **Persisted** | `setting('covers.note_dismissed', 'true')` — one row in the settings table ([`09-erd.md`](../../blueprint/09-erd.md) §1), written on the keystroke like every other mutation in this product. **It is never shown again**, in this session or any other |
| **Then `v` retires** | With the note dismissed on a terminal that draws no images, `v` is **unbound and absent from the footer**, and `Z-04`'s Wide footer returns to its ratified **73**-cell string. The precedent is `c` on `Z-04`, which is live only on an empty library (§9 there) |
| **Nothing is lost** | The fact stays available in `Z-10 Help`, which describes what `v` does and which terminals draw covers |

**Why `v` retires rather than staying bound and doing nothing.** A key that is listed and does
nothing is anti-pattern *"a footer that lies"*; a key that is bound and silent is the pattern
this bundle uses for the reserved Phase 2 keys, which are not listed either. Retiring it is the
same rule applied to a key the player has explicitly finished with.

### 5.5 · RENDER 80×24 — the note dismissed: `Z-04`, byte-identical

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   247 games · 198 not started · 12 in progress · 6 zerado · 31 abandoned       │
│                                                                                │
│     STATE           TITLE                                        HOURS   SRC   │
│     ◐  IN PROGRESS  Baldur's Gate 3                                87h   STM   │
│     ○  NOT STARTED  Blasphemous                                     0h   STM   │
│     ◉  ZERADO       Celeste                                        14h   STM   │
│     ○  NOT STARTED  Chrono Trigger                                   —   PHY   │
│     ◐  IN PROGRESS  Dark Souls III                                 63h   STM   │
│     ⊘  ABANDONED    Disco Elysium                                  22h   STM   │
│     ◉  ZERADO       Hades                                          58h   STM   │
│     ○  NOT STARTED  Hollow Knight                                  41h   STM   │
│     ◐  IN PROGRESS  Outer Wilds                                    12h   STM   │
│   ▌ ◉  ZERADO       Return of the Obra Dinn                         9h   STM   │
│     ⊘  ABANDONED    Sekiro: Shadows Die Twice                       3h   STM   │
│     ○  NOT STARTED  Sid Meier's Civilization VI: Gathering St…      0h   STM   │
│   ROWS  4–15 of 247                                                            │
│   ↑↓ move   ⏎ open   s status   / filter   a add   r sync   ? help   q quit    │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**This is `Z-04` §3, cell for cell** — 12 game rows, the 70-cell summary, the 73-cell footer with
`a add` restored. **The end state of a terminal without images is the ratified library screen and
nothing else.** That is the strongest statement this spec can make about *supported, not
refused*, and it is checkable: diff the two renders and they are identical.

---

## 6 · Visual hierarchy — what the eye reaches, in order

| # | What | Channel that carries it (case → weight → colour role → box drawing → spacing) |
|---|---|---|
| **1** | **The covers, as a block** | The only pictures on any Zerado screen. They are the reason the mode exists and they win by sheer area — `4 × 17 × 4 = 272` cells of a 74 × 16 body. **No Zerado mark competes with them, because none is allowed on them (§8.1)** |
| **2** | **`LIBRARY`** | UPPERCASE + bold + `--z-primary` amber, alone on its row with respiro above and below. It does not move when the mode changes, which is what tells the eye the screen did not change |
| **3** | **The focused tile** | The `▌` marker (position) + **bold** title (weight) + `--z-primary` on the marker (colour). The only bold text in the grid |
| **4** | **The focused game's full title** on the pinned row | `--z-text` **255** — the brightest text on screen, at a fixed column, changing only when the cursor moves |
| **5** | **The state chips** | A caption row under every tile. The `◉ ZERADO` chips are the only cyan Zerado draws, so they read by rarity — **against artwork that may be any colour at all**, which is why the chip sits on the terminal's own ground and never on the picture |
| **6** | The counts | Numerals in `--z-text` against `--z-text-secondary` glyphs and words — a weight step inside one line, no colour spend |
| **7** | Chrome | Breadcrumb, position row, footer, and the tile words `not fetched` / `nothing found` / `added by hand` — all quiet, all `--z-text-secondary` or `--z-text-tertiary` |

**The one thing a player should see first is the wall of covers; the one thing they should see
second is which one they are standing on.** Everything else is the same furniture the list mode
has, in the same places, deliberately.

**No border around a tile — considered and rejected.** A focused tile could be boxed. It was
rejected because a border costs 2 columns and 2 rows *of the tile*, which changes the tile's shape
between focused and unfocused and therefore reflows the grid on every cursor move; because
`--z-border` at 1.53:1 may not mark anything that carries meaning
([`../02-colour-budget.md`](../02-colour-budget.md) §8.1); and because focus already has its
three ratified channels and does not need a fourth invented for one screen
([`../01-design-system.md`](../01-design-system.md) §1.7).

**No caption on top of the artwork — forbidden outright, §8.1.**

---

## 7 · Every applied spacing token, by name, with its value at the rendered tier

| Token | Tiny `<40` | Narrow `40–59` | Standard `60–79` | **Wide `80–119`** | ExtraWide `120+` |
|---|---|---|---|---|---|
| `OuterMarginX` | 0 | 1 | 2 | **2** | 2 |
| `OuterMarginY` | 0 | 1 | 1 | **1** | 1 |
| `InnerPaddingX` | 1 | 1 | 1 | **1** | 2 |
| `InnerPaddingY` | 0 | 1 | 1 | **1** | 1 |
| `InterElementGap` | 1 | 1 | 1 | **1** | 1 |
| `HeaderBand(tier, false)` | 1 | 3 | 3 | **3** | 3 |
| **`leftInset`** | 1 | 2 | 3 | **3** | 4 |
| `BodyRect.w` | 30 | 36 | 54 | **74** | 112 |
| `BodyRect.h` | 21 | 16 | 16 | **16** | 32 |
| **cover mode available?** | **no** — §12 | yes | yes | **yes** | yes |

`hasSubtitle` is **`false`**, here and on every Zerado screen
([`02-composition.md`](../../blueprint/02-composition.md) §1.2).

**Applied, not merely quoted:**

| Surface | Token | Value at Wide |
|---|---|---|
| Frame inset, all four sides | `OuterMarginX` / `OuterMarginY` | 2 cols / 1 row |
| Inside the frame | `InnerPaddingX` / `InnerPaddingY` | 1 col / 1 row |
| breadcrumb → title | `InterElementGap` | 1 row |
| band → body | `InterElementGap` | 1 row |
| Header-left **and** content-left **and** the grid's left edge | `leftInset` | **column 4** |
| pinned block → grid | `InterElementGap` | 1 row (the respiro; **D-15-4** may claim it below 61 body columns) |
| Detail pane inset (ExtraWide) | `BorderInsetX` / `BorderInsetY` | 2 cols / 0 rows |
| Footer | the canon's reserved footer row | 1 row, not stolen from `BodyRect` |

**The non-token numbers on this screen, declared:** the tile's `17 × 6` and its `2`-column
`gutterX` and `1`-or-`0`-row `gutterY` (§3.1, §3.4 — a component geometry, derived from the
ratified 14-column chip and stated once); the `26`-column identity floor inherited from `Z-04`
§11.2; and the grid formula itself, which is a formula rather than a constant. **Zero magic
margins.** No spacing on this screen is picked by eye.

---

## 8 · Colour, glyph and label for every state shown

Ratios are read from the brand manual's measured table (§4.2). **None is estimated.**

| Shown as | Token | Hex | ANSI-256 | 16-colour | Glyph | ASCII | Label | Ratio |
|---|---|---|---|---|---|---|---|---|
| Not started | `--z-state-not-started` | `#A5A29B` | **247** | `white` | `○` U+25CB | `[ ]` | `NOT STARTED` | **7.62** AA |
| In progress | `--z-state-in-progress` | `#FFB000` | **214** | `bright yellow` | `◐` U+25D0 | `[~]` | `IN PROGRESS` | **10.59** AAA |
| Zerado | `--z-state-zerado` | `#19E0FF` | **45** | `bright cyan` | `◉` U+25C9 | `[*]` | `ZERADO` | **12.15** AAA |
| Abandoned | `--z-state-abandoned` | `#C77DFF` | **177** | `bright magenta` | `⊘` U+2298 | `[x]` | `ABANDONED` | **7.21** AA |
| Screen title `LIBRARY` | `--z-primary` | `#FFB000` | **214** | `bright yellow` | — | — | the word | **10.59** AAA |
| Focus marker `▌` | `--z-primary` | `#FFB000` | **214** | `bright yellow` | `▌` U+258C | `>` | — | **10.59** AAA |
| Tile title · the identity row | `--z-text` | `#E9EEF5` | **255** | `bright white` | — | — | the title | **16.65** AAA |
| Breadcrumb, summary words, position row, footer | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | — | — | — | **9.36** AAA |
| Summary numerals | `--z-text` | `#E9EEF5` | **255** | `bright white` | — | — | — | **16.65** AAA |
| Tile words `not fetched` · `nothing found` · `added by hand`, truncation `…` | `--z-text-tertiary` | `#8492A8` | ***underived*** | `white` | — | — | the words | **6.15** AA |
| **The capability note `▌`** | `--z-border-strong` | `#64748B` | **67** | `bright black` | `▌` | `>` | `COVERS` | **4.08** (1.4.11) |
| The capability note's sentence | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | — | — | the sentence | **9.36** AAA |
| Degrade banner `▌`, informational | `--z-border-strong` | `#64748B` | **67** | `bright black` | `▌` | `>` | the label word | **4.08** |
| Degrade banner `▌`, action required | `--z-primary` | `#FFB000` | **214** | `bright yellow` | `▌` | `>` | the label word | **10.59** AAA |
| Detail-pane border, focused (ExtraWide) | `--z-border-strong` | `#64748B` | **67** | `bright black` | `┏━┓` | `+-+` | — | **4.08** |
| Detail-pane border, unfocused | `--z-border` | `#2A3342` | **236** | `black` | `┌─┐` | `+-+` | — | 1.53 — **decoration only** |
| Audio annunciator, unmuted / muted | `--z-primary` / `--z-text-secondary` | `#FFB000` / `#A9B5C7` | **214** / **249** | `bright yellow` / `white` | `▮` / `▯` | *none — label alone* | `AUDIO` / `MUTED` | **10.59** / **9.36** |
| **The artwork** | **none — it is not a Zerado colour** | — | — | — | — | — | — | **§8.1** |

> ***underived*** means exactly that: **no ANSI-256 index has been derived for
> `--z-text-tertiary` `#8492A8`.** Nobody may pick one at the keyboard (brand §10 rule 5).
> **Interim rendering on this screen: the three tile words render uncoloured** — no SGR at all —
> which invents nothing and degrades identically to the `NO_COLOR` path. Owner:
> `fft-brand-architect`; tracked in [`../00-design-brief.md`](../00-design-brief.md) §9.

**The co-render rule holds on every tile at every tier: colour AND glyph AND label, all three.**
The chip is the ratified 14-column component, unchanged and uncompressed, on a caption row of its
own. Remove the colour and the glyph and the word remain; remove the picture entirely and
**nothing about the state is lost, because the state was never in the picture.**

### 8.1 · The artwork is content, not chrome — the rule that makes everything else hold

A cover is uncontrolled, full-colour, third-party artwork sitting inside a product whose entire
discipline is controlled colour. One rule resolves it, and it is binding:

> **Zerado never draws a mark on a cover.** No chip, no glyph, no badge, no title, no focus
> marker, no border, no dimming, no overlay, no tint. Every Zerado mark on this screen is on the
> terminal's own ground, outside every image rectangle.

**Five things this one rule buys:**

1. **Contrast is knowable.** Every ratio in the table above is measured against `--z-surface`
   ([`../01-design-system.md`](../01-design-system.md) §1.4). A chip drawn over artwork would
   have **no** computable ratio — the ground would be whatever the picture is that day — and the
   AA/AAA column would become a fiction. This is the reason, not an aesthetic preference.
2. **The colour budget stays countable.** [`../02-colour-budget.md`](../02-colour-budget.md)
   §3.1 counts cyan by scanning the **ANSI stream** for SGR runs and classifying them by
   payload. An image placed by the graphics protocol emits **no SGR run**, so it cannot register
   as chrome cyan, cannot register as amber, and cannot move the 10 % amber ceiling. The machine
   check is unaffected by artwork **by construction**, and that is falsifiable:
   **no SGR run may overlap an image cell rectangle.**
3. **The human check needs one added line, and only one.** §3.2's screenshot method becomes:
   *ignore the pictures; count Zerado's own marks.* Added to §15 below.
4. **The degrade is free.** A tile whose picture is missing needs no fallback design, because
   the picture was never carrying anything. The type tile of §3.3 is the same tile with the
   image cells empty — which is why *"we show what we can"* is a composition and not a
   consolation.
5. **`NO_COLOR` is coherent.** See §13.

*The alternative considered and rejected:* a state glyph **badged onto the corner of the cover**,
the way a store overlays a discount flash. It buys back one caption row per tile — which is a
whole extra row of tiles at Wide — and it is rejected on all four of: unknowable contrast (1),
uncountable colour (2), an unreadable badge on a busy image (co-render's *colour* channel
disappearing into the art), and the label, which cannot fit in a corner at all and which
[`../01-design-system.md`](../01-design-system.md) §3.6 forbids dropping. Recorded because it is
the first thing anyone will propose.

---

## 9 · The full state table

Every row is a real screen, including the ones nobody writes down.

| # | State | Trigger | Body composition | Cost | Copy |
|---|---|---|---|---|---|
| **C1** | **No image support at all** | detection returned false, **or** `ZERADO_NO_IMAGES` is set | **The text deck** — `Z-04`'s ledger, unaltered — plus the capability note on body row 1 | 1 game row (12 → 11) | §11.1 |
| **C2** | **No image support, note dismissed** | `covers.note_dismissed` | **`Z-04` §3, byte-identical.** `v` is unbound and absent from the footer | none | — |
| **C3** | **Image support, every cover cached** | the ordinary state of an enriched library | §3.2 | — | — |
| **C4** | **Image support, a cover not fetched yet** | `metadata.cover_ref IS NULL` and never attempted | The tile is a **type tile**; the box says `not fetched`. `r` re-syncs | — | §11.2 |
| **C5** | **A fetch failed** | the enrichment run errored for this game | The tile says **`not fetched`** — the true statement about the game. The *failure* is screen-wide and is carried by the banner set (`Z-04` §8.1), never printed per tile — **D-15-10** | 1 game/tile row if a banner is up | `Z-04` §8.1 |
| **C6** | **The provider has no art for this game** | the fetch succeeded and returned nothing | Type tile, `nothing found`. **Not an error and not retryable** | — | §11.2 |
| **C7** | **A game that will never have a cover** | `source = physical` — a disc, a cartridge, a hand-added row | Type tile, `added by hand`. **No word of apology, no retry, nothing wrong.** A shelf has no store art and never will | — | §11.2 |
| **C8** | **Offline, cover cached** | no route / DNS, `cover_ref` resolves on disk | **The picture renders normally.** The cache is local; nothing is fetched at render time ([`09-erd.md`](../../blueprint/09-erd.md) §1 — *"`cover_ref` is a local cache path, never a remote URL"*) | — | — |
| **C9** | **Offline, cover not cached** | no route / DNS, no local blob | Type tile, `not fetched`, **plus the `OFFLINE` banner** — which is the honest carrier of *why* nothing new will arrive. §9.1 | 0 covers (**D-15-5**) | `Z-04` B1 |
| **C10** | **The OS deleted the cover cache** | XDG cache eviction ([`09-erd.md`](../../blueprint/09-erd.md) §3) | Identical to C4 — `not fetched`, and `r` refetches. **This is the designed behaviour, not a fault**: the cache is disposable by choice | — | §11.2 |
| **C11** | **Loading** | reading SQLite on first paint | The summary reads `— games`; the grid renders **nothing**. **No scanner, no per-tile placeholder, no skeleton** — a local read is not an indeterminate wait ([`../01-design-system.md`](../01-design-system.md) §9.3) | — | §11.3 |
| **C12** | **First run — nothing synced, nothing connected** | library empty **and** no provider ever connected | **`v` does nothing and is not in that footer.** There is no deck of nothing, and `Z-04` §8.3's empty state owns the body — the same rule that makes `/` inert on an empty library (`Z-07` F1) | — | `Z-04` §10.1 |
| **C13** | **Empty because the provider returned nothing** | a sync succeeded and returned 0 items | As C12 — the mode is unreachable; `Z-04` §8.4 owns the body | — | `Z-04` §10.2 |
| **C14** | **A degrade banner is up** | any of `Z-04` §8.1 B1–B7 | Banner takes body row 1; the pinned block grows down; **`gutterY` goes to 0 and the cover count is unchanged** — **D-15-5** | **0 covers** | `Z-04` §8.1 |
| **C15** | **Filter active** | `/` or `f`, then `v` — **the two modes compose** | `Z-07`'s bar takes body row 1 and carries the ratio; the identity row takes body row 2; the grid keeps 13 rows at Wide and **loses no tile row**. At Narrow the 3-row bar leaves 11 and the grid drops to 1 row | 0 at Wide · 1 tile row at Narrow | `Z-07` |
| **C16** | **Filter active, zero matches** | filter matches 0 | `Z-07`'s zero-result block owns the body. **No empty grid renders** — a grid of nothing is not a smaller grid | all | `Z-07` §11.5 |
| **C17** | **A pending write** | `Z-06` applied, SQLite not yet confirmed | The tile's chip label is suffixed `…` — `ZERADO…` — until the write confirms ([`../01-design-system.md`](../01-design-system.md) §3.4). The picture does not change | — | `Z-04` §10.5 |
| **C18** | **A transition into *zerado* just landed** | `Z-06` wrote `zerado` | The **identity row yields to the result line** until the next keypress — the same slot, the same rule as `Z-04` **D-04-3**. The summary keeps its glyph key beside it | — | `Z-04` §10.5 |
| **C19** | **412 covers, scrolled to the end** | overflow | R-10(a)(b)(c) all hold — §20 | — | — |
| **C20** | **Detail pane focused** (ExtraWide) | `Tab` or `⏎` | Pane border goes heavy; the grid cursor stays visible in chrome, not amber | — | §14 |
| **C21** | **Below 40 columns** | Tiny | **Cover mode is not available.** The deck falls back to the text deck and `v` leaves the footer — §12 | all | §11.4 |
| **C22** | **Below the refusal floor** | `< 24` cols or `< 8` rows | **Frameless**, one sentence, `exit 2` at start-up — `Z-04` §11.3, inherited unchanged | all | `Z-04` §11.3 |
| **C23** | **`NO_COLOR` set** | env | **The pictures still render**; Zerado emits zero SGR — §13 | — | identical text |
| **C24** | **Absent rows present** | `absent_since IS NOT NULL` on ≥ 1 row | Absent games are **not in the deck**, exactly as they are not in the list, and the summary appends `<n> absent` — `Z-04` **D-04-8**, unchanged. **`absent` is not a fifth state and gets no fifth tile treatment** | — | `Z-04` §10.3 |

**Only one strip is ever on body row 1.** Priority, most specific first:
**`Z-07`'s filter bar > a degrade banner > the capability note.**
A stale or incomplete library outranks a missing picture; a filter the player is typing into
outranks both. The note is not lost — it returns on the next render where nothing outranks it,
and it is dismissed by `x` whenever it is visible.

### 9.1 · RENDER 80×24 — C9, offline with a partly-cached deck

`gutterY = 0` (**D-15-5**), so the banner costs **no covers**. Two tiles have no cached blob and
say so; the other six render from the cache exactly as they do online.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   ▌ OFFLINE   Last synced 3 hours ago. Everything here still works.            │
│   247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31  Return of the Obra Dinn                   │
│                                                                                │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░    not fetched      ░░░░░░░░░░░░░░░░░   │
│   ▌ Return of the…     Sekiro: Shadow…    Sid Meier's Ci…    Signalis          │
│     ◉  ZERADO          ⊘  ABANDONED       ○  NOT STARTED     ◐  IN PROGRESS    │
│   ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░                     ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░    not fetched      ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│     Slay the Spire     Stardew Valley     Subnautica         Tunic             │
│     ◉  ZERADO          ○  NOT STARTED     ⊘  ABANDONED       ◐  IN PROGRESS    │
│   COVERS  17–24 of 247                                                         │
│   ↑↓←→ move  ⏎ open  s status  / filter  r sync  v list  ? help  q quit        │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Counts:** banner **65** · identity row **58** · every tile band row **74** · position row
**20** · footer **69**. Grid rows `2 × 6 + 1 × 0 = 12` of `gridH` 12 — exact.

**The age rule is not violated by the pictures.**
[`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §4 requires any
**network-derived value** to render with its age. A cover is not a value — it is a picture of a
box, and it does not go stale in a way a player can be misled by, unlike a price or a playtime
count. **Where the age does belong is the detail pane**, beside the other provenance: a
`COVER  cached, 2 days ago` row on `Z-05`. That is `Z-05`'s to add and is routed in §17, not
asserted here.

---

## 10 · The key map

Every key active in this mode, and nothing that is not.

| Key | Does | Scope | Note |
|---|---|---|---|
| `↑` `↓` `k` `j` | Move the cursor **one tile row** | mode | The grid follows — R-10(b) |
| `←` `→` `h` `l` | Move the cursor **one tile** | mode | Already the global meaning: *"move between columns"* ([`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §3) |
| `g` / `G` | First tile / last tile | global | |
| `Ctrl-D` / `Ctrl-U` | Half a page of tile rows | global | |
| **`v`** | **Leave cover mode → the list** | mode | And on `Z-04`, **enter** it. Bound whenever the deck is reachable; see the ladder below |
| `⏎` | Open the focused game | mode | ≤ 119: push `Z-05` · ≥ 120: move focus into the pane. **Unchanged from the list** |
| `Tab` / `Shift-Tab` | Next / previous region | global | **ExtraWide only** — `R = 2`, because the pane stays (**D-15-7**) |
| `s` | Set this game's status → `Z-06` | mode | Set-status everywhere, always |
| `/` | Filter and search → `Z-07` | mode | The two modes compose — C15 |
| `f` | Jump to `Z-07`'s state chips | mode | |
| `a` | Add a game by hand → `Z-08` | mode | Bound; the **first hint the footer ladder drops** (§10.1) |
| `r` | Re-sync → `Z-03` | mode | Also what refetches a `not fetched` cover |
| `,` | Settings → `Z-09` | global | |
| **`x`** | **Dismiss the capability note** | screen | **Only while the note is showing**, and in the footer for exactly that long (§5.4) |
| `m` | Mute / unmute | global | Only when audio has been enabled |
| `?` | Help → `Z-10` | global | |
| `Esc` | **Nothing.** This is a mode of the root, at rest | global | Cover mode is not a trap: it has no focus of its own to escape, and `v` returns to the list. `Z-04` §13.4 is unchanged |
| `q` · `Ctrl-C` | Quit | global | |

**Reserved and deliberately inert:** `:` and `Ctrl-K` (palette, `Z-17`), `1`–`9`, `n` / `p`.
Pressing one does **nothing** and shows **no error**.

### 10.1 · When `v` is bound, and when it retires — declared, not improvised

| Terminal / state | `v` bound? | In the footer? |
|---|---|---|
| Draws images · library has rows | **yes** | `v covers` (list) / `v list` (deck) |
| No images · note **not** dismissed | **yes** — it raises the note | `v covers`, then `x dismiss` while the note is up |
| No images · note dismissed | **no** | **no** — and `Z-04`'s ratified 73-cell footer returns (§5.5) |
| `< 40` columns (Tiny) | **no** | **no** — §12 |
| Library empty (C12 / C13) | **no** | **no** |

### 10.2 · Footer strings, exact — and what `v` costs `Z-04`

`Z-04` §9.2's ladder is **declared and unchanged**: separator 3 → 2, then drop in the order
`a add` → `r sync` → `m mute` → `s status` → `/ filter` → `⏎ open` → `↑↓ move`, with `? help` and
`q quit` last. Adding one hint runs that ladder; it does not amend it.

| Screen state | Footer | Cells |
|---|---|---|
| **`Z-04` list, covers available** | `↑↓ move  ⏎ open  s status  / filter  r sync  v covers  ? help  q quit` | **69** |
| **`Z-04` list, the note showing** | `↑↓ move  ⏎ open  s status  / filter  r sync  x dismiss  ? help  q quit` | **70** |
| **`Z-04` list, `v` retired** | `↑↓ move   ⏎ open   s status   / filter   a add   r sync   ? help   q quit` | **73** — *the ratified string, unchanged* |
| **Cover mode, Wide** | `↑↓←→ move  ⏎ open  s status  / filter  r sync  v list  ? help  q quit` | **69** |
| **Cover mode, ExtraWide** | `↑↓←→ move   ⏎ detail   tab pane   s status   / filter   a add   r sync   v list   ? help   q quit` | 97 ≤ 112 |
| **Cover mode, Narrow** | `↑↓←→ ⏎ v   ? help   q quit` | 26 |

**The arithmetic, shown, because it costs a hint.** With `v covers` present the Wide line
composes to **84** at 3-space separators and **76** at 2-space — both over the 74-column body —
so the ladder drops **`a add`**, its own first choice, and lands at **69**. `a` stays bound and
is in `Z-10 Help`; only the hint goes, and only while `v` is offered. This is `Z-04`'s ladder
behaving exactly as specified, and it is why the ladder was declared instead of the strings
being frozen. **`Z-04` §3, §7.1 and §9.1 carry the corrected strings** (§17).

---

## 11 · The exact copy — ready to paste

### 11.1 · The capability note

| Slot | String | Cells |
|---|---|---|
| the note | `▌ COVERS   Cover art needs Ghostty or Kitty. Everything else here works.` | **72** |
| the footer affordance | `x dismiss` | 9 |

`'` is **ASCII U+0027** (Narrow), never U+2019 (Ambiguous) — the same rule the rest of the
bundle's contractions follow.

### 11.2 · The three tile words

| Slot | String | Cells |
|---|---|---|
| not asked for yet | `not fetched` | 11 |
| asked; the provider has none | `nothing found` | 13 |
| a hand-added copy | `added by hand` | 13 |

Lowercase, `--z-text-tertiary`, on the image box's last row at tile column 3. **`added by hand`
is the same phrase `Z-06` §10 already uses**, so it is a word the player has met.

### 11.3 · Loading (C11)

```
— games
```

No spinner, no scanner, no skeleton tile, no per-tile placeholder.

### 11.4 · Below 40 columns (C21)

| Slot | String | Cells |
|---|---|---|
| Narrow footer, in the deck | `↑↓←→ ⏎ v   ? help   q quit` | 26 |
| Tiny — the mode is simply absent | *(no copy — `v` is not offered and nothing is said)* | — |

**Nothing is printed at Tiny.** A line explaining that covers need 40 columns would be a notice
about a thing the player did not ask for, on the tier with the least room to spare. §12 gives
the reasoning.

### 11.5 · Position row

| Slot | String | Cells at 247 |
|---|---|---|
| cover mode | `COVERS  17–24 of 247` | 20 |
| list mode | `ROWS  4–15 of 247` | 17 — `Z-04`'s, unchanged |

**The label word changes with the mode because the unit does.** A grid has no rows to count, and
`ROWS  17–24` over a grid of tiles would be a number that describes nothing on screen.

### 11.6 · Copy notes

- **Casing** — `Zerado` in the breadcrumb; `ZERADO` in the chip; `zerado` in the summary's prose;
  `COVERS` uppercase as a label word and as the position-row label, which is the readout role.
- **Say the number** — `17–24 of 247`, `8 covers`, `2 days ago`. Never *some* and never *a few*.
- **No exclamation marks. No emoji. The user is never a "gamer".**
- **No copy refers to a colour, a shape or a position** (WCAG 1.3.3). Never *"the tile on the
  right"*, never *"the cyan one"*.
- **Nothing here claims an unbuilt capability.** There is no *sinopse*, no mood tag, no price and
  no "coming in Phase 2" line on this screen. Covers are built; nothing else in Phase 2 is.
- **Zerado is a games product and the copy says so.** `COVERS` is cover art for games, and
  `added by hand` means a disc or a cartridge. Where a field name is neutral — `TITLE`,
  `SOURCE` — it is neutral because it names a *field*, not because it is hedging toward some
  other kind of thing.
- **Every string above is a catalogue entry, not a literal**
  ([`../../blueprint/16-i18n.md`](../../blueprint/16-i18n.md) §1: *"No user-facing string literal
  appears in code. Every string comes from a catalogue, by key."*). Two of them carry a
  translation hazard worth flagging at the point they are written rather than at the point they
  break: the **three tile words** live in a 15-column box and the **note** in a 74-column body,
  and neither budget grows in another language. The tile words are the tighter constraint —
  `not fetched` is 11 cells in English and a translator has 15. **A catalogue entry that
  overflows its field is a layout defect, not a copy defect**, so both fields are measured at
  render time (§3.6) and both belong in whatever length budget the i18n catalogue carries.

---

## 12 · 40-column behaviour, and the refusal floor

### 12.1 · RENDER 40×24 — Narrow · `leftInset` 2 · body `36 × 16`

`2 across × 2 down = 4 covers`, and the grid fills the body **exactly**: `2 × 17 + 1 × 2 = 36`.
The identity row takes body row 2, because `36 < 61` (**D-15-4**).

```
┌────────────────────────────────────────┐
│                                        │
│  Zerado ✦ Library                      │
│                                        │
│  LIBRARY                               │
│                                        │
│                                        │
│  247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31     │
│  Return of the Obra Dinn               │
│  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  │
│  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  │
│  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  │
│  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  │
│  ▌ Return of the…     Sekiro: Shadow…  │
│    ◉  ZERADO          ⊘  ABANDONED     │
│                                        │
│  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  │
│  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  │
│  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  │
│  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  │
│    Signalis           Slay the Spire   │
│    ◐  IN PROGRESS     ◉  ZERADO        │
│  COVERS  17–20 of 247                  │
│  ↑↓←→ ⏎ v   ? help   q quit            │
│                                        │
└────────────────────────────────────────┘
```

**Counts:** body 36 · summary form 3 **33** ≤ 36 · identity row 23 · every tile band row **36**
· position row 22 · footer 26.

**Sheds at Narrow:** the respiro (to the identity row, **D-15-4**) · the footer's words (to the
glyph-and-key form `Z-04` §9.1 already uses at this tier).
**Keeps, because [`03-responsive.md`](../../blueprint/03-responsive.md) §4 forbids hiding them:**
the state — glyph **and** label, at the full 14-column chip · the game's title, untruncated, on
the identity row · the pinned summary · the focus marker · the degrade banner · the footer.

**Standard `60 × 24`** is the same composition: body 54, 2 across, 4 covers, identity row on its
own row (`54 < 61`). At 79 columns the same tier takes 3 across and 6 covers, and the identity
row moves up to share the summary row (`73 ≥ 61`). Stated as arithmetic rather than as a
judgement.

### 12.2 · Tiny `< 40` — cover mode is not available, and nothing is said about it

At Tiny the body is `30 × 21`. One tile fits across (`30 − 17 = 13` spare) and three down, so the
deck would show **3 covers of 247** — against the text deck's **9 games** at the same size
([`Z-04-library.md`](./Z-04-library.md) §11.2). Eighty-three screens to walk a library.

> **D-15-8 · Below 40 columns `v` is unbound and absent from the footer; the library renders as
> the text deck.** A one-column grid is a list with pictures in it, and at 30 columns the picture
> would take 17 of the body's 30 while the title field fell from 26 to 15 — starving the identity
> column that R-10(a) protects, on the tier the Spacing Canon already treats as starvation
> territory (`OuterMarginX` and `OuterMarginY` both shed to 0).
>
> A terminal **resized down** into Tiny while the deck is open falls back to the list on the
> `WindowSizeMsg`, keeping the same focused game; resized back up, `v` returns to the footer and
> the deck can be re-entered. **The mode is not lost, only the room for it.** This is the same
> shape as `Z-06`'s overlay becoming a route at Tiny ([`02-composition.md`](../../blueprint/02-composition.md)
> §2.4) — behaviourally continuous, compositionally different.
>
> **And nothing is printed.** A line explaining why covers are absent would be a notice about a
> thing the player did not ask for, on the narrowest terminal, spending a row it cannot spare.
> Silence here is the same judgement §5.1 makes in the other direction: **Zerado never volunteers
> a fact about images.** *Flagged for the founder in §18, because it is the one place this spec
> chooses to say nothing.*

### 12.3 · The refusal floor — 24 columns or 8 rows

Inherited from `Z-04` §11.3 unchanged, because the mode does not exist below 40 anyway.
**Frameless**, one line, exit status **2** at start-up; a running session resized below the floor
shows the same sentence and **keeps running**:

```
Zerado needs at least 24 columns and 8 rows. This terminal is 20 x 6.
```

---

## 13 · `NO_COLOR` — rendered, not asserted

`NO_COLOR` set → **zero SGR sequences** (brand §5.4). Bold is an SGR sequence and is therefore
also gone. The characters are unchanged:

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Library                                                             │
│                                                                                │
│   LIBRARY                                                                      │
│                                                                                │
│                                                                                │
│   247 games  ○ 198  ◐ 12  ◉ 6  ⊘ 31  Return of the Obra Dinn                   │
│                                                                                │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ▌ Return of the…     Sekiro: Shadow…    Sid Meier's Ci…    Signalis          │
│     ◉  ZERADO          ⊘  ABANDONED       ○  NOT STARTED     ◐  IN PROGRESS    │
│                                                                                │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│   ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░   │
│     Slay the Spire     Stardew Valley     Subnautica         Tunic             │
│     ◉  ZERADO          ○  NOT STARTED     ⊘  ABANDONED       ◐  IN PROGRESS    │
│   COVERS  17–24 of 247                                                         │
│   ↑↓←→ move  ⏎ open  s status  / filter  r sync  v list  ? help  q quit        │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

| Information | Channel that survives |
|---|---|
| Each game's state | the **glyph** and the **word**, at full width, on every tile — the co-render rule doing its job |
| Which tile has focus | the `▌` marker — **position**. The one channel that never needed colour |
| Which game the cursor is on | the identity row **says the title**, in text |
| That a tile has no picture | the **word** in the image box — `not fetched`, `nothing found`, `added by hand` |
| That this is the deck and not the list | the position row's label word, `COVERS` |
| Where in 247 the player is | `COVERS  17–24 of 247` |
| That the terminal draws no images | the note, in words, once (§5.3) |

**No information is lost.**

### 13.1 · `NO_COLOR` does not remove the pictures — *design decision*

> **`NO_COLOR` removes Zerado's colour. It does not remove the artwork.**

`NO_COLOR`'s contract is *"Zerado emits no SGR sequences at all"*
([`../03-designer-manual.md`](../03-designer-manual.md) §5.4). An image placed by the graphics
protocol is not an SGR sequence, and greyscaling or suppressing it would be Zerado **adding** a
behaviour the variable does not ask for.

**The reason it is coherent rather than a loophole is §8.1.** The artwork carries no information,
so removing it would lose nothing and keeping it costs nothing: not one contrast ratio in §8 is
measured against a picture, not one cyan mark can hide in one, and the co-render rule is
satisfied entirely on the caption rows. A player who wants a terminal with no pictures in it has
**`ZERADO_NO_IMAGES=1`**, which is the variable for that question and answers it completely.

*The alternative considered and rejected:* treat `NO_COLOR` as *"no colour of any kind"* and
suppress the images. Rejected because it conflates two different requests — *do not encode
meaning in colour* and *do not show me photographs* — and because it would make `NO_COLOR`
delete content on exactly one screen in the product, which is a surprise a general-purpose
environment variable has no business springing.

---

## 14 · The focus model, and how `Esc` behaves

### 14.1 · Regions

| Tier | `R` | Regions |
|---|---|---|
| Narrow · Standard · **Wide** | **1** | the grid. `Tab` is unbound and unlisted |
| ExtraWide | **2** | the grid · the detail pane (**D-15-7**) |

Cover mode **does not change the region count** at any tier, which is what keeps it a body
renderer rather than a composition.

### 14.2 · How focus is shown — three channels, any two sufficient

| Channel | Focused tile | Every other tile |
|---|---|---|
| **Position** | `▌` U+258C in the caption's 2-column focus field — ASCII `>` | two spaces |
| **Weight** | the title is **bold** | normal |
| **Colour** | `--z-primary` amber on the marker | none |

Identical to the ledger row's ([`../01-design-system.md`](../01-design-system.md) §1.7), by
reuse, not by re-derivation. **The cursor is amber, not cyan**, so the deck spends no
focus-ring cyan at all. Under `NO_COLOR` the marker and the bold survive; at 16 colours the
marker survives.

**Focus is never nowhere.** The traversal is over the same row set the list uses, in the same
order (title A→Z, fixed — [`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2),
so **`v` preserves the focused game across the mode change in both directions**, and the deck
scrolls to bring it into view. A player who was on *Return of the Obra Dinn* in the list is on
*Return of the Obra Dinn* in the deck. That is the property that makes `v` feel like a mode and
not a navigation.

**Cursor-following scroll — R-10(b).** The focused tile is always visible; a background refetch
that rebuilds the row set preserves cursor and offset **by game identity**, not by index, exactly
as `Z-04` §12(b) requires. A cover arriving mid-scroll swaps a type tile for a picture **in
place** and moves nothing.

### 14.3 · `Esc`, exhaustively, in this mode

| Situation | `Esc` |
|---|---|
| Cover mode, nothing else open | **Nothing.** This is the root at rest, and `Z-04` §13.4 is unchanged. `v` returns to the list |
| The capability note showing | **Nothing** — `x` dismisses it (§5.4). The note takes no focus, so it is not a trap and 2.1.2 is not engaged |
| Filter mode also active (C15) | `Z-07`'s two-step `Esc`, unchanged — blur the editor, then clear the filter. **Cover mode is not what `Esc` unwinds** |
| The detail pane focused (ExtraWide) | Focus returns to the grid — `Z-04` §13.4 |
| An overlay open (`Z-06`) | The overlay closes; the deck is behind it, unchanged |

**There is no way to get stuck in cover mode**, because there is nothing to be stuck in: it is
one keystroke to the list, the same keystroke that opened it, and every other key on the screen
means what it means in the list.

---

## 15 · Colour budget declaration

**Counted by [`../02-colour-budget.md`](../02-colour-budget.md) §3.1, from the ANSI stream,
classified by payload.**

| Class | Count | Where |
|---|---|---|
| **CHROME cyan** | **0** | The deck urges nothing. *"A screen may spend zero. Cyan is earned; a screen with nothing to urge does not need to urge anything"* (§2.2). The mode is a way of looking, not a call to action |
| **STATE cyan** | unbounded, **not counted** | The `◉` glyph and the `ZERADO` label on every finished game's tile. §2.1 — a player with more cyan has finished more games |
| **Focus-ring cyan** | **0** | The tile cursor is `--z-primary` **amber** (§14.2), so the deck emits no cyan outside the state chips |
| **Amber** | title (allow-list **1**) · the `▌` cursor (**5**) · every `IN PROGRESS` chip (**3**) · `▮ AUDIO` when enabled (**9**) | every mark on the §4.1 allow-list |
| **Red** | **0** | Nothing on this screen is motion or alarm. The capability note is **chrome**, and the `OFFLINE` banner is chrome, for the same published reason (§5.2) |
| **The artwork** | **not counted, and it cannot be** | An image emits no SGR run, so it cannot appear in a cyan or amber scan (§8.1 point 2) |

**Amber ceiling.** At `80 × 24 = 1920` cells: title `LIBRARY` **7** + cursor **1** + the two
`IN PROGRESS` chips in §3.2 at 14 cells each = **36 cells → 1.9 %**. The worst case a deck can
reach is all eight tiles `IN PROGRESS`: `8 × 14 + 7 + 1 = 120 → 6.3 %`. **Below the 10 % ceiling
at maximum**, which is a bound rather than an observation.

**One added line for the reviewer, and it is the only change this screen asks of the budget
document.** §3.2's human method reads *"count what is left."* On this screen:

> **Ignore the pictures entirely. Count Zerado's own marks.** An image is not a colour Zerado
> spent, and no Zerado mark may sit on one — so if a coloured mark appears **inside** an image
> rectangle, that is a §8.1 violation and an automatic fail, whatever its hue.

**The temptation this screen resists.** A cover deck invites a tinted overlay on the focused
tile, a coloured frame to "tie" the selection to the pane, and a cyan glow on finished games.
All three are failure-gallery items 2, 3 and 7, and all three are refused by §8.1 before the
budget is even consulted.

---

## 16 · Reuse verdict, per element

| Element | Verdict | Note |
|---|---|---|
| The grid | **Build fresh** — `lipgloss.JoinHorizontal` over tiles, `JoinVertical` over bands | `bubbles/list` is a **list**: one item per row, its own pagination and its own filter UI. A 2-D cursor over a fixed-pitch grid is not what it models, and bending it would cost more than the ~60 lines this is |
| The scroll region | **`bubbles/viewport`** — direct fit | The same primitive the ledger and `Z-09` use, so cursor-following scroll is written once and behaves identically everywhere |
| The tile | **Build fresh** — a `lipgloss` block, plus the width-aware pad and truncate of §1.2 | It is the ratified chip plus a title plus an image rectangle |
| The state chip | **Reused unchanged** from [`../01-design-system.md`](../01-design-system.md) §3 | Not re-specified here, not restyled, not compressed. This is the whole reason co-render survives the mode |
| Header band, footer, position row, audio annunciator | **Reused unchanged** from `Z-04` | A mode swaps the body renderer |
| The capability note | **Build fresh** — the informational banner's anatomy, with a dismiss | Not the degrade banner itself: it is raised by a **capability**, not a failure, and it is dismissible, which no banner is. It shares the `▌` + label-word + sentence shape so there is one thing to learn |
| Image placement | **The spine's [`17-images.md`](../../blueprint/17-images.md) §4 `Images` seam** | `Cover(id, cells Rect) (Placement, bool)` — the screen hands it the **cell rectangle** and takes a placement or a `false`. Protocol, encoding, cache and lifetime are not a design decision and are deliberately not specified here. **`false` is the type tile of §3.3**, which is why the type tile is the ordinary render and not an error path |
| Capability detection | **The spine's [`17-images.md`](../../blueprint/17-images.md) §2 · §4** | `Capability()` — resolved once at startup, cached for the session, **failing closed**. §5.1 states the design requirements it must satisfy; the mechanism is not this document's. `NullImages` is the no-support implementation and is what the tests run, so the text path is the exercised one |
| Cover cache | **Reused** — the XDG cache directory ([`09-erd.md`](../../blueprint/09-erd.md) §3 · [`17-images.md`](../../blueprint/17-images.md) §5) | Bounded, evicting least-recently-shown, **disposable and refetchable by design** — *"deleting the whole cache must cost nothing but bandwidth."* C10 is that choice behaving correctly, not a fault |
| `harmonica` | **Not used.** No motion on this screen | The scanner is for indeterminate waits and never ambient; a cover fading in would be decoration ([`../01-design-system.md`](../01-design-system.md) §9.3) |

---

## 17 · Upstream findings — contradictions and stale rows

Routed, not fixed here. Each names the document, the line, and the owner.

| # | Finding | Where | Owner |
|---|---|---|---|
| 1 | **`Z-15` is not in the Phase 1 composition table.** [`02-composition.md`](../../blueprint/02-composition.md) §2 lists eleven screens and §3 still carries `Z-15` under *"later phases"*, while [`01-screen-inventory.md`](../../blueprint/01-screen-inventory.md) §2 already says **twelve** Phase 1 screens. The §3 row's verdict — *"a mode swaps the body renderer, not the frame"* — is correct and this spec builds on it; it is only in the wrong table | `02-composition.md` §2 · §3 | `fft-tui-architect` |
| 2 | **`Z-15` has no row in the responsive degrade table.** [`03-responsive.md`](../../blueprint/03-responsive.md) §3 covers `Z-01`–`Z-11`. The values this spec composes to are: ExtraWide *deck ∥ detail, 3 × 4* · Wide *4 × 2* · Standard *2 × 2, or 3 × 2 at 79* · Narrow *2 × 2* · Tiny **not available** (**D-15-8**) | `03-responsive.md` §3 | `fft-tui-architect` |
| 3 | **`Z-04` carries `v` now, and its footer strings move.** §10.2 has the exact replacements for `Z-04` §3's mockup footer, §7.1's audio footer and §9.1's table, plus a `v` row in its §9 key map. The **ladder is unchanged** — this is `Z-04` §9.2 running, and `a add` is its own declared first drop | `Z-04-library.md` §3 · §7.1 · §9 · §9.1 | **applied in this pass** |
| 4 | **`Z-10 Help` does not list `v` or `x`.** Its key map is **generated** from the registry (D-10-1), so binding `v` fixes it automatically — but the spec's drawn goldens are hand-counted and go stale: block 1 gains one line, so the document is **34 lines, not 33**, and every `ROWS  1–15 of 33` readout becomes `of 34`. `x dismiss` carries an availability predicate and appears only while the note is up, so it adds no line to the default render. `Z-10` needs a rev B | `Z-10-help-and-key-map.md` §4 · §5 · §6 · §7 · §9 · §13.1 · §18 | `fft-tui-designer` — **next pass** |
| 5 | **`Z-05` has no `COVER` row**, so a cached cover has nowhere to state its age. Recommended: `COVER  cached, 2 days ago` in the provenance block, with `not fetched` · `nothing found` · `added by hand` as its other values — the same three words §11.2 fixes here. `Z-05` §10 currently states *"no cover-art placeholder"*, which was correct while covers were Phase 2 | `Z-05-game-detail.md` §10 · D-05-6 | `fft-tui-designer` — **next pass** |
| 6 | **`Z-09 Settings § DISPLAY` should carry an `Images` row.** `Glyphs`, `Motion` and `Colour` are already read-only reports of what Zerado resolved about the terminal; image support is exactly that kind of fact, and it is the durable home for the note's content once the note is dismissed. Values: `Kitty graphics protocol` · `Not supported by this terminal` · `Off, ZERADO_NO_IMAGES is set`. **Not done in this pass** — it takes the group from 22 rows to 23 and moves every position readout in `Z-09`, which is a founder-visible change to a GOLDEN spec rather than a sweep | `Z-09-settings.md` §10.4 | founder, then `fft-tui-designer` |
| 7 | **`ZERADO_NO_IMAGES` is adopted but not yet catalogued.** [`17-images.md`](../../blueprint/17-images.md) §4 binds it — *"the implementation used when detection fails, in tests, and under `ZERADO_NO_IMAGES`"* — so §5.1's proposal is **closed, and this spec uses it as adopted**. What is missing is the one place a player or a builder would go looking: there is no single environment-variable table in the spine. `NO_COLOR` lives in the brand manual, `ZERADO_ASCII` in [`03-responsive.md`](../../blueprint/03-responsive.md) §5c, `ZERADO_NO_AUDIO` in [`12-audio.md`](../../blueprint/12-audio.md) §10, `ZERADO_DB` in `Z-09` §10.5 and `ZERADO_NO_IMAGES` in `17-images.md` §4. **Five variables, five homes, no index** | the spine, no single owner today | `fft-tui-architect` |
| 8 | **The design brief's accessibility table says Zerado has no images.** [`../00-design-brief.md`](../00-design-brief.md) §131 reads *"**1.4.5 / 1.4.9** Images of Text — No images. Cover art arrives in Phase 2 and is art, not text."* The **conclusion** survives — a cover is artwork, not an image of text, so 1.4.5 is still satisfied — but the premise is now false, and **SC 1.1.1 Non-text Content becomes live for the first time in this product**. Zerado's answer is already built and should be recorded as the brief's answer: **the caption is the text alternative, it is always present, it is never optional, and it is what remains when the picture is not there** (§8.1, §13). The line needs rewriting rather than deleting | `00-design-brief.md` §131 | `fft-design-architect` |
| 9 | **`ADR-0001` still calls the cover deck a Phase 2 mode.** Its *"not assumed anywhere"* list reads *"the text deck is the default and the cover deck is a **Phase 2 mode** that may never be universally available"* — superseded by the 2026-08-25 direction, and now contradicted by [`17-images.md`](../../blueprint/17-images.md), which the ADR does not cite. **Re-checked at head: the matching line in [`00-index.md`](../../blueprint/00-index.md) has already been fixed**, so the two documents now disagree with each other as well. The ADR's *substance* survives intact and is this screen's design — the text deck really is the default, and image support really is not universal | `ADR-0001`, the *not assumed* list | `fft-tui-architect` |
| 10 | **Cover art is now in both offline-contract tables.** [`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2 lists `Cover art (Z-15) · DEGRADES` under **Phase 1** — correct, and the count sentence beneath it now reads *twelve screens* — and still lists `Cover art, sinopse · DEGRADES` under **Phase 2 and later**. The two rows do not conflict, because *sinopse* really is Phase 2, but the duplicated subject invites a reader to conclude covers are Phase 2. Splitting the row so Phase 2 carries *sinopse* alone would close it | `07-offline-contract.md` §2 | `fft-tui-architect` |

---

## 18 · Open for the founder

1. **ExtraWide: the pane, or twice the covers.** §4 keeps `Z-04`'s ratified `66 ∥ 2 ∥ 44` split
   and shows **12 covers beside the full detail pane**. Dropping the pane gives **24 covers** —
   double — and moves the untruncated title onto the pinned row. It is a composition change
   (`R` 2 → 1, `Tab` unbound, `⏎` changes meaning), so this spec did not take it.
   **Which do you want at 120 columns?**
2. **The note's wording: two programs, or the protocol.** It says *"Cover art needs Ghostty or
   Kitty. Everything else here works."* — the first sentence is [`17-images.md`](../../blueprint/17-images.md)
   §3's own worked example, verbatim; the second is this spec's, echoing the `OFFLINE` banner's
   ratified reassurance. It does **not** name the **Kitty graphics protocol**, which is the term
   a player would need to search. §5.3 puts the program names in the note and the protocol name
   in Help and Settings. **Confirm** — this is the sentence a player on Terminal.app or Alacritty
   will read, and it is the one place in this bundle where tone outranks precision.
3. **At Tiny, covers are absent and nothing is said** (**D-15-8**). A one-column grid shows 3 of
   247 against the text deck's 9, and a line explaining the absence would cost a row on the
   narrowest terminal. Confirm silence, or ask for the line.
4. **Dismissal is permanent** (§5.4). Once `x` is pressed, the note never returns and `v` retires
   on that kind of terminal. There is no *"show it again"* control, because a control for it
   would be a dial nobody would ever find. If you want one, it is the `Images` row of §17
   finding 6.
5. **Five environment variables now, and no single place that lists them.** `NO_COLOR`,
   `ZERADO_ASCII`, `ZERADO_NO_AUDIO`, `ZERADO_REDUCED_MOTION`, `ZERADO_DB` and — since
   [`17-images.md`](../../blueprint/17-images.md) §4 — `ZERADO_NO_IMAGES`. Each is documented
   where it was invented and nowhere together (§17 finding 7). Not this screen's to fix, but
   this screen is the one that took the count past five.
6. **Covers are the store's header capsule, not portrait box art** (**D-15-2**). That is forced
   by the arithmetic, not chosen: portrait art plus the un-droppable state label does not fit at
   80 columns at all. Worth knowing, because the pictures will not look like the box on a shelf.

---

## 19 · Design decisions made in this spec

| # | Decision | Reason |
|---|---|---|
| **D-15-1** | The tile is **17 × 6** at every tier; only the grid's column and row counts change (§3.1) | A box that breathes differently at different widths looks resized, not designed — D-06-1's reasoning, applied. One geometry is also one golden per tier |
| **D-15-2** | The image box is **17 × 4 ≈ 2.125 : 1** — the store header capsule, not portrait box art (§3.1) | The **co-render rule sizes the tile**: the 14-column chip forces a ≥ 16-column tile, and portrait art at that width needs 15 rows for a single band. It would not render at 80 columns at all, which verdict 2 forbids |
| **D-15-3** | **Zerado never draws a mark on a cover** (§8.1) | It is the only way contrast stays measurable, the colour budget stays countable, and the no-image render is a composition rather than a fallback |
| **D-15-4** | The caption is the tile's **label**; the pinned row is the **identity column** (§3.5) | A 15-column title is the exact failure `02-composition.md` §2.1 names for R-10(a). The summary drops to its declared form 3 to make room, and D-04-1's guard holds because every tile spells its state |
| **D-15-5** | The grid's **row gutter is the flexible row**, 1 or 0, and nothing else flexes (§3.4) | Without it a degrade banner costs half the deck and leaves six blank rows. The shape of the rule is D-04-5's, in the other axis |
| **D-15-6** | **The note answers; it never announces** (§5.3) | It is raised only by `v`, said once, and dismissed for good. That is the whole of *"never a recurring scold"*, and it is a behaviour rather than a promise |
| **D-15-7** | At ExtraWide the **detail pane stays**; the deck takes 3 columns (§4) | Dropping it changes `R`, `Tab` and `⏎` — composition, not a body renderer. The alternative is costed and put to the founder |
| **D-15-8** | **Below 40 columns cover mode is unavailable, and nothing is printed about it** (§12.2) | 3 covers of 247 against the text deck's 9, with the title field falling from 26 to 15 on the tier the canon already treats as starvation territory |
| **D-15-9** | **`NO_COLOR` does not remove the pictures** (§13.1) | `NO_COLOR` removes SGR. An image is not SGR, it carries no meaning, and suppressing it would make one environment variable delete content on one screen. `ZERADO_NO_IMAGES` is the variable for that question |
| **D-15-10** | **Three tile words**, and `fetch failed` is not one of them (§3.3) | A failure is a property of the network, not of the game, and it is screen-wide — so it is the banner's. Printing it 247 times would blame the game for the connection |
| **D-15-11** | The position row says **`COVERS`**, not `ROWS` (§11.5) | A grid has no rows to count. The label word changes with the mode because the unit does |

---

## 20 · Screen-specific acceptance criteria

Beyond [`../00-design-brief.md`](../00-design-brief.md) §10 and
[`../02-colour-budget.md`](../02-colour-budget.md) §10, which are the bar.

1. **Every tile carries colour AND glyph AND label.** Render at every tier and grep the raw
   output: the count of `NOT STARTED` + `IN PROGRESS` + `ZERADO` + `ABANDONED` equals the number
   of tiles on screen. **No tile anywhere renders a glyph without its word.**
2. **No SGR run overlaps an image cell rectangle** (§8.1). Capture the raw stream, resolve every
   image placement to its cell rectangle, and assert the intersection is empty. This is the
   machine form of *Zerado never draws on a cover*, and it is a hard fail.
3. **The chrome-cyan count is 0**, and every remaining cyan run's payload is `◉`, `ZERADO` or
   `[*]` (§15).
4. **Amber is ≤ 10 % at the worst case**, not just the drawn one: force all eight Wide tiles to
   `IN PROGRESS` and assert **120 cells of 1920 → 6.3 %**.
5. **The mockups add up.** Every body row in §3.2, §3.3, §4, §5.2, §5.5, §9.1, §12.1 and §13
   measures its stated width with the **width-aware** function, with **Ambiguous = 1** — and the
   same assertion passes with **Ambiguous = 2** for every line that is not inside a padded field
   ([`../01-design-system.md`](../01-design-system.md) §1.2 rules 1 and 2).
6. **A banner costs zero covers** (**D-15-5**). Force each of `Z-04` B1–B7 in cover mode at
   80 × 24 and assert the visible cover count is **8** in every case, with `gutterY = 0`.
7. **The text deck is byte-identical to `Z-04`.** Run with image support forced off and the note
   dismissed; diff the render against `Z-04` §3's golden. **Zero differences**, including the
   73-cell footer with `a add` restored.
8. **The note appears only after `v`, and only once.** Start with no image support: assert the
   first paint has no note and no `x` in the footer; press `v`, assert the note and the 70-cell
   footer; press `x`, assert the note is gone, `covers.note_dismissed` is written, and `v` has
   left the footer; restart the process and assert the note does not return.
9. **The note never outranks a banner or the filter bar** (§9). Force `OFFLINE` + no image
   support + `v` and assert body row 1 is the banner; clear the banner and assert the note
   returns.
10. **`v` preserves the focused game in both directions** (§14.2). Scroll to game 380 of 412 in
    the list, press `v`, assert the same game is focused and visible in the deck; press `v` again
    and assert the list is back on it, at the same offset. **By identity, not by index.**
11. **The three tile words are distinguishable without colour.** Render with `NO_COLOR=1` and
    assert `not fetched`, `nothing found` and `added by hand` are all present and legible, and
    that a hand-added game shows `added by hand` and **no retry affordance**.
12. **A cover arriving mid-scroll moves nothing.** With the deck open at an overflowing count,
    complete an enrichment for a visible game and assert the type tile becomes a picture **in
    place** — same cursor, same offset, same tile origin (R-10(b)).
13. **The ledger triad survives the mode** at 412 covers: the identity column is populated and
    human-readable (the identity row, untruncated) · the focused tile is always visible ·
    `COVERS  n–m of 412` and the summary are **pinned outside the grid** and cannot be scrolled
    off (R-10(a)(b)(c)).
14. **At Tiny, `v` is unbound, absent from the footer, and silent.** Resize from 80 to 32 with
    the deck open and assert the list renders, the same game is focused, and **no line is printed
    about covers**.
15. **`ZERADO_NO_IMAGES=1` draws no image at all** — not fewer, none — and the screen is
    identical to the no-support render.
16. **Artifacts:** the eight of [`03-responsive.md`](../../blueprint/03-responsive.md) §7,
    **plus** `80 × 24` for: the mixed deck, the note, the note dismissed, offline with a partly
    cached deck, a filter active in cover mode, and `NO_COLOR`. Plus one at **412 covers**,
    scrolled to the end, for criterion 13 — which a frozen golden cannot prove.
