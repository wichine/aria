package production

import (
	"aria/hatch/microservice/core"
	pb "aria/hatch/microservice/protocol/production"
	"aria/hatch/microservice/service"
	"context"
	"fmt"
)

type AddProductionService struct {
	*core.Service
	Do func(ctx context.Context, production *service.Production) (int64, error)
}

func AddProductionImpl() *AddProductionService {
	a := &AddProductionService{
		Service: core.NewDefaultService(),
		// true implement
		Do: func(ctx context.Context, production *service.Production) (int64, error) {
			return 0, nil
		},
	}
	a.Endpoint = func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// coercion first
		req := request.(*pb.AddProductionRequest)
		// then convert to 'Do' param
		prod := &service.Production{
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
