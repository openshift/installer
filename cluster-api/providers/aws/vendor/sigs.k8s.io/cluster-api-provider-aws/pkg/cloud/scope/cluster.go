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
	"context"
	"fmt"

	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/throttle"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

// ClusterScopeParams defines the input parameters used to create a new Scope.
type ClusterScopeParams struct {
	Client         client.Client
	Logger         *logr.Logger
	Cluster        *clusterv1.Cluster
	AWSCluster     *infrav1.AWSCluster
	ControllerName string
	Endpoints      []ServiceEndpoint
	Session        awsclient.ConfigProvider
}

// NewClusterScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.AWSCluster == nil {
		return nil, errors.New("failed to generate new scope from nil AWSCluster")
	}

	if params.Logger == nil {
		log := klogr.New()
		params.Logger = &log
	}

	clusterScope := &ClusterScope{
		Logger:         *params.Logger,
		client:         params.Client,
		Cluster:        params.Cluster,
		AWSCluster:     params.AWSCluster,
		controllerName: params.ControllerName,
	}

	session, serviceLimiters, err := sessionForClusterWithRegion(params.Client, clusterScope, params.AWSCluster.Spec.Region, params.Endpoints, *params.Logger)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	helper, err := patch.NewHelper(params.AWSCluster, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	clusterScope.patchHelper = helper
	clusterScope.session = session
	clusterScope.serviceLimiters = serviceLimiters

	return clusterScope, nil
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	Cluster    *clusterv1.Cluster
	AWSCluster *infrav1.AWSCluster

	session         awsclient.ConfigProvider
	serviceLimiters throttle.ServiceLimiters
	controllerName  string
}

// Network returns the cluster network object.
func (s *ClusterScope) Network() *infrav1.NetworkStatus {
	return &s.AWSCluster.Status.Network
}

// VPC returns the cluster VPC.
func (s *ClusterScope) VPC() *infrav1.VPCSpec {
	return &s.AWSCluster.Spec.NetworkSpec.VPC
}

// Subnets returns the cluster subnets.
func (s *ClusterScope) Subnets() infrav1.Subnets {
	return s.AWSCluster.Spec.NetworkSpec.Subnets
}

// IdentityRef returns the cluster identityRef.
func (s *ClusterScope) IdentityRef() *infrav1.AWSIdentityReference {
	return s.AWSCluster.Spec.IdentityRef
}

// SetSubnets updates the clusters subnets.
func (s *ClusterScope) SetSubnets(subnets infrav1.Subnets) {
	s.AWSCluster.Spec.NetworkSpec.Subnets = subnets
}

// CNIIngressRules returns the CNI spec ingress rules.
func (s *ClusterScope) CNIIngressRules() infrav1.CNIIngressRules {
	if s.AWSCluster.Spec.NetworkSpec.CNI != nil {
		return s.AWSCluster.Spec.NetworkSpec.CNI.CNIIngressRules
	}
	return infrav1.CNIIngressRules{}
}

// SecurityGroupOverrides returns the cluster security group overrides.
func (s *ClusterScope) SecurityGroupOverrides() map[infrav1.SecurityGroupRole]string {
	return s.AWSCluster.Spec.NetworkSpec.SecurityGroupOverrides
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *ClusterScope) SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup {
	return s.AWSCluster.Status.Network.SecurityGroups
}

// SecondaryCidrBlock is currently unimplemented for non-managed clusters.
func (s *ClusterScope) SecondaryCidrBlock() *string {
	return nil
}

// Name returns the CAPI cluster name.
func (s *ClusterScope) Name() string {
	return s.Cluster.Name
}

// Namespace returns the cluster namespace.
func (s *ClusterScope) Namespace() string {
	return s.Cluster.Namespace
}

// InfraClusterName returns the AWS cluster name.
func (s *ClusterScope) InfraClusterName() string {
	return s.AWSCluster.Name
}

// Region returns the cluster region.
func (s *ClusterScope) Region() string {
	return s.AWSCluster.Spec.Region
}

// KubernetesClusterName is the name of the Kubernetes cluster. For the cluster
// scope this is the same as the CAPI cluster name.
func (s *ClusterScope) KubernetesClusterName() string {
	return s.Cluster.Name
}

// ControlPlaneLoadBalancer returns the AWSLoadBalancerSpec.
func (s *ClusterScope) ControlPlaneLoadBalancer() *infrav1.AWSLoadBalancerSpec {
	return s.AWSCluster.Spec.ControlPlaneLoadBalancer
}

// ControlPlaneLoadBalancerScheme returns the Classic ELB scheme (public or internal facing).
func (s *ClusterScope) ControlPlaneLoadBalancerScheme() infrav1.ClassicELBScheme {
	if s.ControlPlaneLoadBalancer() != nil && s.ControlPlaneLoadBalancer().Scheme != nil {
		return *s.ControlPlaneLoadBalancer().Scheme
	}
	return infrav1.ClassicELBSchemeInternetFacing
}

func (s *ClusterScope) ControlPlaneLoadBalancerName() *string {
	if s.AWSCluster.Spec.ControlPlaneLoadBalancer != nil {
		return s.AWSCluster.Spec.ControlPlaneLoadBalancer.Name
	}
	return nil
}

func (s *ClusterScope) ControlPlaneEndpoint() clusterv1.APIEndpoint {
	return s.AWSCluster.Spec.ControlPlaneEndpoint
}

func (s *ClusterScope) Bucket() *infrav1.S3Bucket {
	return s.AWSCluster.Spec.S3Bucket
}

