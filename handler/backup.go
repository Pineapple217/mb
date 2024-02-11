package handler

import (
	"log/slog"
	"net/http"

	"github.com/Pineapple217/mb/database"
	"github.com/labstack/echo/v4"
)

func CreateBackup(c echo.Context) error {
	queries := database.GetQueries()
	err := queries.Backup(c.Request().Context())
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// TODO download backups menu
