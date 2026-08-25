---
title: Zerado — SEO & Metadata
discipline: CONTENT
doc-no: ZRD-CONTENT-02
rev: A
date: 2026-08-24
status: final — implement verbatim
---

# Zerado — SEO & metadata

## Title tag

```
Zerado — terminal game library tracker, sorted by mood
```
54 characters. Per `naming.md`'s consequence for this exact name: never a bare "Zerado" in a
title — always paired with a qualifier, because "zerado" alone is a common Portuguese word
(see "Ranking note," below) and carries no meaning to a search engine on its own.

## Meta description

```
Zerado sorts your Steam library by mood and tracks price against the all-time low.
Terminal-first, offline, one SQLite file you own. Join the waitlist.
```
151 characters. Contains the primary keyword, states the two concrete mechanisms (mood sort,
price-vs-low), and ends on the one real CTA the page has — there is nothing else to offer a
pre-launch product.

## Canonical URL

```
https://zerado.app/
```

## Open Graph

| Field | Value |
|---|---|
| `og:type` | `website` |
| `og:site_name` | Zerado |
| `og:title` | Zerado — terminal game library tracker, sorted by mood |
| `og:description` | Zerado sorts your Steam library by mood and tracks price against the all-time low. Terminal-first, offline, one SQLite file you own. |
| `og:url` | `https://zerado.app/` |
| `og:image` | *(see "OG image concept," below — no file exists yet)* |
| `og:locale` | `en_US` |

**`og:description` is intentionally shorter than the meta description** — it drops the CTA
sentence, since a social-share card and a search result earn a click differently; the card's
job is to describe, the search result's job is to invite the click.

## Twitter card

| Field | Value |
|---|---|
| `twitter:card` | `summary_large_image` |
| `twitter:title` | Zerado — terminal game library tracker |
| `twitter:description` | Sorts your Steam library by mood. Tracks price against the all-time low. Terminal-first, offline, one file you own. |
| `twitter:image` | *(same asset as `og:image` — see below)* |

## OG image concept — one line, for the build phase to execute

A dark, amber-on-chassis card (§4/§6 of the brand manual) carrying the wordmark `ZERADO`, the
tagline "Pick a mood. Get a game. Play it tonight.", and a single glyph-and-label state chip
(`◉ ZERADO`, cyan) as the one accent of colour — no screenshot, since none exists yet, and no
fabricated UI.

## JSON-LD — `SoftwareApplication`

Every field below is either verifiably true today or a plain description of the product. **No
`aggregateRating`, no `offers`/price, no `applicationSubCategory` download counts, no
`softwareVersion`** — none of those exist yet, and a schema block is exactly the kind of thing
a crawler treats as a hard factual claim, not marketing copy.

```json
{
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  "name": "Zerado",
  "url": "https://zerado.app/",
  "description": "Zerado is a terminal-first game-library tracker. It reads your game collection, sorts it by mood instead of genre, tracks price against the all-time low, and answers what to play tonight.",
  "applicationCategory": "UtilitiesApplication",
  "author": {
    "@type": "Organization",
    "name": "FlowForgeSoft"
  },
  "sameAs": [
    "https://github.com/JustCode-CruzAlex/Zerado"
  ]
}
```

Deliberately omitted, and why:
- **`operatingSystem`** — not yet confirmed publicly for this build; stating one would be a
  guess dressed as a fact.
- **`offers` / `price`** — nothing is for sale.
- **`aggregateRating` / `review`** — no reviews exist. Fabricating one is the single fastest
  way to lose a rich-result eligibility check, and it's also just untrue.
- **`softwareVersion` / `datePublished`** — pre-launch; there is no version to name.

## Keyword set

**Primary:**
- terminal game library tracker
- game backlog tracker terminal

**Secondary:**
- Steam library tracker TUI
- game mood tracker
- game tracker offline
- game price tracker all-time low

**Long-tail, compound (the ones with the least competition):**
- zerado game tracker
- zerado terminal
- zerado backlog

## Ranking note — read before setting expectations

**"Zerado" alone is not a winnable search term**, and this was accepted knowingly at naming
(see `naming.md`, "The risks we accepted," §3). It is common, everyday Brazilian Portuguese —
"zerado" also means a bank balance at zero, a reset counter, a football score. Ranking for the
bare word in Portuguese-speaking markets is not realistic, and paid ranking for it would mostly
buy traffic that has nothing to do with the product.

What *is* winnable: the domain `zerado.app` is exact-match and unambiguous, and the compound
long-tail queries above have effectively no existing competition. Per `naming.md`'s
consequence for the site, every title and every piece of metadata pairs the bare name with a
qualifier — "terminal game library" — for this exact reason. The launch channel is direct and
community-led (a public GitHub repo, a small terminal-tool audience), not organic search — SEO
here is a long-tail hedge, not the primary acquisition plan.
