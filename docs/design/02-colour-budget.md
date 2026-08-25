---
title: Zerado — The Colour Budget
discipline: DESIGN SYSTEM
doc-no: ZRD-DESIGN-03
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: concept-explainer
---

# Zerado — The Colour Budget

**This is a rule, not a suggestion.** It exists so that a screen with cyan in five places
**fails review instead of shipping**. Every clause below is written to be checked against a
rendered artifact by someone who did not design the screen.

Governing authority: brand manual §4.1. *"A screen with cyan in five places has thrown away the
only signal the palette has."*

---

## 1 · Why a budget exists

The palette carries **product meaning**, not decoration:

| Colour | Role | Frequency |
|---|---|---|
| **Amber** `--z-primary` | The ambient voice — the phosphor readout, the colour the machine speaks in | Common, deliberately |
| **Cyan** `--z-accent` | **Earned.** Completion, and the single most important action on a screen | Rare, measured in single elements |
| **Red** `--z-scanner` | Motion and alarm **only** | Closed list, §5 |
| **Chrome** `--z-state-not-started` / text greys | Structure and the dormant state | Everywhere structural |

Brand §4.1's target distribution on any surface: **60 % chassis · 30 % chrome · 10 % amber**,
with **cyan measured in single elements, not percentages.** This document turns that sentence
into something countable.

---

## 2 · Cyan — the two classes

The single mistake this budget prevents is treating all cyan alike. There are **two classes**,
and only one of them is budgeted.

### 2.1 · STATE cyan — data, unbounded, uncounted

The `ZERADO` state — its glyph `◉`, its label `ZERADO`, its ASCII `[*]` — renders in
`--z-accent` `#19E0FF` (ANSI-256 **45**, 16-colour `bright cyan`) **wherever it appears**.

A library with 40 finished games shows 40 cyan chips. **That is not a violation.** It is the
product working: the colour of completion, applied to the things that are complete. Cyan here is
*data*, and data is not budgeted. A player with more cyan on screen has finished more games,
which is the entire emotional payload of the product.

**Qualifies as STATE cyan — and nothing else does:**
- the `◉` glyph in a state chip
- the `ZERADO` label in a state chip
- the `[*]` ASCII form under `ZERADO_ASCII=1`
- the `[ ◉ ZERADO ]` filter chip **when selected**

### 2.2 · CHROME cyan — the budget

Every other cyan mark on the screen. **The budget is ONE (1) per screen.**

That one is *"the single most important call to action on a screen"* (brand §4.1). Concretely,
on a Zerado screen it is at most one of:

- the primary action hint in the status bar or an empty state — `s  connect Steam`
- the confirming action on a **non**-destructive dialog
- the one key hint the screen most wants pressed

**A screen may spend zero.** Cyan is earned; a screen with nothing to urge does not need to urge
anything. Zero is a pass. **Two is a fail.**

### 2.3 · The exception list — closed

Exactly three things are cyan and **do not count against the budget**:

| Exception | Why it is exempt |
|---|---|
| **STATE cyan** (§2.1) | It is data, not chrome |
| **The focus ring** `--z-focus-ring` | Required by WCAG 2.4.7 and by brand §4.2 (*never removed*). It is singular by definition — exactly one element holds focus — so it cannot multiply |
| **The text cursor** in an input | The terminal's own caret; part of the focus indication |

Nothing else is ever added to this list. If a fourth exception is proposed, it is a change to the
brand and goes through §9.

### 2.4 · Where cyan is forbidden outright

Never a border. Never a heading or a screen title. Never a divider or a rule. Never a
readout label. Never emphasis in body text. Never a progress bar. Never a chart series.
Never the word `zerado` inside a summary sentence — **`247 games · 6 zerado · 12 in progress`
renders in `--z-text-secondary` with the numerals in `--z-text`, and none of it is cyan.**
That last one is the most tempting violation in the whole product, and it is still a violation:
the summary is prose about state, not a state cell.

---

## 3 · How to count cyan — the method

Reproducible, so two reviewers reach the same verdict.

### 3.1 · Machine method — from the ANSI stream

1. Render the screen headless (`freeze` / `vhs`) at the target viewport **and at an overflowing
   row count** (≥ 400 games), capturing the **raw output including escape sequences**.
