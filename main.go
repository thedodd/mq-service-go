package main

import (
	"fmt"

	"gitlab.com/project-leaf/mq-service-go/src/api"
	"gitlab.com/project-leaf/mq-service-go/src/broker"
	"gitlab.com/project-leaf/mq-service-go/src/config"
	"gitlab.com/project-leaf/mq-service-go/src/logging"
)

func main() {
	cfg := config.New()
	log := logging.GetLogger(cfg)
	broker := broker.New(cfg, log)

	// Ensure MQ topology is ready to rock. We will not crash the server here.
	if err := broker.EnsureTopology(); err != nil {
		log.Errorf("Error while ensuring broker topology: %s", err.Error())
	}

	// Boot the API.
	if err := api.New(cfg, log, broker).Listen(); err != nil {
		fmt.Printf("Error from listener: %T: %s", err, err.Error())
	}
}
