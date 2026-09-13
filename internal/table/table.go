package table

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Alignment defines the text alignment within a table column.
type Alignment int

const (
	AlignNone Alignment = iota
	AlignLeft
	AlignCenter
	AlignRight
)

// BorderCharacters defines the glyphs used to construct table frames.
type BorderCharacters struct {
	TopLeft         string
	TopMid          string
	TopRight        string
	MidLeft         string
	MidMid          string
	MidRight        string
	BottomLeft      string
	BottomMid       string
	BottomRight     string
	Horizontal      string
	Vertical        string
	HeaderSeparator string
}

// Predefined border styles.
var (
	BorderRounded = BorderCharacters{
		TopLeft: "╭", TopMid: "┬", TopRight: "╮",
		MidLeft: "├", MidMid: "┼", MidRight: "┤",
		BottomLeft: "╰", BottomMid: "┴", BottomRight: "╯",
		Horizontal: "─", Vertical: "│", HeaderSeparator: "─",
	}

	BorderBox = BorderCharacters{
		TopLeft: "┌", TopMid: "┬", TopRight: "┐",
		MidLeft: "├", MidMid: "┼", MidRight: "┤",
		BottomLeft: "└", BottomMid: "┴", BottomRight: "┘",
		Horizontal: "─", Vertical: "│", HeaderSeparator: "─",
	}

	BorderDouble = BorderCharacters{
		TopLeft: "╔", TopMid: "╦", TopRight: "╗",
		MidLeft: "╠", MidMid: "╬", MidRight: "╣",
		BottomLeft: "╚", BottomMid: "╩", BottomRight: "╝",
		Horizontal: "═", Vertical: "║", HeaderSeparator: "═",
	}

	BorderASCII = BorderCharacters{
		TopLeft: "+", TopMid: "+", TopRight: "+",
		MidLeft: "+", MidMid: "+", MidRight: "+",
		BottomLeft: "+", BottomMid: "+", BottomRight: "+",
		Horizontal: "-", Vertical: "|", HeaderSeparator: "-",
	}

	BorderMarkdown = BorderCharacters{
		TopLeft: "", TopMid: "", TopRight: "",
		MidLeft: "|", MidMid: "|", MidRight: "|",
		BottomLeft: "", BottomMid: "", BottomRight: "",
		Horizontal: "-", Vertical: "|", HeaderSeparator: "-",
	}

	BorderMinimal = BorderCharacters{
		TopLeft: "", TopMid: "", TopRight: "",
		MidLeft: "", MidMid: "  ", MidRight: "",
		BottomLeft: "", BottomMid: "", BottomRight: "",
		Horizontal: "─", Vertical: "  ", HeaderSeparator: "─",
	}
)

// StyleNameToBorders maps CLI style names to BorderCharacters.
func StyleNameToBorders(name string) BorderCharacters {
	switch strings.ToLower(name) {
	case "box", "sharp":
		return BorderBox
	case "double":
		return BorderDouble
	case "ascii":
		return BorderASCII
	case "markdown":
		return BorderMarkdown
	case "minimal", "simple":
		return BorderMinimal
	case "rounded":
		fallthrough
	default:
		return BorderRounded
	}
}

// Table represents a complete tabular data structure ready for terminal rendering.
type Table struct {
	Headers     []string
	Alignments  []Alignment
	Rows        [][]string
	MaxWidth    int
	Border      BorderCharacters
	BorderStyle lipgloss.Style
	HeaderStyle lipgloss.Style
	CellStyle   lipgloss.Style
}

// New creates a new Table instance with sensible defaults.
func New(headers []string) *Table {
	return &Table{
		Headers:     headers,
		Alignments:  make([]Alignment, len(headers)),
		Rows:        [][]string{},
		MaxWidth:    80,
		Border:      BorderRounded,
		BorderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		HeaderStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")),
		CellStyle:   lipgloss.NewStyle(),
	}
}

// AddRow adds a row of data to the table.
func (t *Table) AddRow(row ...string) {
	t.Rows = append(t.Rows, row)
}

// SetAlignments sets column alignments.
func (t *Table) SetAlignments(aligns []Alignment) {
	t.Alignments = aligns
}

