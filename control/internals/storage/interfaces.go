package storage

import (
	"context"

	"github.com/abielalejandro/control/config"
)

type Storage interface {
	Create(ctx context.Context, msg *Message) (bool, error)
	Update(ctx context.Context, id string, result *SentimentalResult) (bool, error)
}

func NewStorage(config *config.Config) Storage {
	switch config.Storage.Type {
	case "generic":
		return NewGenericStorage(config)
	case "cassandra":
		return NewCassandraStorage(config)
	default:
		return NewGenericStorage(config)
	}
}
