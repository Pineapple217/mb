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

	"github.com/Pineapple217/mb/config"
	_ "github.com/mattn/go-sqlite3"
)

const PostsPerPage int = 25

var (
	//go:embed schema.sql
	ddl     string
	queries *Queries
)

func Init(databaseSource string) {
	ctx := context.Background()

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.Mkdir(backupDir, 0755)
		if err != nil {
			slog.Error("Error creating directory:",
				"error",
				err,
			)

		}
		slog.Debug("Created backup directory", "directory", backupDir)
	}

	db, err := sql.Open("sqlite3", databaseSource)
	if err != nil {
		panic(err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}
	queries = New(db)
}

func GetQueries() *Queries {
	return queries
}

const queryPosts = `SELECT id, created_at, tags, content FROM posts`
const countPosts = `SELECT COUNT(*) FROM posts`
const queryPostsOrder = ` ORDER BY created_at DESC`

func (q *Queries) QueryPost(ctx context.Context, tags []string, search string, page int) ([]Post, int, error) {
	var filter strings.Builder
	if search != "" {
		filter.WriteString(fmt.Sprintf(" where (LOWER(content) glob '*%s*' collate nocase)",
			strings.ToLower(search)))
	}

	if len(tags) != 0 {
		if search != "" {
			filter.WriteString(" and (")
		} else {
			filter.WriteString(" where (")
		}

		for i, tag := range tags {
			filter.WriteString(fmt.Sprintf("tags like '%%%s%%' escape '\\'", tag))

			if i+1 < len(tags) {
				filter.WriteString(" and ")
			} else {
				filter.WriteString(")")
			}
		}
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
			&i.ID,
			&i.CreatedAt,
			&i.Tags,
			&i.Content,
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

const queryBackup = `vacuum into '%s.db'`
const backupDir = config.DataDir + `/backups`

func (q *Queries) Backup(ctx context.Context) error {
	timeString := time.Now().Format("2006-01-02_15-04-05")
	file := fmt.Sprintf("%s/backup_%s", backupDir, timeString)
	_, err := q.db.ExecContext(ctx, fmt.Sprintf(queryBackup, file))
	return err
}
