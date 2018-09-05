package core

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
	"net"
)

type Service struct {
	Middleware []endpoint.Middleware
	Endpoint   endpoint.Endpoint
}

func (s *Service) WithMiddleware(m endpoint.Middleware) {
	s.Middleware = append(s.Middleware, m)
}

func (s *Service) Compose() endpoint.Endpoint {
	for _, m := range s.Middleware {
		s.Endpoint = m(s.Endpoint)
	}
	return s.Endpoint
}

func NewDefaultService() *Service {
	return &Service{
		Middleware: []endpoint.Middleware{},
	}
}

type AriaConfig struct {
	GrpcPort string
	HttpPort string
}

type Aria struct {
	GrpcServer   *grpc.Server
	GrpcListener net.Listener
	HttpEngine   *gin.Engine
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

func (a *Aria) RegisterAll(ts ...Transport) {
	for _, s := range ts {
		s.Register(a.GrpcServer)
	}
}
func (a *Aria) ServeGRPC() {
	if err := a.GrpcServer.Serve(a.GrpcListener); err != nil {
		panic(err)
	}
}

type Transport interface {
	Register(s *grpc.Server)
}
