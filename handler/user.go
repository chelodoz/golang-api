package handler

import (
	"encoding/json"
	"golang-api/dto"
	"golang-api/service"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

type UserHandler interface {
	CreateUser(rw http.ResponseWriter, r *http.Request)
	// DeleteUser(rw http.ResponseWriter, r *http.Request)
	GetUser(rw http.ResponseWriter, r *http.Request)
	GetUsers(rw http.ResponseWriter, r *http.Request)
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

func (u *userHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {

	users, err := u.service.GetUsers()
	if err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}

	dto.WriteResponse(rw, http.StatusOK, users)
}

func (u *userHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	userId := getUserID(r)
	user, err := u.service.GetUser(userId)
	if err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}

	dto.WriteResponse(rw, http.StatusOK, user)
}

func (u *userHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	var createUserRequest dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}

	if err := validate.Struct(&createUserRequest); err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}

	user, err := u.service.CreateUser(createUserRequest)
	if err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}

	dto.WriteResponse(rw, http.StatusCreated, user)
}

func getUserID(r *http.Request) uint {
	vars := mux.Vars(r)
	// convert the id into an integer and return
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(id)
}
