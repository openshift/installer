/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/iam"
)

// EKSAPI defines the EKS API interface.
type EKSAPI interface {
	CreateCluster(ctx context.Context, params *eks.CreateClusterInput, optFns ...func(*eks.Options)) (*eks.CreateClusterOutput, error)
	DeleteCluster(ctx context.Context, params *eks.DeleteClusterInput, optFns ...func(*eks.Options)) (*eks.DeleteClusterOutput, error)
	DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error)
	UpdateClusterConfig(ctx context.Context, params *eks.UpdateClusterConfigInput, optFns ...func(*eks.Options)) (*eks.UpdateClusterConfigOutput, error)
	UpdateClusterVersion(ctx context.Context, params *eks.UpdateClusterVersionInput, optFns ...func(*eks.Options)) (*eks.UpdateClusterVersionOutput, error)
	DescribeUpdate(ctx context.Context, params *eks.DescribeUpdateInput, optFns ...func(*eks.Options)) (*eks.DescribeUpdateOutput, error)
	AssociateEncryptionConfig(ctx context.Context, params *eks.AssociateEncryptionConfigInput, optFns ...func(*eks.Options)) (*eks.AssociateEncryptionConfigOutput, error)
	ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error)
	CreateNodegroup(ctx context.Context, params *eks.CreateNodegroupInput, optFns ...func(*eks.Options)) (*eks.CreateNodegroupOutput, error)
	DeleteNodegroup(ctx context.Context, params *eks.DeleteNodegroupInput, optFns ...func(*eks.Options)) (*eks.DeleteNodegroupOutput, error)
	DescribeNodegroup(ctx context.Context, params *eks.DescribeNodegroupInput, optFns ...func(*eks.Options)) (*eks.DescribeNodegroupOutput, error)
	UpdateNodegroupConfig(ctx context.Context, params *eks.UpdateNodegroupConfigInput, optFns ...func(*eks.Options)) (*eks.UpdateNodegroupConfigOutput, error)
	UpdateNodegroupVersion(ctx context.Context, params *eks.UpdateNodegroupVersionInput, optFns ...func(*eks.Options)) (*eks.UpdateNodegroupVersionOutput, error)
	DescribeAddon(ctx context.Context, params *eks.DescribeAddonInput, optFns ...func(*eks.Options)) (*eks.DescribeAddonOutput, error)
	CreateAddon(ctx context.Context, params *eks.CreateAddonInput, optFns ...func(*eks.Options)) (*eks.CreateAddonOutput, error)
	UpdateAddon(ctx context.Context, params *eks.UpdateAddonInput, optFns ...func(*eks.Options)) (*eks.UpdateAddonOutput, error)
	ListAddons(ctx context.Context, params *eks.ListAddonsInput, optFns ...func(*eks.Options)) (*eks.ListAddonsOutput, error)
	DescribeAddonConfiguration(ctx context.Context, params *eks.DescribeAddonConfigurationInput, optFns ...func(*eks.Options)) (*eks.DescribeAddonConfigurationOutput, error)
	DescribeAddonVersions(ctx context.Context, params *eks.DescribeAddonVersionsInput, optFns ...func(*eks.Options)) (*eks.DescribeAddonVersionsOutput, error)
	DeleteAddon(ctx context.Context, params *eks.DeleteAddonInput, optFns ...func(*eks.Options)) (*eks.DeleteAddonOutput, error)
	ListIdentityProviderConfigs(ctx context.Context, params *eks.ListIdentityProviderConfigsInput, optFns ...func(*eks.Options)) (*eks.ListIdentityProviderConfigsOutput, error)
	DescribeIdentityProviderConfig(ctx context.Context, params *eks.DescribeIdentityProviderConfigInput, optFns ...func(*eks.Options)) (*eks.DescribeIdentityProviderConfigOutput, error)
	AssociateIdentityProviderConfig(ctx context.Context, params *eks.AssociateIdentityProviderConfigInput, optFns ...func(*eks.Options)) (*eks.AssociateIdentityProviderConfigOutput, error)
	CreateFargateProfile(ctx context.Context, params *eks.CreateFargateProfileInput, optFns ...func(*eks.Options)) (*eks.CreateFargateProfileOutput, error)
	DeleteFargateProfile(ctx context.Context, params *eks.DeleteFargateProfileInput, optFns ...func(*eks.Options)) (*eks.DeleteFargateProfileOutput, error)
	DescribeFargateProfile(ctx context.Context, params *eks.DescribeFargateProfileInput, optFns ...func(*eks.Options)) (*eks.DescribeFargateProfileOutput, error)
	TagResource(ctx context.Context, params *eks.TagResourceInput, optFns ...func(*eks.Options)) (*eks.TagResourceOutput, error)
	UntagResource(ctx context.Context, params *eks.UntagResourceInput, optFns ...func(*eks.Options)) (*eks.UntagResourceOutput, error)
	DisassociateIdentityProviderConfig(ctx context.Context, params *eks.DisassociateIdentityProviderConfigInput, optFns ...func(*eks.Options)) (*eks.DisassociateIdentityProviderConfigOutput, error)

	// Waiters for EKS Cluster
	WaitUntilClusterActive(ctx context.Context, params *eks.DescribeClusterInput, maxWait time.Duration) error
	WaitUntilClusterDeleted(ctx context.Context, params *eks.DescribeClusterInput, maxWait time.Duration) error
	WaitUntilClusterUpdating(ctx context.Context, params *eks.DescribeClusterInput, maxWait time.Duration) error

	// Waiters for EKS Nodegroup
	WaitUntilNodegroupActive(ctx context.Context, params *eks.DescribeNodegroupInput, maxWait time.Duration) error
	WaitUntilNodegroupDeleted(ctx context.Context, params *eks.DescribeNodegroupInput, maxWait time.Duration) error

	// Waiters for EKS Addon
	WaitUntilAddonDeleted(ctx context.Context, params *eks.DescribeAddonInput, maxWait time.Duration) error
}

