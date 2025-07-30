package brewer_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/brewer"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) NewBasket(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"Form":   map[string]struct{}{},
		"Basket": map[string]struct{}{},
		"edit":   true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	id, _ := server.PathID(r)
	comData["Brewer"] = c.findEspressoBrewer(rw, user, id)
}

type createBasketRequest struct {
	Name  string  `json:"name" validate:"required"`
	Brand string  `json:"brand" validate:"required"`
	Dose  float64 `json:"dose" validate:"required"`
}

func (c Controller) CreateBasket(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"edit": true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	id, _ := server.PathID(r)
	brewerModel := c.findEspressoBrewer(rw, user, id)
	if brewerModel == nil {
		return
	}
	comData["Brewer"] = brewerModel

	var req createBasketRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	basket := brewer.Basket{
		Name:   req.Name,
		Brand:  req.Brand,
		Dose:   req.Dose,
		Brewer: *brewerModel,
	}
	brewerModel.AddBasket(basket)

	if err := c.repo.SaveBrewer(brewerModel); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create basket")
		return
	}

	comData["Basket"] = basket
	comData["edit"] = false
	comData.SetForm(createBasketRequest{})

	ui.Toast(rw, ui.Success, "Basket created")
}
