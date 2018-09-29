package otherservice

import (
	"aria/core/config"
	"aria/core/log"
	"aria/core/svcproxy"
)

var initFuncs = map[string]func(serviceKey string){}
var logger = log.GetLogger("otherservice")

// must call at start of main
func InitOtherService() error {
	for name, key := range config.Config().ServiceProxy {
		if initFunc, ok := initFuncs[name]; ok {
			initFunc(key)
			logger.Debugf("register proxy [%s] to svcproxy.", name)
		} else {
			logger.Errorf("proxy namespace [%s] not found in initFuncs.", name)
		}
	}
	return svcproxy.InitServiceProxy(config.Config().EtcdServers)
}
