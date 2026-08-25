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
      # Strictly OLDER than current, so repeated rollbacks WALK BACK through
      # history. Selecting "newest that isn't current" would toggle: after
      # R3 -> R2, the newest non-current is R3 again, and a second rollback
      # would roll forward into the release you were escaping.
      # Release names are UTC timestamps, so lexicographic order is chronological.
      previous=\$(ls -1 | sort | awk -v c="\$current" '\$0 < c' | tail -1)
      [ -n "\$previous" ] || { echo "No release older than \$current to roll back to."; exit 1; }
      # Validate BEFORE switching: never leave a box pointing at a new release
      # with a config that will not load.
      nginx -t
      ln -sfn "${ROOT}/releases/\$previous" "${ROOT}/current.tmp"
      mv -Tf "${ROOT}/current.tmp" "${ROOT}/current"
      systemctl reload nginx
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
# No --chmod: macOS ships openrsync, which rejects it outright
# ("invalid argument"). Permissions are normalised on the server instead,
# which works with whichever rsync the operator happens to have.
rsync -az --delete site/dist/ "${REMOTE}:${ROOT}/releases/${RELEASE}/"

# --- switch, atomically --------------------------------------------------
say "Switching current -> ${RELEASE}"
ssh "$REMOTE" bash -euo pipefail <<SWITCH
  test -f "${ROOT}/releases/${RELEASE}/index.html" || { echo "Upload incomplete — refusing to switch."; exit 1; }
  find "${ROOT}/releases/${RELEASE}" -type d -exec chmod 755 {} +
  find "${ROOT}/releases/${RELEASE}" -type f -exec chmod 644 {} +
  # Validate BEFORE the switch. Under set -e a failing "nginx -t" after the
  # symlink move would abort with the new release already live and the reload
  # never issued -- a state the atomicity below is supposed to make impossible.
  #
  # NO BACKTICKS ANYWHERE IN THIS HEREDOC. It is unquoted (<<SWITCH, not
  # <<'SWITCH') because it needs ${ROOT} and ${RELEASE} expanded locally -- and
  # that also makes backticks command substitution, executed on the OPERATOR'S
  # machine. A comment that merely QUOTED an nginx command really did run it
  # locally on the first deploy -- including while documenting this very trap.
  nginx -t
  ln -sfn "${ROOT}/releases/${RELEASE}" "${ROOT}/current.tmp"
  mv -Tf "${ROOT}/current.tmp" "${ROOT}/current"
  systemctl reload nginx
  cd "${ROOT}/releases"
  ls -1 | sort -r | tail -n +$((KEEP + 1)) | xargs -r rm -rf
  echo "Live: \$(readlink ${ROOT}/current)"
SWITCH

# --- verify --------------------------------------------------------------
# `|| echo 000` matters: under `set -e` a bare `code=$(curl …)` that exits
# non-zero kills the script at the assignment, so the diagnosis below — and the
# rollback hint — would never print. curl exits 60 on a certificate-name
# mismatch, which is exactly what happens if HOST is a raw IP rather than
# zerado.app, so the failure mode this guards is a likely one, not a theoretical
# one. The release is already live at this point; only the check is in doubt.
say "Verifying the live page"
code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 20 "https://${HOST}/" || echo 000)
if [ "$code" != "200" ]; then
  echo "  https://${HOST}/ returned ${code}"
  [ "$code" = "000" ] && echo "  (000 = no HTTP response: TLS failure, DNS, or a firewall. If HOST is an IP," \
                              "the certificate is for zerado.app and will not match — use the hostname.)"
  die "Verification failed. The release IS live; roll back with: $0 --rollback"
fi

for path in / /index.html; do
  printf '  %-14s ' "$path"
  curl -sI --max-time 20 "https://${HOST}${path}" | grep -i 'cache-control' | tr -d '\r' || echo '(no cache-control!)'
done

say "Deployed ${RELEASE}. Roll back with: $0 --rollback"
