package api

import (
	"fmt"
	"net"

	"gitlab.com/project-leaf/mq-service-go/src/broker"
	"gitlab.com/project-leaf/mq-service-go/src/config"
	"gitlab.com/project-leaf/mq-service-go/src/internalService"
	"gitlab.com/project-leaf/mq-service-go/src/proto/mq"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// API is the API definition for this service.
type API struct {
	// Config is the API's runtime configuration.
	Config *config.Config
	Log    *logrus.Logger

	grpcServer *grpc.Server
}

// New will build and return a new `API` instance.
func New(cfg *config.Config, log *logrus.Logger, broker *broker.Broker) *API {
	// Create the underlying gRPC server for this API.
	grpcServer := grpc.NewServer()

	// Register services.
	internalMQService := internalService.New(cfg, log, broker)
	mq.RegisterInternalMQServiceServer(grpcServer, internalMQService)

	return &API{cfg, log, grpcServer}
}

// Listen will make this API listen on its configured port.
//
// This routine will panic if it can not successfully listen on the specified interface and port.
func (api *API) Listen() error {
	// Initialize the API listener.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", api.Config.Port))
	if err != nil {
		api.Log.Panicf("Failed to initialize the API listener: %v", err) // NOTE: routine may diverge here.
	}

	// Listen for requests.
	api.Log.Infof("MQ service is listening on '0.0.0.0:%d'.", api.Config.Port)
	if err := api.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
