package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

type viewCoffeesData struct {
	ui.PageData
	Roasters []coffee.Roaster
	Coffees  []coffee.Coffee
	Flavours []coffee.FlavourProfile
	Drinks   []types.DrinkType
	Open     bool
}

func (c Controller) ViewCoffees(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	pageData := viewCoffeesData{PageData: ui.NewPageData("Coffees", "coffees", user)}
	pageData.Form = createCoffeeRequest{}
	pageData.Roasters = c.repo.IndexRoastersForUser(user)
	pageData.Coffees = c.repo.IndexCoffeesForUser(user)
	pageData.Flavours = c.repo.IndexFlavourProfiles()
	pageData.Drinks = types.Drinks

	ui.RenderUser(rw, r, pageData)
}

type viewCoffeeData struct {
	ui.PageData
	Coffee   *coffee.Coffee
	Roasters []coffee.Roaster
	Flavours []coffee.FlavourProfile
	Drinks   []types.DrinkType
	Open     bool
}

func (c Controller) ViewCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Coffee Not Found", "404", user))
		return
	}

	coffee, err := c.repo.FindCoffee(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Coffee Not Found", "404", user))
		return
	}

	pageData := viewCoffeeData{PageData: ui.NewPageData("Coffee", "coffee", user)}
	pageData.Coffee = coffee
	pageData.Drinks = types.Drinks
	pageData.Roasters = c.repo.IndexRoastersForUser(user)
	pageData.Flavours = c.repo.IndexFlavourProfiles()
	pageData.Form = createCoffeeRequest{
		Flavours: coffee.FlavourIds(),
	}

	ui.RenderUser(rw, r, pageData)
}
