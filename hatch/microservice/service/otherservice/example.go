package otherservice

import (
	"aria/hatch/microservice/core/svcdiscovery"
	pb "aria/hatch/microservice/protocol/production"
	"aria/hatch/microservice/service/production"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitsd "github.com/go-kit/kit/sd"
)

var services = map[string]*serviceRegistrar{
	"AddProduction":    &serviceRegistrar{"/services/production", production.AddProductionImpl().Proxy(), endpoint.Nop},
	"GetAllProduction": &serviceRegistrar{"/services/production", production.GetAllProductionImpl().Proxy(), endpoint.Nop},
}

type serviceRegistrar struct {
	key      string
	factory  kitsd.Factory
	endpoint endpoint.Endpoint
}

// Important: must call at start of main
func InitOtherService() error {
	// get the global service discovery instance
	sd, err := svcdiscovery.GetEtcdServiceDiscoveryInstance()
	if err != nil {
		return fmt.Errorf("Get etcd service discovery instance error: %s", err)
	}

	// subcribe necessary services from sd
	for name, sr := range services {
		sr.endpoint, err = sd.Subscribe(sr.key, sr.factory)
		if err != nil {
			return fmt.Errorf("subcribe service for [%s] error: %s", name, err)
		}
	}

	return nil
}

// ***********************************************************************************
// implement interface of other services below,each function typically contains an endpoint
// call and a protocol convertion.Make sure no spelling mistake happens.
// ***********************************************************************************

func AddProduction(request *pb.AddProductionRequest) (*pb.AddProductionResponse, error) {
	resp, err := services["AddProduction"].endpoint(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	if response, ok := resp.(*pb.AddProductionResponse); ok {
		return response, nil
	} else {
		return nil, fmt.Errorf("convert response to AddProductionResponse error.")
	}
}

func GetAllProduction(request *pb.GetAllProductionRequest) (*pb.GetAllProductionResponse, error) {
	resp, err := services["GetAllProduction"].endpoint(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	if response, ok := resp.(*pb.GetAllProductionResponse); ok {
		return response, nil
	} else {
		return nil, fmt.Errorf("convert response to AddProductionResponse error.")
	}
}
