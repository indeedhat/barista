package coffee

import (
	"fmt"
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) ViewRoasters(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	roasters := c.repo.IndexRoastersForUser(user)

	pageData := ui.NewPageData("Roasters", "roasters", user)
	pageData.Data["Roasters"] = roasters
	ui.RenderUser(rw, r, pageData)
}

type upsertRoasterRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"omitempty,url"`
}

func (r upsertRoasterRequest) apply(roaster *Roaster) {
	roaster.Name = r.Name
	roaster.Description = r.Description
	roaster.URL = r.URL
}

func (c Controller) CreateRoaster(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Roasters", "roasters", user)
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)
	pageData.Data["open"] = true
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	var req upsertRoasterRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create roaster")
		return
	}

	roaster := Roaster{User: *user}
	req.apply(&roaster)

	if err := c.repo.SaveRoaster(&roaster); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create roaster")
		return
	}

	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)
	pageData.Data["open"] = false
	pageData.Form = upsertRoasterRequest{}

	ui.Toast(rw, ui.Success, "Roaster created")
}

func (c Controller) ViewRoaster(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Roaster Not Found", "404", user))
		return
	}

	roaster, err := c.repo.FindRoaster(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Roaster Not Found", "404", user))
		return
	}

	pageData := ui.NewPageData("Roaster", "roaster", user)
	pageData.Data["Roaster"] = roaster
	ui.RenderUser(rw, r, pageData)
}

func (c Controller) UpdateRoaster(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Roaster", "roaster", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	roaster, err := c.repo.FindRoaster(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	pageData.Title = roaster.Name
	pageData.Data["Roaster"] = roaster

	var req upsertRoasterRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create roaster")
		return
	}

	req.apply(roaster)
	if err := c.repo.SaveRoaster(roaster); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create roaster")
		return
	}

	pageData.Title = roaster.Name
	ui.Toast(rw, ui.Success, "Roaster Updated")
}

func (c Controller) UpdateRoasterImage(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Roaster", "roaster", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	roaster, err := c.repo.FindRoaster(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	pageData.Title = roaster.Name
	pageData.Data["Roaster"] = roaster

	savePath, err := server.UploadFile(r, "image", fmt.Sprint(RoasterImagePath, roaster.ID), &server.UploadProps{
		Ext:  []string{".jpg", ".jpeg", ".png"},
		Mime: []string{"image/png", "image/jpeg"},
	})
	if err != nil {
		ui.Toast(rw, ui.Warning, "Failed to upload image")
		return
	}

	if savePath != "" {
		roaster.Icon = savePath[5:]
		if err := c.repo.SaveRoaster(roaster); err != nil {
			ui.Toast(rw, ui.Warning, "Failed to save image")
			return
		}
	}
}

func (c Controller) DeleteRoaster(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Roaster", "roaster", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	roaster, err := c.repo.FindRoaster(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	pageData.Title = roaster.Name
	pageData.Data["Roaster"] = roaster

	if len(roaster.Coffees) > 0 {
		ui.Toast(rw, ui.Warning, "Roaster cannot be deleted while it still has coffees")
		return
	}

	if err := c.repo.DeleteRoaster(roaster); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete roaster")
		return
	}

	ui.Toast(rw, ui.Success, "Roaster deleted")
	server.Redirect(rw, r, "/roasters")
}
