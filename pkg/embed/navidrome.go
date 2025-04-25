package embed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/Pineapple217/mb/pkg/database"
)

type navidromeTrack struct {
	Id     string `json:"id"`
	Tracks []struct {
		Id     string `json:"id"`
		Title  string `json:"title"`
		Artist string `json:"artist"`
	} `json:"tracks"`
}

var (
	naviRe = regexp.MustCompile(`(?s)window\.__SHARE_INFO__\s*=\s*("(?:\\.|[^"\\])*")`)
)

func NavidromeScrape(ctx context.Context, q *database.Queries, id string) (database.NavidromeCache, error) {
	resp, err := http.Get(config.NavidromePrefix + id)
	if err != nil {
		return database.NavidromeCache{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return database.NavidromeCache{}, err
	}

	matches := naviRe.FindStringSubmatch(string(body))
	if len(matches) != 2 {
		return database.NavidromeCache{}, errors.New("window.__APP_CONFIG__ not found")
	}

	stringJson := strings.ReplaceAll(matches[1], "\\\"", "\"")
	stringJson = stringJson[1 : len(stringJson)-1]
	fmt.Println(stringJson)
	var tracks navidromeTrack
	json.Unmarshal([]byte(stringJson), &tracks)

	if len(tracks.Tracks) == 0 {
		return database.NavidromeCache{}, errors.New("no navidrome tracks found in json")
	}

	nc, err := q.CreateNavidromeCache(ctx, database.CreateNavidromeCacheParams{
		ShareID:    tracks.Id,
		TrackID:    tracks.Tracks[0].Id,
		TrackName:  tracks.Tracks[0].Title,
		ArtistName: tracks.Tracks[0].Artist,
	})
	if err != nil {
		return database.NavidromeCache{}, err
	}

	slog.Info("navidrome scrape", "id", nc.ShareID)
	return nc, nil
}
