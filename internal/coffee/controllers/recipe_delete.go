package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) DeleteRecipe(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("recipe-card", ui.ComponentData{
		"Open":   true,
		"Drinks": types.Drinks,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	coffeeId, _ := server.PathID(r, "coffee_id")
	recipeId, _ := server.PathID(r, "recipe_id")
	if coffeeId == 0 || recipeId == 0 {
		ui.Toast(rw, ui.Warning, "Recipe not found")
		return
	}

	coffee, err := c.repo.FindCoffee(coffeeId, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}
	comData["Coffee"] = coffee

	recipe := coffee.Recipe(recipeId)
	if recipe == nil {
		ui.Toast(rw, ui.Warning, "Recipe not found")
		return
	}
	comData["Recipe"] = recipe

	if err := c.repo.DeleteRecipe(recipe); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete recipe")
		return
	}

	comData["Component"] = ""
	ui.Toast(rw, ui.Success, "Recipe deleted")
}
