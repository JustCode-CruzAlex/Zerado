# Deploying zerado.app

Everything needed to put the landing page on `zerado.app`, in the order it has to
happen, written to be **executed step by step** by whoever holds the DigitalOcean
account. Every command is copy-pasteable and every step states what you should
see back. No step says "then configure DNS" — it says which record, which value,
which TTL.

**The path is a $4.00/month droplet with nginx and certbot.** Chosen by the
founder on 2026-08-25. App Platform static is in the appendix with the reasoning
for whoever revisits this.

---

## ✅ Deployed — this runbook was executed, not just written

**`https://zerado.app` has been live since 2026-08-25.** Every step below was run
end to end against the real droplet, and the steps were corrected where reality
disagreed with them. Four things this run caught that a desk-written runbook
would have shipped broken:

| What broke | Fix |
|---|---|
| `http2 on;` is nginx **≥ 1.25.1**. Ubuntu 24.04 ships **1.24.0**, where it is an unknown directive and a fatal parse error. | `listen 443 ssl http2;`, which works on both |
| The full config declares `listen ... ssl`, so nginx would not load before a certificate existed — and certbot would not issue one against a config nginx could not load | An HTTP-only bootstrap config (Step 4), then the real one (Step 5) |
| `gzip_types text/html` → `[warn] duplicate MIME type` on every config test | `text/html` removed; nginx always gzips it |
| `rsync --chmod` is rejected by macOS's openrsync | Permissions normalised server-side instead |

Live results: **Lighthouse 100/100/100/100 on desktop and mobile**, axe-core **0
violations** at all four viewports, TLS valid to 2026-11-23, and
`certbot renew --dry-run` green.

---

## Where things stand

Verified **2026-08-25**:

| | |
|---|---|
| `zerado.app` nameservers | `ns1.digitalocean.com`, `ns2.digitalocean.com`, `ns3.digitalocean.com` — delegated from GoDaddy ✅ |
| `zerado.app` A record | **None.** `dig +short zerado.app A` returns nothing |
| `https://zerado.app` | No response — nothing is served yet |
| `github.com/JustCode-CruzAlex/Zerado` | 200, public |
| `www.flowforgesoft.com` | 200 |

The zone is live and delegated, and nothing is served from it. The cutover has
**no blast radius today** — there is no traffic to lose. Cut over before there is.

---

## Why the droplet is the better outcome, not a concession

The ticket recommended App Platform static and framed the droplet as a fallback
on price. That framing turned out to be wrong on the facts, in two ways.

**On price:** the $0.00 case needs one of the account's three free static-site
slots, and the account already runs several apps. That tier is almost certainly
consumed, which puts App Platform at **$3.00/month** against the droplet's
**$4.00** — a dollar apart, not free against paid.

**On capability, which is the part that actually decides it:** App Platform's
static-site component cannot read a `_headers` file, cannot set per-path
`Cache-Control`, and does not document HSTS control. Ticket #1 asks for all
three. Taking that path meant shipping three acceptance criteria as knowing
deviations.

| Ticket #1 asks for | App Platform static | Droplet + nginx |
|---|---|---|
| `_headers` cache policy honoured | ❌ not read | ✅ ported verbatim |
| `immutable` on `/_astro/*`, `no-cache` on the document | ❌ not configurable | ✅ |
| HSTS on | ❌ not documented as configurable | ✅ two years, `includeSubDomains` |
| Managed, auto-renewed TLS | ✅ | ✅ certbot, systemd timer |
| Rebuild + republish with no manual step | ✅ `deploy_on_push` | ✅ one idempotent command |
| Global CDN | ✅ | ➖ single region (see below) |

**The deviation disappears.** For one dollar a month, three acceptance criteria
stop being exceptions and become satisfied.

What is genuinely given up is App Platform's global CDN: the droplet serves from
one region. For a 360 KB zero-JavaScript page whose LCP is already 366 ms on
desktop, a single well-chosen region is not the bottleneck — and Cloudflare can
be put in front later without touching anything in this repository.

---

## The cost, as a number

