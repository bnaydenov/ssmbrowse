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

// var tviewApp *tview.Application
var ssmParamTable *tview.Table

var (
	app             *tview.Application
	ssmSearchPrefix *tview.InputField
	pages           *tview.Pages
	ssmTable        *tview.Table
	mainGrid        *tview.Grid
	foundParams     []ssm.Parameter
	startToken      *string
	notFoundModal   *tview.Modal
)

func Entrypoint() {
	fmt.Println("Loading information about your AWS SSM params...")

	// params,  _  = awsutils.GetParemters(aws.String("/qldwflkjfds/"), startToken, params)

	// fmt.Println(len(params))
	// for _, p := range params {
	// 	fmt.Println(*p.Name)
	// }
	// os.Exit(0)
	// for nextToken != nil {
	// 	params, nextToken = aws.GetParemters([]string{"/"}, nextToken, params)
	// }

	// for _, p := range params {
	// 	fmt.Println(*p.Name)
	// }
	// fmt.Println(len(params))

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	app = tview.NewApplication()
	pages = tview.NewPages()

	// main page
	ssmSearchPrefix = createSsmSearchPrefix()

	// paramFilter.SetBorderColor(tcell.ColorDarkOrange).SetBorderPadding(0, 0, 1, 1)

	mainGrid = tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0).
		SetBorders(true)

	mainGrid.AddItem(ssmSearchPrefix, 0, 0, 1, 3, 0, 0, true)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	//  grid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
	// Layout for screens wider than 100 cells.
	// grid.AddItem(main, 1, 1, 1, 1, 0, 100, false)
	mainGrid.AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	pages.AddPage("main", mainGrid, true, true)

	//Error page
	notFoundModal = createNotFoundModal()
	pages.AddPage("error", notFoundModal, true, false)
	
	// pages.SetBorderPadding(0, 0, 1, 1).SetBorderColor(tcell.ColorDarkOrange)
	

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		switch event.Key() {
		case tcell.KeyEsc:
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

//createSsmSearchPrefix creates SsmSearchPrefix
func createSsmSearchPrefix() *tview.InputField {

	ssmSearchPrefix := tview.NewInputField().SetLabel("Enter a param prefix: ").SetFieldBackgroundColor(tcell.ColorDarkOrange)
	ssmSearchPrefix.SetText("/")
	ssmSearchPrefix.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			foundParams = nil
			startToken = nil
			if ssmTable != nil {
				ui.TruncTableRows(ssmTable, ssmTable.GetRowCount())
				mainGrid.RemoveItem(ssmTable)
			}
			foundParams, _ = awsutils.GetParemters(aws.String(ssmSearchPrefix.GetText()), startToken, foundParams)
			// show error is not ssm params found with provided prefix
			if len(foundParams) == 0 {
				notFoundModal.SetText(fmt.Sprintf("Can't find SSM params with preffix: %s", ssmSearchPrefix.GetText()))
				pages.SwitchToPage("error")
				return
			}
			ssmTable = createResultTable(foundParams)
			mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
			app.SetFocus(ssmTable)
		}
	})
	return ssmSearchPrefix
}

//createNotFoundModal is function which creates modal error box if no ssm param found
func createNotFoundModal() *tview.Modal {
	modal := tview.NewModal().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				ssmSearchPrefix.SetText("/")
				pages.SwitchToPage("main")
			}
		})
	return modal
}

//createResultTable is creating main results table
func createResultTable(ssmParams []ssm.Parameter) *tview.Table {

	table := tview.NewTable().
		SetFixed(4, 6).
		SetSelectable(true, false)
	table.
		SetBorder(true).
		SetTitle(" ⌛ SSM parameter browser...").
		SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(tcell.ColorDarkOrange)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(ssmSearchPrefix)
			return nil
		}
		return event
	})

	// table.SetSelectionChangedFunc(func(row, column int) {
	// 	if row == 10 {
	// 		SetBorderColor(tcell.ColorDarkRed)
	// 	}
	// })

	expansions := []int{3, 1, 1, 1, 1, 1}
	alignment := []int{ui.L, ui.L, ui.L, ui.L, ui.L, ui.L}

	headers := []string{"Name", "Tier", "Type", "Description", "Version", "Last modified"}
	ui.AddTableData(table, 0, [][]string{headers}, alignment, expansions, tcell.ColorYellow, false)

	data := funk.Map(foundParams, func(param ssm.Parameter) []string {
		return []string{
			aws.StringValue(param.Name),
			aws.StringValue(param.Type),
			aws.StringValue(param.Type),
			aws.StringValue(param.DataType),
			fmt.Sprintf("%d", *param.Version),
			param.LastModifiedDate.Format("01-01-2021 00:00:00"),
		}
	}).([][]string)

	ui.AddTableData(table, 1, data, alignment, expansions, tcell.ColorWhite, true)
	return table
}
