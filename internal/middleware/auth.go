package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtg "github.com/ashershnyov/gophermart-loyalty-program/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth checks the JWT token for validity.
func JWTAuth(jwtGen *jwtg.Generator) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing auth header", http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "bad authorization header format", http.StatusUnauthorized)
				return
			}

			tok, err := jwtGen.ParseToken(tokenParts[1])
			if err != nil || !tok.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := tok.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userIDfloat, ok := claims[jwtg.UserIDKey].(float64)
			if !ok {
				http.Error(w, "no user id in token", http.StatusUnauthorized)
				return
			}

			userID := int64(userIDfloat)

			ctx := context.WithValue(r.Context(), jwtg.CtxKey, userID)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
