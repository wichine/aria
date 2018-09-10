package otherservice

import (
	"aria/hatch/microservice/core/svcdiscovery"
	pb "aria/hatch/microservice/protocol/example"
	"aria/hatch/microservice/service/exampleservice"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitsd "github.com/go-kit/kit/sd"
)

// FIXME: 添加自己需要的微服务
var services = map[string]*serviceRegistrar{
	"AddProduction":    &serviceRegistrar{"/services/example", exampleservice.AddProductionImpl().Proxy(), endpoint.Nop},
	"GetAllProduction": &serviceRegistrar{"/services/example", exampleservice.GetAllProductionImpl().Proxy(), endpoint.Nop},
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

// FIXME: 添加需要调用的微服务的接口
// ***********************************************************************************
// implement interface of other services below,each function typically contains an endpoint
// call and a protocol convertion.Make sure no spelling mistake happens.
// ***********************************************************************************

func AddProduction(request *pb.AddProductionRequest) (*pb.AddProductionResponse, error) {
	// FIXME: services的key需要修改
	resp, err := services["AddProduction"].endpoint(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	// FIXME: pb.AddProductionResponse改为正确的类型
	if response, ok := resp.(*pb.AddProductionResponse); ok {
		return response, nil
	} else {
		// FIXME: error中的类型别忘记改
		return nil, fmt.Errorf("convert response to AddProductionResponse error.")
	}
}

func GetAllProduction(request *pb.GetAllProductionRequest) (*pb.GetAllProductionResponse, error) {
	// FIXME: services的key需要修改
	resp, err := services["GetAllProduction"].endpoint(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	// FIXME: pb.GetAllProductionResponse改为正确的类型
	if response, ok := resp.(*pb.GetAllProductionResponse); ok {
		return response, nil
	} else {
		// FIXME: error中的类型别忘记改
		return nil, fmt.Errorf("convert response to AddProductionResponse error.")
	}
}
