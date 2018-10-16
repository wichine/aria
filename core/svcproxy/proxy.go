package svcproxy

import (
	"aria/core/log"
	"aria/core/svcdiscovery"
	"aria/core/tools"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitsd "github.com/go-kit/kit/sd"
	"reflect"
)

var logger = log.GetLogger("AriaServiceProxy")

var serviceMap = map[string]*ServiceProxy{}
var ProxyMethodName = "Proxy"

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

func (sr *ServiceProxy) GetServiceFullName() string {
	return fmt.Sprintf("%s -> %s", sr.key, sr.method)
}

func (sr *ServiceProxy) Call(ctx context.Context, request interface{}) (response interface{}, err error) {
	if !sr.initialized {
		err = fmt.Errorf("service [%s] not be initialized.", sr.GetServiceFullName())
		return
	}

	span := tools.GetZipkinSpanFromContext(ctx)
	span.SetName("proxy: call rpc")
	newCtx := tools.SetZipkinSpanToGrpcHeader(ctx, span)
	response, err = sr.endpoint(newCtx, request)
	span.Finish()
	return
}

// Step2: register service proxy object to global map
func RegisterServices(services ...*ServiceProxy) {
	for _, s := range services {
		serviceMap[s.method] = s
	}
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
		serviceNames = append(serviceNames, fmt.Sprintf("\n    %s", sr.GetServiceFullName()))
	}
	serviceNames = append(serviceNames, "\n")
	logger.Infof("Service proxy initialized: %v", serviceNames)
	return nil
}

// namespace: the pointer to a object that will be called as a proxy namespace
// serviceKey: the service key which will be used to subcribe service from service discovery component
// factory: the implement for the service imported used to search Proxy() method to register to svcproxy
func initializeAllServiceInOneNameSpace(namespace interface{}, serviceKey string, factory interface{}) {
	nsPtr := reflect.ValueOf(namespace)
	if nsPtr.Kind().String() != "ptr" {
		panic("namespace should be a ptr of a service proxy namespace object!")
	}
	fPtr := reflect.ValueOf(factory)
	if fPtr.Kind().String() != "ptr" {
		panic("factory should be a ptr of a service implement object!")
	}

	nsValue := nsPtr.Elem()
	nsType := nsValue.Type()
	fValue := fPtr.Elem()
	fType := fValue.Type()
	logger.Debugf("current namespace is: %s", nsType.String())
	logger.Debugf("current factory type is: %s", fType.String())

	// find method with the same name of ns field
	for i := 0; i < nsValue.NumField(); i++ {
		methodToFind := nsType.Field(i).Name
		logger.Debugf("searching method [%s] in %s:", methodToFind, fType.String())
		proxyMethod := reflect.Value{}
		methodFoundField := reflect.Value{}
		for j := 0; j < fType.NumField(); j++ {
			logger.Debugf("  =====> start | find method [%s] in %s -> %s", methodToFind, fType.String(), fType.Field(j).Name)

			methodFound := fValue.Field(j).MethodByName(methodToFind)
			if !methodFound.IsValid() {
				logger.Debugf("  <===== end | method [%s] NOT FOUND in %s -> %s, continue...", methodToFind, fType.String(), fType.Field(j).Name)
				continue
			}

			// method found,then get the proxy method
			methodFoundField = fValue.Field(j)
			proxyMethod = fValue.Field(j).MethodByName(ProxyMethodName)
			logger.Debugf("  <===== end | method [%s] FOUND in %s -> %s, break!", methodToFind, fType.String(), fType.Field(j).Name)
			break
		}
		// can not found method in factory
		if !methodFoundField.IsValid() {
			panic(fmt.Sprintf("can not find method [%s] in [%s]", methodToFind, fType.String()))
		}
		// method found but not has a Proxy method
		if !proxyMethod.IsValid() {
			panic(fmt.Sprintf("[%s] not have the method [%s]!", methodFoundField.Type().String(), ProxyMethodName))
		}

		result := proxyMethod.Call([]reflect.Value{})
		if len(result) != 1 {
			panic(fmt.Sprintf("%s -> %s() should return one result,but returns %d", methodFoundField.Type().String(), ProxyMethodName, len(result)))
		}

		// result[0].Interface().(kitsd.Factory)
		sdFactory, ok := result[0].Interface().(kitsd.Factory)
		if !ok {
			panic(fmt.Sprintf("%s -> %s() should return a sd.Factory,but returns %s", methodFoundField.Type().String(), ProxyMethodName, result[0].Type().String()))
		}
		// create service proxy objects to init value
		nsValue.Field(i).Set(reflect.ValueOf(NewServiceProxy(serviceKey, methodToFind, sdFactory)))
		// register service proxy object to global map
		RegisterServices(nsValue.Field(i).Interface().(*ServiceProxy))
	}
}
