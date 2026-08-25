---
title: Zerado — the cross-check register
discipline: ARCHITECTURE
doc-no: ZRD-SPINE-14
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: status-report
ticket: "#2"
---

# The cross-check register

Every contradiction the three deliverables found in each other, enumerated so the count is
**auditable rather than asserted**.

> **This document exists because a reviewer was right about it.** `00-index.md` claimed *"25
> contradictions were found and closed"* and the number appeared nowhere else, which made it a
> figure the founder had to take on trust — in a bundle whose stated standard is that *"a wrong
> premise that reaches the right answer is still a defect."* Enumerating them corrected the count
> too: it is **29**, not 25. The original figure double-counted two items and omitted the six the
> design architect found.

**35 findings: 29 from the cross-check (all closed), 6 more from building D7's theme gate (4 closed, 2 open upstream).**

Two of them — #15 and #17 — were briefly marked *reopened* and *partial* after `fft-tui-designer`
read this register against head and found it over-claiming. Both were then closed by commit
`2e3a787`, **and this register did not say so for three commits.** A GOLDEN review caught that: the
register had gone from over-claiming to *under*-claiming, which is the same defect wearing the other
face.

The pattern behind it is worth more than the two rows: **the last three commits fixed source
documents without re-running the checkers against head.** The verification layer lagged the content
layer. That is the failure this bundle is least entitled to, since its whole argument is that the
verification layer is trustworthy — so §"How to re-run these checks" exists to make the re-run
cheap, and it is now part of finishing a change rather than a thing done once. Each names who found it, where it was, and where the fix
landed.

---

## A · Found by `fft-design-architect` while writing deliverable B — 6

| # | Finding | Closed |
|---|---|---|
| **1** | **`♪` U+266A does not carry an emoji presentation.** I had asserted it did. Verified against `emoji-data.txt`: no entry — its neighbours `♠ ♣ ♥ ♦ ♨` are listed, the music note is not. The width argument alone carries the decision | `01-design-system.md` §1.2, recorded as a checked-and-refuted claim |
| **2** | **`▪` U+25AA *is* emoji-listed**, and was in use as the Phase 3 price marker — the exact failure the emoji verdict exists to prevent, sitting inside the document that defines the rule | Moved to `▬` U+25AC (Neutral, non-emoji) |
| **3** | **WCAG 1.4.7 was cited for the music bed.** Read at source it is scoped to prerecorded audio-only content with speech in the foreground — inapplicable | Corrected to a checked-and-rejected note; **1.4.2** is the criterion that actually bites |
| **4** | **`03-designer-manual.md` §5.11 still recorded the audio verdict as permanently rejected** after the reversal | Struck through, marked SUPERSEDED, provenance note added |
| **5** | **`00-design-brief.md` §3.3 listed 1.4.2 Audio Control as inapplicable** (*"Zerado plays no audio"*) — the opposite of true after the reversal | Moved to §3.1, the bites-directly table |
| **6** | **Brand manual §9's invariant list is incomplete.** It names five invariants carrying identity across the three render targets; audio is a sixth, and the manual predates it by one day | Recorded as a founder-gate item with `fft-brand-architect` named as owner — **not** edited downstream, because a downstream document extending the brand's own invariant list is the drift §9 exists to prevent |

## B · Found by `fft-tui-designer` (deck cluster) — 12

