package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pineapple217/mb/pkg/database"
	"github.com/Pineapple217/mb/pkg/media"
	"github.com/Pineapple217/mb/pkg/view"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Media(c echo.Context) error {
	mediaFiles, err := h.Q.ListMediafiles(c.Request().Context())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return render(c, view.Media(mediaFiles))
}

func (h *Handler) Mediafile(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	mediafile, err := h.Q.GetMediafile(c.Request().Context(), id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return render(c, view.Mediafile(mediafile))
}

func (h *Handler) MediaUpload(c echo.Context) error {
	file, err := c.FormFile("upload")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = media.SaveFile(
		c.Request().Context(),
		h.Q,
		file,
		strings.TrimSpace(c.FormValue("name")),
	)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/media")
}

func (h *Handler) Thumbnail(c echo.Context) error {
	name := c.Param("name")

	t, err := h.Q.GetMediaThunbnail(
		c.Request().Context(),
		name,
	)
	if err != nil {
		return err
	}

	c.Response().Header().Add("Cache-Control", "max-age=3600")
	c.Blob(http.StatusOK, "image/jpg", t)
	return nil
}

func (h *Handler) MediaDeleteForm(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	media, err := h.Q.GetMediafile(c.Request().Context(), id)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	return render(c, view.DeleteMedia(media))
}

func (h *Handler) MediaDelete(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	m, err := h.Q.GetMediafile(c.Request().Context(), id)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}
	err = media.DeleteFile(m.FilePath)
	if err != nil {
		return err
	}

	err = h.Q.DeleteMediafile(c.Request().Context(), id)
	if err != nil {
		slog.Warn("Deleted mediafile but could not delete db record", "err", err, "id", m.ID)
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/media")
}

func (h *Handler) MediaEditForm(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	media, err := h.Q.GetMediafile(c.Request().Context(), id)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	return render(c, view.EditMedia(media))
}

func (h *Handler) MediaEdit(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	err = h.Q.UpdateMedia(c.Request().Context(), database.UpdateMediaParams{
		FileName: strings.TrimSpace(c.FormValue("name")),
		ID:       id,
	})
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, "/media/"+idStr)
}
