package auth_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) Logout(rw http.ResponseWriter, r *http.Request) {
	http.SetCookie(rw, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		HttpOnly: true,
		Domain:   r.URL.Host,
		Path:     "/",
		MaxAge:   0,
	})
	ui.Redirect(rw, "/login")
}
