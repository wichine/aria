package service

import (
	"aria/hatch/apigateway/middleware"
	"fmt"
	"google.golang.org/grpc"
	"service_generated_by_aria/service/exampleservice"
)

func init() {
	registerFuncs["/service/example"] = RegisterService
}

func NewExampleService(serviceKey string) (*exampleservice.ExampleService, error) {
	as := exampleservice.AddProductionImpl()
	gs := exampleservice.GetAllProductionImpl()

	// Important: must call this method or the proxy function will be invalid
	middleware.WrapMiddleware(serviceKey, as, gs)

	return &exampleservice.ExampleService{
		AddProductionService:    as,
		GetAllProductionService: gs,
	}, nil
}

func RegisterService(serviceKey string, server *grpc.Server) error {
	s, err := NewExampleService(serviceKey)
	if err != nil {
		return fmt.Errorf("create example service error: %s", err)
	}
	s.Register(server)
	return nil
}