// Render formats the table into a multi-line terminal string with borders and wrapped cells.
func (t *Table) Render() string {
	numCols := len(t.Headers)
	if numCols == 0 {
		// Determine numCols from longest row if no headers
		for _, row := range t.Rows {
			if len(row) > numCols {
				numCols = len(row)
			}
		}
	}
	if numCols == 0 {
		return ""
	}

	// Normalize alignments
	alignments := make([]Alignment, numCols)
	for i := 0; i < numCols; i++ {
		if i < len(t.Alignments) && t.Alignments[i] != AlignNone {
			alignments[i] = t.Alignments[i]
		} else {
			alignments[i] = AlignLeft
		}
	}

	// Normalize rows
	normalizedRows := make([][]string, len(t.Rows))
	for rIdx, row := range t.Rows {
		nRow := make([]string, numCols)
		for cIdx := 0; cIdx < numCols; cIdx++ {
			if cIdx < len(row) {
				nRow[cIdx] = row[cIdx]
			}
		}
		normalizedRows[rIdx] = nRow
	}

	// Calculate column widths
	colWidths := t.computeColumnWidths(numCols, normalizedRows)

	var sb strings.Builder

	// Top border
	if t.Border.TopLeft != "" || t.Border.TopRight != "" || t.Border.TopMid != "" {
		sb.WriteString(t.renderBorderRow(
			t.Border.TopLeft,
			t.Border.TopMid,
			t.Border.TopRight,
			t.Border.Horizontal,
			colWidths,
		))
		sb.WriteString("\n")
	}

	// Header row
	if len(t.Headers) > 0 {
		headerCells := make([]string, numCols)
		for i := 0; i < numCols; i++ {
			if i < len(t.Headers) {
				headerCells[i] = t.Headers[i]
			}
		}
		sb.WriteString(t.renderRow(headerCells, colWidths, alignments, true))
		sb.WriteString("\n")

		// Header separator border
		if t.Border.MidLeft != "" || t.Border.MidRight != "" || t.Border.MidMid != "" || t.Border.HeaderSeparator != "" {
			sb.WriteString(t.renderBorderRow(
				t.Border.MidLeft,
				t.Border.MidMid,
				t.Border.MidRight,
				t.Border.HeaderSeparator,
				colWidths,
			))
			sb.WriteString("\n")
		}
	}

	// Data rows
	for rIdx, row := range normalizedRows {
		sb.WriteString(t.renderRow(row, colWidths, alignments, false))
		sb.WriteString("\n")

		// Optional separator between rows for markdown format
		if t.Border.TopLeft == "" && t.Border.BottomLeft == "" && t.Border.Horizontal == "-" && rIdx == -1 {
			// standard markdown
		}
	}

	// Bottom border
	if t.Border.BottomLeft != "" || t.Border.BottomRight != "" || t.Border.BottomMid != "" {
		sb.WriteString(t.renderBorderRow(
			t.Border.BottomLeft,
			t.Border.BottomMid,
			t.Border.BottomRight,
			t.Border.Horizontal,
			colWidths,
		))
	}

	return strings.TrimRight(sb.String(), "\n")
}

func (t *Table) computeColumnWidths(numCols int, rows [][]string) []int {
	naturalWidths := make([]int, numCols)
	minWordWidths := make([]int, numCols)

	// Measure headers
	for i := 0; i < numCols; i++ {
		if i < len(t.Headers) {
			w := VisualWidth(t.Headers[i])
			if w > naturalWidths[i] {
				naturalWidths[i] = w
			}
			for _, word := range strings.Fields(t.Headers[i]) {
				ww := VisualWidth(word)
				if ww > minWordWidths[i] {
					minWordWidths[i] = ww
				}
			}
		}
		if minWordWidths[i] < 3 {
			minWordWidths[i] = 3
		}
	}

	// Measure data rows
	for _, row := range rows {
		for i := 0; i < numCols; i++ {
			cell := row[i]
			w := VisualWidth(cell)
			if w > naturalWidths[i] {
				naturalWidths[i] = w
			}
			for _, word := range strings.Fields(cell) {
				ww := VisualWidth(word)
				if ww > minWordWidths[i] {
					minWordWidths[i] = ww
				}
			}
		}
	}

	// Calculate framing overhead:
	// Each cell has 1 space padding on left, 1 on right = 2 spaces * numCols
	// Vertical borders: if outer borders exist, numCols + 1 vertical characters.
	hasOuterBorders := t.Border.Vertical != "" && t.Border.TopLeft != ""
	overhead := numCols * 2 // cell padding
	if hasOuterBorders {
		overhead += (numCols + 1) * VisualWidth(t.Border.Vertical)
	} else if t.Border.Vertical != "" {
		overhead += (numCols - 1) * VisualWidth(t.Border.Vertical)
	}

	targetWidth := t.MaxWidth
	if targetWidth <= 0 {
		targetWidth = 80
	}

	availableForContent := targetWidth - overhead
	if availableForContent < numCols*3 {
		availableForContent = numCols * 3
	}

	// Total natural width
	totalNatural := 0
	for _, w := range naturalWidths {
		totalNatural += w
	}

	// If natural width fits comfortably within available space, use it!
	if totalNatural <= availableForContent {
		res := make([]int, numCols)
		copy(res, naturalWidths)
		for i := range res {
			if res[i] < 3 {
				res[i] = 3
			}
		}
		return res
	}

	// Cap minWordWidths at naturalWidths and enforce minimum of 3
	for i := 0; i < numCols; i++ {
		if minWordWidths[i] > naturalWidths[i] {
			minWordWidths[i] = naturalWidths[i]
		}
		if minWordWidths[i] < 3 {
			minWordWidths[i] = 3
		}
	}

	totalMin := 0
	for _, mw := range minWordWidths {
		totalMin += mw
	}

	res := make([]int, numCols)
	usedSpace := 0

	if availableForContent >= totalMin {
		// All columns can satisfy their minWordWidth without word-splitting!
		// 1. Allocate minWordWidth to each column
		for i := 0; i < numCols; i++ {
			res[i] = minWordWidths[i]
			usedSpace += res[i]
		}

		// 2. Distribute remaining space proportionally based on deficit (natural - min)
		remaining := availableForContent - usedSpace
		totalDeficit := 0
		for i := 0; i < numCols; i++ {
			deficit := naturalWidths[i] - minWordWidths[i]
			if deficit > 0 {
				totalDeficit += deficit
			}
		}

		if totalDeficit > 0 && remaining > 0 {
			for i := 0; i < numCols; i++ {
				deficit := naturalWidths[i] - minWordWidths[i]
				if deficit > 0 {
					extra := (deficit * remaining) / totalDeficit
					if extra > deficit {
						extra = deficit
					}
					res[i] += extra
					usedSpace += extra
				}
			}
		}

		// 3. Distribute any remaining space to columns still below natural width
		leftover := availableForContent - usedSpace
		for leftover > 0 {
			bestCol := -1
			maxDef := -1
			for i := 0; i < numCols; i++ {
				def := naturalWidths[i] - res[i]
				if def > maxDef {
					maxDef = def
					bestCol = i
				}
			}
			if bestCol == -1 || maxDef <= 0 {
				// All columns reached natural width, distribute remaining cyclically
				for i := 0; i < numCols && leftover > 0; i++ {
					res[i]++
					leftover--
				}
				break
			}
			res[bestCol]++
			leftover--
		}
	} else {
		// Terminal is severely constrained: distribute proportionally from total available
		for i := 0; i < numCols; i++ {
			w := 3
			if totalNatural > 0 {
				w = (naturalWidths[i] * availableForContent) / totalNatural
			}
			if w < 3 {
				w = 3
			}
			res[i] = w
			usedSpace += w
		}
		leftover := availableForContent - usedSpace
		for leftover > 0 {
			bestCol := -1
			maxDef := -1
			for i := 0; i < numCols; i++ {
				def := naturalWidths[i] - res[i]
				if def > maxDef {
					maxDef = def
					bestCol = i
				}
			}
			if bestCol == -1 || maxDef <= 0 {
				for i := 0; i < numCols && leftover > 0; i++ {
					res[i]++
					leftover--
				}
				break
			}
			res[bestCol]++
			leftover--
		}
	}

	return res
}

