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

package openstack

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

// GetSupportedMicroversions returns the minimum and maximum microversion that is supported by the ServiceClient Endpoint.
func GetSupportedMicroversions(client gophercloud.ServiceClient) (string, string, error) {
	type valueResp struct {
		ID         string `json:"id"`
		Status     string `json:"status"`
		Version    string `json:"version"`
		MinVersion string `json:"min_version"`
	}

	type response struct {
		Version valueResp `json:"version"`
	}
	var resp response
	_, err := client.Request(context.TODO(), "GET", client.Endpoint, &gophercloud.RequestOpts{
		JSONResponse: &resp,
		OkCodes:      []int{200, 300},
	})
	if err != nil {
		return "", "", err
	}

	return resp.Version.MinVersion, resp.Version.Version, nil
}

// MicroversionSupported checks if a microversion falls in the supported interval.
// It returns true if the version is within the interval and false otherwise.
func MicroversionSupported(version string, minVersion string, maxVersion string) (bool, error) {
	// Parse the version X.Y into X and Y integers that are easier to compare.
	vMajor, v, err := parseMicroversion(version)
	if err != nil {
		return false, err
	}
	minMajor, minimum, err := parseMicroversion(minVersion)
	if err != nil {
		return false, err
	}
	maxMajor, maximum, err := parseMicroversion(maxVersion)
	if err != nil {
		return false, err
	}

	// Check that the major version number is supported.
	if (vMajor < minMajor) || (vMajor > maxMajor) {
		return false, err
	}

	// Check that the minor version number is supported
	if (v <= maximum) && (v >= minimum) {
		return true, nil
	}

	return false, nil
}

// parseMicroversion parses the version X.Y into separate integers X and Y.
// For example, "2.53" becomes 2 and 53.
func parseMicroversion(version string) (int, int, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid microversion format: %q", version)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return major, minor, nil
}
