---
title: Zerado — Z-03 Sync
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-03
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-03 · Sync

> Fills [`../03-designer-manual.md`](../03-designer-manual.md) §3's 16-section contract.
> Composition binding from [`../../blueprint/02-composition.md`](../../blueprint/02-composition.md) §2 —
> single-pane readout, the scanner on **one** row, `R = 0`.
> **`PARTIAL` is a real terminal state, not an error.** A half-synced library is a valid library.

---

## 1 · Identity

| | |
|---|---|
| **Screen** | `Z-03` · Sync |
| **Phase** | 1 |
| **Kind** | Route |
| **Routes in** | `Z-02` on a successful connect · `r` from `Z-04` / `Z-05` |
| **Routes out** | `Z-04 Library` (`⏎`) · pop on `Esc` from a terminal state · `Z-09 Settings` (`,`) |
| **Offline class** | **REFUSES** ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2) — *"Names the reason, shows when the last sync was, offers `r` to retry"* |
| **Records** | one `sync_run` row: `status` ∈ `ok` · `partial` · `failed` · `cancelled`, plus `items_seen` / `items_new` / `items_changed` ([`09-erd.md`](../../blueprint/09-erd.md) §1) |

---

## 2 · Purpose

**Show honest progress on a wait, then a readout of what changed — and survive a failure halfway
through.**

Screen inventory §5: *"`Z-03` must not report a count it has not finished counting."*

---

## 3 · Mockups at 80 columns

Frame row map as [`Z-01-first-run.md`](./Z-01-first-run.md) §3.1 —
`1 + 3 + 1 + 16 + 1 + 1 + 1 = 24`, `80 − 2 × 3 = 74`, content at **column 4**.

### 3.1 · The two progress components, and the rule that picks between them

| The wait | Component | Label |
|---|---|---|
| the denominator is **not** known yet | **the scanner** ([`../01-design-system.md`](../01-design-system.md) §9) | `WAITING ON STEAM` |
| the denominator **is** known | **the determinate readout** (§8) | `SYNCING STEAM   147 / 247` |
| known, then **no progress for 10 s** | back to the scanner (§8.3) | `WAITING ON STEAM` |

*"When Zerado knows the count — '412 of 1,140' — it shows the count and a determinate bar. The
sweep is for when the answer is honestly 'we don't know yet.'"*
([`03-responsive.md`](../../blueprint/03-responsive.md) §5.) **One at a time, never two, never
ambient.**

### 3.2 · RENDER 80×24 — WAITING (indeterminate; nothing has arrived)

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Sync

   SYNC

   WAITING ON STEAM
   ───────────────────━━━────────────────────────────────────────────────────

   Nothing has arrived yet.

   esc cancels. What has arrived is already saved.









   esc cancel   , settings   ? help   q quit

```

The track is `─` U+2500 spanning the full body width, **74 cells**; the pip is **exactly three**
consecutive `━` U+2501, positioned from elapsed time on the ratified **2400 ms** sinusoid
`cubic-bezier(0.45, 0, 0.55, 1)` at **30 fps**, travelling `0 … 74 − 3 = 71`.

### 3.3 · RENDER 80×24 — SYNCING (determinate)

Bar arithmetic, stated so it can be checked: `round(147 / 247 × 74) = 44` filled, `74 − 44 = 30`
unfilled.

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Sync

   SYNC

   SYNCING STEAM                                                     147 / 247
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━──────────────────────────────

   12 new. 4 changed. 131 unchanged.

   esc cancels. What has arrived is already saved.









   esc cancel   , settings   ? help   q quit

```

**The counts sum to what has been seen:** `12 + 4 + 131 = 147`, and `147 / 247` is on the row
above. Nothing on this screen is a count Zerado has not finished counting — the running figures
describe **processed** items, never projected ones.

### 3.4 · RENDER 80×24 — DONE

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Sync

   SYNC

   DONE

   247 games. 6 finished. Last played: 3 weeks ago.

   12 new. 4 changed. 231 unchanged.










   ⏎ library   r sync again   , settings   ? help   q quit

