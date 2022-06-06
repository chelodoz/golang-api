package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrExpiredToken  = errors.New("token has expired")
	minSecretKeySize = 32
)

// JWTService is an interface for managing tokens
type JWTService interface {
	//CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	//VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*JWTPayload, error)
}

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) (JWTService, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size:must be at least %d characters", minSecretKeySize)
	}
	return &jwtService{secretKey}, nil
}

//CreateToken creates a new token for a specific username and duration
func (jwtSrv *jwtService) CreateToken(username string, duration time.Duration) (string, error) {

	// Set custom and standard claims
	jwtPayload, err := newJWTPayload(username, duration)
	if err != nil {
		return "", err
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)

	// Generate encoded token using the secret signing key
	return token.SignedString([]byte(jwtSrv.secretKey))
}

//VerifyToken checks if the token is valid or not
func (jwtSrv *jwtService) VerifyToken(token string) (*JWTPayload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwtSrv.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	jwtPayload, ok := jwtToken.Claims.(*JWTPayload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return jwtPayload, nil
}

type JWTPayload struct {
	ID        uuid.UUID `json:id`
	Username  string    `json:username`
	IssuedAt  time.Time `json:issuedAt`
	ExpiredAt time.Time `json:expiredAt`
}

func (jwtPayload *JWTPayload) Valid() error {
	if time.Now().After(jwtPayload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func newJWTPayload(username string, duration time.Duration) (*JWTPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	jwtPayload := &JWTPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return jwtPayload, nil
}
