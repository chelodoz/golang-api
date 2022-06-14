package middleware

import (
	"errors"
	"fmt"
	"golang-api/dto"
	"golang-api/util"

	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

type JwtMiddleware struct {
	config util.Config
}

func NewJwtMiddleware(config util.Config) *JwtMiddleware {
	return &JwtMiddleware{config}
}

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func (middleware *JwtMiddleware) AuthorizeJWT() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(authorizationHeaderKey)

			if len(authorizationHeader) == 0 {
				err := errors.New("authorization header is not provided")
				dto.WriteResponse(rw, http.StatusUnauthorized, dto.ServiceError{Message: err.Error()})
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				err := errors.New("invalid authorization header format")
				dto.WriteResponse(rw, http.StatusUnauthorized, dto.ServiceError{Message: err.Error()})
				return
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != authorizationTypeBearer {
				err := fmt.Errorf("unsupported authorization type %s", authorizationType)
				dto.WriteResponse(rw, http.StatusUnauthorized, dto.ServiceError{Message: err.Error()})
				return
			}

			accessToken := fields[1]
			_, err := util.VerifyToken(accessToken, middleware.config.JWTSecretKey)
			if err != nil {
				dto.WriteResponse(rw, http.StatusUnauthorized, dto.ServiceError{Message: err.Error()})
				return
			}
			next.ServeHTTP(rw, r)
		})
	}
}
