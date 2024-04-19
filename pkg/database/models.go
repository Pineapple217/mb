// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"time"
)

type Mediafile struct {
	ID            int64
	FileName      string
	FilePath      string
	FileExtention string
	FileType      string
	Thumbnail     []byte
	UploadedAt    time.Time
}

type Post struct {
	ID        int64
	CreatedAt int64
	Tags      sql.NullString
	Content   string
}

type SpotifyCache struct {
	ID              int64
	TrackID         string
	TrackName       string
	ArtistName      string
	ArtistID        string
	CoverArtUrl     string
	AudioPreviewUrl string
}

type YoutubeCache struct {
	ID        int64
	YtID      string
	Thumb     string
	Title     string
	Author    string
	AuthorUrl string
}