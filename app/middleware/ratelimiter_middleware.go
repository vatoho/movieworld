package middleware

import (
	ratelimiterusecase "kinopoisk/app/ratelimiter/usecase"
	"log"
	"net/http"
)

func RateLimiterMiddleware(rateLimiterUseCases ratelimiterusecase.RateLimiterUseCase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, err := GetLoggerFromContext(r.Context())
		if err != nil {
			log.Printf("can not get logger from context: %s", err)
			WriteNoLoggerResponse(w)
		}
		logger.Infof("ratelimiter middleware check")
		requestAddr := r.RemoteAddr
		canMakeRequest := rateLimiterUseCases.CheckRateLimit(requestAddr)
		if !canMakeRequest {
			return
		}
		next.ServeHTTP(w, r)
	})
}
