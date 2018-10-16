package service

import (
	"aria/core/svcproxy"
)

// must call at start of main
func InitService(serviceMap map[string]string, etcdServers []string) error {
	return svcproxy.InitializeAllServiceInProxyCenterWithConfig(serviceMap, etcdServers)
}
