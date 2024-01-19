package embed

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Pineapple217/mb/database"
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

func YoutubeScrape(ctx context.Context, id string) database.YoutubeCache {
	thumb := youtubeGetThumb(id)
	meta := youtubeGetMetaData(id)
	// TODO: possible parallelization but realistically not needed

	queries := database.GetQueries()
	ytc, err := queries.CreateYoutubebCache(ctx, database.CreateYoutubebCacheParams{
		YtID:      id,
		Thumb:     thumb,
		Title:     meta.Title,
		Author:    meta.Author,
		AuthorUrl: meta.AuthorUrl,
	})

	if err != nil {
		panic(err)
	}

	slog.Info("youtube scrape", "id", ytc.YtID)
	return ytc
}

func youtubeGetThumb(id string) string {
	for _, v := range imgs {
		resp, err := http.Get(fmt.Sprintf("https://i3.ytimg.com/vi/%s/%s", id, v))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			return v
		}
	}
	return imgs[len(imgs)-1]
}

func youtubeGetMetaData(id string) rootMeta {
	s := fmt.Sprintf(youtubeMetaDataUrl, id)
	resp, err := http.Get(s)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var root rootMeta
	err = json.Unmarshal(body, &root)
	if err != nil {
		slog.Error("Error unmarshalling JSON ytmeta", "err", err)
	}
	return root
}
