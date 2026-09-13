package term

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

// Info holds detected terminal capabilities.
type Info struct {
	IsTTY        bool
	Width        int
	Height       int
	IsITerm2     bool
	IsTmux       bool
	TermProgram  string
	TermName     string
	HasTrueColor bool
	HasOSC8      bool
	HasOSC1337   bool
}

// IsTerminal checks if the provided file descriptor is connected to a terminal.
func IsTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

// GetSize returns the current terminal width and height.
// If not a terminal or an error occurs, it falls back to sensible defaults (80x24).
func GetSize() (width, height int, err error) {
	if IsTerminal(os.Stdout) {
		w, h, e := term.GetSize(int(os.Stdout.Fd()))
		if e == nil && w > 0 && h > 0 {
			return w, h, nil
		}
	}
	if IsTerminal(os.Stderr) {
		w, h, e := term.GetSize(int(os.Stderr.Fd()))
		if e == nil && w > 0 && h > 0 {
			return w, h, nil
		}
	}
	return 80, 24, fmt.Errorf("could not determine terminal size, using fallback")
}

// Detect inspects environment variables and terminal file descriptors.
func Detect() Info {
	isTTY := IsTerminal(os.Stdout)
	w, h, _ := GetSize()

	termProgram := os.Getenv("TERM_PROGRAM")
	termName := os.Getenv("TERM")
	lcTerminal := os.Getenv("LC_TERMINAL")
	itermSession := os.Getenv("ITERM_SESSION_ID")
	tmuxEnv := os.Getenv("TMUX")

	isTmux := tmuxEnv != "" || strings.HasPrefix(termName, "tmux") || strings.HasPrefix(termName, "screen")

	// Detect iTerm2 directly or when running inside tmux launched from iTerm2
	isITerm := termProgram == "iTerm.app" ||
		lcTerminal == "iTerm2" ||
		itermSession != ""

	// If inside tmux, TERM_PROGRAM might not be set unless passthrough is enabled,
	// but ITERM_SESSION_ID or LC_TERMINAL often propagates.
	if !isITerm && isTmux {
		if itermSession != "" || lcTerminal == "iTerm2" {
			isITerm = true
		}
	}

	// WezTerm also supports OSC 1337 (iTerm2 image protocol)
	supportsOSC1337 := isITerm || termProgram == "WezTerm"

	// Modern terminals supporting OSC 8 hyperlinks
	supportsOSC8 := isITerm ||
		termProgram == "WezTerm" ||
		termProgram == "vscode" ||
		termProgram == "ghostty" ||
		termProgram == "Apple_Terminal" ||
		os.Getenv("VTE_VERSION") != "" ||
		strings.Contains(termName, "kitty") ||
		strings.Contains(termName, "alacritty")

	// Truecolor detection
	colorTerm := os.Getenv("COLORTERM")
	hasTrueColor := colorTerm == "truecolor" || colorTerm == "24bit" || isITerm || termProgram == "ghostty" || termProgram == "WezTerm"

	return Info{
		IsTTY:        isTTY,
		Width:        w,
		Height:       h,
		IsITerm2:     isITerm,
		IsTmux:       isTmux,
		TermProgram:  termProgram,
		TermName:     termName,
		HasTrueColor: hasTrueColor,
		HasOSC8:      supportsOSC8,
		HasOSC1337:   supportsOSC1337,
	}
}

// FormatHyperlink formats an OSC 8 hyperlink sequence if supported.
func FormatHyperlink(url, text string, enable bool) string {
	if !enable || url == "" {
		if text != "" {
			return text
		}
		return url
	}
	// OSC 8 ;; URL \ ESC \ text ESC ] 8 ;; ESC \
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

// RunPager launches a pager (such as $PAGER or less -R -F -X) and streams content to it.
func RunPager(content string) error {
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}

	parts := strings.Fields(pager)
	cmdName := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	// If using less and no explicit args provided, default to raw control chars (-R),
	// quit if one screen (-F), and do not clear screen on exit (-X).
	if cmdName == "less" && len(args) == 0 {
		args = []string{"-R", "-F", "-X"}
	}

	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		_, _ = io.WriteString(stdin, content)
	}()

	return cmd.Wait()
}
