package middleware

import (
	"context"
	"encoding/hex"
	"net/http"

	"github.com/Pineapple217/mb/auth"
	ct "github.com/Pineapple217/mb/context"
	"github.com/Pineapple217/mb/view"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isAuth := false
		authCookie, err := c.Cookie("auth")

		if err != nil {
			isAuth = false
		} else {
			if authCookie.Value == hex.EncodeToString(auth.SecretCookie[:]) {
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
			c.Response().Writer.WriteHeader(http.StatusForbidden)
			return view.AuthRedirect().Render(c.Request().Context(), c.Response().Writer)
		} else {
			return next(c)
		}
	}
}
