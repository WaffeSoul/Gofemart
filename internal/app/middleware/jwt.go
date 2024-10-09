package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"gofemart/internal/jwt"
)

func JWTInterceptor(jwtM *jwt.JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authorization")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				w.WriteHeader(http.StatusUnauthorized)
				return
			default:
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}
		tokenString := strings.TrimPrefix(cookie.Value, "Bearer ")

		userClaims, err := jwtM.VerifyToken(context.Background(), tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId := userClaims.UserID
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
