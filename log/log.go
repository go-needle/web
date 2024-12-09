package log

import (
	"github.com/go-needle/log"
	"io"
)

// logger global
var logger = log.New()

// log methods
var (
	Debug  = logger.Debug
	Debugf = logger.Debugf
	Info   = logger.Info
	Infof  = logger.Infof
	Warn   = logger.Warn
	Warnf  = logger.Warnf
	Error  = logger.Error
	Errorf = logger.Errorf
	Fatal  = logger.Fatal
	Fatalf = logger.Fatalf
)

// Set controls log level and output which is only for web
func Set(level int, out io.Writer) {
	logger.Set(level, out)
}
