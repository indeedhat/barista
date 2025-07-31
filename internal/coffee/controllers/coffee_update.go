package coffee_controllers

import (
	"fmt"
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type updateCoffeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Roaster  uint   `json:"roaster" validate:"required"`
	Roast    uint8  `json:"roast" validate:"required"`    // TODO: validate level
	Caffeine uint8  `json:"caffeine" validate:"required"` // TODO: validate level
	Rating   uint8  `json:"rating"`
	Notes    string `json:"notes"`
	URL      string `json:"url"`
	Flavours []uint `json:"flavours"`
}

type updateCoffeeData struct {
	ui.PageData
	Roasters []coffee.Roaster
	Coffee   coffee.Coffee
	Coffees  []coffee.Coffee
	Flavours []coffee.FlavourProfile
	Open     bool
}

func (c Controller) UpdateCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := updateCoffeeData{PageData: ui.NewPageData("Coffee", "coffee", user)}
	defer func() {
		ui.RenderUser(rw, r, pageData)
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

	pageData.Title = coffeeModel.Name
	pageData.Coffee = *coffeeModel
	pageData.Roasters = c.repo.IndexRoastersForUser(user)
	pageData.Flavours = c.repo.IndexFlavourProfiles()

	var req updateCoffeeRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update coffee")
		return
	}

	roaster, err := c.repo.FindRoaster(req.Roaster)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	var flavours []coffee.FlavourProfile
	if len(req.Flavours) > 0 {
		if flavours, err = c.repo.FindFlavourProfiles(req.Flavours); err != nil {
			ui.Toast(rw, ui.Warning, "One or more flavours not found")
			return
		}
	}

	coffeeModel.Flavours = flavours
	coffeeModel.Roaster = *roaster
	coffeeModel.Name = req.Name
	coffeeModel.Roast = coffee.RoastLevel(req.Roast)
	coffeeModel.Caffeine = coffee.CaffeineLevel(req.Caffeine)
	coffeeModel.Rating = req.Rating
	coffeeModel.Notes = req.Notes
	coffeeModel.URL = req.URL

	if err := c.repo.SaveCoffee(coffeeModel); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update coffee")
		return
	}

	pageData.Title = coffeeModel.Name
	ui.Toast(rw, ui.Success, "Coffee Updated")
}

func (c Controller) UpdateCoffeeImage(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := updateCoffeeData{PageData: ui.NewPageData("Coffee", "coffee", user)}
	defer func() {
		ui.RenderUser(rw, r, pageData)
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

	pageData.Title = coffee.Name
	pageData.Coffee = *coffee
	pageData.Roasters = c.repo.IndexRoastersForUser(user)
	pageData.Flavours = c.repo.IndexFlavourProfiles()

	savePath, err := server.UploadFile(r, "image", fmt.Sprint(CoffeeImagePath, coffee.ID), &server.UploadProps{
		Ext:  []string{".jpg", ".jpeg", ".png"},
		Mime: []string{"image/png", "image/jpeg"},
	})
	if err != nil {
		ui.Toast(rw, ui.Warning, "Failed to upload image")
		return
	}

	if savePath != "" {
		coffee.Icon = savePath[5:]
		if err := c.repo.SaveCoffee(coffee); err != nil {
			ui.Toast(rw, ui.Warning, "Failed to save image")
			return
		}
	}

	ui.Toast(rw, ui.Success, "Image Updated")
}
