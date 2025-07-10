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

// Package services contains the interfaces for the AWS services.
package services

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

const (
	// TemporaryResourceID is the name used temporarily when creating AWS resources.
	TemporaryResourceID = "temporary-resource-id"
	// AnyIPv4CidrBlock is the CIDR block to match all IPv4 addresses.
	AnyIPv4CidrBlock = "0.0.0.0/0"
	// AnyIPv6CidrBlock is the CIDR block to match all IPv6 addresses.
	AnyIPv6CidrBlock = "::/0"
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
	SuspendProcesses(name string, processes []string) error
	ResumeProcesses(name string, processes []string) error
	SubnetIDs(scope *scope.MachinePoolScope) ([]string, error)
	DescribeLifecycleHooks(asgName string) ([]*expinfrav1.AWSLifecycleHook, error)
	CreateLifecycleHook(ctx context.Context, asgName string, hook *expinfrav1.AWSLifecycleHook) error
	UpdateLifecycleHook(ctx context.Context, asgName string, hook *expinfrav1.AWSLifecycleHook) error
	DeleteLifecycleHook(ctx context.Context, asgName string, hook *expinfrav1.AWSLifecycleHook) error
}

// EC2Interface encapsulates the methods exposed to the machine
// actuator.
type EC2Interface interface {
	InstanceIfExists(id *string) (*infrav1.Instance, error)
	TerminateInstance(id string) error
	CreateInstance(ctx context.Context, scope *scope.MachineScope, userData []byte, userDataFormat string) (*infrav1.Instance, error)
	GetRunningInstanceByTags(scope *scope.MachineScope) (*infrav1.Instance, error)

	GetAdditionalSecurityGroupsIDs(securityGroup []infrav1.AWSResourceReference) ([]string, error)
	GetCoreSecurityGroups(machine *scope.MachineScope) ([]string, error)
	GetInstanceSecurityGroups(instanceID string) (map[string][]string, error)
	UpdateInstanceSecurityGroups(id string, securityGroups []string) error
	UpdateResourceTags(resourceID *string, create, remove map[string]string) error
	ModifyInstanceMetadataOptions(instanceID string, options *infrav1.InstanceMetadataOptions) error

	TerminateInstanceAndWait(instanceID string) error
	DetachSecurityGroupsFromNetworkInterface(groups []string, interfaceID string) error

	DiscoverLaunchTemplateAMI(ctx context.Context, scope scope.LaunchTemplateScope) (*string, error)
	GetLaunchTemplate(id string) (lt *expinfrav1.AWSLaunchTemplate, userDataHash string, userDataSecretKey *apimachinerytypes.NamespacedName, bootstrapDataHash *string, err error)
	GetLaunchTemplateID(id string) (string, error)
	GetLaunchTemplateLatestVersion(id string) (string, error)
	CreateLaunchTemplate(scope scope.LaunchTemplateScope, imageID *string, userDataSecretKey apimachinerytypes.NamespacedName, userData []byte, bootstrapDataHash string) (string, error)
	CreateLaunchTemplateVersion(id string, scope scope.LaunchTemplateScope, imageID *string, userDataSecretKey apimachinerytypes.NamespacedName, userData []byte, bootstrapDataHash string) error
	PruneLaunchTemplateVersions(id string) (*ec2.LaunchTemplateVersion, error)
	DeleteLaunchTemplate(id string) error
	LaunchTemplateNeedsUpdate(scope scope.LaunchTemplateScope, incoming *expinfrav1.AWSLaunchTemplate, existing *expinfrav1.AWSLaunchTemplate) (bool, error)
	DeleteBastion() error
	ReconcileBastion() error
	// ReconcileElasticIPFromPublicPool reconciles the elastic IP from a custom Public IPv4 Pool.
	ReconcileElasticIPFromPublicPool(pool *infrav1.ElasticIPPool, instance *infrav1.Instance) (bool, error)

	// ReleaseElasticIP reconciles the elastic IP from a custom Public IPv4 Pool.
	ReleaseElasticIP(instanceID string) error
}

// MachinePoolReconcileInterface encapsulates high-level reconciliation functions regarding EC2 reconciliation. It is
// separate from EC2Interface so that we can mock AWS requests separately. For example, by not mocking the
// ReconcileLaunchTemplate function, but mocking EC2Interface, we can test which EC2 API operations would have been called.
type MachinePoolReconcileInterface interface {
	ReconcileLaunchTemplate(ctx context.Context, ignitionScope scope.IgnitionScope, scope scope.LaunchTemplateScope, s3Scope scope.S3Scope, ec2svc EC2Interface, objectStoreSvc ObjectStoreInterface, canUpdateLaunchTemplate func() (bool, error), runPostLaunchTemplateUpdateOperation func() error) error
	ReconcileTags(scope scope.LaunchTemplateScope, resourceServicesToUpdate []scope.ResourceServiceToUpdate) error
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
	DeleteLoadbalancers(ctx context.Context) error
	ReconcileLoadbalancers(ctx context.Context) error
	IsInstanceRegisteredWithAPIServerELB(ctx context.Context, i *infrav1.Instance) (bool, error)
	IsInstanceRegisteredWithAPIServerLB(ctx context.Context, i *infrav1.Instance, lb *infrav1.AWSLoadBalancerSpec) ([]string, bool, error)
	DeregisterInstanceFromAPIServerELB(ctx context.Context, i *infrav1.Instance) error
	DeregisterInstanceFromAPIServerLB(ctx context.Context, targetGroupArn string, i *infrav1.Instance) error
	RegisterInstanceWithAPIServerELB(ctx context.Context, i *infrav1.Instance) error
	RegisterInstanceWithAPIServerLB(ctx context.Context, i *infrav1.Instance, lb *infrav1.AWSLoadBalancerSpec) error
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
	DeleteBucket(ctx context.Context) error
	ReconcileBucket(ctx context.Context) error
	Delete(ctx context.Context, m *scope.MachineScope) error
	Create(ctx context.Context, m *scope.MachineScope, data []byte) (objectURL string, err error)
	CreateForMachinePool(ctx context.Context, scope scope.LaunchTemplateScope, data []byte) (objectURL string, err error)
	DeleteForMachinePool(ctx context.Context, scope scope.LaunchTemplateScope, bootstrapDataHash string) error
}

// AWSNodeInterface installs the CNI for EKS clusters.
type AWSNodeInterface interface {
	ReconcileCNI(ctx context.Context) error
}

// IAMAuthenticatorInterface installs aws-iam-authenticator for EKS clusters.
type IAMAuthenticatorInterface interface {
	ReconcileIAMAuthenticator(ctx context.Context) error
}

// KubeProxyInterface installs kube-proxy for EKS clusters.
type KubeProxyInterface interface {
	ReconcileKubeProxy(ctx context.Context) error
}
