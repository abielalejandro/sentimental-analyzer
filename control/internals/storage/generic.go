package storage

import (
	"context"

	"github.com/abielalejandro/control/config"
)

type GenericStorage struct {
	Config *config.Config
}

func NewGenericStorage(config *config.Config) Storage {
	return &GenericStorage{
		Config: config,
	}
}

func (storage *GenericStorage) Create(ctx context.Context, msg *Message) (bool, error) {
	return true, nil
}

func (storage *GenericStorage) Update(ctx context.Context, id string, result *SentimentalResult) (bool, error) {
	return true, nil
}
