#!/usr/bin/env bash
set -e

TAG="$1"

echo "Building release $TAG"

./pipeline/build.sh
#./pipeline/test.sh
./pipeline/deploy.sh "$TAG"

