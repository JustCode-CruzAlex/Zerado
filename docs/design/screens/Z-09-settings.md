---
title: Zerado — Z-09 Settings
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-09
rev: B
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-09 · Settings

> Fills [`../03-designer-manual.md`](../03-designer-manual.md) §3's 16-section contract.
> Composition binding from [`../../blueprint/02-composition.md`](../../blueprint/02-composition.md) §2 —
> single-pane grouped form, `R = 1`.
> **Every current value is visible without opening anything**, and **the credential backing and the
> audio backing are both stated honestly** — a property the player cannot see is one they cannot
> rely on ([`06-data-seams.md`](../../blueprint/06-data-seams.md) §5.4 ·
> [`12-audio.md`](../../blueprint/12-audio.md) §10).
>
> **Rev B** adds one row — `DISPLAY § Images` — and moves the arithmetic that hangs off it. The row
> is not a nicety: it is the **only durable home** for *why this terminal has no cover art* once
> [`Z-15`](./Z-15-cover-deck.md) §5.4's note is dismissed and `v` retires. §17 item 2a, §10.4.

---

## 1 · Identity

| | |
|---|---|
| **Screen** | `Z-09` · Settings |
| **Phase** | 1 |
| **Kind** | Route |
| **Routes in** | `,` from anywhere |
| **Routes out** | `Z-02 Connect a store` (`⏎` on a Steam row) · `Z-03 Sync` (`⏎` on `Last sync`) · pop on `Esc` |
| **Offline class** | **WORKS** ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2) — *"Every dial reads and writes locally. Connecting a store is `Z-02`, a separate screen with its own class — Settings only routes there"*. The **radio** has its own class in the same table — it streams, so offline it stops — and this screen is where that is stated, in the `Radio` row's own value (§10.3). The screen does not change class because one of the things it reports about does |
| **Writes** | `setting(key, value)` ([`09-erd.md`](../../blueprint/09-erd.md) §1). Never a credential — that is the `Vault`'s |

---

## 2 · Purpose

**Show every dial the product has, grouped, with the current value visible without opening
anything.**

Screen inventory §5: *"`Z-09` must not hide a current value behind a submenu."*

---

## 3 · Mockups at 80 columns

Frame row map as [`Z-01-first-run.md`](./Z-01-first-run.md) §3.1 — content at **column 4**,
body `74 × 16`.

### 3.1 · The setting row — column budget at 74

| Field | Cols | Range | Note |
|---|---|---|---|
| focus gutter | **2** | 4–5 | `▌`, width-aware padded |
| label | **18** | 6–23 | longest is `Interface sounds` = 16 |
| gap | 2 | 24–25 | |
| value | **52** | 26–77 | **every value in §10 fits.** The longest is `~/.local/share/zerado/library.db — will be created` = **50**; the longest in the `DISPLAY` group is `Off, no Kitty graphics or iTerm2 images` = **39**. 52 is chosen to hold all of them at Wide, which is why no row wraps at this tier |
| | **74** | | |

Group heads sit at **column 4**, in the gutter's own column, because a head is not focusable.
Group notes are indented **2** — to column 6, the label column — because they are about the group,
not about a row.

### 3.2 · RENDER 80×24 — everything connected, audio on

**23 content rows in a 15-row scroll region.** Row 16 of the body is the position line, pinned
outside the scroll — the same discipline R-10(c) puts on the ledger's summary.

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Settings

   SETTINGS

   STEAM
     Account             76561198012345678
   ▌ API key             In the OS keychain
     Last sync           3 hours ago, 247 games
     Your key is never inside library.db — that file is safe to copy.

   AUDIO
     Audio               On
     Output              CoreAudio, built-in output
     Radio               On
     Radio volume        60
     Interface sounds    On
     FX volume           80

   DISPLAY
   ▄ 1–15 of 23

   ↑↓ move   ⏎ replace   d disconnect   esc back   ? help   q quit

```

### 3.3 · RENDER 80×24 — scrolled to the end

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Settings

   SETTINGS

     Output              CoreAudio, built-in output
     Radio               On
     Radio volume        60
     Interface sounds    On
     FX volume           80

   DISPLAY
     Glyphs              Auto (Unicode)
     Motion              Auto (on)
   ▌ Colour              On, 256 colours
     Images              On, Kitty graphics protocol

   LIBRARY
     File                ~/.local/share/zerado/library.db
     Size                412 KB, 247 games
   ▄ 9–23 of 23

   ↑↓ move   esc back   ? help   q quit

```

`Colour` is focused and the footer offers **no verb for it** — because it has none. A footer that
offered `⏎ change` on a read-only row would be a footer that lies.

### 3.4 · RENDER 80×24 — nothing connected, audio off, no audio device

The honest screen. This is the state a player reaches from `Z-01` with `,`.

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Settings

   SETTINGS

   STEAM
   ▌ Not connected       ⏎ to connect
     No key is stored.

   AUDIO
     Audio               Off
     Output              No audio device. Zerado is silent.
     Radio               On
     Radio volume        60
     Interface sounds    On
     FX volume           80
     Audio is off. These are kept, and take effect when you turn it on.

   DISPLAY
     Glyphs              Auto (Unicode)
   ▄ 1–15 of 22

   ↑↓ move   ⏎ connect   esc back   ? help   q quit

```

The `⏎ connect` hint is where this screen's **one** chrome cyan is spent, and it is the only state
in which any is spent at all (§14).

### 3.5 · RENDER 80×24 — the disconnect confirmation

Destructive, per [`../01-design-system.md`](../01-design-system.md) §13.1 item 2 —
*"Disconnecting a store (drops its synced rows)."* The overlay uses §13.2's anatomy at its drawn
width; **the default is the safe action** and there is no pre-selected `y`.

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Settings

   SETTINGS

   STEAM
     Account             76561198012345678
   ▌ API key             In the OS keychain
   ┌────────────────────────────────────────────┐
   │ ▌ DISCONNECT STEAM                         │
   │                                            │
   │   247 synced games leave the library.      │
   │   Games you added by hand stay. Your key   │
   │   is deleted from the OS keychain.         │
   │                                            │
   │   y  disconnect            esc  keep       │
   └────────────────────────────────────────────┘
     FX volume           80

   DISPLAY
   ▄ 1–15 of 23

   y disconnect   esc keep

```

The box is `--z-border-strong` — a control boundary, so 1.4.11 applies and **4.08** satisfies it;
`--z-border` (1.53) may never do this job. It may paint `--z-surface-overlay` for modality but
**must not depend on it**: at the 16-colour floor the fill vanishes and the border is what is left
(§1.3, §7.2). It is sized and positioned so it **does not entirely obscure the focused row**
(2.4.11): the box occupies body rows **4–12**, so `▌ API key` — the row `d` was pressed on — sits
**directly above it, with its marker intact**, and the `STEAM` head sits above that. Below the box,
body rows 13–15 keep `FX volume`, the gap and the `DISPLAY` head. The confirmation is therefore
adjacent to its subject rather than on top of it, which is what makes 2.4.11 checkable by looking
instead of by assertion (§13).

