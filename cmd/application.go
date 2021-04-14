package cmd

import (
	"fmt"

	"github.com/bnaydenov/ssmbrowse/internal/pkg/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/thoas/go-funk"
)

// var tviewApp *tview.Application
var ssmParamTable *tview.Table

func Entrypoint () {
	fmt.Println("Loading information about your AWS SSM params...")
	
	// box := tview.NewBox().SetBorder(true).SetTitle("ssm browser")
	ssm_table := buildClusterTable()
	
	if err := tview.NewApplication().SetRoot(ssm_table, true).Run(); err != nil {
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

	
	expansions := []int{3, 1, 1, 1}
	alignment := []int{ui.L, ui.L, ui.L, ui.L}

	headers := []string{"Name", "Tier", "Type", "Last modified"}
	ui.AddTableData(table, 0, [][]string{headers}, alignment, expansions, tcell.ColorYellow, false)
	
	params := []int{1,2,3,4,5,6,7,8,9,10}
	
	data := funk.Map(params, func(param int) []string {
		return []string{
			fmt.Sprintf("Name %d", param),
			fmt.Sprintf("Tier %d", param),
			fmt.Sprintf("Type %d", param),
			fmt.Sprintf("Last modified %d", param),
		}
	}).([][]string)
	ui.AddTableData(table, 1, data, alignment, expansions, tcell.ColorWhite, true)

	return table
}