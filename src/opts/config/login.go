package config

import (
	"anicli/config"
	"anicli/opts/utils"
)

var loginCtx utils.Context

func init() {
	loginCtx = utils.NewContext("login", "Authenticate into AniList account", &[]*utils.Context{}, nil)
	loginCtx.AddBoolFlags(
		utils.NewBoolFlag("help", "h", false, "Show login commands", loginCtx.Fs.Usage),
		utils.NewBoolFlag("force", "f", false, "Force re-login", forcedLogin),
	)
	loginCtx.DefaultHandler = login
}

func forcedLogin() {
	config.Login(true)
}

func login() {
	config.Login(false)
}