```

`12 + 4 + 231 = 247` ✓. The first line is the brand manual's own ratified copy (§8, voice example
1), **verbatim**.

### 3.5 · RENDER 80×24 — PARTIAL (the sync that failed halfway)

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Sync

   SYNC

   ▌ PARTIAL

     Steam stopped answering after 138 of 247 games.
     The 138 that arrived are in your library — nothing was lost.

     4 new. 2 changed. 132 unchanged.









   r finish the sync   ⏎ library   , settings   ? help   q quit

```

`4 + 2 + 132 = 138` ✓. **The bar is not redrawn at 100 %** — §8.3: *"The bar never jumps to 100 %
on failure."* It is replaced by the readout, which states the number it actually reached.

### 3.6 · RENDER 80×24 — FAILED (nothing arrived; classifier = 401/403)

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Sync

   SYNC

   ▌ STEAM KEY REJECTED

     Steam rejected that key. Check it hasn't been regenerated.
     Settings → Steam.

     Your library is unchanged — nothing was lost.
     Last synced 3 days ago.







   r try again   , settings   ⏎ library   ? help   q quit

```

---

## 4 · Mockup at `120 × 40`

`leftInset` **4** · body `112 × 32` · content at **column 5** · `1 + 3 + 1 + 32 + 1 + 1 + 1 = 40` ✓
Bar: `round(147 / 247 × 112) = 67` filled, `112 − 67 = 45` unfilled.

```text
0.........1.........2.........3.........4.........5.........6.........7.........8.........9.........A.........B.........

    Zerado ✦ Sync

    SYNC

    SYNCING STEAM                                                                                          147 / 247
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━─────────────────────────────────────────────

    12 new. 4 changed. 131 unchanged.

    esc cancels. What has arrived is already saved.
```

### 4.1 · "Per-provider lines" at ExtraWide, read honestly — *design decision*

[`03-responsive.md`](../../blueprint/03-responsive.md) §3 gives `Z-03` *"Scanner + per-provider
lines + running counts"* at ExtraWide. **In Phase 1 there is exactly one `Syncer`.** `physical` is a
`Provider` and deliberately **not** a `Syncer` ([`06-data-seams.md`](../../blueprint/06-data-seams.md) §2.2),
so it is never a line on this screen — rendering it would be a field for a capability that does not
exist (designer manual §7, anti-pattern 14).

So the per-provider block is **the block already drawn, repeated once per `Syncer`**, and with
`N = 1` the ExtraWide render is the Wide render at 112 columns. When `N > 1` the head becomes
`SYNCING` alone and each provider gets its own label-and-bar pair separated by an
`InterElementGap`:

```text
    SYNCING

    Steam                                                                                                  147 / 247
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━─────────────────────────────────────────────

    GOG                                                                                                       9 / 61
    ━━━━━━━━━━━━━━━━━───────────────────────────────────────────────────────────────────────────────────────────────
