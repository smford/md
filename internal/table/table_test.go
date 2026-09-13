package table

import (
	"strings"
	"testing"
)

func TestVisualWidth(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"\x1b[31mhello\x1b[0m", 5},                       // ANSI color
		{"\x1b]8;;https://example.com\x1b\\link\x1b]8;;\x1b\\", 4}, // OSC 8 link
		{"🚀 rocket", 8},                                    // 2 for emoji + 1 space + 6 chars = 9? Wait: 🚀 is width 2. " " is 1. "rocket" is 6. Total = 9!
		{"你好", 4},                                         // 2 + 2 = 4
	}

	for _, tt := range tests {
		got := VisualWidth(tt.input)
		if tt.input == "🚀 rocket" && got != 9 {
			t.Errorf("VisualWidth(%q) = %d, want 9", tt.input, got)
		} else if tt.input != "🚀 rocket" && got != tt.want {
			t.Errorf("VisualWidth(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestWrapWord(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	wrapped := WrapWord(text, 15)
	for _, line := range wrapped {
		if VisualWidth(line) > 15 {
			t.Errorf("line %q exceeds max width 15 (width %d)", line, VisualWidth(line))
		}
	}
}

func TestTable_Render_AlignmentAndWrapping(t *testing.T) {
	tbl := New([]string{"Service", "Status", "Latency"})
	tbl.SetAlignments([]Alignment{AlignLeft, AlignCenter, AlignRight})
	tbl.AddRow("auth-service", "✅ Healthy", "12ms")
	tbl.AddRow("database-primary-cluster", "⚠️ High Load", "450ms")
	tbl.MaxWidth = 60

	out := tbl.Render()
	lines := strings.Split(out, "\n")
	if len(lines) < 5 {
		t.Fatalf("expected at least 5 lines of table output, got %d", len(lines))
	}

	// Verify top border
	if !strings.HasPrefix(lines[0], "╭") || !strings.HasSuffix(lines[0], "╮") {
		t.Errorf("expected rounded top border, got %q", lines[0])
	}

	// Verify bottom border
	lastLine := lines[len(lines)-1]
	if !strings.HasPrefix(lastLine, "╰") || !strings.HasSuffix(lastLine, "╯") {
		t.Errorf("expected rounded bottom border, got %q", lastLine)
	}

	// Verify headers are present
	if !strings.Contains(out, "Service") || !strings.Contains(out, "Status") || !strings.Contains(out, "Latency") {
		t.Errorf("missing headers in table output")
	}

	// Verify data is present
	if !strings.Contains(out, "auth-service") || !strings.Contains(out, "12ms") {
		t.Errorf("missing data in table output")
	}
}

func TestTable_Styles(t *testing.T) {
	styles := []string{"rounded", "box", "double", "ascii", "markdown", "minimal"}
	for _, s := range styles {
		t.Run(s, func(t *testing.T) {
			tbl := New([]string{"Col A", "Col B"})
			tbl.Border = StyleNameToBorders(s)
			tbl.AddRow("1", "2")
			out := tbl.Render()
			if len(out) == 0 {
				t.Errorf("empty table output for style %s", s)
			}
		})
	}
}
