package coffee

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

const (
	CoffeeImagePath  = "data/uploads/coffee/"
	RoasterImagePath = "data/uploads/roaster/"
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

type flavoursData struct {
	ui.PageData
	Flavours []FlavourProfile
	Open     bool
}

func (c Controller) ViewFlavours(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := flavoursData{PageData: ui.NewPageData("Flavours", "flavours", user)}

	pageData.Form = createFlavourProfileRequest{}
	pageData.Flavours = c.repo.IndexFlavourProfiles()

	ui.RenderUser(rw, r, pageData)
}

type createFlavourProfileRequest struct {
	Name string `json:"name" validate:"required"`
}

func (c Controller) CreateFlavourProfile(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*auth.User)
	pageData := flavoursData{PageData: ui.NewPageData("Flavours", "flavours", user)}
	pageData.Flavours = c.repo.IndexFlavourProfiles()
	pageData.Open = true

	var req createFlavourProfileRequest
	if err := server.UnmarshalBody(r, &req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "The server did not understand the request")
		ui.RenderUser(rw, r, pageData)
		return
	}

	if err := server.ValidateRequest(req, &pageData); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create flavour")
		ui.RenderUser(rw, r, pageData)
		return
	}

	flavour := FlavourProfile{
		Name: req.Name,
	}

	if err := c.repo.SaveFlavourProfile(&flavour); err != nil {
		ui.Toast(rw, ui.Warning, "Failed to create flavour")
		ui.RenderUser(rw, r, pageData)
		return
	}

	pageData.Flavours = c.repo.IndexFlavourProfiles()
	pageData.Open = false
	pageData.Form = createFlavourProfileRequest{}

	ui.Toast(rw, ui.Success, "Flavour created")
	ui.RenderUser(rw, r, pageData)
}
