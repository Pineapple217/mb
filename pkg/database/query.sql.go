// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
)

const createMediafile = `-- name: CreateMediafile :one
INSERT INTO mediafiles (
  uploaded_at, file_name, file_path, file_type, file_extention, thumbnail
) VALUES (
  strftime('%s', 'now'), ?, ?, ?, ?, ?
)
RETURNING id, file_name, file_path, file_extention, file_type, thumbnail, uploaded_at
`

type CreateMediafileParams struct {
	FileName      string
	FilePath      string
	FileType      string
	FileExtention string
	Thumbnail     []byte
}

func (q *Queries) CreateMediafile(ctx context.Context, arg CreateMediafileParams) (Mediafile, error) {
	row := q.db.QueryRowContext(ctx, createMediafile,
		arg.FileName,
		arg.FilePath,
		arg.FileType,
		arg.FileExtention,
		arg.Thumbnail,
	)
	var i Mediafile
	err := row.Scan(
		&i.ID,
		&i.FileName,
		&i.FilePath,
		&i.FileExtention,
		&i.FileType,
		&i.Thumbnail,
		&i.UploadedAt,
	)
	return i, err
}

const createNavidromeCache = `-- name: CreateNavidromeCache :one
INSERT INTO navidrome_cache (
  share_id, track_id, track_name, artist_name
) VALUES (
  ?, ?, ?, ?
)
RETURNING id, share_id, track_id, track_name, artist_name
`

type CreateNavidromeCacheParams struct {
	ShareID    string
	TrackID    string
	TrackName  string
	ArtistName string
}

func (q *Queries) CreateNavidromeCache(ctx context.Context, arg CreateNavidromeCacheParams) (NavidromeCache, error) {
	row := q.db.QueryRowContext(ctx, createNavidromeCache,
		arg.ShareID,
		arg.TrackID,
		arg.TrackName,
		arg.ArtistName,
	)
	var i NavidromeCache
	err := row.Scan(
		&i.ID,
		&i.ShareID,
		&i.TrackID,
		&i.TrackName,
		&i.ArtistName,
	)
	return i, err
}

const createPost = `-- name: CreatePost :one
INSERT INTO posts (
  created_at, tags, content, html, private
) VALUES (
  strftime('%s', 'now'), ?, ?, ?, ?
)
RETURNING created_at, tags, content, html, private
`

type CreatePostParams struct {
	Tags    sql.NullString
	Content string
	Html    string
	Private int64
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Tags,
		arg.Content,
		arg.Html,
		arg.Private,
	)
	var i Post
	err := row.Scan(
		&i.CreatedAt,
		&i.Tags,
		&i.Content,
		&i.Html,
		&i.Private,
	)
	return i, err
}

const createSpotifyCache = `-- name: CreateSpotifyCache :one
INSERT INTO spotify_cache (
  track_id, track_name, artist_name, artist_id, cover_art_url, audio_preview_url
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING id, track_id, track_name, artist_name, artist_id, cover_art_url, audio_preview_url
`

type CreateSpotifyCacheParams struct {
	TrackID         string
	TrackName       string
	ArtistName      string
	ArtistID        string
	CoverArtUrl     string
	AudioPreviewUrl string
}

func (q *Queries) CreateSpotifyCache(ctx context.Context, arg CreateSpotifyCacheParams) (SpotifyCache, error) {
	row := q.db.QueryRowContext(ctx, createSpotifyCache,
		arg.TrackID,
		arg.TrackName,
		arg.ArtistName,
		arg.ArtistID,
		arg.CoverArtUrl,
		arg.AudioPreviewUrl,
	)
	var i SpotifyCache
	err := row.Scan(
		&i.ID,
		&i.TrackID,
		&i.TrackName,
		&i.ArtistName,
		&i.ArtistID,
		&i.CoverArtUrl,
		&i.AudioPreviewUrl,
	)
	return i, err
}

const createYoutubebCache = `-- name: CreateYoutubebCache :one
INSERT INTO youtube_cache (
  yt_id, thumb, title, author, author_url
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING id, yt_id, thumb, title, author, author_url
`

type CreateYoutubebCacheParams struct {
	YtID      string
	Thumb     string
	Title     string
	Author    string
	AuthorUrl string
}

func (q *Queries) CreateYoutubebCache(ctx context.Context, arg CreateYoutubebCacheParams) (YoutubeCache, error) {
	row := q.db.QueryRowContext(ctx, createYoutubebCache,
		arg.YtID,
		arg.Thumb,
		arg.Title,
		arg.Author,
		arg.AuthorUrl,
	)
	var i YoutubeCache
	err := row.Scan(
		&i.ID,
		&i.YtID,
		&i.Thumb,
		&i.Title,
		&i.Author,
		&i.AuthorUrl,
	)
	return i, err
}

const deleteMediafile = `-- name: DeleteMediafile :exec
DELETE FROM mediafiles
WHERE id = ?
`

