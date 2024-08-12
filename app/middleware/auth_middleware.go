package middleware

import (
	"context"
	"fmt"
	"kinopoisk/app/delivery"
	userusecase "kinopoisk/app/users/usecase"
	"log"
	"net/http"
	"strings"
)

type userKey int
type tokenKey int
type loggerKey int

const (
	MyUserKey   userKey   = 1
	MyTokenKey  tokenKey  = 2
	MyLoggerKey loggerKey = 3
)

func AuthMiddleware(uc userusecase.UserUseCase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, err := GetLoggerFromContext(r.Context())
		if err != nil {
			log.Printf("can not get logger from context: %s", err)
			WriteNoLoggerResponse(w)
		}
		logger.Infof("auth middleware start")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			delivery.WriteResponse(logger, w, []byte(`{"message": "there is no access token or it is in wrong format"}`),
				http.StatusUnauthorized)
			return
		}
		tokenValue := strings.TrimPrefix(authHeader, "Bearer ")
		mySession, err := uc.GetSession(tokenValue, logger)
		if err != nil || mySession.ID == "" {
			errText := fmt.Sprintf(`{"message": "there is no session for token %s}`, tokenValue)
			delivery.WriteResponse(logger, w, []byte(errText), http.StatusUnauthorized)
			return
		}
		sessionUser := mySession.User
		ctx := r.Context()
		ctx = context.WithValue(ctx, MyUserKey, sessionUser)
		ctx = context.WithValue(ctx, MyTokenKey, tokenValue)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
