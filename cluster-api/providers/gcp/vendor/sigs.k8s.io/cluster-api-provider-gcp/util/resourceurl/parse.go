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

// Package resourceurl implements resource url parsing utilities.
package resourceurl

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	// ResourcePrefix is the prefix of a resource URL.
	ResourcePrefix = "https://www.googleapis.com/"
	// NumParts is the number of parts in a resource URL.
	NumParts = 8
	// ResourceCategoryIndex is the index of the resource category in resource URL parts.
	ResourceCategoryIndex = 0
	// ProjectIndex is the index of the project in resource URL parts.
	ProjectIndex = 3
	// LocationIndex is the index of the location in resource URL parts.
	LocationIndex = 5
	// SubResourceIndex is the index of the sub resource in resource URL parts.
	SubResourceIndex = 6
	// NameIndex is the index of the name in resource URL parts.
	NameIndex = 7
)

// ResourceURL captures the individual fields of a GCP resource URL.
// An example of GCP resource URL:
// https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-b/instanceGroupManagers/gke-capg-gke-demo-mypool-aa1282e0-grp
type ResourceURL struct {
	// The resource category (e.g. compute)
	ResourceCategory string
	// The project where the resource lives in (e.g. my-project)
	Project string
	// The location where the resource lives in (e.g. us-central1-b)
	Location string
	// The sub-type of the resource (e.g. instanceGroupManagers)
	SubResource string
	// The name of the resource (e.g. gke-capg-gke-demo-mypool-aa1282e0-grp)
	Name string
}

// Parse parses a resource url.
func Parse(url string) (ResourceURL, error) {
	if !strings.HasPrefix(url, ResourcePrefix) {
		return ResourceURL{}, errors.New("invalid resource url")
	}
	parts := strings.Split(url[len(ResourcePrefix):], "/")
	if len(parts) != NumParts {
		return ResourceURL{}, errors.New("invalid resource url")
	}
	return ResourceURL{
		ResourceCategory: parts[ResourceCategoryIndex],
		Project:          parts[ProjectIndex],
		Location:         parts[LocationIndex],
		SubResource:      parts[SubResourceIndex],
		Name:             parts[NameIndex],
	}, nil
}
