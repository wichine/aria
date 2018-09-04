package main

import (
	"aria/hatch/microservice/core"
	"aria/hatch/microservice/service/production"
)

func main() {
	a := core.NewAria(core.AriaConfig{
		GrpcPort: ":9090",
	})
	a.RegisterAll(
		// all service registration here
		production.ServiceImpl(),
	)
	a.ServeGRPC()
}
