package svcdiscovery

import (
	"aria/core/config"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
)

type SvcDiscovery interface {
	Register(serviceKey, instanceAddress string) error
	DeRegister()
	Subscribe(serviceKey string, f sd.Factory) (endpoint.Endpoint, error)
}

var globalSD SvcDiscovery

func GetEtcdServiceDiscoveryInstance() (SvcDiscovery, error) {
	var err error
	if globalSD != nil {
		return globalSD, nil
	}

	globalSD, err = NewEtcdServiceDiscovery(EtcdConfig{
		Servers: config.Config().EtcdServers,
		Options: DefaultEtcdOptions,
	})
	if err != nil {
		globalSD = nil
	}
	return globalSD, err
}
