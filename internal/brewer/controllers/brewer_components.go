package brewer_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) BrewersSelect(rw http.ResponseWriter, r *http.Request) {
	drink := r.URL.Query().Get("drink")
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("brewers-select", ui.ComponentData{
		"Brewers": c.repo.IndexBrewersForUser(user, types.DrinkType(drink).Brewers()...),
		"value":   r.URL.Query().Get("value"),
	})

	ui.RenderComponent(rw, comData)
}
