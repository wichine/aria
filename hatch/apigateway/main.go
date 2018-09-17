package main

import (
	"aria/core/config"
	"aria/core/log"
	"aria/core/svcdiscovery"
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

	port := strings.Split(config.Config().Address, ":")[1]
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	service.RegisterAllService(server)

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
	for _, k := range service.GetAllServiceKeys() {
		key := config.Config().ServiceKey + k
		sd.Register(key, config.Config().Address)
	}

	// start server
	logger.Infof("apigateway start service on port: %s", port)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
