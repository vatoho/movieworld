package interceptor

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	errorapp "kinopoisk/app/errors"
	"log"
	"time"
)

type loggerKey int

const MyLoggerKey loggerKey = 3

func AccessLogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("access log interceptor start")
	logger, err := initLogger()
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		return nil, err
	}
	requestID := uuid.New().String()
	logger = logger.With(zap.String("request-id", requestID))
	ctx = context.WithValue(ctx, MyLoggerKey, logger)
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	logger.Infow("Request result",
		"after incoming call ", info.FullMethod,
		"request ", req,
		"reply ", reply,
		"time of call ", time.Since(start),
		"metadata ", md,
		"error ", err,
	)
	return reply, err
}

func initLogger() (*zap.SugaredLogger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error in logger initialization: %s", err)
		return nil, err
	}
	myLogger := zapLogger.Sugar()
	return myLogger, nil
}

func GetLoggerFromContext(ctx context.Context) (*zap.SugaredLogger, error) {
	myLogger, ok := ctx.Value(MyLoggerKey).(*zap.SugaredLogger)
	if !ok {
		return myLogger, errorapp.ErrorNoLogger
	}
	return myLogger, nil
}
