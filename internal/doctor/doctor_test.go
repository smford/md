package doctor

import (
	"strings"
	"testing"
)

func TestRunDiagnostics(t *testing.T) {
	report := RunDiagnostics()
	if !strings.Contains(report, "Terminal Diagnostic") {
		t.Errorf("expected diagnostic header in report")
	}
	if !strings.Contains(report, "OSC 1337") {
		t.Errorf("expected OSC 1337 check in report")
	}
	if !strings.Contains(report, "OSC 8") {
		t.Errorf("expected OSC 8 check in report")
	}
}
