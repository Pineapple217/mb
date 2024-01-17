package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pineapple217/mb/database"
	"github.com/Pineapple217/mb/view"
	"github.com/labstack/echo/v4"
)

func Post(c echo.Context) error {
	// parse xid
	xidStr := c.Param("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	// get post
	queries := database.GetQueries()
	post, err := queries.GetPost(c.Request().Context(), xid)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	tags, err := queries.GetAllTags(c.Request().Context())
	if err != nil {
		panic(err)
	}
	if tags == nil {
		tags = []database.GetAllTagsRow{}
	}

	return render(c, view.Post(post, tags))
}

func CreatePost(c echo.Context) error {
	// parse tags
	tags := strings.TrimSpace(c.FormValue("tags"))
	tagsNS := sql.NullString{String: tags, Valid: tags != ""}
	// parse content
	content := strings.TrimSpace(c.FormValue("content"))
	if content == "" {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	// create post
	queries := database.GetQueries()
	queries.CreatePost(c.Request().Context(), database.CreatePostParams{
		Tags:    tagsNS,
		Content: content,
	})

	return c.Redirect(http.StatusSeeOther, "/")
}

func EditPostForm(c echo.Context) error {
	// parse xid
	xidStr := c.Param("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	// get post
	queries := database.GetQueries()
	post, err := queries.GetPost(c.Request().Context(), xid)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	return render(c, view.EditPost(post))
}

func EditPost(c echo.Context) error {
	// parse xid
	xidStr := c.FormValue("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}
	// parse tags
	tags := strings.TrimSpace(c.FormValue("tags"))
	tagsNS := sql.NullString{String: tags, Valid: tags != ""}
	// parse content
	content := strings.TrimSpace(c.FormValue("content"))
	if content == "" {
		return c.Redirect(http.StatusSeeOther, "/post/"+xidStr)
	}
	// update in db
	queries := database.GetQueries()
	err = queries.UpdatePost(c.Request().Context(), database.UpdatePostParams{
		Tags:      tagsNS,
		Content:   content,
		CreatedAt: xid,
	})
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, "/post/"+xidStr)
}

func DeletePostForm(c echo.Context) error {
	// parse xid
	xidStr := c.Param("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	// get post
	queries := database.GetQueries()
	post, err := queries.GetPost(c.Request().Context(), xid)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	return render(c, view.DeletePost(post))
}

func DeletePost(c echo.Context) error {
	// parse xid
	xidStr := c.FormValue("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}
	// delete post
	queries := database.GetQueries()
	err = queries.DeletePost(c.Request().Context(), xid)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func Posts(c echo.Context) error {
	const prefix = "tag_"
	qps := c.QueryParams()
	queryTags := make([]string, 0, len(qps))
	for k := range qps {
		if strings.HasPrefix(k, prefix) && qps[k][0] == "on" {
			queryTags = append(queryTags, strings.TrimPrefix(k, prefix))
		}
	}

	queries := database.GetQueries()
	posts, err := queries.QueryPost(c.Request().Context(), queryTags, c.QueryParam("search"))
	if err != nil {
		return err
	}

	tags, err := queries.GetAllTags(c.Request().Context())
	if err != nil {
		panic(err)
	}
	if tags == nil {
		tags = []database.GetAllTagsRow{}
	}

	return render(c, view.Posts(posts, tags))
}
