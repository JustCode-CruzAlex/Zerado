---
title: Zerado — IGDB cover art on the landing page, licence finding
discipline: LEGAL / SOURCING
doc-no: ZRD-LEGAL-01
rev: A
date: 2026-08-25
status: recorded — basis for a founder decision taken the same day
ticket: "#16"
---

# May IGDB cover art appear on the Zerado landing page?

**Read the sources on 2026-08-25. This document records what was actually read,
what it says, what it does not say, and what was decided on that basis.** It is
a record, not an argument. A future reader should be able to see the ground
that was stood on rather than infer it.

---

## 1 · The question, stated narrowly

Two questions were open, and they are not the same question:

1. **Does IGDB's free tier apply to Zerado** — free software, open source,
   donation-supported, no affiliate, generating no revenue?
2. **May cover art served by IGDB appear on a *marketing website*, as distinct
   from inside the application?** Those can be different permissions in the same
   terms document and were deliberately not conflated.

Question 2 is the one this ticket turned on. It is answered in §4.

---

## 2 · Sources read, with retrieval notes

| Source | URL | Read | Note |
|---|---|---|---|
| IGDB API documentation — *Getting Started*, *Images*, *Business FAQ* | `https://api-docs.igdb.com/` | 2026-08-25 | Fetched directly. |
| Twitch Developer Services Agreement ("TDSA") | `https://legal.twitch.com/en/legal/developer-agreement/` | 2026-08-25 | Page states **"Last modified on 12/04/2024"**. `twitch.tv/p/en/legal/developer-agreement/` 302-redirects here. |
| IGDB Content Policy — *Image Policy* | `https://www.igdb.com/content-policy` | 2026-08-25 | The live page is a client-rendered SPA and returns 403 to non-browser clients; the text below was read from the Internet Archive snapshot `20260731220557` of that URL. Flagged so a future reader knows the retrieval route. |

---

## 3 · What the terms say, quoted

### 3.1 IGDB — the tier, and what it is conditioned on

> "The IGDB.com API is free for non-commercial usage under the terms of the
> Twitch Developer Service Agreement."
> — api-docs.igdb.com, *Account Creation*

> "**2. What is the price of the API?** The API is free for both non-commercial
> and commercial projects."
> — api-docs.igdb.com, *Business related FAQ*

> "**1. I want to use the API for a commercial project, is it allowed?** Yes, we
> offer commercial partnerships for users looking to integrate the API in
> monetized products. From our side, as part of the partnership, we ask for user
> facing attribution to IGDB.com from products integrating the IGDB API."
> — *Business related FAQ*

The published test is therefore **monetisation of the project**, and the
attribution obligation is attached to the **commercial partnership**, not to the
free tier.

### 3.2 IGDB — self-hosting is not merely permitted, it is preferred

> "**3. Am I allowed to store/cache the data locally?** Yes. In fact, we prefer
> if you store and serve the data to your end users. You remain in control over
> your user experience, while alleviating pressure on the API itself."
> — *Business related FAQ*

> "Note: Images that are removed or replaced from IGDB.com exist for 30 days
> before they are removed. Keep that in mind when designing cache logic."
> — api-docs.igdb.com, *Images*

This settles the self-hosting constraint from the other direction: the page's
zero-external-request guarantee and the provider's stated preference agree.

Image URL structure, from the same section:

> `https://images.igdb.com/igdb/image/upload/t_{size}/{hash}.jpg`

with `cover_big` documented as **264 × 374** and any size gaining a retina
variant by appending `_2x`.

### 3.3 IGDB — what "fair attribution" means, where it is owed

> "**4. Regarding user facing attribution (relating to the commercial
> partnership), any specific guidelines?** Not really. We expect fair
> attribution, i.e. attribution that is visible to your users and located in a
> static location (e.g. not in a change log)."
> — *Business related FAQ*

### 3.4 The TDSA — the grant, and its shape

A website is an App, and promotion is an enumerated licensed purpose:

> "**1. Apps:** software applications, games, **websites**, channels,
> Extensions, Drops, and other digital products: (a) that you submit to Twitch
> for license, sale, distribution, or promotion; and/or (b) which use any Program
> Materials…"
> — TDSA §II.A.1 (emphasis added)

