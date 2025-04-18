package anime

import (
	"fmt"
	"os"

	"anicli/client"
	"anicli/opts/utils"
)

var listCtx utils.Context

func init() {
	listCtx = utils.NewContext("list", "Manage anime your list", nil)
	listCtx.DefaultHandler = fullList
	listCtx.AddBoolFlags([]utils.BoolFlag{
		utils.NewBoolFlag("help", "h", false, "Show anime list commands", listHelp),
		utils.NewBoolFlag("noId", "", false, "Print anime entries without id", func() { client.MediaId = client.MediaIdFuncs[0] }),
		utils.NewBoolFlag("noStatus", "", false, "Print anime entries without status", func() {client.MediaListFormatType = client.Blank}),
	})
}

func listHelp() {
	fmt.Println("Usage of anicli lists")
	listCtx.Fs.PrintDefaults()
	os.Exit(1)
}

func fullList() {
	id, err := client.GetUserId()

	if err != nil {
		return
	}

	listStr, err := client.GetFullAnimeList(id)
	if err != nil {
		fmt.Println("Could not print anime list")
		return
	}

	fmt.Println(listStr)
}
