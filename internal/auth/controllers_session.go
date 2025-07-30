package auth

import (
	"net/http"
	"time"

	"github.com/indeedhat/barista/internal/server"
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
		Register: envEnableRegister.Get(),
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
		Register: envEnableRegister.Get(),
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

	jwt, err := GenerateUserJwt(user.ID, user.Name, uint8(user.Level), user.JwtKillSwitch)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Failed to process login")
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     sessionCookie,
		Value:    jwt,
		HttpOnly: true,
		Domain:   r.URL.Host,
		Path:     "/",
		MaxAge:   86400 * 30,
	})

	ui.Redirect(rw, "/")
}

func (c Controller) ViewRegister(rw http.ResponseWriter, r *http.Request) {
	if !envEnableRegister.Get() {
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
	if !envEnableRegister.Get() {
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

	hash, err := hashPassword(req.Password)
	if err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		ui.Toast(rw, ui.Warning, "Server Error, Please try again")
		return
	}

	user := User{
		Name:          req.Name,
		Password:      string(hash),
		Level:         LevelMember,
		JwtKillSwitch: time.Now().Unix(),
	}

	if err := c.repo.SaveUser(&user); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	ui.Redirect(rw, "/login?register=true")
}
