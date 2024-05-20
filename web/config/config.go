package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App            `yaml:"app"`
		HTTP           `yaml:"http"`
		Log            `yaml:"logger"`
		EventBus       `yaml:"event_bus"`
		RabbitEventBus `yaml:"rabbit_event_bus"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	EventBus struct {
		Type string `yaml:"type" env-required:"true"  env:"EVENT_BUS_TYPE" env-default: "rabbitmq"`
	}

	RabbitEventBus struct {
		Url                      string `yaml:"url" env:"RABBITMQ_URL" env-default: "amqp://guest:guest@localhost:5672/"`
		Exchange                 string `yaml:"exchange" env:"RABBITMQ_EXCHANGE" env-default: "sentimental"`
		ExchangeType             string `yaml:"exchange_type" env:"RABBITMQ_EXCHANGE_TYPE" env-default: "topic"`
		Queue                    string `yaml:"queue" env:"RABBITMQ_QUEUE" env-default: ""`
		ProducerMasterRoutingKey string `yaml:"producer_master_routing_key" env:"RABBITMQ_PRODUCER_MASTER_ROUTING" env-default: "ws.text.created"`
		ConsumerMasterRoutingKey string `yaml:"consumer_master_routing_key" env:"RABBITMQ_CONSUMER_MASTER_ROUTING" env-default: "master.text.analyzed"`
		AutoAck                  bool   `yaml:"auto_ack" env:"RABBITMQ_AUTO_ACK" env-default: true`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
