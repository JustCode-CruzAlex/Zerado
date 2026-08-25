---
title: Zerado — Z-02 Connect a store
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-02
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-02 · Connect a store

> Fills [`../03-designer-manual.md`](../03-designer-manual.md) §3's 16-section contract.
> Composition binding from [`../../blueprint/02-composition.md`](../../blueprint/02-composition.md) §2 —
> single-pane `huh` form, `R = 1`.
> **The screen is rendered from `Capabilities().Credentials`**
> ([`06-data-seams.md`](../../blueprint/06-data-seams.md) §2.3), not written per store. §3.5 is the
> element-by-element proof.

---

## 1 · Identity

| | |
|---|---|
| **Screen** | `Z-02` · Connect a store — **not** "connect Steam" ([`01-screen-inventory.md`](../../blueprint/01-screen-inventory.md) §3.3) |
| **Phase** | 1 · Steam is the only Phase 1 **instance** |
| **Kind** | Route |
| **Routes in** | `Z-01` door 1 · `Z-09 Settings` → the Steam row |
| **Routes out** | `Z-03 Sync` on success · pop on `Esc` |
| **Offline class** | **REFUSES** ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2) — *"Cannot validate a key without reaching the provider. Says so on submit, keeps what was typed"* |
| **Never routed here** | a provider whose `Capabilities().Credentials` is empty — `physical` has none, so there is nothing to connect |

---

## 2 · Purpose

**Take the player's own credentials for one provider, explain each field where it is entered, and
fail with the actual reason.**

Screen inventory §5: *"`Z-02` must not tell the player their profile is private before it has
evidence — it is the empty result that proves it, and the copy belongs there."*

---

## 3 · Mockup at 80 columns — the design floor

Frame row map as [`Z-01-first-run.md`](./Z-01-first-run.md) §3.1: `1 + 3 + 1 + 16 + 1 + 1 + 1 = 24`,
`80 − 2 × 3 = 74`, content begins at **column 4**.

### 3.1 · The field row — column budget at 74

| Field | Cols | Range | Note |
|---|---|---|---|
| focus gutter | **2** | 4–5 | `▌` U+258C, width-aware padded |
| gap | 1 | 6 | |
| left boundary | **1** | 7 | `│` unfocused → `┃` focused |
| gap | 1 | 8 | content never touches the boundary |
| **text** | **67** | 9–75 | 66 characters plus the caret's cell |
| gap | 1 | 76 | |
| right boundary | **1** | 77 | |
| | **74** | | |

**Fields are full body width, uniformly** — *design decision.* A pasted key of unknown length must
be visible without horizontal scroll, and two fields of two different widths read as an accident
rather than a form.

### 3.2 · RENDER 80×24 — default, nothing typed, field 1 focused

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Connect a store

   CONNECT A STORE

   STEAM

   STEAM API KEY
   ▌  ┃ ▎                                                                   ┃
     Free, and about a minute to make: steamcommunity.com/dev/apikey

   STEAM ID
      │                                                                     │
     Your profile URL works too — Zerado reads the ID out of it.

   Your key goes to the OS keychain, never into library.db.




   tab next   ⏎ connect   esc back   ? help   q quit

```

### 3.3 · RENDER 80×24 — checking (the one indeterminate wait on this screen)

The **scanner** ([`../01-design-system.md`](../01-design-system.md) §9), one row, spanning the body
width. Exactly one on screen — never two ([`03-responsive.md`](../../blueprint/03-responsive.md) §5).
The field help lines yield; the fields and their contents do not.

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Connect a store

   CONNECT A STORE

   STEAM

   STEAM API KEY
      │ ********************************                                    │

   STEAM ID
   ▌  ┃ 76561198012345678                                                   ┃

   CHECKING WITH STEAM
   ───────────────────━━━────────────────────────────────────────────────────

   Nothing is saved yet. esc cancels.

   esc cancel   ? help   q quit

```

### 3.4 · RENDER 80×24 — the key was rejected

The message block replaces the field help. **The fields keep everything the player typed** — the
offline contract's own words for this screen.

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Connect a store

   CONNECT A STORE

   STEAM

   STEAM API KEY
   ▌  ┃ ********************************                                    ┃

   STEAM ID
      │ 76561198012345678                                                   │

   ▌ STEAM KEY REJECTED
     Steam rejected that key. Check it hasn't been regenerated at
     steamcommunity.com/dev/apikey.
     Nothing was saved. The key you typed is still in the field.


   tab next   ⏎ try again   esc back   ? help   q quit

