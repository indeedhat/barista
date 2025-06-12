package internal

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/indeedhat/barista/assets"
	"github.com/indeedhat/barista/assets/templates"
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
	buildApiRoutes(r, coffeeController, authController, authRepo)
	buildUiRoutes(r, authRepo)

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

func buildUiRoutes(r server.Router, authRepo auth.Repository) {
	r.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.Public))))
	r.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	tmpl := template.New("")
	tmpls := template.Must(tmpl.Funcs(template.FuncMap{
		"embed": func(name string, data any) template.HTML {
			var out strings.Builder
			if err := tmpl.ExecuteTemplate(&out, name, data); err != nil {
				log.Println(err)
			}
			return template.HTML(out.String())
		},
	}).ParseFS(templates.FS, "layouts/*", "pages/*"))

	guest := r.Group("", auth.IsGuestMiddleware(auth.UI, authRepo))
	{
		guest.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
			log.Print("/login")
			log.Print(tmpls.ExecuteTemplate(w, "layouts/guest", map[string]any{
				"Page": "pages/login",
			}))
		})
	}

	private := r.Group("", auth.IsLoggedInMiddleware(auth.UI, authRepo))
	{
		private.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			log.Print("/")
			log.Print(tmpls.ExecuteTemplate(w, "layouts/user", map[string]any{
				"Page": "pages/home",
			}))
		})
	}
}
