package coffee_controllers

import (
	"fmt"
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
	Caffeine []kv
	Drinks   []string
	Brewers  []string
	Rating   []kv
}

type viewRecipesData struct {
	ui.PageData
	Recipes []coffee.Recipe
	Filters viewRecipesFilters
}

type kv struct {
	Key   string
	Value string
}

func (c Controller) ViewRecipes(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	recipes := c.repo.IndexRecipesForUser(user)

	ui.RenderUser(rw, r, viewRecipesData{
		PageData: ui.NewPageData("Recipes", "recipes", user),
		Recipes:  recipes,
		Filters: viewRecipesFilters{
			Coffees:  extractRecipe(recipes, func(r coffee.Recipe) *string { return &r.Coffee.Name }),
			Caffeine: extractCaffeineLevels(recipes),
			Drinks:   extractRecipe(recipes, func(r coffee.Recipe) *string { return &r.Drink }),
			Brewers: extractRecipe(recipes, func(r coffee.Recipe) *string {
				if r.Brewer == nil {
					return nil
				}
				return &r.Brewer.Name
			}),
			Rating: extractRatings(recipes),
		},
	})
}

func extractRecipe(recipes []coffee.Recipe, cb func(coffee.Recipe) *string) []string {
	var values []string
	seen := make(map[string]struct{})

	for _, recipe := range recipes {
		value := cb(recipe)
		if value == nil || *value == "" {
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

func extractRatings(recipes []coffee.Recipe) []kv {
	ratings := extractRecipe(recipes, func(r coffee.Recipe) *string {
		return ptr(fmt.Sprint(r.Rating))
	})
	slices.Sort(ratings)

	final := make([]kv, 0, len(ratings))
	for _, r := range ratings {
		final = append(final, kv{r, r + " Stars"})
	}

	return final
}

func extractCaffeineLevels(recipes []coffee.Recipe) []kv {
	levels := extractRecipe(recipes, func(r coffee.Recipe) *string {
		return ptr(fmt.Sprint(r.Coffee.Caffeine))
	})
	slices.Sort(levels)

	final := make([]kv, 0, len(levels))
	for _, l := range levels {
		switch l {
		case "1":
			final = append(final, kv{l, "Caffeinated"})
		case "2":
			final = append(final, kv{l, "Half Caf"})
		case "3":
			final = append(final, kv{l, "Decaf"})
		}
	}

	return final
}
