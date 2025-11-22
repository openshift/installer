/*
Copyright 2020 The ORC Authors.

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

package logging

// These levels are defined to categorize what a given log level means.
// Logs with levels 1-3 are for operators. Log level 4 is for developers.
const (
	// Status logs are always shown. Status logs should be reserved for operational
	// logs about the service itself, e.g.:
	// - startup and shutdown messages
	// - runtime conditions which may indicate something about the state of the service,
	//    e.g. inability to reach kube-apiserver.
	Status = 1
	// Info is the default log level for most deployments. It should log the principal actions
	// of the service, i.e. resource creation and deletion, and 'reconcile complete'
	// (i.e. Progressing=False) messages for success and failure. It should not include actions
	// which happen on every reconcile.
	// Example: "OpenStack resource created"
	Info = 2
	// Verbose logs provide additional context for an administrator trying to understand why an action
	// be occurring or not occurring. It should produce logs on every reconcile attempt.
	// Example: "web-download is not supported because..."
	Verbose = 3
	// Debug logs are very verbose. They should include things that should
	// help with debugging/development.
	// Example: "Got resource"
	Debug = 4
)
