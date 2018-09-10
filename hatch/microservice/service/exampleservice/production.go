package exampleservice

import (
	pb "aria/hatch/microservice/protocol/example"
	"google.golang.org/grpc"
)

// FIXME: 本服务初始化方法，修改为正确的命名
func ServiceImpl() *ExampleService {
	return &ExampleService{
		AddProductionImpl(),
		GetAllProductionImpl(),
	}
}

// FIXME: 实现proto里面的所有rpc接口，接口对象化，每个接口一个文件，包含了service、endpoint、transport三层
type ExampleService struct {
	*AddProductionService
	*GetAllProductionService
}

// FIXME: 调用proto中正确的注册方法
func (ps *ExampleService) Register(server *grpc.Server) {
	pb.RegisterProductionServiceServer(server, ps)
}

// FIXME: 下面定义本服务需要用到的数据结构
type Production struct {
	Type       string
	Code       string
	Name       string
	ValueDate  int64
	DueDate    int64
	AnnualRate int64
}
