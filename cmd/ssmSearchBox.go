package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/bnaydenov/ssmbrowse/internal/pkg/awsutils"
	"github.com/bnaydenov/ssmbrowse/internal/pkg/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//createSsmSearchBox creates SsmSearchPrefix
func createSsmSearchBox() *tview.InputField {

	ssmSearchPrefix := tview.NewInputField().SetLabel("Enter a param prefix: ").SetFieldBackgroundColor(tcell.ColorDarkOrange)
	// ssmSearchPrefix.SetText("/")

	ssmSearchPrefix.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			foundParams = nil
			startToken = nil
			if ssmTable != nil {
				ui.TruncTableRows(ssmTable, ssmTable.GetRowCount())
				mainGrid.RemoveItem(ssmTable)
			}
			
			var err error
			foundParams, nextToken, err = awsutils.SsmDescribeParameters(aws.String(ssmSearchPrefix.GetText()), startToken, foundParams)
			if err != nil {
					errorModal.SetText(fmt.Sprintf("%s", err.Error()))
					pages.SwitchToPage("error")
					return
		    }
	
			// show error is not ssm params found with provided prefix
			if len(foundParams) == 0 {
				errorModal.SetText(fmt.Sprintf("Can't find SSM params containing '%s'", ssmSearchPrefix.GetText()))
				pages.SwitchToPage("error")
				return
			}
			ssmTable = createResultTable(foundParams, true)
			mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
			app.SetFocus(ssmTable)
		}

		if key == tcell.KeyTAB {
			app.SetFocus(ssmTable)
		}
	})
	return ssmSearchPrefix
}