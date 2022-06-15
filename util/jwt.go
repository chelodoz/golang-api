package util

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken         = errors.New("token is invalid")
	ErrExpiredToken         = errors.New("token has expired")
	minSecretKeySize        = 32
	authorizationTypeBearer = "bearer"
)

//CreateToken creates a new token for a specific username and duration
func CreateToken(username string, duration time.Duration, secretKey string) (string, *JWTPayload, error) {
	if err := secretKeyValidation(secretKey); err != nil {
		return "", nil, err
	}
	// Set custom and standard claims
	jwtPayload, err := newJWTPayload(username, duration)
	if err != nil {
		return "", nil, err
	}
	// Create token with claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)

	// Generate encoded token using the secret signing key
	token, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, err
	}
	return token, jwtPayload, nil
}

//VerifyToken checks if the token is valid or not
func VerifyToken(token string, secretKey string) (*JWTPayload, error) {
	if err := secretKeyValidation(secretKey); err != nil {
		return nil, err
	}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
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

func ValidateBearerHeader(authorizationHeader string) (string, error) {
	if len(authorizationHeader) == 0 {
		return "", errors.New("authorization header is not provided")
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		return "", errors.New("invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {

		return "", fmt.Errorf("unsupported authorization type %s", authorizationType)
	}
	return fields[1], nil
}

func secretKeyValidation(secretKey string) error {
	if len(secretKey) < minSecretKeySize {
		return fmt.Errorf("invalid key size:must be at least %d characters", minSecretKeySize)
	}
	return nil
}

type JWTPayload struct {
	ID        uuid.UUID
	UserEmail string
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func (jwtPayload *JWTPayload) Valid() error {
	if time.Now().After(jwtPayload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func newJWTPayload(userEmail string, duration time.Duration) (*JWTPayload, error) {
	tokenID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	jwtPayload := &JWTPayload{
		ID:        tokenID,
		UserEmail: userEmail,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return jwtPayload, nil
}
