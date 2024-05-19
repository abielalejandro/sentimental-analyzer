package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abielalejandro/web/config"
	"github.com/abielalejandro/web/pkg/logger"
	"github.com/abielalejandro/web/pkg/utils"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqBus struct {
	config *config.Config
	l      *logger.Logger
	chW    chan<- string
	chR    <-chan string
}

func NewRabbitMqBus(
	config *config.Config,
	chW chan<- string,
	chR <-chan string) EventBus {
	return &RabbitMqBus{config: config, l: logger.New(config.Log.Level), chW: chW, chR: chR}
}

func (gen *RabbitMqBus) Listen() {

	go func() {
		var forever chan struct{}
		conn, err := amqp.Dial(gen.config.RabbitEventBus.Url)
		utils.FailOnError(err, "Error connecting to rabbit")
		defer conn.Close()

		ch, err := conn.Channel()
		utils.FailOnError(err, "Failed to open a channel")
		defer ch.Close()

		gen.listenMaster(ch)
		gen.readLoopMsgs(ch)
		<-forever
	}()

}

func (gen *RabbitMqBus) readLoopMsgs(ch *amqp.Channel) {

	go func() {
		for elem := range gen.chR {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			gen.l.Info(fmt.Sprintf("receiving msg from ws %s", elem))

			t := &Message{
				Msg: elem,
			}

			event := cloudevents.NewEvent()
			event.SetDataContentType("application/json")
			event.SetSource("sentimental/ws")
			event.SetType(gen.config.RabbitEventBus.ProducerMasterRoutingKey)
			event.SetData(cloudevents.ApplicationJSON, t)

			bytes, _ := json.Marshal(event)

			err := ch.PublishWithContext(ctx,
				gen.config.RabbitEventBus.Exchange,
				gen.config.RabbitEventBus.ProducerMasterRoutingKey,
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
			gen.l.Info(fmt.Sprintf("mensaje enviado al master %v", string(bytes)))
			cancel()
		}
	}()

}

func (gen *RabbitMqBus) listenMaster(ch *amqp.Channel) {
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
	utils.FailOnError(err, "Error declaring the queue for master")

	err = ch.QueueBind(
		q.Name,
		gen.config.RabbitEventBus.ConsumerMasterRoutingKey,
		gen.config.RabbitEventBus.Exchange,
		false,
		nil,
	)
	utils.FailOnError(err, "Error  linking queue/exchange master")

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
			t := string(d.Body[:])
			gen.l.Info(fmt.Sprintf("receiving from master [%v]", t))
			if err != nil {
				gen.l.Error(err.Error())
				continue
			}
		}
	}()
}
