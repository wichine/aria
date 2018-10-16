package main

import (
	"aria/core"
	"aria/core/config"
	"aria/core/log"
	"aria/core/svcdiscovery"
	"aria/core/tools"
	"aria/hatch/apigateway/service"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strings"
)

var logger = log.GetLogger("apigateway")

func main() {
	err := config.InitConfig("./config/config.yaml")
	if err != nil {
		panic(err)
	}
	logger.Infof("Config: %s", core.GetStructString(config.Config()))
	err = tools.InitializeZipkin(
		config.Config().Statistic.Tracing.Zipkin.Url,
		config.Config().ServiceKey,
		config.Config().Address,
		config.Config().Statistic.Enable,
	)
	if err != nil {
		logger.Errorf("Initialize zipkin error: %s", err)
	}

	port := strings.Split(config.Config().Address, ":")[1]
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	service.RegisterAllService(server, config.Config().Service)

	// start etcd if needed
	if config.Config().ServiceDiscovery.Enable {
		// start an local etcd instance if EtcdServerOn is set to true
		if config.Config().ServiceDiscovery.EtcdServerOn {
			err := svcdiscovery.StartEtcdServer()
			if err != nil {
				panic(fmt.Errorf("start local etcd instance error: %s", err))
			}
		}
	}

	// register all service to sd
	sd, err := svcdiscovery.GetEtcdServiceDiscoveryInstance(config.Config().EtcdServers)
	if err != nil {
		panic(err)
	}
	err = sd.Register(config.Config().ServiceKey, config.Config().Address)
	if err != nil {
		panic(err)
	}
	logger.Infof("registered self to sd: key= %s, address= %s", config.Config().ServiceKey, config.Config().Address)

	// start server
	logger.Infof("apigateway start service on port: %s", port)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
