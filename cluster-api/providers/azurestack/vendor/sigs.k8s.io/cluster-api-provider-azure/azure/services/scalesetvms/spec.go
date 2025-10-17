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

package scalesetvms

import (
	"context"
)

// ScaleSetVMSpec defines the specification for a VMSS VM.
type ScaleSetVMSpec struct {
	Name          string
	InstanceID    string
	ResourceGroup string
	ScaleSetName  string
	ProviderID    string
	ResourceID    string
	IsFlex        bool
}

// ResourceName returns the instance ID of the VMSS VM. This is because the it is identified by the instance ID in Azure instead of the name.
func (s *ScaleSetVMSpec) ResourceName() string {
	return s.InstanceID
}

// ResourceGroupName returns the name of the resource group the VMSS that owns this VM.
func (s *ScaleSetVMSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName returns the name of the VMSS that owns this VM.
func (s *ScaleSetVMSpec) OwnerResourceName() string {
	return s.ScaleSetName
}

// Parameters is a no-op for VMSS VMs as this spec is only used to Get().
func (s *ScaleSetVMSpec) Parameters(_ context.Context, _ interface{}) (params interface{}, err error) {
	return nil, nil
}

// VMSSFlexGetter defines the specification for a VMSS flex VM.
type VMSSFlexGetter struct {
	Name          string
	ResourceGroup string
}

// ResourceName returns the name of the flex VM.
func (s *VMSSFlexGetter) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the flex VM.
func (s *VMSSFlexGetter) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for flex VMs.
func (s *VMSSFlexGetter) OwnerResourceName() string {
	return ""
}

// Parameters is a no-op for flex VMs as this spec is only used to Get().
func (s *VMSSFlexGetter) Parameters(_ context.Context, _ interface{}) (params interface{}, err error) {
	return nil, nil
}
