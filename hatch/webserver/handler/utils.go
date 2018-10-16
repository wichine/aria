package handler

import (
	"aria/core/log"
	"aria/core/svcproxy"
	"aria/core/tools"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var logger = log.GetLogger("webserver")

func LogForError(funcName string, pErr *error) {
	err := *pErr
	if err != nil {
		logger.Errorf("[%s] %s", funcName, err)
	}
}

func CallService(c *gin.Context, proxy *svcproxy.ServiceProxy, request interface{}) (interface{}, error) {
	span := tools.GetZipkinSpanFromGinContext(c)
	resp, err := proxy.Call(tools.SetZipkinSpanToContext(context.Background(), span), request)
	logger.Debugf("call service [%s] rpc response: %v", proxy.GetServiceFullName(), resp)
	return resp, err
}

func CallServiceWithResponseWriteToCtx(c *gin.Context, proxy *svcproxy.ServiceProxy, request interface{}) (interface{}, error) {
	span := tools.GetZipkinSpanFromGinContext(c)
	response, err := proxy.Call(tools.SetZipkinSpanToContext(context.Background(), span), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Call service [%s] error: %s", proxy.GetServiceFullName(), err))
		return nil, err
	}
	c.JSON(http.StatusOK, response)
	logger.Debugf("call service [%s] rpc response: %v", proxy.GetServiceFullName(), response)
	return response, nil
}