```

*(Illustrative of the structure only — **`GOG` is not built and must never render in Phase 1.**)*
This spec deliberately ships the `N = 1` form and proves the `N > 1` shape rather than faking a
second provider into a screenshot.

---

## 5 · Visual hierarchy

**The one thing the player must see first: the number.**

| Rank | Element | Channel | Note |
|---|---|---|---|
| 1 | `SYNC` | case + weight + `--z-primary` | the H1 |
| 2 | **the count `147 / 247`** — or, in a terminal state, **the first fact of the readout** | position (right edge of its row, alone) + `--z-text` **255**, the brightest text on screen | *"The number is always shown — the bar is the ornament, the count is the information"* (§8.1) |
| 3 | the bar, or the scanner pip | **motion** (the only moving thing) + stroke **weight** | motion is a channel here and nowhere else in Zerado |
| 4 | the readout label | case + `--z-text-secondary` | `SYNCING STEAM`, `WAITING ON STEAM` |
| 5 | the counts line | `--z-text` | three facts |
| 6 | the standing note | `--z-text-secondary` | `esc cancels…` |
| 7 | breadcrumb, footer | chrome | |

**In a terminal state the annunciator takes rank 2** — `▌` plus an UPPERCASE heading — and the
number moves to rank 3, inside the sentence.

### 5.1 · There is no focus marker on this screen, and that is correct

`R = 0`. Nothing on `Z-03` is focusable, in any state: it is a readout, and every action is a
footer key. There is therefore no focus ring to draw, and **its absence is not a removal**
([`../02-colour-budget.md`](../02-colour-budget.md) §8.2 governs the *focused element*; there is
none). Stated because a reviewer will look for one.

**A consequence worth naming:** `▌` is unambiguous on `Z-03`. It is only ever the terminal-state
annunciator, never a focus marker — unlike on `Z-02` and `Z-08`, where §5.1 of those specs has to
keep the two apart.

---

## 6 · Every applied spacing token, by name

| Token | Wide value | Where Z-03 spends it |
|---|---|---|
| `OuterMarginX` | **2** | frame inset |
| `OuterMarginY` | **1** | rows 1 and 24 |
| `InnerPaddingX` | **1** | inside the frame |
| `InnerPaddingY` | **1** | row 22 |
| `InterElementGap` | **1** | breadcrumb→title · band→body · readout→counts · counts→standing note · annunciator heading→body · between per-provider blocks when `N > 1` |
| `HeaderBandHeight` | **3** | `hasSubtitle = false` |
| `leftInset` | **3** | header-left **==** content-left at column 4 |
| track / bar width | **`BodyRect.w`** — 74 at Wide, 112 at ExtraWide, 36 at Narrow | never a magic number; always the body width |
| message-block indent | `leftInset + 2` = **column 6** | the two columns the `▌` annunciator field occupies, so the body hangs under the heading |

**No in-body key-hint line** — *design decision, applied across this cluster.*
[`../01-design-system.md`](../01-design-system.md) §10.1 and §11.1 draw those components with their
own hint row because they are specified **standalone**. Hosted on a framed screen that already has
the Spacing Canon's **reserved footer row**, a second hint row is two footers, and
[`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §6 gives the footer that
job. The footer therefore changes per terminal state, and the body keeps the row. **Flagged in §17.**

---

## 7 · Colour, glyph and label for every state shown

| State | Token | Hex | ANSI-256 | 16-colour | Ratio | Glyph / structure | Label | `NO_COLOR` |
|---|---|---|---|---|---|---|---|---|
| screen title | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | — | `SYNC` UPPER + bold | **yes** — case |
| readout label | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** | — | `SYNCING STEAM` / `WAITING ON STEAM` | **yes** |
| the count | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** | — | `147 / 247` | **yes** |
| **bar, filled** | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | `━` U+2501 **heavy** | the count states it in words | **yes** — **heavy vs light stroke** |
| **bar, unfilled** | `--z-primary-muted` `#8A5E00` — **ANSI index underived** → **interim: uncoloured** (§8.2) | — | *underived* | `bright black` | — | `─` U+2500 **light** | — | **yes** |
| **scanner track** | `--z-scanner-track` `#5C1414` — **underived** → **interim: uncoloured** (§9.4) | — | *underived* | `black` | — | `─` U+2500 | `WAITING ON STEAM` | **yes** |
| **scanner pip** | `--z-scanner` | `#FF2E2E` | **9** | `bright red` | 5.25 — **motion, not text** | `━` × **exactly 3** | — | **yes** — heavy stroke |
| **`DONE` head** | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | no annunciator — **nothing to annunciate** | `DONE` | **yes** — case |
| **`PARTIAL` annunciator** | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | `▌` U+258C | `PARTIAL` | **yes** — `▌` + case |
| **`CANCELLED` annunciator** | `--z-border-strong` | `#64748B` | **67** | `bright black` | **4.08** | `▌` U+258C | `CANCELLED` | **yes** |
| **`FAILED` annunciator** | `--z-scanner` | `#FF2E2E` | **9** | `bright red` | 5.25 — structure, not text | `▌` U+258C | `STEAM KEY REJECTED` etc. | **yes** |
| **error / refusal text** | **uncoloured + bold** — documented interim; `--z-scanner-300` `#FF6B6B` (**6.99** AA) has no derived index (§11.2) | `#FF6B6B` | *underived* | `bright red` | **6.99** | — | the sentences | **yes** |
| audio indicator | applied from [`../01-design-system.md`](../01-design-system.md) §5.2 — see [`Z-01-first-run.md`](./Z-01-first-run.md) §7.1 | | | | | `▮` / `▯` **Neutral** | `AUDIO` / `MUTED` | **yes** |

### 7.1 · The annunciator ladder — *design decision*

Four terminal states need four different weights of alarm, and the product already owns the
vocabulary for three of them. This spec spends it as a ladder:

