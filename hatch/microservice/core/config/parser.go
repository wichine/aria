package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"reflect"
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
	v := newViper()
	v.SetConfigFile(configPath)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read in config error: %s", err)
	}
	globalConfig = &AriaConfig{}
	if err := v.Unmarshal(globalConfig); err != nil {
		panic(err)
	}
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
	params := Flatten(config)
	var buffer bytes.Buffer
	for i := range params {
		buffer.WriteString("\n\t")
		buffer.WriteString(params[i])
	}
	fmt.Printf("Aria config: %s\n", buffer.String())
}

func Flatten(i interface{}) []string {
	var res []string
	flatten("", &res, reflect.ValueOf(i))
	return res
}

const DELIMITER = "."

func flatten(k string, m *[]string, v reflect.Value) {
	delimiter := DELIMITER
	if k == "" {
		delimiter = ""
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			*m = append(*m, fmt.Sprintf("%s =", k))
			return
		}
		flatten(k, m, v.Elem())
	case reflect.Struct:
		if x, ok := v.Interface().(fmt.Stringer); ok {
			*m = append(*m, fmt.Sprintf("%s = %v", k, x))
			return
		}

		for i := 0; i < v.NumField(); i++ {
			flatten(k+delimiter+v.Type().Field(i).Name, m, v.Field(i))
		}
	case reflect.String:
		// It is useful to quote string values
		*m = append(*m, fmt.Sprintf("%s = \"%s\"", k, v))
	default:
		*m = append(*m, fmt.Sprintf("%s = %v", k, v))
	}
}
