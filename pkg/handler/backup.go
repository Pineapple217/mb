package handler

import (
	"log/slog"
	"net/http"

	"github.com/Pineapple217/mb/pkg/backup"
	"github.com/Pineapple217/mb/pkg/view"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateBackup(c echo.Context) error {
	err := backup.Backup(c.Request().Context(), h.Q)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/backup")
}

func (h *Handler) Backups(c echo.Context) error {
	b, err := backup.GetAllBackups()
	if err != nil {
		return err
	}
	return render(c, view.Backups(b))
}