| # | Finding | Closed |
|---|---|---|
| **7** | **The 74-column game row disagreed with itself** — the spine gave title 42 · playtime 6 · source 4; the design system gave 43 · 5 · 3. Both summed to 74, which is worse than one being wrong: a builder following one and a reviewer following the other would disagree about a correct screen | Design system aligned to the spine, **with the cause recorded** — the title went 43 → 42 when the glyph field became two columns |
| **8** | **"ExtraWide → title 81"** assumed a full-width 112-column ledger the spine does not compose | Corrected to **40**, with the trap named: ExtraWide *narrows* the title |
| **9** | **The status bar was placed in the reserved footer row**, which the navigation model explicitly forbids from carrying status | Moved to body row 1 — which is what makes R-10(c) true |
| **10** | **A "three facts maximum" voice guideline** contradicted the state machine's five-fact summary requirement | Reconciled: the guideline governs sentences, and where a falsifiable rule conflicts, the rule wins |
| **11** | **The 90-day staleness banner named prices** — a Phase 3 capability — in a Phase 1 document | Phase 1 copy names the library; the price wording is deferred with its reason |
| **12** | **`Sort` was listed as a Phase 1 feature** with no key bound and no screen owning it | Replaced with `Ordering`, described honestly: fixed title A→Z, no control |
| **13** | **The `Z-06` consequence line was 36 cells** against a binding 28-cell overlay content width | Shortened; the requirement is the naming, not the wording |
| **14** | **`Z-06`'s `34 × 11` box has no room for key hints** — all nine content rows are spent | Hints live in the frame's reserved footer row; an overlay borrows the one on screen |
| **15** | **The detail pane was drawn at 28 columns**, a composition the spine rejected | **CLOSED** (`2e3a787`). `01-design-system.md` states the pane three ways: §5 has the spine's `66 ∥ 2 ∥ 44`, §6.1 has `ledger 64 · gutter 2 · pane 46`, and §6.2 is still headed *pane 28 cols*. §6.1 and the spine both sum to 112, which is the same worse-than-wrong failure as finding #7 — a builder and a reviewer would each think the screen correct. All three now read the spine's `66 ∥ 2 ∥ 44` |
| **16** | **The bordered-surface inset was drawn but never named** — two screens depended on a number read off a mockup | Named: `BorderInsetX` = 2, `BorderInsetY` = 0 |
| **17** | **`ZERADO_ASCII` covered only the state column**, but box drawing, the focus marker and the scanner are equally Ambiguous and carry the frame itself | **CLOSED** (`2e3a787`). Extended in the spine ([`03-responsive.md`](./03-responsive.md) §5b); `01-design-system.md` §1.2 rule 4 still names only the state column. Its box-drawing argument is about *internal alignment*, not absolute width — and a `34 × 11` overlay whose border rows are **entirely** box-drawing is **68 cells wide** on an `ambiguous-width=double` terminal while its content rows are ~36. Rule 4 now switches the whole vocabulary and carries that reasoning |
| **18** | **Nothing specified where a key's description lives** — dispatch, footer and help could be three strings that drift | One key registry is now a spine requirement; `Z-10` is generated from it |

## C · Found by `fft-tui-designer` (entry and failure cluster) — 11

| # | Finding | Closed |
|---|---|---|
| **19** | **The offline contract contradicted itself** — §5 says Zerado never probes and learns it is offline only when something fails; §6 drew first run with a door already disabled at first paint | Resolved by naming the real distinction: **packet versus no packet**. Reading the kernel's own routing table emits nothing |
| **20** | **The closed red list omitted the error-state annunciator**, so read literally every error state in the product failed a checklist line | List now names both; both are alarms |
| **21** | **My `⚠` rejection cited a wrong premise** — U+26A0 is **Neutral**, not Ambiguous. It *is* emoji-listed, so the conclusion held | Recorded as a correction rather than quietly repaired |
| **22** | **`r to retry` is unusable on a form** — with a text input focused, `r` types `r` (WCAG 2.1.4), so the copy named a key that cannot fire | The message carries the fact; the **screen** supplies the affordance |
| **23** | **`Settings → Steam` pointed at the screen the player is standing on** when it fires on `Z-02` | Destination is the screen's to fill in |
| **24** | **`everything below still works` names a composition the message cannot see** — true on `Z-04`, false on `Z-02`, shapeless on `Z-03` | The clause is the screen's to add where it is true |
| **25** | **The ellipsis is Ambiguous width** and is the character everyone forgets to count *inside* a truncation — so a title truncated to fit is one cell too wide | Its own rule in §1.2: budget = field width − **measured** marker width, at render time |
| **26** | **`BodyRect.h` charges `InnerPaddingY` once but `OuterMarginY` twice** | The asymmetry is the canon's and is now stated, because every mockup adds up against it |
| **27** | **`Z-09` read as master/detail at ExtraWide** would break its own must-not rule by composition rather than by a click, and interleave reading order under 1.3.2 | Never master/detail, at any tier; the width buys a wider gutter |
| **28** | **`schema_migration` could not supply a fact `Z-11`'s copy promised** — a schema number cannot name a binary | Added `written_by`, which cannot be backfilled later |
| **29** | **`state` versus `status` were both in use with no rule** | Named convention: *state* is what the player sees, *status* is what the machine stores and commands |

---

## D · Found by the theme validation, after the founder amendment — 6

Building D7's gate immediately found things a gate is for. Numbered separately because they arrived
after the 29 and are **not** all closed.

