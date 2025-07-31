package brewer_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/brewer"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

type deleteBrewerData struct {
	ui.PageData
	Brewer brewer.Brewer
}

func (c Controller) DeleteBrewer(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := updateBrewerData{PageData: ui.NewPageData("Brewer", "brewer", user)}
	defer func() {
		ui.RenderUser(rw, r, pageData)
	}()

	id, err := server.PathID(r)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return
	}

	brewer, err := c.repo.FindBrewer(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return
	}

	pageData.Title = brewer.Name
	pageData.Brewer = brewer

	if len(brewer.Baskets) > 0 {
		ui.Toast(rw, ui.Warning, "Brewer cannot be deleted while it still has baskets")
		return
	}

	if err := c.repo.DeleteBrewer(brewer); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to delete brewer")
		return
	}

	ui.Toast(rw, ui.Success, "Brewer Deleted")
	server.Redirect(rw, r, "/brewers")
}
