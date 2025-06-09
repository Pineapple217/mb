package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	_ "embed"

	"github.com/Pineapple217/mb/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

const PostsPerPage int = 25

var (
	//go:embed schema.sql
	ddl string
)

func NewQueries(databaseSource string) *Queries {
	slog.Info("Creating database connection")
	ctx := context.Background()

	if _, err := os.Stat(config.BackupDir); os.IsNotExist(err) {
		err := os.Mkdir(config.BackupDir, 0755)
		if err != nil {
			slog.Error("Error creating directory:",
				"error",
				err,
			)

		}
		slog.Debug("Created backup directory", "directory", config.BackupDir)
	}

	db, err := sql.Open("sqlite3", databaseSource)
	if err != nil {
		panic(err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}
	queries := New(db)
	err = queries.Migrate(ctx)
	if err != nil {
		panic(err)
	}
	return queries
}

const queryPosts = `SELECT created_at, tags, content, html, private FROM posts`
const countPosts = `SELECT COUNT(*) FROM posts`
const queryPostsOrder = ` ORDER BY created_at DESC`

func (q *Queries) QueryPost(ctx context.Context, tags []string, search string, private int, page int) ([]Post, int, error) {
	var filter strings.Builder
	firstWhere := false
	if search != "" {
		filter.WriteString(fmt.Sprintf(" where (LOWER(content) glob '*%s*' collate nocase)",
			strings.ToLower(search)))
		firstWhere = true
	}

	if len(tags) != 0 {
		if firstWhere {
			filter.WriteString(" and (")
		} else {
			filter.WriteString(" where (")
			firstWhere = true
		}

		for i, tag := range tags {
			filter.WriteString(fmt.Sprintf("' ' || tags || ' ' like '%% %s %%' escape '\\'", tag))

			if i+1 < len(tags) {
				filter.WriteString(" and ")
			} else {
				filter.WriteString(")")
			}
		}
	}

	if private == 0 {
		if firstWhere {
			filter.WriteString(" and (")
		} else {
			filter.WriteString(" where (")
			firstWhere = true
		}
		filter.WriteString("private = 0)")
	}

	limit := fmt.Sprintf(" limit %d offset %d", PostsPerPage, PostsPerPage*page)
	query := queryPosts + filter.String() + queryPostsOrder + limit
	countQuery := countPosts + filter.String()

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	row := q.db.QueryRowContext(ctx, countQuery)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

const queryBackup = `vacuum into '%s'`

func (q *Queries) Backup(ctx context.Context) (string, error) {
	timeString := time.Now().Format("2006-01-02_15-04-05")
	file := fmt.Sprintf("%s/backup_%s.db", config.BackupDir, timeString)
	_, err := q.db.ExecContext(ctx, fmt.Sprintf(queryBackup, file))
	return file, err
}

func CreateDataDir() {
	if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
		err := os.Mkdir(config.DataDir, 0755)
		if err != nil {
			slog.Error("Failed to create directory",
				"dir", config.DataDir,
				"error", err,
			)
			panic(err)
		}
	}
}