func (t *Table) renderBorderRow(left, mid, right, horiz string, colWidths []int) string {
	var sb strings.Builder
	sb.WriteString(t.BorderStyle.Render(left))

	for i, w := range colWidths {
		// Each cell has 1 leading space and 1 trailing space, so horiz repeats w + 2 times
		sb.WriteString(t.BorderStyle.Render(strings.Repeat(horiz, w+2)))
		if i < len(colWidths)-1 {
			sb.WriteString(t.BorderStyle.Render(mid))
		}
	}

	sb.WriteString(t.BorderStyle.Render(right))
	return sb.String()
}

func (t *Table) renderRow(cells []string, colWidths []int, alignments []Alignment, isHeader bool) string {
	numCols := len(colWidths)
	wrappedCells := make([][]string, numCols)
	maxLines := 1

	for i := 0; i < numCols; i++ {
		content := ""
		if i < len(cells) {
			content = cells[i]
		}
		lines := WrapWord(content, colWidths[i])
		wrappedCells[i] = lines
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	var sb strings.Builder
	for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
		if t.Border.Vertical != "" && (t.Border.TopLeft != "" || t.Border.MidLeft != "") {
			sb.WriteString(t.BorderStyle.Render(t.Border.Vertical))
		}

		for colIdx := 0; colIdx < numCols; colIdx++ {
			w := colWidths[colIdx]
			var text string
			if lineIdx < len(wrappedCells[colIdx]) {
				text = wrappedCells[colIdx][lineIdx]
			}

			// Format cell with padding and alignment
			padded := t.padCell(text, w, alignments[colIdx])

			// Apply style
			if isHeader {
				sb.WriteString(" " + t.HeaderStyle.Render(padded) + " ")
			} else {
				sb.WriteString(" " + t.CellStyle.Render(padded) + " ")
			}

			// Separator
			if colIdx < numCols-1 {
				sb.WriteString(t.BorderStyle.Render(t.Border.Vertical))
			} else if t.Border.Vertical != "" && (t.Border.TopRight != "" || t.Border.MidRight != "") {
				sb.WriteString(t.BorderStyle.Render(t.Border.Vertical))
			}
		}

		if lineIdx < maxLines-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (t *Table) padCell(text string, targetWidth int, align Alignment) string {
	curWidth := VisualWidth(text)
	if curWidth >= targetWidth {
		return text
	}

	diff := targetWidth - curWidth
	switch align {
	case AlignRight:
		return strings.Repeat(" ", diff) + text
	case AlignCenter:
		left := diff / 2
		right := diff - left
		return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
	case AlignLeft:
		fallthrough
	default:
		return text + strings.Repeat(" ", diff)
	}
}
