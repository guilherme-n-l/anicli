package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"anicli/config"
)

type gQLQuery struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

var (
	MediaIdFuncs = []func(int) string{
		func(_ int) string { return "" },
		func(id int) string { return fmt.Sprintf("[%d] ", id) },
	}

	MediaId = MediaIdFuncs[1]
)

type mediaListType string

const (
	completed mediaListType = "Completed"
	dropped   mediaListType = "Dropped"
	watching  mediaListType = "Watching"
	planning  mediaListType = "Planning"
)

type mediaList struct {
	Name    string `json:"name"`
	Entries []struct {
		ID    int `json:"id"`
		Media struct {
			ID    int `json:"id"`
			Title struct {
				Romaji string `json:"romaji"`
			} `json:"title"`
		} `json:"media"`
	} `json:"entries"`
}

type mediaListFormatType int

const (
	Emoji mediaListFormatType = iota
	Letter
	Blank
)

type mediaListFormat struct {
	Type mediaListFormatType
}

var MediaListFormatType = Emoji

func mediaListRune(mlt mediaListType) rune {
	symbols := map[mediaListType][]rune{
		completed: {'C', '✅'},
		dropped:   {'D', '🚮'},
		watching:  {'W', '📺'},
		planning:  {'P', '🔜'},
	}

	switch MediaListFormatType {
	case Letter:
		return symbols[mlt][0]
	case Blank:
		return '\u0000'
	default:
		return symbols[mlt][1]
	}
}

type mediaListResponse struct {
	Data struct {
		MediaListCollection struct {
			Lists []mediaList `json:"lists"`
		} `json:"MediaListCollection"`
	} `json:"data"`
}

type media string

const (
	Anime media = "ANIME"
	Manga media = "MANGA"
)

const API_URL = "https://graphql.anilist.co"

func sendGQLRequest(query string, variables map[string]any) ([]byte, error) {
	data, err := json.Marshal(gQLQuery{
		Query:     query,
		Variables: variables,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.GetAuthToken())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d, response: %s", res.StatusCode, string(body))
	}

	return body, nil
}

func parseMediaLists(lists mediaListResponse) string {
	res := ""

	for _, list := range lists.Data.MediaListCollection.Lists {
		for _, entry := range list.Entries {
			res += fmt.Sprintf(
				"%s%c %s\n",
				MediaId(entry.Media.ID),
				mediaListRune(mediaListType(list.Name)),
				entry.Media.Title.Romaji)
		}
	}

	return res
}

func GetFullMediaList(userId int, media media) (string, error) {
	res, err := sendGQLRequest(
		"query($type:MediaType!,$userId:Int!){MediaListCollection(type:$type,userId:$userId){lists{name,entries{id,media{id,title{romaji}}}}}}",
		map[string]any{
			"type":   media,
			"userId": userId,
		})

	if err != nil {
		return "", err
	}

	var lists mediaListResponse

	err = json.Unmarshal(res, &lists)
	if err != nil {
		return "", err
	}

	return parseMediaLists(lists), nil
}

func GetMediaById(mediaId int) (string, error) {
	res, err := sendGQLRequest(
		"query($id:Int){Media(id:$id){id,title{romaji,english,native}}}",
		map[string]any{
			"id": mediaId,
		})

	if err != nil {
		return "", err
	}

	fmt.Println(string(res))

	return "", nil
}

func GetUserId() (int, error) {
	query := "query{Viewer{id}}"

	res, err := sendGQLRequest(query, nil)
	if err != nil {
		return 0, err
	}

	var typedRes struct {
		Data struct {
			Viewer struct {
				ID int `json:"id"`
			} `json:"Viewer"`
		} `json:"data"`
	}

	err = json.Unmarshal(res, &typedRes)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return typedRes.Data.Viewer.ID, nil
}
