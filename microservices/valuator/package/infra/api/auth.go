package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(h.secret), nil
		})
		fmt.Println("ASDASD2", token.Raw)

		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID := claims["sub"].(string)
		country := claims["c"].(string)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, userCountryKey, country)

		fmt.Println("ASDASD3")
		fmt.Println("ASDASD3.1" + userID)
		fmt.Println("ASDASD3.2" + country)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
