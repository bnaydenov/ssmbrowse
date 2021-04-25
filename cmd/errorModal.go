package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//createErrorModal is function which creates modal error box if no ssm param found
func createErrorModal() *tview.Modal {
	modal := tview.NewModal().
	AddButtons([]string{"OK"}).
	SetBackgroundColor(tcell.ColorDarkOrange).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				//crete empty result ssmTable with header only
				ssmTable = createResultTable(foundParams, false)
	            mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
				pages.SwitchToPage("main")
			}
		})
	return modal
}
