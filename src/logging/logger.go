package logging

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/project-leaf/mq-service-go/src/config"
)

// GetLogger will construct and return a new logger instance based on the given configuration.
func GetLogger(c *config.Config) *logrus.Logger {
	var mode logrus.Level
	if c.LogLevel == config.LevelDebug {
		mode = logrus.DebugLevel
	} else {
		mode = logrus.InfoLevel
	}

	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     mode,
	}
}