---

## 4 · Mockup at `120 × 40`

`leftInset` **4** · body `112 × 32` · content at **column 5** · `1 + 3 + 1 + 32 + 1 + 1 + 1 = 40` ✓

**All 23 rows fit with nine body rows to spare — no scrolling at all.** 23 of 32 body rows, so the
extra height is not merely sufficient, it is unpressured; a `Station` row (§17) or a sixth `DISPLAY`
readout would still fit. That is the whole progressive enhancement of the extra height, and it is
why **the position line disappears when there is nothing to position** — the render below draws the
empty tail rather than cropping, because "room to spare" is a claim a mockup should show.

```text
0.........1.........2.........3.........4.........5.........6.........7.........8.........9.........A.........B.........

    Zerado ✦ Settings

    SETTINGS

    STEAM
      Account             76561198012345678
    ▌ API key             In the OS keychain
      Last sync           3 hours ago, 247 games
      Your key is never inside library.db — that file is safe to copy.

    AUDIO
      Audio               On
      Output              CoreAudio, built-in output
      Radio               On
      Radio volume        60
      Interface sounds    On
      FX volume           80

    DISPLAY
      Glyphs              Auto (Unicode)
      Motion              Auto (on)
      Colour              On, 256 colours
      Images              On, Kitty graphics protocol

    LIBRARY
      File                ~/.local/share/zerado/library.db
      Size                412 KB, 247 games










    ↑↓ move   ⏎ replace   d disconnect   esc back   ? help   q quit

```

### 4.1 · "Group ∥ values" at ExtraWide, read honestly — *design decision*

[`03-responsive.md`](../../blueprint/03-responsive.md) §3 gives `Z-09` *"Group ∥ values"* at
ExtraWide. Read as a **master–detail split** — group names on the left, the selected group's values
on the right — it would put every value of every unselected group behind a selection, and screen
inventory §5 forbids exactly that: *"must not hide a current value behind a submenu."*

Read as **the label column beside the value column**, which is what the row budget already is, it
is satisfied at every tier and the extra width is spent on the thing that actually helps: **all 23
rows on screen at once.**

This spec takes the second reading. A second reason is `1.3.2 Meaningful Sequence` — two
side-by-side *groups* would interleave in byte order, so anything consuming the output stream would
read `Account 76561198012345678   Audio   On` as one line. One column, top to bottom, is the
terminal's only honest structure. **Flagged in §17** so `fft-tui-architect` can confirm the reading
or restate the table.

---

## 5 · Visual hierarchy

**The one thing the player must see first: the group they came for, with its current values already
readable.**

| Rank | Element | Channel | Note |
|---|---|---|---|
| 1 | `SETTINGS` | case + weight + `--z-primary` | the H1 |
| 2 | **the group heads** — `STEAM` `AUDIO` `DISPLAY` `LIBRARY` | case (UPPER) + weight (bold) + `--z-primary` + **position** (flush at column 4, outdented past every row) | the outdent is what makes this a *grouped* form rather than a list. It costs no rows |
| 3 | **the focused row** | `▌` gutter (position) + bold (weight) + amber (colour) | |
| 4 | the **values** | `--z-text` **255** — the brightest text on screen | this screen exists to show values, so the values outrank their labels |
| 5 | the labels | `--z-text-secondary`, sentence case | deliberately *not* uppercase — see §5.1 |
| 6 | group notes | `--z-text-secondary`, indented to the label column | one per group, at most |
| 7 | position line, breadcrumb, footer | chrome | |

### 5.1 · Labels here are sentence case, not the UPPERCASE readout — *design decision*

`Z-02` and `Z-08` set their field labels UPPERCASE, per the readout role
([`../01-design-system.md`](../01-design-system.md) §1.5). `Z-09` does not, and the reason is
countable: this screen has **four** uppercase group heads and **fifteen** labels. Uppercasing all
nineteen would leave nothing to outrank anything, and the group heads — which are the screen's
actual navigation — would stop reading as heads.

*(The 23 rows of §3.2 are not 23 labels: four are group heads, one is a group note and three are
the gaps between groups. Fifteen rows carry a label. An earlier revision of this paragraph
subtracted the four heads from the row total and called the remainder "labels", which counted the
note and the gaps.)*

The readout role is kept where it is load-bearing: **the group heads carry it**, in amber, and they
are the section heads §1.5 assigns amber to. The rows below them are label–value pairs, and a pair
reads better when the value is the loud half.

---

## 6 · Every applied spacing token, by name

| Token | Wide value | Where Z-09 spends it |
|---|---|---|
| `OuterMarginX` | **2** | frame inset |
| `OuterMarginY` | **1** | rows 1 and 24 |
| `InnerPaddingX` | **1** | inside the frame |
| `InnerPaddingY` | **1** | row 22 |
| `InterElementGap` | **1** | breadcrumb→title · band→body · **between groups** · scroll region→position line |
| `HeaderBandHeight` | **3** | `hasSubtitle = false` |
| `leftInset` | **3** | header-left **==** content-left at column 4 |
| group-note indent | `+2` cols → column 6 | the note aligns with the labels it is about, not with its head |
| scroll region | `BodyRect.h − 1` = **15** rows at Wide, **31** at ExtraWide | the position line is pinned outside it |

**Groups are separated by one `InterElementGap` and by the head's outdent — never by a rule and
never by a fill.** A `--z-border` hairline between groups would be decoration at 1.53:1 that
vanishes at 16 colours, and a background fill is forbidden outright: **elevation is carried by
borders and spacing, never by fill** (§1.3, §7.1).

---

## 7 · Colour, glyph and label for every state shown

