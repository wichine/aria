package service

import (
	"aria/core/svcproxy"
	"service_generated_by_aria/service/exampleservice"
)

var Example exampleNameSpace

// field name should be as same as rpc method
type exampleNameSpace struct {
	AddProduction    *svcproxy.ServiceProxy
	GetAllProduction *svcproxy.ServiceProxy
}

func init() {
	// the serviceName must be a key in config.service.serviceProxy
	svcproxy.RegisterServiceDiscriptionToProxyCenter("example", &Example, exampleservice.ServiceImpl())
}