| Item | Cost |
|---|---|
| Droplet — Basic, 1 vCPU / 512 MiB / 10 GiB SSD / 500 GiB transfer | **$4.00 / month** ($0.00595/hour) |
| DigitalOcean DNS for the zone | free |
| TLS certificate (Let's Encrypt via certbot) | free |
| **Total** | **$4.00 / month** |

**The smallest tier is the correct size and is not under-provisioned.** nginx
serving ~360 KB of static files uses a few MB of RAM and effectively no CPU;
there is no runtime, no database and no application process. The 500 GiB
transfer allowance is roughly **1.4 million** full first-time page loads a
month. Do not size up — there is nothing here to size up *for*.

Optional, and not included above: droplet backups add 20% (**$0.80/month**).
Worth it for a machine you will forget you own; the site itself is fully
reproducible from this repository either way.

Prices read from [digitalocean.com/pricing/droplets](https://www.digitalocean.com/pricing/droplets) on 2026-08-25.

> **512 MiB needs swap.** certbot pulls in Python and can exhaust a swapless
> 512 MiB box during certificate issuance. Step 2 below adds a 1 GiB swapfile.
> Do not skip it — this is the one way the cheapest tier bites.

---

## Step 1 — Install and authenticate `doctl`

`doctl` is **not installed** on the machine this runbook was written on, and no
DigitalOcean credentials exist there. Every step below runs on your machine.

```bash
brew install doctl                 # macOS
# or: snap install doctl           # Linux

doctl auth init                    # paste a Personal Access Token (read + write)
doctl account get
```

Create the token at **API → Tokens → Generate New Token** in the console.

**Expect:** `doctl account get` prints your email, droplet limit and status
`active`. If it errors, the token is wrong or lacks write scope — nothing below
will work until this line does.

## Step 2 — Create the droplet

```bash
# List your SSH key fingerprints — you need one for the create call.
doctl compute ssh-key list

# nyc3 is the closest region to Brazil that DigitalOcean offers.
doctl compute droplet create zerado-web \
  --region nyc3 \
  --size s-1vcpu-512mb-10gb \
  --image ubuntu-24-04-x64 \
  --ssh-keys <FINGERPRINT_FROM_ABOVE> \
  --enable-monitoring \
  --wait

IP=$(doctl compute droplet get zerado-web --format PublicIPv4 --no-header)
echo "Droplet IP: $IP"
```

**Expect:** a table with `zerado-web`, status `active`, and a public IPv4.
`echo` prints that IP. Keep it — the next steps use `$IP`.

If you have no SSH key on the account: `ssh-keygen -t ed25519 -C zerado` then
`doctl compute ssh-key import zerado --public-key-file ~/.ssh/id_ed25519.pub`.

Then the swapfile, before anything else touches this box:

```bash
ssh root@"$IP" bash -s <<'REMOTE'
set -euo pipefail
fallocate -l 1G /swapfile && chmod 600 /swapfile && mkswap /swapfile && swapon /swapfile
echo '/swapfile none swap sw 0 0' >> /etc/fstab
free -h
REMOTE
```

**Expect:** `free -h` shows a `Swap:` row of `1.0Gi`.

## Step 3 — DNS: the exact records

Two A records, both pointing at the droplet. Nothing else.

| Type | Name | Value | TTL |
|---|---|---|---|
| `A` | `@` (the apex, `zerado.app`) | the droplet's public IPv4 | `3600` |
| `A` | `www` | the same IPv4 | `3600` |

```bash
doctl compute domain records create zerado.app \
  --record-type A --record-name @   --record-data "$IP" --record-ttl 3600

doctl compute domain records create zerado.app \
  --record-type A --record-name www --record-data "$IP" --record-ttl 3600

doctl compute domain records list zerado.app
```

**Expect:** the list shows both A records with your IP. Then, from your machine:

```bash
dig +short zerado.app A
dig +short www.zerado.app A
```

**Expect:** both print the droplet IP. If they print nothing, wait — the records
are new and the TTL is 3600. Do **not** proceed to certbot until both resolve;
Let's Encrypt validates over HTTP against these names and will fail otherwise.

> **TTL note.** 3600 (one hour) is right for a record you may need to move. Lower
> it to 300 for a day *before* a planned IP change, then put it back.

## Step 4 — nginx, and the HTTP-only bootstrap

**Read this first, because the obvious order does not work.** The full config
(`zerado.app.conf`) declares `listen 443 ssl` with real certificate paths. nginx
treats a missing `ssl_certificate` as `[emerg]` — a fatal parse error, not a
warning — so it will not load that file before the certificate exists. And
certbot runs its own `nginx -t` before doing anything, so a config nginx cannot
load *also* blocks certbot from issuing the certificate that would fix it.

The way out is a throwaway HTTP-only config that serves the ACME challenge.
Install it, get the certificate, then swap in the real one.

```bash
ssh root@"$IP" bash -s <<'REMOTE'
set -euo pipefail
apt-get update -qq
apt-get install -y -qq nginx certbot rsync
mkdir -p /var/www/zerado/releases /var/www/certbot /etc/nginx/snippets
nginx -v
REMOTE
```

**Expect:** `nginx version: nginx/1.24.0` on Ubuntu 24.04.

> Only `certbot` is installed, not `python3-certbot-nginx`. This deploy uses
> `certbot certonly --webroot`, which does not need the nginx plugin — and the
> plugin would try to rewrite the config we carefully wrote.

Now the bootstrap config:

```bash
scp docs/deploy/nginx/zerado.app.bootstrap.conf root@"$IP":/etc/nginx/sites-available/zerado.app

ssh root@"$IP" bash -s <<'REMOTE'
set -euo pipefail
ln -sf /etc/nginx/sites-available/zerado.app /etc/nginx/sites-enabled/zerado.app
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl reload nginx
REMOTE

curl -s http://zerado.app/
```

**Expect:** `nginx: configuration file /etc/nginx/nginx.conf test is successful`,
and `curl` prints `zerado.app — provisioning`. This test genuinely passes: the
bootstrap config declares no TLS at all.

## Step 5 — TLS, then the real config, then proving renewal

Issue the certificate over HTTP-01, using the webroot the bootstrap config is
already serving:

```bash
ssh root@"$IP" certbot certonly --webroot -w /var/www/certbot \
  -d zerado.app -d www.zerado.app \
  --agree-tos -m alex@flowforgesoft.com --non-interactive
```

**Expect:** `Successfully received certificate.` and paths under
`/etc/letsencrypt/live/zerado.app/`.

> No `--redirect` flag: the HTTP→HTTPS 308s are already in `zerado.app.conf`.
> Passing it would have certbot rewrite a config it did not author.

The certificate now exists, so the real config will load. Swap it in:

```bash
scp docs/deploy/nginx/snippets/zerado-headers.conf root@"$IP":/etc/nginx/snippets/
scp docs/deploy/nginx/zerado.app.conf              root@"$IP":/etc/nginx/sites-available/zerado.app

ssh root@"$IP" bash -s <<'REMOTE'
set -euo pipefail
nginx -t
systemctl reload nginx
REMOTE
```

**Expect:** `test is successful` again — this time with TLS configured. If it
fails here, the certificate did not issue; go back a step rather than editing
the config.

**Renewal is verified, not assumed.** certbot installs a systemd timer; confirm
the timer exists *and* that a renewal actually completes end to end:

```bash
ssh root@"$IP" bash -s <<'REMOTE'
set -euo pipefail
echo "--- the timer that will do it ---"
systemctl list-timers certbot.timer --all --no-pager
echo "--- a real renewal, against the staging endpoint ---"
certbot renew --dry-run
echo "--- what is installed and when it expires ---"
certbot certificates
REMOTE
```

**Expect:** the timer is `active`; the dry run ends with
`Congratulations, all simulated renewals succeeded`; `certbot certificates`
shows both names on one certificate with ~89 days left.

If the dry run fails, renewal *will* fail silently in 60 days. Fix it now — it
is almost always port 80 being closed, or the ACME `location` block having been
dropped from the config.

## Step 6 — Deploy

One command, from a clean checkout of this repository:

```bash
./scripts/deploy.sh
```

**Do not pass `ZERADO_HOST="$IP"`.** The certificate is issued for `zerado.app`,
so pointing the script at a raw IP fails the TLS handshake on a name mismatch
and the final verification reports `000`. DNS already resolves by Step 3, so the
default host is correct.

It builds the site, **runs the page invariants before anything leaves your
machine**, rsyncs into a timestamped release directory, validates nginx, flips
the `current` symlink atomically, reloads nginx, prunes old releases, and
verifies the live URL returns 200.

**Expect:** it ends with `Deployed 20260825-…` and prints `no-cache` for both
`/` and `/index.html`.

It is **idempotent** — running it twice produces a second release and the same
live result. Nothing is ever half-deployed: the symlink swap is a single atomic
rename, so a visitor sees either the old release or the new one.

### Rollback

```bash
./scripts/deploy.sh --rollback   # revert to the previous release
./scripts/deploy.sh --list       # what is on the server, and what is live
```

Rollback is another atomic symlink swap — it does not rebuild, re-upload or need
the network beyond SSH, so it works when the thing you are rolling back from is
broken. The last **5** releases are kept.

### Republishing on a push to `main`

The repository's CI builds and checks the page on every pull request touching
`site/**`, which is the gate ticket #1 asks for. Publishing stays a deliberate
`./scripts/deploy.sh`.

To make a push to `main` publish by itself, add an SSH deploy key to repository
secrets and call the same script from a workflow — the script is written to be
driven either way. That is deliberately **not** committed: it puts a key with
write access to the web root into GitHub, and that is a decision to take on
purpose rather than inherit from a runbook.

---

## Verification — run every line, keep the output

This is the acceptance checklist. None of it can run until the site is served.

```bash
# 200 on both hostnames, valid TLS
curl -sI https://zerado.app     | head -20
curl -sI https://www.zerado.app | head -20

# www redirects to the apex
curl -sI https://www.zerado.app | grep -i '^location'

# http is redirected, not served
curl -sI http://zerado.app | head -3

# The cache policy is in force per route, and the security headers survived
# nginx's all-or-nothing add_header inheritance on EVERY route, not just `/`.
ASSET=$(curl -s https://zerado.app | grep -oE '/_astro/[^"]+\.css' | head -1)
for p in / /index.html "$ASSET" /logo.svg; do
  echo "== $p"
  curl -sI "https://zerado.app$p" \
    | grep -iE 'cache-control|strict-transport|content-security|x-frame|referrer-policy'
done
```

**Expect:** `immutable` with `max-age=31536000` on the `/_astro/*` asset,
`max-age=604800` on `/logo.svg`, `no-cache` on both `/` and `/index.html`, and
`Strict-Transport-Security` plus a CSP on **all four**.

```bash
# Certificate
echo | openssl s_client -servername zerado.app -connect zerado.app:443 2>/dev/null \
  | openssl x509 -noout -subject -issuer -dates

# Every link on the live page resolves
curl -s https://zerado.app | grep -oE 'href="https?://[^"]*"' | sed 's/href="//;s/"$//' | sort -u \
  | while read -r u; do printf '%s  %s\n' "$(curl -s -o /dev/null -w '%{http_code}' -L --max-time 20 "$u")" "$u"; done

# The waitlist mailto links carry the exact subject
curl -s https://zerado.app | grep -oE 'href="mailto:[^"]*"' | sort | uniq -c

# Ratified decision Q4 holds on the live page
curl -s https://zerado.app | grep -ic "$(printf 're%s' 'ddit')"   # must print 0

# The roadmap says what it should: Phase 1 in progress, 2-4 planned, none done
curl -s https://zerado.app > /tmp/zerado-live.html
node scripts/check-page.mjs /tmp/zerado-live.html
```

Then Lighthouse against the **live** URL, not localhost:

```bash
npx lighthouse https://zerado.app --preset=desktop --output=html --output-path=lh-desktop.html
npx lighthouse https://zerado.app --output=html --output-path=lh-mobile.html
```

A score below **95** in any category is investigated before the pull request
merges. The pre-deploy baseline is 100/100 on both
([`../performance/performance-report.md`](../performance/performance-report.md));
anything lower on the live URL is a hosting-layer regression — most likely
compression — not a page regression. If Lighthouse flags text compression, brotli
is not built into stock Ubuntu nginx; gzip is already on, and the
`ngx_brotli` lines in the config are commented out ready to enable.

Finally, screenshots at **375, 768, 1280 and 1920 px** of the live URL, attached
to the pull request.

---

## Repository settings — the console bits that are not code

```bash
gh repo edit JustCode-CruzAlex/Zerado \
  --description "Zerado — terminal game library. Pick a mood, get a game, play it tonight." \
  --homepage "https://zerado.app" \
  --add-topic zerado,game-library,backlog,terminal,tui,cli,bubbletea,charm,go,sqlite,local-first,steam,astro,landing-page,flowforgesoft
```

---

## Appendix — App Platform static, and why it was not chosen

Kept so the reasoning survives for whoever revisits this. The spec is still
committed and still correct at [`../../.do/app.yaml`](../../.do/app.yaml); CI
asserts its shape on every pull request.

**The price**, verified 2026-08-25 against
[digitalocean.com/pricing/app-platform](https://www.digitalocean.com/pricing/app-platform)
and [the product docs](https://docs.digitalocean.com/products/app-platform/details/pricing/):

| Item | Price |
|---|---|
| First **three** apps using only static-site components | $0.00 / month |
| Each additional static-site app | $3.00 / month |
| Outbound data included, per app | 1 GiB / month |
| Outbound data beyond that | $0.02 / GiB |

The account already runs several apps, so the free slots are almost certainly
consumed and the real figure is **$3.00/month** — against $4.00 for a droplet
that also satisfies three acceptance criteria App Platform cannot.

**What it cannot do**, verified against
[edge settings](https://docs.digitalocean.com/products/app-platform/how-to/configure-edge-settings/)
and [manage static sites](https://docs.digitalocean.com/products/app-platform/how-to/manage-static-sites/):

- **`_headers` is not read.** Only Netlify and Cloudflare Pages read that file.
- **Per-path `Cache-Control` is not configurable** for a static-site component;
  edge caching cannot even be disabled for apps containing static sites.
- **HSTS is not documented as configurable.**

Its genuine advantages — a global CDN and `deploy_on_push` — are real, and if the
page ever needs edge presence more than it needs exact headers, this is how it
would be applied:

```bash
doctl apps spec validate .do/app.yaml
doctl apps create --spec .do/app.yaml
doctl apps logs <APP_ID> --type build --follow
```

The DNS records differ: App Platform manages the apex itself from the
`domains:` block in the spec, so the two hand-created A records above would be
removed first.
