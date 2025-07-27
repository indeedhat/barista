package internal

import (
	"net/http"

	"github.com/indeedhat/barista/assets"
	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui"
)

func BuildRoutes(
	r server.Router,
	coffeeController coffee.Controller,
	authController auth.Controller,
	authRepo auth.Repository,
) *http.ServeMux {
	buildApiRoutes(r, authController, authRepo)
	buildUiRoutes(r, coffeeController, authController, authRepo)

	return r.ServerMux()
}

func buildApiRoutes(
	r server.Router,
	authController auth.Controller,
	authRepo auth.Repository,
) {
	// TODO: these are all unimplemented in the new htmx based ui
	//       once they have this whole function will be removed
	guest := r.Group("/api", auth.IsGuestMiddleware(auth.API, authRepo))
	{
		guest.HandleFunc("POST /auth/login", authController.Login)
	}

	private := r.Group("/api", auth.IsLoggedInMiddleware(auth.API, authRepo))
	{
		private.HandleFunc("POST /me", authController.GetLoggedInUser)
	}

	self := r.Group("/api", auth.AdminOrSelfMiddleware(auth.API, authRepo))
	{
		self.HandleFunc("PATCH /user/{id}", authController.UpdateUser)
		self.HandleFunc("POST /user/{id}/force-logout", authController.ForceLogoutUser)
	}

	admin := r.Group("/api", auth.UserHasPermissionMiddleware(auth.API, auth.LevelAdmin, authRepo))
	{
		admin.HandleFunc("POST /user", authController.CreateUser)
	}
}

func buildUiRoutes(
	r server.Router,
	coffeeController coffee.Controller,
	authController auth.Controller,
	authRepo auth.Repository,
) {
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

		private.HandleFunc("POST /logout", authController.Logout)
	}
}
