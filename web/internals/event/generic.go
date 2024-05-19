package event

import (
	"github.com/abielalejandro/web/config"
	"github.com/abielalejandro/web/pkg/logger"
)

type GenericBus struct {
	config *config.Config
	chW    chan<- SentimentalResult
	chR    <-chan string
}

func NewGenericBus(
	config *config.Config,
	chW chan<- SentimentalResult,
	chR <-chan string,
) EventBus {

	return &GenericBus{config: config, chW: chW, chR: chR}
}

func (gen *GenericBus) Listen() {
	l := logger.New(gen.config.Log.Level)
	l.Info("Listen for messages")
}
