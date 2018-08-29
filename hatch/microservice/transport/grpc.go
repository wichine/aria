package transport

import (
	"aria/hatch/microservice/endpoint"
	pb "aria/hatch/microservice/protocol/production"
	"context"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	grpccontext "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type grpcServer struct {
	addProductionHandler    grpctransport.Handler
	getAllProductionHandler grpctransport.Handler
}

// TODO:add grpc transport middleware arguments for grpc server
func NewGrpcServer(epts endpoint.Endpoints) pb.ProductionServiceServer {
	// add middlerware here, refer to Go kit gRPC Interceptor
	opts := []grpctransport.ServerOption{}

	return &grpcServer{
		addProductionHandler: grpctransport.NewServer(
			epts.AddProductionEndpoint,
			decodeAddProductionRequest,
			encodeAddproductionResponse,
			opts...,
		),
		getAllProductionHandler: grpctransport.NewServer(
			epts.GetAllProductionEndpoint,
			decodeGetAllProductionRequest,
			encodeGetAllproductionResponse,
			opts...,
		),
	}
}

func (g *grpcServer) AddProduction(ctx grpccontext.Context, request *pb.AddProductionRequest) (*pb.AddProductionResponse, error) {
	_, resp, err := g.addProductionHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AddProductionResponse), nil
}
func (g *grpcServer) GetAllProduction(ctx grpccontext.Context, request *pb.GetAllProductionRequest) (*pb.GetAllProductionResponse, error) {
	_, resp, err := g.getAllProductionHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetAllProductionResponse), nil
}

func decodeAddProductionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.AddProductionRequest)
	if !ok {
		return nil, fmt.Errorf("Error translate [request] to [pb.AddProductionRequest]")
	}
	return req, nil
}

func encodeAddproductionResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*pb.AddProductionResponse)
	if !ok {
		return nil, fmt.Errorf("Error translate [response] to [pb.AddProductionResponse]")
	}
	return resp, nil
}

func decodeGetAllProductionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.GetAllProductionRequest)
	if !ok {
		return nil, fmt.Errorf("Error translate [request] to [pb.GetAllProductionRequest]")
	}
	return req, nil
}

func encodeGetAllproductionResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*pb.GetAllProductionResponse)
	if !ok {
		return nil, fmt.Errorf("Error translate [response] to [pb.GetAllProductionResponse]")
	}
	return resp, nil
}

// Interceptor in "github.com/go-kit/kit/transport/grpc" has an error
func Interceptor(
	ctx grpccontext.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	ctx = grpccontext.WithValue(ctx, grpctransport.ContextKeyRequestMethod, info.FullMethod)
	return handler(ctx, req)
}
