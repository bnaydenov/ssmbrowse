package cmd

import "github.com/rivo/tview"

func updateFooterItem(item *tview.TextView, text string, alignment int) {
    item.SetText(text).SetTextAlign(alignment)
}