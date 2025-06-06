package internal

import (
	"net/http"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/server"
)

func BuildRoutes(
	r server.Router,
	coffeeController coffee.Controller,
	authController auth.Controller,
	authRepo auth.Repository,
) *http.ServeMux {
	guest := r.Group("/api", auth.IsGuestMiddleware)
	{
		guest.HandleFunc("POST /auth/login", authController.Login)
	}

	private := r.Group("/api", auth.IsLoggedInMiddleware(authRepo))
	{
		private.HandleFunc("POST /me", authController.GetLoggedInUser)

		private.HandleFunc("POST /coffee", coffeeController.CreateCoffee)
		private.HandleFunc("PUT /coffee/{id}", coffeeController.UpdateCoffee)
		private.HandleFunc("DELETE /coffee/{id}", coffeeController.DeleteCoffee)

		private.HandleFunc("POST /roaster", coffeeController.CreateRoaster)
		private.HandleFunc("PUT /roaster/{id}", coffeeController.UpdateRoaster)
		private.HandleFunc("DELETE /roaster/{id}", coffeeController.DeleteRoaster)

		private.HandleFunc("POST /flavour", coffeeController.CreateFlavourProfile)
		private.HandleFunc("PUT /flavour/{id}", coffeeController.UpdateFlavourProfile)
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

	return r.ServerMux()
}
