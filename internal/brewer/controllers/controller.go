package brewer_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/brewer"
	"github.com/indeedhat/barista/internal/types"
	"github.com/indeedhat/barista/internal/ui"
)

const BrewerImagePath = "data/uploads/brewer/"

type Controller struct {
	repo brewer.Repository
}

func New(repo brewer.Repository) Controller {
	return Controller{repo}
}

func (c Controller) findEspressoBrewer(rw http.ResponseWriter, user *auth.User, id uint) *brewer.Brewer {
	if id == 0 {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return nil
	}

	brewer, err := c.repo.FindBrewer(id, user.ID)
	if err != nil {
		ui.Toast(rw, ui.Warning, "Brewer not found")
		return nil
	}

	if brewer.Type != types.BrewerEspresso {
		ui.Toast(rw, ui.Warning, "Only espresso machines can have baskets")
		return nil
	}

	return brewer
}
