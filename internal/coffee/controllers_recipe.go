package coffee

import (
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type createRecipeRequest struct {
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
	comData := ui.NewComponentData("recipe-card", ui.ComponentData{
		"edit": true,
	})

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderComponent(rw, comData)
		return
	}

	coffee, err := c.repo.FindCoffee(id)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderComponent(rw, comData)
		return
	}
	comData["Coffee"] = coffee

	var req createRecipeRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		log.Print("unmarshal", err)
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		ui.RenderComponent(rw, comData)
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		log.Print(err)
		spew.Dump(comData)
		ui.Toast(rw, ui.Warning, "Bad request")
		ui.RenderComponent(rw, comData)
		return
	}

	recipe := Recipe{
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
	}
	assignSteps(&recipe, req.Steps)

	coffee.Recipes = append(coffee.Recipes, recipe)
	if err := c.repo.SaveCoffee(coffee); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create coffee")
		ui.RenderComponent(rw, comData)
		return
	}

	comData["Recipe"] = recipe
	comData["edit"] = false
	comData.SetForm(createRecipeRequest{})

	ui.Toast(rw, ui.Success, "Recipe created")
	ui.RenderComponent(rw, comData)
}

func (c Controller) NewRecipe(rw http.ResponseWriter, r *http.Request) {
	log.Print("got to controller")
	comData := ui.NewComponentData("recipe-card", ui.ComponentData{
		"Form":   map[string]struct{}{},
		"Recipe": map[string]struct{}{},
		"edit":   true,
	})

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderComponent(rw, comData)
		return
	}

	coffee, err := c.repo.FindCoffee(id)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderComponent(rw, comData)
		return
	}
	comData["Coffee"] = coffee

	ui.RenderComponent(rw, comData)
}

func assignSteps(recipe *Recipe, steps []recipeStepRequest) {
	recipe.Steps = RecipeSteps{}
	spew.Dump(steps)

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
