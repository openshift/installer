package terraform

import (
	"github.com/sirupsen/logrus"
)

type printfer struct {
	logger *logrus.Logger
	level  logrus.Level
}

func (t *printfer) Printf(format string, v ...interface{}) {
	t.logger.Logf(t.level, format, v...)
}

func newPrintfer() *printfer {
	return &printfer{
		logger: logrus.StandardLogger(),
		level:  logrus.DebugLevel,
	}
}
