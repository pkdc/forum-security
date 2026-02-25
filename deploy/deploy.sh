#!/usr/bin/env bash
# Deploy / update the forum application on the VPS.
# Run as root (or via sudo) after setup-vps.sh has been run.
#
# Usage:
#   sudo bash deploy.sh [repo-url]
#
# Examples:
#   sudo bash deploy.sh https://github.com/you/forum-security.git
#   sudo bash deploy.sh          # re-deploys from existing /opt/forum/src

set -euo pipefail

REPO_URL="${1:-}"
APP_USER="forum"
APP_DIR="/opt/forum"
SRC_DIR="$APP_DIR/src"
BINARY="$APP_DIR/forum"

export PATH=$PATH:/usr/local/go/bin

echo "==> Deploying forum"

# ── 1. Fetch source ──────────────────────────────────────────────────────────
if [ -n "$REPO_URL" ]; then
    echo "==> Cloning / updating from $REPO_URL"
    if [ -d "$SRC_DIR/.git" ]; then
        git -C "$SRC_DIR" pull
    else
        git clone "$REPO_URL" "$SRC_DIR"
    fi
elif [ ! -d "$SRC_DIR" ]; then
    echo "ERROR: No repo URL provided and $SRC_DIR does not exist."
    echo "Usage: sudo bash deploy.sh <repo-url>"
    exit 1
fi

# ── 2. Build ─────────────────────────────────────────────────────────────────
echo "==> Building binary"
cd "$SRC_DIR"
go mod download
CGO_ENABLED=1 go build -o "$BINARY" ./cmd/web
chown "$APP_USER:$APP_USER" "$BINARY"

# ── 3. Copy static assets and templates ──────────────────────────────────────
echo "==> Syncing assets and templates"
rsync -a --delete "$SRC_DIR/assets/"    "$APP_DIR/assets/"
rsync -a --delete "$SRC_DIR/templates/" "$APP_DIR/templates/"
chown -R "$APP_USER:$APP_USER" "$APP_DIR/assets" "$APP_DIR/templates"

# ── 4. Restart service ───────────────────────────────────────────────────────
echo "==> Restarting forum service"
systemctl restart forum
systemctl status forum --no-pager

echo ""
echo "==> Deployment complete."
echo "    Logs: journalctl -u forum -f"
