# Deploying zerado.app

Everything needed to put the landing page on `zerado.app`, in the order it has
to happen. Written to be run by someone with access to the DigitalOcean account
that holds the `zerado.app` zone.

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

The zone is live and delegated, and nothing is served from it. That means the
cutover has **no blast radius today** — there is no traffic to lose. Cut over
before there is.

---

## The price — confirmed, not estimated

From DigitalOcean's own pricing page and its product documentation, both read on
**2026-08-25**:

| Item | Price |
|---|---|
| First **three** apps that use only static-site components | **$0.00 / month** |
| Each additional static-site app beyond those three | **$3.00 / month** |
| Outbound data included, per app | **1 GiB / month** |
| Outbound data beyond that | **$0.02 / GiB** |
| Managed TLS certificate (Let's Encrypt) | included |
| DigitalOcean DNS for the zone | free |

**So `zerado.app` costs `$0.00`/month**, provided the account is not already
using its three free static-site slots. The whole page is ~360 KB served; 1 GiB
of transfer is roughly **2,900 full first-time page loads a month** before a
single cent of overage, and repeat visits come from CDN cache.

> **The one thing to confirm in the console, because it depends on the account
> and not on the price list:** how many static-site apps already exist. If three
> are already in use, this app is **$3.00 / month**, not $0.00. Check under
> **Apps** before creating it.

Sources: [digitalocean.com/pricing/app-platform](https://www.digitalocean.com/pricing/app-platform) ·
[docs.digitalocean.com — App Platform pricing](https://docs.digitalocean.com/products/app-platform/details/pricing/)

---

## The decision — App Platform static, and what it costs you

**Recommended: App Platform static site.** The page has no runtime, no database,
no environment variables and no server-side code. A droplet buys none of that
back and costs a machine to patch, an nginx to configure and a certificate to
renew. App Platform serves the same bytes with managed TLS, a global CDN, and a
git-push deploy, and keeps DNS in the account that already holds the zone. At
$0.00/month it is also free.

**But it cannot do two things this repository's `site/public/_headers` asks
for**, and that is worth knowing before rather than after:

| Asked for | App Platform static |
|---|---|
| `_headers` file honoured | ❌ Not read. Only Netlify and Cloudflare Pages read that file natively. |
| Per-path `Cache-Control` — `immutable` on `/_astro/*`, `no-cache` on `/` | ❌ Not configurable for a static-site component. App Platform applies its own CDN policy; edge caching cannot be turned off for apps with static sites. |
| `Strict-Transport-Security` (HSTS) | ❌ Not documented as configurable. |
| Managed TLS, auto-renewed | ✅ |
| Global CDN | ✅ |
| Rebuild + republish on push to `main` | ✅ `deploy_on_push` |

Verified against [App Platform edge settings](https://docs.digitalocean.com/products/app-platform/how-to/configure-edge-settings/)
and [manage static sites](https://docs.digitalocean.com/products/app-platform/how-to/manage-static-sites/), 2026-08-25.

**The trade, stated plainly.** For an unauthenticated, cookieless, zero-JavaScript
marketing page, the practical cost of both gaps is small: the CDN still caches,
the page still scores what it scores, and there is no session for a downgrade
attack to steal. The cost is real but low. **It is a founder call, not a
worker's**, because it is a knowing deviation from what the ticket asked for.

**If that trade is refused**, the fallback is committed and ready: a smallest
droplet with nginx and certbot, using
[`nginx/zerado.app.conf`](nginx/zerado.app.conf) — which ports `_headers`
verbatim, sends HSTS, and adds a content-security policy. It costs about
**$4–6/month** plus the ops burden, and the DNS steps below are the same in
shape. Nothing else in this repository changes.

---

## Path A — App Platform static (recommended)

### 1. Install and authenticate `doctl`

`doctl` is **not installed** on the machine this was prepared on.

```bash
brew install doctl                 # macOS
# or: snap install doctl           # Linux

doctl auth init                    # paste a Personal Access Token with read+write
doctl account get                  # confirms the token works
```

Create the token at **API → Tokens → Generate New Token** in the DigitalOcean
console. It needs **read and write** scope.

### 2. Check the free-tier slot before creating anything

```bash
doctl apps list --format ID,Spec.Name,DefaultIngress
```

Count the existing static-only apps. Fewer than three ⇒ this app is $0.00/month.

### 3. Connect GitHub, then create the app

App Platform pulls from GitHub through its own GitHub App, so **no DigitalOcean
token is ever stored in this repository**. Install it once:

**Console → Apps → Create App → GitHub → Manage Access**, and grant it
`JustCode-CruzAlex/Zerado`.

Then:

```bash
# Validate the spec BEFORE creating anything — this is the deploy's single
# entry point and the one artefact that cannot be smoke-tested without
# credentials. CI asserts its shape; this asserts it against the real API.
doctl apps spec validate .do/app.yaml

doctl apps create --spec .do/app.yaml
doctl apps list --format ID,Spec.Name          # note the APP_ID
doctl apps logs <APP_ID> --type build --follow # watch the first build
```

The spec builds with `npm ci && npm run build` from `/site` and publishes
`/dist`. Later changes are applied with:

```bash
doctl apps update <APP_ID> --spec .do/app.yaml
```

### 4. DNS

The spec already declares both hostnames, so App Platform creates the records in
the DigitalOcean-managed zone itself:

- `zerado.app` — **PRIMARY**
- `www.zerado.app` — **ALIAS**

App Platform points the apex at its edge with an `ALIAS`/`A` record it manages;
do not hand-create a conflicting A record. Confirm:

```bash
doctl compute domain records list zerado.app
dig +short zerado.app A
dig +short www.zerado.app
```

### 5. TLS

Managed Let's Encrypt certificates are issued automatically once the records
resolve — typically within a few minutes, occasionally up to an hour. Watch
**Settings → Domains** in the console until both hostnames read *Active*.

### 6. Deploys run themselves from here

`deploy_on_push: true` means every push to `main` that touches the repository
rebuilds and republishes with **no manual step**. Demonstrate it once and keep
the link:

```bash
doctl apps list-deployments <APP_ID> --format ID,Phase,Cause,Created
```

---

## Path B — droplet + nginx + certbot (fallback)

Only if Path A's header trade is refused.

```bash
# 1. Smallest droplet, Ubuntu 24.04.
doctl compute droplet create zerado-web \
  --region nyc3 --size s-1vcpu-512mb-10gb --image ubuntu-24-04-x64 \
  --ssh-keys <YOUR_KEY_FINGERPRINT> --wait

# 2. Point DNS at it.
IP=$(doctl compute droplet get zerado-web --format PublicIPv4 --no-header)
doctl compute domain records create zerado.app --record-type A --record-name @   --record-data "$IP" --record-ttl 3600
doctl compute domain records create zerado.app --record-type A --record-name www --record-data "$IP" --record-ttl 3600

# 3. On the droplet.
apt update && apt install -y nginx certbot python3-certbot-nginx
mkdir -p /var/www/zerado /var/www/certbot

# 4. From your machine — build and ship.
cd site && npm ci && npm run build
rsync -av --delete dist/ root@"$IP":/var/www/zerado/

# 5. Config and certificate. The snippet is NOT optional: nginx drops every
#    inherited add_header in any location that sets one of its own, so each
#    location re-includes the full header set from this file.
ssh root@"$IP" 'mkdir -p /etc/nginx/snippets'
scp docs/deploy/nginx/snippets/zerado-headers.conf root@"$IP":/etc/nginx/snippets/
scp docs/deploy/nginx/zerado.app.conf root@"$IP":/etc/nginx/sites-available/zerado.app
ssh root@"$IP" 'ln -sf /etc/nginx/sites-available/zerado.app /etc/nginx/sites-enabled/ \
  && rm -f /etc/nginx/sites-enabled/default && nginx -t && systemctl reload nginx'
ssh root@"$IP" 'certbot --nginx -d zerado.app -d www.zerado.app --agree-tos -m alex@flowforgesoft.com --redirect'
ssh root@"$IP" 'certbot renew --dry-run'
```

Automating step 4 needs a GitHub Actions job with an SSH deploy key in repository
secrets — which Path A avoids entirely. That job is not committed, because the
recommendation is Path A.

---

## Verification — run every line, keep the output

This is the acceptance checklist. None of it can be run until the site is
actually served.

```bash
# Path B only — confirm the headers survived on EVERY route, not just `/`.
# This is the check that catches nginx's all-or-nothing add_header inheritance.
for p in / /index.html /logo.svg; do
  echo "== $p"
  curl -sI "https://zerado.app$p" \
    | grep -iE 'cache-control|strict-transport|content-security|x-frame|referrer-policy'
done

# 200 on both hostnames, valid TLS
curl -sI https://zerado.app     | head -20
curl -sI https://www.zerado.app | head -20

# What the cache policy ACTUALLY is on the hashed bundle — observe, do not assume.
ASSET=$(curl -s https://zerado.app | grep -oE '/_astro/[^"]+\.css' | head -1)
curl -sI "https://zerado.app${ASSET}" | grep -i 'cache-control\|age\|cf-cache\|x-cache'
curl -sI https://zerado.app | grep -i 'cache-control\|strict-transport-security'

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
```

Then, against the **live** URL, not localhost:

```bash
npx lighthouse https://zerado.app --preset=desktop --output=html --output-path=lh-desktop.html
npx lighthouse https://zerado.app --output=html --output-path=lh-mobile.html
```

A score below **95** in any category is investigated before the pull request
merges. The baseline from the pre-deploy run is 100/100 on both
([`../performance/performance-report.md`](../performance/performance-report.md));
anything lower on the live URL is a hosting-layer regression — most likely
compression or cache policy — not a page regression.

Finally, screenshots at **375, 768, 1280 and 1920 px** of the live URL, attached
to the pull request.

---

## Repository settings — the console bits that are not code

Set once, under **Settings** and the **About** panel of the repository:

- **Description:** `Zerado — terminal game library. Pick a mood, get a game, play it tonight.`
- **Website:** `https://zerado.app`
- **Topics:** `zerado`, `game-library`, `backlog`, `terminal`, `tui`, `cli`,
  `bubbletea`, `charm`, `go`, `sqlite`, `local-first`, `steam`, `astro`,
  `landing-page`, `flowforgesoft`

Or from the command line:

```bash
gh repo edit JustCode-CruzAlex/Zerado \
  --description "Zerado — terminal game library. Pick a mood, get a game, play it tonight." \
  --homepage "https://zerado.app" \
  --add-topic zerado,game-library,backlog,terminal,tui,cli,bubbletea,charm,go,sqlite,local-first,steam,astro,landing-page,flowforgesoft
```
