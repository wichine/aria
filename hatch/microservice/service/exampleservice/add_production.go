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

// FIXME: 命名修改
type AddProductionService struct {
	*core.Service

	// FIXME: 每个接口的核心逻辑，参数格式根据需要修改
	Do func(ctx context.Context, production *Production) (int64, error)
}

// FIXME: 命名修改
func AddProductionImpl() *AddProductionService {
	a := &AddProductionService{
		Service: core.NewDefaultService(),

		// FIXME: service层核心逻辑在这里实现
		Do: func(ctx context.Context, production *Production) (int64, error) {
			return 0, nil
		},
	}

	a.Endpoint = func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// coercion first
		// FIXME: 协议转换格式修改为正确的格式
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
		// 调用service具体逻辑
		rawResp, err := a.Do(ctx, prod)
		if err != nil {
			return nil, err
		}
		// encode to true res
		// FIXME: 返回格式修改
		return &pb.AddProductionResponse{Status: 0, Msg: fmt.Sprintf("Total count : %d", rawResp)}, nil
	}
	return a
}

// FIXME: 实现正确的grpc接口
func (aps *AddProductionService) AddProduction(ctx context.Context, request *pb.AddProductionRequest) (*pb.AddProductionResponse, error) {
	r, err := aps.Compose()(ctx, request)
	if err != nil {
		return nil, err
	}
	// FIXME: 返回格式修改
	// coercion response
	res := r.(*pb.AddProductionResponse)
	return res, nil
}

// 实现grpc client的接口，instance为从服务发现中获取的提供服务的实际地址
func (aps *AddProductionService) Proxy() sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		address := instance
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, nil, fmt.Errorf("dial grpc error: %s", err)
		}
		// FIXME: 使用正确的client接口
		client := pb.NewProductionServiceClient(conn)
		ep := func(ctx context.Context, request interface{}) (interface{}, error) {
			// FIXME: 类型修改
			req, ok := request.(*pb.AddProductionRequest)
			if !ok {
				// FIXME: error中的类型修改
				return nil, fmt.Errorf("can not convert request to pb.AddProductionRequest")
			}
			// FIXME: 调用正确的方法
			return client.AddProduction(ctx, req)
		}
		return ep, conn, nil
	}
}
