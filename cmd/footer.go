package cmd

import (
	"fmt"

	awsutils "github.com/bnaydenov/ssmbrowse/internal/pkg/awsutils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func updateFooterItem(item *tview.TextView, text string, alignment int, color tcell.Color) {
	item.SetText(text).SetTextAlign(alignment).SetTextColor(color)
}

func updateRightFooter() {
	accountID, awsRegion, err := awsutils.GetAwsSessionDetails()
	if err != nil {
		errorModal.SetText(fmt.Sprintf("%s", err.Error()))
		pages.SwitchToPage("error")
		return
	}
	updateFooterItem(rightFooterItem, fmt.Sprintf("AWS ID:%s | Region:%s | ver:%s", *accountID, *awsRegion, version), tview.AlignRight, tcell.ColorBlue)
}
