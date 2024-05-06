package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/Pineapple217/mb/pkg/database"
)

func CleanCache(q *database.Queries) {
	slog.Info("Starting clean databases chaches")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	c, err := q.RemoveUnusedYoutubeCache(ctx)
	if err != nil {
		slog.Warn("Failed to clean database cache", "type", "youtube", "error", err)
	}
	slog.Info("Cleaned cache", "type", "youtube", "count", c)

	c, err = q.RemoveUnusedSpotifyCache(ctx)
	if err != nil {
		slog.Warn("Failed to clean database cache", "type", "spotify", "error", err)
	}
	slog.Info("Cleaned cache", "type", "spotify", "count", c)
}
