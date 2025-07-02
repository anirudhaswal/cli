#!/usr/bin/env bash

# Exit if any command fails
set -e

INPUT_SCRIPT="main.ts"

# Confirm the input script exists
if [ ! -f "$INPUT_SCRIPT" ]; then
  echo "❌ Error: $INPUT_SCRIPT not found."
  exit 1
fi

# Architectures your system supports
ARCH=$(uname -m)

if [[ "$ARCH" != "x86_64" ]]; then
  echo "❌ This script is only for x86_64 hosts. Detected: $ARCH"
  exit 1
fi

# Supported targets for x86_64
TARGETS=(
  "x86_64-unknown-linux-gnu"    # Linux
  "x86_64-apple-darwin"         # macOS Intel
  "x86_64-pc-windows-msvc"      # Windows
)

# Loop over targets
for TARGET in "${TARGETS[@]}"; do
  echo ""
  echo "⚙️ Compiling for $TARGET..."

  # Determine output filename
  case $TARGET in
    *windows*)
      OUT="typemorph-${TARGET}.exe"
      ;;
    *darwin*)
      OUT="typemorph-${TARGET}-mac"
      ;;
    *linux*)
      OUT="typemorph-${TARGET}-linux"
      ;;
    *)
      OUT="typemorph-${TARGET}"
      ;;
  esac

  deno compile \
    --target "$TARGET" \
    --allow-read \
    --allow-write \
    --allow-net \
    --output "$OUT" \
    "$INPUT_SCRIPT"

  echo "✅ Built $OUT"
done

echo ""
echo "🎉 All binaries compiled!"
