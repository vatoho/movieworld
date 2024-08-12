package main

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"kinopoisk/service_auth/interceptor"
	auth "kinopoisk/service_auth/proto"
	userrepo "kinopoisk/service_auth/repo/mysql"
	sessionrepo "kinopoisk/service_auth/repo/redis"
	authserviceusecase "kinopoisk/service_auth/usecase"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	maxDBConnections  = 10
	maxPingDBAttempts = 60
)

func openMySQLConnection() (*sql.DB, error) {
	dsn := "root:"
	mysqlPassword := os.Getenv("pass")
	dsn += mysqlPassword
	dsn += "@tcp(mysql:3306)/golang?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxDBConnections)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	attemptsNumber := 0
	for range ticker.C {
		err = db.Ping()
		attemptsNumber++
		if err == nil {
			break
		}
		if attemptsNumber == maxPingDBAttempts {
			return nil, err
		}
	}
	return db, nil
}

func openRedis() (redis.Conn, error) {
	c, err := redis.DialURL("redis://user:@redis:6379/0")
	if err != nil {
		return nil, err
	}
	return c, nil

}

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Printf("error in logger start")
		return
	}
	logger := zapLogger.Sugar()
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Printf("error in logger sync")
		}
	}()
	envFilePath := "./.env"
	err = godotenv.Load(envFilePath)
	if err != nil {
		logger.Fatalf("Error loading .env file: %s", err)
	}
	mySQLDb, err := openMySQLConnection()
	if err != nil {
		logger.Errorf("error in connection to mysql: %s", err)
		return
	}
	logger.Infof("connected to mysql")
	defer func() {
		err = mySQLDb.Close()
		if err != nil {
			logger.Errorf("error in close connection to mysql: %s", err)
		}
	}()

	redisConn, err := openRedis()
	if err != nil {
		logger.Infof("error on connection to redis: %s", err.Error())
	}
	defer func(redisConn redis.Conn) {
		err = redisConn.Close()
		if err != nil {
			logger.Infof("error on redis close: %s", err.Error())
		}
	}(redisConn)
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Fatalf("can not listen port 8082: %s", err)
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AccessLogInterceptor),
	)
	userRepo := userrepo.NewUserRepoMySQL(mySQLDb)
	sessionRepo := sessionrepo.NewSessionRepoRedis(redisConn)
	auth.RegisterAuthMakerServer(server, authserviceusecase.NewAuthGRPCServer(userRepo, sessionRepo))
	logger.Info("starting server at :8082")
	err = server.Serve(lis)
	if err != nil {
		logger.Fatalf("error in serving server on port 8082 %s", err)
	}
}
