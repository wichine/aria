package svcproxy

import (
	"aria/core/config"
	"aria/core/middleware"
	"aria/core/svcdiscovery"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
)

type Connector interface {
	AddMiddleware(...endpoint.Middleware)
	Proxy() sd.Factory
}

func ConvertServiceToProxy(serviceKey string, connectors ...Connector) error {
	for _, connector := range connectors {
		// default middleware
		dmw, err := makeDefaultMiddleware(serviceKey, connector.Proxy())
		if err != nil {
			return fmt.Errorf("add middleware error: %s", err)
		}
		// TODO: add other middlewares
		connector.AddMiddleware(middleware.LogMiddleware(logger))
		connector.AddMiddleware(middleware.ZipkinMiddleware)
		// must add default middleware at the end
		connector.AddMiddleware(dmw)
	}
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

type nilConnector struct{}

func (*nilConnector) AddMiddleware(...endpoint.Middleware) {}
func (*nilConnector) Proxy() sd.Factory                    { return nil }
