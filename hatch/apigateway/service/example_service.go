package service

import (
	"aria/core"
	"aria/core/svcproxy"
	"service_generated_by_aria/service/exampleservice"
)

func init() {
	// FIXME: the key of map should be one key of config.service.serviceProxy map
	servicesFactory["example"] = func(key string) core.Transport {
		// serviceImpl should be replaced by the corresponding one
		return svcproxy.ConvertServiceToProxyServer(key, exampleservice.ServiceImpl())
	}
}
