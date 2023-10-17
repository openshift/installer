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

package providerid

import (
	"errors"
	"fmt"
	"path"

	"sigs.k8s.io/cluster-api-provider-gcp/util/resourceurl"
)

const (
	// Prefix is the gce provider id prefix.
	Prefix = "gce://"
)

// ProviderID represents the id for a GCP cluster.
type ProviderID interface {
	Project() string
	Location() string
	Name() string
	fmt.Stringer
}

// NewFromResourceURL creates a provider from a GCP resource url.
func NewFromResourceURL(url string) (ProviderID, error) {
	resourceURL, err := resourceurl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("parsing resource url %s: %w", url, err)
	}

	return New(resourceURL.Project, resourceURL.Location, resourceURL.Name)
}

// New creates a new provider id.
func New(project, location, name string) (ProviderID, error) {
	if project == "" {
		return nil, errors.New("project required for provider id")
	}
	if location == "" {
		return nil, errors.New("location required for provider id")
	}
	if name == "" {
		return nil, errors.New("name required for provider id")
	}

	return &providerID{
		project:  project,
		location: location,
		name:     name,
	}, nil
}

type providerID struct {
	project  string
	location string
	name     string
}

func (p *providerID) Project() string {
	return p.project
}

func (p *providerID) Location() string {
	return p.location
}

func (p *providerID) Name() string {
	return p.name
}

func (p *providerID) String() string {
	return Prefix + path.Join(p.project, p.location, p.name)
}
