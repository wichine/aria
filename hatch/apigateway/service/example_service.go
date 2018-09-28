package service

import (
	"aria/core/svcproxy"
	"fmt"
	"google.golang.org/grpc"
	"service_generated_by_aria/service/exampleservice"
)

func init() {
	// FIXME: the name should be one key of config.service.serviceProxy map
	servicesFactory["example"] = NewExampleService
}

type ExampleServie struct {
	// serviceName should be unique
	serviceKey string
}

func NewExampleService(key string) Service {
	return &ExampleServie{key}
}

func (es *ExampleServie) newService() (interface{}, error) {
	// FIXME: create a service
	s := exampleservice.ServiceImpl()

	// Important: must call this method or the proxy function will be invalid
	err := svcproxy.ConvertServiceToProxy(es.serviceKey, s.AddProductionService, s.GetAllProductionService)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (es *ExampleServie) RegisterService(server *grpc.Server) error {
	s, err := es.newService()
	if err != nil {
		return fmt.Errorf("RegisterService create new service error: %s", err)
	}
	// FIXME: convert interface to concrete type
	if rs, ok := s.(*exampleservice.ExampleService); ok {
		rs.Register(server)
		return nil
	}
	return fmt.Errorf("RegisterService convert type error")
}

func (es *ExampleServie) GetKey() string {
	return es.serviceKey
}
