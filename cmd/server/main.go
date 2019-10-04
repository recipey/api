package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // import so the init() function runs to setup postgres driver

	"github.com/recipey/api/postgres"
	"github.com/recipey/api/server"
	"github.com/recipey/api/user_registration"
)

// TODO: impl server command
// 1. orchestrates all dependency injection here
// 2. be sure to write subpackages to follow hexagon DI pattern
// 3 subpackages should be able to work on DI args passed in and implement the interface how it sees fit

// context pkgs define a service for callers to use
// routes define routing for an endpoint to another
// main will define a root route to use routes provided by other contexts
// pass in args to define what the repository will use for db conn
func main() {
	// TODO: use env vars
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		`db`, `recipey`, `recipey`, `recipey_dev`, `disable`)
	db, _ := sql.Open("postgres", psqlInfo)
	defer db.Close()

	var (
		usersRepo, _ = postgres.NewUserRepository(db)
	)

	// UserRegistration service
	userRegistrationService := userregistration.NewUserRegistrationService(usersRepo)

	// Build a server instance with dependencies
	recipeyServer := server.New(userRegistrationService)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      recipeyServer.Router,
	}

	// Run server in goroutine so it doesn't block
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Recipey API is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// Block until signal received
	<-sc

	// Create a deadline to wait for
	gracefulTimeout := time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	// Doesn't block if srv has no connections
	srv.Shutdown(ctx)
	log.Println("Shutting down server...")
	os.Exit(0)
}