| Element / state | Token | Hex | ANSI-256 | 16-colour | Ratio | Glyph / structure | Label | `NO_COLOR` |
|---|---|---|---|---|---|---|---|---|
| screen title | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | — | `SETTINGS` | **yes** — case |
| group head | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | — | `STEAM` etc., UPPER + bold, outdented | **yes** — case + outdent |
| row label | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** | — | `API key` | **yes** |
| row value | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** | — | `In the OS keychain` | **yes** |
| **focused row** | `--z-primary` on the marker | `#FFB000` | **214** | `bright yellow` | **10.59** | `▌` U+258C (ASCII `>`) + **bold row** | — | **yes** — marker **and** bold |
| read-only row | identical to a focusable row | | | | | **no glyph difference** | the **footer** offers it no verb | **yes** — the footer is the channel |
| group note | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** | — | the sentence | **yes** |
| position line | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** | `▄` U+2584 — **Ambiguous**, flowing, must be width-measured | `1–15 of 23` | **yes** |
| confirmation border | `--z-border-strong` | `#64748B` | **67** | `bright black` | **4.08** — meets 1.4.11 | `┌ ─ │ ┘` | — | **yes** — box drawn |
| confirmation annunciator | `--z-scanner` | `#FF2E2E` | **9** | `bright red` | 5.25 — **structure, not text** | `▌` | `DISCONNECT STEAM` | **yes** — `▌` + case |
| write-error annunciator | `--z-scanner` | `#FF2E2E` | **9** | `bright red` | 5.25 | `▌` | `NOT SAVED` | **yes** |
| write-error text | **uncoloured + bold** — documented interim (§11.2) | `#FF6B6B` | *underived* | `bright red` | **6.99** | — | the sentences | **yes** |
| audio indicator | from [`../01-design-system.md`](../01-design-system.md) §5.2 | `#FFB000` / `#A9B5C7` | **214** / **249** | | **10.59** / **9.36** | `▮` / `▯` — **Neutral** | `AUDIO` / `MUTED` | **yes** — filled vs hollow |

### 7.1 · How "read-only" co-renders without a glyph — *design decision*

**Six** rows on this screen are pure readouts — `Account`, `Output` when it is a failure sentence,
`Colour`, **`Images`**, `File` and `Size`. They must be distinguishable from actionable rows
**without colour**, and the product has no spare glyph.

The channel is **the footer**, which is already the screen's honest statement of what works right
now (nav §6). Focus a read-only row and the verbs disappear from the footer; focus an actionable
one and they are there. That is a text channel, it survives `NO_COLOR` and a screenshot, and it
requires no new vocabulary.

**Rejected:** dimming the row (colour-only, fails 1.4.1), a padlock or bullet glyph (every narrow
candidate is either taken or Ambiguous-width — §7.2 of `Z-01`), and skipping read-only rows in the
traversal (makes `↑↓` feel broken, and the player still needs to reach a long value to read it).

---

## 8 · The full state table

| # | State | Trigger | Renders | Copy |
|---|---|---|---|---|
| 1 | **First run** — nothing connected, empty library, reached from `Z-01` with `,` | | §3.4 | `Not connected` / `No key is stored.` · `Size` reads `0 games` and `File` reads the path **that will be created**, marked as such — §10 |
| 2 | **Default** — connected, audio on and available | | §3.2 | as drawn |
| 3 | **Credential backing = keychain** | `Vault.Backing() == "keychain"` | `API key   In the OS keychain` | §10 |
| 4 | **Credential backing = file** | `Vault.Backing() == "file"` | `API key   In credentials.json, mode 0600` · the group note changes too | §10 |
| 5 | **Key missing from the vault** | connection row exists, `Vault.Get` misses | `API key   Not found. ⏎ to add it again.` | §10 |
| 6 | **Never synced** | `LastSync(steam) == nil` | `Last sync   Never` | §10 |
| 7 | **Sync was partial** | `last_sync_status == 'partial'` | `Last sync   3 hours ago, partially — 138 of 247` | §10 |
| 8 | **Audio off** *(the default)* | `Enabled == false` | §3.4's AUDIO group + its note | `Audio is off. These are kept, and take effect when you turn it on.` |
| 9 | **Audio unavailable** — five reasons | `AudioState.Reason` | the `Output` row states it, **once**, and nothing else changes | §10 · the `Output` matrix |
| 10 | **`ZERADO_NO_AUDIO` set** | env | `Audio   On, overridden by ZERADO_NO_AUDIO` — **overridden, not off**, because those are different facts ([`12-audio.md`](../../blueprint/12-audio.md) §10) | §10 |
| **10a** | **Radio on, no network** | the stream cannot open | The `Radio` row reads `On. No network, so no music.` **The switch still says the player's choice; the sentence says what is happening.** No banner, no cue, no retry key — a stream that ends when the network does is behaving correctly | §10.3 |
| 11 | **Build has no audio player** | default build, no `-tags audio` | the AUDIO group renders **its head and one line, and no dials** — §8.1 | `This build was made without audio. Nothing here can make a sound.` |
| 12 | **`NO_COLOR` set** | env | `Colour   Off, NO_COLOR is set` — read-only, and it says **who** set it | §10 |
| 13 | **`ZERADO_ASCII=1` set** | env | `Glyphs   ASCII, forced by ZERADO_ASCII` | §10 |
| 14 | **`ZERADO_REDUCED_MOTION` set** | env | `Motion   Reduced, forced by ZERADO_REDUCED_MOTION` | §10 |
| **14a** | **Images resolved** | `Images.Capability()` returns `kitty` or `iterm2` ([`17-images.md`](../../blueprint/17-images.md) §4) | `Images   On, Kitty graphics protocol` / `On, iTerm2 inline images` — drawn in §3.3 and §4 | §10.4 |
| **14b** | **No image protocol** | `Capability() == none`, including the fail-closed timeout | `Images   Off, no Kitty graphics or iTerm2 images`. **Nothing else changes** — a terminal without images is a supported configuration, not a degrade ([`17-images.md`](../../blueprint/17-images.md) §3), so there is no banner, no amber and no cue. This row is the **only durable home** for the fact once [`Z-15`](./Z-15-cover-deck.md) §5.4's note is dismissed and `v` retires | §10.4 |
| **14c** | **`ZERADO_NO_IMAGES` set** | env | `Images   Off, ZERADO_NO_IMAGES is set` — the variable is named, like `NO_COLOR` two rows up. It **wins over a real capability**: a terminal that draws images still reads `Off` here, because the row reports what Zerado resolved, not what the terminal could do | §10.4 |
| 15 | **Disconnect confirmation** | `d` on a Steam row | §3.5 | §10 · **D** |
| 16 | **Setting write failed** | `SetSetting` errors | message block below the group; **the dial reverts to its stored value**, it does not show a change that did not happen | §10 · **W** |
| 17 | **Offline** | network off | **nothing changes.** `Z-09` is `WORKS`; only `Z-02`, which it routes to, is not | — |
| 18 | **Loading** | first paint | **N/A.** Every value is a local read already resolved before the route is pushed. `AudioState` is read from the seam, which never blocks ([`12-audio.md`](../../blueprint/12-audio.md) §6) | — |
| 19 | **Empty / partial** | **N/A** — the group set is fixed and compiled in | — | — |
| 20 | **`NO_COLOR` render** | env | §12 | identical text |
| 21 | **Below `24 × 8`** | | never renders — [`Z-11-fatal-error.md`](./Z-11-fatal-error.md) | the refusal sentence |

### 8.1 · A dial that cannot do anything is not shown — *design decision*

When the binary was built without the audio player, the five audio dials can never affect anything.
Rendering them would be **anti-pattern 14** — *"A field shown for a capability that doesn't exist
yet. Omit the block; don't label it empty."*

