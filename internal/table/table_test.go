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
		{"\x1b[31mhello\x1b[0m", 5},                                // ANSI color
		{"\x1b]8;;https://example.com\x1b\\link\x1b]8;;\x1b\\", 4},  // OSC 8 link
		{"🚀 rocket", 9},                                           // 2 for emoji + 1 space + 6 chars = 9
		{"你好", 4},                                                  // 2 + 2 = 4
		{"⚠️", 2},                                                   // Emoji with variation selector 16 (\u26a0\ufe0f)
		{"⚠️ Degraded", 11},                                         // 2 + 1 + 8 = 11
		{"✅ Healthy", 10},                                          // 2 + 1 + 7 = 10
		{"🛑 Paused", 9},                                            // 2 + 1 + 6 = 9
		{"🚀 Optimal", 10},                                          // 2 + 1 + 7 = 10
	}

	for _, tt := range tests {
		got := VisualWidth(tt.input)
		if got != tt.want {
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

func TestTable_MicroserviceSLAMatrix(t *testing.T) {
	tbl := New([]string{"Service Name", "Cluster ID", "Health", "Uptime SLA", "p50 Latency", "p99 Latency", "Error Rate"})
	tbl.SetAlignments([]Alignment{AlignLeft, AlignCenter, AlignCenter, AlignCenter, AlignRight, AlignRight, AlignRight})
	tbl.AddRow("api-gateway", "us-east-1a", "✅ Healthy", "99.99%", "1.2ms", "3.4ms", "0.001%")
	tbl.AddRow("auth-service", "us-east-1b", "✅ Healthy", "99.95%", "8.5ms", "18.2ms", "0.012%")
	tbl.AddRow("payment-processor", "us-west-2a", "⚠️ Degraded", "99.99%", "45.0ms", "182.4ms", "0.350%")
	tbl.AddRow("search-indexing", "eu-west-1c", "🛑 Paused", "99.90%", "210.0ms", "940.0ms", "2.100%")
	tbl.AddRow("cache-redis-l1", "us-east-1a", "🚀 Optimal", "99.999%", "0.2ms", "0.7ms", "0.000%")
	tbl.MaxWidth = 120

	out := tbl.Render()
	lines := strings.Split(out, "\n")
	expectedWidth := VisualWidth(lines[0])

	for i, line := range lines {
		w := VisualWidth(line)
		if w != expectedWidth {
			t.Errorf("line %d visual width %d != expected width %d:\n%s", i, w, expectedWidth, line)
		}
	}
}

