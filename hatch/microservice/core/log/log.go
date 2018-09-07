package log

import (
	"github.com/go-kit/kit/log"
	"io"
	"os"
)

var DefaultLogWriter io.Writer = os.Stdout
var DefaultGoKitLogger = log.NewLogfmtLogger(os.Stdout)