// ControlPlaneConfigMapName returns the name of the ConfigMap used to
// coordinate the bootstrapping of control plane nodes.
func (s *ClusterScope) ControlPlaneConfigMapName() string {
	return fmt.Sprintf("%s-controlplane", s.Cluster.UID)
}

// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
func (s *ClusterScope) ListOptionsLabelSelector() client.ListOption {
	return client.MatchingLabels(map[string]string{
		clusterv1.ClusterLabelName: s.Cluster.Name,
	})
}

// PatchObject persists the cluster configuration and status.
func (s *ClusterScope) PatchObject() error {
	// Always update the readyCondition by summarizing the state of other conditions.
	// A step counter is added to represent progress during the provisioning process (instead we are hiding during the deletion process).
	applicableConditions := []clusterv1.ConditionType{
		infrav1.VpcReadyCondition,
		infrav1.SubnetsReadyCondition,
		infrav1.ClusterSecurityGroupsReadyCondition,
		infrav1.LoadBalancerReadyCondition,
	}

	if s.VPC().IsManaged(s.Name()) {
		applicableConditions = append(applicableConditions,
			infrav1.InternetGatewayReadyCondition,
			infrav1.NatGatewaysReadyCondition,
			infrav1.RouteTablesReadyCondition)

		if s.AWSCluster.Spec.Bastion.Enabled {
			applicableConditions = append(applicableConditions, infrav1.BastionHostReadyCondition)
		}
	}

	conditions.SetSummary(s.AWSCluster,
		conditions.WithConditions(applicableConditions...),
		conditions.WithStepCounterIf(s.AWSCluster.ObjectMeta.DeletionTimestamp.IsZero()),
		conditions.WithStepCounter(),
	)

	return s.patchHelper.Patch(
		context.TODO(),
		s.AWSCluster,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			infrav1.VpcReadyCondition,
			infrav1.SubnetsReadyCondition,
			infrav1.InternetGatewayReadyCondition,
			infrav1.NatGatewaysReadyCondition,
			infrav1.RouteTablesReadyCondition,
			infrav1.ClusterSecurityGroupsReadyCondition,
			infrav1.BastionHostReadyCondition,
			infrav1.LoadBalancerReadyCondition,
			infrav1.PrincipalUsageAllowedCondition,
		}})
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() error {
	return s.PatchObject()
}

// AdditionalTags returns AdditionalTags from the scope's AWSCluster. The returned value will never be nil.
func (s *ClusterScope) AdditionalTags() infrav1.Tags {
	if s.AWSCluster.Spec.AdditionalTags == nil {
		s.AWSCluster.Spec.AdditionalTags = infrav1.Tags{}
	}

	return s.AWSCluster.Spec.AdditionalTags.DeepCopy()
}

// APIServerPort returns the APIServerPort to use when creating the load balancer.
func (s *ClusterScope) APIServerPort() int32 {
	if s.Cluster.Spec.ClusterNetwork != nil && s.Cluster.Spec.ClusterNetwork.APIServerPort != nil {
		return *s.Cluster.Spec.ClusterNetwork.APIServerPort
	}
	return 6443
}

// SetFailureDomain sets the infrastructure provider failure domain key to the spec given as input.
func (s *ClusterScope) SetFailureDomain(id string, spec clusterv1.FailureDomainSpec) {
	if s.AWSCluster.Status.FailureDomains == nil {
		s.AWSCluster.Status.FailureDomains = make(clusterv1.FailureDomains)
	}
	s.AWSCluster.Status.FailureDomains[id] = spec
}

// InfraCluster returns the AWS infrastructure cluster or control plane object.
func (s *ClusterScope) InfraCluster() cloud.ClusterObject {
	return s.AWSCluster
}

// ClusterObj returns the cluster object.
func (s *ClusterScope) ClusterObj() cloud.ClusterObject {
	return s.Cluster
}

// Session returns the AWS SDK session. Used for creating clients.
func (s *ClusterScope) Session() awsclient.ConfigProvider {
	return s.session
}

// ServiceLimiter returns the AWS SDK session. Used for creating clients.
func (s *ClusterScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// Bastion returns the bastion details.
func (s *ClusterScope) Bastion() *infrav1.Bastion {
	return &s.AWSCluster.Spec.Bastion
}

// SetBastionInstance sets the bastion instance in the status of the cluster.
func (s *ClusterScope) SetBastionInstance(instance *infrav1.Instance) {
	s.AWSCluster.Status.Bastion = instance
}

// SSHKeyName returns the SSH key name to use for instances.
func (s *ClusterScope) SSHKeyName() *string {
	return s.AWSCluster.Spec.SSHKeyName
}

// ControllerName returns the name of the controller that
// created the ClusterScope.
func (s *ClusterScope) ControllerName() string {
	return s.controllerName
}

// ImageLookupFormat returns the format string to use when looking up AMIs.
func (s *ClusterScope) ImageLookupFormat() string {
	return s.AWSCluster.Spec.ImageLookupFormat
}

// ImageLookupOrg returns the organization name to use when looking up AMIs.
func (s *ClusterScope) ImageLookupOrg() string {
	return s.AWSCluster.Spec.ImageLookupOrg
}

// ImageLookupBaseOS returns the base operating system name to use when looking up AMIs.
func (s *ClusterScope) ImageLookupBaseOS() string {
	return s.AWSCluster.Spec.ImageLookupBaseOS
}
