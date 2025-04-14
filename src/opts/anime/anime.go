package anime

import (
	"flag"
	"fmt"

	"anicli/opts/utils"
)

var (
	fs      = flag.NewFlagSet("anime", flag.ExitOnError)
	Context = utils.FlagContext{Fs: fs, Flags: map[*bool]func(){
		utils.NewBoolFlag("help", "h", false, "Show application commands", fs): func() {fmt.Println("Hello from anime")},
	}}
)
