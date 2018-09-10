package log

import (
	"github.com/go-kit/kit/log"
	"github.com/op/go-logging"
	"io"
	"os"
)

var DefaultLogWriter io.Writer = os.Stdout
var DefaultGoKitLogger = log.NewLogfmtLogger(os.Stdout)

func GetLogger(module string) *logging.Logger {
	logger := logging.MustGetLogger(module)
	loggerBackend := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stdout, "", 0),
		logging.MustStringFormatter(`%{color}[%{time:2006-01-02 15:04:05.000}] %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`),
	)
	levelBackEnd := logging.AddModuleLevel(loggerBackend)
	logger.SetBackend(levelBackEnd)
	return logger
}
