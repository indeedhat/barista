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
	"github.com/indeedhat/barista/internal/brewer"
	"github.com/indeedhat/barista/internal/brewer/controllers"
	"github.com/indeedhat/barista/internal/coffee"
	"github.com/indeedhat/barista/internal/database"
	"github.com/indeedhat/barista/internal/server"
	_ "github.com/indeedhat/dotenv/autoload"
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
		coffee.Recipe{},
		auth.User{},
		brewer.Brewer{},
		brewer.Basket{},
	)

	authRepo := auth.NewSqliteRepo(db)
	coffeeRepo := coffee.NewSqliteRepo(db)
	brewerRepo := brewer.NewSqliteRepo(db)

	authController := auth.NewController(authRepo)
	coffeeController := coffee.NewController(coffeeRepo)
	brewerController := brewer_controllers.New(brewerRepo)

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
		brewerController,
		authRepo,
	)

	svr := &http.Server{
		Addr:    ":8087",
		Handler: mux,
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
