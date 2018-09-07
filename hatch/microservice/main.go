package main

import (
	"aria/hatch/microservice/core"
	"aria/hatch/microservice/core/config"
	"aria/hatch/microservice/core/svcdiscovery"
	"aria/hatch/microservice/service/otherservice"
	"aria/hatch/microservice/service/production"
	"fmt"
	"strings"
)

func main() {
	// init config
	err := config.InitConfig("./config/config.yaml")
	if err != nil {
		panic(err)
	}

	// init other service instance
	err = otherservice.InitOtherService()
	if err != nil {
		panic(err)
	}

	// create service
	grpcAddress := config.Config().Address
	port := strings.Split(grpcAddress, ":")[1]
	a, err := core.NewAria(core.AriaConfig{
		GrpcPort: fmt.Sprintf(":%s", port),
	})
	if err != nil {
		panic(err)
	}
	productionService := production.ServiceImpl()
	a.RegisterAll(
		// all service registration here
		productionService,
	)

	// register to etcd service discovery center
	sd, err := svcdiscovery.GetEtcdServiceDiscoveryInstance()
	if err != nil {
		panic(err)
	}
	err = sd.Register(config.Config().ServiceKey, grpcAddress)
	if err != nil {
		panic(err)
	}
	defer sd.DeRegister()

	// start server
	fmt.Println("Service started!")
	err = a.ServeGRPC()
	if err != nil {
		panic(err)
	}
}