2. Find every foreground-cyan SGR run. Cyan is emitted in one of three encodings depending on
   the terminal's reported depth:

   | Depth | Sequence |
   |---|---|
   | truecolor | `ESC[38;2;25;224;255m` |
   | 256-colour | `ESC[38;5;45m` |
   | 16-colour | `ESC[96m` (bright cyan foreground) |

3. **Classify each run by its payload**, not by its position:
   - payload contains `◉`, `ZERADO`, or `[*]` → **STATE cyan** → *not counted*
   - payload is a focused **control's** ring (a text input, a form field) → **focus ring** →
     *not counted*. Note the ledger's row cursor is `▌` in **amber**, not cyan, so it never
     appears in a cyan scan at all
   - anything else → **CHROME cyan** → *counted*
4. **Pass if the chrome-cyan count is 0 or 1. Fail at 2 or more.**

This is scriptable and belongs in the Screen-Quality Gate as an automated check.

### 3.2 · Human method — from the screenshot

For a reviewer looking at the founder-facing image:

> Ignore every cyan mark that sits in a state column beside the word `ZERADO`. Ignore the focus
> cursor. **Count what is left. More than one, the screen fails.**

### 3.3 · The `NO_COLOR` cross-check

Run the same screen with `NO_COLOR=1`. **If any information disappears, the screen was encoding
meaning in colour** and fails 1.4.1 regardless of the count. The budget and the co-render rule
are checked together, always.

---

## 4 · Amber — role and ceiling

**Amber is the ambient voice.** It is meant to be everywhere the machine speaks. But "the
ambient voice" is not "the default colour", and an all-amber screen is a light show, not a
cockpit.

### 4.1 · The allow-list — where amber may appear

1. The **screen title** in the header band.
2. **Readout labels** and section heads.
3. The **`IN PROGRESS` state** — glyph, label, ASCII `[~]`.
4. **Determinate progress fill** (`━`), and its unlit track (`--z-primary-muted`).
5. **Key hints** in an empty state or a status bar — *except* the one that spends the cyan.
6. The **terminal mark** `[0]` on the splash and in `--help`.
7. The **action-required degrade banner** — and only the action-required one.
8. The **filter mode sigil** `/`.

Anything not on this list is not amber. Body text is `--z-text`. Structure is chrome.

### 4.2 · The ceiling — the backstop

Amber-coloured cells must not exceed **10 % of the rendered viewport**, counting **all** cells
including blanks as the denominator — the faithful translation of brand §4.1's 60/30/10, where
the chassis is the unpainted ground.

At 80 × 28 = 2240 cells, that is **224 amber cells**. A title row, a handful of readout labels
and a set of `IN PROGRESS` chips sit far below it. A screen approaching the ceiling has stopped
using amber as a voice and started using it as a wash.

Measured the same way as §3.1, matching `ESC[38;5;214m` / `ESC[38;2;255;176;0m` / `ESC[93m`.

---

## 5 · Red — motion and alarm only, enumerated

**A closed list of three. Nothing else in Zerado is ever red.**

| # | Use | Token | Hex | ANSI-256 | Ratio | Note |
|---|---|---|---|---|---|---|
| 1 | The **scanner pip**, during a genuinely indeterminate wait | `--z-scanner` | `#FF2E2E` | **9** | 5.25 | Motion. Never ambient (`01-design-system.md` §9.3) |
| 2 | The **destructive-confirmation annunciator** `▌` | `--z-scanner` | `#FF2E2E` | **9** | 5.25 | Alarm as **structure**, not as text |
| 3 | **Error text** | `--z-scanner-300` | `#FF6B6B` | *underived* | **6.99** AA | **The readable red.** Interim: uncoloured + bold until the index is derived |

### 5.1 · The distinction that must not blur

**`--z-scanner` `#FF2E2E` is 5.25 and is `motion/alarm, not text`.** It draws the moving pip and
the annunciator bar. **`--z-scanner-300` `#FF6B6B` is 6.99 and is the only red that may set
words.** Using `--z-scanner` for a sentence violates the brand's own stated rule even though the
ratio happens to clear 4.5.

### 5.2 · Red is forbidden for

**Offline.** The product publicly promises it works with the network off — colouring a kept
promise red calls it a fault. Offline is **chrome** (`01-design-system.md` §12).
Also never: the `ABANDONED` state (that is orchid), a destructive action's *body text*, a
low-price flag (a good price is not an alarm), a heading, a border, or the terminal mark.

