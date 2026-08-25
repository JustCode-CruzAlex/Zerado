---
title: Zerado — email to IGDB, free tier and image use
discipline: LEGAL / SOURCING
doc-no: ZRD-LEGAL-02
rev: A
date: 2026-08-25
status: drafted — awaiting send from a founder-controlled address
ticket: "#16"
---

# Email to IGDB — one message, both questions

Two questions were owed to IGDB and they are folded into a single message, as
ticket #16 required. Neither is blocking: the reading in
[`igdb-image-licence-finding.md`](./igdb-image-licence-finding.md) stands on the
published terms, and the founder took the decision on 2026-08-25. This is to get
the provider's own answer on record.

**To:** `partner@igdb.com`
**Cc:** the address the Twitch developer application is registered under
**Subject:** Free tier for a non-commercial open-source project, and cover art on our project page

> Hello,
>
> I maintain **Zerado**, an open-source game-library manager. I would like to
> confirm two things with you before we go further, and to have your answer on
> record rather than rely on my own reading of the published terms.
>
> **About the project.** Zerado is free software under an open-source licence.
> It charges nothing, has no premium tier, carries no advertising and no
> affiliate links, and generates no revenue of any kind. If it is ever funded it
> will be by voluntary donations that cover server costs; there is no donation
> control on the site today. The source is public at
> `https://github.com/JustCode-CruzAlex/Zerado` and the project page is
> `https://zerado.app`.
>
> **Question 1 — the free tier.** Your documentation states that the API is free
> for non-commercial usage, and the FAQ describes the commercial partnership
> track as being for monetised products. On that test I read Zerado as
> non-commercial: it generates no revenue. Would you confirm that a
> donation-funded, zero-revenue, open-source project of this kind falls within
> the free tier as you apply it in practice?
>
> **Question 2 — cover art on the project page.** This is the question I most
> want your answer on, because I do not think the documentation addresses it
> directly. We display game cover art fetched from IGDB **inside** the
> application, which is the obvious case. We would also like to show a small
> number of covers on the project's **public landing page** — the page that
> describes what the software does — as an illustration of the library view.
>
> I read the Twitch Developer Services Agreement §II.B.2 as covering this, since
> it licenses use of Program Materials to "develop, test, **promote**, measure,
> support, and operate Your Services", and §II.A.1 includes websites in the
> definition of Apps. But promoting a product is a different activity from
> running it, and I would rather ask than assume. **Do you treat cover art shown
> on a project's own marketing or landing page the same as cover art shown inside
> the application?** If you would prefer we did not, we will take them off that
> page — we have an art-directed fallback already built and it costs us nothing
> to go back to it.
>
> **Attribution.** We are not asking to skip it. Your FAQ attaches user-facing
> attribution to the commercial partnership, and we are not in one, but we credit
> IGDB.com visibly and in a static location anyway — under the cover grid itself
> and in the site footer. If you have a preferred wording, mark or placement, tell
> me and I will match it exactly.
>
> **Image handling.** We do not hotlink your CDN. We fetch through the API and
> serve the images from our own host, which I believe is what you prefer per FAQ
> item 3, and we have noted the 30-day window on removed or replaced images when
> designing our cache.
>
> Thank you for the database — it is the reason a project like this is feasible
> at all.
>
> Best regards,
> Alexandre Reis Correa Cruz
> Zerado / FlowForgeSoft

---

## Why this is not a blocker

Ticket #16 originally made the answer to Question 2 a gate. The founder lifted
that gate on 2026-08-25 after seeing the finding, and the reasoning is in
[`igdb-image-licence-finding.md`](./igdb-image-licence-finding.md) §5. The email
still goes, because a provider's own answer is worth more than a careful reading
of their published terms, and because the reversal path in §7 of that document is
cheap: a row with no image renders the art-directed tile it replaced.

## Recording the reply

When IGDB replies, append the reply verbatim below this line with the date
received, and update §7 of the licence finding if the answer narrows anything.

**Reply received:** *(none yet as of 2026-08-25)*
