package config

type Service struct {
	Type       string
	Address    string
	ServiceKey string
}

type ServiceDiscovery struct {
	EtcdServerOn   bool
	EtcdPeerConfig EtcdConfig
	EtcdServers    []string
}

type EtcdConfig struct {
	Name                     string
	DataDir                  string
	ListenPeerUrls           string
	ListenClientUrls         string
	InitialAdvertisePeerUrls string
	InitialCluster           string
	InitialClusterState      string
	InitialClusterToken      string
	AdvertiseClientUrls      string
}

type MainConfig struct {
	Service
	ServiceDiscovery
}
