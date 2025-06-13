package coffee

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
)

const (
	CoffeeImagePath  = "uploads/coffee/"
	RoasterImagePath = "uploads/roaster/"
)

type Controller struct {
	repo Repository
}

func NewController(repo Repository) Controller {
	return Controller{repo}
}

type createSuccessResponse struct {
	ID uint `json:"id"`
}

type createCoffeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Roast    uint8  `json:"roast" validate:"required"` // TODO: validate level
	Rating   *uint8 `json:"rating"`
	Roaster  uint   `json:"roaster" validate:"required"`
	Flavours []uint `json:"flavours"`
}

func (c Controller) CreateCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	var req createCoffeeRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	roaster, err := c.repo.FindRoaster(req.Roaster)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("Roaster not found"))
		return
	}

	var flavours []FlavourProfile
	if len(req.Flavours) > 0 {
		if flavours, err = c.repo.FindFlavourProfiles(req.Flavours); err != nil {
			server.WriteResponse(rw, http.StatusNotFound, errors.New("One or more flavours could not be found"))
			return
		}
	}

	coffee := Coffee{
		Name:      req.Name,
		Roast:     RoastLevel(req.Roast),
		Rating:    req.Rating,
		Roaster:   *roaster,
		Flaviours: flavours,
		User:      *user,
	}

	if err := c.repo.SaveCoffee(&coffee); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusCreated, createSuccessResponse{coffee.ID})
}

func (c Controller) UpdateCoffeeImage(rw http.ResponseWriter, r *http.Request) {
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

	// TODO: i should probably find a nice way of doing this check with middleware but i cant think
	//       of a good way of doing it right now without loading in the model twice
	if coffee.User.ID != user.ID {
		server.WriteResponse(rw, http.StatusForbidden, nil)
		return
	}

	if _, err := server.UploadFile(r, "image", fmt.Sprint(CoffeeImagePath, coffee.ID), &server.UploadProps{
		Ext:  []string{".jpg", ".jpeg", ".png"},
		Mime: []string{"image/png", "image/jpeg"},
	}); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
}

type updateCoffeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Roast    uint8  `json:"roast" validate:"required"` // TODO: validate level
	Rating   *uint8 `json:"rating"`
	Flavours []uint `json:"flavours"`
	Roaster  uint   `json:"roaster" validate:"required"`
}

func (c Controller) UpdateCoffee(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	var req updateCoffeeRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
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

	roaster, err := c.repo.FindRoaster(req.Roaster)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("Roaster not found"))
		return
	}

	var flavours []FlavourProfile
	if len(req.Flavours) > 0 {
		if flavours, err = c.repo.FindFlavourProfiles(req.Flavours); err != nil {
			server.WriteResponse(rw, http.StatusNotFound, errors.New("One or more flavours could not be found"))
			return
		}
	}

	coffee.Name = req.Name
	coffee.Roast = RoastLevel(req.Roast)
	coffee.Rating = req.Rating
	coffee.Roaster = *roaster
	coffee.Flaviours = flavours

	if err := c.repo.SaveCoffee(coffee); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
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

type createFlavourProfileRequest struct {
	Name string `json:"name" validate:"required"`
}

func (c Controller) CreateFlavourProfile(rw http.ResponseWriter, r *http.Request) {
	var req createFlavourProfileRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	flavour := FlavourProfile{
		Name: req.Name,
	}

	if err := c.repo.SaveFlavourProfile(&flavour); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusCreated, createSuccessResponse{flavour.ID})
}

type updateFlavourProfileRequest struct {
	Name string `json:"name" validate:"required"`
}

func (c Controller) DeleteFlavourProfile(rw http.ResponseWriter, r *http.Request) {
	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	flavour, err := c.repo.FindFlavourProfile(id)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("Coffee not found"))
		return
	}

	if len(flavour.Coffees) > 0 {
		server.WriteResponse(
			rw,
			http.StatusNotFound,
			errors.New("Cannot delete a flavour profile that still has assigned coffees"),
		)
		return
	}

	if err := c.repo.DeleteFlavourProfile(flavour); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
}
