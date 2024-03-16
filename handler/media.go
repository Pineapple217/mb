package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pineapple217/mb/database"
	"github.com/Pineapple217/mb/media"
	"github.com/Pineapple217/mb/view"
	"github.com/labstack/echo/v4"
)

func Media(c echo.Context) error {
	queries := database.GetQueries()
	mediaFiles, err := queries.ListMediafiles(c.Request().Context())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return render(c, view.Media(mediaFiles))
}

func Mediafile(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	queries := database.GetQueries()
	mediafile, err := queries.GetMediafile(c.Request().Context(), id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return render(c, view.Mediafile(mediafile))
}

// TODO: audio and video
func MediaUpload(c echo.Context) error {
	file, err := c.FormFile("upload")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = media.SaveFile(
		c.Request().Context(),
		file,
		strings.TrimSpace(c.FormValue("name")),
	)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/media")
}

func Thumbnail(c echo.Context) error {
	name := c.Param("name")

	queries := database.GetQueries()
	t, err := queries.GetMediaThunbnail(
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

func MediaDeleteForm(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	queries := database.GetQueries()
	media, err := queries.GetMediafile(c.Request().Context(), id)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	return render(c, view.DeleteMedia(media))
}

func MediaDelete(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	queries := database.GetQueries()
	err = queries.DeleteMediafile(c.Request().Context(), id)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, "/media")
}

func MediaEditForm(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	queries := database.GetQueries()
	media, err := queries.GetMediafile(c.Request().Context(), id)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	return render(c, view.EditMedia(media))
}

func MediaEdit(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	queries := database.GetQueries()
	err = queries.UpdateMedia(c.Request().Context(), database.UpdateMediaParams{
		FileName: strings.TrimSpace(c.FormValue("name")),
		ID:       id,
	})
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, "/media/"+idStr)
}