func (q *Queries) DeleteMediafile(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteMediafile, id)
	return err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE created_at = ?
`

func (q *Queries) DeletePost(ctx context.Context, createdAt int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, createdAt)
	return err
}

const getAllTags = `-- name: GetAllTags :many
WITH split (
    tag,
    tags_remaining
)
AS (-- Initial query
    SELECT '',
           tags || ' '
      FROM posts
      WHERE private <= ?
    UNION ALL
    SELECT trim(substr(tags_remaining, 0, instr(tags_remaining, ' ') ) ),
           substr(tags_remaining, instr(tags_remaining, ' ') + 1) 
      FROM split
     WHERE tags_remaining != ''
)
SELECT MIN(tag) as tag,
       COUNT( * ) AS tag_count
FROM split
WHERE tag != ''
GROUP BY tag
ORDER BY tag_count DESC
`

type GetAllTagsRow struct {
	Tag      interface{}
	TagCount int64
}

func (q *Queries) GetAllTags(ctx context.Context, private int64) ([]GetAllTagsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllTags, private)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllTagsRow
	for rows.Next() {
		var i GetAllTagsRow
		if err := rows.Scan(&i.Tag, &i.TagCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMediaThunbnail = `-- name: GetMediaThunbnail :one
SELECT thumbnail FROM mediafiles
WHERE file_path = ? LIMIT 1
`

func (q *Queries) GetMediaThunbnail(ctx context.Context, filePath string) ([]byte, error) {
	row := q.db.QueryRowContext(ctx, getMediaThunbnail, filePath)
	var thumbnail []byte
	err := row.Scan(&thumbnail)
	return thumbnail, err
}

const getMediafile = `-- name: GetMediafile :one
SELECT id, file_name, file_path, file_extention, file_type, thumbnail, uploaded_at FROM mediafiles
WHERE id = ? LIMIT 1
`

func (q *Queries) GetMediafile(ctx context.Context, id int64) (Mediafile, error) {
	row := q.db.QueryRowContext(ctx, getMediafile, id)
	var i Mediafile
	err := row.Scan(
		&i.ID,
		&i.FileName,
		&i.FilePath,
		&i.FileExtention,
		&i.FileType,
		&i.Thumbnail,
		&i.UploadedAt,
	)
	return i, err
}

const getNavidromeCache = `-- name: GetNavidromeCache :one
SELECT id, share_id, track_id, track_name, artist_name FROM navidrome_cache
WHERE share_id = ? LIMIT 1
`

func (q *Queries) GetNavidromeCache(ctx context.Context, shareID string) (NavidromeCache, error) {
	row := q.db.QueryRowContext(ctx, getNavidromeCache, shareID)
	var i NavidromeCache
	err := row.Scan(
		&i.ID,
		&i.ShareID,
		&i.TrackID,
		&i.TrackName,
		&i.ArtistName,
	)
	return i, err
}

const getPost = `-- name: GetPost :one
SELECT created_at, tags, content, html, private FROM posts
WHERE created_at = ? LIMIT 1
`

func (q *Queries) GetPost(ctx context.Context, createdAt int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, createdAt)
	var i Post
	err := row.Scan(
		&i.CreatedAt,
		&i.Tags,
		&i.Content,
		&i.Html,
		&i.Private,
	)
	return i, err
}

const getPostCount = `-- name: GetPostCount :one
SELECT COUNT(*)
FROM posts
WHERE private <= ?
`

func (q *Queries) GetPostCount(ctx context.Context, private int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPostCount, private)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getPostLatest = `-- name: GetPostLatest :one
SELECT created_at, tags, content, html, private FROM posts
WHERE private = 0
ORDER BY created_at DESC LIMIT 1
`

func (q *Queries) GetPostLatest(ctx context.Context) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostLatest)
	var i Post
	err := row.Scan(
		&i.CreatedAt,
		&i.Tags,
		&i.Content,
		&i.Html,
		&i.Private,
	)
	return i, err
}

const getPostPage = `-- name: GetPostPage :one
SELECT 
CAST(
    CASE 
        WHEN EXISTS (SELECT 1 FROM posts WHERE posts.created_at = ?1)
        THEN CEIL((SELECT COUNT(*) FROM posts WHERE (posts.created_at >= (SELECT posts.created_at FROM posts WHERE created_at = ?1)) and posts.private <= ?2)  / 25.0) - 1
        ELSE -1
    END AS INT
) AS page_number
`

type GetPostPageParams struct {
	ID int64
	P  int64
}

func (q *Queries) GetPostPage(ctx context.Context, arg GetPostPageParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPostPage, arg.ID, arg.P)
	var page_number int64
	err := row.Scan(&page_number)
	return page_number, err
}

const getSpotifyCache = `-- name: GetSpotifyCache :one
SELECT id, track_id, track_name, artist_name, artist_id, cover_art_url, audio_preview_url FROM spotify_cache
WHERE track_id = ? LIMIT 1
`

func (q *Queries) GetSpotifyCache(ctx context.Context, trackID string) (SpotifyCache, error) {
	row := q.db.QueryRowContext(ctx, getSpotifyCache, trackID)
	var i SpotifyCache
	err := row.Scan(
		&i.ID,
		&i.TrackID,
		&i.TrackName,
		&i.ArtistName,
		&i.ArtistID,
		&i.CoverArtUrl,
		&i.AudioPreviewUrl,
	)
	return i, err
}

const getYoutubeCache = `-- name: GetYoutubeCache :one
SELECT id, yt_id, thumb, title, author, author_url FROM youtube_cache
WHERE yt_id = ? LIMIT 1
`

func (q *Queries) GetYoutubeCache(ctx context.Context, ytID string) (YoutubeCache, error) {
	row := q.db.QueryRowContext(ctx, getYoutubeCache, ytID)
	var i YoutubeCache
	err := row.Scan(
		&i.ID,
		&i.YtID,
		&i.Thumb,
		&i.Title,
		&i.Author,
		&i.AuthorUrl,
	)
	return i, err
}

const listMediafiles = `-- name: ListMediafiles :many
SELECT id, file_name, file_path, file_extention, file_type, thumbnail, uploaded_at FROM mediafiles
ORDER BY uploaded_at DESC
`

func (q *Queries) ListMediafiles(ctx context.Context) ([]Mediafile, error) {
	rows, err := q.db.QueryContext(ctx, listMediafiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Mediafile
	for rows.Next() {
		var i Mediafile
		if err := rows.Scan(
			&i.ID,
			&i.FileName,
			&i.FilePath,
			&i.FileExtention,
			&i.FileType,
			&i.Thumbnail,
			&i.UploadedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPublicPosts = `-- name: ListPublicPosts :many
