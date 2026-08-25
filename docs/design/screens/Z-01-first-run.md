---
title: Zerado — Z-01 First run
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-01
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-01 · First run

> Fills the 16-section contract in [`../03-designer-manual.md`](../03-designer-manual.md) §3.
> Canon travels with it: brand manual rev A · `naming.md` · `tokens.css`/`tokens.json` ·
> `ratification/decisions.md` · `content/landing-copy.md` · [`../00-design-brief.md`](../00-design-brief.md) ·
> [`../01-design-system.md`](../01-design-system.md) · [`../02-colour-budget.md`](../02-colour-budget.md) ·
> FlowForge TUI Design Manual #2371 · Spacing Canon #2435.
> Composition is binding from [`../../blueprint/02-composition.md`](../../blueprint/02-composition.md) §2 —
> single-pane, content block **top-aligned**, three doors, `R = 1`.

---

## 1 · Identity

| | |
|---|---|
| **Screen** | `Z-01` · First run |
| **Phase** | 1 |
| **Kind** | **Route** — pushed as the **root** instead of `Z-04` when the library is empty *and* no provider has ever been connected ([`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §2.1) |
| **Route in** | Start-up only. The condition is checked **once**; it is never re-checked |
| **Routes out** | `Z-02 Connect a store` · `Z-08 Add a game by hand` · `Z-04 Library` · `Z-09 Settings` (`,`) · `Z-10 Help` (`?`) |
| **Stack position** | Bottom. Nothing is below it, so `Esc` has nothing to pop |
| **Offline class** | **WORKS** ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2) — two of the three doors are fully functional with the network off |

---

## 2 · Purpose

**Say what this program is, and open three doors — one of which works on a plane.**

`Z-01` **collects nothing.** It has no field, no prompt, no question. Screen inventory §5:
*"`Z-01` must not ask for anything. It offers doors; it does not collect."*

---

## 3 · Mockup at 80 columns — the design floor

**Terminal `80 × 24` · Wide · `leftInset` 3 · body `74 × 16` · content begins at column 4.**
Every row below is **exactly 80 columns**; trailing blanks are not drawn. The ruler's digit marks
columns 1, 11, 21 … 71.

### 3.1 · The frame row map (all framed Zerado screens, `80 × 24`)

Derived from [`02-composition.md`](../../blueprint/02-composition.md) §1.3; not re-derived here.

| Row | Owner | Token that buys it |
|---|---|---|
| 1 | blank | `OuterMarginY` = 1 (top) |
| 2 | breadcrumb | `HeaderBandHeight` = 3, row 1 of 3 |
| 3 | blank | `InterElementGap` = 1, inside the band |
| 4 | title | `HeaderBandHeight` row 3 of 3 |
| 5 | blank | `InterElementGap` = 1, band → body |
| 6–21 | **body — 16 rows** | `BodyRect.h` |
| 22 | blank | `InnerPaddingY` = 1 (bottom) |
| 23 | footer | the Canon's reserved footer row |
| 24 | blank | `OuterMarginY` = 1 (bottom) |

`1 + 3 + 1 + 16 + 1 + 1 + 1 = 24` ✓ · `80 − 2 × 3 = 74` ✓

> **Reading recorded.** `BodyRect.h` charges `InnerPaddingY` **once**, not twice. The top inner
> padding is absorbed inside the band's own three rows; the single charged row is the one between
> the body and the footer. Stated so the mockups add up. If the spine intends otherwise, this is
> a one-line correction there, not here.

### 3.2 · RENDER 80×24 — default (doors all open)

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ First run

   FIRST RUN

   Zerado keeps your library in one file on this machine.
   No account. No server. The only traffic is the stores you connect.

   ▌ Connect a store         Steam today. PlayStation, GOG and EA later

     Add a game by hand      Discs, cartridges, anything the stores miss

     Look around first       An empty library, and every key that works

   Radio and interface sounds are off until you turn them on in Settings.







   ↑↓ move   ⏎ choose   , settings   ? help   q quit

```

### 3.3 · RENDER 80×24 — no network (the third door is closed, not hidden)

Refines the drawn mockup in [`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §6.
**The door keeps its position**; only the focus and the reason change. See §5.4 for why.

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ First run

   FIRST RUN

   Zerado keeps your library in one file on this machine.
   No account. No server. The only traffic is the stores you connect.

     Connect a store         OFFLINE   No network. r to check again.

   ▌ Add a game by hand      Discs, cartridges, anything the stores miss

     Look around first       An empty library, and every key that works

   Radio and interface sounds are off until you turn them on in Settings.







   ↑↓ move   ⏎ choose   r check again   , settings   ? help   q quit

```

### 3.4 · The door row — column budget at 74

| Field | Cols | Range (80-col terminal) | Note |
|---|---|---|---|
| focus field | **2** | 4–5 | fixed, width-aware padded — `▌` U+258C is **Ambiguous** ([`01-design-system.md`](../01-design-system.md) §1.2) |
| gap | 1 | 6 | |
| door title | **22** | 7–28 | longest is `Add a game by hand` = 18 |
| gap | 2 | 29–30 | |
| reason | **47** | 31–77 | longest drawn is 43 |
| | **74** | | |

---

## 4 · Mockup at the primary breakpoint — `120 × 40`

**Terminal `120 × 40` · ExtraWide · `leftInset` 4 · body `112 × 32` · content begins at column 5.**
Row map: 1 blank · 2–4 band · 5 blank · **6–37 body (32)** · 38 blank · 39 footer · 40 blank
= `1 + 3 + 1 + 32 + 1 + 1 + 1 = 40` ✓ · `120 − 2 × 4 = 112` ✓

**The block does not stretch to 112.** Prose is capped at the brand's measure, `--z-measure` =
**68 characters** (brand §6.3), and the door block keeps the Wide field budget (`2 / 1 / 22 / 2 / 47`
= 74). *Design decision:* 43 characters of text followed by 69 columns of air is a stretched
layout, not a breathing one; the measure cap exists for exactly this.

```text
0.........1.........2.........3.........4.........5.........6.........7.........8.........9.........A.........B.........

    Zerado ✦ First run

    FIRST RUN

    Zerado keeps your library in one file on this machine.
    No account. No server. The only traffic is the stores you connect.

    ▌ Connect a store         Steam today. PlayStation, GOG and EA later

      Add a game by hand      Discs, cartridges, anything the stores miss

      Look around first       An empty library, and every key that works

    Radio and interface sounds are off until you turn them on in Settings.
```

*(rows 20–37 blank — the block is **top-aligned**, per [`02-composition.md`](../../blueprint/02-composition.md) §2; ruler digits mark columns 1, 11 … 111, `A` = 101, `B` = 111.)*

---

## 5 · Visual hierarchy

**The one thing the player must see first: the focused door.**

| Rank | Element | Channel carrying it | Why it ranks here |
|---|---|---|---|
| 1 | `FIRST RUN` | **case** (UPPER) + **weight** (bold) + **colour role** (`--z-primary` amber) + **position** (top-left, own row) | The display role, §1.5. It is first because terminal reading order is top-left; it is *earned* by being the only amber block on screen |
| 2 | **the focused door** | **position** (`▌` in the 2-col gutter) + **weight** (bold row) + **colour** (amber marker) | Three channels, any two sufficient (§1.7). This is where the eye is *sent* after the title |
| 3 | the two promise lines | **spacing** (their own block, one `InterElementGap` from the doors) + `--z-text` | Two sentences, then out of the way |
| 4 | the other two doors + the reason column | **case** (sentence) + `--z-text` title / `--z-text-secondary` reason | Reachable, not shouting |
| 5 | the audio line | `--z-text-secondary`, one gap below the doors | Deliberately quiet — it is a note, not a prompt |
| 6 | breadcrumb, footer | `--z-text-secondary`, edges | Chrome |

**Hierarchy comes from case, weight, colour role, box drawing and spacing — in that order**
(§1.1). Z-01 uses **no box drawing at all**: it has one region, and a border around the only
thing on screen is a frame around a frame.

### 5.4 · Why the closed door keeps its position — *design decision*

The drawn mockup in `07-offline-contract.md` §6 puts `Add a game by hand` first when the network
is down. This spec **keeps the ratified door order** (`Connect a store` · `Add a game by hand` ·
`Look around first`, the order in [`01-screen-inventory.md`](../../blueprint/01-screen-inventory.md) §2)
in both states and **moves the focus instead**.

- The load-bearing insight of the §6 mockup is *where the focus sits*, and that is kept exactly.
- A list that reorders itself between two renders reads as two different screens.
- The closed door is the **most informative** row for a player who came here to connect Steam.
  Demoting it buries the explanation.

---

## 6 · Every applied spacing token, by name

No magic numbers. Values from [`00-design-brief.md`](../00-design-brief.md) §5.2 = Spacing Canon #2435 §4.

| Token | Tiny | Narrow | Standard | **Wide** | ExtraWide | Where Z-01 spends it |
|---|---|---|---|---|---|---|
| `OuterMarginX` | 0 | 1 | 2 | **2** | 2 | left/right inset of the whole frame |
| `OuterMarginY` | 0 | 1 | 1 | **1** | 1 | blank rows 1 and 24 |
| `InnerPaddingX` | 1 | 1 | 1 | **1** | 2 | inside the frame |
| `InnerPaddingY` | 0 | 1 | 1 | **1** | 1 | blank row 22 (body → footer) |
| `InterElementGap` | 1 | 1 | 1 | **1** | 1 | breadcrumb→title · band→body · prose→doors · **between doors** · doors→audio line |
| `HeaderBandHeight` | **1** | 3 | 3 | **3** | 3 | `hasSubtitle = false`, always (§2.4) |
| `HeaderBand(tier, false)` | 1 | 3 | 3 | **3** | 3 | the effective band |
| **`leftInset`** | **1** | **2** | **3** | **3** | **4** | header-left **==** content-left, verified at column 4 |
| `--z-measure` (brand §6.3) | — | — | — | — | **68** | prose cap at ExtraWide |

**`InterElementGap` is spent once between doors** — *design decision.* Each door is a
destination with its own consequence line; run together they read as a settings menu, spaced they
read as three doors. There is room: `74 × 16` with three doors is not tight.

---

## 7 · Colour, glyph and label for every state shown

Every value read from [`01-design-system.md`](../01-design-system.md) §1.4, which reads the brand's
measured table. Nothing estimated.

### 7.1 · The door states — the co-render map

| Door state | Token | Hex | ANSI-256 | 16-colour | Ratio | Glyph / structural mark | Label channel | Survives `NO_COLOR`? |
|---|---|---|---|---|---|---|---|---|
| **Focused** | `--z-primary` (marker) + bold row | `#FFB000` | **214** | `bright yellow` | **10.59** AAA | `▌` U+258C in the 2-col gutter (ASCII `>`) | the door title | **yes** — marker **and** bold both survive |
| **Open, unfocused** | `--z-text` (title) | `#E9EEF5` | **255** | `bright white` | **16.65** AAA | gutter blank (2 spaces) | title + reason | **yes** |
| **Closed** (no network) | `--z-text-secondary` (title **and** reason) | `#A9B5C7` | **249** | `white` | **9.36** AAA | **none** — the uppercase label word replaces the glyph, per [`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §3 | `OFFLINE` **bold**, then `No network. r to check again.` | **yes** — the word is the state |
| reason text (all doors) | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA | — | — | **yes** |
| breadcrumb `✦` | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** | `✦` U+2726 — **Neutral**, 1 cell | — | **yes** |
| screen title | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA | — | `FIRST RUN` UPPER + bold | **yes** — case + bold |
| audio note | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** | — | the sentence | **yes** |
| **audio indicator — never enabled** *(default)* | — | — | — | — | — | **absent** | **absent**, and `m` absent from the footer | **yes** — nothing to lose |
| **audio indicator — enabled, unmuted** | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** AAA | `▮` U+25AE — **Neutral**, 1 cell | `AUDIO` | **yes** — filled-vs-hollow carries it |
| **audio indicator — enabled, muted** | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** AAA | `▯` U+25AF — **Neutral**, 1 cell | `MUTED` | **yes** |

Audio rows applied from [`01-design-system.md`](../01-design-system.md) §5.2, not re-derived.
**Placement on Z-01:** the screen has no pinned summary row, so the indicator takes the
right-hand slot of the **reserved footer row** — the position §5.1's anatomy draws it in — with a
minimum two-column gap from the key run. It costs **no body row**. When the two would collide the
key run sheds from the right (`? help` and `q quit` last, nav §6) and **the audio glyph is the
last thing dropped, never the first** (§5.5).

```text
   ↑↓ move   ⏎ choose   , settings   ? help   q quit          ▮ AUDIO
```

### 7.2 · Why "closed" carries no extra glyph — *design decision*

`⊘` is taken (`ABANDONED`). `⚠` is rejected by [`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §3.
Every plausible alternative (`·` U+00B7, `×` U+00D7, `—` U+2014) is **East-Asian-Ambiguous** and
would shear a fixed-width field — verified with `unicodedata`, Unicode **16.0.0**, on 2026-08-25.

The offline contract already settles it: *"An uppercase label word — the word carries the state,
so it survives with zero colour and needs no warning glyph."* `OFFLINE` in the reason column is
that word, in the product's own ratified vocabulary, and it is greppable, speakable, and
case-distinct from the two sentence-case reasons beside it.

### 7.3 · The closed door is **not** amber — *design decision*

[`01-design-system.md`](../01-design-system.md) §12.2: *"Amber appears only when the player must do
something."* On `Z-01` the player must **not** do anything — two doors work and the product is
fully usable. The `OFFLINE` word therefore renders in `--z-text-secondary` chrome, **bold**, and
the amber on this screen stays spent on exactly one thing: the focus marker.

### 7.4 · Deliberately not used on Z-01

`--z-text-tertiary` `#8492A8` has **no derived ANSI-256 index** ([`00-design-brief.md`](../00-design-brief.md) §9).
Rather than ship the documented uncoloured interim for the most-read text on the first screen a
player ever sees, **de-emphasis on Z-01 is carried by `--z-text-secondary` (249, derived,
measured 9.36)**. Nothing is invented and nothing renders uncoloured. *This substitution is used
across all six specs in this cluster and is recorded once here.*

---

## 8 · The full state table

| # | State | When | What renders | Copy |
|---|---|---|---|---|
| 1 | **First run** *(this whole screen)* | library empty **and** no provider ever connected | §3.2 | as drawn |
| 2 | **Loading** | **N/A — and stated, not skipped.** Z-01 reads only `Store.Games` count and `Store.Connections`, both local and both already resolved before the route is pushed. There is no async fetch, therefore no spinner and no scanner | — | — |
| 3 | **Empty** | **N/A.** Z-01 *is* the product's empty state ([`01-design-system.md`](../01-design-system.md) §10.1) | — | — |
| 4 | **Partial** | **N/A.** Z-01 has no data set to be partial | — | — |
| 5 | **Error** | **N/A.** Z-01 reads nothing that can fail. If `library.db` is unreadable the program never reaches Z-01 — that is **`Z-11`** | — | — |
| 6 | **Offline — link state reports no route** | see §8.1 | §3.3 — door 1 closed, focus on door 2, footer gains `r check again` | `OFFLINE   No network. r to check again.` |
| 7 | **Offline — link state unknown / unavailable** | the OS gives no answer | **all doors open** (§3.2). The honest default | as §3.2 |
| 8 | **Re-check ran, still no route** | `r` or `⏎` on the closed door | door stays closed; reason line unchanged | `OFFLINE   No network. r to check again.` |
| 9 | **Re-check ran, route returned** | `r` or `⏎` on the closed door | door opens, focus **stays on the door the player was on** | reverts to the open reason |
| 10 | **Audio enabled, unmuted** | `Z-09` → Audio = On | `▮ AUDIO` right-aligned on the **reserved footer row**; footer key run gains `m mute` | — |
| 11 | **Audio enabled, muted** | `m` | `▯ MUTED` right-aligned on the footer row; key run shows `m unmute` | — |
| 12 | **Audio never enabled** *(the default)* | always, until opted in | **no indicator, and `m` is not in the footer** ([`01-design-system.md`](../01-design-system.md) §5.2) | — |
| 13 | **Tiny `< 40`** | | door **titles only**, no prose, no audio line; band collapses to the title row | see §11.3 |
| 14 | **`NO_COLOR`** | env set | §12 | identical text |
| 15 | **Below `24 × 8`** | | Z-01 never renders — see `Z-11` | `Zerado needs at least 24 columns and 8 rows. This terminal is 20 x 6.` |

### 8.1 · How Z-01 knows the door is closed — *design decision, and it needs the founder*

[`07-offline-contract.md`](../../blueprint/07-offline-contract.md) **§5** says Zerado *"does not
guess, and it never probes"* and *"does not know it is offline until you ask it to do something
that needs the network."* **§6** of the same document draws a first-run screen with the Connect
door already disabled **at first paint**. Those two sections cannot both be satisfied by a
request-classification model alone.

**Resolution taken here:** Z-01 reads the **local network-interface / default-route state** — a
purely local fact, obtained from the OS, that **emits no packet**. It is not a reachability probe
and cannot say whether Steam is up; it can only say the machine has no route at all.

| Signal | Door |
|---|---|
| no interface up / no default route | **closed**, with its reason |
| route present | **open** |
| the OS gives no usable answer | **open** — the honest default; the failure then lands on `Z-02` with the ratified refusal copy |

This keeps §5's promise (no traffic, no heartbeat, no guessing) and §6's screen. **Flagged in
§16 as needing the founder's word**, because it reads a system fact §5's sentence did not
anticipate.

---

## 9 · The key map

Global keys from [`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §3.
The **footer** column says whether the key is *listed* — a footer that lies is worse than no footer.

| Key | Does | In the footer? |
|---|---|---|
| `↑` `↓` · `k` `j` | move between doors | **yes** — `↑↓ move` |
| `g` / `G` | first / last door | no — three items; listing it would be noise |
| `⏎` | open the focused door. On the **closed** door it re-reads the link state (§8.1) | **yes** — `⏎ choose` |
| `r` | re-read the link state | **only when a door is closed** — `r check again` |
| `,` | Settings (`Z-09`) | **yes** — `, settings` |
| `m` | toggle mute | **only when audio has been enabled** — `m mute` / `m unmute` |
| `?` | Help (`Z-10`) | **yes** — `? help` |
| `q` | quit | **yes** — `q quit` |
| `Ctrl-C` | quit | no — universal, never listed |
| `Esc` | **nothing.** Z-01 is the root; there is nothing below it to pop | no |
| `Tab` / `Shift-Tab` | **unbound** — `R = 1`, so there is no second region | no |
| `s` `a` `/` | **unbound.** They belong to `Z-04`. The doors are the affordance here; a hidden duplicate key on a screen whose whole job is to show three doors is noise | no |
| `:` `Ctrl-K` `1`–`9` `n` `p` | **reserved and unbound** (§3.1). Pressing one does nothing and shows no error | no |

**`r` on Z-01 is a design decision.** The nav model binds `r` to *"re-sync, from `Z-04` and
`Z-05`"* and does not forbid it elsewhere. `r` here means the same thing it means everywhere —
*try the network thing again* — and a degrade with no way out is a dead end
([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §3). Flagged in §16 for the
architect to fold into the global table.

> **Width caution for the implementer.** `↑` U+2191 and `↓` U+2193 are **East-Asian-Ambiguous**
> (verified, Unicode 16.0.0). `⏎` U+23CE is **Neutral**. The footer is a flowing line, not a
> column, so an ambiguous-double terminal changes only its length — but it **must** be measured
> with the width-aware function of [`01-design-system.md`](../01-design-system.md) §1.2 before
> truncation, never with `len()`. ASCII fallback under `ZERADO_ASCII=1`: `up/dn move`, `enter choose`.

---

## 10 · The exact copy — ready to paste

| Slot | String |
|---|---|
| breadcrumb | `Zerado ✦ First run` |
| title | `FIRST RUN` |
| intro line 1 (Standard+) | `Zerado keeps your library in one file on this machine.` |
| intro line 2 (Standard+) | `No account. No server. The only traffic is the stores you connect.` |
| intro (Narrow and below) | `One file on this machine.` / `No account, no server.` |
| door 1 title | `Connect a store` |
| door 1 reason (Standard+) | `Steam today. PlayStation, GOG and EA later` |
| door 1 reason (Narrow−) | `Steam today` |
| door 1 reason, **closed** | `OFFLINE   No network. r to check again.` |
| door 2 title | `Add a game by hand` |
| door 2 reason (Standard+) | `Discs, cartridges, anything the stores miss` |
| door 2 reason (Narrow−) | `Discs and cartridges` |
| door 3 title | `Look around first` |
| door 3 reason (Standard+) | `An empty library, and every key that works` |
| door 3 reason (Narrow−) | `An empty library` |
| audio note (Standard+) | `Radio and interface sounds are off until you turn them on in Settings.` |
| audio note (Narrow−) | `Sound is off. See Settings.` |
| footer, default | `↑↓ move   ⏎ choose   , settings   ? help   q quit` |
| footer, door closed | `↑↓ move   ⏎ choose   r check again   , settings   ? help   q quit` |
| footer, audio on | append `   m mute` (or `   m unmute`) before `? help` |

**Two forms per string, not one wrapped string** — *design decision.* A hand-set short form
breaks at a sentence boundary; a greedy wrapper breaks mid-clause (`No account, no / server.`).
The provider-descriptor pattern of [`06-data-seams.md`](../../blueprint/06-data-seams.md) §2.3 is
applied to copy: long form and short form, chosen by tier.

**Voice check.** No exclamation marks · no emoji · the reader is never called a "gamer" ·
the numbers that exist are stated · nothing unbuilt is presented as working (`Steam today.
PlayStation, GOG and EA later` matches `landing-copy.md` §06 and FAQ §14 exactly) ·
casing per `naming.md` (`Zerado` the product, `Settings` the screen).

---

## 11 · 40-column behaviour, and the refusal floor

### 11.1 · RENDER 40×24 — Narrow · `leftInset` 2 · body `36 × 16` · content begins at column 3

Row map: 1 blank · 2–4 band · 5 blank · **6–21 body (16)** · 22 blank · 23 footer · 24 blank.
`40 − 2 × 2 = 36` ✓

```text
0.........1.........2.........3.........

  Zerado ✦ First run

  FIRST RUN

  One file on this machine.
  No account, no server.

  ▌ Connect a store
      Steam today

    Add a game by hand
      Discs and cartridges

    Look around first
      An empty library

  Sound is off. See Settings.

  ↑↓ move  ⏎ choose  ? help  q quit

```

The door becomes **two lines** at Narrow: title on line 1, short reason on line 2 indented to the
title column. This is the responsive table's *"Prose to 2 lines per door"*
([`03-responsive.md`](../../blueprint/03-responsive.md) §3), resolved to one short line plus room.
The footer sheds from the right: `, settings` goes first, then `r check again`; **`? help` and
`q quit` are the last two to go** (nav §6).

### 11.2 · Standard `60 × 24` · body `54 × 16`

Identical to Wide with the reason column at `54 − 2 − 1 − 22 − 2 = 27`; the **short** reasons are
used because the long ones exceed 27. Everything else holds.

### 11.3 · Tiny `< 40` — `32 × 24` · `leftInset` 1 · body `30 × 21`

Row map at Tiny: `OuterMarginY` = 0, `HeaderBandHeight` = **1**, `InnerPaddingY` = 0 →
row 1 title · row 2 gap · rows 3–23 body (**21**) · row 24 footer. `24 − 1 − 1 − 0 − 0 − 1 = 21` ✓

```text
0.........1.........2.........3

 FIRST RUN

 ▌ Connect a store

   Add a game by hand

   Look around first

 ↑↓ move  ⏎ choose  q quit
```

**Door titles only, no prose, no audio line** — the responsive table's Tiny column. The reason
column is gone, so the closed door at Tiny carries the word alone on its own line:
`   OFFLINE  no network`.

### 11.4 · The refusal floor — below `24 × 8`

Z-01 does not render a degraded interface. `Zerado` prints one line and exits with status `2`
([`03-responsive.md`](../../blueprint/03-responsive.md) §6). Specified in
[`Z-11-fatal-error.md`](./Z-11-fatal-error.md) §3.4 as **Z-11-EXIT**, case 4:

```text
Zerado needs at least 24 columns and 8 rows. This terminal is 20 x 6.
```

---

## 12 · `NO_COLOR` rendering — shown, not asserted

`NO_COLOR` set → **zero SGR sequences** (brand §5.4, settled §5.4 of the designer manual).
Bold is an SGR sequence and is therefore **also gone**. This is the render:

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ First run

   FIRST RUN

   Zerado keeps your library in one file on this machine.
   No account. No server. The only traffic is the stores you connect.

     Connect a store         OFFLINE   No network. r to check again.

   ▌ Add a game by hand      Discs, cartridges, anything the stores miss

     Look around first       An empty library, and every key that works

   Radio and interface sounds are off until you turn them on in Settings.

   ↑↓ move   ⏎ choose   r check again   , settings   ? help   q quit

```

**Byte-for-byte the same characters as §3.3.** What each channel does with colour removed:

| Information | Channel that survives |
|---|---|
| which door has focus | `▌` in the gutter — **position** |
| the screen title's rank | UPPERCASE + its own row + the `InterElementGap` above it — **case and spacing** |
| the door is closed | the word `OFFLINE`, then a sentence naming the cause and the key |
| primary vs secondary text | position and spacing — the reason column is a *column* |

**No information is lost.** Colour-budget §10 line 18 passes.
Under `NO_COLOR`, reduced motion is also implied ([`03-responsive.md`](../../blueprint/03-responsive.md) §5) —
Z-01 has no motion, so nothing changes.

---

## 13 · Focus model on this screen

| | |
|---|---|
| **Focusable regions** | **1** — the door list. `Tab` is therefore unbound |
| **Focusable items** | the three doors, **including the closed one** |
| **Initial focus** | door 1 (`Connect a store`) when open; **the first open door** when door 1 is closed |
| **Traversal** | `↑` `↓` / `k` `j`, wrapping at the ends |
| **Focus is never nowhere** | there are always three doors; the list cannot empty |
| **Restored on pop** | returning from `Z-02`/`Z-08`/`Z-10` puts the player back on the door they left (nav §4.1 invariant 4) |
| **`Esc`** | **does nothing.** Z-01 is the stack bottom; there is nothing to go back to. The footer does not hint `esc`, because it would be a lie |

**The closed door stays focusable** — *design decision.* It is not permanently disabled, it is
*closed right now*; `⏎` on it re-reads the link state, which is the way out the offline contract
requires. A skipped item also makes `↑↓` feel broken.

**Focus ring:** three channels, any two sufficient — `▌` (position) + bold (weight) +
`--z-primary` amber on the marker (colour). **Never removed**, in any state
([`02-colour-budget.md`](../02-colour-budget.md) §8.2).

---

## 14 · Colour budget declaration

| Class | Count | Where |
|---|---|---|
| **CHROME cyan** | **0** | *Spent nowhere.* Zero is a pass ([`02-colour-budget.md`](../02-colour-budget.md) §2.2) |
| **STATE cyan** | 0 | no state chip is rendered on Z-01 |
| **Focus ring cyan** | 0 | no *control* holds focus here — the ledger-style row cursor is **amber**, so Z-01 emits no cyan at all |
| **Amber** | 2 marks by default; **3** when audio is on and unmuted | the title `FIRST RUN` (allow-list **1**) · the `▌` focus marker (§1.7 / allow-list **5**) · `▮ AUDIO` (allow-list **9** — ambient voice, not an action, spends no cyan) |
| **Red** | 0 | no scanner, no destructive annunciator, no error text |

### 14.1 · Why Z-01 spends zero chrome cyan — *design decision*

`Z-01` is the first screen a player ever sees, and they have completed **nothing**. Cyan is *the
colour of completion*. Spending the earned colour to advertise a door — a door that is
**closed** in one of this screen's two main states, which would make it a cyan dead control —
cheapens the one signal the palette has. Position and focus carry primacy instead, and they carry
it under `NO_COLOR` too.

**Amber ceiling:** 2 marks (`FIRST RUN` = 9 cells, `▌` = 1 cell) against `80 × 24 = 1920` cells
→ **0.5 %**, far under the 10 % ceiling (§4.2).

---

## 15 · Reuse verdict per element

| Element | Verdict | Note |
|---|---|---|
| Header band | **Build fresh** against #2435, behind the single `Frame` wrapper enforced at the router | [`01-design-system.md`](../01-design-system.md) §2.8 — a screen cannot render frameless by construction |
| Door list | **`bubbles/list`** — *rejected.* Build fresh from `lipgloss` | `bubbles/list` brings its own filtering, pagination, status bar and title chrome, none of which Z-01 wants; three static items with a two-column layout is a `lipgloss.JoinVertical` and a width-aware pad. Forcing the shelf primitive here would cost more styling-out than building |
| Focus marker | **Build fresh** | §1.7 — the shared 2-column gutter primitive, used identically by Z-08 and Z-09 |
| The two prose lines | **Build fresh** — `lipgloss` text layout | `glamour` is for markdown; this is two strings |
| Empty-state composition | **Build fresh** | [`01-design-system.md`](../01-design-system.md) §10.1 is Z-01's own component; Z-01 *is* case (a) |
| Footer | **Build fresh** | one `lipgloss` row with width-aware right-shed |
| Audio indicator | **Build fresh** | one `lipgloss` right-aligned group on the reserved footer row, per [`01-design-system.md`](../01-design-system.md) §5.2. Shared verbatim by Z-02, Z-03, Z-08 and Z-09 — built once |
| Audio cue | **None on Z-01** | [`01-design-system.md`](../01-design-system.md) §15.3 is a closed list of three events — sync complete, error, a game becomes `ZERADO`. Z-01 raises none of them. **No cue on navigation, keystrokes or focus movement**, and none on app start |
| `harmonica` | **Not used** | there is no motion on this screen |

---

## 16 · Screen-specific acceptance criteria

Beyond [`00-design-brief.md`](../00-design-brief.md) §10 and
[`02-colour-budget.md`](../02-colour-budget.md) §10. Each is falsifiable from a rendered artifact.

1. **Z-01 contains no input field, no prompt and no question.** Grep the render for `:` at the end
   of a line and for a text cursor; there must be none.
2. **The third door is present, disabled and reasoned in the offline render** — not hidden. Hiding
   it would make the player believe Zerado cannot connect to stores at all.
3. **The door order is identical in both renders.** Diff §3.2 against §3.3: only the gutter, the
   door-1 reason and the footer differ.
4. **The closed door is reachable with `↑↓` and `⏎` on it re-checks.** There is a way out.
5. **Header-left equals content-left at column 4** on an ANSI-stripped `80 × 24` render — measured,
   not eyeballed.
6. **The block is top-aligned**, never vertically centred ([`02-composition.md`](../../blueprint/02-composition.md) §2).
7. **Zero cyan SGR runs** in the `80 × 24` capture, at every colour depth.
8. **`NO_COLOR` render is character-identical to the coloured render** (§12). Byte-diff after
   stripping SGR must be empty.
9. **The audio note renders only as a note**: it is not focusable, has no key hint beside it, and
   does not appear at Tiny. It names *where* audio lives (Settings), never asks a question, and
   never claims audio is on.
10. **`m` is absent from the footer until audio has been enabled**, and the audio indicator is
    absent with it. On the default first run — audio off — the `80 × 24` capture contains no `▮`,
    no `▯` and no `m`.
10b. **Z-01 plays no sound**, including on app start and on every keypress
    ([`01-design-system.md`](../01-design-system.md) §15.3).
11. **`Esc` produces no visible change and no error.**
12. **Nothing on screen contradicts `landing-copy.md`** — specifically §06 (`Steam syncs today …
    PlayStation, GOG and EA are planned next`), §08 (one file, no account) and §14 (`Steam is
    built and works`).
13. Eight artifacts per [`03-responsive.md`](../../blueprint/03-responsive.md) §7:
    `24×8` · `32×24` · `40×24` · `60×24` · **`80×24`** · `120×40` · `80×24 NO_COLOR=1` ·
    `80×24` forced 16-colour.

---

## 17 · Open for the founder

1. **How Z-01 knows the door is closed (§8.1).** `07-offline-contract.md` §5 says Zerado never
   probes; §6 draws a first-run screen that already knows. The resolution here — read the local
   default-route state, which emits no packet — keeps both, but it is a system fact §5's sentence
   did not anticipate. **Confirm, or accept that the Connect door is always drawn open and the
   failure lands on `Z-02`.**
2. **`r` bound on `Z-01`** (§9). It means what it means everywhere — *try the network thing
   again* — but the global table lists `r` for `Z-04`/`Z-05` only. Confirm, and let
   `fft-tui-architect` fold it into `04-navigation-and-focus.md` §3.
3. **The door order in the offline state (§5.4).** This spec keeps the ratified order and moves
   the focus, where `07-offline-contract.md` §6's drawing reorders. Confirm the refinement.
4. **Zero chrome cyan on the first screen the product ever shows (§14.1).** Confirm that the
   earned colour stays unspent until the player has earned something.
