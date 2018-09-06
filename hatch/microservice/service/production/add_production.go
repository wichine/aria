package production

import (
	"aria/hatch/microservice/core"
	pb "aria/hatch/microservice/protocol/production"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"google.golang.org/grpc"
	"io"
)

type AddProductionService struct {
	*core.Service
	Do func(ctx context.Context, production *Production) (int64, error)
}

func AddProductionImpl() *AddProductionService {
	a := &AddProductionService{
		Service: core.NewDefaultService(),
		// true implement
		Do: func(ctx context.Context, production *Production) (int64, error) {
			return 0, nil
		},
	}
	a.Endpoint = func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// coercion first
		req := request.(*pb.AddProductionRequest)
		// then convert to 'Do' param
		prod := &Production{
			Type:       req.Type,
			Code:       req.Code,
			Name:       req.Name,
			ValueDate:  req.ValueDate,
			DueDate:    req.DueDate,
			AnnualRate: req.AnnualRate,
		}
		rawResp, err := a.Do(ctx, prod)
		if err != nil {
			return nil, err
		}
		// encode to true res
		return &pb.AddProductionResponse{Status: 0, Msg: fmt.Sprintf("Total count : %d", rawResp)}, nil
	}
	return a
}

// Serve grpc handler
func (aps *AddProductionService) AddProduction(ctx context.Context, request *pb.AddProductionRequest) (*pb.AddProductionResponse, error) {
	r, err := aps.Compose()(ctx, request)
	if err != nil {
		return nil, err
	}
	// coercion response
	res := r.(*pb.AddProductionResponse)
	return res, nil
}

func (aps *AddProductionService) Proxy() sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		address := instance
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, nil, fmt.Errorf("dial grpc error: %s", err)
		}
		client := pb.NewProductionServiceClient(conn)
		ep := func(ctx context.Context, request interface{}) (interface{}, error) {
			req, ok := request.(*pb.AddProductionRequest)
			if !ok {
				return nil, fmt.Errorf("can not convert request to pb.AddProductionRequest")
			}
			return client.AddProduction(ctx, req)
		}
		return ep, conn, nil
	}
}
