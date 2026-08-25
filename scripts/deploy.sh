#!/usr/bin/env bash
#
# Deploy the Zerado landing page to the droplet. One command, idempotent,
# with a real rollback.
#
#   ./scripts/deploy.sh                     # build + deploy
#   ./scripts/deploy.sh --rollback          # revert to the previous release
#   ./scripts/deploy.sh --list              # show releases on the server
#   ZERADO_HOST=1.2.3.4 ./scripts/deploy.sh # override the target
#
# Layout on the server — nginx serves the `current` symlink, so the switch is a
# single atomic rename. Nothing is ever half-deployed, and the previous release
# stays on disk until it is pruned.
#
#   /var/www/zerado/releases/20260825-143000/   ← rsynced here first
#   /var/www/zerado/releases/20260825-141500/
#   /var/www/zerado/current -> releases/20260825-143000
#
set -euo pipefail

HOST="${ZERADO_HOST:-zerado.app}"
USER="${ZERADO_USER:-root}"
ROOT="/var/www/zerado"
KEEP="${ZERADO_KEEP:-5}"
REMOTE="${USER}@${HOST}"

say() { printf '\033[1m==>\033[0m %s\n' "$*"; }
die() { printf '\033[31mERROR\033[0m %s\n' "$*" >&2; exit 1; }

case "${1:-deploy}" in
  --list)
    ssh "$REMOTE" "ls -1 ${ROOT}/releases 2>/dev/null | sort; echo '--- current ---'; readlink ${ROOT}/current"
    exit 0
    ;;
  --rollback)
    say "Rolling back on ${HOST}"
    ssh "$REMOTE" bash -euo pipefail <<ROLLBACK
      cd "${ROOT}/releases"
      current=\$(basename "\$(readlink "${ROOT}/current")")
      previous=\$(ls -1 | sort | grep -v "^\${current}\$" | tail -1)
      [ -n "\$previous" ] || { echo "No previous release to roll back to."; exit 1; }
      ln -sfn "${ROOT}/releases/\$previous" "${ROOT}/current.tmp"
      mv -Tf "${ROOT}/current.tmp" "${ROOT}/current"
      nginx -t && systemctl reload nginx
      echo "Rolled back: \$current -> \$previous"
ROLLBACK
    exit 0
    ;;
  deploy) ;;
  *) die "Unknown argument: $1" ;;
esac

command -v rsync >/dev/null || die "rsync is not installed locally."

# --- build ---------------------------------------------------------------
say "Building the site"
( cd site && npm ci --no-audit --no-fund && npm run build )

say "Checking the page invariants before anything leaves this machine"
node scripts/check-page.mjs site/dist/index.html

[ -f site/dist/index.html ] || die "site/dist/index.html missing after build."
[ -f site/dist/_headers ]   || die "site/dist/_headers missing after build."

# --- ship ----------------------------------------------------------------
RELEASE="$(date -u +%Y%m%d-%H%M%S)"
say "Shipping release ${RELEASE} to ${HOST}"

ssh "$REMOTE" "mkdir -p ${ROOT}/releases/${RELEASE}"
rsync -az --delete --chmod=D755,F644 site/dist/ "${REMOTE}:${ROOT}/releases/${RELEASE}/"

# --- switch, atomically --------------------------------------------------
say "Switching current -> ${RELEASE}"
ssh "$REMOTE" bash -euo pipefail <<SWITCH
  test -f "${ROOT}/releases/${RELEASE}/index.html" || { echo "Upload incomplete — refusing to switch."; exit 1; }
  ln -sfn "${ROOT}/releases/${RELEASE}" "${ROOT}/current.tmp"
  mv -Tf "${ROOT}/current.tmp" "${ROOT}/current"
  nginx -t
  systemctl reload nginx
  cd "${ROOT}/releases"
  ls -1 | sort -r | tail -n +$((KEEP + 1)) | xargs -r rm -rf
  echo "Live: \$(readlink ${ROOT}/current)"
SWITCH

# --- verify --------------------------------------------------------------
say "Verifying the live page"
code=$(curl -s -o /dev/null -w '%{http_code}' "https://${HOST}/")
[ "$code" = "200" ] || die "https://${HOST}/ returned ${code}. Roll back with: $0 --rollback"

for path in / /index.html; do
  printf '  %-14s ' "$path"
  curl -sI "https://${HOST}${path}" | grep -i 'cache-control' | tr -d '\r' || echo '(no cache-control!)'
done

say "Deployed ${RELEASE}. Roll back with: $0 --rollback"
