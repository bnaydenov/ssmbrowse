package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//TableInfo is Contains information for configuring a tview table
type TableInfo struct {
	Table      *tview.Table
	Alignment  []int
	Expansions []int
	Selectable bool
}

//AddTableConfigData is .......
func AddTableConfigData(table *TableInfo, startRow int, data [][]string, color tcell.Color) {
	AddTableData(table.Table, startRow, data, table.Alignment, table.Expansions, color, table.Selectable)
}

// L - Left align
const L = tview.AlignLeft

// C - Center align
const C = tview.AlignCenter

// R - Right align
const R = tview.AlignRight

// AddTableData is ......
func AddTableData(table *tview.Table, startRow int, data [][]string, alignment []int, expansions []int, color tcell.Color, selectable bool) {

	if len(expansions) < len(data[0]) {
		log.Printf("warning: expansions (%d) not aligned with data [%d][%d]\n", len(expansions), len(data), len(data[0]))
	}
	if len(alignment) < len(data[0]) {
		log.Printf("warning: alignment (%d) not aligned with data [%d][%d]\n", len(alignment), len(data), len(data[0]))
	}

	for row, line := range data {
		for col, text := range line {
			cell := tview.NewTableCell(text).
				SetAlign(alignment[col]).
				SetExpansion(expansions[col]).
				SetTextColor(color).
				SetSelectable(selectable)
			table.SetCell(row+startRow, col, cell)
		}
	}
}

//SetColumnStyle style for column
func SetColumnStyle(table *tview.Table, col int, rowStart int, style tcell.Style) {
	for row := rowStart; row <= table.GetRowCount()-1; row++ {
		table.GetCell(row, col).SetStyle(style)
	}
}

// TruncTableRows - truncate  all row from table
func TruncTableRows(table *tview.Table, maxRows int) {
	for table.GetRowCount() > maxRows {
		table.RemoveRow(maxRows)
	}
}