// EKSClient is a wrapper over eks.Client for implementing custom methods of EKSAPI.
type EKSClient struct {
	*eks.Client
}

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     *scope.ManagedControlPlaneScope
	EC2Client ec2iface.EC2API
	EKSClient EKSAPI
	iam.IAMService
	STSClient stsiface.STSAPI
}

// ServiceOpts defines the functional arguments for the service.
type ServiceOpts func(s *Service)

// WithIAMClient creates an access spec with a custom http client.
func WithIAMClient(client *http.Client) ServiceOpts {
	return func(s *Service) {
		s.IAMService.Client = client
	}
}

// NewService returns a new service given the api clients.
func NewService(controlPlaneScope *scope.ManagedControlPlaneScope, opts ...ServiceOpts) *Service {
	s := &Service{
		scope:     controlPlaneScope,
		EC2Client: scope.NewEC2Client(controlPlaneScope, controlPlaneScope, controlPlaneScope, controlPlaneScope.ControlPlane),
		EKSClient: &EKSClient{
			Client: scope.NewEKSClient(controlPlaneScope, controlPlaneScope, controlPlaneScope, controlPlaneScope.ControlPlane),
		},
		IAMService: iam.IAMService{
			Wrapper:   &controlPlaneScope.Logger,
			IAMClient: scope.NewIAMClient(controlPlaneScope, controlPlaneScope, controlPlaneScope, controlPlaneScope.ControlPlane),
			Client:    http.DefaultClient,
		},
		STSClient: scope.NewSTSClient(controlPlaneScope, controlPlaneScope, controlPlaneScope, controlPlaneScope.ControlPlane),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// NodegroupService holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type NodegroupService struct {
	scope             *scope.ManagedMachinePoolScope
	ASGService        services.ASGInterface
	AutoscalingClient *autoscaling.Client
	EKSClient         EKSAPI
	iam.IAMService
	STSClient stsiface.STSAPI
}

// NewNodegroupService returns a new service given the api clients.
func NewNodegroupService(machinePoolScope *scope.ManagedMachinePoolScope) *NodegroupService {
	return &NodegroupService{
		scope:             machinePoolScope,
		AutoscalingClient: scope.NewASGClient(machinePoolScope, machinePoolScope, machinePoolScope, machinePoolScope.ManagedMachinePool),
		EKSClient: &EKSClient{
			Client: scope.NewEKSClient(machinePoolScope, machinePoolScope, machinePoolScope, machinePoolScope.ManagedMachinePool),
		},
		IAMService: iam.IAMService{
			Wrapper:   &machinePoolScope.Logger,
			IAMClient: scope.NewIAMClient(machinePoolScope, machinePoolScope, machinePoolScope, machinePoolScope.ManagedMachinePool),
		},
		STSClient: scope.NewSTSClient(machinePoolScope, machinePoolScope, machinePoolScope, machinePoolScope.ManagedMachinePool),
	}
}

// FargateService holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
type FargateService struct {
	scope     *scope.FargateProfileScope
	EKSClient EKSAPI
	iam.IAMService
	STSClient stsiface.STSAPI
}

// NewFargateService returns a new service given the api clients.
func NewFargateService(fargatePoolScope *scope.FargateProfileScope) *FargateService {
	return &FargateService{
		scope: fargatePoolScope,
		EKSClient: &EKSClient{
			Client: scope.NewEKSClient(fargatePoolScope, fargatePoolScope, fargatePoolScope, fargatePoolScope.FargateProfile),
		},
		IAMService: iam.IAMService{
			Wrapper:   &fargatePoolScope.Logger,
			IAMClient: scope.NewIAMClient(fargatePoolScope, fargatePoolScope, fargatePoolScope, fargatePoolScope.FargateProfile),
		},
		STSClient: scope.NewSTSClient(fargatePoolScope, fargatePoolScope, fargatePoolScope, fargatePoolScope.FargateProfile),
	}
}
