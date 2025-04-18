package opts

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"anicli/client"
	"anicli/opts/anime"
	"anicli/opts/utils"
)

const versionFile = "VERSION"

var (
	Ctx utils.Context

	// mangaFs    = manga.Fs
	// userFs     = user.Fs
	// configFs   = config.Fs

	appVersion string
)

func init() {
	Ctx = utils.NewContext("", "",
		&[]*utils.Context{
			&anime.Ctx,
		},
	)
	
	Ctx.DefaultHandler = help

	Ctx.AddBoolFlags([]utils.BoolFlag{
		utils.NewBoolFlag("help", "h", false, "Show application commands", help),
		utils.NewBoolFlag("version", "v", false, "Show application version", version),
		utils.NewBoolFlag("login", "l", false, "Login into AniList", login),
	})
}

func getVersion() error {
	if appVersion != "" {
		return nil
	}

	file, err := os.Open(versionFile)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		appVersion = scanner.Text()
	} else {
		return fmt.Errorf("Empty version file")
	}

	if scanner.Scan() {
		log.Println("WARN: Version file contains more than one line. Reading first line only")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func version() {
	if err := getVersion(); err != nil {
		panic(err)
	}

	fmt.Printf("AniCLI version: %s\n", appVersion)
}

func help() {
	fmt.Println("Usage of anicli")
	Ctx.Fs.PrintDefaults()
	Ctx.PrintSubContexts()
}

func login() {
	if err := client.Login(); err != nil {
		panic(err)
	}
}
