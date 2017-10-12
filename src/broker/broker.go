package broker

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"gitlab.com/project-leaf/mq-service-go/src/config"
)

const (
	exchangeTypeTopic = "topic"

	exchangeEvents = "events"

	queueEventsImageScanUploaded    = "events.image_scan_uploaded"
	queueEventsImageScanUploadedKey = "events.image_scan_uploaded"
)

var (
	// Ten minutes is our current SLA for message processing.
	//
	// Probably need to lower this threshold in the future.
	ttlSLA int64 = 1000 * 60 * 10

	ttlTable amqp.Table
)

// Broker exposes an interface for managing connections to the backend message broker.
type Broker struct {
	config     *config.Config
	log        *logrus.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
}

// New will build and return a `Broker` instance.
func New(cfg *config.Config, log *logrus.Logger) *Broker {
	return &Broker{cfg, log, nil, nil}
}

// EnsureTopology will ensure the needed topology is in place in the broker.
//
// This routine should only be called once when the service is first started.
func (broker *Broker) EnsureTopology() error {
	broker.log.Info("Ensuring broker topology.")
	chn, _, chnErr := broker.getChannel()
	if chnErr != nil {
		return broker.HandleError(chnErr)
	}

	// Ensure needed exchanges.
	if err := chn.ExchangeDeclare(exchangeEvents, exchangeTypeTopic, true, false, false, false, nil); err != nil {
		return broker.HandleError(err)
	}

	// Ensure needed queues.
	if _, err := chn.QueueDeclare(queueEventsImageScanUploaded, true, false, false, false, amqp.Table{"x-message-ttl": ttlSLA}); err != nil {
		return broker.HandleError(err)
	}
	if err := chn.QueueBind(queueEventsImageScanUploaded, queueEventsImageScanUploadedKey, exchangeEvents, false, nil); err != nil {
		return broker.HandleError(err)
	}

	return nil
}

// HandleError will mutate the receiver so that the service can recover from broker errors.
//
// This routine will also accept an error and return it at the end. Read the TODO below.
// TODO: update this routine to take an error instance from amqp lib and return an
// error type specific to this service. Probably just a `core.Error`.
func (broker *Broker) HandleError(err error) error {
	if broker.channel != nil {
		broker.channel.Close()
		broker.channel = nil
	}
	if broker.connection != nil {
		broker.connection.Close()
		broker.connection = nil
	}
	return err
}

///////////////////////
// Private Interface //

// getConnection ensure a live connection to the broker exists and will return it.
//
// **This routine will mutate the receiver if a connection is established successfully.**
// In addition to mutating the receiver, it will also return the established connection
// to the caller.
func (broker *Broker) getConnection() (*amqp.Connection, error) {
	// If a live connection already exists, then return in.
	if broker.connection != nil {
		return broker.connection, nil
	}

	// Dial a new connection.
	broker.log.Info("Establishing broker connection.")
	conn, dialErr := amqp.Dial(broker.config.BrokerConnectionString)
	if dialErr != nil {
		return nil, dialErr
	}

	// Mutate receiver by updating its `connection` field, and return.
	broker.connection = conn
	return conn, nil
}

// getChannel will ensure an open channel to the broker exists and will return it.
//
// **This routine will mutate the receiver if a channel is established successfully.**
// In addition to mutating the receiver, it will also return the opened channel and the
// parent connection to the caller.
func (broker *Broker) getChannel() (*amqp.Channel, *amqp.Connection, error) {
	// Ensure we have a working connection.
	conn, connErr := broker.getConnection()
	if connErr != nil {
		return nil, nil, connErr
	}

	// If an open channel already exists, use it.
	if broker.channel != nil {
		return broker.channel, conn, nil
	}

	// Else, open a new channel.
	broker.log.Info("Opening broker channel.")
	chn, chnErr := conn.Channel()
	if chnErr != nil {
		conn.Close()
		broker.connection = nil // Mutates receiver.
		return nil, nil, chnErr
	}

	// Mutate receiver by updating its `channel` field, and return.
	broker.channel = chn
	return chn, conn, nil
}
