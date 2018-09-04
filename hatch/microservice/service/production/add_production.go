package production

import (
	"aria/hatch/microservice/core"
	pb "aria/hatch/microservice/protocol/production"
	"aria/hatch/microservice/service"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
)

type AddProductionService struct {
	EndpointType core.EndpointType
	Decode       func(ctx context.Context, request interface{}) (interface{}, error)
	Encode       func(ctx context.Context, response interface{}) (interface{}, error)
	Do           func(ctx context.Context, production *service.Production) (int64, error)
}

func AddProductionImpl(epType core.EndpointType) *AddProductionService {
	return &AddProductionService{
		EndpointType: epType,
		Do: func(ctx context.Context, production *service.Production) (int64, error) {
			return 0, nil
		},
		Encode: func(ctx context.Context, response interface{}) (interface{}, error) {
			resp, ok := response.(*pb.AddProductionResponse)
			if !ok {
				return nil, fmt.Errorf("Error translate [response] to [pb.AddProductionResponse]")
			}
			return resp, nil
		},
		Decode: func(ctx context.Context, request interface{}) (interface{}, error) {
			resp, ok := request.(*pb.AddProductionRequest)
			if !ok {
				return nil, fmt.Errorf("Error translate [response] to [pb.AddProductionResponse]")
			}
			return resp, nil
		},
	}
}

func (aps *AddProductionService) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, error error) {
		req := request.(*pb.AddProductionRequest)
		prod := &service.Production{
			Type:       req.Type,
			Code:       req.Code,
			Name:       req.Name,
			ValueDate:  req.ValueDate,
			DueDate:    req.DueDate,
			AnnualRate: req.AnnualRate,
		}
		count, err := aps.Do(ctx, prod)
		if err != nil {
			return nil, err
		}
		return pb.AddProductionResponse{Status: 0, Msg: fmt.Sprintf("Total count： %d", count)}, nil

	}
}

func (aps *AddProductionService) Transport() core.AriaTransport {
	// default use endpoint
	ep := aps.Endpoint()
	return core.AriaTransport{
		Grpc: gokitgrpc.NewServer(
			ep,
			aps.Decode,
			aps.Encode,
		),
	}
}

func (aps *AddProductionService) Proxy() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}
