package coffee

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) ViewCoffees(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	roasters := c.repo.IndexRoastersForUser(user)
	coffees := c.repo.IndexCoffeesForUser(user)

	pageData := ui.NewPageData("Coffees", "coffees", user)
	pageData.Form = createCoffeeRequest{}
	pageData.Data["Roasters"] = roasters
	pageData.Data["Coffees"] = coffees

	ui.RenderUser(rw, r, pageData)
}

type createCoffeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Roaster  uint   `json:"roaster" validate:"required"`
	Roast    uint8  `json:"roast" validate:"required"` // TODO: validate level
	Rating   uint8  `json:"rating"`
	Notes    string `json:"notes"`
	URL      string `json:"url"`
	Flavours []uint `json:"flavours"`
}

func (c Controller) CreateCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Coffees", "coffees", user)
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)
	pageData.Data["Coffees"] = c.repo.IndexCoffeesForUser(user)
	pageData.Data["open"] = true

	var req createCoffeeRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		ui.RenderUser(rw, r, pageData)
		return
	}

	spew.Dump(pageData.Form)

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		ui.RenderUser(rw, r, pageData)
		return
	}

	roaster, err := c.repo.FindRoaster(req.Roaster)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		ui.RenderUser(rw, r, pageData)
		return
	}

	var flavours []FlavourProfile
	if len(req.Flavours) > 0 {
		if flavours, err = c.repo.FindFlavourProfiles(req.Flavours); err != nil {
			ui.Toast(rw, ui.Warning, "One or more flavours not found")
			ui.RenderUser(rw, r, pageData)
			return
		}
	}

	coffee := Coffee{
		Name:     req.Name,
		Roaster:  *roaster,
		Roast:    RoastLevel(req.Roast),
		Rating:   req.Rating,
		Notes:    req.Notes,
		URL:      req.URL,
		Flavours: flavours,
		User:     *user,
	}

	if err := c.repo.SaveCoffee(&coffee); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create coffee")
		ui.RenderUser(rw, r, pageData)
		return
	}

	pageData.Data["Coffees"] = c.repo.IndexCoffeesForUser(user)
	pageData.Data["open"] = false

	ui.Toast(rw, ui.Success, "Coffee created")
	ui.RenderUser(rw, r, pageData)
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

	ui.RenderUser(rw, r, pageData)
}

func (c Controller) UpdateCoffeeImage(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Coffee", "coffee", user)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderUser(rw, r, pageData)
		return
	}

	coffee, err := c.repo.FindCoffee(id)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderUser(rw, r, pageData)
		return
	}

	pageData.Title = coffee.Name
	pageData.Data["Coffee"] = coffee
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)

	// TODO: i should probably find a nice way of doing this check with middleware but i cant think
	//       of a good way of doing it right now without loading in the model twice
	if coffee.UserID != user.ID {
		ui.Toast(rw, ui.Warning, "Coffee does not belong to you")
		ui.RenderUser(rw, r, pageData)
		return
	}

	savePath, err := server.UploadFile(r, "image", fmt.Sprint(CoffeeImagePath, coffee.ID), &server.UploadProps{
		Ext:  []string{".jpg", ".jpeg", ".png"},
		Mime: []string{"image/png", "image/jpeg"},
	})
	if err != nil {
		ui.Toast(rw, ui.Warning, "Failed to upload image")
		ui.RenderUser(rw, r, pageData)
		return
	}

	if savePath != "" {
		coffee.Icon = savePath
		if err := c.repo.SaveCoffee(coffee); err != nil {
			ui.Toast(rw, ui.Warning, "Failed to save image")
		}
	}

	ui.RenderUser(rw, r, pageData)
}

type updateCoffeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Roaster  uint   `json:"roaster" validate:"required"`
	Roast    uint8  `json:"roast" validate:"required"` // TODO: validate level
	Rating   uint8  `json:"rating"`
	Notes    string `json:"notes"`
	URL      string `json:"url"`
	Flavours []uint `json:"flavours"`
}

func (c Controller) UpdateCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Coffee", "coffee", user)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderUser(rw, r, pageData)
		return
	}

	coffee, err := c.repo.FindCoffee(id)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderUser(rw, r, pageData)
		return
	}

	pageData.Title = coffee.Name
	pageData.Data["Coffee"] = coffee
	pageData.Data["Roasters"] = c.repo.IndexRoastersForUser(user)

	var req updateCoffeeRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		ui.RenderUser(rw, r, pageData)
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update coffee")
		ui.RenderUser(rw, r, pageData)
		return
	}

	if coffee.UserID != user.ID {
		ui.Toast(rw, ui.Warning, "Coffee does not belong to you")
		ui.RenderUser(rw, r, pageData)
		return
	}

	roaster, err := c.repo.FindRoaster(req.Roaster)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Coffee not found")
		ui.RenderUser(rw, r, pageData)
		return
	}

	var flavours []FlavourProfile
	if len(req.Flavours) > 0 {
		if flavours, err = c.repo.FindFlavourProfiles(req.Flavours); err != nil {
			ui.Toast(rw, ui.Warning, "One or more flavours not found")
			ui.RenderUser(rw, r, pageData)
			return
		}
	}

	coffee.Name = req.Name
	coffee.Roaster = *roaster
	coffee.Roast = RoastLevel(req.Roast)
	coffee.Rating = req.Rating
	coffee.Notes = req.Notes
	coffee.URL = req.URL
	coffee.Flavours = flavours

	if err := c.repo.SaveCoffee(coffee); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update coffee")
		ui.RenderUser(rw, r, pageData)
		return
	}

	pageData.Title = coffee.Name
	ui.Toast(rw, ui.Success, "Coffee Updated")
	ui.RenderUser(rw, r, pageData)
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
