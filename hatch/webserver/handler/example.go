package handler

import (
	"aria/hatch/webserver/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"service_generated_by_aria/protocol/example"
)

type AddProductionRequest struct {
	// 类型
	Type string `json:"type,omitempty"`
	// 编号
	Code       string `json:"code,omitempty"`
	Name       string `json:"name,omitempty"`
	ValueDate  int64  `json:"valueDate,omitempty"`
	DueDate    int64  `json:"dueDate,omitempty"`
	AnnualRate int64  `json:"annualRate,omitempty"`
}

type Production AddProductionRequest

type AddProductionResponse struct {
	// 返回状态 200/500
	Status int `json:"status,omitempty"`
	// 成功/错误 信息
	Msg string `json:"msg,omitempty"`
}

type GetProductionResponse struct {
	// 返回状态 200/500
	Status int `json:"status,omitempty"`
	// 成功/错误 信息
	Msg string `json:"msg,omitempty"`
	// 成功返回产品信息
	Payload []Production `json:payload,omitempty`
}

// @Tags Production
// @Summary 添加产品
// @Param Request body handler.AddProductionRequest true "请求Payload"
// Accept application/json
// Produce application/json
// @Success 200 {object} handler.AddProductionResponse
// @Failure 500 {object} handler.AddProductionResponse
// @Router /production/add [post]
func AddProduction(c *gin.Context) {
	request := &AddProductionRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, AddProductionResponse{500, err.Error()})
		return
	}
	resp, err := service.AddProduction.Call(&example.AddProductionRequest{
		Type:       request.Type,
		Code:       request.Code,
		Name:       request.Name,
		ValueDate:  request.ValueDate,
		DueDate:    request.DueDate,
		AnnualRate: request.AnnualRate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, AddProductionResponse{500, fmt.Sprintf("call prduction service error: %s", err)})
		return
	}
	c.JSON(http.StatusOK, AddProductionResponse{200, resp.(*example.AddProductionResponse).Msg})
}

// @Tags Production
// @Summary 获取产品
// @Param id path string true "请求产品编号,all获取全部"
// Produce application/json
// @Success 200 {object} handler.GetProductionResponse
// @Failure 500 {object} handler.GetProductionResponse
// @Router /production/get/{id} [get]
func GetProduction(c *gin.Context) {
	id := c.Param("id")
	if id == "all" {
		resp, err := service.GetAllProduction.Call(&example.GetAllProductionRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("call production service error: %s", err))
			return
		}
		payload := []Production{}
		for _, p := range resp.(*example.GetAllProductionResponse).Production {
			payload = append(payload, Production{
				p.Type,
				p.Code,
				p.Name,
				p.ValueDate,
				p.DueDate,
				p.AnnualRate,
			})
		}
		c.JSON(http.StatusOK, GetProductionResponse{200, "get production success", payload})
		return
	}
	c.JSON(http.StatusBadRequest, GetProductionResponse{500, fmt.Sprintf("param [%s] not supported.", id), nil})
}
