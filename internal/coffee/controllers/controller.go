package coffee_controllers

import (
	"time"

	"github.com/indeedhat/barista/internal/coffee"
)

const (
	CoffeeImagePath  = "data/uploads/coffee/"
	RoasterImagePath = "data/uploads/roaster/"
)

type Controller struct {
	repo coffee.Repository
}

func New(repo coffee.Repository) Controller {
	return Controller{repo}
}

type createSuccessResponse struct {
	ID uint `json:"id"`
}

type recipeSteps struct {
	Time         int    `json:"time"`
	Title        string `json:"title"`
	Instructions string `json:"instructions"`
}

func (s recipeSteps) empty() bool {
	return s.Time == 0 &&
		s.Title == "" &&
		s.Instructions == ""
}

func assignSteps(recipe *coffee.Recipe, steps []recipeSteps) {
	recipe.Steps = coffee.RecipeSteps{}
	recipe.Time = 0

	for _, step := range steps {
		if step.empty() {
			continue
		}

		rstep := coffee.RecipeStep{Instructions: step.Instructions}
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

func ptr[T any](v T) *T {
	return &v
}
