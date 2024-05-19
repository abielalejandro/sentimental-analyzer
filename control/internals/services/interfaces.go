package services

import (
	"context"

	"github.com/abielalejandro/control/config"
	"github.com/abielalejandro/control/internals/storage"
)

type Service interface {
	ProcessMsg(
		ctx context.Context,
		msg string) (string, error)
	UpdateSentimentalMsg(
		ctx context.Context,
		id string,
		result *storage.SentimentalResult) error
}

func NewControlService(
	config *config.Config,
	storage storage.Storage) Service {

	return &ControlService{
		config:  config,
		storage: storage,
	}
}
