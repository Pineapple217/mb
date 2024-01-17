CREATE TABLE IF NOT EXISTS posts (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at INTEGER NOT NULL,
  tags       TEXT,
  content    TEXT    NOT NULL
);


CREATE TABLE IF NOT EXISTS spotify_cache (
  id                INTEGER PRIMARY KEY AUTOINCREMENT,
  track_id          TEXT NOT NULL,
  track_name        TEXT NOT NULL,
  artist_name       TEXT NOT NULL,
  artist_id         TEXT NOT NULL,
  cover_art_url     TEXT NOT NULL,
  audio_preview_url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS yt_thumb_cache (
  id       INTEGER PRIMARY KEY AUTOINCREMENT,
  yt_id    TEXT    NOT NULL,
  yt_thumb TEXT    NOT NULL
);

