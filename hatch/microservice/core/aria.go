package core

import (
	"aria/hatch/microservice/core/config"
	"aria/hatch/microservice/core/svcdiscovery"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
	"strings"
)

type Aria struct {
	Config       *config.AriaConfig
	GrpcServer   *grpc.Server
	GrpcListener net.Listener
	HttpEngine   *gin.Engine
}

func New(config *config.AriaConfig) *Aria {
	port := strings.Split(config.Address, ":")[1]
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	return &Aria{
		Config:       config,
		GrpcListener: lis,
		GrpcServer:   grpc.NewServer(),
	}
}

func (a *Aria) RegisterAll(ts ...Transport) {
	for _, s := range ts {
		s.Register(a.GrpcServer)
	}
}
func (a *Aria) ServeGRPC() error {
	if err := a.GrpcServer.Serve(a.GrpcListener); err != nil {
		return err
	}
	return nil
}
func (a *Aria) Run() error {
	// register service
	if a.Config.ServiceDiscovery.Enable {
		// register to etcd service discovery center
		s, err := svcdiscovery.GetEtcdServiceDiscoveryInstance()
		if err != nil {
			panic(err)
		}
		err = s.Register(a.Config.ServiceKey, a.Config.Address)
		if err != nil {
			panic(err)
		}
		defer s.DeRegister()
	}
	// serve grpc / http
	if err := a.GrpcServer.Serve(a.GrpcListener); err != nil {
		return err
	}
	return nil
}

type Transport interface {
	Register(s *grpc.Server)
}
