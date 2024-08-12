package ratelimiterusecase

import (
	ratelimiterrepo "kinopoisk/app/ratelimiter/repo/redis"
	"log"
	"net"
	"sync"
	"time"
)

type RateLimiterUseCase interface {
	CheckRateLimit(userAddr string) bool
}

type RateLimiterUseCaseStruct struct {
	mu              *sync.RWMutex
	RateLimiterRepo ratelimiterrepo.RateLimiterRepo
}

func NewRateLimiterUseCaseStruct(repo ratelimiterrepo.RateLimiterRepo) *RateLimiterUseCaseStruct {
	return &RateLimiterUseCaseStruct{
		mu:              &sync.RWMutex{},
		RateLimiterRepo: repo,
	}
}

func (rl *RateLimiterUseCaseStruct) CheckRateLimit(userAddr string) bool {
	currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
	rl.mu.RLock()
	host, _, err := net.SplitHostPort(userAddr)
	if err != nil {
		log.Printf("error in spliting addres: %s", err)
		return true
	}
	numOfRequests := rl.RateLimiterRepo.CheckRateLimitRepo(host, currentTimeMillis-2000)
	rl.mu.RUnlock()
	var canMakeRequest = false
	if numOfRequests < 3 {
		canMakeRequest = true
		rl.mu.Lock()
		rl.RateLimiterRepo.AddRateRepo(host, currentTimeMillis)
		rl.mu.Unlock()
	} else {
		log.Printf("too much queries from %s\n", host)
	}
	return canMakeRequest
}
