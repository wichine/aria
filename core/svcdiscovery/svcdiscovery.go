package svcdiscovery

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
)

type SvcDiscovery interface {
	Register(serviceKey, instanceAddress string) error
	DeRegister()
	Subscribe(serviceKey string, f sd.Factory) (endpoint.Endpoint, error)
}

var globalSD SvcDiscovery

func GetEtcdServiceDiscoveryInstance(servers []string) (SvcDiscovery, error) {
	var err error
	if globalSD != nil {
		return globalSD, nil
	}

	globalSD, err = NewEtcdServiceDiscovery(EtcdConfig{
		Servers: servers,
		Options: DefaultEtcdOptions,
	})
	if err != nil {
		globalSD = nil
	}
	return globalSD, err
}
