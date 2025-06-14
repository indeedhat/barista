package coffee

import (
	"errors"
	"net/http"

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