| State | Annunciator | Why that rung |
|---|---|---|
| `DONE` | **none** | nothing happened that needs announcing. A success annunciator is decoration |
| `CANCELLED` | **chrome** `--z-border-strong` | the player did it on purpose. Nothing is wrong. This is the **informational** class of the degrade banner (§12.1) |
| `PARTIAL` | **amber** `--z-primary` | *"Amber appears only when the player must do something"* (§12.2) — and the copy asks for `r` |
| `FAILED` | **red** `--z-scanner` | the error state (§11.1). Alarm as **structure**, never as text |

**`PARTIAL` is deliberately not red.** *"A half-synced library is a valid library"*
([`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §3.2). Colouring a valid
outcome red would be the same defect as colouring `OFFLINE` red — calling a designed behaviour a
fault ([`../02-colour-budget.md`](../02-colour-budget.md) §5.2, failure gallery item 6).

---

## 8 · The full state table

| # | State | Trigger | Renders | Copy |
|---|---|---|---|---|
| 1 | **First run — the very first sync** | pushed from `Z-02` | §3.2 → §3.3 → §3.4 | on failure, *"Nothing has ever been synced."* replaces *"Last synced 3 days ago."* — see §8.1 |
| 2 | **Waiting** *(indeterminate)* | request sent, nothing streamed | §3.2 — **the scanner** | `WAITING ON STEAM` / `Nothing has arrived yet.` |
| 3 | **Syncing** *(determinate)* | first item arrives → the denominator is known | §3.3 — the bar | `SYNCING STEAM` + `147 / 247` |
| 4 | **Stalled** | ≥ 10 s with no progress (§8.3) | back to the scanner, bar removed | `WAITING ON STEAM` |
| 5 | **`DONE`** | channel closes, no error | §3.4 · `sync_run.status = ok` · cue `sync.done` | §10 · **DONE** |
| 6 | **`PARTIAL`** | items arrived, then the stream errored | §3.5 · `status = partial` · **no cue** | §10 · **PARTIAL** |
| 7 | **`CANCELLED`** | `esc` or `q` during 2–4 | §3.5 shape · `status = cancelled` · **no cue** | §10 · **CANCELLED** |
| 8 | **`FAILED` — no route / DNS** | classifier | §3.6 shape · `status = failed` · cue `sync.failed` | §10 · **F-A** |
| 9 | **`FAILED` — timeout / 5xx** | classifier | same | §10 · **F-B** — ratified |
| 10 | **`FAILED` — 401 / 403** | classifier | §3.6 | §10 · **F-C** — ratified |
| 11 | **`FAILED` — 200 + empty** | channel closed with zero items | same | §10 · **F-D** — ratified. **The library is not emptied** — §8.2 |
| 12 | **Empty** | **N/A** — `Z-03` has no list. The *result* being empty is state 11, which is a refusal, not an empty state | — | — |
| 13 | **Offline banner** | **N/A on `Z-03`.** The screen's whole body is the message; a banner above it would say the same thing twice. `Z-04` carries the banner ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §3) | — | — |
| 14 | **`q` during a sync** | | cancel the context, **then** quit, saying so in one line on the way out (nav §3.2) | `Sync cancelled. 138 of 247 games were saved.` |
| 15 | **Audio** | | rows as [`Z-01-first-run.md`](./Z-01-first-run.md) §8 rows 10–12 | — |
| 16 | **`NO_COLOR`** | env | §12 · **and the pip parks** | identical text |
| 17 | **Below `24 × 8`** | | never renders — [`Z-11-fatal-error.md`](./Z-11-fatal-error.md) | the refusal sentence |

### 8.1 · The first-run row, spelled out

Every `FAILED` and `PARTIAL` copy ends with an age line. On the **first ever sync** there is no
previous `sync_run`, so `Last synced 3 days ago.` would be a fabricated fact.

| `LastSync(steam)` | Line rendered |
|---|---|
| a previous run exists | `Last synced 3 days ago.` |
| **`nil`** | `Nothing has ever been synced.` |

The age is said the way people say it — `3 days ago`, `2 hours ago`, `just now`, `last June` —
never an RFC 3339 timestamp ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §4.1).

### 8.2 · A sync that returns zero items **never empties the library** — *design decision*

State 11 can occur on a **re-sync** of a 247-game library whose owner has just made their profile
private. A naive "the provider's view is the truth" upsert would delete every row and the screen
would then be honestly reporting a catastrophe it caused.

**Rule:** a `Sync` that yields zero items is a **refusal**, not an empty result set. No deletion, no
`UpsertBatch`, and the `sync_run` records `failed` with the classified reason. The copy's
`Your library is unchanged — nothing was lost.` is then true, which is the only condition under
which it may be printed. **Flagged in §17** because it is a `Store` contract, not a screen one.

---

## 9 · The key map

No text input exists on `Z-03`, so every single-key shortcut is live.

| Key | Does | In the footer? |
|---|---|---|
| `esc` | **running:** cancel the sync → `CANCELLED`. **terminal state:** pop the route | **yes** — `esc cancel` while running; not listed at rest, where `⏎ library` is the way out |
| `q` | **running:** cancel, then quit, with one line on the way out. **otherwise:** quit | **yes** — `q quit` |
| `Ctrl-C` | quit | no |
| `⏎` | open the library (`Z-04`) | **only in a terminal state** — `⏎ library` |
| `r` | start the sync again. On `PARTIAL`, **finish** it | **only in a terminal state** — `r sync again` / `r finish the sync` / `r try again` |
| `,` | Settings (`Z-09`) | **yes** — `, settings`. It is the destination the ratified `Settings → Steam` copy names, so it must be reachable from the screen that prints it |
| `?` | Help | **yes** |
| `m` | toggle mute | only when audio is enabled |
| `↑` `↓` `g` `G` `Tab` | **unbound** — `R = 0`, nothing to move between | no |
| `s` `a` `/` | **unbound** — they belong to `Z-04` | no |

**Footer per state** — the footer is the hint block (§6):

| State | Footer |
|---|---|
| waiting / syncing | `esc cancel   , settings   ? help   q quit` |
| `DONE` | `⏎ library   r sync again   , settings   ? help   q quit` |
| `PARTIAL` | `r finish the sync   ⏎ library   , settings   ? help   q quit` |
| `CANCELLED` | `r start again   ⏎ library   , settings   ? help   q quit` |
| `FAILED` | `r try again   , settings   ⏎ library   ? help   q quit` |

---

## 10 · The exact copy — ready to paste

### 10.1 · Chrome and running

| Slot | String |
|---|---|
| breadcrumb | `Zerado ✦ Sync` |
| title | `SYNC` |
| waiting label | `WAITING ON STEAM` |
| waiting body | `Nothing has arrived yet.` |
| syncing label | `SYNCING STEAM` |
| syncing count | `147 / 247` |
| running counts | `12 new. 4 changed. 131 unchanged.` |
| standing note | `esc cancels. What has arrived is already saved.` |
| quit-mid-sync line | `Sync cancelled. 138 of 247 games were saved.` |

### 10.2 · The four terminal states

**DONE** — line 1 is the brand manual's ratified copy, verbatim
```text
   DONE

   247 games. 6 finished. Last played: 3 weeks ago.

   12 new. 4 changed. 231 unchanged.
```
When the provider does not report `LastPlayed` (`Capabilities.LastPlayed == false`) the third fact
is **omitted**, not rendered empty: `247 games. 6 finished.` A field for a capability that does not
exist is anti-pattern 14.

**PARTIAL**
```text
   ▌ PARTIAL

     Steam stopped answering after 138 of 247 games.
     The 138 that arrived are in your library — nothing was lost.

     4 new. 2 changed. 132 unchanged.
```

**CANCELLED**
```text
   ▌ CANCELLED

     You stopped the sync after 138 of 247 games.
     The 138 that arrived are in your library — nothing was lost.

     4 new. 2 changed. 132 unchanged.
```

**FAILED · F-A — no route / DNS**
```text
   ▌ SYNC FAILED

     No network. Nothing was fetched and your library is unchanged.
     Last synced 3 days ago.
```

**FAILED · F-B — timeout / 5xx** — sentence one ratified, verbatim
```text
   ▌ STEAM UNREACHABLE

     Steam didn't answer. Not your key — their end, or the connection.
     Your library is unchanged — nothing was lost.
     Last synced 3 days ago.
```

**FAILED · F-C — 401 / 403** — ratified, verbatim, wrapped at the body width
```text
   ▌ STEAM KEY REJECTED

     Steam rejected that key. Check it hasn't been regenerated.
     Settings → Steam.

     Your library is unchanged — nothing was lost.
     Last synced 3 days ago.
```

**FAILED · F-D — 200 + empty** — ratified, verbatim
```text
   ▌ STEAM PROFILE PRIVATE

     Steam returned an empty library.
     Game details are private on your profile — Steam won't share the
     list until that's public. Settings → Privacy.

     Your library is unchanged — nothing was lost.
```

### 10.3 · Recorded departure from one ratified string

| Ratified | On `Z-03` | Why |
|---|---|---|
| `No network. Last synced 3 days ago — everything below still works. r to retry.` | **F-A** above | *"everything below"* names a library **below the banner**, which is the `Z-04` composition. On `Z-03` there is nothing below — the whole body is the message. The two facts it carries (no network, the age) are both kept, and `r` is in the footer where the nav model puts it |

`r to retry` is kept as a **footer key** rather than in prose on `Z-03`, because unlike `Z-02` there
is no text input here and `r` genuinely is a live single-key shortcut.

**Voice check.** No exclamation marks · no emoji · never "gamer" · every number stated and every
set of counts sums · each refusal names what, why, next action and what happened to the data ·
none says *"Something went wrong."*

---

## 11 · 40-column behaviour, and the refusal floor

### 11.1 · RENDER 40×24 — Narrow · `leftInset` 2 · body `36 × 16`

Responsive table: *"Scanner + one count line."* The determinate form stacks label above bar
(§8.4's own 40-column drawing). Bar: `round(147 / 247 × 36) = 21` filled, `15` unfilled.

```text
0.........1.........2.........3.........

  Zerado ✦ Sync

  SYNC

  SYNCING STEAM        147/247
  ━━━━━━━━━━━━━━━━━━━━━───────────────

  12 new. 4 changed.

  esc cancels. What has
  arrived is already saved.




  esc cancel  ? help  q quit

```

The count contracts to `147/247` (§8.4's own narrow form). The counts line drops its third fact
before it drops a whole word — **whole facts, never mid-word** (§5.5).

### 11.2 · Standard `60 × 24` · body `54 × 16`

Responsive table: *"Counts only, no per-provider lines."* With `N = 1` there is nothing to shed, so
Standard renders Wide's composition at 54 columns. Bar 54 cells.

### 11.3 · Tiny `< 40` — `32 × 24` · body `30 × 21`

**The scanner is dropped, not shrunk** ([`03-responsive.md`](../../blueprint/03-responsive.md) §5):
*"a three-cell pip on a ~28-cell track is a blinking dash, which reads as a glitch rather than as
hardware. The count line replaces it, which is more information in less space."*

```text
0.........1.........2.........3

 SYNC

 SYNCING STEAM

 147 / 247

 12 new. 4 changed.

 esc cancels.

 esc  ? help  q quit
```

The determinate **bar** is also dropped at Tiny for the same reason and by the same rule — the
count is the information, the bar is the ornament (§8.1). In the WAITING state at Tiny the body is
the label and one line: `WAITING ON STEAM` / `Nothing has arrived yet.`

### 11.4 · The refusal floor — below `24 × 8`

`Z-03` never renders; see [`Z-11-fatal-error.md`](./Z-11-fatal-error.md) §3.4. **A sync already
running when the terminal is dragged below the floor is not killed** — the context is untouched,
the screen is replaced by the refusal sentence, and dragging back reveals the sync still going
([`03-responsive.md`](../../blueprint/03-responsive.md) §6: *"it does not exit mid-session"*).

---

## 12 · `NO_COLOR` rendering — shown, not asserted

Zero SGR, **and `NO_COLOR` is also the reduced-motion signal**
([`03-responsive.md`](../../blueprint/03-responsive.md) §5). The pip therefore **parks at the centre
of the track at full weight** and does not travel: `pipLeft = round((74 − 3) / 2) = 36`. It is
deliberately **not hidden** — the lit slot is an identity element, the travel is the decoration.

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Sync

   SYNC

   WAITING ON STEAM
   ────────────────────────────────────━━━───────────────────────────────────

   Nothing has arrived yet.

   esc cancels. What has arrived is already saved.









   esc cancel   , settings   ? help   q quit

```

And the determinate bar, with zero SGR — the fill is carried entirely by **stroke weight**:

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   SYNCING STEAM                                                     147 / 247
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━──────────────────────────────

   12 new. 4 changed. 131 unchanged.
```

| Information | Channel that survives |
|---|---|
| how far the sync has got | **the count, in digits** — the bar is the ornament (§8.1) |
| the bar's fill | **heavy `━` against light `─`** — one ink, two weights. This is why `█`/`░` are forbidden as a pair: they are different width classes (§1.2) |
| the pip | the same heavy-against-light contrast |
| that it is a wait, not a hang | the parked pip is still **visible**; the label says `WAITING ON STEAM` |
| which terminal state | `▌` plus an UPPERCASE word — `PARTIAL`, `CANCELLED`, `STEAM KEY REJECTED` |
| that a partial sync kept its data | the sentence `The 138 that arrived are in your library — nothing was lost.` |

**No information is lost.** And with `ZERADO_NO_AUDIO=1` set as well, nothing is lost either: the
`DONE` readout is the signal, and the cue was only ever the second one
([`12-audio.md`](../../blueprint/12-audio.md) §9).

---

## 13 · Focus model on this screen

| | |
|---|---|
| **Regions** | **0**. `Z-03` is a readout; every action is a footer key |
| **Focusable items** | none, in every state. `Tab`, `↑`, `↓`, `g`, `G` are unbound |
| **Focus ring** | **none is drawn, because no element is focusable.** This is not a removal — §8.2 of the colour budget governs *the focused element*, and there is none |
| **`Esc` — running** | **cancel the sync.** The `context` is cancelled, in-flight I/O aborts, **what was written stays written**. Lands on `CANCELLED` |
| **`Esc` — terminal state** | **pop the route** |
| **`q` — running** | cancel, then quit, with one line on the way out |

This matches [`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §5 exactly:
*"`Z-03` while a sync is running → Cancel the sync → Lands on `Z-03`'s `CANCELLED` state, then
`Esc` again to leave."* **No keyboard trap** (2.1.2).

---

## 14 · Colour budget declaration

| Class | Count | Where |
|---|---|---|
| **CHROME cyan** | **exactly 1, and 0 while running** | the primary footer hint of the current terminal state: `⏎ library` (`DONE`, `CANCELLED`) · `r finish the sync` (`PARTIAL`) · `r try again` (`FAILED`). **While the sync runs there is nothing to urge, so the budget is 0** |
| **STATE cyan** | **0** | `Z-03` renders no state chip. In particular `6 finished` in the `DONE` line is **prose about state and is not cyan** ([`../02-colour-budget.md`](../02-colour-budget.md) §2.4 — *"the most tempting violation in the whole product"*) |
| **Focus-ring cyan** | 0 | nothing is focusable |
| **Amber** | running: label + bar fill + `▮ AUDIO`. Terminal: title + head/annunciator | allow-list **1** (title) · **2** (readout labels) · **4** (determinate fill and its unlit track) · **7** (the action-required annunciator on `PARTIAL`, and only that one) · **9** (audio) |
| **Red** | **1**, never more | the scanner pip (§5 list **1**) **or** the `FAILED` annunciator (§5 list **2**, see [`Z-02-connect-a-store.md`](./Z-02-connect-a-store.md) §14.1). They are never on screen at the same time |
| **Error text red** | **0** | the documented interim — uncoloured + bold. `--z-scanner` sets no words |

**Amber ceiling, the one screen where it is worth computing.** The determinate bar is the largest
amber run in the product: at `120 × 40` a 95 %-full bar is **106 cells**, plus label 13 and title 4
→ **123 of 4800 = 2.6 %**. At `80 × 24`: 70 + 13 + 4 = **87 of 1920 = 4.5 %**. Both comfortably
under the **10 %** ceiling (§4.2), and the ceiling is never approached because the bar is one row.

---

## 15 · Reuse verdict per element

| Element | Verdict | Note |
|---|---|---|
| Determinate readout | **Build fresh. `bubbles/progress` does not fit** | §8.5 — its default ramp is a **truecolor gradient**, which spends colour the budget does not have and cannot resolve at the 16-colour floor. A progress bar is a division, a repeat, and the width-aware fill from §1.2 |
| Scanner | **Build fresh. Not `harmonica`** | §9.8 — harmonica models a **damped spring converging on a target**; the scanner is an **undamped, infinite, alternating sinusoid**. Evaluate the bezier directly. Recorded because reaching for the shelf primitive here is a defensible-looking mistake |
| The 30 fps ticker | **Build fresh** — a `tea.Tick`, torn down with the view | §9.6 — runs **only** while an indeterminate wait is on screen; **no leaked goroutine, no timer surviving the view**; a tick must never block input |
| Terminal-state block | **Build fresh. Not `charmbracelet/log`** | §11.5 — `log` is the structured **developer** logging role and its level colours are its own palette |
| Header band, footer, audio indicator | **Build fresh**, shared | §2.8 · §5.2 |
| Audio cue | `sync.done` on `DONE`; `sync.failed` on `FAILED`; **no cue on `PARTIAL` or `CANCELLED`** | the cue list is **closed** (§15.3) and neither is on it. A sync the player just cancelled must not make a noise. **The cue is never the only signal** — the `DONE` readout renders identically with audio off, absent, or never enabled ([`12-audio.md`](../../blueprint/12-audio.md) §9) |

---

## 16 · Screen-specific acceptance criteria

1. **The counts always sum.** In every render: running `new + changed + unchanged == seen`;
   `DONE` `new + changed + unchanged == total`; `PARTIAL`/`CANCELLED` the same against the reached
   figure. A render where they do not is a fail, not a nit.
2. **No count is ever projected.** Grep the running render for any figure that is not a completed
   tally — the only ratio on screen is `seen / total`, and `total` is the provider's own reported
   size, not an estimate.
3. **The bar never reaches 100 % on a failure** (§8.3). Force a mid-stream error and assert the
   last drawn fill matches the reached count.
4. **`PARTIAL` is not red**, and `CANCELLED` is not red. Scan the capture for `ESC[38;5;9m` /
   `ESC[91m` in both states; there must be none.
5. **Exactly one moving thing.** Never a scanner and a determinate bar together; never two
   scanners; never a scanner in a terminal state.
6. **The scanner is dropped at Tiny**, not shrunk (§11.3), and **the bar with it**.
7. **A zero-item sync does not delete a row** (§8.2). Sync a 247-game library against a provider
   that returns nothing and assert the row count is still 247.
8. **A cancelled sync keeps what arrived.** Cancel at item 138 and assert 138 rows are present and
   `sync_run.status = cancelled`.
9. **No goroutine survives the screen.** `goleak` across a full push → sync → cancel → pop cycle,
   including the 30 fps ticker.
10. **The first-run row is real** (§8.1). With no prior `sync_run`, no render anywhere may contain
    the string `Last synced`.
11. **`Settings → Steam` is reachable from the screen that prints it** — `,` is bound and in the
    footer.
12. **`NO_COLOR` parks the pip at the centre** at full weight, and does not hide it (§12).
13. **Chrome-cyan count is 0 while running and exactly 1 in each terminal state.**
14. Eight artifacts per [`03-responsive.md`](../../blueprint/03-responsive.md) §7, **plus** one per
    terminal state at `80 × 24` — `DONE`, `PARTIAL`, `CANCELLED`, and all four `FAILED` variants.
    These are the screens this spec exists for.

---

## 17 · Open for the founder

1. **A zero-item sync must never empty the library (§8.2).** It is a `Store`/sync contract, not a
   screen one, and without it the copy `Your library is unchanged — nothing was lost.` can be a
   lie. Route to `fft-tui-architect` for the seam, and confirm.
2. **No in-body key-hint line on framed screens (§6).** The design system draws its empty-state and
   error-state components with their own hint row; hosted on a framed screen that already reserves
   a footer, that is two footers. Confirm the footer wins, and let the component sections say so.
3. **`PARTIAL`'s annunciator is amber, not red (§7.1).** It is a valid outcome that asks for one
   key. Confirm the ladder chrome → amber → red.
4. **The `DONE` line drops its third fact** when the provider does not report `LastPlayed` (§10.2)
   rather than rendering it empty. Confirm — it means the ratified three-fact sentence is
   two facts for a provider that cannot report the third.
