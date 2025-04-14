package opts

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"anicli/client"
	"anicli/opts/anime"
	"anicli/opts/utils"
)

const versionFile = "VERSION"

var (
	fs          = flag.CommandLine
	mainContext = utils.FlagContext{Fs: fs, Flags: map[*bool]func(){
		utils.NewBoolFlag("help", "h", false, "Show application commands", nil):   Help,
		utils.NewBoolFlag("version", "v", false, "Show application version", nil): version,
		utils.NewBoolFlag("login", "l", false, "Login into AniList", nil):         login,
	}}
	animeContext = anime.Context
	// mangaFs    = manga.Fs
	// userFs     = user.Fs
	// configFs   = config.Fs
	appVersion string
)

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
		fmt.Fprintln(os.Stderr, "WARN: Version file contains more than one line. Reading first line only")
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

func Help() {
	fmt.Println("Usage of anicli")
	flag.PrintDefaults()
}

func login() {
	if err := client.Login(); err != nil {
		panic(err)
	}
}

func getContext() utils.FlagContext {
	if len(os.Args) < 2 {
		return mainContext
	}

	switch os.Args[1] {
	case "anime":
		return animeContext
		// case "manga": return mangaContext
		// case "user": return userContext
		// case "config": return configContext
	default:
		return mainContext
	}
}

func preventInvalidArgs(context utils.FlagContext) {
	if context.Fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Invalid arg provided not allowed: %s\n", context.Fs.Arg(0))
		os.Exit(1)
	}
}

func ParseArgs() {
	context := getContext()

	if context.Fs == fs {
		flag.Parse()
	} else {
		context.Fs.Parse(os.Args[2:])
	}

	preventInvalidArgs(context)

	for f, handler := range context.Flags {
		if *f {
			handler()
		}
	}
}
