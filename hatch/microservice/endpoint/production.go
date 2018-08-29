package endpoint

import (
	pb "aria/hatch/microservice/protocol/production"
	"aria/hatch/microservice/service"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

func MakeAddProductionEndpoint(s service.ProductionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.AddProductionRequest)
		prod := &service.Production{
			Type:       req.Type,
			Code:       req.Code,
			Name:       req.Name,
			ValueDate:  req.ValueDate,
			DueDate:    req.DueDate,
			AnnualRate: req.AnnualRate,
		}
		count, err := s.AddProduction(ctx, prod)
		if err != nil {
			return nil, err
		}
		return &pb.AddProductionResponse{Status: 0, Msg: fmt.Sprintf("Total count： %d", count)}, nil
	}
}

func MakeGetAllProductionEndpoint(s service.ProductionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ps, err := s.GetAllProduction(ctx)
		if err != nil {
			return nil, err
		}
		result := []*pb.Production{}
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