> "**2. License from Twitch.** Subject to the terms and conditions in this
> Agreement, Twitch hereby grants you a limited, non-exclusive, worldwide,
> royalty-free, non-transferable, **non-sublicensable**, revocable license during
> the Term to: i. Use and reproduce the Program Materials solely to develop,
> test, **promote**, measure, support, and operate Your Services."
> — TDSA §II.B.2 (emphasis added)

> "**8. Your Services:** your Apps, content, ads, services, technology, data,
> digital goods, and all other materials included in or made available through
> your Apps."
> — TDSA §II.A.8

### 3.5 The TDSA — what Twitch expressly does *not* give

> "**A. Ownership of Program Materials.** The Program Materials are the
> intellectual property of Twitch **or its licensors**. The Program Materials are
> licensed, not sold, and Twitch retains and reserves all rights not expressly
> granted in this Agreement."
> — TDSA §III.A (emphasis added)

> "**D. Warranty Disclaimer.** THE PROGRAM MATERIALS ARE PROVIDED TO YOU "AS IS",
> "WHERE IS", WITH ALL FAULTS AND EACH PARTY DISCLAIMS ALL WARRANTIES, WHETHER
> EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE."
> — TDSA §X.D

> "**B.** … Before providing Twitch or any end user Your Services, you will have
> obtained the rights necessary for the exercise of all rights granted under this
> Agreement, and you will be solely responsible for and will pay any licensors or
> co-owners any royalties or other monies due to them related to Your Services;"
> — TDSA §XI.B, *Developer Representation and Warranties*

> "**E.** None of Your Services or the sale, distribution, or **promotion**
> thereof will violate any law; … or violate or infringe any intellectual
> property, proprietary, or other rights of any person or entity (including
> contractual rights, copyrights, trademarks, patents, trade dress, trade secret,
> common law rights, rights of publicity or privacy, or moral rights);"
> — TDSA §XI.E (emphasis added)

> "**E. Indemnification.** You will defend, hold harmless, and indemnify Twitch
> from any claims … to the extent arising out of: … (iii) the performance,
> **promotion**, sale, or distribution of Your Services…"
> — TDSA §X.E (emphasis added)

### 3.6 IGDB — how the artwork got into IGDB in the first place

> "**Uploading of images policy.** Before you upload an image, make sure that the
> image falls in one of the four categories: **Own work:** You own all rights to
> the image, usually meaning that you created it entirely yourself. **Freely
> licensed:** The image has been released under an acceptable free license.
> **Public domain:** The image is in the public domain, i.e. free of all
> copyrights. **Fair use:** Fair use is a limitation and exception to the
> exclusive right granted by copyright law to the author of a creative work…"
> — igdb.com/content-policy, *Image Policy*

This is the load-bearing sentence in the whole finding, and it is why the answer
below is in two parts rather than one.

---

## 4 · The finding

**On the narrow question the ticket asked — marketing website versus inside the
application — the terms answer favourably, and explicitly.** §II.A.1 defines an
App to include a **website**, and §II.B.2.i licenses use of Program Materials to
"develop, test, **promote**, measure, support, and operate Your Services".
Promotion is an enumerated licensed purpose, not an unstated one. There is no
clause anywhere in the TDSA or the IGDB documentation confining use to the
interior of a shipped application. **The marketing-page/in-app distinction that
this ticket was right to insist on does not, in fact, cut against us.**

**On the tier question, the answer is also favourable.** IGDB's published test is
whether the project is monetised. Zerado charges nothing, carries no affiliate,
runs no premium tier and generates no revenue. It sits on the free tier on the
provider's own stated test.

**The distinction that actually matters is a different one, and it is not
resolved by any of these documents.** Twitch licenses *its own* Program
Materials. That grant is expressly **non-sublicensable** (§II.B.2), the materials
are the property of "Twitch **or its licensors**" with all rights not expressly
granted reserved (§III.A), and they are supplied **AS IS with all warranties
disclaimed** (§X.D). Meanwhile §XI.B and §XI.E put the burden of having obtained
every necessary third-party right, and the warranty of non-infringement,
squarely on Zerado — and §X.E(iii) has Zerado indemnify Twitch for claims
arising from exactly the **promotion** of Your Services.

