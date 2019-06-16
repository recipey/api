package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	// "github.com/recipey/api/user_registration"
	"github.com/recipey/api"
	"github.com/recipey/api/postgres"
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
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		`db`, `recipey`, `recipey`, `recipey_dev`, `disable`)
	db, _ := sql.Open("postgres", psqlInfo)
	defer db.Close()

	var (
		users, _ = postgres.NewUserRepository(db)
	)

	fmt.Println("INSERT NEW USER")
	user := api.NewUser("foo", "foo@bar.com")
	users.Store(user)
	fmt.Println("DONE INSERT NEW USER")

	fmt.Println("ALL USERS!!")
	fmt.Println(users.FindAll())
	fmt.Println("DONE ALL USERS!!")

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Recipey API is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
