package production

import (
	"aria/hatch/microservice/core"
	pb "aria/hatch/microservice/protocol/production"
	"context"
	"google.golang.org/grpc"
)

func ServiceImpl() *ProductionService {
	return &ProductionService{
		AddProductionImpl(core.Common),
		GetAllProductionImpl(core.Common),
	}
}

type ProductionService struct {
	*AddProductionService
	*GetAllProductionService
}

func (ps *ProductionService) AddProduction(ctx context.Context, request *pb.AddProductionRequest) (*pb.AddProductionResponse, error) {
	_, resp, err := ps.AddProductionService.Transport().Grpc.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AddProductionResponse), nil
}

func (ps *ProductionService) GetAllProduction(ctx context.Context, request *pb.GetAllProductionRequest) (*pb.GetAllProductionResponse, error) {
	_, resp, err := ps.GetAllProductionService.Transport().Grpc.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetAllProductionResponse), nil
}

func (ps *ProductionService) Register(server *grpc.Server) {
	pb.RegisterProductionServiceServer(server, ps)
}