SELECT created_at, tags, content, html, private FROM posts
WHERE private = 0
ORDER BY created_at DESC
`

func (q *Queries) ListPublicPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, listPublicPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.CreatedAt,
			&i.Tags,
			&i.Content,
			&i.Html,
			&i.Private,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeUnusedNavidromeCache = `-- name: RemoveUnusedNavidromeCache :execrows
DELETE FROM navidrome_cache
WHERE id IN (
    SELECT navidrome_cache.id
    FROM navidrome_cache
    LEFT JOIN posts ON instr(posts.content, navidrome_cache.share_id) > 0
    WHERE instr(posts.content, navidrome_cache.share_id) IS NULL
)
`

func (q *Queries) RemoveUnusedNavidromeCache(ctx context.Context) (int64, error) {
	result, err := q.db.ExecContext(ctx, removeUnusedNavidromeCache)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const removeUnusedSpotifyCache = `-- name: RemoveUnusedSpotifyCache :execrows
DELETE FROM spotify_cache
WHERE id IN (
    SELECT spotify_cache.id
    FROM spotify_cache
    LEFT JOIN posts ON instr(posts.content, spotify_cache.track_id) > 0
    WHERE instr(posts.content, spotify_cache.track_id) IS NULL
)
`

func (q *Queries) RemoveUnusedSpotifyCache(ctx context.Context) (int64, error) {
	result, err := q.db.ExecContext(ctx, removeUnusedSpotifyCache)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const removeUnusedYoutubeCache = `-- name: RemoveUnusedYoutubeCache :execrows
DELETE FROM youtube_cache
WHERE id IN (
    SELECT youtube_cache.id
    FROM youtube_cache
    LEFT JOIN posts ON instr(posts.content, youtube_cache.yt_id) > 0
    WHERE instr(posts.content, youtube_cache.yt_id) IS NULL
)
`

func (q *Queries) RemoveUnusedYoutubeCache(ctx context.Context) (int64, error) {
	result, err := q.db.ExecContext(ctx, removeUnusedYoutubeCache)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateMedia = `-- name: UpdateMedia :exec
UPDATE mediafiles
set file_name = ?
WHERE id = ?
`

type UpdateMediaParams struct {
	FileName string
	ID       int64
}

func (q *Queries) UpdateMedia(ctx context.Context, arg UpdateMediaParams) error {
	_, err := q.db.ExecContext(ctx, updateMedia, arg.FileName, arg.ID)
	return err
}

const updatePost = `-- name: UpdatePost :exec
UPDATE posts
set tags = ?,
    content = ?,
    html = ?,
    private = ?
WHERE created_at = ?
`

type UpdatePostParams struct {
	Tags      sql.NullString
	Content   string
	Html      string
	Private   int64
	CreatedAt int64
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) error {
	_, err := q.db.ExecContext(ctx, updatePost,
		arg.Tags,
		arg.Content,
		arg.Html,
		arg.Private,
		arg.CreatedAt,
	)
	return err
}
