package main

import (
	"database/sql"
	"flag"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	ratingservicerepo "kinopoisk/service_rating/repo/mysql"
	ratingserviceusecase "kinopoisk/service_rating/usecase"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ChangeRatingQueueName = "change_rating"
	maxDBConnections      = 10
	maxPingDBAttempts     = 60
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
	var rabbitConn *amqp.Connection
	var rabbitChan *amqp.Channel
	rabbitAddr := flag.String("addr", "amqp://guest:guest@rabbitmq:5672/", "rabbit addr")
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

	_, err = rabbitChan.QueueDeclare(
		ChangeRatingQueueName, // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		logger.Fatalf("can not init queue: %s", err)
	}
	err = rabbitChan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		logger.Fatalf("can not set QoS: %s", err)
	}
	tasks, err := rabbitChan.Consume(
		ChangeRatingQueueName, // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		logger.Fatalf("can not init queue consumer: %s", err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
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
	ratingChangerDB := ratingservicerepo.NewRatingChangerMySQL(mySQLDb, logger)
	ratingChanger := ratingserviceusecase.NewRatingChangerApp(logger, ratingChangerDB)
	go ratingChanger.ChangeRating(tasks)
	logger.Infof("worker started")
	wg.Wait()
}