```

### 3.5 · Every rendered element, and the descriptor member it comes from

**Nothing on this screen switches on `ProviderID`.** This table is the spec's proof of that.

| Rendered element | Source | Steam's value in Phase 1 |
|---|---|---|
| the screen exists at all | `len(Capabilities().Credentials) > 0` | 2 fields |
| body block 1, the section head | `Provider.Display()` | `Steam` → rendered `STEAM` |
| number of field blocks | `len(Credentials)` | **2** |
| field order top to bottom | slice order | key, then ID |
| field 1 label | `Credentials[0].Label` | `Steam API key` → rendered `STEAM API KEY` |
| field 1 masked | `Credentials[0].Secret == true` | masked with `*` (§7.3) |
| field 1 destination | `Secret == true` → the **`Vault`** | never `library.db` |
| field 1 help | `Credentials[0].Help` | `Free, and about a minute to make:` |
| field 1 help URL | `Credentials[0].HelpURL` | `steamcommunity.com/dev/apikey` |
| field 1 inline error | `Credentials[0].Validate(s)` | see §8.2 |
| field 2 label | `Credentials[1].Label` | `Steam ID` → `STEAM ID` |
| field 2 unmasked | `Credentials[1].Secret == false` | plain |
| field 2 destination | `Secret == false` → `provider_connection.account_ref` — *"an identifier, never a secret"* ([`09-erd.md`](../../blueprint/09-erd.md) §1) | the Steam ID |
| field 2 help | `Credentials[1].Help` | `Your profile URL works too — Zerado reads the ID out of it.` |
| the vault line | `Vault.Backing()` | `keychain` \| `file` — §10 |
| the refusal copy | the classifier ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §5), **not** the provider | §8 |

**Adding GOG is: implement the interfaces, declare the fields, register.** This screen changes by
zero characters. §16 line 1 makes that testable.

---

## 4 · Mockup at `120 × 40`

`leftInset` **4** · body `112 × 32` · content begins at **column 5** ·
`1 + 3 + 1 + 32 + 1 + 1 + 1 = 40` ✓

The responsive table's ExtraWide column for `Z-02` is **"Field + inline help beside it"**
([`03-responsive.md`](../../blueprint/03-responsive.md) §3). The field keeps a **74-column** budget
(identical to Wide) and the help moves into a **34-column** column beside it, separated by a
4-column gutter: `74 + 4 + 34 = 112` ✓

```text
0.........1.........2.........3.........4.........5.........6.........7.........8.........9.........A.........B.........

    Zerado ✦ Connect a store

    CONNECT A STORE

    STEAM

    STEAM API KEY
    ▌  ┃ ▎                                                                   ┃    Free, and about a minute
                                                                                  to make:
                                                                                  steamcommunity.com/dev/apikey

    STEAM ID
       │                                                                     │    Your profile URL works too
                                                                                  — Zerado reads the ID out
                                                                                  of it.

    Your key goes to the OS keychain, never into library.db.
