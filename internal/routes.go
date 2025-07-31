package internal

import (
	"net/http"

	"github.com/indeedhat/barista/assets"
	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/auth/controllers"
	"github.com/indeedhat/barista/internal/brewer/controllers"
	"github.com/indeedhat/barista/internal/coffee/controllers"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func BuildRoutes(
	r server.Router,
	coffeeController coffee_controllers.Controller,
	authController auth_controllers.Controller,
	brewerController brewer_controllers.Controller,
	authRepo auth.Repository,
) *http.ServeMux {
	r.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.Public))))

	guest := r.Group("", auth.IsGuestMiddleware(auth.UI, authRepo))
	{
		guest.HandleFunc("GET /login", authController.ViewLogin)
		guest.HandleFunc("POST /login", authController.Login)

		guest.HandleFunc("GET /register", authController.ViewRegister)
		guest.HandleFunc("POST /register", authController.Register)
	}

	private := r.Group("", auth.IsLoggedInMiddleware(auth.UI, authRepo))
	{
		private.Handle("GET /uploads/",
			http.StripPrefix("/uploads/", http.FileServer(http.Dir("data/uploads"))),
		)

		private.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" && r.URL.Path != "/home" {
				ui.Toast(w, ui.Warning, "Not Found")
				w.WriteHeader(http.StatusNotFound)
				ui.RenderUser(w, r,
					ui.NewPageData("404 Not Found", "404", r.Context().Value("user").(*auth.User)),
				)
				return
			}

			coffeeController.ViewRecipes(w, r)
		})

		private.HandleFunc("GET /user/settings", authController.ViewSettings)
		private.HandleFunc("POST /user/change-password", authController.ChangePassword)

		private.HandleFunc("GET /coffees", coffeeController.ViewCoffees)
		private.HandleFunc("POST /coffees", coffeeController.CreateCoffee)
		private.HandleFunc("GET /coffees/{id}", coffeeController.ViewCoffee)
		private.HandleFunc("PUT /coffees/{id}", coffeeController.UpdateCoffee)
		private.HandleFunc("POST /coffees/{id}/icon", coffeeController.UpdateCoffeeImage)
		private.HandleFunc("DELETE /coffee/{id}", coffeeController.DeleteCoffee)

		private.HandleFunc("GET /coffees/{id}/recipes", coffeeController.NewRecipe)
		private.HandleFunc("POST /coffees/{id}/recipes", coffeeController.CreateRecipe)
		private.HandleFunc("PUT /coffees/{coffee_id}/recipes/{recipe_id}", coffeeController.UpdateRecipe)
		private.HandleFunc("DELETE /coffees/{coffee_id}/recipes/{recipe_id}", coffeeController.DeleteRecipe)

		private.HandleFunc("GET /recipes", coffeeController.ViewRecipes)

		private.HandleFunc("GET /flavours", coffeeController.ViewFlavours)
		private.HandleFunc("POST /flavours", coffeeController.CreateFlavourProfile)

		private.HandleFunc("GET /roasters", coffeeController.ViewRoasters)
		private.HandleFunc("POST /roasters", coffeeController.CreateRoaster)
		private.HandleFunc("GET /roasters/{id}", coffeeController.ViewRoaster)
		private.HandleFunc("PUT /roasters/{id}", coffeeController.UpdateRoaster)
		private.HandleFunc("POST /roasters/{id}/icon", coffeeController.UpdateRoasterImage)
		private.HandleFunc("DELETE /roaster/{id}", coffeeController.DeleteRoaster)

		private.HandleFunc("GET /baskets/select", brewerController.BasketsSelect)
		private.HandleFunc("GET /brewers/select", brewerController.BrewersSelect)

		private.HandleFunc("GET /brewers", brewerController.ViewBrewers)
		private.HandleFunc("POST /brewers", brewerController.CreateBrewer)
		private.HandleFunc("GET /brewers/{id}", brewerController.ViewBrewer)
		private.HandleFunc("PUT /brewers/{id}", brewerController.UpdateBrewer)
		private.HandleFunc("POST /brewers/{id}/icon", brewerController.UpdateBrewerImage)
		private.HandleFunc("DELETE /brewers/{id}", brewerController.DeleteBrewer)

		private.HandleFunc("GET /brewers/{id}/baskets", brewerController.NewBasket)
		private.HandleFunc("POST /brewers/{id}/baskets", brewerController.CreateBasket)
		private.HandleFunc("PUT /brewers/{brewer_id}/baskets/{basket_id}", brewerController.UpdateBasket)
		private.HandleFunc("DELETE /brewers/{brewer_id}/baskets/{basket_id}", brewerController.DeleteBasket)

		private.HandleFunc("POST /logout", authController.Logout)
	}

	return r.ServerMux()
}
