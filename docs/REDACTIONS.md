# Redactions in the published deliverable

This repository is **public**. The ForgePlay landing-page deliverable it was
built from was not. Two kinds of change were made on the way in, and both are
recorded here rather than made silently.

Nothing about the **page itself** was changed. Its copy, design, section order
and build output are byte-for-byte what the reviewer marked GOLDEN.

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
| `.github/workflows/site.yml` | The CI guardrail builds the same pattern with `printf`, for the same reason. |

`docs/pdf/blueprint.pdf`, `docs/pdf/qa-report.pdf` and `docs/pdf/review.pdf`
were **re-rendered from the redacted sources** so the PDFs and the Markdown do
not disagree. The other six PDFs were not touched.

**No measurement, result, score or verdict was altered.** The QA report still
records zero occurrences, the review verdict is still GOLDEN with 0 blocking and
0 major, and the CI guardrail in `.github/workflows/site.yml` now enforces the
same rule on every build.

The day that source is named publicly, this is a one-line change in each file.
