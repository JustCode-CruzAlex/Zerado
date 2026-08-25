---
title: Zerado — Landing Page Copy
discipline: CONTENT
doc-no: ZRD-CONTENT-01
rev: A
date: 2026-08-24
status: final — implement verbatim
source: forgeplay-output/landing-page/ratification/decisions.md, mock.outline.md, brief.md
        forgeplay-output/landing-page/brand/brand-manual.md, naming.md
---

# Zerado — landing page copy

**This is final copy, not a draft.** Every string below is labelled so it can be wired
without guessing. Sixteen sections, in the ratified order — nothing added, nothing dropped,
nothing reordered. Language: English, with `zerado` and `sinopse` kept in Portuguese and each
explained once, in a clause, the first time a reader meets it (Q5-b).

Casing, per `naming.md`: the product is always **`Zerado`**; the command/domain is always
lowercase `zerado`; the in-app status, in prose, is lowercase and *italicised on its first
appearance on the page* — that first appearance is in §05, below — and plain lowercase after.

No dates, no version numbers, no download links, no prices, no user counts, no testimonials
appear anywhere in this document. The roadmap (§12) marks every phase **Planned**. Nothing is
downloadable yet — the CTA is a waitlist, everywhere it appears.

---

## 01 · hero

**Headline:**
> Pick a mood. Get a game. Play it tonight.

**Subhead:**
> `Zerado` reads your entire library — Steam first, then PlayStation, GOG, EA, and the
> cartridges still on your shelf — and answers the only question that actually matters at
> 11pm: what do I play? It runs in your terminal, works with the network off, and everything
> it knows lives in one file on your machine.

**CTA primary:** Join the waitlist
→ `mailto:alex@flowforgesoft.com?subject=Zerado-WaitList`

**CTA secondary:** View on GitHub
→ `https://github.com/JustCode-CruzAlex/Zerado`

**Microcopy, under the CTA pair:**
> Opens your mail client, addressed and subject-lined already. No form, no account.

**Image slot — full width, below the fold line:**
The program's library view, running.
**Alt text:** "An illustration of Zerado's terminal library view. This is a mockup of the
planned interface, not a screenshot — the program isn't runnable by visitors yet."

---

## 02 · maker-line

**Body — one line, no heading needed:**
> `Zerado` is the second product from FlowForgeSoft, the company behind FlowForge.

---

## 03 · the-problem

**Heading (optional, small):** The problem

**Body — three short lines, no imagery:**
> You own 400 games.
> You've finished six.
> Most nights, you scroll for twenty minutes and start something you've already seen.

---

## 04 · moods

**Heading:** Sorted by mood, not genre

**Intro line:**
> Genre tells you what a game is. Mood tells you what tonight is for. `Zerado` sorts by the
> second one.

**Card 1**
- Title: Mindless grind
- Body: Repetitive, low-stakes, doesn't ask for your attention. For the night you want to
  zone out, not think.

**Card 2**
- Title: Story rich, kind of sad
- Body: Something to feel, not just clear. For the night you want to sit with a story, not
  race through one.

**Card 3**
- Title: Quick fifteen minutes
- Body: Starts fast, ends clean, fits in the window you actually have before bed.

**Card 4**
- Title: Tactical, full focus
- Body: Wants your whole attention and rewards it. For the night you have the attention to
  spend.

**Image slot:**
The mood picker on screen.
**Alt text:** "Mockup of Zerado's mood picker in the terminal, showing mood categories —
including 'Mindless grind' and 'Quick fifteen minutes' — next to the games that match each
one."

---

## 05 · states

**Heading:** Where each game stands

**Intro line:**
> Four states, one for every game in the library.

**State 1**
- Glyph + label: `○` NOT STARTED
- Line: In the library. Untouched.

**State 2**
- Glyph + label: `◐` IN PROGRESS
- Line: Currently being played.

