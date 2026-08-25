---
name: project-zerado-theme-system
description: Zerado's light state palette FAILS its own CVD gate and blocks every light theme until brand governance repairs it — plus the two source errors found on 2026-08-25.
metadata:
  type: project
---

**Zerado's light/paper state set does not pass Zerado's own four-state gate**, measured
2026-08-25 while authoring `docs/design/05-theme-system.md` (ticket #2, founder amendment).

**Why:** brand §4.5 claims the paper colours are *"the same hues carried to ink weight"*. Four of
five are (within 9.5°). **`not started` is rotated 173°** — from the dark set's warm `h 92.97°` to
`h 265.99°`, landing 1.1° from `#9FB0C6`, the blue-cast steel brand §4.4 **explicitly rejected**
on the dark side for collapsing against the cyan. Result: `not-started × zerado` separates by only
**ΔE 5.41** under protanopia and `zerado × abandoned` by **8.91** under deuteranopia, against a
required 10.0.

**How to apply:** **no Zerado light theme ships — terminal or phone — until this is repaired**, and
the repair is `fft-brand-architect`'s through brand governance §10 (tokens.css + tokens.json +
manual, one commit, ANSI-256 re-derived). A minimum-motion feasibility proof exists in
`05-theme-system.md` §3.4 (carry the hue; darken `abandoned` by ΔL\* −8 → floor 10.83) and is
explicitly **non-binding** — never quote its hexes as canon.

**Two errors found at source, both still open:**
1. `tokens.css` §10 records `7.30:1` for `--z-state-abandoned` `#6D3D93`; it measures **7.67:1**.
   Nothing shipped broken (the true value is higher). Correction routed to brand governance.
2. The corpus summary handed to me listed `solitude` and `ethereal` as light themes — **both are
   dark** (`mode = "dark"`; bg luminance 0.006 / 0.004) — and omitted `lupine` and `rose-pine`
   (the bundled `rose-pine` is the **Dawn** light variant), which are.

**Also settled in that document:** the CVD floor is **ΔE 10.0**, not the default's own 11.9/11.81 —
a floor set at the thing being measured is a tautology. The method is now **pinned** in
`05-theme-system.md` §2.1 (matrices, D65, linear-RGB clamp, CIEDE2000 kL=kC=kH=1) because two
implementations of the manual's unpinned method disagree in the third significant figure.

Related: [[reference-zerado-source-corpora]], [[feedback-verify-at-source-not-the-brief]]
