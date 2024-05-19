package event

import (
	"github.com/abielalejandro/web/config"
)

const (
	BUS_RABBIT  = "rabbitmq"
	BUS_GENERIC = "generic"
)

type EventBus interface {
	Listen()
}

func NewEventBus(
	config *config.Config,
	chW chan<- SentimentalResult,
	chR <-chan string) EventBus {

	switch config.EventBus.Type {
	case BUS_RABBIT:
		return NewRabbitMqBus(config, chW, chR)
	default:
		return NewGenericBus(config, chW, chR)
	}
}
