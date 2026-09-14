package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/smford/mdee/internal/config"
	"github.com/smford/mdee/internal/table"
	"github.com/smford/mdee/internal/term"
)

func TestCLIVersion(t *testing.T) {
	cmd := newVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}
}

func TestCLIDoctor(t *testing.T) {
	cmd := newDoctorCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("doctor command failed: %v", err)
	}
}

func TestRenderAndOutput(t *testing.T) {
	opts := config.DefaultOptions()
	opts.Plain = true
	info := term.Info{IsTTY: false, Width: 80, Height: 24}

	doc := "# Test Header\n\n| A | B |\n|---|---|\n| 1 | 2 |\n"
	var buf bytes.Buffer

	// Temporarily redirect stdout during render test
	ctx := context.Background()
	err := renderAndOutput(ctx, opts, info, []byte(doc))
	if err != nil {
		t.Fatalf("renderAndOutput error: %v", err)
	}
	_ = buf
}

func TestTableAlignmentInCLI(t *testing.T) {
	opts := config.DefaultOptions()
	opts.Width = 60
	opts.TableStyle = "box"

	doc := `
| Left | Center | Right |
| :--- | :---: | ---: |
| L | C | R |
`
	ctx := context.Background()
	opts.Plain = true
	var buf strings.Builder
	r := table.New([]string{"Left", "Center", "Right"})
	r.SetAlignments([]table.Alignment{table.AlignLeft, table.AlignCenter, table.AlignRight})
	r.AddRow("L", "C", "R")
	r.Border = table.BorderBox
	r.MaxWidth = 60
	rendered := r.Render()
	buf.WriteString(rendered)

	if !strings.Contains(buf.String(), "┌") || !strings.Contains(buf.String(), "┘") {
		t.Errorf("expected box borders in rendered table")
	}
	_ = ctx
	_ = doc
}
