package ncurses

import (
	"fmt"
	"github.com/gbin/goncurses"
	"strings"
)

type Table struct {
	ColWidths []int
}

func NewTable(colWidths []int) *Table {
	return &Table{
		ColWidths: colWidths,
	}
}

func (t *Table) Render(win *goncurses.Window, x int, y int, vals [][]string) {
	for i := range vals {
		row := vals[i]
		truncated := make([]string, len(row))
		for j := range row {
			col := row[j]
			colWidth := t.ColWidths[j]
			if len(col) >= colWidth {
				truncated[j] = col[:colWidth]
			} else {
				truncated[j] = fmt.Sprintf("%*s", colWidth, col)
			}
		}

		rowStr := strings.Join(truncated, " | ")
		win.MovePrint(y+i, x, rowStr)
	}
}
