package dto

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=6"`
}

type LoginResponse struct {
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	Email                 string    `json:"email"`
}

type TokenRequest struct {
	AccessToken  string `json:"accessToken"  binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}
type AccessTokenRequest struct {
	AccessToken string `json:"accessToken"  binding:"required"`
}
