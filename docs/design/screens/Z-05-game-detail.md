---
title: Zerado — Z-05 Game detail
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-05
rev: A
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

**This is the normal Phase 1 state, not a degraded one.** Phase 1 has no cover art, no *sinopse*
and no mood tags, so this composition — five facts about the game, two about its status, one
about the data's age — is the screen, designed as such.

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
│                                                                                │
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
| 18–21 | 12–15 | blank |
| 22 | 16 | **block 3 — the age of the data**, bottom-anchored |
| 23 | — | footer, the reserved frame row |
| 24 | — | blank · `OuterMarginY` = 1 |

### 3.2 · The three blocks, and why there are exactly three

| Block | Keys | The question it answers |
|---|---|---|
| **1 · the game** | `PLAYTIME` · `LAST PLAYED` · `ADDED` · `SOURCE` | *What is this thing?* |
| **2 · the status** | `SET BY` · `<PROVIDER> SAYS` | *Why does it say `ZERADO`, and what would happen if I cleared that?* |
| **3 · the age** | `LAST SYNCED` | *How old is what I am reading?* — `07-offline-contract.md` §4, and it is mandatory |

Block 2 is the reason `Z-06`'s fifth item is comprehensible: it names, in advance, exactly what
*Clear override* will restore. Block 3 is the age rule, which is *"the rule most likely to be
lost during a build, because dropping the age always makes the layout tidier."*

> **D-05-4 · Block 3 is bottom-anchored** to the view's last row when at least one blank row
> would otherwise separate it from the flowed content; otherwise it flows immediately after
> block 2 with one `InterElementGap`. One rule, both hosts. In host B it turns 17 rows of empty
> box into a readout footer; in host A it reads as a status line, which is what it is.

---

## 4 · Mockup — host B, the pane at 120 × 40

44 × 32 (border 1 + inset 2 + **content 38** + inset 2 + border 1), sitting in `Z-04`'s
`112 × 32` body after a 2-column gutter. **Identical content, identical block order, identical
label field. Only the rectangle changed.**

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
│                                          │
│                                          │
│  LAST SYNCED  3 hours ago                │
└──────────────────────────────────────────┘
```

**When the pane holds focus** — `Tab` or `⏎` from the ledger. The border goes **heavy**; the
pane title goes amber. Under `NO_COLOR` the box weight alone still says which region has focus
(`04-navigation-and-focus.md` §4.2).

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
> resized rather than designed. **Proposed as a named token `BorderInsetX = 2` — see §18.**

> **The pane is 30 content rows and Phase 1 fills 12 of them.** That is not an unfinished
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

Ratios read from the brand manual's measured table (§4.2). None estimated.

| Element | Token | Hex | ANSI-256 | 16-colour | Ratio |
|---|---|---|---|---|---|
| The game's title | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** AAA |
| Readout values | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** AAA |
| Readout labels (UPPER) | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA |
| Explanatory prose | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA |
| `—` · `not tracked` · `never` | `--z-text-tertiary` | `#8492A8` | ***underived*** | `white` | **6.15** AA |
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

**Block 3 is absent.** A provider that is not a `Syncer` has no last-sync time, and labelling one
`never` would invite a retry that cannot exist. *Omit the block; do not label it empty.*
**`r` is not in the footer** on this state, for the same reason.

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
│                                                                                │
│   LAST SYNCED  never                                                           │
│   ↑↓ scroll   r sync   s status   esc back   ? help   q quit                   │
│                                                                                │
└────────────────────────────────────────────────────────────────────────────────┘
```

Banner **64** cells. **`r sync` is in the footer here and is not in §9.1's** — that single
difference is the whole point of the distinction: the footer offers the retry exactly when a
retry exists.

### 9.3 · The five renderings of *not a value* — the table this screen exists for

| Fact state | Value renders | Colour | Block 3 | Retry offered? | Why |
|---|---|---|---|---|---|
| **Known** | `9h` · `2 Aug 2026` | `--z-text` | present | — | |
| **Known to be empty** | `0h` · `never played` | `--z-text-secondary` | present | **no** | The provider answered. There is nothing to fetch |
| **Not fetched yet** | **`—`** | `--z-text-tertiary` | `never` | **yes — `r`** | Nobody has asked. This is the one the player can act on |
| **Unavailable right now** | `—` and the banner names it | `--z-text-tertiary` | the age | **yes — `r`** | The last value, with its age, or nothing plus the reason |
| **Not trackable at all** | **`not tracked`** | `--z-text-tertiary` | **absent** | **no** | The provider has no such capability. Not a fault, not a gap |

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
│                                │
│ esc back  ? help  q quit       │
└────────────────────────────────┘
```

