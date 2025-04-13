package opts

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"anicli/client"
)

const versionFile = "VERSION"

var version string

func newBoolFlag(maximized string, minimized string, value bool, desc string) *bool {
	newFlag := flag.Bool(maximized, value, desc)

	if minimized != "" {
		flag.BoolVar(newFlag, minimized, value, desc)
	}

	return newFlag
}

func getVersion() error {
	if version != "" {
		return nil
	}

	file, err := os.Open(versionFile)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		version = scanner.Text()
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

func Version() {
	if err := getVersion(); err != nil {
		panic(err)
	}

	fmt.Printf("AniCLI version: %s\n", version)
}

func help() {
	Version()
	fmt.Println("AniCLI Usage")
	flag.PrintDefaults()
}

func login() {
	if err := client.Login(); err != nil {
		panic(err)
	}
}

func ParseFlags() {
	flags := map[*bool]func(){
		newBoolFlag("help", "h", false, "Show application commands"):   help,
		newBoolFlag("version", "v", false, "Show application version"): Version,
		newBoolFlag("login", "l", false, "Login into AniList"):         login,
	}

	flag.Parse()

	for f, handler := range flags {
		if *f {
			handler()
		}
	}
}
