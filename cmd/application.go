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
	foundParams     []ssm.ParameterMetadata
	startToken      *string
	notFoundModal   *tview.Modal
	ssmParamForm *tview.Form
	nextToken *string
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
	
    //crete empty result ssmTable with header only
	ssmTable = createResultTable(foundParams, false)
	mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
	
	mainGrid.AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	pages.AddPage("main", mainGrid, true, true)

	//Error page
	pages.AddPage("error", notFoundModal, true, false)
	


   //SSM Param form
   ssmParamForm = tview.NewForm().
    SetFieldBackgroundColor(tcell.ColorDarkOrange).SetFieldTextColor(tcell.ColorBlack).
    SetButtonsAlign(tview.AlignCenter).
    AddButton("OK", func() {
			ssmParamForm.Clear(false)
			pages.SwitchToPage("main")
			app.SetFocus(ssmTable)
		})
	
	ssmParamForm.SetBorder(true).SetTitle("Set ssm parm name as title").SetTitleAlign(tview.AlignLeft)
	pages.AddPage("ssmParam", ssmParamForm, true, false)


	// pages.SetBorderPadding(0, 0, 1, 1).SetBorderColor(tcell.ColorDarkOrange)
	

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		if ssmParamForm.HasFocus() {
            if event.Key() == tcell.KeyRune {
				key := event.Rune()
				if key == 'c' {
					// fmt.Println("YYYYYY")
				}
			}
		}
		
		
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
			
			var err error
			foundParams, nextToken, err = awsutils.SsmDescribeParameters(aws.String(ssmSearchPrefix.GetText()), startToken, foundParams)
			if err != nil {
					notFoundModal.SetText(fmt.Sprintf("%s", err.Error()))
					pages.SwitchToPage("error")
					return
		    }
	
			// show error is not ssm params found with provided prefix
			if len(foundParams) == 0 {
				notFoundModal.SetText(fmt.Sprintf("Can't find SSM params with preffix: %s", ssmSearchPrefix.GetText()))
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

//createNotFoundModal is function which creates modal error box if no ssm param found
func createNotFoundModal() *tview.Modal {
	modal := tview.NewModal().
	AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				ssmSearchPrefix.SetText("/")
                
				//crete empty result ssmTable with header only
				ssmTable = createResultTable(foundParams, false)
	            mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
				
				pages.SwitchToPage("main")
			}
		})
	return modal
}

//createResultTable is creating main results table
func createResultTable(ssmParams []ssm.ParameterMetadata, withData bool) *tview.Table {

	table := tview.NewTable().
		SetFixed(1, 6).
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
	
	table.SetSelectionChangedFunc(func(row, column int) {
		
		currentRowCount := len(foundParams)
		if row == len(foundParams) {
            if nextToken != nil {
                
				var err error
				foundParams, nextToken,  err = awsutils.SsmDescribeParameters(aws.String(ssmSearchPrefix.GetText()), startToken, foundParams)
				if err != nil {
						notFoundModal.SetText(fmt.Sprintf("%s", err.Error()))
						pages.SwitchToPage("error")
				}
				
				ssmTable = createResultTable(foundParams, true)
			    mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
			    ssmTable.Select(currentRowCount,0)
			    app.SetFocus(ssmTable)
			}
		}
	})

	table.SetSelectedFunc(func(row int, column int) {
		
		ssmParam := table.GetCell(row, column).GetReference().(ssm.ParameterMetadata)
		ssmParamForm.SetTitle(*ssmParam.Name).SetTitleAlign(tview.AlignCenter)
		
		// if *ssmParam.Type =="SecureString" {
			secureSsmParam := awsutils.GetParameter(*ssmParam.Name)
			ssmParamForm.AddInputField("Value:", *secureSsmParam.Parameter.Value ,100, nil, nil)
		    ssmParamForm.AddInputField("Version:", fmt.Sprintf("%d",*secureSsmParam.Parameter.Version), 100, nil, nil)
		    ssmParamForm.AddInputField("ARN:", *secureSsmParam.Parameter.ARN, 100, nil, nil)
		    ssmParamForm.AddInputField("Last Modified Date:", secureSsmParam.Parameter.LastModifiedDate.Local().String() , 100, nil, nil)
		// } else  {
        //     ssmParamForm.AddInputField("Value:", *ssmParam.Value,100, nil, nil)
		//     ssmParamForm.AddInputField("Version:", fmt.Sprintf("%d",*ssmParam.Version), 100, nil, nil)
		//     ssmParamForm.AddInputField("ARN:", *ssmParam.ARN, 100, nil, nil)
		//     ssmParamForm.AddInputField("Last Modified Date:", ssmParam.LastModifiedDate.Local().String() , 100, nil, nil)
		// }
		
		
		ssmParamForm.SetFocus(4)
		pages.SwitchToPage("ssmParam")
		
		// fmt.Println(*ssmParam.Name)
		// table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		// table.SetSelectable(true, false)
	})

	// table.SetSelectionChangedFunc(func(row, column int) {
	// 	if row == 10 {
	// 		SetBorderColor(tcell.ColorDarkRed)
	// 	}
	// })

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
