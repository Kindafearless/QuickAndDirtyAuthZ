package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"AuthZ/health"
	"AuthZ/middleware"
	"AuthZ/authentication"
	"AuthZ/resource"
)

// InitializeRoutes configures all of the routing handlers for endpoints within the application.
func (a *App) InitializeRoutes() error {

	// Health
	a.Router.Handle("/health", middleware.Middleware(http.HandlerFunc(health.Health))).Methods("GET", "OPTIONS")

	// Authentication (Pretending to send user to identity server)
	// Resource Endpoints
	authDatabaseConnector := authentication.DatabaseConnector{DB: a.DB, DatabaseName: a.UserTableName}

	a.Router.Handle("/login", middleware.Middleware(http.HandlerFunc(authDatabaseConnector.Login))).Methods("POST", "OPTIONS")
	a.Router.Handle("/users/list", middleware.Middleware(http.HandlerFunc(authDatabaseConnector.List))).Methods("GET", "OPTIONS")

	// Resource Endpoints
	resourceCommonLib := resource.CommonLib{DB: a.DB, DatabaseName: a.DataTableName}

	a.Router.Handle("/resource/get/{id}", middleware.Middleware(http.HandlerFunc(resourceCommonLib.Get))).Methods("GET", "OPTIONS")
	a.Router.Handle("/resource/delete/{id}", middleware.Middleware(http.HandlerFunc(resourceCommonLib.Delete))).Methods("DELETE", "OPTIONS")
	a.Router.Handle("/resource/list", middleware.Middleware(http.HandlerFunc(resourceCommonLib.List))).Methods("GET", "OPTIONS")
	a.Router.Handle("/resource/discretion", middleware.Middleware(http.HandlerFunc(resourceCommonLib.ListWithDiscretion))).Methods("GET", "OPTIONS")

	return a.ReportRoutes()
}

func (a *App) ReportRoutes() error {
	// To make our lives easier: Iterates through all the routes and prints to console
	err := a.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		// Reports active endpoints on initialization
		fmt.Println(t)

		return nil
	})

	return err
}

// Run executes the configuration of the web service on specified port
func (a *App) Run(port string) error {

	headersOk := handlers.AllowedHeaders([]string{
		"Authorization",
		"Accept",
		"Accept-Language",
		"Allow",
		"Content-Type",
		"X-Requested-With",
		"X-Forwarded-For",
		"X-Forwarded-Host",
		"X-Forwarded-Proto",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Credentials"})
	exposeHeaders := handlers.ExposedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	credentialsOk := handlers.AllowCredentials()
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "DELETE", "PUT", "OPTIONS"})

	fmt.Println("AuthN/Z Sample API Started on port: ", port)

	// Serve the application unencrypted
	err := http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk, exposeHeaders)(a.Router))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
