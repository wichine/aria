package svcproxy

import (
	"aria/core/config"
	"aria/core/svcdiscovery"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
)

type Connector interface {
	AddMiddleware(...endpoint.Middleware)
	Proxy() sd.Factory
}

var isServiceWrapMiddleware = map[string]bool{}

func ConvertServiceToProxy(serviceKey string, connectors ...Connector) error {
	if isServiceWrapMiddleware[serviceKey] {
		return fmt.Errorf("service [%s] already wrap middleware", serviceKey)
	}
	for _, connector := range connectors {
		// default middleware
		dmw, err := makeDefaultMiddleware(serviceKey, connector.Proxy())
		if err != nil {
			return fmt.Errorf("add middleware error: %s", err)
		}
		// TODO: add other middlewares
		connector.AddMiddleware(logMiddleware)
		// must add default middleware at the end
		connector.AddMiddleware(dmw)
	}
	isServiceWrapMiddleware[serviceKey] = true
	return nil
}

func makeDefaultMiddleware(serviceKey string, factory sd.Factory) (endpoint.Middleware, error) {
	sd, err := svcdiscovery.GetEtcdServiceDiscoveryInstance(config.Config().EtcdServers)
	if err != nil {
		return nil, fmt.Errorf("getProxyEndpoint get etcd instance error: %s", err)
	}
	proxyEndpoint, err := sd.Subscribe(serviceKey, factory)
	if err != nil {
		return nil, fmt.Errorf("getProxyEndpoint subscribe for %s error: %s", serviceKey, err)
	}
	mw := func(ep endpoint.Endpoint) endpoint.Endpoint {
		return proxyEndpoint
	}
	return mw, nil
}

func logMiddleware(ep endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = ep(ctx, request)
		logger.Debugf("[Proxy Log] mothod: %v, request: %v, response: %v, error: %s", ctx.Value("FullMethod"), request, response, err)
		return
	}
}
