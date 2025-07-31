package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type createRoasterRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"omitempty,url"`
}

type createRoasterData struct {
	ui.PageData
	Roaster  *coffee.Roaster
	Roasters []coffee.Roaster
	Open     bool
}

func (c Controller) CreateRoaster(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := createRoasterData{PageData: ui.NewPageData("Roasters", "roasters", user)}
	pageData.Roasters = c.repo.IndexRoastersForUser(user)
	pageData.Open = true
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	var req createRoasterRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create roaster")
		return
	}

	roaster := coffee.Roaster{
		User:        *user,
		Name:        req.Name,
		Description: req.Description,
		URL:         req.URL,
	}

	if err := c.repo.SaveRoaster(&roaster); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create roaster")
		return
	}

	pageData.Roasters = c.repo.IndexRoastersForUser(user)
	pageData.Open = false
	pageData.Form = createRoasterRequest{}

	ui.Toast(rw, ui.Success, "Roaster created")
}
