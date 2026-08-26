---
title: Zerado — the CLI as a versioned API
discipline: API
doc-no: ZRD-API-03
rev: A
date: 2026-08-25
status: draft — for review
archetype: implementation-plan
ticket: "#6"
---

# The CLI verb surface

A CLI is an API the moment anyone types it. A verb in somebody's shell history, a script parsing an
exit code, a `cron` line — each is a caller with no upgrade path.

So the surface is declared as **data** in `internal/cli`, with a golden test over it, and the
stability policy is written down before there is anything to break.

---

## 1 · Phase 1

| Verb | Arguments | Flags | Class |
|---|---|---|---|
| *(none)* | | | the TUI. Interactive |
| `sync` | `[provider]` | `--all` | **NEEDS THE NETWORK** |
| `list` | | `--state --source --search --absent --limit` | WORKS |
| `mark` | `<game> <state>` | `--clear` | WORKS |
| `add` | | `--title --platform --state --owned-since` | WORKS |
| `game` | `<game>` | | WORKS. Interactive — the deep link |
| `doctor` | | | WORKS |
| `version` | | | WORKS |
| `help` | `[verb]` | | WORKS |

Global: `--json --quiet --db=path --no-color --version --help`.

`TestOfflineClassIsDeclared` asserts that **only `sync`** claims to need the network. Everything else
is a read of a local file, which is
[`../blueprint/07-offline-contract.md`](../blueprint/07-offline-contract.md)'s classification
expressed where a scripted caller can see it *before* running anything.

### Two verbs that exist because a screen needed them

**`mark --clear`.** Clearing an override is a different action from setting `NOT STARTED`
([`../blueprint/05-state-machine.md`](../blueprint/05-state-machine.md) §5) — clearing on a game with
playtime makes it `IN PROGRESS` immediately, while choosing `NOT STARTED` stores a value that
sticks. A CLI without `--clear` could not do what `Z-06` does.

Its arity needs `Arg.RequiredUnless`: `state` is required **unless `--clear` is present**. An earlier
revision declared `state` flatly required alongside the flag, which meant the surface — read as data,
which is the whole point of declaring it as data — rejected the invocation this paragraph says must
work. A surface that contradicts its own stated requirement is worse than an incomplete one, because
only whoever writes the parser discovers it. `TestTheSurfaceCanExpressEveryDocumentedInvocation` now
runs the documented invocations against the declared arities, so the two cannot drift apart without
one of them failing. *(Found by the review at `4484d9a`.)*

**`game <id>`.** [`../blueprint/04-navigation-and-focus.md`](../blueprint/04-navigation-and-focus.md)
notes that deep-linking requires the route stack to be constructible from a descriptor rather than
only by pushing. That is a **Phase 1 constraint caused by this verb**, and it is recorded next to
the verb that causes it.

---

## 2 · Reserved, and not half-working

`tonight` · `price` · `watch` · `tag` · `enrich` · `devices` · `export` · `import`.

The names are claimed so a Phase 1 user's script cannot come to depend on `zerado tonight` meaning
something else, and so the eventual feature does not have to settle for a worse word.

They declare **no arguments and no flags** — `TestReservedVerbsPromiseNothing` asserts it — because a
reserved verb has no shape yet, and because copy that mentions a capability the build does not have
presents something unbuilt as working.

---

## 3 · Exit codes

`0` ok · `1` internal · `2` usage · `3` offline · `4` unreachable · `5` unauthorized · `6` empty ·
`7` rate-limited · `8` not found · `9` malformed · `10` state · **`130` cancelled**.

`130` borrows the shell's own convention for a process killed by `SIGINT`, so a player who pressed
Ctrl-C sees the number they would have seen anyway.

The mapping is total over the taxonomy and asserted by `TestExitCodesAreTotal`.

---

## 4 · `--json`

One shape for success and failure, because a caller that has to parse two different top-level
structures depending on an exit code will get it wrong on the day it matters.

```json
{"api":1,"ok":true,"data":{"games":247}}
{"api":1,"ok":false,"error":{"kind":"offline","op":"steam.Sync","subject":"Steam",
                             "message_key":"fault.offline","retry_after_seconds":0}}
```

**The envelope carries no message.** That is not an omission: a message is user-facing text rendered
from the catalogue in the player's language, and a pre-rendered English sentence in a JSON field is
the hardest D9 violation to spot. A consumer that wants a sentence renders one from the key; a
consumer that wants to branch uses the kind, which is what it should have been using anyway.

`retry_after_seconds` is a number rather than `"1m30s"`, because every JSON consumer can compare a
number and not all of them can parse a duration string.

---

## 5 · The stability policy

- `api` is the contract's **major** version and appears in every envelope.
- **Fields are added, never removed and never repurposed.** A consumer that ignores unknown fields
  keeps working across every minor change.
- **A field's meaning is fixed.** Changing what `unchanged` counts is a breaking change even though
  the JSON looks identical, and it bumps the version.
- **Error kinds are the taxonomy's stable names.** Adding one is additive; a consumer that does not
  recognise a kind must treat it as a failure rather than as success — which is why `ok` exists
  separately from the kind.
- **Exit codes and kinds move together** and are documented together.
- **A verb is never renamed** within a major version. `TestVerbSurfaceIsStable` holds a golden list,
  so a rename has to be a deliberate act that updates the test and says why in the diff.
- Status values (`not_started` · `in_progress` · `zerado` · `abandoned`), provider ids, run statuses
  and backing kinds are all **stored values as well as printed ones**. They are API twice over.
