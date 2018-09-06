package svcdiscovery

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"os"
	"strings"
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
	etcdEnv := os.Getenv("ETCD_SERVERS")
	servers := strings.Split(etcdEnv, ",")
	globalSD, err = NewEtcdServiceDiscovery(EtcdConfig{
		Servers: servers,
		Options: DefaultEtcdOptions,
	})
	if err != nil {
		globalSD = nil
	}
	return globalSD, err
}
