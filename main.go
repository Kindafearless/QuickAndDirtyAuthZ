package main

import (
	"QuickAndDirtyAuthZ/database"
	"github.com/gorilla/mux"
)

// App represents any parameters used in the global application scope.
type App struct {
	Router *mux.Router
	DB     *database.InMemory
	UsersDB     *database.InMemory
	DataTableName string
	UserTableName string
}

func main() {

	// Initialize an in-memory database for use
	schema := database.CreateSchema("data", "users")
	db, err := database.CreateDatabase(schema)
	if err != nil {
		panic(err)
	}

	// Setup application objects
	im := database.InMemory{DB: db}

	a := App{
		Router: mux.NewRouter(),
		DB: &im,
		DataTableName: "data",
		UserTableName: "users",
	}

	// Seed the database with values to play with
	// Seed the sample data table
	err = a.SeedData()
	if err != nil {
		panic(err)
	}

	// Seed the sample user table
	err = a.SeedUsers()
	if err != nil {
		panic(err)
	}

	err = a.InitializeRoutes()
	if err != nil {
		panic(err)
	}

	err = a.Run(":7777")
	if err != nil {
		panic(err)
	}
}