So state 11 renders:

```text
   AUDIO
     This build was made without audio. Nothing here can make a sound.
```

and nothing else. Note the difference from state 9, where the dials **are** shown: there, audio
exists and a device is missing, so the settings are real and will take effect the moment one
appears. *Absent capability* and *absent device* are different facts and get different screens.

### 8.2 · Audio never raises a banner, and this is the one place it speaks — *applied, not decided*

[`../01-design-system.md`](../01-design-system.md) §15.5 is explicit: *"Audio never blocks, never
warns, and never fails a run. A missing sound device is not a degrade worth a banner… Reporting it
would be noise about noise."*

That is compatible with stating it here, and the distinction is worth naming because it looks like
a contradiction: **§15.5 forbids a banner; `Z-09` is a readout.** The player who wonders *"why is
it silent?"* comes here and gets one sentence. Nothing follows them around. The rule this screen
obeys is *say it once, in the place someone would look* — the same rule the credential backing
obeys two groups above.

---

## 9 · The key map

No text input exists on `Z-09` — every value is edited in place with `⏎` or `←→` — so **every
single-key shortcut is live**, and the footer can be honest row by row.

| Key | Does | In the footer? |
|---|---|---|
| `↑` `↓` · `k` `j` | move between rows, **skipping group heads and notes** | **yes** — `↑↓ move` |
| `g` / `G` | first / last row | no |
| `Ctrl-D` / `Ctrl-U` | half page | no |
| `⏎` | the focused row's action, if it has one — `replace` · `sync now` · `connect` · `toggle` · `look again` | **only when the focused row has one** |
| `←` `→` · `h` `l` | adjust the focused row — volume **± 5**, or cycle a choice | **only on an adjustable row** — `←→ volume` / `←→ change` |
| `d` | **disconnect the focused store** — raises the confirmation (§3.5). *Design decision, §17* | **only on a `STEAM` row while connected** — `d disconnect` |
| `y` | confirm, **only while the confirmation is open** | **yes, only then** |
| `esc` | close the confirmation → otherwise **pop the route** | **yes** — `esc keep` while confirming, `esc back` otherwise |
| `,` | **unbound here.** Pushing a route already on the stack unwinds to it (nav §1 rule 3), so `,` on `Z-09` would be a key that does nothing visible | no |
| `?` | Help | **yes** |
| `q` · `Ctrl-C` | quit | `q quit` |
| `m` | toggle mute | **only when audio is enabled** |
| `s` `a` `/` `r` | **unbound** — they belong to `Z-04` | no |

**Footer per focused row** — the footer is the hint block, and on this screen it is also the
read-only channel (§7.1):

| Focused row | Footer |
|---|---|
| `Account`, `File`, `Size`, `Colour`, `Images`, `Output` (when it is a failure sentence) | `↑↓ move   esc back   ? help   q quit` |
| `API key` | `↑↓ move   ⏎ replace   d disconnect   esc back   ? help   q quit` |
| `Last sync` | `↑↓ move   ⏎ sync now   d disconnect   esc back   ? help   q quit` |
| `Not connected` | `↑↓ move   ⏎ connect   esc back   ? help   q quit` |
| `Audio`, `Radio`, `Interface sounds` | `↑↓ move   ⏎ toggle   esc back   ? help   q quit` |
| `Radio volume`, `FX volume` | `↑↓ move   ←→ volume   esc back   ? help   q quit` |
| `Glyphs`, `Motion` | `↑↓ move   ←→ change   esc back   ? help   q quit` |
| confirmation open | `y disconnect   esc keep` |

---

## 10 · The exact copy — ready to paste

### 10.1 · Chrome and groups

| Slot | String |
|---|---|
| breadcrumb | `Zerado ✦ Settings` |
| title | `SETTINGS` |
| group heads | `STEAM` · `AUDIO` · `DISPLAY` · `LIBRARY` |
| position line | `▄ 1–15 of 23` — `of 22` when nothing is connected (§3.4), `of 38` at Narrow and `of 40` at Tiny, where a row is two lines (§11) |

### 10.2 · `STEAM`

| Row | Value |
|---|---|
| `Account` | `76561198012345678` |
| `Account`, not connected | the row is replaced by `Not connected` / `⏎ to connect` |
| `API key`, keychain | `In the OS keychain` |
| `API key`, keychain, named | `In the macOS Keychain` · `In the GNOME keyring` · `In Windows Credential Manager` — whichever `Vault.Backing()` reports |
| `API key`, file | `In credentials.json, mode 0600` |
| `API key`, missing | `Not found. ⏎ to add it again.` |
| `Last sync` | `3 hours ago, 247 games` |
| `Last sync`, never | `Never` |
| `Last sync`, partial | `3 hours ago, partially — 138 of 247` |
| group note, keychain | `Your key is never inside library.db — that file is safe to copy.` |
| group note, file | `Your key is beside library.db, not inside it. Copy the .db freely.` |
| group note, not connected | `No key is stored.` |

**Why the note is there at all.** *"Settings shows the backing because a security property the
player cannot see is a security property they cannot rely on"*
([`06-data-seams.md`](../../blueprint/06-data-seams.md) §5.4). The published promise is that the
library file is *"a single SQLite file you can back up, move, or delete"* — this line is what makes
that safe to act on.

### 10.3 · `AUDIO`

| Row | Value |
|---|---|
| `Audio` | `On` · `Off` |
| `Audio`, env override | `On, overridden by ZERADO_NO_AUDIO` |
| `Output`, available | whatever `AudioState.Backing` reports, e.g. `CoreAudio, built-in output` |
| `Output`, `no_device` | `No audio device. Zerado is silent.` |
| `Output`, `ssh` | `No audio over SSH. Zerado is silent.` |
| `Output`, `ci` | `No audio on this machine. Zerado is silent.` |
| `Output`, `env_disabled` | `Silent. ZERADO_NO_AUDIO asked for no sound.` |
| `Output`, `init_failed` | `The audio device did not start. Zerado is silent.` |
| `Radio` · `Interface sounds` | `On` · `Off` |
| **`Radio`, on with no network** | **`On. No network, so no music.`** — 28 cells. Not an error, not amber, and it raises nothing anywhere else |
| `Radio volume` · `FX volume` | `60` · `80` — a number, `0` to `100` |
| group note, audio off | `Audio is off. These are kept, and take effect when you turn it on.` |
| group note, no audio in the build | `This build was made without audio. Nothing here can make a sound.` |

**Each of these is said once, in Settings, and never anywhere else.** No banner, no repeated
notice, no error a screen has to handle — `Cue()` cannot fail because it has no error return
([`12-audio.md`](../../blueprint/12-audio.md) §4).

