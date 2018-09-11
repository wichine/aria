package svcdiscovery

import (
	"aria/hatch/microservice/core/config"
	"testing"
)

func Test_isEtcdImageExist(t *testing.T) {
	initDockerClient()
	exist := isEtcdImageExist()
	t.Log(exist)
}

func Test_getEtcdImages(t *testing.T) {
	initDockerClient()
	if err := getEtcdImages(); err != nil {
		t.Error(err)
	}
}

func Test_All(t *testing.T) {
	config.InitConfig("../../config/config.yaml")
	initDockerClient()
	if err := StartEtcdServer(); err != nil {
		t.Error(err)
	}
}
