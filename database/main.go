package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed schema.sql
	ddl     string
	queries *Queries
)

func Init(databaseSource string) {
	ctx := context.Background()
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
const queryPostsOrder = ` ORDER BY created_at DESC`

func (q *Queries) QueryPost(ctx context.Context, tags []string, search string) ([]Post, error) {
	query := queryPosts

	if search != "" {
		query += fmt.Sprintf(" where (content glob '*%s*' collate nocase)", search)
	}

	if len(tags) != 0 {
		if search != "" {
			query += " and ("
		} else {
			query += " where ("
		}

		for i, tag := range tags {
			// tag := strings.NewReplacer("_", "\\_", "%", "\\%").Replace(t)
			query += fmt.Sprintf("tags like '%%%s%%' escape '\\'", tag)

			if i+1 < len(tags) {
				query += " and "
			} else {
				query += ")"
			}
		}
	}
	query += queryPostsOrder

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
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
