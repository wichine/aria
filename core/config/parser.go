package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var globalConfig *AriaConfig
var defaultConfigPath = "./config/config.yaml"
var envPrefix = "SERVICE"

func Config() *AriaConfig {
	if globalConfig == nil {
		InitConfig(defaultConfigPath)
	}
	return globalConfig
}

func InitConfig(configPath string) error {
	globalConfig = &AriaConfig{}
	return ParseConfig(envPrefix, configPath, globalConfig)
}

func newViper(envPrefix string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	return v
}

func ParseConfig(envPrefix, ymlFile string, object interface{}) error {
	v := newViper(envPrefix)
	v.SetConfigFile(ymlFile)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read in config error: %s", err)
	}
	err = v.Unmarshal(object)
	if err != nil {
		return fmt.Errorf("unmarshal config to object error: %s", err)
	}
	return nil
}
