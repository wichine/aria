package svcdiscovery

//
// import (
// 	"fmt"
// 	docker "github.com/fsouza/go-dockerclient"
// )
//
// var client *docker.Client
// var etcdDockerImageName = ""
//
// func StartEtcdServer() {
//
// }
//
// func initDockerClient() error {
// 	var err error
// 	client, err = docker.NewClient("unix:///var/run/docker.sock")
// 	if err != nil {
// 		return fmt.Errorf("init docker client error: %s", err)
// 	}
// 	return nil
// }
//
// func isEtcdImageExist() bool {
// 	client.InspectImage("")
// 	return false
// }
