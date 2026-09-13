package table

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/rivo/uniseg"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\]8;;.*?(?:\x1b\\|\a)|\x1b\]8;;(?:\x1b\\|\a)`)

// StripANSI removes ANSI escape sequences from a string.
func StripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// VisualWidth calculates the terminal display width of a string,
// accounting for double-width characters (CJK, emojis, variation selectors, ZWJ sequences)
// and ignoring ANSI escape codes.
func VisualWidth(s string) int {
	clean := StripANSI(s)
	return uniseg.StringWidth(clean)
}

type ansiToken struct {
	text    string
	isANSI  bool
	isSpace bool
	width   int
}

func tokenizeParagraph(p string) []ansiToken {
	var tokens []ansiToken
	matches := ansiRegex.FindAllStringIndex(p, -1)
	lastIdx := 0

	addPlain := func(plain string) {
		runes := []rune(plain)
		i := 0
		for i < len(runes) {
			if unicode.IsSpace(runes[i]) {
				tokens = append(tokens, ansiToken{text: " ", isSpace: true, width: 1})
				for i < len(runes) && unicode.IsSpace(runes[i]) {
					i++
				}
			} else {
				start := i
				for i < len(runes) && !unicode.IsSpace(runes[i]) {
					i++
				}
				word := string(runes[start:i])
				tokens = append(tokens, ansiToken{text: word, isSpace: false, isANSI: false, width: uniseg.StringWidth(word)})
			}
		}
	}

	for _, m := range matches {
		if m[0] > lastIdx {
			addPlain(p[lastIdx:m[0]])
		}
		tokens = append(tokens, ansiToken{text: p[m[0]:m[1]], isANSI: true, width: 0})
		lastIdx = m[1]
	}
	if lastIdx < len(p) {
		addPlain(p[lastIdx:])
	}

	return tokens
}

// WrapWord wraps text into lines of at most maxWidth visual columns,
// preserving ANSI escape sequences without breaking or leaking color.
func WrapWord(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	// Normalize newlines in cell
	text = strings.ReplaceAll(text, "\r\n", "\n")
	paragraphs := strings.Split(text, "\n")
	var result []string

	for _, p := range paragraphs {
		trimmed := strings.TrimSpace(StripANSI(p))
		if trimmed == "" {
			result = append(result, "")
			continue
		}

		tokens := tokenizeParagraph(p)
		if len(tokens) == 0 {
			result = append(result, "")
			continue
		}

		var currentLine strings.Builder
		currentWidth := 0
		activeStyle := ""
		pendingSpace := false

		flushLine := func() {
			if currentLine.Len() > 0 {
				lineStr := currentLine.String()
				if activeStyle != "" && !strings.HasSuffix(lineStr, "\x1b[0m") {
					lineStr += "\x1b[0m"
				}
				result = append(result, lineStr)
				currentLine.Reset()
				currentWidth = 0
			}
		}

		startNewLine := func() {
			flushLine()
			if activeStyle != "" {
				currentLine.WriteString(activeStyle)
			}
			pendingSpace = false
		}

		for _, tok := range tokens {
			if tok.isANSI {
				currentLine.WriteString(tok.text)
				if tok.text == "\x1b[0m" {
					activeStyle = ""
				} else if strings.HasPrefix(tok.text, "\x1b[") && strings.HasSuffix(tok.text, "m") {
					if tok.text == "\x1b[m" {
						activeStyle = ""
					} else {
						activeStyle = tok.text
					}
				}
				continue
			}

			if tok.isSpace {
				if currentWidth > 0 {
					pendingSpace = true
				}
				continue
			}

			// Text word
			word := tok.text
			wLen := tok.width

			if wLen > maxWidth {
				// Single word exceeds maxWidth - flush current line if any, then break by graphemes
				if currentWidth > 0 {
					startNewLine()
				}
				g := uniseg.NewGraphemes(word)
				for g.Next() {
					cluster := g.Str()
					cw := g.Width()
					if currentWidth+cw > maxWidth && currentWidth > 0 {
						startNewLine()
					}
					currentLine.WriteString(cluster)
					currentWidth += cw
				}
				pendingSpace = false
				continue
			}

			spaceWidth := 0
			if pendingSpace && currentWidth > 0 {
				spaceWidth = 1
			}

			if currentWidth+spaceWidth+wLen <= maxWidth {
				if spaceWidth > 0 {
					currentLine.WriteByte(' ')
					currentWidth++
				}
				currentLine.WriteString(word)
				currentWidth += wLen
				pendingSpace = false
			} else {
				startNewLine()
				currentLine.WriteString(word)
				currentWidth = wLen
				pendingSpace = false
			}
		}

		flushLine()
	}

	if len(result) == 0 {
		return []string{""}
	}

	return result
}
