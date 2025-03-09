package database

import (
	"context"
	"log/slog"
)

type TableInfo struct {
	CID          int64
	Name         string
	Type         string
	NotNull      bool
	DefaultValue any
	PrimaryKey   bool
}

func (q *Queries) GetTableInfo(ctx context.Context, table string) ([]TableInfo, error) {
	rows, err := q.db.QueryContext(ctx, "PRAGMA table_info("+table+");")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TableInfo
	for rows.Next() {
		var i TableInfo
		if err := rows.Scan(
			&i.CID,
			&i.Name,
			&i.Type,
			&i.NotNull,
			&i.DefaultValue,
			&i.PrimaryKey,
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

func (q *Queries) Migrate(ctx context.Context) error {
	tis, err := q.GetTableInfo(ctx, "posts")
	if err != nil {
		return err
	}
	htmlExits := false
	for _, ti := range tis {
		if ti.Name == "html" {
			htmlExits = true
		}
	}
	if !htmlExits { // Remove me for 1.0
		slog.Info("HTML column not found, adding column")
		_, err = q.db.ExecContext(ctx, "ALTER TABLE posts ADD COLUMN html TEXT NOT NULL DEFAULT 'ERROR NO HTML';")
		if err != nil {
			return err
		}
	}
	return nil
}
