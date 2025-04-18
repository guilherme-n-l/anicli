package client

func GetFullAnimeList(userId int) (string, error) {
	return getFullMediaList(userId, anime)
}
