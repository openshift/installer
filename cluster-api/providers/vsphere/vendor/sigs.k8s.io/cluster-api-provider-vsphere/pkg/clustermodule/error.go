/*
Copyright 2022 The Kubernetes Authors.

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

package clustermodule

import (
	"fmt"
	"strings"
)

const errString = "not a compute cluster"

// IncompatibleOwnerError represents the error condition wherein the resource pool in use
// during cluster module creation is not owned by compute cluster resource but owned by
// a standalone host.
type IncompatibleOwnerError struct {
	resource string
}

func (e IncompatibleOwnerError) Error() string {
	return fmt.Sprintf("%s %s", e.resource, errString)
}

// NewIncompatibleOwnerError creates an instance of the IncompatibleOwnerError struct.
// This is introduced for testing purposes only.
func NewIncompatibleOwnerError(name string) IncompatibleOwnerError {
	return IncompatibleOwnerError{resource: name}
}

// IsIncompatibleOwnerError checks if the passed error is an IncompatibleOwnerError.
func IsIncompatibleOwnerError(err error) bool {
	return strings.HasSuffix(err.Error(), errString)
}
