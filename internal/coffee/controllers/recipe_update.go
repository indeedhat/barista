package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

type updateRecipeRequest struct {
	Name         string        `json:"name" validate:"required"`
	Dose         float64       `json:"dose" validate:"required"`
	WeightOut    float64       `json:"weight_out" validate:"required"`
	Drink        string        `json:"drink" validate:"required"`
	Declump      string        `json:"declump"`
	RDT          uint8         `json:"rdt"`
	Frozen       bool          `json:"frozen"`
	GrindSetting float64       `json:"grind_setting" validate:"required"`
	Grinder      string        `json:"grinder" validate:"required"`
	Steps        []recipeSteps `json:"steps"`
	Rating       uint8         `json:"rating"`
	Basket       *uint         `json:"basket"`
	Brewer       *uint         `json:"brewer"`
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

	var req updateRecipeRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	recipe.Name = req.Name
	recipe.Dose = req.Dose
	recipe.WeightOut = req.WeightOut
	recipe.Drink = req.Drink
	recipe.Declump = req.Declump
	recipe.RDT = req.RDT
	recipe.Frozen = req.Frozen
	recipe.GrindSetting = req.GrindSetting
	recipe.Grinder = req.Grinder
	recipe.Rating = req.Rating
	recipe.BrewerID = req.Brewer
	recipe.BasketID = req.Basket
	assignSteps(recipe, req.Steps)

	if err := c.repo.SaveRecipe(recipe); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to save recipe")
		return
	}

	coffee.AddRecipe(*recipe)
	comData["Recipe"] = recipe
	comData["edit"] = false
	comData.SetForm(updateRecipeRequest{})

	ui.Toast(rw, ui.Success, "Recipe updated")
}
