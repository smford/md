package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines the visual styling rules for markdown elements.
type Theme struct {
	Name             string
	H1               lipgloss.Style
	H2               lipgloss.Style
	H3               lipgloss.Style
	H4               lipgloss.Style
	H5               lipgloss.Style
	H6               lipgloss.Style
	Paragraph        lipgloss.Style
	Bold             lipgloss.Style
	Italic           lipgloss.Style
	Strikethrough    lipgloss.Style
	Blockquote       lipgloss.Style
	BlockquoteBorder lipgloss.Style
	CodeSpan         lipgloss.Style
	CodeBlockBorder  lipgloss.Style
	CodeBlockTitle   lipgloss.Style
	ListBullet       lipgloss.Style
	TaskDone         lipgloss.Style
	TaskTodo         lipgloss.Style
	Link             lipgloss.Style
	LinkURL          lipgloss.Style
	HorizontalRule   lipgloss.Style
	TableBorder      lipgloss.Style
	TableHeader      lipgloss.Style
	TableCell        lipgloss.Style
	ChromaTheme      string
}

// AvailableThemes lists supported theme names.
func AvailableThemes() []string {
	return []string{"dark", "light", "dracula", "monokai", "solarized-dark", "solarized-light", "plain"}
}

// GetTheme resolves a theme by name or returns plain/dark default.
func GetTheme(name string, plain bool) *Theme {
	if plain || strings.ToLower(name) == "plain" {
		return plainTheme()
	}

	switch strings.ToLower(name) {
	case "light":
		return lightTheme()
	case "dracula":
		return draculaTheme()
	case "monokai":
		return monokaiTheme()
	case "solarized-dark":
		return solarizedDarkTheme()
	case "solarized-light":
		return solarizedLightTheme()
	case "dark":
		fallthrough
	default:
		return darkTheme()
	}
}

func plainTheme() *Theme {
	empty := lipgloss.NewStyle()
	return &Theme{
		Name:             "plain",
		H1:               empty,
		H2:               empty,
		H3:               empty,
		H4:               empty,
		H5:               empty,
		H6:               empty,
		Paragraph:        empty,
		Bold:             empty,
		Italic:           empty,
		Strikethrough:    empty,
		Blockquote:       empty,
		BlockquoteBorder: empty,
		CodeSpan:         empty,
		CodeBlockBorder:  empty,
		CodeBlockTitle:   empty,
		ListBullet:       empty,
		TaskDone:         empty,
		TaskTodo:         empty,
		Link:             empty,
		LinkURL:          empty,
		HorizontalRule:   empty,
		TableBorder:      empty,
		TableHeader:      empty,
		TableCell:        empty,
		ChromaTheme:      "bw",
	}
}

func darkTheme() *Theme {
	return &Theme{
		Name:             "dark",
		H1:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")),
		H2:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75")),
		H3:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("111")),
		H4:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("147")),
		H5:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("183")),
		H6:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("246")),
		Paragraph:        lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		Bold:             lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")),
		Italic:           lipgloss.NewStyle().Italic(true),
		Strikethrough:    lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("243")),
		Blockquote:       lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("248")),
		BlockquoteBorder: lipgloss.NewStyle().Foreground(lipgloss.Color("39")),
		CodeSpan:         lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("236")).Padding(0, 1),
		CodeBlockBorder:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		CodeBlockTitle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75")),
		ListBullet:       lipgloss.NewStyle().Foreground(lipgloss.Color("39")),
		TaskDone:         lipgloss.NewStyle().Foreground(lipgloss.Color("78")),
		TaskTodo:         lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Link:             lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Underline(true),
		LinkURL:          lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		HorizontalRule:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		TableBorder:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		TableHeader:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75")),
		TableCell:        lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		ChromaTheme:      "dracula",
	}
}

func lightTheme() *Theme {
	return &Theme{
		Name:             "light",
		H1:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("27")),
		H2:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("26")),
		H3:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("25")),
		H4:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("24")),
		H5:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("23")),
		H6:               lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("240")),
		Paragraph:        lipgloss.NewStyle().Foreground(lipgloss.Color("235")),
		Bold:             lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("16")),
		Italic:           lipgloss.NewStyle().Italic(true),
		Strikethrough:    lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("244")),
		Blockquote:       lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("239")),
		BlockquoteBorder: lipgloss.NewStyle().Foreground(lipgloss.Color("31")),
		CodeSpan:         lipgloss.NewStyle().Foreground(lipgloss.Color("161")).Background(lipgloss.Color("254")).Padding(0, 1),
		CodeBlockBorder:  lipgloss.NewStyle().Foreground(lipgloss.Color("250")),
		CodeBlockTitle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("26")),
		ListBullet:       lipgloss.NewStyle().Foreground(lipgloss.Color("27")),
		TaskDone:         lipgloss.NewStyle().Foreground(lipgloss.Color("28")),
		TaskTodo:         lipgloss.NewStyle().Foreground(lipgloss.Color("243")),
		Link:             lipgloss.NewStyle().Foreground(lipgloss.Color("26")).Underline(true),
		LinkURL:          lipgloss.NewStyle().Foreground(lipgloss.Color("246")),
		HorizontalRule:   lipgloss.NewStyle().Foreground(lipgloss.Color("250")),
		TableBorder:      lipgloss.NewStyle().Foreground(lipgloss.Color("248")),
		TableHeader:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("26")),
		TableCell:        lipgloss.NewStyle().Foreground(lipgloss.Color("236")),
		ChromaTheme:      "github",
	}
}

func draculaTheme() *Theme {
	t := darkTheme()
	t.Name = "dracula"
	t.H1 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	t.H2 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff79c6"))
	t.H3 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	t.H4 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	t.BlockquoteBorder = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9"))
	t.CodeSpan = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c")).Background(lipgloss.Color("#282a36")).Padding(0, 1)
	t.ListBullet = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9"))
	t.TableHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff79c6"))
	t.Link = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Underline(true)
	t.ChromaTheme = "dracula"
	return t
}

func monokaiTheme() *Theme {
	t := darkTheme()
	t.Name = "monokai"
	t.H1 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fd971f"))
	t.H2 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#a6e22e"))
	t.H3 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#66d9ef"))
	t.H4 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f92672"))
	t.BlockquoteBorder = lipgloss.NewStyle().Foreground(lipgloss.Color("#fd971f"))
	t.CodeSpan = lipgloss.NewStyle().Foreground(lipgloss.Color("#e6db74")).Background(lipgloss.Color("#272822")).Padding(0, 1)
	t.ListBullet = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e22e"))
	t.TableHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#66d9ef"))
	t.Link = lipgloss.NewStyle().Foreground(lipgloss.Color("#66d9ef")).Underline(true)
	t.ChromaTheme = "monokai"
	return t
}

func solarizedDarkTheme() *Theme {
	t := darkTheme()
	t.Name = "solarized-dark"
	t.H1 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#268bd2"))
	t.H2 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2aa198"))
	t.H3 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#b58900"))
	t.ChromaTheme = "solarized-dark"
	return t
}

func solarizedLightTheme() *Theme {
	t := lightTheme()
	t.Name = "solarized-light"
	t.H1 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#268bd2"))
	t.H2 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2aa198"))
	t.H3 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#b58900"))
	t.ChromaTheme = "solarized-light"
	return t
}
