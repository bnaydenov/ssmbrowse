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

	
	expansions := []int{3, 1, 1, 1,1,1}
	alignment := []int{ui.L, ui.L, ui.L, ui.L,ui.L,ui.L}

	headers := []string{"Name", "Tier", "Type", "Description", "Version","Last modified"}
	ui.AddTableData(table, 0, [][]string{headers}, alignment, expansions, tcell.ColorYellow, false)
	
	var startToken *string
	var params []ssm.ParameterMetadata
	
	params, _  = aws.GetParemters([]string{"/"}, startToken, params)
    
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