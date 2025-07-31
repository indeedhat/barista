package auth_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/cookie"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type loginPageData struct {
	ui.PageData
	Register bool
}

func (c Controller) ViewLogin(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("register") != "" {
		ui.Toast(rw, ui.Success, "User Created, you may now login")
	}

	ui.RenderGuest(rw, r, loginPageData{
		PageData: ui.NewPageData("Login", "login"),
		Register: auth.EnvEnableRegister.Get(),
	})
}

type loginRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login handles user login attempts
func (c Controller) Login(rw http.ResponseWriter, r *http.Request) {
	pageData := loginPageData{
		PageData: ui.NewPageData("Login", "login"),
		Register: auth.EnvEnableRegister.Get(),
	}
	defer func() {
		ui.RenderGuest(rw, r, pageData)
	}()

	var req loginRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		pageData.FieldErrors = server.ExtractFIeldErrors(err).Fields
		pageData.Form = req
		ui.Toast(rw, ui.Warning, "Login failed")
		return
	}

	user, err := c.repo.FindUserByLogin(req.Name, req.Password)
	if user == nil {
		pageData.Form = req
		ui.Toast(rw, ui.Warning, "Login failed")
		return
	}

	jwt, err := auth.GenerateUserJwt(user.ID, user.Name, uint8(user.Level), user.JwtKillSwitch)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Failed to process login")
		return
	}

	cookie.Set(rw, r, cookie.SessionKey, jwt)

	ui.Redirect(rw, "/")
}
