package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	actorrepo "kinopoisk/app/actors/repo/mysql"
	actorusecase "kinopoisk/app/actors/usecase"
	"kinopoisk/app/delivery/handlers"
	filmrepo "kinopoisk/app/films/repo/mysql"
	filmusecase "kinopoisk/app/films/usecase"
	"kinopoisk/app/middleware"
	ratelimiterrepo "kinopoisk/app/ratelimiter/repo/redis"
	ratelimiterusecase "kinopoisk/app/ratelimiter/usecase"
	reviewusecase "kinopoisk/app/reviews/usecase"
	searchrepo "kinopoisk/app/search/repo/mysql"
	searchusecase "kinopoisk/app/search/usecase"
	userusecase "kinopoisk/app/users/usecase"
	auth "kinopoisk/service_auth/proto"
	review "kinopoisk/service_review/proto"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

	grpcConnReview, err := grpc.Dial(
		"service_review:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf("cant connect to grpc review service")
	}
	defer func(grcpConnReview *grpc.ClientConn) {
		err = grcpConnReview.Close()
		if err != nil {
			logger.Errorf("can not stop grpc review server")
		}
	}(grpcConnReview)

	grpcConnAuth, err := grpc.Dial(
		"service_auth:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf("cant connect to grpc auth service")
	}
	defer func(grcpConnAuth *grpc.ClientConn) {
		err = grcpConnAuth.Close()
		if err != nil {
			logger.Errorf("can not stop grpc user server")
		}
	}(grpcConnAuth)

	filmRepo := filmrepo.NewFilmRepoMySQL(mySQLDb, logger)
	filmUseCase := filmusecase.NewFilmUseCaseStruct(filmRepo)

	authGRPCClient := auth.NewAuthMakerClient(grpcConnAuth)
	authUseCase := userusecase.NewAuthGRPCClient(authGRPCClient)

	reviewGRPCClient := review.NewReviewMakerClient(grpcConnReview)
	reviewUseCase := reviewusecase.NewReviewGRPCClient(reviewGRPCClient, filmRepo)

	actorRepo := actorrepo.NewActorRepoMySQL(mySQLDb, logger)
	actorUseCase := actorusecase.NewActorUseCaseStruct(actorRepo)

	rateLimiterRepo := ratelimiterrepo.NewRateLimiterRepoRedis(redisConn, logger)
	rateLimiterUseCase := ratelimiterusecase.NewRateLimiterUseCaseStruct(rateLimiterRepo)

	searchRepo := searchrepo.NewSearchRepoMySQL(mySQLDb, logger)
	searchUseCase := searchusecase.NewSearchUseCaseStruct(searchRepo)

	authHandler := handlers.NewUserHandler(authUseCase)
	reviewHandler := handlers.NewReviewHandler(reviewUseCase)
	filmHandler := handlers.NewFilmHandler(filmUseCase)
	actorHandler := handlers.NewActorHandler(actorUseCase)
	searchHandler := handlers.NewSearchHandler(searchUseCase)

	router := mux.NewRouter()
	router.HandleFunc("/actors", actorHandler.GetActors).Methods(http.MethodGet)
	router.HandleFunc("/actor/{ACTOR_ID}", actorHandler.GetActorByID).Methods(http.MethodGet)

	router.HandleFunc("/films", filmHandler.GetFilms).Methods(http.MethodGet)
	router.HandleFunc("/films/by/{ACTOR_ID}", filmHandler.GetFilmsByActor).Methods(http.MethodGet)

	router.HandleFunc("/film/{FILM_ID}", filmHandler.GetFilmByID).Methods(http.MethodGet)
	router.HandleFunc("/films/soon", filmHandler.GetFilmsSoon).Methods(http.MethodGet)

	router.HandleFunc("/film/{FILM_ID}/actors", filmHandler.GetFilmActors).Methods(http.MethodGet)
	router.HandleFunc("/film/{FILM_ID}/genres", filmHandler.GetFilmGenres).Methods(http.MethodGet)

	router.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)

	router.HandleFunc("/review/{FILM_ID}", reviewHandler.GetReviewsForFilm).Methods(http.MethodGet)

	router.HandleFunc("/search/{DATA}", searchHandler.MakeSearch).Methods(http.MethodGet)

	checkAuthRouter := mux.NewRouter()
	router.Handle("/films/favourite", middleware.AuthMiddleware(authUseCase, checkAuthRouter)).Methods(http.MethodGet)
	router.Handle("/films/favourite/{FILM_ID}", middleware.AuthMiddleware(authUseCase, checkAuthRouter)).Methods(http.MethodPost)
	router.Handle("/films/favourite/{FILM_ID}", middleware.AuthMiddleware(authUseCase, checkAuthRouter)).Methods(http.MethodDelete)

	router.Handle("/review/{FILM_ID}", middleware.AuthMiddleware(authUseCase, checkAuthRouter)).Methods(http.MethodPost)
	router.Handle("/review/{REVIEW_ID}", middleware.AuthMiddleware(authUseCase, checkAuthRouter)).Methods(http.MethodDelete)
	router.Handle("/review/{REVIEW_ID}", middleware.AuthMiddleware(authUseCase, checkAuthRouter)).Methods(http.MethodPut)

	checkAuthRouter.HandleFunc("/films/favourite", filmHandler.GetFavouriteFilms).Methods(http.MethodGet)
	checkAuthRouter.HandleFunc("/films/favourite/{FILM_ID}", filmHandler.AddFavouriteFilm).Methods(http.MethodPost)
	checkAuthRouter.HandleFunc("/films/favourite/{FILM_ID}", filmHandler.DeleteFavouriteFilm).Methods(http.MethodDelete)

	checkAuthRouter.HandleFunc("/review/{FILM_ID}", reviewHandler.AddReview).Methods(http.MethodPost)
	checkAuthRouter.HandleFunc("/review/{REVIEW_ID}", reviewHandler.DeleteReview).Methods(http.MethodDelete)
	checkAuthRouter.HandleFunc("/review/{REVIEW_ID}", reviewHandler.UpdateReview).Methods(http.MethodPut)

	accessLogRouter := middleware.AccessLog(router)
	errorLogRouter := middleware.ErrorLog(accessLogRouter)
	rateLimiterRouter := middleware.RateLimiterMiddleware(rateLimiterUseCase, errorLogRouter)
	mux := middleware.RequestInitMiddleware(rateLimiterRouter)

	addr := ":8080"

	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		logger.Fatalf("errror in server start")
	}
}
