---
name: feedback-verify-at-source-not-the-brief
description: On Zerado, re-verify every figure at source even when the dispatching brief summarises it — the summary has been wrong — and attribute each number as the manual's or as computed.
metadata:
  type: feedback
---

Re-read the primary source before asserting anything, **including facts the dispatching brief
already summarised for you**, and label every printed figure either *the manual's* or *computed by
me, with the method named*. Never estimate.

**Why:** the controller's own instruction on 2026-08-25 was *"Read them at source before writing.
Do not trust this summary for anything you assert"* — and the summary did in fact contain two
errors (its light-theme list named two dark themes and missed two light ones). Separately, one
figure in `tokens.css` was wrong at source. Both were only findable by recomputation. The
attribution discipline exists because the brand manual's **contrast** digits reproduce exactly
while its **ΔE** digits do not, so treating them the same way silently launders an unreproducible
number into a spec.

**How to apply:** on any Zerado design document — recompute the corpus's own published figures as a
self-check before using them (34 of 35 contrast ratios reproduced; the mismatch was the finding);
state the method when you compute; and where canon does not settle something, mark it explicitly as
a **design decision** with the reason. Voice on this project is dry and concrete: no exclamation
marks, no emoji, no hedging. Where a fork is genuinely consequential, present the rejected option
in a table with its cost rather than burying it (`05-theme-system.md` §6 is the shape that landed).

Related: [[project-zerado-theme-system]], [[reference-zerado-source-corpora]]
