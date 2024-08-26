package terminal

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

const (
	topLeft     = "┌"
	topRight    = "┐"
	bottomLeft  = "└"
	bottomRight = "┘"
	horizontal  = "─"
	vertical    = "│"
	leftT       = "├"
	rightT      = "┤"
	topT        = "┬"
	bottomT     = "┴"
	cross       = "┼"
)

// Alignment represents the text alignment within a table cell.
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

// Column represents a single column in a table.
type Column struct {
	Header    string    // The header text for the column
	Field     string    // The field name to extract data from
	MinWidth  int       // The minimum width of the column
	MaxWidth  *int      // The maximum width of the column
	Wrap      bool      // Whether to wrap text in this column
	Alignment Alignment // The text alignment for this column
}

// Table represents a table to be rendered in the terminal.
type Table struct {
	terminal   *Terminal                // The terminal to render the table in
	Columns    []Column                 // The columns of the table
	Data       []map[string]interface{} // The data to be displayed in the table
	Percentage float64                  // The percentage of terminal width to use
}

// NewTable creates a new Table instance.
func (t *Terminal) NewTable(percentage float64) *Table {
	return &Table{
		terminal:   t,
		Percentage: percentage,
	}
}

// AddColumn adds a new column to the table.
func (t *Table) AddColumn(header, field string, minWidth int, maxWidth *int, wrap bool, alignment Alignment) {
	t.Columns = append(t.Columns, Column{
		Header:    header,
		Field:     field,
		MinWidth:  minWidth,
		MaxWidth:  maxWidth,
		Wrap:      wrap,
		Alignment: alignment,
	})
}

// AddRow adds a new row of data to the table.
func (t *Table) AddRow(row map[string]interface{}) {
	t.Data = append(t.Data, row)
}

// Render draws the table in the terminal.
func (t *Table) Render() {
	// Calculate table width based on percentage
	tableWidth := int(float64(t.terminal.width) * t.Percentage / 100)

	// Calculate column widths
	columnWidths := t.calculateColumnWidths(tableWidth)

	t.printBorder(columnWidths, topLeft, topRight, topT, horizontal)

	// Print header
	t.printHeader(columnWidths)

	t.printBorder(columnWidths, leftT, rightT, cross, horizontal)

	// Print rows
	for i, row := range t.Data {
		t.printRow(row, columnWidths)
		if i < len(t.Data)-1 {
			t.printBorder(columnWidths, leftT, rightT, cross, horizontal)
		}
	}
	t.printBorder(columnWidths, bottomLeft, bottomRight, bottomT, horizontal)
}

// calculateColumnWidths determines the width of each column based on the table's total width.
func (t *Table) calculateColumnWidths(tableWidth int) []int {
	// Account for all separators, including start and end
	remainingWidth := tableWidth - (len(t.Columns) * 3)
	columnWidths := make([]int, len(t.Columns))

	// First pass: allocate minimum widths
	for i, col := range t.Columns {
		columnWidths[i] = col.MinWidth
		remainingWidth -= col.MinWidth
	}

	// Second pass: distribute remaining width, respecting MaxWidth if set
	if remainingWidth > 0 {
		for i := range t.Columns {
			if t.Columns[i].MaxWidth != nil && columnWidths[i] < *t.Columns[i].MaxWidth {
				additionalWidth := min(remainingWidth, *t.Columns[i].MaxWidth-columnWidths[i])
				columnWidths[i] += additionalWidth
				remainingWidth -= additionalWidth
			}
		}
	}

	// Third pass: distribute any remaining width to columns without MaxWidth
	if remainingWidth > 0 {
		columnsWithoutMax := 0
		for _, col := range t.Columns {
			if col.MaxWidth == nil {
				columnsWithoutMax++
			}
		}

		if columnsWithoutMax > 0 {
			widthPerColumn := remainingWidth / columnsWithoutMax
			for i, col := range t.Columns {
				if col.MaxWidth == nil {
					additionalWidth := min(widthPerColumn, remainingWidth)
					columnWidths[i] += additionalWidth
					remainingWidth -= additionalWidth
				}
			}
		}
	}

	return columnWidths
}

