package middleware

import (
	"context"

	ct "github.com/Pineapple217/mb/pkg/context"
	"github.com/Pineapple217/mb/pkg/database"
	"github.com/labstack/echo/v4"
)

func Stats(next echo.HandlerFunc, q *database.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := ct.IsAuth(c.Request().Context())
		private := 0
		if auth {
			private = 1
		}
		// post count
		postCount, err := q.GetPostCount(c.Request().Context(), int64(private))
		if err != nil {
			postCount = -1
		}
		ctx := context.WithValue(c.Request().Context(), ct.PostCountContextKey, postCount)

		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
