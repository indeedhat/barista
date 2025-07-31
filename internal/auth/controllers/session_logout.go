package auth_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/cookie"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) Logout(rw http.ResponseWriter, r *http.Request) {
	cookie.Delete(rw, r, cookie.SessionKey)
	ui.Redirect(rw, "/login")
}
