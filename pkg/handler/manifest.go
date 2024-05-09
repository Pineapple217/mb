package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Pineapple217/mb/pkg/static"
	"github.com/labstack/echo/v4"
)

type Icon struct {
	Src   string `json:"src"`
	Type  string `json:"type"`
	Sizes string `json:"sizes"`
}

type IconSet struct {
	Icons []Icon `json:"icons"`
}

var manfest = IconSet{
	Icons: []Icon{
		{Src: static.Icon192Png, Type: "image/png", Sizes: "192x192"},
		{Src: static.Icon512Png, Type: "image/png", Sizes: "512x512"},
	},
}

func (h *Handler) Manifest(c echo.Context) error {
	jsonData, err := json.Marshal(manfest)
	if err != nil {
		return err
	}
	c.Response().Header().Set("Content-Type", "application/manifest+json")
	return c.String(http.StatusOK, string(jsonData))
}