And IGDB's own Image Policy shows why that burden is not empty: IGDB accepts
cover uploads on a **fair use** basis. A database that accepts artwork under fair
use is not asserting ownership of it and is in no position to pass a copyright
licence down the chain. **The publisher's copyright in each cover is not cleared
by anybody in this chain, and the TDSA says in terms that this is the
developer's problem.**

So, stated plainly and without varnish:

> **IGDB and Twitch permit this use, including on a marketing page. Neither of
> them clears the publisher's copyright in the artwork, neither claims to, and
> the TDSA expressly assigns that residual to Zerado.**

What is *not* being claimed here: that non-commercial use makes this
automatically fine. It does not. Non-commercial, non-revenue, open-source status
materially improves the position — it satisfies the provider's own tier test and
it strengthens the first and fourth fair-use factors — but it is an improvement
to the position, not a permission, and this document does not pretend otherwise.

---

## 5 · The decision taken, and by whom

**Founder, 2026-08-25, verbatim:** *"we can use the images, and you're right
about the physical, so get everything from IGDB, and let's replace the
placeholders."*

The founder was shown the shape of §4 — that the provider permits it and that
the publisher-copyright residual remains and sits with Zerado — and decided to
proceed on the basis of Zerado's non-commercial, zero-revenue, donation-only,
open-source posture. **That is a risk-acceptance decision, and it is his to
make.** It is recorded here rather than left implicit precisely so that it reads
as a decision taken with the residual in view, not as an oversight.

The same exchange settled the physical-copy question: **one provider, one licence
question.** `PHYSICAL` is a provenance label on the row, not a different image
source — a cartridge or disc has an IGDB entry exactly as a Steam title does. No
image on this page comes from a web search or any source of unknown provenance.

---

## 6 · Attribution — what is owed, and what we do

**Owed:** nothing contractually. IGDB attaches user-facing attribution to the
**commercial partnership** (§3.1), and Zerado is not in one.

**Done anyway:** visible, user-facing attribution to IGDB.com in a **static
location** — the shape IGDB describes as "fair attribution" in §3.3 — discharged
in two places:

- the caption directly under the cover grid in §06 of the landing page, naming
  IGDB as the source of the cover art; and
- the site footer's credits line.

Doing the thing the provider says it values, while it is free to do so, is worth
more than the letter of an obligation we do not currently carry. If Zerado ever
becomes monetised, this attribution is already in place and the tier question is
the only thing that reopens.

---

## 7 · What would reverse this

Any one of these, and the art-directed CSS tiles come back — they are preserved
in git history and the component still renders them when a row carries no image:

- IGDB or Twitch answering the email in §`docs/legal/igdb-email-2026-08-25.md`
  in a way that narrows the free tier or the promotional use;
- Zerado acquiring any revenue stream, which fails IGDB's own published tier test
  and moves the project to the partnership track;
- a rights-holder objecting to a specific cover, in which case that row loses its
  image and the tile treatment renders in its place — a data-only change.

---

## 8 · Provenance of every image on the page

**Recorded per file in `docs/legal/cover-provenance.json`**, written by
`scripts/fetch-covers.mjs` on the run that produced the images: for each of the
twelve, the IGDB game id, the IGDB game slug, the name and release year IGDB
holds, the cover image id, and the exact source URL the bytes came from — the
image URL structure quoted in §3.2. Every file under `site/public/covers/`
appears there, and there is no other image source.

`site/src/data/coverGrid.ts` carries the row's *identity* — the name and release
year the fetch is required to match, and the platform the row is tagged with. It
is what pins which game a row means; it is not the record of what was fetched.

*(Corrected 2026-08-25, ticket #16. This section previously said the IGDB slug
and cover image hash were recorded in `coverGrid.ts`. They are not, and never
were — that file has no such fields. The claim was harmless while the grid was
twelve CSS gradients and there were no files to trace; it stopped being harmless
the moment real third-party artwork shipped, because this is the section a
rights-holder question gets answered from. The provenance file is generated by
the fetch rather than maintained by hand, so it cannot drift from what was
actually downloaded.)*
