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

func TestWrapWord_ANSI(t *testing.T) {
	boldText := "\x1b[1;38;5;255mDistributed Consensus\x1b[0m"
	wrapped := WrapWord(boldText, 10)
	for i, l := range wrapped {
		w := VisualWidth(l)
		if w > 10 {
			t.Errorf("line %d width %d exceeds max 10: %q", i, w, l)
		}
		stripped := StripANSI(l)
		if strings.Contains(stripped, "5m") || strings.Contains(stripped, "[") {
			t.Errorf("line %d has leaked escape sequence: %q", i, l)
		}
	}

	// When maxWidth is large enough for individual words (15), words should not be split
	wrapped15 := WrapWord(boldText, 15)
	if len(wrapped15) != 2 {
		t.Fatalf("expected 2 lines for width 15, got %d", len(wrapped15))
	}
	if StripANSI(wrapped15[0]) != "Distributed" {
		t.Errorf("expected first line 'Distributed', got %q", StripANSI(wrapped15[0]))
	}
	if StripANSI(wrapped15[1]) != "Consensus" {
		t.Errorf("expected second line 'Consensus', got %q", StripANSI(wrapped15[1]))
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

func TestTable_WordWrapping_Section13(t *testing.T) {
	for _, width := range []int{80, 100, 120, 150} {
		tbl := New([]string{"Component", "Responsibility & Architecture", "Failure Modes & Risk", "Mitigation Runbook"})
		tbl.SetAlignments([]Alignment{AlignLeft, AlignLeft, AlignLeft, AlignLeft})
		tbl.AddRow(
			"\x1b[1;38;5;255mIngress Controller\x1b[0m",
			"Terminates external TLS connections and directs HTTP traffic across internal Kubernetes service pods.",
			"High connection concurrency causing epoll thread pool starvation and dropped SYN packets.",
			"Scale horizontal replicas and tune net.core.somaxconn and worker connections.",
		)
		tbl.AddRow(
			"\x1b[1;38;5;255mDistributed Consensus\x1b[0m",
			"Raft-based distributed key-value storage maintaining cluster configuration state and leader elections.",
			"Split-brain partition during network transit degradation between availability zones.",
			"Ensure odd quorum voting nodes and verify heartbeat election timeouts.",
		)
		tbl.AddRow(
			"\x1b[1;38;5;255mTime-Series Engine\x1b[0m",
			"High-throughput metrics ingestion pipeline storing Prometheus metrics and alerting telemetry.",
			"Disk IOPS saturation during high-cardinality metric spikes causing ingest backpressure.",
			"Enable write-ahead log compression and drop high-cardinality label dimensions.",
		)
		tbl.MaxWidth = width

		out := tbl.Render()
		lines := strings.Split(out, "\n")
		expectedWidth := VisualWidth(lines[0])

		for i, line := range lines {
			w := VisualWidth(line)
			if w != expectedWidth {
				t.Errorf("width %d: line %d visual width %d != expected width %d:\n%s", width, i, w, expectedWidth, line)
			}
			stripped := StripANSI(line)
			if strings.Contains(stripped, "5m") {
				t.Errorf("width %d: line %d leaked ANSI escape fragment '5m': %s", width, i, stripped)
			}
		}
	}
}

