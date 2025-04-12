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
	"fmt"

	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

// ManagedMachinePoolScopeParams defines the input parameters used to create a new Scope.
type ManagedMachinePoolScopeParams struct {
	Client             client.Client
	Logger             *logger.Logger
	Cluster            *clusterv1.Cluster
	ControlPlane       *ekscontrolplanev1.AWSManagedControlPlane
	ManagedMachinePool *expinfrav1.AWSManagedMachinePool
	MachinePool        *expclusterv1.MachinePool
	ControllerName     string
	Endpoints          []ServiceEndpoint
	Session            awsclient.ConfigProvider

	EnableIAM            bool
	AllowAdditionalRoles bool

	InfraCluster EC2Scope
}

// NewManagedMachinePoolScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewManagedMachinePoolScope(params ManagedMachinePoolScopeParams) (*ManagedMachinePoolScope, error) {
	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil AWSManagedMachinePool")
	}
	if params.MachinePool == nil {
		return nil, errors.New("failed to generate new scope from nil MachinePool")
	}
	if params.ManagedMachinePool == nil {
		return nil, errors.New("failed to generate new scope from nil ManagedMachinePool")
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

	ammpHelper, err := patch.NewHelper(params.ManagedMachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init AWSManagedMachinePool patch helper")
	}
	mpHelper, err := patch.NewHelper(params.MachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init MachinePool patch helper")
	}

	return &ManagedMachinePoolScope{
		Logger:                     *params.Logger,
		Client:                     params.Client,
		patchHelper:                ammpHelper,
		capiMachinePoolPatchHelper: mpHelper,

		Cluster:              params.Cluster,
		ControlPlane:         params.ControlPlane,
		ManagedMachinePool:   params.ManagedMachinePool,
		MachinePool:          params.MachinePool,
		EC2Scope:             params.InfraCluster,
		session:              session,
		serviceLimiters:      serviceLimiters,
		controllerName:       params.ControllerName,
		enableIAM:            params.EnableIAM,
		allowAdditionalRoles: params.AllowAdditionalRoles,
	}, nil
}

// ManagedMachinePoolScope defines the basic context for an actuator to operate upon.
type ManagedMachinePoolScope struct {
	logger.Logger
	client.Client
	patchHelper                *patch.Helper
	capiMachinePoolPatchHelper *patch.Helper

	Cluster            *clusterv1.Cluster
	ControlPlane       *ekscontrolplanev1.AWSManagedControlPlane
	ManagedMachinePool *expinfrav1.AWSManagedMachinePool
	MachinePool        *expclusterv1.MachinePool
	EC2Scope           EC2Scope

	session         awsclient.ConfigProvider
	serviceLimiters throttle.ServiceLimiters
	controllerName  string

	enableIAM            bool
	allowAdditionalRoles bool
}

// ManagedPoolName returns the managed machine pool name.
func (s *ManagedMachinePoolScope) ManagedPoolName() string {
	return s.ManagedMachinePool.Name
}

