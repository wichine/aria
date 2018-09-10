package main

import (
	aria "aria/hatch/microservice/core"
	"aria/hatch/microservice/core/config"
	"aria/hatch/microservice/core/log"
	"aria/hatch/microservice/core/svcdiscovery"
	"aria/hatch/microservice/service/exampleservice"
	"aria/hatch/microservice/service/otherservice"
	"fmt"
	"github.com/op/go-logging"
)

var logger *logging.Logger

func init() {
	logger = log.GetLogger("main")
}

func main() {
	// init config
	if err := config.InitConfig("./config/config.yaml"); err != nil {
		panic(err)
	}

	if config.Config().ServiceDiscovery.Enable {
		// start an local etcd instance if EtcdServerOn is set to true
		if config.Config().ServiceDiscovery.EtcdServerOn {
			err := svcdiscovery.StartEtcdServer()
			if err != nil {
				panic(fmt.Errorf("start local etcd instance error: %s", err))
			}
		}

		// init other service instance
		if err := otherservice.InitOtherService(); err != nil {
			panic(err)
		}
	}

	// create service
	a := aria.New(config.Config())
	a.RegisterAll(
		// all service registration here
		exampleservice.ServiceImpl(),
	)
	logger.Info("=====================")
	// start server
	if err := a.Run(); err != nil {
		panic(err)
	}
}
