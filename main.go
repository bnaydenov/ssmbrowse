package main

import "github.com/bnaydenov/ssmbrowse/cmd"

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
    builtBy = "unknown"
)

func main() {
	
	buildData := map[string]interface{}{
		"version": version,
		"commit": commit,
		"date": date,
		"builtBy": builtBy,
	}

	// fmt.Printf("%s", buildData["version"])
	cmd.Entrypoint(buildData)
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"strings"

// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// )

// func main() {
// 	// Initialize application
// 	app := tview.NewApplication()

// 	// Create label
// 	label := tview.NewTextView().SetText("Please enter your name:")

// 	// Create input field
// 	input := tview.NewInputField()

// 	// Create submit button
// 	btn := tview.NewButton("Submit")

// 	// Create empty Box to pad each side of appGrid
// 	bx := tview.NewBox()

// 	// Create Grid containing the application's widgets
// 	appGrid := tview.NewGrid().
// 		SetColumns(-1, 24, 16, -1).
// 		SetRows(-1, 2, 3, -1).
// 		AddItem(bx, 0, 0, 3, 1, 0, 0, false). // Left - 3 rows
// 		AddItem(bx, 0, 1, 1, 1, 0, 0, false). // Top - 1 row
// 		AddItem(bx, 0, 3, 3, 1, 0, 0, false). // Right - 3 rows
// 		AddItem(bx, 3, 1, 1, 1, 0, 0, false). // Bottom - 1 row
// 		AddItem(label, 1, 1, 1, 1, 0, 0, false).
// 		AddItem(input, 1, 2, 1, 1, 0, 0, false).
// 		AddItem(btn, 2, 1, 1, 2, 0, 0, false)

// 	// submittedName is toggled each time Enter is pressed
// 	var submittedName bool
// 	var pages *tview.Pages

// 	// Capture user input
// 	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		// Anything handled here will be executed on the main thread
// 		switch event.Key() {
// 		case tcell.KeyEnter:
// 			submittedName = !submittedName

// 			if submittedName {
// 				name := input.GetText()
// 				if strings.TrimSpace(name) == "" {
// 					name = "Anonymous"
// 				}

// 				// Create a modal dialog
// 				m := tview.NewModal().
// 					SetText(fmt.Sprintf("Greetings, %s!", name)).
// 					AddButtons([]string{"Hello"})

// 				// Display and focus the dialog
// 				app.SetRoot(m, true).SetFocus(m)
// 			} else {
// 				// Clear the input field
// 				input.SetText("")

// 				// Display appGrid and focus the input field
// 				app.SetRoot(appGrid, true).SetFocus(input)
// 			}
// 			return nil
// 		case tcell.KeyEsc:
// 			// Exit the application
// 			app.Stop()
// 			return nil
// 		}

// 		return event
// 	})

// 	pages = tview.NewPages()

// 	pages.AddPage("main", appGrid, true, true)
// 	// Set the grid as the application root and focus the input field
// 	app.SetRoot(pages, true).SetFocus(input)

// 	// Run the application
// 	err := app.Run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
