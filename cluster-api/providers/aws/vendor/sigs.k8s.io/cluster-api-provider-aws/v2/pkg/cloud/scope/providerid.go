/*
Copyright 2018 The Kubernetes Authors.

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

package scope

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Copied from https://github.com/kubernetes-sigs/cluster-api/blob/bda002f52575eeaff68da1ba33c8ef27d5b1014c/controllers/noderefutil/providerid.go
// As this is removed by https://github.com/kubernetes-sigs/cluster-api/pull/9136
var (
	// ErrEmptyProviderID means that the provider id is empty.
	//
	// Deprecated: This var is going to be removed in a future release.
	ErrEmptyProviderID = errors.New("providerID is empty")

	// ErrInvalidProviderID means that the provider id has an invalid form.
	//
	// Deprecated: This var is going to be removed in a future release.
	ErrInvalidProviderID = errors.New("providerID must be of the form <cloudProvider>://<optional>/<segments>/<provider id>")
)

// ProviderID is a struct representation of a Kubernetes ProviderID.
// Format: cloudProvider://optional/segments/etc/id
type ProviderID struct {
	original      string
	cloudProvider string
	id            string
}

/*
- must start with at least one non-colon
- followed by ://
- followed by any number of characters
- must end with a non-slash.
*/
var providerIDRegex = regexp.MustCompile("^[^:]+://.*[^/]$")

// NewProviderID parses the input string and returns a new ProviderID.
func NewProviderID(id string) (*ProviderID, error) {
	if id == "" {
		return nil, ErrEmptyProviderID
	}

	if !providerIDRegex.MatchString(id) {
		return nil, ErrInvalidProviderID
	}

	colonIndex := strings.Index(id, ":")
	cloudProvider := id[0:colonIndex]

	lastSlashIndex := strings.LastIndex(id, "/")
	instance := id[lastSlashIndex+1:]

	res := &ProviderID{
		original:      id,
		cloudProvider: cloudProvider,
		id:            instance,
	}

	if !res.Validate() {
		return nil, ErrInvalidProviderID
	}

	return res, nil
}

// CloudProvider returns the cloud provider portion of the ProviderID.
//
// Deprecated: This method is going to be removed in a future release.
func (p *ProviderID) CloudProvider() string {
	return p.cloudProvider
}

// ID returns the identifier portion of the ProviderID.
//
// Deprecated: This method is going to be removed in a future release.
func (p *ProviderID) ID() string {
	return p.id
}

// Equals returns true if this ProviderID string matches another ProviderID string.
//
// Deprecated: This method is going to be removed in a future release.
func (p *ProviderID) Equals(o *ProviderID) bool {
	return p.String() == o.String()
}

// String returns the string representation of this object.
//
// Deprecated: This method is going to be removed in a future release.
func (p ProviderID) String() string {
	return p.original
}

// Validate returns true if the provider id is valid.
//
// Deprecated: This method is going to be removed in a future release.
func (p *ProviderID) Validate() bool {
	return p.CloudProvider() != "" && p.ID() != ""
}

// IndexKey returns the required level of uniqueness
// to represent and index machines uniquely from their node providerID.
//
// Deprecated: This method is going to be removed in a future release.
func (p *ProviderID) IndexKey() string {
	return p.String()
}

// ProviderIDPrefix is the prefix of AWS resource IDs to form the Kubernetes Provider ID.
// NOTE: this format matches the 2 slashes format used in cloud-provider and cluster-autoscaler.
const ProviderIDPrefix = "aws://"

// GenerateProviderID generates a valid AWS Node/Machine ProviderID field.
//
// By default, the last id provided is used as identifier (last part).
func GenerateProviderID(ids ...string) string {
	return fmt.Sprintf("%s/%s", ProviderIDPrefix, strings.Join(ids, "/"))
}
