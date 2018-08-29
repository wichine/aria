package endpoint

import (
	"aria/hatch/microservice/service"
	"github.com/go-kit/kit/endpoint"
)

// 定义endpoints接口

type Endpoints struct {
	AddProductionEndpoint    endpoint.Endpoint
	GetAllProductionEndpoint endpoint.Endpoint
}

func MakeAllEndpoints(s service.ProductionService) *Endpoints {
	return &Endpoints{
		AddProductionEndpoint:    addEndpointMiddleware(MakeAddProductionEndpoint(s), rateLimiter, circuitBreaker),
		GetAllProductionEndpoint: addEndpointMiddleware(MakeGetAllProductionEndpoint(s), rateLimiter, circuitBreaker),
	}
}
