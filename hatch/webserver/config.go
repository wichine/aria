package main

import (
	"aria/core"
	"aria/core/config"
	"fmt"
)

type Swagger struct {
	GenerateDocs    bool
	SwaggerServerOn bool
}

type Server struct {
	Port uint16
}

type ServiceDiscovery struct {
	EtcdServers []string
}

type ServerConfig struct {
	Server
	Swagger
	ServiceDiscovery
}

var globalConfig *ServerConfig

const defaultConfigFile = "./config/config.yaml"
const envPrefix = "WEB"

func InitConfig(configFile string) error {
	globalConfig = &ServerConfig{}
	err := config.ParseConfig(envPrefix, configFile, globalConfig)
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
