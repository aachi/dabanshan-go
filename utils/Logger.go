package utils

import (
	"os"

	"github.com/go-kit/kit/log"
)

// NewLogger ...
func NewLogger() log.Logger {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	return logger
}
