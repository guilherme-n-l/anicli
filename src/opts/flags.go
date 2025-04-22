package opts

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"anicli/opts/anime"
	"anicli/opts/config"
	"anicli/opts/utils"
)

const versionFile = "VERSION"

var (
	Ctx        utils.Context
	appVersion string
)

func init() {
	Ctx = utils.NewContext("", "", &[]*utils.Context{&anime.Ctx, &config.Ctx}, nil)

	Ctx.AddBoolFlags(
		utils.NewBoolFlag("help", "h", false, "Show application commands", Ctx.Fs.Usage),
		utils.NewBoolFlag("version", "v", false, "Show application version", version),
	)
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
