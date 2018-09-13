package middleware

import (
	"aria/core/svcdiscovery"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
)

type Connector interface {
	AddMiddleware(...endpoint.Middleware)
	Proxy() sd.Factory
}

var isServiceWrapMiddleware = map[string]bool{}

func WrapMiddleware(serviceKey string, connectors ...Connector) error {
	if isServiceWrapMiddleware[serviceKey] {
		return fmt.Errorf("service [%s] already wrap middleware", serviceKey)
	}
	for _, connector := range connectors {
		// default middleware
		ep, err := makeDefaultMiddleware(serviceKey, connector.Proxy())
		if err != nil {
			return fmt.Errorf("add middleware error: %s", err)
		}
		connector.AddMiddleware(ep)
		// TODO: add other middlewares
	}
	isServiceWrapMiddleware[serviceKey] = true
	return nil
}

func makeDefaultMiddleware(serviceKey string, factory sd.Factory) (endpoint.Middleware, error) {
	sd, err := svcdiscovery.GetEtcdServiceDiscoveryInstance()
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
