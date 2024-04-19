package handler

import (
	"net/http"
	"strings"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/labstack/echo/v4"
)

const robotFile = `User-agent: *
Allow: /
Allow: /media/t/*
Allow: /media/*$
Allow: /post/*
Allow: /?page=*
Disallow: /?*$
Disallow: /backup
Disallow: /auth
Disallow: /media
Disallow: /media/*/*
Disallow: /post/*/*

User-agent: Googlebot-Image
Disallow: /

Sitemap: ##SITEMAP##
`

func (h *Handler) RobotTxt(c echo.Context) error {
	sitemap := config.Host + "/index.xml"
	r := strings.Replace(robotFile, "##SITEMAP##", sitemap, -1)
	return c.String(http.StatusOK, r)
}
