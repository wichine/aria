package main

import (
	"aria/hatch/microservice/endpoint"
	pb "aria/hatch/microservice/protocol/production"
	"aria/hatch/microservice/service"
	"aria/hatch/microservice/transport"
	"google.golang.org/grpc"
	"net"
)

func main() {
	ps := service.NewProductionService()
	endpoints := endpoint.MakeAllEndpoints(ps)
	go transport.StartHttpServer(*endpoints)

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}
	productionGrpcServer := transport.NewGrpcServer(*endpoints)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(transport.Interceptor))
	pb.RegisterProductionServiceServer(grpcServer, productionGrpcServer)
	grpcServer.Serve(lis)
}
