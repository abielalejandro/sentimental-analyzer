package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App            `yaml:"app"`
		Log            `yaml:"logger"`
		Storage        `yaml:"storage"`
		EventBus       `yaml:"event_bus"`
		RabbitEventBus `yaml:"rabbit_event_bus"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Storage struct {
		Type     string `yaml:"type" env-required:"true"  env:"STORAGE_TYPE" env-default: "generic"`
		Addr     string `yaml:"addr"  env-required:"true" env:"DB_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env-required:"true"  env:"DB_PORT" env-default:"7000"`
		Password string `yaml:"password" env-required:"true" env:"DB_PWD" env-default:"admin"`
		Db       string `yaml:"db"  env-required:"true" env:"DB_NAME" env-default:"shortener"`
		Username string `yaml:"user" env-required:"true"  env:"DB_USER" env-default:"admin"`
		Ttl      int64  `yaml:"ttl" env-required:"true"  env:"DB_TTL" env-default:10`
	}

	EventBus struct {
		Type string `yaml:"type" env-required:"true"  env:"EVENT_BUS_TYPE" env-default: "rabbitmq"`
	}

	RabbitEventBus struct {
		Url                        string `yaml:"url" env:"RABBITMQ_URL" env-default: "amqp://guest:guest@localhost:5672/"`
		Exchange                   string `yaml:"exchange" env:"RABBITMQ_EXCHANGE" env-default: "sentimental"`
		ExchangeType               string `yaml:"exchange_type" env:"RABBITMQ_EXCHANGE_TYPE" env-default: "topic"`
		Queue                      string `yaml:"queue" env:"RABBITMQ_QUEUE" env-default: ""`
		ProducerWsRoutingKey       string `yaml:"producer_ws_routing_key" env:"RABBITMQ_PRODUCER_WS_ROUTING" env-default: "master.text.analyzed"`
		ProducerAnalizerRoutingKey string `yaml:"producer_analizer_routing_key" env:"RABBITMQ_PRODUCER_ANALYZER_ROUTING" env-default: "master.text.created"`
		ConsumerWsRoutingKey       string `yaml:"consumer_ws_routing_key" env:"RABBITMQ_CONSUMER_WS_ROUTING" env-default: "ws.text.created"`
		ConsumerAnalizerRoutingKey string `yaml:"consumer_analizer_routing_key" env:"RABBITMQ_CONSUMER_ANALYZER_ROUTING" env-default: "analyzer.text.analyzed"`
		AutoAck                    bool   `yaml:"auto_ack" env:"RABBITMQ_AUTO_ACK" env-default: true`
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
