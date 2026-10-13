package powervs

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	suppressDuration = time.Minute * 5
)

// errorTracker holds a history of errors.
type errorTracker struct {
	history map[string]time.Time
}

// suppressWarning logs errors WARN once every duration and the rest to DEBUG.
func (o *errorTracker) suppressWarning(identifier string, err error, logger logrus.FieldLogger) {
	if o.history == nil {
		o.history = map[string]time.Time{}
	}
	if firstSeen, ok := o.history[identifier]; ok {
		if time.Since(firstSeen) > suppressDuration {
			logger.Warn(err)
			o.history[identifier] = time.Now() // reset the clock
		} else {
			logger.Debug(err)
		}
	} else { // first error for this identifier
		o.history[identifier] = time.Now()
		logger.Debug(err)
	}
}
