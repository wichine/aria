package main

import (
	"aria/hatch/microservice/core"
	"aria/hatch/microservice/service/otherservice"
	"aria/hatch/microservice/service/production"
	"aria/hatch/microservice/svcdiscovery"
	"fmt"
	"os"
)

func main() {
	// mock the ECTD_SERVERS env, should read from env at start
	os.Setenv("ETCD_SERVERS", "127.0.0.1:2379")

	// init other service instance
	err := otherservice.InitOtherService()
	if err != nil {
		panic(err)
	}

	// create service
	a, err := core.NewAria(core.AriaConfig{
		GrpcPort: ":9090",
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
	err = sd.Register("/services/example", "127.0.0.1:9090")
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
