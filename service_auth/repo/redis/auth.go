package sessionrepo

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	auth "kinopoisk/service_auth/proto"
)

type SessionRepo interface {
	CreateSessionRepo(session *auth.Session) error
	GetSessionRepo(token string) (*auth.Session, error)
	DeleteSessionRepo(token string) (bool, error)
}

type SessionRepoRedis struct {
	redisConn  redis.Conn
	expireTime int
}

func NewSessionRepoRedis(redisConn redis.Conn) *SessionRepoRedis {
	return &SessionRepoRedis{
		redisConn:  redisConn,
		expireTime: 24 * 60 * 60,
	}
}

func (s *SessionRepoRedis) CreateSessionRepo(session *auth.Session) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}
	result, err := redis.String(s.redisConn.Do("SET", session.ID, sessionJSON, "EX", s.expireTime))
	if err != nil || result != "OK" {
		return err
	}
	return nil
}

func (s *SessionRepoRedis) GetSessionRepo(token string) (*auth.Session, error) {
	sess := &auth.Session{}
	sessFromRedis, err := redis.String(s.redisConn.Do("GET", token))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(sessFromRedis), sess)
	if err != nil {
		return nil, err
	}
	return sess, nil

}

func (s *SessionRepoRedis) DeleteSessionRepo(token string) (bool, error) {
	exists, err := redis.Bool(s.redisConn.Do("EXISTS", token))
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	_, err = s.redisConn.Do("DEL", token)
	if err != nil {
		return false, err
	}
	return true, nil
}
