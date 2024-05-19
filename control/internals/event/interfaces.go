package event

import (
	"github.com/abielalejandro/control/config"
	"github.com/abielalejandro/control/internals/services"
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
	svc services.Service) EventBus {

	switch config.EventBus.Type {
	case BUS_RABBIT:
		return NewRabbitMqBus(config, svc)
	default:
		return NewGenericBus(config, svc)
	}
}
