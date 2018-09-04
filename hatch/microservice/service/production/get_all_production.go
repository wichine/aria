package production

import (
	"aria/hatch/microservice/core"
	pb "aria/hatch/microservice/protocol/production"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
)

// model
type Production struct {
	Type       string
	Code       string
	Name       string
	ValueDate  int64
	DueDate    int64
	AnnualRate int64
}

type GetAllProductionService struct {
	EndpointType core.EndpointType
	Decode       func(ctx context.Context, request interface{}) (interface{}, error)
	Encode       func(ctx context.Context, response interface{}) (interface{}, error)
	Do           func(ctx context.Context) ([]*Production, error)
}

func GetAllProductionImpl(epType core.EndpointType) *GetAllProductionService {
	return &GetAllProductionService{
		EndpointType: epType,
		// true implement here
		Do: func(ctx context.Context) ([]*Production, error) {
			return []*Production{
				{
					Type: "test",
				},
			}, nil
		},
		Decode: func(_ context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(*pb.GetAllProductionRequest)
			if !ok {
				return nil, fmt.Errorf("Error translate [request] to [pb.GetAllProductionRequest]")
			}
			return req, nil
		},
		Encode: func(_ context.Context, response interface{}) (interface{}, error) {
			resp, ok := response.(*pb.GetAllProductionResponse)
			if !ok {
				return nil, fmt.Errorf("Error translate [response] to [pb.GetAllProductionResponse]")
			}
			return resp, nil
		},
	}
}

func (gas *GetAllProductionService) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ps, err := gas.Do(ctx)
		if err != nil {
			return nil, err
		}
		var result []*pb.Production
		for _, p := range ps {
			prod := &pb.Production{
				Type:       p.Type,
				Code:       p.Code,
				Name:       p.Name,
				ValueDate:  p.ValueDate,
				DueDate:    p.DueDate,
				AnnualRate: p.AnnualRate,
			}
			result = append(result, prod)
		}
		return &pb.GetAllProductionResponse{Status: 0, Production: result}, nil
	}
}

func (gas *GetAllProductionService) Transport() core.AriaTransport {
	ep := gas.Endpoint()
	return core.AriaTransport{
		Grpc: gokitgrpc.NewServer(
			ep,
			gas.Decode,
			gas.Encode,
		),
	}
}

func (gas *GetAllProductionService) Proxy() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}
