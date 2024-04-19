package handler

import (
	"net/http"

	"github.com/Pineapple217/mb/pkg/view"
	"github.com/labstack/echo/v4"
)

func NotFound(c echo.Context) error {
	c.Response().Writer.WriteHeader(http.StatusNotFound)
	return render(c, view.NotFound())
}
