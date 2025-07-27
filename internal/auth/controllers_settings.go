package auth

import (
	"net/http"

	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) ViewSettings(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	ui.RenderUser(rw, r, ui.NewPageData("User Settings", "user-settings", user))
}

type changePasswordRequest struct {
	Existing        string `json:"existing" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_conf" validate:"required"`
}

func (c Controller) ChangePassword(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	pageData := ui.NewPageData("User Settings", "user-settings", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	var req changePasswordRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	pageData.Form = req

	if err := server.ValidateRequest(req); err != nil {
		pageData.FieldErrors = server.ExtractFIeldErrors(err).Fields
		ui.Toast(rw, ui.Warning, "Change Password Failed")
		return
	}

	if u, _ := c.repo.FindUserByLogin(user.Name, req.Password); u == nil || user.ID != u.ID {
		pageData.FieldErrors["existing"] = []string{"Does not match current password"}
		ui.Toast(rw, ui.Warning, "Change Password Failed")
		return
	}

	if req.Password != req.PasswordConfirm {
		pageData.Form = req
		pageData.FieldErrors["password"] = append(pageData.FieldErrors["password"], "Passwords do not match")
		pageData.FieldErrors["password_conf"] = append(pageData.FieldErrors["password_conf"], "Passwords do not match")
		ui.Toast(rw, ui.Warning, "Change Password Failed")
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		ui.Toast(rw, ui.Warning, "Server Error, Please try again")
		return
	}

	if err := c.repo.UpdateUserPassword(user, hash); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to save new password")
		return
	}

	ui.Toast(rw, ui.Success, "Password Updated")
}
