package embed

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Pineapple217/mb/pkg/database"
)

const youtubeMetaDataUrl string = "https://www.youtube.com/oembed?url=http%%3A//youtube.com/watch%%3Fv%%3D%s&format=json"

var imgs = []string{
	"maxresdefault.jpg",
	"mqdefault.jpg",
	"0.jpg", // always present
}

type rootMeta struct {
	Title     string `json:"title"`
	Author    string `json:"author_name"`
	AuthorUrl string `json:"author_url"`
}

func YoutubeScrape(ctx context.Context, q *database.Queries, id string) (database.YoutubeCache, error) {
	thumb, err := youtubeGetThumb(id)
	if err != nil {
		return database.YoutubeCache{}, err
	}
	meta, err := youtubeGetMetaData(id)
	if err != nil {
		return database.YoutubeCache{}, err
	}
	// TODO: possible parallelization but realistically not needed

	ytc, err := q.CreateYoutubebCache(ctx, database.CreateYoutubebCacheParams{
		YtID:      id,
		Thumb:     thumb,
		Title:     meta.Title,
		Author:    meta.Author,
		AuthorUrl: meta.AuthorUrl,
	})
	if err != nil {
		return database.YoutubeCache{}, err
	}

	slog.Info("youtube scrape", "id", ytc.YtID)
	return ytc, nil
}

func youtubeGetThumb(id string) (string, error) {
	for _, v := range imgs {
		resp, err := http.Get(fmt.Sprintf("https://i3.ytimg.com/vi/%s/%s", id, v))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			return v, nil
		}
	}
	return imgs[len(imgs)-1], nil
}

func youtubeGetMetaData(id string) (rootMeta, error) {
	s := fmt.Sprintf(youtubeMetaDataUrl, id)
	resp, err := http.Get(s)
	if err != nil {
		return rootMeta{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rootMeta{}, err
	}
	var root rootMeta
	err = json.Unmarshal(body, &root)
	if err != nil {
		return rootMeta{}, err
	}
	return root, err
}