// ServiceLimiter returns the AWS SDK session. Used for creating clients.
func (s *ManagedMachinePoolScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ClusterName returns the cluster name.
func (s *ManagedMachinePoolScope) ClusterName() string {
	return s.ControlPlane.Spec.EKSClusterName
}

// EnableIAM indicates that reconciliation should create IAM roles.
func (s *ManagedMachinePoolScope) EnableIAM() bool {
	return s.enableIAM
}

// AllowAdditionalRoles indicates if additional roles can be added to the created IAM roles.
func (s *ManagedMachinePoolScope) AllowAdditionalRoles() bool {
	return s.allowAdditionalRoles
}

// Partition returns the machine pool subnet IDs.
func (s *ManagedMachinePoolScope) Partition() string {
	return system.GetPartitionFromRegion(s.ControlPlane.Spec.Region)
}

// IdentityRef returns the cluster identityRef.
func (s *ManagedMachinePoolScope) IdentityRef() *infrav1.AWSIdentityReference {
	return s.ControlPlane.Spec.IdentityRef
}

// AdditionalTags returns AdditionalTags from the scope's ManagedMachinePool
// The returned value will never be nil.
func (s *ManagedMachinePoolScope) AdditionalTags() infrav1.Tags {
	tags := make(infrav1.Tags)

	// Start with the cluster-wide tags...
	tags.Merge(s.EC2Scope.AdditionalTags())
	// ... and merge in the Machine's
	tags.Merge(s.ManagedMachinePool.Spec.AdditionalTags)

	return tags
}

// RoleName returns the node group role name.
func (s *ManagedMachinePoolScope) RoleName() string {
	return s.ManagedMachinePool.Spec.RoleName
}

// Version returns the nodegroup Kubernetes version.
func (s *ManagedMachinePoolScope) Version() *string {
	return s.MachinePool.Spec.Template.Spec.Version
}

// ControlPlaneSubnets returns the control plane subnets.
func (s *ManagedMachinePoolScope) ControlPlaneSubnets() infrav1.Subnets {
	return s.ControlPlane.Spec.NetworkSpec.Subnets
}

// SubnetIDs returns the machine pool subnet IDs.
func (s *ManagedMachinePoolScope) SubnetIDs() ([]string, error) {
	strategy, err := newDefaultSubnetPlacementStrategy(&s.Logger)
	if err != nil {
		return []string{}, fmt.Errorf("getting subnet placement strategy: %w", err)
	}

	return strategy.Place(&placementInput{
		SpecSubnetIDs:           s.ManagedMachinePool.Spec.SubnetIDs,
		SpecAvailabilityZones:   s.ManagedMachinePool.Spec.AvailabilityZones,
		ParentAvailabilityZones: s.MachinePool.Spec.FailureDomains,
		ControlplaneSubnets:     s.ControlPlaneSubnets(),
		SubnetPlacementType:     s.ManagedMachinePool.Spec.AvailabilityZoneSubnetType,
	})
}

// NodegroupReadyFalse marks the ready condition false using warning if error isn't
// empty.
func (s *ManagedMachinePoolScope) NodegroupReadyFalse(reason string, err string) error {
	severity := clusterv1.ConditionSeverityWarning
	if err == "" {
		severity = clusterv1.ConditionSeverityInfo
	}
	conditions.MarkFalse(
		s.ManagedMachinePool,
		expinfrav1.EKSNodegroupReadyCondition,
		reason,
		severity,
		"%s",
		err,
	)
	if err := s.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to mark nodegroup not ready")
	}
	return nil
}

// IAMReadyFalse marks the ready condition false using warning if error isn't
// empty.
func (s *ManagedMachinePoolScope) IAMReadyFalse(reason string, err string) error {
	severity := clusterv1.ConditionSeverityWarning
	if err == "" {
		severity = clusterv1.ConditionSeverityInfo
	}
	conditions.MarkFalse(
		s.ManagedMachinePool,
		expinfrav1.IAMNodegroupRolesReadyCondition,
		reason,
		severity,
		"%s",
		err,
	)
	if err := s.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to mark nodegroup role not ready")
	}
	return nil
}

// PatchObject persists the control plane configuration and status.
func (s *ManagedMachinePoolScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.ManagedMachinePool,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			expinfrav1.EKSNodegroupReadyCondition,
			expinfrav1.IAMNodegroupRolesReadyCondition,
		}})
}

// PatchCAPIMachinePoolObject persists the capi machinepool configuration and status.
func (s *ManagedMachinePoolScope) PatchCAPIMachinePoolObject(ctx context.Context) error {
	return s.capiMachinePoolPatchHelper.Patch(
		ctx,
		s.MachinePool,
	)
}

// Close closes the current scope persisting the control plane configuration and status.
func (s *ManagedMachinePoolScope) Close() error {
	return s.PatchObject()
}

// InfraCluster returns the AWS infrastructure cluster or control plane object.
func (s *ManagedMachinePoolScope) InfraCluster() cloud.ClusterObject {
	return s.ControlPlane
}

// ClusterObj returns the cluster object.
func (s *ManagedMachinePoolScope) ClusterObj() cloud.ClusterObject {
	return s.Cluster
}

// Session returns the AWS SDK session. Used for creating clients.
func (s *ManagedMachinePoolScope) Session() awsclient.ConfigProvider {
	return s.session
}

