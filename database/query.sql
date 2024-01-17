-- name: GetPost :one
SELECT * FROM posts
WHERE created_at = ? LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY created_at DESC;

-- name: CreatePost :one
INSERT INTO posts (
  created_at, tags, content
) VALUES (
  strftime('%s', 'now'), ?, ?
)
RETURNING *;

-- name: UpdatePost :exec
UPDATE posts
set tags = ?,
    content = ?
WHERE created_at = ?;

-- name: DeletePost :exec
DELETE FROM posts
WHERE created_at = ?;

-- name: GetPostCount :one
SELECT COUNT(*)
FROM posts;

-- name: CreateSpotifyCache :one
INSERT INTO spotify_cache (
  track_id, track_name, artist_name, artist_id, cover_art_url, audio_preview_url
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetSpotifyCache :one
SELECT * FROM spotify_cache
WHERE track_id = ? LIMIT 1;

-- name: CreateTYThumbCache :one
INSERT INTO yt_thumb_cache (
  yt_id, yt_thumb
) VALUES (
  ?, ?
)
RETURNING *;

-- name: GetYoutubeCache :one
SELECT * FROM yt_thumb_cache
WHERE yt_id = ? LIMIT 1;

-- name: GetTagsCount :one
WITH split(tag, tags_remaining) AS (
  -- Initial query
  SELECT 
    '',
    tags || ' '
  FROM posts
  UNION ALL
  SELECT
    trim(substr(tags_remaining, 0, instr(tags_remaining, ' '))),
    substr(tags_remaining, instr(tags_remaining, ' ') + 1)
  FROM split
  WHERE tags_remaining != ''
)
SELECT COUNT(DISTINCT tag) AS unique_tag_count
FROM split
WHERE tag != '';

-- name: GetAllTags :many
WITH split (
    tag,
    tags_remaining
)
AS (-- Initial query
    SELECT '',
           tags || ' '
      FROM posts
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
ORDER BY tag_count DESC;
