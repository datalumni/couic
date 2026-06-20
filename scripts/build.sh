#!/usr/bin/env bash
set -euo pipefail
TARGET="${1:-linux}"
case "$TARGET" in
  linux|windows) ;;
  darwin)
    echo "macOS builds must run natively (see AGENTS.md for commands)." >&2
    exit 1
    ;;
  *)
    echo "Usage: $0 {linux|windows}" >&2
    exit 1
    ;;
esac

ext=""
[ "$TARGET" = "windows" ] && ext=".exe"

docker build -f scripts/Dockerfile --build-arg TARGET="$TARGET" -t couic-builder .
docker create --name tmp couic-builder
docker cp "tmp:/couic${ext}" "./couic${ext}"
docker rm tmp
echo "Built: couic${ext}"
