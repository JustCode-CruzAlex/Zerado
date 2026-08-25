## What and why

<!-- What changed, and more importantly why. The diff already shows what. -->

Closes #

## Type

- [ ] Bug fix
- [ ] Feature
- [ ] Documentation
- [ ] Chore / dependency / CI
- [ ] Design or brand system

## Evidence

<!--
Anything visual needs proof, not assurance. For a change under site/**:
screenshots at 375, 768, 1280 and 1920 px. For a build or CI change: the output.
Verify with `npm run preview` — NEVER by opening dist/index.html over file://.
-->

| Width | Screenshot |
|---|---|
| 375 | |
| 768 | |
| 1280 | |
| 1920 | |

## The review bar

<!-- Tick what applies. Leave a line unticked and say why — that is more useful than a tick that is not true. -->

- [ ] CI is green (`npm run build` passes).
- [ ] The change does what the title says, and nothing unrelated.
- [ ] **Zero client-side JavaScript** and **zero external network requests** still hold on the page.
- [ ] **Zero raster images** in the render path (`og-card.png` is for social scrapers, not the page).
- [ ] Lighthouse ≥ 95 in every category; axe-core reports **0 violations** at all four breakpoints.
- [ ] **Colour is never the only signal** — every state co-renders colour + glyph + label.
- [ ] **Nothing dishonest ships**: no phase marked done before it is done, no mockup presented as a screenshot, no capability described in the present tense before it exists.
- [ ] No placeholder text (`lorem ipsum`, `TODO`, `[TBD]`) in shipped output.
- [ ] Ratified copy, design and section order are unchanged — or the change was agreed in an issue first (link it).
- [ ] Commit messages follow `CONTRIBUTING.md`, and **carry no AI attribution**.
- [ ] Tests accompany the change (once there is Go code to test).

## Notes for the reviewer

<!-- Anything you are unsure about, deliberately left out, or want a second opinion on. -->
