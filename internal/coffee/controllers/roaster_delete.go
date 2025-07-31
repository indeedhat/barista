package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type deleteRoasterData struct {
	ui.PageData
	Roaster *coffee.Roaster
}

func (c Controller) DeleteRoaster(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := deleteRoasterData{PageData: ui.NewPageData("Roaster", "roaster", user)}
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	roaster, err := c.repo.FindRoaster(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster not found")
		return
	}

	pageData.Title = roaster.Name
	pageData.Roaster = roaster

	if len(roaster.Coffees) > 0 {
		ui.Toast(rw, ui.Warning, "Roaster cannot be deleted while it still has coffees")
		return
	}

	if err := c.repo.DeleteRoaster(roaster); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete roaster")
		return
	}

	ui.Toast(rw, ui.Success, "Roaster deleted")
	server.Redirect(rw, r, "/roasters")
}
