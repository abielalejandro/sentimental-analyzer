package event

import (
	"github.com/abielalejandro/control/config"
	"github.com/abielalejandro/control/internals/services"
	"github.com/abielalejandro/control/pkg/logger"
)

type GenericBus struct {
	config *config.Config
	svc    services.Service
}

func NewGenericBus(
	config *config.Config,
	svc services.Service,
) EventBus {

	return &GenericBus{config: config, svc: svc}
}

func (gen *GenericBus) Listen() {
	l := logger.New(gen.config.Log.Level)
	l.Info("Listen for messages")
}