**State 3 — the payoff line**
- Glyph + label: `◉` ZERADO
- Line: Finished. *`Zerado`* is Portuguese for beaten — cleared, 100%'d — from the arcade
  counter that ran past its last digit and rolled back to zero. The product is called
  `Zerado`. Marking a game *zerado*, in `Zerado`, is the whole point of it.

**State 4**
- Glyph + label: `⊘` ABANDONED
- Line: Started, set down, not coming back to it.

---

## 06 · one-collection

**Heading:** One collection, everything you own

**Connected stores:**
> Steam syncs today — your library, playtime and status, kept current on its own. PlayStation,
> GOG and EA are planned next: connect an account, see it in the same list.

**The physical shelf:**
> Add the discs and cartridges next to your desk by hand — title, platform, done. Each one
> gets the same mood tag and status as anything synced automatically, and, once cover art and
> *sinopse* — Portuguese for synopsis, a short plot summary — are wired in, the same page as
> everything else. A physical copy isn't a second-class row in the list.

**Image slot:**
Cover-art grid, mixed synced and hand-added titles, one marked Physical.
**Alt text:** "Mockup of a cover-art grid showing synced and manually-added games side by
side, with one cover labelled 'Physical' to show a hand-entered disc or cartridge sitting in
the same collection as everything else."

---

## 07 · price-intelligence

**Heading:** When to buy, not just what to buy

**Intro line:**
> Say what you can spend. `Zerado` checks today's price against the lowest that game has ever
> gone, and tells you plainly whether today is a good day to buy.

**Worked example — labelled as an example, not live data:**
> A game already in your library is $15 today — inside budget. Its all-time low is $8, three
> months ago. The verdict: **"$15, and in budget. It's been $8. Maybe wait."**

**Caption under the example:**
> Illustrative numbers. Price data comes from IsThereAnyDeal once `Zerado` is connected —
> nothing above is a live quote.

**Image slot:**
A single price card with its history.
**Alt text:** "Illustration of a price card: current price $15, all-time low $8 reached three
months ago, verdict 'maybe wait.' Example numbers for illustration, not a live quote."

---

## 08 · yours-and-offline

**Heading:** It's yours, and it's offline

**Three points:**
> One file on your machine. Your whole library lives in a single SQLite file you can back up,
> move, or delete — nobody else holds a copy.

> No account to use it. Nothing to sign up for, nothing to log into, before it works.

> Runs with the network off. Once your library is synced, `Zerado` doesn't need the internet
> to answer what to play tonight.

---

## 09 · terminal-first

**Heading:** Built for the terminal, on purpose

**Body:**
> It's a text program. It starts instantly, holds your whole library in memory, and fetches
> from every connected source at once instead of making you wait on each one in turn. No
> mouse required, no page to load, no telemetry running in the background.

**Image slot:**
A second screenshot — one game's detail page.
**Alt text:** "Mockup of a game's detail page in the terminal, showing its status, mood tags
and price history in one screen."

---

## 10 · and-on-your-phone

**Heading:** And on your phone, later

**Body:**
> Native iOS and Android apps are planned for Phase 4 — the same library, the same states, the
> same mood sort, in your pocket instead of your terminal. The old form and the new form, one
> product.

**Image slot:**
Two phone frames.
**Alt text:** "An illustration of two phone frames representing the planned iOS and Android
apps. Not built yet — Phase 4."

---

## 11 · community

**Heading:** The community layer — Phase 4, not built yet

**Body:**
> Comments, reviews and public profiles. These run on servers, and servers cost money — so
> this part will need a premium account or a donation once it exists. Nothing is decided about
> price, tiers, or a date, and nothing on this page is asking for either. This section shows
> what's planned, stamped for what it is: a later phase.

**Image slot:**
Example community screens, watermarked as a later phase.
**Alt text:** "Example mockup of the community layer — a comment thread, a review, and a
public profile — watermarked 'Phase 4, not available yet.'"

---

## 12 · roadmap

**Heading:** Roadmap

**Intro line:**
> Four phases, in order. Phase 1 is under way — nothing here is marked done, because that
> would be the one dishonest line on this page.

