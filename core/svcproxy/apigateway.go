package svcproxy

import (
	"aria/core"
	"fmt"
	"reflect"
)

// convert a service implement (grpc server handler) to a proxy server (commonly use as apigateway)
func ConvertServiceToProxyServer(serviceKey string, serviceImpl interface{}) core.Transport {
	sPtr := reflect.ValueOf(serviceImpl)
	if sPtr.Kind().String() != "ptr" {
		panic("serviceImpl should be a ptr of a service implement object!")
	}
	sValue := sPtr.Elem()
	sType := sValue.Type()
	connectorType := reflect.TypeOf((*Connector)(nil)).Elem()

	// convert service to proxy
	connectors := []Connector{}
	for i := 0; i < sValue.NumField(); i++ {
		if sValue.Field(i).Type().Implements(connectorType) {
			logger.Debugf("converting service method to proxy connector: %s -> %s", sType.String(), sType.Field(i).Name)
			connectors = append(connectors, sValue.Field(i).Interface().(Connector))
		} else {
			logger.Warningf("can not convert field [%s -> %s] to type [%s], continue...", sType.String(), sType.Field(i).Name, connectorType.String())
			continue
		}
	}
	err := ConvertServiceToProxy(serviceKey, connectors...)
	if err != nil {
		panic(fmt.Sprintf("convert service to proxy errror: %s", err))
	}

	// register service to grpc server
	ariaTransportType := reflect.TypeOf((*core.Transport)(nil)).Elem()
	if sPtr.Type().Implements(ariaTransportType) {
		logger.Debugf(`convert service [ subcribe from: "%s", service proxy implement: "%s" ] to %s`, serviceKey, sPtr.Type().String(), ariaTransportType.String())
		return serviceImpl.(core.Transport)
	} else {
		panic(fmt.Sprintf("can not convert service %s to %s", sPtr.Type().String(), ariaTransportType.String()))
	}
}
