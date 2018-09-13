package service

import "google.golang.org/grpc"

var registerFuncs = map[string]func(serviceKey string, server *grpc.Server) error{}

func RegisterAllService(server *grpc.Server) error {
	for serviceKey, register := range registerFuncs {
		err := register(serviceKey, server)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllServiceKeys() []string {
	keys := []string{}
	for key, _ := range registerFuncs {
		keys = append(keys, key)
	}
	return keys
}
