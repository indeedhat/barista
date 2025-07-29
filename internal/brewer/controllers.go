package brewer

import (
	"fmt"
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

const BrewerImagePath = "data/uploads/brewer/"

type Controller struct {
	repo Repository
}

func NewController(repo Repository) Controller {
	return Controller{repo}
}

func (c Controller) ViewBrewers(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	pageData := ui.NewPageData("Brewers", "brewers", user)
	pageData.Data["Brewers"] = c.repo.IndexBrewersForUser(user)
	pageData.Data["BrewerTypes"] = types.Brewers
	pageData.Form = createBrewerRequest{}

	ui.RenderUser(rw, r, pageData)
}

type createBrewerRequest struct {
	Name        string `json:"name" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	ModelNumber string `json:"model" validate:"required"`
	Type        string `json:"type" validate:"required"`
}

func (r createBrewerRequest) apply(brewer *Brewer) {
	brewer.Name = r.Name
	brewer.Brand = r.Brand
	brewer.ModelNumber = r.ModelNumber
	brewer.Type = types.BrewerType(r.Type)
}

func (c Controller) CreateBrewer(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Brewers", "brewers", user)
	pageData.Data["open"] = true
	pageData.Data["BrewerTypes"] = types.Brewers
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	var req createBrewerRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	brewer := Brewer{User: *user}
	req.apply(&brewer)

	if err := c.repo.SaveBrewer(&brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create brewer")
		return
	}

	pageData.Data["Brewers"] = c.repo.IndexBrewersForUser(user)
	pageData.Data["open"] = false
	pageData.Form = createBrewerRequest{}

	ui.Toast(rw, ui.Success, "Brewer created")
}

func (c Controller) ViewBrewer(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Brewer Not Found", "404", user))
		return
	}

	brewer, err := c.repo.FindBrewer(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Brewer Not Found", "404", user))
		return
	}

	pageData := ui.NewPageData("Brewer", "brewer", user)
	pageData.Data["Brewer"] = brewer
	pageData.Form = updateBrewerRequest{}

	ui.RenderUser(rw, r, pageData)
}

type updateBrewerRequest struct {
	Name        string `json:"name" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	ModelNumber string `json:"model" validate:"required"`
}

func (r updateBrewerRequest) apply(brewer *Brewer) {
	brewer.Name = r.Name
	brewer.Brand = r.Brand
	brewer.ModelNumber = r.ModelNumber
}

func (c Controller) UpdateBrewerImage(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Brewer", "brewer", user)
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
	pageData.Data["Brewer"] = brewer

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

	pageData.Data["Brewers"] = c.repo.IndexBrewersForUser(user)

	ui.Toast(rw, ui.Success, "Image Updated")
}

func (c Controller) UpdateBrewer(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Brewer", "brewer", user)
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
	pageData.Data["Brewer"] = brewer

	var req updateBrewerRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update brewer")
		return
	}

	req.apply(brewer)
	if err := c.repo.SaveBrewer(brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update brewer")
		return
	}

	pageData.Title = brewer.Name
	ui.Toast(rw, ui.Success, "Brewer Updated")
}

func (c Controller) DeleteBrewer(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Brewer", "brewer", user)
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
	pageData.Data["Brewer"] = brewer

	if len(brewer.Baskets) > 0 {
		ui.Toast(rw, ui.Warning, "Brewer cannot be deleted while it still has baskets")
		return
	}

	if err := c.repo.DeleteBrewer(brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete brewer")
		return
	}

	ui.Toast(rw, ui.Success, "Brewer Deleted")
	server.Redirect(rw, r, "/brewers")
}

func (c Controller) BrewersSelect(rw http.ResponseWriter, r *http.Request) {
	drink := r.URL.Query().Get("drink")
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("brewers-select", ui.ComponentData{
		"Brewers": c.repo.IndexBrewersForUser(user, types.DrinkType(drink).Brewers()...),
		"value":   r.URL.Query().Get("value"),
	})

	ui.RenderComponent(rw, comData)
}
