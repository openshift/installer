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

package session

import (
	"fmt"
	"strings"
)

const (
	errString = "vCenter version cannot be identified"
)

// TODO (srm09): figure out a common place for custom errors.
type unidentifiedVCenterVersion struct {
	version string
}

func (e unidentifiedVCenterVersion) Error() string {
	return fmt.Sprintf("%s: %s", errString, e.version)
}

// IsUnidentifiedVCenterVersion identies an error as an unidentifiedVCenterVersion error.
func IsUnidentifiedVCenterVersion(err error) bool {
	return strings.HasPrefix(err.Error(), errString)
}
