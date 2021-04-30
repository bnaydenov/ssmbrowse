package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func updateFooterItem(item *tview.TextView, text string, alignment int, color tcell.Color) {
    item.SetText(text).SetTextAlign(alignment).SetTextColor(color)
}