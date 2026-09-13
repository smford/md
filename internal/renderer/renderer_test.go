package renderer

import (
	"context"
	"strings"
	"testing"

	"github.com/smford/md/internal/config"
	"github.com/smford/md/internal/table"
)

func TestRenderer_BasicElements(t *testing.T) {
	opts := config.DefaultOptions()
	opts.Width = 80
	r := New(opts)

	doc := `
# Main Heading

This is a paragraph with **bold**, *italic*, and ` + "`inline code`" + `.

## Features

- Fast
- Reliable
- Graceful degradation

1. First step
2. Second step

- [x] Completed task
- [ ] Pending task

> SRE reliability is paramount.
> Always measure before optimizing.

---

` + "```go\npackage main\n\nfunc main() {}\n```"

	out, err := r.Render(context.Background(), []byte(doc))
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	t.Logf("OUTPUT:\n%s\n", out)

	cleanOut := table.StripANSI(out)
	if !strings.Contains(cleanOut, "Main Heading") {
		t.Errorf("expected output to contain 'Main Heading'")
	}
	if !strings.Contains(cleanOut, "Features") {
		t.Errorf("expected output to contain 'Features'")
	}
	if !strings.Contains(cleanOut, "Fast") || !strings.Contains(cleanOut, "Reliable") {
		t.Errorf("expected list items in output")
	}
	if !strings.Contains(cleanOut, "First step") {
		t.Errorf("expected ordered list items in output")
	}
	if !strings.Contains(cleanOut, "SRE reliability is paramount") {
		t.Errorf("expected blockquote in output")
	}
	if !strings.Contains(cleanOut, "package main") {
		t.Errorf("expected code block in output")
	}
}

func TestRenderer_Table(t *testing.T) {
	opts := config.DefaultOptions()
	opts.Width = 80
	opts.TableStyle = "rounded"
	r := New(opts)

	doc := `
| Service | Status | Latency (p99) |
| :--- | :---: | ---: |
| frontend-proxy | Healthy | 4.2ms |
| database-primary | Degraded | 185.0ms |
| cache-redis | Healthy | 0.8ms |
`

	out, err := r.Render(context.Background(), []byte(doc))
	if err != nil {
		t.Fatalf("Render table failed: %v", err)
	}

	if !strings.Contains(out, "frontend-proxy") {
		t.Errorf("expected 'frontend-proxy' in table output")
	}
	if !strings.Contains(out, "185.0ms") {
		t.Errorf("expected '185.0ms' in table output")
	}
	if !strings.Contains(out, "╭") || !strings.Contains(out, "╯") {
		t.Errorf("expected rounded table borders")
	}
}

func TestRenderer_Image(t *testing.T) {
	opts := config.DefaultOptions()
	opts.ImageMode = "never" // Fallback box
	r := New(opts)

	doc := `![Architecture Diagram](https://example.com/arch.png "Arch")`
	out, err := r.Render(context.Background(), []byte(doc))
	if err != nil {
		t.Fatalf("Render image failed: %v", err)
	}

	if !strings.Contains(out, "Architecture Diagram") || !strings.Contains(out, "https://example.com/arch.png") {
		t.Errorf("expected image fallback box in output: %s", out)
	}
}

func TestRenderer_Plain(t *testing.T) {
	opts := config.DefaultOptions()
	opts.Plain = true
	r := New(opts)

	doc := `# Title
- [x] Done
`
	out, err := r.Render(context.Background(), []byte(doc))
	if err != nil {
		t.Fatalf("Render plain failed: %v", err)
	}

	if strings.Contains(out, "\x1b[") {
		t.Errorf("expected plain output without ANSI escape sequences: %q", out)
	}
	if !strings.Contains(out, "[x]") {
		t.Errorf("expected '[x]' in plain task list: %s", out)
	}
}

func TestRenderer_Frontmatter(t *testing.T) {
	opts := config.DefaultOptions()
	r := New(opts)

	doc := `---
title: Post Title
date: 2026-09-13
---

# Content
`
	out, err := r.Render(context.Background(), []byte(doc))
	if err != nil {
		t.Fatalf("Render frontmatter failed: %v", err)
	}
	if !strings.Contains(out, "Metadata") || !strings.Contains(out, "Post Title") {
		t.Errorf("expected frontmatter card: %s", out)
	}

	// Verify lines inside metadata box start with "│" and end with "│"
	lines := strings.Split(out, "\n")
	foundTitle := false
	for _, l := range lines {
		if strings.Contains(l, "Post Title") {
			foundTitle = true
			trimmed := strings.TrimSpace(l)
			if !strings.HasPrefix(trimmed, "│") || !strings.HasSuffix(trimmed, "│") {
				t.Errorf("expected metadata row enclosed in vertical bars, got: %q", l)
			}
			if strings.HasPrefix(l, "       ") {
				t.Errorf("detected excessive indentation on metadata line: %q", l)
			}
		}
	}
	if !foundTitle {
		t.Errorf("title row not found in output: %s", out)
	}
}
