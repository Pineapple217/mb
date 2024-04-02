package handler

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Pineapple217/mb/config"
	"github.com/labstack/echo/v4"
)

func RobotTxt(c echo.Context) error {
	data, err := os.ReadFile("static/private/robot.txt")
	if err != nil {
		slog.Warn("could not server robot file", "error", err)
		return err
	}
	sitemap := config.Host + "/index.xml"
	r := strings.Replace(string(data), "##SITEMAP##", sitemap, -1)
	return c.String(http.StatusOK, r)
}
