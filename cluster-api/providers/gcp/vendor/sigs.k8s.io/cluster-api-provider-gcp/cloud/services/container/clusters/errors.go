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
	"fmt"

	"github.com/pkg/errors"
)

var (
	// ErrAutopilotClusterMachinePoolsNotAllowed is used when there are machine pools specified for an autopilot enabled cluster.
	ErrAutopilotClusterMachinePoolsNotAllowed = errors.New("cannot use machine pools with an autopilot enabled cluster")
)

// NewErrUnexpectedClusterStatus creates a new error for an unexpected cluster status.
func NewErrUnexpectedClusterStatus(status string) error {
	return &errUnexpectedClusterStatus{status}
}

type errUnexpectedClusterStatus struct {
	status string
}

func (e *errUnexpectedClusterStatus) Error() string {
	return fmt.Sprintf("unexpected error status: %s", e.status)
}
