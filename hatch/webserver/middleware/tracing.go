package middleware

import (
	"aria/core/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

func ZipkinTracing(url, serviceName, address string, on bool) gin.HandlerFunc {
	if err := tools.InitializeZipkin(url, serviceName, address, on); err != nil {
		panic(fmt.Sprintf("init zipkin middleware error: %s", err))
	}
	return func(c *gin.Context) {
		span := tools.GetZipkinSpanFromGinContext(c)
		span.SetName(fmt.Sprintf("handler: %s", c.HandlerName()))
		tools.SetZipkinSpanToGinContext(c, span)
		c.Next()
		span.Finish()
	}
}
