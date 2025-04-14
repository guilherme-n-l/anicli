package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"anicli/config"
	"anicli/platform/darwin"
	"anicli/platform/linux"
	"anicli/platform/windows"
)

const (
	authCodeBaseURL    = "https://anilist.co/api/v2/oauth"
	authTokenBaseURL   = "https://anilist.co/api/v2/oauth/token"
	authTokenGrantType = "authorization_code"
	JSON               = "application/json"
)

func openBrowser(url string) error {
	var (
		runArgs []string
		runCmd  string
	)

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

func read(msg string) (string, error) {
	fmt.Print(msg)

	reader := bufio.NewReader(os.Stdin)

	res, err := reader.ReadString('\n')

	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(res), nil
}

func inputAuthValues(m map[string]string) error {
	for _, k := range sortedKeys(m) {
		input, err := read(fmt.Sprintf("Enter %s: ", k))

		if err != nil {
			return err
		}

		m[k] = input
	}

	return nil
}

func getCode(m map[string]string) error {
	params := url.Values{}
	params.Add("client_id", m["client_id"])
	params.Add("redirect_uri", m["redirect_uri"])
	params.Add("response_type", "code")

	url := fmt.Sprintf("%s/authorize?%s", authCodeBaseURL, params.Encode())

	openBrowser(url)

	input, err := read("Enter code: ")

	if err != nil {
		return err
	}

	m["code"] = input

	return nil
}

func reqAccessToken(m map[string]string) (*http.Response, error) {
	reqBody, err := json.Marshal(m)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(reqBody))

	res, err := http.Post(authTokenBaseURL, JSON, bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	return res, nil
}

func parseAccessTokenResponse(body []byte) (string, error) {

	var resBodyMap map[string]any

	if err := json.Unmarshal(body, &resBodyMap); err != nil {
		return "", err
	}

	token, ok := resBodyMap["access_token"].(string)

	if !ok {
		return "", fmt.Errorf("Could not parse token from response body")
	}

	return token, nil
}

func getAccessToken(m map[string]string) (string, error) {
	res, err := reqAccessToken(m)

	if err != nil {
		return "", nil
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return "", nil
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d: %s", res.StatusCode, string(resBody))
	}

	return parseAccessTokenResponse(resBody)
}

func Login() error {
	cfg, err := config.GetUserConfig()

	if err != nil {
		return nil
	}

	if cfg.GetAuthToken() != "" {
		fmt.Println("User already logged in.")
		return nil
	}

	var tokenRequest = map[string]string{
		"client_id":     "",
		"client_secret": "",
		"redirect_uri":  "",
	}

	if err := inputAuthValues(tokenRequest); err != nil {
		return err
	}

	if err := getCode(tokenRequest); err != nil {
		return err
	}

	for k, v := range tokenRequest {
		fmt.Printf("%s: %s\n", k, v)
	}

	tokenRequest["grant_type"] = authTokenGrantType

	token, err := getAccessToken(tokenRequest)

	if err != nil {
		return err
	}

	cfg.SetAuthToken(token)

	return nil
}
