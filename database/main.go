package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "embed"

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
	// TODO string building via byte buffer
	var filter string
	if search != "" {
		filter += fmt.Sprintf(" where (LOWER(content) glob '*%s*' collate nocase)",
			strings.ToLower(search))
	}

	if len(tags) != 0 {
		if search != "" {
			filter += " and ("
		} else {
			filter += " where ("
		}

		for i, tag := range tags {
			filter += fmt.Sprintf("tags like '%%%s%%' escape '\\'", tag)

			if i+1 < len(tags) {
				filter += " and "
			} else {
				filter += ")"
			}
		}
	}
	limit := fmt.Sprintf(" limit %d offset %d", PostsPerPage, PostsPerPage*page)
	query := queryPosts + filter + queryPostsOrder + limit
	countQuery := countPosts + filter

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
