package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/bnaydenov/ssmbrowse/internal/pkg/aws"
	"github.com/bnaydenov/ssmbrowse/internal/pkg/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/thoas/go-funk"
)

// var tviewApp *tview.Application
var ssmParamTable *tview.Table

var (
    app            *tview.Application
	paramFilter      *tview.InputField
	paramFilterFlexBox     *tview.Flex
	ssmParamsFlexBox *tview.Flex
	mainFlexBox    *tview.Flex
	pages    *tview.Pages
	ssmTable *tview.Table
	
)

func Entrypoint () {
	fmt.Println("Loading information about your AWS SSM params...")
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	// box := tview.NewBox().SetBorder(true).SetTitle("ssm browser")
	app = tview.NewApplication()

	pages = tview.NewPages()

	ssmTable = buildClusterTable()

	paramFilter = tview.NewInputField().SetLabel("Enter a param prefix: ").SetFieldBackgroundColor(tcell.ColorDarkOrange)
	paramFilter.SetDoneFunc(func(key tcell.Key) {
		// if key == tcell.KeyEnter {
			// paramFilter.SetText("XXXXX")
		// }
		app.SetFocus(ssmTable)
	})

	
	
	// paramFilter.SetBorderColor(tcell.ColorDarkOrange).SetBorderPadding(0, 0, 1, 1)
	
	grid := tview.NewGrid().
	SetRows(1,0,1).
	SetColumns(0).
	SetBorders(true)

	grid.AddItem(paramFilter, 0, 0, 1, 3, 0, 0, true)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
     grid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
    // Layout for screens wider than 100 cells.
    // grid.AddItem(main, 1, 1, 1, 1, 0, 100, false)
    grid.AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)
	
	
	pages.AddPage("main", grid, true, true)
	pages.SetBorderPadding(0, 0, 1, 1).SetBorderColor(tcell.ColorDarkOrange)

	
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

	if err := app.SetRoot(pages, true).SetFocus(paramFilter).Run(); err != nil {
		panic(err)
	}

}


func buildClusterTable() *tview.Table {

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
			// Exit the application
			app.SetFocus(paramFilter)
			return nil
		}
		return event 
	})


	// table.SetSelectionChangedFunc(func(row, column int) {
	// 	if row == 10 {	
	// 		SetBorderColor(tcell.ColorDarkRed)
	// 	}
	// })

	
	expansions := []int{3, 1, 1, 1,1,1}
	alignment := []int{ui.L, ui.L, ui.L, ui.L,ui.L,ui.L}

	headers := []string{"Name", "Tier", "Type", "Description", "Version","Last modified"}
	ui.AddTableData(table, 0, [][]string{headers}, alignment, expansions, tcell.ColorYellow, false)
	
	var startToken  *string
	var params []ssm.ParameterMetadata
	
	params,  _  = aws.GetParemters([]string{"/"}, startToken, params)
    
	// for nextToken != nil {
	// 	params, nextToken = aws.GetParemters([]string{"/"}, nextToken, params)
	// }
	
	// for _, p := range params {
	// 	fmt.Println(*p.Name)
	// }
	// fmt.Println(len(params))
	
	
	data := funk.Map(params, func(param ssm.ParameterMetadata) []string {
		return []string{
			derefString(param.Name),
			derefString(param.Tier),
			derefString(param.Type),
			derefString(param.Description),
			fmt.Sprintf("%d", *param.Version),
			param.LastModifiedDate.Format("01-01-2021 00:00:00"),
		}
	}).([][]string)
	
	ui.AddTableData(table, 1, data, alignment, expansions, tcell.ColorWhite, true)
	return table
}

func derefString(s *string) string {
    if s != nil {
        return *s
    }
    return ""
}