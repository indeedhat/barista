package brewer_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/brewer"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

type createBrewerRequest struct {
	Name        string `json:"name" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	ModelNumber string `json:"model" validate:"required"`
	Type        string `json:"type" validate:"required"`
}

type createBrewerData struct {
	ui.PageData
	BrewerTypes []types.BrewerType
	Brewers     []brewer.Brewer
	Open        bool
}

func (c Controller) CreateBrewer(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := createBrewerData{PageData: ui.NewPageData("Brewers", "brewers", user)}
	pageData.Open = true
	pageData.BrewerTypes = types.Brewers
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

	brewer := brewer.Brewer{
		Name:        req.Name,
		Brand:       req.Brand,
		ModelNumber: req.ModelNumber,
		Type:        types.BrewerType(req.Type),
		User:        *user,
	}

	if err := c.repo.SaveBrewer(&brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create brewer")
		return
	}

	pageData.Brewers = c.repo.IndexBrewersForUser(user)
	pageData.Open = false
	pageData.Form = createBrewerRequest{}

	ui.Toast(rw, ui.Success, "Brewer created")
}
