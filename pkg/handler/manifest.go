package handler

import (
	"encoding/json"
	"net/http"

	s "github.com/Pineapple217/mb/pkg/static"
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

func createManifest() []byte {
	manifest := IconSet{
		Icons: []Icon{
			{Src: s.StaticMap["/static/icon-192.png"], Type: "image/png", Sizes: "192x192"},
			{Src: s.StaticMap["/static/icon-512.png"], Type: "image/png", Sizes: "512x512"},
		},
	}
	jsonData, err := json.Marshal(manifest)
	if err != nil {
		panic(err)
	}
	return jsonData
}

func (h *Handler) Manifest(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/manifest+json")
	c.Response().Header().Add("Cache-Control", "public, max-age=3600, must-revalidate")
	return c.JSONBlob(http.StatusOK, h.manifest)
}
