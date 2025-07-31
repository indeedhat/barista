package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type viewRoastersData struct {
	ui.PageData
	Roasters []coffee.Roaster
	Open     bool
}

func (c Controller) ViewRoasters(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	roasters := c.repo.IndexRoastersForUser(user)

	pageData := viewRoastersData{PageData: ui.NewPageData("Roasters", "roasters", user)}
	pageData.Roasters = roasters
	ui.RenderUser(rw, r, pageData)
}

type viewRoasterData struct {
	ui.PageData
	Roaster *coffee.Roaster
}

func (c Controller) ViewRoaster(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Roaster Not Found", "404", user))
		return
	}

	roaster, err := c.repo.FindRoaster(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Roaster Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Roaster Not Found", "404", user))
		return
	}

	pageData := viewRoasterData{PageData: ui.NewPageData("Roaster", "roaster", user)}
	pageData.Roaster = roaster
	ui.RenderUser(rw, r, pageData)
}
