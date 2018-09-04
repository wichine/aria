package core

import (
	"context"
	"github.com/gin-gonic/gin"
	gokitEndpoint "github.com/go-kit/kit/endpoint"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"net"
)

type EndpointType int

const (
	Common EndpointType = iota
	Proxy
)

type AriaCommon struct {
	EpType EndpointType
	Decode func(ctx context.Context, request interface{}) (interface{}, error)
	Encode func(ctx context.Context, response interface{}) (interface{}, error)
}
type AriaConfig struct {
	GrpcPort string
	HttpPort string
}

func NewAria(config AriaConfig) *Aria {
	lis, err := net.Listen("tcp", config.GrpcPort)
	if err != nil {
		panic(err)
	}
	return &Aria{
		GrpcListener: lis,
		GrpcServer:   grpc.NewServer(),
	}
}

type Aria struct {
	GrpcServer   *grpc.Server
	GrpcListener net.Listener
}

func (a *Aria) RegisterAll(services ...AriaService) {
	for _, s := range services {
		s.Register(a.GrpcServer)
	}
}
func (a *Aria) ServeGRPC() {
	if err := a.GrpcServer.Serve(a.GrpcListener); err != nil {
		panic(err)
	}
}

type AriaServiceHandler interface {
	Endpoint() gokitEndpoint.Endpoint
	Proxy() gokitEndpoint.Endpoint
	Transport() AriaTransport
}

type AriaService interface {
	Register(s *grpc.Server)
}

type AriaTransport struct {
	Http gin.HandlerFunc
	Grpc gokitgrpc.Handler
}
