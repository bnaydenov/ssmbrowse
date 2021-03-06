package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	awsutils "github.com/bnaydenov/ssmbrowse/internal/pkg/awsutils"
	"github.com/bnaydenov/ssmbrowse/internal/pkg/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/thoas/go-funk"
)

//createResultTable is creating main results table
func createResultTable(ssmParams []ssm.ParameterMetadata, withData bool) *tview.Table {

	table := tview.NewTable().
		SetFixed(1, 6).
		SetSelectable(true, false)
	table.
		SetBorder(true).
		SetTitle("  SSM parameter browser ").
		SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(tcell.ColorDarkOrange)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(ssmSearchBox)
			return nil
		}
		return event
	})

	table.SetSelectionChangedFunc(func(row, column int) {

		currentRowCount := len(foundParams)
		if row == len(foundParams) {
			if nextToken != nil {

				var err error
				foundParams, nextToken, err = awsutils.SsmDescribeParameters(aws.String(ssmSearchBox.GetText()), nextToken, foundParams)
				if err != nil {
					errorModal.SetText(fmt.Sprintf("%s", err.Error()))
					pages.SwitchToPage("error")
				}

				// if there is no more item remove  centerFooterItem text
				if nextToken == nil {
					updateFooterItem(centerFooterItem, "", tview.AlignCenter, tcell.ColorWhite)
				}

				ssmTable = createResultTable(foundParams, true)
				mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
				ssmTable.Select(currentRowCount, 0)
				app.SetFocus(ssmTable)
			}
		}
	})

	table.SetSelectedFunc(func(row int, column int) {

		ssmParam := table.GetCell(row, column).GetReference().(ssm.ParameterMetadata)
		ssmParamDetailsForm.SetTitle(fmt.Sprintf(" %s ", *ssmParam.Name)).SetTitleAlign(tview.AlignCenter)

		secureSsmParam, err := awsutils.GetParameter(*ssmParam.Name)
		if err != nil {
			errorModal.SetText(fmt.Sprintf("%s", err.Error()))
			pages.SwitchToPage("error")
		}

		ssmParamDetailsForm.AddInputField("Value:", *secureSsmParam.Parameter.Value, 100, nil, nil)
		ssmParamDetailsForm.AddInputField("Version:", fmt.Sprintf("%d", *secureSsmParam.Parameter.Version), 100, nil, nil)
		ssmParamDetailsForm.AddInputField("ARN:", *secureSsmParam.Parameter.ARN, 100, nil, nil)
		ssmParamDetailsForm.AddInputField("Last Modified Date:", secureSsmParam.Parameter.LastModifiedDate.Local().String(), 100, nil, nil)

		ssmParamDetailsForm.SetFocus(4)
		pages.SwitchToPage("ssmParam")

	})

	expansions := []int{3, 1, 1, 1}
	alignment := []int{ui.L, ui.L, ui.L, ui.L, ui.L, ui.L}

	headers := []string{"Name", "Type", "Version", "Last modified"}
	ui.AddTableData(table, 0, [][]string{headers}, alignment, expansions, tcell.ColorYellow, false)
	if withData {
		data := funk.Map(ssmParams, func(param ssm.ParameterMetadata) []string {
			return []string{
				aws.StringValue(param.Name),
				aws.StringValue(param.Type),
				fmt.Sprintf("%d", *param.Version),
				param.LastModifiedDate.Local().String(),
			}
		}).([][]string)

		ui.AddTableData(table, 1, data, alignment, expansions, tcell.ColorWhite, true)
		// Add a reference to the data to column 0 in each row for easy access later on
		for row, ssmParam := range ssmParams {
			table.GetCell(row+1, 0).SetReference(ssmParam)
		}

	}

	return table
}