---

## 6 · Chrome — structure

Chrome is the brushed-stainless register: text, boundaries, and the dormant state.

| Use | Token | ANSI-256 | 16-col |
|---|---|---|---|
| Primary text | `--z-text` | 255 | `bright white` |
| Secondary text, summaries, readout values | `--z-text-secondary` | 249 | `white` |
| Tertiary — de-emphasised, unavailable, unselected | `--z-text-tertiary` | *underived* | `white` |
| `NOT STARTED` state | `--z-state-not-started` | 247 | `white` |
| Control boundaries | `--z-border-strong` | 67 | `bright black` |
| Decorative hairlines | `--z-border` | 236 | `black` |
| Offline / informational annunciator | `--z-border-strong` | 67 | `bright black` |

**The warm grey `#A5A29B` is engineering, not taste.** The blue-cast alternative `#9FB0C6`
collapsed against the cyan at **ΔE 8.8 under deuteranopia**; the warm grey measures **25.8**.
Never "correct" it back toward blue.

---

## 7 · The 16-colour floor — the consequence that changes layouts

At 16 colours **`--z-surface` `#0B0D14` and `--z-surface-raised` `#141A24` both collapse to
`black`.** The surface ramp flattens to nothing.

### 7.1 · The binding rule

> **Elevation is carried by borders and spacing. Never by fill.**
> No region, panel, pane, row or selection may be distinguished from its neighbour by background
> colour. Design so that **no layout depends on surface fill to separate regions.**

This is why Zerado does not paint a background at all (`01-design-system.md` §1.3), why the
focused row is marked by a gutter glyph plus bold plus the ring rather than by a highlight bar,
and why the detail pane is separated by a **border and a two-column gutter** rather than by a
raised fill.

### 7.2 · The two sanctioned fills

The confirmation overlay and the command palette may paint `--z-surface-overlay` for modality —
but **must also carry a `--z-border-strong` border**, so they remain legible when the fill
vanishes at 16 colours.

### 7.3 · The check

Render the screen at 16 colours. **If two regions become one, the screen fails.**

---

## 8 · Borders and the focus ring — two rules that are not negotiable

### 8.1 · `--z-border` is decoration and may never mark a control

`--z-border` `#2A3342` measures **1.53:1**. It may draw a hairline between blocks of content. It
may **never** mark the edge of an input, a button, a dialog, a filter chip, or any control whose
boundary carries meaning. Those use **`--z-border-strong` `#64748B`**, which measures **4.08**
and satisfies WCAG 1.4.11 (Non-text Contrast, 3:1, AA).

**The check:** find every box on the screen. If it contains or *is* a control, its border must be
`--z-border-strong`. A `--z-border` box around an input is an automatic fail.

### 8.2 · The focus ring is never removed

`--z-focus-ring` `#19E0FF` (ANSI **45** / `bright cyan`), **on every interactive element, in
every state, on every screen.** Not in read-only panes, not in "preview" modes, not while
syncing, not for aesthetics.

Brand §4.2 states the reason plainly: a keyboard-first audience is the launch audience, and
removing focus rings for aesthetics would be a self-inflicted wound on exactly the people the
product is for. In a terminal there is **no pointer at all** — the ring is not an accessibility
courtesy, it is the only way to know where you are.

Focus is carried by three channels so it survives `NO_COLOR` and 16 colours: the `▌` gutter
marker (position), bold (weight), and colour — `--z-primary` amber on a ledger row cursor,
`--z-focus-ring` cyan on a focused control. Any two of the three are enough.

---

## 9 · Changing any of this

Colour changes go through brand governance (§10), not through a screen PR:

1. Change the token in **`tokens.css` and `tokens.json`** — same commit.
2. Update the affected section of the **brand manual** — same commit.
3. If it touches a colour used for text, **recompute the contrast ratio.** Do not estimate it.
4. If it touches a **state** colour, **re-run the CVD simulation on all six pairs.** The
   **ΔE 11.9** minimum (zerado × abandoned under deuteranopia) is the floor to protect.
5. If it touches a colour used in the terminal, **re-derive the ANSI-256 index** by
   nearest-neighbour search in CIELAB — never adjust the old one by hand.

