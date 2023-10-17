/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tele

import (
	"github.com/go-logr/logr"
)

type compositeLogSink struct {
	logSinks []logr.LogSink
}

func (c *compositeLogSink) Init(info logr.RuntimeInfo) {
	// we change the depth of the stack so that we can get the real
	// line where the log statement was called. We need to do this because the composite logger adds to the
	// call stack due to wrapping the internal logger.
	info.CallDepth += 2
	for _, l := range c.logSinks {
		l.Init(info)
	}
}

func (c *compositeLogSink) Enabled(v int) bool {
	for _, l := range c.logSinks {
		if !l.Enabled(v) {
			return false
		}
	}
	return true
}

func (c *compositeLogSink) iter(fn func(l logr.LogSink)) {
	for _, l := range c.logSinks {
		fn(l)
	}
}

func (c *compositeLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	c.iter(func(l logr.LogSink) {
		l.Info(level, msg, keysAndValues...)
	})
}

func (c *compositeLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	c.iter(func(l logr.LogSink) {
		l.Error(err, msg, keysAndValues...)
	})
}

func (c *compositeLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	var logSinks = make([]logr.LogSink, len(c.logSinks))
	for i, l := range c.logSinks {
		logSinks[i] = l.WithValues(keysAndValues...)
	}

	return &compositeLogSink{
		logSinks: logSinks,
	}
}

func (c *compositeLogSink) WithName(name string) logr.LogSink {
	var logSinks = make([]logr.LogSink, len(c.logSinks))
	for i, l := range c.logSinks {
		logSinks[i] = l.WithName(name)
	}

	return &compositeLogSink{
		logSinks: logSinks,
	}
}

// NewCompositeLogger is the main entry-point to this implementation.
func NewCompositeLogger(logSinks []logr.LogSink) logr.Logger {
	sink := &compositeLogSink{
		logSinks: logSinks,
	}
	return logr.New(sink)
}
