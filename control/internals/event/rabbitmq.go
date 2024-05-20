package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abielalejandro/control/config"
	"github.com/abielalejandro/control/internals/services"
	"github.com/abielalejandro/control/internals/storage"
	"github.com/abielalejandro/control/pkg/logger"
	"github.com/abielalejandro/control/pkg/utils"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqBus struct {
	config *config.Config
	svc    services.Service
	l      *logger.Logger
}

func NewRabbitMqBus(
	config *config.Config,
	svc services.Service) EventBus {
	return &RabbitMqBus{config: config, svc: svc, l: logger.New(config.Log.Level)}
}

func (gen *RabbitMqBus) Listen() {
	var forever chan struct{}
	conn, err := amqp.Dial(gen.config.RabbitEventBus.Url)
	utils.FailOnError(err, "Error connecting to rabbit")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	gen.listenWs(ch)
	gen.listenAnalyzer(ch)
	<-forever
}

func (gen *RabbitMqBus) listenWs(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		gen.config.RabbitEventBus.Exchange,
		gen.config.RabbitEventBus.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Error declaring the exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Error declaring the queue for ws")

	err = ch.QueueBind(
		q.Name,
		gen.config.RabbitEventBus.ConsumerWsRoutingKey,
		gen.config.RabbitEventBus.Exchange,
		false,
		nil,
	)
	utils.FailOnError(err, "Error  linking queue/exchange ws")

	msgs, err := ch.Consume(
		q.Name,
		"",
		gen.config.RabbitEventBus.AutoAck,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for d := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			eventReceived := cloudevents.NewEvent()
			err := json.Unmarshal(d.Body, &eventReceived)
			if err != nil {
				gen.l.Error(err.Error())
				continue
			}

			var t Message
			err = json.Unmarshal(eventReceived.Data(), &t)
			if err != nil {
				gen.l.Error(err.Error())
				continue
			}

			gen.l.Info(fmt.Sprintf("mensaje recibido from ws %v", t))

			id, err := gen.svc.ProcessMsg(ctx, t.Msg)
			if err != nil {
				gen.l.Error(err.Error())
				continue
			}

			toAnalyze := &MessageToAnalyze{
				Msg: t.Msg,
				Id:  id,
			}

			event := cloudevents.NewEvent()
			event.SetID(id)
			event.SetDataContentType("application/json")
			event.SetSource("sentimental/control")
			event.SetType(gen.config.RabbitEventBus.ProducerAnalizerRoutingKey)
			event.SetData(cloudevents.ApplicationJSON, toAnalyze)

			bytes, _ := json.Marshal(event)

			err = ch.PublishWithContext(ctx,
				gen.config.RabbitEventBus.Exchange,
				gen.config.RabbitEventBus.ProducerAnalizerRoutingKey,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        bytes,
				},
			)

			if err != nil {
				gen.l.Error(err.Error())
			}
			gen.l.Info(fmt.Sprintf("mensaje enviado a analizar %v", string(bytes)))
			cancel()
		}
	}()
}

func (gen *RabbitMqBus) listenAnalyzer(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		gen.config.RabbitEventBus.Exchange,
		gen.config.RabbitEventBus.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Error declaring the exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	utils.FailOnError(err, "Error declaring the queue for the analyzer")

	err = ch.QueueBind(
		q.Name,
		gen.config.RabbitEventBus.ConsumerAnalizerRoutingKey,
		gen.config.RabbitEventBus.Exchange,
		false,
		nil,
	)
	utils.FailOnError(err, "Error  linking queue/exchangei analyzer")

	msgs, err := ch.Consume(
		q.Name,
		"",
		gen.config.RabbitEventBus.AutoAck,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for d := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			eventReceived := cloudevents.NewEvent()
			err := json.Unmarshal(d.Body, &eventReceived)
			if err != nil {
				gen.l.Error(err.Error())
				continue
			}

			gen.l.Info(fmt.Sprintf("receiving crudo [%v]", string(eventReceived.Data())))
			var t storage.SentimentalResult
			err = json.Unmarshal(eventReceived.Data(), &t)
			if err != nil {
				gen.l.Error(err.Error())
				continue
			}

			gen.l.Info(fmt.Sprintf("receiving in an [%v]", t))
			err = gen.svc.UpdateSentimentalMsg(ctx, eventReceived.ID(), &t)

			if err != nil {
				gen.l.Error(err.Error())
				continue
			}

			event := cloudevents.NewEvent()
			event.SetDataContentType("application/json")
			event.SetSource("sentimental/control")
			event.SetType(gen.config.RabbitEventBus.ProducerWsRoutingKey)
			event.SetData(cloudevents.ApplicationJSON, t)

			bytes, _ := json.Marshal(event)

			err = ch.PublishWithContext(ctx,
				gen.config.RabbitEventBus.Exchange,
				gen.config.RabbitEventBus.ProducerWsRoutingKey,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        bytes,
				},
			)

			if err != nil {
				gen.l.Error(err.Error())
			}

			gen.l.Info(fmt.Sprintf("mensaje enviado %v", event.String()))
			cancel()
		}
	}()
}
