package svcdiscovery

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitsd "github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"os"
	"time"
)

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(os.Stderr)
	log.With(logger, "time", log.DefaultTimestampUTC)
	log.With(logger, "caller", log.DefaultCaller)
}

type EtcdConfig struct {
	Servers []string
	Options etcdv3.ClientOptions
}

var DefaultEtcdOptions = etcdv3.ClientOptions{
	// Path to trusted ca file
	CACert: "",
	// Path to certificate
	Cert: "",
	// Path to private key
	Key: "",
	// Username if required
	Username: "",
	// Password if required
	Password: "",
	// If DialTimeout is 0, it defaults to 3s
	DialTimeout: time.Second * 3,
	// If DialKeepAlive is 0, it defaults to 3s
	DialKeepAlive: time.Second * 3,
}

type etcdv3SD struct {
	config    EtcdConfig
	registrar *etcdv3.Registrar
	client    etcdv3.Client
}

func NewEtcdServiceDiscovery(config EtcdConfig) (SvcDiscovery, error) {
	sd := &etcdv3SD{
		config: config,
	}
	client, err := etcdv3.NewClient(context.Background(), sd.config.Servers, sd.config.Options)
	if err != nil {
		return nil, fmt.Errorf("create new client of etcdv3 error: %s", err)
	}
	sd.client = client
	return sd, nil
}

func (sd *etcdv3SD) Register(serviceKey, instanceAddress string) error {
	service := etcdv3.Service{
		Key:   fmt.Sprintf("%s/%s", serviceKey, instanceAddress),
		Value: instanceAddress,
		TTL:   etcdv3.NewTTLOption(3*time.Second, 10*time.Second),
	}
	registrar := etcdv3.NewRegistrar(sd.client, service, logger)
	sd.registrar = registrar
	registrar.Register()
	return nil
}

func (sd *etcdv3SD) DeRegister() {
	sd.registrar.Deregister()
}

func (sd *etcdv3SD) Subscribe(serviceKey string, f kitsd.Factory) (endpoint.Endpoint, error) {
	instancer, err := etcdv3.NewInstancer(sd.client, serviceKey, logger)
	if err != nil {
		return nil, fmt.Errorf("create new instancer error: %s", err)
	}
	endpointer := kitsd.NewEndpointer(instancer, f, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 500*time.Millisecond, balancer)
	return retry, nil
}
