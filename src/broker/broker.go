package broker

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"gitlab.com/project-leaf/mq-service-go/src/config"
	"gitlab.com/project-leaf/mq-service-go/src/proto/core"
	"gitlab.com/project-leaf/mq-service-go/src/proto/mq"
)

const (
	exchangeTypeTopic = "topic"

	// ExchangeEvents is the exchange where event messages are published.
	ExchangeEvents = "events"

	queueEventsScanPhotoSampled     = "events.scan.photo.sampled"
	queueEventsScanPhotoSampledKey  = "events.scan.photo.sampled"
	queueEventsScanPhotoUploaded    = "events.scan.photo.uploaded"
	queueEventsScanPhotoUploadedKey = "events.scan.photo.uploaded"
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
	config *config.Config
	log    *logrus.Logger

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
func (broker *Broker) EnsureTopology() *core.Error {
	broker.log.Info("Ensuring broker topology.")
	chn, _, chnErr := broker.getChannel()
	if chnErr != nil {
		return broker.handleError(chnErr)
	}

	// Ensure needed exchanges.
	if err := chn.ExchangeDeclare(ExchangeEvents, exchangeTypeTopic, true, false, false, false, nil); err != nil {
		return broker.handleError(err)
	}

	// Ensure PhotoScanUploaded queue.
	if _, err := chn.QueueDeclare(queueEventsScanPhotoUploaded, true, false, false, false, amqp.Table{"x-message-ttl": ttlSLA}); err != nil {
		return broker.handleError(err)
	}
	if err := chn.QueueBind(queueEventsScanPhotoUploaded, queueEventsScanPhotoUploadedKey, ExchangeEvents, false, nil); err != nil {
		return broker.handleError(err)
	}

	// Ensure PhotoScanSampled queue.
	if _, err := chn.QueueDeclare(queueEventsScanPhotoSampled, true, false, false, false, amqp.Table{"x-message-ttl": ttlSLA}); err != nil {
		return broker.handleError(err)
	}
	if err := chn.QueueBind(queueEventsScanPhotoSampled, queueEventsScanPhotoSampledKey, ExchangeEvents, false, nil); err != nil {
		return broker.handleError(err)
	}

	broker.log.Info("Broker topology is ready.")
	return nil
}

// PublishEvent will publish the given `SystemEventMessage` to the `events` exchange.
func (broker *Broker) PublishEvent(message mq.SystemEventMessage, ctx *core.Context) *core.Error {
	chn, _, chnErr := broker.getChannel()
	if chnErr != nil {
		broker.log.Errorf("Error getting channel: %T: %s", chnErr, chnErr.Error())
		return broker.handleError(chnErr)
	}

	// Build the event wrapper.
	event := &mq.SystemEvent{
		Context: ctx,
		Event:   message,
	}

	// Marshal the given protobuf message to bytes.
	data, dataErr := proto.Marshal(event)
	if dataErr != nil {
		broker.log.Errorf("Error marshalling protobuf message: %T: %s", dataErr, dataErr.Error())
		return core.NewError500()
	}

	// Construct the AMQP segment to be sent over the wire.
	msg := amqp.Publishing{
		ContentType:  "application/protobuf",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Type:         message.RoutingKey(),
		AppId:        "mq-service",
		Body:         data,
	}

	// Publish the event.
	if err := chn.Publish(ExchangeEvents, message.RoutingKey(), true, false, msg); err != nil {
		broker.log.Errorf("Error publishing event: %T: %s", err, err.Error())
		return broker.handleError(err)
	}

	return nil
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

// handleError will mutate the receiver so that the service can recover from broker errors.
//
// Public interface methods which use the internal connection &| channel should call
// this method any time an error is returned from a method related to a connection &|
// channel. This will ensure that connections can be re-eastablished and channels re-opened.
//
// This routine will also take the given error and construct an error from it which can be
// more directly used in this service.
func (broker *Broker) handleError(err error) *core.Error {
	// Put the broker back into a pristine state so that it
	// can handle connection &| channel issues.
	if broker.channel != nil {
		broker.channel.Close()
		broker.channel = nil
	}
	if broker.connection != nil {
		broker.connection.Close()
		broker.connection = nil
	}

	return core.New500FromError(err, broker.log)
}
