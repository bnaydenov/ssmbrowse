package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// var tviewApp *tview.Application
var ssmParamTable *tview.Table

var (
	app             *tview.Application
	ssmSearchBox *tview.InputField
	pages           *tview.Pages
	ssmTable        *tview.Table
	mainGrid        *tview.Grid
	foundParams     []ssm.ParameterMetadata
	startToken      *string
	errorModal   *tview.Modal
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
	ssmSearchBox = createSsmSearchBox()

	// paramFilter.SetBorderColor(tcell.ColorDarkOrange).SetBorderPadding(0, 0, 1, 1)

	mainGrid = tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0).
		SetBorders(true)

	mainGrid.AddItem(ssmSearchBox, 0, 0, 1, 3, 0, 0, true)
	
    //crete empty result ssmTable with header only
	ssmTable = createResultTable(foundParams, false)
	mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)
	
	mainGrid.AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	pages.AddPage("main", mainGrid, true, true)

	//Error page
	errorModal = createErrorModal()
	pages.AddPage("error", errorModal, true, false)
	


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
