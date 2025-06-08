package internal

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/indeedhat/barista/assets"
	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
	"github.com/indeedhat/barista/internal/ui/pages"
)

func BuildRoutes(
	r server.Router,
	coffeeController coffee.Controller,
	authController auth.Controller,
	authRepo auth.Repository,
) *http.ServeMux {
	buildApiRoutes(r, coffeeController, authController, authRepo)
	buildUiRoutes(r)

	return r.ServerMux()
}

func buildApiRoutes(
	r server.Router,
	coffeeController coffee.Controller,
	authController auth.Controller,
	authRepo auth.Repository,
) {
	guest := r.Group("/api", auth.IsGuestMiddleware)
	{
		guest.HandleFunc("POST /auth/login", authController.Login)
	}

	private := r.Group("/api", auth.IsLoggedInMiddleware(authRepo))
	{
		private.HandleFunc("POST /auth/logout", authController.Logout)

		private.HandleFunc("POST /me", authController.GetLoggedInUser)

		private.HandleFunc("POST /coffee", coffeeController.CreateCoffee)
		private.HandleFunc("PUT /coffee/{id}", coffeeController.UpdateCoffee)
		private.HandleFunc("POST /coffee/{id}/image", coffeeController.UpdateCoffeeImage)
		private.HandleFunc("DELETE /coffee/{id}", coffeeController.DeleteCoffee)

		private.HandleFunc("POST /roaster", coffeeController.CreateRoaster)
		private.HandleFunc("PUT /roaster/{id}", coffeeController.UpdateRoaster)
		private.HandleFunc("POST /roaster/{id}/image", coffeeController.UpdateRoasterImage)
		private.HandleFunc("DELETE /roaster/{id}", coffeeController.DeleteRoaster)

		private.HandleFunc("POST /flavour", coffeeController.CreateFlavourProfile)
		private.HandleFunc("DELETE /flavour/{id}", coffeeController.DeleteFlavourProfile)
	}

	self := r.Group("/api", auth.AdminOrSelfMiddleware(authRepo))
	{
		self.HandleFunc("PATCH /user/{id}", authController.UpdateUser)
		self.HandleFunc("POST /user/{id}/change-password", authController.ChangePassword)
		self.HandleFunc("POST /user/{id}/force-logout", authController.ForceLogoutUser)
	}

	admin := r.Group("/api", auth.UserHasPermissionMiddleware(auth.LevelAdmin, authRepo))
	{
		admin.HandleFunc("POST /user", authController.CreateUser)
	}
}

func buildUiRoutes(r server.Router) {
	r.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.Assets))))
	r.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	r.Handle("GET /login", templ.Handler(pages.Login()))
}
