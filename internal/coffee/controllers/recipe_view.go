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
	Coffees  []string
	Caffeine map[int]types.CaffeineLevel
	Drinks   []types.DrinkType
	Brewers  []string
	Rating   map[int]string
}

type viewRecipesData struct {
	ui.PageData
	Recipes []coffee.Recipe
	Filters viewRecipesFilters
}

func (c Controller) ViewRecipes(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	recipes := c.repo.IndexRecipesForUser(user)

	ui.RenderUser(rw, r, viewRecipesData{
		PageData: ui.NewPageData("Recipes", "recipes", user),
		Recipes:  recipes,
		Filters: viewRecipesFilters{
			Coffees: extractRecipe(recipes, func(r coffee.Recipe) *string { return &r.Coffee.Name }),
			Caffeine: map[int]types.CaffeineLevel{
				1: types.CafLevelFull,
				2: types.CafLevelHalf,
				3: types.CafLevelDecaf,
			},
			Drinks: types.Drinks,
			Brewers: extractRecipe(recipes, func(r coffee.Recipe) *string {
				if r.Brewer == nil {
					return nil
				}
				return &r.Brewer.Name
			}),
			Rating: map[int]string{
				1: "1 Star",
				2: "2 Stars",
				3: "3 Stars",
				4: "4 Stars",
				5: "5 Stars",
			},
		},
	})
}

func extractRecipe(recipes []coffee.Recipe, cb func(coffee.Recipe) *string) []string {
	var values []string
	seen := make(map[string]struct{})

	for _, recipe := range recipes {
		value := cb(recipe)
		if value == nil {
			continue
		}
		if _, found := seen[*value]; found {
			continue
		}

		seen[*value] = struct{}{}
		values = append(values, *value)
	}

	slices.Sort(values)
	return values
}
