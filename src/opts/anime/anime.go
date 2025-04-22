package anime

import (
	"fmt"
	"log"
	"strconv"

	"anicli/client"
	"anicli/opts/utils"
)

var Ctx utils.Context

func init() {
	Ctx = utils.NewContext("anime", "Manage anime in your list or from AniList", &[]*utils.Context{&listCtx}, utils.NewSlugSet(utils.NumberSlug))
	Ctx.DefaultHandler = getAnime
	Ctx.AddBoolFlags(
		utils.NewBoolFlag("help", "h", false, "Show anime commands", Ctx.Fs.Usage),
	)
}

func getAnime() {
	if Ctx.Fs.NArg() == 0 {
		Ctx.Fs.Usage()
	}

	animeId, err := strconv.Atoi(Ctx.Fs.Arg(0))

	if err != nil {
		log.Println(err)
		return
	}

	res, err := client.GetAnimebyId(animeId)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(res)
}
