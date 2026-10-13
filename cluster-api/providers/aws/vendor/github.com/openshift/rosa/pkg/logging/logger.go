/*
Copyright (c) 2020 Red Hat, Inc.

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

// This file contains the code to build the default loggers used by the project.

package logging

import (
	"github.com/sirupsen/logrus"

	"github.com/openshift/rosa/pkg/debug"
)

// NewLogger creates a new logger with the default config for the project
func NewLogger() (result *logrus.Logger) {
	// Create the logger:
	result = logrus.New()
	result.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		DisableQuote:  true,
		FullTimestamp: true,
	})

	// Enable the debug level if needed:
	if debug.Enabled() {
		result.SetLevel(logrus.DebugLevel)
	}

	return
}
