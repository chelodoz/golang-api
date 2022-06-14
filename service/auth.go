package service

import (
	"context"
	"golang-api/entity"
	"golang-api/util"
)

type AuthService interface {
	Login(email string, password string) bool
	Logout(ctx context.Context, email string) error
	CreateTokens(ctx context.Context, email string, prevTokenID string) (*entity.TokenDetails, error)
}

type authService struct {
	userRepository  entity.UserRepository
	tokenRepository entity.TokenRepository
	config          util.Config
}

func NewAuthService(userRepository entity.UserRepository, tokenRepository entity.TokenRepository, config util.Config) AuthService {
	return &authService{
		userRepository,
		tokenRepository,
		config,
	}
}

func (authService *authService) CreateTokens(ctx context.Context, email string, prevTokenID string) (*entity.TokenDetails, error) {
	accessToken, accessJwtPayload, err := util.CreateToken(email, authService.config.AccessTokenDuration, authService.config.JWTSecretKey)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshJwtPayload, err := util.CreateToken(email, authService.config.RefreshTokenDuration, authService.config.JWTSecretKey)
	if err != nil {
		return nil, err
	}

	if err := authService.tokenRepository.SetRefreshToken(ctx, email, refreshToken, refreshJwtPayload.ExpiredAt); err != nil {
		return nil, err
	}

	return &entity.TokenDetails{
		SessionUuid:           refreshJwtPayload.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessJwtPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshJwtPayload.ExpiredAt,
	}, nil

}

func (authService *authService) Logout(ctx context.Context, email string) error {
	return authService.tokenRepository.DeleteUserRefreshTokens(ctx, email)
}

func (service *authService) Login(email string, password string) bool {
	user, err := service.userRepository.GetUserByEmail(email)
	if err != nil {
		return false
	}

	err = util.CheckPassword(password, user.Password)
	return err == nil
}
