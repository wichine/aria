package transport

import (
	"aria/hatch/microservice/endpoint"
	"aria/hatch/microservice/protocol"
	pb "aria/hatch/microservice/protocol/production"
	"fmt"
	"github.com/gin-gonic/gin"
	kitendpoint "github.com/go-kit/kit/endpoint"
)

const (
	ERR_CODE = 500
	SUC_CODE = 200
)

func StartHttpServer(eps endpoint.Endpoints) {
	eg := gin.New()
	eg.Use(gin.Logger(), gin.Recovery())
	eg.POST("/production/v1/add", newGinHandler(
		eps.AddProductionEndpoint,
		httpDecodeAddProductionRequest,
		httpEncodeAddProductionResponse,
		processError,
	))
	eg.GET("/production/v1/getall", newGinHandler(
		eps.GetAllProductionEndpoint,
		httpDecodeGetAllProductionRequest,
		httpEncodeGetAllProductionResponse,
		processError,
	))

	eg.Run(":8090")
}

func processError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}
	errResponse := protocol.ErrorResponse{int64(ERR_CODE), fmt.Sprintf("Error: %s", err)}
	ctx.AbortWithError(ERR_CODE, err)
	ctx.JSON(ERR_CODE, errResponse)
}

func httpDecodeAddProductionRequest(ctx *gin.Context) (interface{}, error) {
	req := &protocol.AddProductionRequest{}
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		return nil, err
	}
	endpointReq := &pb.AddProductionRequest{
		Type:       req.Type,
		Code:       req.Code,
		Name:       req.Name,
		ValueDate:  req.ValueDate,
		DueDate:    req.DueDate,
		AnnualRate: req.AnnualRate,
	}
	return endpointReq, nil
}

func httpEncodeAddProductionResponse(response interface{}) (interface{}, error) {
	resp := &protocol.AddProductionResponse{
		Code: 200,
		Msg:  response.(*pb.AddProductionResponse).Msg,
	}
	return resp, nil
}

func httpDecodeGetAllProductionRequest(_ *gin.Context) (interface{}, error) {
	return nil, nil
}

func httpEncodeGetAllProductionResponse(response interface{}) (interface{}, error) {
	ps := []protocol.Production{}
	for _, p := range response.(*pb.GetAllProductionResponse).Production {
		ps = append(ps, protocol.Production{
			Type:       p.Type,
			Code:       p.Code,
			Name:       p.Name,
			ValueDate:  p.ValueDate,
			DueDate:    p.DueDate,
			AnnualRate: p.AnnualRate,
		})
	}
	resp := &protocol.GetAllProductionResponse{
		Code:       200,
		Production: ps,
	}
	return resp, nil
}

type decodeRequestFunc func(*gin.Context) (interface{}, error)
type encodeResponseFunc func(interface{}) (interface{}, error)
type errorResponseFunc func(*gin.Context, error)

func newGinHandler(ep kitendpoint.Endpoint, decodeRequest decodeRequestFunc, encodeResponse encodeResponseFunc, processError errorResponseFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req, err := decodeRequest(ctx)
		if err != nil {
			processError(ctx, err)
			return
		}
		resp, err := ep(ctx, req)
		if err != nil {
			processError(ctx, err)
			return
		}
		result, err := encodeResponse(resp)
		if err != nil {
			processError(ctx, err)
			return
		}
		ctx.JSON(SUC_CODE, result)
	}
}
