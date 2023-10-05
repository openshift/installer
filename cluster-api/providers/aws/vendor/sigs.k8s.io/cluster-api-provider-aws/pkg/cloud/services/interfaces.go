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

package services

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

const (
	// TemporaryResourceID is the name used temporarily when creating AWS resources.
	TemporaryResourceID = "temporary-resource-id"
	// AnyIPv4CidrBlock is the CIDR block to match all IPv4 addresses.
	AnyIPv4CidrBlock = "0.0.0.0/0"
)

// ASGInterface encapsulates the methods exposed to the machinepool
// actuator.
type ASGInterface interface {
	ASGIfExists(id *string) (*expinfrav1.AutoScalingGroup, error)
	GetASGByName(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error)
	CreateASG(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error)
	UpdateASG(scope *scope.MachinePoolScope) error
	StartASGInstanceRefresh(scope *scope.MachinePoolScope) error
	CanStartASGInstanceRefresh(scope *scope.MachinePoolScope) (bool, error)
	UpdateResourceTags(resourceID *string, create, remove map[string]string) error
	DeleteASGAndWait(id string) error
}

// EC2Interface encapsulates the methods exposed to the machine
// actuator.
type EC2Interface interface {
	InstanceIfExists(id *string) (*infrav1.Instance, error)
	TerminateInstance(id string) error
	CreateInstance(scope *scope.MachineScope, userData []byte, userDataFormat string) (*infrav1.Instance, error)
	GetRunningInstanceByTags(scope *scope.MachineScope) (*infrav1.Instance, error)

	GetAdditionalSecurityGroupsIDs(securityGroup []infrav1.AWSResourceReference) ([]string, error)
	GetCoreSecurityGroups(machine *scope.MachineScope) ([]string, error)
	GetInstanceSecurityGroups(instanceID string) (map[string][]string, error)
	UpdateInstanceSecurityGroups(id string, securityGroups []string) error
	UpdateResourceTags(resourceID *string, create, remove map[string]string) error

	TerminateInstanceAndWait(instanceID string) error
	DetachSecurityGroupsFromNetworkInterface(groups []string, interfaceID string) error

	DiscoverLaunchTemplateAMI(scope *scope.MachinePoolScope) (*string, error)
	GetLaunchTemplate(id string) (lt *expinfrav1.AWSLaunchTemplate, userDataHash string, err error)
	GetLaunchTemplateID(id string) (string, error)
	CreateLaunchTemplate(scope *scope.MachinePoolScope, imageID *string, userData []byte) (string, error)
	CreateLaunchTemplateVersion(scope *scope.MachinePoolScope, imageID *string, userData []byte) error
	PruneLaunchTemplateVersions(id string) error
	DeleteLaunchTemplate(id string) error
	LaunchTemplateNeedsUpdate(scope *scope.MachinePoolScope, incoming *expinfrav1.AWSLaunchTemplate, existing *expinfrav1.AWSLaunchTemplate) (bool, error)
	DeleteBastion() error
	ReconcileBastion() error
}

// SecretInterface encapsulated the methods exposed to the
// machine actuator.
type SecretInterface interface {
	Delete(m *scope.MachineScope) error
	Create(m *scope.MachineScope, data []byte) (string, int32, error)
	UserData(secretPrefix string, chunks int32, region string, endpoints []scope.ServiceEndpoint) ([]byte, error)
}

// ELBInterface encapsulates the methods exposed to the cluster and machine
// controller.
type ELBInterface interface {
	DeleteLoadbalancers() error
	ReconcileLoadbalancers() error
	IsInstanceRegisteredWithAPIServerELB(i *infrav1.Instance) (bool, error)
	DeregisterInstanceFromAPIServerELB(i *infrav1.Instance) error
	RegisterInstanceWithAPIServerELB(i *infrav1.Instance) error
}

// NetworkInterface encapsulates the methods exposed to the cluster
// controller.
type NetworkInterface interface {
	DeleteNetwork() error
	ReconcileNetwork() error
}

// SecurityGroupInterface encapsulates the methods exposed to the cluster
// controller.
type SecurityGroupInterface interface {
	DeleteSecurityGroups() error
	ReconcileSecurityGroups() error
}

// ObjectStoreInterface encapsulates the methods exposed to the machine actuator.
type ObjectStoreInterface interface {
	DeleteBucket() error
	ReconcileBucket() error
	Delete(m *scope.MachineScope) error
	Create(m *scope.MachineScope, data []byte) (objectURL string, err error)
}
