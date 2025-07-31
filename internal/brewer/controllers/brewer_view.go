package brewer_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/brewer"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) ViewBrewers(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	pageData := ui.NewPageData("Brewers", "brewers", user)
	pageData.Form = createBrewerRequest{}

	ui.RenderUser(rw, r, pageData)
}

type viewBrewerData struct {
	ui.PageData
	Brewer brewer.Brewer
}

func (c Controller) ViewBrewer(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Brewer Not Found", "404", user))
		return
	}

	brewer, err := c.repo.FindBrewer(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer Not Found")
		ui.RenderUser(rw, r, ui.NewPageData("Brewer Not Found", "404", user))
		return
	}

	pageData := viewBrewerData{
		ui.NewPageData("Brewer", "brewer", user),
		*brewer,
	}
	pageData.Form = updateBrewerRequest{}

	ui.RenderUser(rw, r, pageData)
}
