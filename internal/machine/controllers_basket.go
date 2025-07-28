package machine

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) NewBasket(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"Form":   map[string]struct{}{},
		"Basket": map[string]struct{}{},
		"edit":   true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Machine not found")
		return
	}

	machine, err := c.repo.FindMachine(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Machine not found")
		return
	}

	comData["Machine"] = machine
}

type upsertBasketRequest struct {
	Name  string  `json:"name" validate:"required"`
	Brand string  `json:"brand" validate:"required"`
	Dose  float64 `json:"dose" validate:"required"`
}

func (r upsertBasketRequest) apply(basket *Basket) {
	basket.Name = r.Name
	basket.Brand = r.Brand
	basket.Dose = r.Dose
}

func (c Controller) CreateBasket(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"edit": true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Machine not found")
		return
	}

	machine, err := c.repo.FindMachine(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Machine not found")
		return
	}
	comData["Machine"] = machine

	var req upsertBasketRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	basket := Basket{}
	req.apply(&basket)

	machine.AddBasket(basket)
	if err := c.repo.SaveMachine(machine); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create basket")
		return
	}

	comData["Basket"] = basket
	comData["edit"] = false
	comData.SetForm(upsertBasketRequest{})

	ui.Toast(rw, ui.Success, "Basket created")
}

func (c Controller) UpdateBasket(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"edit": true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	machineId, _ := server.PathID(r, "machine_id")
	basketId, _ := server.PathID(r, "basket_id")
	if machineId == 0 || basketId == 0 {
		ui.Toast(rw, ui.Warning, "Basket not found")
		return
	}

	machine, err := c.repo.FindMachine(machineId, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Machine not found")
		return
	}
	comData["Machine"] = machine

	basket := machine.Basket(basketId)
	if basket == nil {
		ui.Toast(rw, ui.Warning, "Basket not found")
		return
	}
	comData["Basket"] = basket

	var req upsertBasketRequest
	if err := server.UnmarshalBody(r, &req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &comData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	req.apply(basket)
	machine.AddBasket(*basket)

	if err := c.repo.SaveMachine(machine); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to save basket")
		return
	}

	comData["Basket"] = basket
	comData["edit"] = false
	comData.SetForm(upsertBasketRequest{})

	ui.Toast(rw, ui.Success, "Basket updated")
}

func (c Controller) DeleteBasket(rw http.ResponseWriter, r *http.Request) {
	comData := ui.NewComponentData("basket-card", ui.ComponentData{
		"open": true,
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	machineId, _ := server.PathID(r, "machine_id")
	basketId, _ := server.PathID(r, "basket_id")
	if machineId == 0 || basketId == 0 {
		ui.Toast(rw, ui.Warning, "Basket not found")
		return
	}

	machine, err := c.repo.FindMachine(machineId)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Machine not found")
		return
	}
	comData["Machine"] = machine

	basket := machine.Basket(basketId)
	if basket == nil {
		ui.Toast(rw, ui.Warning, "Basket not found")
		return
	}
	comData["Basket"] = basket

	machine.RemoveBasket(*basket)

	if err := c.repo.SaveMachine(machine); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete basket")
		return
	}

	comData["Component"] = ""
	ui.Toast(rw, ui.Success, "Basket deleted")
}
