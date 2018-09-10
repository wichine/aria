package main

import (
	aria "aria/hatch/microservice/core"
	"aria/hatch/microservice/core/config"
	"aria/hatch/microservice/core/svcdiscovery"
	"aria/hatch/microservice/service/otherservice"
	"aria/hatch/microservice/service/production"
	"fmt"
	"github.com/op/go-logging"
	"os"
)

var logger *logging.Logger

func init() {
	stdoutBackend := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stdout, "", 0),
		logging.MustStringFormatter(`%{color}[%{time:2006-01-02 15:04:05.000}] %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`),
	)
	logging.SetBackend(stdoutBackend)
	logger = logging.MustGetLogger("main")
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
		production.ServiceImpl(),
	)
	// start server
	if err := a.Run(); err != nil {
		panic(err)
	}
}
