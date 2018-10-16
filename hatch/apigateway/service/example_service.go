package service

import (
	"aria/core/svcproxy"
	"service_generated_by_aria/service/exampleservice"
)

func init() {
	// serviceName should be one key of config.service.serviceProxy map
	svcproxy.RegisterServiceDiscriptionToServicesFactory("example", exampleservice.ServiceImpl())
}
