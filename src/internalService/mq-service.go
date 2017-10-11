package internalService

import (
	"context"

	"gitlab.com/project-leaf/mq-service-go/src/proto/mq"
)

// InternalMQService is the type which implements our `mq-service.proto::InternalMQService`.
type InternalMQService struct{}

// New will build and return an `InternalMQService` instance.
func New() *InternalMQService {
	return &InternalMQService{}
}

// PubImageScanUploaded will publish an `ImageScanUploaded` event to the central event bus according to the given request.
func (service *InternalMQService) PubImageScanUploaded(ctx context.Context, req *mq.PubImageScanUploadedRequest) (*mq.PubImageScanUploadedResponse, error) {
	println("Called PubImageScanUploaded handler.")
	return nil, nil
}
