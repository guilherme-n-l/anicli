package client

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"

	"anicli/platform/darwin"
	"anicli/platform/linux"
	"anicli/platform/windows"
)

const authTokenBaseURL = "https://anilist.co/api/v2/oauth/token"

func openBrowser(url string) error {
	var runCmd string
	var runArgs []string
	switch runtime.GOOS {
	case "windows":
		runCmd, runArgs = windows.OpenCmd, append(windows.BrowserArgs, url)
	case "linux":
		runCmd, runArgs = linux.OpenCmd, []string{url}
	case "darwin":
		runCmd, runArgs = darwin.OpenCmd, []string{url}
	default:
		return fmt.Errorf("Running in unsupported platform %s\n", runtime.GOOS)
	}

	return exec.Command(runCmd, runArgs...).Start()
}

func sortedKeys(m map[string]string) []string {
	var keys []string

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func inputMapValues(m map[string]string) error {
	for _, k := range sortedKeys(m) {
		fmt.Printf("Enter %s: ", k)

		reader := bufio.NewReader(os.Stdin)

		input, err := reader.ReadString('\n')

		if err != nil {
			return err
		}

		m[k] = input
	}

	return nil
}

func Login() error {
	var tokenRequest = map[string]string{
		"client_id":     "",
		"client_secret": "",
		"redirect_uri":  "",
		"code":          "",
	}

	if err := inputMapValues(tokenRequest); err != nil {
		return err
	}

	for k, v := range tokenRequest {
		fmt.Printf("%s: %s", k, v)
	}

	return nil
}
