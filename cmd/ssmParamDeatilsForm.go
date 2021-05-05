package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//SSM Param form
func createSsmParamDetailsForm() *tview.Form {

	ssmParamDetailsForm = tview.NewForm().
		SetFieldBackgroundColor(tcell.ColorDarkOrange).SetFieldTextColor(tcell.ColorBlack).
		SetButtonsAlign(tview.AlignCenter).
		AddButton("OK", func() {
			ssmParamDetailsForm.Clear(false)
			pages.SwitchToPage("main")
			app.SetFocus(ssmTable)
		})

	ssmParamDetailsForm.SetBorder(true).SetTitle("Set ssm parm name as title").SetTitleAlign(tview.AlignLeft)

	return ssmParamDetailsForm
}
