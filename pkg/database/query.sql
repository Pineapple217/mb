-- name: GetPost :one
SELECT * FROM posts
WHERE created_at = ? LIMIT 1;

-- name: GetPostLatest :one
SELECT * FROM posts
WHERE private = 0
ORDER BY created_at DESC LIMIT 1;

-- name: ListPublicPosts :many
SELECT * FROM posts
WHERE private = 0
ORDER BY created_at DESC;

-- name: CreatePost :one
INSERT INTO posts (
  created_at, tags, content, html, private
) VALUES (
  strftime('%s', 'now'), ?, ?, ?, ?
)
RETURNING *;

-- name: UpdatePost :exec
UPDATE posts
set tags = ?,
    content = ?,
    html = ?,
    private = ?
WHERE created_at = ?;

-- name: DeletePost :exec
DELETE FROM posts
WHERE created_at = ?;

-- name: GetPostCount :one
SELECT COUNT(*)
FROM posts
WHERE private <= ?;

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

-- name: CreateYoutubebCache :one
INSERT INTO youtube_cache (
  yt_id, thumb, title, author, author_url
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetYoutubeCache :one
SELECT * FROM youtube_cache
WHERE yt_id = ? LIMIT 1;

-- name: RemoveUnusedYoutubeCache :execrows
DELETE FROM youtube_cache
WHERE id IN (
    SELECT youtube_cache.id
    FROM youtube_cache
    LEFT JOIN posts ON instr(posts.content, youtube_cache.yt_id) > 0
    WHERE instr(posts.content, youtube_cache.yt_id) IS NULL
);

-- name: RemoveUnusedSpotifyCache :execrows
DELETE FROM spotify_cache
WHERE id IN (
    SELECT spotify_cache.id
    FROM spotify_cache
    LEFT JOIN posts ON instr(posts.content, spotify_cache.track_id) > 0
    WHERE instr(posts.content, spotify_cache.track_id) IS NULL
);

-- name: GetAllTags :many
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
ORDER BY tag_count DESC;

-- name: GetPostPage :one
SELECT 
CAST(
    CASE 
        WHEN EXISTS (SELECT 1 FROM posts WHERE posts.created_at = :id)
        THEN CEIL((SELECT COUNT(*) FROM posts WHERE (posts.created_at >= (SELECT posts.created_at FROM posts WHERE created_at = :id)) and posts.private <= :p)  / 25.0) - 1
        ELSE -1
    END AS INT
) AS page_number;


-- name: ListMediafiles :many
SELECT * FROM mediafiles
ORDER BY uploaded_at DESC;

-- name: CreateMediafile :one
INSERT INTO mediafiles (
  uploaded_at, file_name, file_path, file_type, file_extention, thumbnail
) VALUES (
  strftime('%s', 'now'), ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetMediaThunbnail :one
SELECT thumbnail FROM mediafiles
WHERE file_path = ? LIMIT 1;

-- name: GetMediafile :one
SELECT * FROM mediafiles
WHERE id = ? LIMIT 1;

-- name: DeleteMediafile :exec
DELETE FROM mediafiles
WHERE id = ?;

-- name: UpdateMedia :exec
UPDATE mediafiles
set file_name = ?
WHERE id = ?;