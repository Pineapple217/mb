package handler

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, component templ.Component) error {
	ctx := c.Request().Context()
	return component.Render(ctx, c.Response().Writer)
}

func constructUrlQuery(search string, tags []string) string {
	var sb strings.Builder
	if search != "" {
		sb.WriteString("&search=")
		sb.WriteString(search)
	}
	if len(tags) != 0 {
		sb.WriteString("&")
		for i, tag := range tags {
			sb.WriteString("tag_")
			sb.WriteString(tag)
			sb.WriteString("=on")
			if i+1 < len(tags) {
				sb.WriteString("&")
			}
		}
	}
	return sb.String()
}
