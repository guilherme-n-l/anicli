package config

import (
	"log"

	"anicli/config/utils"
)

var hasAuthToken = len(utils.GetUserConfig().Authentication.AuthToken) != 0

func GetAuthToken() string {
	if !hasAuthToken {
		log.Fatalln("Please login using `anicli config login`")
	}

	return utils.GetUserConfig().Authentication.AuthToken
}
