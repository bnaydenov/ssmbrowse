package cmd

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app                  *tview.Application
	ssmParamTable        *tview.Table
	ssmSearchBox         *tview.InputField
	pages                *tview.Pages
	ssmTable             *tview.Table
	mainGrid, ssmParamGrid             *tview.Grid
	foundParams          []ssm.ParameterMetadata
	errorModal           *tview.Modal
	ssmParamDetailsForm  *tview.Form
	nextToken            *string
	leftFooterItem       *tview.TextView
	centerFooterItem     *tview.TextView
	rightFooterItem      *tview.TextView
	accountID, awsRegion *string
	version              string
)

//Entrypoint is
func Entrypoint(buildData map[string]interface{}) {
	// fmt.Println("Loading information about your AWS SSM params...")

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

	// newPrimitive := func(text string) tview.Primitive {
	// 	return tview.NewTextView().
	// 		SetTextAlign(tview.AlignCenter).
	// 		SetText(text)
	// }
	version = buildData["version"].(string)

	app = tview.NewApplication()
	pages = tview.NewPages()

	ssmSearchBox = createSsmSearchBox()

	mainGrid = tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0)

	mainGrid.AddItem(ssmSearchBox, 0, 0, 1, 3, 0, 0, true).SetBorder(true).SetBorderColor(tcell.ColorDarkOrange)
	// mainGrid.SetBorder(true).SetBorderColor(tcell.ColorDarkOrange)
	//crete empty result ssmTable with header only
	ssmTable = createResultTable(foundParams, false)
	mainGrid.AddItem(ssmTable, 1, 0, 1, 3, 0, 0, false)

	leftFooterItem = tview.NewTextView()
	mainGrid.AddItem(leftFooterItem, 2, 0, 1, 1, 0, 0, false)
	updateFooterItem(leftFooterItem, "ESC/CTRL+C=Exit | TAB=Change focus | ENTER=Details", tview.AlignLeft, tcell.ColorWhite)

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
    
	ssmParamGrid = tview.NewGrid().SetRows(0,0,0,0,1).SetColumns(0,0,0,0,0)
	ssmParamGrid.AddItem(ssmParamDetailsForm,1,1,2,3,0,0,true)

	ssmParamFooterItemLeft := tview.NewTextView().SetText("").SetTextAlign(tview.AlignLeft)
	ssmParamGrid.AddItem(ssmParamFooterItemLeft,4,0,1,1,0,0,false)
	pages.AddPage("ssmParam",ssmParamGrid , true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Anything handled here will be executed on the main thread
		if ssmParamDetailsForm.HasFocus() {
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
