package handler

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, component templ.Component) error {
	ctx := c.Request().Context()
	return component.Render(ctx, c.Response().Writer)
}

func constructUrlQuery(search string, tags []string) string {
	var buffer bytes.Buffer
	if search != "" {
		buffer.WriteString("&search=")
		buffer.WriteString(search)
	}
	if len(tags) != 0 {
		buffer.WriteString("&")
		for i, tag := range tags {
			buffer.WriteString("tag_")
			buffer.WriteString(tag)
			buffer.WriteString("=on")
			if i+1 < len(tags) {
				buffer.WriteString("&")
			}
		}
	}
	return buffer.String()
}
