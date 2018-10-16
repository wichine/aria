package service

import (
	"aria/core/svcproxy"
	"google.golang.org/grpc"
)

func RegisterAllService(server *grpc.Server, serviceMap map[string]string) {
	svcproxy.RegisterAllServiceInServicesFactory(server, serviceMap)
}
