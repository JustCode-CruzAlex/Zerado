---
title: Zerado — Z-11 Fatal error
discipline: SCREEN SPEC
doc-no: ZRD-SCREEN-11
rev: A
date: 2026-08-25
status: draft — for founder ratification
archetype: component-variants
ticket: "#2"
---

# Z-11 · Fatal error

> Fills [`../03-designer-manual.md`](../03-designer-manual.md) §3's 16-section contract — with
> **§6 and §11 answered by exemption**, because this screen has no frame and therefore no spacing
> tokens. §6 says so explicitly rather than leaving the section blank.
>
> Composition binding from [`../../blueprint/02-composition.md`](../../blueprint/02-composition.md) §2:
> **FRAMELESS. Plain text, left-aligned, no chrome. `R = 0`.**
>
> The one rule this screen exists to keep:
> *"`Z-11` must not depend on anything that could be what broke."*
> ([`01-screen-inventory.md`](../../blueprint/01-screen-inventory.md) §5.)

---

## 1 · Identity

| | |
|---|---|
| **Screen** | `Z-11` · Fatal error |
| **Phase** | 1 |
| **Kind** | Route — **and deliberately off the route graph.** *"It replaces everything, from anywhere, and there is no way back into the program from it"* (nav §2) |
| **Routes in** | any point at which the program cannot continue, including before the TUI starts |
| **Routes out** | **none.** The program ends, or — in `HOLD` mode — the terminal is resized |
| **Offline class** | **WORKS** ([`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §2) — *"It reports a local failure and depends on nothing. It is the one screen that is more reliable offline, because it reaches for nothing at all"* |
| **Stream** | **stderr** in `EXIT` mode; the program's own stdout view in `HOLD` mode — §3.2 |

---

## 2 · Purpose

**When the program cannot continue, say what broke, where the file is, and what to try — in the
plainest renderer that can exist.**

A fatal error may be a failure of the very subsystem the frame is built on. **A crash screen that
itself crashes tells the player nothing.**

---

## 3 · What this screen has — and the exact list of what it does not touch

### 3.1 · The dependency list, stated as a prohibition so it can be grepped

`Z-11` renders with `fmt.Fprintln` and a word-wrapper. That is the whole toolchain.

| Not used | Because |
|---|---|
| `lipgloss` / any layout engine | it is a dependency, and layout is one of the things that can be what broke |
| the theme, the token table, `space`, `frame`, `chrome` | the frame and its sizer are the subsystem most likely to be implicated |
| **any SGR sequence, ever** | zero colour, zero bold — with or without `NO_COLOR`. §12 |
| **any non-ASCII character** | no `▌`, no `─`, no `✦`, no `○`, no `…`, no `→`, no box drawing. §3.3 |
| the width-aware measurement table | unnecessary **by construction** once the output is ASCII, which is `Na — Narrow` in every terminal |
| `Store` / the database | the library may be exactly what failed. `Z-11` reads **nothing** |
| `Vault` / the keychain | same |
| the network | there is nothing to ask anybody |
| **the audio subsystem** | **`Z-11` is silent, always** — no `Cue()`, in any mode, for any case. A crash screen must not try to play a sound, and `error` is on the cue list precisely so that *ordinary* errors can carry one. This is not an ordinary error |
| a spinner, a scanner, a ticker, any motion | there is nothing indeterminate about a program that has stopped |
| an input loop | in `EXIT` mode. See §3.2 |

### 3.2 · Two modes, and the difference is whether the program is still alive

The nav model, the responsive floor and the ticket's four cases only reconcile if `Z-11` is
**one screen with two modes**. Naming them makes every other section unambiguous.

| | **`Z-11`-EXIT** | **`Z-11`-HOLD** |
|---|---|---|
| When | cases 1, 2, 3, 5 — and case 4 **at start-up** | case 4 **mid-session**: a running terminal dragged below the floor |
| The program is | **ending** | **healthy, and still running** |
| Renders to | **stderr** | the program's own view (stdout) |
| Input loop | **none.** There is nothing to wait for, and an input loop is a dependency on the input subsystem | the existing loop; **only `WindowSizeMsg` changes anything** |
| Keys | none — the process has exited | `q` and `Ctrl-C` quit. **`Esc` does nothing** |
| Alt screen | left first, if it was ever entered, so the text lands in the real scrollback | still in it; the view is replaced |
| Exit | non-zero, per §8 | **does not exit.** *"The player is probably dragging a divider and will drag it back"* ([`03-responsive.md`](../../blueprint/03-responsive.md) §6) |

This is the reading that makes nav §5's row —
*"`Z-11 Fatal error` → `Esc` does Nothing → there is nothing to go back to; only `q` and
`Ctrl-C` work"* — true of the one variant that persists long enough for a key to matter.
**Flagged in §17.**

### 3.3 · ASCII only, and it is the same argument as the frame

Every glyph the product otherwise loves is a liability here. `▌` U+258C, `─` U+2500 and the whole
box-drawing family are **East-Asian-Ambiguous** (verified with `unicodedata`, Unicode **16.0.0**,
2026-08-25) — so rendering them correctly needs the width-aware measurement of
[`../01-design-system.md`](../01-design-system.md) §1.2, which needs a terminal-capability
detection, which is one more thing that can be what broke. `…` U+2026 is Ambiguous too, which is
why **nothing on this screen is ever truncated.**

ASCII is `Na — Narrow` in every terminal, in every locale, under every `ambiguous-width` setting.
The escape hatch the design system reserves for a terminal that reports nothing — the ratified
ASCII column — is here the **only** vocabulary, and it is not a degrade: it is the correct
rendering.

### 3.4 · The five cases

| # | Case | Mode | Exit | The two facts it must name |
|---|---|---|---|---|
| 1 | the library file cannot be opened or read | EXIT | **1** | the **full path**, and the OS's own reason |
| 2 | the schema was written by a **newer** binary | EXIT | **3** | **both versions** — the file's and this build's |
| 3 | a **migration failed** | EXIT | **4** | which migration and why, and **the `library.db.pre-<version>` backup path** |
| 4 | the terminal is **below `24 × 8`** | **EXIT at start-up · HOLD mid-session** | **2** — fixed by [`03-responsive.md`](../../blueprint/03-responsive.md) §6 | the required size and the actual size |
| 5 | an unexpected **panic** — *added by this spec, §17* | EXIT | **5** | the panic value and the top frame |

**Exit code 2 is fixed by the spine.** Codes 1, 3, 4 and 5 are this spec's proposal — distinct
codes cost nothing and let a script tell "your file is broken" apart from "your terminal is small",
which is the whole reason a text program has exit codes.

---

## 4 · Mockup at 80 columns

`Z-11` is frameless, so there is no ruler-relative inset to check — but the drawings are still
drawn to an exact width. **Every line below is ASCII, starts at column 1 or column 3, and no line
exceeds 72 columns.** The ruler's digit marks columns 1, 11, 21 … 71.

### 4.1 · The fixed shape — every case, four blocks

```text
0.........1.........2.........3.........4.........5.........6.........7.........
zerado: <what broke, lower case, no full stop>

  <the subject: a full path, or the terminal size>
  <the reason, verbatim from the OS or the migration>

<what happened to the player's data - one sentence>
<what to do next - one sentence>
```

`zerado: ` is the classic Unix diagnostic prefix and it is doing three jobs: it names the program
in a scrollback full of other output, it is greppable, and it uses the **command** casing that
`naming.md` fixes for the binary. It is the only hierarchy device available besides case, indent
and blank lines, and on this screen those four are the whole toolkit.

### 4.2 · RENDER 80×24 — case 1, the library file cannot be read

```text
0.........1.........2.........3.........4.........5.........6.........7.........
zerado: cannot open the library file

  /home/alex/.local/share/zerado/library.db
  permission denied

Nothing was read and nothing was changed.
Make the file readable, or set ZERADO_DB to another path.
```

### 4.3 · RENDER 80×24 — case 2, a newer Zerado wrote this file

```text
0.........1.........2.........3.........4.........5.........6.........7.........
zerado: this library was written by a newer Zerado

  /home/alex/.local/share/zerado/library.db
  schema version 7; this build understands version 5

Nothing was read and nothing was changed.
Update Zerado, or point ZERADO_DB at a different file.
```

**Both versions are named**, which is the ADR's own requirement:
*"A database written by a newer binary is a fatal error that names both versions — never a silent
downgrade"* ([`ADR-0001`](../../adr/ADR-0001-zerado-foundational-architecture.md)). Silently
proceeding against a schema you do not understand is how one file becomes two incompatible files.

### 4.4 · RENDER 80×24 — case 3, a migration did not finish

```text
0.........1.........2.........3.........4.........5.........6.........7.........
zerado: a migration did not finish

  /home/alex/.local/share/zerado/library.db
  migration 6 failed: no such column: sort_title

Your library as it was before the migration is beside it, untouched:

  /home/alex/.local/share/zerado/library.db.pre-6

Nothing in that file was changed. Do not delete it.
Report this at github.com/JustCode-CruzAlex/Zerado.
```

**The backup path is named in full**, which is [`09-erd.md`](../../blueprint/09-erd.md) §6's
requirement: *"The file is backed up before any migration that drops or rewrites a column, to
`library.db.pre-<version>` beside it, and `Z-11` names that path if the migration fails."*
`<version>` is **the migration that was being applied** — `pre-6` is the state before 6 ran.

This is the only case with a **five**-block shape, because it has a second path to name, and a
path that is not on screen is a path the player cannot act on.

### 4.5 · RENDER 80×24 — case 5, an unexpected stop

```text
0.........1.........2.........3.........4.........5.........6.........7.........
zerado: stopped unexpectedly

  panic: runtime error: index out of range [12] with length 3
  internal/tui/library.go:214

Every change you had already made was saved when you made it.
Report this at github.com/JustCode-CruzAlex/Zerado.
```

**The data sentence is the product's own invariant, restated at the worst possible moment.**
*"Every mutation Zerado makes is committed when it is made — a status change writes to SQLite
before the overlay closes. There is no unsaved state"*
([`04-navigation-and-focus.md`](../../blueprint/04-navigation-and-focus.md) §3.2). That is why this
sentence can be printed truthfully after a crash, and it is the single most reassuring true thing
the product can say here.

**The panic value and the top frame, and nothing else.** A full stack trace is developer output,
not player output — the same distinction `sync_run.error` draws: *"the classified failure, not a
stack trace"*. Two lines are enough to file a useful report.

---

## 5 · Visual hierarchy

`Z-11` has **no colour, no weight, no box drawing and no glyphs.** Of the five channels
[`../01-design-system.md`](../01-design-system.md) §1.1 names — case, weight, colour role, box
drawing, spacing — **two survive**, plus one this screen adds.

| Rank | Element | Channel | Note |
|---|---|---|---|
| 1 | `zerado: <what broke>` | **position** (column 1, row 1, flush) + the prefix | the only line that is not indented. Flush-left *is* the emphasis here |
| 2 | the indented fact block | **indent** (2 columns) + **spacing** (a blank row above and below) | the path and the reason: the two things that get copied into a shell or a bug report |
| 3 | the data sentence | position (back to column 1) | what happened to the player's file |
| 4 | the next action | position, last | |

**The whole hierarchy is: flush-left outranks indented, and a blank row means "a different
thing."** That is the spacing canon's own `InterElementGap` idea surviving in a screen that has no
tokens — not because a token was applied, but because the idea is right and needs no library.

### 5.1 · Flush against the terminal edge — the one sanctioned exception, and why

R-2 and [`00-design-brief.md`](../00-design-brief.md) §10 line 1 are unambiguous: *"Nothing is
flush against any terminal edge at any tier ≥ Narrow."* `Z-11` is flush at column 1, on every
tier, deliberately.

1. **There is no frame to inset.** `OuterMarginX` is *"the inset between the framed surface and the
   terminal edge"* — with no framed surface, the token has nothing to apply to. This is an
   exemption, not a violation, and §6 records it token by token.
2. **The margin would be spent at exactly the size where columns are scarcest.** Case 4 fires at
   `20 × 6`; two columns of inset is 10 % of that terminal.
3. **A diagnostic is not a screen.** It is a line of program output, and program output starts at
   column 1. `zerado: cannot open the library file` should look exactly like `ls: cannot access` —
   because that is what it is, and looking like it is the point.

---

## 6 · Spacing tokens — **none apply, and here is each one saying so**

The contract asks for every applied token by name. `Z-11` applies **none**, and the honest answer
is the table, not a blank section.

| Token | On `Z-11` | Why |
|---|---|---|
| `OuterMarginX` | **not applied — 0** | there is no framed surface to inset. §5.1 |
| `OuterMarginY` | **not applied — 0** | the message begins on row 1 |
| `InnerPaddingX` | **not applied** | no frame boundary to pad inside of. The 2-column indent of the fact block is **not** this token — it is a literal two spaces, chosen because it is what a diagnostic looks like |
| `InnerPaddingY` | **not applied** | |
| `InterElementGap` | **not applied as a token**; the *idea* survives as one blank line between blocks (§5) | using the token would mean importing the package that computes it |
| `HeaderBandHeight` | **not applied — there is no header band** | no breadcrumb, no title row. `Zerado ✦ Fatal error` is never drawn: a breadcrumb implies somewhere to go back to, and there is nowhere |
| `HeaderBand(tier, false)` | **not applied** | |
| `leftInset` | **not applied — content begins at column 1** | and therefore *"header-left equals content-left"* is vacuously true: there is no header |
| **the reserved footer row** | **not applied — there is no footer** | in `EXIT` mode there are no keys to hint. In `HOLD` mode the advice line is body text, printed only when it fits (§11.2) |
| **wrap width** | `min(terminalWidth, 72)` | the one number this screen computes. **72**, not 68, because `--z-measure` is a *prose* measure and this is diagnostic output; 72 is the conventional terminal-message width and leaves room for a quoting `> ` prefix in a bug report |
| **paths** | **never wrapped, never truncated** | a path broken at a space cannot be copied; a path with `…` in it cannot be used. Paths are printed on their own line and the terminal soft-wraps them. This is the only place hard-wrapping yields |

**Paths are printed in full and are never tilde-shortened.** `Z-09 Settings` shortens
`~/.local/share/...` because it is being read; `Z-11` prints the absolute path because it is being
**copied into a shell**.

---

## 7 · Colour, glyph and label — the whole table is empty, and that is the spec

| Element | Colour token | Hex | ANSI-256 | 16-colour | Glyph | Label | `NO_COLOR` |
|---|---|---|---|---|---|---|---|
| every line, every case, every mode | **none — no SGR is emitted** | — | — | — | **none — ASCII only** | the words themselves | **identical, always** |

**Co-render holds trivially and completely.** The co-render rule is *colour and glyph and label*;
`Z-11` carries **label only** — and that is the leg the rule exists to protect. *"A TTY has no
accessibility API… the only way state survives is if it is written as words"*
([`00-design-brief.md`](../00-design-brief.md) §3.2, SC 4.1.2). `Z-11` is the product's proof that
the words alone are sufficient, because on this one screen the words are all there is.

**The 16-colour floor, `NO_COLOR`, a monochrome terminal, a screenshot in a bug report, a
`ambiguous-width=double` CJK terminal, a terminal that reports no capabilities at all, and a pipe
into a file — `Z-11` renders identically in all eight.** No other screen in the product can say
that, and it is the reason this one is built the way it is.

---

## 8 · The full state table

| # | State | Mode | Exit | Renders | Copy |
|---|---|---|---|---|---|
| 1 | **library unreadable** | EXIT | 1 | §4.2 | §10 · **C1** |
| 2 | **newer schema** | EXIT | 3 | §4.3 | §10 · **C2** |
| 3 | **migration failed** | EXIT | 4 | §4.4 | §10 · **C3** |
| 4a | **terminal below `24 × 8` at start-up** | EXIT | **2** | one line | §10 · **C4** — ratified, verbatim |
| 4b | **terminal dragged below `24 × 8` mid-session** | **HOLD** | — | §11.2 | **C4** + the advice line **when it fits** |
| 4c | **the terminal is dragged back above the floor** | — | — | the previous screen returns, **with its focus and scroll intact** | — |
| 5 | **panic** | EXIT | 5 | §4.5 | §10 · **C5** |
| 6 | **First run** | either | | **N/A, and it matters:** `Z-11` has no first-run variant, because it reads nothing and has no state that could be uninitialised. Case 1 at first run says the same sentence it says at the thousandth run — and *"Nothing was read and nothing was changed"* is true either way | — |
| 7 | **Loading** | — | | **N/A.** `Z-11` fetches nothing. There is no frame in which it could be loading | — |
| 8 | **Empty / partial** | — | | **N/A.** No data set | — |
| 9 | **Offline** | either | | **no difference whatsoever.** `Z-11` is `WORKS` and reaches for nothing | — |
| 10 | **Audio** | either | | **silent, in every mode, in every case.** No indicator, no `m`, no cue — there is no footer to carry an indicator and no running screen to mute | — |
| 11 | **`NO_COLOR`** | either | | **identical.** §12 | — |
| 12 | **`ZERADO_ASCII`** | either | | **identical** — `Z-11` is already ASCII (§3.3) | — |
| 13 | **stdout redirected to a file** | EXIT | | the message still appears **on the terminal**, because it goes to **stderr** | — |

---

## 9 · The key map

### 9.1 · `Z-11`-EXIT — there is no key map

The process has ended. Nothing is listening. This is not an omission: **an input loop is a
dependency on the input subsystem**, and §3.1's whole argument is that this screen depends on
nothing that could be what broke.

### 9.2 · `Z-11`-HOLD

| Key | Does |
|---|---|
| `q` | quit |
| `Ctrl-C` | quit |
| **`Esc`** | **nothing.** *"There is nothing to go back to"* (nav §5) |
| every other key | nothing, and **shows nothing** — no beep, no flash, no error line. The same treatment §3.1 of the nav model gives the reserved-but-unbound keys |
| `?` `,` `m` `r` `s` `a` `/` | **unbound.** Help, Settings and mute are all screens, and there is no screen to return to |
| — | the only *message* that changes anything is `WindowSizeMsg` |

**Not a keyboard trap** (2.1.2): no element takes focus, and `q` and `Ctrl-C` both leave. The
advice line names `q` whenever there is room for it (§11.2); where there is not, `Ctrl-C` is a
standard exit method and needs no advice.

---

## 10 · The exact copy — ready to paste

**C1 · the library file cannot be opened or read** — exit `1`
```text
zerado: cannot open the library file

  <full path>
  <the OS error, verbatim>

Nothing was read and nothing was changed.
Make the file readable, or set ZERADO_DB to another path.
```

**C2 · the schema is from a newer Zerado** — exit `3`
```text
zerado: this library was written by a newer Zerado

  <full path>
  schema version <file>; this build understands version <binary>

Nothing was read and nothing was changed.
Update Zerado, or point ZERADO_DB at a different file.
```

**C3 · a migration did not finish** — exit `4`
```text
zerado: a migration did not finish

  <full path>
  migration <n> failed: <the error, verbatim>

Your library as it was before the migration is beside it, untouched:

  <full path>.pre-<n>

Nothing in that file was changed. Do not delete it.
Report this at github.com/JustCode-CruzAlex/Zerado.
```

**C4 · the terminal is below the floor** — exit `2` at start-up.
**Ratified by [`03-responsive.md`](../../blueprint/03-responsive.md) §6. Reproduce verbatim.**
```text
Zerado needs at least 24 columns and 8 rows. This terminal is 20 x 6.
```
In `HOLD` mode, and **only when it fits entirely in the rows left after a blank line** (§11.2):
```text
Make the window bigger. q quits.
```

**C5 · an unexpected stop** — exit `5`
```text
zerado: stopped unexpectedly

  <the panic value>
  <the top frame: file:line>

Every change you had already made was saved when you made it.
Report this at github.com/JustCode-CruzAlex/Zerado.
```

### 10.1 · Two recorded copy decisions

**C4 does not carry the `zerado:` prefix**, and every other case does. It is reproduced verbatim
because [`03-responsive.md`](../../blueprint/03-responsive.md) §6 fixes it, and it reads correctly
as prose because it names the product rather than the command — which is the right casing for a
sentence a person reads (`naming.md`). **Flagged in §17** in case the founder prefers one shape for
all five.

**The repository URL is the ratified one**, `github.com/JustCode-CruzAlex/Zerado`, verified public
and answering `200` anonymously (`ratification/decisions.md` Q2). No issue-tracker path is
asserted, because the repository URL is the fact that was checked.

**Voice check.** Every case names **what happened**, **why**, **what happened to the player's
data**, and **the next action** — the four things
[`../01-design-system.md`](../01-design-system.md) §11.3 requires. No exclamation marks. No emoji.
The reader is never called a "gamer". **No case says *"Something went wrong"*** — which is the one
sentence a terminal user cannot act on, and the sentence this screen exists to never print.

---

## 11 · Behaviour at 40 columns and below, and the refusal floor

`Z-11` **is** the refusal floor's screen. It has no tiers, no breakpoints and no degrade table:
the same five cases render the same five messages at every width, hard-wrapped to
`min(terminalWidth, 72)` at word boundaries — with paths exempt (§6).

### 11.1 · RENDER 24×8 — case 1 at the narrowest terminal Zerado will still render

```text
0.........1.........2...
zerado: cannot open the
library file

  /home/alex/.local/shar
e/zerado/library.db

Nothing was read and
nothing was changed.
Make the file readable,
or set ZERADO_DB to
another path.
```

**The path overruns 24 columns and the terminal soft-wraps it — drawn above exactly as it
lands, mid-word, on the second row. That is correct.** A path broken
at a space cannot be pasted; a path ending in `…` cannot be used. And the message running past
eight rows is correct too: `EXIT` mode has already left the alternate screen, so the text lands in
the real scrollback where it can be scrolled to and copied. **A diagnostic is allowed to be longer
than the window. A view is not** — which is exactly the distinction §3.2 draws between the modes.

### 11.2 · RENDER 20×6 — case 4, `HOLD` mode, the smallest thing that renders

```text
0.........1.........
Zerado needs at
least 24 columns and
8 rows. This
terminal is 20 x 6.
```

Four lines of the ratified sentence, hard-wrapped at 20. **Two rows remain and the advice line
needs three** — a blank plus `Make the window bigger. q quits.` wrapped to two lines — so it is
**omitted**. The rule, stated so it can be implemented: *the advice line is printed only when a
blank line plus the whole wrapped advice fits in the rows remaining.* The refusal sentence is never
pushed off the screen to make room for advice about it.

At `30 × 8` the sentence wraps to three lines and the advice fits:

```text
0.........1.........2.........
Zerado needs at least 24
columns and 8 rows. This
terminal is 30 x 8.

Make the window bigger. q
quits.
```

### 11.3 · Below `20 × 6`, and below that again

There is no lower floor. The sentence keeps wrapping and the terminal keeps scrolling. At a width
where even a single word does not fit, the terminal breaks it — and that is still more useful than
anything Zerado could do about it. **`Z-11` never refuses to render**, because it is what the
refusal *is*.

---

## 12 · `NO_COLOR` rendering — shown, not asserted

```text
0.........1.........2.........3.........4.........5.........6.........7.........
zerado: cannot open the library file

  /home/alex/.local/share/zerado/library.db
  permission denied

Nothing was read and nothing was changed.
Make the file readable, or set ZERADO_DB to another path.
```

**Byte-identical to §4.2 — because `Z-11` never emits an SGR sequence in the first place.**

`NO_COLOR` set or unset, `ZERADO_ASCII` set or unset, `ZERADO_REDUCED_MOTION` set or unset,
truecolor or 256 or 16 colours or none, a terminal that reports its capabilities or one that
reports nothing: **one rendering, eight environments.** There is no second render to compare
against, no channel to lose, and nothing to check — which is what makes this the one screen whose
`NO_COLOR` conformance is true by construction rather than by test.

The colour-budget checklist ([`../02-colour-budget.md`](../02-colour-budget.md) §10) passes line by
line, vacuously and on purpose: chrome cyan **0** · every other cyan **0** · amber **0** · red
**0** · no region separated by fill · every control boundary — there are no controls · the focus
ring — there is no focusable element · co-render — **label only, and the label is the whole
screen**.

---

## 13 · Focus model on this screen

| | |
|---|---|
| **Regions** | **0** |
| **Focusable items** | **none**, in either mode |
| **Focus ring** | **none is drawn, because nothing is focusable.** This is not a removal — [`../02-colour-budget.md`](../02-colour-budget.md) §8.2 governs *the focused element*, and there is no element |
| **`Tab`** | unbound |
| **`Esc`** | **nothing**, in both modes. Nav §5, verbatim: *"There is nothing to go back to."* And the screen **does not hint `esc`**, because hinting a key that does nothing is a footer that lies — in `EXIT` mode there is no footer at all |
| **The way out** | `EXIT`: the process has already ended. `HOLD`: resize, or `q`, or `Ctrl-C` |
| **2.1.2 No Keyboard Trap** | satisfied — nothing takes focus, and two standard exits work |
| **2.4.11 Focus Not Obscured** | vacuous — nothing is focused and nothing overlays anything |

---

## 14 · Colour budget declaration

| Class | Count | Where |
|---|---|---|
| **CHROME cyan** | **0** | |
| **STATE cyan** | **0** | |
| **Focus-ring cyan** | **0** | nothing is focusable |
| **Text-cursor cyan** | **0** | no input |
| **Amber** | **0** | not even the terminal mark. `[0] ZERADO` belongs to the first-run splash and to `--help` ([`../01-design-system.md`](../01-design-system.md) §1.6); putting the brand on a crash screen would be the brand taking a bow at the worst possible moment |
| **Red** | **0** | **and this is deliberate.** `--z-scanner` is on the closed red list for error text's sake, and this is the product's most serious error — but rendering it red would require an SGR sequence, which requires knowing the terminal's depth, which is a dependency. **The words carry the alarm.** Red would add nothing a reader of `zerado: cannot open the library file` does not already have |
| **Total SGR sequences emitted** | **zero** | in every case, in both modes, in every environment |

**Amber ceiling:** 0 of 1920 → **0 %**.

---

## 15 · Reuse verdict per element

| Element | Verdict | Note |
|---|---|---|
| The whole screen | **Build fresh — and build it in one file with no imports beyond `fmt`, `os` and `strings`** | this is the reuse verdict. Every dependency is a liability here, and the screen is about forty lines |
| Word wrapper | **Build fresh** — a `strings.Fields` loop | `lipgloss.NewStyle().Width()` would wrap correctly and would also import the layout engine. Forty lines against a dependency the screen exists to avoid |
| Header band / footer / frame | **Not used** | §6 |
| `bubbles/*` · `glamour` · `huh` · `harmonica` · `charmbracelet/log` | **Not used** | `log` is worth naming twice: its role is structured **developer** logging (§11.5), its level colours are its own palette, and it would emit SGR |
| Audio | **Not used** — no `Cue()`, in any mode | §3.1 |
| `HOLD` mode's view | **Build fresh** — a `View()` that returns the same plain string | it runs inside the existing Bubble Tea program, which is alive by definition in that mode; it still touches no theme, no frame and no colour |

---

## 16 · Screen-specific acceptance criteria

These are unusually mechanical, because this screen's whole value is that it cannot fail.

1. **Zero SGR bytes.** Capture the output of all five cases, with and without `NO_COLOR`, and grep
   for `\x1b[`. There must be **no match**, in any case, in either mode.
2. **Zero non-ASCII bytes.** `LC_ALL=C grep -P '[^\x00-\x7F]'` over every capture must be empty.
3. **The dependency prohibition is a grep-level review rule.** The fatal-path package must not
   import `lipgloss`, the theme, `space`, `frame`, `chrome`, the store, the vault, `net/http`, or
   the audio seam. This is the same shape of rule
   [`07-offline-contract.md`](../../blueprint/07-offline-contract.md) §7 uses to keep the offline
   classes from drifting, and for the same reason: it survives a year of changes.
4. **The crash screen cannot crash.** Inject a panic into the theme, the frame **and** the store
   packages in turn; each must still produce case 5's message and exit `5`.
5. **It renders with the terminal reporting nothing.** `TERM=` unset, no size reported, stdout not
   a TTY — the message still appears, in full.
6. **It renders when the library is unreadable**, which is case 1's own precondition — proving the
   fatal path reads nothing.
7. **Exit codes** are 1 / 2 / 3 / 4 / 5 as §3.4, and **2 for the below-floor start-up case matches
   [`03-responsive.md`](../../blueprint/03-responsive.md) §6 exactly**.
8. **stderr, not stdout, in `EXIT` mode.** `zerado > /dev/null` still shows the message.
9. **Case 2 names both versions.** Case 3 names the migration number **and** the
   `library.db.pre-<n>` path, in full and unshortened.
10. **No path is ever wrapped or truncated.** Grep every capture for a path containing a space or
    a `…`; there must be no match.
11. **`HOLD` does not exit.** Drag a running terminal below `24 × 8`, then back above it: the
    previous screen returns **with its focus and scroll position intact**, and the process ID is
    unchanged.
12. **`HOLD` never pushes the refusal sentence off screen** to make room for the advice line
    (§11.2). Test at `20 × 6`, `24 × 8` and `30 × 8`.
13. **No sound is played**, in any case, in either mode.
14. **No case contains the string "Something went wrong"**, and every case contains a sentence
    naming what happened to the player's data.
15. Artifacts: all five cases at `80 × 24`; case 1 at `24 × 8`; case 4 `HOLD` at `20 × 6`, `24 × 8`
    and `30 × 8`; and case 1 again with `NO_COLOR=1`, which must be **byte-identical** to the
    capture without it.

---

## 17 · Open for the founder

1. **The two modes (§3.2).** `EXIT` prints and leaves; `HOLD` keeps running and waits for a resize.
   This is what reconciles nav §5's *"only `q` and `Ctrl-C` work"* with
   [`03-responsive.md`](../../blueprint/03-responsive.md) §6's *"keeps running… exit is only for the
   start-up case"*. Confirm the split, and let `fft-tui-architect` name it in the `Esc` table.
2. **Case 5, the panic handler, is added by this spec.** The ticket named four cases. A program
   that stops without saying so is worse than any of the four, and *"a crash screen that itself
   crashes tells the player nothing"* only means something if there is a handler to be robust.
   Confirm the fifth case, its copy, and exit code `5`.
3. **Exit codes 1, 3, 4 and 5 are proposed here (§3.4).** `2` is the spine's. Distinct codes cost
   nothing and let a script tell a broken file from a small terminal. Confirm, and put them in
   `--help`.
4. **C4 is the only case without the `zerado:` prefix (§10.1)**, because the spine fixes its
   wording. Confirm the inconsistency is acceptable — or restate C4 with the prefix, which is a
   one-line change to `03-responsive.md` §6 and would make all five cases one shape.
5. **The schema version is named, not the binary version (§4.3).** `schema_migration(version,
   applied_at)` does not record which build applied it, so *"which version wrote it"* can only be
   answered as a schema number. Naming the **binary** too would need a `written_by` column added in
   Phase 1 — one column, decided now or not at all. Route to `fft-tui-architect`.
