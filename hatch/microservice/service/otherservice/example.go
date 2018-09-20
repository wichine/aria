package otherservice

import (
	"aria/core/config"
	"aria/core/svcproxy"
	"aria/hatch/microservice/service/exampleservice"
)

// Step1: create service proxy objects
var exampleKey = "/services/example"
var AddProduction = svcproxy.NewServiceProxy(exampleKey, "AddProduction", exampleservice.AddProductionImpl().Proxy())
var GetAllProduction = svcproxy.NewServiceProxy(exampleKey, "GetAllProduction", exampleservice.GetAllProductionImpl().Proxy())

func init() {
	// Step2: register service proxy object to global map
	svcproxy.RegisterServices(AddProduction, GetAllProduction)
}

// Step3: call this method to subcribe from service discovery component
// Important: must call at start of main
func InitOtherService() error {
	return svcproxy.InitServiceProxy(config.Config().EtcdServers)
}
