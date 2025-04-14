package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	AuthToken string `toml:"auth_token"`
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
	file, err := getFile()

	if err != nil {
		return err
	}

	defer file.Close()

	if err = toml.NewEncoder(file).Encode(&Config{}); err != nil {
		return err
	}

	return nil
}

func GetUserConfig() (*Config, error) {
	if userConfig != nil {
		return userConfig, nil
	}

	if err := ensureDirExists(); err != nil {
		return nil, err
	}

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

	var cfg Config

	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, err
	}

	userConfig = &cfg
	return userConfig, nil
}

func (cfg *Config) GetAuthToken() string {
	return cfg.AuthToken
}

func (cfg *Config) SetAuthToken(token string) error {
	file, err := getFile()

	if err != nil {
		return nil
	}

	defer file.Close()

	cfg.AuthToken = token

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		return err
	}

	userConfig = cfg
	return nil
}
