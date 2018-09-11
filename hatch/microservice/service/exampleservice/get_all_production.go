package exampleservice

import (
	"aria/core"
	pb "aria/hatch/microservice/protocol/example"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"google.golang.org/grpc"
	"io"
)

type GetAllProductionService struct {
	*core.Service
	Do func(ctx context.Context) ([]*Production, error)
}

func GetAllProductionImpl() *GetAllProductionService {
	s := &GetAllProductionService{
		Service: core.NewDefaultService(),
		// true implement here
		Do: func(ctx context.Context) ([]*Production, error) {
			return []*Production{
				{
					Type: "test",
				},
			}, nil
		},
	}
	s.Endpoint = func(ctx context.Context, request interface{}) (response interface{}, err error) {
		rawResp, err := s.Do(ctx)
		if err != nil {
			return nil, err
		}
		var productions []*pb.Production
		for _, p := range rawResp {
			prod := &pb.Production{
				Type:       p.Type,
				Code:       p.Code,
				Name:       p.Name,
				ValueDate:  p.ValueDate,
				DueDate:    p.DueDate,
				AnnualRate: p.AnnualRate,
			}
			productions = append(productions, prod)
		}
		// encode to true res
		return &pb.GetAllProductionResponse{Status: 0, Production: productions}, nil
	}
	return s
}

// Serve grpc handler
func (gas *GetAllProductionService) GetAllProduction(ctx context.Context, request *pb.GetAllProductionRequest) (response *pb.GetAllProductionResponse, err error) {
	r, err := gas.Compose()(ctx, request)
	if err != nil {
		return nil, err
	}
	// coercion response
	res := r.(*pb.GetAllProductionResponse)
	return res, nil
}

func (gas *GetAllProductionService) Proxy() sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		address := instance
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, nil, fmt.Errorf("dial grpc error: %s", err)
		}
		client := pb.NewProductionServiceClient(conn)
		ep := func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(*pb.GetAllProductionRequest)
			if !ok {
				return nil, fmt.Errorf("can not convert request to pb.AddProductionRequest")
			}
			return client.GetAllProduction(ctx, req)
		}
		return ep, conn, nil
	}
}
