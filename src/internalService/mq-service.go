package internalService

import (
	"context"

	"github.com/sirupsen/logrus"

	"gitlab.com/project-leaf/mq-service-go/src/broker"
	"gitlab.com/project-leaf/mq-service-go/src/config"
	"gitlab.com/project-leaf/mq-service-go/src/proto/mq"
)

// InternalMQService is the type which implements our `mq-service.proto::InternalMQService`.
type InternalMQService struct {
	config *config.Config
	log    *logrus.Logger
	broker *broker.Broker
}

// New will build and return an `InternalMQService` instance.
func New(cfg *config.Config, log *logrus.Logger, broker *broker.Broker) *InternalMQService {
	return &InternalMQService{cfg, log, broker}
}

// PubImageScanUploaded will publish an `ImageScanUploaded` event to the central event bus according to the given request.
func (service *InternalMQService) PubImageScanUploaded(ctx context.Context, req *mq.PubImageScanUploadedRequest) (*mq.PubImageScanUploadedResponse, error) {
	response := &mq.PubImageScanUploadedResponse{Error: nil}

	// Build the event object and send it to the broker.
	event := &mq.SystemEvent_ImageScanUploaded{
		ImageScanUploaded: &mq.EventImageScanUploaded{
			Id: req.GetId(),
		},
	}
	if err := service.broker.PublishEvent(event, req.GetContext()); err != nil {
		response.Error = err
		return response, nil
	}

	return response, nil
}
