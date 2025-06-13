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
	buildApiRoutes(r, coffeeController, authController, authRepo)
	buildUiRoutes(r, coffeeController, authController, authRepo)

	return r.ServerMux()
}

func buildApiRoutes(
	r server.Router,
	coffeeController coffee.Controller,
	authController auth.Controller,
	authRepo auth.Repository,
) {
	guest := r.Group("/api", auth.IsGuestMiddleware(auth.API, authRepo))
	{
		guest.HandleFunc("POST /auth/login", authController.Login)
	}

	private := r.Group("/api", auth.IsLoggedInMiddleware(auth.API, authRepo))
	{
		private.HandleFunc("POST /me", authController.GetLoggedInUser)

		private.HandleFunc("POST /coffee", coffeeController.CreateCoffee)
		private.HandleFunc("PUT /coffee/{id}", coffeeController.UpdateCoffee)
		private.HandleFunc("POST /coffee/{id}/image", coffeeController.UpdateCoffeeImage)
		private.HandleFunc("DELETE /coffee/{id}", coffeeController.DeleteCoffee)

		private.HandleFunc("DELETE /roaster/{id}", coffeeController.DeleteRoaster)

		private.HandleFunc("POST /flavour", coffeeController.CreateFlavourProfile)
		private.HandleFunc("DELETE /flavour/{id}", coffeeController.DeleteFlavourProfile)
	}

	self := r.Group("/api", auth.AdminOrSelfMiddleware(auth.API, authRepo))
	{
		self.HandleFunc("PATCH /user/{id}", authController.UpdateUser)
		self.HandleFunc("POST /user/{id}/change-password", authController.ChangePassword)
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
	r.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	guest := r.Group("", auth.IsGuestMiddleware(auth.UI, authRepo))
	{
		guest.HandleFunc("GET /login", authController.ViewLogin)
		guest.HandleFunc("POST /login", authController.Login)
	}

	private := r.Group("", auth.IsLoggedInMiddleware(auth.UI, authRepo))
	{
		private.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" && r.URL.Path != "/home" {
				pageData := ui.NewPageData("404 Not Found", "404")
				pageData.User = r.Context().Value("user").(*auth.User)
				ui.Toast(w, ui.Warning, "Not Found")
				w.WriteHeader(http.StatusNotFound)
				ui.RenderUser(w, r, pageData)
				return
			}

			pageData := ui.NewPageData("Home", "home")
			pageData.User = r.Context().Value("user").(*auth.User)
			ui.RenderUser(w, r, pageData)
		})

		private.HandleFunc("GET /roasters", coffeeController.ViewRoasters)
		private.HandleFunc("POST /roasters", coffeeController.CreateRoaster)
		private.HandleFunc("GET /roasters/{id}", coffeeController.ViewRoaster)
		private.HandleFunc("POST /roasters/{id}", coffeeController.UpdateRoaster)
		private.HandleFunc("POST /roasters/{id}/icon", coffeeController.UpdateRoasterImage)

		private.HandleFunc("POST /logout", authController.Logout)
	}
}