**Radio and Interface sounds are independent, deliberately.** *"Someone may want the keyclicks
without the soundtrack, or the soundtrack without the keyclicks, and neither is the odd request"*
([`12-audio.md`](../../blueprint/12-audio.md) §1). They are never collapsed into one switch — which
is also what satisfies **WCAG 1.4.2 Audio Control** three times over: opt-in only, independently
mutable, independently volumed.

**They are also independent about the network, and that is the point of the split.** The radio
**streams**, so with no connection there is no music — and that is fine, not a fault, and it
raises nothing. Interface sounds are local cues and keep working offline like everything else in
this product. **Nothing is bundled**, so there is no fallback track and nothing to apologise for:
the row simply says what is true. Founder direction, 2026-08-25 —
*"let's skip the bundle music, if the user is offline no music, that's it."*

**Volume is a number, not a bar** — *design decision.* The count is the information and the bar is
the ornament ([`../01-design-system.md`](../01-design-system.md) §8.1); a bar here would also be
amber outside the allow-list, and four little bars on a settings screen is a light show, not a
cockpit.

### 10.4 · `DISPLAY`

| Row | Value |
|---|---|
| `Glyphs` | `Auto (Unicode)` · `Auto (ASCII)` · `Unicode` · `ASCII` |
| `Glyphs`, env override | `ASCII, forced by ZERADO_ASCII` |
| `Motion` | `Auto (on)` · `On` · `Reduced` |
| `Motion`, env override | `Reduced, forced by ZERADO_REDUCED_MOTION` |
| `Colour` — **read-only** | `On, truecolor` · `On, 256 colours` · `On, 16 colours` |
| `Colour`, `NO_COLOR` set | `Off, NO_COLOR is set` |
| `Images` — **read-only** | `On, Kitty graphics protocol` **27** · `On, iTerm2 inline images` **24** |
| `Images`, no protocol | `Off, no Kitty graphics or iTerm2 images` **39** |
| `Images`, env override | `Off, ZERADO_NO_IMAGES is set` **28** |

**The four values are the three `Capability()` results and the override**, named exactly as
[`17-images.md`](../../blueprint/17-images.md) §2 names them — *Kitty graphics*, *iTerm2 inline
images* — because **the protocol name is this row's whole job**.
[`Z-15`](./Z-15-cover-deck.md) §5.3 settles that split: the dismissible note carries the two
**program** names, since it is an interruption and must land in the second it is read; the places a
player goes **looking** carry the **protocol**, because that is the term they would have to search
to understand why one terminal draws covers and another does not. This is one of those places, so
it says the searchable thing.

**Three vocabulary rules this row inherits.** It never says *support* — that is the vocabulary of a
compatibility error and this is not one; it never says *your terminal*, which makes an absence the
player's fault; and it never opens with a word like `NO`, which reads an absence as a failure. All
three are [`Z-15`](./Z-15-cover-deck.md) §5.3's rejections and
[`17-images.md`](../../blueprint/17-images.md) §3's tone rule — *"a terminal without image support
is a supported configuration"* — applied to a readout instead of a note.

> **One proposed value was not taken, and it is worth saying which.**
> [`Z-15`](./Z-15-cover-deck.md) §17 item 6 sketched this row's values as `Kitty graphics protocol`
> · **`Not supported by this terminal`** · `Off, ZERADO_NO_IMAGES is set`. The middle one is the
> sentence `Z-15`'s own §5.3 rejects two paragraphs earlier — *support* is compatibility-error
> vocabulary and *this terminal* makes an absence the player's setup's fault. A proposal in a
> remediation table is a lead, not a ruling; the rejection table outranks it. The value here is
> `Off, no Kitty graphics or iTerm2 images` — no verdict on the terminal, just the two things
> Zerado looked for and did not find.

**The env form is `is set`, not `forced by`** — the same distinction `DISPLAY` already draws.
`Glyphs` and `Motion` say `forced by`, because the environment is overriding a **dial the player
owns**. `Colour` and `Images` say `is set`, because there is no dial to override; the variable is
simply the reason the readout reads what it reads. Same grammar, different fact.

**`Images` is required, not optional, and it is the only durable home for this fact.**
[`Z-15`](./Z-15-cover-deck.md) §5.4's recommendation note retires once dismissed, and `v` retires
with it on a terminal that draws no covers. `Z-10 Help` is **generated from the key registry**, so a
retired key cannot appear there — which means that after dismissal there is nowhere else a player
can learn why they have no cover art. This row is that place. It is read-only for the same reason
`Colour` is: the capability belongs to the terminal, and Zerado's job is to say what it resolved.

That is also why the row is a **row** and not a group note. A note is about its group
(§3.1) and would be dropped the moment `DISPLAY` gained a second thing worth saying; a row is part
of the compiled-in set (§8 state 19) and cannot be dropped. The one durable fact in the product
gets the one structure that cannot go missing. **Flagged in §17**, and checkable at §16 item 15.

**`Colour` is read-only in Phase 1** — *design decision.* `NO_COLOR` belongs to the environment and
Zerado's job is to obey it; the depth belongs to the terminal. An in-app colour switch is a
capability nobody has specified, and a dial for it would claim something unbuilt (anti-pattern 14).
What the row does instead is the honest and more useful thing: **say what Zerado actually
resolved**, so a player debugging a screenshot can see it. **Flagged in §17.**

### 10.5 · `LIBRARY`

| Row | Value |
|---|---|
| `File` | `~/.local/share/zerado/library.db` — the resolved path, tilde-shortened for the home directory |
| `File`, before it exists | `~/.local/share/zerado/library.db — will be created` |
| `File`, `$ZERADO_DB` set | the override path, plus `, set by ZERADO_DB` |
| `Size` | `412 KB, 247 games` |
| `Size`, empty library | `4 KB, 0 games` |

### 10.6 · **D** · the disconnect confirmation

```text
   ▌ DISCONNECT STEAM

     247 synced games leave the library.
     Games you added by hand stay. Your key is deleted from the
     OS keychain.

     y  disconnect            esc  keep
```

Names **what happens**, **what survives**, and **what happens to the secret**. The default is the
safe action; `esc` keeps; there is no pre-selected `y`
([`../01-design-system.md`](../01-design-system.md) §13.2).

### 10.7 · **W** · a setting could not be written

```text
   ▌ NOT SAVED

     Zerado could not write to the library file.
     ~/.local/share/zerado/library.db: read-only file system

     The setting is unchanged. Fix the file and try again.
```

**The dial reverts to its stored value.** A settings screen that shows a change it failed to make
is lying about the state of the machine, and this screen's whole job is to not do that.

**Voice check.** No exclamation marks · no emoji · never "gamer" · the number is said · nothing
contradicts `landing-copy.md` §08 (one file, no account) or §14 (`Nothing about your library is
sent to a Zerado-run server, because there isn't one`).

---

## 11 · 40-column behaviour, and the refusal floor

