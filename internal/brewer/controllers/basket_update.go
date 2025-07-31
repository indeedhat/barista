package brewer_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type updateBasketRequest struct {
	Name  string  `json:"name" validate:"required"`
	Brand string  `json:"brand" validate:"required"`
	Dose  float64 `json:"dose" validate:"required"`
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

	var req updateBasketRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	basket.Name = req.Name
	basket.Brand = req.Brand
	basket.Dose = req.Dose

	brewer.AddBasket(*basket)

	if err := c.repo.SaveBrewer(brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to save basket")
		return
	}

	comData["Basket"] = basket
	comData["edit"] = false
	comData.SetForm(updateBasketRequest{})

	ui.Toast(rw, ui.Success, "Basket updated")
}
