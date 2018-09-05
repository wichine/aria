package production

import (
	"aria/hatch/microservice/core"
	pb "aria/hatch/microservice/protocol/production"
	"context"
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
