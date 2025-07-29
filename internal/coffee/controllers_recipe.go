package coffee

import (
	"net/http"
	"time"

	"github.com/indeedhat/barista/internal/auth"
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

func (c Controller) ViewRecipes(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Recipes", "recipes", user)
	pageData.Data["Recipes"] = c.repo.IndexRecipesForUser(user)
	pageData.Data["Drinks"] = types.Drinks

	ui.RenderUser(rw, r, pageData)
}

type upsertRecipeRequest struct {
	Name         string              `json:"name" validate:"required"`
	Dose         float64             `json:"dose" validate:"required"`
	WeightOut    float64             `json:"weight_out" validate:"required"`
	Drink        string              `json:"drink" validate:"required"`
	Declump      string              `json:"declump"`
	RDT          uint8               `json:"rdt"`
	Frozen       bool                `json:"frozen"`
	GrindSetting float64             `json:"grind_setting" validate:"required"`
	Grinder      string              `json:"grinder" validate:"required"`
	Steps        []recipeStepRequest `json:"steps"`
	Rating       uint8               `json:"rating"`
	Basket       *uint               `json:"basket"`
	Brewer       *uint               `json:"brewer"`
}

func (r upsertRecipeRequest) apply(recipe *Recipe) {
	recipe.Name = r.Name
	recipe.Dose = r.Dose
	recipe.WeightOut = r.WeightOut
	recipe.Drink = r.Drink
	recipe.Declump = r.Declump
	recipe.RDT = r.RDT
	recipe.Frozen = r.Frozen
	recipe.GrindSetting = r.GrindSetting
	recipe.Grinder = r.Grinder
	recipe.Rating = r.Rating
	recipe.BrewerID = r.Brewer
	recipe.BasketID = r.Basket

	assignSteps(recipe, r.Steps)
}

type recipeStepRequest struct {
	Time         int    `json:"time"`
	Title        string `json:"title"`
	Instructions string `json:"instructions"`
}

func (s recipeStepRequest) empty() bool {
	return s.Time == 0 &&
		s.Title == "" &&
		s.Instructions == ""
}

func (c Controller) CreateRecipe(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("recipe-card", ui.ComponentData{
		"edit":   true,
		"Drinks": types.Drinks,
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

	var req upsertRecipeRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	recipe := Recipe{
		User:   *user,
		Coffee: *coffee,
	}
	req.apply(&recipe)

	coffee.Recipes = append(coffee.Recipes, recipe)
	if err := c.repo.SaveRecipe(&recipe); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create recipe")
		return
	}

	comData["Recipe"] = recipe
	comData["edit"] = false
	comData.SetForm(upsertRecipeRequest{})

	ui.Toast(rw, ui.Success, "Recipe created")
}

func (c Controller) UpdateRecipe(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("recipe-card", ui.ComponentData{
		"edit":   true,
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

	var req upsertRecipeRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	req.apply(recipe)

	if err := c.repo.SaveRecipe(recipe); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to save recipe")
		return
	}

	coffee.AddRecipe(*recipe)
	comData["Recipe"] = recipe
	comData["edit"] = false
	comData.SetForm(upsertRecipeRequest{})

	ui.Toast(rw, ui.Success, "Recipe updated")
}

func (c Controller) DeleteRecipe(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("recipe-card", ui.ComponentData{
		"open":   true,
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

func assignSteps(recipe *Recipe, steps []recipeStepRequest) {
	recipe.Steps = RecipeSteps{}
	recipe.Time = 0

	for _, step := range steps {
		if step.empty() {
			continue
		}

		rstep := RecipeStep{Instructions: step.Instructions}
		if step.Time > 0 {
			recipe.Time += time.Duration(step.Time) * time.Second
			rstep.Time = ptr(time.Duration(step.Time) * time.Second)
		}
		if step.Title != "" {
			rstep.Title = &step.Title
		}

		recipe.Steps = append(recipe.Steps, rstep)
	}
}