// ControllerName returns the name of the controller that
// created the ManagedMachinePool.
func (s *ManagedMachinePoolScope) ControllerName() string {
	return s.controllerName
}

// KubernetesClusterName is the name of the EKS cluster name.
func (s *ManagedMachinePoolScope) KubernetesClusterName() string {
	return s.ControlPlane.Spec.EKSClusterName
}

// NodegroupName is the name of the EKS nodegroup.
func (s *ManagedMachinePoolScope) NodegroupName() string {
	return s.ManagedMachinePool.Spec.EKSNodegroupName
}

// Name returns the name of the AWSManagedMachinePool.
func (s *ManagedMachinePoolScope) Name() string {
	return s.ManagedMachinePool.Name
}

// Namespace returns the namespace of the AWSManagedMachinePool.
func (s *ManagedMachinePoolScope) Namespace() string {
	return s.ManagedMachinePool.Namespace
}

// GetRawBootstrapData returns the raw bootstrap data from the linked Machine's bootstrap.dataSecretName.
func (s *ManagedMachinePoolScope) GetRawBootstrapData() ([]byte, *types.NamespacedName, error) {
	if s.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName == nil {
		return nil, nil, errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: s.Namespace(), Name: *s.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName}

	if err := s.Client.Get(context.TODO(), key, secret); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to retrieve bootstrap data secret for AWSManagedMachinePool %s/%s", s.Namespace(), s.Name())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, nil, errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, &key, nil
}

// GetObjectMeta returns the ObjectMeta for the AWSManagedMachinePool.
func (s *ManagedMachinePoolScope) GetObjectMeta() *metav1.ObjectMeta {
	return &s.ManagedMachinePool.ObjectMeta
}

// GetSetter returns the condition setter.
func (s *ManagedMachinePoolScope) GetSetter() conditions.Setter {
	return s.ManagedMachinePool
}

// GetEC2Scope returns the EC2Scope.
func (s *ManagedMachinePoolScope) GetEC2Scope() EC2Scope {
	return s.EC2Scope
}

// IsEKSManaged returns true if the control plane is managed by EKS.
func (s *ManagedMachinePoolScope) IsEKSManaged() bool {
	return true
}

// GetLaunchTemplateIDStatus returns the launch template ID status.
func (s *ManagedMachinePoolScope) GetLaunchTemplateIDStatus() string {
	if s.ManagedMachinePool.Status.LaunchTemplateID != nil {
		return *s.ManagedMachinePool.Status.LaunchTemplateID
	}
	return ""
}

// SetLaunchTemplateIDStatus sets the launch template ID status.
func (s *ManagedMachinePoolScope) SetLaunchTemplateIDStatus(id string) {
	s.ManagedMachinePool.Status.LaunchTemplateID = &id
}

// GetLaunchTemplateLatestVersionStatus returns the launch template latest version status.
func (s *ManagedMachinePoolScope) GetLaunchTemplateLatestVersionStatus() string {
	if s.ManagedMachinePool.Status.LaunchTemplateVersion != nil {
		return *s.ManagedMachinePool.Status.LaunchTemplateVersion
	}
	return ""
}

// SetLaunchTemplateLatestVersionStatus sets the launch template latest version status.
func (s *ManagedMachinePoolScope) SetLaunchTemplateLatestVersionStatus(version string) {
	s.ManagedMachinePool.Status.LaunchTemplateVersion = &version
}

// GetLaunchTemplate returns the launch template.
func (s *ManagedMachinePoolScope) GetLaunchTemplate() *expinfrav1.AWSLaunchTemplate {
	return s.ManagedMachinePool.Spec.AWSLaunchTemplate
}

// GetMachinePool returns the machine pool.
func (s *ManagedMachinePoolScope) GetMachinePool() *expclusterv1.MachinePool {
	return s.MachinePool
}

// LaunchTemplateName returns the launch template name.
func (s *ManagedMachinePoolScope) LaunchTemplateName() string {
	return fmt.Sprintf("%s-%s", s.ControlPlane.Name, s.ManagedMachinePool.Name)
}

// GetRuntimeObject returns the AWSManagedMachinePool, in runtime.Object form.
func (s *ManagedMachinePoolScope) GetRuntimeObject() runtime.Object {
	return s.ManagedMachinePool
}
