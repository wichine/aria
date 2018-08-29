package service

import (
	"context"
)

type ProductionService interface {
	AddProduction(ctx context.Context, prod *Production) (int64, error)
	GetAllProduction(ctx context.Context) ([]Production, error)
}

// service定义自己的请求、响应格式以便降低耦合;
// 上一层需要根据下一层的数据格式增加适配函数

type Production struct {
	Type       string
	Code       string
	Name       string
	ValueDate  int64
	DueDate    int64
	AnnualRate int64
}