### 11.1 · RENDER 40×24 — Narrow · `leftInset` 2 · body `36 × 16`

Responsive table: *"Values below labels."* The 18 + 2 + 52 budget does not fit in 36, so the row
becomes **two lines** — label on the first, value indented on the second — the same trade the game
row makes at this tier, and for the same reason: **at 40 columns a player is reading, not
scanning.**

```text
0.........1.........2.........3.........

  Zerado ✦ Settings

  SETTINGS

  STEAM
    Account
      76561198012345678
  ▌ API key
      In the OS keychain
    Last sync
      3 hours ago, 247 games
    Your key is not in library.db.

  AUDIO
    Audio
      On
    Output
      CoreAudio, built-in output
    Radio
  ▄ 1–15 of 38

  ↑↓ move  ⏎ replace  esc  ?

```

**The value line has 32 columns** — the label starts at column 5, the value is indented two more to
column 7, and the body ends at column 38. A value longer than 32 **wraps**; it is never truncated. A
setting whose value is cut off is a setting the player cannot verify, which is this screen's one
failure mode.

**Two-line rows grow the list from 23 rows to 38 lines**, which the position line states honestly.
The arithmetic: **15** label-and-value rows at two lines each = 30, plus the **4** group heads, the
**1** group note and the **3** gaps between groups at one line each = 8. In the state drawn above no
value exceeds 32, so nothing wraps a third line — the longest here is `~/.local/share/zerado/library.db`
at exactly **32**, and `Images` reads `On, Kitty graphics protocol` at **27**. On a terminal with no
image protocol that value is **39** and takes a third line, which makes the list **39**; the position
line says which. And
the group notes take their **short form** at Narrow and below — the same long-form/short-form
descriptor pattern [`Z-01-first-run.md`](./Z-01-first-run.md) §10 records, for the same reason: a
hand-set short sentence breaks where a person would break it, and a greedy wrapper breaks
mid-clause.

| Note | Standard and above | Narrow | Tiny |
|---|---|---|---|
| `STEAM` | `Your key is never inside library.db — that file is safe to copy.` | `Your key is not in library.db.` | `Key is not in library.db.` |
| `AUDIO`, off | `Audio is off. These are kept, and take effect when you turn it on.` | `Off. Settings are kept.` | `Off. Settings are kept.` |

### 11.2 · Standard `60 × 24` · body `54 × 16`

Responsive table: *"Same"* as Wide. Label 18 + gap 2 + value **32**.

**Ten of §10's values exceed 32 and wrap to a second line at this tier.** In descending order:
`~/.local/share/zerado/library.db — will be created` **50** · `The audio device did not start.
Zerado is silent.` **49** · `No audio on this machine. Zerado is silent.` **43** ·
`Silent. ZERADO_NO_AUDIO asked for no sound.` **43** · `Reduced, forced by ZERADO_REDUCED_MOTION`
**40** · `Off, no Kitty graphics or iTerm2 images` **39** · `No audio over SSH. Zerado is silent.`
**36** · `3 hours ago, partially — 138 of 247` **35** · `No audio device. Zerado is silent.` **34** ·
`On, overridden by ZERADO_NO_AUDIO` **33**. Every one of the remaining values fits.

**These are all failure and override sentences, and that is the whole point:** the values that
overflow are the ones that had something to explain, and at Standard they are given a second line
rather than a truncation. Each wrap adds one line to the list, which the position line states.
Stated as arithmetic rather than as a judgement.

> **Correction, recorded rather than quietly fixed.** An earlier revision of this section named
> `On, overridden by ZERADO_NO_AUDIO`, 33, as *"the longest value"* and said *"every other value
> fits."* Both were false at the time they were written — nine other values were already over 32,
> the longest by eighteen columns. Adding `Images` did not create the defect; measuring for `Images`
> found it. The lesson is the one this bundle keeps relearning: **measure the string, do not trust
> the sentence beside it.**

### 11.3 · Tiny `< 40` — `32 × 24` · body `30 × 21`

Responsive table: *"Values below labels."* Same two-line row, band collapsed to the title.

```text
0.........1.........2.........3.
 SETTINGS

 STEAM
   Account
     76561198012345678
 ▌ API key
     In the OS keychain
   Last sync
     3 hours ago, 247 games
   Key is not in library.db.

 AUDIO
   Audio
     On
   Output
     CoreAudio, built-in output
   Radio
     On
   Radio volume
     60
   Interface sounds
     On
 ▄ 1–20 of 40
 ↑↓  ⏎  q
```

At Tiny the list is **40 lines** in a 20-line region and the position line says so. **Nothing is
dropped and nothing is truncated** — the screen simply becomes longer, which is the correct trade
for a surface whose entire job is showing values.

**Where the 40 comes from.** The value line has **26** columns here — label at column 4, value
indented two more to column 6, body ending at column 31. Fifteen label-and-value rows at two lines
each = 30, plus 4 heads + 1 note + 3 gaps = 8, plus **two third lines** for the two values that
exceed 26: `~/.local/share/zerado/library.db` at **32** and `Images`' `On, Kitty graphics protocol`
at **27**. 30 + 8 + 2 = **40**. The `Images` value misses by a single column and still wraps rather
than truncate — the rule is the rule, and Tiny is the tier where it earns its keep.

At Tiny the **confirmation becomes a route**, not an overlay — the 46-column box does not fit
inside `30 × 21` with its margins ([`02-composition.md`](../../blueprint/02-composition.md) §2.4).
Behaviourally identical; only the composition changes.

### 11.4 · The refusal floor — below `24 × 8`

`Z-09` never renders; see [`Z-11-fatal-error.md`](./Z-11-fatal-error.md) §3.4.

---

## 12 · `NO_COLOR` rendering — shown, not asserted

Zero SGR; bold goes with it. The characters are unchanged:

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Settings

   SETTINGS

   STEAM
     Account             76561198012345678
   ▌ API key             In the OS keychain
     Last sync           3 hours ago, 247 games
     Your key is never inside library.db — that file is safe to copy.

   AUDIO
     Audio               On
     Output              CoreAudio, built-in output
     Radio               On
     Radio volume        60
     Interface sounds    On
     FX volume           80

   DISPLAY
   ▄ 1–15 of 23

   ↑↓ move   ⏎ replace   d disconnect   esc back   ? help   q quit