// printBorder prints a horizontal border of the table.
func (t *Table) printBorder(columnWidths []int, left, right, separator, fill string) {
	fmt.Fprint(t.terminal.writer, left)
	for i, width := range columnWidths {
		fmt.Fprint(t.terminal.writer, strings.Repeat(fill, width+2))
		if i < len(columnWidths)-1 {
			fmt.Fprint(t.terminal.writer, separator)
		}
	}
	fmt.Fprintln(t.terminal.writer, right)
}

// printHeader prints the header row of the table.
func (t *Table) printHeader(columnWidths []int) {
	fmt.Fprint(t.terminal.writer, vertical)
	for i, col := range t.Columns {
		fmt.Fprint(t.terminal.writer, " "+alignText(col.Header, columnWidths[i], col.Alignment)+" ")
		fmt.Fprint(t.terminal.writer, vertical)
	}
	fmt.Fprintln(t.terminal.writer)
}

// printRow prints a single row of the table.
func (t *Table) printRow(row map[string]interface{}, columnWidths []int) {
	maxLines := 1
	cellContents := make([][]string, len(t.Columns))

	for i, col := range t.Columns {
		value := fmt.Sprintf("%v", row[col.Field])
		if col.Wrap {
			cellContents[i] = wrapText(value, columnWidths[i])
			if len(cellContents[i]) > maxLines {
				maxLines = len(cellContents[i])
			}
		} else {
			cellContents[i] = []string{truncateString(value, columnWidths[i])}
		}
	}

	for line := 0; line < maxLines; line++ {
		fmt.Fprint(t.terminal.writer, vertical)
		for i, content := range cellContents {
			if line < len(content) {
				fmt.Fprint(t.terminal.writer, " "+alignText(content[line], columnWidths[i], t.Columns[i].Alignment)+" ")
			} else {
				fmt.Fprint(t.terminal.writer, " "+strings.Repeat(" ", columnWidths[i])+" ")
			}
			fmt.Fprint(t.terminal.writer, vertical)
		}
		fmt.Fprintln(t.terminal.writer)
	}
}

// alignText aligns the given text within the specified width according to the alignment.
func alignText(s string, width int, alignment Alignment) string {
	sWidth := runewidth.StringWidth(s)
	if sWidth >= width {
		return runewidth.Truncate(s, width, "")
	}

	switch alignment {
	case AlignLeft:
		return runewidth.FillRight(s, width)
	case AlignRight:
		return runewidth.FillLeft(s, width)
	case AlignCenter:
		leftPad := (width - sWidth) / 2
		rightPad := width - sWidth - leftPad
		return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
	default:
		return runewidth.FillRight(s, width)
	}
}

// wrapText wraps the given text to fit within the specified width.
func wrapText(s string, width int) []string {
	var lines []string
	var currentLine string
	currentLineWidth := 0

	words := strings.Fields(s)
	for _, word := range words {
		wordWidth := runewidth.StringWidth(word)

		if currentLineWidth+wordWidth+1 > width {
			if currentLine != "" {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = ""
				currentLineWidth = 0
			}

			if wordWidth > width {
				for len(word) > 0 {
					truncated := runewidth.Truncate(word, width, "")
					lines = append(lines, truncated)
					word = strings.TrimPrefix(word, truncated)
				}
			} else {
				currentLine = word
				currentLineWidth = wordWidth
			}
		} else {
			if currentLine != "" {
				currentLine += " "
				currentLineWidth++
			}
			currentLine += word
			currentLineWidth += wordWidth
		}
	}

	if currentLine != "" {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return lines
}

// truncateString truncates the given string to fit within the specified width, adding an ellipsis if necessary.
func truncateString(s string, width int) string {
	if runewidth.StringWidth(s) <= width {
		return s
	}
	return runewidth.Truncate(s, width-4, "....")
}
