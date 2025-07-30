package auth_controllers

import (
	"net/http"
	"time"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) ViewRegister(rw http.ResponseWriter, r *http.Request) {
	if !auth.EnvEnableRegister.Get() {
		ui.Redirect(rw, "/")
		return
	}

	ui.RenderGuest(rw, r, ui.PageData{
		Title: "Register",
		Page:  "pages/register",
		Form:  registerRequest{},
	})
}

type registerRequest struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_conf" validate:"required"`
}

// Register handles user register attempts
func (c Controller) Register(rw http.ResponseWriter, r *http.Request) {
	if !auth.EnvEnableRegister.Get() {
		ui.Redirect(rw, "/")
		return
	}

	pageData := ui.NewPageData("Register", "register")
	pageData.Form = registerRequest{}
	defer func() {
		ui.RenderGuest(rw, r, pageData)
	}()

	var req registerRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		pageData.FieldErrors = server.ExtractFIeldErrors(err).Fields
		pageData.Form = registerRequest{Name: req.Name}
		ui.Toast(rw, ui.Warning, "Register failed")
		return
	}

	if user, _ := c.repo.FindUserByLogin(req.Name, req.Password); user != nil {
		pageData.Form = registerRequest{Name: req.Name}
		ui.Toast(rw, ui.Warning, "Name in use")
		return
	}

	if req.Password != req.PasswordConfirm {
		pageData.Form = req
		ui.Toast(rw, ui.Warning, "Passwords do not match")
		return
	}

	user := auth.User{
		Name:          req.Name,
		Password:      req.Password,
		Level:         auth.LevelMember,
		JwtKillSwitch: time.Now().Unix(),
	}

	if err := c.repo.SaveUser(&user); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	ui.Redirect(rw, "/login?register=true")
}
