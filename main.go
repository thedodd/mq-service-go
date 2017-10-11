package main

import (
	"fmt"

	"gitlab.com/project-leaf/mq-service-go/src/api"
)

func main() {
	if err := api.New().Listen(); err != nil {
		fmt.Printf("Error from listener: %T: %s", err, err.Error())
	}
}
