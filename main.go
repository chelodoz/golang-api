package main

import (
	"context"
	"golang-api/database"
	"golang-api/handler"
	"golang-api/middleware"
	"golang-api/repository"
	"golang-api/service"
	"golang-api/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	startHTTPServer(config)
}

func startHTTPServer(config util.Config) {

	db, _ := database.DBConnection(config)

	userRepository := repository.NewUserRepository(db)
	tokenRepository := repository.NewRedisCache(config.CACHEHost, 0, 14000)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userRepository, tokenRepository, config)
	jwtMiddleware := middleware.NewJwtMiddleware(config)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService, config)

	router := mux.NewRouter()

	secure := router.NewRoute().PathPrefix("/secure").Subrouter()
	secure.Use(jwtMiddleware.AuthorizeJWT())
	secure.HandleFunc("/users", userHandler.GetUsers).Methods(http.MethodGet)
	secure.HandleFunc("/users", userHandler.CreateUser).Methods(http.MethodPost)
	secure.HandleFunc("/users", userHandler.DeleteUser).Methods(http.MethodDelete)
	secure.HandleFunc("/users/{id}", userHandler.GetUser).Methods(http.MethodGet)
	secure.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods(http.MethodPatch)

	router.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	// CORS
	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	log := log.New(os.Stdout, "golang-api ", log.LstdFlags)
	// create a new server
	server := http.Server{
		Addr:         config.HTTPServerAddress, // configure the bind address
		Handler:      cors(router),             // set the default handler
		ErrorLog:     log,                      // set the logger for the server
		ReadTimeout:  5 * time.Second,          // max time to read request from the client
		WriteTimeout: 10 * time.Second,         // max time to write response to the client
		IdleTimeout:  120 * time.Second,        // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		log.Printf("Starting server on port: %v", config.HTTPServerAddress)

		err := http.ListenAndServe(config.HTTPServerAddress, router)
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-ch
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)

}
