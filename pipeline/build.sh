#!/usr/bin/env bash
set -euo pipefail

APP_NAME="slog-simple-blog"
OUTPUT_DIR="dist"
MAIN_PKG="./cmd/app"

echo "Building $APP_NAME (linux/amd64)"

# clean old build
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# environment for deterministic builds
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# optional metadata
GIT_COMMIT=$(git rev-parse --short HEAD)
GIT_TAG=$(git describe --tags --dirty --always 2>/dev/null || echo "dev")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# build flags
LDFLAGS=(
  "-s -w"
  "-X main.version=$GIT_TAG"
  "-X main.commit=$GIT_COMMIT"
  "-X main.date=$BUILD_DATE"
)

go build \
  -trimpath \
  -ldflags "${LDFLAGS[*]}" \
  -o "$OUTPUT_DIR/$APP_NAME" \
  "$MAIN_PKG"

echo "Build complete: $OUTPUT_DIR/$APP_NAME"