Body **30 × 21**. `OuterMarginX` and `OuterMarginY` shed to **0**; the band is the title row
alone. The chip does **not** shrink and the label is **never** dropped.

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

LAST SYNCED  never
```

### 10.4 · D6 — known to be empty

```
PLAYTIME     0h
LAST PLAYED  never played
```

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
```

The **title and the chip render immediately** — both came from the ledger row the player was
already looking at. No spinner, no scanner: reading a local row is not an indeterminate wait.

### 10.7 · D10 — a pending write

```
◉  ZERADO…
```

The label carries the `…` until SQLite confirms. **Never an optimistic silent change.**

### 10.8 · Copy notes

- **Dates** are said the way people say them — `12 Mar 2026`, `3 hours ago`, `just now`,
  `last June`. Never `2026-08-22T04:11:09Z`. The exact timestamp belongs in this view for anyone
  who wants it, and Phase 1's answer is that `ADDED` and `LAST PLAYED` are already the precise
  dates; only `LAST SYNCED` is relative, because relative is what "how stale is this" means.
- **Say the number** — `9h`, `41 hours`, `147 of 247`.
- **No exclamation marks. No emoji. The user is never a "gamer".**
- **Casing** — `Zerado` the product in prose and in the breadcrumb; `ZERADO` the chip; `zerado`
  never appears on this screen because there is no summary sentence here.
- **Nothing on this screen claims an unbuilt capability.** There is no `MOOD` label, no
  `SINOPSE` label, no cover-art placeholder and no "coming in Phase 2" line. The roadmap lives
  on the landing page.
- **Type-neutral where equally natural** — `DETAIL`, `SOURCE`, `ADDED`, `Nothing is tracked`.
  Phase 1 says *game* where a game is what the player is looking at.

---

## 11 · 40-column behaviour, and the refusal floor

**40 × 24** is §9.4: body `36 × 16`, readout **stacked** (D-05-2, `36 < 38`), title wrapped, the
chip unchanged at its full 14 columns, blocks 2 and 3 reachable by scrolling.

**Sheds at Narrow and below:** the beside arrangement (the label moves above the value).
**Never sheds, at any width:** the state — glyph **and** label · the game's title · the
degrade banner when one is active · the footer key line.

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
│                                                                                │
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
| Which region has focus (host B) | `┏━┓` against `┌─┐` — box-drawing **weight** |
| The banner class | The label word |
| Screen title | UPPERCASE + bold, alone on its row |

**Under `NO_COLOR` this screen loses nothing**, which is the §3.3 cross-check: run it with
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
| `r` | Re-sync → `Z-03` | ✓ | ✓ | **Listed only when a retry exists** — never on a hand-added copy (§9.1) |
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
| Host B, pane focused | as D2 | the pane's border is `--z-border-strong`, **not cyan** | **0** | **PASS** |

**`Z-05` spends ZERO chrome cyan in every state.** A detail view has nothing to urge — the
player already chose this game. Zero is a pass (`02-colour-budget.md` §2.2), and it is the right
answer here rather than a missed opportunity.

**Amber allow-list entries used:** 1 the screen title (host A) / the focused pane title
(host B) · 2 readout labels · 3 the `IN PROGRESS` state · 7 the action-required degrade banner.
**Not used:** 4 progress fill · 5 key hints beyond the footer · 6 the terminal mark · 8 the
filter sigil.

**Amber ceiling:** at 80 × 24 = 1920 cells, ceiling **192**. §3's render spends 11
(`GAME DETAIL`) + 55 (readout labels) ≈ **66**.

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