Owner: `fft-brand-architect`. A screen may not introduce a colour; it may only spend one.

---

## 10 · The reviewer checklist — literal pass/fail

Run against the founder-facing rendered artifact. Every line is pass or fail; there is no
"mostly".

**Cyan**
1. ☐ Chrome-cyan count, by §3, is **0 or 1**.
2. ☐ Every remaining cyan mark is a `ZERADO` state cell, the focus ring, or the text cursor.
3. ☐ No cyan on any border, heading, title, divider, readout label, or progress bar.
4. ☐ The word `zerado` in the summary row is **not** cyan.

**Amber**
5. ☐ Every amber mark is on the §4.1 allow-list.
6. ☐ Amber cells are **≤ 10 %** of the viewport.
7. ☐ Body text is not amber.

**Red**
8. ☐ Red appears **only** as the scanner pip, a destructive annunciator, or error text.
9. ☐ `--z-scanner` (5.25) sets **no words**; error text uses `--z-scanner-300` (6.99) or the
   documented interim (uncoloured + bold).
10. ☐ Offline is **not** red.
11. ☐ The scanner is on screen **only** during a genuinely indeterminate wait — never ambient.

**Structure**
12. ☐ No region is separated from another by background fill.
13. ☐ Rendered at 16 colours, no two regions merge.
14. ☐ Every control boundary is `--z-border-strong`; `--z-border` marks no control.
15. ☐ The focus ring is present on the focused element — in this state, on this screen.
16. ☐ Focus survives `NO_COLOR`: the `▌` gutter marker and bold still identify it.

**Co-render**
17. ☐ Every state shows **colour and glyph and label**.
18. ☐ With `NO_COLOR=1`, **no information is lost** — same screen, same meaning, zero SGR.
19. ☐ No copy refers to a colour, a shape or a position (`"the cyan one"`, `"the row on the
    right"`) — WCAG 1.3.3.

**Any unchecked box is a failed review.** Not a note, not a nit — a fail. That is what makes
this a budget rather than an opinion.

---

## 11 · The failure gallery — what a violation actually looks like

Real ways this goes wrong, so they are recognised rather than argued about.

| # | The screen does this | Verdict |
|---|---|---|
| 1 | Section headings in cyan "because they're important" | **Fail** — chrome cyan × 4. Headings are amber (readout) or chrome |
| 2 | The detail pane bordered in cyan to "tie it to the selection" | **Fail** — cyan is never a border |
| 3 | Every success message in cyan | **Fail** — cyan means *completion of a game*, not "this went fine" |
| 4 | A cyan sync progress bar | **Fail** — progress is amber; cyan is not a process colour |
| 5 | `6 zerado` in the summary row rendered cyan | **Fail** — prose about state is not a state cell (§2.4) |
| 6 | `OFFLINE` in red | **Fail** — offline is a kept promise, rendered chrome (§5.2) |
| 7 | Panes separated by a raised background fill | **Fail** — invisible at 16 colours (§7) |
| 8 | The Steam-key input outlined in `--z-border` | **Fail** — 1.53:1 on a control; violates WCAG 1.4.11 |
| 9 | Focus ring dropped in the "read-only" detail pane | **Fail** — the ring is never removed (§8.2) |
| 10 | A scanner sweeping under the header as an ambient flourish | **Fail** — the scanner is for indeterminate waits, never decoration |
| 11 | Amber title, amber labels, amber values, amber hints, amber borders | **Fail** — the ceiling exists to catch exactly this |
| 12 | `ABANDONED` rendered red instead of orchid | **Fail** — the states are ratified and CVD-verified; red would also collapse against the scanner |

---

## 12 · Open for the founder

1. **The one-per-screen chrome-cyan budget** is my reading of *"the single most important call to
   action on a screen"* made countable. Confirm the number is **1** and not, say, 2 for a screen
   with a genuine primary and secondary action.
2. **The two-class model** (§2) — state cyan unbounded, chrome cyan budgeted — is a design
   decision this document introduces. It is what makes the rule survive a library where 40 games
   are finished. Confirm it reads correctly against the brand's intent.
3. **The 10 % amber ceiling** (§4.2) translates brand §4.1's 60/30/10 to a rendered viewport with
   blanks in the denominator. Confirm that reading.
