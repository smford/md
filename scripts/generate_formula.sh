#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-1.0.0}"
TAG_NAME="${2:-v${VERSION}}"
DIST_DIR="${3:-dist}"
OUT_FILE="${4:-${DIST_DIR}/mdee.rb}"

get_sha() {
  local file="$1"
  if [ -f "$file" ]; then
    sha256sum "$file" | cut -d' ' -f1
  else
    echo "REPLACE_WITH_SHA256"
  fi
}

SHA_DARWIN_ARM64=$(get_sha "${DIST_DIR}/mdee-${TAG_NAME}-darwin-arm64.tar.gz")
SHA_DARWIN_AMD64=$(get_sha "${DIST_DIR}/mdee-${TAG_NAME}-darwin-amd64.tar.gz")
SHA_LINUX_ARM64=$(get_sha "${DIST_DIR}/mdee-${TAG_NAME}-linux-arm64.tar.gz")
SHA_LINUX_AMD64=$(get_sha "${DIST_DIR}/mdee-${TAG_NAME}-linux-amd64.tar.gz")

mkdir -p "$(dirname "$OUT_FILE")"

cat <<EOF > "$OUT_FILE"
# typed: false
# frozen_string_literal: true

# This formula was auto-generated for mdee (https://github.com/smford/mdee).
class Mdee < Formula
  desc "Terminal Markdown viewer for macOS iTerm2"
  homepage "https://github.com/smford/mdee"
  version "${VERSION}"
  license "AGPL-3.0-or-later"

  on_macos do
    on_arm do
      url "https://github.com/smford/mdee/releases/download/v#{version}/mdee-v#{version}-darwin-arm64.tar.gz"
      sha256 "${SHA_DARWIN_ARM64}"
    end
    on_intel do
      url "https://github.com/smford/mdee/releases/download/v#{version}/mdee-v#{version}-darwin-amd64.tar.gz"
      sha256 "${SHA_DARWIN_AMD64}"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/smford/mdee/releases/download/v#{version}/mdee-v#{version}-linux-arm64.tar.gz"
      sha256 "${SHA_LINUX_ARM64}"
    end
    on_intel do
      url "https://github.com/smford/mdee/releases/download/v#{version}/mdee-v#{version}-linux-amd64.tar.gz"
      sha256 "${SHA_LINUX_AMD64}"
    end
  end

  def install
    bin.install "mdee"
  end

  test do
    assert_match "mdee version", shell_output("#{bin}/mdee --version")
  end
end
EOF

echo "Generated Homebrew formula at ${OUT_FILE}"
