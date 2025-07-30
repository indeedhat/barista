package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

type createCoffeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Roaster  uint   `json:"roaster" validate:"required"`
	Roast    uint8  `json:"roast" validate:"required"`    // TODO: validate level
	Caffeine uint8  `json:"caffeine" validate:"required"` // TODO: validate level
	Rating   uint8  `json:"rating"`
	Notes    string `json:"notes"`
	URL      string `json:"url"`
	Flavours []uint `json:"flavours"`
}

type createCoffeeData struct {
	ui.PageData
	Roasters []coffee.Roaster
	Coffee   coffee.Coffee
	Coffees  []coffee.Coffee
	Flavours []coffee.FlavourProfile
	Drinks   []types.DrinkType
	Open     bool
}

func (c Controller) CreateCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := createCoffeeData{PageData: ui.NewPageData("Coffees", "coffees", user)}
	pageData.Roasters = c.repo.IndexRoastersForUser(user)
	pageData.Coffees = c.repo.IndexCoffeesForUser(user)
	pageData.Flavours = c.repo.IndexFlavourProfiles()
	pageData.Drinks = types.Drinks
	pageData.Open = true
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	var req createCoffeeRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	roaster, err := c.repo.FindRoaster(req.Roaster)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	var flavours []coffee.FlavourProfile
	if len(req.Flavours) > 0 {
		if flavours, err = c.repo.FindFlavourProfiles(req.Flavours); err != nil {
			ui.Toast(rw, ui.Warning, "One or more flavours not found")
			return
		}
	}

	coffee := coffee.Coffee{
		Roaster:  *roaster,
		Flavours: flavours,
		User:     *user,
		Name:     req.Name,
		Roast:    coffee.RoastLevel(req.Roast),
		Caffeine: coffee.CaffeineLevel(req.Caffeine),
		Rating:   req.Rating,
		Notes:    req.Notes,
		URL:      req.URL,
	}

	if err := c.repo.SaveCoffee(&coffee); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create coffee")
		return
	}

	pageData.Coffees = c.repo.IndexCoffeesForUser(user)
	pageData.Open = false
	pageData.Form = createCoffeeRequest{}

	ui.Toast(rw, ui.Success, "Coffee created")
}
