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

package util

import (
	"strconv"
	"strings"
)

// LessThan compares the integer values of the supplied VMX versions
// and returns whether the first version is less than the second.
// It returns an error if an invalid vmx version is passed.
func LessThan(version1, version2 string) (bool, error) {
	v1, err := getIntVersion(version1)
	if err != nil {
		return false, err
	}

	v2, err := getIntVersion(version2)
	if err != nil {
		return false, err
	}

	return v1 < v2, nil
}

func getIntVersion(version string) (int, error) {
	versionStr := strings.TrimPrefix(version, "vmx-")
	return strconv.Atoi(versionStr)
}
