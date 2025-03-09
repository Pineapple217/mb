CREATE TABLE IF NOT EXISTS posts (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at INTEGER NOT NULL,
  tags       TEXT,
  content    TEXT    NOT NULL,
  html       TEXT    NOT NULL,
  private    INTEGER NOT NULL DEFAULT 0
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

CREATE TABLE IF NOT EXISTS youtube_cache (
  id       INTEGER PRIMARY KEY AUTOINCREMENT,
  yt_id    TEXT    NOT NULL,
  thumb    TEXT    NOT NULL,
  title    TEXT    NOT NULL,
  author   TEXT    NOT NULL,
  author_url TEXT  NOT NULL
);

CREATE TABLE IF NOT EXISTS mediafiles (
  id       INTEGER PRIMARY KEY AUTOINCREMENT,
  file_name TEXT NOT NULL,
  file_path TEXT NOT NULL,
  file_extention TEXT NOT NULL,
  file_type TEXT NOT NULL,
  thumbnail BLOB NOT NULL,
  uploaded_at TIMESTAMP NOT NULL
)
