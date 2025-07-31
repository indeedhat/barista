package brewer_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) DeleteBasket(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"Open": true,
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
