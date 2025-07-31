package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

type createRecipeRequest struct {
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

	coffeeModel, err := c.repo.FindCoffee(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}
	comData["Coffee"] = coffeeModel

	var req createRecipeRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	recipe := coffee.Recipe{
		User:         *user,
		Coffee:       *coffeeModel,
		Name:         req.Name,
		Dose:         req.Dose,
		WeightOut:    req.WeightOut,
		Drink:        req.Drink,
		Declump:      req.Declump,
		RDT:          req.RDT,
		Frozen:       req.Frozen,
		GrindSetting: req.GrindSetting,
		Grinder:      req.Grinder,
		Rating:       req.Rating,
		BrewerID:     req.Brewer,
		BasketID:     req.Basket,
	}
	assignSteps(&recipe, req.Steps)

	coffeeModel.Recipes = append(coffeeModel.Recipes, recipe)
	if err := c.repo.SaveRecipe(&recipe); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create recipe")
		return
	}

	comData["Recipe"] = recipe
	comData["edit"] = false
	comData.SetForm(createRecipeRequest{})

	ui.Toast(rw, ui.Success, "Recipe created")
}
