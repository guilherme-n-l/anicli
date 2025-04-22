package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"anicli/config/utils"
	"anicli/platform/darwin"
	"anicli/platform/linux"
	"anicli/platform/windows"
	appUtils "anicli/utils"
)

const (
	authCodeBaseURL       = "https://anilist.co/api/v2/oauth"
	authTokenBaseURL      = "https://anilist.co/api/v2/oauth/token"
	authTokenGrantType    = "authorization_code"
	authTokenResponseType = "code"
	JSON                  = "application/json"
)

var (
	redirect_url = "http://localhost"
	codeListener net.Listener
	userCodeChan chan string
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

func inputAuthValues() (appUtils.Map[string, string], error) {
	m := appUtils.Map[string, string]{
		"client_id":     "",
		"client_secret": "",
		"redirect_uri":  redirect_url,
		"grant_type":    authTokenGrantType,
		"response_type": authTokenResponseType,
	}

	for _, k := range appUtils.Sort(m.Keys()) {
		if len(m[k]) != 0 {
			continue
		}

		input, err := appUtils.Read(fmt.Sprintf("Enter %s: ", k))
		if err != nil {
			return m, err
		}

		m[k] = input
	}

	return m, nil
}

func handleCode(w http.ResponseWriter, r *http.Request) {
	userCode := r.URL.Query().Get("code")

	if userCode == "" {
		http.Error(w, "Missing 'code' query parameter", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	userCodeChan <- userCode
}

func authServer() (string, error) {
	server := &http.Server{
		Addr: ":0",
	}

	http.HandleFunc("/callback", handleCode)

	var err error
	codeListener, err = net.Listen("tcp", server.Addr)
	if err != nil {
		return "", err
	}

	go func() {
		if err := server.Serve(codeListener); err != nil {
			log.Fatalln(err)
		}
	}()

	return codeListener.Addr().String()[4:], nil
}

func getCode(m map[string]string) {
	url := fmt.Sprintf("%s/authorize?", authCodeBaseURL)
	for _, i := range []string{"client_id", "redirect_uri", "response_type"} {
		url += fmt.Sprintf("%s=%s&", i, m[i])
	}

	openBrowser(url[:len(url)-1])
}

func reqAccessToken(m map[string]string) (*http.Response, error) {
	reqBody, err := json.Marshal(m)

	if err != nil {
		return nil, err
	}

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
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d: %s", res.StatusCode, string(resBody))
	}

	return parseAccessTokenResponse(resBody)
}

func login() error {
	userCodeChan = make(chan string)

	port, err := authServer()
	if err != nil {
		return err
	}

	redirect_url += port + "/callback"

	for _, msg := range []string{
		"Step by step guide into authenticating AniList account\n",
		"\t1. Login to your account by accessing https://anilist.co/\n",
		"\t2. Go to: https://anilist.co/settings/developer\n",
		"\t3. Select the `Create New Client` button\n",
		"\t4. Fill in the form with the following information",
		"\t- Name: anicli (Note: Can be any name of your choice)",
		fmt.Sprintf("\t- Redirect URL: %s\n", redirect_url),
		"\t5. Select the `Save button`\n",
		"\t6. Enter the following info:",
	} {
		fmt.Println(msg)
	}

	reqMap, err := inputAuthValues()
	if err != nil {
		return err
	}

	getCode(reqMap)

	fmt.Println("Awaiting AniList response arrival...")

	reqMap["code"] = <-userCodeChan
	if len(reqMap["code"]) == 0 {
		return fmt.Errorf("Error: Empty user code")
	}

	token, err := getAccessToken(reqMap)
	if err != nil {
		return err
	}

	fmt.Println("Got access token. Saving to config file...")
	utils.SetAuthToken(token)
	fmt.Println("User authenticated successfully.")

	return nil
}

func Login(forced bool) {
	if !forced && hasAuthToken {
		fmt.Println("User already logged in. use `--force` flag to re-login")
		os.Exit(1)
	}

	err := login()

	if err != nil {
		log.Fatalln(err)
	}
}
