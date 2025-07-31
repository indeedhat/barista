package coffee_controllers

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/ui"
)

type viewFlavoursData struct {
	ui.PageData
	Flavours []coffee.FlavourProfile
	Open     bool
}

func (c Controller) ViewFlavours(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := viewFlavoursData{PageData: ui.NewPageData("Flavours", "flavours", user)}

	pageData.Form = createFlavourProfileRequest{}
	pageData.Flavours = c.repo.IndexFlavourProfiles()

	ui.RenderUser(rw, r, pageData)
}
