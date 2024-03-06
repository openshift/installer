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
	"net/http"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/iam"
)

// EKSAPI defines the EKS API interface.
type EKSAPI interface {
	eksiface.EKSAPI
	WaitUntilClusterUpdating(input *eks.DescribeClusterInput, opts ...request.WaiterOption) error
}

// EKSClient defines a wrapper over EKS API.
type EKSClient struct {
	eksiface.EKSAPI
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
		EKSClient: EKSClient{
			EKSAPI: scope.NewEKSClient(controlPlaneScope, controlPlaneScope, controlPlaneScope, controlPlaneScope.ControlPlane),
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
	AutoscalingClient autoscalingiface.AutoScalingAPI
	EKSClient         eksiface.EKSAPI
	iam.IAMService
	STSClient stsiface.STSAPI
}

// NewNodegroupService returns a new service given the api clients.
func NewNodegroupService(machinePoolScope *scope.ManagedMachinePoolScope) *NodegroupService {
	return &NodegroupService{
		scope:             machinePoolScope,
		AutoscalingClient: scope.NewASGClient(machinePoolScope, machinePoolScope, machinePoolScope, machinePoolScope.ManagedMachinePool),
		EKSClient:         scope.NewEKSClient(machinePoolScope, machinePoolScope, machinePoolScope, machinePoolScope.ManagedMachinePool),
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
	EKSClient eksiface.EKSAPI
	iam.IAMService
	STSClient stsiface.STSAPI
}

// NewFargateService returns a new service given the api clients.
func NewFargateService(fargatePoolScope *scope.FargateProfileScope) *FargateService {
	return &FargateService{
		scope:     fargatePoolScope,
		EKSClient: scope.NewEKSClient(fargatePoolScope, fargatePoolScope, fargatePoolScope, fargatePoolScope.FargateProfile),
		IAMService: iam.IAMService{
			Wrapper:   &fargatePoolScope.Logger,
			IAMClient: scope.NewIAMClient(fargatePoolScope, fargatePoolScope, fargatePoolScope, fargatePoolScope.FargateProfile),
		},
		STSClient: scope.NewSTSClient(fargatePoolScope, fargatePoolScope, fargatePoolScope, fargatePoolScope.FargateProfile),
	}
}
