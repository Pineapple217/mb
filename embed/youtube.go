package embed

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Pineapple217/mb/database"
)

var imgs = []string{
	"maxresdefault.jpg",
	"mqdefault.jpg",
	"0.jpg", // always present
}

func YoutubeScrape(ctx context.Context, id string) database.YtThumbCache {
	thumb := youtubeGetThumb(id)

	queries := database.GetQueries()
	ytc, err := queries.CreateTYThumbCache(ctx, database.CreateTYThumbCacheParams{
		YtID:    id,
		YtThumb: thumb,
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
		if resp.StatusCode != http.StatusNotFound {
			return v
		}
	}
	return imgs[len(imgs)-1]
}
