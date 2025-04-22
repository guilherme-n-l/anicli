package client

import "anicli/client/utils"

func GetFullMangaList(userId int) (string, error) {
	return utils.GetFullMediaList(userId, utils.Manga)
}

var GetMangabyId = utils.GetMediaById
