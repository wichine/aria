package main

import (
	"aria/core"
	ariaCfg "aria/core/config"
	"fmt"
)

type Swagger struct {
	WithSwagger bool
	Path        string
}

type Server struct {
	Port uint16
}

type ServiceDiscovery struct {
	EtcdServers []string
}

type Services map[string]string

type ServerConfig struct {
	Server
	Swagger
	Services
	ServiceDiscovery
	ariaCfg.Statistic
}

var globalConfig *ServerConfig

const defaultConfigFile = "./config/config.yaml"
const envPrefix = "WEB"

func InitConfig(configFile string) error {
	globalConfig = &ServerConfig{}
	err := ariaCfg.ParseConfig(envPrefix, configFile, globalConfig)
	if err != nil {
		return fmt.Errorf("init config parsing config error: %s", err)
	}

	fmt.Println("================ Config ================")
	for _, v := range core.Flatten(globalConfig) {
		fmt.Println("  ", v)
	}
	fmt.Println("========================================")

	return nil
}

func Config() *ServerConfig {
	if globalConfig != nil {
		return globalConfig
	}
	err := InitConfig(defaultConfigFile)
	if err != nil {
		panic(fmt.Errorf("config not be initialized,init default config error: %s", err))
	}
	return globalConfig
}
