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

**29 distinct findings. 27 closed, 1 reopened, 1 partial** — see #15 and #17, both caught by
`fft-tui-designer` reading *this register* against head and finding it over-claimed. A register that
asserts closure it has not verified is the same defect it exists to record. Each names who found it, where it was, and where the fix
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
| **15** | **The detail pane was drawn at 28 columns**, a composition the spine rejected | **REOPENED.** `01-design-system.md` states the pane three ways: §5 has the spine's `66 ∥ 2 ∥ 44`, §6.1 has `ledger 64 · gutter 2 · pane 46`, and §6.2 is still headed *pane 28 cols*. §6.1 and the spine both sum to 112, which is the same worse-than-wrong failure as finding #7 — a builder and a reviewer would each think the screen correct. Being fixed to the spine's split |
| **16** | **The bordered-surface inset was drawn but never named** — two screens depended on a number read off a mockup | Named: `BorderInsetX` = 2, `BorderInsetY` = 0 |
| **17** | **`ZERADO_ASCII` covered only the state column**, but box drawing, the focus marker and the scanner are equally Ambiguous and carry the frame itself | **PARTIAL.** Extended in the spine ([`03-responsive.md`](./03-responsive.md) §5b); `01-design-system.md` §1.2 rule 4 still names only the state column. Its box-drawing argument is about *internal alignment*, not absolute width — and a `34 × 11` overlay whose border rows are **entirely** box-drawing is **68 cells wide** on an `ambiguous-width=double` terminal while its content rows are ~36. Being fixed |
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
