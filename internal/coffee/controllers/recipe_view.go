package coffee_controllers

import (
	"net/http"
	"slices"

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
	Caffiene    map[int]types.CaffieneLevel
	DrinkTypes  []types.DrinkType
	BrewerTypes []types.BrewerType
	Rating      map[int]string
}

type viewRecipesData struct {
	ui.PageData
	Recipes []coffee.Recipe
	Filters viewRecipesFilters
}

func (c Controller) ViewRecipes(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	recipes := c.repo.IndexRecipesForUser(user)
	pageData := viewRecipesData{
		PageData: ui.NewPageData("Recipes", "recipes", user),
		Recipes:  recipes,
		Filters: viewRecipesFilters{
			Coffees: extractRecipeCoffees(recipes),
			Caffiene: map[int]types.CaffieneLevel{
				1: types.CafLevelFull,
				2: types.CafLevelHalf,
				3: types.CafLevelDecaf,
			},
			DrinkTypes:  types.Drinks,
			BrewerTypes: types.Brewers,
			Rating: map[int]string{
				1: "1 Star",
				2: "2 Stars",
				3: "3 Stars",
				4: "4 Stars",
				5: "5 Stars",
			},
		},
	}

	ui.RenderUser(rw, r, pageData)
}

func extractRecipeCoffees(recipes []coffee.Recipe) []string {
	var seen map[string]struct{}
	var coffees []string

	for _, recipe := range recipes {
		if _, found := seen[recipe.Coffee.Name]; found {
			continue
		}

		seen[recipe.Coffee.Name] = struct{}{}
		coffees = append(coffees, recipe.Coffee.Name)
	}

	slices.Sort(coffees)
	return coffees
}
