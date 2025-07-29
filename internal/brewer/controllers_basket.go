package brewer

import (
	"net/http"
	"strconv"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
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

type upsertBasketRequest struct {
	Name  string  `json:"name" validate:"required"`
	Brand string  `json:"brand" validate:"required"`
	Dose  float64 `json:"dose" validate:"required"`
}

func (r upsertBasketRequest) apply(basket *Basket) {
	basket.Name = r.Name
	basket.Brand = r.Brand
	basket.Dose = r.Dose
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
	brewer := c.findEspressoBrewer(rw, user, id)
	if brewer == nil {
		return
	}
	comData["Brewer"] = brewer

	var req upsertBasketRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	basket := Basket{Brewer: *brewer}
	req.apply(&basket)

	brewer.AddBasket(basket)
	if err := c.repo.SaveBrewer(brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create basket")
		return
	}

	comData["Basket"] = basket
	comData["edit"] = false
	comData.SetForm(upsertBasketRequest{})

	ui.Toast(rw, ui.Success, "Basket created")
}

func (c Controller) UpdateBasket(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"edit": true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	brewerId, _ := server.PathID(r, "brewer_id")
	basketId, _ := server.PathID(r, "basket_id")
	if brewerId == 0 || basketId == 0 {
		ui.Toast(rw, ui.Warning, "Basket not found")
		return
	}

	brewer := c.findEspressoBrewer(rw, user, brewerId)
	if brewer == nil {
		return
	}
	comData["Brewer"] = brewer

	basket := brewer.Basket(basketId)
	if basket == nil {
		ui.Toast(rw, ui.Warning, "Basket not found")
		return
	}
	comData["Basket"] = basket

	var req upsertBasketRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	req.apply(basket)
	brewer.AddBasket(*basket)

	if err := c.repo.SaveBrewer(brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to save basket")
		return
	}

	comData["Basket"] = basket
	comData["edit"] = false
	comData.SetForm(upsertBasketRequest{})

	ui.Toast(rw, ui.Success, "Basket updated")
}

func (c Controller) DeleteBasket(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"open": true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	brewerId, _ := server.PathID(r, "brewer_id")
	basketId, _ := server.PathID(r, "basket_id")
	if brewerId == 0 || basketId == 0 {
		ui.Toast(rw, ui.Warning, "Basket not found")
		return
	}

	brewer := c.findEspressoBrewer(rw, user, brewerId)
	if brewer == nil {
		return
	}
	comData["Brewer"] = brewer

	basket := brewer.Basket(basketId)
	if basket == nil {
		ui.Toast(rw, ui.Warning, "Basket not found")
		return
	}
	comData["Basket"] = basket

	brewer.RemoveBasket(*basket)

	if err := c.repo.SaveBrewer(brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete basket")
		return
	}

	comData["Component"] = ""
	ui.Toast(rw, ui.Success, "Basket deleted")
}

func (c Controller) BasketsSelect(rw http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("brewer.int")
	if idStr == "" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	user := r.Context().Value("user").(*auth.User)
	brewerId, _ := strconv.Atoi(idStr)
	brewer, err := c.repo.FindBrewer(uint(brewerId), user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return
	}
	if brewer.Type != "Espresso" {
		return
	}

	comData := ui.NewComponentData("baskets-select", ui.ComponentData{
		"value": r.URL.Query().Get("value"),
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	comData["Baskets"] = brewer.Baskets
}

func (c Controller) findEspressoBrewer(rw http.ResponseWriter, user *auth.User, id uint) *Brewer {
	if id == 0 {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return nil
	}

	brewer, err := c.repo.FindBrewer(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return nil
	}

	if brewer.Type != types.BrewerEspresso {
		ui.Toast(rw, ui.Warning, "Only espresso machines can have baskets")
		return nil
	}

	return brewer
}
