# Security Policy

## Reporting a vulnerability

**Do not open a public issue for a security problem.** Report it privately:

- **Preferred:** [open a private security advisory](../../security/advisories/new)
  on this repository (GitHub → Security → Advisories → Report a vulnerability).
- **Or by email:** **alex@flowforgesoft.com**, subject line starting with
  `Zerado-Security`.

Please include what you found, how to reproduce it, what an attacker could do
with it, and any proof-of-concept you have. If you would like credit in the fix,
say so and tell us the name to use.

### What to expect

| | |
|---|---|
| Acknowledgement of your report | within **3 business days** |
| An initial assessment — accepted, needs more information, or not a vulnerability | within **10 business days** |
| Fix and disclosure | coordinated with you; we will not publish before you have had a chance to review |

Zerado is a small project run by one company. These are honest targets, not a
funded SLA. If a deadline slips you will hear why rather than hear nothing.

There is **no bug bounty**. We will not pay for reports, and we will not pretend
otherwise to attract them.

## Supported versions

Zerado is **pre-alpha**. There is no released binary, no version tag, and no
supported release line yet. Today the attack surface is:

| Component | Status |
|---|---|
| The landing page at `zerado.app` (`site/`) | Live — in scope |
| Repository contents and CI workflows | In scope |
| The `zerado` CLI/TUI | Does not exist yet |
| Any Zerado-operated API or server | Does not exist |

When Phase 1 ships, this table gains a supported-versions row. Until then,
`main` is the only thing there is, and it is the thing that gets fixed.

## Out of scope

- Findings against **third-party services** Zerado will talk to (Steam,
  price-data providers) — report those to the service.
- Missing hardening headers that have no exploitable consequence on a static,
  zero-JavaScript, same-origin page.
- Automated scanner output with no demonstrated impact.
- Social engineering, physical attacks, or denial of service against
  infrastructure we do not run.
- **The waitlist address being publicly visible.** `alex@flowforgesoft.com`
  appears on the landing page on purpose: the waitlist is a `mailto:` link, not a
  form, so there is no backend to breach and no address list to leak. It will be
  scraped. That trade was made deliberately and is not a vulnerability.

## API keys are yours, not ours

This is the most important security property of Zerado's design, so it is stated
plainly rather than left to be inferred:

- **Zerado has no server**, and there is no Zerado account. Your library lives in
  a single SQLite file on your own machine.
- **Every credential is user-held.** The Steam Web API key you generate is yours;
  it is stored locally, used locally, and **never transmitted to any
  FlowForgeSoft-operated endpoint** — none exists to transmit it to.
- **We cannot leak a key we never receive.** If your key is compromised, it was
  compromised on your machine or at Steam, and you revoke it at
  [steamcommunity.com/dev/apikey](https://steamcommunity.com/dev/apikey).
- Treat your key like a password: it is not a secret Zerado can rotate for you.

When the Phase 4 community layer exists, it will be the first Zerado-operated
service, and this document gains the section it needs. That has not happened.

## Supply chain

The landing page has **zero runtime dependencies** — the shipped page is static
HTML, CSS, self-hosted WOFF2 fonts and inline SVG, with no client-side
JavaScript and no external network requests. Astro is a build-time dependency
only, pinned by `site/package-lock.json`. Dependency changes are reviewed like
any other change.
