/*
Copyright 2021 The Kubernetes Authors.

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

package disks

import "context"

// DiskSpec defines the specification for a disk.
type DiskSpec struct {
	Name          string
	ResourceGroup string
}

// ResourceName returns the name of the disk.
func (s *DiskSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *DiskSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for disks.
func (s *DiskSpec) OwnerResourceName() string {
	return ""
}

// Parameters is a no-op for disks.
func (s *DiskSpec) Parameters(_ context.Context, _ interface{}) (params interface{}, err error) {
	return nil, nil
}
