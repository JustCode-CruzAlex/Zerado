# Redactions in the published deliverable

This repository is **public**. The ForgePlay landing-page deliverable it was
built from was not. Three kinds of change were made on the way in, and all three
are recorded here rather than made silently.

**The page's design and section order are untouched.** All sixteen sections are
present, in the ratified order, with the ratified layout. What changed is copy —
each change recorded in §3 below, each one either a founder instruction or a
correction of something that was not true — so the build is **no longer
byte-for-byte** what the reviewer marked GOLDEN, and this document says so
rather than repeating a claim that has stopped being true.

| | ratified build | current |
|---|---|---|
| `dist/index.html` sha256 | `6dafafea…01366a` | `47fdfe02…273cd` |
| bytes | 56,457 | 57,387 |

---

## 1. Documents that were not published at all

The `ratification/` folder, the delivery summary and the play notes stayed out
of this repository. They are the founder-facing record of the play — the six
ratification questions and answers, commercial-licensing exposure on a
metadata provider, trademark-clearance status, and the negotiating position on
an unnamed community source.

None of it is documentation an open-source contributor needs, and some of it is
material a company does not publish about itself. It remains in the ForgePlay
output alongside the run.

The five PDFs rendered from those documents (`brief`, `decisions`, `qna`,
`mock.outline`, `summary`) were removed with them. The nine that remain
correspond one-to-one to documents that **are** published here.

## 2. One word, removed everywhere

**Ratified decision Q4** bans a specific community-source name from every Zerado
surface: it is not named on the page, in an icon, in a screenshot or in the FAQ,
because naming it would not currently be true. The ticket extends that ban to
this repository's public documentation.

That creates a small paradox: the documents that *enforce* the ban are the ones
that quoted the word, in order to say it must not appear. A contributor running

```bash
grep -ri <the-name> .
```

would get hits — and every hit would be a rule saying the word must not be
there. So the word was removed from the enforcing documents too, and each check
was kept working:

| File | What changed |
|---|---|
| `design/blueprint.md` | The prose constraint now names the rule, not the word. |
| `design/blueprint.tokens.json` | `forbiddenStrings` holds the sentinel `$Q4_BANNED_SOURCE_NAME` instead of the literal, plus a `forbiddenStringsNote` explaining why. |
| `qa/harness/run-content.mjs` | The forbidden-token regex is assembled at runtime from fragments. **The check is unchanged and still fails on a real occurrence** — only the literal is gone from the source file. |
| `qa/qa-report.md` | Three references reworded; the recorded results (`0` hits) are untouched. |
| `review/review.md` | One row in the ratification-compliance table reworded; the verdict is untouched. |
| `.github/workflows/guardrails.yml` | The CI guardrail builds the same pattern with `printf`, for the same reason. |
| `scripts/check-page.mjs` | The page invariant assembles it the same way. |

`docs/pdf/blueprint.pdf`, `docs/pdf/qa-report.pdf` and `docs/pdf/review.pdf`
were **re-rendered from the redacted sources** so the PDFs and the Markdown do
not disagree. The other six PDFs were not touched.

**No measurement, result, score or verdict was altered.** The QA report still
records zero occurrences, and the review verdict is still GOLDEN with 0 blocking
and 0 major.

The rule is now enforced in CI by [`.github/workflows/guardrails.yml`](../.github/workflows/guardrails.yml),
which runs on **every** pull request and every push to `main` with no `paths:`
filter — scanning every tracked text file *and* extracting text from all nine
PDFs. It deliberately carries no path filter: when these checks lived inside the
path-filtered site workflow, a pull request touching only `README.md` or
`docs/**` never triggered them at all.

The day that source is named publicly, this is a one-line change in each file.

## 3. Copy amended after ratification

Five changes, all after the page was ratified on 2026-08-24, each recorded in
`content/landing-copy.md` beside the copy it changes, and each pinned by
`scripts/check-page.mjs` so it cannot silently regress. Three carry a founder
instruction, quoted; one corrects a false statement; one is an accessibility fix
with no copy change at all.

### 3.1 — Phase 1 reads *In progress* (§12)

Founder instruction, 2026-08-25: *"we need to show Phase 1 'In progress'"*.
Phase 1 work started that day. Phases 2–4 are unchanged, no phase is marked
done, and no done-equivalent status exists to mark one with. The blueprint's
change-control clause (§9) names roadmap status as a thing the build may not
change without coming back to it, so it was amended properly —
`design/blueprint.md` §1.0, **Amendment 1**, rev A → rev B.

