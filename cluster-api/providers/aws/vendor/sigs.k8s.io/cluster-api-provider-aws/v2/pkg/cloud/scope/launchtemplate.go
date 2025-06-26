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

package scope

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// LaunchTemplateScope defines a scope defined around a launch template.
type LaunchTemplateScope interface {
	GetMachinePool() *expclusterv1.MachinePool
	GetLaunchTemplate() *expinfrav1.AWSLaunchTemplate
	LaunchTemplateName() string
	GetLaunchTemplateIDStatus() string
	SetLaunchTemplateIDStatus(id string)
	GetLaunchTemplateLatestVersionStatus() string
	SetLaunchTemplateLatestVersionStatus(version string)
	GetRawBootstrapData() ([]byte, string, *types.NamespacedName, error)

	IsEKSManaged() bool
	AdditionalTags() infrav1.Tags

	GetObjectMeta() *metav1.ObjectMeta
	GetSetter() conditions.Setter
	PatchObject() error
	GetEC2Scope() EC2Scope

	client.Client
	logger.Wrapper
}

// ResourceServiceToUpdate is a struct that contains the resource ID and the resource service to update.
type ResourceServiceToUpdate struct {
	ResourceID      *string
	ResourceService ResourceService
}

// ResourceService defines the interface for resources.
type ResourceService interface {
	UpdateResourceTags(resourceID *string, create, remove map[string]string) error
}
