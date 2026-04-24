package output

import (
	"fmt"
	"os"
	"strings"
)

func Success(msg string) { fmt.Printf("  OK  %s\n", msg) }
func Info(msg string)    { fmt.Printf(" INFO %s\n", msg) }
func Warn(msg string)    { fmt.Printf(" WARN %s\n", msg) }
func Error(msg string)   { fmt.Fprintf(os.Stderr, "ERROR %s\n", msg) }
func Fatal(msg string) {
	Error(msg)
	os.Exit(1)
}

func Table(headers []string, rows [][]string) {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var b strings.Builder
	printRow := func(cols []string) {
		for i, col := range cols {
			if i < len(widths) {
				b.WriteString(pad(col, widths[i]))
			} else {
				b.WriteString(col)
			}
			if i < len(cols)-1 {
				b.WriteString("  ")
			}
		}
		b.WriteByte('\n')
	}

	printRow(headers)
	sep := make([]string, len(headers))
	for i, w := range widths {
		sep[i] = strings.Repeat("-", w)
	}
	printRow(sep)
	for _, row := range rows {
		printRow(row)
	}
	fmt.Print(b.String())
}

func pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
