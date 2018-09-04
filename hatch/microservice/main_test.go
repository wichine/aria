package main

import (
	pb "aria/hatch/microservice/protocol/production"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"testing"
)

func Test_main(t *testing.T) {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewProductionServiceClient(conn)
	resp, err := c.GetAllProduction(context.Background(), &pb.GetAllProductionRequest{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.Production)
}
