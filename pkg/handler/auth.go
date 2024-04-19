package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/Pineapple217/mb/pkg/auth"
	"github.com/Pineapple217/mb/pkg/view"
	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthForm(c echo.Context) error {
	r, err := url.QueryUnescape(c.QueryParam("r"))
	if err != nil || r == "" {
		r = "/"
	}
	return render(c, view.AuthForm(r))
}

// TODO: keep redirect afther giving wrong password

func (h *Handler) Auth(c echo.Context) error {
	pw := c.FormValue("auth")
	if pw != auth.SecretPassword {
		return c.Redirect(http.StatusSeeOther, "/auth")
	}
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	encoded := base64.URLEncoding.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(encoded + auth.SecretPassword))
	hashStr := base64.RawStdEncoding.EncodeToString(hash[:])
	auth.SecretCookie = hashStr
	cookie := http.Cookie{
		Name:     "auth",
		SameSite: http.SameSiteStrictMode,
		Value:    hashStr,
	}
	c.SetCookie(&cookie)
	r := c.FormValue("redirect")
	return c.Redirect(http.StatusSeeOther, r)
}
