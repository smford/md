# typed: false
# frozen_string_literal: true

class Md < Formula
  desc "Terminal Markdown viewer for macOS iTerm2"
  homepage "https://github.com/smford/md"
  version "1.0.0"
  license "AGPL-3.0-or-later"

  on_macos do
    on_arm do
      url "https://github.com/smford/md/releases/download/v#{version}/md-v#{version}-darwin-arm64.tar.gz"
      sha256 "9f47af4c5f2fe44475fd9bd8efc68f5b0e26a92b77748709abecb6fc2465019e"
    end
    on_intel do
      url "https://github.com/smford/md/releases/download/v#{version}/md-v#{version}-darwin-amd64.tar.gz"
      sha256 "17a3b4ddda5f1ac96445a567fdbf39707bd32f29bf3d7865b757f5c5b0967d5d"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/smford/md/releases/download/v#{version}/md-v#{version}-linux-arm64.tar.gz"
      sha256 "5f79430b220e0f4ab208f8683821ef470ba36a69167dea80da1bac2e0f8e702d"
    end
    on_intel do
      url "https://github.com/smford/md/releases/download/v#{version}/md-v#{version}-linux-amd64.tar.gz"
      sha256 "44590a5dfce5338295728811482c89ddd859f6e1b058af7837d76457d23411f7"
    end
  end

  def install
    bin.install "md"
  end

  test do
    assert_match "md version", shell_output("#{bin}/md --version")
  end
end
