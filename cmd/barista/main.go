package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/indeedhat/barista/internal"
	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/database"
	"github.com/indeedhat/barista/internal/server"
	_ "github.com/indeedhat/dotenv/autoload"
	"github.com/rs/cors"
)

func main() {
	firstRun := !database.Exists()

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		coffee.Coffee{},
		coffee.Roaster{},
		coffee.FlavourProfile{},
		auth.User{},
	)

	authRepo := auth.NewSqliteRepo(db)
	coffeeRepo := coffee.NewSqliteRepo(db)

	authController := auth.NewController(authRepo)
	coffeeController := coffee.NewController(coffeeRepo)

	if firstRun {
		if err := authRepo.CreateRootUser(); err != nil {
			log.Fatalf("Failed to create root user: %s", err)
		}
	}

	router := server.NewRouter(server.ServerConfig{
		MaxBodySize: 1 << 20,
	})

	mux := internal.BuildRoutes(
		router,
		coffeeController,
		authController,
		authRepo,
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{server.CorsAllowHost.Get()},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
		ExposedHeaders:   []string{"Auth_token"},
	})

	svr := &http.Server{
		Addr:    ":8087",
		Handler: c.Handler(mux),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("ListenAndServer: %v", svr.ListenAndServe())
	}()

	<-quit
	log.Print("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Print("Server forced to shutdown after timeout")
	}
}
