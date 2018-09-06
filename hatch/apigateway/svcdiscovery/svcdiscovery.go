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
