package coffee_controllers

import (
	"fmt"
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type updateRoasterRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"omitempty,url"`
}

type updateRoasterData struct {
	ui.PageData
	Roaster *coffee.Roaster
}

func (c Controller) UpdateRoaster(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := updateRoasterData{PageData: ui.NewPageData("Roaster", "roaster", user)}
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
	pageData.Roaster = roaster

	var req updateRoasterRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create roaster")
		return
	}

	roaster.Name = req.Name
	roaster.Description = req.Description
	roaster.URL = req.URL

	if err := c.repo.SaveRoaster(roaster); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create roaster")
		return
	}

	pageData.Title = roaster.Name
	ui.Toast(rw, ui.Success, "Roaster Updated")
}

func (c Controller) UpdateRoasterImage(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := updateRoasterData{PageData: ui.NewPageData("Roaster", "roaster", user)}
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
	pageData.Roaster = roaster

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
