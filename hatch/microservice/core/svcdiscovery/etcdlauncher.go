package svcdiscovery

import (
	"aria/hatch/microservice/core/config"
	"context"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"net/http"
	"os"
	"path/filepath"
)

var client *docker.Client
var etcdDockerImageRepository = "gcr.io/etcd-development/etcd"
var etcdDockerImageTag = "v3.3.9"
var etcdDockerImageTarUrl = "http://192.168.9.251:9000/project/tfinance/docker-images/etcd_v3.3.9_docker_image.tar"
var etcdContainerName = "etcdv3"

func StartEtcdServer() error {
	err := initDockerClient()
	if err != nil {
		return fmt.Errorf("start etcd error: %s", err)
	}
	if isEtcdContainerRunning() {
		return nil
	}
	if !isEtcdImageExist() {
		err = getEtcdImages()
		if err != nil {
			return fmt.Errorf("start etcd error: %s", err)
		}
	}
	err = startEtcd()
	if err != nil {
		return fmt.Errorf("start etcd error: %s", err)
	}
	return nil
}

func initDockerClient() error {
	if client != nil {
		return nil
	}
	var err error
	client, err = docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return fmt.Errorf("init docker client error: %s", err)
	}
	return nil
}

func isEtcdImageExist() bool {
	image, err := client.InspectImage(fmt.Sprintf("%s:%s", etcdDockerImageRepository, etcdDockerImageTag))
	if err != nil || image == nil {
		return false
	}
	return true
}

func getEtcdImages() error {
	// try to get tar from internal file server
	resp, err := http.Get(etcdDockerImageTarUrl)
	defer resp.Body.Close()
	if err != nil {
		// try to get from official registry
		err := client.PullImage(docker.PullImageOptions{
			Repository: etcdDockerImageRepository,
			Tag:        etcdDockerImageTag,
		}, docker.AuthConfiguration{})
		if err != nil {
			return fmt.Errorf("get etcd docker image error: %s", err)
		}
		return nil
	}

	err = client.LoadImage(docker.LoadImageOptions{
		InputStream:  resp.Body,
		OutputStream: os.Stderr,
		Context:      context.Background(),
	})
	if err != nil {
		return fmt.Errorf("load etcd docker image error: %s", err)
	}
	return nil
}

func isEtcdContainerRunning() bool {
	container, err := client.InspectContainer(etcdContainerName)
	if err != nil {
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    etcdContainerName,
			Force: true,
		})
		return false
	}
	if container.State.Running {
		return true
	}
	client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    etcdContainerName,
		Force: true,
	})
	return false
}

func startEtcd() error {
	cfg := config.Config().ServiceDiscovery.EtcdPeerConfig
	dataSourcePath, err := filepath.Abs(cfg.DataDir)
	if err != nil {
		return fmt.Errorf("start etcd: get abs of datadir error: %s", err)
	}
	err = os.MkdirAll(dataSourcePath, 0755)
	if err != nil {
		return fmt.Errorf("start etcd: make dir error: %s", err)
	}
	_, err = client.CreateContainer(docker.CreateContainerOptions{
		Name: etcdContainerName,
		Config: &docker.Config{
			Image: fmt.Sprintf("%s:%s", etcdDockerImageRepository, etcdDockerImageTag),
			Cmd: []string{
				"/usr/local/bin/etcd",
				"--name", cfg.Name,
				"--data-dir", "/etcd-data",
				"--listen-client-urls", cfg.ListenClientUrls,
				"--advertise-client-urls", cfg.AdvertiseClientUrls,
				"--listen-peer-urls", cfg.ListenPeerUrls,
				"--initial-advertise-peer-urls", cfg.InitialAdvertisePeerUrls,
				"--initial-cluster", cfg.InitialCluster,
				"--initial-cluster-token", cfg.InitialClusterToken,
				"--initial-cluster-state", cfg.InitialClusterState,
			},
		},
		HostConfig: &docker.HostConfig{
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port("2379/tcp"): []docker.PortBinding{
					docker.PortBinding{HostIP: "", HostPort: "2379"},
				},
				docker.Port("2380/tcp"): []docker.PortBinding{
					docker.PortBinding{HostIP: "", HostPort: "2380"},
				},
			},
			Mounts: []docker.HostMount{
				docker.HostMount{
					Type:   "bind",
					Source: dataSourcePath,
					Target: "/etcd-data",
				},
			},
		},
		NetworkingConfig: &docker.NetworkingConfig{},
		Context:          context.Background(),
	})
	if err != nil {
		return fmt.Errorf("create etcd container error: %s", err)
	}

	err = client.StartContainer(etcdContainerName, nil)
	if err != nil {
		return fmt.Errorf("start etcd container error: %s", err)
	}
	return nil
}
