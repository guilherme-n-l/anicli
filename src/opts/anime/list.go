package anime

import (
	"fmt"

	"anicli/client"
	clientUtils "anicli/client/utils"
	optUtils "anicli/opts/utils"
)

var listCtx optUtils.Context

func init() {
	listCtx = optUtils.NewContext("list", "Manage anime your list", nil, nil)
	listCtx.DefaultHandler = fullList
	listCtx.AddBoolFlags(
		optUtils.NewBoolFlag("help", "h", false, "Show anime list commands", listCtx.Fs.Usage),
		optUtils.NewBoolFlag("noId", "", false, "Print anime entries without id", func() { clientUtils.MediaId = clientUtils.MediaIdFuncs[0] }),
		optUtils.NewBoolFlag("noStatus", "", false, "Print anime entries without status", func() { clientUtils.MediaListFormatType = clientUtils.Blank }),
	)
}

func fullList() {
	id, err := clientUtils.GetUserId()

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
