package ratelimiterrepo

import (
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

type RateLimiterRepo interface {
	CheckRateLimitRepo(userAddr string, minTime int64) int
	AddRateRepo(userAddr string, currentTime int64)
}

type RateLimiterRepoRedis struct {
	RedisConn redis.Conn
	logger    *zap.SugaredLogger
}

func NewRateLimiterRepoRedis(redisConn redis.Conn, logger *zap.SugaredLogger) *RateLimiterRepoRedis {
	return &RateLimiterRepoRedis{
		RedisConn: redisConn,
		logger:    logger,
	}
}

func (rl *RateLimiterRepoRedis) CheckRateLimitRepo(userAddr string, minTime int64) int {
	count, err := redis.Int(rl.RedisConn.Do("ZCOUNT", userAddr, minTime, "+inf"))
	if err != nil {
		rl.logger.Errorf("error in count requests in ratelimiter")
		return 0
	}
	_, err = rl.RedisConn.Do("ZREMRANGEBYSCORE", userAddr, "-inf", minTime)
	if err != nil {
		rl.logger.Errorf("error in deleting old  requests in ratelimiter")
	}
	return count
}

func (rl *RateLimiterRepoRedis) AddRateRepo(userAddr string, currentTime int64) {
	_, err := rl.RedisConn.Do("ZADD", userAddr, currentTime, currentTime)
	if err != nil {
		rl.logger.Errorf("error in adding request in ratelimiter: %s", err)
	}
}
