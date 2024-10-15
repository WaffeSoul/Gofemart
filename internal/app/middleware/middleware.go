package middleware

import (
	"gofemart/internal/jwt"
	"gofemart/internal/logger"
	"net/http"
)

func Middleware(h http.Handler) http.Handler {
	return logger.WithLogging(GzipMiddleware(h))
}

func MiddlewareWithJWT(jwtM *jwt.JWTManager, h http.Handler) http.Handler {
	return logger.WithLogging(GzipMiddleware(JWTInterceptor(jwtM, h)))
}
