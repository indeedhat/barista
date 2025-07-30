package coffee_controllers

import (
	"github.com/indeedhat/barista/internal/coffee"
)

const (
	CoffeeImagePath  = "data/uploads/coffee/"
	RoasterImagePath = "data/uploads/roaster/"
)

type Controller struct {
	repo coffee.Repository
}

func NewController(repo coffee.Repository) Controller {
	return Controller{repo}
}

type createSuccessResponse struct {
	ID uint `json:"id"`
}
