package config

import "testing"

func Test_Config(t *testing.T) {
	err := InitConfig("../../config/config.yaml")
	if err != nil {
		t.Error(err)
	}
	config := Config()
	if config == nil {
		t.Error("config is nil")
	}
	t.Log(config)
}
