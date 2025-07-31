package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) NewRecipe(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("recipe-card", ui.ComponentData{
		"Form":   map[string]struct{}{},
		"Recipe": map[string]struct{}{},
		"Drinks": types.Drinks,
		"edit":   true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
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

	comData["Coffee"] = coffee
}

type viewRecipesFilters struct {
	Coffees     []string
	Cafiene     []string
	DrinkTypes  []string
	BrewerTypes []string
	Rating      []string
}

type viewRecipesData struct {
	ui.PageData
	Recipes []coffee.Recipe
	Filters viewRecipesFilters
}

func (c Controller) ViewRecipes(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := viewRecipesData{
		PageData: ui.NewPageData("Recipes", "recipes", user),
		Recipes:  c.repo.IndexRecipesForUser(user),
		Filters: viewRecipesFilters{
			Cafiene: []string{},
		},
	}

	ui.RenderUser(rw, r, pageData)
}