```

**1.3.2 Meaningful Sequence, checked.** Two side-by-side columns interleave in byte order, so
side-by-side is only sanctioned where each fragment is self-contained. Here the field and its help
are **one label's worth of content**, and the help never carries a fact the field's label does not
already introduce. *No two independent fields are ever placed side by side on Z-02* — that would
interleave two unrelated inputs and is the failure the SC exists to prevent.

---

## 5 · Visual hierarchy

**The one thing the player must see first: the focused field.**

| Rank | Element | Channel | Note |
|---|---|---|---|
| 1 | `CONNECT A STORE` | case + weight + `--z-primary` | the screen's H1 |
| 2 | **the focused field** | `▌` gutter (position) + **heavy** `┃` boundary (weight) + `--z-focus-ring` cyan (colour) | three channels, any two sufficient |
| 3 | `STEAM` | case + weight + `--z-primary` | *which* store — the one word that changes per provider |
| 4 | field labels | UPPERCASE + `--z-text-secondary` | the readout role (§1.5) |
| 5 | field help | sentence case + `--z-text-secondary`, indented to the boundary column | subordinate by position |
| 6 | the vault line | `--z-text-secondary`, one `InterElementGap` below | a fact, not an instruction |
| 7 | breadcrumb, footer | chrome | |

**When a message is showing, it takes rank 2** — the `▌` annunciator plus an UPPERCASE heading is
the strongest mark on the screen, and it is directly below the field it is about.

### 5.1 · Two `▌` on one screen, and how they stay apart — *design decision*

`▌` is the focus marker (§1.7) **and** the error annunciator (§11.1). On `Z-02` they can coexist.
They are separated by three channels, and **two of them survive `NO_COLOR`**:

| | focus marker | error annunciator |
|---|---|---|
| colour | `--z-primary` amber **214** | `--z-scanner` red **9** |
| what follows it | a **boundary glyph** and a field | an **UPPERCASE heading** |
| position | on a field row | in the message block, below a gap |

They never appear on the same row. Recorded because it is the sort of collision that is obvious
in a spec and invisible in a review.

---

## 6 · Every applied spacing token, by name

| Token | Wide value | Where Z-02 spends it |
|---|---|---|
| `OuterMarginX` | **2** | frame inset |
| `OuterMarginY` | **1** | rows 1 and 24 |
| `InnerPaddingX` | **1** | inside the frame — **and again inside the field boundary**, which is why content never touches `│` |
| `InnerPaddingY` | **1** | row 22 |
| `InterElementGap` | **1** | breadcrumb→title · band→body · `STEAM`→field 1 · field-1 block→field-2 block · field 2→vault line · fields→message block |
| `HeaderBandHeight` | **3** | `hasSubtitle = false` |
| `leftInset` | **3** | header-left **==** content-left at column 4 |
| ExtraWide help gutter | **4** cols | §4 — `74 + 4 + 34 = 112` |

**A field block is label · field · help, with no gap inside it** — *design decision.* `huh` gives
each field a title/field/help triple with its own breathing room; the three lines are **one
object** and an internal gap would break the pairing. The `InterElementGap` goes *between* blocks,
which is exactly what the token is for.

---

## 7 · Colour, glyph and label for every state shown

Values read from [`../01-design-system.md`](../01-design-system.md) §1.4. Nothing estimated.

| State | Token | Hex | ANSI-256 | 16-colour | Ratio | Glyph / structure | Label | `NO_COLOR` |
|---|---|---|---|---|---|---|---|---|
| screen title | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | — | `CONNECT A STORE` UPPER + bold | **yes** — case |
| provider head | `--z-primary` | `#FFB000` | **214** | `bright yellow` | **10.59** | — | `STEAM` UPPER + bold | **yes** — case |
| field label | `--z-text-secondary` | `#A9B5C7` | **249** | `white` | **9.36** | — | e.g. `STEAM API KEY` | **yes** |
| field value | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** | — | the typed text | **yes** |
| **field, unfocused** | `--z-border-strong` | `#64748B` | **67** | `bright black` | **4.08** — meets 1.4.11 | `│` U+2502 **light** | — | **yes** — boundary present |
| **field, focused** | `--z-focus-ring` | `#19E0FF` | **45** | `bright cyan` | **12.15** | `┃` U+2503 **heavy** + `▌` gutter in `--z-primary` **214** | label rendered **bold** | **yes** — heavy weight + `▌` |
| text caret | the terminal's own | — | — | — | — | drawn `▎` in mockups | — | **yes** |
| **scanner track** | `--z-scanner-track` `#5C1414` — **ANSI index underived** → **interim: uncoloured** ([`../01-design-system.md`](../01-design-system.md) §9.4) | — | *underived* | `black` | — | `─` U+2500 | `CHECKING WITH STEAM` above it | **yes** — light stroke |
| **scanner pip** | `--z-scanner` | `#FF2E2E` | **9** | `bright red` | 5.25 — **motion, not text** | `━` U+2501 × **exactly 3** | — | **yes** — **heavy** stroke is the pip |
| **message annunciator** | `--z-scanner` | `#FF2E2E` | **9** | `bright red` | 5.25 — structure, not text | `▌` U+258C | — | **yes** |
| **message heading** | **uncoloured + bold** — the documented interim, because `--z-scanner-300` `#FF6B6B` (**6.99** AA) has no derived ANSI index (§11.2) | `#FF6B6B` | *underived* | `bright red` | **6.99** | — | `STEAM KEY REJECTED` etc. | **yes** — case + `▌` |
| message body | `--z-text` | `#E9EEF5` | **255** | `bright white` | **16.65** | — | the sentences | **yes** |
| audio indicator | see [`Z-01-first-run.md`](./Z-01-first-run.md) §7.1 — applied unchanged from [`../01-design-system.md`](../01-design-system.md) §5.2 | | | | | `▮` / `▯` (both **Neutral**) | `AUDIO` / `MUTED` | **yes** |

### 7.3 · The mask character is ASCII `*` — *design decision, and it is a width decision*

