package handler

import (
	"encoding/json"
	"golang-api/dto"
	"golang-api/service"
	"golang-api/util"
	"net/http"
)

type AuthHandler interface {
	Login(rw http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService service.AuthService
	config      util.Config
}

func NewAuthHandler(authService service.AuthService, config util.Config) AuthHandler {
	return &authHandler{
		authService,
		config,
	}
}

func (handler *authHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}

	if err := validate.Struct(&loginRequest); err != nil {
		dto.WriteResponse(rw, http.StatusBadRequest, dto.ServiceError{Message: err.Error()})
		return
	}
	isAuthenticated := handler.authService.Login(loginRequest.Email, loginRequest.Password)
	if !isAuthenticated {
		dto.WriteResponse(rw, http.StatusUnauthorized, dto.ServiceError{Message: "Unauthorized"})
		return
	}

	tokenDetails, err := handler.authService.CreateTokens(r.Context(), loginRequest.Email, "")

	if err != nil {
		dto.WriteResponse(rw, http.StatusInternalServerError, dto.ServiceError{Message: "Internal Server Error"})
	}
	dto.WriteResponse(rw, http.StatusOK, dto.LoginResponse{
		AccessToken:           tokenDetails.AccessToken,
		AccessTokenExpiresAt:  tokenDetails.AccessTokenExpiresAt,
		RefreshToken:          tokenDetails.RefreshToken,
		RefreshTokenExpiresAt: tokenDetails.RefreshTokenExpiresAt,
		Email:                 loginRequest.Email,
	})
}
