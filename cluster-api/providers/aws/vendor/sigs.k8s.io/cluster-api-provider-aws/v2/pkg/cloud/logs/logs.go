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

package logs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-logr/logr"
)

const (
	logWithHTTPHeader = 9
	logWithHTTPBody   = 10
)

// GetAWSLogLevel will return the log level of an AWS Logger.
func GetAWSLogLevel(logger logr.Logger) aws.LogLevelType {
	if logger.V(logWithHTTPBody).Enabled() {
		return aws.LogDebugWithHTTPBody
	}

	if logger.V(logWithHTTPHeader).Enabled() {
		return aws.LogDebug
	}

	return aws.LogOff
}

// NewWrapLogr will create an AWS Logger wrapper.
func NewWrapLogr(logger logr.Logger) aws.Logger {
	return &logrWrapper{
		log: logger,
	}
}

type logrWrapper struct {
	log logr.Logger
}

func (l *logrWrapper) Log(msgs ...interface{}) {
	switch len(msgs) {
	case 0:
		return
	case 1:
		l.log.Info(msgs[0].(string))
	default:
		l.log.Info(msgs[0].(string), msgs[:1]...)
	}
}