**Phase 1**
- Number/name: Phase 1 — CLI/TUI MVP
- Line: Your library, your statuses, stored locally.
- Status: In progress

  *(Changed from "Planned" on 2026-08-25, on the founder's instruction — "we need
  to show Phase 1 'In progress'" — when Phase 1 work started. In progress is not
  done: no phase is marked done, and nothing on the page claims a capability
  exists. The marker co-renders colour + glyph + label per brand-manual.md §3,
  using the same `--z-state-in-progress` token and `◐` glyph as the product's own
  IN PROGRESS game state.)*

**Phase 2**
- Number/name: Phase 2 — Enrichment
- Line: Covers, synopsis, prices, moods.
- Status: Planned

**Phase 3**
- Number/name: Phase 3 — Recommendations & Budget
- Line: What to buy, and whether to wait.
- Status: Planned

**Phase 4**
- Number/name: Phase 4 — Social & Mobile
- Line: Sync, community, the phone apps.
- Status: Planned

*(No dates render on any phase — undated, ordered phases with a status marker only,
until told otherwise. Whether the roadmap carries dates is a separate open question
and has not been answered.)*

---

## 13 · after-phase-4

**Heading:** After Phase 4 — ideas, not promises

**Body:**
> A few directions under discussion, none scheduled and none committed: more storefronts
> beyond the four above, deeper backlog analytics, and mood tagging that improves the more you
> use it. Nothing here has a date, and nothing here is a promise.

---

## 14 · faq

**Heading:** Questions

**Q: Is it free?**
> Yes. `Zerado` itself costs nothing. The only money in the picture, both disclosed below: an
> affiliate commission when a price link leads to a purchase, and, once the Phase 4 community
> layer exists, an optional premium account or donation to help cover its server cost.

**Q: Do I need anything from Steam? Does my profile have to be public?**
> A free Steam API key, which takes a minute to generate. And yes — Steam only returns a
> library if the profile is public, or if the key belongs to that same account. If a profile
> is private, Steam won't share the list; that's a Steam setting, not something `Zerado` can
> work around.

**Q: Does my data leave my machine?**
> Your library lives in one SQLite file on your machine, and nothing needs an account to use
> it. The only network traffic is `Zerado` reaching out to the services you've connected —
> Steam, price data — using your own keys. Nothing about your library is sent to a
> `Zerado`-run server, because there isn't one.

**Q: Which stores work today, and which are later?**
> Steam is built and works. PlayStation, GOG and EA are planned, not yet built. Physical discs
> and cartridges are entered by hand from day one — no store integration required for those.

**Q: Can I run my own server?**
> There's no server to run. `Zerado` is local-first: it reads your library and stores the
> result in one file on your own machine. The one thing that will eventually run centrally is
> the Phase 4 community layer — comments, reviews, public profiles — because sharing a comment
> with someone else requires a server somewhere.

**Q: How is this paid for?**
> Two ways, both disclosed and neither a fee to use `Zerado` itself: an affiliate commission
> when a price link you follow leads to a purchase, and, starting in Phase 4, an optional
> premium account or donation toward the cost of running the community servers.

---

## 15 · closing-cta

**Line above the button:**
> Nothing to install yet. Tell us you're in, and we'll tell you when there is.

**CTA:** Join the waitlist
→ `mailto:alex@flowforgesoft.com?subject=Zerado-WaitList`

**Microcopy:**
> Opens your mail client, addressed and subject-lined already. No form, no account.

---

## 16 · footer

**Contact line:**
> Contact: `alex@flowforgesoft.com`

**Affiliate disclosure:**
> `Zerado` earns a commission when you buy a game through a price link on this page. No cost
> to you, and no other funding ask anywhere on this page.

**Company line:**
> A FlowForgeSoft product.

**Mark caption, under the FlowForge logo:**
> Powered by FlowForge.

**Repeated links (footer nav, no new copy):**
> Join the waitlist · View on GitHub
