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

// PubPhotoScanUploaded will publish an `PhotoScanUploaded` event to the central event bus according to the given request.
func (service *InternalMQService) PubPhotoScanUploaded(ctx context.Context, req *mq.PubPhotoScanUploadedRequest) (*mq.PubPhotoScanUploadedResponse, error) {
	response := &mq.PubPhotoScanUploadedResponse{Error: nil}
	service.log.Debug("Handling request to publish a PhotoScanUploaded event.")

	// Build the event object and send it to the broker.
	event := &mq.SystemEvent_PhotoScanUploaded{
		PhotoScanUploaded: &mq.EventPhotoScanUploaded{
			Id: req.GetId(),
		},
	}
	if err := service.broker.PublishEvent(event, req.GetContext()); err != nil {
		response.Error = err
		return response, nil
	}

	service.log.Debug("PhotoScanUploaded event successfully published.")
	return response, nil
}

// PubPhotoScanSampled will publish an `PhotoScanSampled` event to the central event bus according to the given request.
func (service *InternalMQService) PubPhotoScanSampled(ctx context.Context, req *mq.PubPhotoScanSampledRequest) (*mq.PubPhotoScanSampledResponse, error) {
	response := &mq.PubPhotoScanSampledResponse{Error: nil}
	service.log.Debug("Handling request to publish a PhotoScanSampled event.")

	// Build the event object and send it to the broker.
	event := &mq.SystemEvent_PhotoScanSampled{
		PhotoScanSampled: &mq.EventPhotoScanSampled{
			Id: req.GetId(),
		},
	}
	if err := service.broker.PublishEvent(event, req.GetContext()); err != nil {
		response.Error = err
		return response, nil
	}

	service.log.Debug("PhotoScanSampled event successfully published.")
	return response, nil
}
