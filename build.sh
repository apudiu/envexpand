#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

OUT_DIR="$SCRIPT_DIR/bin"
mkdir -p "$OUT_DIR"

echo "Building envexpand for Linux, macOS, and Windows (amd64 + arm64)..."

targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
  "windows arm64"
)

for target in "${targets[@]}"; do
  read -r goos goarch <<<"$target"

  outfile="${OUT_DIR}/envexpand_${goos}_${goarch}"
  if [[ "$goos" == "windows" ]]; then
    outfile+=".exe"
  fi

  echo "  -> ${goos}/${goarch}"
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags "-s -w -buildid=" -o "$outfile" .
done

echo "Done:"
echo "  $OUT_DIR/envexpand_linux_amd64"
echo "  $OUT_DIR/envexpand_linux_arm64"
echo "  $OUT_DIR/envexpand_darwin_amd64"
echo "  $OUT_DIR/envexpand_darwin_arm64"
echo "  $OUT_DIR/envexpand_windows_amd64.exe"
echo "  $OUT_DIR/envexpand_windows_arm64.exe"
