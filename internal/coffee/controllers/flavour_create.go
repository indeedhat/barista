package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type createFlavourProfileRequest struct {
	Name string `json:"name" validate:"required"`
}

type createFlavoursData struct {
	ui.PageData
	Flavours []coffee.FlavourProfile
	Open     bool
}

func (c Controller) CreateFlavourProfile(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := createFlavoursData{PageData: ui.NewPageData("Flavours", "flavours", user)}
	pageData.Flavours = c.repo.IndexFlavourProfiles()
	pageData.Open = true

	var req createFlavourProfileRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		ui.RenderUser(rw, r, pageData)
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create flavour")
		ui.RenderUser(rw, r, pageData)
		return
	}

	flavour := coffee.FlavourProfile{
		Name: req.Name,
	}

	if err := c.repo.SaveFlavourProfile(&flavour); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create flavour")
		ui.RenderUser(rw, r, pageData)
		return
	}

	pageData.Flavours = c.repo.IndexFlavourProfiles()
	pageData.Open = false
	pageData.Form = createFlavourProfileRequest{}

	ui.Toast(rw, ui.Success, "Flavour created")
	ui.RenderUser(rw, r, pageData)
}

type createFlavourFromComponentRequest struct {
	Name     string `json:"new_flavour" validate:"required"`
	Existing []uint `json:"flavours" validate:"required"`
}

func (c Controller) CreateFlavourFromComponent(rw http.ResponseWriter, r *http.Request) {
	comData := ui.NewComponentData("flavours-input", ui.ComponentData{
		"flavours": c.repo.IndexFlavourProfiles(),
	})

	var req createFlavourFromComponentRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create flavour")
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	comData["existing"] = req.Existing

	flavour := coffee.FlavourProfile{
		Name: req.Name,
	}

	if err := c.repo.SaveFlavourProfile(&flavour); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create flavour")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	comData["flavours"] = c.repo.IndexFlavourProfiles()
	comData["existing"] = append(comData["existing"].([]uint), flavour.ID)

	ui.Toast(rw, ui.Success, "Flavour created")
	ui.RenderComponent(rw, comData)
}
