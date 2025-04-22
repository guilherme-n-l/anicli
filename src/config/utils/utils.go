package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type nameStyle string

const (
	romanji nameStyle = "romanji"
	native  nameStyle = "native"
)

type statusStyle string

const (
	emoji  statusStyle = "emoji"
	letter statusStyle = "letter"
	blanks statusStyle = "blank"
)

type Config struct {
	Preferences struct {
		NameStyle nameStyle `toml:"name_style"`
	} `toml:"Preferences"`
	Authentication struct {
		AuthToken string `toml:"auth_token"`
	} `toml:"Authentication"`
}

const CONFIG_DIRNAME = ".anicli"
const CONFIG_FILENAME = "config.toml"

var userConfig *Config

func getDirPath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(home, CONFIG_DIRNAME), nil
}

func getFilePath() (string, error) {
	dirPath, err := getDirPath()

	if err != nil {
		return "", err
	}

	return filepath.Join(dirPath, CONFIG_FILENAME), nil
}

func getFile() (*os.File, error) {
	filePath, err := getFilePath()

	if err != nil {
		return nil, err
	}

	file, err := os.Create(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func ensureDirExists() error {
	dirPath, err := getDirPath()

	if err != nil {
		return err
	}

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func createDefaultConfig() error {
	if err := (&Config{}).writeToFile(); err != nil {
		return err
	}

	return nil
}

func readFromConfigFile() (*Config, error) {
	var cfg Config

	filePath, err := getFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)

	if err != nil || len(data) == 0 {
		if os.IsNotExist(err) || len(data) == 0 {
			if err = createDefaultConfig(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, err
	}

	return &cfg, err
}

func (cfg *Config) writeToFile() error {
	file, err := getFile()
	if err != nil {
		return nil
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		return err
	}

	return nil
}

func CommitChanges() error {
	return GetUserConfig().writeToFile()
}

func GetUserConfig() *Config {
	if userConfig != nil {
		return userConfig
	}

	if err := ensureDirExists(); err != nil {
		log.Fatalln(err)
	}

	cfg, err := readFromConfigFile()
	if err != nil {
		log.Fatalln(err)
	}

	userConfig = cfg
	return userConfig
}

func SetAuthToken(token string) {
	GetUserConfig().Authentication.AuthToken = token

	if err := CommitChanges(); err != nil {
		log.Fatalln(err)
	}
}
