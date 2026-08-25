---
title: Zerado — Z-05 Game detail
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-05
rev: B
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-05 · Game detail

> **One view, two hosts.** A route below 120 columns; a pane inside `Z-04` at 120 and above.
> *"This is the single most important composition decision in the bundle, because it is what
> makes the same detail spec buildable once and mounted twice — and it is why `Z-05`'s spec must
> be written host-agnostic."* (`02-composition.md` §2.1.)

**Canon that governs this screen, as pointers:** `00-design-brief.md` §10 ·
`01-design-system.md` §6 (detail pane) · §3 (chip) · §1.5 (readout role) ·
`02-colour-budget.md` §10 · `03-designer-manual.md` §3 · `02-composition.md` §2.1 · §2.3 ·
`03-responsive.md` §3 · `04-navigation-and-focus.md` §5 · `05-state-machine.md` §1, §2.1, §5,
§6 · `06-data-seams.md` §2 · `07-offline-contract.md` §2, §4 ·
`01-screen-inventory.md` §5 — *"`Z-05` must not render an empty field as if it were an empty
value. Not-fetched-yet and known-to-be-empty are different facts and the player can act on only
one of them."*

---

## 0 · The host contract — read this before anything else

**This view never assumes its own width, never draws a header band, never owns a global key and
never decides what `Esc` means.** It is handed a rectangle and a focus flag, and it renders a
single column into it. Everything host-specific is in the table below and nowhere else in this
document.

| | **Host A — route** (Tiny · Narrow · Standard · Wide) | **Host B — pane** (ExtraWide `120+`) |
|---|---|---|
| Who mounts it | the route stack, as its own frame | `Z-04`'s body, right column |
| Content rectangle | `BodyRect` — `74 × 16` at 80 × 24 | **38 × 30** inside a 44 × 32 bordered box |
| Header band | the frame's — `Zerado ✦ Game detail` / `GAME DETAIL` | **none.** `Z-04`'s band, reading `LIBRARY`, is above it |
| A border | **none** | `┌─┐` unfocused · **`┏━┓`** focused, titled `DETAIL` |
| Footer | the frame's reserved row, listing this screen's keys | `Z-04`'s footer |
| `⏎` from the ledger | pushes this route | **moves focus into this pane** |
| **`Esc`** | **pop the route**, back to `Z-04` | **return focus to the ledger**; the pane stays rendered |
| Audio annunciator | on the frame's title row | on `Z-04`'s title row |
| Degrade banner | this screen's body row 1 | `Z-04`'s body row 1, above the whole split |

**What the view exposes to its host, and nothing more:** its natural height, whether it has
scrollable overflow, and *"I have nothing of my own to dismiss."* The host maps `Esc`. If you
find yourself writing `if isRoute` inside the view, the contract has been broken.

---

## 1 · Identity

| | |
|---|---|
| **ID** | `Z-05` |
| **Name** | Game detail |
| **Phase** | 1 |
| **Kind** | **Route ≤ 119 cols · Pane ≥ 120 cols** |
| **Route in** | `⏎` on a focused row of `Z-04` (≤ 119). At ≥ 120 nothing is pushed — focus moves |
| **Route out** | `Esc` → `Z-04` · `s` → `Z-06` · `r` → `Z-03` · `,` → `Z-09` · `?` → `Z-10` |
| **Shape** | **Detail view** — shape 2 of the five (`02-composition.md` §3) |
| **Offline class** | **WORKS** (`07-offline-contract.md` §2) — every locally-known field |
| **Displayed screen name** | breadcrumb `Zerado ✦ Game detail` · title `GAME DETAIL` · pane title `DETAIL` |

> **Why the header band names the *screen* and not the game.** Setting the game's own title in
> the display role — amber uppercase — would be beautiful in host A and **impossible** in host B,
> where `Z-04`'s band reads `LIBRARY`. A view that is host-agnostic cannot own the band. The
> game's title is therefore **body row 1 in both hosts**, in `--z-text`, sentence case, wrapped
> and never truncated. Recorded so the prettier option is not re-proposed at build time.

---

## 2 · Purpose

**One sentence:** everything Zerado actually knows about one game, with everything it does not
know marked as *which kind* of not-knowing it is.

---

## 3 · Mockup — host A, 80 × 24, the design floor

`tier = Wide` · `leftInset = 3` · **body = 74 × 16** · every visible row begins at **column 4**.

**This is the normal Phase 1 state, not a degraded one.** Phase 1 has no *sinopse* and no mood
tags, so this composition — four facts about the game, two about where its status came from, two
about the age of what is on screen — is the screen, designed as such.

