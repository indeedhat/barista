package machine

import (
	"fmt"
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

const MachineImagePath = "data/uploads/machine/"

type Controller struct {
	repo Repository
}

func NewController(repo Repository) Controller {
	return Controller{repo}
}

func (c Controller) ViewMachines(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	pageData := ui.NewPageData("Machines", "machines", user)
	pageData.Data["Machines"] = c.repo.IndexMachinesForUser(user)
	pageData.Form = upsertMachineRequest{}

	ui.RenderUser(rw, r, pageData)
}

type upsertMachineRequest struct {
	Name        string `json:"name" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	ModelNumber string `json:"model" validate:"required"`
}

func (r upsertMachineRequest) apply(machine *Machine) {
	machine.Name = r.Name
	machine.Brand = r.Brand
	machine.ModelNumber = r.ModelNumber
}

func (c Controller) CreateMachine(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Machines", "machines", user)
	pageData.Data["open"] = true
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	var req upsertMachineRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Bad request")
		return
	}

	machine := Machine{}
	req.apply(&machine)

	if err := c.repo.SaveMachine(&machine); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create machine")
		return
	}

	pageData.Data["Machines"] = c.repo.IndexMachinesForUser(user)
	pageData.Data["open"] = false
	pageData.Form = upsertMachineRequest{}

	ui.Toast(rw, ui.Success, "Machine created")
}

func (c Controller) ViewMachine(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Machine Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Machine Not Found", "404", user))
		return
	}

	machine, err := c.repo.FindMachine(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Machine Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Machine Not Found", "404", user))
		return
	}

	pageData := ui.NewPageData("Machine", "machine", user)
	pageData.Data["Machine"] = machine
	pageData.Form = upsertMachineRequest{}

	ui.RenderUser(rw, r, pageData)
}

func (c Controller) UpdateMachineImage(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Machine", "machine", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
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

	pageData.Title = machine.Name
	pageData.Data["Machine"] = machine

	savePath, err := server.UploadFile(r, "image", fmt.Sprint(MachineImagePath, machine.ID), &server.UploadProps{
		Ext:  []string{".jpg", ".jpeg", ".png"},
		Mime: []string{"image/png", "image/jpeg"},
	})
	if err != nil {
		ui.Toast(rw, ui.Warning, "Failed to upload image")
		return
	}

	if savePath != "" {
		machine.Icon = savePath[5:]
		if err := c.repo.SaveMachine(machine); err != nil {
			ui.Toast(rw, ui.Warning, "Failed to save image")
			return
		}
	}

	ui.Toast(rw, ui.Success, "Image Updated")
}

func (c Controller) UpdateMachine(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Machine", "machine", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
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

	pageData.Title = machine.Name
	pageData.Data["Machine"] = machine

	var req upsertMachineRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update machine")
		return
	}

	req.apply(machine)
	if err := c.repo.SaveMachine(machine); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to update machine")
		return
	}

	pageData.Title = machine.Name
	ui.Toast(rw, ui.Success, "Machine Updated")
}

func (c Controller) DeleteMachine(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := ui.NewPageData("Machine", "machine", user)
	defer func() {
		ui.RenderUser(rw, r, pageData)
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

	pageData.Title = machine.Name
	pageData.Data["Machine"] = machine

	if len(machine.Baskets) > 0 {
		ui.Toast(rw, ui.Warning, "Machine cannot be deleted while it still has baskets")
		return
	}

	if err := c.repo.DeleteMachine(machine); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete machine")
		return
	}

	ui.Toast(rw, ui.Success, "Machine Deleted")
	server.Redirect(rw, r, "/machines")
}
