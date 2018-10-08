package service

import (
	"aria/core"
	"aria/core/config"
	"aria/core/log"
	"google.golang.org/grpc"
)

var servicesFactory = map[string]serviceFactory{}
var logger = log.GetLogger("apigateway")

type serviceFactory func(key string) core.Transport

func RegisterAllService(server *grpc.Server, cfg config.Service) error {
	for name, key := range cfg.ServiceProxy {
		if factory, ok := servicesFactory[name]; ok {
			service := factory(key)
			service.Register(server)
			logger.Infof("register service [%s] to grpc server.", name)
		} else {
			logger.Errorf("service factory of [%s] not found in service factory map!", name)
		}
	}
	return nil
}
