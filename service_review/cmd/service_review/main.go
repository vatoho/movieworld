package main

import (
	"database/sql"
	"flag"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"kinopoisk/service_review/interceptor"
	review "kinopoisk/service_review/proto"
	reviewservicerepo "kinopoisk/service_review/repo/mysql"
	reviewserviceusecse "kinopoisk/service_review/usecase"
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
	rabbitAddr := flag.String("addr", "amqp://guest:guest@127.0.0.1:5672/", "rabbit addr")
	var rabbitConn *amqp.Connection
	var rabbitChan *amqp.Channel
	flag.Parse()
	rabbitConn, err = amqp.Dial(*rabbitAddr)
	if err != nil {
		logger.Fatalf("can not connect to rabbit: %s", err)
	}
	rabbitChan, err = rabbitConn.Channel()
	if err != nil {
		logger.Fatalf("can not make rabbit message channel: %s", err)
	}
	defer func(rabbitChan *amqp.Channel) {
		err = rabbitChan.Close()
		if err != nil {
			logger.Errorf("can not close rabbit channel: %s", err)
		}
	}(rabbitChan)

	q, err := rabbitChan.QueueDeclare(
		reviewserviceusecse.ChangeRatingQueueName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logger.Fatalf("can not init queue: %s", err)
	}
	logger.Infof("queue %s have %d msg and %d consumers\n",
		q.Name, q.Messages, q.Consumers)

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

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatalf("can not listen port 8081: %s", err)
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AccessLogInterceptor),
	)
	reviewRepo := reviewservicerepo.NewReviewRepoMySQL(mySQLDb, logger)
	review.RegisterReviewMakerServer(server, reviewserviceusecse.NewReviewGRPCServer(reviewRepo, rabbitChan))
	logger.Info("starting server at :8081")
	err = server.Serve(lis)
	if err != nil {
		logger.Fatalf("error in serving server on port 8081 %s", err)
	}
}
