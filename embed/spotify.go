package embed

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"encoding/base64"
	"encoding/json"

	"log/slog"

	"github.com/Pineapple217/mb/database"
)

type RootCoverArt struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type RootArtist struct {
	ID      string `json:"id"`
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
}

type RootTrack struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	URI          string `json:"uri"`
	AlbumOfTrack struct {
		CoverArt struct {
			Sources []RootCoverArt `json:"sources"`
		} `json:"coverArt"`
	} `json:"albumOfTrack"`

	Previews struct {
		AudioPreviews struct {
			Items []struct {
				URL string `json:"url"`
			} `json:"items"`
		} `json:"audioPreviews"`
	} `json:"previews"`

	FirstArtist struct {
		Items []RootArtist `json:"items"`
	} `json:"firstArtist"`
}

type Root struct {
	Entities struct {
		Items map[string]RootTrack `json:"items"`
	} `json:"entities"`
}

const query = `<script\s+id="initial-state"\s+type="text/plain">([^<]+)</script>`

func SpotifyScrape(ctx context.Context, url string) database.SpotifyCache {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(query)
	s := re.FindStringSubmatch(string(body))
	if len(s) == 1 {
		panic("not matches found")
	}
	js, err := base64.StdEncoding.DecodeString(string(s[1]))
	if err != nil {
		panic(err)
	}
	var root Root
	err = json.Unmarshal(js, &root)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	keys := make([]string, 0, len(root.Entities.Items))
	for k := range root.Entities.Items {
		keys = append(keys, k)
	}

	track := root.Entities.Items[keys[0]]
	profile := track.FirstArtist.Items[0]
	queries := database.GetQueries()
	sc, err := queries.CreateSpotifyCache(ctx, database.CreateSpotifyCacheParams{
		TrackID:         track.ID,
		TrackName:       track.Name,
		ArtistName:      profile.Profile.Name,
		ArtistID:        profile.ID,
		CoverArtUrl:     track.AlbumOfTrack.CoverArt.Sources[1].URL,
		AudioPreviewUrl: track.Previews.AudioPreviews.Items[0].URL,
	})
	if err != nil {
		panic(err)
	}
	slog.Info("spotify scrape", "id", sc.TrackID)
	return sc
}
