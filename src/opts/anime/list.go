package anime

import (
	"anicli/opts/utils"
	"fmt"
)

var listCtx utils.Context

func init() {
	listCtx = utils.NewContext("list", "Manage anime your list", nil)
	listCtx.DefaultHandler = listHelp
	listCtx.AddBoolFlags([]utils.BoolFlag{
		utils.NewBoolFlag("help", "h", false, "Show anime list commands", listHelp),
	})
}

func listHelp() {
	fmt.Println("Help from list")
}
