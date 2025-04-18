package anime

import (
	"anicli/opts/utils"
	"fmt"
)

var Ctx utils.Context

func init() {
	Ctx = utils.NewContext("anime", "Manage anime in your list or from AniList", &[]*utils.Context{&listCtx})
	Ctx.DefaultHandler = help
	Ctx.AddBoolFlags([]utils.BoolFlag{
		utils.NewBoolFlag("help", "h", false, "Show anime commands", help),
	})
}

func help() {
	fmt.Println("Help from anime")
}
