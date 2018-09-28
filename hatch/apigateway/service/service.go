package service

import (
	"aria/core/config"
	"aria/core/log"
	"google.golang.org/grpc"
)

var servicesFactory = map[string]func(string) Service{}
var logger = log.GetLogger("apigateway")

type Service interface {
	newService() (interface{}, error)
	RegisterService(server *grpc.Server) error
	GetKey() string
}

func RegisterAllService(server *grpc.Server, cfg config.Service) error {
	for name, key := range cfg.ServiceProxy {
		if factory, ok := servicesFactory[name]; ok {
			service := factory(key)
			err := service.RegisterService(server)
			if err != nil {
				logger.Errorf("register service [%s] to grpc error: %s", name, err)
				return err
			}
			logger.Infof("register service [%s] to grpc server.", service.GetKey())
		} else {
			logger.Errorf("service factory of [%s] not found in service factory map!", name)
		}
	}
	return nil
}
