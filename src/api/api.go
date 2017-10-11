package api

import (
	"fmt"
	"net"

	"gitlab.com/project-leaf/mq-service-go/src/internalService"
	"gitlab.com/project-leaf/mq-service-go/src/proto/mq"

	"google.golang.org/grpc"
)

// API is the API definition for this service.
type API struct {
	grpcServer *grpc.Server
}

// New will build and return a new `API` instance.
func New() *API {
	// Create the underlying gRPC server for this API.
	grpcServer := grpc.NewServer()

	// Register services.
	internalMQService := internalService.New()
	mq.RegisterInternalMQServiceServer(grpcServer, internalMQService)

	return &API{
		grpcServer: grpcServer,
	}
}

// Listen will make this API listen on its configured port.
func (api *API) Listen() error {
	// Initialize the API listener.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4004)) // TODO: update from config.
	if err != nil {
		println("Failed to initialize the API listener: %v", err)
	}

	// Listen for requests.
	if err := api.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
