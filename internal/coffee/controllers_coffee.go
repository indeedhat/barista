package coffee

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) ViewCoffees(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	pageData := ui.NewPageData("Coffees", "coffees", user)
	pageData.Form = upsertCoffeeRequest{}
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)
	pageData.Data["Coffees"] = c.repo.IndexCoffeesForUser(user)
	pageData.Data["Flavours"] = c.repo.IndexFlavourProfiles()

	ui.RenderUser(rw, r, pageData)
}

type upsertCoffeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Roaster  uint   `json:"roaster" validate:"required"`
	Roast    uint8  `json:"roast" validate:"required"`    // TODO: validate level
	Caffeine uint8  `json:"caffeine" validate:"required"` // TODO: validate level
	Rating   uint8  `json:"rating"`
	Notes    string `json:"notes"`
	URL      string `json:"url"`
	Flavours []uint `json:"flavours"`
}

func (r upsertCoffeeRequest) apply(coffee *Coffee) {
	coffee.Name = r.Name
	coffee.Roast = RoastLevel(r.Roast)
	coffee.Caffeine = CaffeineLevel(r.Caffeine)
	coffee.Rating = r.Rating
	coffee.Notes = r.Notes
	coffee.URL = r.URL
}

func (c Controller) CreateCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Coffees", "coffees", user)
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)
	pageData.Data["Coffees"] = c.repo.IndexCoffeesForUser(user)
	pageData.Data["Flavours"] = c.repo.IndexFlavourProfiles()
	pageData.Data["open"] = true
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	var req upsertCoffeeRequest
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

	var flavours []FlavourProfile
	if len(req.Flavours) > 0 {
		if flavours, err = c.repo.FindFlavourProfiles(req.Flavours); err != nil {
			ui.Toast(rw, ui.Warning, "One or more flavours not found")
			return
		}
	}

	coffee := Coffee{
		Roaster:  *roaster,
		Flavours: flavours,
		User:     *user,
	}
	req.apply(&coffee)

	if err := c.repo.SaveCoffee(&coffee); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create coffee")
		return
	}

	pageData.Data["Coffees"] = c.repo.IndexCoffeesForUser(user)
	pageData.Data["open"] = false
	pageData.Form = upsertCoffeeRequest{}

	ui.Toast(rw, ui.Success, "Coffee created")
}

func (c Controller) ViewCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Coffee Not Found", "404", user))
		return
	}

	coffee, err := c.repo.FindCoffee(id)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Coffee Not Found", "404", user))
		return
	}

	pageData := ui.NewPageData("Coffee", "coffee", user)
	pageData.Data["Coffee"] = coffee
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)
	pageData.Data["Flavours"] = c.repo.IndexFlavourProfiles()
	pageData.Form = upsertCoffeeRequest{
		Flavours: coffee.FlavourIds(),
	}

	ui.RenderUser(rw, r, pageData)
}

func (c Controller) UpdateCoffeeImage(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Coffee", "coffee", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	coffee, err := c.repo.FindCoffee(id)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	pageData.Title = coffee.Name
	pageData.Data["Coffee"] = coffee
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)
	pageData.Data["Flavours"] = c.repo.IndexFlavourProfiles()

	// TODO: i should probably find a nice way of doing this check with middleware but i cant think
	//       of a good way of doing it right now without loading in the model twice
	if coffee.UserID != user.ID {
		ui.Toast(rw, ui.Warning, "Coffee does not belong to you")
		return
	}

	savePath, err := server.UploadFile(r, "image", fmt.Sprint(CoffeeImagePath, coffee.ID), &server.UploadProps{
		Ext:  []string{".jpg", ".jpeg", ".png"},
		Mime: []string{"image/png", "image/jpeg"},
	})
	if err != nil {
		ui.Toast(rw, ui.Warning, "Failed to upload image")
		return
	}

	if savePath != "" {
		coffee.Icon = savePath
		if err := c.repo.SaveCoffee(coffee); err != nil {
			ui.Toast(rw, ui.Warning, "Failed to save image")
			return
		}
	}

	ui.Toast(rw, ui.Success, "Image Updated")
}

func (c Controller) UpdateCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Coffee", "coffee", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	coffee, err := c.repo.FindCoffee(id)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	pageData.Title = coffee.Name
	pageData.Data["Coffee"] = coffee
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)
	pageData.Data["Flavours"] = c.repo.IndexFlavourProfiles()

	var req upsertCoffeeRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update coffee")
		return
	}

	if coffee.UserID != user.ID {
		ui.Toast(rw, ui.Warning, "Coffee does not belong to you")
		return
	}

	roaster, err := c.repo.FindRoaster(req.Roaster)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		return
	}

	var flavours []FlavourProfile
	if len(req.Flavours) > 0 {
		if flavours, err = c.repo.FindFlavourProfiles(req.Flavours); err != nil {
			ui.Toast(rw, ui.Warning, "One or more flavours not found")
			return
		}
	}

	coffee.Flavours = flavours
	coffee.Roaster = *roaster
	req.apply(coffee)
	if err := c.repo.SaveCoffee(coffee); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update coffee")
		return
	}

	pageData.Title = coffee.Name
	ui.Toast(rw, ui.Success, "Coffee Updated")
}

func (c Controller) DeleteCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	coffee, err := c.repo.FindCoffee(id)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("Coffee not found"))
		return
	}

	if coffee.User.ID != user.ID {
		server.WriteResponse(rw, http.StatusForbidden, nil)
		return
	}

	if err := c.repo.DeleteCoffee(coffee); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
}
