package jwt

import (
	"context"
	"errors"
	"gofemart/internal/model"
	"gofemart/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	errUnexpectedSigningMethod = errors.New("unexpected signing method")
	errInvalidToken            = errors.New("invalid token")
)

type JWTManager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTManager(secretKey, accessTokenDurationString, refreshTokenDurationString string) (*JWTManager, error) {
	accessTokenDuration, err := time.ParseDuration(accessTokenDurationString)
	if err != nil {
		return nil, err
	}

	refreshTokenDuration, err := time.ParseDuration(refreshTokenDurationString)
	if err != nil {
		return nil, err
	}

	return &JWTManager{
		secretKey:            secretKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}, nil
}

func (manager *JWTManager) GenerateTokens(ctx context.Context, id int, store storage.Store) (accessTokenString, refreshTokenString string, err error) {
	accessExpirationTime := time.Now().Add(manager.accessTokenDuration)
	// TODO: check if user/team mode. Then get team data if team mode
	accessClaims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		},
		UserID: id,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err = accessToken.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", "", err
	}

	// refresh token
	refreshExpirationTime := time.Now().Add(manager.refreshTokenDuration)

	refreshClaims := &model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
		UserID: id,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err = refreshToken.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (manager *JWTManager) VerifyToken(ctx context.Context, tokenString string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSigningMethod
		}
		return []byte(manager.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(*model.UserClaims); ok && token.Valid {
		return token.Claims.(*model.UserClaims), nil
	}

	return nil, errInvalidToken
}
