package otherservice

import (
	"aria/core/svcproxy"
	"aria/hatch/microservice/service/exampleservice"
)

var Example exampleNameSpace

// field name should be as same as rpc method
type exampleNameSpace struct {
	AddProduction    *svcproxy.ServiceProxy
	GetAllProduction *svcproxy.ServiceProxy
}

func init() {
	// the map key must be a key in config.service.serviceProxy
	initFuncs["example"] = func(serviceKey string) {
		// use the global namespace var as parameter
		svcproxy.InitializeAllServiceInOneNameSpace(&Example, serviceKey, exampleservice.ServiceImpl())
	}
}
