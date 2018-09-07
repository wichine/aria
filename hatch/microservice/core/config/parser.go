package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var globalConfig *MainConfig
var defaultConfigPath = "./config/config.yaml"
var envPrefix = "SERVICE"

func Config() *MainConfig {
	if globalConfig == nil {
		InitConfig(defaultConfigPath)
	}
	return globalConfig
}

func InitConfig(configPath string) error {
	v := newViper()
	v.SetConfigFile(configPath)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read in config error: %s", err)
	}
	globalConfig = &MainConfig{}
	v.UnmarshalKey("service", &globalConfig.Service)
	v.UnmarshalKey("serviceDiscovery", &globalConfig.ServiceDiscovery)
	printConfig(globalConfig)
	return nil
}

func newViper() *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	return v
}

func printConfig(config interface{}) {
	fmt.Printf("Config: %v\n", config)
}
