package service

import (
	"context"
	"encoding/json"
	"net/http"

	"gofemart/internal/crypto"
	"gofemart/internal/model"
)

// var (
// 	// errInvalidNameOrPassword = errors.New("invalid name or password")
// 	// errWrongSignInToken      = errors.New("wrong sign in token")
// )

func (s *Service) SignUp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userReq model.User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, _ := s.store.Users().FindByName(userReq.Username)
		if user != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
		hashedPassword, err := crypto.HashPassword(userReq.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userReq.Password = hashedPassword
		err = s.store.Users().Create(&userReq)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		accessToken, refreshToken, err := s.JwtManager.GenerateTokens(context.Background(), userReq.ID, s.store)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:  "authorization",
			Value: "Bearer " + accessToken,
		}
		http.SetCookie(w, &cookie)
		cookie = http.Cookie{
			Name:  "refresh",
			Value: "Bearer " + refreshToken,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Service) SignIn() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userReq model.User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := s.store.Users().FindByName(userReq.Username)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user == nil || !crypto.IsPasswordCorrect(user.Password, userReq.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		accessToken, refreshToken, err := s.JwtManager.GenerateTokens(context.Background(), user.ID, s.store)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:  "authorization",
			Value: "Bearer " + accessToken,
		}
		http.SetCookie(w, &cookie)
		cookie = http.Cookie{
			Name:  "refresh",
			Value: "Bearer " + refreshToken,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Service) Refresh(ctx context.Context, refresh string) (*string, *string, error) {
	claims, err := s.JwtManager.VerifyToken(ctx, refresh)
	if err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, err := s.JwtManager.GenerateTokens(ctx, claims.UserID, s.store)
	if err != nil {
		return nil, nil, err
	}

	return &accessToken, &refreshToken, nil
}
