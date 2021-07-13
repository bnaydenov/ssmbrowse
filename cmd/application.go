package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	awsutils "github.com/bnaydenov/ssmbrowse/internal/pkg/awsutils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

var (
	app                    *tview.Application
	ssmParamTable          *tview.Table
	ssmSearchBox           *tview.InputField
	pages                  *tview.Pages
	ssmTable               *tview.Table
	mainGrid, ssmParamGrid *tview.Grid
	foundParams            []ssm.ParameterMetadata
	errorModal             *tview.Modal
	ssmParamDetailsForm    *tview.Form
	nextToken              *string
	leftFooterItem         *tview.TextView
	centerFooterItem       *tview.TextView
	rightFooterItem        *tview.TextView
	accountID, awsRegion   *string
	version                string
)

//Entrypoint is
func Entrypoint(buildData map[string]interface{}) {

	version = buildData["version"].(string)

	app = tview.NewApplication()
	pages = tview.NewPages()

	ssmSearchBox = createSsmSearchBox()

	mainGrid = tview.NewGrid().
		SetRows(1, 0, 2).
		SetColumns(0)

	mainGrid.AddItem(ssmSearchBox, 0, 0, 1, 3, 0, 0, true).SetBorder(true).SetBorderColor(tcell.ColorDarkOrange)

	//crete empty result ssmTable with header only
	ssmTable = createResultTable(foundParams, false)
	mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)

	leftFooterItem = tview.NewTextView()
	mainGrid.AddItem(leftFooterItem, 2, 0, 1, 1, 0, 0, false)
	updateFooterItem(leftFooterItem, "ESC/CTRL+C=Exit | TAB=Switch focus | ENTER=See details \nC=Copy value to clipboard | X=Copy name to clipboard", tview.AlignLeft, tcell.ColorWhite)

	centerFooterItem = tview.NewTextView()
	mainGrid.AddItem(centerFooterItem, 2, 1, 1, 1, 0, 0, false)
	// updateFooterItem(centerFooterItem,"", tview.AlignCenter)

	rightFooterItem = tview.NewTextView()
	mainGrid.AddItem(rightFooterItem, 2, 2, 1, 1, 0, 0, false)
	updateFooterItem(rightFooterItem, version, tview.AlignRight, tcell.ColorBlue)

	pages.AddPage("main", mainGrid, true, true)

	//Error modal
	errorModal = createErrorModal()
	pages.AddPage("error", errorModal, true, false)

	//SSM Param details form
	ssmParamDetailsForm = createSsmParamDetailsForm()

	ssmParamGrid = tview.NewGrid().SetRows(0, 0, 0, 0, 1).SetColumns(0, 0, 0, 0, 0)
	ssmParamGrid.AddItem(ssmParamDetailsForm, 1, 1, 2, 3, 0, 0, true)

	ssmParamFooterItemLeft := tview.NewTextView().SetText("").SetTextAlign(tview.AlignLeft)
	ssmParamGrid.AddItem(ssmParamFooterItemLeft, 4, 0, 1, 1, 0, 0, false)
	pages.AddPage("ssmParam", ssmParamGrid, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		if ssmTable.HasFocus() {
			if event.Key() == tcell.KeyRune {
				key := event.Rune()
				if key == 'c' || key == 'x' || key == 'C' || key == 'X' {
					
					selectedRow, selectedCol := ssmTable.GetSelection()
                    selectedSSMParam := ssmTable.GetCell(selectedRow, selectedCol).GetReference().(ssm.ParameterMetadata)

                    switch string(key) {
					// Copy SSM value to clipboard
					case "c", "C":
						selectedSSMParamDetails, err := awsutils.GetParameter(*selectedSSMParam.Name)
                        if err != nil {
                            errorModal.SetText(fmt.Sprintf("%s", err.Error()))
                            pages.SwitchToPage("error")
                        }
						// write/read text format data of the clipboard, and
                        // the byte buffer regarding the text are UTF8 encoded.
						clipboard.Write(clipboard.FmtText,[]byte(*selectedSSMParamDetails.Parameter.Value))
						updateFooterItem(centerFooterItem, fmt.Sprintf("Value of '%s' is copied to clipboard.",*selectedSSMParamDetails.Parameter.Name ), tview.AlignCenter, tcell.ColorDarkOrange)

					// Copy SSM name to clipboard
					case "x", "X":
						// write/read text format data of the clipboard, and
                        // the byte buffer regarding the text are UTF8 encoded.
						clipboard.Write(clipboard.FmtText,[]byte(*selectedSSMParam.Name))
						updateFooterItem(centerFooterItem, fmt.Sprintf("SSM name '%s' is copied to clipboard.",*selectedSSMParam.Name), tview.AlignCenter, tcell.ColorDarkOrange)
					}
				}
			}
		}
		
		switch event.Key() {
		case tcell.KeyEsc:
			if  ssmParamDetailsForm.HasFocus() {
                ssmParamDetailsForm.Clear(false)
                pages.SwitchToPage("main")
                app.SetFocus(ssmTable)
				return nil
			}
			// Exit the application
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}