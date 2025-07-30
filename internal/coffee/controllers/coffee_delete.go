package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) DeleteCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := updateCoffeeData{PageData: ui.NewPageData("Coffee", "coffee", user)}
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	coffee, err := c.repo.FindCoffee(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	pageData.Title = coffee.Name
	pageData.Coffee = *coffee
	pageData.Roasters = c.repo.IndexRoastersForUser(user)
	pageData.Flavours = c.repo.IndexFlavourProfiles()

	if len(coffee.Recipes) > 0 {
		ui.Toast(rw, ui.Warning, "Coffee cannot be deleted while it still has recipes")
		return
	}

	if err := c.repo.DeleteCoffee(coffee); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete coffee")
		return
	}

	ui.Toast(rw, ui.Success, "Coffee Deleted")
	server.Redirect(rw, r, "/coffees")
}