> **Rev B: Phase 1 has cover art, and this screen gained the `COVER` row.**
> [`Z-15-cover-deck.md`](./Z-15-cover-deck.md) §9.1 routed it here with the reason:
> [`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §4 requires every
> **network-derived value** to render with its age, and **a cover's age is not `LAST SYNCED`'s**
> — the row synced three hours ago and the picture was fetched two days ago. One age row could
> not tell the truth about both. **`Z-15` draws no age on a tile and never will** (`Z-15` §8.1:
> *Zerado never draws a mark on a cover*), so the detail view is the only surface where the
> cover's provenance can be stated at all. §9.7 maps every one of `Z-15`'s seven picture states
> onto what this row says.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Game detail                                                         │
│                                                                                │
│   GAME DETAIL                                                                  │
│                                                                                │
│                                                                                │
│   Return of the Obra Dinn                                                      │
│                                                                                │
│   ◉  ZERADO                                                                    │
│                                                                                │
│   PLAYTIME     9h                                                              │
│   LAST PLAYED  2 Aug 2026                                                      │
│   ADDED        12 Mar 2026                                                     │
│   SOURCE       Steam                                                           │
│                                                                                │
│   SET BY       you, 12 Aug 2026                                                │
│   STEAM SAYS   IN PROGRESS                                                     │
│                                                                                │
│                                                                                │
│                                                                                │
│   COVER        cached, 2 days ago                                              │
│   LAST SYNCED  3 hours ago                                                     │
│   ↑↓ scroll   s status   r sync   esc back   ? help   q quit                   │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Character counts:** terminal **80 × 24** · body **74 × 16** · readout label field **11** +
gutter **2**, so every value begins at body column **14** · footer **58** · title row **23**
(`Return of the Obra Dinn`) · state chip **14** (`◉` in a 2-cell field + 1 + an 11-cell label).

### 3.1 · Row map — 80 × 24

| Terminal row | Body row | Content |
|---|---|---|
| 1 | — | blank · `OuterMarginY` = 1 |
| 2–4 | — | header band: breadcrumb · `InterElementGap` · title |
| 5–6 | — | blank · `InnerPaddingY` = 1, then `InterElementGap` = 1 |
| 7 | 1 | **the game's title** — wrapped, never truncated |
| 8 | 2 | `InterElementGap` |
| 9 | 3 | **the state chip** |
| 10 | 4 | `InterElementGap` |
| 11–14 | 5–8 | **block 1 — the game** |
| 15 | 9 | `InterElementGap` |
| 16–17 | 10–11 | **block 2 — where the status came from** |
| 18–20 | 12–14 | blank |
| 21–22 | 15–16 | **block 3 — the age of what is on screen**, bottom-anchored: `COVER`, then `LAST SYNCED` |
| 23 | — | footer, the reserved frame row |
| 24 | — | blank · `OuterMarginY` = 1 |

### 3.2 · The three blocks, and why there are exactly three

| Block | Keys | The question it answers |
|---|---|---|
| **1 · the game** | `PLAYTIME` · `LAST PLAYED` · `ADDED` · `SOURCE` | *What is this thing?* |
| **2 · the status** | `SET BY` · `<PROVIDER> SAYS` | *Why does it say `ZERADO`, and what would happen if I cleared that?* |
| **3 · the age** | `COVER` · `LAST SYNCED` | *How old is what I am reading — and what I am looking at?* — `07-offline-contract.md` §4, and both rows are mandatory when the source is live |

Block 2 is the reason `Z-06`'s fifth item is comprehensible: it names, in advance, exactly what
*Clear override* will restore. Block 3 is the age rule, which is *"the rule most likely to be
lost during a build, because dropping the age always makes the layout tidier."*

> **There are still three blocks, and `COVER` did not make a fourth.** It answers block 3's
> question — *how old is this?* — about the one thing on this screen that has an age of its own.
> A fourth block would have cost a fourth `InterElementGap` and put the mandatory `LAST SYNCED`
> row into competition with it; joining block 3 costs **one row**, taken from the three blank
> rows that separate block 2 from the anchored block, and leaves three of them at four. And it
> puts the picture's age and the row's age on adjacent lines, where the difference between
> *two days* and *three hours* is the point.
>
> **`COVER` sits above `LAST SYNCED`, never below.** `LAST SYNCED` stays the view's last line in
> every state, exactly as in rev A, so **D-05-4's anchor did not move** — the block grew upward.
> A builder who bottom-anchors block 3 gets the same final row they got before.

> **D-05-4 · Block 3 is bottom-anchored** to the view's last row when at least one blank row
> would otherwise separate it from the flowed content; otherwise it flows immediately after
> block 2 with one `InterElementGap`. One rule, both hosts. In host B it turns 17 rows of empty
> box into a readout footer; in host A it reads as a status line, which is what it is.

---

## 4 · Mockup — host B, the pane at 120 × 40

44 × 32 (border 1 + inset 2 + **content 38** + inset 2 + border 1), sitting in `Z-04`'s
`112 × 32` body after a 2-column gutter. **Identical content, identical block order, identical
label field. Only the rectangle changed.**

*Rev B redrew this render for two reasons: it gained the `COVER` row, and it had **31** content
rows where §6's table says **30**. It is now 30 — border 1 + 30 + border 1 = the 32 rows `Z-04`'s
body has — with `COVER` and `LAST SYNCED` on rows 29 and 30 and seventeen blank rows above them.*

```
┌ DETAIL ──────────────────────────────────┐
│  Return of the Obra Dinn                 │
│                                          │
│  ◉  ZERADO                               │
│                                          │
│  PLAYTIME     9h                         │
│  LAST PLAYED  2 Aug 2026                 │
│  ADDED        12 Mar 2026                │
│  SOURCE       Steam                      │
│                                          │
│  SET BY       you, 12 Aug 2026           │
│  STEAM SAYS   IN PROGRESS                │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│                                          │
│  COVER        cached, 2 days ago         │
│  LAST SYNCED  3 hours ago                │
└──────────────────────────────────────────┘
```

**When the pane holds focus** — `Tab` or `⏎` from the ledger. The border goes **heavy**; the
pane title goes amber. Under `NO_COLOR` the box weight alone still says which region has focus
(`04-navigation-and-focus.md` §4.2). **This is an excerpt, not a golden — the top 12 of the
pane's 30 content rows.** Every row below it, including the bottom-anchored block 3, is
byte-identical to the render above; the border weight and the title's colour are the whole of
the difference.

```
┏ DETAIL ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃  Return of the Obra Dinn                 ┃
┃                                          ┃
┃  ◉  ZERADO                               ┃
┃                                          ┃
┃  PLAYTIME     9h                         ┃
┃  LAST PLAYED  2 Aug 2026                 ┃
┃  ADDED        12 Mar 2026                ┃
┃  SOURCE       Steam                      ┃
┃                                          ┃
┃  SET BY       you, 12 Aug 2026           ┃
┃  STEAM SAYS   IN PROGRESS                ┃
┃                                          ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

> **D-06-1 (shared with `Z-06`) · A bordered surface is inset by 2 columns each side and 0 rows.**
> Read from `01-design-system.md` §6.2's own ratified anatomy — `│  Return of the Obra Dinn │` —
> where the first content row sits directly under the top border and the left inset is two
> spaces. It is **fixed at 2 at every tier** rather than tracking `InnerPaddingX`, because a
> bordered box that breathes differently at different terminal widths is a box that looks
> resized rather than designed. **Now a spine token:** `02-composition.md` names `BorderInsetX`
> = **2** columns each side and `BorderInsetY` = **0** rows (`14-contradictions-closed.md` #16).

> **The pane is 30 content rows and Phase 1 fills 13 of them.** That is not an unfinished
> screen; it is the honest amount Phase 1 knows, with block 3 anchored to the last row so the
> box reads as a bounded readout rather than a half-filled container. Phase 2's *sinopse* lands
> in exactly this space, rendered by `glamour` **restyled to Zerado tokens** — never glamour's
> default theme, which carries its own palette and would break the colour budget.

---

## 5 · Visual hierarchy

| # | What | Channel |
|---|---|---|
| **1** | **The game's title** | Body row 1, `--z-text` (16.65 — the brightest text on screen), alone with a respiro under it, at `leftInset`. The only sentence-case line above the readouts |
| **2** | **The state chip** | The only coloured mark in the body. `◉ ZERADO` in cyan reads as an event because there is nothing else competing |
| **3** | **The readout values** | `--z-text`, right of a fixed 11-cell label field — a straight vertical edge at column 14 that the eye runs down |
| **4** | **The readout labels** | UPPERCASE, `--z-text-secondary` — the cockpit annotation register, present but subordinate |
| **5** | **Block boundaries** | `InterElementGap`, one row. **Spacing, not rules** — no hairlines in either host |
| **6** | Chrome | breadcrumb, footer, the pane border |

**The one thing the player should see first is the game's title; the second is the state.** That
ordering is the whole reason the title is not in the header band: in host B the band belongs to
`Z-04`, and a hierarchy that inverts between hosts is not a hierarchy.

**No letterspacing on the readout labels.** The brand's readout style tracks out to 0.18em; a
character grid has no sub-cell tracking, so "tracking" would mean inserting literal spaces —
which changes the text content, doubles every label's width at exactly the tiers the canon is
protecting, and hands `L A S T   P L A Y E D` to anything reading the output stream.
UPPERCASE alone carries the role.

---

## 6 · Every applied spacing token, by name

| Token | Tiny | Narrow | Standard | **Wide** | ExtraWide (host B) |
|---|---|---|---|---|---|
| `OuterMarginX` | 0 | 1 | 2 | **2** | 2 |
| `OuterMarginY` | 0 | 1 | 1 | **1** | 1 |
| `InnerPaddingX` | 1 | 1 | 1 | **1** | 2 |
| `InnerPaddingY` | 0 | 1 | 1 | **1** | 1 |
| `InterElementGap` | 1 | 1 | 1 | **1** | 1 |
| `HeaderBand(tier, false)` | **1** | 3 | 3 | **3** | 3 — **`Z-04`'s, not this view's** |
| **`leftInset`** | **1** | 2 | 3 | **3** | 4 — **`Z-04`'s** |
| View content width | 30 | 36 | 54 | **74** | **38** |
| View content height | 21 | 16 | 16 | **16** | **30** |

**Applied:**

| Surface | Token | Value |
|---|---|---|
| title → chip, chip → block 1, between blocks | `InterElementGap` | 1 row, every time |
| Content left edge, host A | `leftInset` | column 4 at Wide |
| Content left edge, host B | `Z-04`'s `leftInset` + list 66 + gutter 2 + border 1 + `BorderInsetX` 2 | column 76 |
| Inside the pane border | `BorderInsetX` (D-06-1) | 2 cols each side, 0 rows |
| Ledger ∥ pane separation | a 2-column gutter **and** the pane's border | never a fill |
| Footer | the canon's reserved footer row | host A only |

**Zero magic numbers.** The readout's 11-cell label field and its 2-cell gutter are this
component's declared geometry (§7.1); `ReadoutBesideMin = 38` is a named threshold (D-05-2).

---

## 7 · Colour, glyph and label for everything shown

Ratios read from the brand manual's measured table (`brand-manual.md` §4.2). None estimated.

| Element | Token | Hex | ANSI-256 | 16-colour | Ratio |
|---|---|---|---|---|---|
| The game's title | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** AAA |
| Readout values | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** AAA |
| Readout labels (UPPER) | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA |
| Explanatory prose | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA |
| `—` · `not tracked` · `never` | `--z-text-tertiary` | `#8492A8` | ***underived*** | `white` | **6.15** AA |
| **`COVER` value — `cached, 2 days ago`** | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** AAA |
| **`COVER` value — `not fetched` · `nothing found`** | `--z-text-tertiary` | `#8492A8` | ***underived*** | `white` | **6.15** AA |
| Screen title `GAME DETAIL` (host A) | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA |
| Pane title `DETAIL`, focused | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA |
| Pane title `DETAIL`, unfocused | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA |
| Pane border, focused | `--z-border-strong` | `#64748B` | **67** | `bright black` | **4.08** — meets 1.4.11 |
| Pane border, unfocused | `--z-border` | `#2A3342` | **236** | `black` | 1.53 — decorative, and it bounds no control |
| Banner `▌`, informational | `--z-border-strong` | `#64748B` | **67** | `bright black` | **4.08** |
| Banner `▌`, action required | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA |

**The state chip — the four states, ratified and unchanged:**

| State | Token | Hex | ANSI-256 | 16-col | Glyph | ASCII | Label | Ratio |
|---|---|---|---|---|---|---|---|---|
| Not started | `--z-state-not-started` | `#A5A29B` | **247** | `white` | `○` U+25CB | `[ ]` | `NOT STARTED` | **7.62** |
| In progress | `--z-state-in-progress` | `#FFB000` | **214** | `bright yellow` | `◐` U+25D0 | `[~]` | `IN PROGRESS` | **10.59** |
| Zerado | `--z-state-zerado` | `#19E0FF` | **45** | `bright cyan` | `◉` U+25C9 | `[*]` | `ZERADO` | **12.15** |
| Abandoned | `--z-state-abandoned` | `#C77DFF` | **177** | `bright magenta` | `⊘` U+2298 | `[x]` | `ABANDONED` | **7.21** |

> ***underived*** — `--z-text-tertiary` `#8492A8` has **no derived ANSI-256 index**. **Interim
> on this screen: tertiary text renders uncoloured**, which invents nothing and degrades exactly
> as `NO_COLOR` does. Nobody may pick an index at the keyboard. Owner: `fft-brand-architect`.
> This screen uses tertiary more than any other in the bundle — it is the colour of *not knowing*
> — so it is the screen most affected by the gap.

**The warm grey `#A5A29B` on `NOT STARTED` is load-bearing engineering, not taste.** The
blue-cast `#9FB0C6` collapsed against the cyan at **ΔE 8.8 under deuteranopia**; the warm grey
measures **25.8**. Never corrected back toward blue.

### 7.1 · The readout geometry

```
PLAYTIME     9h
│            │
│            └ value field — flex, --z-text, starts at content column 14
└ label field — 11 cols, UPPERCASE, --z-text-secondary, then a 2-col gutter
```

11 is the longest fixed label (`LAST PLAYED`, `LAST SYNCED`). A provider-named label —
`STEAM SAYS`, and later `PLAYSTATION SAYS` — may exceed 11; when it does the **field grows to
fit the longest label in the block**, and every value in that block moves with it, so the
vertical edge is never broken. The field never differs between blocks on one render.

### 7.2 · Audio — the fourth channel

The annunciator is on the **header band's title row** (D-A1, `Z-04-library.md` §7.1) — which in
host B is `Z-04`'s row and in host A is this frame's. **This view never draws it**; the host
does. `▮ AUDIO` `--z-primary` / `▯ MUTED` `--z-text-secondary`, both **Neutral**-width glyphs
verified against UCD 16.0.0.

**`Z-05` has no sound event of its own.** Nothing on this screen changes without a keypress, so
there is nothing here that could make a noise. **If a sound is ever attached to arriving on this
screen, it must have a visible carrier in this spec first** — audio is never the only channel.

---

## 8 · The full state table

The whole screen is a state table; §9 draws the ones that matter.

| # | State | Trigger | Composition | Copy |
|---|---|---|---|---|
| **D1** | **First run — unreachable by construction** | Library empty | **`Z-05` cannot be reached.** `⏎` on an empty ledger does nothing; there is no row to open. `Z-04`'s empty state owns the screen | — |
| **D2** | **The normal Phase 1 state — no enrichment exists yet** | Always, in Phase 1 | §3. Three blocks, all populated. **The Phase 2 blocks are absent, not empty-labelled** | §10.1 |
| **D3** | **Loading** | The row is known, the detail is being read | Title and chip render **immediately** — both are already in memory from the ledger row. Values render `—` until the read returns. **Never a spinner** | §10.6 |
| **D4** | **A hand-added copy — no playtime SIGNAL** | `SOURCE = physical`, `Capabilities.Progress = false` | `PLAYTIME` and `LAST PLAYED` read **`not tracked`**; block 3 is **absent**; a two-line prose block explains why | §10.2 |
| **D5** | **Not fetched yet** | The row arrived in a sync that ended before details did | `PLAYTIME` and `LAST PLAYED` read **`—`**; `LAST SYNCED` reads **`never`**; the banner offers `r` | §10.3 |
| **D6** | **Known to be empty** | The provider answered, and the answer was nothing | `LAST PLAYED` reads **`never played`**; `PLAYTIME` reads **`0h`**. No retry is offered, because there is nothing to retry | §10.4 |
| **D7** | **Offline / a failed sync** | Classified failure (`07-offline-contract.md` §5) | Degrade banner at body row 1; every value renders normally — **this screen is `WORKS`**; block 3 carries the age | `Z-04` §8.1 |
| **D8** | **A manual override is set** | `status_manual IS NOT NULL` | Block 2 shows `SET BY  you, <date>` **and** `<PROVIDER> SAYS  <STATE>` | §10.1 |
| **D9** | **No override — the state is derived** | `status_manual IS NULL` | Block 2 shows `SET BY  Zerado, from 9h played` and **omits** `<PROVIDER> SAYS`, because the provider's answer *is* the state | §10.5 |
| **D10** | **A pending write** | `Z-06` applied, SQLite not yet confirmed | The chip label is suffixed `…` — `ZERADO…` | §10.7 |
| **D11** | **The title overflows the host** | Any host | **Wraps at a word boundary. Never truncated.** It is the identity (R-10(a)) and, unlike a ledger row, there is room | §9.3 |
| **D12** | **Content taller than the host** | Narrow · Tiny · a long title | `bubbles/viewport` scrolls it. The footer says `↑↓ scroll` | — |
| **D13** | **Below the refusal floor** | `< 24` cols or `< 8` rows | The frame refuses for the whole program — `Z-04-library.md` §11.3 | — |
| **D14** | **The game is absent — a sync stopped returning it** | `absent_since IS NOT NULL` (`06-data-seams.md` §2.4) | §9.6. `SOURCE` carries `<Provider> — absent since <date>`; a **two-line prose block** after block 2 says why. **No banner, no retry, no red, no fifth state.** Every other value renders exactly as D2, and **block 3 is `LAST SYNCED` alone** — D-05-8 | §10.8 |
| **D15** | **A cover is cached** | `metadata.cover_ref` resolves on disk | `COVER  cached, 2 days ago` — **the age, in the register `LAST SYNCED` uses**. `--z-text`, because it is a value Zerado has | §9.7 |
| **D16** | **A cover has not been fetched** | `cover_ref IS NULL` and never attempted · a fetch errored · the OS evicted the cache | `COVER  not fetched` — `--z-text-tertiary`. **`r` is offered**, and this is the only cover state the player can act on | §9.7 |
| **D17** | **The provider has no art for this game** | the fetch succeeded and returned nothing | `COVER  nothing found` — `--z-text-tertiary`. **No retry**, no apology, nothing wrong. The provider answered | §9.7 |

> **D15–D17 are three facts, and the player can act on exactly one of them.** That is the same
> shape as `—` versus `0h` versus `not tracked` in §9.3, and it is carried the same way: **in
> words**, in the same three words the deck's tiles use, so the vocabulary is one product's and
> not one screen's. **A fourth fact — a game that will never have a cover — has no `COVER` row
> at all**, because `SOURCE  Added by hand` already is that row (D4, and **D-05-8**).

> **D1 is the honest first-run row, and it is a *nothing*.** The manual demands a first-run row
> on every screen; the truthful one here is that this screen has no first-run appearance,
> because it needs a game and first run has none. Recorded rather than faked. **`⏎` on an empty
> ledger must do nothing and show no error** — it must not push an empty detail view.

---

## 9 · The honest mockups

### 9.1 · D4 — a hand-added copy: no playtime SIGNAL, which is not a missing value

`physical` is a **provider in its own right** with `Capabilities{Sync: false, Playtime: false,
LastPlayed: false}` (`06-data-seams.md` §2), not a flag on a Steam-shaped row. So playtime here
is not unknown, not zero, and not pending: **there is nothing that could report it.** Rendering
`0h` would be a lie; rendering `—` would promise that a sync could fix it.

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Game detail                                                         │
│                                                                                │
│   GAME DETAIL                                                                  │
│                                                                                │
│                                                                                │
│   Chrono Trigger                                                               │
│                                                                                │
│   ○  NOT STARTED                                                               │
│                                                                                │
│   PLAYTIME     not tracked                                                     │
│   LAST PLAYED  not tracked                                                     │
│   ADDED        4 Feb 2026                                                      │
│   SOURCE       Added by hand                                                   │
│                                                                                │
│   SET BY       nobody yet                                                      │
│                                                                                │
│   Steam reports playtime; a copy you added by hand has no store                │
│   behind it to report one. Zerado does not guess, so every state               │
│   on this one is yours to set.                                                 │
│                                                                                │
│                                                                                │
│   ↑↓ scroll   s status   esc back   ? help   q quit                            │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

**Block 3 is absent — both rows of it.** A provider that is not a `Syncer` has no last-sync time,
and labelling one `never` would invite a retry that cannot exist. *Omit the block; do not label
it empty.* **`r` is not in the footer** on this state, for the same reason.

**And that is where the deck's third word goes.** `Z-15` §11.2 puts `added by hand` on the tile,
because a 17 × 6 tile has a caption and a chip and nothing else — it has no `SOURCE` row, so the
word is the only way it can say what this is. **This screen already said it**, four rows up, in
`SOURCE  Added by hand`. A `COVER  added by hand` row beneath it would be the same fact twice in
one 16-row body, and the prose block already explains the whole of it: *a copy you added by hand
has no store behind it.* **One fact, one row, and the row that owns it is the one that names the
source** — which is D-05-7's rule, running in the other direction.

**`SOURCE  Added by hand`** — the ratified phrasing (`01-design-system.md` §6.4). *"A physical
copy isn't a second-class row in the list"* is a published promise; this screen keeps it by
giving a hand-added copy the **same three blocks minus the one that cannot exist**, not a
smaller screen.

### 9.2 · D5 — not fetched yet, which the player CAN act on

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Game detail                                                         │
│                                                                                │
│   GAME DETAIL                                                                  │
│                                                                                │
│                                                                                │
│   ▌ SYNC INCOMPLETE   Steam has not been asked about this one yet.             │
│                                                                                │
│   Metroid Dread                                                                │
│                                                                                │
│   ○  NOT STARTED                                                               │
│                                                                                │
│   PLAYTIME     —                                                               │
│   LAST PLAYED  —                                                               │
│   ADDED        9 Aug 2026                                                      │
│   SOURCE       Steam                                                           │
│                                                                                │
│   SET BY       nobody yet                                                      │
│                                                                                │
│                                                                                │
│   COVER        not fetched                                                     │
│   LAST SYNCED  never                                                           │
│   ↑↓ scroll   r sync   s status   esc back   ? help   q quit                   │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

Banner **64** cells. **`r sync` is in the footer here and is not in §9.1's** — that single
difference is the whole point of the distinction: the footer offers the retry exactly when a
retry exists.

**Block 3 says `not fetched` and `never` on the same two lines, and neither is redundant.**
`LAST SYNCED  never` says *nobody has asked Steam about this row*; `COVER  not fetched` says
*nobody has fetched the picture*. They agree here because one sync would fix both — but they are
independent facts and they part company constantly: a library that syncs hourly and an
enrichment run that failed gives `LAST SYNCED  3 hours ago` beside `COVER  not fetched`, which is
**D16 arriving on an otherwise complete row**. That is why `r` is now offered on two conditions
and not one (§14).

### 9.3 · The five renderings of *not a value* — the table this screen exists for

| Fact state | Value renders | Colour | Block 3 | Retry offered? | Why |
|---|---|---|---|---|---|
| **Known** | `9h` · `2 Aug 2026` | `--z-text` | present | — | |
| **Known to be empty** | `0h` · `never played` | `--z-text-secondary` | present | **no** | The provider answered. There is nothing to fetch |
| **Not fetched yet** | **`—`** | `--z-text-tertiary` | `never` | **yes — `r`** | Nobody has asked. This is the one the player can act on |
| **Unavailable right now** | `—` and the banner names it | `--z-text-tertiary` | the age | **yes — `r`** | The last value, with its age, or nothing plus the reason |
| **Not trackable at all** | **`not tracked`** | `--z-text-tertiary` | **absent** | **no** | The provider has no such capability. Not a fault, not a gap |

**The sixth rendering is not a value at all — it is a fact about the row** (D14, §9.6). Every
value above still renders normally; what changes is that `SOURCE` carries `— absent since <date>`
and a prose block explains it. It belongs beside this table rather than in it, because the other
five answer *"what does Zerado know about this field?"* and this one answers *"is this row still
being reported?"*

**And the seventh is not a value either — it is a fact about an ASSET** (D15–D17, §9.7). The
`COVER` row does not answer *what does Zerado know about this field?*; it answers *does Zerado
have this file, and how old is it?* **That is why it does not use `—`.** A `—` is the rendering
for a value with a shape — a number, a date — that nobody has asked for yet; a cover is a file
that either exists on disk or does not, and the deck already has three words for which kind of
not-having it is. **The words win**, for the same reason `—` and `not tracked` beat a colour
difference: `Z-15`'s tile cannot render `—` in a 15-column box without it reading as a dead
cell, so if this screen used `—` the product would have two vocabularies for one fact. **One
word, two surfaces.**

**`—` versus `not tracked` is the load-bearing distinction on this screen**, and it is carried in
**words**, so it survives `NO_COLOR`, a screenshot, and a screen reader reading the raw stream.
A colour difference alone would not.

### 9.4 · Narrow — 40 × 24, labels above values

```
┌────────────────────────────────────────┐
│                                        │
│  Zerado ✦ Game detail                  │
│                                        │
│  GAME DETAIL                           │
│                                        │
│                                        │
│  Sid Meier's Civilization VI:          │
│  Gathering Storm                       │
│                                        │
│  ○  NOT STARTED                        │
│                                        │
│  PLAYTIME                              │
│  0h                                    │
│                                        │
│  LAST PLAYED                           │
│  never                                 │
│                                        │
│  ADDED                                 │
│  2 Jan 2026                            │
│                                        │
│  SOURCE                                │
│  Steam                                 │
│  ↑↓ scroll  esc back  ? help           │
│                                        │
└────────────────────────────────────────┘
```

Body **36 × 16**. The title wraps at a word boundary and is **never truncated**. Blocks 2 and 3
are below the fold; `↑↓` scrolls to them, and the footer says so.

> **D-05-2 · The readout arrangement is chosen by the host's content width, at one named
> threshold: `ReadoutBesideMin = 38`.** At or above it, label and value share a row (host B 38 ✓,
> Standard 54 ✓, Wide 74 ✓). Below it, the label sits on its own row above the value (Narrow 36,
> Tiny 30) — which is exactly what `03-responsive.md` §3 specifies for Narrow.
> **Why a width and not a tier:** this view has two hosts and the pane's content width (38) is
> *narrower than Standard's body* (54) while being *wider than Narrow's* (36). A tier-based rule
> is not expressible in a host-agnostic view; a width-based one is, and it agrees with the
> spine's table at the tested breakpoint.

### 9.5 · Tiny — 32 × 24

```
┌────────────────────────────────┐
│ GAME DETAIL                    │
│                                │
│ Return of the Obra Dinn        │
│                                │
│ ◉  ZERADO                      │
│                                │
│ PLAYTIME                       │
│ 9h                             │
│                                │
│ LAST PLAYED                    │
│ 2 Aug 2026                     │
│                                │
│ ADDED                          │
│ 12 Mar 2026                    │
│                                │
│ SOURCE                         │
│ Steam                          │
│                                │
│ SET BY                         │
│ you, 12 Aug 2026               │
│                                │
│ STEAM SAYS                     │
│ IN PROGRESS                    │
│ esc back  ? help  q quit       │
└────────────────────────────────┘
```

Body **30 × 21**. `OuterMarginX` and `OuterMarginY` shed to **0**; the band is the title row
alone. The chip does **not** shrink and the label is **never** dropped.

*Rev B redrew this render: it had **20** body rows where the table above says **21**, and it
spent the last two on blanks. The twenty-first row is `STEAM SAYS  IN PROGRESS`'s value, so the
fold now falls between block 2 and block 3 rather than inside a run of empty rows.* **`COVER`
and `LAST SYNCED` are the next two readouts below it**, reached by `↓` — the same place block 3
has always been at this width, and neither row is dropped (§11).

### 9.6 · D14 — the game is absent, which is a fact about the SOURCE

**This is the screen `06-data-seams.md` §2.4 names.** A game a sync stops returning is
**tombstoned, never deleted**: `absent_since` is set and the row stays, because *"the provider
stopped returning it"* and *"the player no longer owns it"* are not the same fact and only the
first is observable. `Z-04` excludes the row from the default view; `Z-07`'s `[ABSENT]` chip is
how the player got here; **this screen is the one that says plainly why.**

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Game detail                                                         │
│                                                                                │
│   GAME DETAIL                                                                  │
│                                                                                │
│                                                                                │
│   Alan Wake                                                                    │
│                                                                                │
│   ◉  ZERADO                                                                    │
│                                                                                │
│   PLAYTIME     19h                                                             │
│   LAST PLAYED  2 Jul 2026                                                      │
│   ADDED        12 Mar 2026                                                     │
│   SOURCE       Steam — absent since 14 Aug 2026                                │
│                                                                                │
│   SET BY       you, 3 Aug 2026                                                 │
│   STEAM SAYS   IN PROGRESS                                                     │
│                                                                                │
│   Zerado cannot tell whether you still own it, so it changed nothing.          │
│   The row and its status stay, out of the list, until Steam returns it.        │
│                                                                                │
│   LAST SYNCED  3 hours ago                                                     │
│   ↑↓ scroll   s status   r sync   esc back   ? help   q quit                   │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

> **D-05-7 · The absence is a value on the `SOURCE` row, not a banner — and the prose block does
> the explaining.**
>
> **Why not a banner.** A degrade banner has four mandatory parts (`07-offline-contract.md` §3):
> an uppercase label word, **what is unavailable**, **how stale it is**, and **the key that
> retries it**. An absent row has none of them. Nothing is unavailable — every field on this
> screen renders normally. Nothing is stale — the last sync was three hours ago and it succeeded.
> And there is no retry: `r` re-syncs the library, and if Steam has not started returning the game
> it will not, so a banner offering it would be the screen promising a fix it does not have. A
> banner would also cost two body rows (itself plus its respiro) and compete for `Z-04` §8.1's
> one-banner-at-a-time slot with degradations that *are* urgent.
>
> **Why the `SOURCE` row.** Because that is the fact. `SOURCE` is the row that says where this
> game comes from; *the source stopped returning it* is a statement about the source, and it
> lands in the line the eye is already reading for exactly that. **It costs zero rows**, which is
> what lets block 3 keep its respiro on the one screen where the age matters most.
>
> **Rendering.** `Steam` stays `--z-text` (16.65); the clause `— absent since 14 Aug 2026` is
> `--z-text-secondary` (`#A9B5C7` / 249 / `white`, 9.36) — a weight step inside one value, the
> same technique as the summary's numerals against its words, and **no new token**. The date is
> said the way people say them (§10.8's rule), never an ISO timestamp. Under `NO_COLOR` the
> clause is unchanged: **it is carried entirely by words.**
>
> **`—` U+2014 is Ambiguous** (verified, UCD 16.0.0). It sits inside a flex value that is
> **measured, not counted**, so it cannot shear a column.

**Row arithmetic, both forms — 16 body rows, and it closes without scrolling:**

| | With an override (drawn above) | No override (D9 form) |
|---|---|---|
| title · gap · chip · gap | 1–4 | 1–4 |
| block 1 — four rows, `SOURCE` carrying the clause | 5–8 | 5–8 |
| gap | 9 | 9 |
| block 2 | 10–11 | 10 |
| gap | 12 | 11 |
| the prose block | 13–14 | 12–13 |
| gap | 15 | 14–15 |
| **block 3, bottom-anchored (D-05-4)** — **`LAST SYNCED` alone** | **16** | **16** |

> **This render has no `COVER` row, and the reason is the age rule rather than the arithmetic —
> D-05-8.** On an absent row the `SOURCE` line already carries an age, and it is the **outer**
> one: `absent since 14 Aug 2026` bounds every fact on this screen. Nothing Steam holds about
> this game has reached Zerado since that date — not the playtime, not the state, and not the
> picture. `COVER  cached, 5 months ago` beneath it would state a second, narrower age inside
> one the screen has already given.
>
> **The arithmetic is real too, and hiding it would be worse than stating it.** Rev A's D14 used
> all 16 body rows with no slack — the row this table calls *spare* is the `InterElementGap`
> between the prose and block 3, not a spare row at all. A second block-3 row therefore measures
> **17 in a 16-row view**, and the line that falls below the fold is `LAST SYNCED` — *"the rule
> most likely to be lost during a build."* **A spec that spends the mandatory row on a
> redundant one has lost the argument twice.** Both reasons point the same way, which is why
> this is a decision and not a concession.

**Neither form scrolls at 80 × 24**, which is the whole reason the banner was declined. In host B
(38 content columns) the prose re-wraps to four lines in a 32-row pane and is **never truncated**,
exactly as D4's block is.

**`absent` is not a fifth state and this screen renders no fifth glyph.** *Alan Wake* is
`◉ ZERADO` and its chip is byte-identical to any other `ZERADO` chip in the product — which is
`06-data-seams.md` §2.4's point: *"a game you finished and no longer own is exactly the row you would be angriest to
lose."* `STEAM SAYS  IN PROGRESS` still renders, because the provider's last answer is still a
fact about where the status came from.

**When it comes back, nothing announces it.** `absent_since` is cleared silently (`06-data-seams.md` §2.4): the
`SOURCE` clause stops rendering, the prose block disappears, the row returns to `Z-04`'s default
view, and the screen is D2 again. **No banner, no result line, no celebration** — the *zerado*
line (`Z-04-library.md` §10.5) is the only moment Zerado announces, and this is not it.

### 9.7 · The `COVER` row — every one of `Z-15`'s seven picture states, and what this screen says

The deck and the detail view must never disagree about one game's cover, so the mapping is
written out rather than left to be inferred. **The three words are
[`Z-15-cover-deck.md`](./Z-15-cover-deck.md) §11.2's, character for character** — this screen
adopts them, it does not coin a fourth.

| `Z-15` state | What the tile shows | **This screen's `COVER` row** | Colour | `r` offered? |
|---|---|---|---|---|
| **C3** image support, cover cached | the picture | **`cached, 2 days ago`** | `--z-text` | — |
| **C4 / C10** not fetched yet · the OS evicted the cache | `not fetched` | **`not fetched`** | `--z-text-tertiary` | **yes** |
| **C5** a fetch failed | `not fetched` | **`not fetched`** | `--z-text-tertiary` | **yes** |
| **C6** the provider has no art | `nothing found` | **`nothing found`** | `--z-text-tertiary` | **no** |
| **C7** a hand-added copy | `added by hand` | **no row at all** — `SOURCE  Added by hand` is it, and block 3 is absent (§9.1) | — | **no** |
| **C8** offline, cover cached | the picture | **`cached, 2 days ago`** — unchanged; the cache is local | `--z-text` | — |
| **C9** offline, cover not cached | `not fetched` | **`not fetched`**, with the `OFFLINE` banner naming why nothing new will arrive | `--z-text-tertiary` | **yes** |

**`Z-15` has no picture state for an absent row**, because absence is a fact about the row and
not about the picture — its **C24** says so outright (*"`absent` is not a fifth state and gets no
fifth tile treatment"*). This screen has the same answer in the other direction: **D14 renders no
`COVER` row**, because `SOURCE  Steam — absent since 14 Aug 2026` is already the outer age that
bounds the picture's. §9.6 and **D-05-8**.

**Three things this table settles, each of which someone would otherwise get wrong:**

1. **A failed fetch says `not fetched`, not `failed`.** `Z-15` **D-15-10**: *"a failure is a
   property of the network, not of the game"*, and it is screen-wide, so it is the banner's.
   The `COVER` row states what is true **about this game** — the picture is not here — and never
   blames the game for the connection. **C4, C5, C9 and C10 are one string** for exactly that
   reason, and the four are told apart by the banner, not by the row.
2. **The row is identical on a terminal that cannot draw images.** Whether *this* terminal can
   paint a cached blob is a fact about the terminal, and it lives in the capability note and in
   `Z-09 § DISPLAY` — never here. A blob on disk is cached whether or not anything can render
   it, so **`COVER  cached, 2 days ago` renders on Terminal.app exactly as it does on Ghostty.**
   This falls straight out of `Z-15` §8.1: because Zerado never draws a mark on a cover, a
   cover's provenance is carried entirely in **text**, and text renders everywhere.
3. **This screen never draws the picture.** It gained a **row**, not an image region. Whether
   the detail pane ever paints a cover is a composition question — a new region inside a
   bordered surface, in a view that is mounted in two hosts of different widths — and
   composition belongs to `fft-tui-architect`. `Z-15` §4 keeps the pane at 44 × 32 with the
   deck beside it and draws no cover in it; this spec follows that and proposes nothing.

**And one thing it deliberately leaves alone: `Z-15` §12.2's silence at Tiny.** Cover *mode* is
unavailable below 40 columns and nothing is printed about it. **The `COVER` row is not covered by
that silence and renders at every tier**, because the two are different facts: the note is about
the terminal — a thing the player cannot act on without changing both program and window — while
this row is about the data, and `r` works at 32 columns exactly as it works at 120. At Tiny the
row is below the fold in the stacked arrangement (§9.5) and one `↓` reaches it, which is where
every block-2 and block-3 row already is.

---

## 10 · The exact copy — ready to paste

### 10.1 · D2 / D8 — the normal state, an override set

```
Return of the Obra Dinn

◉  ZERADO

PLAYTIME     9h
LAST PLAYED  2 Aug 2026
ADDED        12 Mar 2026
SOURCE       Steam

SET BY       you, 12 Aug 2026
STEAM SAYS   IN PROGRESS

COVER        cached, 2 days ago
LAST SYNCED  3 hours ago
```

### 10.2 · D4 — a hand-added copy

```
PLAYTIME     not tracked
LAST PLAYED  not tracked
ADDED        4 Feb 2026
SOURCE       Added by hand

SET BY       nobody yet

Steam reports playtime; a copy you added by hand has no store
behind it to report one. Zerado does not guess, so every state
on this one is yours to set.
```

At 38 columns (host B) the prose re-wraps to five lines; it is **never truncated**.

### 10.3 · D5 — not fetched yet

```
▌ SYNC INCOMPLETE   Steam has not been asked about this one yet.

PLAYTIME     —
LAST PLAYED  —

COVER        not fetched
LAST SYNCED  never
```

### 10.4 · D6 — known to be empty

```
PLAYTIME     0h
LAST PLAYED  never played

COVER        nothing found
```

**`nothing found` is D6's shape applied to an asset** — the provider answered, the answer was
nothing, and there is nothing to retry. It is the cover's `never played`: two words that say a
question was asked, so that silence does not have to be read as a gap.

`never played` rather than `never`: `05-state-machine.md` §6 is explicit that a null
`last_played_at` means *not reported*, not *never played*, so the two words are what separate the
fact from the absence of one.

### 10.5 · D9 — no override; the state is derived

```
SET BY       Zerado, from 9h played
```
```
SET BY       Zerado, no playtime reported
```
```
SET BY       nobody yet
```

The third form is for a provider that cannot report playtime at all — there is no derivation
input, so nothing derived it either. `<PROVIDER> SAYS` is **omitted** in all three: when there is
no override, the provider's answer *is* the state and repeating it below the chip would be
furniture.

### 10.6 · D3 — loading

```
Return of the Obra Dinn

◉  ZERADO

PLAYTIME     —
LAST PLAYED  —
ADDED        —
SOURCE       —

COVER        —
LAST SYNCED  —
```

**`COVER` renders `—` while loading and only while loading** — the one moment the row does use
the value vocabulary, because during a read Zerado genuinely does not yet know which of D15,
D16 or D17 this game is in. The moment the read returns it becomes a word.

The **title and the chip render immediately** — both came from the ledger row the player was
already looking at. No spinner, no scanner: reading a local row is not an indeterminate wait.

### 10.7 · D10 — a pending write

```
◉  ZERADO…
```

The label carries the `…` until SQLite confirms. **Never an optimistic silent change.**

### 10.8 · D14 — the game is absent

**The copy `06-data-seams.md` §2.4 authorises.** Rev A of this spec refused to write it, on the
grounds that a decision about a player's own data is not a screen's to make. `06-data-seams.md` §2.4 makes it — *"a
game a sync stops returning is tombstoned, never deleted"* — and names the refusal as the right
call. This is the copy.

```
SOURCE       Steam — absent since 14 Aug 2026
```
```
Zerado cannot tell whether you still own it, so it changed nothing.
The row and its status stay, out of the list, until Steam returns it.
```

**45 cells for the readout line** (label field 11 + gutter 2 + a 32-cell value) **· 67 and 69
for the prose.** At 38 content columns (host B) the prose re-wraps to four lines; it is
**never truncated**, the same rule as D4's block.

**What each clause is doing, because every one of them is a decision:**

| Clause | Why it is there | What it deliberately does not say |
|---|---|---|
| `absent since 14 Aug 2026` | **What happened, and when.** The date is the one fact the player can check against their own memory | Not `removed`, not `deleted`, not `missing` — none of those is observable |
| `Zerado cannot tell whether you still own it` | **Why.** *"The provider stopped returning it"* and *"the player no longer owns it"* are not the same fact and only the first one is observable (`06-data-seams.md` §2.4) | It does not guess which. It does not list delistings, region changes and family shares — naming causes for *this* game would be a guess wearing a fact's clothes |
| `so it changed nothing` | **The consequence, stated as an action the product did not take.** *"Deletion is irreversible. Tombstoning is not"* | It does not apologise, and it does not ask the player to do anything |
| `The row and its status stay` | **What it means for their data**, which is the sentence that had to be right | It says *its* status rather than *your* status, because a derived status is not theirs and this line must be true in both cases |
| `out of the list` | Why they cannot find it on `Z-04` — the thing they will actually notice | |
| `until Steam returns it` | It is **reversible, and reverses by itself.** `absent_since` is cleared the moment the game comes back | It does not offer a key. `r` will not fix this, and a copy line pointing at one would be the screen promising a fix it does not have |

**The register is the sync refusals'** (`07-offline-contract.md` §3): the fact, then the
consequence, in two short sentences. **No exclamation mark, no apology, no emoji, and the player
is never told they did anything** — because they did not.

**Nothing here claims an unbuilt capability.** There is no *"remove it"* affordance offered:
`06-data-seams.md` §2.4 reserves deletion for when the **player** asks, and Phase 1 binds no key that asks. Offering
one in copy would be anti-pattern 14.

### 10.9 · Copy notes

- **Dates** are said the way people say them — `12 Mar 2026`, `3 hours ago`, `just now`,
  `last June`. Never `2026-08-22T04:11:09Z`. The exact timestamp belongs in this view for anyone
  who wants it, and Phase 1's answer is that `ADDED` and `LAST PLAYED` are already the precise
  dates; only `LAST SYNCED` is relative, because relative is what "how stale is this" means.
- **Say the number** — `9h`, `41 hours`, `147 of 247`, `2 days ago`.
- **The three cover words are `Z-15` §11.2's, verbatim, and this screen coins nothing.**
  `not fetched` · `nothing found` · `added by hand`. The tile and the detail row say the same
  string about the same game, which is the `absent` lesson applied a second time: **one word,
  every surface it appears on.** Their i18n hazard is the tile's, not this screen's — a
  translator has 15 columns in a tile and 60 here — so the length budget that governs them lives
  with the tighter field.
- **No exclamation marks. No emoji. The user is never a "gamer".**
- **Casing** — `Zerado` the product in prose and in the breadcrumb; `ZERADO` the chip; `zerado`
  never appears on this screen because there is no summary sentence here.
- **Nothing on this screen claims an unbuilt capability.** There is no `MOOD` label, no
  `SINOPSE` label, no cover-art placeholder and no "coming in Phase 2" line. The roadmap lives
  on the landing page.
- **Zerado is a games product and the copy says so.** `PLAYTIME`, `LAST PLAYED`,
  `Steam says`, `Added by hand`. `DETAIL`, `SOURCE` and `ADDED` are neutral because each names a
  **field**, not because the copy is leaving room for another kind of thing
  (`11-media-model.md` rev B).

---

## 11 · 40-column behaviour, and the refusal floor

**40 × 24** is §9.4: body `36 × 16`, readout **stacked** (D-05-2, `36 < 38`), title wrapped, the
chip unchanged at its full 14 columns, blocks 2 and 3 reachable by scrolling.

**Sheds at Narrow and below:** the beside arrangement (the label moves above the value).
**Never sheds, at any width:** the state — glyph **and** label · the game's title · the
degrade banner when one is active · the footer key line · **both rows of block 3**, `COVER` and
`LAST SYNCED`, which go below the fold at Narrow and Tiny and are reached by `↓` — never dropped,
never abbreviated, and never merged into one row.

**32 × 24** is §9.5: body `30 × 21`, `leftInset` 1, band collapsed to the title row.

**The refusal floor is the program's, not this screen's** — below **24 columns or 8 rows**
Zerado renders one frameless sentence and, at start-up, exits `2`:

```
Zerado needs at least 24 columns and 8 rows. This terminal is 20 x 6.
```

A running session resized below the floor shows the same sentence and **keeps running**.

---

## 12 · `NO_COLOR` — rendered, not asserted

Zero SGR sequences. The §3 screen, character for character:

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                                                                                │
│   Zerado ✦ Game detail                                                         │
│                                                                                │
│   GAME DETAIL                                                                  │
│                                                                                │
│                                                                                │
│   Return of the Obra Dinn                                                      │
│                                                                                │
│   ◉  ZERADO                                                                    │
│                                                                                │
│   PLAYTIME     9h                                                              │
│   LAST PLAYED  2 Aug 2026                                                      │
│   ADDED        12 Mar 2026                                                     │
│   SOURCE       Steam                                                           │
│                                                                                │
│   SET BY       you, 12 Aug 2026                                                │
│   STEAM SAYS   IN PROGRESS                                                     │
│                                                                                │
│                                                                                │
│                                                                                │
│   COVER        cached, 2 days ago                                              │
│   LAST SYNCED  3 hours ago                                                     │
│   ↑↓ scroll   s status   r sync   esc back   ? help   q quit                   │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

| Information | Carried without colour by |
|---|---|
| The state | `◉` **and** `ZERADO` |
| Label vs value | **Case** — `PLAYTIME` against `9h` — plus the fixed column edge at 14 |
| Block boundaries | `InterElementGap`, spacing, not colour |
| **Not-fetched vs not-trackable** | The **words** `—` and `not tracked`, plus the presence or absence of block 3 and of `r sync` in the footer |
| **Which kind of cover Zerado has** | The **words** `cached, 2 days ago` · `not fetched` · `nothing found`, plus the absence of the row entirely on a hand-added copy. **Colour carries nothing here** — the value/tertiary split is emphasis, and every distinction survives without it |
| Which region has focus (host B) | `┏━┓` against `┌─┐` — box-drawing **weight** |
| The banner class | The label word |
| Screen title | UPPERCASE + bold, alone on its row |

**Under `NO_COLOR` this screen loses nothing**, which is `02-colour-budget.md` §3.3's cross-check: run it with
`NO_COLOR=1` and if any information disappears, the screen was encoding meaning in colour.

**`ZERADO_ASCII=1`:**

```
[*]  ZERADO

PLAYTIME     9h
```

---

## 13 · The focus model, and `Esc`

### 13.1 · What can hold focus

**In host A:** the view is a single scrollable region. There is exactly one focusable thing —
the viewport — and it holds focus from the moment the route is pushed. There is **no `Tab`**,
and the footer does not list one.

**In host B:** the view is one of `Z-04`'s two regions. It takes focus by `Tab` or by `⏎` on the
ledger, and gives it back with `Tab`, `Shift-Tab` or `Esc`.

### 13.2 · How focus is shown

| Host | Focused | Not focused |
|---|---|---|
| **A** — route | the whole screen is the focused region; the viewport scrolls | n/a |
| **B** — pane | **`┏━┓`** `--z-border-strong` border + **amber** pane title | `┌─┐` `--z-border` border + `--z-text-secondary` title |

Two channels survive `NO_COLOR`: the **box-drawing weight**, and the ledger's row weight going
from bold to normal on the other side of the gutter. **The focus ring is never removed** — not
in "read-only" panes, not while syncing, not for aesthetics. In a terminal there is no pointer;
the ring is the only way to know where you are.

### 13.3 · `Esc` — the one thing this view must not decide

| Host | `Esc` | Then |
|---|---|---|
| **A — route** | **Pop the route** | `Z-04`, focus restored to the row the player left, **not** the top of the list |
| **B — pane** | **Return focus to the ledger** | The pane stays rendered and keeps showing the same game |

**The view does not implement either.** It reports that it has nothing of its own to dismiss and
the host maps the key. This is the one place a host-agnostic view is most likely to be broken by
a convenience `if`.

`Esc` **always leaves.** There is no keyboard trap here (WCAG 2.1.2). This screen has no text
input, so 2.1.4 is not engaged — but the keys it does own are listed in §14 and nowhere does a
single key destroy anything.

### 13.4 · Focus restoration is the most-felt detail in a 400-row library

`04-navigation-and-focus.md` §4.1 rule 4: the route stack stores each frame's focus state, so
returning from `Z-05` puts the player back on the row they left. Rule 5: a rebuild never moves
focus — cursor and offset are preserved by **game identity**, not by index.

---

## 14 · The key map

| Key | Does | Host A | Host B | Note |
|---|---|---|---|---|
| `↑` `↓` `k` `j` | Scroll the view | ✓ | ✓ | Only when there is overflow; listed only then |
| `Ctrl-D` / `Ctrl-U` | Half a page | ✓ | ✓ | |
| `g` / `G` | Top / bottom | ✓ | ✓ | |
| `s` | Set this game's status → `Z-06` | ✓ | ✓ | The overlay is drawn over whichever host is rendering |
| `r` | Re-sync → `Z-03` | ✓ | ✓ | **Listed only when a retry exists** — never on a hand-added copy (§9.1). **Rev B: two conditions, not one** — a value that is `—`, **or** a `COVER` row reading `not fetched` (`Z-15` §10: `r` is *"also what refetches a `not fetched` cover"*). Never for `nothing found`, which is an answer |
| `Tab` / `Shift-Tab` | Move between regions | **not bound** | ✓ | Host A has one region |
| `⏎` | — | **not bound** | **not bound** | There is nothing to activate here in Phase 1 |
| `,` | Settings → `Z-09` | ✓ | ✓ | Global |
| `m` | Mute / unmute | ✓ | ✓ | Global; **only when audio has been enabled** |
| `?` | Help → `Z-10` | ✓ | ✓ | `Z-10` lists **this** screen's keys, because it is where the player came from |
| `Esc` | §13.3 | ✓ | ✓ | **The host decides what it means** |
| `q` | Quit | ✓ | ✓ | Immediate; nothing to confirm |
| `Ctrl-C` | Quit | ✓ | ✓ | Always |

**Reserved and inert in Phase 1:** `n` / `p` — *next / previous game from within the detail
view*. They are reserved precisely for this screen and deliberately unbound until Phase 2 gives
the detail view enough content to make paging worthwhile. Pressing one does nothing and shows no
error. **`Z-10` lists them under RESERVED**, so a player who tries them learns why.

**Footer strings, exact:**

| State | Footer | Cells |
|---|---|---|
| Normal, host A | `↑↓ scroll   s status   r sync   esc back   ? help   q quit` | 58 |
| Hand-added copy (no retry) | `↑↓ scroll   s status   esc back   ? help   q quit` | 49 |
| Not fetched yet | `↑↓ scroll   r sync   s status   esc back   ? help   q quit` | 58 |
| Narrow 40 | `↑↓ scroll  esc back  ? help` | 27 |
| Tiny 32 | `esc back  ? help  q quit` | 24 |
| Host B | **`Z-04`'s footer** | — |

Separator 3 spaces, tightening to 2 before any hint drops; `? help` and `q quit` are last.

---

## 15 · Colour budget declaration

| State | STATE cyan (uncounted) | Focus ring (exempt) | **CHROME CYAN** | Verdict |
|---|---|---|---|---|
| D2 normal, `ZERADO` | 1 chip — `◉` + `ZERADO` | none — this screen has no focusable control | **0** | **PASS** |
| D2 normal, any other state | 0 | none | **0** | **PASS** |
| D4 hand-added | 0 | none | **0** | **PASS** |
| D5 not fetched | 0 | none | **0** | **PASS** |
| D7 offline | as D2 | none | **0** | **PASS** |
| **D14 absent** | as D2 — the chip is unchanged | none | **0** | **PASS** |
| **D15–D17 the `COVER` row** | 0 — a readout row has no chip | none | **0** | **PASS** |
| Host B, pane focused | as D2 | the pane's border is `--z-border-strong`, **not cyan** | **0** | **PASS** |

**`Z-05` spends ZERO chrome cyan in every state.** A detail view has nothing to urge — the
player already chose this game. Zero is a pass (`02-colour-budget.md` §2.2), and it is the right
answer here rather than a missed opportunity.

**Amber allow-list entries used:** 1 the screen title (host A) / the focused pane title
(host B) · 3 the `IN PROGRESS` state · 7 the action-required degrade banner.
**Not used:** 2 readout labels · 4 progress fill · 5 key hints beyond the footer · 6 the terminal
mark · 8 the filter sigil.

> **Rev B struck entry 2 from the *used* list, and corrected the count below it.** The
> `02-colour-budget.md` §4.1 allow-list says where amber **may** appear, not where it must — and
> §5 and §7 of this spec both put readout labels in `--z-text-secondary`, *"the cockpit
> annotation register, present but subordinate."* Rev A's §15 counted them as amber anyway,
> which contradicted its own token table twice over. **The verdict does not move** — it was a
> PASS by a wide margin and is now a PASS by a wider one — but a number that disagrees with the
> token table is a number a builder will implement.

**Amber ceiling:** at 80 × 24 = 1920 cells, ceiling **192**. §3's render spends **11** —
`GAME DETAIL`, and nothing else. The eight readout labels, `COVER` among them, are
`--z-text-secondary`; the chip is cyan; the footer is chrome. **The `COVER` row costs zero amber
and zero cyan**, which is what a readout row should cost.

**Red: none, in any state.** No scanner — there is no indeterminate wait on this screen. No
destructive confirmation — nothing here deletes anything, and marking `ABANDONED` is reversible
and must never raise one. No error text — a failure is a *banner*, and `OFFLINE` is chrome.

**`--z-border` appears exactly once**, as the unfocused pane border, where it bounds **no
control**. The focused pane border — which is a focus indicator and therefore does carry meaning
— is `--z-border-strong` at **4.08**, satisfying WCAG 1.4.11.

---

## 16 · Reuse verdict, per element

| Element | Verdict | Why |
|---|---|---|
| **The detail view itself** | **Build once, mount twice.** Nothing in it may assume a border, a pane width or a surrounding route | The §0 host contract |
| Scroll body | **`bubbles/viewport` — direct fit** | Content overflows at Narrow, at Tiny and on a long title |
| Border + layout | `lipgloss` | Border, join, alignment maths |
| State chip | The shared chip from `01-design-system.md` §3 — **not a second implementation** | A `ZERADO` chip must look and read identically everywhere (WCAG 3.2.4) |
| Readout block | Build fresh · `lipgloss` + the width-aware pad | Two labelled columns; no `bubbles` primitive fits |
| Degrade banner | The shared banner — one `lipgloss` row | |
| Title wrapping | The width-aware wrap, at word boundaries | Never `len()`, never rune count |
| *Sinopse* (Phase 2) | `glamour`, **restyled to Zerado tokens** | Its default theme carries its own palette and would break the colour budget on sight |
| Error rendering | **Not `charmbracelet/log`** | Developer logging role, its own level colours |

---

## 17 · Upstream findings

**Re-checked against head on 2026-08-25.** Three of rev A's five findings are **closed** and are
struck from this table. Finding 2 — the bordered-surface inset drawn but never named — landed as
spine tokens `BorderInsetX = 2` / `BorderInsetY = 0` (`02-composition.md`;
`14-contradictions-closed.md` #16). Finding 5 — the audio verdict — is struck through and marked
SUPERSEDED in `03-designer-manual.md` §5.11 (#4). **Finding 4 is closed by decision, not by
correction, and it is the reason this revision exists:** it read *"no document says what happens
to a detail view whose game a sync removed"* and refused to invent an answer.
`06-data-seams.md` §2.4 now decides it, names the refusal as the right call, and this spec
carries the resulting state row (**D14**), mockup (§9.6) and copy (§10.8).

| # | Finding | Where | Owner |
|---|---|---|---|
| 1 | **`01-design-system.md` §6 still describes a detail pane the spine does not compose, in two different ways.** §6.2 is headed *"Anatomy — Wide tier, pane 28 cols"*, a composition §6.1 and `02-composition.md` §2.1 both rule out below 120 columns; and §6.1's own ExtraWide split reads `ledger 64 · gutter 2 · pane 46` where `02-composition.md` §2.3 binds **66 ∥ 2 ∥ 44**. Both sum to 112, which is the harder kind of disagreement — a builder following one and a reviewer following the other would each think the screen correct. **This spec follows the spine: pane 44, content 38.** `14-contradictions-closed.md` #15 already records this as closed with *"real pane is 44 wide, 38 content"*, so the design system is the document that has not caught up | `01-design-system.md` §6.1, §6.2 | `fft-design-architect` |
| 2 | **`03-responsive.md` §3 specifies Z-05's readout arrangement by tier**, which a two-host view cannot use — the pane's 38 content columns fall *between* Narrow's 36 and Standard's 54. Resolved as the width threshold **D-05-2**, which agrees with the spine at the tested breakpoint | `03-responsive.md` §3 | `fft-tui-architect` |
| 3 | **`Z-15` §9.1 asks for `COVER  cached, 2 days ago` "beside the other provenance", which reads as block 2.** It is in **block 3** instead, and the reason is `Z-15`'s own: the row was routed here because a cover is a **network-derived value that needs its age** (`07-offline-contract.md` §4), and block 3 is the age block. Block 2 is `SET BY` / `<PROVIDER> SAYS` and **D-05-3** makes it the block that renders *Clear override* predictable — a picture inside that argument would break it. Recorded because `Z-15` §17 finding 5 says *provenance block* and a reviewer holding both specs will look in the wrong one | `Z-15-cover-deck.md` §9.1 · §17 finding 5 | `fft-tui-designer` — **`Z-15` next pass**, wording only |
| 4 | **`Z-15` §17 finding 5 also says this spec "currently states *no cover-art placeholder*", and it still does — correctly.** §10.9's copy note and **D-05-6** forbid a *placeholder*: an empty box, a `[cover]` label, a "coming in Phase 2" line. The `COVER` row is **not** a placeholder — it is a readout with a value in every state, and in the one state where it would have nothing true to say it is **absent** (D-05-8). The two rules do not collide and D-05-6 stands unamended; saying so here stops someone striking it | this spec §10.9 · D-05-6 | **closed in this pass** |
| 5 | **Whether the detail pane ever *draws* a cover is undecided, and is not this spec's to decide.** `Z-15` §4 keeps the pane at 44 × 32 with no image region in it, and this spec adds a **row**, not a picture. An image region inside a bordered surface, in a view mounted in two hosts of different content widths (38 and 74), is a composition change of the same class as **D-15-7**. Flagged rather than proposed | `02-composition.md` §2.1 · §2.3 · `Z-15-cover-deck.md` §4 | `fft-tui-architect` |

---

## 18 · Open for the founder

> **Two items closed since rev A, and neither needs the founder now.**
>
> **Naming the bordered-surface inset** (item 1) landed: `02-composition.md` carries
> `BorderInsetX` = **2** columns each side and `BorderInsetY` = **0** rows as spine tokens
> (`14-contradictions-closed.md` #16). D-06-1 is no longer a number read off a mockup.
>
> **What happens to a game a sync no longer returns** (item 2) was routed to
> `fft-tui-architect` and **decided** in `06-data-seams.md` §2.4: tombstoned, never deleted;
> `absent_since` set on the first complete, successful sync that omits it and cleared silently
> when it returns; excluded from `Z-04`'s default view; findable through `Z-07`'s facet; deleted
> only when the player asks. The copy this spec refused to write before the decision is now
> written — **§10.8**, with D14 and §9.6. *The refusal was the right call and the seam says so;
> the sentence it produced is better than the one rev A sketched, because `06-data-seams.md` §2.4 decided the
> reason and not just the behaviour.*

1. **The nine underived ANSI-256 indices.** This screen leans on `--z-text-tertiary` more than
   any other in the bundle — it is the colour of *not knowing* — and ships it **uncoloured** as
   the documented interim. Confirm the derivation lands before `Z-05` is built.

---

## 19 · Design decisions made in this spec

| # | Decision | Reason |
|---|---|---|
| **D-05-1** | The header band names the **screen**, not the game; the game's title is body row 1 in both hosts (§1) | A host-agnostic view cannot own a band it does not have in host B, and a hierarchy that inverts between hosts is not a hierarchy |
| **D-05-2** | Readout arrangement switches at the named threshold `ReadoutBesideMin = 38` content columns (§9.4) | The pane's width falls between two tiers, so a tier-based rule is not expressible here; the threshold agrees with `03-responsive.md` at the tested breakpoint |
| **D-05-3** | The status-provenance block (`SET BY` / `<PROVIDER> SAYS`) is a permanent block, present in every state, with the provider row omitted when no override exists (§3.2, §10.5) | It is the fact that makes `Z-06`'s *Clear override* predictable, and `05-state-machine.md` §5 requires that the consequence be nameable in advance |
| **D-05-4** | Block 3 is bottom-anchored when a blank row would otherwise separate it (§3.2) | Turns 17 rows of empty box into a readout footer in host B and a status line in host A, with one rule for both |
| **D-05-5** | Five distinct renderings of *not a value*, carried in **words**: `9h` · `0h` / `never played` · `—` · `—` + banner · `not tracked` (§9.3) | `01-screen-inventory.md` §5 requires the distinction and the player can act on only two of the five; a colour-only difference would not survive `NO_COLOR`, a screenshot, or a screen reader |
| **D-05-6** | No Phase 2 field labels, no "coming in Phase 2" line, no cover placeholder (§10.9) | Anti-pattern 14 — *omit the block, do not label it empty.* The roadmap is on the landing page, not in a detail view |
| **D-05-7** | An absent row is announced on the **`SOURCE` value** plus a two-line prose block — **not** a banner (§9.6, §10.8) | A degrade banner's four mandatory parts (`07-offline-contract.md` §3) all fail here: nothing is unavailable, nothing is stale, and there is no retry. *The source stopped returning it* is a fact about the source, so it lands on the row that names the source — at **zero row cost**, which is what lets block 3 keep its respiro on the screen where the age matters most |
| **D-05-8** | **The `COVER` row is block 3's *first* row, and it renders only when the source is live** (§3.2, §9.7). Present on a syncable row that the provider is still returning (D2, D3, D5, D6, D7, D10, D15–D17); **absent** on a hand-added copy (D4) and on an absent row (D14) | Three reasons, and they agree. **(1)** A cover is a network-derived value and `07-offline-contract.md` §4 requires its age — an age `LAST SYNCED` cannot give, because the row and the picture are fetched by different runs. **(2)** Where there is no live source, the `SOURCE` row already carries the fact and carries it better: `Added by hand` says no store will ever have art, and `absent since <date>` is an **outer** age that bounds the picture's along with everything else — a second, narrower age inside it is furniture. **(3)** D14 has no spare row: rev A used all 16 and the row that would fall below the fold is `LAST SYNCED`, the mandatory one. **`SOURCE` owns a cover's provenance when there is no live source; `COVER` owns its age when there is** |
| **D-06-1** | A bordered surface is inset 2 columns each side and 0 rows, fixed at every tier (§4) | Read from `01-design-system.md` §6.2's ratified anatomy; **now a spine token** — `02-composition.md` names `BorderInsetX` = 2 and `BorderInsetY` = 0 |

---

## 20 · Screen-specific acceptance criteria

Beyond `00-design-brief.md` §10 and `02-colour-budget.md` §10, both of which apply in full.

1. **The same binary renders both hosts from one view.** A reviewer diffs the block order, the
   label field, the copy and the state renderings between the 80 × 24 route and the 120 × 40
   pane: **they are identical**. Any `if host == route` inside the view is a fail.
2. **The view never draws a header band, never draws a footer in host B, and never handles
   `Esc` itself.**
3. **The game's title is never truncated in either host** — it wraps at a word boundary, at 74,
   at 38, at 36 and at 30 columns.
4. **`—` and `not tracked` are visually and textually distinct**, and only the `—` state offers
   `r` in the footer. Verified on a hand-added copy and on a not-yet-fetched Steam row, side by
   side.
5. **Block 3 is absent on a hand-added copy** and present on every syncable one. `LAST SYNCED`
   never reads `never` for a provider that cannot sync.
5b. **The `COVER` row says exactly what the tile says, for the same game, in every one of
   `Z-15`'s seven picture states** (§9.7). Render the deck and the detail view side by side and
   assert the strings match — `not fetched`, `nothing found` — with two declared exceptions:
   a cached cover, where the tile shows the picture and this row shows its **age**; and a
   hand-added copy, where the tile says `added by hand` and this screen has **no `COVER` row**.
   **No third vocabulary anywhere.**
5c. **`r sync` is in the footer if and only if a retry exists** — now on two conditions: a `—`
   value, or `COVER  not fetched`. Verified on a complete row whose enrichment failed
   (`LAST SYNCED  3 hours ago` **and** `COVER  not fetched`, `r` present) and on `nothing found`
   (`r` absent unless something else needs it).
5d. **D14 and D4 render no `COVER` row**, in either host, at every tier — and D14 still closes
   in 16 body rows without scrolling, in both the override and the no-override form
   (criterion 14).
5e. **The `COVER` row is byte-identical on a terminal with no image support.** Run with
   `ZERADO_NO_IMAGES=1` and diff this screen against the same game on a terminal that draws
   images: **zero differences.** A cover's provenance is text, and text does not need a
   graphics protocol.
6. **`⏎` on an empty ledger does not push this screen** and shows no error.
7. **The chip renders before the values do**, on a cold open — title and state come from the
   ledger row already in memory.
8. **Chrome-cyan count is 0 in every state**, by `02-colour-budget.md` §3.1.
9. **`NO_COLOR=1` loses no information**, including the not-fetched / not-trackable distinction.
10. **At the 16-colour floor the pane is still a distinct region** — its border and the
    2-column gutter carry it; no fill is involved.
11. **The focused pane is identifiable with colour off**, by box weight alone.
12. **No scanner and no progress bar render on this screen in any state.**
13. **An absent row renders its four block-1 values, its chip and its block 2 exactly as a
    present row of the same state does** — the only differences are the `SOURCE` clause and the
    prose block. Verified by diffing a `ZERADO` row against the same row tombstoned.
14. **D14 does not scroll at 80 × 24**, with an override and without one — the 16-row
    arithmetic of §9.6 closes in both forms, and `LAST SYNCED` is on screen in both.
15. **No banner, no red, no retry key and no fifth state glyph appear on an absent row**, in
    either host, at any tier.
16. **The absent prose is never truncated** — at 74, 38, 36 and 30 content columns it re-wraps.
17. **Clearing `absent_since` returns the screen to D2 with no announcement** — no banner, no
    result line, no dwell.
18. **Founder-validated screenshot before merge**, at all six viewports of
    `03-responsive.md` §7 plus `NO_COLOR=1` and forced-16-colour at 80 × 24 — and in **both
    hosts** at 120 × 40. No screenshot → not GOLDEN → no merge.
