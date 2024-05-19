package services

import (
	"context"
	"fmt"
	"time"

	"github.com/abielalejandro/control/config"
	"github.com/abielalejandro/control/internals/storage"
	"github.com/abielalejandro/control/pkg/logger"
)

type logMiddleware struct {
	next Service
	log  *logger.Logger
}

func (middleware *logMiddleware) ProcessMsg(ctx context.Context,
	msg string) (id string, err error) {
	defer func(start time.Time) {
		middleware.log.Info(fmt.Sprintf("Executing  ProcessMsg %s takes %v", msg, time.Since(start)))
	}(time.Now())

	return middleware.next.ProcessMsg(ctx, msg)
}

func (middleware *logMiddleware) UpdateSentimentalMsg(
	ctx context.Context,
	id string,
	result *storage.SentimentalResult) (err error) {
	defer func(start time.Time) {
		middleware.log.Info(fmt.Sprintf("Executing  UpdateSentimentalMsg %s takes %v", id, time.Since(start)))
	}(time.Now())

	return middleware.next.UpdateSentimentalMsg(ctx, id, result)

}

func NewLogMiddleware(config *config.Config, next Service) Service {
	return &logMiddleware{
		next: next,
		log:  logger.New(config.Log.Level),
	}
}
