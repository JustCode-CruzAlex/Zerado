---
title: Zerado — Naming Rationale
discipline: BRAND IDENTITY
dwg-no: ZRD-BRAND-02
rev: A
date: 2026-08-24
status: ratified
subtitle: Why the name is Zerado, what it cost, and what we agreed to carry
---

# Zerado — the name

**Locked in ratification round 1.** The name goes in the headline, the browser
tab, the logo, the binary, the domain, and the command people type.

---

## What the word means

*Zerado* is Portuguese. Applied to a game it means **beaten — finished, cleared,
100%'d.** A player says *"zerei"* the way an English speaker says *"I finished
it"*, except it carries more: not merely reaching the credits but having
**closed the thing out**.

The origin is the part that matters here. It comes from **arcade cabinets** —
the score counter running past its last digit and rolling back to zero. You had
beaten the machine so thoroughly that its own display gave up. The word is not a
metaphor borrowed from the eighties; it is **literally an artifact of the
eighties**, minted by hardware that no longer exists.

That is why the identity is built on it rather than merely citing it. A counter
rolling to zero is the mark (see `brand-manual.md` §3). The name is not
decoration on top of the design — it *is* the design's argument.

---

## Why it won

### 1. Personality

The alternatives described a feature. This one describes **a feeling a player
already has a word for**. It is short, it is spoken, and it is not a compound of
two English nouns that a hundred other tools have already compounded.

### 2. The origin is the product

Zerado tracks a library and asks what you should play tonight. Its emotional
payload is the moment a game moves into the *finished* column. The name is that
moment. No other candidate had a story that resolved into the product's own
core interaction.

### 3. Typing comfort — a real criterion for a terminal tool

This is a program people invoke by hand, many times a day. That makes the
keyboard a design surface.

| Property | Result |
|---|---|
| Length | 6 letters |
| Shift key required | No |
| Digits or punctuation | None |
| Hand alternation | z-e-r-a-d-o alternates left/right cleanly |
| Ambiguous letters (l/1, O/0) | None |
| Reads as a command | `zerado sync` · `zerado tonight` · `zerado mark 42 zerado` |

`zerado sync` reads like a real command from a real tool — the same shape as
`git push` or `brew update`. That was the test, and it passed.

---

## The runners-up — recorded honestly

Two other names reached the final round. Neither was rejected for being bad.

### `gigawatt`

**For:** an unambiguous *Back to the Future* citation, instantly legible to the
target audience, easy for English speakers to say and spell.

**Against:** it is a **borrowed** eighties reference rather than an **owned**
one — a quotation from a film, not a word from the domain. It says nothing
about games, libraries, or finishing anything. It also invites a licensing
conversation nobody wants, and "1.21 gigawatts" is a joke the internet has
already used to exhaustion.

### `backlit`

**For:** genuinely good English. Evokes a CRT, a keyboard, a screen in a dark
room. Short, typeable, pronounceable, and quietly on-theme.

**Against:** it is **atmosphere without meaning**. It describes the lighting of
the room the product runs in, not the thing the product does. It is also a
common adjective in consumer electronics, which makes it weak to own and weak
to search. And it has no second life as a verb or a status.

**Why `zerado` beat both:** it is the only candidate that is simultaneously the
product's *name*, the product's *core action*, and an *authentic eighties
artifact*. The other two had one of those at most.

---

## The risks we accepted

These were known before the decision and were **accepted, not overlooked**. They
are recorded here so nobody re-litigates them later or is surprised by them.

### 1. English speakers cannot pronounce or guess it

An English speaker seeing `zerado` cold does not know whether it is
*zeh-RAH-doo*, *ZEH-ra-do*, or *zuh-RAY-doh*, and cannot infer the meaning at
all. This is a genuine cost at the top of the funnel.

> **Accepted because:** the launch audience is small, technical, and reads
> before it clicks. The landing page has room to explain the word once, in a
> clause, at the moment it first appears — which the ratified copy plan does.
>
> **Mitigation, and it is real:** the page must **never** use the word as if it
> were self-evident. First contact always carries the gloss.
>
> Reference pronunciation: **zeh-RAH-doo** (Brazilian Portuguese).

