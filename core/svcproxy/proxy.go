package svcproxy

import (
	"aria/core/log"
	"aria/core/svcdiscovery"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitsd "github.com/go-kit/kit/sd"
)

var logger = log.GetLogger("ariaServiceProxy")

var serviceMap = map[string]*ServiceProxy{}

// Step2: register service proxy object to global map
func RegisterServices(services ...*ServiceProxy) {
	for _, s := range services {
		serviceMap[s.method] = s
	}
}

type ServiceProxy struct {
	key         string
	method      string
	factory     kitsd.Factory
	endpoint    endpoint.Endpoint
	initialized bool
}

// Step1: create a new service proxy object
func NewServiceProxy(key, method string, f kitsd.Factory) *ServiceProxy {
	return &ServiceProxy{
		key:         key,
		method:      method,
		factory:     f,
		endpoint:    endpoint.Nop,
		initialized: false,
	}
}

func (sr *ServiceProxy) initialize(sd svcdiscovery.SvcDiscovery) error {
	var err error
	sr.endpoint, err = sd.Subscribe(sr.key, sr.factory)
	if err != nil {
		return fmt.Errorf("subcribe service for [%s] error: %s", sr.method, err)
	}
	sr.initialized = true
	return nil
}

func (sr *ServiceProxy) serviceFullName() string {
	return fmt.Sprintf("%s -> %s", sr.key, sr.method)
}

func (sr *ServiceProxy) Call(request interface{}) (response interface{}, err error) {
	if !sr.initialized {
		err = fmt.Errorf("service [%s] not be initialized.", sr.serviceFullName())
		return
	}
	response, err = sr.endpoint(context.TODO(), request)
	return
}

// Step3: call this method to subcribe from service discovery component
// Notice: should be called at start of main
func InitServiceProxy(servers []string) error {
	// get the global service discovery instance
	sd, err := svcdiscovery.GetEtcdServiceDiscoveryInstance(servers)
	if err != nil {
		return fmt.Errorf("InitServiceProxy get etcd service discovery instance error: %s", err)
	}

	// subcribe necessary services from sd
	serviceNames := []string{}
	for _, sr := range serviceMap {
		err = sr.initialize(sd)
		if err != nil {
			return fmt.Errorf("init service proxy error: %s", err)
		}
		serviceNames = append(serviceNames, fmt.Sprintf("\n    %s", sr.serviceFullName()))
	}
	serviceNames = append(serviceNames, "\n")
	logger.Infof("Service proxy initialized: %v", serviceNames)
	return nil
}
