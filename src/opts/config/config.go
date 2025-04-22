package config

import (
	"anicli/opts/utils"
)

var Ctx utils.Context

func init() {
	Ctx = utils.NewContext("config", "Manage user configuration", &[]*utils.Context{&loginCtx}, nil)
	Ctx.AddBoolFlags(utils.NewBoolFlag("help", "h", false, "Show config commands", Ctx.Fs.Usage))
}
