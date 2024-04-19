package handler

import "github.com/Pineapple217/mb/pkg/database"

type Handler struct {
	Q *database.Queries
}

func NewHandler(q *database.Queries) *Handler {
	return &Handler{
		Q: q,
	}
}
