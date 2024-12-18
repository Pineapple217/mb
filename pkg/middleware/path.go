package middleware

import (
	"context"

	ct "github.com/Pineapple217/mb/pkg/context"
	"github.com/labstack/echo/v4"
)

func Path(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), ct.PathContextKey, c.Request().URL.RequestURI())

		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
