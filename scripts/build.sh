#!/usr/bin/env bash
set -euo pipefail

# Set working directory to the location of .goreleaser.yaml
pwd

# Usage: ./build.sh [snapshot|release]
# Default: snapshot

MODE="${1:-snapshot}"

# Step 1: Generate CLI documentation
DOCS_DIR="docs"
echo "[INFO] Generating CLI documentation into $DOCS_DIR..."
go run ./cmd/suprsend/main.go gendocs "$DOCS_DIR"
echo "[INFO] Documentation generated."

# Step 2: Build with goreleaser
if ! command -v goreleaser &> /dev/null; then
  echo "[ERROR] goreleaser not found. Please install it first." >&2
  exit 1
fi

if [[ "$MODE" == "release" ]]; then
  echo "[INFO] Running goreleaser for RELEASE..."
  goreleaser release --clean
else
  echo "[INFO] Running goreleaser for SNAPSHOT..."
  goreleaser release --snapshot --clean
fi

echo "[SUCCESS] Build complete." 