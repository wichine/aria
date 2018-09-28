package main

import (
	"context"
	"google.golang.org/grpc"
	"service_generated_by_aria/protocol/example"
	"testing"
)

func TestGrpc(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	client := example.NewProductionServiceClient(conn)
	resp, err := client.AddProduction(context.Background(), &example.AddProductionRequest{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
	resp1, err := client.GetAllProduction(context.Background(), &example.GetAllProductionRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp1)
}
