package table

import (
	"regexp"
	"strings"

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

// WrapWord wraps text into lines of at most maxWidth visual columns.
func WrapWord(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	// Normalize newlines in cell
	text = strings.ReplaceAll(text, "\r\n", "\n")
	paragraphs := strings.Split(text, "\n")
	var result []string

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			result = append(result, "")
			continue
		}

		words := strings.Fields(p)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}

		var currentLine strings.Builder
		currentWidth := 0

		for _, word := range words {
			wLen := VisualWidth(word)

			// If single word exceeds maxWidth, break it by grapheme clusters
			if wLen > maxWidth {
				if currentLine.Len() > 0 {
					result = append(result, currentLine.String())
					currentLine.Reset()
					currentWidth = 0
				}

				var subWord strings.Builder
				subWidth := 0
				g := uniseg.NewGraphemes(word)
				for g.Next() {
					cluster := g.Str()
					cw := g.Width()
					if subWidth+cw > maxWidth && subWidth > 0 {
						result = append(result, subWord.String())
						subWord.Reset()
						subWidth = 0
					}
					subWord.WriteString(cluster)
					subWidth += cw
				}
				if subWord.Len() > 0 {
					currentLine.WriteString(subWord.String())
					currentWidth = subWidth
				}
				continue
			}

			// Check if word fits on current line with a space
			spaceWidth := 0
			if currentWidth > 0 {
				spaceWidth = 1
			}

			if currentWidth+spaceWidth+wLen <= maxWidth {
				if currentWidth > 0 {
					currentLine.WriteByte(' ')
					currentWidth++
				}
				currentLine.WriteString(word)
				currentWidth += wLen
			} else {
				// Flush current line and start new line with word
				result = append(result, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
				currentWidth = wLen
			}
		}

		if currentLine.Len() > 0 {
			result = append(result, currentLine.String())
		}
	}

	if len(result) == 0 {
		return []string{""}
	}

	return result
}
