package event

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"sync"
)

// Dispatcher interface must be implemented by Queue
type Dispatcher interface {
	DispatchEvent(routingKey string, payload interface{})
}

type queueCommand struct {
	routingKey string
	payload    interface{}
}

// RabbitMQEventDispatcher is an event dispatcher that sends event to the RabbitMQ Exchange
type RabbitMQEventDispatcher struct {
	logger                 *zerolog.Logger
	exchangeName           string
	connection             *amqp.Connection
	channel                *amqp.Channel
	sendChannel            chan *queueCommand
	connectionCloseChannel chan *amqp.Error
	connectionMutex        sync.Mutex
}

// NewRabbitMQEventDispatcher create and returns a new RabbitMQEventDispatcher
func NewRabbitMQEventDispatcher(logger *zerolog.Logger, connectionString string) (*RabbitMQEventDispatcher, error) {
	sendChannel := make(chan *queueCommand, 200)
	connectionCloseChannel := make(chan *amqp.Error)

	ctxLogger := logger.With().Str("module", "RabbitMQEventDispatcher").Logger()

	dispatcher := &RabbitMQEventDispatcher{
		logger:                 &ctxLogger,
		exchangeName:           "xm_exchange",
		sendChannel:            sendChannel,
		connectionCloseChannel: connectionCloseChannel,
	}

	dispatcher.start(connectionString)

	return dispatcher, nil
}

// DispatchEvent dispatches events to the message queue
func (eventDispatcher *RabbitMQEventDispatcher) DispatchEvent(routingKey string, payload interface{}) {
	eventDispatcher.sendChannel <- &queueCommand{routingKey: routingKey, payload: payload}
}

func (eventDispatcher *RabbitMQEventDispatcher) start(connectionString string) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		eventDispatcher.logger.Fatal().Err(err).Msg("cannot connect to rabbitmq")
	}

	eventDispatcher.logger.Info().Msg("rabbitMQ connected")

	ch, err := conn.Channel()
	if err != nil {
		eventDispatcher.logger.Fatal().Err(err).Msg("failed to open rabbitmq connection channel")
	}

	err = ch.ExchangeDeclare(eventDispatcher.exchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		eventDispatcher.logger.Fatal().Err(err).Msg("failed to declare rabbitmq exchange")
	}

	go func() {
		for {
			var command *queueCommand

			// ensure that connection process is not going on
			eventDispatcher.connectionMutex.Lock()
			eventDispatcher.connectionMutex.Unlock()

			select {
			case commandFromSendChannel := <-eventDispatcher.sendChannel:
				command = commandFromSendChannel
			}

			body, err := json.Marshal(command.payload)
			if err != nil {
				eventDispatcher.logger.Error().Msg("Failed to convert payload to JSON" + ": " + err.Error())
				continue
			}

			err = eventDispatcher.channel.Publish(eventDispatcher.exchangeName, command.routingKey, false, false, amqp.Publishing{ContentType: "application/json", Body: body})
			if err != nil {
				eventDispatcher.logger.Error().Msg("Failed to publish to an Exchange" + ": " + err.Error())
			} else {
				eventDispatcher.logger.Trace().Msg("Sent message to queue")
			}
		}
	}()
}
