package app

import (
	"fmt"

	"github.com/abielalejandro/control/config"
	"github.com/abielalejandro/control/internals/event"
	"github.com/abielalejandro/control/internals/services"
	"github.com/abielalejandro/control/internals/storage"
	"github.com/abielalejandro/control/pkg/logger"
)

type App struct {
	event.EventBus
	*config.Config
}

func NewApp(config *config.Config) *App {
	db := storage.NewStorage(config)
	svc := services.NewControlService(config, db)
	return &App{
		Config:   config,
		EventBus: event.NewEventBus(config, services.NewLogMiddleware(config, svc)),
	}
}

func (app *App) Run() {
	l := logger.New(app.Config.Log.Level)
	l.Info("App Running ")
	l.Info(fmt.Sprintf("Config %v", app.Config))
	app.EventBus.Listen()
}
