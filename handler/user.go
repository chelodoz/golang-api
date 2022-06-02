package handler

import (
	"encoding/json"
	"golang-api/dto"
	"golang-api/service"
	"net/http"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

type UserHandler interface {
	CreateUser(rw http.ResponseWriter, r *http.Request)
	// DeleteUser(rw http.ResponseWriter, r *http.Request)
	// GetUser(rw http.ResponseWriter, r *http.Request)
	// GetUsers(rw http.ResponseWriter, r *http.Request)
	// UpdateUser(rw http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	validate = validator.New()
	return &userHandler{
		service: service,
	}
}

func (u *userHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	var createUserRequest dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(dto.ServiceError{Message: err.Error()})
		return
	}

	err := validate.Struct(&createUserRequest)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(dto.ServiceError{Message: err.Error()})
		return
	}

	_, err = u.service.CreateUser(createUserRequest)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(dto.ServiceError{Message: err.Error()})
		return
	}

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(dto.ServiceError{Message: "Success!"})
}
