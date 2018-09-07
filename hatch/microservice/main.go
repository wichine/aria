package main

import (
	aria "aria/hatch/microservice/core"
	"aria/hatch/microservice/core/config"
	"aria/hatch/microservice/service/otherservice"
	"aria/hatch/microservice/service/production"
	"fmt"
)

func main() {
	// init config
	if err := config.InitConfig("./config/config.yaml"); err != nil {
		panic(err)
	}

	// init other service instance
	if err := otherservice.InitOtherService(); err != nil {
		panic(err)
	}

	// create service
	a := aria.New(config.Config())
	a.RegisterAll(
		// all service registration here
		production.ServiceImpl(),
	)
	// start server
	fmt.Println("Service started!")
	if err := a.Run(); err != nil {
		panic(err)
	}
}
