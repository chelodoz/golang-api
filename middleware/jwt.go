package middleware

import (
	"golang-api/dto"
	"golang-api/util"

	"net/http"
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
			authorizationHeader := r.Header.Get("authorization")

			accessToken, err := util.ValidateBearerHeader(authorizationHeader)

			if err != nil {
				dto.WriteResponse(rw, http.StatusUnauthorized, dto.ServiceError{Message: err.Error()})
				return
			}

			_, err = util.VerifyToken(accessToken, middleware.config.JWTSecretKey)
			if err != nil {
				dto.WriteResponse(rw, http.StatusUnauthorized, dto.ServiceError{Message: err.Error()})
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
