package log

import (
	logger "github.com/sirupsen/logrus"
)

func Initlogger() {
	customFormatter := new(logger.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logger.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}
func LogInfo(format string, args ...interface{}) {
	Initlogger()
	logger.Infof(format, args...)
}

func LogError(format string, args ...interface{}) {
	Initlogger()
	logger.Errorf(format, args...)
}

func LogDebug(format string, args ...interface{}) {
	Initlogger()
	logger.Debugf(format, args...)
}

func LogWarning(format string, args ...interface{}) {
	Initlogger()
	logger.Warnf(format, args...)
}
