package embed

import (
	"context"
	"errors"
	"io"
	"net/http"
	"regexp"

	"encoding/base64"
	"encoding/json"

	"log/slog"

	"github.com/Pineapple217/mb/pkg/database"
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

const SpotifyUrlPrefix = "https://open.spotify.com/track/"

func SpotifyScrape(ctx context.Context, q *database.Queries, id string) (database.SpotifyCache, error) {
	resp, err := http.Get(SpotifyUrlPrefix + id)
	if err != nil {
		return database.SpotifyCache{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return database.SpotifyCache{}, err
	}
	re := regexp.MustCompile(query)
	s := re.FindStringSubmatch(string(body))
	if len(s) < 2 {
		return database.SpotifyCache{}, errors.New("track not found")
	}
	js, err := base64.StdEncoding.DecodeString(string(s[1]))
	if err != nil {
		return database.SpotifyCache{}, err
	}
	var root Root
	err = json.Unmarshal(js, &root)
	if err != nil {
		return database.SpotifyCache{}, err
	}
	keys := make([]string, 0, len(root.Entities.Items))
	for k := range root.Entities.Items {
		keys = append(keys, k)
	}

	track := root.Entities.Items[keys[0]]
	profile := track.FirstArtist.Items[0]
	sc, err := q.CreateSpotifyCache(ctx, database.CreateSpotifyCacheParams{
		TrackID:         track.ID,
		TrackName:       track.Name,
		ArtistName:      profile.Profile.Name,
		ArtistID:        profile.ID,
		CoverArtUrl:     track.AlbumOfTrack.CoverArt.Sources[1].URL,
		AudioPreviewUrl: track.Previews.AudioPreviews.Items[0].URL,
	})
	if err != nil {
		return database.SpotifyCache{}, err
	}
	slog.Info("spotify scrape", "id", sc.TrackID)
	return sc, nil
}
