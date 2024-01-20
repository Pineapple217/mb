package middleware

import (
	"context"

	ct "github.com/Pineapple217/mb/context"
	"github.com/Pineapple217/mb/database"
	"github.com/labstack/echo/v4"
)

func Stats(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		queries := database.GetQueries()
		// post count
		postCount, err := queries.GetPostCount(c.Request().Context())
		if err != nil {
			postCount = -1
		}
		ctx := context.WithValue(c.Request().Context(), ct.PostCountContextKey, postCount)

		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
