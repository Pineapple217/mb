package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	ct "github.com/Pineapple217/mb/pkg/context"
	"github.com/Pineapple217/mb/pkg/database"
	"github.com/Pineapple217/mb/pkg/view"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Post(c echo.Context) error {
	auth := ct.IsAuth(c.Request().Context())
	private := 0
	if auth {
		private = 1
	}
	// parse xid
	xidStr := c.Param("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	// get post
	post, err := h.Q.GetPost(c.Request().Context(), xid)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	if post.Private > 0 && private == 0 {
		return echo.NotFoundHandler(c)
	}
	tags, err := h.Q.GetAllTags(c.Request().Context(), int64(private))
	if err != nil {
		panic(err)
	}
	if tags == nil {
		tags = []database.GetAllTagsRow{}
	}

	return render(c, view.Post(post, tags, h.Q))
}

func (h *Handler) CreatePost(c echo.Context) error {
	// parse tags
	tags := strings.TrimSpace(c.FormValue("tags"))
	tagsNS := sql.NullString{String: tags, Valid: tags != ""}
	// parse content
	content := strings.TrimSpace(c.FormValue("content"))
	if content == "" {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	privateStr := strings.TrimSpace(c.FormValue("private"))
	private := 0
	if privateStr == "on" {
		private = 1
	}
	// create post
	h.Q.CreatePost(c.Request().Context(), database.CreatePostParams{
		Tags:    tagsNS,
		Content: content,
		Private: int64(private),
	})

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) EditPostForm(c echo.Context) error {
	// parse xid
	xidStr := c.Param("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	// get post
	post, err := h.Q.GetPost(c.Request().Context(), xid)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	return render(c, view.EditPost(post))
}

func (h *Handler) EditPost(c echo.Context) error {
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
	privateStr := strings.TrimSpace(c.FormValue("private"))
	private := 0
	if privateStr == "on" {
		private = 1
	}
	// update in db
	err = h.Q.UpdatePost(c.Request().Context(), database.UpdatePostParams{
		Tags:      tagsNS,
		Content:   content,
		Private:   int64(private),
		CreatedAt: xid,
	})
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, "/?p="+xidStr+"##")
}

func (h *Handler) DeletePostForm(c echo.Context) error {
	// parse xid
	xidStr := c.Param("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	// get post
	post, err := h.Q.GetPost(c.Request().Context(), xid)
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	return render(c, view.DeletePost(post, h.Q))
}

func (h *Handler) DeletePost(c echo.Context) error {
	// parse xid
	xidStr := c.FormValue("xid")
	xid, err := strconv.ParseInt(xidStr, 10, 64)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}
	// delete post
	err = h.Q.DeletePost(c.Request().Context(), xid)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

const prefix = "tag_"

func (h *Handler) Posts(c echo.Context) error {
	auth := ct.IsAuth(c.Request().Context())
	private := 0
	if auth {
		private = 1
	}
	// TODO: could be cleaner
	pStr := c.QueryParam("p")
	if pStr != "" {
		p, err := strconv.ParseInt(pStr, 10, 64)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		page, err := h.Q.GetPostPage(c.Request().Context(), database.GetPostPageParams{
			ID: p,
			P:  int64(private),
		})
		if err != nil {
			return err
		}
		posts, postCount, err := h.Q.QueryPost(
			c.Request().Context(),
			nil,
			"",
			private,
			int(page),
		)
		if err != nil {
			return err
		}

		tags, err := h.Q.GetAllTags(c.Request().Context(), int64(private))
		if err != nil || tags == nil {
			tags = []database.GetAllTagsRow{}
		}
		maxPage := (postCount - 1) / database.PostsPerPage
		nav := view.Nav(int(page), maxPage, "")

		return render(c, view.Posts(posts, tags, nav, p, h.Q))
	} else {
		qps := c.QueryParams()
		queryTags := make([]string, 0, len(qps))
		for k := range qps {
			if strings.HasPrefix(k, prefix) && qps[k][0] == "on" {
				queryTags = append(queryTags, strings.TrimPrefix(k, prefix))
			}
		}
		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			page = 0
		}
		posts, postCount, err := h.Q.QueryPost(
			c.Request().Context(),
			queryTags,
			c.QueryParam("search"),
			private,
			page,
		)
		if err != nil {
			return err
		}

		tags, err := h.Q.GetAllTags(c.Request().Context(), int64(private))
		if err != nil || tags == nil {
			tags = []database.GetAllTagsRow{}
		}
		maxPage := (postCount - 1) / database.PostsPerPage
		urlQuery := constructUrlQuery(c.QueryParam("search"), queryTags)
		nav := view.Nav(page, maxPage, urlQuery)

		return render(c, view.Posts(posts, tags, nav, -1, h.Q))
	}
}

// Displays the latest PUBLIC post
func (h *Handler) PostLatest(c echo.Context) error {
	auth := ct.IsAuth(c.Request().Context())
	private := 0
	if auth {
		private = 1
	}

	post, err := h.Q.GetPostLatest(c.Request().Context())
	if errors.Is(err, sql.ErrNoRows) {
		return NotFoundMsg(c, "No posts available yet")
	}
	if err != nil {
		return err
	}

	tags, err := h.Q.GetAllTags(c.Request().Context(), int64(private))
	if err != nil {
		return err
	}
	return render(c, view.Post(post, tags, h.Q))
}
