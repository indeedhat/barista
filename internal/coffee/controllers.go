package coffee

import (
	"errors"
	"net/http"

	"github.com/indeedhat/barista/internal/server"
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
	}

	if err := c.repo.SaveCoffee(&coffee); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusCreated, createSuccessResponse{coffee.ID})
}

type updateCoffeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Roast    uint8  `json:"roast" validate:"required"` // TODO: validate level
	Rating   *uint8 `json:"rating"`
	Flavours []uint `json:"flavours"`
	Roaster  uint   `json:"roaster" validate:"required"`
}

func (c Controller) UpdateCoffee(rw http.ResponseWriter, r *http.Request) {
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

	if err := c.repo.DeleteCoffee(coffee); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
}

type createRoasterRequest struct {
	Name string `json:"name" validate:"required"`
}

func (c Controller) CreateRoaster(rw http.ResponseWriter, r *http.Request) {
	var req createRoasterRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	roaster := Roaster{
		Name: req.Name,
	}

	if err := c.repo.SaveRoaster(&roaster); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusCreated, createSuccessResponse{roaster.ID})
}

type updateRoasterRequest struct {
	Name string `json:"name" validate:"required"`
}

func (c Controller) UpdateRoaster(rw http.ResponseWriter, r *http.Request) {
	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	var req updateRoasterRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	roaster, err := c.repo.FindRoaster(id)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("Roaster not found"))
		return
	}

	roaster.Name = req.Name

	if err := c.repo.SaveRoaster(roaster); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
}

func (c Controller) DeleteRoaster(rw http.ResponseWriter, r *http.Request) {
	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	roaster, err := c.repo.FindRoaster(id)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("Roaster not found"))
		return
	}

	if len(roaster.Coffees) > 0 {
		server.WriteResponse(
			rw,
			http.StatusNotFound,
			errors.New("Cannot delete a roaster that still has assigned coffees"),
		)
		return
	}

	if err := c.repo.DeleteRoaster(roaster); err != nil {
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

func (c Controller) UpdateFlavourProfile(rw http.ResponseWriter, r *http.Request) {
	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	var req updateFlavourProfileRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	flavour, err := c.repo.FindFlavourProfile(id)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("Flavour Profile not found"))
		return
	}

	flavour.Name = req.Name

	if err := c.repo.SaveFlavourProfile(flavour); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
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
