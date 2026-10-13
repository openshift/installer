/*
Copyright 2023 The Kubernetes Authors.

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

package clusters

import (
	"github.com/pkg/errors"
)

// ErrAutopilotClusterMachinePoolsNotAllowed is used when there are machine pools specified for an autopilot enabled cluster.
var ErrAutopilotClusterMachinePoolsNotAllowed = errors.New("cannot use machine pools with an autopilot enabled cluster")

// NewErrUnexpectedClusterStatus creates a new error for an unexpected cluster status.
func NewErrUnexpectedClusterStatus(status string) error {
	return &UnexpectedClusterStatusError{status}
}

// UnexpectedClusterStatusError is the error struct
type UnexpectedClusterStatusError struct {
	status string
}

func (e *UnexpectedClusterStatusError) Error() string {
	return "unexpected error status: " + e.status
}
