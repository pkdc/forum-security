#!/usr/bin/env bash
# VPS initial setup script for the forum application.
# Run once as root on a fresh Ubuntu/Debian VPS.
#
# Usage:
#   curl -fsSL https://yourrepo/deploy/setup-vps.sh | sudo bash
#   OR copy to server and run: sudo bash setup-vps.sh
#
# Before running, set:
#   DOMAIN=yourdomain.com  (must already point to this server's IP)

set -euo pipefail

DOMAIN="${DOMAIN:-yourdomain.com}"
APP_USER="forum"
APP_DIR="/opt/forum"
SERVICE_FILE="/etc/systemd/system/forum.service"
GO_VERSION="1.24.0"

echo "==> Setting up forum on domain: $DOMAIN"

# ── 1. System packages ──────────────────────────────────────────────────────
apt-get update -q
apt-get install -y --no-install-recommends \
    git curl ca-certificates ufw build-essential gcc sqlite3

# ── 2. Install Go (if not present or wrong version) ─────────────────────────
if ! go version 2>/dev/null | grep -q "$GO_VERSION"; then
    echo "==> Installing Go $GO_VERSION"
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
        | tar -C /usr/local -xz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile.d/go.sh
fi
export PATH=$PATH:/usr/local/go/bin

# ── 3. Create dedicated user and app directory ───────────────────────────────
if ! id "$APP_USER" &>/dev/null; then
    echo "==> Creating system user: $APP_USER"
    useradd --system --no-create-home --shell /usr/sbin/nologin "$APP_USER"
fi

mkdir -p "$APP_DIR/forum/certs"
chown -R "$APP_USER:$APP_USER" "$APP_DIR"

# ── 4. Firewall: open only SSH, HTTP, HTTPS ──────────────────────────────────
echo "==> Configuring UFW firewall"
ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp   comment 'SSH'
ufw allow 80/tcp   comment 'HTTP  (ACME challenges)'
ufw allow 443/tcp  comment 'HTTPS'
ufw --force enable

# ── 5. Install systemd service ───────────────────────────────────────────────
echo "==> Installing systemd service"
cp "$(dirname "$0")/forum.service" "$SERVICE_FILE"
# Inject the real domain into the service file
sed -i "s/yourdomain.com/$DOMAIN/" "$SERVICE_FILE"

systemctl daemon-reload
systemctl enable forum

echo ""
echo "==> Setup complete."
echo "    Next: run deploy.sh to build and start the application."