`•` U+2022, `●` U+25CF and `·` U+00B7 are all **East-Asian-Ambiguous** (verified with
`unicodedata`, Unicode **16.0.0**, 2026-08-25). A 32-character mask of an ambiguous glyph is
**64 cells** on a double-ambiguity terminal — it would push the field's right boundary clean off
the row. The mask is the one run on this screen whose length is attacker-controlled by the
player's paste, so it must be immune: **`*`, `Na — Narrow`, immune by construction**, the same
argument that gives `[ ] [~] [*] [x]` its escape-hatch role (§1.2).

### 7.4 · Boundary weight carries focus, so focus survives `NO_COLOR`

`--z-border-strong` (**4.08**) on every field boundary satisfies WCAG 1.4.11; `--z-border`
(**1.53**) never marks a control ([`../02-colour-budget.md`](../02-colour-budget.md) §8.1). The
focused field additionally swaps to the **heavy** box-drawing weight — the same mechanism
[`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §4.2 uses for a focused
pane (`┏━┓` against `┌─┐`), applied to a control.

---

## 8 · The full state table

| # | State | Trigger | Renders | Copy |
|---|---|---|---|---|
| 1 | **First run** — reached from `Z-01` door 1 | | §3.2, both fields empty, focus on field 1 | as drawn. `Esc` pops back to `Z-01` |
| 2 | **Re-connect** — a connection already exists | from `Z-09` | §3.2 **plus** a line under `STEAM`: | `Connected as 76561198012345678. A new key replaces the old one.` |
| 3 | **Typing** | | focus ring on the active field; **no single-key shortcut fires** (2.1.4) | — |
| 4 | **Empty submit** | `⏎` with a field blank | the empty field takes focus; its help row becomes the inline error | `Steam needs this one too.` |
| 5 | **Checking** *(the only indeterminate wait)* | `⏎` | §3.3 — the scanner, one row, help lines yielded, fields intact | `CHECKING WITH STEAM` / `Nothing is saved yet. esc cancels.` |
| 6 | **Cancelled check** | `esc` during 5 | back to §3.2 with every field's contents intact | no message — the player did it on purpose |
| 7 | **Offline — no route / DNS** | classifier ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §5) | message block | §10 · **A** |
| 8 | **Unreachable — timeout / 5xx** | classifier | message block | §10 · **B** |
| 9 | **Key rejected — 401 / 403** | classifier | §3.4 | §10 · **C** |
| 10 | **Empty library — 200 + zero items** | classifier | message block; **credentials are saved**, because they authenticated | §10 · **D** — ratified copy |
| 11 | **Success** | first item arrives | credentials written to the `Vault` and `provider_connection`; **push `Z-03`** | no message — the next screen is the message |
| 12 | **Loading** | **N/A** — `Z-02` renders from the descriptor, which is compiled in. There is nothing to fetch before first paint | — | — |
| 13 | **Empty / partial** | **N/A** — `Z-02` has no data set | — | — |
| 14 | **Audio** | rows identical to [`Z-01-first-run.md`](./Z-01-first-run.md) §8 rows 10–12 | indicator on the reserved footer row; `m` in the key run only when enabled | — |
| 15 | **`NO_COLOR`** | env | §12 | identical text |
| 16 | **Below `24 × 8`** | | never renders — [`Z-11-fatal-error.md`](./Z-11-fatal-error.md) | the refusal sentence |

### 8.1 · How `Z-02` gets its evidence — *design decision*

The inventory forbids claiming a private profile **before there is evidence**. So the submit is not
a bespoke "validate key" call:

> **`⏎` runs one real `Syncer.Sync(ctx, creds)` and reads until the first `Item` or the channel
> closes, then cancels its context.**

| Outcome | Verdict |
|---|---|
| an `Item` arrives | the key works **and** the library is non-empty → save, push `Z-03` (which runs its own sync) |
| the channel closes with **zero** items | 200-and-empty → the **private profile** refusal, **here**, with evidence |
| `error`, classified | the matching refusal, here |

The cost is one extra provider call — for Steam, one HTTP request. The gain is that the player is
never routed to a sync screen only to be refused there, and the ratified copy lands exactly where
the player hit it. **Flagged in §17** because it constrains the seam slightly.

### 8.2 · What `Validate` may reject — *design decision, deliberately conservative*

Steam's `Validate` rejects **empty and whitespace-only, and nothing else.** It asserts no length,
no character class and no format.

A format assertion is a claim about somebody else's API. If it is wrong — or right today and wrong
next year — it locks a player out of a key that works, and the product would be refusing a valid
credential on its own authority. Steam's verdict is the only authority that matters, and §8.1
already goes and gets it. Recorded so nobody adds a regex later as a "polish" pass. Reopening it
needs the format **read at source** and written into the descriptor, not remembered.

---

## 9 · The key map

While a text input holds focus, **every printable key is literal text** (WCAG 2.1.4 · nav §3).
That is the whole reason this table looks different from `Z-04`'s.

| Key | Does | In the footer? |
|---|---|---|
| any printable key | types into the focused field | no |
| `tab` / `shift-tab` · `↓` `↑` | next / previous **field** | **yes** — `tab next` |
| `⏎` | **submit** — from any field | **yes** — `⏎ connect`, and `⏎ try again` once a message is showing |
| `Ctrl-R` | reveal / hide the masked field while it holds focus | **only while the secret field has focus** — `^R show` / `^R hide` |
| `Ctrl-U` / `Ctrl-W` | clear field / delete word — the terminal's own editing keys | no |
| `esc` | cancel a running check → otherwise **pop the route** | **yes** — `esc cancel` while checking, `esc back` otherwise |
| `Ctrl-C` | quit | no |
| `q` | **types `q`.** It does **not** quit while a field has focus | no |
| `?` | **types `?`** while a field has focus. Help is reachable by blurring first — see §13 | **yes** — `? help`, and it is honest because the footer is drawn for the screen, not the field |
| `m` | **types `m`** while a field has focus | only when audio is enabled |
| `,` | **types `,`** | no |
| `r` | **types `r`** — which is why the ratified `r to retry` becomes `⏎ to try again` on this screen (§10) | no |
| `Tab` as a region key | **unbound** — `R = 1` | no |

> **`Ctrl-R` reveals the key** — *design decision.* A masked 32-character paste that Steam rejects
> is otherwise unverifiable, and the player's only recourse is to retype it. A **modified** key is
> mandatory here: a single-key toggle would type a character instead (2.1.4). `Ctrl-R` has no
> conflicting binding inside the app.
>
> **Paste must work** — WCAG 3.3.8. Blocking paste on an API-key field manufactures exactly the
> barrier the criterion exists to prevent.

> **`?` and `q` in the footer while a field has focus.** The footer describes **the screen**, and
> `?`/`q` do work on it — after `esc`. This is the one place the footer's "keys that work here
> right now" rule is read at screen granularity rather than field granularity; the alternative is a
> footer that empties itself the moment anyone types, which is worse. **Flagged in §17.**

---

## 10 · The exact copy — ready to paste

### 10.1 · Chrome and fields

| Slot | String |
|---|---|
| breadcrumb | `Zerado ✦ Connect a store` |
| title | `CONNECT A STORE` |
| provider head | `STEAM` *(from `Display()`, uppercased at render)* |
| already-connected line | `Connected as 76561198012345678. A new key replaces the old one.` |
| field 1 label | `Steam API key` → rendered `STEAM API KEY` |
| field 1 help | `Free, and about a minute to make: steamcommunity.com/dev/apikey` |
| field 1 inline error | `Steam needs this one too.` |
| field 2 label | `Steam ID` → rendered `STEAM ID` |
| field 2 help | `Your profile URL works too — Zerado reads the ID out of it.` |
| field 2 inline error | `Steam needs this one too.` |
| vault line — keychain | `Your key goes to the OS keychain, never into library.db.` |
| vault line — file | `Your key goes to credentials.json (mode 0600), never into library.db.` |
| checking label | `CHECKING WITH STEAM` |
| checking note | `Nothing is saved yet. esc cancels.` |
| footer, default | `tab next   ⏎ connect   esc back   ? help   q quit` |
| footer, secret focused | `tab next   ⏎ connect   ^R show   esc back   ? help   q quit` |
| footer, checking | `esc cancel   ? help   q quit` |
| footer, message showing | `tab next   ⏎ try again   esc back   ? help   q quit` |

> **Verify at source before ship:** `steamcommunity.com/dev/apikey` is the descriptor's declared
> `HelpURL`, reproduced here for the mockup. It is **not** asserted from memory as correct —
> confirm it, and the key's format, at Steam's own documentation, and put the confirmed values in
> the provider descriptor where they belong (§3.5). This spec deliberately asserts **no** key
> length or character class (§8.2).

### 10.2 · The four refusals

Each names **what happened · why · the next action · what happened to the player's data**
([`../01-design-system.md`](../01-design-system.md) §11.3).

**A · no route / DNS**
```text
▌ OFFLINE
  Zerado could not reach Steam, and nothing was sent.
  Everything you typed is still here. ⏎ to try again.
```

**B · timeout / 5xx** — sentence one is the ratified copy, verbatim
```text
▌ STEAM UNREACHABLE
  Steam didn't answer. Not your key — their end, or the connection.
  Nothing was saved. Everything you typed is still here. ⏎ to try again.
```

**C · 401 / 403**
```text
▌ STEAM KEY REJECTED
  Steam rejected that key. Check it hasn't been regenerated at
  steamcommunity.com/dev/apikey.
  Nothing was saved. The key you typed is still in the field.
```

**D · 200 + empty** — the ratified copy, verbatim, plus the data-fate line
```text
▌ STEAM PROFILE PRIVATE
  Steam returned an empty library.
  Game details are private on your profile — Steam won't share the
  list until that's public. Settings → Privacy.
  Your key is saved and works. ⏎ to try again once the profile is public.
```

### 10.3 · Two recorded departures from the ratified refusal strings

| Ratified string | On `Z-02` | Why |
|---|---|---|
| `… r to retry.` | `… ⏎ to try again.` | **`r` cannot be a shortcut on this screen.** A text input has focus, so `r` types `r` — WCAG 2.1.4, and it is anti-pattern 5 in the designer manual §7. Naming a key that would type a letter instead would be a footer that lies, in prose |
| `Steam rejected that key. Check it hasn't been regenerated. Settings → Steam.` | `… regenerated at steamcommunity.com/dev/apikey.` | The ratified pointer sends the player to the screen **they are standing on**. The field's own `HelpURL` is the real next action |
| `No network. Last synced 3 days ago — everything below still works. r to retry.` | **not used here** — copy **A** replaces it | On `Z-02` nothing has synced and there is no library "below". Reproducing it would state a false fact |

`Settings → Privacy` in **D** is kept verbatim because it points at **Steam's** settings, which is
correct and off-product.

**Voice check.** No exclamation marks · no emoji · never "gamer" · every refusal names what, why,
next action and what happened to the data · none says *"Something went wrong."*

---

## 11 · 40-column behaviour, and the refusal floor

### 11.1 · RENDER 40×24 — Narrow · `leftInset` 2 · body `36 × 16`

Responsive table: *"Help collapses to `?` — press to expand."* The field budget becomes
`2 + 1 + 1 + 1 + 27 + 1 + 1 = 34`… which leaves the boundary flush against the body's edge, so at
Narrow **the boundary glyphs are dropped and the field is drawn as an underline instead** — a
control boundary in `--z-border-strong` that costs a row but no columns.

```text
0.........1.........2.........3.........

  Zerado ✦ Connect a store

  CONNECT A STORE

  STEAM

  STEAM API KEY                    ?
  ▌ ▎
    ────────────────────────────────────

  STEAM ID                         ?
    76561198012345678
    ────────────────────────────────────

  tab next  ⏎ connect  esc  ? help

```

- **The `?` beside the label expands that field's help** into the two rows below it, pushing the
  next block down; `esc` collapses it. Nothing is hidden that cannot be reached in one keystroke,
  and the key is on screen.
- The underline is `--z-border-strong` (**4.08**), never `--z-border` (1.53) — it is a control
  boundary and 1.4.11 applies.
- Focus at Narrow: `▌` gutter + **bold value** + the underline in `--z-focus-ring` cyan. Three
  channels still.

### 11.2 · Standard `60 × 24` · body `54 × 16`

Wide's composition unchanged, field text 47 cols. Help stays **below** the field (the responsive
table's Standard column is *"Same"* as Wide).

### 11.3 · Tiny `< 40` — `32 × 24` · body `30 × 21`

Responsive table: *"Same; one field visible at a time."* The band is the title row only; the body
shows **one field block** and a position line, `tab` moves to the next:

```text
0.........1.........2.........3

 CONNECT A STORE

 STEAM

 STEAM API KEY              ?
 ▌ ▎
   ────────────────────────────

 field 1 of 2

 tab next  ⏎ ok  q quit
```

The scanner is **dropped at Tiny**, not shrunk ([`03-responsive.md`](../../blueprint/03-responsive.md) §5).
The checking state at Tiny shows the label alone: `CHECKING WITH STEAM`.

### 11.4 · The refusal floor — below `24 × 8`

`Z-02` never renders. See [`Z-11-fatal-error.md`](./Z-11-fatal-error.md) §3.4.

---

## 12 · `NO_COLOR` rendering — shown, not asserted

Zero SGR. Bold is an SGR sequence and is gone with the rest. The characters are unchanged:

```text
0.........1.........2.........3.........4.........5.........6.........7.........

   Zerado ✦ Connect a store

   CONNECT A STORE

   STEAM

   STEAM API KEY
   ▌  ┃ ********************************                                    ┃

   STEAM ID
      │ 76561198012345678                                                   │

   ▌ STEAM KEY REJECTED
     Steam rejected that key. Check it hasn't been regenerated at
     steamcommunity.com/dev/apikey.
     Nothing was saved. The key you typed is still in the field.


   tab next   ⏎ try again   esc back   ? help   q quit

```

| Information | Channel that survives |
|---|---|
| which field has focus | `▌` gutter **and** the heavy `┃` against the light `│` — two channels, and §1.7 needs two |
| that a control has a boundary | the boundary glyph is drawn, not implied by colour |
| that something failed | `▌` + the UPPERCASE heading + four sentences of plain text |
| that the key is secret | it is masked — masking is not a colour |
| the scanner's pip | the **heavy** `━` against the **light** `─`. This is why the primitive uses stroke weight as its primary channel (§9.7) |

**No information is lost.** `NO_COLOR` also implies reduced motion
([`03-responsive.md`](../../blueprint/03-responsive.md) §5): the scanner pip **parks at the centre
of the track at full weight** and does not travel — `pipLeft = round((74 − 3) / 2) = 36`. It is
deliberately not hidden.

---

## 13 · Focus model on this screen

| | |
|---|---|
| **Regions** | **1** — the field group. `Tab` is a *field* key here, not a region key |
| **Items** | `len(Credentials)` fields, in slice order |
| **Initial focus** | field 1. On re-entry after a rejection, **the field the refusal is about** |
| **Traversal** | `tab` / `shift-tab`, `↓` / `↑`; wraps |
| **During a check** | focus **stays** on the field it was on and **the ring stays drawn** — the ring is never removed, in any state ([`../02-colour-budget.md`](../02-colour-budget.md) §8.2). Input is ignored; `esc` cancels |
| **Focus is never nowhere** | the field list is fixed and non-empty by construction (`len(Credentials) > 0` is the condition for routing here) |

### 13.1 · `Esc` on Z-02 — one press, two contexts

| Context | `Esc` |
|---|---|
| a check is running | **cancel the check.** Every field keeps its contents |
| otherwise | **pop the route** |

**Recorded departure.** [`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §5
has *"a text input has focus → blur the input, keeping what was typed."* On `Z-02` a blur state
would be a **dead intermediate**: `R = 1`, so there is no other region to return to, and there are
no screen-level single-key shortcuts to re-enable — every one of them types a character here (§9).
The player could not tell they were in it, and it would cost a second `Esc` to leave a form that
saves nothing. Blur is kept where it earns its place (`Z-07`'s filter, where blurring re-enables a
real key map). **Flagged in §17** for `fft-tui-architect` to fold into the table.

**No keyboard trap** (2.1.2): `esc` always leaves, from every state including mid-check.

---

## 14 · Colour budget declaration

| Class | Count | Where |
|---|---|---|
| **CHROME cyan** | **1** | `⏎ connect` in the footer — *"the one key hint the screen most wants pressed"* ([`../02-colour-budget.md`](../02-colour-budget.md) §2.2). It becomes `⏎ try again` when a message is showing, and it is **still the same one mark** |
| **Focus-ring cyan** | exempt | the focused field's `┃` boundary — §2.3, singular by definition |
| **Text-cursor cyan** | exempt | the terminal's own caret — §2.3 |
| **STATE cyan** | 0 | no state chip renders on Z-02 |
| **Amber** | 3 marks, +1 with audio on | title (allow-list **1**) · `STEAM` head (**2**) · `▌` focus gutter (**5**) · `▮ AUDIO` (**9**) |
| **Red** | ≤ 2, never at once | the scanner pip (§5 list **1**) *or* the message `▌` annunciator (§5 list **2**, see below). The scanner and a message are **never on screen together** |
| **Error text red** | **0** | `--z-scanner-300` has no derived index → the documented interim: **uncoloured + bold** (§11.2). `--z-scanner` sets no words |

**Amber ceiling:** `CONNECT A STORE` 15 + `STEAM` 5 + `▌` 1 = **21 cells** of `80 × 24 = 1920`
→ **1.1 %**, far under 10 %.

### 14.1 · A gap in `02-colour-budget.md` §5, flagged not worked around

The colour budget's **closed** red list names *"the destructive-confirmation annunciator `▌`"* as
item 2. The **error state**'s `▌` annunciator ([`../01-design-system.md`](../01-design-system.md)
§11.1) is also `--z-scanner`, and it is **not on that list** — so read literally, every error state
in the product fails the budget's own checklist line 8. This spec renders it per §11.1 and
**flags item 2 for amendment** to read *"a destructive-confirmation **or error-state**
annunciator."* See §17.

---

## 15 · Reuse verdict per element

| Element | Verdict | Note |
|---|---|---|
| The form | **`huh` — fits**, restyled to Zerado tokens | [`../01-design-system.md`](../01-design-system.md) §13.5's verdict for confirmations applies here for the same reason: `huh` inherits deliberate title→field→help breathing room. **`huh`'s default theme carries its own palette and must not ship** |
| The text inputs | **`bubbles/textinput`** — a direct fit, and it is what `huh` uses underneath | §7.7's verdict, same primitive |
| Masking | `textinput`'s `EchoMode` — **but with `*`**, not the default `•` (§7.3) | the default echo character is Ambiguous-width |
| Field boundary | **Build fresh** — `lipgloss` border, `--z-border-strong` / `--z-focus-ring` | `huh`'s default field chrome is its own theme |
| Scanner | **Build fresh. Not `harmonica`** | §9.8 — harmonica models a damped spring converging on a target; the scanner is an undamped infinite alternating sinusoid. Evaluate the bezier directly |
| Message block | **Build fresh. Not `charmbracelet/log`** | §11.5 — `log`'s role is structured developer logging and its level colours are its own palette |
| Header band, footer, audio indicator | **Build fresh**, shared across all screens | §2.8 · §5.2 |
| Audio cue | `error` on states 7–10 ([`12-audio.md`](../../blueprint/12-audio.md) §4) — **fire-and-forget, never blocking, and never the only signal**: the message block renders whether or not anything is audible | |

---

## 16 · Screen-specific acceptance criteria

1. **The provider-descriptor test.** Register a throwaway provider with **three** credential
   fields — one secret, two not — and route to `Z-02`. It must render three field blocks, in slice
   order, with the right masking and the right help, **with zero changes to this screen's code.**
   This is the single most important acceptance line on the screen, because it is the whole reason
   the screen is called "connect a store".
2. **No `ProviderID` switch anywhere in `Z-02`.** Grep the screen package for `ProviderID` — the
   only legal appearance is as an opaque value passed through.
3. **A failed submit loses nothing.** Type into both fields, force each of the four classifier
   outcomes, and assert both fields still hold exactly what was typed.
4. **The private-profile copy never renders before there is evidence.** Grep every pre-submit
   render for `private`; there must be no match.
5. **`r` is not bound.** With a field focused, pressing `r` inserts `r`. Same for `q`, `?`, `,`,
   `s`, `a`, `/` and `m` (2.1.4).
6. **Paste works** on the masked field (3.3.8), and `Ctrl-R` reveals it.
7. **The mask is `*`.** Grep the render for `•`; there must be no match.
8. **Exactly one scanner**, only while a check is running, never beside a message block, never at
   Tiny (§9.3).
9. **The focus ring is present during the check** — assert the boundary is still `┃` mid-flight.
10. **Every field boundary is `--z-border-strong` or `--z-focus-ring`.** `--z-border` (1.53) marks
    no control anywhere on this screen.
11. **The vault line matches `Vault.Backing()`** — force both backings and diff the string.
12. **`esc` leaves from every state**, including mid-check (2.1.2).
13. **Chrome-cyan count is exactly 1** by [`../02-colour-budget.md`](../02-colour-budget.md) §3.1,
    at every colour depth, in every one of the four refusal states.
14. Eight artifacts per [`03-responsive.md`](../../blueprint/03-responsive.md) §7, **and one extra**:
    the four refusal states at `80 × 24`, because they are the screens this spec exists for.

---

## 17 · Open for the founder

1. **`Z-02` submits a real `Sync` to get its evidence (§8.1).** It is what makes the ratified
   private-profile copy land where the player hits it, and it costs one extra provider call.
   Confirm — and let `fft-tui-architect` record it against the `Syncer` seam.
2. **`Esc` does not blur on `Z-02` (§13.1).** A departure from the exhaustive `Esc` table, with its
   reason. Confirm, or accept a dead intermediate state on every form.
3. **`? help` and `q quit` stay in the footer while a field has focus (§9)** even though those keys
   type characters at that instant. The alternative is a footer that empties itself when anyone
   types. Confirm the screen-granularity reading.
4. **`02-colour-budget.md` §5 item 2 needs one word (§14.1)** — the error-state annunciator is red
   and is not on the closed list. Route to `fft-brand-architect`.
5. **`steamcommunity.com/dev/apikey` and the Steam key format are unverified here (§10.1).** They
   belong in the provider descriptor, read at source. This spec deliberately asserts no format
   (§8.2). Confirm someone owns that check before ship.
