package service

import (
	"aria/core/config"
	"aria/core/log"
	"aria/core/svcproxy"
)

var initFuncs = map[string]func(serviceKey string){}
var logger = log.GetLogger("service")

// must call at start of main
func InitService(serviceMap map[string]string) error {
	if serviceMap == nil || len(serviceMap) == 0 {
		logger.Warningf("no service found in service map.")
		return nil
	}
	for name, key := range serviceMap {
		if initFunc, ok := initFuncs[name]; ok {
			initFunc(key)
			logger.Debugf("register proxy [%s] to svcproxy.", name)
		} else {
			logger.Errorf("proxy namespace [%s] not found in initFuncs.", name)
		}
	}
	return svcproxy.InitServiceProxy(config.Config().EtcdServers)
}
