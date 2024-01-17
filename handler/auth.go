package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/Pineapple217/mb/auth"
	"github.com/Pineapple217/mb/view"
	"github.com/labstack/echo/v4"
)

func AuthForm(c echo.Context) error {

	return render(c, view.AuthForm())
}

func Auth(c echo.Context) error {
	pw := c.FormValue("auth")
	if pw != auth.SecretPassword {
		return c.Redirect(http.StatusSeeOther, "/auth")
	}
	// TODO: replace unix time with crypto safe random
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d-%s", time.Now().Unix(), auth.SecretPassword)))
	auth.SecretCookie = hash
	cookie := http.Cookie{
		Name:     "auth",
		SameSite: http.SameSiteStrictMode,
		Value:    hex.EncodeToString(hash[:]),
	}
	c.SetCookie(&cookie)
	return c.Redirect(http.StatusSeeOther, "/")
}
