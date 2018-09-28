package otherservice

import (
	"aria/core/svcproxy"
	"aria/hatch/microservice/service/exampleservice"
)

var exampleKey = "/services/example"
var Example exampleNameSpace

type exampleNameSpace struct {
	AddProduction    *svcproxy.ServiceProxy
	GetAllProduction *svcproxy.ServiceProxy
}

func init() {
	// Step1: create service proxy objects
	Example = exampleNameSpace{
		svcproxy.NewServiceProxy(exampleKey, "AddProduction", exampleservice.AddProductionImpl().Proxy()),
		svcproxy.NewServiceProxy(exampleKey, "GetAllProduction", exampleservice.GetAllProductionImpl().Proxy()),
	}
	// Step2: register service proxy object to global map
	svcproxy.RegisterServices(
		Example.AddProduction,
		Example.GetAllProduction,
	)
}
