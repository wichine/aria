package otherservice

import (
	"aria/core/config"
	"aria/core/svcproxy"
)

// must call at start of main
func InitOtherService() error {
	return svcproxy.InitializeAllServiceInProxyCenterWithConfig(config.Config().ServiceProxy, config.Config().EtcdServers)
}