### 2. It collides with its own in-app status

The product is called Zerado, and marking a game *zerado* is a thing you do
inside it. `zerado mark 42 zerado` is a real command.

> **Accepted because this is a feature, not a defect.** A tool named after the
> act it exists to enable is a tool whose name explains itself the moment you
> use it. The landing page turns the collision into the payoff line rather than
> hiding it.
>
> **The one rule that keeps it from becoming confusing:** the product is always
> **`Zerado`** (capitalised) and the status is always **`zerado`** (lowercase,
> and in prose, italicised on first use). See the casing convention below.

### 3. It is a common Portuguese word, so it is hard to search for

"zerado" returns bank balances at zero, reset counters, football scores, and
general Portuguese usage. SEO for the bare term is effectively unwinnable in
Portuguese-speaking markets.

> **Accepted because:** the domain is `zerado.app`, which is exact-match and
> unambiguous; the launch channel is direct and community-led, not search; and
> the useful long-tail queries are compound anyway — *"zerado game tracker"*,
> *"zerado terminal"*, *"zerado backlog"* — where there is no competition at all.
>
> **Consequence for the site:** the page must always pair the name with a
> qualifier in titles and metadata — **"Zerado — terminal game library"**, never
> a bare "Zerado".

---

## Collision check

Recorded from the naming round:

| Checked for | Result |
|---|---|
| An existing **game** named Zerado | None found |
| An existing **game tracker / backlog tool** named Zerado | None found |
| A **relevant trademark** in the software or games class | None found |
| The domain `zerado.app` | Available and **verified resolving** |
| The repository `github.com/JustCode-CruzAlex/Zerado` | **Verified public, HTTP 200 anonymously** |

> **Stated plainly, because it matters:** this was a *conflict scan*, not a legal
> clearance. It is enough to proceed with a landing page and a public
> repository. It is **not** a registrability opinion, and it should not be
> treated as one. Before the name carries revenue — and as of 2026-08-25 there
> is none, the affiliate model having been dropped — it needs a real trademark
> search by someone qualified to give that opinion, in the classes and
> territories that will actually matter.

---

## Casing convention — settled

This has three answers because there are three different objects. Getting them
confused is the single most likely way to make the name read as sloppy.

| Context | Form | Example |
|---|---|---|
| **The product**, in prose and titles | `Zerado` | "Zerado reads your Steam library." |
| **The command**, the binary, the domain, any code | `zerado` | `zerado sync` · `zerado.app` |
| **The status** inside the product, in prose | `zerado`, italic on first use | "Six games are *zerado* this year." |
| **The status** in the interface | `ZERADO` | the state chip, the filter, the column |
| **The wordmark** | `ZERADO` | the logotype only — see below |

### On the wordmark being uppercase

The logotype is drawn as **`ZERADO`** in constructed geometric capitals. This
does **not** license uppercase in prose. A logotype is a drawn object with its
own rules; it is not an instruction about how to write the word.

Two reasons it is capitals:

1. **Period correctness.** Wide, tightly-drawn, generously tracked capitals are
   the eighties corporate-future register — the DeLorean's own logotype, the
   Knight Industries badge, every dashboard readout of the era. Lowercase would
   read as 2016 startup, not 1984 machine.
2. **Construction.** The wordmark is built from the same rounded-rectangle
   skeleton as the mark, and the final **O** of `ZERADO` is geometrically the
   same shape as the mark itself. That rhyme only works in capitals, and it is
   what makes the lockup read as one designed system rather than an icon placed
   next to some type.

### Never

- **`ZeRaDo`**, **`zeraDO`**, or any camel/stylised casing.
- **`Zerado.app`** — the domain is lowercase, always.
- Writing the product as `zerado` in prose, or the command as `Zerado` in a
  code block. The distinction is the whole mechanism that stops the
  name/status collision from becoming ambiguity.
