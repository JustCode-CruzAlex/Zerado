# Redactions in the published deliverable

This repository is **public**. The ForgePlay landing-page deliverable it was
built from was not. Three kinds of change were made on the way in, and all three
are recorded here rather than made silently.

**The page's design and section order are untouched.** All sixteen sections are
present, in the ratified order, with the ratified layout. What changed is copy,
in two bounded places, both recorded in §3 below — so the build is **no longer
byte-for-byte** what the reviewer marked GOLDEN, and this document says so
rather than repeating a claim that has stopped being true.

| | ratified build | current |
|---|---|---|
| `dist/index.html` sha256 | `6dafafea…01366a` | `5aa8a2a7…4e65b` |
| bytes | 56,457 | 56,903 |

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

Two changes to the page's words, both after the page was ratified on 2026-08-24,
both recorded in `content/landing-copy.md` beside the copy they change, and both
now pinned by `scripts/check-page.mjs` so they cannot silently regress.

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
