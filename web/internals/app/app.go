package app

import (
	"fmt"

	"github.com/abielalejandro/web/api"
	"github.com/abielalejandro/web/config"
	"github.com/abielalejandro/web/internals/event"
	"github.com/abielalejandro/web/pkg/logger"
)

type App struct {
	event.EventBus
	api.Api
	*config.Config
}

func NewApp(config *config.Config) *App {
	chOut := make(chan string)
	chIn := make(chan string)
	return &App{
		Config:   config,
		EventBus: event.NewEventBus(config, chIn, chOut),
		Api:      api.NewApi(config, chOut, chIn),
	}
}

func (app *App) Run() {
	l := logger.New(app.Config.Log.Level)
	l.Info("App Running ")
	l.Info(fmt.Sprintf("Config %v", app.Config))
	app.EventBus.Listen()
	app.Api.Run()
}