```

| Information | Channel that survives |
|---|---|
| which row has focus | the `▌` gutter — **position** |
| that `STEAM` is a group head | it is **outdented two columns past every row** and it is UPPERCASE — position and case |
| which rows are actionable | **the footer**, which lists the focused row's verbs and lists none for a readout (§7.1) |
| where the key is stored | the sentence, in words, in two places — the value and the group note |
| why it is silent | the `Output` row's sentence |
| **why there is no cover art** | the `Images` row's value, which **names the protocol** — the searchable term, in words, on the only screen that still carries it once `v` retires (§10.4) |
| that a setting is overridden by the environment | the value **names the variable** — `NO_COLOR`, `ZERADO_ASCII`, `ZERADO_NO_AUDIO`, `ZERADO_NO_IMAGES` |
| that there is more below | `▄ 1–15 of 23` |

**No information is lost.** And with `ZERADO_NO_AUDIO=1` set as well, nothing is lost — the screen
is where audio's whole story is told in text.

---

## 13 · Focus model on this screen

| | |
|---|---|
| **Regions** | **1** — the settings list. `Tab` is unbound |
| **Focusable items** | every **row**. Group heads and group notes are skipped by the traversal |
| **Initial focus** | the first focusable row of the first group |
| **Traversal** | `↑` `↓` / `k` `j`, no wrap at the ends (a settings list with a bottom should feel like it has one) |
| **Scrolling** | **cursor-following** — the focused row is always visible, and the position line always says where it is. The same discipline R-10(b) puts on the ledger |
| **Focus is never nowhere** | the row set is compiled in; the only row that can appear or disappear is the `STEAM` group's shape, and when it changes focus moves to the nearest surviving row (nav §4.1 invariant 3) |
| **Focus ring** | present on the focused row in **every** state, including while the confirmation is open — the row behind it keeps its marker, which is what makes 2.4.11 checkable |
| **`Esc`** | confirmation open → **cancel, the safe branch, always**. Otherwise **pop the route** |
| **Restored on pop** | returning from `Z-02` or `Z-03` puts the player back on the row they left |

**Focus is trapped inside the confirmation while it is open** — which is correct, because there is
a documented way out (`esc`), and it is what 2.1.2 asks for.

---

## 14 · Colour budget declaration

| Class | Count | Where |
|---|---|---|
| **CHROME cyan** | **0 normally · exactly 1 when nothing is connected** | `⏎ connect` in the footer (§3.4). *"A screen may spend zero. Cyan is earned; a screen with nothing to urge does not need to urge anything"* ([`../02-colour-budget.md`](../02-colour-budget.md) §2.2). Settings is a readout — until there is no store, and then connecting one is unambiguously the most important thing on it |
| **STATE cyan** | **0** | no state chip renders on `Z-09` |
| **Focus-ring cyan** | **0** | the row cursor is **amber**, not cyan (§1.7), so `Z-09` emits no cyan at all in its normal state |
| **Amber** | 5 marks, +1 with audio on | title (allow-list **1**) · four group heads (**2**, section heads) · `▌` row cursor (**5**) · `▮ AUDIO` (**9**) |
| **Red** | **0 normally** · **1** while the confirmation or a write error is showing | the `▌` annunciator (§5 list **2**) |
| **Error text red** | **0** | uncoloured + bold, the documented interim |

**Amber ceiling:** title 8 + heads `STEAM` 5 + `AUDIO` 5 + `DISPLAY` 7 + `LIBRARY` 7 + cursor 1 =
**33 cells** of `80 × 24 = 1920` → **1.7 %**. At ExtraWide, where all four heads are on screen at
once: 33 of `120 × 40 = 4800` → **0.7 %**.

**The temptation this screen resists.** `Audio  On` and `Radio  On` are the kind of value that
invites a "success" colour. They are `--z-text`, like every other value. Cyan means *a game is
finished*, not *this switch is up* — failure-gallery item 3.

---

## 15 · Reuse verdict per element

| Element | Verdict | Note |
|---|---|---|
| The grouped list | **Build fresh** — `bubbles/viewport` for the scroll region + `lipgloss` for rows | `huh` groups are **rejected here**: `huh` is a *form* — it collects a set of answers and submits them. `Z-09` has no submit; every row commits on the keystroke, which is the product's rule everywhere (nav §3.2, *"every mutation is committed when it is made"*). Wrapping that in a form primitive would invent a save step that must not exist |
| The scroll region | **`bubbles/viewport`** — direct fit | the same primitive the ledger uses, so cursor-following scroll is written once |
| Position line | **Build fresh** — one `lipgloss` row pinned outside the viewport | R-10(c)'s discipline, applied to a non-ledger |
| Row cursor | **Build fresh** — the shared 2-column gutter | §1.7, identical to `Z-01`, `Z-08` |
| Confirmation | **`huh` — fits**, restyled to Zerado tokens | [`../01-design-system.md`](../01-design-system.md) §13.5. **`huh`'s default theme carries its own palette and must not ship** |
| Volume control | **Build fresh** — a number and two keys | **not** `bubbles/progress` (truecolor gradient, §8.5) and **not** a slider widget. §10.3 |
| Header band, footer, audio indicator | **Build fresh**, shared across the cluster | §2.8 · §5.2 |
| Audio state | **read from the `Audio` seam's `State()`** ([`12-audio.md`](../../blueprint/12-audio.md) §4) — never from a device query on the render path | `State()` is a struct read; device init happened once, off the render path, with a timeout |
| Audio cue | **none on this screen.** Toggling audio on does **not** play a demonstration cue | the cue list is closed (§15.3) and a settings change is not on it. A screen that makes a noise to prove it can make a noise is exactly the ambient decoration the brand rules out |

---

## 16 · Screen-specific acceptance criteria

1. **No value is behind a submenu.** Every value in §10 is readable by scrolling and pressing
   nothing. Grep the render at `120 × 40`: all **23** rows present, no position line, no scrolling
   required — and, because the tail is drawn, nine empty body rows after `Size`.
2. **The credential backing is stated and correct.** Force `Vault.Backing()` to `keychain` and to
   `file`; the `API key` row and the group note both change, and both name the real location.
3. **The audio backing is stated and correct.** Force each `AudioState.Reason` — `none`,
   `no_device`, `ssh`, `ci`, `env_disabled`, `init_failed` — and assert the `Output` row's sentence.
4. **`ZERADO_NO_AUDIO` shows as overridden, not off** (§10.3), and `NO_COLOR` shows as `Off,
   NO_COLOR is set`, naming the variable.
4a. **The radio states the network honestly and raises nothing.** With `Radio` on, cut the
   network: assert the row reads `On. No network, so no music.`, that the switch still reads the
   player's choice, that `Interface sounds` is unaffected and still fires, and that **no banner
   appears on any screen** — grep every render for `radio`, `stream` and `music`. `Z-09` is the
   only place it is mentioned.
5. **A build without the audio player shows no dials** (§8.1) — the AUDIO group is a head and one
   sentence.
6. **Audio raises no banner, on any screen** — grep every other screen's render for `audio`,
   `silent`, `device`. `Z-09` is the only place it is mentioned.
7. **The footer never offers a verb the focused row does not have** (§7.1). Walk every row and
   diff the footer against the table in §9.
8. **A failed write reverts the dial** (§10.7). Force `SetSetting` to fail and assert the rendered
   value is the stored one, not the attempted one.
9. **Disconnect is confirmed, and the default is safe.** `d` raises the box; `esc` keeps; there is
   no pre-selected `y`; the copy names what leaves, what stays and what happens to the key.
10. **Marking a game `ABANDONED` raises no confirmation anywhere** — the closed destructive list is
    three items and that is not one of them (§13.1). Checked here because `Z-09` owns the only
    confirmation in this cluster.
11. **The confirmation does not entirely obscure the focused row** (2.4.11), and it is legible at
    the 16-colour floor where its fill vanishes and only the `--z-border-strong` border remains.
12. **No value is ever truncated at any tier** (§11.1) — it wraps.
13. **Chrome-cyan count is 0 in the default state and exactly 1 in the not-connected state.**
14. Eight artifacts per [`03-responsive.md`](../../blueprint/03-responsive.md) §7, **plus**
    `80 × 24` for: not-connected, `file` backing, each audio-unavailable reason, the confirmation,
    and **`Capability() == none`**.
15. **The `Images` row survives dismissal, and it is the only thing that does.** Force
    `Capability()` to `kitty`, `iterm2` and `none`, and set `ZERADO_NO_IMAGES`; assert the four
    values of §10.4 and that the row is read-only in all four (no verb in the footer, §7.1). Then
    the durable-home test, which is the one that matters: on a `none` terminal, write
    `setting('covers.note_dismissed','true')`, confirm `v` is unbound and **absent from `Z-10`**
    (which is generated from the key registry and therefore cannot list it), and assert that
    `Z-09 § DISPLAY` still names the protocol. **Grep the whole render set for `Kitty graphics`:
    after dismissal, `Z-09` is the only screen that has it.** This is the check that would have
    caught [`Z-15`](./Z-15-cover-deck.md) §5.4's cross-reference before it shipped.
16. **The row count is 23, and every readout of it agrees.** Grep every render for `of ` and diff
    against §10.1: `of 23` connected · `of 22` not connected · `of 38` at Narrow · `of 40` at Tiny.
    A position line is a claim about a list, and it is the one claim on this screen that can go
    stale without anything else looking wrong.

---

## 17 · Open for the founder

1. **`d` disconnects (§9).** It is a new binding — the global map does not mention `d`, and the TUI
   manual's R-9 warns against inheriting FlowForge's `f`/`d` as cargo cult. `d` here is Zerado's
   own, on one row, behind a confirmation. Confirm, and let `fft-tui-architect` fold it into
   `04-navigation-and-focus.md` §3.
2. **`Colour` is read-only (§10.4).** Zerado obeys `NO_COLOR` and the terminal's reported depth and
   reports what it resolved. Confirm — or specify an in-app colour switch, which would be a new
   capability rather than a new row.
2a. **`DISPLAY` has gained an `Images` row (§10.4), which is what took this spec to rev B.** It is
   the change [`Z-15`](./Z-15-cover-deck.md) §17 item 6 deferred as *"a founder-visible change to a
   GOLDEN spec rather than a sweep"*, and it is no longer optional: `Z-15` §5.4 originally claimed
   the no-covers fact stayed available in `Z-10 Help`, and it cannot — `Z-10` is generated from the
   key registry, so when `v` retires on dismissal the fact retires with it. **This row is the only
   durable home left.** Two things to confirm, and one that is not open:

   **(a) The row itself, and its cost.** 22 rows became 23, which moved every position readout on
   the screen and grew the Narrow list to 38 lines and the Tiny list to 40. Nothing broke: 23 still
   fit `120 × 40` with nine body rows to spare, no tier gained a scrollbar it did not have, and the
   value fits the Wide row budget without wrapping. The arithmetic is in §3.1, §4, §11.1 and §11.3
   and is checkable at §16 item 16.

   **(b) The wording — the protocol, not the programs.** The row says `Kitty graphics protocol`
   where the dismissible note says *"Ghostty or Kitty"*. That split is `Z-15` §5.3's, and
   [`Z-15`](./Z-15-cover-deck.md) *Open for the founder* item 2 puts it to you: **the interruption
   carries the program names, the place you go
   looking carries the searchable term.** If you decide the note should name the protocol too, or
   that this row should name the programs, this row's four values change with it — they are §10.4
   and nothing else depends on them.

   **Not open: that it is read-only.** It is the same case as item 2. There is no dial because the
   capability is the terminal's; the only lever a player has is `ZERADO_NO_IMAGES`, which is an
   environment variable by design ([`17-images.md`](../../blueprint/17-images.md) §2 —
   *"never a config flag the player must find"*).
3. **"Group ∥ values" at ExtraWide (§4.1)** is read as *label column beside value column*, not as a
   master–detail split, because a split would hide values behind a selection. Confirm the reading.
4. **The Audio group's row set (§10.3)** is six rows: the opt-in, the honest backing, two channel
   switches and two volumes. The **bundled-versus-user-supplied music question is closed** —
   founder direction, 2026-08-25, made the music a **stream**, so there is no local directory to
   point at and no seventh `Music folder` row. *(An earlier revision of this
   spec flagged `12-audio.md` §1, §3 and §7 as still describing a bundled subsystem. **They have
   since been repaired** — §3 is now headed *"Streamed, never bundled"* and §7 *"Stations are data —
   and the licensing question is closed"*. The flag is withdrawn.)*

   **What is genuinely open is the station.** The direction names *"synthwave / 80s stations"* —
   plural — and this spec has **not** invented a picker, because a control nobody asked for is
   the anti-pattern this bundle rejects. If the player is to choose, it is a seventh row,
   `Station`, with the same row shape and a value like `Synthwave` cycled by `←→`, and this spec
   revises to rev B. If Zerado picks one station and never says so, the row set stays at six.
5. **Labels are sentence case on this screen and UPPERCASE on the two forms (§5.1).** The reason is
   countable — nineteen uppercase strings would flatten the group heads. Confirm the exception.
6. **The ratified network promise now has to cover a radio stream — and that is genuinely open.**
   *"The only network traffic is `Zerado` reaching out to the services you've connected"*
   ([`../03-designer-manual.md`](../03-designer-manual.md) §5.9) is a **published** line. An opt-in
   station is defensibly a service the player connected — they turned it on and can stop it in one
   keystroke — but the sentence names Steam and price data specifically, and a station uses **no key
   of the player's**. Defensible is not the same as ratified, so this is **routed to the founder
   rather than reinterpreted by a designer**, and it is on the gate list
   ([`13-handoffs.md`](../../blueprint/13-handoffs.md) §5).

   *(This item previously also reported `12-audio.md` as stale and assigned an owner. That upstream
   was repaired and **the flag outlived its repair** — it was reaching the founder as an open item
   directing a specialist at a defect that no longer existed. Withdrawn; only the promise question,
   which is real, remains.)*
