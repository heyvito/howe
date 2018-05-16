package helpers

import (
	"github.com/hashicorp/go-syslog"
)

var logger gsyslog.Syslogger
var loggerRunning bool

func init() {
	var err error
	logger, err = gsyslog.NewLogger(gsyslog.LOG_ERR, "USER", "howe")
	loggerRunning = (err == nil)
}

// ReportError is used by Widgets to report internal errors to the syslog
func ReportError(data string) {
	if !loggerRunning {
		return
	}
	logger.Write([]byte(data))
}
