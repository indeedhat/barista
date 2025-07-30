package brewer_controllers

import (
	"net/http"
	"strconv"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/ui"
)

func (c Controller) BasketsSelect(rw http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("brewer.int")
	if idStr == "" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	user := r.Context().Value("user").(*auth.User)
	brewerId, _ := strconv.Atoi(idStr)
	brewer, err := c.repo.FindBrewer(uint(brewerId), user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return
	}
	if brewer.Type != "Espresso" {
		return
	}

	comData := ui.NewComponentData("baskets-select", ui.ComponentData{
		"value": r.URL.Query().Get("value"),
	})
	defer func() {
		ui.RenderComponent(rw, comData)
	}()

	comData["Baskets"] = brewer.Baskets
}