### 3.2 — the Steam capability claims moved to the future tense (§06, §14)

At rev A the page said **"Steam is built and works"** (§14), **"Steam syncs
today"** (§06), and rendered the Steam store row as `status="live"`.

**There is no Go code in this repository**, so all three were false in the
present tense — which is the exact thing this project's own review bar forbids
(`.github/PULL_REQUEST_TEMPLATE.md`: *"no capability described in the present
tense before it exists"*), and which the founder's own statement of 2026-08-25
confirms: nothing is built besides the website.

The change is **tense only**. No section moved, no design changed, no claim was
added. `StoreRow`'s `live` variant is deliberately kept, unused, because it is
what Phase 1 flips on the day Steam sync ships — at which point this section
records the reversal.

This one was **not** covered by a founder instruction. It was raised by the
independent review across two rounds, escalated to the founder twice without an
answer, and then applied on the reasoning that the ticket's out-of-scope clause
protects the page's ratified *creative* decisions — its design, sections, voice
and section order — and not a statement of fact the founder has since
contradicted. Reverting it is a four-string change.

### 3.3 — no premium account (§11, §14)

Founder instruction, 2026-08-25: *"let's not have premium account for now and
only donation, so we're not commercial app."* **Amends ratified decision Q3.**
The Phase 4 layer is now stated as donation-supported; every mention of a
premium account or paid tier is gone.

Ratified decision Q6 still governs and is unchanged: **disclosure is not an
ask.** There is still no donate button, no sponsor button and no funding CTA —
CI asserts zero funding controls on every build.

### 3.4 — the affiliate model is dropped (§14, §16)

Founder's condition, 2026-08-25: *"affiliate links I don't think is commercial
since I'm not charging customers and just using the money to maintain servers.
But if is a must to fit the non-commercial, drop it, and keep only the
donation."* **Amends ratified decision Q6.**

IGDB's published test for its free tier is whether the **project generates
revenue** — not whether users are charged, and not what the money is spent on.
An affiliate commission is revenue generated by the project, so the condition is
met and the model is dropped. The footer's affiliate disclosure is replaced by a
truthful money statement; both §14 money answers are rewritten.

**The price-intelligence feature (§07) is untouched** — all-time low, current
price and the wait-or-buy verdict all stay exactly as ratified. Only the
affiliate tag on an outbound link goes, and the page never described that tag,
so §07 needed no copy change at all.

`content/seo.md` and the page's JSON-LD were checked and carry no monetization
language, no `offers` and no `price`. Nothing to amend there.

### 3.5 — an accessibility fix, not a copy change (all six mockup figures)

`role="img"` sat on six `<figure>` elements, which ARIA-in-HTML disallows —
Lighthouse's `aria-allowed-role` failed on all six. The role moved to a wrapper
`<div>` where it is valid (and was simply removed from §06's figure, since
`CoverGrid` already carries it correctly on its own `<div>`).

Verified layout-neutral: page height identical at 375 / 768 / 1280 / 1920, and
0.002–0.004% of pixels differ — a 2px band, the animated scanner caught at a
different phase. `aria-allowed-role` violations went 6 → 0 and live Lighthouse
accessibility went 99 → 100.

**Honest limit:** what this fixes for certain is ARIA validity. Whether a real
screen reader now prunes the mockup's fabricated library from its subtree is
**not verified** — Playwright's aria snapshot is not the platform accessibility
tree, and it still showed the contents. The pre-existing disclosure stands: every
accessibility check here is programmatic, and no testing with real assistive
technology has been done.

### The delta, in full

The complete visible-text difference between the ratified build and the current
one, word-diffed — nothing else on the page moved:

```
§06   - syncs today                    §12   - Every one of them is planned — none
      + is what Phase 1 builds               + Phase 1 is under way — nothing here
      - are planned next:                    - ○ Planned
      + come after:                          + ◐ In progress
      - ◉  →  + ○
                                       §14   + None of them yet — nothing is runnable.
      + "— planned" ×4                       - built and works.  → + what Phase 1 builds;
        (visually hidden; the store           - are planned, not yet built.  → + follow.
         status reaching a screen              - are  → + will be
         reader, which it previously
         did not)
```
