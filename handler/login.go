package handler

import (
	"encoding/json"
	"golang-api/dto"
	"golang-api/service"
	"golang-api/util"
	"net/http"
)

type LoginHandler interface {
	Login(rw http.ResponseWriter, r *http.Request)
}

type loginHandler struct {
	loginService service.LoginService
	jWtService   service.JWTService
	config       util.Config
}

func NewLoginHandler(loginService service.LoginService, jWtService service.JWTService, config util.Config) LoginHandler {
	return &loginHandler{
		loginService: loginService,
		jWtService:   jWtService,
		config:       config,
	}
}

func (handler *loginHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}

	if err := validate.Struct(&loginRequest); err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}
	isAuthenticated := handler.loginService.Login(loginRequest.Email, loginRequest.Password)
	if !isAuthenticated {
		dto.WriteResponse(rw, http.StatusUnauthorized, dto.ServiceError{Message: "Unauthorized"})
		return
	}

	token, err := handler.jWtService.CreateToken(loginRequest.Email, handler.config.AccessTokenDuration)

	if err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}

	dto.WriteResponse(rw, http.StatusOK, dto.TokenResponse{Token: token})
}
