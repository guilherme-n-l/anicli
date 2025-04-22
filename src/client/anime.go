package client

import "anicli/client/utils"

func GetFullAnimeList(userId int) (string, error) {
	return utils.GetFullMediaList(userId, utils.Anime)
}

var GetAnimebyId = utils.GetMediaById
