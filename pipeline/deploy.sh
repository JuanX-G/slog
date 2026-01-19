#!/usr/bin/env bash
set -euo pipefail

PROJECT="slog-simple-blog"
TAG="${1:?release tag required}"

PI_HOST="$RELEASE_REMOTE_USER@$RELEASE_REMOTE_ADDR"
BASE_DIR="/var/www/releases/$PROJECT"
RELEASE_DIR="$BASE_DIR/$TAG"
CURRENT_LINK="$BASE_DIR/latest"

echo "Deploying $PROJECT release $TAG"

# sanity checks
if [[ ! -d dist ]]; then
  echo "dist/ not found — did build fail?"
  exit 1
fi

if [[ -z "$(ls -A dist)" ]]; then
  echo "dist/ is empty — refusing to deploy"
  exit 1
fi

# create release dir on Pi
ssh "$PI_HOST" "mkdir -p '$RELEASE_DIR'"

# upload artifacts
rsync -avz --delete \
  dist/ \
  "$PI_HOST:$RELEASE_DIR/"

# atomically update 'latest' symlink
ssh "$PI_HOST" <<EOF
ln -sfn "$RELEASE_DIR" "$CURRENT_LINK"
EOF

echo "Release $TAG deployed successfully"
