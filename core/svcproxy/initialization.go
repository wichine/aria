package svcproxy

import (
	"aria/core"
	"google.golang.org/grpc"
	"strings"
)

//**********************************************
//    used to call other service
//**********************************************

var proxyCenter = map[string]proxyInitFactory{}

type proxyInitFactory func(serviceKey string)

// serviceName: key of config.serviceProxy
// namespace: ptr used to call grpc method,field name should be the same as service in protobuf
// serviceImpl: ptr of the service grpc server handler
func RegisterServiceDiscriptionToProxyCenter(serviceName string, namespace, serviceImpl interface{}) {
	proxyCenter[strings.ToLower(serviceName)] = func(serviceKey string) {
		initializeAllServiceInOneNameSpace(namespace, serviceKey, serviceImpl)
	}
}

func InitializeAllServiceInProxyCenterWithConfig(serviceMapInConfig map[string]string, sdServers []string) error {
	if serviceMapInConfig == nil || len(serviceMapInConfig) == 0 {
		logger.Warningf("no service found in service map.")
		return nil
	}
	for name, key := range serviceMapInConfig {
		if initFunc, ok := proxyCenter[strings.ToLower(name)]; ok {
			initFunc(key)
			logger.Debugf("register proxy [%s] to svcproxy.", name)
		} else {
			logger.Errorf("proxy namespace [%s] not found in proxyCenter.", name)
		}
	}
	logger.Debug("subscribing service in svcproxy...")
	return InitServiceProxy(sdServers)
}

//**********************************************
//    used by apigateway
//**********************************************

var servicesFactory = map[string]serviceFactory{}

type serviceFactory func(serviceKey string) core.Transport

// serviceName: key of config.serviceProxy
// serviceImpl: ptr of the service grpc server handler
func RegisterServiceDiscriptionToServicesFactory(serviceName string, serviceImpl interface{}) {
	servicesFactory[strings.ToLower(serviceName)] = func(serviceKey string) core.Transport {
		return ConvertServiceToProxyServer(serviceKey, serviceImpl)
	}
}

func RegisterAllServiceInServicesFactory(server *grpc.Server, serviceMap map[string]string) {
	if servicesFactory == nil || len(servicesFactory) == 0 {
		logger.Warningf("no service found in service map.")
		return
	}
	for name, key := range serviceMap {
		if factory, ok := servicesFactory[strings.ToLower(name)]; ok {
			service := factory(key)
			service.Register(server)
			logger.Infof("register service [%s] to grpc server.", name)
		} else {
			logger.Errorf("service factory of [%s] not found in service factory map!", name)
		}
	}
	return
}
