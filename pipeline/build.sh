#!/usr/bin/env bash
set -euo pipefail

APP_NAME="slog-simple-blog"
OUTPUT_DIR="dist"
MAIN_PKG="./cmd/app.go"

echo "Building $APP_NAME (linux/amd64)"

# clean old build
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# environment for deterministic builds
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

go build \
  -trimpath \
  -ldflags "${LDFLAGS[*]}" \
  -o "$OUTPUT_DIR/$APP_NAME" \
  "$MAIN_PKG"

echo "Build complete: $OUTPUT_DIR/$APP_NAME"

