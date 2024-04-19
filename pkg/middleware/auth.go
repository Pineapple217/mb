package middleware

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Pineapple217/mb/pkg/auth"
	ct "github.com/Pineapple217/mb/pkg/context"
	"github.com/Pineapple217/mb/pkg/view"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isAuth := false
		authCookie, err := c.Cookie("auth")

		if err != nil {
			isAuth = false
		} else {
			if authCookie.Value == auth.SecretCookie {
				isAuth = true
			}
		}

		ctx := context.WithValue(c.Request().Context(), ct.AuthContextKey, isAuth)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !ct.IsAuth(c.Request().Context()) {
			c.Response().Writer.WriteHeader(http.StatusUnauthorized)
			r := url.QueryEscape(c.Request().URL.Path)
			return view.AuthRedirect(r).Render(c.Request().Context(), c.Response().Writer)
		} else {
			return next(c)
		}
	}
}
