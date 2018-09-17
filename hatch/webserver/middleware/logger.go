package middleware

import (
	"aria/core/log"
	"github.com/gin-gonic/gin"
)

var logModule = "GinLogger"
var logger = log.GetLogger(logModule)

type loggerBackend struct {
}

func Logger() gin.HandlerFunc {
	gin.DisableConsoleColor()
	return gin.LoggerWithWriter(loggerBackend{})
}

func (l loggerBackend) Write(p []byte) (n int, err error) {
	logger.Notice(string(p))
	return len(p), nil
}
