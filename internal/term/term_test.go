package term

import (
	"os"
	"testing"
)

func TestDetect(t *testing.T) {
	// Set dummy env
	os.Setenv("TERM_PROGRAM", "iTerm.app")
	defer os.Unsetenv("TERM_PROGRAM")

	info := Detect()
	if !info.IsITerm2 {
		t.Errorf("expected IsITerm2 to be true when TERM_PROGRAM=iTerm.app")
	}
	if !info.HasOSC1337 {
		t.Errorf("expected HasOSC1337 to be true for iTerm2")
	}
	if !info.HasOSC8 {
		t.Errorf("expected HasOSC8 to be true for iTerm2")
	}
	if info.Width <= 0 {
		t.Errorf("expected positive width, got %d", info.Width)
	}
}

func TestFormatHyperlink(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		text    string
		enable  bool
		want    string
	}{
		{
			name:   "disabled",
			url:    "https://example.com",
			text:   "Example",
			enable: false,
			want:   "Example",
		},
		{
			name:   "enabled",
			url:    "https://example.com",
			text:   "Example",
			enable: true,
			want:   "\x1b]8;;https://example.com\aExample\x1b]8;;\a",
		},
		{
			name:   "empty url",
			url:    "",
			text:   "Just Text",
			enable: true,
			want:   "Just Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatHyperlink(tt.url, tt.text, tt.enable)
			if got != tt.want {
				t.Errorf("FormatHyperlink() = %q, want %q", got, tt.want)
			}
		})
	}
}
