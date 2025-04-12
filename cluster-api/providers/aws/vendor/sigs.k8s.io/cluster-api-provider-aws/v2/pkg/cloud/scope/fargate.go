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

package scope

import (
	"context"

	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/throttle"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

// FargateProfileScopeParams defines the input parameters used to create a new Scope.
type FargateProfileScopeParams struct {
	Client         client.Client
	Logger         *logger.Logger
	Cluster        *clusterv1.Cluster
	ControlPlane   *ekscontrolplanev1.AWSManagedControlPlane
	FargateProfile *expinfrav1.AWSFargateProfile
	ControllerName string
	Endpoints      []ServiceEndpoint
	Session        awsclient.ConfigProvider

	EnableIAM bool
}

// NewFargateProfileScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewFargateProfileScope(params FargateProfileScopeParams) (*FargateProfileScope, error) {
	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil AWSFargateProfile")
	}
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	managedScope := &ManagedControlPlaneScope{
		Logger:         *params.Logger,
		Client:         params.Client,
		Cluster:        params.Cluster,
		ControlPlane:   params.ControlPlane,
		controllerName: params.ControllerName,
	}

	session, serviceLimiters, err := sessionForClusterWithRegion(params.Client, managedScope, params.ControlPlane.Spec.Region, params.Endpoints, params.Logger)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	helper, err := patch.NewHelper(params.FargateProfile, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &FargateProfileScope{
		Logger:          *params.Logger,
		Client:          params.Client,
		Cluster:         params.Cluster,
		ControlPlane:    params.ControlPlane,
		FargateProfile:  params.FargateProfile,
		patchHelper:     helper,
		session:         session,
		serviceLimiters: serviceLimiters,
		controllerName:  params.ControllerName,
		enableIAM:       params.EnableIAM,
	}, nil
}

// FargateProfileScope defines the basic context for an actuator to operate upon.
type FargateProfileScope struct {
	logger.Logger
	Client      client.Client
	patchHelper *patch.Helper

	Cluster        *clusterv1.Cluster
	ControlPlane   *ekscontrolplanev1.AWSManagedControlPlane
	FargateProfile *expinfrav1.AWSFargateProfile

	session         awsclient.ConfigProvider
	serviceLimiters throttle.ServiceLimiters
	controllerName  string

	enableIAM bool
}

// ManagedPoolName returns the managed machine pool name.
func (s *FargateProfileScope) ManagedPoolName() string {
	return s.FargateProfile.Name
}

// ServiceLimiter returns the AWS SDK session. Used for creating clients.
func (s *FargateProfileScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ClusterName returns the cluster name.
func (s *FargateProfileScope) ClusterName() string {
	return s.Cluster.Name
}

// EnableIAM indicates that reconciliation should create IAM roles.
func (s *FargateProfileScope) EnableIAM() bool {
	return s.enableIAM
}

// AdditionalTags returns AdditionalTags from the scope's FargateProfile
// The returned value will never be nil.
func (s *FargateProfileScope) AdditionalTags() infrav1.Tags {
	if s.FargateProfile.Spec.AdditionalTags == nil {
		s.FargateProfile.Spec.AdditionalTags = infrav1.Tags{}
	}

	return s.FargateProfile.Spec.AdditionalTags.DeepCopy()
}

// RoleName returns the node group role name.
func (s *FargateProfileScope) RoleName() string {
	return s.FargateProfile.Spec.RoleName
}

// ControlPlaneSubnets returns the control plane subnets.
func (s *FargateProfileScope) ControlPlaneSubnets() *infrav1.Subnets {
	return &s.ControlPlane.Spec.NetworkSpec.Subnets
}

// SubnetIDs returns the machine pool subnet IDs.
func (s *FargateProfileScope) SubnetIDs() []string {
	return s.FargateProfile.Spec.SubnetIDs
}

// Partition returns the machine pool subnet IDs.
func (s *FargateProfileScope) Partition() string {
	if s.ControlPlane.Spec.Partition == "" {
		s.ControlPlane.Spec.Partition = system.GetPartitionFromRegion(s.ControlPlane.Spec.Region)
	}
	return s.ControlPlane.Spec.Partition
}

// IAMReadyFalse marks the ready condition false using warning if error isn't
// empty.
func (s *FargateProfileScope) IAMReadyFalse(reason string, err string) error {
	severity := clusterv1.ConditionSeverityWarning
	if err == "" {
		severity = clusterv1.ConditionSeverityInfo
	}
	conditions.MarkFalse(
		s.FargateProfile,
		expinfrav1.IAMFargateRolesReadyCondition,
		reason,
		severity,
		"%s",
		err,
	)
	if err := s.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to mark role not ready")
	}
	return nil
}

// PatchObject persists the control plane configuration and status.
func (s *FargateProfileScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.FargateProfile,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			expinfrav1.EKSFargateProfileReadyCondition,
			expinfrav1.EKSFargateCreatingCondition,
			expinfrav1.EKSFargateDeletingCondition,
			expinfrav1.IAMFargateRolesReadyCondition,
		}})
}

// Close closes the current scope persisting the control plane configuration and status.
func (s *FargateProfileScope) Close() error {
	return s.PatchObject()
}

// InfraCluster returns the AWS infrastructure cluster or control plane object.
func (s *FargateProfileScope) InfraCluster() cloud.ClusterObject {
	return s.ControlPlane
}

// ClusterObj returns the cluster object.
func (s *FargateProfileScope) ClusterObj() cloud.ClusterObject {
	return s.Cluster
}

// Session returns the AWS SDK session. Used for creating clients.
func (s *FargateProfileScope) Session() awsclient.ConfigProvider {
	return s.session
}

// ControllerName returns the name of the controller that
// created the FargateProfile.
func (s *FargateProfileScope) ControllerName() string {
	return s.controllerName
}

// KubernetesClusterName is the name of the EKS cluster name.
func (s *FargateProfileScope) KubernetesClusterName() string {
	return s.ControlPlane.Spec.EKSClusterName
}