| # | Finding | Where | Owner |
|---|---|---|---|
| 1 | **`01-design-system.md` §6.2 draws the detail pane at "Wide tier, pane 28 cols"**, but §6.1 and `02-composition.md` §2.1 both establish that **no detail pane exists below 120 columns**. The 28-column anatomy describes a composition the spine rejected; the real pane is **44** wide with **38** of content | `01-design-system.md` §6.2 | `fft-design-architect` |
| 2 | **The bordered-surface inset is drawn but never named.** §6.2's own anatomy uses 2 columns and 0 rows, which does not equal `InnerPaddingX` at any tier except ExtraWide. This spec adopts the drawn value as **D-06-1** and proposes naming it | `01-design-system.md` §6.2 | `fft-design-architect` |
| 3 | **`03-responsive.md` §3 specifies Z-05's readout arrangement by tier**, which a two-host view cannot use — the pane's 38 content columns fall *between* Narrow's 36 and Standard's 54. Resolved as the width threshold **D-05-2**, which agrees with the spine at the tested breakpoint | `03-responsive.md` §3 | `fft-tui-architect` |
| 4 | **No document says what happens to a detail view whose game a sync removed.** `06-data-seams.md` §2.4 says a sync streams and a cancel leaves *"a valid partial library"*, but never states whether a row absent from a provider's response is deleted, tombstoned or kept. This spec does **not** invent an answer — see §18 item 2 | `06-data-seams.md` | `fft-tui-architect` |
| 5 | `03-designer-manual.md` §5.11 verdict 3 still reads as a permanent rejection of the audio subsystem, superseded by founder direction relayed 2026-08-25 | `03-designer-manual.md` §5.11 | `fft-design-architect` |

---

## 18 · Open for the founder

1. **Name the bordered-surface inset as a token.** D-06-1 fixes it at 2 columns / 0 rows on the
   authority of `01-design-system.md` §6.2's own drawing, and both the detail pane and `Z-06`'s
   overlay depend on it. It is currently a number read off a mockup. **Proposal:
   `space.BorderInsetX = 2`, fixed at every tier**, added to the token table so no screen ever
   types a 2.
2. **What happens to a game a sync no longer returns?** Kept with its status, tombstoned, or
   deleted? It decides real copy on this screen — *"This is no longer in your Steam library.
   Zerado kept the row and your status."* is a sentence I can write the moment the behaviour is
   decided, and it is a sentence I refuse to write before then, because the product's promise is
   about a file the player owns. **This spec deliberately has no such state row.** Route to
   `fft-tui-architect` (`06-data-seams.md`).
3. **The nine underived ANSI-256 indices.** This screen leans on `--z-text-tertiary` more than
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
| **D-05-6** | No Phase 2 field labels, no "coming in Phase 2" line, no cover placeholder (§10.8) | Anti-pattern 14 — *omit the block, do not label it empty.* The roadmap is on the landing page, not in a detail view |
| **D-06-1** | A bordered surface is inset 2 columns each side and 0 rows, fixed at every tier (§4) | Read from `01-design-system.md` §6.2's ratified anatomy; a box that breathes differently at different widths looks resized rather than designed |

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
6. **`⏎` on an empty ledger does not push this screen** and shows no error.
7. **The chip renders before the values do**, on a cold open — title and state come from the
   ledger row already in memory.
8. **Chrome-cyan count is 0 in every state**, by `02-colour-budget.md` §3.1.
9. **`NO_COLOR=1` loses no information**, including the not-fetched / not-trackable distinction.
10. **At the 16-colour floor the pane is still a distinct region** — its border and the
    2-column gutter carry it; no fill is involved.
11. **The focused pane is identifiable with colour off**, by box weight alone.
12. **No scanner and no progress bar render on this screen in any state.**
13. **Founder-validated screenshot before merge**, at all six viewports of
    `03-responsive.md` §7 plus `NO_COLOR=1` and forced-16-colour at 80 × 24 — and in **both
    hosts** at 120 × 40. No screenshot → not GOLDEN → no merge.
