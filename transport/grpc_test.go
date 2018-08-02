package transport

import (
	pb "aria/protocol/production"
	grpccontext "golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestGrpc(t *testing.T) {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewProductionServiceClient(conn)
	for i := 0; i < 10; i++ {
		resp, err := client.AddProduction(grpccontext.TODO(), &pb.AddProductionRequest{
			Type:       "test",
			Code:       "00001",
			Name:       "测试产品",
			ValueDate:  time.Now().Unix(),
			DueDate:    time.Now().Unix() + int64(24*3600*100*time.Second),
			AnnualRate: 4,
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(resp)
	}

	resp1, err := client.GetAllProduction(grpccontext.TODO(), &pb.GetAllProductionRequest{})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp1)
}
