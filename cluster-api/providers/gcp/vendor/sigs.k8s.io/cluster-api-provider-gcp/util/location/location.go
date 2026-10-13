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

// Package location implements location parsing utilities.
package location

import (
	"strings"

	"github.com/pkg/errors"
)

// Location captures the region and zone of a GCP location.
// Examples of GCP location:
// us-central1 (region).
// us-central1-c (region with zone).
type Location struct {
	Region string
	Zone   *string
}

// Parse parses a location string.
func Parse(location string) (Location, error) {
	parts := strings.Split(location, "-")
	if len(parts) < 2 {
		return Location{}, errors.New("invalid location")
	}
	region := strings.Join(parts[:2], "-")
	var zone *string
	if len(parts) == 3 {
		zone = &parts[2]
	}
	return Location{
		Region: region,
		Zone:   zone,
	}, nil
}
