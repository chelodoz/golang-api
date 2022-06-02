package main

import (
	"golang-api/database"
	"golang-api/handler"
	"golang-api/repository"
	"golang-api/router"
	"golang-api/service"
	"golang-api/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	db, _ := database.DBConnection(config)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	httpRouter := router.NewMuxRouter()
	httpRouter.POST("/user", userHandler.CreateUser)
	httpRouter.SERVE(config)
}
