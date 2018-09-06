package production

import (
	pb "aria/hatch/microservice/protocol/production"
	"google.golang.org/grpc"
)

func ServiceImpl() *ProductionService {
	return &ProductionService{
		AddProductionImpl(),
		GetAllProductionImpl(),
	}
}

type ProductionService struct {
	*AddProductionService
	*GetAllProductionService
}

func (ps *ProductionService) Register(server *grpc.Server) {
	pb.RegisterProductionServiceServer(server, ps)
}
