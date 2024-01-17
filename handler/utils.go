package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, component templ.Component) error {
	ctx := c.Request().Context()
	return component.Render(ctx, c.Response().Writer)
}
