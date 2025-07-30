package brewer_controllers

import (
	"fmt"
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/brewer"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type updateBrewerRequest struct {
	Name        string `json:"name" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	ModelNumber string `json:"model" validate:"required"`
}

type updateBrewerData struct {
	ui.PageData
	Brewer *brewer.Brewer
}

func (c Controller) UpdateBrewer(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := updateBrewerData{PageData: ui.NewPageData("Brewer", "brewer", user)}
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return
	}

	brewer, err := c.repo.FindBrewer(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return
	}

	pageData.Title = brewer.Name
	pageData.Brewer = brewer

	var req updateBrewerRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update brewer")
		return
	}

	brewer.Name = req.Name
	brewer.Brand = req.Brand
	brewer.ModelNumber = req.ModelNumber

	if err := c.repo.SaveBrewer(brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update brewer")
		return
	}

	pageData.Title = brewer.Name
	ui.Toast(rw, ui.Success, "Brewer Updated")
}

func (c Controller) UpdateBrewerImage(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := updateBrewerData{PageData: ui.NewPageData("Brewer", "brewer", user)}
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return
	}

	brewer, err := c.repo.FindBrewer(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return
	}

	pageData.Title = brewer.Name
	pageData.Brewer = brewer

	savePath, err := server.UploadFile(r, "image", fmt.Sprint(BrewerImagePath, brewer.ID), &server.UploadProps{
		Ext:  []string{".jpg", ".jpeg", ".png"},
		Mime: []string{"image/png", "image/jpeg"},
	})
	if err != nil {
		ui.Toast(rw, ui.Warning, "Failed to upload image")
		return
	}

	if savePath != "" {
		brewer.Icon = savePath[5:]
		if err := c.repo.SaveBrewer(brewer); err != nil {
			ui.Toast(rw, ui.Warning, "Failed to save image")
			return
		}
	}

	ui.Toast(rw, ui.Success, "Image Updated")
}