| # | Finding | State |
|---|---|---|
| **30** | **The light state set fails its own gate** — `not started × zerado` **5.41** (protan), `zerado × abandoned` **8.91** (deutan), floor 10.0. Brand §4.5's light `not started` is rotated 173° and lands 1.1° from `#9FB0C6`, *the blue-cast steel §4.4 rejected on the dark side.* The dark correction was never carried to paper | **OPEN** — `fft-brand-architect`, and no light theme ships until it is repaired |
| **31** | **`tokens.css` §10 records `7.30:1` for `--z-state-abandoned`; it measures `7.67:1`** — 1 of 35 brand contrast figures that did not reproduce | **OPEN** — upstream |
| **32** | **My light-theme list was wrong in four places.** I relayed `solitude` and `ethereal` as light (both are `mode = "dark"`) and omitted `lupine` and `rose-pine` (both light). **Mine — I passed a list on without checking it**, in a bundle whose standard is that claims are verified at source | **CLOSED** — corrected from the files |
| **33** | **The theme corpus is two dialects, not one.** 32 files are flat-ANSI; 3 use omarchy's semantic schema and have no `colorN` at all. Both define `background`/`foreground`/`accent`, so a semantic file **parses clean and yields sixteen empty slots** — the exact defect FlowForge shipped and documented as `#5368` | **CLOSED** — the loader must reject, not silently accept |
| **34** | **No harvested light theme reaches a usable tier.** `catppuccin-latte`, `flexoki-light` and `rose-pine` each clear 4.5:1 on only **2 of 18** entries; `catppuccin-latte`'s yellow is **2.31:1**. Systematic, not coincidence — omarchy light palettes are authored for syntax highlighting, not body text | **CLOSED** — Zerado's light offering must be its own authored expression |
| **35** | **The ΔE floor was a tautology.** 11.9 is the default palette's own measurement under an unpinned method, so the floor moved whenever the palette did. Replaced with the manual's own *"≥ 10 is distinct"* | **CLOSED** — floor 10.0, method pinned |

**Finding 32 is mine and is the one worth keeping.** A list arrived in a relay, I passed it to a
specialist as fact, and four of its entries were wrong. Nothing downstream broke, because the
specialist checked the files rather than trusting me — which is the same reason two of my width
assertions were caught earlier. The lesson is not *be more careful*; it is that **the habit of
verifying at source is what makes a chain of agents safe to run at all.**

---

## How to re-run these checks at the next gate

The reviewer's fair objection to the first version of this bundle was that its verification numbers
were **load-bearing but unreproducible**. It suggested committing the checkers under `tools/`. That
would add executable files to a ticket whose out-of-scope line is *"All Go code… not one line of
implementation"*, and whose no-code property has now been verified twice — so the method is written
down instead. Each of these is a one-liner against a clean checkout.

| Claim | How to reproduce it |
|---|---|
| **No code in the diff** | `git diff --name-only main..HEAD \| grep -vE '\.(md\|toml\|svg\|excalidraw)$'` → empty |
| **Broken relative links** | Walk every `docs/**/*.md`, resolve each inline-link target that is not `http`, against the file's own directory |
| **Every doc carries an archetype** | `grep -L '^archetype:' docs/**/*.md` → empty |
| **Offline invariant** | In `07-offline-contract.md` §2's Phase 1 table, extract `` `Z-NN` `` from column 1 and the class word from column 2 → **12 rows, 12 distinct, 9 `WORKS` · 2 `NEEDS THE NETWORK` · 1 `DEGRADES`** |
| **Mockup widths** | For every ```` ```text ```` block bound to a `RENDER W×H` heading, measure each non-blank line with `east_asian_width`, counting `W`/`F` as 2 and **Ambiguous as 1** (the design system §1.2 rule) → no line exceeds its declared `W` |
| **Title-block overrun** | For each rendered SVG, find the `[titleblock].description` text run and the **cell** rules (short vertical lines inside the title band — *not* the cyanotype theme's full-height graticule) → the run ends before the next cell rule |
| **Charts are live, not stale exports** | Re-run `flowforge chart render` on all ten specs and diff the SVGs; the only permitted difference is the `RENDERED FROM` path prefix, which reflects the invocation directory |

**The one number that was not reproducible has been retired.** An earlier PR body claimed *"82
mockups machine-verified"*; that was two agents' self-reports summed, and they counted different
things. The reproducible figures are: **72** fenced `text` blocks across the eleven specs, of which
**16** are bound to an explicit `RENDER W×H` heading — **158 lines measured, 0 over.** An
independent re-measurement of all render-shaped blocks against each block's own ruler also found
**0 overruns**.

---

## What the register is actually evidence of

Not that the bundle was written carelessly — every one of these was found **before** the bundle was
surfaced, by a specialist reading another specialist's output rather than their own.

Three are worth reading twice:

- **#6, #21 and #1** are corrections to *upstream* work, including two of mine. Both of my width
  assertions reached the right conclusion on a wrong premise, which is the failure mode hardest to
  catch, because nothing downstream misbehaves.
- **#18** was found by a designer noticing that its own screen's promise — *"every key that does
  anything"* — had no mechanism behind it, and escalating rather than writing a spec that could not
  be true.
- The **absent-row propagation** (`06-data-seams.md` §2.4 reaching `Z-04`, `Z-05`, `Z-07` and sheet
  01) was **not** caught by this cross-check and was found by the independent reviewer instead. It
  is the register's own honest limit: the cross-check ran when the specs were written, and §2.4 was
  added afterwards. A cross-check is a snapshot, not a standing guarantee.
