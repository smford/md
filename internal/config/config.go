package config

// Options defines runtime rendering configuration.
type Options struct {
	Width       int    // Available width (0 means auto-detect terminal width)
	Theme       string // Color theme name (dark, light, dracula, monokai, solarized-dark, solarized-light, plain)
	TableStyle  string // Border style: rounded, box, double, ascii, markdown, minimal
	ImageMode   string // "auto", "always", "never"
	ImageWidth  string // "auto", "100%", "80", "400px"
	ImageHeight string // "auto", "20", "300px"
	LineNumbers bool   // Display line numbers in code blocks
	Hyperlinks  bool   // Enable OSC 8 clickable terminal hyperlinks
	Pager       bool   // Enable pager for long output in interactive TTY
	Plain       bool   // Plain text output with no ANSI escape sequences
	Debug       bool   // Output diagnostic / debug logs to stderr
	BasePath    string // Directory of input markdown file for resolving relative assets
}

// DefaultOptions returns standard SRE defaults.
func DefaultOptions() Options {
	return Options{
		Width:       0,
		Theme:       "dark",
		TableStyle:  "rounded",
		ImageMode:   "auto",
		ImageWidth:  "auto",
		ImageHeight: "auto",
		LineNumbers: false,
		Hyperlinks:  true,
		Pager:       false,
		Plain:       false,
		Debug:       false,
		BasePath:    "",
	}
}
